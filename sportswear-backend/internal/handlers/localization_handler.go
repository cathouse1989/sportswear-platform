package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// LocalizationHandler 本地化处理器（SEO + 货币 + 地理定位 + 单位）
type LocalizationHandler struct {
	seoService      *services.SEOService
	currencyService *services.CurrencyService
	geoService      *services.GeoService
	unitService     *services.UnitService
}

// NewLocalizationHandler 创建本地化处理器
func NewLocalizationHandler(seoService *services.SEOService, currencyService *services.CurrencyService, geoService *services.GeoService, unitService *services.UnitService) *LocalizationHandler {
	return &LocalizationHandler{
		seoService:      seoService,
		currencyService: currencyService,
		geoService:      geoService,
		unitService:     unitService,
	}
}

// ==================== SEO 管理 ====================

// UpsertSEO 创建/更新 SEO 配置
func (h *LocalizationHandler) UpsertSEO(c *gin.Context) {
	var req services.SEOMultiLangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	seo, err := h.seoService.UpsertSEO(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, seo)
}

// ==================== 货币管理 ====================

// ListCurrencies 货币列表（公开）
func (h *LocalizationHandler) ListCurrencies(c *gin.Context) {
	currencies, err := h.currencyService.ListCurrencies(true)
	if err != nil {
		utils.InternalError(c, "获取货币列表失败")
		return
	}
	utils.Success(c, currencies)
}

// ConvertCurrency 货币转换（公开）
// GET /api/v1/public/currencies/convert?amount=100&from=USD&to=EUR
func (h *LocalizationHandler) ConvertCurrency(c *gin.Context) {
	amount, _ := strconv.ParseFloat(c.DefaultQuery("amount", "0"), 64)
	to := c.DefaultQuery("to", "USD")

	// from 固定基准货币 USD
	converted, err := h.currencyService.Convert(amount, to)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	formatted, _ := h.currencyService.FormatPrice(converted, to)

	utils.Success(c, gin.H{
		"amount":    amount,
		"from":      "USD",
		"to":        to,
		"converted": converted,
		"formatted": formatted,
	})
}

// ==================== 地理位置/本地化信息（公开，阶段四） ====================

// GetLocaleInfo 获取本地化信息（语言+货币+时区）
func (h *LocalizationHandler) GetLocaleInfo(c *gin.Context) {
	lang := middleware.GetLang(c)
	tz := middleware.GetTimezone(c)

	var currencyCode string
	if c := c.Query("currency"); c != "" {
		currencyCode = c
	} else {
		def, _ := h.currencyService.GetDefaultCurrency()
		if def != nil {
			currencyCode = def.Code
		}
	}

	utils.Success(c, gin.H{
		"language": lang,
		"timezone": tz,
		"currency": currencyCode,
	})
}

// GetGeolocation 地理定位：根据访客国家推荐语言/货币/时区
func (h *LocalizationHandler) GetGeolocation(c *gin.Context) {
	// 从请求头解析国家（Cloudflare CF-IPCountry / 反代 X-Country / X-Geo-Country）
	country := c.GetHeader("CF-IPCountry")
	if country == "" {
		country = c.GetHeader("X-Country")
	}
	if country == "" {
		country = c.GetHeader("X-Geo-Country")
	}
	if country == "" {
		country = c.Query("country")
	}

	locale := h.geoService.ResolveRecommended(country)

	utils.Success(c, gin.H{
		"country_iso2": locale.CountryISO2,
		"country":      locale.Country,
		"language":     locale.Language,
		"currency":     locale.Currency,
		"timezone":     locale.Timezone,
	})
}

// GetLanguages 获取启用语言列表（前端语言切换器）
func (h *LocalizationHandler) GetLanguages(c *gin.Context) {
	languages, err := h.geoService.GetActiveLanguages()
	if err != nil {
		utils.InternalError(c, "获取语言列表失败")
		return
	}
	utils.Success(c, languages)
}

// ==================== 单位本地化 ====================

// GetSizeChart 获取尺寸对照表
func (h *LocalizationHandler) GetSizeChart(c *gin.Context) {
	utils.Success(c, h.unitService.GetSizeChart())
}

// ConvertSize 尺码转换
// GET /api/v1/public/units/convert-size?size=M&from=US&to=EU
func (h *LocalizationHandler) ConvertSize(c *gin.Context) {
	size := c.Query("size")
	from := c.DefaultQuery("from", "US")
	to := c.DefaultQuery("to", "EU")

	result, err := h.unitService.ConvertSize(size, from, to)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{
		"size":   size,
		"from":   from,
		"to":     to,
		"result": result,
	})
}

