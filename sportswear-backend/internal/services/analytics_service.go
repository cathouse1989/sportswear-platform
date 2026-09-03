package services

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"sportswear-backend/internal/database"
	"sportswear-backend/internal/models"
)

// AnalyticsService 流量分析服务
type AnalyticsService struct {
	db *gorm.DB
}

// NewAnalyticsService 创建流量分析服务
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// MaxQueryDays 最大查询时间范围（90 天），避免跨过多月度表导致性能问题
const MaxQueryDays = 90

// normalizeRange 规范化时间范围：限制最大跨度不超过 MaxQueryDays
func normalizeRange(start, end time.Time) (time.Time, time.Time) {
	if start.IsZero() {
		start = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if end.IsZero() {
		end = time.Now()
	}
	if end.Before(start) {
		start, end = end, start
	}
	// 限制最大查询范围
	maxStart := end.AddDate(0, 0, -MaxQueryDays)
	if start.Before(maxStart) {
		start = maxStart
	}
	return start, end
}

// visitTables 返回查询涉及的 visit_logs 表：存量基础表（历史数据）+ 时间范围内的月度表（仅已存在）
func (s *AnalyticsService) visitTables(start, end time.Time) ([]string, error) {
	start, end = normalizeRange(start, end)
	existing, err := database.ExistingTables(s.db, "visit_logs")
	if err != nil {
		return nil, err
	}
	exMap := make(map[string]bool, len(existing))
	for _, n := range existing {
		exMap[n] = true
	}
	candidates := database.MonthTablesBetween("visit_logs", start, end)
	tables := []string{"visit_logs"} // 存量基础表
	for _, n := range candidates {
		if exMap[n] {
			tables = append(tables, n)
		}
	}
	return tables, nil
}

// GetTrafficOverview 流量总览（访问量、独立IP、转化率、每日趋势）——跨季度分表组合查询
func (s *AnalyticsService) GetTrafficOverview(days int) (map[string]interface{}, error) {
	if days <= 0 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days)
	end := time.Now()

	tables, err := s.visitTables(since, end)
	if err != nil {
		return nil, err
	}

	var totalVisits, uniqueIPs, uniqueVisitors, botVisits, productViews, leadCount int64

	// 全量页面浏览（只统计 page_view，排除 api_call，解决重复统计问题）
	// 按 visitor_id + path 去重（同一访客同一页面只算一次浏览）
	visitsSelect := `SELECT visitor_id, path, created_at, ip, device, entity_type FROM __T__ WHERE created_at >= ? AND visit_type = 'page_view'`
	visitsSQL, visitsArgs := database.UnionAll(tables, visitsSelect, []interface{}{since})
	if err := s.db.Raw(fmt.Sprintf(`SELECT COUNT(DISTINCT visitor_id || ':' || path) FROM (%s) t`, visitsSQL), visitsArgs...).Scan(&totalVisits).Error; err != nil {
		return nil, err
	}
	// 独立 IP
	if err := s.db.Raw(fmt.Sprintf(`SELECT COUNT(DISTINCT ip) FROM (%s) t WHERE visit_type = 'page_view'`, visitsSQL), visitsArgs...).Scan(&uniqueIPs).Error; err != nil {
		return nil, err
	}
	// 独立访客（按 visitor_id 去重，比 IP 更准确）
	if err := s.db.Raw(fmt.Sprintf(`SELECT COUNT(DISTINCT visitor_id) FROM (%s) t WHERE visit_type = 'page_view' AND visitor_id != 'unknown' AND visitor_id != 'anonymous'`, visitsSQL), visitsArgs...).Scan(&uniqueVisitors).Error; err != nil {
		return nil, err
	}
	// 机器人流量
	botSelect := `SELECT created_at FROM __T__ WHERE created_at >= ? AND device = 'bot'`
	botSQL, botArgs := database.UnionAll(tables, botSelect, []interface{}{since})
	if err := s.db.Raw(fmt.Sprintf(`SELECT COUNT(*) FROM (%s) t`, botSQL), botArgs...).Scan(&botVisits).Error; err != nil {
		return nil, err
	}
	// 产品浏览（只统计 page_view）
	prodSelect := `SELECT visitor_id, path FROM __T__ WHERE created_at >= ? AND entity_type = 'product' AND visit_type = 'page_view'`
	prodSQL, prodArgs := database.UnionAll(tables, prodSelect, []interface{}{since})
	if err := s.db.Raw(fmt.Sprintf(`SELECT COUNT(DISTINCT visitor_id || ':' || path) FROM (%s) t`, prodSQL), prodArgs...).Scan(&productViews).Error; err != nil {
		return nil, err
	}
	// 询盘
	s.db.Model(&models.Lead{}).Where("created_at >= ?", since).Count(&leadCount)

	// 询盘转化率（独立访客 → 询盘）
	conversionRate := float64(0)
	if uniqueVisitors > 0 {
		conversionRate = float64(leadCount) / float64(uniqueVisitors) * 100
	}

	// 每日访问趋势（只统计 page_view，按 visitor_id + path 去重）
	var dailyTrend []map[string]interface{}
	trendSQL := fmt.Sprintf(`SELECT DATE(created_at) AS date, COUNT(DISTINCT visitor_id || ':' || path) AS visits, COUNT(DISTINCT ip) AS unique_ips FROM (%s) t WHERE visit_type = 'page_view' GROUP BY DATE(created_at) ORDER BY date ASC`, visitsSQL)
	if err := s.db.Raw(trendSQL, visitsArgs...).Scan(&dailyTrend).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"days":            days,
		"total_visits":    totalVisits,
		"unique_ips":      uniqueIPs,
		"unique_visitors": uniqueVisitors,
		"bot_visits":      botVisits,
		"human_visits":    totalVisits - botVisits,
		"product_views":   productViews,
		"leads":           leadCount,
		"conversion_rate": conversionRate, // 询盘转化率 %
		"daily_trend":     dailyTrend,
	}, nil
}

