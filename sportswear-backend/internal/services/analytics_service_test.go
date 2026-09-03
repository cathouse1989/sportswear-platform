package services

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
)

// newMockAnalyticsService 创建使用 sqlmock 的 AnalyticsService
func newMockAnalyticsService(t *testing.T) (*AnalyticsService, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock 创建失败: %v", err)
	}
	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm 打开失败: %v", err)
	}
	svc := NewAnalyticsService(gormDB)
	return svc, mock, mockDB
}

// mockVisitTables 模拟 visitTables 的基础表查询
// 返回 visit_logs 基础表 + 无月度分表
func mockVisitTables(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename LIKE $1 ESCAPE '|'`)).
		WithArgs("visit_logs|_%").
		WillReturnRows(sqlmock.NewRows([]string{"tablename"}))
}

// ============ GetTrafficOverview ============

func TestGetTrafficOverview(t *testing.T) {
	svc, mock, db := newMockAnalyticsService(t)
	defer db.Close()

	mockVisitTables(mock)

	// 模拟 UNION ALL 子查询结果
	rows := sqlmock.NewRows([]string{"count"}).AddRow(42)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id \|\| ':' \|\| path\)`).
		WillReturnRows(rows)

	rows2 := sqlmock.NewRows([]string{"count"}).AddRow(10)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT ip\)`).
		WillReturnRows(rows2)

	rows3 := sqlmock.NewRows([]string{"count"}).AddRow(8)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id\)`).
		WillReturnRows(rows3)

	rows4 := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\)`).
		WillReturnRows(rows4)

	rows5 := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id \|\| ':' \|\| path\)`).
		WillReturnRows(rows5)

	// 询盘数
	mock.ExpectQuery(`SELECT count\(\*\) FROM "leads" WHERE created_at`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	// 每日趋势
	mock.ExpectQuery(`SELECT DATE`).
		WillReturnRows(sqlmock.NewRows([]string{"date", "visits", "unique_ips"}).
			AddRow("2026-09-01", 10, 5).
			AddRow("2026-09-02", 15, 7))

	result, err := svc.GetTrafficOverview(7)
	if err != nil {
		t.Fatalf("GetTrafficOverview 失败: %v", err)
	}

	if totalVisits, ok := result["total_visits"]; !ok || totalVisits.(int64) != 42 {
		t.Errorf("total_visits = %v, want 42", totalVisits)
	}
	if uniqueIPs, ok := result["unique_ips"]; !ok || uniqueIPs.(int64) != 10 {
		t.Errorf("unique_ips = %v, want 10", uniqueIPs)
	}
	if leads, ok := result["leads"]; !ok || leads.(int64) != 3 {
		t.Errorf("leads = %v, want 3", leads)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}
func TestGetTrafficOverview_DaysZero(t *testing.T) {
	svc, mock, db := newMockAnalyticsService(t)
	defer db.Close()

	mockVisitTables(mock)

	// 当 days <= 0 时，默认使用 7 天
	zeroRow := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id \|\| ':' \|\| path\)`).
		WillReturnRows(zeroRow)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT ip\)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id\)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT COUNT\(\*\)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id \|\| ':' \|\| path\)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "leads" WHERE created_at`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT DATE`).
		WillReturnRows(sqlmock.NewRows([]string{"date", "visits", "unique_ips"}))

	result, err := svc.GetTrafficOverview(0)
	if err != nil {
		t.Fatalf("GetTrafficOverview(0) 失败: %v", err)
	}
	if days, ok := result["days"]; !ok || days.(int) != 7 {
		t.Errorf("days = %v, want 7", days)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}

// ============ GetTopPages ============

func TestGetTopPages(t *testing.T) {
	svc, mock, db := newMockAnalyticsService(t)
	defer db.Close()

	mockVisitTables(mock)

	mock.ExpectQuery(`SELECT path, entity_type, entity_slug, COUNT\(\*\) AS views`).
		WillReturnRows(sqlmock.NewRows([]string{"path", "entity_type", "entity_slug", "views", "unique_visitors"}).
			AddRow("/products/custom-leggings", "product", "custom-leggings", 15, 8).
			AddRow("/blogs/fabric-guide", "blog", "fabric-guide", 10, 6))

	pages, err := svc.GetTopPages(7, 5)
	if err != nil {
		t.Fatalf("GetTopPages 失败: %v", err)
	}
	if len(pages) != 2 {
		t.Errorf("len(pages) = %d, want 2", len(pages))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}

// ============ GetConversionFunnel ============

func TestGetConversionFunnel(t *testing.T) {
	svc, mock, db := newMockAnalyticsService(t)
	defer db.Close()

	mockVisitTables(mock)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(100)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id \|\| ':' \|\| path\)`).
		WillReturnRows(rows)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id \|\| ':' \|\| path\)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "leads" WHERE created_at`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	funnel, err := svc.GetConversionFunnel(30)
	if err != nil {
		t.Fatalf("GetConversionFunnel 失败: %v", err)
	}
	if funnel.PageViews != 100 {
		t.Errorf("PageViews = %d, want 100", funnel.PageViews)
	}
	if funnel.ProductViews != 100 {
		t.Errorf("ProductViews = %d, want 100", funnel.ProductViews)
	}
	if funnel.Leads != 5 {
		t.Errorf("Leads = %d, want 5", funnel.Leads)
	}
	if funnel.ProductRate != 100.0 {
		t.Errorf("ProductRate = %f, want 100.0", funnel.ProductRate)
	}
	if funnel.LeadRate != 5.0 {
		t.Errorf("LeadRate = %f, want 5.0", funnel.LeadRate)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}

// ============ GetConsentInsights ============

func TestGetConsentInsights(t *testing.T) {
	svc, mock, db := newMockAnalyticsService(t)
	defer db.Close()

	mockVisitTables(mock)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(80)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id\) FROM`).
		WillReturnRows(rows)
	rows2 := sqlmock.NewRows([]string{"count"}).AddRow(20)
	mock.ExpectQuery(`SELECT COUNT\(DISTINCT visitor_id\) FROM`).
		WillReturnRows(rows2)

	// 同意用户的询盘数
	mock.ExpectQuery(`SELECT count\(\*\) FROM "leads" WHERE \(created_at`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	// 拒绝用户的询盘数
	mock.ExpectQuery(`SELECT count\(\*\) FROM "leads" WHERE \(created_at`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	insights, err := svc.GetConsentInsights(30)
	if err != nil {
		t.Fatalf("GetConsentInsights 失败: %v", err)
	}
	if insights.GrantedCount != 80 {
		t.Errorf("GrantedCount = %d, want 80", insights.GrantedCount)
	}
	if insights.DeniedCount != 20 {
		t.Errorf("DeniedCount = %d, want 20", insights.DeniedCount)
	}
	if insights.TotalVisitors != 100 {
		t.Errorf("TotalVisitors = %d, want 100", insights.TotalVisitors)
	}
	if insights.ConsentRate != 80.0 {
		t.Errorf("ConsentRate = %f, want 80.0", insights.ConsentRate)
	}
	if insights.GrantedConversion != 12.5 {
		t.Errorf("GrantedConversion = %f, want 12.5", insights.GrantedConversion)
	}
	if insights.DeniedConversion != 10.0 {
		t.Errorf("DeniedConversion = %f, want 10.0", insights.DeniedConversion)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}
// ============ GetIPLeads ============

func TestGetIPLeads(t *testing.T) {
	svc, mock, db := newMockAnalyticsService(t)
	defer db.Close()

	mock.ExpectQuery(`SELECT id, name, company, email, country, status, score, score_level, source, medium, campaign, visitor_id, ip, created_at FROM "leads" WHERE`).
		WithArgs("192.168.1.1", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "ip", "created_at"}).
			AddRow(1, "张三", "192.168.1.1", time.Now()))

	leads, err := svc.GetIPLeads("192.168.1.1", 30)
	if err != nil {
		t.Fatalf("GetIPLeads 失败: %v", err)
	}
	if len(leads) != 1 {
		t.Errorf("len(leads) = %d, want 1", len(leads))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}

func TestGetIPLeads_EmptyIP(t *testing.T) {
	svc, mock, db := newMockAnalyticsService(t)
	defer db.Close()

	// 空 IP 应正常返回空列表（不报错）
	mock.ExpectQuery(`SELECT id, name, company, email, country, status, score, score_level, source, medium, campaign, visitor_id, ip, created_at FROM "leads" WHERE`).
		WithArgs("", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "ip", "created_at"}))

	leads, err := svc.GetIPLeads("", 30)
	if err != nil {
		t.Fatalf("GetIPLeads('') 失败: %v", err)
	}
	if len(leads) != 0 {
		t.Errorf("len(leads) = %d, want 0", len(leads))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}

