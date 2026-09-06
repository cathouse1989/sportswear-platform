package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"sportswear-backend/internal/config"
	"sportswear-backend/internal/database"
	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/router"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

func main() {
	// 全局时区统一为 UTC（面向全球网站，存储层统一使用 UTC）
	os.Setenv("TZ", "UTC")
	time.Local = time.UTC

	// 加载 .env 文件
	_ = godotenv.Load()

	// 加载配置
	cfg := config.Load()

	// 初始化分级日志系统（按日期分目录 + 自动轮转清理）
	utils.InitLogger(cfg.Server.Env)
	defer utils.SyncLogs()

	utils.Logger.Infow("服务器启动中",
		"env", cfg.Server.Env,
		"port", cfg.Server.Port,
	)

	// 初始化数据库
	db, err := database.Init(cfg)
	if err != nil {
		utils.Logger.Errorw("数据库初始化失败", "error", err)
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置访问监测的数据库连接
	middleware.SetAnalyticsDB(db)

	// 初始化 Redis 缓存（不可用时自动降级）
	redisClient, redisErr := database.InitRedis(cfg)
	var cacheService *services.CacheService
	if redisErr != nil {
		utils.Logger.Warnw("Redis 不可用，缓存降级", "error", redisErr.Error())
		cacheService = services.NewCacheService(nil)
	} else {
		cacheService = services.NewCacheService(redisClient)
	}

	// 绑定 DB：加载门户缓存业务开关（portal_cache_enabled）
	cacheService.BindDB(db)

	// 初始化默认数据
	authService := services.NewAuthService(db, cfg)
	if err := authService.InitDefaultData(); err != nil {
		utils.Logger.Errorw("初始化默认数据失败", "error", err)
		log.Fatalf("初始化默认数据失败: %v", err)
	}

	// 初始化默认主题配置
	portalService := services.NewPortalService(db, cacheService)
	if err := portalService.InitDefaultTheme(); err != nil {
		utils.Logger.Warnw("初始化默认主题失败", "error", err.Error())
	}

	// 初始化 i18n 默认词条
	i18nService := services.NewI18nService(db)
	if err := i18nService.InitDefaults(); err != nil {
		utils.Logger.Warnw("初始化 i18n 词条失败", "error", err.Error())
	}

	// 初始化枚举字典（数据字典：值域 + 多语言翻译，幂等）
	enumService := services.NewEnumService(db)
	if err := enumService.InitDefaults(); err != nil {
		utils.Logger.Warnw("初始化枚举字典失败", "error", err.Error())
	}

	// 初始化默认货币
	currencyService := services.NewCurrencyService(db)
	if err := currencyService.InitDefaults(); err != nil {
		utils.Logger.Warnw("初始化默认货币失败", "error", err.Error())
	}

	// 初始化国家本地化映射
	geoService := services.NewGeoService(db)
	if err := geoService.InitGeoLocales(); err != nil {
		utils.Logger.Warnw("初始化国家映射失败", "error", err.Error())
	}

	// 初始化默认存储源
	storageSourceService := services.NewStorageSourceService(db)
	if err := storageSourceService.InitDefault(); err != nil {
		utils.Logger.Warnw("初始化默认存储源失败", "error", err.Error())
	}

	// 初始化示例数据（仅首次启动时填充）
	services.Seed(db)

	// 初始化默认自媒体账号（幂等，自媒体管理为唯一配置入口）
	services.SeedSelfMedias(db)

	// 初始化 IP 地理库演示数据（幂等，供询盘 IP 国家解析）
	services.SeedIPGeoRanges(db)

	// 启动数据保留清理任务（每日自动匿名化过期 PII，GDPR/PIPL 合规）
	privacyService := services.NewPrivacyService(db)
	privacyService.StartRetentionSweeper()

	// 设置路由
	r := router.Setup(cfg, db, cacheService)

	// 门户公开接口热点缓存预热：服务启动即生成门户缓存（对应"部署时生成缓存"），
	// 不阻塞启动；后续内容变更由各写操作经 InvalidateContentCache 保持最新。
	go services.WarmUpPublicCache(db, cacheService)

	// 启动服务器
	addr := ":" + cfg.Server.Port
	utils.Logger.Infow("服务器启动完成", "addr", addr)
	if err := r.Run(addr); err != nil {
		utils.Logger.Errorw("服务器启动失败", "error", err)
		log.Fatalf("服务器启动失败: %v", err)
	}
}
