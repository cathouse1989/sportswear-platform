package services

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// MediaService 媒体服务
type MediaService struct {
	db *gorm.DB
}

// NewMediaService 创建媒体服务
func NewMediaService(db *gorm.DB) *MediaService {
	return &MediaService{db: db}
}

// ListMedia 媒体列表
func (s *MediaService) ListMedia(page, pageSize int, category, mediaType, keyword string) ([]models.Media, int64, error) {
	var media []models.Media
	var total int64

	query := s.db.Model(&models.Media{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if mediaType != "" {
		query = query.Where("type = ?", mediaType)
	}
	if keyword != "" {
		query = query.Where("original_name LIKE ? OR title LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.
		Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&media).Error

	return media, total, err
}

// GetMedia 获取媒体
func (s *MediaService) GetMedia(id string) (*models.Media, error) {
	var media models.Media
	err := s.db.First(&media, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("媒体资源不存在")
	}
	return &media, nil
}

// CreateMedia 创建媒体记录
func (s *MediaService) CreateMedia(req *MediaRequest, userID uuid.UUID) (*models.Media, error) {
	media := models.Media{
		OriginalName: req.OriginalName,
		FileName:     req.FileName,
		FileType:     req.FileType,
		FileSize:     req.FileSize,
		Width:        req.Width,
		Height:       req.Height,
		URL:          req.URL,
		Thumbnail:    req.Thumbnail,
		WebPURL:      req.WebPURL,
		AVIFURL:      req.AVIFURL,
		Alt:          req.Alt,
		Title:        req.Title,
		Description:  req.Description,
		Category:     req.Category,
		Type:         req.Type,
		UploadedBy:   &userID,
		IsPublic:     req.IsPublic,
	}

	if err := s.db.Create(&media).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

// MediaRequest 媒体请求
type MediaRequest struct {
	OriginalName string               `json:"original_name" binding:"required"`
	FileName     string               `json:"file_name" binding:"required"`
	FileType     string               `json:"file_type" binding:"required"`
	FileSize     int64                `json:"file_size"`
	Width        int                  `json:"width"`
	Height       int                  `json:"height"`
	URL          string               `json:"url" binding:"required"`
	Thumbnail    string               `json:"thumbnail"`
	WebPURL      string               `json:"webp_url"`
	AVIFURL      string               `json:"avif_url"`
	Alt          string               `json:"alt"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Category     models.MediaCategory `json:"category"`
	Type         models.MediaType     `json:"type"`
	IsPublic     bool                 `json:"is_public"`
}

// UpdateMedia 更新媒体
func (s *MediaService) UpdateMedia(id string, req *UpdateMediaRequest) (*models.Media, error) {
	var media models.Media
	if err := s.db.First(&media, "id = ?", id).Error; err != nil {
		return nil, errors.New("媒体资源不存在")
	}

	updates := map[string]interface{}{}
	if req.OriginalName != "" {
		updates["original_name"] = req.OriginalName
	}
	if req.Alt != "" {
		updates["alt"] = req.Alt
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Thumbnail != "" {
		updates["thumbnail"] = req.Thumbnail
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if len(updates) > 0 {
		if err := s.db.Model(&media).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return &media, nil
}

// UpdateMediaRequest 更新媒体请求
type UpdateMediaRequest struct {
	OriginalName string               `json:"original_name"`
	Alt          string               `json:"alt"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Category     models.MediaCategory `json:"category"`
	Thumbnail    string               `json:"thumbnail"`
	IsPublic     *bool                `json:"is_public"`
}

// DeleteMedia 删除媒体
func (s *MediaService) DeleteMedia(id string) error {
	return s.db.Delete(&models.Media{}, "id = ?", id).Error
}

// CreateMediaDirect 直接创建媒体记录（用于文件上传后）
func (s *MediaService) CreateMediaDirect(media *models.Media) error {
	return s.db.Create(media).Error
}

// DetectMediaType 检测媒体类型
func DetectMediaType(filename string) models.MediaType {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif", ".svg", ".bmp":
		return models.MediaTypeImage
	case ".mp4", ".webm", ".mov", ".avi", ".mkv":
		return models.MediaTypeVideo
	default:
		return models.MediaTypeFile
	}
}

// DetectFileType 检测文件 MIME 类型
func DetectFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".avif":
		return "image/avif"
	case ".svg":
		return "image/svg+xml"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".pdf":
		return "application/pdf"
	case ".doc", ".docx":
		return "application/msword"
	case ".xls", ".xlsx":
		return "application/vnd.ms-excel"
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/x-rar-compressed"
	default:
		return "application/octet-stream"
	}
}
