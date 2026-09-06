package services

import (
	"strings"

	"github.com/google/uuid"

	"sportswear-backend/internal/models"
)

// ==================== 门户路由注册表 & 健康检查 ====================
// 以「路由（route）」为唯一锚点，串联 页面(slug) / 导航(url) / SEO(route)，
// 收敛此前散落在前端多处硬编码的「页面类型 → 路径」映射，形成闭环。
// 后台 SEO 管理、导航管理、健康检查页均消费本文件提供的能力。

// PagePath 根据页面类型与 slug 派生门户访问路径（唯一事实来源）。
// 与前台 Nuxt 页面路由对齐：列表页/固定页映射到固定路径，其余走 /:slug 兜底页（[slug].vue）。
func PagePath(pageType, slug string) string {
	switch models.PageType(pageType) {
	case models.PageTypeHome:
		return "/"
	case models.PageTypeProduct, models.PageTypeProductCategory:
		return "/products"
	case models.PageTypeBlog:
		return "/blog"
	case models.PageTypeCase:
		return "/cases"
	case models.PageTypeFAQ:
		return "/faq"
	case models.PageTypeContact:
		return "/contact"
	default:
		s := strings.TrimSpace(slug)
		if s == "" {
			return "/"
		}
		return "/" + strings.TrimLeft(s, "/")
	}
}

// RouteKey 根据页面类型与 slug 派生路由 key（与 seo 表 entity_type='route' 的 route 维度一致）。
func RouteKey(pageType models.PageType, slug string) string {
	switch pageType {
	case models.PageTypeHome:
		return "home"
	case models.PageTypeProduct, models.PageTypeProductCategory:
		return "products"
	case models.PageTypeBlog:
		return "blog"
	case models.PageTypeCase:
		return "cases"
	case models.PageTypeFAQ:
		return "faq"
	case models.PageTypeContact:
		return "contact"
	default:
		return slug
	}
}

// RouteFromPage 页面 → 路由 key（列表页/固定页映射到固定 key，其余落地页以 slug 作为 route key）。
func RouteFromPage(pg models.Page) string {
	return RouteKey(pg.Type, pg.Slug)
}

// presetRoute 系统预置路由（门户固定页面）
type presetRoute struct {
	Route string
	Path  string
	Label string
}

// presetRoutes 与门户固定页面一一对应。
var presetRoutes = []presetRoute{
	{"home", "/", "首页"},
	{"products", "/products", "产品中心"},
	{"cases", "/cases", "案例展示"},
	{"about", "/about", "关于我们"},
	{"blog", "/blog", "博客"},
	{"faq", "/faq", "常见问题"},
	{"contact", "/contact", "联系我们"},
	{"privacy", "/privacy-policy", "隐私政策"},
}

// PortalRouteEntry 路由条目（单一事实来源，供 SEO 下拉 / 导航派生 / 健康检查共用）
type PortalRouteEntry struct {
	Route         string `json:"route"`          // 短 key，SEO 管理用
	Path          string `json:"path"`           // 门户访问路径
	Label         string `json:"label"`          // 显示名
	Source        string `json:"source"`         // system / page
	PageID        string `json:"page_id,omitempty"`
	PageStatus    string `json:"page_status,omitempty"`
	SEOConfigured bool   `json:"seo_configured"` // 该路由是否已配置 route SEO（title/description 非空）
	NavCount      int    `json:"nav_count"`      // 关联的导航条数
}

// ListPortalRoutes 返回路由注册表：系统预置路由 + 已发布页面派生路由（去重），
// 并填充每个路由的 SEO 配置状态与导航关联数。
func (s *CMSService) ListPortalRoutes() ([]PortalRouteEntry, error) {
	entries := make([]PortalRouteEntry, 0, len(presetRoutes)+16)
	seen := make(map[string]bool)

	for _, p := range presetRoutes {
		entries = append(entries, PortalRouteEntry{Route: p.Route, Path: p.Path, Label: p.Label, Source: "system"})
		seen[p.Route] = true
	}

	var pages []models.Page
	if err := s.db.Where("status = ?", models.ContentStatusPublished).Find(&pages).Error; err != nil {
		return nil, err
	}
	for _, pg := range pages {
		route := RouteFromPage(pg)
		if route == "" || seen[route] {
			continue
		}
		seen[route] = true
		entries = append(entries, PortalRouteEntry{
			Route:      route,
			Path:       PagePath(string(pg.Type), pg.Slug),
			Label:      pg.Title,
			Source:     "page",
			PageID:     pg.ID.String(),
			PageStatus: string(pg.Status),
		})
	}

	for i := range entries {
		entries[i].SEOConfigured = s.routeSEOConfigured(entries[i].Route)
		entries[i].NavCount = s.navCountForRoute(entries[i])
	}

	return entries, nil
}

