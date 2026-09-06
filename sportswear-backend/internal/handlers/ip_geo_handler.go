package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// IPGeoHandler IP 地理库处理器（维护 IP 段 → 国家映射）
type IPGeoHandler struct {
	geoIPService *services.GeoIPService
}

// NewIPGeoHandler 创建 IP 地理库处理器
func NewIPGeoHandler(geoIPService *services.GeoIPService) *IPGeoHandler {
	return &IPGeoHandler{geoIPService: geoIPService}
}

// ListIPGeoRanges 分页查询 IP 地理库
func (h *IPGeoHandler) ListIPGeoRanges(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	items, total, err := h.geoIPService.ListIPGeoRanges(page, pageSize)
	if err != nil {
		utils.InternalError(c, "获取 IP 地理库失败")
		return
	}
	utils.SuccessPage(c, items, page, pageSize, total)
}

// ImportIPGeoRanges 批量导入 IP 段数据（按 start_ip 幂等 upsert）
func (h *IPGeoHandler) ImportIPGeoRanges(c *gin.Context) {
	var req struct {
		Items []services.IPGeoRangeItem `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	imported, err := h.geoIPService.ImportIPGeoRanges(req.Items)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"imported": imported})
}

// DeleteIPGeoRange 删除单条 IP 段
func (h *IPGeoHandler) DeleteIPGeoRange(c *gin.Context) {
	if err := h.geoIPService.DeleteIPGeoRange(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除 IP 段失败")
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}
