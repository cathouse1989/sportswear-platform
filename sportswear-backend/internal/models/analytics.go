package models

import (
	"time"
)

// VisitLog 访问日志（流量监测，商业价值分析）
// 客户维度信息：设备（PC/移动/Pad）、IP、国家、语言、时间、来源归因等。
// 按季度分表存储：visit_logs_YYYYMMDD（YYYYMMDD=季度起始日），见 database/partition.go。
type VisitLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	// 访问者信息（客户维度）
	IP          string `gorm:"type:varchar(50);index" json:"ip"`
	Country     string `gorm:"type:varchar(100);index" json:"country"` // 国家（ISO2 代码 / 国家名）
	Language    string `gorm:"type:varchar(10);index" json:"language"` // 客户端语言（en/zh/es/...）
	Device      string `gorm:"type:varchar(20)" json:"device"`         // desktop, mobile, tablet, bot
	DeviceModel string `gorm:"type:varchar(100)" json:"device_model"`  // 设备型号（iPhone / iPad / SM-xxx / Windows PC ...）
	Browser     string `gorm:"type:varchar(50)" json:"browser"`
	OS          string `gorm:"type:varchar(50)" json:"os"`
	UserAgent   string `gorm:"type:varchar(500)" json:"user_agent"`
	// 访问内容
	Method     string `gorm:"type:varchar(10)" json:"method"`
	Path       string `gorm:"type:varchar(500);index" json:"path"`       // /api/v1/public/products/custom-yoga-leggings
	EntityType string `gorm:"type:varchar(50);index" json:"entity_type"` // product, page, blog, case, home, social_click
	EntitySlug string `gorm:"type:varchar(200);index" json:"entity_slug"`
	EntityID   string `gorm:"type:varchar(100)" json:"entity_id"`     // 实体 UUID（如产品 ID）
	EntityName string `gorm:"type:varchar(255)" json:"entity_name"`   // 实体名称（如产品名，便于后台检索）
	Status     int    `json:"status"`                                 // HTTP 状态码
	LatencyMs  int64  `json:"latency_ms"`
	// 来源追踪（SEO/广告归因，用于广告转化分析）
	Referer    string `gorm:"type:varchar(500)" json:"referer"`
	Source     string `gorm:"type:varchar(100);index" json:"source"` // google, bing, direct, social, referral...
	Medium     string `gorm:"type:varchar(100)" json:"medium"`       // organic, cpc, referral, click
	Keyword    string `gorm:"type:varchar(200)" json:"keyword"`      // 搜索关键词
	UtmSource  string `gorm:"type:varchar(100);index" json:"utm_source"`
	UtmMedium  string `gorm:"type:varchar(100)" json:"utm_medium"`
	UtmCampaign string `gorm:"type:varchar(200);index" json:"utm_campaign"`
	UtmContent string `gorm:"type:varchar(200)" json:"utm_content"`
	UtmTerm    string `gorm:"type:varchar(200)" json:"utm_term"`
	SessionID  string `gorm:"type:varchar(64);index" json:"session_id"`
}

// TableName 表名（基础表；实际写入按季度分表，见 database.QuarterTableName）
func (VisitLog) TableName() string {
	return "visit_logs"
}
