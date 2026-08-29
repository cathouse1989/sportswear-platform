package services

import (
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// GeoService 地理定位服务
// 面向全球站点：根据访客 IP 所在国家，自动推荐语言与货币
//
// 国家识别来源（按优先级）：
// 1. CF-IPCountry 头（Cloudflare）
// 2. X-Country 头（其他网关/反代）
// 3. X-Geo-Country 头
// 4. 默认美国（en/USD）
type GeoService struct {
	db *gorm.DB
}

// NewGeoService 创建地理定位服务
func NewGeoService(db *gorm.DB) *GeoService {
	return &GeoService{db: db}
}

// CountryLocale 国家本地化映射
type CountryLocale struct {
	Country     string `json:"country"`
	CountryISO2 string `json:"country_iso2"`
	Language    string `json:"language"`
	Currency    string `json:"currency"`
	Timezone    string `json:"timezone"`
}

// 内置国家 → 语言/货币/时区映射（常见贸易国家）
var countryLocaleMap = map[string]CountryLocale{
	"US": {Country: "United States", CountryISO2: "US", Language: "en", Currency: "USD", Timezone: "America/New_York"},
	"GB": {Country: "United Kingdom", CountryISO2: "GB", Language: "en", Currency: "GBP", Timezone: "Europe/London"},
	"CN": {Country: "China", CountryISO2: "CN", Language: "zh", Currency: "CNY", Timezone: "Asia/Shanghai"},
	"HK": {Country: "Hong Kong", CountryISO2: "HK", Language: "zh", Currency: "USD", Timezone: "Asia/Hong_Kong"},
	"TW": {Country: "Taiwan", CountryISO2: "TW", Language: "zh", Currency: "USD", Timezone: "Asia/Taipei"},
	"JP": {Country: "Japan", CountryISO2: "JP", Language: "ja", Currency: "JPY", Timezone: "Asia/Tokyo"},
	"KR": {Country: "South Korea", CountryISO2: "KR", Language: "ko", Currency: "USD", Timezone: "Asia/Seoul"},
	"DE": {Country: "Germany", CountryISO2: "DE", Language: "de", Currency: "EUR", Timezone: "Europe/Berlin"},
	"FR": {Country: "France", CountryISO2: "FR", Language: "fr", Currency: "EUR", Timezone: "Europe/Paris"},
	"ES": {Country: "Spain", CountryISO2: "ES", Language: "es", Currency: "EUR", Timezone: "Europe/Madrid"},
	"IT": {Country: "Italy", CountryISO2: "IT", Language: "it", Currency: "EUR", Timezone: "Europe/Rome"},
	"NL": {Country: "Netherlands", CountryISO2: "NL", Language: "en", Currency: "EUR", Timezone: "Europe/Amsterdam"},
	"RU": {Country: "Russia", CountryISO2: "RU", Language: "ru", Currency: "USD", Timezone: "Europe/Moscow"},
	"AE": {Country: "United Arab Emirates", CountryISO2: "AE", Language: "en", Currency: "USD", Timezone: "Asia/Dubai"},
	"SA": {Country: "Saudi Arabia", CountryISO2: "SA", Language: "en", Currency: "USD", Timezone: "Asia/Riyadh"},
	"IN": {Country: "India", CountryISO2: "IN", Language: "en", Currency: "USD", Timezone: "Asia/Kolkata"},
	"AU": {Country: "Australia", CountryISO2: "AU", Language: "en", Currency: "AUD", Timezone: "Australia/Sydney"},
	"NZ": {Country: "New Zealand", CountryISO2: "NZ", Language: "en", Currency: "NZD", Timezone: "Pacific/Auckland"},
	"CA": {Country: "Canada", CountryISO2: "CA", Language: "en", Currency: "CAD", Timezone: "America/Toronto"},
	"BR": {Country: "Brazil", CountryISO2: "BR", Language: "pt", Currency: "USD", Timezone: "America/Sao_Paulo"},
	"MX": {Country: "Mexico", CountryISO2: "MX", Language: "es", Currency: "USD", Timezone: "America/Mexico_City"},
	"ZA": {Country: "South Africa", CountryISO2: "ZA", Language: "en", Currency: "USD", Timezone: "Africa/Johannesburg"},
	"SG": {Country: "Singapore", CountryISO2: "SG", Language: "en", Currency: "USD", Timezone: "Asia/Singapore"},
	"MY": {Country: "Malaysia", CountryISO2: "MY", Language: "en", Currency: "USD", Timezone: "Asia/Kuala_Lumpur"},
	"TH": {Country: "Thailand", CountryISO2: "TH", Language: "th", Currency: "USD", Timezone: "Asia/Bangkok"},
	"VN": {Country: "Vietnam", CountryISO2: "VN", Language: "vi", Currency: "USD", Timezone: "Asia/Ho_Chi_Minh"},
	"ID": {Country: "Indonesia", CountryISO2: "ID", Language: "en", Currency: "USD", Timezone: "Asia/Jakarta"},
	"PH": {Country: "Philippines", CountryISO2: "PH", Language: "en", Currency: "USD", Timezone: "Asia/Manila"},
	"TR": {Country: "Turkey", CountryISO2: "TR", Language: "tr", Currency: "USD", Timezone: "Europe/Istanbul"},
	"PL": {Country: "Poland", CountryISO2: "PL", Language: "pl", Currency: "EUR", Timezone: "Europe/Warsaw"},
}

// DetectLocale 根据国家 ISO2 代码推荐本地化信息
func (s *GeoService) DetectLocale(countryISO2 string) CountryLocale {
	if locale, ok := countryLocaleMap[countryISO2]; ok {
		return locale
	}
	// 默认
	return CountryLocale{Country: "United States", CountryISO2: "US", Language: "en", Currency: "USD", Timezone: "America/New_York"}
}

// ResolveRecommended 根据国家代码推荐语言/货币（自动激活的语言、货币存在性校验）
func (s *GeoService) ResolveRecommended(countryISO2 string) CountryLocale {
	locale := s.DetectLocale(countryISO2)

	// 校验语言是否启用
	var langCount int64
	s.db.Model(&models.Language{}).Where("code = ? AND is_active = ?", locale.Language, true).Count(&langCount)
	if langCount == 0 {
		locale.Language = "en"
	}

	// 校验货币是否启用
	var curCount int64
	s.db.Model(&models.Currency{}).Where("code = ? AND is_active = ?", locale.Currency, true).Count(&curCount)
	if curCount == 0 {
		def, _ := s.GetDefaultCurrencyCode()
		if def != "" {
			locale.Currency = def
		} else {
			locale.Currency = "USD"
		}
	}

	return locale
}

// GetDefaultCurrencyCode 获取默认货币代码
func (s *GeoService) GetDefaultCurrencyCode() (string, error) {
	var currency models.Currency
	err := s.db.Where("is_default = ? AND is_active = ?", true, true).First(&currency).Error
	if err != nil {
		return "", err
	}
	return currency.Code, nil
}

// GetActiveLanguages 获取启用语言列表（供前端语言切换器）
func (s *GeoService) GetActiveLanguages() ([]models.Language, error) {
	var languages []models.Language
	err := s.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&languages).Error
	return languages, err
}