// GetTopPages 热门页面（跨季度分表组合查询）
func (s *AnalyticsService) GetTopPages(days, limit int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 7
	}
	if limit <= 0 {
		limit = 10
	}
	since := time.Now().AddDate(0, 0, -days)
	tables, err := s.visitTables(since, time.Now())
	if err != nil {
		return nil, err
	}

	selectFmt := `SELECT path, entity_type, entity_slug, COUNT(*) AS views, COUNT(DISTINCT ip) AS unique_visitors FROM __T__ WHERE created_at >= ? AND device != 'bot' GROUP BY path, entity_type, entity_slug`
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, []interface{}{since})

	var topPages []map[string]interface{}
	sql := fmt.Sprintf(`SELECT path, entity_type, entity_slug, SUM(views) AS views, SUM(unique_visitors) AS unique_visitors FROM (%s) t GROUP BY path, entity_type, entity_slug ORDER BY views DESC LIMIT ?`, unionSQL)
	err = s.db.Raw(sql, append(unionArgs, limit)...).Scan(&topPages).Error
	return topPages, err
}

// GetTopProducts 热门产品（商业价值：哪些产品最受关注）
func (s *AnalyticsService) GetTopProducts(days, limit int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 7
	}
	if limit <= 0 {
		limit = 10
	}
	since := time.Now().AddDate(0, 0, -days)
	tables, err := s.visitTables(since, time.Now())
	if err != nil {
		return nil, err
	}

	selectFmt := `SELECT entity_slug AS slug, MAX(entity_name) AS name, COUNT(*) AS views, COUNT(DISTINCT ip) AS unique_visitors FROM __T__ WHERE created_at >= ? AND entity_type = 'product' AND device != 'bot' GROUP BY entity_slug`
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, []interface{}{since})

	var topProducts []map[string]interface{}
	sql := fmt.Sprintf(`SELECT slug, MAX(name) AS name, SUM(views) AS views, SUM(unique_visitors) AS unique_visitors FROM (%s) t GROUP BY slug ORDER BY views DESC LIMIT ?`, unionSQL)
	err = s.db.Raw(sql, append(unionArgs, limit)...).Scan(&topProducts).Error
	return topProducts, err
}

