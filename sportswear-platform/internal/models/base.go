package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 生成 UUID
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// TranslationStatus 翻译状态
type TranslationStatus string

const (
	TranslationStatusUntranslated TranslationStatus = "untranslated"
	TranslationStatusDraft        TranslationStatus = "draft"
	TranslationStatusTranslated   TranslationStatus = "translated"
	TranslationStatusPending      TranslationStatus = "pending"
	TranslationStatusPublished    TranslationStatus = "published"
)

// ContentStatus 内容状态
type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusReview    ContentStatus = "review"
	ContentStatusPublished ContentStatus = "published"
	ContentStatusOffline   ContentStatus = "offline"
	ContentStatusArchived  ContentStatus = "archived"
)

// PageType 页面类型
type PageType string

const (
	PageTypeHome            PageType = "home"
	PageTypeNormal          PageType = "normal"
	PageTypeProduct         PageType = "product"
	PageTypeProductCategory PageType = "product_category"
	PageTypeOEM             PageType = "oem"
	PageTypeODM             PageType = "odm"
	PageTypePrivateLabel    PageType = "private_label"
	PageTypeFactory         PageType = "factory"
	PageTypeProduction      PageType = "production"
	PageTypeBlog            PageType = "blog"
	PageTypeCase            PageType = "case"
	PageTypeFAQ             PageType = "faq"
	PageTypeContact         PageType = "contact"
	PageTypeSEO             PageType = "seo_landing"
)

// LeadStatus 询盘状态
type LeadStatus string

const (
	LeadStatusNew        LeadStatus = "new"
	LeadStatusContacted  LeadStatus = "contacted"
	LeadStatusConfirmed  LeadStatus = "confirmed"
	LeadStatusQuoted     LeadStatus = "quoted"
	LeadStatusSampling   LeadStatus = "sampling"
	LeadStatusNegotiating LeadStatus = "negotiating"
	LeadStatusWon        LeadStatus = "won"
	LeadStatusLost       LeadStatus = "lost"
	LeadStatusSpam       LeadStatus = "spam"
)

// LeadScore 询盘评分等级
type LeadScore string

const (
	LeadScoreHigh   LeadScore = "high"
	LeadScoreMedium LeadScore = "medium"
	LeadScoreNormal LeadScore = "normal"
	LeadScoreLow    LeadScore = "low"
)

// MediaType 媒体类型
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeFile  MediaType = "file"
)

// MediaCategory 媒体分类
type MediaCategory string

const (
	MediaCategoryHome       MediaCategory = "home"
	MediaCategoryProduct    MediaCategory = "product"
	MediaCategoryOEM        MediaCategory = "oem"
	MediaCategoryODM        MediaCategory = "odm"
	MediaCategoryFactory    MediaCategory = "factory"
	MediaCategoryProduction MediaCategory = "production"
	MediaCategoryCase       MediaCategory = "case"
	MediaCategoryBlog       MediaCategory = "blog"
	MediaCategoryCert       MediaCategory = "certification"
	MediaCategoryPublic     MediaCategory = "public"
)

// Gender 适用性别
type Gender string

const (
	GenderUnisex Gender = "unisex"
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderKids   Gender = "kids"
)

// ProductStatus 产品状态
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusPublished ProductStatus = "published"
	ProductStatusOffline   ProductStatus = "offline"
)

// ProductType 产品类型
type ProductType string

const (
	ProductTypeOEM ProductType = "oem"
	ProductTypeODM ProductType = "odm"
	ProductTypeBoth ProductType = "both"
)

// FollowUpMethod 跟进方式
type FollowUpMethod string

const (
	FollowUpEmail     FollowUpMethod = "email"
	FollowUpWhatsApp  FollowUpMethod = "whatsapp"
	FollowUpPhone     FollowUpMethod = "phone"
	FollowUpMeeting   FollowUpMethod = "meeting"
	FollowUpOther     FollowUpMethod = "other"
)

// NotificationType 通知类型
type NotificationType string

const (
	NotificationNewLead     NotificationType = "new_lead"
	NotificationLeadUpdate  NotificationType = "lead_update"
	NotificationQuote       NotificationType = "quote"
	NotificationSystem      NotificationType = "system"
)

// OperationType 操作类型
type OperationType string

const (
	OperationCreate OperationType = "create"
	OperationUpdate OperationType = "update"
	OperationDelete OperationType = "delete"
	OperationPublish OperationType = "publish"
	OperationUnpublish OperationType = "unpublish"
	OperationLogin   OperationType = "login"
	OperationLogout  OperationType = "logout"
	OperationUpload  OperationType = "upload"
	OperationOther   OperationType = "other"
)