package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/models"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// PublicHandler 公开 API 处理器
type PublicHandler struct {
	productService *services.ProductService
	cmsService     *services.CMSService
	leadService    *services.LeadService
	cache          *services.CacheService
}

// NewPublicHandler 创建公开 API 处理器
func NewPublicHandler(
	productService *services.ProductService,
	cmsService *services.CMSService,
	leadService *services.LeadService,
	cacheService *services.CacheService,
) *PublicHandler {
	return &PublicHandler{
		productService: productService,
		cmsService:     cmsService,
		leadService:    leadService,
		cache:          cacheService,
	}
}

// cached 缓存旁路辅助（cache-aside）：
// 优先读取 Redis 缓存，命中直接返回；未命中时执行 loader 回源 DB，并将结果写入缓存。
// Redis 不可用时自动降级为直查 DB（CacheService.Get/Set 内部已处理降级）。
// useCache=false 时完全跳过缓存读写（预览模式），始终回源 DB 获取最新数据。
func cached[T any](h *PublicHandler, key string, ttl time.Duration, useCache bool, load func() (T, error)) (T, error) {
	var zero T
	if useCache && h.cache != nil && h.cache.IsEnabled() {
		var got T
		if h.cache.Get(key, &got) {
			return got, nil
		}
	}
	v, err := load()
	if err != nil {
		return zero, err
	}
	if useCache && h.cache != nil && h.cache.IsEnabled() {
		h.cache.Set(key, v, ttl)
	}
	return v, nil
}

// previewMode 是否处于"预览"模式（URL 带 preview=1/true）。
// 管理后台"门户预览"强制绕过缓存直查 DB，确保预览始终反映数据库最新状态，
// 不依赖缓存失效时序，也不污染线上缓存。
func (h *PublicHandler) previewMode(c *gin.Context) bool {
	switch c.Query("preview") {
	case "1", "true", "yes":
		return true
	}
	return false
}

// cachedPage 缓存分页列表结果的载体（Redis 中保存 items + total）
type cachedPage[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}

// GetHome 获取首页数据（并发查询优化 + Redis 缓存）
func (h *PublicHandler) GetHome(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:home:" + lang

	// 首页属于高频聚合查询，命中缓存直接返回，未命中回源 DB 并写缓存
	data, err := cached[gin.H](h, key, services.CacheTTLShort, !h.previewMode(c), func() (gin.H, error) {
		// 获取首页页面（必须成功）
		page, err := h.cmsService.GetPageBySlug("home")
		if err != nil {
			return nil, errors.New("首页不存在")
		}
		services.LocalizePage(page, lang)

		// 并发查询其余数据（goroutine + WaitGroup）
		var (
			wg             sync.WaitGroup
			products       []models.Product
			categories     []models.Category
			blogs          []models.Blog
			cases          []models.Case
			certifications []models.Certification
			processes      []models.ProductionProcess
			factories      []models.Factory
		)

		wg.Add(7)
		go func() { defer wg.Done(); products, _ = h.productService.ListFeaturedProducts(8) }()
		go func() { defer wg.Done(); categories, _ = h.productService.ListCategories() }()
		go func() { defer wg.Done(); blogs, _, _ = h.cmsService.ListBlogs(1, 4, "", "published") }()
		go func() { defer wg.Done(); cases, _, _ = h.cmsService.ListPublishedCases(1, 4, "") }()
		go func() { defer wg.Done(); certifications, _ = h.cmsService.ListPublishedCertifications() }()
		go func() { defer wg.Done(); processes, _ = h.cmsService.ListPublishedProductionProcesses() }()
		go func() { defer wg.Done(); factories, _ = h.cmsService.ListPublishedFactories() }()
		wg.Wait()

		// 本地化处理
		services.LocalizeProducts(products, lang)
		services.LocalizeBlogs(blogs, lang)
		services.LocalizeCases(cases, lang)

		return gin.H{
			"page":                 page,
			"hero_slides":          extractHeroSlides(page),
			"hero_settings":        extractHeroSettings(page),
			"featured_products":    products,
			"categories":           categories,
			"blogs":                blogs,
			"cases":                cases,
			"certifications":       certifications,
			"production_processes": processes,
			"factories":            factories,
		}, nil
	})
	if err != nil {
		utils.NotFound(c, "首页不存在")
		return
	}
	utils.Success(c, data)
}

