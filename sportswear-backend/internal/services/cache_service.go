package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"sportswear-backend/internal/utils"
)

// CacheService Redis 缓存服务
// 高性能查询保障：热点内容缓存 + 发布时主动失效
type CacheService struct {
	client  *redis.Client
	enabled bool
}

// NewCacheService 创建缓存服务（Redis 不可用时自动降级为直查 DB）
func NewCacheService(client *redis.Client) *CacheService {
	enabled := false
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err == nil {
			enabled = true
			utils.Logger.Infow("Redis 缓存已启用")
		} else {
			utils.Logger.Warnw("Redis 连接失败，缓存降级为直查 DB", "error", err.Error())
		}
	}
	return &CacheService{client: client, enabled: enabled}
}

// IsEnabled 缓存是否可用
func (s *CacheService) IsEnabled() bool {
	return s.enabled
}

// Get 从缓存获取数据（JSON 反序列化）
func (s *CacheService) Get(key string, dest interface{}) bool {
	if !s.enabled {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		// 缓存数据损坏，删除该键
		s.Delete(key)
		return false
	}
	return true
}

// Set 写入缓存
func (s *CacheService) Set(key string, value interface{}, ttl time.Duration) {
	if !s.enabled {
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
	if !s.enabled || len(keys) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	s.client.Del(ctx, keys...)
}

// InvalidatePattern 按模式批量失效缓存（如 cache:product:*）
func (s *CacheService) InvalidatePattern(pattern string) {
	if !s.enabled {
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

// ========== 缓存键规范 ==========

const (
	cacheTTLShort  = 5 * time.Minute  // 内容详情
	cacheTTLMedium = 10 * time.Minute // 列表数据
	cacheTTLLong   = 30 * time.Minute // 导航等低频变更
)

// CacheKeys 缓存键生成器
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

// InvalidateContentCache 内容发布/更新/下线时主动失效相关缓存
// 确保数据一致性：后台更新 → 缓存立即失效 → 下次请求回源 DB → 写入新缓存
// slug 为空时按实体整体失效（例如后台只持有 id，无法拿到 slug 时）
func (s *CacheService) InvalidateContentCache(entityType, slug string) {
	// 首页聚合了产品/博客/案例/认证/生产流程/工厂/分类/导航等多类内容，
	// 任一内容变更都可能影响首页，统一一并失效最稳妥
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
	case "production-process":
		s.InvalidatePattern("cache:production-processes:*")
		invalidateHome()
	case "navigation":
		s.InvalidatePattern("cache:navigations:*")
	default:
		if slug != "" {
			s.InvalidatePattern("cache:" + entityType + ":" + slug + ":*")
		}
		s.InvalidatePattern("cache:" + entityType + ":*")
	}
}

// 导出的缓存 TTL 常量（供 handlers 等外部包在写缓存时使用）
const (
	CacheTTLShort  = cacheTTLShort  // 5 分钟：内容详情
	CacheTTLMedium = cacheTTLMedium // 10 分钟：列表数据
	CacheTTLLong   = cacheTTLLong   // 30 分钟：导航等低频变更
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
