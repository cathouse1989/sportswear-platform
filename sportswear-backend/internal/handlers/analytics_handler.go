package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// AnalyticsHandler 流量分析处理器
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

// NewAnalyticsHandler 创建流量分析处理器
func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetTrafficOverview 流量总览（访问量、独立IP、转化率、每日趋势）
func (h *AnalyticsHandler) GetTrafficOverview(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	data, err := h.analyticsService.GetTrafficOverview(days)
	if err != nil {
		utils.InternalError(c, "获取流量总览失败")
		return
	}
	utils.Success(c, data)
}

// GetTopPages 热门页面
func (h *AnalyticsHandler) GetTopPages(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := h.analyticsService.GetTopPages(days, limit)
	if err != nil {
		utils.InternalError(c, "获取热门页面失败")
		return
	}
	utils.Success(c, data)
}

// GetTopProducts 热门产品（商业价值：最受关注的产品）
func (h *AnalyticsHandler) GetTopProducts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := h.analyticsService.GetTopProducts(days, limit)
	if err != nil {
		utils.InternalError(c, "获取热门产品失败")
		return
	}
	utils.Success(c, data)
}

// GetUtmCampaignAnalysis 广告归因分析（utm_campaign 维度）
func (h *AnalyticsHandler) GetUtmCampaignAnalysis(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	data, err := h.analyticsService.GetUtmCampaignAnalysis(days, limit)
	if err != nil {
		utils.InternalError(c, "获取广告归因分析失败")
		return
	}
	utils.Success(c, data)
}

// GetSourceAnalysis 来源分析（SEO/广告归因）
func (h *AnalyticsHandler) GetSourceAnalysis(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := h.analyticsService.GetSourceAnalysis(days)
	if err != nil {
		utils.InternalError(c, "获取来源分析失败")
		return
	}
	utils.Success(c, data)
}

// GetCountryAnalysis 国家/地区分布
func (h *AnalyticsHandler) GetCountryAnalysis(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))

	data, err := h.analyticsService.GetCountryAnalysis(days, limit)
	if err != nil {
		utils.InternalError(c, "获取地区分布失败")
		return
	}
	utils.Success(c, data)
}

// GetDeviceAnalysis 设备/浏览器分析
func (h *AnalyticsHandler) GetDeviceAnalysis(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := h.analyticsService.GetDeviceAnalysis(days)
	if err != nil {
		utils.InternalError(c, "获取设备分析失败")
		return
	}
	utils.Success(c, data)
}

// GetSocialClickAnalysis 社交媒体点击分析（外部链接跳转跟踪）
func (h *AnalyticsHandler) GetSocialClickAnalysis(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := h.analyticsService.GetSocialClickAnalysis(days)
	if err != nil {
		utils.InternalError(c, "获取社交媒体点击分析失败")
		return
	}
	utils.Success(c, data)
}

// ListVisitLogs 门户访问日志明细（分页 + IP/国家/来源/实体/语言/设备/时间范围筛选）
func (h *AnalyticsHandler) ListVisitLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	start, end := ParseDateRange(c.Query("start_date"), c.Query("end_date"))

	f := services.VisitLogFilter{
		IP:         c.Query("ip"),
		Country:    c.Query("country"),
		Source:     c.Query("source"),
		EntityType: c.Query("entity_type"),
		EntitySlug: c.Query("entity_slug"),
		Language:   c.Query("language"),
		Device:     c.Query("device"),
		Keyword:    c.Query("keyword"),
		Start:      start,
		End:        end,
	}

	rows, total, err := h.analyticsService.ListVisitLogs(page, pageSize, f)
	if err != nil {
		utils.InternalError(c, "获取访问日志明细失败")
		return
	}
	utils.SuccessPage(c, rows, page, pageSize, total)
}

// GetVisitorJourney 获取访客旅程（某个访客的完整访问路径）
func (h *AnalyticsHandler) GetVisitorJourney(c *gin.Context) {
	visitorID := c.Query("visitor_id")
	if visitorID == "" {
		utils.BadRequest(c, "缺少 visitor_id 参数")
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := h.analyticsService.GetVisitorJourney(visitorID, days)
	if err != nil {
		utils.InternalError(c, "获取访客旅程失败")
		return
	}
	utils.Success(c, data)
}

// GetIPLeads 获取 IP 关联的询盘列表
func (h *AnalyticsHandler) GetIPLeads(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		utils.BadRequest(c, "缺少 ip 参数")
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := h.analyticsService.GetIPLeads(ip, days)
	if err != nil {
		utils.InternalError(c, "获取 IP 询盘关联失败")
		return
	}
	utils.Success(c, data)
}

// GetConversionFunnel 获取转化漏斗（访问 → 产品浏览 → 询盘）
func (h *AnalyticsHandler) GetConversionFunnel(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := h.analyticsService.GetConversionFunnel(days)
	if err != nil {
		utils.InternalError(c, "获取转化漏斗失败")
		return
	}
	utils.Success(c, data)
}

// GetConsentInsights 获取隐私合规洞察（同意/未同意用户的转化对比）
func (h *AnalyticsHandler) GetConsentInsights(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := h.analyticsService.GetConsentInsights(days)
	if err != nil {
		utils.InternalError(c, "获取隐私合规洞察失败")
		return
	}
	utils.Success(c, data)
}
