package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// I18nHandler 国际化处理器
type I18nHandler struct {
	i18nService *services.I18nService
}

// NewI18nHandler 创建国际化处理器
func NewI18nHandler(i18nService *services.I18nService) *I18nHandler {
	return &I18nHandler{i18nService: i18nService}
}

// GetDictionary 获取指定语言的词条字典（公开接口，前端一键切换语言）
func (h *I18nHandler) GetDictionary(c *gin.Context) {
	lang := middleware.GetLang(c)
	dict, err := h.i18nService.GetEntriesByLang(lang)
	if err != nil {
		utils.InternalError(c, "获取词条失败")
		return
	}
	utils.Success(c, gin.H{
		"language":   lang,
		"dictionary": dict,
	})
}

// ListEntries 词条列表（后台管理）
func (h *I18nHandler) ListEntries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	language := c.Query("language")
	module := c.Query("module")
	keyword := c.Query("keyword")

	entries, total, err := h.i18nService.ListEntries(page, pageSize, language, module, keyword)
	if err != nil {
		utils.InternalError(c, "获取词条列表失败")
		return
	}
	utils.SuccessPage(c, entries, page, pageSize, total)
}

// UpsertEntry 创建或更新词条
func (h *I18nHandler) UpsertEntry(c *gin.Context) {
	var req services.I18nEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	entry, err := h.i18nService.UpsertEntry(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, entry)
}

// DeleteEntry 删除词条
func (h *I18nHandler) DeleteEntry(c *gin.Context) {
	if err := h.i18nService.DeleteEntry(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除词条失败")
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}
