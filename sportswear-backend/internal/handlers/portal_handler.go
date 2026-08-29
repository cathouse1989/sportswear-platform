package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

type PortalHandler struct {
	portalService *services.PortalService
	cache         *services.CacheService
	db            *gorm.DB
}

func NewPortalHandler(portalService *services.PortalService, cache *services.CacheService, db *gorm.DB) *PortalHandler {
	return &PortalHandler{portalService: portalService, cache: cache, db: db}
}

func (h *PortalHandler) SaveDraft(c *gin.Context) {
	var req struct {
		Snapshot interface{} `json:"snapshot" binding:"required"`
		Note     string      `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid params: "+err.Error())
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

func (h *PortalHandler) PublishVersion(c *gin.Context) {
	if err := h.portalService.PublishVersion(c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"published": true})
}

func (h *PortalHandler) ListVersions(c *gin.Context) {
	versions, err := h.portalService.ListVersions(c.Param("id"))
	if err != nil {
		utils.InternalError(c, "list versions failed")
		return
	}
	utils.Success(c, versions)
}

func (h *PortalHandler) RollbackVersion(c *gin.Context) {
	if err := h.portalService.RollbackVersion(c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"rolled_back": true})
}

func (h *PortalHandler) ListThemeConfigs(c *gin.Context) {
	configs, err := h.portalService.ListThemeConfigs()
	if err != nil {
		utils.InternalError(c, "list theme failed")
		return
	}
	utils.Success(c, configs)
}

func (h *PortalHandler) UpdateThemeConfig(c *gin.Context) {
	var req struct {
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid params: "+err.Error())
		return
	}
	if err := h.portalService.UpdateThemeConfig(c.Param("key"), req.Value); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"updated": true})
}

func (h *PortalHandler) GetPortalCacheStatus(c *gin.Context) {
	if h.cache == nil {
		utils.Success(c, gin.H{"redis_ok": false, "enabled": false, "using_cache": false, "last_refresh_at": nil})
		return
	}
	utils.Success(c, h.cache.Status())
}

func (h *PortalHandler) SetPortalCacheEnabled(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid params: "+err.Error())
		return
	}
	if h.cache == nil {
		utils.BadRequest(c, "cache service unavailable")
		return
	}
	if err := h.cache.SetPolicyEnabled(req.Enabled, req.Enabled); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, h.cache.Status())
}

func (h *PortalHandler) RefreshPortalCache(c *gin.Context) {
	if h.cache == nil || !h.cache.IsRedisOK() {
		utils.BadRequest(c, "redis unavailable")
		return
	}
	if !h.cache.RefreshPublicCache(h.db) {
		utils.BadRequest(c, "refresh failed")
		return
	}
	utils.Success(c, h.cache.Status())
}

func (h *PortalHandler) PublishPortalCache(c *gin.Context) {
	if h.cache == nil || !h.cache.IsRedisOK() {
		utils.BadRequest(c, "redis unavailable")
		return
	}
	if !h.cache.RefreshPublicCache(h.db) {
		utils.BadRequest(c, "publish failed")
		return
	}
	utils.Success(c, h.cache.Status())
}

func (h *PortalHandler) GetThemeConfig(c *gin.Context) {
	config, err := h.portalService.GetThemeConfig()
	if err != nil {
		utils.InternalError(c, "get theme failed")
		return
	}
	utils.Success(c, config)
}
