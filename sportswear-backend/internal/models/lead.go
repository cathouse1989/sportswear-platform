package models

import (
	"time"

	"github.com/google/uuid"
)

// Lead 询盘
type Lead struct {
	BaseModel
	// 基础信息
	Name           string `gorm:"type:varchar(100);not null" json:"name"`
	Company        string `gorm:"type:varchar(200)" json:"company"`
	Email          string `gorm:"type:varchar(255);not null" json:"email"`
	Phone          string `gorm:"type:varchar(50)" json:"phone"`
	WhatsApp       string `gorm:"type:varchar(50)" json:"whatsapp"`
	Country        string `gorm:"type:varchar(100)" json:"country"`
	CompanyWebsite string `gorm:"type:varchar(500)" json:"company_website"`
	// 项目需求
	ProjectType     string     `gorm:"type:varchar(50)" json:"project_type"` // oem, odm, private_label, product, quote, contact
	ProductCategory string     `gorm:"type:varchar(200)" json:"product_category"`
	ProductID       *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Quantity        int        `gorm:"default:0" json:"quantity"`
	Budget          string     `gorm:"type:varchar(100)" json:"budget"`
	TargetDate      string     `gorm:"type:varchar(50)" json:"target_date"`
	Message         string     `gorm:"type:text" json:"message"`
	// 附件
	Attachments string `gorm:"type:text" json:"attachments"` // JSON 数组
	// 来源追踪
	Source      string `gorm:"type:varchar(100)" json:"source"`
	Medium      string `gorm:"type:varchar(100)" json:"medium"`
	Campaign    string `gorm:"type:varchar(100)" json:"campaign"`
	Keyword     string `gorm:"type:varchar(200)" json:"keyword"`
	FirstPage   string `gorm:"type:varchar(500)" json:"first_page"`
	LandingPage string `gorm:"type:varchar(500)" json:"landing_page"`
	Device      string `gorm:"type:varchar(50)" json:"device"`
	Language    string `gorm:"type:varchar(10)" json:"language"`
	IP          string `gorm:"type:varchar(50)" json:"ip"`
	// 状态
	Status         LeadStatus `gorm:"type:varchar(20);default:new" json:"status"`
	Score          int        `gorm:"default:0" json:"score"`
	ScoreLevel     LeadScore  `gorm:"type:varchar(20);default:low" json:"score_level"`
	AssignedTo     *uuid.UUID `gorm:"type:uuid" json:"assigned_to"`
	AssignedToUser *User      `gorm:"-" json:"assigned_to_user,omitempty"`
	NextFollowUp   *time.Time `json:"next_follow_up"`
	// 关联
	FollowUps []LeadFollowUp `json:"follow_ups,omitempty"`
}

// LeadFollowUp 询盘跟进
type LeadFollowUp struct {
	BaseModel
	LeadID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"lead_id"`
	UserID       uuid.UUID      `gorm:"type:uuid" json:"user_id"`
	User         *User          `gorm:"-" json:"user,omitempty"`
	Method       FollowUpMethod `gorm:"type:varchar(20);default:email" json:"method"`
	Content      string         `gorm:"type:text" json:"content"`
	NextFollowUp *time.Time     `json:"next_follow_up"`
}

// Quote 报价（第二阶段）
type Quote struct {
	BaseModel
	LeadID         uuid.UUID  `gorm:"type:uuid;index;not null" json:"lead_id"`
	QuoteNumber    string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"quote_number"`
	ProductID      *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Quantity       int        `gorm:"default:0" json:"quantity"`
	UnitPrice      float64    `gorm:"type:decimal(10,2)" json:"unit_price"`
	TotalPrice     float64    `gorm:"type:decimal(10,2)" json:"total_price"`
	Currency       string     `gorm:"type:varchar(10);default:USD" json:"currency"`
	MOQ            int        `gorm:"default:0" json:"moq"`
	LeadTime       string     `gorm:"type:varchar(100)" json:"lead_time"`
	PaymentTerms   string     `gorm:"type:varchar(200)" json:"payment_terms"`
	ShippingMethod string     `gorm:"type:varchar(200)" json:"shipping_method"`
	ValidUntil     *time.Time `json:"valid_until"`
	Status         string     `gorm:"type:varchar(20);default:draft" json:"status"` // draft, sent, viewed, accepted, rejected, expired
	Notes          string     `gorm:"type:text" json:"notes"`
}

// Notification 通知
type Notification struct {
	BaseModel
	Type       NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Title      string           `gorm:"type:varchar(200);not null" json:"title"`
	Content    string           `gorm:"type:text" json:"content"`
	UserID     *uuid.UUID       `gorm:"type:uuid;index" json:"user_id"`
	IsRead     bool             `gorm:"default:false" json:"is_read"`
	ReadAt     *time.Time       `json:"read_at"`
	EntityType string           `gorm:"type:varchar(50)" json:"entity_type"`
	EntityID   *uuid.UUID       `gorm:"type:uuid" json:"entity_id"`
}

// OperationLog 操作日志
type OperationLog struct {
	BaseModel
	UserID      uuid.UUID     `gorm:"type:uuid;index" json:"user_id"`
	User        *User         `gorm:"-" json:"user,omitempty"`
	Operation   OperationType `gorm:"type:varchar(20);not null" json:"operation"`
	Module      string        `gorm:"type:varchar(50)" json:"module"`
	EntityType  string        `gorm:"type:varchar(50)" json:"entity_type"`
	EntityID    string        `gorm:"type:varchar(100)" json:"entity_id"`
	Before      string        `gorm:"type:text" json:"before"`
	After       string        `gorm:"type:text" json:"after"`
	IP          string        `gorm:"type:varchar(50)" json:"ip"`
	UserAgent   string        `gorm:"type:varchar(500)" json:"user_agent"`
	Description string        `gorm:"type:varchar(500)" json:"description"`
}
