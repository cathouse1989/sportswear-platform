package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sportswear-backend/internal/config"
	"sportswear-backend/internal/handlers"
	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// Setup 设置路由
func Setup(cfg *config.Config, db *gorm.DB, cacheService *services.CacheService) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.Timezone()) // 时区解析：从 X-Timezone 头或 tz 参数解析访客时区
	r.Use(middleware.Language()) // 语言解析：从 lang 参数或 Accept-Language 头解析客户端语言
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 上传文件读取路由（/uploads/**）在下方 mediaHandler 初始化后注册，
	// 由 ServeUpload 按存储驱动（本地磁盘 / MinIO）流式读取。

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 时间信息（全球网站：前端获取服务器 UTC 时间和时区列表）
	r.GET("/api/v1/time", func(c *gin.Context) {
		tz := middleware.GetTimezone(c)
		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"utc_now":          utils.FormatRFC3339(utils.NowUTC()),
				"local_now":        utils.FormatRFC3339In(utils.NowUTC(), tz),
				"timezone":         tz,
				"common_timezones": utils.CommonTimezones,
			},
			"requestId": c.GetString("request_id"),
		})
	})

	// 初始化服务
	authService := services.NewAuthService(db, cfg)
	productService := services.NewProductService(db)
	cmsService := services.NewCMSService(db)
	leadService := services.NewLeadService(db)
	// IP 地理定位（基于数据库表 ip_geo_ranges，数据由后台「IP 地理库」接口导入维护）
	geoIPService := services.NewGeoIPService(db)
	leadService.SetGeoIP(geoIPService)
	subscriptionService := services.NewSubscriptionService(db)
	mediaService := services.NewMediaService(db)
	trashService := services.NewTrashService(db)
	analyticsService := services.NewAnalyticsService(db)
	portalService := services.NewPortalService(db, cacheService)
	i18nService := services.NewI18nService(db)
	enumService := services.NewEnumService(db)
	seoService := services.NewSEOService(db)
	currencyService := services.NewCurrencyService(db)
	geoService := services.NewGeoService(db)
	storageSourceService := services.NewStorageSourceService(db)
	uploadService := services.NewUploadService(cfg)

	// 确保对象存储桶存在（minio 驱动下，幂等）；失败仅告警，上传时再暴露错误
	if err := uploadService.EnsureBucket(); err != nil {
		utils.Logger.Warnw("初始化对象存储桶失败", "driver", uploadService.Driver(), "error", err.Error())
	}

	unitService := services.NewUnitService()
	notificationService := services.NewNotificationService(db)
	operationLogService := services.NewOperationLogService(db)
	quoteService := services.NewQuoteService(db)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService, cacheService)
	cmsHandler := handlers.NewCMSHandler(cmsService, cacheService)
	leadHandler := handlers.NewLeadHandler(leadService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	mediaHandler := handlers.NewMediaHandler(mediaService, uploadService)

	// 上传文件读取路由（/uploads/**）：local 读本地磁盘，minio 读对象存储，
	// 与上传时落库的 url（/uploads/<path>）一致，形成展示闭环。
	r.GET("/uploads/*filepath", mediaHandler.ServeUpload)

	publicHandler := handlers.NewPublicHandler(productService, cmsService, leadService, cacheService, uploadService)
	trashHandler := handlers.NewTrashHandler(trashService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	portalHandler := handlers.NewPortalHandler(portalService, cacheService, db)
	i18nHandler := handlers.NewI18nHandler(i18nService)
	enumHandler := handlers.NewEnumHandler(enumService)
	localizationHandler := handlers.NewLocalizationHandler(seoService, currencyService, geoService, unitService)
	ipGeoHandler := handlers.NewIPGeoHandler(geoIPService)
	storageSourceHandler := handlers.NewStorageSourceHandler(storageSourceService)
	systemHandler := handlers.NewSystemHandler(notificationService, operationLogService, quoteService, cmsService, productService)

	// API 路由
	api := r.Group("/api/v1")

	// ==================== 公开 API（无需认证，带频率限制防滥用） ====================
	// 注意：公开接口只返回已发布（published）内容，草稿/下线内容不可见
	public := api.Group("/public")
	public.Use(middleware.RateLimit()) // IP 频率限制：每分钟 3000 次/路径
	public.Use(middleware.Analytics()) // 访问监测：仅统计门户（门户）流量，后台流量无商业价值
	{
		public.GET("/home", publicHandler.GetHome)
		public.GET("/pages", publicHandler.ListLandingPages)
		public.GET("/pages/:slug", publicHandler.GetPage)
		public.GET("/products", publicHandler.ListProducts)
		public.GET("/products/:slug", publicHandler.GetProduct)
		public.GET("/categories", publicHandler.ListCategories)
		public.GET("/series", publicHandler.ListSeries)
		public.GET("/blogs", publicHandler.ListBlogs)
		public.GET("/blogs/:slug", publicHandler.GetBlog)
		public.GET("/cases", publicHandler.ListCases)
		public.GET("/cases/:slug", publicHandler.GetCase)
		public.GET("/faqs", publicHandler.ListFAQs)
		public.GET("/fabrics", publicHandler.ListFabrics)
		public.GET("/factories", publicHandler.ListFactories)
		public.GET("/certifications", publicHandler.ListCertifications)
		public.GET("/production-processes", publicHandler.ListProductionProcesses)
		public.GET("/self-medias", publicHandler.ListSelfMedias)
		public.GET("/navigations", publicHandler.ListNavigations)
		public.POST("/leads", publicHandler.CreateLead)                             // 询盘提交（受频率限制保护）
		public.POST("/subscribe", subscriptionHandler.Subscribe)                    // 订阅更新（Newsletter，受频率限制保护）
		public.POST("/uploads/lead-attachment", publicHandler.UploadLeadAttachment) // 询盘附件上传（图片/文档/压缩包，默认≤20MB，受频率限制保护）
		public.POST("/click-track", publicHandler.TrackClick)                       // 外部链接点击跟踪（社交媒体跳转）

		// 隐私合规 API（GDPR / PIPL / CCPA / LGPD）
		public.POST("/data-export", publicHandler.RequestDataExport)     // 数据主体请求导出个人数据
		public.POST("/data-deletion", publicHandler.RequestDataDeletion) // 数据主体请求删除个人数据（被遗忘权）

		public.GET("/theme", portalHandler.GetThemeConfig)                     // 主题配置（前端渲染视觉风格）
		public.GET("/i18n", i18nHandler.GetDictionary)                         // i18n 词条字典（前端一键切换语言）
		public.GET("/enums", enumHandler.GetPublicEnums)                       // 枚举字典（数据字典，前端枚举标签翻译）
		public.GET("/currencies", localizationHandler.ListCurrencies)          // 货币列表
		public.GET("/currencies/convert", localizationHandler.ConvertCurrency) // 货币转换
		public.GET("/locale", localizationHandler.GetLocaleInfo)               // 本地化信息（语言+货币+时区）
		public.GET("/geo", localizationHandler.GetGeolocation)                 // 地理定位（IP→国家→语言/货币/时区）
		public.GET("/seo", localizationHandler.GetRouteSEO)                    // 路由 SEO（门户列表页/落地页）
		public.GET("/languages", localizationHandler.GetLanguages)             // 启用语言列表
		public.GET("/media/:id", mediaHandler.GetMedia)                        // 媒体公开访问（通过 ID 获取媒体信息）
		public.GET("/units/sizes", localizationHandler.GetSizeChart)           // 尺寸对照表
		public.GET("/units/convert-size", localizationHandler.ConvertSize)     // 尺码转换
		public.GET("/units/convert-weight", localizationHandler.ConvertWeight) // 重量单位转换
	}

	// sitemap（公开，SEO 多语言）
	r.GET("/sitemap.xml", systemHandler.GetSitemap)

	// ==================== 后台 API（JWT 认证 + RBAC 权限硬控） ====================
	admin := api.Group("/admin")
	{
		// 登录接口（无需认证，自带频率限制）
		admin.POST("/auth/login", middleware.RateLimit(), authHandler.Login)
		// 登出接口（无需认证：仅清除会话 Cookie，前端同时清理本地 token）
		admin.POST("/auth/logout", middleware.RateLimit(), authHandler.Logout)

		// 需要认证的路由：验证 JWT + 从数据库加载角色权限
		auth := admin.Group("")
		auth.Use(middleware.Auth(cfg, authService))
		auth.Use(middleware.Audit(operationLogService)) // 审计：记录后台写操作
		{
			// 当前用户信息（所有登录用户可访问）
			auth.GET("/auth/profile", authHandler.GetProfile)

			// ---------- 用户管理 ----------
			auth.GET("/users",
				middleware.RequirePermission("user:view"), authHandler.ListUsers)
			auth.POST("/users",
				middleware.RequirePermission("user:create"), authHandler.CreateUser)
			auth.PUT("/users/:id",
				middleware.RequirePermission("user:update"), authHandler.UpdateUser)
			auth.DELETE("/users/:id",
				middleware.RequirePermission("user:delete"), authHandler.DeleteUser)
			auth.GET("/roles",
				middleware.RequireAnyPermission("role:manage", "user:create", "user:update"), authHandler.ListRoles)
			auth.GET("/roles/page",
				middleware.RequireAnyPermission("role:manage", "user:create", "user:update"), authHandler.ListRolesPage)
			auth.POST("/roles",
				middleware.RequirePermission("role:manage"), authHandler.CreateRole)
			auth.PUT("/roles/:id",
				middleware.RequirePermission("role:manage"), authHandler.UpdateRole)
			auth.PUT("/roles/:id/status",
				middleware.RequirePermission("role:manage"), authHandler.UpdateRoleStatus)
			auth.DELETE("/roles/:id",
				middleware.RequirePermission("role:manage"), authHandler.DeleteRole)
			auth.GET("/permissions",
				middleware.RequirePermission("role:manage"), authHandler.ListPermissions)

			// ---------- 应用维度管理 ----------
			auth.GET("/apps",
				middleware.RequirePermission("setting:manage"), authHandler.ListApps)
			auth.POST("/apps",
				middleware.RequirePermission("setting:manage"), authHandler.CreateApp)
			auth.PUT("/apps/:id/status",
				middleware.RequirePermission("setting:manage"), authHandler.UpdateAppStatus)
			auth.POST("/apps/:id/users",
				middleware.RequirePermission("setting:manage"), authHandler.AssignUserToApp)
			auth.DELETE("/apps/:id/users",
				middleware.RequirePermission("setting:manage"), authHandler.RemoveUserFromApp)

			// ---------- 产品管理 ----------
			auth.GET("/products/stats",
				middleware.RequirePermission("product:view"), productHandler.GetProductStats)
			auth.GET("/products",
				middleware.RequirePermission("product:view"), productHandler.ListProducts)
			auth.GET("/products/:id",
				middleware.RequirePermission("product:view"), productHandler.GetProduct)
			auth.POST("/products",
				middleware.RequirePermission("product:create"), productHandler.CreateProduct)
			auth.PUT("/products/:id",
				middleware.RequirePermission("product:update"), productHandler.UpdateProduct)
			auth.DELETE("/products/:id",
				middleware.RequirePermission("product:delete"), productHandler.DeleteProduct)
			auth.POST("/products/:id/publish",
				middleware.RequirePermission("product:publish"), productHandler.PublishProduct)
			auth.POST("/products/:id/unpublish",
				middleware.RequirePermission("product:publish"), productHandler.UnpublishProduct)

			// 分类管理
			auth.GET("/categories",
				middleware.RequireAnyPermission("category:manage", "product:view"), productHandler.ListCategories)
			auth.POST("/categories",
				middleware.RequirePermission("category:manage"), productHandler.CreateCategory)
			auth.PUT("/categories/:id",
				middleware.RequirePermission("category:manage"), productHandler.UpdateCategory)
			auth.DELETE("/categories/:id",
				middleware.RequirePermission("category:manage"), productHandler.DeleteCategory)

			// 系列管理
			auth.GET("/series",
				middleware.RequireAnyPermission("series:manage", "product:view"), productHandler.ListSeries)
			auth.POST("/series",
				middleware.RequirePermission("series:manage"), productHandler.CreateSeries)
			auth.POST("/series/:id/publish",
				middleware.RequirePermission("product:publish"), productHandler.PublishSeries)
			auth.POST("/series/:id/unpublish",
				middleware.RequirePermission("product:publish"), productHandler.UnpublishSeries)

			// 面料管理
			auth.GET("/fabrics",
				middleware.RequireAnyPermission("fabric:manage", "product:view"), productHandler.ListFabrics)
			auth.POST("/fabrics",
				middleware.RequirePermission("fabric:manage"), productHandler.CreateFabric)
			auth.POST("/fabrics/:id/publish",
				middleware.RequirePermission("product:publish"), productHandler.PublishFabric)
			auth.POST("/fabrics/:id/unpublish",
				middleware.RequirePermission("product:publish"), productHandler.UnpublishFabric)

			// ---------- CMS 页面管理 ----------
			auth.GET("/pages",
				middleware.RequirePermission("page:view"), cmsHandler.ListPages)
			auth.GET("/pages/:id",
				middleware.RequirePermission("page:view"), cmsHandler.GetPage)
			auth.PUT("/pages/:id",
				middleware.RequirePermission("page:update"), cmsHandler.UpdatePage)
			auth.PUT("/pages/:id/hero",
				middleware.RequirePermission("page:update"), cmsHandler.UpdateHeroSlides)
			auth.POST("/pages/:id/publish",
				middleware.RequirePermission("page:publish"), cmsHandler.PublishPage)
			auth.POST("/pages/:id/unpublish",
				middleware.RequirePermission("page:publish"), cmsHandler.UnpublishPage)

			// 导航管理（结构由代码/seed 维护，仅开放编辑/排序/显隐，不开放增删）
			auth.GET("/navigations",
				middleware.RequireAnyPermission("navigation:manage", "page:view"), cmsHandler.ListNavigations)
			auth.PUT("/navigations/:id",
				middleware.RequirePermission("navigation:manage"), cmsHandler.UpdateNavigation)

			auth.PUT("/navigations/sort",
				middleware.RequirePermission("navigation:manage"), cmsHandler.BatchSortNavigations)

		// 门户路由注册表 & 健康检查（供 SEO 管理、导航管理、健康检查页共用）
		auth.GET("/portal/routes",
			middleware.RequireAnyPermission("navigation:manage", "seo:manage", "page:view"), cmsHandler.ListPortalRoutes)
		auth.GET("/portal/health",
			middleware.RequireAnyPermission("navigation:manage", "seo:manage", "page:view"), cmsHandler.PortalHealthCheck)

			// 博客管理
			auth.GET("/blogs",
				middleware.RequireAnyPermission("blog:manage", "page:view"), cmsHandler.ListBlogs)
			auth.GET("/blogs/:id",
				middleware.RequireAnyPermission("blog:manage", "page:view"), cmsHandler.GetBlog)
			auth.POST("/blogs",
				middleware.RequirePermission("blog:manage"), cmsHandler.CreateBlog)
			auth.PUT("/blogs/:id",
				middleware.RequirePermission("blog:manage"), cmsHandler.UpdateBlog)
			auth.DELETE("/blogs/:id",
				middleware.RequirePermission("blog:manage"), cmsHandler.DeleteBlog)
			auth.POST("/blogs/:id/publish",
				middleware.RequirePermission("blog:manage"), cmsHandler.PublishBlog)
			auth.POST("/blogs/:id/unpublish",
				middleware.RequirePermission("blog:manage"), cmsHandler.UnpublishBlog)

			// 案例管理
			auth.GET("/cases",
				middleware.RequireAnyPermission("case:manage", "page:view"), cmsHandler.ListCases)
			auth.GET("/cases/:id",
				middleware.RequireAnyPermission("case:manage", "page:view"), cmsHandler.GetCase)
			auth.POST("/cases",
				middleware.RequirePermission("case:manage"), cmsHandler.CreateCase)
			auth.PUT("/cases/:id",
				middleware.RequirePermission("case:manage"), cmsHandler.UpdateCase)
			auth.DELETE("/cases/:id",
				middleware.RequirePermission("case:manage"), cmsHandler.DeleteCase)
			auth.POST("/cases/:id/publish",
				middleware.RequirePermission("case:manage"), cmsHandler.PublishCase)
			auth.POST("/cases/:id/unpublish",
				middleware.RequirePermission("case:manage"), cmsHandler.UnpublishCase)

			// FAQ 管理
			auth.GET("/faqs",
				middleware.RequireAnyPermission("faq:manage", "page:view"), cmsHandler.ListFAQs)
			auth.POST("/faqs",
				middleware.RequirePermission("faq:manage"), cmsHandler.CreateFAQ)
			auth.PUT("/faqs/:id",
				middleware.RequirePermission("faq:manage"), cmsHandler.UpdateFAQ)
			auth.DELETE("/faqs/:id",
				middleware.RequirePermission("faq:manage"), cmsHandler.DeleteFAQ)

			// 工厂管理
			auth.GET("/factories",
				middleware.RequireAnyPermission("factory:manage", "page:view"), cmsHandler.ListFactories)
			auth.POST("/factories",
				middleware.RequirePermission("factory:manage"), cmsHandler.CreateFactory)
			auth.PUT("/factories/:id",
				middleware.RequirePermission("factory:manage"), cmsHandler.UpdateFactory)
			auth.DELETE("/factories/:id",
				middleware.RequirePermission("factory:manage"), cmsHandler.DeleteFactory)
			auth.POST("/factories/:id/publish",
				middleware.RequirePermission("factory:manage"), cmsHandler.PublishFactory)
			auth.POST("/factories/:id/unpublish",
				middleware.RequirePermission("factory:manage"), cmsHandler.UnpublishFactory)

			// 认证管理
			auth.GET("/certifications",
				middleware.RequireAnyPermission("certification:manage", "page:view"), cmsHandler.ListCertifications)
			auth.POST("/certifications",
				middleware.RequirePermission("certification:manage"), cmsHandler.CreateCertification)
			auth.PUT("/certifications/:id",
				middleware.RequirePermission("certification:manage"), cmsHandler.UpdateCertification)
			auth.DELETE("/certifications/:id",
				middleware.RequirePermission("certification:manage"), cmsHandler.DeleteCertification)
			auth.POST("/certifications/:id/publish",
				middleware.RequirePermission("certification:manage"), cmsHandler.PublishCertification)
			auth.POST("/certifications/:id/unpublish",
				middleware.RequirePermission("certification:manage"), cmsHandler.UnpublishCertification)

			// 生产流程管理
			auth.GET("/production-processes",
				middleware.RequireAnyPermission("production:manage", "page:view"), cmsHandler.ListProductionProcesses)
			auth.POST("/production-processes",
				middleware.RequirePermission("production:manage"), cmsHandler.CreateProductionProcess)
			auth.PUT("/production-processes/:id",
				middleware.RequirePermission("production:manage"), cmsHandler.UpdateProductionProcess)
			auth.DELETE("/production-processes/:id",
				middleware.RequirePermission("production:manage"), cmsHandler.DeleteProductionProcess)
			auth.POST("/production-processes/:id/publish",
				middleware.RequirePermission("production:manage"), cmsHandler.PublishProductionProcess)
			auth.POST("/production-processes/:id/unpublish",
				middleware.RequirePermission("production:manage"), cmsHandler.UnpublishProductionProcess)

			// 自媒体管理
			auth.GET("/self-medias",
				middleware.RequireAnyPermission("selfmedia:manage", "page:view"), cmsHandler.ListSelfMedias)
			auth.GET("/self-medias/config",
				middleware.RequireAnyPermission("selfmedia:manage", "page:view"), cmsHandler.GetSelfMediaConfig)
			auth.PUT("/self-medias/config",
				middleware.RequirePermission("selfmedia:manage"), cmsHandler.UpdateSelfMediaConfig)
			auth.POST("/self-medias",
				middleware.RequirePermission("selfmedia:manage"), cmsHandler.CreateSelfMedia)
			auth.PUT("/self-medias/:id",
				middleware.RequirePermission("selfmedia:manage"), cmsHandler.UpdateSelfMedia)
			auth.DELETE("/self-medias/:id",
				middleware.RequirePermission("selfmedia:manage"), cmsHandler.DeleteSelfMedia)
			auth.POST("/self-medias/:id/publish",
				middleware.RequirePermission("selfmedia:manage"), cmsHandler.PublishSelfMedia)
			auth.POST("/self-medias/:id/unpublish",
				middleware.RequirePermission("selfmedia:manage"), cmsHandler.UnpublishSelfMedia)

			// ---------- 媒体管理 ----------
			auth.GET("/media",
				middleware.RequireAnyPermission("media:manage", "media:upload"), mediaHandler.ListMedia)
			auth.GET("/media/:id",
				middleware.RequireAnyPermission("media:manage", "media:upload"), mediaHandler.GetMedia)
			auth.POST("/media",
				middleware.RequirePermission("media:upload"), mediaHandler.CreateMedia)
			auth.POST("/media/upload",
				middleware.RequirePermission("media:upload"), mediaHandler.UploadFile)
			auth.PUT("/media/:id",
				middleware.RequirePermission("media:manage"), mediaHandler.UpdateMedia)
			auth.DELETE("/media/:id",
				middleware.RequirePermission("media:manage"), mediaHandler.DeleteMedia)

			// ---------- 询盘管理 ----------
			auth.GET("/leads",
				middleware.RequirePermission("lead:view"), leadHandler.ListLeads)
			auth.GET("/leads/:id",
				middleware.RequirePermission("lead:view"), leadHandler.GetLead)
			auth.PUT("/leads/:id",
				middleware.RequireAnyPermission("lead:update", "lead:followup"), leadHandler.UpdateLead)
			auth.DELETE("/leads/:id",
				middleware.RequirePermission("lead:update"), leadHandler.DeleteLead)
			auth.POST("/leads/:id/followups",
				middleware.RequirePermission("lead:followup"), leadHandler.AddFollowUp)
			// 回填历史询盘的 IP 国家（补齐 ip_country，依赖 GEOIP_DB_PATH 离线库）
			auth.POST("/leads/backfill-ip-country",
				middleware.RequirePermission("lead:update"), leadHandler.BackfillIPCountries)

			// ---------- 订阅管理（门户"订阅更新"） ----------
			auth.GET("/subscribers/stats",
				middleware.RequirePermission("subscription:view"), subscriptionHandler.GetStats)
			auth.GET("/subscribers",
				middleware.RequirePermission("subscription:view"), subscriptionHandler.ListSubscribers)
			auth.PUT("/subscribers/:id",
				middleware.RequirePermission("subscription:update"), subscriptionHandler.UpdateSubscriber)
			auth.DELETE("/subscribers/:id",
				middleware.RequirePermission("subscription:update"), subscriptionHandler.DeleteSubscriber)

			// ---------- 数据统计 ----------
			auth.GET("/dashboard",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), leadHandler.GetDashboardStats)

			// ---------- 回收站（软删除管理） ----------
			auth.GET("/trash",
				middleware.RequirePermission("setting:manage"), trashHandler.ListTrash)
			auth.POST("/trash/:id/restore",
				middleware.RequirePermission("setting:manage"), trashHandler.RestoreItem)
			auth.POST("/trash/:id/purge",
				middleware.RequirePermission("setting:manage"), trashHandler.PurgeItem)
			auth.DELETE("/trash",
				middleware.RequirePermission("setting:manage"), trashHandler.EmptyTrash)

			// ---------- 流量分析（商业价值监测） ----------
			auth.GET("/analytics/overview",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetTrafficOverview)
			auth.GET("/analytics/top-pages",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetTopPages)
			auth.GET("/analytics/top-products",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetTopProducts)
			auth.GET("/analytics/sources",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetSourceAnalysis)
			auth.GET("/analytics/countries",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetCountryAnalysis)
			auth.GET("/analytics/devices",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetDeviceAnalysis)
			auth.GET("/analytics/social-clicks",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetSocialClickAnalysis)
			// 广告归因（utm_campaign 维度，衡量各广告 Campaign 转化来源）
			auth.GET("/analytics/utm-campaigns",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetUtmCampaignAnalysis)
			// 门户访问日志明细（IP/国家/来源/实体/时间范围筛选，跨月度分表组合查询）
			auth.GET("/analytics/visit-logs",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.ListVisitLogs)
			// 访客旅程（某个访客的完整访问路径）
			auth.GET("/analytics/visitor-journey",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetVisitorJourney)
			// IP-询盘关联（某个 IP 提交的询盘列表）
			auth.GET("/analytics/ip-leads",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetIPLeads)
			// IP 访问记录（某个 IP 的访问概览 + 访问明细：现在/过去/历史）
			auth.GET("/analytics/ip-visits",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetIPVisits)
			// 转化漏斗（访问 → 产品浏览 → 询盘）
			auth.GET("/analytics/conversion-funnel",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetConversionFunnel)
			// 隐私合规洞察（同意/未同意用户的转化对比）
			auth.GET("/analytics/consent-insights",
				middleware.RequireAnyPermission("dashboard:view", "lead:view"), analyticsHandler.GetConsentInsights)

			// ---------- 门户展示配置（版本化发布 + 主题配置） ----------
			// 页面版本：草稿与线上分离，编辑不影响线上，发布后立即生效
			auth.POST("/pages/:id/versions",
				middleware.RequirePermission("page:update"), portalHandler.SaveDraft)
			auth.GET("/pages/:id/versions",
				middleware.RequirePermission("page:view"), portalHandler.ListVersions)
			auth.POST("/versions/:id/publish",
				middleware.RequirePermission("page:publish"), portalHandler.PublishVersion)
			auth.POST("/versions/:id/rollback",
				middleware.RequirePermission("page:publish"), portalHandler.RollbackVersion)

			// 主题配置
			// GET 仅需登录：主题/Logo 属于非敏感展示配置，所有后台用户界面均需使用（如顶栏 Logo）
			auth.GET("/theme", portalHandler.ListThemeConfigs)
			auth.PUT("/theme/:key",
				middleware.RequirePermission("setting:manage"), portalHandler.UpdateThemeConfig)

			// portal cache: status / toggle / refresh / publish
			auth.GET("/portal-cache",
				middleware.RequireAnyPermission("setting:manage", "page:view"), portalHandler.GetPortalCacheStatus)
			auth.PUT("/portal-cache/enabled",
				middleware.RequireAnyPermission("setting:manage", "page:update"), portalHandler.SetPortalCacheEnabled)
			auth.POST("/portal-cache/refresh",
				middleware.RequireAnyPermission("setting:manage", "page:update"), portalHandler.RefreshPortalCache)
			auth.POST("/portal-cache/publish",
				middleware.RequireAnyPermission("setting:manage", "page:publish"), portalHandler.PublishPortalCache)

			// ---------- 国际化词条管理（i18n 词典） ----------
			auth.GET("/i18n/entries",
				middleware.RequirePermission("language:manage"), i18nHandler.ListEntries)
			auth.POST("/i18n/entries",
				middleware.RequirePermission("language:manage"), i18nHandler.UpsertEntry)
			auth.DELETE("/i18n/entries/:id",
				middleware.RequirePermission("language:manage"), i18nHandler.DeleteEntry)

			// ---------- 枚举字典管理（数据字典：值域 + 多语言翻译） ----------
			// 只读接口：登录即可读（枚举标签为通用配置数据，产品/内容等页面均需展示）
			auth.GET("/enums/types", enumHandler.ListTypes)
			auth.GET("/enums/items", enumHandler.ListItems)
			auth.POST("/enums/types",
				middleware.RequirePermission("setting:manage"), enumHandler.UpsertType)
			auth.DELETE("/enums/types/:id",
				middleware.RequirePermission("setting:manage"), enumHandler.DeleteType)
			auth.POST("/enums/items",
				middleware.RequirePermission("setting:manage"), enumHandler.UpsertItem)
			auth.DELETE("/enums/items/:id",
				middleware.RequirePermission("setting:manage"), enumHandler.DeleteItem)

			// ---------- SEO 管理（多语言） ----------
			auth.POST("/seo",
				middleware.RequirePermission("seo:manage"), localizationHandler.UpsertSEO)
			auth.POST("/seo/routes", middleware.RequirePermission("seo:manage"), localizationHandler.UpsertRoutesSEO)
			auth.GET("/seo/routes", middleware.RequirePermission("seo:manage"), localizationHandler.ListRouteSEO)
			auth.POST("/seo/entity", middleware.RequirePermission("seo:manage"), localizationHandler.UpsertEntitySEO)
			auth.GET("/seo/entity", middleware.RequirePermission("seo:manage"), localizationHandler.ListEntitySEO)

			// ---------- 国家本地化映射（区域设置） ----------
			auth.GET("/geo-locales", middleware.RequirePermission("setting:manage"), localizationHandler.ListGeoLocales)
			auth.POST("/geo-locales", middleware.RequirePermission("setting:manage"), localizationHandler.UpsertGeoLocale)
			auth.DELETE("/geo-locales/:id", middleware.RequirePermission("setting:manage"), localizationHandler.DeleteGeoLocale)

			// ---------- IP 地理库（IP 段 → 国家映射，供询盘 IP 国家解析） ----------
			auth.GET("/ip-geo-ranges", middleware.RequirePermission("setting:manage"), ipGeoHandler.ListIPGeoRanges)
			auth.POST("/ip-geo-ranges/import", middleware.RequirePermission("setting:manage"), ipGeoHandler.ImportIPGeoRanges)
			auth.DELETE("/ip-geo-ranges/:id", middleware.RequirePermission("setting:manage"), ipGeoHandler.DeleteIPGeoRange)

			auth.POST("/ai/translate", middleware.RequirePermission("language:manage"), handlers.AITranslate)

			// ---------- 存储源管理 ----------
			auth.GET("/storage-sources",
				middleware.RequirePermission("setting:manage"), storageSourceHandler.List)
			auth.GET("/storage-sources/:id",
				middleware.RequirePermission("setting:manage"), storageSourceHandler.Get)
			auth.POST("/storage-sources",
				middleware.RequirePermission("setting:manage"), storageSourceHandler.Create)
			auth.PUT("/storage-sources/:id",
				middleware.RequirePermission("setting:manage"), storageSourceHandler.Update)
			auth.DELETE("/storage-sources/:id",
				middleware.RequirePermission("setting:manage"), storageSourceHandler.Delete)

			// ---------- 通知系统 ----------
			auth.GET("/notifications",
				middleware.RequirePermission("lead:view"), systemHandler.ListNotifications)
			auth.GET("/notifications/unread-count",
				middleware.RequirePermission("lead:view"), systemHandler.UnreadCount)
			auth.POST("/notifications/:id/read",
				middleware.RequirePermission("lead:view"), systemHandler.MarkRead)
			auth.POST("/notifications/read-all",
				middleware.RequirePermission("lead:view"), systemHandler.MarkAllRead)

			// ---------- 操作日志（审计查询） ----------
			auth.GET("/operation-logs",
				middleware.RequirePermission("setting:manage"), systemHandler.ListOperationLogs)

			// ---------- 报价管理 ----------
			auth.GET("/quotes",
				middleware.RequirePermission("lead:view"), systemHandler.ListQuotes)
			auth.GET("/quotes/:id",
				middleware.RequirePermission("lead:view"), systemHandler.GetQuote)
			auth.POST("/quotes",
				middleware.RequirePermission("lead:update"), systemHandler.CreateQuote)
			auth.PUT("/quotes/:id",
				middleware.RequirePermission("lead:update"), systemHandler.UpdateQuote)
			auth.DELETE("/quotes/:id",
				middleware.RequirePermission("lead:update"), systemHandler.DeleteQuote)
		}
	}

	return r
}
