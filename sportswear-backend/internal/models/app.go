package models

import (
	"github.com/google/uuid"
)

// App 子应用（应用维度权限隔离的基础）
// 例如：admin 后台、客户中心、移动端、小程序、第三方集成等
type App struct {
	BaseModel
	Code        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // admin, customer_portal, mobile, mini_program, third_party
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	APIKey      string `gorm:"type:varchar(100);uniqueIndex;not null" json:"api_key"` // 子应用接入密钥
	Secret      string `gorm:"type:varchar(255)" json:"-"`                            // 应用密钥（签名用）
	CallbackURL string `gorm:"type:varchar(500)" json:"callback_url"`
	IsActive    bool   `gorm:"default:true;index" json:"is_active"` // 应用启用/禁用
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	Desc        string `gorm:"type:varchar(500)" json:"desc"`
}

// AppUser 用户与应用的关联（同一用户在不同应用可拥有不同角色）
type AppUser struct {
	BaseModel
	AppID  uuid.UUID `gorm:"type:uuid;index;not null" json:"app_id"`
	UserID uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	User   *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	App    *App      `gorm:"foreignKey:AppID" json:"app,omitempty"`
	Status string    `gorm:"type:varchar(20);default:active" json:"status"` // 该用户在此应用中的状态 active/disabled
}
