package models

import (
	"github.com/google/uuid"
)

// MediaSource 媒体来源类型
type MediaSource string

const (
	MediaSourceLocal    MediaSource = "local"    // 本地磁盘存储
	MediaSourceMinIO    MediaSource = "minio"    // MinIO/S3 对象存储
	MediaSourceExternal MediaSource = "external" // 外部链接（YouTube、第三方 CDN 等）
)

// VideoProvider 视频平台提供方
type VideoProvider string

const (
	VideoProviderSelf     VideoProvider = "self"     // 自托管视频（mp4/webm）
	VideoProviderYouTube  VideoProvider = "youtube"  // YouTube 视频
	VideoProviderVimeo    VideoProvider = "vimeo"    // Vimeo 视频
	VideoProviderBilibili VideoProvider = "bilibili" // B 站视频
)

// Media 媒体资源
type Media struct {
	BaseModel
	OriginalName   string        `gorm:"type:varchar(255);not null" json:"original_name"`
	FileName       string        `gorm:"type:varchar(255);not null" json:"file_name"`
	FileType       string        `gorm:"type:varchar(50);not null" json:"file_type"` // image/jpeg, image/png, video/mp4, application/pdf
	FileSize       int64         `gorm:"default:0" json:"file_size"`
	Width          int           `gorm:"default:0" json:"width"`
	Height         int           `gorm:"default:0" json:"height"`
	URL            string        `gorm:"type:varchar(1000);not null" json:"url"` // 最终可访问 URL（解析后）
	Path           string        `gorm:"type:varchar(500)" json:"path"`          // 存储路径（local/minio 时，如 products/2026-08/uuid.jpg）
	Thumbnail      string        `gorm:"type:varchar(1000)" json:"thumbnail"`
	WebPURL        string        `gorm:"type:varchar(1000)" json:"webp_url"`
	AVIFURL        string        `gorm:"type:varchar(1000)" json:"avif_url"`
	Alt            string        `gorm:"type:varchar(200)" json:"alt"`
	Title          string        `gorm:"type:varchar(200)" json:"title"`
	Description    string        `gorm:"type:text" json:"description"`
	Category       MediaCategory `gorm:"type:varchar(50);default:public" json:"category"`
	Type           MediaType     `gorm:"type:varchar(20);default:image" json:"type"`
	Source         MediaSource   `gorm:"type:varchar(20);default:local" json:"source"`  // local / minio / external
	Provider       VideoProvider `gorm:"type:varchar(20);default:self" json:"provider"` // self / youtube / vimeo / bilibili
	ExternalID     string        `gorm:"type:varchar(100)" json:"external_id"`          // 视频平台视频 ID（youtube/vimeo）
	EmbedURL       string        `gorm:"type:varchar(1000)" json:"embed_url"`           // 嵌入 URL（视频平台 iframe 地址）
	UploadedBy     *uuid.UUID    `gorm:"type:uuid" json:"uploaded_by"`
	UploadedByUser *User         `gorm:"-" json:"uploaded_by_user,omitempty"`
	IsPublic       bool          `gorm:"default:true" json:"is_public"`
}

// ResolveAccessURL 解析媒体最终访问 URL（结合存储源配置）
// source=local/minio: baseURL + path；source=external: 直接用 URL
func (m *Media) ResolveAccessURL(baseURL string) string {
	if m.Source == MediaSourceExternal {
		return m.URL
	}
	// local/minio：baseURL + path，若 path 为空则直接用 URL
	if m.Path != "" {
		return baseURL + "/" + m.Path
	}
	return m.URL
}
