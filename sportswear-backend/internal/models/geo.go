package models

// GeoLocale 国家本地化映射（后台「区域设置」可配置）
// 取代 geo_service 中的内置 countryLocaleMap，运营可在后台维护国家 → 语言/货币/时区映射。
type GeoLocale struct {
	BaseModel
	Country     string `gorm:"type:varchar(100);not null" json:"country"`
	CountryISO2 string `gorm:"type:varchar(10);uniqueIndex;not null" json:"country_iso2"`
	Language    string `gorm:"type:varchar(10)" json:"language"`
	Currency    string `gorm:"type:varchar(10)" json:"currency"`
	Timezone    string `gorm:"type:varchar(50)" json:"timezone"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}
