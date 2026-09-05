package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// SubscriptionHandler 订阅处理器（门户"订阅更新"Newsletter）
type SubscriptionHandler struct {
	service *services.SubscriptionService
}

// NewSubscriptionHandler 创建订阅处理器
func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// Subscribe 公开订阅（附 UTM/语言/落地页等归因信息，供广告转化分析）
func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var req services.SubscribeRequest
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

	sub, err := h.service.Subscribe(&req, ip, visitorID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, sub)
}

// ListSubscribers 订阅列表
func (h *SubscriptionHandler) ListSubscribers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	subs, total, err := h.service.List(page, pageSize, status, keyword)
	if err != nil {
		utils.InternalError(c, "获取订阅列表失败")
		return
	}
	utils.SuccessPage(c, subs, page, pageSize, total)
}

// UpdateSubscriber 更新订阅状态
func (h *SubscriptionHandler) UpdateSubscriber(c *gin.Context) {
	var req services.UpdateSubscriberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	sub, err := h.service.UpdateStatus(c.Param("id"), req.Status)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, sub)
}

// DeleteSubscriber 删除订阅
func (h *SubscriptionHandler) DeleteSubscriber(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除订阅失败")
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

// GetStats 订阅统计
func (h *SubscriptionHandler) GetStats(c *gin.Context) {
	stats, err := h.service.Stats()
	if err != nil {
		utils.InternalError(c, "获取统计数据失败")
		return
	}
	utils.Success(c, stats)
}