// routeSEOConfigured 判断指定路由是否已配置 route SEO（以英文 title/description 任一非空为准）
func (s *CMSService) routeSEOConfigured(route string) bool {
	var seo models.SEO
	err := s.db.Where("entity_type = ? AND entity_id = ? AND language = ?", "route", RouteSEOID(route), "en").First(&seo).Error
	if err != nil {
		return false
	}
	return strings.TrimSpace(seo.Title) != "" || strings.TrimSpace(seo.Description) != ""
}

// entitySEOConfiguredSet 批量判断实体是否已配置英文 SEO（title/description 任一非空）。
// 返回已配置的实体 ID 集合，供 missing_entity_seo 健康检查使用。
func (s *CMSService) entitySEOConfiguredSet(entityType string, ids []uuid.UUID) map[uuid.UUID]bool {
	result := make(map[uuid.UUID]bool, len(ids))
	if len(ids) == 0 {
		return result
	}
	var seos []models.SEO
	if err := s.db.Select("entity_id, title, description").
		Where("entity_type = ? AND entity_id IN ? AND language = ?", entityType, ids, "en").
		Find(&seos).Error; err != nil {
		return result
	}
	for _, seo := range seos {
		if strings.TrimSpace(seo.Title) != "" || strings.TrimSpace(seo.Description) != "" {
			result[seo.EntityID] = true
		}
	}
	return result
}

// entitySEOItem 详情实体 SEO 健康检查的引用项（ID + Slug）。
type entitySEOItem struct {
	ID   uuid.UUID
	Slug string
}

// appendMissingEntitySEO 将「已发布但缺英文源语言 SEO」的详情实体追加到健康报告。
func (s *CMSService) appendMissingEntitySEO(report *PortalHealthReport, entityType, pathPrefix, issueTitle string, rows []entitySEOItem) {
	ids := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.ID)
	}
	set := s.entitySEOConfiguredSet(entityType, ids)
	for _, r := range rows {
		if set[r.ID] {
			continue
		}
		report.Issues = append(report.Issues, HealthIssue{
			Type:   "missing_entity_seo",
			Title:  issueTitle,
			Detail: "「" + r.Slug + "」已发布但未配置 SEO 标题/描述，门户将回退到默认标题。",
			Route:  entityType,
			Path:   pathPrefix + r.Slug,
		})
	}
}

// navCountForRoute 统计指向某路由对应路径 / 关联页面的导航条数
func (s *CMSService) navCountForRoute(e PortalRouteEntry) int {
	q := s.db.Model(&models.Navigation{}).Where("url = ?", e.Path)
	if e.PageID != "" {
		if pid, err := uuid.Parse(e.PageID); err == nil {
			q = s.db.Model(&models.Navigation{}).Where("url = ? OR page_id = ?", e.Path, pid)
		}
	}
	var count int64
	_ = q.Count(&count).Error
	return int(count)
}

// ==================== 健康检查 ====================

// HealthIssue 单个闭环断点
type HealthIssue struct {
	Type   string `json:"type"`             // dead_link / missing_nav / missing_seo / missing_entity_seo / nav_drift
	Title  string `json:"title"`            // 简短描述
	Detail string `json:"detail,omitempty"` // 补充说明
	Route  string `json:"route,omitempty"`
	Path   string `json:"path,omitempty"`
	PageID string `json:"page_id,omitempty"`
	NavID  string `json:"nav_id,omitempty"`
}

