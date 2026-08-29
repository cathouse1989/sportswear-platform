package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-platform/internal/services"
	"sportswear-platform/internal/utils"
)

// TrashHandler 回收站处理器
type TrashHandler struct {
	trashService *services.TrashService
}

// NewTrashHandler 创建回收站处理器
func NewTrashHandler(trashService *services.TrashService) *TrashHandler {
	return &TrashHandler{trashService: trashService}
}

// ListTrash 回收站列表
func (h *TrashHandler) ListTrash(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	entityType := c.Query("entity_type")

	items, total, err := h.trashService.ListTrash(entityType, page, pageSize)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.SuccessPage(c, items, page, pageSize, total)
}

// RestoreItem 恢复回收站数据
func (h *TrashHandler) RestoreItem(c *gin.Context) {
	var req struct {
		EntityType string `json:"entity_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.trashService.Restore(req.EntityType, c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"restored": true})
}

// PurgeItem 彻底删除（不可恢复）
func (h *TrashHandler) PurgeItem(c *gin.Context) {
	var req struct {
		EntityType string `json:"entity_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.trashService.Purge(req.EntityType, c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"purged": true})
}

// EmptyTrash 清空回收站
func (h *TrashHandler) EmptyTrash(c *gin.Context) {
	entityType := c.Query("entity_type")

	count, err := h.trashService.EmptyTrash(entityType)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"deleted_count": count})
}
