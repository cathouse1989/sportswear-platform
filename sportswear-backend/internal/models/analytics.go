package models

import (
	"time"
)

// VisitLog 访问日志（流量监测，商业价值分析）
type VisitLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	// 访问者信息
	IP        string `gorm:"type:varchar(50);index" json:"ip"`
	Country   string `gorm:"type:varchar(100);index" json:"country"`
	Device    string `gorm:"type:varchar(20)" json:"device"` // desktop, mobile, tablet, bot
	Browser   string `gorm:"type:varchar(50)" json:"browser"`
	OS        string `gorm:"type:varchar(50)" json:"os"`
	UserAgent string `gorm:"type:varchar(500)" json:"user_agent"`
	// 访问内容
	Method     string `gorm:"type:varchar(10)" json:"method"`
	Path       string `gorm:"type:varchar(500);index" json:"path"`       // /api/v1/public/products/custom-yoga-leggings
	EntityType string `gorm:"type:varchar(50);index" json:"entity_type"` // product, page, blog, case, home
	EntitySlug string `gorm:"type:varchar(200);index" json:"entity_slug"`
	Status     int    `json:"status"` // HTTP 状态码
	LatencyMs  int64  `json:"latency_ms"`
	// 来源追踪（SEO/广告归因）
	Referer     string `gorm:"type:varchar(500)" json:"referer"`
	Source      string `gorm:"type:varchar(100);index" json:"source"` // google, bing, direct, social
	Medium      string `gorm:"type:varchar(100)" json:"medium"`       // organic, cpc, referral
	Keyword     string `gorm:"type:varchar(200)" json:"keyword"`      // 搜索关键词
	UtmCampaign string `gorm:"type:varchar(200)" json:"utm_campaign"`
	SessionID   string `gorm:"type:varchar(64);index" json:"session_id"`
}

// TableName 表名
func (VisitLog) TableName() string {
	return "visit_logs"
}
