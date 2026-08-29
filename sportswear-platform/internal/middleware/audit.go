package middleware

import (
	"bytes"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sportswear-platform/internal/models"
	"sportswear-platform/internal/services"
)

// Audit 审计中间件：自动记录后台写操作（POST/PUT/DELETE）到 OperationLog
// 需在 Auth 中间件之后使用（依赖 context 中的用户 ID）
func Audit(opLogService *services.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// 只记录写操作
		if method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" {
			c.Next()
			return
		}

		// 记录请求体（限制大小，避免过大）
		var bodyStr string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(io.LimitReader(c.Request.Body, 64*1024))
			bodyStr = string(bodyBytes)
			// 恢复 body 供后续 handler 使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 记录操作日志（异步，不阻塞请求）
		userID := GetUserID(c)
		path := c.Request.URL.Path
		module := parseModule(path)
		operation := mapMethodToOperation(method)

		c.Next() // 先执行 handler

		// 异步记录（goroutine，避免阻塞主流程）
		go func() {
			status := c.Writer.Status()
			// 只记录成功的写操作（2xx）
			if status >= 200 && status < 300 {
				_ = opLogService.Record(&models.OperationLog{
					UserID:      userID,
					Operation:   operation,
					Module:      module,
					EntityType:  parseAuditEntity(path),
					EntityID:    parseAuditEntityID(path),
					After:       truncate(bodyStr, 2000),
					IP:          c.ClientIP(),
					UserAgent:   truncate(c.Request.UserAgent(), 500),
					Description: method + " " + path,
				})
			}
		}()
	}
}

// parseModule 从路径解析模块名
func parseModule(path string) string {
	path = strings.TrimPrefix(path, "/api/v1/admin/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "unknown"
}

// parseAuditEntity 从路径解析实体类型
func parseAuditEntity(path string) string {
	path = strings.TrimPrefix(path, "/api/v1/admin/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "unknown"
}

// parseAuditEntityID 从路径解析实体 ID
func parseAuditEntityID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		// 检查第二段是否为 UUID
		if _, err := uuid.Parse(parts[1]); err == nil {
			return parts[1]
		}
	}
	return ""
}

// mapMethodToOperation 将 HTTP 方法映射为操作类型
func mapMethodToOperation(method string) models.OperationType {
	switch method {
	case "POST":
		return models.OperationCreate
	case "PUT", "PATCH":
		return models.OperationUpdate
	case "DELETE":
		return models.OperationDelete
	default:
		return models.OperationOther
	}
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
