package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局日志实例
var (
	Logger    *zap.SugaredLogger // 业务日志
	AccessLog *zap.SugaredLogger // 访问日志
	AuditLog  *zap.SugaredLogger // 审计日志
)

// logManager 日志目录管理器：按日期分目录 + 自动轮转 + 自动清理
type logManager struct {
	mu          sync.Mutex
	baseDir     string // logs/
	retainDays  int    // 保留天数
	currentDate string
	files       map[string]*os.File // name -> file (app/access/audit)
}

var globalLogManager *logManager

// InitLogger 初始化分级日志系统
// 目录结构：logs/yyyy-mm-dd/app.log, access.log, audit/audit.log
// 自动能力：每日自动切换新目录；自动清理超过保留期的旧目录
func InitLogger(env string) {
	retainDays := 30 // 默认保留 30 天
	if env == "development" {
		retainDays = 7
	}

	globalLogManager = &logManager{
		baseDir:    "logs",
		retainDays: retainDays,
		files:      make(map[string]*os.File),
	}
	globalLogManager.openToday()

	// 启动后台协程：每日 0 点自动轮转 + 清理过期日志
	go globalLogManager.autoRotate()
}

// openToday 打开今天的日志文件
func (m *logManager) openToday() {
	m.mu.Lock()
	defer m.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today == m.currentDate {
		return
	}

	// 关闭旧文件
	for _, f := range m.files {
		_ = f.Close()
	}
	m.files = make(map[string]*os.File)

	// 创建日期目录：logs/2026-08-22/
	dayDir := filepath.Join(m.baseDir, today)
	auditDir := filepath.Join(dayDir, "audit")
	_ = os.MkdirAll(auditDir, 0755)

	// 编码器配置
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// 根据环境确定最低级别
	level := zapcore.InfoLevel
	if envLevel() == "development" {
		level = zapcore.DebugLevel
	}
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	// ========== 业务日志 ==========
	appFile, err := os.OpenFile(filepath.Join(dayDir, "app.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		m.files["app"] = appFile
		bizCore := zapcore.NewTee(
			zapcore.NewCore(jsonEncoder, zapcore.AddSync(appFile), levelEnabler),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), levelEnabler),
		)
		Logger = zap.New(bizCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
	}

	// ========== 访问日志 ==========
	accessFile, err := os.OpenFile(filepath.Join(dayDir, "access.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		m.files["access"] = accessFile
		AccessLog = zap.New(zapcore.NewCore(jsonEncoder, zapcore.AddSync(accessFile), zapcore.InfoLevel)).Sugar()
	}

	// ========== 审计日志 ==========
	auditFile, err := os.OpenFile(filepath.Join(auditDir, "audit.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		m.files["audit"] = auditFile
		AuditLog = zap.New(zapcore.NewCore(jsonEncoder, zapcore.AddSync(auditFile), zapcore.InfoLevel)).Sugar()
	}

	m.currentDate = today
}

// autoRotate 每日 0 点自动轮转日志目录，并清理过期日志
func (m *logManager) autoRotate() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 5, 0, now.Location())
		time.Sleep(next.Sub(now))

		m.openToday()
		m.cleanExpired()
	}
}

// cleanExpired 清理超过保留期的日志目录
func (m *logManager) cleanExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -m.retainDays)
	cutoffStr := cutoff.Format("2006-01-02")

	var dirs []string
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) == 10 && e.Name()[4] == '-' {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)

	for _, d := range dirs {
		if d < cutoffStr {
			_ = os.RemoveAll(filepath.Join(m.baseDir, d))
		}
	}
}

// GetLogDirs 获取现有日志目录列表（供管理接口使用）
func GetLogDirs() []string {
	entries, err := os.ReadDir("logs")
	if err != nil {
		return nil
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() && len(e.Name()) == 10 && strings.Contains(e.Name(), "-") {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dirs)))
	return dirs
}

// LogRotationInfo 日志轮转信息
func LogRotationInfo() map[string]interface{} {
	info := map[string]interface{}{
		"base_dir":     "logs",
		"structure":    "logs/yyyy-mm-dd/{app.log, access.log, audit/audit.log}",
		"current_date": globalLogManager.currentDate,
		"retain_days":  globalLogManager.retainDays,
		"auto_rotate":  "daily at 00:00",
		"auto_clean":   fmt.Sprintf("dirs older than %d days removed automatically", globalLogManager.retainDays),
		"dirs":         GetLogDirs(),
	}
	return info
}

func envLevel() string {
	env := os.Getenv("SERVER_ENV")
	if env == "" {
		return "development"
	}
	return env
}

// Sync 刷新日志缓冲
func SyncLogs() {
	if Logger != nil {
		_ = Logger.Sync()
	}
	if AccessLog != nil {
		_ = AccessLog.Sync()
	}
	if AuditLog != nil {
		_ = AuditLog.Sync()
	}
}
