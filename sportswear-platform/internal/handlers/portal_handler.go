package handlers

import (
	"github.com/gin-gonic/gin"

	"sportswear-platform/internal/middleware"
	"sportswear-platform/internal/services"
	"sportswear-platform/internal/utils"
)

// PortalHandler 门户展示配置处理器
type PortalHandler struct {
	portalService *services.PortalService
}

// NewPortalHandler 创建门户配置处理器
func NewPortalHandler(portalService *services.PortalService) *PortalHandler {
	return &PortalHandler{portalService: portalService}
}

// ==================== 后台：版本化发布 ====================

// SaveDraft 保存页面草稿版本（不影响线上）
func (h *PortalHandler) SaveDraft(c *gin.Context) {
	var req struct {
		Snapshot interface{} `json:"snapshot" binding:"required"`
		Note     string      `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	version, err := h.portalService.SaveDraft(c.Param("id"), req.Snapshot, userID.String(), req.Note)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, version)
}

// PublishVersion 发布指定版本（前台立即生效）
func (h *PortalHandler) PublishVersion(c *gin.Context) {
	if err := h.portalService.PublishVersion(c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"published": true})
}

// ListVersions 页面版本历史
func (h *PortalHandler) ListVersions(c *gin.Context) {
	versions, err := h.portalService.ListVersions(c.Param("id"))
	if err != nil {
		utils.InternalError(c, "获取版本历史失败")
		return
	}
	utils.Success(c, versions)
}

// RollbackVersion 回滚到历史版本
func (h *PortalHandler) RollbackVersion(c *gin.Context) {
	if err := h.portalService.RollbackVersion(c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"rolled_back": true})
}

// ==================== 后台：主题配置 ====================

// ListThemeConfigs 后台：获取全部主题配置项列表（含 key/value/name/group）
func (h *PortalHandler) ListThemeConfigs(c *gin.Context) {
	configs, err := h.portalService.ListThemeConfigs()
	if err != nil {
		utils.InternalError(c, "获取主题配置失败")
		return
	}
	utils.Success(c, configs)
}

// UpdateThemeConfig 更新主题配置项
func (h *PortalHandler) UpdateThemeConfig(c *gin.Context) {
	var req struct {
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.portalService.UpdateThemeConfig(c.Param("key"), req.Value); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"updated": true})
}

// ==================== 公开：主题配置 ====================

// GetThemeConfig 获取主题配置（前端渲染视觉风格）
func (h *PortalHandler) GetThemeConfig(c *gin.Context) {
	config, err := h.portalService.GetThemeConfig()
	if err != nil {
		utils.InternalError(c, "获取主题配置失败")
		return
	}
	utils.Success(c, config)
}
