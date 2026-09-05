package models

import (
	"time"

	"github.com/google/uuid"
)

// Page 页面
type Page struct {
	BaseModel
	Title        string            `gorm:"type:varchar(200);not null" json:"title"`
	Slug         string            `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	Type         PageType          `gorm:"type:varchar(50);default:normal" json:"type"`
	Status       ContentStatus     `gorm:"type:varchar(20);default:draft" json:"status"`
	Template     string            `gorm:"type:varchar(100)" json:"template"`
	SortOrder    int               `gorm:"default:0" json:"sort_order"`
	PublishedAt  *time.Time        `json:"published_at"`
	Modules      []PageModule      `json:"modules,omitempty"`
	Translations []PageTranslation `json:"translations,omitempty"`
	SEO          *SEO              `gorm:"-" json:"seo,omitempty"`
}

// PageTranslation 页面翻译
type PageTranslation struct {
	BaseModel
	PageID   uuid.UUID         `gorm:"type:uuid;index;not null" json:"page_id"`
	Language string            `gorm:"type:varchar(10);index;not null" json:"language"`
	Title    string            `gorm:"type:varchar(200)" json:"title"`
	Content  string            `gorm:"type:text" json:"content"`
	Status   TranslationStatus `gorm:"type:varchar(20);default:untranslated" json:"status"`
}

// PageModule 页面模块
type PageModule struct {
	BaseModel
	PageID    uuid.UUID `gorm:"type:uuid;index;not null" json:"page_id"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"` // banner, product_recommend, oem, odm, factory, production, case, certification, blog, contact, text, image, video
	Title     string    `gorm:"type:varchar(200)" json:"title"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	IsVisible bool      `gorm:"default:true" json:"is_visible"`
	Config    string    `gorm:"type:jsonb" json:"config"` // JSONB 模块配置
}

// Navigation 导航
type Navigation struct {
	BaseModel
	Name      string       `gorm:"type:varchar(100);not null" json:"name"`
	Type      string       `gorm:"type:varchar(20);default:header" json:"type"` // header, footer
	URL       string       `gorm:"type:varchar(500)" json:"url"`
	Target    string       `gorm:"type:varchar(20);default:_self" json:"target"`
	SortOrder int          `gorm:"default:0" json:"sort_order"`
	IsVisible bool         `gorm:"default:true" json:"is_visible"`
	ParentID  *uuid.UUID   `gorm:"type:uuid" json:"parent_id"`
	PageID    *uuid.UUID   `gorm:"type:uuid;index" json:"page_id,omitempty"` // 关联页面（可空外键，实现 Page↔Nav 关联）
	Page      *Page        `gorm:"foreignKey:PageID" json:"page,omitempty"`  // 预加载关联页面信息
	Children  []Navigation `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// Blog 博客
type Blog struct {
	BaseModel
	Title        string            `gorm:"type:varchar(200);not null" json:"title"`
	Slug         string            `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	Category     string            `gorm:"type:varchar(50)" json:"category"` // oem_guide, odm_guide, fabric, trend, sourcing, brand, production, industry
	Tags         string            `gorm:"type:varchar(500)" json:"tags"`
	Author       string            `gorm:"type:varchar(100)" json:"author"`
	CoverImage   string            `gorm:"type:varchar(500)" json:"cover_image"`
	Content      string            `gorm:"type:text" json:"content"`
	Status       ContentStatus     `gorm:"type:varchar(20);default:draft" json:"status"`
	PublishedAt  *time.Time        `json:"published_at"`
	Summary      string            `gorm:"-" json:"summary,omitempty"` // 列表摘要（运行时计算，非 DB 列）
	Translations []BlogTranslation `json:"translations,omitempty"`
	SEO          *SEO              `gorm:"-" json:"seo,omitempty"`
}

// BlogTranslation 博客翻译
type BlogTranslation struct {
	BaseModel
	BlogID   uuid.UUID         `gorm:"type:uuid;index;not null" json:"blog_id"`
	Language string            `gorm:"type:varchar(10);index;not null" json:"language"`
	Title    string            `gorm:"type:varchar(200)" json:"title"`
	Content  string            `gorm:"type:text" json:"content"`
	Status   TranslationStatus `gorm:"type:varchar(20);default:untranslated" json:"status"`
}

// Case 案例
type Case struct {
	BaseModel
	Title          string            `gorm:"type:varchar(200);not null" json:"title"`
	Slug           string            `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	ClientIndustry string            `gorm:"type:varchar(100)" json:"client_industry"`
	ProjectType    string            `gorm:"type:varchar(50)" json:"project_type"` // oem, odm, private_label
	Products       string            `gorm:"type:varchar(500)" json:"products"`
	ClientNeed     string            `gorm:"type:text" json:"client_need"`
	Problem        string            `gorm:"type:text" json:"problem"`
	Solution       string            `gorm:"type:text" json:"solution"`
	Process        string            `gorm:"type:text" json:"process"`
	Result         string            `gorm:"type:text" json:"result"`
	CoverImage     string            `gorm:"type:varchar(500)" json:"cover_image"`
	Status         ContentStatus     `gorm:"type:varchar(20);default:draft" json:"status"`
	PublishedAt    *time.Time        `json:"published_at"`
	Translations   []CaseTranslation `json:"translations,omitempty"`
	SEO            *SEO              `gorm:"-" json:"seo,omitempty"`
}

// CaseTranslation 案例翻译
type CaseTranslation struct {
	BaseModel
	CaseID     uuid.UUID         `gorm:"type:uuid;index;not null" json:"case_id"`
	Language   string            `gorm:"type:varchar(10);index;not null" json:"language"`
	Title      string            `gorm:"type:varchar(200)" json:"title"`
	ClientNeed string            `gorm:"type:text" json:"client_need"`
	Problem    string            `gorm:"type:text" json:"problem"`
	Solution   string            `gorm:"type:text" json:"solution"`
	Process    string            `gorm:"type:text" json:"process"`
	Result     string            `gorm:"type:text" json:"result"`
	Status     TranslationStatus `gorm:"type:varchar(20);default:untranslated" json:"status"`
}

// FAQ FAQ
type FAQ struct {
	BaseModel
	Question     string           `gorm:"type:text;not null" json:"question"`
	Answer       string           `gorm:"type:text" json:"answer"`
	Category     string           `gorm:"type:varchar(50)" json:"category"` // moq, oem, odm, sample, payment, production, logistics, fabric, quality, certification
	Language     string           `gorm:"type:varchar(10);default:en" json:"language"`
	SortOrder    int              `gorm:"default:0" json:"sort_order"`
	IsActive     bool             `gorm:"default:true" json:"is_active"`
	Translations []FAQTranslation `json:"translations,omitempty"`
}

// FAQTranslation FAQ 翻译
type FAQTranslation struct {
	BaseModel
	FAQID    uuid.UUID         `gorm:"type:uuid;index;not null" json:"faq_id"`
	Language string            `gorm:"type:varchar(10);index;not null" json:"language"`
	Question string            `gorm:"type:text" json:"question"`
	Answer   string            `gorm:"type:text" json:"answer"`
	Status   TranslationStatus `gorm:"type:varchar(20);default:untranslated" json:"status"`
}

// Factory 工厂
type Factory struct {
	BaseModel
	Name                 string        `gorm:"type:varchar(200);not null" json:"name"`
	Status               ProductStatus `gorm:"type:varchar(20);default:draft;index" json:"status"` // 前端展示状态标识
	Location             string        `gorm:"type:varchar(200)" json:"location"`
	Area                 string        `gorm:"type:varchar(100)" json:"area"`
	Employees            int           `gorm:"default:0" json:"employees"`
	ProductionLines      int           `gorm:"default:0" json:"production_lines"`
	Equipment            string        `gorm:"type:text" json:"equipment"`
	MonthlyCapacity      string        `gorm:"type:varchar(100)" json:"monthly_capacity"`
	AnnualCapacity       string        `gorm:"type:varchar(100)" json:"annual_capacity"`
	Warehouse            string        `gorm:"type:varchar(200)" json:"warehouse"`
	QualityManagement    string        `gorm:"type:text" json:"quality_management"`
	ProductionCapability string        `gorm:"type:text" json:"production_capability"`
	Description          string        `gorm:"type:text" json:"description"`
	Image                string        `gorm:"type:varchar(500)" json:"image"`
	IsActive             bool          `gorm:"default:true" json:"is_active"`
}


// Certification 认证
type Certification struct {
	BaseModel
	Name        string        `gorm:"type:varchar(200);not null" json:"name"`
	Status      ProductStatus `gorm:"type:varchar(20);default:draft;index" json:"status"` // 前端展示状态标识
	Code        string        `gorm:"type:varchar(100)" json:"code"`
	IssueDate   string        `gorm:"type:varchar(50)" json:"issue_date"`
	ExpiryDate  string        `gorm:"type:varchar(50)" json:"expiry_date"`
	Image       string        `gorm:"type:varchar(500)" json:"image"`
	PDF         string        `gorm:"type:varchar(500)" json:"pdf"`
	Description string        `gorm:"type:text" json:"description"`
	IsActive    bool          `gorm:"default:true" json:"is_active"`
}

// ProductionProcess 生产流程
type ProductionProcess struct {
	BaseModel
	Name        string        `gorm:"type:varchar(200);not null" json:"name"`
	Status      ProductStatus `gorm:"type:varchar(20);default:draft;index" json:"status"` // 前端展示状态标识
	Description string        `gorm:"type:text" json:"description"`
	Image       string        `gorm:"type:varchar(500)" json:"image"`
	Video       string        `gorm:"type:varchar(500)" json:"video"`
	SortOrder   int           `gorm:"default:0" json:"sort_order"`
	IsActive    bool          `gorm:"default:true" json:"is_active"`
}
