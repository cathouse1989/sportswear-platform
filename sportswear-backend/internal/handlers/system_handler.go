package handlers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// SystemHandler 系统处理器（通知 + 操作日志 + 报价 + sitemap）
type SystemHandler struct {
	notificationService *services.NotificationService
	operationLogService *services.OperationLogService
	quoteService        *services.QuoteService
	cmsService          *services.CMSService
	productService      *services.ProductService
}

// NewSystemHandler 创建系统处理器
func NewSystemHandler(
	notificationService *services.NotificationService,
	operationLogService *services.OperationLogService,
	quoteService *services.QuoteService,
	cmsService *services.CMSService,
	productService *services.ProductService,
) *SystemHandler {
	return &SystemHandler{
		notificationService: notificationService,
		operationLogService: operationLogService,
		quoteService:        quoteService,
		cmsService:          cmsService,
		productService:      productService,
	}
}

// ==================== 通知 ====================

// ListNotifications 通知列表
func (h *SystemHandler) ListNotifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	onlyUnread := c.Query("only_unread") == "true"

	userID := middleware.GetUserID(c)
	notifications, total, err := h.notificationService.List(userID, page, pageSize, onlyUnread)
	if err != nil {
		utils.InternalError(c, "获取通知列表失败")
		return
	}
	utils.SuccessPage(c, notifications, page, pageSize, total)
}

// UnreadCount 未读数量
func (h *SystemHandler) UnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	count, err := h.notificationService.UnreadCount(userID)
	if err != nil {
		utils.InternalError(c, "获取未读数量失败")
		return
	}
	utils.Success(c, gin.H{"count": count})
}

// MarkRead 标记已读
func (h *SystemHandler) MarkRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.notificationService.MarkRead(c.Param("id"), userID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"read": true})
}

// MarkAllRead 全部标记已读
func (h *SystemHandler) MarkAllRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.notificationService.MarkAllRead(userID); err != nil {
		utils.InternalError(c, "标记已读失败")
		return
	}
	utils.Success(c, gin.H{"read": true})
}

// ==================== 操作日志 ====================

// ListOperationLogs 操作日志列表
func (h *SystemHandler) ListOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	userID := c.Query("user_id")
	module := c.Query("module")
	operation := c.Query("operation")

	logs, total, err := h.operationLogService.List(page, pageSize, userID, module, operation)
	if err != nil {
		utils.InternalError(c, "获取操作日志失败")
		return
	}
	utils.SuccessPage(c, logs, page, pageSize, total)
}

// ==================== 报价 ====================

// ListQuotes 报价列表
func (h *SystemHandler) ListQuotes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	leadID := c.Query("lead_id")
	status := c.Query("status")

	quotes, total, err := h.quoteService.ListQuotes(page, pageSize, leadID, status)
	if err != nil {
		utils.InternalError(c, "获取报价列表失败")
		return
	}
	utils.SuccessPage(c, quotes, page, pageSize, total)
}

// GetQuote 获取报价
func (h *SystemHandler) GetQuote(c *gin.Context) {
	quote, err := h.quoteService.GetQuote(c.Param("id"))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, quote)
}

// CreateQuote 创建报价
func (h *SystemHandler) CreateQuote(c *gin.Context) {
	var req services.QuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	quote, err := h.quoteService.CreateQuote(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, quote)
}

// UpdateQuote 更新报价
func (h *SystemHandler) UpdateQuote(c *gin.Context) {
	var req services.QuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	quote, err := h.quoteService.UpdateQuote(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, quote)
}

// DeleteQuote 删除报价
func (h *SystemHandler) DeleteQuote(c *gin.Context) {
	if err := h.quoteService.DeleteQuote(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除报价失败")
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

// ==================== Sitemap ====================

// GetSitemap 生成 sitemap.xml
func (h *SystemHandler) GetSitemap(c *gin.Context) {
	baseURL := c.DefaultQuery("base_url", "https://example.com")
	lang := c.DefaultQuery("lang", "en")

	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sb.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

	// 首页
	sb.WriteString(services.BuildSitemapEntry(baseURL, lang, "", "") + "\n")

	// 页面（已发布）
	pages, _, _ := h.cmsService.ListPages(1, 1000, "", "published")
	for _, p := range pages {
		sb.WriteString(services.BuildSitemapEntry(baseURL, lang, "pages/"+p.Slug, p.UpdatedAt.Format("2006-01-02")) + "\n")
	}

	// 产品（已发布）
	products, _, _ := h.productService.ListProducts(1, 1000, "", "", "published")
	for _, p := range products {
		sb.WriteString(services.BuildSitemapEntry(baseURL, lang, "products/"+p.Slug, p.UpdatedAt.Format("2006-01-02")) + "\n")
	}

	// 博客（已发布）
	blogs, _, _ := h.cmsService.ListBlogs(1, 1000, "", "published")
	for _, b := range blogs {
		sb.WriteString(services.BuildSitemapEntry(baseURL, lang, "blogs/"+b.Slug, b.UpdatedAt.Format("2006-01-02")) + "\n")
	}

	sb.WriteString("</urlset>")

	c.Header("Content-Type", "application/xml")
	c.String(200, sb.String())
}