// GetSourceAnalysis 来源分析（SEO/广告归因，商业价值核心）
func (s *AnalyticsService) GetSourceAnalysis(days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)
	tables, err := s.visitTables(since, time.Now())
	if err != nil {
		return nil, err
	}

	selectFmt := `SELECT source, medium, COUNT(*) AS visits, COUNT(DISTINCT ip) AS unique_visitors FROM __T__ WHERE created_at >= ? AND device != 'bot' AND source != '' GROUP BY source, medium`
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, []interface{}{since})

	var sources []map[string]interface{}
	sql := fmt.Sprintf(`SELECT source, medium, SUM(visits) AS visits, SUM(unique_visitors) AS unique_visitors FROM (%s) t GROUP BY source, medium ORDER BY visits DESC`, unionSQL)
	err = s.db.Raw(sql, unionArgs...).Scan(&sources).Error
	return sources, err
}

// GetCountryAnalysis 国家/地区分布（直接使用 visit_logs 记录的国家字段）
func (s *AnalyticsService) GetCountryAnalysis(days, limit int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	if limit <= 0 {
		limit = 15
	}
	since := time.Now().AddDate(0, 0, -days)
	tables, err := s.visitTables(since, time.Now())
	if err != nil {
		return nil, err
	}

	selectFmt := `SELECT COALESCE(NULLIF(country, ''), 'unknown') AS country, COUNT(*) AS visits, COUNT(DISTINCT ip) AS unique_visitors FROM __T__ WHERE created_at >= ? AND device != 'bot' GROUP BY country`
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, []interface{}{since})

	var countries []map[string]interface{}
	sql := fmt.Sprintf(`SELECT country, SUM(visits) AS visits, SUM(unique_visitors) AS unique_visitors FROM (%s) t GROUP BY country ORDER BY visits DESC LIMIT ?`, unionSQL)
	err = s.db.Raw(sql, append(unionArgs, limit)...).Scan(&countries).Error
	return countries, err
}

// GetUtmCampaignAnalysis 广告归因分析（按 utm_campaign 分组，衡量各广告 Campaign 带来的流量）
func (s *AnalyticsService) GetUtmCampaignAnalysis(days, limit int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	if limit <= 0 {
		limit = 20
	}
	since := time.Now().AddDate(0, 0, -days)
	tables, err := s.visitTables(since, time.Now())
	if err != nil {
		return nil, err
	}

	selectFmt := `SELECT utm_campaign AS campaign, MAX(source) AS source, MAX(medium) AS medium, COUNT(*) AS visits, COUNT(DISTINCT ip) AS unique_visitors FROM __T__ WHERE created_at >= ? AND device != 'bot' AND utm_campaign != '' GROUP BY utm_campaign`
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, []interface{}{since})

	var campaigns []map[string]interface{}
	sql := fmt.Sprintf(`SELECT campaign, MAX(source) AS source, MAX(medium) AS medium, SUM(visits) AS visits, SUM(unique_visitors) AS unique_visitors FROM (%s) t GROUP BY campaign ORDER BY visits DESC LIMIT ?`, unionSQL)
	err = s.db.Raw(sql, append(unionArgs, limit)...).Scan(&campaigns).Error
	return campaigns, err
}

// GetSocialClickAnalysis 社交媒体点击分析（外部链接跳转跟踪）
func (s *AnalyticsService) GetSocialClickAnalysis(days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)
	tables, err := s.visitTables(since, time.Now())
	if err != nil {
		return nil, err
	}

	selectFmt := `SELECT entity_slug AS platform, COUNT(*) AS clicks, COUNT(DISTINCT ip) AS unique_visitors FROM __T__ WHERE created_at >= ? AND entity_type = 'social_click' AND device != 'bot' GROUP BY entity_slug`
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, []interface{}{since})

	var socialClicks []map[string]interface{}
	sql := fmt.Sprintf(`SELECT platform, SUM(clicks) AS clicks, SUM(unique_visitors) AS unique_visitors FROM (%s) t GROUP BY platform ORDER BY clicks DESC`, unionSQL)
	err = s.db.Raw(sql, unionArgs...).Scan(&socialClicks).Error
	return socialClicks, err
}

