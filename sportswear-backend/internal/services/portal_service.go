package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// PortalService 门户展示配置服务
// 体系化方案：页面模块化 + 版本化发布 + 主题配置
type PortalService struct {
	db    *gorm.DB
	cache *CacheService
}

// NewPortalService 创建门户配置服务
func NewPortalService(db *gorm.DB, cache *CacheService) *PortalService {
	return &PortalService{db: db, cache: cache}
}

// ==================== 页面版本化发布 ====================

// SaveDraft 保存草稿版本（不影响线上）
// 工作流：后台编辑 → 保存草稿 → 预览 → 发布
func (s *PortalService) SaveDraft(pageID string, snapshot interface{}, userID string, note string) (*models.PageVersion, error) {
	pid, err := uuid.Parse(pageID)
	if err != nil {
		return nil, errors.New("页面 ID 无效")
	}

	// 获取当前最大版本号
	var maxVersion int
	s.db.Model(&models.PageVersion{}).Where("page_id = ?", pid).
		Select("COALESCE(MAX(version), 0)").Scan(&maxVersion)

	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, errors.New("快照序列化失败")
	}

	version := models.PageVersion{
		PageID:   pid,
		Version:  maxVersion + 1,
		Status:   "draft",
		Snapshot: string(snapshotJSON),
		Note:     note,
	}
	if uid, err := uuid.Parse(userID); err == nil {
		version.CreatedBy = &uid
	}

	if err := s.db.Create(&version).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// PublishVersion 发布指定版本（前台立即生效，旧版本归档）
func (s *PortalService) PublishVersion(versionID string) error {
	vid, err := uuid.Parse(versionID)
	if err != nil {
		return errors.New("版本 ID 无效")
	}

	var version models.PageVersion
	if err := s.db.First(&version, "id = ?", vid).Error; err != nil {
		return errors.New("版本不存在")
	}
	if version.Status == "published" {
		return errors.New("该版本已发布")
	}

	now := time.Now()
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		// 归档当前已发布版本
		if err := tx.Model(&models.PageVersion{}).
			Where("page_id = ? AND status = ?", version.PageID, "published").
			Update("status", "archived").Error; err != nil {
			return err
		}
		// 发布新版本
		if err := tx.Model(&version).Updates(map[string]interface{}{
			"status":       "published",
			"published_at": now,
		}).Error; err != nil {
			return err
		}
		// 同步更新页面主表状态
		if err := tx.Model(&models.Page{}).Where("id = ?", version.PageID).
			Updates(map[string]interface{}{
				"status":       models.ContentStatusPublished,
				"published_at": now,
			}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	// 页面发布后主动失效 page 与首页缓存，保证门户立即生效
	if s.cache != nil {
		s.cache.InvalidateContentCache("page", "")
	}
	return nil
}

// ListVersions 页面版本历史
func (s *PortalService) ListVersions(pageID string) ([]models.PageVersion, error) {
	var versions []models.PageVersion
	err := s.db.Where("page_id = ?", pageID).Order("version DESC").Find(&versions).Error
	return versions, err
}

// GetPublishedVersion 获取已发布版本（公开接口使用）
func (s *PortalService) GetPublishedVersion(pageID string) (*models.PageVersion, error) {
	var version models.PageVersion
	err := s.db.Where("page_id = ? AND status = ?", pageID, "published").
		Order("version DESC").First(&version).Error
	if err != nil {
		return nil, errors.New("无已发布版本")
	}
	return &version, nil
}

// RollbackVersion 回滚到历史版本（复制为新的草稿并发布）
func (s *PortalService) RollbackVersion(versionID string) error {
	var source models.PageVersion
	if err := s.db.First(&source, "id = ?", versionID).Error; err != nil {
		return errors.New("版本不存在")
	}

	// 创建新草稿版本（内容来自历史版本）
	draft, err := s.SaveDraft(source.PageID.String(), json.RawMessage(source.Snapshot), "", "回滚自 v"+itoaInt(source.Version))
	if err != nil {
		return err
	}
	return s.PublishVersion(draft.ID.String())
}

// ==================== 主题配置 ====================

// ListThemeConfigs 获取全部主题配置列表（后台管理用，返回数组含 key/value/name/group）
func (s *PortalService) ListThemeConfigs() ([]models.ThemeConfig, error) {
	var configs []models.ThemeConfig
	if err := s.db.Where("is_active = ?", true).
		// 注意：group/key 均为 PostgreSQL 保留字，必须显式加引号
		Order(`"group" ASC, "key" ASC`).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetThemeConfig 获取全部主题配置（公开接口：前端渲染视觉风格）
func (s *PortalService) GetThemeConfig() (map[string]interface{}, error) {
	// 尝试缓存
	if s.cache != nil && s.cache.IsEnabled() {
		var cached map[string]interface{}
		if s.cache.Get("cache:theme", &cached) {
			return cached, nil
		}
	}

	var configs []models.ThemeConfig
	if err := s.db.Where("is_active = ?", true).Find(&configs).Error; err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, c := range configs {
		var val interface{}
		if err := json.Unmarshal([]byte(c.Value), &val); err == nil {
			result[c.Key] = val
		} else {
			result[c.Key] = c.Value
		}
	}

	// 写入缓存
	if s.cache != nil && s.cache.IsEnabled() {
		s.cache.Set("cache:theme", result, cacheTTLLong)
	}
	return result, nil
}

// UpdateThemeConfig 更新主题配置项
func (s *PortalService) UpdateThemeConfig(key string, value string) error {
	result := s.db.Model(&models.ThemeConfig{}).Where("key = ?", key).Update("value", value)
	if result.RowsAffected == 0 {
		return errors.New("配置项不存在")
	}
	// 失效主题缓存
	if s.cache != nil {
		s.cache.Delete("cache:theme")
	}
	return nil
}

// InitDefaultTheme 初始化默认主题配置（幂等：仅补充缺失的配置项，不覆盖已有值）。
// 逐键检查而非整表判断，保证后续版本新增配置项时，存量数据库也能在启动时自动补齐。
func (s *PortalService) InitDefaultTheme() error {
	for i := range models.DefaultThemeConfigs {
		item := models.DefaultThemeConfigs[i]
		var count int64
		if err := s.db.Model(&models.ThemeConfig{}).Where("\"key\" = ?", item.Key).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := s.db.Create(&item).Error; err != nil {
				return err
			}
		}
	}

	// 社交/自媒体链接已迁移至「自媒体管理」：停用 theme_configs 中历史遗留的相关键，
	// 避免「其他主题配置项」与「自媒体管理」两处重复出现配置（只保留自媒体管理一处）。
	deprecatedKeys := []string{
		"social_youtube", "social_instagram", "social_xiaohongshu",
		"social_facebook", "social_twitter", "social_linkedin",
		"self_media_max_display",
	}
	if err := s.db.Model(&models.ThemeConfig{}).Where("\"key\" IN ?", deprecatedKeys).Update("is_active", false).Error; err != nil {
		return err
	}
	return nil
}

func itoaInt(n int) string {
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
