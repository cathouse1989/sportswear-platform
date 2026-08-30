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

// ==================== 季度分表基础设施 ====================
//
// 需求：操作日志、门户访问流量日志，每季度一张新表，表名后缀带上季度起始日（年月日）。
// 例如 2026 年 Q3（7月1日）→ operation_logs_20260701 / visit_logs_20260701。
// 写入时按记录的 createdAt 路由到对应季度表；查询时按时间范围组合（UNION ALL）多张季度表 + 存量旧表。

// QuarterStart 返回 t 所在季度的起始日（UTC 零点）
func QuarterStart(t time.Time) time.Time {
	year, month, _ := t.Date()
	qIndex := (int(month) - 1) / 3 // 0=Q1, 1=Q2, 2=Q3, 3=Q4
	return time.Date(year, time.Month(qIndex*3+1), 1, 0, 0, 0, 0, time.UTC)
}

// NextQuarter 返回 t 的下一个季度起始日
func NextQuarter(t time.Time) time.Time {
	return QuarterStart(t).AddDate(0, 3, 0)
}

// QuarterTableName 返回 base 表对应 t 所在季度的表名（后缀 YYYYMMDD）
func QuarterTableName(base string, t time.Time) string {
	return fmt.Sprintf("%s_%s", base, QuarterStart(t).Format("20060102"))
}

// QuarterTablesBetween 返回 [start,end] 覆盖的所有季度表名（含尚未创建的表）
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
func EnsureTable(db *gorm.DB, tableName string, model interface{}) error {
	oi, _ := ensureOnce.LoadOrStore(tableName, &sync.Once{})
	once := oi.(*sync.Once)
	var err error
	once.Do(func() {
		err = db.Table(tableName).AutoMigrate(model)
	})
	return err
}

// EnsureCurrentQuarterTables 启动时预建当前季度的操作日志/访问日志表
func EnsureCurrentQuarterTables(db *gorm.DB) error {
	now := time.Now()
	specs := []struct {
		base  string
		model interface{}
	}{
		{"operation_logs", &models.OperationLog{}},
		{"visit_logs", &models.VisitLog{}},
	}
	for _, spec := range specs {
		name := QuarterTableName(spec.base, now)
		if err := EnsureTable(db, name, spec.model); err != nil {
			return fmt.Errorf("创建季度表 %s 失败: %w", name, err)
		}
	}
	return nil
}

// ExistingQuarterTables 列出实际存在的 base 前缀季度表（按表名字典序）
func ExistingQuarterTables(db *gorm.DB, base string) ([]string, error) {
	var names []string
	if err := db.Raw(
		"SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE ?",
		base+"_%",
	).Scan(&names).Error; err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
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