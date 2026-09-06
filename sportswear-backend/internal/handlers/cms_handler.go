package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// CMSHandler CMS 处理器
type CMSHandler struct {
	cmsService *services.CMSService
	cache      *services.CacheService
}

// NewCMSHandler 创建 CMS 处理器
func NewCMSHandler(cmsService *services.CMSService, cacheService *services.CacheService) *CMSHandler {
	return &CMSHandler{cmsService: cmsService, cache: cacheService}
}

// invalidateCache 后台写操作后主动失效对应门户缓存（缓存不可用时无副作用）
func (h *CMSHandler) invalidateCache(entityType, slug string) {
	if h.cache != nil {
		h.cache.InvalidateContentCache(entityType, slug)
	}
}

// ==================== 页面 ====================

// ListPages 页面列表
func (h *CMSHandler) ListPages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	pages, total, err := h.cmsService.ListPages(page, pageSize, keyword, status)
	if err != nil {
		utils.InternalError(c, "获取页面列表失败")
		return
	}
	utils.SuccessPage(c, pages, page, pageSize, total)
}

// GetPage 获取页面
func (h *CMSHandler) GetPage(c *gin.Context) {
	page, err := h.cmsService.GetPage(c.Param("id"))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, page)
}

// CreatePage 创建页面
func (h *CMSHandler) CreatePage(c *gin.Context) {
	var req services.PageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	page, err := h.cmsService.CreatePage(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("page", "")
	utils.Created(c, page)
}

// UpdatePage 更新页面
func (h *CMSHandler) UpdatePage(c *gin.Context) {
	var req services.PageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	page, err := h.cmsService.UpdatePage(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("page", "")
	utils.Success(c, page)
}

// UpdateHeroSlides 更新页面轮播图配置（首页 banner 模块）
// 请求体：{ "slides": [...], "settings"?: { autoplay, interval_ms, transition, show_dots, show_arrows, pause_on_hover } }
func (h *CMSHandler) UpdateHeroSlides(c *gin.Context) {
	var req struct {
		Slides   []services.HeroSlide   `json:"slides" binding:"required"`
		Settings *services.HeroSettings `json:"settings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if len(req.Slides) == 0 {
		utils.BadRequest(c, "至少需要一张轮播图")
		return
	}
	if err := h.cmsService.UpdateHeroSlides(c.Param("id"), req.Slides, req.Settings); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("page", "")
	h.invalidateCache("theme", "")
	utils.Success(c, gin.H{"updated": true})
}

// DeletePage 删除页面
func (h *CMSHandler) DeletePage(c *gin.Context) {
	if err := h.cmsService.DeletePage(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除页面失败")
		return
	}
	h.invalidateCache("page", "")
	h.invalidateCache("navigation", "")
	// 自动同步：被删页面的关联导航自动隐藏
	if _, err := h.cmsService.SyncNavVisibilityWithPages(); err != nil {
		// 同步失败不阻断主流程，但记录以便排查
		h.invalidateCache("navigation", "")
	}
	utils.Success(c, gin.H{"deleted": true})
}

// PublishPage 发布页面
func (h *CMSHandler) PublishPage(c *gin.Context) {
	if err := h.cmsService.PublishPage(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布页面失败")
		return
	}
	h.invalidateCache("page", "")
	h.invalidateCache("navigation", "")
	// 自动同步：重新发布后关联导航自动恢复可见
	if _, err := h.cmsService.SyncNavVisibilityWithPages(); err != nil {
		h.invalidateCache("navigation", "")
	}
	utils.Success(c, gin.H{"published": true})
}

// UnpublishPage 下线页面
func (h *CMSHandler) UnpublishPage(c *gin.Context) {
	if err := h.cmsService.UnpublishPage(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线页面失败")
		return
	}
	h.invalidateCache("page", "")
	h.invalidateCache("navigation", "")
	// 自动同步：下线页面的关联导航自动隐藏
	if _, err := h.cmsService.SyncNavVisibilityWithPages(); err != nil {
		h.invalidateCache("navigation", "")
	}
	utils.Success(c, gin.H{"unpublished": true})
}

// ==================== 导航 ====================

// ListNavigations 导航列表
func (h *CMSHandler) ListNavigations(c *gin.Context) {
	navType := c.Query("type")
	navigations, err := h.cmsService.ListNavigations(navType)
	if err != nil {
		utils.InternalError(c, "获取导航列表失败")
		return
	}
	utils.Success(c, navigations)
}

// UpdateNavigation 更新导航
func (h *CMSHandler) UpdateNavigation(c *gin.Context) {
	var req services.NavigationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	nav, err := h.cmsService.UpdateNavigation(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("navigation", "")
	utils.Success(c, nav)
}

// BatchSortNavigations 批量更新导航排序
func (h *CMSHandler) BatchSortNavigations(c *gin.Context) {
	var req services.SortNavigationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if len(req.Items) == 0 {
		utils.BadRequest(c, "排序列表不能为空")
		return
	}
	if err := h.cmsService.BatchSortNavigations(&req); err != nil {
		utils.InternalError(c, "批量排序失败")
		return
	}
	h.invalidateCache("navigation", "")
	utils.Success(c, gin.H{"sorted": true})
}

// SyncNavVisibilityWithPages 同步导航可见性与页面状态
func (h *CMSHandler) SyncNavVisibilityWithPages(c *gin.Context) {
	updated, err := h.cmsService.SyncNavVisibilityWithPages()
	if err != nil {
		utils.InternalError(c, "同步导航状态失败")
		return
	}
	h.invalidateCache("navigation", "")
	utils.Success(c, gin.H{"updated": updated})
}

// ==================== 门户路由注册表 & 健康检查 ====================

// ListPortalRoutes 返回路由注册表（系统预置 + 已发布页面派生），供 SEO 管理下拉、导航派生与健康检查使用
func (h *CMSHandler) ListPortalRoutes(c *gin.Context) {
	routes, err := h.cmsService.ListPortalRoutes()
	if err != nil {
		utils.InternalError(c, "获取路由注册表失败")
		return
	}
	utils.Success(c, routes)
}

// PortalHealthCheck 门户闭环健康检查（死链 / 缺导航 / 缺 SEO / 路径漂移）
func (h *CMSHandler) PortalHealthCheck(c *gin.Context) {
	report, err := h.cmsService.PortalHealthCheck()
	if err != nil {
		utils.InternalError(c, "健康检查失败")
		return
	}
	utils.Success(c, report)
}

// ==================== 博客 ====================

// ListBlogs 博客列表
func (h *CMSHandler) ListBlogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	category := c.Query("category")

	blogs, total, err := h.cmsService.ListBlogs(page, pageSize, keyword, status, category)
	if err != nil {
		utils.InternalError(c, "获取博客列表失败")
		return
	}
	utils.SuccessPage(c, blogs, page, pageSize, total)
}

// GetBlog 获取博客
func (h *CMSHandler) GetBlog(c *gin.Context) {
	blog, err := h.cmsService.GetBlog(c.Param("id"))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, blog)
}

// CreateBlog 创建博客
func (h *CMSHandler) CreateBlog(c *gin.Context) {
	var req services.BlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	blog, err := h.cmsService.CreateBlog(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("blog", "")
	utils.Created(c, blog)
}

// UpdateBlog 更新博客
func (h *CMSHandler) UpdateBlog(c *gin.Context) {
	var req services.BlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	blog, err := h.cmsService.UpdateBlog(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("blog", "")
	utils.Success(c, blog)
}

// DeleteBlog 删除博客
func (h *CMSHandler) DeleteBlog(c *gin.Context) {
	if err := h.cmsService.DeleteBlog(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除博客失败")
		return
	}
	h.invalidateCache("blog", "")
	utils.Success(c, gin.H{"deleted": true})
}

// PublishBlog 发布博客
func (h *CMSHandler) PublishBlog(c *gin.Context) {
	if err := h.cmsService.PublishBlog(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布博客失败")
		return
	}
	h.invalidateCache("blog", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishBlog 下线博客
func (h *CMSHandler) UnpublishBlog(c *gin.Context) {
	if err := h.cmsService.UnpublishBlog(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线博客失败")
		return
	}
	h.invalidateCache("blog", "")
	utils.Success(c, gin.H{"unpublished": true})
}

// ==================== 案例 ====================

// ListCases 案例列表
func (h *CMSHandler) ListCases(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	projectType := c.Query("project_type")

	cases, total, err := h.cmsService.ListCases(page, pageSize, keyword, projectType)
	if err != nil {
		utils.InternalError(c, "获取案例列表失败")
		return
	}
	utils.SuccessPage(c, cases, page, pageSize, total)
}

// GetCase 获取案例
func (h *CMSHandler) GetCase(c *gin.Context) {
	caseItem, err := h.cmsService.GetCase(c.Param("id"))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, caseItem)
}

// CreateCase 创建案例
func (h *CMSHandler) CreateCase(c *gin.Context) {
	var req services.CaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	caseItem, err := h.cmsService.CreateCase(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("case", "")
	utils.Created(c, caseItem)
}

// UpdateCase 更新案例
func (h *CMSHandler) UpdateCase(c *gin.Context) {
	var req services.CaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	caseItem, err := h.cmsService.UpdateCase(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("case", "")
	utils.Success(c, caseItem)
}

// DeleteCase 删除案例
func (h *CMSHandler) DeleteCase(c *gin.Context) {
	if err := h.cmsService.DeleteCase(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除案例失败")
		return
	}
	h.invalidateCache("case", "")
	utils.Success(c, gin.H{"deleted": true})
}

// PublishCase 发布案例
func (h *CMSHandler) PublishCase(c *gin.Context) {
	if err := h.cmsService.PublishCase(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布案例失败")
		return
	}
	h.invalidateCache("case", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishCase 下线案例
func (h *CMSHandler) UnpublishCase(c *gin.Context) {
	if err := h.cmsService.UnpublishCase(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线案例失败")
		return
	}
	h.invalidateCache("case", "")
	utils.Success(c, gin.H{"unpublished": true})
}

// ==================== FAQ ====================

// ListFAQs FAQ 列表
func (h *CMSHandler) ListFAQs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	faqs, total, err := h.cmsService.ListFAQs(page, pageSize, category, keyword)
	if err != nil {
		utils.InternalError(c, "获取 FAQ 列表失败")
		return
	}
	utils.SuccessPage(c, faqs, page, pageSize, total)
}

// CreateFAQ 创建 FAQ
func (h *CMSHandler) CreateFAQ(c *gin.Context) {
	var req services.FAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	faq, err := h.cmsService.CreateFAQ(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("faq", "")
	utils.Created(c, faq)
}

// UpdateFAQ 更新 FAQ
func (h *CMSHandler) UpdateFAQ(c *gin.Context) {
	var req services.FAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	faq, err := h.cmsService.UpdateFAQ(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("faq", "")
	utils.Success(c, faq)
}

// DeleteFAQ 删除 FAQ
func (h *CMSHandler) DeleteFAQ(c *gin.Context) {
	if err := h.cmsService.DeleteFAQ(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除 FAQ 失败")
		return
	}
	h.invalidateCache("faq", "")
	utils.Success(c, gin.H{"deleted": true})
}

// ==================== 询盘问题模板 ====================

// ListInquiryTemplates 询盘问题模板列表
func (h *CMSHandler) ListInquiryTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	items, total, err := h.cmsService.ListInquiryTemplates(page, pageSize, category, keyword)
	if err != nil {
		utils.InternalError(c, "获取询盘模板列表失败")
		return
	}
	utils.SuccessPage(c, items, page, pageSize, total)
}

// CreateInquiryTemplate 创建询盘问题模板
func (h *CMSHandler) CreateInquiryTemplate(c *gin.Context) {
	var req services.InquiryTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	item, err := h.cmsService.CreateInquiryTemplate(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("inquiry-template", "")
	utils.Created(c, item)
}

// UpdateInquiryTemplate 更新询盘问题模板
func (h *CMSHandler) UpdateInquiryTemplate(c *gin.Context) {
	var req services.InquiryTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	item, err := h.cmsService.UpdateInquiryTemplate(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("inquiry-template", "")
	utils.Success(c, item)
}

// DeleteInquiryTemplate 删除询盘问题模板
func (h *CMSHandler) DeleteInquiryTemplate(c *gin.Context) {
	if err := h.cmsService.DeleteInquiryTemplate(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除询盘模板失败")
		return
	}
	h.invalidateCache("inquiry-template", "")
	utils.Success(c, gin.H{"deleted": true})
}

// ==================== 工厂 ====================

// ListFactories 工厂列表
func (h *CMSHandler) ListFactories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	factories, total, err := h.cmsService.ListFactories(page, pageSize, keyword)
	if err != nil {
		utils.InternalError(c, "获取工厂列表失败")
		return
	}
	utils.SuccessPage(c, factories, page, pageSize, total)
}

// CreateFactory 创建工厂
func (h *CMSHandler) CreateFactory(c *gin.Context) {
	var req services.FactoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	factory, err := h.cmsService.CreateFactory(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("factory", "")
	utils.Created(c, factory)
}

// UpdateFactory 更新工厂
func (h *CMSHandler) UpdateFactory(c *gin.Context) {
	var req services.FactoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	factory, err := h.cmsService.UpdateFactory(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("factory", "")
	utils.Success(c, factory)
}

// DeleteFactory 删除工厂
func (h *CMSHandler) DeleteFactory(c *gin.Context) {
	if err := h.cmsService.DeleteFactory(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除工厂失败")
		return
	}
	h.invalidateCache("factory", "")
	utils.Success(c, gin.H{"deleted": true})
}

// PublishFactory 发布工厂（前端可见）
func (h *CMSHandler) PublishFactory(c *gin.Context) {
	if err := h.cmsService.PublishFactory(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布工厂失败")
		return
	}
	h.invalidateCache("factory", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishFactory 下线工厂（前端不可见）
func (h *CMSHandler) UnpublishFactory(c *gin.Context) {
	if err := h.cmsService.UnpublishFactory(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线工厂失败")
		return
	}
	h.invalidateCache("factory", "")
	utils.Success(c, gin.H{"unpublished": true})
}

// ==================== 认证 ====================

// ListCertifications 认证列表
func (h *CMSHandler) ListCertifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	certifications, total, err := h.cmsService.ListCertifications(page, pageSize, keyword)
	if err != nil {
		utils.InternalError(c, "获取认证列表失败")
		return
	}
	utils.SuccessPage(c, certifications, page, pageSize, total)
}

// CreateCertification 创建认证
func (h *CMSHandler) CreateCertification(c *gin.Context) {
	var req services.CertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	cert, err := h.cmsService.CreateCertification(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("certification", "")
	utils.Created(c, cert)
}

// UpdateCertification 更新认证
func (h *CMSHandler) UpdateCertification(c *gin.Context) {
	var req services.CertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	cert, err := h.cmsService.UpdateCertification(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("certification", "")
	utils.Success(c, cert)
}

// DeleteCertification 删除认证
func (h *CMSHandler) DeleteCertification(c *gin.Context) {
	if err := h.cmsService.DeleteCertification(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除认证失败")
		return
	}
	h.invalidateCache("certification", "")
	utils.Success(c, gin.H{"deleted": true})
}

// PublishCertification 发布认证（前端可见）
func (h *CMSHandler) PublishCertification(c *gin.Context) {
	if err := h.cmsService.PublishCertification(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布认证失败")
		return
	}
	h.invalidateCache("certification", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishCertification 下线认证（前端不可见）
func (h *CMSHandler) UnpublishCertification(c *gin.Context) {
	if err := h.cmsService.UnpublishCertification(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线认证失败")
		return
	}
	h.invalidateCache("certification", "")
	utils.Success(c, gin.H{"unpublished": true})
}

// ==================== 生产流程 ====================

// ListProductionProcesses 生产流程列表
func (h *CMSHandler) ListProductionProcesses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	processes, total, err := h.cmsService.ListProductionProcesses(page, pageSize, keyword)
	if err != nil {
		utils.InternalError(c, "获取生产流程列表失败")
		return
	}
	utils.SuccessPage(c, processes, page, pageSize, total)
}

// CreateProductionProcess 创建生产流程
func (h *CMSHandler) CreateProductionProcess(c *gin.Context) {
	var req services.ProductionProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	p, err := h.cmsService.CreateProductionProcess(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("production_process", "")
	utils.Created(c, p)
}

// UpdateProductionProcess 更新生产流程
func (h *CMSHandler) UpdateProductionProcess(c *gin.Context) {
	var req services.ProductionProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	p, err := h.cmsService.UpdateProductionProcess(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("production_process", "")
	utils.Success(c, p)
}

// DeleteProductionProcess 删除生产流程
func (h *CMSHandler) DeleteProductionProcess(c *gin.Context) {
	if err := h.cmsService.DeleteProductionProcess(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除生产流程失败")
		return
	}
	h.invalidateCache("production_process", "")
	utils.Success(c, gin.H{"deleted": true})
}

// PublishProductionProcess 发布生产流程
func (h *CMSHandler) PublishProductionProcess(c *gin.Context) {
	if err := h.cmsService.PublishProductionProcess(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布生产流程失败")
		return
	}
	h.invalidateCache("production_process", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishProductionProcess 下线生产流程
func (h *CMSHandler) UnpublishProductionProcess(c *gin.Context) {
	if err := h.cmsService.UnpublishProductionProcess(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线生产流程失败")
		return
	}
	h.invalidateCache("production_process", "")
	utils.Success(c, gin.H{"unpublished": true})
}

// ==================== 自媒体 ====================

// ListSelfMedias 自媒体列表
func (h *CMSHandler) ListSelfMedias(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	items, total, err := h.cmsService.ListSelfMedias(page, pageSize, keyword)
	if err != nil {
		utils.InternalError(c, "获取自媒体列表失败")
		return
	}
	utils.SuccessPage(c, items, page, pageSize, total)
}

// CreateSelfMedia 创建自媒体
func (h *CMSHandler) CreateSelfMedia(c *gin.Context) {
	var req services.SelfMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	item, err := h.cmsService.CreateSelfMedia(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("self_media", "")
	utils.Created(c, item)
}

// UpdateSelfMedia 更新自媒体
func (h *CMSHandler) UpdateSelfMedia(c *gin.Context) {
	var req services.SelfMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	item, err := h.cmsService.UpdateSelfMedia(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("self_media", "")
	utils.Success(c, item)
}

// DeleteSelfMedia 删除自媒体
func (h *CMSHandler) DeleteSelfMedia(c *gin.Context) {
	if err := h.cmsService.DeleteSelfMedia(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除自媒体失败")
		return
	}
	h.invalidateCache("self_media", "")
	utils.Success(c, gin.H{"deleted": true})
}

// PublishSelfMedia 发布自媒体（前端可见）
func (h *CMSHandler) PublishSelfMedia(c *gin.Context) {
	if err := h.cmsService.PublishSelfMedia(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布自媒体失败")
		return
	}
	h.invalidateCache("self_media", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishSelfMedia 下线自媒体（前端不可见）
func (h *CMSHandler) UnpublishSelfMedia(c *gin.Context) {
	if err := h.cmsService.UnpublishSelfMedia(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线自媒体失败")
		return
	}
	h.invalidateCache("self_media", "")
	utils.Success(c, gin.H{"unpublished": true})
}

// GetSelfMediaConfig 读取自媒体展示配置（最多展示数量）
func (h *CMSHandler) GetSelfMediaConfig(c *gin.Context) {
	utils.Success(c, gin.H{"max_display": h.cmsService.SelfMediaMaxDisplay()})
}

// UpdateSelfMediaConfig 更新自媒体展示配置（最多展示数量）
func (h *CMSHandler) UpdateSelfMediaConfig(c *gin.Context) {
	var req struct {
		MaxDisplay int `json:"max_display"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.cmsService.SetSelfMediaMaxDisplay(req.MaxDisplay); err != nil {
		utils.BadRequest(c, "保存自媒体配置失败")
		return
	}
	h.invalidateCache("self_media", "")
	utils.Success(c, gin.H{"max_display": h.cmsService.SelfMediaMaxDisplay()})
}
