package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"sportswear-backend/internal/config"
	"sportswear-backend/internal/models"
)

// UploadService 文件上传服务（MinIO/S3 对象存储）
type UploadService struct {
	cfg      *config.Config
	endpoint string
	useSSL   bool
}

// NewUploadService 创建上传服务
func NewUploadService(cfg *config.Config) *UploadService {
	return &UploadService{
		cfg:      cfg,
		endpoint: cfg.MinIO.Endpoint,
		useSSL:   cfg.MinIO.UseSSL,
	}
}

// UploadResult 上传结果
type UploadResult struct {
	URL       string           `json:"url"`
	FileName  string           `json:"file_name"`
	FileSize  int64            `json:"file_size"`
	FileType  string           `json:"file_type"`
	MediaType models.MediaType `json:"media_type"`
}

// allowedExtensions 允许的文件扩展名（安全白名单）
var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".webp": true, ".avif": true, ".svg": true,
	".mp4": true, ".webm": true, ".mov": true,
	".pdf": true, ".doc": true, ".docx": true,
	".xls": true, ".xlsx": true, ".zip": true,
}

// ValidateFile 校验文件安全：扩展名白名单 + 大小限制
func (s *UploadService) ValidateFile(fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}
	maxSize := s.cfg.Upload.MaxSize * 1024 * 1024 // MB -> bytes
	if fileHeader.Size > maxSize {
		return fmt.Errorf("文件超过大小限制 %dMB", s.cfg.Upload.MaxSize)
	}
	return nil
}

// GenerateStoragePath 生成存储路径与随机文件名（防路径遍历 + 防重名）
func (s *UploadService) GenerateStoragePath(originalName string, category string) (storagePath, fileName string) {
	ext := filepath.Ext(originalName)
	now := time.Now()
	fileName = fmt.Sprintf("%s%s", uuid.New().String(), ext)
	// 按分类/年月组织：products/2026-08/uuid.jpg
	storagePath = fmt.Sprintf("%s/%s/%s", category, now.Format("2006-01"), fileName)
	return storagePath, fileName
}

// PublicURL 生成公开访问 URL
// 生产环境建议通过 CDN 域名访问，本地通过 MinIO 网关
func (s *UploadService) PublicURL(storagePath string) string {
	scheme := "http"
	if s.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.endpoint, s.cfg.MinIO.Bucket, storagePath)
}

// SaveToLocal 本地保存（MinIO 未部署时的降级方案）
// 文件保存到 uploads/ 目录，通过静态文件服务对外提供
func (s *UploadService) SaveToLocal(file io.Reader, storagePath string) error {
	localPath := fmt.Sprintf("uploads/%s", storagePath)
	if err := ensureDir(localPath); err != nil {
		return err
	}
	return saveFile(file, localPath)
}

// LocalURL 本地文件的公开访问 URL
func (s *UploadService) LocalURL(storagePath string) string {
	return "/uploads/" + storagePath
}

// GetMediaType 检测媒体类型
func GetMediaType(filename string) models.MediaType {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif", ".svg":
		return models.MediaTypeImage
	case ".mp4", ".webm", ".mov":
		return models.MediaTypeVideo
	default:
		return models.MediaTypeFile
	}
}

// GetFileType 检测 MIME 类型
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeMap := map[string]string{
		".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
		".gif": "image/gif", ".webp": "image/webp", ".avif": "image/avif",
		".svg": "image/svg+xml", ".mp4": "video/mp4", ".webm": "video/webm",
		".pdf": "application/pdf", ".zip": "application/zip",
	}
	if t, ok := mimeMap[ext]; ok {
		return t
	}
	return "application/octet-stream"
}

func ensureDir(path string) error {
	idx := strings.LastIndex(path, "/")
	if idx > 0 {
		return mkdirAll(path[:idx])
	}
	return nil
}
