package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
	"sportswear-backend/internal/utils"
)

// CMSService CMS 服务
type CMSService struct {
	db *gorm.DB
}

// NewCMSService 创建 CMS 服务
func NewCMSService(db *gorm.DB) *CMSService {
	return &CMSService{db: db}
}

// ==================== 页面管理 ====================

// ListPages 页面列表
func (s *CMSService) ListPages(page, pageSize int, keyword, status string) ([]models.Page, int64, error) {
	var pages []models.Page
	var total int64

	query := s.db.Model(&models.Page{})
	if keyword != "" {
		query = query.Where("title LIKE ? OR slug LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Preload("Translations").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("sort_order ASC, created_at DESC").Find(&pages).Error

	return pages, total, err
}

// GetPage 获取页面
func (s *CMSService) GetPage(id string) (*models.Page, error) {
	var page models.Page
	err := s.db.Preload("Modules").Preload("Translations").
		First(&page, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("页面不存在")
	}
	return &page, nil
}

// GetPageBySlug 通过 Slug 获取页面
func (s *CMSService) GetPageBySlug(slug string) (*models.Page, error) {
	var page models.Page
	err := s.db.Preload("Modules").Preload("Translations").
		First(&page, "slug = ? AND status = ?", slug, models.ContentStatusPublished).Error
	if err != nil {
		return nil, errors.New("页面不存在")
	}
	return &page, nil
}

// ListPublishedLandingPages 列出所有已发布的「落地页」（走 /{slug} 兜底页的页面类型）。
// 固定列表页（home/product/blog/case/faq/contact 等）已由门户静态 sitemap 覆盖，此处排除。
func (s *CMSService) ListPublishedLandingPages() ([]models.Page, error) {
	fixedTypes := []models.PageType{
		models.PageTypeHome,
		models.PageTypeProduct,
		models.PageTypeProductCategory,
		models.PageTypeBlog,
		models.PageTypeCase,
		models.PageTypeFAQ,
		models.PageTypeContact,
	}
	var pages []models.Page
	err := s.db.Where("status = ? AND type NOT IN ?", models.ContentStatusPublished, fixedTypes).
		Order("sort_order ASC").Find(&pages).Error
	return pages, err
}

// CreatePage 创建页面
func (s *CMSService) CreatePage(req *PageRequest) (*models.Page, error) {
	page := models.Page{
		Title:     req.Title,
		Slug:      req.Slug,
		Type:      req.Type,
		Status:    req.Status,
		Template:  req.Template,
		SortOrder: req.SortOrder,
	}

	if err := s.db.Create(&page).Error; err != nil {
		return nil, err
	}

	// 创建模块
	if len(req.Modules) > 0 {
		for _, m := range req.Modules {
			module := models.PageModule{
				PageID:    page.ID,
				Type:      m.Type,
				Title:     m.Title,
				SortOrder: m.SortOrder,
				IsVisible: m.IsVisible,
				Config:    m.Config,
			}
			s.db.Create(&module)
		}
	}

	// 创建翻译
	if len(req.Translations) > 0 {
		for _, t := range req.Translations {
			s.db.Create(&models.PageTranslation{
				PageID:   page.ID,
				Language: t.Language,
				Title:    t.Title,
				Content:  t.Content,
				Status:   models.TranslationStatusPublished,
			})
		}
	}

	// 创建 SEO（未显式提供时，自动用页面标题生成默认值，幂等不覆盖）
	if req.SEO != nil {
		s.db.Create(&models.SEO{
			EntityType:    "page",
			EntityID:      page.ID,
			Language:      req.SEO.Language,
			Title:         req.SEO.Title,
			Description:   req.SEO.Description,
			Keywords:      req.SEO.Keywords,
			Canonical:     req.SEO.Canonical,
			Robots:        req.SEO.Robots,
			OGTitle:       req.SEO.OGTitle,
			OGDescription: req.SEO.OGDescription,
			OGImage:       req.SEO.OGImage,
			SchemaData:    req.SEO.SchemaData,
		})
	} else {
		if err := EnsureEntitySEO(s.db, "page", page.ID, req.Title, ""); err != nil {
			return nil, err
		}
	}

	return s.GetPage(page.ID.String())
}

// PageRequest 页面请求
type PageRequest struct {
	Title        string                   `json:"title" binding:"required"`
	Slug         string                   `json:"slug" binding:"required"`
	Type         models.PageType          `json:"type"`
	Status       models.ContentStatus     `json:"status"`
	Template     string                   `json:"template"`
	SortOrder    int                      `json:"sort_order"`
	Modules      []PageModuleRequest      `json:"modules"`
	Translations []PageTranslationRequest `json:"translations"`
	SEO          *SEORequest              `json:"seo"`
}

// PageModuleRequest 页面模块请求
type PageModuleRequest struct {
	Type      string `json:"type" binding:"required"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
	IsVisible bool   `json:"is_visible"`
	Config    string `json:"config"`
}

// PageTranslationRequest 页面翻译请求
type PageTranslationRequest struct {
	Language string `json:"language" binding:"required"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// UpdatePage 更新页面
func (s *CMSService) UpdatePage(id string, req *PageRequest) (*models.Page, error) {
	var page models.Page
	if err := s.db.First(&page, "id = ?", id).Error; err != nil {
		return nil, errors.New("页面不存在")
	}
	oldSlug := page.Slug
	oldType := page.Type

	// status/template 为空时不覆盖（编辑保存只传基础字段，避免清空发布状态）
	updates := map[string]interface{}{
		"title":      req.Title,
		"slug":       req.Slug,
		"type":       req.Type,
		"sort_order": req.SortOrder,
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Template != "" {
		updates["template"] = req.Template
	}
	if err := s.db.Model(&page).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 更新模块（先删除再重建）
	if req.Modules != nil {
		s.db.Where("page_id = ?", page.ID).Delete(&models.PageModule{})
		for _, m := range req.Modules {
			s.db.Create(&models.PageModule{
				PageID:    page.ID,
				Type:      m.Type,
				Title:     m.Title,
				SortOrder: m.SortOrder,
				IsVisible: m.IsVisible,
				Config:    m.Config,
			})
		}
	}

	// 更新翻译
	if len(req.Translations) > 0 {
		for _, t := range req.Translations {
			var translation models.PageTranslation
			err := s.db.Where("page_id = ? AND language = ?", page.ID, t.Language).First(&translation).Error
			if err != nil {
				s.db.Create(&models.PageTranslation{
					PageID:   page.ID,
					Language: t.Language,
					Title:    t.Title,
					Content:  t.Content,
					Status:   models.TranslationStatusPublished,
				})
			} else {
				s.db.Model(&translation).Updates(map[string]interface{}{
					"title":   t.Title,
					"content": t.Content,
					"status":  models.TranslationStatusPublished,
				})
			}
		}
	}

	// 更新 SEO
	if req.SEO != nil {
		var seo models.SEO
		err := s.db.Where("entity_type = ? AND entity_id = ?", "page", page.ID).First(&seo).Error
		if err != nil {
			s.db.Create(&models.SEO{
				EntityType:    "page",
				EntityID:      page.ID,
				Language:      req.SEO.Language,
				Title:         req.SEO.Title,
				Description:   req.SEO.Description,
				Keywords:      req.SEO.Keywords,
				Canonical:     req.SEO.Canonical,
				Robots:        req.SEO.Robots,
				OGTitle:       req.SEO.OGTitle,
				OGDescription: req.SEO.OGDescription,
				OGImage:       req.SEO.OGImage,
				SchemaData:    req.SEO.SchemaData,
			})
		} else {
			s.db.Model(&seo).Updates(map[string]interface{}{
				"language":       req.SEO.Language,
				"title":          req.SEO.Title,
				"description":    req.SEO.Description,
				"keywords":       req.SEO.Keywords,
				"canonical":      req.SEO.Canonical,
				"robots":         req.SEO.Robots,
				"og_title":       req.SEO.OGTitle,
				"og_description": req.SEO.OGDescription,
				"og_image":       req.SEO.OGImage,
				"schema_data":    req.SEO.SchemaData,
			})
		}
	}

	// slug/type 变更联动：级联更新关联导航 url 与 route SEO，避免死链与 SEO 失效
	if oldSlug != req.Slug || oldType != req.Type {
		s.SyncPagePathEffects(page.ID, oldSlug, req.Slug, oldType, req.Type)
	}

	return s.GetPage(id)
}

// HeroSlide 首页轮播图（hero banner）配置项
// 门户首页顶部的可滑动大图。后台在"页面管理 → 轮播图"中编辑，
// 通过 UpdateHeroSlides 以 JSONB 形式保存到首页 type=banner 模块的 config：
//
//	{"settings":{...},"slides":[{"image":"...","title":"...","subtitle":"...","button_text":"...","button_url":"..."}, ...]}
type HeroSlide struct {
	Image      string `json:"image"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	ButtonText string `json:"button_text"`
	ButtonURL  string `json:"button_url"`
}

// HeroSettings 首页轮播图展示参数（与 slides 一并写入 banner.config）
type HeroSettings struct {
	Autoplay     *bool  `json:"autoplay,omitempty"`
	IntervalMs   int    `json:"interval_ms,omitempty"`
	Transition   string `json:"transition,omitempty"` // fade | slide
	ShowDots     *bool  `json:"show_dots,omitempty"`
	ShowArrows   *bool  `json:"show_arrows,omitempty"`
	PauseOnHover *bool  `json:"pause_on_hover,omitempty"`
}

// DefaultHeroSettings 门户默认轮播行为（后台未配置时使用）
func DefaultHeroSettings() HeroSettings {
	autoplay := true
	showDots := true
	showArrows := true
	pauseOnHover := true
	return HeroSettings{
		Autoplay:     &autoplay,
		IntervalMs:   5000,
		Transition:   "fade",
		ShowDots:     &showDots,
		ShowArrows:   &showArrows,
		PauseOnHover: &pauseOnHover,
	}
}

// UnmarshalPageModuleConfig 解析页面模块 JSONB config。
// 兼容：对象 JSON、被二次编码成 JSON 字符串的对象、前后空白。
func UnmarshalPageModuleConfig(raw string) map[string]interface{} {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &out); err == nil && out != nil {
		return out
	}
	var inner string
	if err := json.Unmarshal([]byte(raw), &inner); err == nil {
		inner = strings.TrimSpace(inner)
		if inner != "" && json.Unmarshal([]byte(inner), &out) == nil && out != nil {
			return out
		}
	}
	return nil
}

func boolVal(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}

// heroSettingsPersist 写入 JSONB 时不用 omitempty，确保 false 也能落库
type heroSettingsPersist struct {
	Autoplay     bool   `json:"autoplay"`
	IntervalMs   int    `json:"interval_ms"`
	Transition   string `json:"transition"`
	ShowDots     bool   `json:"show_dots"`
	ShowArrows   bool   `json:"show_arrows"`
	PauseOnHover bool   `json:"pause_on_hover"`
}

func persistHeroSettings(s HeroSettings) heroSettingsPersist {
	return heroSettingsPersist{
		Autoplay:     boolVal(s.Autoplay, true),
		IntervalMs:   s.IntervalMs,
		Transition:   s.Transition,
		ShowDots:     boolVal(s.ShowDots, true),
		ShowArrows:   true,
		PauseOnHover: boolVal(s.PauseOnHover, true),
	}
}

func (s *CMSService) syncHeroThemeSettings(settings HeroSettings) {
	p := persistHeroSettings(settings)
	pairs := []struct{ key, val string }{
		{"hero_autoplay", strconv.FormatBool(p.Autoplay)},
		{"hero_interval_ms", strconv.Itoa(p.IntervalMs)},
		{"hero_transition", p.Transition},
		{"hero_show_dots", strconv.FormatBool(p.ShowDots)},
		{"hero_pause_on_hover", strconv.FormatBool(p.PauseOnHover)},
	}
	for _, item := range pairs {
		s.db.Model(&models.ThemeConfig{}).Where(`"key" = ?`, item.key).Update("value", item.val)
	}
}

// NormalizeHeroSettings 合并缺省值并校正非法字段
func NormalizeHeroSettings(in *HeroSettings) HeroSettings {
	out := DefaultHeroSettings()
	if in == nil {
		return out
	}
	if in.Autoplay != nil {
		out.Autoplay = in.Autoplay
	}
	if in.IntervalMs > 0 {
		ms := in.IntervalMs
		if ms < 2000 {
			ms = 2000
		}
		if ms > 30000 {
			ms = 30000
		}
		out.IntervalMs = ms
	}
	if in.Transition == "fade" || in.Transition == "slide" {
		out.Transition = in.Transition
	}
	if in.ShowDots != nil {
		out.ShowDots = in.ShowDots
	}
	if in.ShowArrows != nil {
		out.ShowArrows = in.ShowArrows
	}
	if in.PauseOnHover != nil {
		out.PauseOnHover = in.PauseOnHover
	}
	return out
}

// UpdateHeroSlides 更新页面轮播图配置（仅 upsert type=banner 模块，不影响页面其它模块）
// settings 可为 nil：此时保留已有 settings，若无则写入默认值。
func (s *CMSService) UpdateHeroSlides(pageID string, slides []HeroSlide, settings *HeroSettings) error {
	pid, err := uuid.Parse(pageID)
	if err != nil {
		return errors.New("页面 ID 无效")
	}
	if len(slides) == 0 {
		return errors.New("至少需要一张轮播图")
	}
	for i, slide := range slides {
		if strings.TrimSpace(slide.Image) == "" {
			return fmt.Errorf("第 %d 张轮播图缺少图片地址", i+1)
		}
	}

	var page models.Page
	if err := s.db.First(&page, "id = ?", pid).Error; err != nil {
		return errors.New("页面不存在")
	}

	var module models.PageModule
	err = s.db.Where("page_id = ? AND type = ?", pid, "banner").First(&module).Error

	finalSettings := DefaultHeroSettings()
	if err == nil && module.Config != "" {
		if existing := UnmarshalPageModuleConfig(module.Config); existing != nil {
			if raw, ok := existing["settings"]; ok {
				b, _ := json.Marshal(raw)
				var parsed HeroSettings
				if json.Unmarshal(b, &parsed) == nil {
					finalSettings = NormalizeHeroSettings(&parsed)
				}
			}
		}
	}
	if settings != nil {
		finalSettings = NormalizeHeroSettings(settings)
	}

	configJSON, err2 := json.Marshal(map[string]interface{}{
		"settings": persistHeroSettings(finalSettings),
		"slides":   slides,
	})
	if err2 != nil {
		return errors.New("配置序列化失败")
	}

	s.syncHeroThemeSettings(finalSettings)

	if err == nil {
		if uerr := s.db.Model(&module).Updates(map[string]interface{}{
			"config":     gorm.Expr("?::jsonb", string(configJSON)),
			"is_visible": true,
		}).Error; uerr != nil {
			return uerr
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.db.Create(&models.PageModule{
		PageID:    pid,
		Type:      "banner",
		Title:     "Hero Banner",
		SortOrder: 1,
		IsVisible: true,
		Config:    string(configJSON),
	}).Error
}

// DeletePage 删除页面（级联软删除模块与翻译）

func (s *CMSService) DeletePage(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("page_id = ?", id).Delete(&models.PageModule{}).Error; err != nil {
			return err
		}
		if err := tx.Where("page_id = ?", id).Delete(&models.PageTranslation{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Page{}, "id = ?", id).Error
	})
}

// PublishPage 发布页面
func (s *CMSService) PublishPage(id string) error {
	now := time.Now()
	return s.db.Model(&models.Page{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.ContentStatusPublished,
		"published_at": now,
	}).Error
}

// UnpublishPage 下线页面
func (s *CMSService) UnpublishPage(id string) error {
	return s.db.Model(&models.Page{}).Where("id = ?", id).Update("status", models.ContentStatusOffline).Error
}

// ==================== 导航管理 ====================

// ListNavigations 导航列表
func (s *CMSService) ListNavigations(navType string) ([]models.Navigation, error) {
	var navigations []models.Navigation
	query := s.db.Where("parent_id IS NULL")
	if navType != "" {
		query = query.Where("type = ?", navType)
	}
	err := query.Preload("Children").Preload("Page").Preload("Translations").Preload("Children.Translations").Order("sort_order ASC").Find(&navigations).Error
	return navigations, err
}

// CreateNavigation 创建导航
func (s *CMSService) CreateNavigation(req *NavigationRequest) (*models.Navigation, error) {
	nav := models.Navigation{
		Name:      req.Name,
		Type:      req.Type,
		URL:       req.URL,
		Target:    req.Target,
		SortOrder: req.SortOrder,
		IsVisible: req.IsVisible,
		ParentID:  utils.StringPtrToUUIDPtr(req.ParentID),
		PageID:    utils.StringPtrToUUIDPtr(req.PageID),
	}
	if err := s.db.Create(&nav).Error; err != nil {
		return nil, err
	}
	return &nav, nil
}

// NavigationRequest 导航请求
type NavigationRequest struct {
	Name         string             `json:"name" binding:"required"`
	Type         string             `json:"type"`
	URL          string             `json:"url"`
	Target       string             `json:"target"`
	SortOrder    int                `json:"sort_order"`
	IsVisible    bool               `json:"is_visible"`
	ParentID     *string            `json:"parent_id"`
	PageID       *string            `json:"page_id"` // 关联页面 ID
	Translations []TranslationInput `json:"translations"`
}

// UpdateNavigation 更新导航
func (s *CMSService) UpdateNavigation(id string, req *NavigationRequest) (*models.Navigation, error) {
	var nav models.Navigation
	if err := s.db.First(&nav, "id = ?", id).Error; err != nil {
		return nil, errors.New("导航不存在")
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"type":       req.Type,
		"url":        req.URL,
		"target":     req.Target,
		"sort_order": req.SortOrder,
		"is_visible": req.IsVisible,
	}
	if req.ParentID != nil {
		updates["parent_id"] = utils.StringPtrToUUIDPtr(req.ParentID)
	}
	// page_id always present: frontend sends string or null; null clears the association
	updates["page_id"] = utils.StringPtrToUUIDPtr(req.PageID)
	if err := s.db.Model(&nav).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.saveNavigationTranslations(tx, nav.ID, req.Translations)
	}); err != nil {
		return nil, err
	}
	return &nav, nil
}

// DeleteNavigation 删除导航（子级自动提升为顶级）
func (s *CMSService) DeleteNavigation(id string) error {
	pid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 子级导航的 parent_id 置 NULL，提升为顶级
		if err := tx.Model(&models.Navigation{}).Where("parent_id = ?", pid).Update("parent_id", nil).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Navigation{}, "id = ?", pid).Error
	})
}

// ListNavigationsByPageIDs 批量查询指定页面的导航关联（用于列表页状态列）
func (s *CMSService) ListNavigationsByPageIDs(pageIDs []uuid.UUID) ([]models.Navigation, error) {
	if len(pageIDs) == 0 {
		return nil, nil
	}
	var navs []models.Navigation
	err := s.db.Where("page_id IN ?", pageIDs).Find(&navs).Error
	return navs, err
}

// GetNavigationByPageID 查询某页面是否已关联导航（轻量判断）
func (s *CMSService) GetNavigationByPageID(pageID string) ([]models.Navigation, error) {
	var navs []models.Navigation
	pid, err := uuid.Parse(pageID)
	if err != nil {
		return nil, err
	}
	err = s.db.Where("page_id = ?", pid).Find(&navs).Error
	return navs, err
}

// SortNavigationRequest 批量排序请求
type SortNavigationRequest struct {
	Items []SortNavigationItem `json:"items" binding:"required"`
}

// SortNavigationItem 单条排序项
type SortNavigationItem struct {
	ID        string `json:"id" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

// BatchSortNavigations 批量更新导航排序
func (s *CMSService) BatchSortNavigations(req *SortNavigationRequest) error {
	for _, item := range req.Items {
		if err := s.db.Model(&models.Navigation{}).Where("id = ?", item.ID).Update("sort_order", item.SortOrder).Error; err != nil {
			return err
		}
	}
	return nil
}

// SyncNavVisibilityWithPages 同步导航可见性与页面状态
// 规则：关联页面已发布的导航保持可见，未发布/已下线页面的导航自动隐藏
// 返回：更新的导航数量
func (s *CMSService) SyncNavVisibilityWithPages() (int, error) {
	// 获取所有关联了页面的导航
	var navs []models.Navigation
	if err := s.db.Where("page_id IS NOT NULL").Find(&navs).Error; err != nil {
		return 0, err
	}

	// 收集所有关联的页面ID
	pageIDs := make([]uuid.UUID, 0, len(navs))
	for _, n := range navs {
		if n.PageID != nil {
			pageIDs = append(pageIDs, *n.PageID)
		}
	}
	if len(pageIDs) == 0 {
		return 0, nil
	}

	// 查询这些页面的发布状态
	var pages []models.Page
	if err := s.db.Select("id, status").Where("id IN ?", pageIDs).Find(&pages).Error; err != nil {
		return 0, err
	}

	// 构建页面状态映射
	pageStatusMap := make(map[uuid.UUID]string)
	for _, p := range pages {
		pageStatusMap[p.ID] = string(p.Status)
	}

	// 同步导航可见性
	updated := 0
	for _, n := range navs {
		if n.PageID == nil {
			continue
		}
		status, exists := pageStatusMap[*n.PageID]
		if !exists {
			// 页面不存在，隐藏导航
			if n.IsVisible {
				if err := s.db.Model(&models.Navigation{}).Where("id = ?", n.ID).Update("is_visible", false).Error; err != nil {
					return updated, err
				}
				updated++
			}
			continue
		}
		shouldVisible := status == "published"
		if n.IsVisible != shouldVisible {
			if err := s.db.Model(&models.Navigation{}).Where("id = ?", n.ID).Update("is_visible", shouldVisible).Error; err != nil {
				return updated, err
			}
			updated++
		}
	}

	return updated, nil
}

// ==================== 博客管理 ====================

// ListBlogs 博客列表
func (s *CMSService) ListBlogs(page, pageSize int, keyword, status, category string) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	query := s.db.Model(&models.Blog{})
	if keyword != "" {
		query = query.Where("title LIKE ? OR slug LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	err := query.Preload("Translations").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&blogs).Error

	return blogs, total, err
}

// GetBlog 获取博客
func (s *CMSService) GetBlog(id string) (*models.Blog, error) {
	var blog models.Blog
	err := s.db.Preload("Translations").First(&blog, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("博客不存在")
	}
	return &blog, nil
}

// GetBlogBySlug 通过 Slug 获取博客
func (s *CMSService) GetBlogBySlug(slug string) (*models.Blog, error) {
	var blog models.Blog
	err := s.db.Preload("Translations").
		First(&blog, "slug = ? AND status = ?", slug, models.ContentStatusPublished).Error
	if err != nil {
		return nil, errors.New("博客不存在")
	}
	return &blog, nil
}

// CreateBlog 创建博客
func (s *CMSService) CreateBlog(req *BlogRequest) (*models.Blog, error) {
	blog := models.Blog{
		Title:      req.Title,
		Slug:       req.Slug,
		Category:   req.Category,
		Tags:       req.Tags,
		Author:     req.Author,
		CoverImage: req.CoverImage,
		Content:    req.Content,
		Status:     req.Status,
	}

	if err := s.db.Create(&blog).Error; err != nil {
		return nil, err
	}

	// 创建翻译（跳过标题与正文均为空的空翻译，避免产生无效记录）
	for _, t := range req.Translations {
		if t.Title == "" && t.Content == "" {
			continue
		}
		s.db.Create(&models.BlogTranslation{
			BlogID:   blog.ID,
			Language: t.Language,
			Title:    t.Title,
			Content:  t.Content,
			Status:   models.TranslationStatusPublished,
		})
	}

	// 创建 SEO（未显式提供时，自动用标题 + 正文摘要生成默认值，幂等不覆盖）
	if req.SEO != nil {
		s.db.Create(&models.SEO{
			EntityType:    "blog",
			EntityID:      blog.ID,
			Language:      req.SEO.Language,
			Title:         req.SEO.Title,
			Description:   req.SEO.Description,
			Keywords:      req.SEO.Keywords,
			Canonical:     req.SEO.Canonical,
			Robots:        req.SEO.Robots,
			OGTitle:       req.SEO.OGTitle,
			OGDescription: req.SEO.OGDescription,
			OGImage:       req.SEO.OGImage,
			SchemaData:    req.SEO.SchemaData,
		})
	} else {
		if err := EnsureEntitySEO(s.db, "blog", blog.ID, req.Title, truncateText(req.Content, 200)); err != nil {
			return nil, err
		}
	}

	return s.GetBlog(blog.ID.String())
}

// BlogRequest 博客请求
type BlogRequest struct {
	Title        string                   `json:"title" binding:"required"`
	Slug         string                   `json:"slug" binding:"required"`
	Category     string                   `json:"category"`
	Tags         string                   `json:"tags"`
	Author       string                   `json:"author"`
	CoverImage   string                   `json:"cover_image"`
	Content      string                   `json:"content"`
	Status       models.ContentStatus     `json:"status"`
	Translations []BlogTranslationRequest `json:"translations"`
	SEO          *SEORequest              `json:"seo"`
}

// BlogTranslationRequest 博客翻译请求
type BlogTranslationRequest struct {
	Language string `json:"language" binding:"required"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

// UpdateBlog 更新博客
func (s *CMSService) UpdateBlog(id string, req *BlogRequest) (*models.Blog, error) {
	var blog models.Blog
	if err := s.db.First(&blog, "id = ?", id).Error; err != nil {
		return nil, errors.New("博客不存在")
	}

	// status 为空时不覆盖（编辑保存只传基础字段，避免清空发布状态；对齐 UpdatePage 语义）
	updates := map[string]interface{}{
		"title":       req.Title,
		"slug":        req.Slug,
		"category":    req.Category,
		"tags":        req.Tags,
		"author":      req.Author,
		"cover_image": req.CoverImage,
		"content":     req.Content,
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if err := s.db.Model(&blog).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 更新翻译：标题与正文均为空视为删除该语言翻译，否则 upsert
	for _, t := range req.Translations {
		if t.Title == "" && t.Content == "" {
			s.db.Where("blog_id = ? AND language = ?", blog.ID, t.Language).Delete(&models.BlogTranslation{})
			continue
		}
		var translation models.BlogTranslation
		err := s.db.Where("blog_id = ? AND language = ?", blog.ID, t.Language).First(&translation).Error
		if err != nil {
			s.db.Create(&models.BlogTranslation{
				BlogID:   blog.ID,
				Language: t.Language,
				Title:    t.Title,
				Content:  t.Content,
				Status:   models.TranslationStatusPublished,
			})
		} else {
			s.db.Model(&translation).Updates(map[string]interface{}{
				"title":   t.Title,
				"content": t.Content,
				"status":  models.TranslationStatusPublished,
			})
		}
	}

	// 更新 SEO
	if req.SEO != nil {
		var seo models.SEO
		err := s.db.Where("entity_type = ? AND entity_id = ?", "blog", blog.ID).First(&seo).Error
		if err != nil {
			s.db.Create(&models.SEO{
				EntityType:    "blog",
				EntityID:      blog.ID,
				Language:      req.SEO.Language,
				Title:         req.SEO.Title,
				Description:   req.SEO.Description,
				Keywords:      req.SEO.Keywords,
				Canonical:     req.SEO.Canonical,
				Robots:        req.SEO.Robots,
				OGTitle:       req.SEO.OGTitle,
				OGDescription: req.SEO.OGDescription,
				OGImage:       req.SEO.OGImage,
				SchemaData:    req.SEO.SchemaData,
			})
		} else {
			s.db.Model(&seo).Updates(map[string]interface{}{
				"language":       req.SEO.Language,
				"title":          req.SEO.Title,
				"description":    req.SEO.Description,
				"keywords":       req.SEO.Keywords,
				"canonical":      req.SEO.Canonical,
				"robots":         req.SEO.Robots,
				"og_title":       req.SEO.OGTitle,
				"og_description": req.SEO.OGDescription,
				"og_image":       req.SEO.OGImage,
				"schema_data":    req.SEO.SchemaData,
			})
		}
	}

	return s.GetBlog(id)
}

// DeleteBlog 删除博客（级联软删除翻译）
func (s *CMSService) DeleteBlog(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("blog_id = ?", id).Delete(&models.BlogTranslation{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Blog{}, "id = ?", id).Error
	})
}

// PublishBlog 发布博客
func (s *CMSService) PublishBlog(id string) error {
	now := time.Now()
	return s.db.Model(&models.Blog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.ContentStatusPublished,
		"published_at": now,
	}).Error
}

// UnpublishBlog 下线博客（前台不可见）
func (s *CMSService) UnpublishBlog(id string) error {
	return s.db.Model(&models.Blog{}).Where("id = ?", id).Update("status", models.ContentStatusOffline).Error
}

// ==================== 案例管理 ====================

// ListCases 案例列表（后台：全部状态）
func (s *CMSService) ListCases(page, pageSize int, keyword, projectType string) ([]models.Case, int64, error) {
	var cases []models.Case
	var total int64

	query := s.db.Model(&models.Case{})
	if keyword != "" {
		query = query.Where("title LIKE ? OR slug LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if projectType != "" {
		query = query.Where("project_type = ?", projectType)
	}

	query.Count(&total)
	err := query.Preload("Translations").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&cases).Error

	return cases, total, err
}

// ListPublishedCases 已发布案例列表（前台公开接口专用）
func (s *CMSService) ListPublishedCases(page, pageSize int, keyword, projectType string) ([]models.Case, int64, error) {
	var cases []models.Case
	var total int64

	query := s.db.Model(&models.Case{}).Where("status = ?", models.ContentStatusPublished)
	if keyword != "" {
		query = query.Where("title LIKE ? OR slug LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if projectType != "" {
		query = query.Where("project_type = ?", projectType)
	}

	query.Count(&total)
	err := query.Preload("Translations").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&cases).Error

	return cases, total, err
}

// GetCase 获取案例
func (s *CMSService) GetCase(id string) (*models.Case, error) {
	var caseItem models.Case
	err := s.db.Preload("Translations").First(&caseItem, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("案例不存在")
	}
	s.loadCaseSEO(&caseItem)
	return &caseItem, nil
}

// GetCaseBySlug 通过 Slug 获取已发布案例（前台公开接口专用）
func (s *CMSService) GetCaseBySlug(slug string) (*models.Case, error) {
	var caseItem models.Case
	err := s.db.Preload("Translations").
		First(&caseItem, "slug = ? AND status = ?", slug, models.ContentStatusPublished).Error
	if err != nil {
		return nil, errors.New("案例不存在")
	}
	s.loadCaseSEO(&caseItem)
	return &caseItem, nil
}

// loadCaseSEO 加载案例 SEO 配置（gorm:"-" 字段需手动查询）
func (s *CMSService) loadCaseSEO(caseItem *models.Case) {
	var seo models.SEO
	if err := s.db.Where("entity_type = ? AND entity_id = ?", "case", caseItem.ID).First(&seo).Error; err == nil {
		caseItem.SEO = &seo
	}
}

// CreateCase 创建案例
func (s *CMSService) CreateCase(req *CaseRequest) (*models.Case, error) {
	caseItem := models.Case{
		Title:          req.Title,
		Slug:           req.Slug,
		ClientIndustry: req.ClientIndustry,
		ProjectType:    req.ProjectType,
		Products:       req.Products,
		ClientNeed:     req.ClientNeed,
		Problem:        req.Problem,
		Solution:       req.Solution,
		Process:        req.Process,
		Result:         req.Result,
		CoverImage:     req.CoverImage,
		Status:         req.Status,
	}

	if err := s.db.Create(&caseItem).Error; err != nil {
		return nil, err
	}

	// 创建翻译
	if len(req.Translations) > 0 {
		for _, t := range req.Translations {
			s.db.Create(&models.CaseTranslation{
				CaseID:     caseItem.ID,
				Language:   t.Language,
				Title:      t.Title,
				ClientNeed: t.ClientNeed,
				Problem:    t.Problem,
				Solution:   t.Solution,
				Process:    t.Process,
				Result:     t.Result,
				Status:     models.TranslationStatusPublished,
			})
		}
	}

	// 创建 SEO（未显式提供时，自动用标题 + 摘要生成默认值，幂等不覆盖）
	if req.SEO != nil {
		s.db.Create(&models.SEO{
			EntityType:    "case",
			EntityID:      caseItem.ID,
			Language:      req.SEO.Language,
			Title:         req.SEO.Title,
			Description:   req.SEO.Description,
			Keywords:      req.SEO.Keywords,
			Canonical:     req.SEO.Canonical,
			Robots:        req.SEO.Robots,
			OGTitle:       req.SEO.OGTitle,
			OGDescription: req.SEO.OGDescription,
			OGImage:       req.SEO.OGImage,
			SchemaData:    req.SEO.SchemaData,
		})
	} else {
		summary := req.Solution
		if strings.TrimSpace(summary) == "" {
			summary = req.Problem
		}
		if err := EnsureEntitySEO(s.db, "case", caseItem.ID, req.Title, truncateText(summary, 200)); err != nil {
			return nil, err
		}
	}

	return s.GetCase(caseItem.ID.String())
}

// CaseRequest 案例请求
type CaseRequest struct {
	Title          string                   `json:"title" binding:"required"`
	Slug           string                   `json:"slug" binding:"required"`
	ClientIndustry string                   `json:"client_industry"`
	ProjectType    string                   `json:"project_type"`
	Products       string                   `json:"products"`
	ClientNeed     string                   `json:"client_need"`
	Problem        string                   `json:"problem"`
	Solution       string                   `json:"solution"`
	Process        string                   `json:"process"`
	Result         string                   `json:"result"`
	CoverImage     string                   `json:"cover_image"`
	Status         models.ContentStatus     `json:"status"`
	Translations   []CaseTranslationRequest `json:"translations"`
	SEO            *SEORequest              `json:"seo"`
}

// CaseTranslationRequest 案例翻译请求
type CaseTranslationRequest struct {
	Language   string `json:"language" binding:"required"`
	Title      string `json:"title"`
	ClientNeed string `json:"client_need"`
	Problem    string `json:"problem"`
	Solution   string `json:"solution"`
	Process    string `json:"process"`
	Result     string `json:"result"`
}

// UpdateCase 更新案例
func (s *CMSService) UpdateCase(id string, req *CaseRequest) (*models.Case, error) {
	var caseItem models.Case
	if err := s.db.First(&caseItem, "id = ?", id).Error; err != nil {
		return nil, errors.New("案例不存在")
	}

	// status 为空时不覆盖（编辑保存只传基础字段，避免清空发布状态；对齐 UpdatePage 语义）
	updates := map[string]interface{}{
		"title":           req.Title,
		"slug":            req.Slug,
		"client_industry": req.ClientIndustry,
		"project_type":    req.ProjectType,
		"products":        req.Products,
		"client_need":     req.ClientNeed,
		"problem":         req.Problem,
		"solution":        req.Solution,
		"process":         req.Process,
		"result":          req.Result,
		"cover_image":     req.CoverImage,
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if err := s.db.Model(&caseItem).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 更新翻译：先删后建（对齐 Page 模块重建语义），保证被移除/清空的语种彻底删除
	if err := s.db.Where("case_id = ?", caseItem.ID).Delete(&models.CaseTranslation{}).Error; err != nil {
		return nil, err
	}
	for _, t := range req.Translations {
		if t.Language == "" {
			continue
		}
		s.db.Create(&models.CaseTranslation{
			CaseID:     caseItem.ID,
			Language:   t.Language,
			Title:      t.Title,
			ClientNeed: t.ClientNeed,
			Problem:    t.Problem,
			Solution:   t.Solution,
			Process:    t.Process,
			Result:     t.Result,
			Status:     models.TranslationStatusPublished,
		})
	}

	// 更新 SEO
	if req.SEO != nil {
		var seo models.SEO
		err := s.db.Where("entity_type = ? AND entity_id = ?", "case", caseItem.ID).First(&seo).Error
		if err != nil {
			s.db.Create(&models.SEO{
				EntityType:    "case",
				EntityID:      caseItem.ID,
				Language:      req.SEO.Language,
				Title:         req.SEO.Title,
				Description:   req.SEO.Description,
				Keywords:      req.SEO.Keywords,
				Canonical:     req.SEO.Canonical,
				Robots:        req.SEO.Robots,
				OGTitle:       req.SEO.OGTitle,
				OGDescription: req.SEO.OGDescription,
				OGImage:       req.SEO.OGImage,
				SchemaData:    req.SEO.SchemaData,
			})
		} else {
			s.db.Model(&seo).Updates(map[string]interface{}{
				"language":       req.SEO.Language,
				"title":          req.SEO.Title,
				"description":    req.SEO.Description,
				"keywords":       req.SEO.Keywords,
				"canonical":      req.SEO.Canonical,
				"robots":         req.SEO.Robots,
				"og_title":       req.SEO.OGTitle,
				"og_description": req.SEO.OGDescription,
				"og_image":       req.SEO.OGImage,
				"schema_data":    req.SEO.SchemaData,
			})
		}
	}

	return s.GetCase(id)
}

// DeleteCase 删除案例（级联软删除翻译）
func (s *CMSService) DeleteCase(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("case_id = ?", id).Delete(&models.CaseTranslation{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Case{}, "id = ?", id).Error
	})
}

// PublishCase 发布案例
func (s *CMSService) PublishCase(id string) error {
	now := time.Now()
	return s.db.Model(&models.Case{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.ContentStatusPublished,
		"published_at": now,
	}).Error
}

// UnpublishCase 下线案例（前台不可见）
func (s *CMSService) UnpublishCase(id string) error {
	return s.db.Model(&models.Case{}).Where("id = ?", id).Update("status", models.ContentStatusOffline).Error
}

// ==================== FAQ 管理 ====================

// ListFAQs FAQ 列表（后台管理：一条 FAQ 一条记录，多语言在 translations 中）
func (s *CMSService) ListFAQs(page, pageSize int, category, keyword string) ([]models.FAQ, int64, error) {
	var faqs []models.FAQ
	var total int64

	query := s.db.Model(&models.FAQ{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("question LIKE ? OR answer LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Preload("Translations").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("sort_order ASC").Find(&faqs).Error

	return faqs, total, err
}

// ListPublishedFAQs 已启用 FAQ 列表（前台公开接口专用，语言由 LocalizeFAQ 覆盖）
func (s *CMSService) ListPublishedFAQs(page, pageSize int, category string) ([]models.FAQ, int64, error) {
	var faqs []models.FAQ
	var total int64

	query := s.db.Model(&models.FAQ{}).Where("is_active = ?", true)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	err := query.Preload("Translations").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("sort_order ASC").Find(&faqs).Error

	return faqs, total, err
}

// CreateFAQ 创建 FAQ（主表存英文源，其他语言写入翻译表）
func (s *CMSService) CreateFAQ(req *FAQRequest) (*models.FAQ, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	faq := models.FAQ{
		Question:  req.Question,
		Answer:    req.Answer,
		Category:  req.Category,
		Language:  "en",
		SortOrder: req.SortOrder,
		IsActive:  isActive,
	}
	if err := s.db.Create(&faq).Error; err != nil {
		return nil, err
	}

	for _, t := range req.Translations {
		if t.Language == "" || t.Language == "en" {
			continue
		}
		if t.Question == "" && t.Answer == "" {
			continue
		}
		s.db.Create(&models.FAQTranslation{
			FAQID:    faq.ID,
			Language: t.Language,
			Question: t.Question,
			Answer:   t.Answer,
			Status:   models.TranslationStatusPublished,
		})
	}
	return &faq, nil
}

// FAQRequest FAQ 请求
type FAQRequest struct {
	Question     string                  `json:"question" binding:"required"`
	Answer       string                  `json:"answer"`
	Category     string                  `json:"category"`
	Language     string                  `json:"language"` // 兼容保留；主表固定英文源
	SortOrder    int                     `json:"sort_order"`
	IsActive     *bool                   `json:"is_active"`
	Translations []FAQTranslationRequest `json:"translations"`
}

// FAQTranslationRequest FAQ 翻译请求
type FAQTranslationRequest struct {
	Language string `json:"language" binding:"required"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// UpdateFAQ 更新 FAQ（主表英文源 + 重建翻译表）
func (s *CMSService) UpdateFAQ(id string, req *FAQRequest) (*models.FAQ, error) {
	var faq models.FAQ
	if err := s.db.First(&faq, "id = ?", id).Error; err != nil {
		return nil, errors.New("FAQ 不存在")
	}

	updates := map[string]interface{}{
		"question":   req.Question,
		"answer":     req.Answer,
		"category":   req.Category,
		"language":   "en",
		"sort_order": req.SortOrder,
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if err := s.db.Model(&faq).Updates(updates).Error; err != nil {
		return nil, err
	}

	if req.Translations != nil {
		if err := s.db.Where("faq_id = ?", faq.ID).Delete(&models.FAQTranslation{}).Error; err != nil {
			return nil, err
		}
		for _, t := range req.Translations {
			if t.Language == "" || t.Language == "en" {
				continue
			}
			if t.Question == "" && t.Answer == "" {
				continue
			}
			s.db.Create(&models.FAQTranslation{
				FAQID:    faq.ID,
				Language: t.Language,
				Question: t.Question,
				Answer:   t.Answer,
				Status:   models.TranslationStatusPublished,
			})
		}
	}
	return &faq, nil
}

// DeleteFAQ 删除 FAQ（级联软删除翻译）
func (s *CMSService) DeleteFAQ(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("faq_id = ?", id).Delete(&models.FAQTranslation{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.FAQ{}, "id = ?", id).Error
	})
}

// ==================== 工厂管理 ====================

// ListFactories 工厂列表（后台：全部状态，服务端分页 + 名称搜索）
func (s *CMSService) ListFactories(page, pageSize int, keyword string) ([]models.Factory, int64, error) {
	var factories []models.Factory
	var total int64

	query := s.db.Model(&models.Factory{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("name ASC").
		Offset((page - 1) * pageSize).Limit(pageSize).Preload("Translations").Find(&factories).Error

	return factories, total, err
}

// ListPublishedFactories 已发布工厂列表（前台）
func (s *CMSService) ListPublishedFactories() ([]models.Factory, error) {
	var factories []models.Factory
	err := s.db.Where("status = ?", models.ProductStatusPublished).
		Order("name ASC").Preload("Translations").Find(&factories).Error
	return factories, err
}

// PublishFactory 发布工厂（前端可见）
func (s *CMSService) PublishFactory(id string) error {
	return s.db.Model(&models.Factory{}).Where("id = ?", id).
		Update("status", models.ProductStatusPublished).Error
}

// UnpublishFactory 下线工厂（前端不可见）
func (s *CMSService) UnpublishFactory(id string) error {
	return s.db.Model(&models.Factory{}).Where("id = ?", id).
		Update("status", models.ProductStatusOffline).Error
}

// CreateFactory 创建工厂
func (s *CMSService) CreateFactory(req *FactoryRequest) (*models.Factory, error) {
	factory := models.Factory{
		Name:                 req.Name,
		Status:               req.Status,
		Location:             req.Location,
		Area:                 req.Area,
		Employees:            req.Employees,
		ProductionLines:      req.ProductionLines,
		Equipment:            req.Equipment,
		MonthlyCapacity:      req.MonthlyCapacity,
		AnnualCapacity:       req.AnnualCapacity,
		Warehouse:            req.Warehouse,
		QualityManagement:    req.QualityManagement,
		ProductionCapability: req.ProductionCapability,
		Description:          req.Description,
		Image:                req.Image,
		IsActive:             req.IsActive,
	}
	if err := s.db.Create(&factory).Error; err != nil {
		return nil, err
	}
	return &factory, nil
}

// FactoryRequest 工厂请求
type FactoryRequest struct {
	Name                 string               `json:"name" binding:"required"`
	Location             string               `json:"location"`
	Area                 string               `json:"area"`
	Employees            int                  `json:"employees"`
	ProductionLines      int                  `json:"production_lines"`
	Equipment            string               `json:"equipment"`
	MonthlyCapacity      string               `json:"monthly_capacity"`
	AnnualCapacity       string               `json:"annual_capacity"`
	Warehouse            string               `json:"warehouse"`
	QualityManagement    string               `json:"quality_management"`
	ProductionCapability string               `json:"production_capability"`
	Description          string               `json:"description"`
	Image                string               `json:"image"`
	Status               models.ProductStatus `json:"status"`
	IsActive             bool                 `json:"is_active"`
	Translations         []TranslationInput   `json:"translations"`
}

// UpdateFactory 更新工厂
func (s *CMSService) UpdateFactory(id string, req *FactoryRequest) (*models.Factory, error) {
	var factory models.Factory
	if err := s.db.First(&factory, "id = ?", id).Error; err != nil {
		return nil, errors.New("工厂不存在")
	}

	updates := map[string]interface{}{
		"name":                  req.Name,
		"location":              req.Location,
		"area":                  req.Area,
		"employees":             req.Employees,
		"production_lines":      req.ProductionLines,
		"equipment":             req.Equipment,
		"monthly_capacity":      req.MonthlyCapacity,
		"annual_capacity":       req.AnnualCapacity,
		"warehouse":             req.Warehouse,
		"quality_management":    req.QualityManagement,
		"production_capability": req.ProductionCapability,
		"description":           req.Description,
		"image":                 req.Image,
		"is_active":             req.IsActive,
	}
	if err := s.db.Model(&factory).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.saveFactoryTranslations(tx, factory.ID, req.Translations)
	}); err != nil {
		return nil, err
	}
	return &factory, nil
}

// DeleteFactory 删除工厂
func (s *CMSService) DeleteFactory(id string) error {
	return s.db.Delete(&models.Factory{}, "id = ?", id).Error
}

// ==================== 认证管理 ====================

// ListCertifications 认证列表（后台：全部状态，服务端分页 + 名称/编号搜索）
func (s *CMSService) ListCertifications(page, pageSize int, keyword string) ([]models.Certification, int64, error) {
	var certifications []models.Certification
	var total int64

	query := s.db.Model(&models.Certification{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("name ASC").
		Offset((page - 1) * pageSize).Limit(pageSize).Preload("Translations").Find(&certifications).Error

	return certifications, total, err
}

// ListPublishedCertifications 已发布认证列表（前台）
func (s *CMSService) ListPublishedCertifications() ([]models.Certification, error) {
	var certifications []models.Certification
	err := s.db.Where("status = ?", models.ProductStatusPublished).
		Order("name ASC").Preload("Translations").Find(&certifications).Error
	return certifications, err
}

// PublishCertification 发布认证（前端可见）
func (s *CMSService) PublishCertification(id string) error {
	return s.db.Model(&models.Certification{}).Where("id = ?", id).
		Update("status", models.ProductStatusPublished).Error
}

// UnpublishCertification 下线认证（前端不可见）
func (s *CMSService) UnpublishCertification(id string) error {
	return s.db.Model(&models.Certification{}).Where("id = ?", id).
		Update("status", models.ProductStatusOffline).Error
}

// CreateCertification 创建认证
func (s *CMSService) CreateCertification(req *CertificationRequest) (*models.Certification, error) {
	cert := models.Certification{
		Name:        req.Name,
		Status:      req.Status,
		Code:        req.Code,
		IssueDate:   req.IssueDate,
		ExpiryDate:  req.ExpiryDate,
		Image:       req.Image,
		PDF:         req.PDF,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
	if err := s.db.Create(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

// CertificationRequest 认证请求
type CertificationRequest struct {
	Name         string               `json:"name" binding:"required"`
	Code         string               `json:"code"`
	IssueDate    string               `json:"issue_date"`
	ExpiryDate   string               `json:"expiry_date"`
	Image        string               `json:"image"`
	PDF          string               `json:"pdf"`
	Description  string               `json:"description"`
	Status       models.ProductStatus `json:"status"`
	IsActive     bool                 `json:"is_active"`
	Translations []TranslationInput   `json:"translations"`
}

// UpdateCertification 更新认证
func (s *CMSService) UpdateCertification(id string, req *CertificationRequest) (*models.Certification, error) {
	var cert models.Certification
	if err := s.db.First(&cert, "id = ?", id).Error; err != nil {
		return nil, errors.New("认证不存在")
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"code":        req.Code,
		"issue_date":  req.IssueDate,
		"expiry_date": req.ExpiryDate,
		"image":       req.Image,
		"pdf":         req.PDF,
		"description": req.Description,
		"is_active":   req.IsActive,
	}
	if err := s.db.Model(&cert).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.saveCertificationTranslations(tx, cert.ID, req.Translations)
	}); err != nil {
		return nil, err
	}
	return &cert, nil
}

// DeleteCertification 删除认证
func (s *CMSService) DeleteCertification(id string) error {
	return s.db.Delete(&models.Certification{}, "id = ?", id).Error
}

// ==================== 生产流程管理 ====================

// ListProductionProcesses 生产流程列表（后台：全部状态，服务端分页 + 名称搜索）
func (s *CMSService) ListProductionProcesses(page, pageSize int, keyword string) ([]models.ProductionProcess, int64, error) {
	var processes []models.ProductionProcess
	var total int64

	query := s.db.Model(&models.ProductionProcess{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("sort_order ASC").
		Offset((page - 1) * pageSize).Limit(pageSize).Preload("Translations").Find(&processes).Error

	return processes, total, err
}

// ListPublishedProductionProcesses 已发布生产流程列表（前台）
func (s *CMSService) ListPublishedProductionProcesses() ([]models.ProductionProcess, error) {
	var processes []models.ProductionProcess
	err := s.db.Where("status = ?", models.ProductStatusPublished).
		Order("sort_order ASC").Preload("Translations").Find(&processes).Error
	return processes, err
}

// PublishProductionProcess 发布生产流程（前端可见）
func (s *CMSService) PublishProductionProcess(id string) error {
	return s.db.Model(&models.ProductionProcess{}).Where("id = ?", id).
		Update("status", models.ProductStatusPublished).Error
}

// UnpublishProductionProcess 下线生产流程（前端不可见）
func (s *CMSService) UnpublishProductionProcess(id string) error {
	return s.db.Model(&models.ProductionProcess{}).Where("id = ?", id).
		Update("status", models.ProductStatusOffline).Error
}

// CreateProductionProcess 创建生产流程
func (s *CMSService) CreateProductionProcess(req *ProductionProcessRequest) (*models.ProductionProcess, error) {
	p := models.ProductionProcess{
		Name:        req.Name,
		Status:      req.Status,
		Description: req.Description,
		Image:       req.Image,
		Video:       req.Video,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	}
	if err := s.db.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// ProductionProcessRequest 生产流程请求
type ProductionProcessRequest struct {
	Name         string               `json:"name" binding:"required"`
	Description  string               `json:"description"`
	Image        string               `json:"image"`
	Video        string               `json:"video"`
	SortOrder    int                  `json:"sort_order"`
	Status       models.ProductStatus `json:"status"`
	IsActive     bool                 `json:"is_active"`
	Translations []TranslationInput   `json:"translations"`
}

// UpdateProductionProcess 更新生产流程
func (s *CMSService) UpdateProductionProcess(id string, req *ProductionProcessRequest) (*models.ProductionProcess, error) {
	var p models.ProductionProcess
	if err := s.db.First(&p, "id = ?", id).Error; err != nil {
		return nil, errors.New("生产流程不存在")
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"image":       req.Image,
		"video":       req.Video,
		"sort_order":  req.SortOrder,
		"is_active":   req.IsActive,
	}
	if err := s.db.Model(&p).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return s.saveProcessTranslations(tx, p.ID, req.Translations)
	}); err != nil {
		return nil, err
	}
	return &p, nil
}

// DeleteProductionProcess 删除生产流程
func (s *CMSService) DeleteProductionProcess(id string) error {
	return s.db.Delete(&models.ProductionProcess{}, "id = ?", id).Error
}

// ==================== 自媒体管理 ====================

// SelfMediaMaxDisplay 读取「自媒体最多展示数量」配置（theme_configs），缺省 6
func (s *CMSService) SelfMediaMaxDisplay() int {
	const defaultMax = 6
	var cfg models.ThemeConfig
	if err := s.db.Where("\"key\" = ?", "self_media_max_display").First(&cfg).Error; err != nil {
		return defaultMax
	}
	// Value 可能是 `6` 或 `"6"`（历史主题配置可能以 JSON 字符串存储），做兼容解析
	raw := strings.TrimSpace(cfg.Value)
	raw = strings.Trim(raw, `"`)
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return defaultMax
	}
	return n
}

// SetSelfMediaMaxDisplay 保存「自媒体最多展示数量」配置（theme_configs，is_active=false 以隐藏于主题配置表）
func (s *CMSService) SetSelfMediaMaxDisplay(n int) error {
	if n < 1 {
		n = 1
	}
	var cfg models.ThemeConfig
	err := s.db.Where("\"key\" = ?", "self_media_max_display").First(&cfg).Error
	if err != nil {
		cfg = models.ThemeConfig{
			Key:      "self_media_max_display",
			Value:    strconv.Itoa(n),
			Group:    "self_media",
			Name:     "自媒体最多展示数量",
			IsActive: false,
		}
		return s.db.Create(&cfg).Error
	}
	return s.db.Model(&cfg).Update("value", strconv.Itoa(n)).Error
}

// ListSelfMedias 自媒体列表（后台：全部状态，服务端分页 + 名称/平台/账号搜索）
func (s *CMSService) ListSelfMedias(page, pageSize int, keyword string) ([]models.SelfMedia, int64, error) {
	var items []models.SelfMedia
	var total int64

	query := s.db.Model(&models.SelfMedia{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR platform LIKE ? OR account LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("sort_order ASC, created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error

	return items, total, err
}

// ListPublishedSelfMedias 已发布自媒体列表（前台，按排序截取 limit 条）
func (s *CMSService) ListPublishedSelfMedias(limit int) ([]models.SelfMedia, error) {
	var items []models.SelfMedia
	q := s.db.Where("status = ?", models.ProductStatusPublished).
		Order("sort_order ASC, created_at ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	// 计算前台跳转链接：自定义 URL > 平台主页 + 账号 > 平台主页（账号未配置时跳平台主页）
	for i := range items {
		items[i].URL = items[i].EffectiveURL()
	}
	return items, nil
}

// PublishSelfMedia 发布自媒体（前端可见）
func (s *CMSService) PublishSelfMedia(id string) error {
	return s.db.Model(&models.SelfMedia{}).Where("id = ?", id).
		Update("status", models.ProductStatusPublished).Error
}

// UnpublishSelfMedia 下线自媒体（前端不可见）
func (s *CMSService) UnpublishSelfMedia(id string) error {
	return s.db.Model(&models.SelfMedia{}).Where("id = ?", id).
		Update("status", models.ProductStatusOffline).Error
}

// CreateSelfMedia 创建自媒体
func (s *CMSService) CreateSelfMedia(req *SelfMediaRequest) (*models.SelfMedia, error) {
	item := models.SelfMedia{
		Name:        req.Name,
		Platform:    req.Platform,
		Account:     req.Account,
		URL:         req.URL,
		Image:       req.Image,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		IsActive:    req.IsActive,
	}
	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// SelfMediaRequest 自媒体请求
type SelfMediaRequest struct {
	Name        string               `json:"name" binding:"required"`
	Platform    string               `json:"platform"`
	Account     string               `json:"account"`
	URL         string               `json:"url"`
	Image       string               `json:"image"`
	Description string               `json:"description"`
	SortOrder   int                  `json:"sort_order"`
	Status      models.ProductStatus `json:"status"`
	IsActive    bool                 `json:"is_active"`
}

// UpdateSelfMedia 更新自媒体
func (s *CMSService) UpdateSelfMedia(id string, req *SelfMediaRequest) (*models.SelfMedia, error) {
	var item models.SelfMedia
	if err := s.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, errors.New("自媒体不存在")
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"platform":    req.Platform,
		"account":     req.Account,
		"url":         req.URL,
		"image":       req.Image,
		"description": req.Description,
		"sort_order":  req.SortOrder,
		"is_active":   req.IsActive,
	}
	if err := s.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// DeleteSelfMedia 删除自媒体
func (s *CMSService) DeleteSelfMedia(id string) error {
	return s.db.Delete(&models.SelfMedia{}, "id = ?", id).Error
}

// TranslationInput 翻译输入（通用：各实体按需使用字段）
type TranslationInput struct {
	Language             string `json:"language"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	QualityManagement    string `json:"quality_management"`
	ProductionCapability string `json:"production_capability"`
}

func (s *CMSService) saveFactoryTranslations(tx *gorm.DB, id uuid.UUID, items []TranslationInput) error {
	if err := tx.Where("factory_id = ?", id).Delete(&models.FactoryTranslation{}).Error; err != nil {
		return err
	}
	for _, it := range items {
		if it.Language == "" || it.Language == "en" {
			continue
		}
		if it.Name == "" && it.Description == "" && it.QualityManagement == "" && it.ProductionCapability == "" {
			continue
		}
		if err := tx.Create(&models.FactoryTranslation{FactoryID: id, Language: it.Language, Name: it.Name, Description: it.Description, QualityManagement: it.QualityManagement, ProductionCapability: it.ProductionCapability, Status: models.TranslationStatusPublished}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *CMSService) saveCertificationTranslations(tx *gorm.DB, id uuid.UUID, items []TranslationInput) error {
	if err := tx.Where("certification_id = ?", id).Delete(&models.CertificationTranslation{}).Error; err != nil {
		return err
	}
	for _, it := range items {
		if it.Language == "" || it.Language == "en" {
			continue
		}
		if it.Name == "" && it.Description == "" {
			continue
		}
		if err := tx.Create(&models.CertificationTranslation{CertificationID: id, Language: it.Language, Name: it.Name, Description: it.Description, Status: models.TranslationStatusPublished}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *CMSService) saveProcessTranslations(tx *gorm.DB, id uuid.UUID, items []TranslationInput) error {
	if err := tx.Where("process_id = ?", id).Delete(&models.ProductionProcessTranslation{}).Error; err != nil {
		return err
	}
	for _, it := range items {
		if it.Language == "" || it.Language == "en" {
			continue
		}
		if it.Name == "" && it.Description == "" {
			continue
		}
		if err := tx.Create(&models.ProductionProcessTranslation{ProcessID: id, Language: it.Language, Name: it.Name, Description: it.Description, Status: models.TranslationStatusPublished}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *CMSService) saveNavigationTranslations(tx *gorm.DB, id uuid.UUID, items []TranslationInput) error {
	if err := tx.Where("navigation_id = ?", id).Delete(&models.NavigationTranslation{}).Error; err != nil {
		return err
	}
	for _, it := range items {
		if it.Language == "" || it.Language == "en" {
			continue
		}
		if it.Name == "" {
			continue
		}
		if err := tx.Create(&models.NavigationTranslation{NavigationID: id, Language: it.Language, Name: it.Name, Status: models.TranslationStatusPublished}).Error; err != nil {
			return err
		}
	}
	return nil
}
