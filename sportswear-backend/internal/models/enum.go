package models

import "github.com/google/uuid"

// SysEnumType 枚举类型（数据字典目录）
// 作为数据库枚举字段的「唯一事实源」：集中管理枚举值域、排序、启停与多语言翻译。
// 门户/后台的枚举标签均由此驱动，取代散落在 Go 常量、静态语言包与前端硬编码中的三处重复定义。
type SysEnumType struct {
	BaseModel
	Code        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`         // 语义化编码，如 product.type / product.gender
	Name        string `gorm:"type:varchar(100);not null" json:"name"`                     // 展示名
	Module      string `gorm:"type:varchar(50);index;default:product" json:"module"`       // 分组：product/order/system...
	I18nPrefix  string `gorm:"type:varchar(100)" json:"i18n_prefix"`                       // 门户 localizedEnum 的 key 前缀，如 product.type_options
	Description string `gorm:"type:text" json:"description"`
	IsSystem    bool   `gorm:"default:false" json:"is_system"`                             // 系统内置：禁删、禁改 code
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	Items       []SysEnumItem `gorm:"foreignKey:TypeID" json:"items,omitempty"`
}

// SysEnumItem 枚举项（数据字典值）
// value 为存储值（如 oem / odm），label 为默认语言（en）文案，translations 为其余语言覆盖。
type SysEnumItem struct {
	BaseModel
	TypeID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_enum_items_type_value" json:"type_id"`
	Value        string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_enum_items_type_value" json:"value"`
	Label        string    `gorm:"type:varchar(200);not null" json:"label"`  // 默认语言（en）文案
	Translations string    `gorm:"type:jsonb" json:"translations"`           // {"zh":"...","es":"...","fr":"..."}
	Color        string    `gorm:"type:varchar(20)" json:"color"`            // 后台 tag 颜色：blue/green/red/gold...
	IsDefault    bool      `gorm:"default:false" json:"is_default"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	SortOrder    int       `gorm:"default:0" json:"sort_order"`
}
