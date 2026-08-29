package handlers

import (
	"github.com/gin-gonic/gin"

	"sportswear-platform/internal/services"
	"sportswear-platform/internal/utils"
)

// StorageSourceHandler 存储源处理器
type StorageSourceHandler struct {
	storageSourceService *services.StorageSourceService
}

// NewStorageSourceHandler 创建存储源处理器
func NewStorageSourceHandler(storageSourceService *services.StorageSourceService) *StorageSourceHandler {
	return &StorageSourceHandler{storageSourceService: storageSourceService}
}

// List 存储源列表
func (h *StorageSourceHandler) List(c *gin.Context) {
	sources, err := h.storageSourceService.List()
	if err != nil {
		utils.InternalError(c, "获取存储源列表失败")
		return
	}
	utils.Success(c, sources)
}

// Get 获取存储源
func (h *StorageSourceHandler) Get(c *gin.Context) {
	source, err := h.storageSourceService.Get(c.Param("id"))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, source)
}

// Create 创建存储源
func (h *StorageSourceHandler) Create(c *gin.Context) {
	var req services.StorageSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	source, err := h.storageSourceService.Create(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, source)
}

// Update 更新存储源
func (h *StorageSourceHandler) Update(c *gin.Context) {
	var req services.StorageSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	source, err := h.storageSourceService.Update(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, source)
}

// Delete 删除存储源
func (h *StorageSourceHandler) Delete(c *gin.Context) {
	if err := h.storageSourceService.Delete(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除存储源失败")
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}
