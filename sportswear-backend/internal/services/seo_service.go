package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// SEOService SEO 国际化服务
// 提供多语言 SEO 配置、hreflang 标注、多语言 sitemap 生成
type SEOService struct {
	db *gorm.DB
}

// NewSEOService 创建 SEO 服务
func NewSEOService(db *gorm.DB) *SEOService {
	return &SEOService{db: db}
}

// GetSEO 获取指定实体的 SEO 配置（按语言，无则回退英文）
func (s *SEOService) GetSEO(entityType string, entityID uuid.UUID, lang string) (*models.SEO, error) {
	var seo models.SEO
	// 先查目标语言
	err := s.db.Where("entity_type = ? AND entity_id = ? AND language = ?", entityType, entityID, lang).
		First(&seo).Error
	if err != nil && lang != "en" {
		// 回退英文
		err = s.db.Where("entity_type = ? AND entity_id = ? AND language = ?", entityType, entityID, "en").
			First(&seo).Error
	}
	if err != nil {
		return nil, err
	}
	return &seo, nil
}

// UpsertSEO 创建或更新 SEO 配置（按 entity_type + entity_id + language 唯一）
func (s *SEOService) UpsertSEO(req *SEOMultiLangRequest) (*models.SEO, error) {
	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		return nil, fmt.Errorf("实体 ID 无效")
	}

	var seo models.SEO
	err = s.db.Where("entity_type = ? AND entity_id = ? AND language = ?", req.EntityType, entityID, req.Language).
		First(&seo).Error
	if err == nil {
		// 更新
		seo.Title = req.Title
		seo.Description = req.Description
		seo.Keywords = req.Keywords
		seo.Canonical = req.Canonical
		seo.Robots = req.Robots
		seo.OGTitle = req.OGTitle
		seo.OGDescription = req.OGDescription
		seo.OGImage = req.OGImage
		seo.SchemaData = req.SchemaData
		if err := s.db.Save(&seo).Error; err != nil {
			return nil, err
		}
		return &seo, nil
	}

	seo = models.SEO{
		EntityType:    req.EntityType,
		EntityID:      entityID,
		Language:      req.Language,
		Title:         req.Title,
		Description:   req.Description,
		Keywords:      req.Keywords,
		Canonical:     req.Canonical,
		Robots:        req.Robots,
		OGTitle:       req.OGTitle,
		OGDescription: req.OGDescription,
		OGImage:       req.OGImage,
		SchemaData:    req.SchemaData,
	}
	if err := s.db.Create(&seo).Error; err != nil {
		return nil, err
	}
	return &seo, nil
}

// SEOMultiLangRequest SEO 多语言请求
type SEOMultiLangRequest struct {
	EntityType    string `json:"entity_type" binding:"required"`
	EntityID      string `json:"entity_id" binding:"required"`
	Language      string `json:"language" binding:"required"`
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

// BuildHreflangLink 生成 hreflang 标注
// 示例：<link rel="alternate" hreflang="zh" href="https://site.com/zh/page-slug" />
func BuildHreflangLink(baseURL, lang, path string) string {
	return fmt.Sprintf(`<link rel="alternate" hreflang="%s" href="%s/%s/%s" />`, lang, strings.TrimRight(baseURL, "/"), lang, strings.TrimLeft(path, "/"))
}

// BuildXDefaultHreflang 生成 x-default 标注
func BuildXDefaultHreflang(baseURL, path string) string {
	return fmt.Sprintf(`<link rel="alternate" hreflang="x-default" href="%s/%s" />`, strings.TrimRight(baseURL, "/"), strings.TrimLeft(path, "/"))
}

// BuildSitemapEntry 生成 sitemap 条目
func BuildSitemapEntry(baseURL, lang, path string, lastmod string) string {
	loc := fmt.Sprintf("%s/%s/%s", strings.TrimRight(baseURL, "/"), lang, strings.TrimLeft(path, "/"))
	lastmodPart := ""
	if lastmod != "" {
		lastmodPart = fmt.Sprintf("\n    <lastmod>%s</lastmod>", lastmod)
	}
	return fmt.Sprintf("  <url>\n    <loc>%s</loc>%s\n  </url>", loc, lastmodPart)
}

// ==================== 路由（列表页/落地页）SEO ====================

// RouteSEOID 为路由生成确定性的实体 ID（UUID v5），使路由 SEO 复用 seo 表
// 的 entity_type='route' 维度，无需新增表。同一路由每次生成相同 ID。
func RouteSEOID(route string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte("route:"+route))
}

// RouteSEOEntryReq 路由 SEO 单语言请求（路由 + 语言矩阵，后台 SEO 管理页使用）
type RouteSEOEntryReq struct {
	Language      string `json:"language" binding:"required"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Keywords      string `json:"keywords"`
	OGTitle       string `json:"og_title"`
	OGDescription string `json:"og_description"`
	OGImage       string `json:"og_image"`
}

// GetRouteSEO 获取指定路由的 SEO 配置（按语言，无则回退英文）
func (s *SEOService) GetRouteSEO(route, lang string) (*models.SEO, error) {
	return s.GetSEO("route", RouteSEOID(route), lang)
}

// UpsertRouteSEO 创建或更新路由 SEO（按 route + language 唯一）
func (s *SEOService) UpsertRouteSEO(route string, req *RouteSEOEntryReq) (*models.SEO, error) {
	eid := RouteSEOID(route)
	var seo models.SEO
	err := s.db.Where("entity_type = ? AND entity_id = ? AND language = ?", "route", eid, req.Language).First(&seo).Error
	if err == nil {
		seo.Title = req.Title
		seo.Description = req.Description
		seo.Keywords = req.Keywords
		seo.OGTitle = req.OGTitle
		seo.OGDescription = req.OGDescription
		seo.OGImage = req.OGImage
		if err := s.db.Save(&seo).Error; err != nil {
			return nil, err
		}
		return &seo, nil
	}

	seo = models.SEO{
		EntityType:    "route",
		EntityID:      eid,
		Language:      req.Language,
		Title:         req.Title,
		Description:   req.Description,
		Keywords:      req.Keywords,
		OGTitle:       req.OGTitle,
		OGDescription: req.OGDescription,
		OGImage:       req.OGImage,
	}
	if err := s.db.Create(&seo).Error; err != nil {
		return nil, err
	}
	return &seo, nil
}

// ListRouteSEO 获取指定路由的全部语言 SEO 配置（后台管理）
func (s *SEOService) ListRouteSEO(route string) ([]models.SEO, error) {
	var seos []models.SEO
	err := s.db.Where("entity_type = ? AND entity_id = ?", "route", RouteSEOID(route)).Find(&seos).Error
	return seos, err
}
