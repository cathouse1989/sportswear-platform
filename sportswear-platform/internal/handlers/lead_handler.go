package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-platform/internal/middleware"
	"sportswear-platform/internal/services"
	"sportswear-platform/internal/utils"
)

// LeadHandler 询盘处理器
type LeadHandler struct {
	leadService *services.LeadService
}

// NewLeadHandler 创建询盘处理器
func NewLeadHandler(leadService *services.LeadService) *LeadHandler {
	return &LeadHandler{leadService: leadService}
}

// ListLeads 询盘列表
func (h *LeadHandler) ListLeads(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	leads, total, err := h.leadService.ListLeads(page, pageSize, status, keyword)
	if err != nil {
		utils.InternalError(c, "获取询盘列表失败")
		return
	}
	utils.SuccessPage(c, leads, page, pageSize, total)
}

// GetLead 获取询盘
func (h *LeadHandler) GetLead(c *gin.Context) {
	lead, err := h.leadService.GetLead(c.Param("id"))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, lead)
}

// UpdateLead 更新询盘
func (h *LeadHandler) UpdateLead(c *gin.Context) {
	var req services.UpdateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	lead, err := h.leadService.UpdateLead(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, lead)
}

// DeleteLead 删除询盘
func (h *LeadHandler) DeleteLead(c *gin.Context) {
	if err := h.leadService.DeleteLead(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除询盘失败")
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

// AddFollowUp 添加跟进记录
func (h *LeadHandler) AddFollowUp(c *gin.Context) {
	var req services.FollowUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	followUp, err := h.leadService.AddFollowUp(c.Param("id"), userID.String(), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, followUp)
}

// GetDashboardStats 获取仪表盘统计
func (h *LeadHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.leadService.GetDashboardStats()
	if err != nil {
		utils.InternalError(c, "获取统计数据失败")
		return
	}
	utils.Success(c, stats)
}
