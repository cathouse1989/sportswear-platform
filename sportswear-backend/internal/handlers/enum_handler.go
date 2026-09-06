package handlers

import (
	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// EnumHandler 枚举字典处理器（数据字典：值域 + 多语言翻译的后台管理）
type EnumHandler struct {
	enumService *services.EnumService
}

// NewEnumHandler 创建枚举字典处理器
func NewEnumHandler(enumService *services.EnumService) *EnumHandler {
	return &EnumHandler{enumService: enumService}
}

// GetPublicEnums 获取指定语言的枚举字典（公开接口，门户/后台通用，平铺 key）
func (h *EnumHandler) GetPublicEnums(c *gin.Context) {
	lang := middleware.GetLang(c)
	enums, err := h.enumService.GetPublicEnums(lang)
	if err != nil {
		utils.InternalError(c, "获取枚举字典失败")
		return
	}
	groups, err := h.enumService.GetPublicEnumGroups()
	if err != nil {
		groups = map[string][]string{}
	}
	utils.Success(c, gin.H{"language": lang, "enums": enums, "groups": groups})
}

// ListTypes 枚举类型列表（后台管理）
func (h *EnumHandler) ListTypes(c *gin.Context) {
	types, err := h.enumService.ListTypes()
	if err != nil {
		utils.InternalError(c, "获取枚举类型失败")
		return
	}
	utils.Success(c, types)
}

// UpsertType 创建或更新枚举类型
func (h *EnumHandler) UpsertType(c *gin.Context) {
	var req services.EnumTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	t, err := h.enumService.UpsertType(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, t)
}

// DeleteType 删除枚举类型
func (h *EnumHandler) DeleteType(c *gin.Context) {
	if err := h.enumService.DeleteType(c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

// ListItems 指定类型的枚举项列表
func (h *EnumHandler) ListItems(c *gin.Context) {
	typeID := c.Query("type_id")
	if typeID == "" {
		utils.BadRequest(c, "缺少 type_id 参数")
		return
	}
	items, err := h.enumService.ListItems(typeID)
	if err != nil {
		utils.InternalError(c, "获取枚举项失败")
		return
	}
	utils.Success(c, items)
}

// UpsertItem 创建或更新枚举项
func (h *EnumHandler) UpsertItem(c *gin.Context) {
	var req services.EnumItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	item, err := h.enumService.UpsertItem(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, item)
}

// DeleteItem 删除枚举项
func (h *EnumHandler) DeleteItem(c *gin.Context) {
	if err := h.enumService.DeleteItem(c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}
