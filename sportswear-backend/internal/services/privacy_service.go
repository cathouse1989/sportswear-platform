package services

import (
	"strings"
	"time"

	"gorm.io/gorm"

	"sportswear-backend/internal/database"
	"sportswear-backend/internal/utils"
)

// ==================== 数据保留与隐私清理服务 ====================
//
// 实现 GDPR/PIPL/CCPA/LGPD 等法规的数据最小化与保留期限要求：
// - 访问日志（visit_logs_*）：超过 90 天后匿名化 IP / UA / Referer / 设备型号等个人数据，
//   仅保留国家、语言、路径、来源等聚合分析所需字段；
// - 操作日志（operation_logs_*）：超过 90 天后匿名化 IP / UA；
// - 询盘（leads）：超过 2 年且已结束（won/lost/spam）的询盘，匿名化联系方式，
//   仅保留用于业务统计的维度字段。
//
// 通过后台定时任务每日执行一次（StartRetentionSweeper），也可手动调用 runRetention。

// PrivacyService 数据保留清理服务
type PrivacyService struct {
	db *gorm.DB
}

// NewPrivacyService 创建数据保留清理服务
func NewPrivacyService(db *gorm.DB) *PrivacyService {
	return &PrivacyService{db: db}
}

// 保留期限常量（天）
const (
	// VisitLogRetentionDays 访问日志个人数据保留期：90 天
	VisitLogRetentionDays = 90
	// OperationLogRetentionDays 操作日志个人数据保留期：90 天
	OperationLogRetentionDays = 90
	// LeadRetentionYears 询盘个人数据保留期：2 年（业务完成后）
	LeadRetentionYears = 2
)

// StartRetentionSweeper 启动后台数据保留清理任务（每日凌晨 3 点执行一次）
func (s *PrivacyService) StartRetentionSweeper() {
	go func() {
		// 启动 30 秒后先执行一次，避免与服务器启动高峰冲突
		time.Sleep(30 * time.Second)
		s.runRetention()
		for {
			now := time.Now()
			// 计算距离下一个凌晨 3 点的时间
			next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
			if !next.After(now) {
				next = next.AddDate(0, 0, 1)
			}
			time.Sleep(time.Until(next))
			s.runRetention()
		}
	}()
	utils.Logger.Infow("数据保留清理任务已启动", "visit_log_days", VisitLogRetentionDays, "lead_years", LeadRetentionYears)
}

// runRetention 执行一轮数据保留清理
func (s *PrivacyService) runRetention() {
	start := time.Now()
	s.anonymizeOldVisitLogs()
	s.anonymizeOldOperationLogs()
	s.anonymizeOldLeads()
	utils.Logger.Infow("数据保留清理完成", "duration_ms", time.Since(start).Milliseconds())
}

// anonymizeOldVisitLogs 匿名化超过保留期的访问日志个人数据（遍历所有季度表 + 存量基础表）
func (s *PrivacyService) anonymizeOldVisitLogs() {
	cutoff := time.Now().AddDate(0, 0, -VisitLogRetentionDays)

	tables, err := database.ExistingQuarterTables(s.db, "visit_logs")
	if err != nil {
		utils.Logger.Warnw("列出访问日志季度表失败", "error", err.Error())
		return
	}
	// 存量基础表
	tables = append(tables, "visit_logs")

	// 每季度表独立执行，避免跨表事务
	for _, table := range tables {
		res := s.db.Table(table).
			Where("created_at < ? AND (user_agent IS NOT NULL AND user_agent != ?)", cutoff, "anonymous").
			Updates(map[string]interface{}{
				"ip":           "0.0.0.0",
				"user_agent":   "anonymous",
				"referer":      "",
				"device_model": "",
			})
		if res.Error != nil {
			utils.Logger.Warnw("匿名化访问日志失败", "table", table, "error", res.Error.Error())
			continue
		}
		if res.RowsAffected > 0 {
			utils.Logger.Infow("访问日志已匿名化", "table", table, "rows", res.RowsAffected)
		}
	}
}

// anonymizeOldOperationLogs 匿名化超过保留期的操作日志个人数据
func (s *PrivacyService) anonymizeOldOperationLogs() {
	cutoff := time.Now().AddDate(0, 0, -OperationLogRetentionDays)

	tables, err := database.ExistingQuarterTables(s.db, "operation_logs")
	if err != nil {
		utils.Logger.Warnw("列出操作日志季度表失败", "error", err.Error())
		return
	}
	tables = append(tables, "operation_logs")

	for _, table := range tables {
		res := s.db.Table(table).
			Where("created_at < ? AND (user_agent IS NOT NULL AND user_agent != ?)", cutoff, "anonymous").
			Updates(map[string]interface{}{
				"ip":         "0.0.0.0",
				"user_agent": "anonymous",
			})
		if res.Error != nil {
			utils.Logger.Warnw("匿名化操作日志失败", "table", table, "error", res.Error.Error())
			continue
		}
		if res.RowsAffected > 0 {
			utils.Logger.Infow("操作日志已匿名化", "table", table, "rows", res.RowsAffected)
		}
	}
}

// anonymizeOldLeads 匿名化超过保留期且已结束的询盘联系方式
func (s *PrivacyService) anonymizeOldLeads() {
	cutoff := time.Now().AddDate(0, -LeadRetentionYears, 0)

	// 已结束状态：won（成交）/ lost（流失）/ spam（垃圾）
	finishedStatuses := []string{"won", "lost", "spam"}

	var leads []struct {
		ID        string
		CreatedAt time.Time
	}
	if err := s.db.Table("leads").
		Select("id, created_at").
		Where("status IN ? AND created_at < ?", finishedStatuses, cutoff).
		Scan(&leads).Error; err != nil {
		utils.Logger.Warnw("查询过期询盘失败", "error", err.Error())
		return
	}

	for _, lead := range leads {
		res := s.db.Table("leads").
			Where("id = ?", lead.ID).
			Updates(map[string]interface{}{
				"name":            "[已过期]",
				"company":         "[已过期]",
				"email":           "anonymized_" + strings.ReplaceAll(lead.ID, "-", "") + "@deleted.local",
				"phone":           "[已过期]",
				"whatsapp":        "[已过期]",
				"company_website": "[已过期]",
				"message":         "[已过期]",
				"ip":              "0.0.0.0",
			})
		if res.Error != nil {
			utils.Logger.Warnw("匿名化过期询盘失败", "lead_id", lead.ID, "error", res.Error.Error())
		}
	}
	if len(leads) > 0 {
		utils.Logger.Infow("过期询盘已匿名化", "count", len(leads))
	}
}
