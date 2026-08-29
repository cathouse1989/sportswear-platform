package middleware

import (
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
	"sportswear-backend/internal/utils"
)

// visitLogChan 异步写入通道（不阻塞请求）
var visitLogChan chan models.VisitLog

func init() {
	visitLogChan = make(chan models.VisitLog, 1000)
	go func() {
		var db *gorm.DB
		for log := range visitLogChan {
			if db == nil {
				db = getDB()
				if db == nil {
					continue
				}
			}
			db.Create(&log)
		}
	}()
}

var analyticsDB *gorm.DB

// SetAnalyticsDB 设置分析数据库连接
func SetAnalyticsDB(db *gorm.DB) {
	analyticsDB = db
}

func getDB() *gorm.DB {
	return analyticsDB
}

// UA 解析正则
var (
	botRegex    = regexp.MustCompile(`(?i)(bot|crawler|spider|slurp|bingpreview|facebookexternalhit)`)
	mobileRegex = regexp.MustCompile(`(?i)(mobile|android|iphone|ipad|phone)`)
)

// ParseUserAgent 解析 User-Agent（导出，供其他包记录点击事件）
func ParseUserAgent(ua string) (device, browser, os string) {
	if ua == "" {
		return "unknown", "unknown", "unknown"
	}
	if botRegex.MatchString(ua) {
		return "bot", "bot", "bot"
	}
	device = "desktop"
	if mobileRegex.MatchString(ua) {
		device = "mobile"
	}

	switch {
	case strings.Contains(ua, "Edg/"):
		browser = "Edge"
	case strings.Contains(ua, "Chrome/"):
		browser = "Chrome"
	case strings.Contains(ua, "Firefox/"):
		browser = "Firefox"
	case strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome"):
		browser = "Safari"
	default:
		browser = "Other"
	}

	switch {
	case strings.Contains(ua, "Windows"):
		os = "Windows"
	case strings.Contains(ua, "Mac OS"):
		os = "macOS"
	case strings.Contains(ua, "Android"):
		os = "Android"
	case strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad"):
		os = "iOS"
	case strings.Contains(ua, "Linux"):
		os = "Linux"
	default:
		os = "Other"
	}
	return
}

// parseSource 解析流量来源（SEO/广告归因）
func parseSource(referer string, utmSource, utmMedium, utmCampaign string) (source, medium, keyword string) {
	source = "direct"
	medium = "none"

	if utmSource != "" {
		source = utmSource
	} else if referer != "" {
		switch {
		case strings.Contains(referer, "google."):
			source = "google"
			medium = "organic"
		case strings.Contains(referer, "bing.com"):
			source = "bing"
			medium = "organic"
		case strings.Contains(referer, "baidu.com"):
			source = "baidu"
			medium = "organic"
		case strings.Contains(referer, "facebook.com"):
			source = "facebook"
			medium = "social"
		case strings.Contains(referer, "linkedin.com"):
			source = "linkedin"
			medium = "social"
		case strings.Contains(referer, "twitter.com") || strings.Contains(referer, "x.com"):
			source = "twitter"
			medium = "social"
		case strings.Contains(referer, "youtube.com"):
			source = "youtube"
			medium = "social"
		default:
			source = "referral"
			medium = "referral"
		}
	}

	if utmMedium != "" {
		medium = utmMedium
	}
	return source, medium, keyword
}

// Analytics 访问监测中间件：记录门户（公开）接口的访问数据到 DB + 访问日志
// 仅统计门户流量用于商业价值分析：热门页面、热门产品、来源归因、国家分布、设备分布
// 后台流量没有商业价值，不在此记录。请将此中间件注册到 /api/v1/public 路由组。
func Analytics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录 GET 请求的页面浏览行为
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start).Milliseconds()
		ip := GetClientIP(c)
		ua := GetUserAgent(c)
		device, browser, osName := ParseUserAgent(ua)
		referer := c.GetHeader("Referer")
		source, medium, _ := parseSource(
			referer,
			c.Query("utm_source"),
			c.Query("utm_medium"),
			c.Query("utm_campaign"),
		)

		// 解析实体类型和 slug
		entityType, entitySlug := parseEntity(path)

		log := models.VisitLog{
			CreatedAt:   time.Now(),
			IP:          ip,
			Device:      device,
			Browser:     browser,
			OS:          osName,
			UserAgent:   ua,
			Method:      c.Request.Method,
			Path:        path,
			EntityType:  entityType,
			EntitySlug:  entitySlug,
			Status:      c.Writer.Status(),
			LatencyMs:   latency,
			Referer:     referer,
			Source:      source,
			Medium:      medium,
			UtmCampaign: c.Query("utm_campaign"),
			SessionID:   c.GetString("request_id"),
		}

		// 异步写入 DB（不阻塞响应）
		RecordVisitLog(log)

		// 写入访问日志文件
		if utils.AccessLog != nil {
			utils.AccessLog.Infow("access",
				"ip", ip,
				"method", c.Request.Method,
				"path", path,
				"status", c.Writer.Status(),
				"latency_ms", latency,
				"device", device,
				"browser", browser,
				"os", osName,
				"source", source,
				"medium", medium,
				"referer", referer,
				"entity_type", entityType,
				"entity_slug", entitySlug,
			)
		}
	}
}

// RecordVisitLog 异步写入访问日志（导出，供其他包记录点击事件）
func RecordVisitLog(log models.VisitLog) {
	select {
	case visitLogChan <- log:
	default: // 队列满则丢弃，保护主流程
	}
}

// parseEntity 从路径解析实体类型和 slug
func parseEntity(path string) (entityType, entitySlug string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// /api/v1/public/{type}/{slug?}
	if len(parts) < 4 {
		return "", ""
	}
	entityType = parts[3]
	if len(parts) >= 5 {
		entitySlug = parts[4]
	}
	switch entityType {
	case "products":
		entityType = "product"
	case "blogs":
		entityType = "blog"
	case "pages":
		entityType = "page"
	case "cases":
		entityType = "case"
	case "home":
		entityType = "home"
	}
	return entityType, entitySlug
}
