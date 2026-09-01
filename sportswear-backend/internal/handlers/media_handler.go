package handlers

import (
	"strconv"

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

// DeleteMedia 删除媒体
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	if err := h.mediaService.DeleteMedia(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除媒体失败")
		return
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

	// 打开文件流并保存到本地
	src, err := fileHeader.Open()
	if err != nil {
		utils.InternalError(c, "读取文件失败")
		return
	}
	defer src.Close()

	if err := h.uploadService.SaveToLocal(src, storagePath); err != nil {
		utils.InternalError(c, "保存文件失败")
		return
	}

	// 创建媒体记录
	mediaType := services.DetectMediaType(fileHeader.Filename)
	fileType := services.DetectFileType(fileHeader.Filename)
	userID := middleware.GetUserID(c)

	media := models.Media{
		OriginalName: fileHeader.Filename,
		FileName:     fileName,
		FileType:     fileType,
		FileSize:     fileHeader.Size,
		URL:          h.uploadService.LocalURL(storagePath),
		Path:         storagePath,
		Category:     models.MediaCategory(category),
		Type:         mediaType,
		Source:       models.MediaSourceLocal,
		IsPublic:     true,
		UploadedBy:   &userID,
	}

	if err := h.mediaService.CreateMediaDirect(&media); err != nil {
		utils.InternalError(c, "创建媒体记录失败")
		return
	}

	utils.Created(c, media)
}
