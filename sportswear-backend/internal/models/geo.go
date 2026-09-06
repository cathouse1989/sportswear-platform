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

// IPGeoRange IP 段 → 国家映射（IP 地理定位库，存储在数据库中）
// IPv4 段用 32 位整数表示（start_ip ~ end_ip 闭区间），便于范围查询命中。
type IPGeoRange struct {
	BaseModel
	StartIP     int64  `gorm:"type:bigint;not null;index:idx_ip_geo_start" json:"start_ip"`
	EndIP       int64  `gorm:"type:bigint;not null" json:"end_ip"`
	Country     string `gorm:"type:varchar(10);not null;index" json:"country"` // ISO2 代码，如 US
	CountryName string `gorm:"type:varchar(100)" json:"country_name"`          // 国家名（中文/英文）
}
