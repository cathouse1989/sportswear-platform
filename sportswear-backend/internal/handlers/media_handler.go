package handlers

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/models"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// MediaHandler 媒体处理器
type MediaHandler struct {
	mediaService  *services.MediaService
	uploadService *services.UploadService
}

// NewMediaHandler 创建媒体处理器
func NewMediaHandler(mediaService *services.MediaService, uploadService *services.UploadService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService, uploadService: uploadService}
}

// ListMedia 媒体列表
func (h *MediaHandler) ListMedia(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	category := c.Query("category")
	mediaType := c.Query("type")
	keyword := c.Query("keyword")

	media, total, err := h.mediaService.ListMedia(page, pageSize, category, mediaType, keyword)
	if err != nil {
		utils.InternalError(c, "获取媒体列表失败")
		return
	}
	utils.SuccessPage(c, media, page, pageSize, total)
}

// GetMedia 获取媒体
func (h *MediaHandler) GetMedia(c *gin.Context) {
	media, err := h.mediaService.GetMedia(c.Param("id"))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, media)
}

// CreateMedia 创建媒体记录
func (h *MediaHandler) CreateMedia(c *gin.Context) {
	var req services.MediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	media, err := h.mediaService.CreateMedia(&req, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, media)
}

// UpdateMedia 更新媒体
func (h *MediaHandler) UpdateMedia(c *gin.Context) {
	var req services.UpdateMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	media, err := h.mediaService.UpdateMedia(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, media)
}

// DeleteMedia 删除媒体（DB 记录 + 物理对象，local 删磁盘文件 / minio 删对象 / external 仅删记录）
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	id := c.Param("id")

	// 先取记录拿到存储路径，用于同步删除物理对象
	media, err := h.mediaService.GetMedia(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	// 引用检查：防止删除被产品封面/图集引用的图片导致死链
	if media.Source != models.MediaSourceExternal && media.Path != "" {
		if refCount, rerr := h.mediaService.CountProductReferences(media.Path); rerr == nil && refCount > 0 {
			utils.BadRequest(c, fmt.Sprintf("该媒体被 %d 处产品封面/图集引用，请先解除引用再删除", refCount))
			return
		}
	}

	if err := h.mediaService.DeleteMedia(id); err != nil {
		utils.BadRequest(c, "删除媒体失败")
		return
	}

	// 同步删除物理对象（external 来源仅删元数据；删除失败不阻断，仅告警）
	if media.Source != models.MediaSourceExternal {
		if err := h.uploadService.Delete(media.Path); err != nil {
			utils.Logger.Warnw("删除媒体物理对象失败", "id", id, "path", media.Path, "error", err.Error())
		}
		// 同步删除缩略图（若存在）
		if media.Type == models.MediaTypeImage && media.Path != "" {
			if err := h.uploadService.Delete(services.ThumbPath(media.Path)); err != nil {
				utils.Logger.Warnw("删除媒体缩略图失败", "id", id, "path", services.ThumbPath(media.Path), "error", err.Error())
			}
		}
	}

	utils.Success(c, gin.H{"deleted": true})
}

// UploadFile 上传媒体文件（multipart/form-data，字段 file）
// 将文件保存到本地 uploads 目录，并创建媒体记录
func (h *MediaHandler) UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 校验文件安全（扩展名白名单 + 大小限制）
	if err := h.uploadService.ValidateFile(fileHeader); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 分类：multipart 字段优先，其次 query 参数（便于 <el-upload :action> 带 category），默认 public
	category := c.PostForm("category")
	if category == "" {
		category = c.Query("category")
	}
	if category == "" {
		category = string(models.MediaCategoryPublic)
	}

	// 生成存储路径
	storagePath, fileName := h.uploadService.GenerateStoragePath(fileHeader.Filename, category)

	// 检测媒体类型与 MIME（供落库与 MinIO Content-Type 使用）
	mediaType := services.DetectMediaType(fileHeader.Filename)
	fileType := services.DetectFileType(fileHeader.Filename)

	// 图片解析宽高（数据完整性：JPEG/PNG/GIF/WebP，AVIF 暂跳过）
	var width, height int
	if mediaType == models.MediaTypeImage {
		if w, h, ok := services.ProbeImageDimensions(fileHeader); ok {
			width, height = w, h
		}
	}

	// 打开文件流并按存储驱动保存（local 磁盘 / minio 对象）
	src, err := fileHeader.Open()
	if err != nil {
		utils.InternalError(c, "读取文件失败")
		return
	}
	defer src.Close()

	if err := h.uploadService.Save(src, storagePath, fileType, fileHeader.Size); err != nil {
		utils.InternalError(c, "保存文件失败")
		return
	}

	// 图片生成缩略图（JPEG，最长边 320px），与原图同 driver 存储，纳入一致性闭环
	var thumbnail string
	if mediaType == models.MediaTypeImage {
		if tsrc, terr := fileHeader.Open(); terr == nil {
			if thumbBytes, gerr := services.GenerateThumbnail(tsrc, services.ThumbnailMaxWidth); gerr == nil {
				thumbPath := services.ThumbPath(storagePath)
				if serr := h.uploadService.Save(bytes.NewReader(thumbBytes), thumbPath, "image/jpeg", int64(len(thumbBytes))); serr == nil {
					thumbnail = h.uploadService.URL(thumbPath)
				}
			}
			tsrc.Close()
		}
	}

	// 创建媒体记录（source 由存储驱动决定，与真实存储一致）
	userID := middleware.GetUserID(c)

	media := models.Media{
		OriginalName: fileHeader.Filename,
		FileName:     fileName,
		FileType:     fileType,
		FileSize:     fileHeader.Size,
		Width:        width,
		Height:       height,
		URL:          h.uploadService.URL(storagePath),
		Thumbnail:    thumbnail,
		Path:         storagePath,
		Category:     models.MediaCategory(category),
		Type:         mediaType,
		Source:       h.uploadService.MediaSource(),
		IsPublic:     true,
		UploadedBy:   &userID,
	}

	if err := h.mediaService.CreateMediaDirect(&media); err != nil {
		utils.InternalError(c, "创建媒体记录失败")
		return
	}

	utils.Created(c, media)
}

// ServeUpload 提供上传文件的读取（/uploads/** 路由）
// 按存储驱动从本地磁盘或 MinIO 流式返回，展示层（后台/门户）统一走相对路径 /uploads/<path>，
// 与上传时落库的 url 字段一致，形成「DB + 存储 + 展示」闭环。
func (h *MediaHandler) ServeUpload(c *gin.Context) {
	path := strings.TrimPrefix(c.Param("filepath"), "/")
	if path == "" || strings.Contains(path, "..") {
		utils.NotFound(c, "文件不存在")
		return
	}

	reader, size, contentType, err := h.uploadService.Open(path)
	if err != nil {
		utils.NotFound(c, "文件不存在")
		return
	}
	defer reader.Close()

	if contentType != "" {
		c.Header("Content-Type", contentType)
	}
	c.Header("Content-Length", strconv.FormatInt(size, 10))
	c.Header("Cache-Control", "public, max-age=86400")

	if _, err := io.Copy(c.Writer, reader); err != nil {
		utils.Logger.Warnw("上传文件读取失败", "path", path, "error", err.Error())
	}
}