// extractHeroSlides 从首页页面的 banner 模块中提取轮播图（hero）配置。
// 兼容两种 config 结构：
//   - 新版："slides": [{"image","title","subtitle","button_text","button_url"}, ...]
//   - 旧版：单对象 {"image","title","subtitle","button_text","button_url"}（自动包装为单张轮播图）
// 返回统一结构供门户 / 管理后台门户预览渲染轮播图；无有效配置时返回 nil。
func extractHeroSlides(page *models.Page) []gin.H {
	if page == nil {
		return nil
	}
	for _, m := range page.Modules {
		if m.Type != "banner" || !m.IsVisible {
			continue
		}
		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(m.Config), &raw); err != nil {
			continue
		}
		// 新版：slides 数组
		if list, ok := raw["slides"].([]interface{}); ok {
			var slides []gin.H
			for _, it := range list {
				obj, ok := it.(map[string]interface{})
				if !ok {
					continue
				}
				slides = append(slides, slideFromMap(obj))
			}
			if len(slides) > 0 {
				return slides
			}
		} else if _, hasImage := raw["image"]; hasImage {
			// 旧版：单张配置
			return []gin.H{slideFromMap(raw)}
		}
	}
	return nil
}

// slideFromMap 将 config 中的单张轮播图字段规整为统一结构（缺失字段置空）
func slideFromMap(m map[string]interface{}) gin.H {
	str := func(k string) string {
		if v, ok := m[k].(string); ok {
			return v
		}
		return ""
	}
	return gin.H{
		"image":       str("image"),
		"title":       str("title"),
		"subtitle":    str("subtitle"),
		"button_text": str("button_text"),
		"button_url":  str("button_url"),
	}
}

// extractHeroSettings 从首页 banner 模块提取轮播展示参数；无配置时返回默认值。
func extractHeroSettings(page *models.Page) gin.H {
	settings := services.DefaultHeroSettings()
	if page != nil {
		for _, m := range page.Modules {
			if m.Type != "banner" || !m.IsVisible {
				continue
			}
			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(m.Config), &raw); err != nil {
				continue
			}
			if sraw, ok := raw["settings"]; ok {
				b, _ := json.Marshal(sraw)
				var parsed services.HeroSettings
				if json.Unmarshal(b, &parsed) == nil {
					settings = services.NormalizeHeroSettings(&parsed)
				}
			}
			break
		}
	}
	boolOr := func(p *bool, def bool) bool {
		if p == nil {
			return def
		}
		return *p
	}
	return gin.H{
		"autoplay":       boolOr(settings.Autoplay, true),
		"interval_ms":    settings.IntervalMs,
		"transition":     settings.Transition,
		"show_dots":      boolOr(settings.ShowDots, true),
		"show_arrows":    boolOr(settings.ShowArrows, true),
		"pause_on_hover": boolOr(settings.PauseOnHover, true),
	}
}

// GetPage 获取页面（Redis 缓存）

func (h *PublicHandler) GetPage(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:page:" + c.Param("slug") + ":" + lang

	page, err := cached[*models.Page](h, key, services.CacheTTLShort, !h.previewMode(c), func() (*models.Page, error) {
		p, err := h.cmsService.GetPageBySlug(c.Param("slug"))
		if err != nil {
			return nil, err
		}
		services.LocalizePage(p, lang)
		return p, nil
	})
	if err != nil {
		utils.NotFound(c, "页面不存在")
		return
	}
	utils.Success(c, page)
}

