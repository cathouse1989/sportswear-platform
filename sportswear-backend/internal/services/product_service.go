package services

import (
	"errors"

	"gorm.io/gorm"

	"sportswear-backend/internal/models"
	"sportswear-backend/internal/utils"
)

// ProductService 产品服务
type ProductService struct {
	db *gorm.DB
}

// NewProductService 创建产品服务
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// ListProducts 产品列表（增强版：支持多维度筛选）
func (s *ProductService) ListProducts(page, pageSize int, keyword, categoryID, status string) ([]models.Product, int64, error) {
	return s.listProducts(page, pageSize, keyword, categoryID, status, false, "")
}

// ListProductsV2 增强产品列表（支持 gender/type/material/series/fabric 等多维度筛选）
func (s *ProductService) ListProductsV2(params *ProductListParams) ([]models.Product, int64, error) {
	return s.listProductsV2(params)
}

// ProductListParams 增强产品列表查询参数
type ProductListParams struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	Keyword    string `json:"keyword"`
	CategoryID string `json:"category_id"`
	Status     string `json:"status"`
	Gender     string `json:"gender"`
	Type       string `json:"type"`
	Material   string `json:"material"`
	SeriesID   string `json:"series_id"`
	FabricID   string `json:"fabric_id"`
	IsFeatured *bool  `json:"is_featured"`
	IsNew      *bool  `json:"is_new"`
	SortBy     string `json:"sort_by"`    // sort_order, created_at, sku
	SortOrder  string `json:"sort_order"` // asc, desc
	// Language 语言维度动态排序：不同地区/语种市场有各自的特色运动与运动服装，
	// 门户按访问语言返回该语言下配置的展示顺序（翻译表 sort_order，0 回退全局）。
	Language string `json:"language"`
}

// ListFeaturedProducts 推荐产品列表（仅 is_featured=true，按语言动态排序）
func (s *ProductService) ListFeaturedProducts(limit int, lang string) ([]models.Product, error) {
	featured := true
	products, _, err := s.listProductsV2(&ProductListParams{
		Page:       1,
		PageSize:   limit,
		Status:     "published",
		IsFeatured: &featured,
		Language:   lang,
	})
	return products, err
}

// listProducts 内部通用查询（兼容旧版）
func (s *ProductService) listProducts(page, pageSize int, keyword, categoryID, status string, featuredOnly bool, gender string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{})
	if keyword != "" {
		query = query.Where("sku LIKE ? OR slug LIKE ? OR brief LIKE ? OR description LIKE ? OR material LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if gender != "" {
		query = query.Where("gender = ?", gender)
	}
	if featuredOnly {
		query = query.Where("is_featured = ?", true)
	}
	query.Count(&total)
	err := query.Preload("Category").Preload("Translations").Preload("Images").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("sort_order ASC, created_at DESC").Find(&products).Error

	return products, total, err
}

// listProductsV2 增强版查询
func (s *ProductService) listProductsV2(params *ProductListParams) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// 防御性参数校验：page/pageSize 至少为 1，避免 Offset 为负数
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}

	query := s.db.Model(&models.Product{})

	if params.Keyword != "" {
		query = query.Where("sku LIKE ? OR slug LIKE ? OR brief LIKE ? OR description LIKE ? OR material LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.CategoryID != "" {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Gender != "" {
		query = query.Where("gender = ?", params.Gender)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Material != "" {
		query = query.Where("material LIKE ?", "%"+params.Material+"%")
	}
	if params.IsFeatured != nil {
		query = query.Where("is_featured = ?", *params.IsFeatured)
	}
	if params.IsNew != nil {
		query = query.Where("is_new = ?", *params.IsNew)
	}

	// 系列关联查询
	if params.SeriesID != "" {
		query = query.Joins("JOIN product_series ON product_series.product_id = products.id").
			Where("product_series.series_id = ?", params.SeriesID)
	}
	// 面料关联查询
	if params.FabricID != "" {
		query = query.Joins("JOIN product_fabrics ON product_fabrics.product_id = products.id").
			Where("product_fabrics.fabric_id = ?", params.FabricID)
	}

	query.Count(&total)

	// 排序
	orderClause := "sort_order ASC, created_at DESC"
	if params.SortBy != "" {
		// 白名单校验，防止 ORDER BY 注入
		allowedSortBy := map[string]bool{
			"sort_order": true, "created_at": true, "sku": true,
			"sample_moq": true, "production_moq": true,
		}
		if allowedSortBy[params.SortBy] {
			dir := "ASC"
			if params.SortOrder == "desc" {
				dir = "DESC"
			}
			orderClause = params.SortBy + " " + dir
		} else {
			params.SortBy = ""
		}
	}

	db := query.Preload("Category").Preload("Translations").Preload("Images").
		Preload("Series").Preload("Fabrics").
		Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)

	// 语言维度动态排序：不同地区/语种市场有各自的特色运动与运动服装，
	// 优先取该语言翻译上配置的 sort_order（0/未配置回退全局 sort_order）。
	// 使用参数化子查询，防止语言值注入 ORDER BY。
	if params.Language != "" && params.SortBy == "" {
		db = db.Order(gorm.Expr(
			"COALESCE(NULLIF((SELECT pt.sort_order FROM product_translations pt "+
				"WHERE pt.product_id = products.id AND pt.language = ? AND pt.deleted_at IS NULL LIMIT 1), 0), "+
				"products.sort_order) ASC, products.sort_order ASC, products.created_at DESC",
			params.Language,
		))
	} else {
		db = db.Order(orderClause)
	}

	err := db.Find(&products).Error

	return products, total, err
}