// GetDeviceAnalysis 设备/浏览器分析
func (s *AnalyticsService) GetDeviceAnalysis(days int) (map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)
	tables, err := s.visitTables(since, time.Now())
	if err != nil {
		return nil, err
	}

	var devices, browsers []map[string]interface{}

	deviceSelect := `SELECT device, COUNT(*) AS visits FROM __T__ WHERE created_at >= ? AND device != 'bot' AND device != '' GROUP BY device`
	deviceSQL, deviceArgs := database.UnionAll(tables, deviceSelect, []interface{}{since})
	if err := s.db.Raw(fmt.Sprintf(`SELECT device, SUM(visits) AS visits FROM (%s) t GROUP BY device ORDER BY visits DESC`, deviceSQL), deviceArgs...).Scan(&devices).Error; err != nil {
		return nil, err
	}

	browserSelect := `SELECT browser, COUNT(*) AS visits FROM __T__ WHERE created_at >= ? AND device != 'bot' AND browser != '' GROUP BY browser`
	browserSQL, browserArgs := database.UnionAll(tables, browserSelect, []interface{}{since})
	if err := s.db.Raw(fmt.Sprintf(`SELECT browser, SUM(visits) AS visits FROM (%s) t GROUP BY browser ORDER BY visits DESC`, browserSQL), browserArgs...).Scan(&browsers).Error; err != nil {
		return nil, err
	}

	// 计算移动端占比（供前端统计卡直接使用）
	totalVisits := 0.0
	mobileVisits := 0.0
	for _, d := range devices {
		v, _ := d["visits"].(int64)
		totalVisits += float64(v)
		if device, _ := d["device"].(string); device == "mobile" || device == "tablet" {
			mobileVisits += float64(v)
		}
	}
	mobileRatio := 0.0
	if totalVisits > 0 {
		mobileRatio = mobileVisits / totalVisits
	}

	return map[string]interface{}{
		"devices":      devices,
		"browsers":     browsers,
		"total_visits": int64(totalVisits),
		"mobile_ratio": mobileRatio,
	}, nil
}

// VisitLogFilter 访问日志明细筛选条件
type VisitLogFilter struct {
	IP         string
	Country    string
	Source     string
	EntityType string
	EntitySlug string
	Language   string
	Device     string
	Keyword    string
	Start      time.Time
	End        time.Time
}