// ListProducts 产品列表（增强版：支持多维度筛选 + Redis 缓存）
func (h *PublicHandler) ListProducts(c *gin.Context) {
	lang := middleware.GetLang(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	keyword := c.Query("keyword")
	search := c.Query("search")
	if search != "" && keyword == "" {
		keyword = search
	}

	// 缓存 key 需包含全部筛选维度，保证不同查询条件互不串用
	filters := strings.Join([]string{
		keyword,
		c.Query("category_id"),
		c.Query("gender"),
		c.Query("type"),
		c.Query("material"),
		c.Query("series_id"),
		c.Query("fabric_id"),
		c.Query("sort_by"),
		c.Query("sort_order"),
	}, "|")
	key := fmt.Sprintf("cache:products:%s:%d:%d:%s", lang, page, pageSize, filters)

	result, err := cached[cachedPage[models.Product]](h, key, services.CacheTTLMedium, !h.previewMode(c), func() (cachedPage[models.Product], error) {
		// 使用增强版查询，支持多维度筛选
		params := &services.ProductListParams{
			Page:       page,
			PageSize:   pageSize,
			Keyword:    keyword,
			CategoryID: c.Query("category_id"),
			Status:     "published",
			Gender:     c.Query("gender"),
			Type:       c.Query("type"),
			Material:   c.Query("material"),
			SeriesID:   c.Query("series_id"),
			FabricID:   c.Query("fabric_id"),
			SortBy:     c.Query("sort_by"),
			SortOrder:  c.Query("sort_order"),
		}

		products, total, err := h.productService.ListProductsV2(params)
		if err != nil {
			return cachedPage[models.Product]{}, err
		}
		services.LocalizeProducts(products, lang)
		return cachedPage[models.Product]{Items: products, Total: total}, nil
	})
	if err != nil {
		utils.InternalError(c, "获取产品列表失败")
		return
	}
	utils.SuccessPage(c, result.Items, page, pageSize, result.Total)
}

// GetProduct 获取产品详情（Redis 缓存）
func (h *PublicHandler) GetProduct(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:product:" + c.Param("slug") + ":" + lang

	product, err := cached[*models.Product](h, key, services.CacheTTLShort, !h.previewMode(c), func() (*models.Product, error) {
		p, err := h.productService.GetProductBySlug(c.Param("slug"))
		if err != nil {
			return nil, err
		}
		services.LocalizeProduct(p, lang)
		return p, nil
	})
	if err != nil {
		utils.NotFound(c, "产品不存在")
		return
	}
	utils.Success(c, product)
}

// ListCategories 分类列表（Redis 缓存）
func (h *PublicHandler) ListCategories(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:categories:" + lang

	categories, err := cached[[]models.Category](h, key, services.CacheTTLMedium, !h.previewMode(c), func() ([]models.Category, error) {
		return h.productService.ListCategories()
	})
	if err != nil {
		utils.InternalError(c, "获取分类列表失败")
		return
	}
	utils.Success(c, categories)
}

// ListBlogs 博客列表（Redis 缓存）
func (h *PublicHandler) ListBlogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	lang := middleware.GetLang(c)
	key := fmt.Sprintf("cache:blogs:%s:%d:%d:%s", lang, page, pageSize, keyword)

	result, err := cached[cachedPage[models.Blog]](h, key, services.CacheTTLMedium, !h.previewMode(c), func() (cachedPage[models.Blog], error) {
		blogs, total, err := h.cmsService.ListBlogs(page, pageSize, keyword, "published")
		if err != nil {
			return cachedPage[models.Blog]{}, err
		}
		services.LocalizeBlogs(blogs, lang)
		return cachedPage[models.Blog]{Items: blogs, Total: total}, nil
	})
	if err != nil {
		utils.InternalError(c, "获取博客列表失败")
		return
	}
	utils.SuccessPage(c, result.Items, page, pageSize, result.Total)
}

// GetBlog 获取博客详情（Redis 缓存）
func (h *PublicHandler) GetBlog(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:blog:" + c.Param("slug") + ":" + lang

	blog, err := cached[*models.Blog](h, key, services.CacheTTLShort, !h.previewMode(c), func() (*models.Blog, error) {
		b, err := h.cmsService.GetBlogBySlug(c.Param("slug"))
		if err != nil {
			return nil, err
		}
		services.LocalizeBlog(b, lang)
		return b, nil
	})
	if err != nil {
		utils.NotFound(c, "博客不存在")
		return
	}
	utils.Success(c, blog)
}

// ListCases 案例列表（Redis 缓存）
func (h *PublicHandler) ListCases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	lang := middleware.GetLang(c)
	key := fmt.Sprintf("cache:cases:%s:%d:%d", lang, page, pageSize)

	result, err := cached[cachedPage[models.Case]](h, key, services.CacheTTLMedium, !h.previewMode(c), func() (cachedPage[models.Case], error) {
		cases, total, err := h.cmsService.ListPublishedCases(page, pageSize, "")
		if err != nil {
			return cachedPage[models.Case]{}, err
		}
		services.LocalizeCases(cases, lang)
		return cachedPage[models.Case]{Items: cases, Total: total}, nil
	})
	if err != nil {
		utils.InternalError(c, "获取案例列表失败")
		return
	}
	utils.SuccessPage(c, result.Items, page, pageSize, result.Total)
}

// ListFAQs FAQ 列表（Redis 缓存）
func (h *PublicHandler) ListFAQs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	category := c.Query("category")
	language := c.DefaultQuery("language", "en")
	lang := middleware.GetLang(c)

	key := fmt.Sprintf("cache:faqs:%s:%d:%d:%s:%s", language, page, pageSize, category, lang)

	result, err := cached[cachedPage[models.FAQ]](h, key, services.CacheTTLMedium, !h.previewMode(c), func() (cachedPage[models.FAQ], error) {
		faqs, total, err := h.cmsService.ListFAQs(page, pageSize, category, language)
		if err != nil {
			return cachedPage[models.FAQ]{}, err
		}
		return cachedPage[models.FAQ]{Items: faqs, Total: total}, nil
	})
	if err != nil {
		utils.InternalError(c, "获取 FAQ 列表失败")
		return
	}
	utils.SuccessPage(c, result.Items, page, pageSize, result.Total)
}

