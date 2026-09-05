package models

import (
	"time"
)

// SubscriberStatus 订阅状态
type SubscriberStatus string

const (
	SubscriberStatusSubscribed   SubscriberStatus = "subscribed"
	SubscriberStatusUnsubscribed SubscriberStatus = "unsubscribed"
)

// Subscriber 订阅用户（门户"订阅更新"Newsletter 邮箱收集）
type Subscriber struct {
	BaseModel
	Email    string           `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Status   SubscriberStatus `gorm:"type:varchar(20);default:subscribed;index" json:"status"`
	Language string           `gorm:"type:varchar(10)" json:"language"`
	// 来源追踪（与 Lead 对齐，供广告转化/归因分析）
	Source      string `gorm:"type:varchar(100)" json:"source"`
	Medium      string `gorm:"type:varchar(100)" json:"medium"`
	Campaign    string `gorm:"type:varchar(100)" json:"campaign"`
	Keyword     string `gorm:"type:varchar(200)" json:"keyword"`
	LandingPage string `gorm:"type:varchar(500)" json:"landing_page"`
	Device      string `gorm:"type:varchar(50)" json:"device"`
	IP          string `gorm:"type:varchar(50)" json:"ip"`
	VisitorID   string `gorm:"type:varchar(64);index" json:"visitor_id"`
	// 时间戳
	SubscribedAt   *time.Time `json:"subscribed_at"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at"`
}