// ============ normalizeRange ============

func TestNormalizeRange_ZeroStart(t *testing.T) {
	start, _ := normalizeRange(time.Time{}, time.Now())
	if start.IsZero() {
		t.Error("start 不应为 zero")
	}
}

func TestNormalizeRange_EndBeforeStart(t *testing.T) {
	now := time.Now()
	start, end := normalizeRange(now, now.AddDate(0, 0, -10))
	if !end.Before(start) {
		t.Logf("交换后 start=%v, end=%v", start, end)
	}
}

func TestNormalizeRange_MaxDays(t *testing.T) {
	start, end := normalizeRange(time.Now().AddDate(0, 0, -200), time.Now())
	diff := end.Sub(start).Hours() / 24
	if diff > float64(MaxQueryDays) {
		t.Errorf("范围 %f 天超过最大 %d 天", diff, MaxQueryDays)
	}
}

// ============ visitTables ============

func TestVisitTables_OnlyBase(t *testing.T) {
	svc, mock, db := newMockAnalyticsService(t)
	defer db.Close()

	mockVisitTables(mock)

	now := time.Now()
	tables, err := svc.visitTables(now.AddDate(0, 0, -7), now)
	if err != nil {
		t.Fatalf("visitTables 失败: %v", err)
	}
	if len(tables) < 1 {
		t.Error("应至少返回 visit_logs 基础表")
	}
	if tables[0] != "visit_logs" {
		t.Errorf("tables[0] = %s, want visit_logs", tables[0])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}