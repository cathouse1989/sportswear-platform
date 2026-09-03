package middleware

import (
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sportswear-backend/internal/database"
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
			// 按月分表写入：visit_logs_YYYYMM（YYYYMM=月份）
			if log.CreatedAt.IsZero() {
				log.CreatedAt = time.Now()
			}
			tableName := database.MonthTableName("visit_logs", log.CreatedAt)
			if err := database.EnsureTable(db, tableName, &models.VisitLog{}); err == nil {
				db.Table(tableName).Create(&log)
			}
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
	tabletRegex = regexp.MustCompile(`(?i)(ipad|tablet)`)
	modelRegex  = regexp.MustCompile(`(?i)\((iphone ipro|iphone|ipad|ipod touch|macintosh|windows nt [\d.]+|linux|sm-[a-z0-9]+|pixel [\d]+|redmi[ a-z0-9]*|mi [a-z0-9]+|huawei[ a-z0-9]*|oppo[ a-z0-9]*|vivo[ a-z0-9]*|galaxy[ a-z0-9]*|oneplus[ a-z0-9]*)`)
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
	if tabletRegex.MatchString(ua) {
		device = "tablet"
	} else if mobileRegex.MatchString(ua) {
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

// ParseDeviceModel 解析设备型号（移动端 / PC / Pad 的可读名称）
func ParseDeviceModel(ua string) string {
	if ua == "" {
		return ""
	}
	upper := strings.ToLower(ua)
	switch {
	case strings.Contains(upper, "iphone"):
		return "iPhone"
	case strings.Contains(upper, "ipad"):
		return "iPad"
	case strings.Contains(upper, "ipod"):
		return "iPod"
	case strings.Contains(upper, "windows nt 11"):
		return "Windows 11 PC"
	case strings.Contains(upper, "windows nt 10"):
		return "Windows 10 PC"
	case strings.Contains(upper, "windows"):
		return "Windows PC"
	case strings.Contains(upper, "macintosh") || strings.Contains(upper, "mac os"):
		return "Mac"
	case strings.Contains(upper, "linux"):
		return "Linux Device"
	case strings.Contains(upper, "android"):
		// 尝试从 UA 提取型号：如 SM-S928B、Pixel 8、Redmi Note 13
		if m := modelRegex.FindStringSubmatch(ua); len(m) > 1 {
			return strings.TrimSpace(m[1])
		}
		return "Android Device"
	}
	return ""
}

// ResolveCountry 解析访客国家（优先网关/反代头，其次查询参数）
func ResolveCountry(c *gin.Context) string {
	for _, h := range []string{"CF-IPCountry", "X-Country", "X-Geo-Country"} {
		if v := c.GetHeader(h); v != "" {
			return v
		}
	}
	return c.Query("country")
}

// FirstNonEmptyHeader 返回 query 值与请求头序列中第一个非空值
func FirstNonEmptyHeader(c *gin.Context, queryVal string, headers ...string) string {
	if queryVal != "" {
		return queryVal
	}
	for _, hname := range headers {
		if v := c.GetHeader(hname); v != "" {
			return v
		}
	}
	return ""
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
		case strings.Contains(referer, "instagram.com"):
			source = "instagram"
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

// firstNonEmpty 返回第一个非空值
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// Analytics 访问监测中间件：记录门户（公开）接口的访问数据到 DB + 访问日志
// 仅统计门户流量用于商业价值分析：热门页面、热门产品、来源归因、国家/设备/语言分布
// 后台流量没有商业价值，不在此记录。请将此中间件注册到 /api/v1/public 路由组。
// 数据按月分表写入 visit_logs_YYYYMM。
func Analytics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只记录 GET 请求的页面浏览行为（POST 由各 handler 自行记录，如 click-track）
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		// 隐私合规检查：仅当用户同意 Analytics Cookie 时才记录包含个人数据的完整日志
		consentAnalytics := c.GetHeader("X-Consent-Analytics")
		hasConsent := consentAnalytics == "true"

		// 访客唯一标识（前端生成，存 localStorage）
		visitorID := c.GetHeader("X-Visitor-ID")
		if visitorID == "" {
			visitorID = "unknown"
		}

		latency := time.Since(start).Milliseconds()
		ip := GetClientIP(c)
		ua := GetUserAgent(c)
		device, browser, osName := ParseUserAgent(ua)
		referer := c.GetHeader("Referer")

		// UTM 参数：优先取 URL Query，其次取浏览器/Nuxt SSR 透传的自定义头
		utmSource := firstNonEmpty(c.Query("utm_source"), c.GetHeader("X-UTM-Source"), c.GetHeader("UTM-Source"))
		utmMedium := firstNonEmpty(c.Query("utm_medium"), c.GetHeader("X-UTM-Medium"), c.GetHeader("UTM-Medium"))
		utmCampaign := firstNonEmpty(c.Query("utm_campaign"), c.GetHeader("X-UTM-Campaign"), c.GetHeader("UTM-Campaign"))
		utmContent := firstNonEmpty(c.Query("utm_content"), c.GetHeader("X-UTM-Content"), c.GetHeader("UTM-Content"))
		utmTerm := firstNonEmpty(c.Query("utm_term"), c.GetHeader("X-UTM-Term"), c.GetHeader("UTM-Term"))

		source, medium, _ := parseSource(referer, utmSource, utmMedium, utmCampaign)

		// 解析实体类型和 slug
		entityType, entitySlug := parseEntity(path)

		// 判断访问类型：实体内容请求 → page_view，辅助 API → api_call
		visitType := classifyVisit(path)

		// 未同意时匿名化个人数据（GDPR 合规）
		anonymizedIP := ip
		anonymizedUA := ua
		anonymizedReferer := referer
		anonymizedVisitorID := visitorID
		if !hasConsent {
			anonymizedIP = AnonymizeIP(ip)
			anonymizedUA = "anonymous"
			anonymizedReferer = ""
			anonymizedVisitorID = "anonymous"
		}

		consentStatus := "granted"
		if !hasConsent {
			consentStatus = "denied"
		}

		log := models.VisitLog{
			CreatedAt:     time.Now(),
			VisitType:     visitType,
			VisitorID:     anonymizedVisitorID,
			IP:            anonymizedIP,
			Country:       ResolveCountry(c),
			Language:      GetLang(c),
			Device:        device,
			DeviceModel:   ParseDeviceModel(anonymizedUA),
			Browser:       browser,
			OS:            osName,
			UserAgent:     anonymizedUA,
			Method:        c.Request.Method,
			Path:          path,
			EntityType:    entityType,
			EntitySlug:    entitySlug,
			EntityID:      c.GetString("visit_entity_id"),
			EntityName:    c.GetString("visit_entity_name"),
			Status:        c.Writer.Status(),
			LatencyMs:     latency,
			Referer:       anonymizedReferer,
			Source:        source,
			Medium:        medium,
			UtmSource:     utmSource,
			UtmMedium:     utmMedium,
			UtmCampaign:   utmCampaign,
			UtmContent:    utmContent,
			UtmTerm:       utmTerm,
			SessionID:     c.GetString("request_id"),
			ConsentStatus: consentStatus,
		}

		// 异步写入 DB（不阻塞响应）
		RecordVisitLog(log)

		// 写入访问日志文件
		if utils.AccessLog != nil {
			utils.AccessLog.Infow("access",
				"ip", anonymizedIP,
				"country", log.Country,
				"lang", log.Language,
				"method", c.Request.Method,
				"path", path,
				"status", c.Writer.Status(),
				"latency_ms", latency,
				"device", device,
				"device_model", log.DeviceModel,
				"browser", browser,
				"os", osName,
				"source", source,
				"medium", medium,
				"referer", anonymizedReferer,
				"utm_source", utmSource,
				"utm_campaign", utmCampaign,
				"entity_type", entityType,
				"entity_slug", entitySlug,
				"entity_name", log.EntityName,
				"consent", hasConsent,
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

// classifyVisit 判断访问类型
// page_view：实体内容请求（代表用户浏览了一个页面）
// api_call：辅助 API 请求（导航、主题、语言等，不代表页面浏览）
func classifyVisit(path string) string {
	// 实体内容路由前缀 → page_view
	pageViewPrefixes := []string{
		"/api/v1/public/products",
		"/api/v1/public/blogs",
		"/api/v1/public/pages",
		"/api/v1/public/cases",
		"/api/v1/public/home",
		"/api/v1/public/about",
		"/api/v1/public/faq",
		"/api/v1/public/contact",
		"/api/v1/public/categories",
		"/api/v1/public/series",
		"/api/v1/public/fabrics",
		"/api/v1/public/factories",
		"/api/v1/public/certifications",
	}
	for _, prefix := range pageViewPrefixes {
		if strings.HasPrefix(path, prefix) {
			return "page_view"
		}
	}
	// 其他 API（navigations, theme, currencies, languages, geo, i18n 等）→ api_call
	return "api_call"
}

// AnonymizeIP 匿名化 IP 地址（GDPR 合规）
// IPv4: 将最后一组替换为 0（如 192.168.1.42 → 192.168.1.0）
// IPv6: 将后缀 64 位清零
func AnonymizeIP(ip string) string {
	if ip == "" || ip == "unknown" {
		return "anonymous"
	}
	// IPv4 处理
	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			parts[3] = "0"
			return strings.Join(parts, ".")
		}
	}
	// IPv6 处理：保留前 48 位，后续清零
	if strings.Contains(ip, ":") {
		// 简化的 IPv6 匿名化：保留前 3 组
		parts := strings.Split(ip, ":")
		if len(parts) >= 3 {
			return strings.Join(parts[:3], ":") + "::/0"
		}
	}
	return "anonymous"
}
