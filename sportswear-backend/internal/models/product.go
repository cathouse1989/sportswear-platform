package models

import (
	"github.com/google/uuid"
)

// Product 产品
type Product struct {
	BaseModel
	SKU        string        `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku"`
	Slug       string        `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	Name       string        `gorm:"-" json:"name"` // 本地化名称（按语言解析后填充）
	CategoryID *uuid.UUID    `gorm:"type:uuid" json:"category_id"`
	Category   *Category     `json:"category,omitempty"`
	Type       ProductType   `gorm:"type:varchar(20);default:both" json:"type"`
	Gender     Gender        `gorm:"type:varchar(20);default:unisex" json:"gender"`
	Status     ProductStatus `gorm:"type:varchar(20);default:draft" json:"status"`
	IsFeatured bool          `gorm:"default:false" json:"is_featured"`
	IsNew      bool          `gorm:"default:false" json:"is_new"`
	SortOrder  int           `gorm:"default:0" json:"sort_order"`
	CoverImage string        `gorm:"type:varchar(500)" json:"cover_image"`
	// 详细信息
	Brief        string `gorm:"type:text" json:"brief"`
	Description  string `gorm:"type:text" json:"description"`
	Features     string `gorm:"type:text" json:"features"`
	Usage        string `gorm:"type:text" json:"usage"`
	Material     string `gorm:"type:varchar(200)" json:"material"`
	Composition  string `gorm:"type:varchar(200)" json:"composition"`
	Weight       string `gorm:"type:varchar(50)" json:"weight"`
	Elasticity   string `gorm:"type:varchar(50)" json:"elasticity"`
	Fit          string `gorm:"type:varchar(50)" json:"fit"`
	SupportLevel string `gorm:"type:varchar(50)" json:"support_level"`
	Season       string `gorm:"type:varchar(50)" json:"season"`
	SizeRange    string `gorm:"type:varchar(100)" json:"size_range"`
	// MOQ
	SampleMOQ     int `gorm:"default:1" json:"sample_moq"`
	ProductionMOQ int `gorm:"default:300" json:"production_moq"`
	ColorMOQ      int `gorm:"default:100" json:"color_moq"`
	SizeMOQ       int `gorm:"default:100" json:"size_moq"`
	// 关联
	Translations   []ProductTranslation   `json:"translations,omitempty"`
	Images         []ProductImage         `json:"images,omitempty"`
	Videos         []ProductVideo         `json:"videos,omitempty"`
	Specs          []ProductSpec          `json:"specs,omitempty"`
	Customizations []ProductCustomization `json:"customizations,omitempty"`
	Series         []Series               `gorm:"many2many:product_series;" json:"series,omitempty"`
	Fabrics        []Fabric               `gorm:"many2many:product_fabrics;" json:"fabrics,omitempty"`
	SEO            *SEO                   `gorm:"-" json:"seo,omitempty"`
}

// ProductTranslation 产品翻译
type ProductTranslation struct {
	BaseModel
	ProductID   uuid.UUID         `gorm:"type:uuid;index;not null" json:"product_id"`
	Language    string            `gorm:"type:varchar(10);index;not null" json:"language"`
	Name        string            `gorm:"type:varchar(200);not null" json:"name"`
	Brief       string            `gorm:"type:text" json:"brief"`
	Description string            `gorm:"type:text" json:"description"`
	Features    string            `gorm:"type:text" json:"features"`
	Usage       string            `gorm:"type:text" json:"usage"`
	// 规格类字段的多语言翻译（主表存英文源，此处按语言覆盖）
	Material     string `gorm:"type:varchar(200)" json:"material"`
	Composition  string `gorm:"type:varchar(200)" json:"composition"`
	Weight       string `gorm:"type:varchar(50)" json:"weight"`
	Elasticity   string `gorm:"type:varchar(50)" json:"elasticity"`
	Fit          string `gorm:"type:varchar(50)" json:"fit"`
	SupportLevel string `gorm:"type:varchar(50)" json:"support_level"`
	Season       string `gorm:"type:varchar(50)" json:"season"`
	SizeRange    string `gorm:"type:varchar(100)" json:"size_range"`
	// SortOrder 语言维度的展示排序权重（0 = 跟随产品全局 sort_order）。
	// 场景：不同地区/语种市场有各自的特色运动与运动服装，
	// 各语言站点可配置不同的产品展示顺序（越小越靠前）。
	SortOrder int               `gorm:"default:0" json:"sort_order"`
	Status    TranslationStatus `gorm:"type:varchar(20);default:untranslated" json:"status"`
}

// Category 产品分类（无限级）
type Category struct {
	BaseModel
	ParentID  *uuid.UUID `gorm:"type:uuid" json:"parent_id"`
	Name      string     `gorm:"type:varchar(200);not null" json:"name"`
	Slug      string     `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	SortOrder int        `gorm:"default:0" json:"sort_order"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	Image     string     `gorm:"type:varchar(500)" json:"image"`
	Children  []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Products  []Product  `json:"-"`
}

// Series 产品系列
type Series struct {
	BaseModel
	Name      string        `gorm:"type:varchar(200);not null" json:"name"`
	Slug      string        `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	SortOrder int           `gorm:"default:0" json:"sort_order"`
	Status    ProductStatus `gorm:"type:varchar(20);default:draft;index" json:"status"` // 前端展示状态标识
	IsActive  bool          `gorm:"default:true" json:"is_active"`
	Products  []Product     `gorm:"many2many:product_series;" json:"-"`
}