// PortalHealthReport 门户闭环健康报告
type PortalHealthReport struct {
	TotalIssues int           `json:"total_issues"`
	Routes      int           `json:"routes"`
	NavCount    int           `json:"nav_count"`
	PageCount   int           `json:"page_count"`
	Issues      []HealthIssue `json:"issues"`
}

// PortalHealthCheck 聚合检测 4 类闭环断点：
//  1. dead_link    导航指向不存在的门户路径（内部链接且无对应已发布页面）
//  2. missing_nav  已发布页面（非首页）未关联任何导航
//  3. missing_seo  有路由但未配置 route SEO
//  4. nav_drift    导航 url 与关联页面 slug 派生的路径不一致
func (s *CMSService) PortalHealthCheck() (*PortalHealthReport, error) {
	report := &PortalHealthReport{Issues: []HealthIssue{}}

	var publishedPages []models.Page
	if err := s.db.Where("status = ?", models.ContentStatusPublished).Find(&publishedPages).Error; err != nil {
		return nil, err
	}
	report.PageCount = len(publishedPages)
	slugSet := make(map[string]bool, len(publishedPages))
	for _, pg := range publishedPages {
		slugSet[pg.Slug] = true
	}

	var navs []models.Navigation
	if err := s.db.Preload("Page").Find(&navs).Error; err != nil {
		return nil, err
	}
	report.NavCount = len(navs)

	// 1. dead_link
	for _, n := range navs {
		if isExternalURL(n.URL) || s.isValidInternalPath(n.URL, slugSet) {
			continue
		}
		report.Issues = append(report.Issues, HealthIssue{
			Type:   "dead_link",
			Title:  "导航指向不存在的路径",
			Detail: "导航「" + n.Name + "」的 URL " + n.URL + " 无对应门户路由，且未关联任何已发布页面。",
			Path:   n.URL,
			NavID:  n.ID.String(),
		})
	}

	// 4. nav_drift
	for _, n := range navs {
		if n.PageID == nil || n.Page == nil || isExternalURL(n.URL) {
			continue
		}
		expected := PagePath(string(n.Page.Type), n.Page.Slug)
		if strings.TrimRight(n.URL, "/") != strings.TrimRight(expected, "/") {
			report.Issues = append(report.Issues, HealthIssue{
				Type:   "nav_drift",
				Title:  "导航 URL 与页面路径不一致",
				Detail: "导航「" + n.Name + "」当前 URL 为 " + n.URL + "，而关联页面「" + n.Page.Title + "」派生路径为 " + expected + "。",
				Route:  RouteFromPage(*n.Page),
				Path:   n.URL,
				PageID: n.PageID.String(),
				NavID:  n.ID.String(),
			})
		}
	}

	// 2. missing_nav：已发布页面未出现在导航中。
	// 判据：① 有 page_id 关联导航；② 或存在导航 URL 等于该页面的派生路径（固定页面通过骨架导航 URL 匹配）。
	navPageIDs := make(map[string]bool, len(navs))
	navPaths := make(map[string]bool, len(navs))
	for _, n := range navs {
		if n.PageID != nil {
			navPageIDs[n.PageID.String()] = true
		}
		navPaths[strings.TrimRight(n.URL, "/")] = true
	}
	for _, pg := range publishedPages {
		if pg.Type == models.PageTypeHome || pg.Slug == "home" || navPageIDs[pg.ID.String()] {
			continue
		}
		if navPaths[strings.TrimRight(PagePath(string(pg.Type), pg.Slug), "/")] {
			continue
		}
		report.Issues = append(report.Issues, HealthIssue{
			Type:   "missing_nav",
			Title:  "已发布页面未加入导航",
			Detail: "页面「" + pg.Title + "」（" + PagePath(string(pg.Type), pg.Slug) + "）已发布但未关联任何导航项。",
			Route:  RouteFromPage(pg),
			Path:   PagePath(string(pg.Type), pg.Slug),
			PageID: pg.ID.String(),
		})
	}

	// 3. missing_seo
	routes, err := s.ListPortalRoutes()
	if err != nil {
		return nil, err
	}
	report.Routes = len(routes)
	for _, r := range routes {
		if r.SEOConfigured {
			continue
		}
		report.Issues = append(report.Issues, HealthIssue{
			Type:   "missing_seo",
			Title:  "路由未配置 SEO",
			Detail: "路由「" + r.Route + "」（" + r.Path + "）尚未配置 SEO 标题/描述，门户将回退到页面默认值。",
			Route:  r.Route,
			Path:   r.Path,
			PageID: r.PageID,
		})
	}

	// 5. missing_entity_seo：详情实体（产品/博客/案例）缺英文源语言 SEO
	var productRefs []entitySEOItem
	if err := s.db.Model(&models.Product{}).Select("id, slug").Where("status = ?", models.ProductStatusPublished).Scan(&productRefs).Error; err != nil {
		return nil, err
	}
	s.appendMissingEntitySEO(report, "product", "/products/", "产品详情未配置 SEO", productRefs)

	var blogRefs []entitySEOItem
	if err := s.db.Model(&models.Blog{}).Select("id, slug").Where("status = ?", models.ContentStatusPublished).Scan(&blogRefs).Error; err != nil {
		return nil, err
	}
	s.appendMissingEntitySEO(report, "blog", "/blog/", "博客详情未配置 SEO", blogRefs)

	var caseRefs []entitySEOItem
	if err := s.db.Model(&models.Case{}).Select("id, slug").Where("status = ?", models.ContentStatusPublished).Scan(&caseRefs).Error; err != nil {
		return nil, err
	}
	s.appendMissingEntitySEO(report, "case", "/cases/", "案例详情未配置 SEO", caseRefs)

	report.TotalIssues = len(report.Issues)
	return report, nil
}

