package utils

import (
	"fmt"
	"time"
)

// ==================== 全球时区统一方案 ====================
//
// 面向全球的网站，时间处理遵循以下原则：
//
// 1. 存储统一 UTC：数据库所有时间字段统一使用 UTC，避免多时区混乱
// 2. API 统一 RFC3339 UTC：对外接口统一返回带 Z 后缀的 ISO8601 时间，
//    客户端（浏览器/App）根据用户本地时区自动转换显示
// 3. 时区转换：提供工具函数，需要时按访客时区转换
// 4. 无歧义：RFC3339 格式（如 "2026-08-22T07:30:00Z"）天然无歧义

const (
	// DefaultTimezone 服务器默认业务时区（用于后台展示）
	DefaultTimezone = "Asia/Shanghai"
	// UTC 时区标识
	UTC = "UTC"
)

// TimeFormatRFC3339 统一时间格式（ISO8601 + UTC）
const TimeFormatRFC3339 = time.RFC3339

// Now 返回 UTC 当前时间（全局统一时间源，存储层使用）
func NowUTC() time.Time {
	return time.Now().UTC()
}

// NowLocal 返回指定时区的当前时间
func NowIn(tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	return time.Now().In(loc)
}

// ToUTC 将任意时间转换为 UTC
func ToUTC(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.UTC()
}

// ToTimezone 将时间转换到指定时区（用于按访客时区展示）
func ToTimezone(t time.Time, tz string) time.Time {
	if t.IsZero() {
		return t
	}
	if tz == "" || tz == "UTC" {
		return t.UTC()
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return t.UTC()
	}
	return t.In(loc)
}

// FormatRFC3339 将时间格式化为 RFC3339 UTC（带 Z 后缀）
// 示例：2026-08-22T07:30:00Z
func FormatRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// FormatRFC3339In 将时间转换到指定时区后格式化为 RFC3339
// 示例：FormatRFC3339In(t, "America/New_York") -> 2026-08-22T03:30:00-04:00
func FormatRFC3339In(t time.Time, tz string) string {
	if t.IsZero() {
		return ""
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return FormatRFC3339(t)
	}
	return t.In(loc).Format(time.RFC3339)
}

// ParseRFC3339 解析 RFC3339 时间字符串（自动处理时区偏移）
// 统一转换为 UTC 存储
func ParseRFC3339(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("无效的时间格式: %s (期望 RFC3339)", s)
	}
	return t.UTC(), nil
}

// ParseFlexibleTime 灵活解析时间（支持多种常见格式）
// 支持：RFC3339、日期、日期时间等
func ParseFlexibleTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析时间: %s", s)
}

// GetTimezone 从请求解析时区，默认返回 UTC
// 优先级：显式 tz 参数 > 默认 UTC
func ResolveTimezone(tz string) *time.Location {
	if tz == "" {
		return time.UTC
	}
	if loc, err := time.LoadLocation(tz); err == nil {
		return loc
	}
	return time.UTC
}

// FormatDateOnly 格式化日期（YYYY-MM-DD），按指定时区
// 用于按访客时区显示"今天/昨天"等日期
func FormatDateOnly(t time.Time, tz string) string {
	if t.IsZero() {
		return ""
	}
	return ToTimezone(t, tz).Format("2006-01-02")
}

// FormatHumanTime 人性化时间格式（按指定时区）
// 示例：2026-08-22 15:30（24 小时制）
func FormatHumanTime(t time.Time, tz string) string {
	if t.IsZero() {
		return ""
	}
	return ToTimezone(t, tz).Format("2006-01-02 15:04")
}

// StartOfDayUTC 返回 UTC 某天的开始时间
func StartOfDayUTC(t time.Time) time.Time {
	utc := t.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

// EndOfDayUTC 返回 UTC 某天的结束时间
func EndOfDayUTC(t time.Time) time.Time {
	utc := t.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 23, 59, 59, 999999999, time.UTC)
}

// IsToday 判断时间是否为今天（按指定时区）
func IsToday(t time.Time, tz string) bool {
	now := NowIn(tz)
	local := ToTimezone(t, tz)
	return now.Year() == local.Year() &&
		now.Month() == local.Month() &&
		now.Day() == local.Day()
}

// CommonTimezones 常用时区映射（用于前端展示时区选择）
var CommonTimezones = []struct {
	Value  string
	Label  string
	Offset string
}{
	{"Asia/Shanghai", "中国（北京）", "+08:00"},
	{"Asia/Tokyo", "日本（东京）", "+09:00"},
	{"Asia/Seoul", "韩国（首尔）", "+09:00"},
	{"Asia/Hong_Kong", "香港", "+08:00"},
	{"Asia/Singapore", "新加坡", "+08:00"},
	{"Asia/Dubai", "迪拜", "+04:00"},
	{"Europe/London", "英国（伦敦）", "+00:00"},
	{"Europe/Paris", "法国（巴黎）", "+01:00"},
	{"Europe/Berlin", "德国（柏林）", "+01:00"},
	{"Europe/Moscow", "俄罗斯（莫斯科）", "+03:00"},
	{"America/New_York", "美国（纽约）", "-05:00"},
	{"America/Chicago", "美国（芝加哥）", "-06:00"},
	{"America/Los_Angeles", "美国（洛杉矶）", "-08:00"},
	{"America/Sao_Paulo", "巴西（圣保罗）", "-03:00"},
	{"America/Mexico_City", "墨西哥（墨西哥城）", "-06:00"},
	{"Australia/Sydney", "澳大利亚（悉尼）", "+10:00"},
	{"Pacific/Auckland", "新西兰（奥克兰）", "+12:00"},
	{"Africa/Cairo", "埃及（开罗）", "+02:00"},
	{"Africa/Johannesburg", "南非（约翰内斯堡）", "+02:00"},
	{"UTC", "协调世界时（UTC）", "+00:00"},
}