// ProductImage 产品图片
type ProductImage struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	Type      string    `gorm:"type:varchar(50)" json:"type"` // main, gallery, model, detail, fabric, craft, size, packaging
	URL       string    `gorm:"type:varchar(500);not null" json:"url"`
	Thumbnail string    `gorm:"type:varchar(500)" json:"thumbnail"`
	Alt       string    `gorm:"type:varchar(200)" json:"alt"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
}

// ProductVideo 产品视频
type ProductVideo struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	Type      string    `gorm:"type:varchar(50)" json:"type"` // product, production, craft, usage
	URL       string    `gorm:"type:varchar(500);not null" json:"url"`
	Cover     string    `gorm:"type:varchar(500)" json:"cover"`
	Title     string    `gorm:"type:varchar(200)" json:"title"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
}

// ProductSpec 产品规格
type ProductSpec struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Value     string    `gorm:"type:varchar(200)" json:"value"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
}

// ProductCustomization 产品定制能力
type ProductCustomization struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;index;not null" json:"product_id"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"` // logo, color, fabric, pattern, print, embroidery, label, hangtag, packaging, size, fit, zipper, button, belt, accessory
	IsEnabled bool      `gorm:"default:true" json:"is_enabled"`
	Note      string    `gorm:"type:varchar(500)" json:"note"`
}

// Fabric 面料
type Fabric struct {
	BaseModel
	Name            string        `gorm:"type:varchar(200);not null" json:"name"`
	Code            string        `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Status          ProductStatus `gorm:"type:varchar(20);default:draft;index" json:"status"` // 前端展示状态标识
	Composition     string        `gorm:"type:varchar(200)" json:"composition"`
	Weight          string        `gorm:"type:varchar(50)" json:"weight"`
	Elasticity      string        `gorm:"type:varchar(50)" json:"elasticity"`
	Breathability   string        `gorm:"type:varchar(50)" json:"breathability"`
	MoistureWicking string        `gorm:"type:varchar(50)" json:"moisture_wicking"`
	Softness        string        `gorm:"type:varchar(50)" json:"softness"`
	Compression     string        `gorm:"type:varchar(50)" json:"compression"`
	UVProtection    string        `gorm:"type:varchar(50)" json:"uv_protection"`
	EcoFriendly     string        `gorm:"type:varchar(50)" json:"eco_friendly"`
	Description     string        `gorm:"type:text" json:"description"`
	IsActive        bool          `gorm:"default:true" json:"is_active"`
	Products        []Product     `gorm:"many2many:product_fabrics;" json:"-"`
}

// SEO SEO 配置（支持多语言）
type SEO struct {
	BaseModel
	EntityType    string    `gorm:"type:varchar(50);index;not null" json:"entity_type"` // product, page, blog, case, category
	EntityID      uuid.UUID `gorm:"type:uuid;index;not null" json:"entity_id"`
	Language      string    `gorm:"type:varchar(10);index;default:en" json:"language"` // 语言维度
	Title         string    `gorm:"type:varchar(200)" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	Keywords      string    `gorm:"type:varchar(500)" json:"keywords"`
	Canonical     string    `gorm:"type:varchar(500)" json:"canonical"`
	Robots        string    `gorm:"type:varchar(100)" json:"robots"`
	OGTitle       string    `gorm:"type:varchar(200)" json:"og_title"`
	OGDescription string    `gorm:"type:text" json:"og_description"`
	OGImage       string    `gorm:"type:varchar(500)" json:"og_image"`
	SchemaData    string    `gorm:"type:text" json:"schema_data"`
}
