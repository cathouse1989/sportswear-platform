package services

import (
	"errors"

	"gorm.io/gorm"

	"sportswear-platform/internal/models"
)

// StorageSourceService 媒体存储源管理服务
// 支持配置 IP+PORT，图片/视频通过「存储源地址 + 路径」即可路由访问
type StorageSourceService struct {
	db *gorm.DB
}

// NewStorageSourceService 创建存储源服务
func NewStorageSourceService(db *gorm.DB) *StorageSourceService {
	return &StorageSourceService{db: db}
}

// List 存储源列表
func (s *StorageSourceService) List() ([]models.StorageSource, error) {
	var sources []models.StorageSource
	err := s.db.Order("is_default DESC, created_at ASC").Find(&sources).Error
	return sources, err
}

// Get 获取存储源
func (s *StorageSourceService) Get(id string) (*models.StorageSource, error) {
	var source models.StorageSource
	if err := s.db.First(&source, "id = ?", id).Error; err != nil {
		return nil, errors.New("存储源不存在")
	}
	return &source, nil
}

// GetByCode 通过 code 获取存储源
func (s *StorageSourceService) GetByCode(code string) (*models.StorageSource, error) {
	var source models.StorageSource
	if err := s.db.Where("code = ? AND is_active = ?", code, true).First(&source).Error; err != nil {
		return nil, errors.New("存储源不存在")
	}
	return &source, nil
}

// GetDefault 获取默认存储源
func (s *StorageSourceService) GetDefault() (*models.StorageSource, error) {
	var source models.StorageSource
	err := s.db.Where("is_default = ? AND is_active = ?", true, true).First(&source).Error
	if err != nil {
		// 回退到第一个启用的存储源
		err = s.db.Where("is_active = ?", true).First(&source).Error
	}
	if err != nil {
		return nil, errors.New("无可用存储源")
	}
	return &source, nil
}

// Create 创建存储源
func (s *StorageSourceService) Create(req *StorageSourceRequest) (*models.StorageSource, error) {
	source := models.StorageSource{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Protocol:    req.Protocol,
		Host:        req.Host,
		Port:        req.Port,
		BasePath:    req.BasePath,
		Bucket:      req.Bucket,
		AccessKey:   req.AccessKey,
		SecretKey:   req.SecretKey,
		IsActive:    req.IsActive,
		Description: req.Description,
	}

	// 若设为默认，取消其他默认
	if req.IsDefault {
		s.db.Model(&models.StorageSource{}).Where("is_default = ?", true).Update("is_default", false)
		source.IsDefault = true
	}

	if err := s.db.Create(&source).Error; err != nil {
		return nil, err
	}
	return &source, nil
}

// Update 更新存储源
func (s *StorageSourceService) Update(id string, req *StorageSourceRequest) (*models.StorageSource, error) {
	var source models.StorageSource
	if err := s.db.First(&source, "id = ?", id).Error; err != nil {
		return nil, errors.New("存储源不存在")
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"code":        req.Code,
		"type":        req.Type,
		"protocol":    req.Protocol,
		"host":        req.Host,
		"port":        req.Port,
		"base_path":   req.BasePath,
		"bucket":      req.Bucket,
		"access_key":  req.AccessKey,
		"description": req.Description,
		"is_active":   req.IsActive,
	}
	if req.SecretKey != "" {
		updates["secret_key"] = req.SecretKey
	}
	if req.IsDefault {
		s.db.Model(&models.StorageSource{}).Where("is_default = ?", true).Update("is_default", false)
		updates["is_default"] = true
	}
	if err := s.db.Model(&source).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &source, nil
}

// Delete 删除存储源
func (s *StorageSourceService) Delete(id string) error {
	return s.db.Delete(&models.StorageSource{}, "id = ?", id).Error
}

// ResolveBaseURL 解析存储源基础 URL
func (s *StorageSourceService) ResolveBaseURL(code string) (string, error) {
	source, err := s.GetByCode(code)
	if err != nil {
		return "", err
	}
	return source.BuildBaseURL(), nil
}

// InitDefault 初始化默认存储源（本地）
func (s *StorageSourceService) InitDefault() error {
	var count int64
	s.db.Model(&models.StorageSource{}).Count(&count)
	if count > 0 {
		return nil
	}
	return s.db.Create(&models.StorageSource{
		Name:        "本地存储",
		Code:        "local",
		Type:        models.MediaSourceLocal,
		Protocol:    "http",
		Host:        "localhost",
		Port:        "8080",
		BasePath:    "uploads",
		IsDefault:   true,
		IsActive:    true,
		Description: "默认本地文件存储",
	}).Error
}

// StorageSourceRequest 存储源请求
type StorageSourceRequest struct {
	Name        string             `json:"name" binding:"required"`
	Code        string             `json:"code" binding:"required"`
	Type        models.MediaSource `json:"type" binding:"required"`
	Protocol    string             `json:"protocol"`
	Host        string             `json:"host" binding:"required"`
	Port        string             `json:"port"`
	BasePath    string             `json:"base_path"`
	Bucket      string             `json:"bucket"`
	AccessKey   string             `json:"access_key"`
	SecretKey   string             `json:"secret_key"`
	IsDefault   bool               `json:"is_default"`
	IsActive    bool               `json:"is_active"`
	Description string             `json:"description"`
}
