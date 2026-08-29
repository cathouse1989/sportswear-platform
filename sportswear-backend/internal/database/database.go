package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"sportswear-backend/internal/config"
	"sportswear-backend/internal/models"
	"sportswear-backend/internal/utils"
)

// slowQueryThreshold 慢查询阈值
const slowQueryThreshold = 200 * time.Millisecond

// Init 初始化数据库连接（含连接池优化）
func Init(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: newGormLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// ========== 连接池优化 ==========
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(50)                  // 最大打开连接数
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大存活时间（防止连接老化）
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 空闲连接最大存活时间

	// 自动迁移
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	return db, nil
}

// newGormLogger GORM 日志：慢查询记录 WARN，错误记录 ERROR
func newGormLogger() gormlogger.Interface {
	logLevel := gormlogger.Warn
	if os.Getenv("SERVER_ENV") == "development" {
		logLevel = gormlogger.Info
	}

	return gormlogger.New(
		gormLogWriter{},
		gormlogger.Config{
			SlowThreshold:             slowQueryThreshold,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

// gormLogWriter 将 GORM 日志桥接到 zap 分级日志
type gormLogWriter struct{}

func (w gormLogWriter) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if utils.Logger != nil {
		// GORM 的慢查询和错误都通过这里输出
		utils.Logger.Warnw("database", "detail", msg)
	} else {
		log.Println(msg)
	}
}

// AutoMigrate 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		// 应用维度
		&models.App{},
		&models.AppUser{},
		// 用户与权限
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		// 多语言
		&models.Language{},
		// 系统配置
		&models.Setting{},
		// CMS
		&models.Page{},
		&models.PageTranslation{},
		&models.PageModule{},
		&models.Navigation{},
		// 产品
		&models.Product{},
		&models.ProductTranslation{},
		&models.Category{},
		&models.Series{},
		&models.ProductImage{},
		&models.ProductVideo{},
		&models.ProductSpec{},
		&models.ProductCustomization{},
		&models.Fabric{},
		// 内容
		&models.Blog{},
		&models.BlogTranslation{},
		&models.Case{},
		&models.CaseTranslation{},
		&models.FAQ{},
		&models.FAQTranslation{},
		&models.Factory{},
		&models.Certification{},
		// SEO
		&models.SEO{},
		// 媒体
		&models.Media{},
		// 询盘
		&models.Lead{},
		&models.LeadFollowUp{},
		&models.Quote{},
		// 系统
		&models.Notification{},
		&models.OperationLog{},
		// 流量分析
		&models.VisitLog{},
		// 门户展示配置
		&models.PageVersion{},
		&models.ThemeConfig{},
		// 国际化词条
		&models.I18nEntry{},
		// 媒体存储源配置
		&models.StorageSource{},
		// 多货币
		&models.Currency{},
	)
	if err != nil {
		return err
	}

	// ========== 性能复合索引 ==========
	indexes := []string{
		// 产品：状态+排序（公开列表高频查询）
		"CREATE INDEX IF NOT EXISTS idx_products_status_sort ON products (status, sort_order)",
		"CREATE INDEX IF NOT EXISTS idx_products_category_status ON products (category_id, status)",
		// 页面：状态+slug（公开访问高频查询）
		"CREATE INDEX IF NOT EXISTS idx_pages_status_slug ON pages (status, slug)",
		// 博客：状态+时间（列表按时间倒序）
		"CREATE INDEX IF NOT EXISTS idx_blogs_status_published ON blogs (status, published_at DESC)",
		// 案例：状态+时间
		"CREATE INDEX IF NOT EXISTS idx_cases_status_published ON cases (status, published_at DESC)",
		// SEO：实体定位（唯一查询路径）
		"CREATE INDEX IF NOT EXISTS idx_seos_entity ON seos (entity_type, entity_id)",
		// 翻译表：实体+语言（多语言查询主路径）
		"CREATE INDEX IF NOT EXISTS idx_product_translations_pid_lang ON product_translations (product_id, language)",
		"CREATE INDEX IF NOT EXISTS idx_page_translations_pid_lang ON page_translations (page_id, language)",
		"CREATE INDEX IF NOT EXISTS idx_blog_translations_bid_lang ON blog_translations (blog_id, language)",
		// 访问日志：时间+实体（流量分析聚合查询）
		"CREATE INDEX IF NOT EXISTS idx_visit_logs_created_entity ON visit_logs (created_at, entity_type)",
		"CREATE INDEX IF NOT EXISTS idx_visit_logs_source ON visit_logs (source, created_at)",
		// 询盘：状态+评分（销售看板）
		"CREATE INDEX IF NOT EXISTS idx_leads_status_score ON leads (status, score DESC)",
		"CREATE INDEX IF NOT EXISTS idx_leads_assigned ON leads (assigned_to, next_follow_up)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			log.Printf("创建索引失败: %v", err)
		}
	}

	log.Println("数据库迁移完成")
	return nil
}

// HealthCheck 数据库健康检查
func HealthCheck(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
