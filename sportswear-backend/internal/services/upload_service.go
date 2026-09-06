package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"sportswear-backend/internal/config"
	"sportswear-backend/internal/models"
)

// UploadService 文件上传服务（local 本地磁盘 / MinIO 对象存储，按 STORAGE_DRIVER 切换）
type UploadService struct {
	cfg      *config.Config
	driver   string // local / minio
	endpoint string
	useSSL   bool

	minioClient *minio.Client
	minioOnce   sync.Once
	minioErr    error
}

// NewUploadService 创建上传服务
func NewUploadService(cfg *config.Config) *UploadService {
	return &UploadService{
		cfg:      cfg,
		driver:   cfg.StorageDriver,
		endpoint: cfg.MinIO.Endpoint,
		useSSL:   cfg.MinIO.UseSSL,
	}
}

// Driver 返回当前存储驱动（local / minio）
func (s *UploadService) Driver() string {
	if s.driver == "" {
		return "local"
	}
	return s.driver
}

// MediaSource 返回当前驱动对应的媒体来源（local / minio），用于落库 source 字段与存储真实一致
func (s *UploadService) MediaSource() models.MediaSource {
	if s.Driver() == "minio" {
		return models.MediaSourceMinIO
	}
	return models.MediaSourceLocal
}

// getClient 懒加载 MinIO 客户端（仅 minio 驱动时使用，线程安全）
func (s *UploadService) getClient() (*minio.Client, error) {
	s.minioOnce.Do(func() {
		client, err := minio.New(s.endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(s.cfg.MinIO.AccessKey, s.cfg.MinIO.SecretKey, ""),
			Secure: s.useSSL,
		})
		if err != nil {
			s.minioErr = fmt.Errorf("初始化 MinIO 客户端失败: %w", err)
			return
		}
		s.minioClient = client
	})
	if s.minioErr != nil {
		return nil, s.minioErr
	}
	if s.minioClient == nil {
		return nil, fmt.Errorf("MinIO 客户端未初始化")
	}
	return s.minioClient, nil
}

// EnsureBucket 确保对象存储桶存在（minio 驱动下启动时调用，幂等）
func (s *UploadService) EnsureBucket() error {
	if s.Driver() != "minio" {
		return nil
	}
	client, err := s.getClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, s.cfg.MinIO.Bucket)
	if err != nil {
		return fmt.Errorf("检查 MinIO 桶失败: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, s.cfg.MinIO.Bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("创建 MinIO 桶失败: %w", err)
		}
	}
	return nil
}

// Save 统一保存入口：按存储驱动分发到本地磁盘或 MinIO
func (s *UploadService) Save(reader io.Reader, storagePath, contentType string, size int64) error {
	if s.Driver() == "minio" {
		return s.SaveToMinIO(reader, storagePath, contentType, size)
	}
	return s.SaveToLocal(reader, storagePath)
}

// SaveToMinIO 保存文件到 MinIO 对象存储
func (s *UploadService) SaveToMinIO(reader io.Reader, storagePath, contentType string, size int64) error {
	client, err := s.getClient()
	if err != nil {
		return err
	}
	if contentType == "" {
		contentType = GetFileType(storagePath)
	}
	_, err = client.PutObject(context.Background(), s.cfg.MinIO.Bucket, storagePath, reader, size,
		minio.PutObjectOptions{ContentType: contentType})
	return err
}

// URL 生成统一的展示访问 URL（相对路径 /uploads/<path>）
// local 与 minio 均返回相同相对路径，由后端 /uploads 读取路由按驱动分发，
// 保证 DB 中 url 字段跨存储、跨环境一致，形成「DB + 存储 + 展示」闭环。
func (s *UploadService) URL(storagePath string) string {
	return "/uploads/" + storagePath
}

// Delete 删除物理对象：minio 驱动下对象存储与本地磁盘均尝试删除（兼容存量 local 文件），
// local 驱动仅删本地磁盘；对象/文件不存在视为成功，空路径忽略。
func (s *UploadService) Delete(storagePath string) error {
	if storagePath == "" {
		return nil
	}
	var firstErr error
	if s.Driver() == "minio" {
		if client, err := s.getClient(); err == nil {
			// MinIO RemoveObject 对不存在的对象返回 nil，可安全调用
			if err := client.RemoveObject(context.Background(), s.cfg.MinIO.Bucket, storagePath,
				minio.RemoveObjectOptions{}); err != nil {
				firstErr = err
			}
		} else {
			firstErr = err
		}
	}
	// 本地磁盘一并尝试删除（双读兼容：存量 local 文件在 minio 模式下也能被清理）
	if err := os.Remove("uploads/" + storagePath); err != nil && !os.IsNotExist(err) {
		if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Open 按存储驱动打开文件流，返回 reader、文件大小、Content-Type（供 /uploads 读取路由）。
// minio 驱动下若对象不存在则回退本地磁盘，兼容存量 local 文件（双读）。
func (s *UploadService) Open(storagePath string) (io.ReadCloser, int64, string, error) {
	if s.Driver() == "minio" {
		if client, err := s.getClient(); err == nil {
			ctx := context.Background()
			if info, serr := client.StatObject(ctx, s.cfg.MinIO.Bucket, storagePath, minio.StatObjectOptions{}); serr == nil {
				if obj, gerr := client.GetObject(ctx, s.cfg.MinIO.Bucket, storagePath, minio.GetObjectOptions{}); gerr == nil {
					return obj, info.Size, info.ContentType, nil
				}
			}
		}
		// 回退本地磁盘（对象不存在 / 连接异常时仍保证可读）
	}
	return s.openLocal(storagePath)
}

// openLocal 从本地磁盘读取文件
func (s *UploadService) openLocal(storagePath string) (io.ReadCloser, int64, string, error) {
	f, err := os.Open("uploads/" + storagePath)
	if err != nil {
		return nil, 0, "", err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, "", err
	}
	return f, info.Size(), GetFileType(storagePath), nil
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

// leadAttachmentExtensions 询盘附件允许的扩展名（对齐行业联系表单：图片 + 压缩包 + 文档）
// 注意：不含 .svg（公开未登录上传有 XSS 风险），管理后台上传仍走 allowedExtensions
var leadAttachmentExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	".pdf": true, ".doc": true, ".docx": true,
	".xls": true, ".xlsx": true,
	".zip": true, ".rar": true,
}

// ValidateLeadAttachment 校验询盘附件安全：附件专用白名单（含 .rar）+ 大小限制（MAX_UPLOAD_SIZE，默认 20MB）
func (s *UploadService) ValidateLeadAttachment(fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !leadAttachmentExtensions[ext] {
		return fmt.Errorf("不支持的附件类型: %s（仅支持图片、PDF、Word、Excel、zip/rar 压缩包）", ext)
	}
	maxSize := s.cfg.Upload.MaxSize * 1024 * 1024 // MB -> bytes
	if fileHeader.Size > maxSize {
		return fmt.Errorf("附件超过大小限制 %dMB", s.cfg.Upload.MaxSize)
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
