package database

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// ==================== 分表基础设施 ====================
//
// 需求：操作日志、门户访问流量日志，按月一张新表，表名后缀带上月份（年月）。
// 例如 2026 年 7 月 → operation_logs_202607 / visit_logs_202607。
// 写入时按记录的 createdAt 路由到对应月表；查询时按时间范围组合（UNION ALL）多张月表 + 存量旧表。
//
// 月度分表优势：单表数据量可控（~15万条/月），查询 90 天 = 3 张表 UNION ALL，性能优于季度分表。

// MonthTableName 返回 base 表对应 t 所在月份的表名（后缀 YYYYMM）
func MonthTableName(base string, t time.Time) string {
	return fmt.Sprintf("%s_%s", base, t.Format("200601"))
}

// MonthTablesBetween 返回 [start,end] 覆盖的所有月表名（含尚未创建的表）
func MonthTablesBetween(base string, start, end time.Time) []string {
	if start.After(end) {
		return nil
	}
	var names []string
	y, m := start.Year(), int(start.Month())
	endY, endM := end.Year(), int(end.Month())
	for y < endY || (y == endY && m <= endM) {
		names = append(names, fmt.Sprintf("%s_%d%02d", base, y, m))
		m++
		if m > 12 {
			m = 1
			y++
		}
	}
	return names
}

// QuarterStart 返回 t 所在季度的起始日（UTC 零点）—— 保留兼容
func QuarterStart(t time.Time) time.Time {
	year, month, _ := t.Date()
	qIndex := (int(month) - 1) / 3 // 0=Q1, 1=Q2, 2=Q3, 3=Q4
	return time.Date(year, time.Month(qIndex*3+1), 1, 0, 0, 0, 0, time.UTC)
}

// NextQuarter 返回 t 的下一个季度起始日 —— 保留兼容
func NextQuarter(t time.Time) time.Time {
	return QuarterStart(t).AddDate(0, 3, 0)
}

// QuarterTableName 返回 base 表对应 t 所在季度的表名（后缀 YYYYMMDD）—— 保留兼容
func QuarterTableName(base string, t time.Time) string {
	return fmt.Sprintf("%s_%s", base, QuarterStart(t).Format("20060102"))
}

// QuarterTablesBetween 返回 [start,end] 覆盖的所有季度表名（含尚未创建的表）—— 保留兼容
func QuarterTablesBetween(base string, start, end time.Time) []string {
	if start.After(end) {
		return nil
	}
	q := QuarterStart(start)
	endQ := QuarterStart(end)
	var names []string
	for !q.After(endQ) {
		names = append(names, fmt.Sprintf("%s_%s", base, q.Format("20060102")))
		q = q.AddDate(0, 3, 0)
	}
	return names
}

var ensureOnce sync.Map // tableName -> *sync.Once

// EnsureTable 确保指定表存在（懒创建，进程内幂等）。
// 使用 GORM AutoMigrate 按模型建表/补列/建索引；已在进程内确认过的表直接跳过。
// 建表后显式创建复合索引，保障流量分析聚合查询性能。
func EnsureTable(db *gorm.DB, tableName string, model interface{}) error {
	oi, _ := ensureOnce.LoadOrStore(tableName, &sync.Once{})
	once := oi.(*sync.Once)
	var err error
	once.Do(func() {
		if err = db.Table(tableName).AutoMigrate(model); err != nil {
			return
		}
		// 显式创建复合索引（流量分析高频查询路径）
		// 注意：索引名包含表名，不同表独立，避免冲突
		safeName := strings.ReplaceAll(tableName, ".", "_")
		db.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_created_type ON %s (created_at, visit_type)`, safeName, tableName))
		db.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_visitor ON %s (visitor_id, created_at)`, safeName, tableName))
		db.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_entity ON %s (entity_type, entity_slug, created_at)`, safeName, tableName))
		db.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_ip ON %s (ip, created_at)`, safeName, tableName))
		db.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_source ON %s (source, created_at)`, safeName, tableName))
	})
	return err
}

// EnsureCurrentTables 启动时预建当前月份的操作日志/访问日志表
func EnsureCurrentTables(db *gorm.DB) error {
	now := time.Now()
	specs := []struct {
		base  string
		model interface{}
	}{
		{"operation_logs", &models.OperationLog{}},
		{"visit_logs", &models.VisitLog{}},
	}
	for _, spec := range specs {
		name := MonthTableName(spec.base, now)
		if err := EnsureTable(db, name, spec.model); err != nil {
			return fmt.Errorf("创建月度表 %s 失败: %w", name, err)
		}
	}
	return nil
}

// EnsureCurrentQuarterTables 启动时预建当前月份的操作日志/访问日志表（兼容旧名）
func EnsureCurrentQuarterTables(db *gorm.DB) error {
	return EnsureCurrentTables(db)
}

// ExistingTables 列出实际存在的 base 前缀月度/季度表（按表名字典序）
func ExistingTables(db *gorm.DB, base string) ([]string, error) {
	var names []string
	if err := db.Raw(
		`SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE ? ESCAPE '|'`,
		base+"|_%",
	).Scan(&names).Error; err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

// ExistingQuarterTables 列出实际存在的 base 前缀季度表（兼容旧名）
func ExistingQuarterTables(db *gorm.DB, base string) ([]string, error) {
	return ExistingTables(db, base)
}

// UnionAll 对每张表应用同一 SELECT 模板并 UNION ALL。
// selectFmt 中的 "__T__" 会被替换为表名（可在任意位置使用，例如 src 列 '...'__T__'...'）。
// argsPerTable 为每张表重复追加的参数（where 占位符参数）。
func UnionAll(tables []string, selectFmt string, argsPerTable []interface{}) (string, []interface{}) {
	var sb strings.Builder
	var out []interface{}
	for i, t := range tables {
		if i > 0 {
			sb.WriteString(" UNION ALL ")
		}
		sb.WriteString(strings.ReplaceAll(selectFmt, "__T__", t))
		out = append(out, argsPerTable...)
	}
	return sb.String(), out
}