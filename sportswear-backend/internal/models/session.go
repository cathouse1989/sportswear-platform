package models

import (
	"time"

	"github.com/google/uuid"
)

// VisitorSession 访客会话（用于访客旅程分析）
// 一个访客的一次访问会话（30 分钟无活动则新会话）
type VisitorSession struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	VisitorID      string    `gorm:"type:varchar(64);index" json:"visitor_id"`
	SessionID      string    `gorm:"type:varchar(64);index" json:"session_id"`
	IP             string    `gorm:"type:varchar(50);index" json:"ip"`
	Country        string    `gorm:"type:varchar(100)" json:"country"`
	Device         string    `gorm:"type:varchar(20)" json:"device"`
	Browser        string    `gorm:"type:varchar(50)" json:"browser"`
	OS             string    `gorm:"type:varchar(50)" json:"os"`
	StartedAt      time.Time `json:"started_at"`
	EndedAt        time.Time `json:"ended_at"`
	PageViewsCount int       `json:"page_views_count"`
	IsConverted    bool      `json:"is_converted"`
	LeadID         *uuid.UUID `gorm:"type:uuid;index" json:"lead_id"`
	ConsentStatus  string    `gorm:"type:varchar(20)" json:"consent_status"`
}

// TableName 表名
func (VisitorSession) TableName() string {
	return "visitor_sessions"
}