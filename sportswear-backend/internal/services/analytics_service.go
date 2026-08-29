package services

import (
	"time"

	"gorm.io/gorm"

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

// GetTrafficOverview 流量总览
func (s *AnalyticsService) GetTrafficOverview(days int) (map[string]interface{}, error) {
	if days <= 0 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days)

	var totalVisits, uniqueIPs, botVisits, productViews, leadCount int64

	s.db.Model(&models.VisitLog{}).Where("created_at >= ?", since).Count(&totalVisits)
	s.db.Model(&models.VisitLog{}).Where("created_at >= ?", since).Distinct("ip").Count(&uniqueIPs)
	s.db.Model(&models.VisitLog{}).Where("created_at >= ? AND device = ?", since, "bot").Count(&botVisits)
	s.db.Model(&models.VisitLog{}).Where("created_at >= ? AND entity_type = ?", since, "product").Count(&productViews)
	s.db.Model(&models.Lead{}).Where("created_at >= ?", since).Count(&leadCount)

	// 询盘转化率（访问 IP → 询盘）
	conversionRate := float64(0)
	if uniqueIPs > 0 {
		conversionRate = float64(leadCount) / float64(uniqueIPs) * 100
	}

	// 每日访问趋势
	var dailyTrend []map[string]interface{}
	s.db.Model(&models.VisitLog{}).
		Select("DATE(created_at) as date, COUNT(*) as visits, COUNT(DISTINCT ip) as unique_ips").
		Where("created_at >= ?", since).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&dailyTrend)

	return map[string]interface{}{
		"days":            days,
		"total_visits":    totalVisits,
		"unique_ips":      uniqueIPs,
		"bot_visits":      botVisits,
		"human_visits":    totalVisits - botVisits,
		"product_views":   productViews,
		"leads":           leadCount,
		"conversion_rate": conversionRate, // 询盘转化率 %
		"daily_trend":     dailyTrend,
	}, nil
}

// GetTopPages 热门页面
func (s *AnalyticsService) GetTopPages(days, limit int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 7
	}
	if limit <= 0 {
		limit = 10
	}
	since := time.Now().AddDate(0, 0, -days)

	var topPages []map[string]interface{}
	err := s.db.Model(&models.VisitLog{}).
		Select("path, entity_type, entity_slug, COUNT(*) as views, COUNT(DISTINCT ip) as unique_visitors").
		Where("created_at >= ? AND device != ?", since, "bot").
		Group("path, entity_type, entity_slug").
		Order("views DESC").
		Limit(limit).
		Scan(&topPages).Error
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

	var topProducts []map[string]interface{}
	err := s.db.Model(&models.VisitLog{}).
		Select("entity_slug as slug, COUNT(*) as views, COUNT(DISTINCT ip) as unique_visitors").
		Where("created_at >= ? AND entity_type = ? AND device != ?", since, "product", "bot").
		Group("entity_slug").
		Order("views DESC").
		Limit(limit).
		Scan(&topProducts).Error
	return topProducts, err
}

// GetSourceAnalysis 来源分析（SEO/广告归因，商业价值核心）
func (s *AnalyticsService) GetSourceAnalysis(days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	var sources []map[string]interface{}
	err := s.db.Model(&models.VisitLog{}).
		Select("source, medium, COUNT(*) as visits, COUNT(DISTINCT ip) as unique_visitors").
		Where("created_at >= ? AND device != ?", since, "bot").
		Group("source, medium").
		Order("visits DESC").
		Scan(&sources).Error
	return sources, err
}

// GetCountryAnalysis 国家/地区分布
func (s *AnalyticsService) GetCountryAnalysis(days, limit int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	if limit <= 0 {
		limit = 15
	}
	since := time.Now().AddDate(0, 0, -days)

	var countries []map[string]interface{}
	err := s.db.Raw(`
		SELECT COALESCE(l.country, v.ip) as region, COUNT(*) as visits, COUNT(DISTINCT v.ip) as unique_visitors
		FROM visit_logs v
		LEFT JOIN leads l ON l.ip = v.ip
		WHERE v.created_at >= ? AND v.device != ?
		GROUP BY COALESCE(l.country, v.ip)
		ORDER BY visits DESC
		LIMIT ?
	`, since, "bot", limit).Scan(&countries).Error
	return countries, err
}

// GetSocialClickAnalysis 社交媒体点击分析（外部链接跳转跟踪）
func (s *AnalyticsService) GetSocialClickAnalysis(days int) ([]map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	var socialClicks []map[string]interface{}
	err := s.db.Model(&models.VisitLog{}).
		Select("entity_slug as platform, COUNT(*) as clicks, COUNT(DISTINCT ip) as unique_visitors").
		Where("created_at >= ? AND entity_type = ? AND device != ?", since, "social_click", "bot").
		Group("entity_slug").
		Order("clicks DESC").
		Scan(&socialClicks).Error
	return socialClicks, err
}

// GetDeviceAnalysis 设备/浏览器分析
func (s *AnalyticsService) GetDeviceAnalysis(days int) (map[string]interface{}, error) {
	if days <= 0 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	var devices, browsers []map[string]interface{}

	s.db.Model(&models.VisitLog{}).
		Select("device, COUNT(*) as visits").
		Where("created_at >= ? AND device != ?", since, "bot").
		Group("device").Order("visits DESC").
		Scan(&devices)

	s.db.Model(&models.VisitLog{}).
		Select("browser, COUNT(*) as visits").
		Where("created_at >= ? AND device != ?", since, "bot").
		Group("browser").Order("visits DESC").
		Scan(&browsers)

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