// isExternalURL 判断导航 URL 是否为外部链接 / 锚点（此类不参与死链判定）
func isExternalURL(url string) bool {
	u := strings.TrimSpace(url)
	if u == "" {
		return true
	}
	return strings.HasPrefix(u, "http://") ||
		strings.HasPrefix(u, "https://") ||
		strings.HasPrefix(u, "mailto:") ||
		strings.HasPrefix(u, "tel:") ||
		strings.HasPrefix(u, "#")
}

// isValidInternalPath 判断内部路径是否命中门户有效路由：
// 预设固定路径 / 动态详情页前缀 / 单段 slug（走 [slug].vue 兜底，需存在对应已发布页面）。
func (s *CMSService) isValidInternalPath(path string, publishedSlugs map[string]bool) bool {
	exact := map[string]bool{
		"/": true, "/products": true, "/blog": true, "/cases": true,
		"/faq": true, "/contact": true, "/about": true, "/privacy-policy": true,
	}
	if exact[path] {
		return true
	}
	prefixes := []string{"/products/", "/blog/", "/cases/", "/categories/", "/series/"}
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	rest := strings.Trim(path, "/")
	if rest != "" && !strings.Contains(rest, "/") {
		return publishedSlugs[rest]
	}
	return false
}

// ==================== slug 联动 ====================

// SyncPagePathEffects 页面 slug/type 变更后，级联更新关联导航 url 与 route SEO，避免死链与 SEO 失效。
// 仅更新「未手动改过」（url 仍等于旧派生路径）的导航；route SEO 仅在目标不存在时迁移，避免覆盖。
// 此方法为尽力而为的副作用同步，失败不阻断主流程。
func (s *CMSService) SyncPagePathEffects(pageID uuid.UUID, oldSlug, newSlug string, oldType, newType models.PageType) {
	oldPath := PagePath(string(oldType), oldSlug)
	newPath := PagePath(string(newType), newSlug)
	if oldPath != newPath {
		s.db.Model(&models.Navigation{}).
			Where("page_id = ? AND url = ?", pageID, oldPath).
			Update("url", newPath)
	}

	oldRoute := RouteKey(oldType, oldSlug)
	newRoute := RouteKey(newType, newSlug)
	if oldRoute != newRoute {
		var target models.SEO
		if err := s.db.Where("entity_type = ? AND entity_id = ?", "route", RouteSEOID(newRoute)).First(&target).Error; err == nil {
			return // 目标路由已有 SEO，保留现状，避免覆盖
		}
		s.db.Model(&models.SEO{}).
			Where("entity_type = ? AND entity_id = ?", "route", RouteSEOID(oldRoute)).
			Update("entity_id", RouteSEOID(newRoute))
	}
}


