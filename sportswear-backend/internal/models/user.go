package models

import (
	"time"

	"github.com/google/uuid"
)

// User 用户
type User struct {
	BaseModel
	Email       string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	Name        string     `gorm:"type:varchar(100)" json:"name"`
	Avatar      string     `gorm:"type:varchar(500)" json:"avatar"`
	Phone       string     `gorm:"type:varchar(50)" json:"phone"`
	Status      string     `gorm:"type:varchar(20);default:active" json:"status"` // active, disabled
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `gorm:"type:varchar(50)" json:"last_login_ip"`
	Roles       []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// Role 角色
type Role struct {
	BaseModel
	Name        string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Code        string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	AppID       *uuid.UUID   `gorm:"type:uuid;index" json:"app_id"` // 归属应用，NULL=全局角色（跨应用可用）
	App         *App         `json:"app,omitempty"`
	Description string       `gorm:"type:varchar(500)" json:"description"`
	IsActive    bool         `gorm:"default:true;index" json:"is_active"` // 角色启用/禁用状态
	SortOrder   int          `gorm:"default:0" json:"sort_order"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"-"`
}

// Permission 权限
type Permission struct {
	BaseModel
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Code        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`
	AppID       *uuid.UUID `gorm:"type:uuid;index" json:"app_id"` // 权限归属应用，NULL=平台级权限
	App         *App       `json:"app,omitempty"`
	Module      string     `gorm:"type:varchar(50);not null" json:"module"`
	Description string     `gorm:"type:varchar(500)" json:"description"`
	Roles       []Role     `gorm:"many2many:role_permissions;" json:"-"`
}

// UserRole 用户角色关联
type UserRole struct {
	UserID uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	RoleID uuid.UUID  `gorm:"type:uuid;primaryKey" json:"role_id"`
	AppID  *uuid.UUID `gorm:"type:uuid;index" json:"app_id"` // 授权范围：NULL=所有应用生效，否则仅该应用生效
}

// RolePermission 角色权限关联
type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey" json:"permission_id"`
}

// Language 语言
type Language struct {
	BaseModel
	Code       string `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"` // en, zh, es, fr, de
	Name       string `gorm:"type:varchar(100);not null" json:"name"`
	NativeName string `gorm:"type:varchar(100)" json:"native_name"`
	Direction  string `gorm:"type:varchar(3);default:ltr" json:"direction"` // ltr / rtl（阿拉伯语等从右到左）
	IsDefault  bool   `gorm:"default:false" json:"is_default"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
	Flag       string `gorm:"type:varchar(10)" json:"flag"`
}

// Setting 系统配置
type Setting struct {
	BaseModel
	Key         string `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Group       string `gorm:"type:varchar(50)" json:"group"`
	Description string `gorm:"type:varchar(500)" json:"description"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
}
