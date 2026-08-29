package models

// StorageSource 媒体存储源配置
// 支持配置 IP+PORT，图片/视频通过「存储源地址 + 路径」即可路由访问
type StorageSource struct {
	BaseModel
	Name        string      `gorm:"type:varchar(100);not null" json:"name"`
	Code        string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // 唯一标识，如 local / minio / cdn
	Type        MediaSource `gorm:"type:varchar(20);default:local" json:"type"`        // local / minio / external
	Protocol    string      `gorm:"type:varchar(10);default:http" json:"protocol"`     // http / https
	Host        string      `gorm:"type:varchar(255);not null" json:"host"`            // IP 或域名
	Port        string      `gorm:"type:varchar(10)" json:"port"`                      // 端口，可空
	BasePath    string      `gorm:"type:varchar(255)" json:"base_path"`                // 基础路径（如 /bucket 或 /uploads）
	Bucket      string      `gorm:"type:varchar(100)" json:"bucket"`                   // 对象存储桶名（MinIO/S3）
	AccessKey   string      `gorm:"type:varchar(255)" json:"access_key"`
	SecretKey   string      `gorm:"type:varchar(255)" json:"-"`
	IsDefault   bool        `gorm:"default:false" json:"is_default"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
	Description string      `gorm:"type:varchar(500)" json:"description"`
}

// BuildBaseURL 构建存储源基础 URL（协议 + IP/域名 + 端口 + 基础路径）
// 示例：http://192.168.1.100:9000/sportswear
func (s *StorageSource) BuildBaseURL() string {
	scheme := s.Protocol
	if scheme == "" {
		scheme = "http"
	}
	host := s.Host
	if host == "" {
		return ""
	}
	base := scheme + "://" + host
	if s.Port != "" {
		base += ":" + s.Port
	}
	if s.BasePath != "" {
		if s.BasePath[0] != '/' {
			base += "/"
		}
		base += s.BasePath
	}
	if s.Bucket != "" {
		if s.BasePath == "" {
			base += "/"
		}
		if s.BasePath != "" && s.BasePath[len(s.BasePath)-1] != '/' {
			base += "/"
		}
		base += s.Bucket
	}
	return base
}