// GetProduct 获取产品详情
func (s *ProductService) GetProduct(id string) (*models.Product, error) {
	var product models.Product
	err := s.db.Preload("Category").Preload("Translations").Preload("Images").
		Preload("Videos").Preload("Specs").Preload("Customizations").
		Preload("Series").Preload("Fabrics").
		First(&product, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("产品不存在")
	}
	return &product, nil
}

// GetProductBySlug 通过 Slug 获取产品
func (s *ProductService) GetProductBySlug(slug string) (*models.Product, error) {
	var product models.Product
	err := s.db.Preload("Category").Preload("Translations").Preload("Images").
		Preload("Videos").Preload("Specs").Preload("Customizations").
		Preload("Series").Preload("Fabrics").
		First(&product, "slug = ? AND status = ?", slug, models.ProductStatusPublished).Error
	if err != nil {
		return nil, errors.New("产品不存在")
	}
	return &product, nil
}

// CreateProduct 创建产品
func (s *ProductService) CreateProduct(req *ProductRequest) (*models.Product, error) {
	product := models.Product{
		SKU:           req.SKU,
		Slug:          req.Slug,
		CategoryID:    utils.StringPtrToUUIDPtr(req.CategoryID),
		Type:          req.Type,
		Gender:        req.Gender,
		Status:        req.Status,
		IsFeatured:    req.IsFeatured,
		IsNew:         req.IsNew,
		SortOrder:     req.SortOrder,
		CoverImage:    req.CoverImage,
		Brief:         req.Brief,
		Description:   req.Description,
		Features:      req.Features,
		Usage:         req.Usage,
		Material:      req.Material,
		Composition:   req.Composition,
		Weight:        req.Weight,
		Elasticity:    req.Elasticity,
		Fit:           req.Fit,
		SupportLevel:  req.SupportLevel,
		Season:        req.Season,
		SizeRange:     req.SizeRange,
		SampleMOQ:     req.SampleMOQ,
		ProductionMOQ: req.ProductionMOQ,
		ColorMOQ:      req.ColorMOQ,
		SizeMOQ:       req.SizeMOQ,
	}

	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}

	// 创建翻译
	if len(req.Translations) > 0 {
		for _, t := range req.Translations {
			translation := models.ProductTranslation{
				ProductID:    product.ID,
				Language:     t.Language,
				Name:         t.Name,
				Brief:        t.Brief,
				Description:  t.Description,
				Features:     t.Features,
				Usage:        t.Usage,
				Material:     t.Material,
				Composition:  t.Composition,
				Weight:       t.Weight,
				Elasticity:   t.Elasticity,
				Fit:          t.Fit,
				SupportLevel: t.SupportLevel,
				Season:       t.Season,
				SizeRange:    t.SizeRange,
				SortOrder:    t.SortOrder,
				Status:       models.TranslationStatusPublished,
			}
			s.db.Create(&translation)
		}
	}

	// 创建 SEO
	if req.SEO != nil {
		seo := models.SEO{
			EntityType:    "product",
			EntityID:      product.ID,
			Language:      req.SEO.Language,
			Title:         req.SEO.Title,
			Description:   req.SEO.Description,
			Keywords:      req.SEO.Keywords,
			Canonical:     req.SEO.Canonical,
			Robots:        req.SEO.Robots,
			OGTitle:       req.SEO.OGTitle,
			OGDescription: req.SEO.OGDescription,
			OGImage:       req.SEO.OGImage,
			SchemaData:    req.SEO.SchemaData,
		}
		s.db.Create(&seo)
	}

	// 关联系列
	if len(req.SeriesIDs) > 0 {
		for _, sid := range req.SeriesIDs {
			s.db.Exec("INSERT INTO product_series (product_id, series_id) VALUES (?, ?) ON CONFLICT DO NOTHING", product.ID, sid)
		}
	}

	// 关联面料
	if len(req.FabricIDs) > 0 {
		for _, fid := range req.FabricIDs {
			s.db.Exec("INSERT INTO product_fabrics (product_id, fabric_id) VALUES (?, ?) ON CONFLICT DO NOTHING", product.ID, fid)
		}
	}

	// 同步子资源（图片/视频/规格/定制能力）
	if err := s.syncProductChildren(&product, req); err != nil {
		return nil, err
	}

	return s.GetProduct(product.ID.String())
}