// ListVisitLogs 门户访问日志明细（按时间范围跨季度分表组合查询）
func (s *AnalyticsService) ListVisitLogs(page, pageSize int, f VisitLogFilter) ([]map[string]interface{}, int64, error) {
	tables, err := s.visitTables(f.Start, f.End)
	if err != nil {
		return nil, 0, err
	}

	var conds []string
	var args []interface{}
	if f.IP != "" {
		conds = append(conds, "ip = ?")
		args = append(args, f.IP)
	}
	if f.Country != "" {
		conds = append(conds, "country = ?")
		args = append(args, f.Country)
	}
	if f.Source != "" {
		conds = append(conds, "source = ?")
		args = append(args, f.Source)
	}
	if f.EntityType != "" {
		conds = append(conds, "entity_type = ?")
		args = append(args, f.EntityType)
	}
	if f.EntitySlug != "" {
		conds = append(conds, "entity_slug = ?")
		args = append(args, f.EntitySlug)
	}
	if f.Language != "" {
		conds = append(conds, "language = ?")
		args = append(args, f.Language)
	}
	if f.Device != "" {
		conds = append(conds, "device = ?")
		args = append(args, f.Device)
	}
	if f.Keyword != "" {
		k := "%" + f.Keyword + "%"
		conds = append(conds, "(entity_name ILIKE ? OR entity_slug ILIKE ? OR path ILIKE ?)")
		args = append(args, k, k, k)
	}
	if !f.Start.IsZero() {
		conds = append(conds, "created_at >= ?")
		args = append(args, f.Start)
	}
	if !f.End.IsZero() {
		conds = append(conds, "created_at <= ?")
		args = append(args, f.End)
	}
	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}

	selectFmt := `SELECT "id","created_at","ip","country","language","device","device_model","browser","os","user_agent","method","path","entity_type","entity_slug","entity_id","entity_name","status","latency_ms","referer","source","medium","keyword","utm_source","utm_medium","utm_campaign","utm_content","utm_term","session_id",'__T__' AS src FROM __T__` + where
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, args)

	var total int64
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM (%s) t`, unionSQL)
	if err := s.db.Raw(countSQL, unionArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	listArgs := make([]interface{}, 0, len(unionArgs)+2)
	listArgs = append(listArgs, unionArgs...)
	listArgs = append(listArgs, pageSize, (page-1)*pageSize)
	listSQL := fmt.Sprintf(`SELECT * FROM (%s) t ORDER BY created_at DESC LIMIT ? OFFSET ?`, unionSQL)

	var rows []map[string]interface{}
	if err := s.db.Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// GetVisitorJourney 获取访客旅程（某个访客的完整访问路径）
func (s *AnalyticsService) GetVisitorJourney(visitorID string, days int) ([]map[string]interface{}, error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)
	end := time.Now()

	tables, err := s.visitTables(since, end)
	if err != nil {
		return nil, err
	}

	selectFmt := `SELECT "created_at","path","entity_type","entity_slug","entity_name","visit_type","source","ip","country","device","session_id" FROM __T__ WHERE visitor_id = ? AND created_at >= ? AND visit_type = 'page_view'`
	args := []interface{}{visitorID, since}
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, args)

	listSQL := fmt.Sprintf(`SELECT * FROM (%s) t ORDER BY created_at ASC`, unionSQL)
	var rows []map[string]interface{}
	if err := s.db.Raw(listSQL, unionArgs...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetIPLeads 获取 IP 关联的询盘列表
func (s *AnalyticsService) GetIPLeads(ip string, days int) ([]map[string]interface{}, error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	var rows []map[string]interface{}
	err := s.db.Model(&models.Lead{}).
		Select("id, name, company, email, country, status, score, score_level, source, medium, campaign, visitor_id, ip, created_at").
		Where("ip = ? AND created_at >= ?", ip, since).
		Order("created_at DESC").
		Scan(&rows).Error
	return rows, err
}

// ConversionFunnel 转化漏斗数据
type ConversionFunnel struct {
	PageViews    int64   `json:"page_views"`
	ProductViews int64   `json:"product_views"`
	Leads        int64   `json:"leads"`
	ProductRate  float64 `json:"product_rate"`  // 页面浏览 → 产品浏览
	LeadRate     float64 `json:"lead_rate"`     // 产品浏览 → 询盘
}

// GetConversionFunnel 获取转化漏斗（访问 → 产品浏览 → 询盘）
func (s *AnalyticsService) GetConversionFunnel(days int) (*ConversionFunnel, error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)
	end := time.Now()

	tables, err := s.visitTables(since, end)
	if err != nil {
		return nil, err
	}

	var pageViews, productViews int64

	// 页面浏览量
	pvSelect := `SELECT visitor_id, path FROM __T__ WHERE created_at >= ? AND visit_type = 'page_view'`
	pvSQL, pvArgs := database.UnionAll(tables, pvSelect, []interface{}{since})
	s.db.Raw(fmt.Sprintf(`SELECT COUNT(DISTINCT visitor_id || ':' || path) FROM (%s) t`, pvSQL), pvArgs...).Scan(&pageViews)

	// 产品浏览量
	prodSelect := `SELECT visitor_id, path FROM __T__ WHERE created_at >= ? AND visit_type = 'page_view' AND entity_type = 'product'`
	prodSQL, prodArgs := database.UnionAll(tables, prodSelect, []interface{}{since})
	s.db.Raw(fmt.Sprintf(`SELECT COUNT(DISTINCT visitor_id || ':' || path) FROM (%s) t`, prodSQL), prodArgs...).Scan(&productViews)

	// 询盘数
	var leads int64
	s.db.Model(&models.Lead{}).Where("created_at >= ?", since).Count(&leads)

	funnel := &ConversionFunnel{
		PageViews:    pageViews,
		ProductViews: productViews,
		Leads:        leads,
	}
	if pageViews > 0 {
		funnel.ProductRate = float64(productViews) / float64(pageViews) * 100
	}
	if productViews > 0 {
		funnel.LeadRate = float64(leads) / float64(productViews) * 100
	}

	// 避免未使用 end 的编译警告
	_ = end

	return funnel, nil
}

// ConsentInsights 隐私合规洞察数据
type ConsentInsights struct {
	TotalVisitors   int64   `json:"total_visitors"`    // 总独立访客
	GrantedCount    int64   `json:"granted_count"`     // 同意分析的人数
	DeniedCount     int64   `json:"denied_count"`      // 拒绝分析的人数
	ConsentRate     float64 `json:"consent_rate"`      // 同意率 %
	GrantedLeads    int64   `json:"granted_leads"`     // 同意用户的询盘数
	DeniedLeads     int64   `json:"denied_leads"`      // 拒绝用户的询盘数
	GrantedConversion float64 `json:"granted_conversion"` // 同意用户转化率 %
	DeniedConversion  float64 `json:"denied_conversion"`  // 拒绝用户转化率 %
	AnonymizedRatio float64 `json:"anonymized_ratio"`  // 匿名化数据占比 %
}

// GetConsentInsights 获取隐私合规洞察（同意/未同意用户的转化对比）
func (s *AnalyticsService) GetConsentInsights(days int) (*ConsentInsights, error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)
	end := time.Now()

	tables, err := s.visitTables(since, end)
	if err != nil {
		return nil, err
	}

	var grantedVisitors, deniedVisitors, totalVisitors int64

	// 同意分析的访客数
	grantedSelect := `SELECT visitor_id FROM __T__ WHERE created_at >= ? AND visit_type = 'page_view' AND consent_status = 'granted' AND visitor_id != 'unknown'`
	grantedSQL, grantedArgs := database.UnionAll(tables, grantedSelect, []interface{}{since})
	s.db.Raw(fmt.Sprintf(`SELECT COUNT(DISTINCT visitor_id) FROM (%s) t`, grantedSQL), grantedArgs...).Scan(&grantedVisitors)

	// 拒绝分析的访客数
	deniedSelect := `SELECT visitor_id FROM __T__ WHERE created_at >= ? AND visit_type = 'page_view' AND consent_status = 'denied' AND visitor_id != 'unknown'`
	deniedSQL, deniedArgs := database.UnionAll(tables, deniedSelect, []interface{}{since})
	s.db.Raw(fmt.Sprintf(`SELECT COUNT(DISTINCT visitor_id) FROM (%s) t`, deniedSQL), deniedArgs...).Scan(&deniedVisitors)

	totalVisitors = grantedVisitors + deniedVisitors

	// 同意用户的询盘数（通过 visitor_id 关联）
	var grantedLeads int64
	s.db.Model(&models.Lead{}).
		Where("created_at >= ? AND visitor_id != '' AND visitor_id != 'unknown' AND visitor_id != 'anonymous'", since).
		Count(&grantedLeads)

	// 拒绝用户的询盘数（通过 IP 关联，排除已有 visitor_id 的）
	var deniedLeads int64
	s.db.Model(&models.Lead{}).
		Where("created_at >= ? AND (visitor_id = '' OR visitor_id = 'unknown' OR visitor_id = 'anonymous')", since).
		Count(&deniedLeads)

	insights := &ConsentInsights{
		TotalVisitors: totalVisitors,
		GrantedCount:  grantedVisitors,
		DeniedCount:   deniedVisitors,
		GrantedLeads:  grantedLeads,
		DeniedLeads:   deniedLeads,
	}

	if totalVisitors > 0 {
		insights.ConsentRate = float64(grantedVisitors) / float64(totalVisitors) * 100
		insights.AnonymizedRatio = float64(deniedVisitors) / float64(totalVisitors) * 100
	}
	if grantedVisitors > 0 {
		insights.GrantedConversion = float64(grantedLeads) / float64(grantedVisitors) * 100
	}
	if deniedVisitors > 0 {
		insights.DeniedConversion = float64(deniedLeads) / float64(deniedVisitors) * 100
	}

	// 避免未使用 end 的编译警告
	_ = end

	return insights, nil
}
