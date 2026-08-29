package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

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
