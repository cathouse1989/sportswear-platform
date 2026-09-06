package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newMockGeoIPService 创建使用 sqlmock 的 GeoIPService（正则匹配 SQL，避免依赖 GORM 生成细节）
func newMockGeoIPService(t *testing.T) (*GeoIPService, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock 创建失败: %v", err)
	}
	dialector := postgres.New(postgres.Config{Conn: mockDB, DriverName: "postgres"})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm 打开失败: %v", err)
	}
	return NewGeoIPService(gormDB), mock, mockDB
}

// TestIPv4ToUint32 验证 IPv4 → 整数转换与内网/保留地址过滤。
func TestIPv4ToUint32(t *testing.T) {
	cases := []struct {
		ip   string
		want int64
		ok   bool
	}{
		{"8.8.8.8", 134744072, true},
		{"1.2.3.4", 16909060, true},
		{"255.255.255.255", 4294967295, true},
		{"192.168.1.1", 0, false}, // 内网
		{"10.0.0.1", 0, false},    // 内网
		{"172.16.0.1", 0, false},  // 内网
		{"127.0.0.1", 0, false},   // 回环
		{"169.254.1.1", 0, false}, // 链路本地
		{"::1", 0, false},         // IPv6
		{"not-an-ip", 0, false},
		{"unknown", 0, false},
		{"anonymous", 0, false},
		{"", 0, false},
	}
	for _, c := range cases {
		got, ok := ipv4ToUint32(c.ip)
		if ok != c.ok || got != c.want {
			t.Errorf("ipv4ToUint32(%q) = (%d, %v), want (%d, %v)", c.ip, got, ok, c.want, c.ok)
		}
	}
}

// TestGeoIPService_LookupCountry 命中 IP 段时返回国家。
func TestGeoIPService_LookupCountry(t *testing.T) {
	svc, mock, db := newMockGeoIPService(t)
	defer db.Close()

	// 8.8.8.8 → 134744072，落在 [134744064, 134744319] 段（Google）
	mock.ExpectQuery(`SELECT \* FROM "ip_geo_ranges" WHERE start_ip <= \$1 AND "ip_geo_ranges"\."deleted_at" IS NULL ORDER BY start_ip DESC`).
		WithArgs(int64(134744072), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "start_ip", "end_ip", "country", "country_name"}).
			AddRow("11111111-1111-1111-1111-111111111111", time.Now(), time.Now(), nil, int64(134744064), int64(134744319), "US", "美国"))

	iso2, name := svc.LookupCountry("8.8.8.8")
	if iso2 != "US" || name != "美国" {
		t.Errorf("LookupCountry = (%q, %q), want (US, 美国)", iso2, name)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未满足的 mock 期望: %v", err)
	}
}

// TestGeoIPService_LookupCountry_PrivateIP 内网 IP 直接返回空且不触发 DB 查询。
func TestGeoIPService_LookupCountry_PrivateIP(t *testing.T) {
	svc, mock, db := newMockGeoIPService(t)
	defer db.Close()

	iso2, name := svc.LookupCountry("192.168.1.1")
	if iso2 != "" || name != "" {
		t.Errorf("内网 IP 应返回空, got (%q, %q)", iso2, name)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("内网 IP 不应触发查询: %v", err)
	}
}