// ListFabrics 面料列表（仅已发布 + Redis 缓存）
func (h *PublicHandler) ListFabrics(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:fabrics:" + lang

	fabrics, err := cached[[]models.Fabric](h, key, services.CacheTTLMedium, !h.previewMode(c), func() ([]models.Fabric, error) {
		return h.productService.ListPublishedFabrics()
	})
	if err != nil {
		utils.InternalError(c, "获取面料列表失败")
		return
	}
	utils.Success(c, fabrics)
}

// ListSeries 系列列表（仅已发布 + Redis 缓存）
func (h *PublicHandler) ListSeries(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:series:" + lang

	series, err := cached[[]models.Series](h, key, services.CacheTTLMedium, !h.previewMode(c), func() ([]models.Series, error) {
		return h.productService.ListPublishedSeries()
	})
	if err != nil {
		utils.InternalError(c, "获取系列列表失败")
		return
	}
	utils.Success(c, series)
}

// ListFactories 工厂列表（仅已发布 + Redis 缓存）
func (h *PublicHandler) ListFactories(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:factories:" + lang

	factories, err := cached[[]models.Factory](h, key, services.CacheTTLMedium, !h.previewMode(c), func() ([]models.Factory, error) {
		return h.cmsService.ListPublishedFactories()
	})
	if err != nil {
		utils.InternalError(c, "获取工厂列表失败")
		return
	}
	utils.Success(c, factories)
}

// ListCertifications 认证列表（仅已发布 + Redis 缓存）
func (h *PublicHandler) ListCertifications(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:certifications:" + lang

	certifications, err := cached[[]models.Certification](h, key, services.CacheTTLMedium, !h.previewMode(c), func() ([]models.Certification, error) {
		return h.cmsService.ListPublishedCertifications()
	})
	if err != nil {
		utils.InternalError(c, "获取认证列表失败")
		return
	}
	utils.Success(c, certifications)
}

// ListProductionProcesses 生产流程列表（仅已发布 + Redis 缓存）
func (h *PublicHandler) ListProductionProcesses(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:production-processes:" + lang

	processes, err := cached[[]models.ProductionProcess](h, key, services.CacheTTLMedium, !h.previewMode(c), func() ([]models.ProductionProcess, error) {
		return h.cmsService.ListPublishedProductionProcesses()
	})
	if err != nil {
		utils.InternalError(c, "获取生产流程列表失败")
		return
	}
	utils.Success(c, processes)
}

// ListNavigations 导航列表（Redis 缓存，长 TTL）
func (h *PublicHandler) ListNavigations(c *gin.Context) {
	navType := c.DefaultQuery("type", "header")
	key := "cache:navigations:" + navType

	navigations, err := cached[[]models.Navigation](h, key, services.CacheTTLLong, !h.previewMode(c), func() ([]models.Navigation, error) {
		return h.cmsService.ListNavigations(navType)
	})
	if err != nil {
		utils.InternalError(c, "获取导航列表失败")
		return
	}
	utils.Success(c, navigations)
}

// CreateLead 提交询盘
func (h *PublicHandler) CreateLead(c *gin.Context) {
	var req services.LeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	ip := middleware.GetClientIP(c)
	lead, err := h.leadService.CreateLead(&req, ip)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, lead)
}

// TrackClick 记录外部链接点击（社交媒体跳转跟踪）
// 门户 Footer 的社交图标点击时调用，将点击事件写入 visit_logs（entity_type=social_click）
func (h *PublicHandler) TrackClick(c *gin.Context) {
	var req struct {
		Target string `json:"target" binding:"required"` // youtube, instagram, xiaohongshu, facebook, twitter, linkedin
		URL    string `json:"url" binding:"required"`    // 目标链接
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 校验 target 白名单，防止任意数据写入
	validTargets := map[string]bool{
		"youtube": true, "instagram": true, "xiaohongshu": true,
		"facebook": true, "twitter": true, "linkedin": true,
	}
	if !validTargets[req.Target] {
		utils.BadRequest(c, "无效的点击目标")
		return
	}

	ip := middleware.GetClientIP(c)
	ua := middleware.GetUserAgent(c)
	device, browser, osName := middleware.ParseUserAgent(ua)

	log := models.VisitLog{
		CreatedAt:  time.Now(),
		IP:         ip,
		Device:     device,
		Browser:    browser,
		OS:         osName,
		UserAgent:  ua,
		Method:     "POST",
		Path:       "/api/v1/public/click-track",
		EntityType: "social_click",
		EntitySlug: req.Target,
		Status:     200,
		Referer:    c.GetHeader("Referer"),
		Source:     "social",
		Medium:     "click",
		SessionID:  c.GetString("request_id"),
	}

	// 异步写入 DB（不阻塞响应）
	middleware.RecordVisitLog(log)

	utils.Success(c, gin.H{"tracked": true})
}
