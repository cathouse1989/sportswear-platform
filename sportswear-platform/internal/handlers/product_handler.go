package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-platform/internal/services"
	"sportswear-platform/internal/utils"
)

// ProductHandler 产品处理器
type ProductHandler struct {
	productService *services.ProductService
	cache          *services.CacheService
}

// NewProductHandler 创建产品处理器
func NewProductHandler(productService *services.ProductService, cacheService *services.CacheService) *ProductHandler {
	return &ProductHandler{productService: productService, cache: cacheService}
}

// invalidateCache 后台写操作后主动失效对应门户缓存（缓存不可用时无副作用）
func (h *ProductHandler) invalidateCache(entityType, slug string) {
	if h.cache != nil {
		h.cache.InvalidateContentCache(entityType, slug)
	}
}

// ListProducts 产品列表（增强版：支持多维度筛选）
func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 使用增强版查询
	params := &services.ProductListParams{
		Page:       page,
		PageSize:   pageSize,
		Keyword:    c.Query("keyword"),
		CategoryID: c.Query("category_id"),
		Status:     c.Query("status"),
		Gender:     c.Query("gender"),
		Type:       c.Query("type"),
		Material:   c.Query("material"),
		SeriesID:   c.Query("series_id"),
		FabricID:   c.Query("fabric_id"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
	}

	// 解析布尔参数
	if v := c.Query("is_featured"); v != "" {
		b := v == "true"
		params.IsFeatured = &b
	}
	if v := c.Query("is_new"); v != "" {
		b := v == "true"
		params.IsNew = &b
	}

	products, total, err := h.productService.ListProductsV2(params)
	if err != nil {
		utils.InternalError(c, "获取产品列表失败")
		return
	}
	utils.SuccessPage(c, products, page, pageSize, total)
}

// GetProductStats 产品状态统计（全量口径，供统计卡使用）
func (h *ProductHandler) GetProductStats(c *gin.Context) {
	stats, err := h.productService.GetProductStats()
	if err != nil {
		utils.InternalError(c, "获取产品统计失败")
		return
	}
	utils.Success(c, stats)
}

// GetProduct 获取产品
func (h *ProductHandler) GetProduct(c *gin.Context) {
	product, err := h.productService.GetProduct(c.Param("id"))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.Success(c, product)
}

// CreateProduct 创建产品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req services.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("product", "")
	utils.Created(c, product)
}

// UpdateProduct 更新产品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req services.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	product, err := h.productService.UpdateProduct(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("product", "")
	utils.Success(c, product)
}

// DeleteProduct 删除产品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	if err := h.productService.DeleteProduct(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除产品失败")
		return
	}
	h.invalidateCache("product", "")
	utils.Success(c, gin.H{"deleted": true})
}

// PublishProduct 发布产品
func (h *ProductHandler) PublishProduct(c *gin.Context) {
	if err := h.productService.PublishProduct(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布产品失败")
		return
	}
	h.invalidateCache("product", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishProduct 下架产品
func (h *ProductHandler) UnpublishProduct(c *gin.Context) {
	if err := h.productService.UnpublishProduct(c.Param("id")); err != nil {
		utils.BadRequest(c, "下架产品失败")
		return
	}
	h.invalidateCache("product", "")
	utils.Success(c, gin.H{"unpublished": true})
}

// ListCategories 分类列表
func (h *ProductHandler) ListCategories(c *gin.Context) {
	categories, err := h.productService.ListCategories()
	if err != nil {
		utils.InternalError(c, "获取分类列表失败")
		return
	}
	utils.Success(c, categories)
}

// CreateCategory 创建分类
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req services.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	category, err := h.productService.CreateCategory(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("category", "")
	utils.Created(c, category)
}

// UpdateCategory 更新分类
func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	var req services.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	category, err := h.productService.UpdateCategory(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("category", "")
	utils.Success(c, category)
}

// DeleteCategory 删除分类
func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	if err := h.productService.DeleteCategory(c.Param("id")); err != nil {
		utils.BadRequest(c, "删除分类失败")
		return
	}
	h.invalidateCache("category", "")
	utils.Success(c, gin.H{"deleted": true})
}

// ListSeries 系列列表
func (h *ProductHandler) ListSeries(c *gin.Context) {
	series, err := h.productService.ListSeries()
	if err != nil {
		utils.InternalError(c, "获取系列列表失败")
		return
	}
	utils.Success(c, series)
}

// CreateSeries 创建系列
func (h *ProductHandler) CreateSeries(c *gin.Context) {
	var req services.SeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	series, err := h.productService.CreateSeries(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("series", "")
	utils.Created(c, series)
}

// PublishSeries 发布系列（前端可见）
func (h *ProductHandler) PublishSeries(c *gin.Context) {
	if err := h.productService.PublishSeries(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布系列失败")
		return
	}
	h.invalidateCache("series", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishSeries 下线系列（前端不可见）
func (h *ProductHandler) UnpublishSeries(c *gin.Context) {
	if err := h.productService.UnpublishSeries(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线系列失败")
		return
	}
	h.invalidateCache("series", "")
	utils.Success(c, gin.H{"unpublished": true})
}

// ListFabrics 面料列表
func (h *ProductHandler) ListFabrics(c *gin.Context) {
	fabrics, err := h.productService.ListFabrics()
	if err != nil {
		utils.InternalError(c, "获取面料列表失败")
		return
	}
	utils.Success(c, fabrics)
}

// CreateFabric 创建面料
func (h *ProductHandler) CreateFabric(c *gin.Context) {
	var req services.FabricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	fabric, err := h.productService.CreateFabric(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	h.invalidateCache("fabric", "")
	utils.Created(c, fabric)
}

// PublishFabric 发布面料（前端可见）
func (h *ProductHandler) PublishFabric(c *gin.Context) {
	if err := h.productService.PublishFabric(c.Param("id")); err != nil {
		utils.BadRequest(c, "发布面料失败")
		return
	}
	h.invalidateCache("fabric", "")
	utils.Success(c, gin.H{"published": true})
}

// UnpublishFabric 下线面料（前端不可见）
func (h *ProductHandler) UnpublishFabric(c *gin.Context) {
	if err := h.productService.UnpublishFabric(c.Param("id")); err != nil {
		utils.BadRequest(c, "下线面料失败")
		return
	}
	h.invalidateCache("fabric", "")
	utils.Success(c, gin.H{"unpublished": true})
}
