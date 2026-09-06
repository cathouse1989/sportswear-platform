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
	uploadService  *services.UploadService
}

// NewPublicHandler 创建公开 API 处理器
func NewPublicHandler(
	productService *services.ProductService,
	cmsService *services.CMSService,
	leadService *services.LeadService,
	cacheService *services.CacheService,
	uploadService *services.UploadService,
) *PublicHandler {
	return &PublicHandler{
		productService: productService,
		cmsService:     cmsService,
		leadService:    leadService,
		cache:          cacheService,
		uploadService:  uploadService,
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
			factories      []models.Factory
		)

		wg.Add(6)
		go func() { defer wg.Done(); products, _ = h.productService.ListFeaturedProducts(8, lang) }()
		go func() { defer wg.Done(); categories, _ = h.productService.ListActiveCategories() }()
		go func() { defer wg.Done(); blogs, _, _ = h.cmsService.ListBlogs(1, 4, "", "published", "") }()
		go func() { defer wg.Done(); cases, _, _ = h.cmsService.ListPublishedCases(1, 4, "", "") }()
		go func() { defer wg.Done(); certifications, _ = h.cmsService.ListPublishedCertifications() }()
		go func() { defer wg.Done(); factories, _ = h.cmsService.ListPublishedFactories() }()
		wg.Wait()

		// 本地化处理
		services.LocalizeProducts(products, lang)
		services.LocalizeBlogs(blogs, lang)
		services.LocalizeCases(cases, lang)
		services.LocalizeCertifications(certifications, lang)
		services.LocalizeFactories(factories, lang)

		return gin.H{
			"page":              page,
			"hero_slides":       extractHeroSlides(page),
			"hero_settings":     extractHeroSettings(page),
			"featured_products": products,
			"categories":        categories,
			"blogs":             blogs,
			"cases":             cases,
			"certifications":    certifications,
			"factories":         factories,
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
//
// 返回统一结构供门户 / 管理后台门户预览渲染轮播图；无有效配置时返回 nil。
func extractHeroSlides(page *models.Page) []gin.H {
	if page == nil {
		return nil
	}
	for _, m := range page.Modules {
		if m.Type != "banner" || !m.IsVisible {
			continue
		}
		raw := services.UnmarshalPageModuleConfig(m.Config)
		if raw == nil {
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
			raw := services.UnmarshalPageModuleConfig(m.Config)
			if raw == nil {
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
		"show_arrows":    true,
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
	// 供流量中间件记录实体信息
	c.Set("visit_entity_id", page.ID.String())
	c.Set("visit_entity_name", page.Title)
	utils.Success(c, page)
}

// ListLandingPages 列出所有已发布落地页（供 sitemap 生成使用，Redis 缓存）。
// 固定列表页（home/products/blog/cases/faq/contact）由门户静态 sitemap 覆盖。
func (h *PublicHandler) ListLandingPages(c *gin.Context) {
	pages, err := cached[[]models.Page](h, "cache:landing-pages", services.CacheTTLMedium, !h.previewMode(c), func() ([]models.Page, error) {
		return h.cmsService.ListPublishedLandingPages()
	})
	if err != nil {
		utils.InternalError(c, "获取页面列表失败")
		return
	}
	utils.Success(c, pages)
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
			// 按访问语言动态排序：不同语种市场有各自的特色运动与运动服装
			Language: lang,
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
	// 供流量中间件记录实体信息（产品维度）
	c.Set("visit_entity_id", product.ID.String())
	c.Set("visit_entity_name", product.Name)
	utils.Success(c, product)
}

// ListCategories 分类列表（Redis 缓存）
func (h *PublicHandler) ListCategories(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:categories:" + lang

	categories, err := cached[[]models.Category](h, key, services.CacheTTLMedium, !h.previewMode(c), func() ([]models.Category, error) {
		return h.productService.ListActiveCategories()
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
	category := c.Query("category")

	lang := middleware.GetLang(c)
	key := fmt.Sprintf("cache:blogs:%s:%d:%d:%s:%s", lang, page, pageSize, keyword, category)

	result, err := cached[cachedPage[models.Blog]](h, key, services.CacheTTLMedium, !h.previewMode(c), func() (cachedPage[models.Blog], error) {
		blogs, total, err := h.cmsService.ListBlogs(page, pageSize, keyword, "published", category)
		if err != nil {
			return cachedPage[models.Blog]{}, err
		}
		services.LocalizeBlogs(blogs, lang)
		// 列表瘦身：不下发正文全文与翻译（详情接口再取），仅保留纯文本摘要，控制列表 payload 与带宽
		for i := range blogs {
			blogs[i].Summary = services.ExcerptHTML(blogs[i].Content, 200)
			blogs[i].Content = ""
			blogs[i].Translations = nil
		}
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
	// 供流量中间件记录实体信息
	c.Set("visit_entity_id", blog.ID.String())
	c.Set("visit_entity_name", blog.Title)
	utils.Success(c, blog)
}

// ListCases 案例列表（Redis 缓存）
func (h *PublicHandler) ListCases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	projectType := c.Query("project_type")

	lang := middleware.GetLang(c)
	key := fmt.Sprintf("cache:cases:%s:%d:%d:%s", lang, page, pageSize, projectType)

	result, err := cached[cachedPage[models.Case]](h, key, services.CacheTTLMedium, !h.previewMode(c), func() (cachedPage[models.Case], error) {
		cases, total, err := h.cmsService.ListPublishedCases(page, pageSize, "", projectType)
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

// GetCase 获取案例详情（Redis 缓存）
func (h *PublicHandler) GetCase(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:case:" + c.Param("slug") + ":" + lang

	caseItem, err := cached[*models.Case](h, key, services.CacheTTLShort, !h.previewMode(c), func() (*models.Case, error) {
		item, err := h.cmsService.GetCaseBySlug(c.Param("slug"))
		if err != nil {
			return nil, err
		}
		services.LocalizeCase(item, lang)
		return item, nil
	})
	if err != nil {
		utils.NotFound(c, "案例不存在")
		return
	}
	// 供流量中间件记录实体信息
	c.Set("visit_entity_id", caseItem.ID.String())
	c.Set("visit_entity_name", caseItem.Title)
	utils.Success(c, caseItem)
}

// ListFAQs FAQ 列表（Redis 缓存）
func (h *PublicHandler) ListFAQs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	category := c.Query("category")
	lang := middleware.GetLang(c)

	key := fmt.Sprintf("cache:faqs:%s:%d:%d:%s", lang, page, pageSize, category)

	result, err := cached[cachedPage[models.FAQ]](h, key, services.CacheTTLMedium, !h.previewMode(c), func() (cachedPage[models.FAQ], error) {
		faqs, total, err := h.cmsService.ListPublishedFAQs(page, pageSize, category)
		if err != nil {
			return cachedPage[models.FAQ]{}, err
		}
		services.LocalizeFAQs(faqs, lang)
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
		fabrics, err := h.productService.ListPublishedFabrics()
		if err != nil {
			return nil, err
		}
		services.LocalizeFabrics(fabrics, lang)
		return fabrics, nil
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
	services.LocalizeFactories(factories, lang)
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
	services.LocalizeCertifications(certifications, lang)
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
	services.LocalizeProductionProcesses(processes, lang)
	utils.Success(c, processes)
}

// ListSelfMedias 自媒体列表（仅已发布 + Redis 缓存 + 最多展示数量控制）
func (h *PublicHandler) ListSelfMedias(c *gin.Context) {
	lang := middleware.GetLang(c)
	key := "cache:self-medias:" + lang

	// 兼容控制：显式 ?limit= 优先，否则读取后台「自媒体最多展示数量」配置
	limit := 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if limit <= 0 {
		limit = h.cmsService.SelfMediaMaxDisplay()
	}

	items, err := cached[[]models.SelfMedia](h, key, services.CacheTTLMedium, !h.previewMode(c), func() ([]models.SelfMedia, error) {
		return h.cmsService.ListPublishedSelfMedias(limit)
	})
	if err != nil {
		utils.InternalError(c, "获取自媒体列表失败")
		return
	}
	utils.Success(c, items)
}

// ListNavigations 导航列表（Redis 缓存，长 TTL）
func (h *PublicHandler) ListNavigations(c *gin.Context) {
	navType := c.DefaultQuery("type", "header")
	lang := middleware.GetLang(c)
	key := "cache:navigations:" + navType

	navigations, err := cached[[]models.Navigation](h, key, services.CacheTTLLong, !h.previewMode(c), func() ([]models.Navigation, error) {
		return h.cmsService.ListNavigations(navType)
	})
	if err != nil {
		utils.InternalError(c, "获取导航列表失败")
		return
	}
	services.LocalizeNavigations(navigations, lang)
	utils.Success(c, navigations)
}

// CreateLead 提交询盘（附 UTM/国家/语言/落地页等归因信息，供广告转化分析）
func (h *PublicHandler) CreateLead(c *gin.Context) {
	var req services.LeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 客户端未显式提供时，从 URL/请求头补齐归因字段（门户 SSR 透传 X-UTM-*）
	if req.Source == "" {
		req.Source = middleware.FirstNonEmptyHeader(c, c.Query("utm_source"), "X-UTM-Source", "UTM-Source")
	}
	if req.Medium == "" {
		req.Medium = middleware.FirstNonEmptyHeader(c, c.Query("utm_medium"), "X-UTM-Medium", "UTM-Medium")
	}
	if req.Campaign == "" {
		req.Campaign = middleware.FirstNonEmptyHeader(c, c.Query("utm_campaign"), "X-UTM-Campaign", "UTM-Campaign")
	}
	if req.Keyword == "" {
		req.Keyword = middleware.FirstNonEmptyHeader(c, c.Query("utm_term"), "X-UTM-Term", "UTM-Term")
	}
	if req.Country == "" {
		req.Country = middleware.ResolveCountry(c)
	}
	if req.Language == "" {
		req.Language = middleware.GetLang(c)
	}
	if req.Device == "" {
		ua := middleware.GetUserAgent(c)
		device, _, _ := middleware.ParseUserAgent(ua)
		req.Device = device
	}
	if req.LandingPage == "" {
		req.LandingPage = c.Request.URL.Path
	}

	ip := middleware.GetClientIP(c)
	visitorID := c.GetHeader("X-Visitor-ID")
	if visitorID == "" {
		visitorID = "unknown"
	}
	lead, err := h.leadService.CreateLead(&req, ip, visitorID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, lead)

	// 异步回写 visit_logs 的 lead_id（关联该访客的所有访问记录）
	go h.leadService.BackfillLeadID(visitorID, lead.ID)
}

// UploadLeadAttachment 门户询盘附件上传（multipart/form-data，字段 file）
// 公开未登录接口：附件专用白名单（图片/PDF/Word/Excel/zip/rar，不含 SVG 防 XSS）
// + 大小限制（MAX_UPLOAD_SIZE，默认 20MB）+ RateLimit 中间件保护
// 仅保存文件返回访问 URL；媒体记录随询盘 attachments 字段落库（避免污染媒体库）
func (h *PublicHandler) UploadLeadAttachment(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 校验附件安全（询盘附件专用白名单 + 大小限制）
	if err := h.uploadService.ValidateLeadAttachment(fileHeader); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 生成存储路径（leads/2026-01/uuid.ext）
	storagePath, _ := h.uploadService.GenerateStoragePath(fileHeader.Filename, "leads")

	// 打开文件流并保存到本地
	src, err := fileHeader.Open()
	if err != nil {
		utils.InternalError(c, "读取文件失败")
		return
	}
	defer src.Close()

	if err := h.uploadService.SaveToLocal(src, storagePath); err != nil {
		utils.InternalError(c, "保存文件失败")
		return
	}

	utils.Created(c, gin.H{
		"url":       h.uploadService.LocalURL(storagePath),
		"file_name": fileHeader.Filename,
		"file_size": fileHeader.Size,
		"file_type": services.DetectFileType(fileHeader.Filename),
	})
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
		CreatedAt:   time.Now(),
		IP:          ip,
		Country:     middleware.ResolveCountry(c),
		Language:    middleware.GetLang(c),
		Device:      device,
		DeviceModel: middleware.ParseDeviceModel(ua),
		Browser:     browser,
		OS:          osName,
		UserAgent:   ua,
		Method:      "POST",
		Path:        "/api/v1/public/click-track",
		EntityType:  "social_click",
		EntitySlug:  req.Target,
		Status:      200,
		Referer:     c.GetHeader("Referer"),
		Source:      "social",
		Medium:      "click",
		SessionID:   c.GetString("request_id"),
	}

	// 异步写入 DB（不阻塞响应）
	middleware.RecordVisitLog(log)

	utils.Success(c, gin.H{"tracked": true})
}

// RequestDataExport 数据主体请求导出个人数据（GDPR Art.15 / PIPL 第44-45条）
// 用户通过提供邮箱确认身份，系统返回对应邮箱关联的所有询盘数据
func (h *PublicHandler) RequestDataExport(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请提供有效的邮箱地址")
		return
	}

	data, err := h.leadService.ExportByEmail(req.Email)
	if err != nil {
		if err.Error() == "未找到对应的个人数据" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"exported_at": time.Now(),
		"data":        data,
		"note":        "根据适用的数据保护法规，您有权获取我们持有的您的个人数据副本。",
	})
}

// RequestDataDeletion 数据主体请求删除个人数据（GDPR Art.17 "被遗忘权" / PIPL 第47条）
// 用户通过提供邮箱确认身份，系统将对应邮箱关联的所有询盘数据匿名化或删除
func (h *PublicHandler) RequestDataDeletion(c *gin.Context) {
	var req struct {
		Email  string `json:"email" binding:"required,email"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请提供有效的邮箱地址")
		return
	}

	err := h.leadService.AnonymizeByEmail(req.Email, req.Reason)
	if err != nil {
		if err.Error() == "未找到对应的个人数据" {
			utils.NotFound(c, "未找到与 "+req.Email+" 关联的数据，可能已被删除")
			return
		}
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"message": "您的个人数据删除请求已处理。如有疑问，请联系我们。",
		"email":   req.Email,
	})
}