// syncProductChildren 整体替换产品的子资源（图片/视频/规格/定制能力）。
// 只有请求中「显式携带」的数组才会被处理（nil 表示保持不变）：
//   - 物理删除旧的关联记录，再按前端提交的完整列表重建，
//     保证「整体 PUT」语义：前端提交即为最终结果。
func (s *ProductService) syncProductChildren(product *models.Product, req *ProductRequest) error {
	// ---------- 图片 ----------
	if req.Images != nil {
		if err := s.db.Unscoped().Where("product_id = ?", product.ID).Delete(&models.ProductImage{}).Error; err != nil {
			return err
		}
		for _, img := range req.Images {
			item := models.ProductImage{
				ProductID: product.ID,
				Type:      img.Type,
				URL:       img.URL,
				Thumbnail: img.Thumbnail,
				Alt:       img.Alt,
				SortOrder: img.SortOrder,
			}
			if item.Type == "" {
				item.Type = "gallery"
			}
			if err := s.db.Create(&item).Error; err != nil {
				return err
			}
		}
	}

	// ---------- 视频 ----------
	if req.Videos != nil {
		if err := s.db.Unscoped().Where("product_id = ?", product.ID).Delete(&models.ProductVideo{}).Error; err != nil {
			return err
		}
		for _, video := range req.Videos {
			item := models.ProductVideo{
				ProductID: product.ID,
				Type:      video.Type,
				URL:       video.URL,
				Cover:     video.Cover,
				Title:     video.Title,
				SortOrder: video.SortOrder,
			}
			if item.Type == "" {
				item.Type = "product"
			}
			if err := s.db.Create(&item).Error; err != nil {
				return err
			}
		}
	}

	// ---------- 规格 ----------
	if req.Specs != nil {
		if err := s.db.Unscoped().Where("product_id = ?", product.ID).Delete(&models.ProductSpec{}).Error; err != nil {
			return err
		}
		for i, spec := range req.Specs {
			item := models.ProductSpec{
				ProductID: product.ID,
				Name:      spec.Name,
				Value:     spec.Value,
				SortOrder: spec.SortOrder,
			}
			if spec.SortOrder == 0 {
				item.SortOrder = i
			}
			if err := s.db.Create(&item).Error; err != nil {
				return err
			}
		}
	}

	// ---------- 定制能力 ----------
	if req.Customizations != nil {
		if err := s.db.Unscoped().Where("product_id = ?", product.ID).Delete(&models.ProductCustomization{}).Error; err != nil {
			return err
		}
		for _, cus := range req.Customizations {
			item := models.ProductCustomization{
				ProductID:    product.ID,
				Type:         cus.Type,
				IsEnabled:    cus.IsEnabled,
				Note:         cus.Note,
				Translations: cus.Translations,
			}
			if err := s.db.Create(&item).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// ProductRequest 产品请求（完整版）
type ProductRequest struct {
	SKU           string                      `json:"sku" binding:"required"`
	Slug          string                      `json:"slug" binding:"required"`
	CategoryID    *string                     `json:"category_id"`
	Type          models.ProductType          `json:"type"`
	Gender        models.Gender               `json:"gender"`
	Status        models.ProductStatus        `json:"status"`
	IsFeatured    bool                        `json:"is_featured"`
	IsNew         bool                        `json:"is_new"`
	SortOrder     int                         `json:"sort_order"`
	CoverImage    string                      `json:"cover_image"`
	Brief         string                      `json:"brief"`
	Description   string                      `json:"description"`
	Features      string                      `json:"features"`
	Usage         string                      `json:"usage"`
	Material      string                      `json:"material"`
	Composition   string                      `json:"composition"`
	Weight        string                      `json:"weight"`
	Elasticity    string                      `json:"elasticity"`
	Fit           string                      `json:"fit"`
	SupportLevel  string                      `json:"support_level"`
	Season        string                      `json:"season"`
	SizeRange     string                      `json:"size_range"`
	SampleMOQ     int                         `json:"sample_moq"`
	ProductionMOQ int                         `json:"production_moq"`
	ColorMOQ      int                         `json:"color_moq"`
	SizeMOQ       int                         `json:"size_moq"`
	Translations  []ProductTranslationRequest `json:"translations"`
	SEO           *SEORequest                 `json:"seo"`
	SeriesIDs     []string                    `json:"series_ids"`
	FabricIDs     []string                    `json:"fabric_ids"`
	// 子资源（整体 PUT 时一并提交，不传表示保持不变，传 null/[] 表示清空）
	Images         []ProductImageRequest         `json:"images"`
	Videos         []ProductVideoRequest         `json:"videos"`
	Specs          []ProductSpecRequest          `json:"specs"`
	Customizations []ProductCustomizationRequest `json:"customizations"`
}

// ProductImageRequest 产品图片请求
type ProductImageRequest struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Alt       string `json:"alt"`
	SortOrder int    `json:"sort_order"`
}

// ProductVideoRequest 产品视频请求
type ProductVideoRequest struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	Cover     string `json:"cover"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
}

// ProductSpecRequest 产品规格请求
type ProductSpecRequest struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	SortOrder int    `json:"sort_order"`
}

// ProductCustomizationRequest 产品定制能力请求
type ProductCustomizationRequest struct {
	Type         string `json:"type"`
	IsEnabled    bool   `json:"is_enabled"`
	Note         string `json:"note"`
	Translations string `json:"translations"`
}

// ProductTranslationRequest 产品翻译请求
type ProductTranslationRequest struct {
	Language    string `json:"language" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Brief       string `json:"brief"`
	Description string `json:"description"`
	Features    string `json:"features"`
	Usage       string `json:"usage"`
	// 规格类字段翻译（主表存英文源，此处按语言覆盖）
	Material     string `json:"material"`
	Composition  string `json:"composition"`
	Weight       string `json:"weight"`
	Elasticity   string `json:"elasticity"`
	Fit          string `json:"fit"`
	SupportLevel string `json:"support_level"`
	Season       string `json:"season"`
	SizeRange    string `json:"size_range"`
	// SortOrder 语言维度的展示排序权重（0 = 跟随全局排序）
	SortOrder int `json:"sort_order"`
}

// SEORequest SEO 请求
type SEORequest struct {
	Language      string `json:"language"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Keywords      string `json:"keywords"`
	Canonical     string `json:"canonical"`
	Robots        string `json:"robots"`
	OGTitle       string `json:"og_title"`
	OGDescription string `json:"og_description"`
	OGImage       string `json:"og_image"`
	SchemaData    string `json:"schema_data"`
}

// UpdateProduct 更新产品
func (s *ProductService) UpdateProduct(id string, req *ProductRequest) (*models.Product, error) {
	var product models.Product
	if err := s.db.First(&product, "id = ?", id).Error; err != nil {
		return nil, errors.New("产品不存在")
	}

	updates := map[string]interface{}{
		"sku":            req.SKU,
		"slug":           req.Slug,
		"category_id":    req.CategoryID,
		"type":           req.Type,
		"gender":         req.Gender,
		"status":         req.Status,
		"is_featured":    req.IsFeatured,
		"is_new":         req.IsNew,
		"sort_order":     req.SortOrder,
		"cover_image":    req.CoverImage,
		"brief":          req.Brief,
		"description":    req.Description,
		"features":       req.Features,
		"usage":          req.Usage,
		"material":       req.Material,
		"composition":    req.Composition,
		"weight":         req.Weight,
		"elasticity":     req.Elasticity,
		"fit":            req.Fit,
		"support_level":  req.SupportLevel,
		"season":         req.Season,
		"size_range":     req.SizeRange,
		"sample_moq":     req.SampleMOQ,
		"production_moq": req.ProductionMOQ,
		"color_moq":      req.ColorMOQ,
		"size_moq":       req.SizeMOQ,
	}

	if err := s.db.Model(&product).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 更新翻译
	if len(req.Translations) > 0 {
		for _, t := range req.Translations {
			var translation models.ProductTranslation
			err := s.db.Where("product_id = ? AND language = ?", product.ID, t.Language).First(&translation).Error
			if err != nil {
				s.db.Create(&models.ProductTranslation{
					ProductID:    product.ID,
					Language:     t.Language,
					Name:         t.Name,
					Brief:        t.Brief,
					Description:  t.Description,
					Features:     t.Features,
					Usage:        t.Usage,
					Material:     t.Material,
					Composition:  t.Composition,
					Weight:       t.Weight,
					Elasticity:   t.Elasticity,
					Fit:          t.Fit,
					SupportLevel: t.SupportLevel,
					Season:       t.Season,
					SizeRange:    t.SizeRange,
					SortOrder:    t.SortOrder,
					Status:       models.TranslationStatusPublished,
				})
			} else {
				s.db.Model(&translation).Updates(map[string]interface{}{
					"name":          t.Name,
					"brief":         t.Brief,
					"description":   t.Description,
					"features":      t.Features,
					"usage":         t.Usage,
					"material":      t.Material,
					"composition":   t.Composition,
					"weight":        t.Weight,
					"elasticity":    t.Elasticity,
					"fit":           t.Fit,
					"support_level": t.SupportLevel,
					"season":        t.Season,
					"size_range":    t.SizeRange,
					"sort_order":    t.SortOrder,
					"status":        models.TranslationStatusPublished,
				})
			}
		}
	}

	// 更新 SEO
	if req.SEO != nil {
		var seo models.SEO
		err := s.db.Where("entity_type = ? AND entity_id = ?", "product", product.ID).First(&seo).Error
		if err != nil {
			s.db.Create(&models.SEO{
				EntityType:    "product",
				EntityID:      product.ID,
				Language:      req.SEO.Language,
				Title:         req.SEO.Title,
				Description:   req.SEO.Description,
				Keywords:      req.SEO.Keywords,
				Canonical:     req.SEO.Canonical,
				Robots:        req.SEO.Robots,
				OGTitle:       req.SEO.OGTitle,
				OGDescription: req.SEO.OGDescription,
				OGImage:       req.SEO.OGImage,
				SchemaData:    req.SEO.SchemaData,
			})
		} else {
			s.db.Model(&seo).Updates(map[string]interface{}{
				"language":       req.SEO.Language,
				"title":          req.SEO.Title,
				"description":    req.SEO.Description,
				"keywords":       req.SEO.Keywords,
				"canonical":      req.SEO.Canonical,
				"robots":         req.SEO.Robots,
				"og_title":       req.SEO.OGTitle,
				"og_description": req.SEO.OGDescription,
				"og_image":       req.SEO.OGImage,
				"schema_data":    req.SEO.SchemaData,
			})
		}
	}

	// 更新系列关联
	if req.SeriesIDs != nil {
		s.db.Exec("DELETE FROM product_series WHERE product_id = ?", product.ID)
		for _, sid := range req.SeriesIDs {
			s.db.Exec("INSERT INTO product_series (product_id, series_id) VALUES (?, ?) ON CONFLICT DO NOTHING", product.ID, sid)
		}
	}

	// 更新面料关联
	if req.FabricIDs != nil {
		s.db.Exec("DELETE FROM product_fabrics WHERE product_id = ?", product.ID)
		for _, fid := range req.FabricIDs {
			s.db.Exec("INSERT INTO product_fabrics (product_id, fabric_id) VALUES (?, ?) ON CONFLICT DO NOTHING", product.ID, fid)
		}
	}

	// 同步子资源（图片/视频/规格/定制能力）
	if err := s.syncProductChildren(&product, req); err != nil {
		return nil, err
	}

	return s.GetProduct(id)
}

// DeleteProduct 删除产品（级联软删除关联数据）
func (s *ProductService) DeleteProduct(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&models.ProductTranslation{}).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&models.ProductImage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&models.ProductVideo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&models.ProductSpec{}).Error; err != nil {
			return err
		}
		if err := tx.Where("product_id = ?", id).Delete(&models.ProductCustomization{}).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM product_series WHERE product_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM product_fabrics WHERE product_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Product{}, "id = ?", id).Error
	})
}

// PublishProduct 发布产品
func (s *ProductService) PublishProduct(id string) error {
	return s.db.Model(&models.Product{}).Where("id = ?", id).Update("status", models.ProductStatusPublished).Error
}

// ProductStats 产品状态统计（全量口径，不受列表分页影响）
type ProductStats struct {
	Total     int64 `json:"total"`
	Published int64 `json:"published"`
	Draft     int64 `json:"draft"`
	Offline   int64 `json:"offline"`
	Featured  int64 `json:"featured"`
}

// GetProductStats 统计各状态产品数量（供管理后台统计卡使用）
func (s *ProductService) GetProductStats() (*ProductStats, error) {
	type statusCount struct {
		Status models.ProductStatus `json:"status"`
		Count  int64                `json:"count"`
	}
	var counts []statusCount
	if err := s.db.Model(&models.Product{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	stats := &ProductStats{}
	for _, sc := range counts {
		switch sc.Status {
		case models.ProductStatusPublished:
			stats.Published = sc.Count
		case models.ProductStatusDraft:
			stats.Draft = sc.Count
		case models.ProductStatusOffline:
			stats.Offline = sc.Count
		}
		stats.Total += sc.Count
	}

	if err := s.db.Model(&models.Product{}).
		Where("is_featured = ?", true).Count(&stats.Featured).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

// UnpublishProduct 下架产品
func (s *ProductService) UnpublishProduct(id string) error {
	return s.db.Model(&models.Product{}).Where("id = ?", id).Update("status", models.ProductStatusOffline).Error
}

// ListCategories 分类列表
func (s *ProductService) ListCategories() ([]models.Category, error) {
	var categories []models.Category
	err := s.db.Preload("Children").Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

// CreateCategory 创建分类
func (s *ProductService) CreateCategory(req *CategoryRequest) (*models.Category, error) {
	category := models.Category{
		ParentID:  utils.StringPtrToUUIDPtr(req.ParentID),
		Name:      req.Name,
		Slug:      req.Slug,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
		Image:     req.Image,
	}
	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// CategoryRequest 分类请求
type CategoryRequest struct {
	ParentID  *string `json:"parent_id"`
	Name      string  `json:"name" binding:"required"`
	Slug      string  `json:"slug" binding:"required"`
	SortOrder int     `json:"sort_order"`
	IsActive  bool    `json:"is_active"`
	Image     string  `json:"image"`
}

// UpdateCategory 更新分类
func (s *ProductService) UpdateCategory(id string, req *CategoryRequest) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, "id = ?", id).Error; err != nil {
		return nil, errors.New("分类不存在")
	}

	updates := map[string]interface{}{
		"parent_id":  req.ParentID,
		"name":       req.Name,
		"slug":       req.Slug,
		"sort_order": req.SortOrder,
		"is_active":  req.IsActive,
		"image":      req.Image,
	}
	if err := s.db.Model(&category).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// DeleteCategory 删除分类（检查子分类与产品引用）
func (s *ProductService) DeleteCategory(id string) error {
	var childCount int64
	s.db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return errors.New("该分类下存在子分类，无法删除")
	}

	var productCount int64
	s.db.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
	if productCount > 0 {
		return errors.New("该分类下存在产品，无法删除")
	}

	return s.db.Delete(&models.Category{}, "id = ?", id).Error
}

// ListSeries 系列列表（后台：全部状态）
func (s *ProductService) ListSeries() ([]models.Series, error) {
	var series []models.Series
	err := s.db.Order("sort_order ASC").Find(&series).Error
	return series, err
}

// ListPublishedSeries 已发布系列列表（前台）
func (s *ProductService) ListPublishedSeries() ([]models.Series, error) {
	var series []models.Series
	err := s.db.Where("status = ?", models.ProductStatusPublished).
		Order("sort_order ASC").Find(&series).Error
	return series, err
}

// CreateSeries 创建系列
func (s *ProductService) CreateSeries(req *SeriesRequest) (*models.Series, error) {
	series := models.Series{
		Name:      req.Name,
		Slug:      req.Slug,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		IsActive:  req.IsActive,
	}
	if err := s.db.Create(&series).Error; err != nil {
		return nil, err
	}
	return &series, nil
}

// PublishSeries 发布系列（前端可见）
func (s *ProductService) PublishSeries(id string) error {
	return s.db.Model(&models.Series{}).Where("id = ?", id).
		Update("status", models.ProductStatusPublished).Error
}

// UnpublishSeries 下线系列（前端不可见）
func (s *ProductService) UnpublishSeries(id string) error {
	return s.db.Model(&models.Series{}).Where("id = ?", id).
		Update("status", models.ProductStatusOffline).Error
}

// SeriesRequest 系列请求
type SeriesRequest struct {
	Name      string               `json:"name" binding:"required"`
	Slug      string               `json:"slug" binding:"required"`
	SortOrder int                  `json:"sort_order"`
	Status    models.ProductStatus `json:"status"`
	IsActive  bool                 `json:"is_active"`
}

// ListFabrics 面料列表（后台：全部状态）
func (s *ProductService) ListFabrics() ([]models.Fabric, error) {
	var fabrics []models.Fabric
	err := s.db.Order("name ASC").Find(&fabrics).Error
	return fabrics, err
}

// ListPublishedFabrics 已发布面料列表（前台）
func (s *ProductService) ListPublishedFabrics() ([]models.Fabric, error) {
	var fabrics []models.Fabric
	err := s.db.Where("status = ?", models.ProductStatusPublished).
		Order("name ASC").Find(&fabrics).Error
	return fabrics, err
}

// CreateFabric 创建面料
func (s *ProductService) CreateFabric(req *FabricRequest) (*models.Fabric, error) {
	fabric := models.Fabric{
		Name:            req.Name,
		Code:            req.Code,
		Status:          req.Status,
		Composition:     req.Composition,
		Weight:          req.Weight,
		Elasticity:      req.Elasticity,
		Breathability:   req.Breathability,
		MoistureWicking: req.MoistureWicking,
		Softness:        req.Softness,
		Compression:     req.Compression,
		UVProtection:    req.UVProtection,
		EcoFriendly:     req.EcoFriendly,
		Description:     req.Description,
		IsActive:        req.IsActive,
	}
	if err := s.db.Create(&fabric).Error; err != nil {
		return nil, err
	}
	return &fabric, nil
}

// FabricRequest 面料请求
type FabricRequest struct {
	Name            string               `json:"name" binding:"required"`
	Code            string               `json:"code" binding:"required"`
	Composition     string               `json:"composition"`
	Weight          string               `json:"weight"`
	Elasticity      string               `json:"elasticity"`
	Breathability   string               `json:"breathability"`
	MoistureWicking string               `json:"moisture_wicking"`
	Softness        string               `json:"softness"`
	Compression     string               `json:"compression"`
	UVProtection    string               `json:"uv_protection"`
	EcoFriendly     string               `json:"eco_friendly"`
	Description     string               `json:"description"`
	Status          models.ProductStatus `json:"status"`
	IsActive        bool                 `json:"is_active"`
}

// PublishFabric 发布面料（前端可见）
func (s *ProductService) PublishFabric(id string) error {
	return s.db.Model(&models.Fabric{}).Where("id = ?", id).
		Update("status", models.ProductStatusPublished).Error
}

// UnpublishFabric 下线面料（前端不可见）
func (s *ProductService) UnpublishFabric(id string) error {
	return s.db.Model(&models.Fabric{}).Where("id = ?", id).
		Update("status", models.ProductStatusOffline).Error
}
