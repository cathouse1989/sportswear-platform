package services

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
	"sportswear-backend/internal/utils"
)

// 门户缓存业务开关持久化 key（theme_configs）
const PortalCacheEnabledKey = "portal_cache_enabled"

// CacheService Redis 缓存服务
// 业务开关 portal_cache_enabled：
//   - 开启：公开接口走缓存；内容变更自动失效
//   - 关闭：公开接口直查 DB；不自动更新缓存；可手动刷新
type CacheService struct {
	client        *redis.Client
	redisOK       bool
	policyEnabled bool
	mu            sync.RWMutex
	db            *gorm.DB
	lastRefreshAt *time.Time
}

// NewCacheService 创建缓存服务（Redis 不可用时自动降级为直查 DB）
func NewCacheService(client *redis.Client) *CacheService {
	redisOK := false
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err == nil {
			redisOK = true
			utils.Logger.Infow("Redis 缓存已启用")
		} else {
			utils.Logger.Warnw("Redis 连接失败，缓存降级为直查 DB", "error", err.Error())
		}
	}
	return &CacheService{
		client:        client,
		redisOK:       redisOK,
		policyEnabled: true,
	}
}

// BindDB 绑定数据库并加载门户缓存业务开关（启动时调用）
func (s *CacheService) BindDB(db *gorm.DB) {
	s.mu.Lock()
	s.db = db
	s.mu.Unlock()
	s.loadPolicyFromDB()
}

func (s *CacheService) loadPolicyFromDB() {
	s.mu.RLock()
	db := s.db
	s.mu.RUnlock()
	if db == nil {
		return
	}
	var cfg models.ThemeConfig
	err := db.Where(`"key" = ?`, PortalCacheEnabledKey).First(&cfg).Error
	if err != nil {
		_ = db.Create(&models.ThemeConfig{
			Key:      PortalCacheEnabledKey,
			Value:    "true",
			Group:    "system",
			Name:     "门户内容缓存开关",
			IsActive: true,
		}).Error
		s.mu.Lock()
		s.policyEnabled = true
		s.mu.Unlock()
		return
	}
	enabled := parseBoolConfig(cfg.Value, true)
	s.mu.Lock()
	s.policyEnabled = enabled
	s.mu.Unlock()
}

func parseBoolConfig(raw string, def bool) bool {
	v := raw
	if len(v) >= 2 && v[0] == '"' {
		var s string
		if json.Unmarshal([]byte(v), &s) == nil {
			v = s
		}
	} else {
		var b bool
		if json.Unmarshal([]byte(v), &b) == nil {
			return b
		}
	}
	switch v {
	case "1", "true", "TRUE", "True", "yes", "on":
		return true
	case "0", "false", "FALSE", "False", "no", "off":
		return false
	default:
		return def
	}
}

// IsRedisOK Redis 是否可用
func (s *CacheService) IsRedisOK() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.redisOK
}

// IsPolicyEnabled 业务开关是否开启
func (s *CacheService) IsPolicyEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.policyEnabled
}

// IsEnabled 公开读路径是否使用缓存（Redis 可用且业务开关开启）
func (s *CacheService) IsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.redisOK && s.policyEnabled
}

// Status 后台展示用状态
func (s *CacheService) Status() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := map[string]interface{}{
		"redis_ok":        s.redisOK,
		"enabled":         s.policyEnabled,
		"using_cache":     s.redisOK && s.policyEnabled,
		"last_refresh_at": nil,
	}
	if s.lastRefreshAt != nil {
		out["last_refresh_at"] = s.lastRefreshAt.UTC().Format(time.RFC3339)
	}
	return out
}

// SetPolicyEnabled 设置业务开关并持久化；开启时自动预热
func (s *CacheService) SetPolicyEnabled(enabled bool, warmUp bool) error {
	s.mu.Lock()
	s.policyEnabled = enabled
	db := s.db
	s.mu.Unlock()

	if db != nil {
		val := "false"
		if enabled {
			val = "true"
		}
		res := db.Model(&models.ThemeConfig{}).Where(`"key" = ?`, PortalCacheEnabledKey).Update("value", val)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			if err := db.Create(&models.ThemeConfig{
				Key:      PortalCacheEnabledKey,
				Value:    val,
				Group:    "system",
				Name:     "门户内容缓存开关",
				IsActive: true,
			}).Error; err != nil {
				return err
			}
		}
	}

	if enabled && warmUp && s.IsRedisOK() && db != nil {
		WarmUpPublicCache(db, s)
		s.MarkRefreshed()
	}
	return nil
}

// MarkRefreshed 记录最近一次全量刷新时间
func (s *CacheService) MarkRefreshed() {
	now := time.Now().UTC()
	s.mu.Lock()
	s.lastRefreshAt = &now
	s.mu.Unlock()
}

