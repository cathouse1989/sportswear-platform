package models

// Currency 货币
// 面向全球站点，支持多货币展示与汇率管理
type Currency struct {
	BaseModel
	Code               string  `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"` // USD/EUR/GBP/CNY/JPY...
	Name               string  `gorm:"type:varchar(100);not null" json:"name"`
	Symbol             string  `gorm:"type:varchar(10)" json:"symbol"`           // $ / € / £ / ¥
	Rate               float64 `gorm:"type:decimal(18,6);default:1" json:"rate"` // 对基准货币（USD）的汇率
	IsBase             bool    `gorm:"default:false" json:"is_base"`             // 是否基准货币
	IsDefault          bool    `gorm:"default:false" json:"is_default"`          // 是否默认展示货币
	IsActive           bool    `gorm:"default:true" json:"is_active"`
	DecimalPlaces      int     `gorm:"default:2" json:"decimal_places"`                        // 小数位数
	SymbolPosition     string  `gorm:"type:varchar(10);default:prefix" json:"symbol_position"` // prefix / suffix
	ThousandsSeparator string  `gorm:"type:varchar(10);default:," json:"thousands_separator"`
	DecimalSeparator   string  `gorm:"type:varchar(10);default:." json:"decimal_separator"`
	SortOrder          int     `gorm:"default:0" json:"sort_order"`
}

// DefaultCurrencies 默认货币
var DefaultCurrencies = []Currency{
	{Code: "USD", Name: "US Dollar", Symbol: "$", Rate: 1, IsBase: true, IsDefault: true, IsActive: true, DecimalPlaces: 2, SymbolPosition: "prefix", ThousandsSeparator: ",", DecimalSeparator: ".", SortOrder: 1},
	{Code: "EUR", Name: "Euro", Symbol: "€", Rate: 0.92, IsActive: true, DecimalPlaces: 2, SymbolPosition: "prefix", ThousandsSeparator: ",", DecimalSeparator: ".", SortOrder: 2},
	{Code: "GBP", Name: "British Pound", Symbol: "£", Rate: 0.79, IsActive: true, DecimalPlaces: 2, SymbolPosition: "prefix", ThousandsSeparator: ",", DecimalSeparator: ".", SortOrder: 3},
	{Code: "CNY", Name: "Chinese Yuan", Symbol: "¥", Rate: 7.25, IsActive: true, DecimalPlaces: 2, SymbolPosition: "prefix", ThousandsSeparator: ",", DecimalSeparator: ".", SortOrder: 4},
	{Code: "JPY", Name: "Japanese Yen", Symbol: "¥", Rate: 150.5, IsActive: true, DecimalPlaces: 0, SymbolPosition: "prefix", ThousandsSeparator: ",", DecimalSeparator: ".", SortOrder: 5},
	{Code: "AUD", Name: "Australian Dollar", Symbol: "A$", Rate: 1.52, IsActive: true, DecimalPlaces: 2, SymbolPosition: "prefix", ThousandsSeparator: ",", DecimalSeparator: ".", SortOrder: 6},
	{Code: "CAD", Name: "Canadian Dollar", Symbol: "C$", Rate: 1.36, IsActive: true, DecimalPlaces: 2, SymbolPosition: "prefix", ThousandsSeparator: ",", DecimalSeparator: ".", SortOrder: 7},
}