// ConvertWeight 重量单位转换
// GET /api/v1/public/units/convert-weight?value=10&from=kg&to=lb
func (h *LocalizationHandler) ConvertWeight(c *gin.Context) {
	value, _ := strconv.ParseFloat(c.DefaultQuery("value", "0"), 64)
	from := c.DefaultQuery("from", "kg")
	to := c.DefaultQuery("to", "lb")

	result, err := h.unitService.ConvertWeight(value, services.WeightUnit(from), services.WeightUnit(to))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{
		"value":  value,
		"from":   from,
		"to":     to,
		"result": result,
	})
}

// ==================== 路由 SEO（列表页/落地页） ====================

// GetRouteSEO 获取指定路由的 SEO（公开接口，门户列表页/落地页使用）
// GET /api/v1/public/seo?route=products&lang=en
func (h *LocalizationHandler) GetRouteSEO(c *gin.Context) {
	route := c.Query("route")
	if route == "" {
		utils.BadRequest(c, "缺少 route 参数")
		return
	}
	lang := middleware.GetLang(c)
	seo, err := h.seoService.GetRouteSEO(route, lang)
	if err != nil {
		utils.Success(c, gin.H{"route": route, "language": lang, "seo": nil})
		return
	}
	utils.Success(c, gin.H{"route": route, "language": lang, "seo": seo})
}

// UpsertRoutesSEO 批量保存路由 SEO（后台「SEO 管理」页，按路由 × 语言矩阵）
func (h *LocalizationHandler) UpsertRoutesSEO(c *gin.Context) {
	var req struct {
		Route   string                      `json:"route" binding:"required"`
		Entries []services.RouteSEOEntryReq `json:"entries" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	for _, e := range req.Entries {
		if _, err := h.seoService.UpsertRouteSEO(req.Route, &e); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	}
	utils.Success(c, gin.H{"route": req.Route, "saved": len(req.Entries)})
}

// ListRouteSEO 获取指定路由的全部语言 SEO（后台 SEO 管理页）
func (h *LocalizationHandler) ListRouteSEO(c *gin.Context) {
	route := c.Query("route")
	if route == "" {
		utils.BadRequest(c, "缺少 route 参数")
		return
	}
	seos, err := h.seoService.ListRouteSEO(route)
	if err != nil {
		utils.InternalError(c, "获取路由 SEO 失败")
		return
	}
	utils.Success(c, gin.H{"route": route, "items": seos})
}

// ListEntitySEO 获取指定实体的全部语言 SEO（后台 SEO 管理页，实体详情维度）
func (h *LocalizationHandler) ListEntitySEO(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")
	if entityType == "" || entityID == "" {
		utils.BadRequest(c, "缺少 entity_type 或 entity_id 参数")
		return
	}
	id, err := uuid.Parse(entityID)
	if err != nil {
		utils.BadRequest(c, "entity_id 无效")
		return
	}
	seos, err := h.seoService.ListEntitySEO(entityType, id)
	if err != nil {
		utils.InternalError(c, "获取实体 SEO 失败")
		return
	}
	utils.Success(c, gin.H{"entity_type": entityType, "entity_id": entityID, "items": seos})
}

// UpsertEntitySEO 批量保存实体 SEO（后台 SEO 管理页，按实体 × 语言矩阵）
func (h *LocalizationHandler) UpsertEntitySEO(c *gin.Context) {
	var req struct {
		EntityType string                         `json:"entity_type" binding:"required"`
		EntityID   string                         `json:"entity_id" binding:"required"`
		Entries    []services.SEOMultiLangRequest `json:"entries" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	id, err := uuid.Parse(req.EntityID)
	if err != nil {
		utils.BadRequest(c, "entity_id 无效")
		return
	}
	for i := range req.Entries {
		req.Entries[i].EntityType = req.EntityType
		req.Entries[i].EntityID = id.String()
		if _, err := h.seoService.UpsertSEO(&req.Entries[i]); err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	}
	utils.Success(c, gin.H{"entity_type": req.EntityType, "entity_id": req.EntityID, "saved": len(req.Entries)})
}

// ==================== 国家本地化映射（区域设置） ====================

// ListGeoLocales 国家映射列表
func (h *LocalizationHandler) ListGeoLocales(c *gin.Context) {
	items, err := h.geoService.ListGeoLocales()
	if err != nil {
		utils.InternalError(c, "获取国家映射失败")
		return
	}
	utils.Success(c, items)
}

// UpsertGeoLocale 创建或更新国家映射
func (h *LocalizationHandler) UpsertGeoLocale(c *gin.Context) {
	var req services.GeoLocaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	gl, err := h.geoService.UpsertGeoLocale(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gl)
}

// DeleteGeoLocale 删除国家映射
func (h *LocalizationHandler) DeleteGeoLocale(c *gin.Context) {
	if err := h.geoService.DeleteGeoLocale(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除失败")
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}