// Get 从缓存获取数据（仅业务开关开启时命中）
func (s *CacheService) Get(key string, dest interface{}) bool {
	if !s.IsEnabled() {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		s.deleteKeys(key)
		return false
	}
	return true
}

// Set 写入缓存（Redis 可用即可写，供手动刷新在开关关闭时预热）
func (s *CacheService) Set(key string, value interface{}, ttl time.Duration) {
	if !s.IsRedisOK() {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		return
	}
	s.client.Set(ctx, key, data, ttl)
}

// Delete 删除缓存键
func (s *CacheService) Delete(keys ...string) {
	if !s.IsRedisOK() || len(keys) == 0 {
		return
	}
	s.deleteKeys(keys...)
}

func (s *CacheService) deleteKeys(keys ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	s.client.Del(ctx, keys...)
}

// InvalidatePattern 按模式批量失效缓存
func (s *CacheService) InvalidatePattern(pattern string) {
	if !s.IsRedisOK() {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	iter := s.client.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		s.client.Del(ctx, keys...)
	}
}

const (
	cacheTTLShort  = 5 * time.Minute
	cacheTTLMedium = 10 * time.Minute
	cacheTTLLong   = 30 * time.Minute
)

var CacheKeys = struct {
	Page     func(slug, lang string) string
	Product  func(slug, lang string) string
	Products func(lang string, page int) string
	Nav      func(navType string) string
	Home     func(lang string) string
}{
	Page: func(slug, lang string) string {
		return "cache:page:" + slug + ":" + lang
	},
	Product: func(slug, lang string) string {
		return "cache:product:" + slug + ":" + lang
	},
	Products: func(lang string, page int) string {
		return "cache:products:" + lang + ":p" + itoa(page)
	},
	Nav: func(navType string) string {
		return "cache:navigations:" + navType
	},
	Home: func(lang string) string {
		return "cache:home:" + lang
	},
}

// InvalidateContentCache 内容变更时失效缓存；仅业务开关开启时执行
func (s *CacheService) InvalidateContentCache(entityType, slug string) {
	if !s.IsEnabled() {
		return
	}
	invalidateHome := func() { s.InvalidatePattern("cache:home:*") }

	switch entityType {
	case "page":
		if slug != "" {
			s.InvalidatePattern("cache:page:" + slug + ":*")
		} else {
			s.InvalidatePattern("cache:page:*")
		}
		invalidateHome()
	case "product":
		if slug != "" {
			s.InvalidatePattern("cache:product:" + slug + ":*")
		}
		s.InvalidatePattern("cache:products:*")
		invalidateHome()
	case "category":
		s.InvalidatePattern("cache:categories:*")
		invalidateHome()
	case "series":
		s.InvalidatePattern("cache:series:*")
		invalidateHome()
	case "fabric":
		s.InvalidatePattern("cache:fabrics:*")
		invalidateHome()
	case "blog":
		if slug != "" {
			s.InvalidatePattern("cache:blog:" + slug + ":*")
		}
		s.InvalidatePattern("cache:blogs:*")
		invalidateHome()
	case "case":
		if slug != "" {
			s.InvalidatePattern("cache:case:" + slug + ":*")
		}
		s.InvalidatePattern("cache:cases:*")
		invalidateHome()
	case "faq":
		s.InvalidatePattern("cache:faqs:*")
		invalidateHome()
	case "factory":
		s.InvalidatePattern("cache:factories:*")
		invalidateHome()
	case "certification":
		s.InvalidatePattern("cache:certifications:*")
		invalidateHome()
	case "self_media":
		s.InvalidatePattern("cache:self-medias:*")
		invalidateHome()
	case "navigation":
		s.InvalidatePattern("cache:navigations:*")
	case "theme":
		s.Delete("cache:theme")
	default:
		if slug != "" {
			s.InvalidatePattern("cache:" + entityType + ":" + slug + ":*")
		}
		s.InvalidatePattern("cache:" + entityType + ":*")
	}
}

// RefreshPublicCache 手动刷新/发布门户缓存
func (s *CacheService) RefreshPublicCache(db *gorm.DB) bool {
	if s == nil || !s.IsRedisOK() || db == nil {
		return false
	}
	for _, pattern := range []string{
		"cache:home:*", "cache:page:*", "cache:product:*", "cache:products:*",
		"cache:categories:*", "cache:series:*", "cache:fabrics:*",
		"cache:blogs:*", "cache:blog:*", "cache:cases:*", "cache:case:*",
		"cache:faqs:*", "cache:factories:*", "cache:certifications:*",
		"cache:self-medias:*",
		"cache:navigations:*",
	} {
		s.InvalidatePattern(pattern)
	}
	s.Delete("cache:theme")
	WarmUpPublicCache(db, s)
	s.MarkRefreshed()
	return true
}

const (
	CacheTTLShort  = cacheTTLShort
	CacheTTLMedium = cacheTTLMedium
	CacheTTLLong   = cacheTTLLong
)

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
