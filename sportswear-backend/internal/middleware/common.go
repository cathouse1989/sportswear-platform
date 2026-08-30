package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sportswear-backend/internal/utils"
)

// RequestID 请求 ID 中间件
// 每个请求生成唯一 ID，贯穿：响应体 requestId、响应头 X-Request-ID、业务日志、访问日志、审计日志
// 支持全链路追踪：客户端可携带 X-Request-ID 头保持追踪链
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-App-Key, X-UTM-Source, X-UTM-Medium, X-UTM-Campaign, X-UTM-Content, X-UTM-Term, UTM-Source, UTM-Medium, UTM-Campaign, UTM-Content, UTM-Term")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger 请求日志中间件：按状态分级打印，带 request_id 全链路追踪
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		requestID := c.GetString("request_id")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if utils.Logger != nil {
			entry := []interface{}{
				"request_id", requestID,
				"method", c.Request.Method,
				"path", path,
				"status", status,
				"ip", c.ClientIP(),
				"latency_ms", latency.Milliseconds(),
			}
			switch {
			case status >= 500:
				utils.Logger.Errorw("request", entry...)
			case status >= 400:
				utils.Logger.Warnw("request", entry...)
			default:
				utils.Logger.Infow("request", entry...)
			}
		}
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		AbortWithError(c, 500, "INTERNAL_ERROR", "服务器内部错误")
	})
}

// Timezone 时区解析中间件
// 面向全球网站：从请求头 X-Timezone 或查询参数 tz 解析访客时区，
// 存储到 context，处理器可据此进行时区转换展示。
// 默认 UTC，无值时不报错。
func Timezone() gin.HandlerFunc {
	return func(c *gin.Context) {
		tz := c.GetHeader("X-Timezone")
		if tz == "" {
			tz = c.Query("tz")
		}
		if tz == "" {
			tz = utils.UTC
		}
		// 校验时区有效性，无效则回退 UTC
		loc := utils.ResolveTimezone(tz)
		c.Set("timezone", loc.String())
		c.Next()
	}
}

// GetTimezone 从 context 获取访客时区（处理器中使用）
func GetTimezone(c *gin.Context) string {
	if tz, ok := c.Get("timezone"); ok {
		if s, ok := tz.(string); ok {
			return s
		}
	}
	return utils.UTC
}

// Language 语言解析中间件
// 面向全球网站：从 lang 查询参数或 Accept-Language 头解析客户端语言，
// 存储到 context，处理器据此返回对应语言内容。
// 默认 en，无值时不报错。
func Language() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		if lang == "" {
			lang = c.GetHeader("X-Lang")
		}
		if lang == "" {
			// 从 Accept-Language 头解析首选语言
			if al := c.GetHeader("Accept-Language"); al != "" {
				lang = parseAcceptLanguage(al)
			}
		}
		if lang == "" {
			lang = "en"
		}
		c.Set("lang", lang)
		c.Next()
	}
}

// GetLang 从 context 获取客户端语言（处理器中使用）
func GetLang(c *gin.Context) string {
	if l, ok := c.Get("lang"); ok {
		if s, ok := l.(string); ok {
			return s
		}
	}
	return "en"
}

// parseAcceptLanguage 从 Accept-Language 头解析首选语言
// 示例：zh-CN,zh;q=0.9,en;q=0.8 -> zh-CN
func parseAcceptLanguage(header string) string {
	if header == "" {
		return ""
	}
	// 取第一个逗号前的部分
	for i := 0; i < len(header); i++ {
		if header[i] == ',' || header[i] == ';' {
			return header[:i]
		}
	}
	return header
}
