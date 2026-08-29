package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sportswear-backend/internal/config"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// Auth 认证中间件：验证 JWT 并从数据库加载用户角色与权限
func Auth(cfg *config.Config, authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], cfg.JWT.Secret)
		if err != nil {
			utils.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		userID := claims.UserID.String()

		// 从数据库实时加载用户角色与权限（确保权限变更立即生效）
		roleCodes, err := authService.GetUserRoleCodes(userID)
		if err != nil {
			utils.InternalError(c, "加载用户角色失败")
			c.Abort()
			return
		}

		permissionCodes, err := authService.GetUserPermissions(userID)
		if err != nil {
			utils.InternalError(c, "加载用户权限失败")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("user_roles", roleCodes)
		c.Set("user_permissions", permissionCodes)
		c.Next()
	}
}

// RequirePermission 权限中间件：基于数据库 RBAC 硬控
// super_admin 角色拥有所有权限，其他角色必须拥有对应权限码
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("user_roles")
		if !exists {
			utils.Unauthorized(c, "未认证")
			c.Abort()
			return
		}

		roleList, ok := roles.([]string)
		if !ok {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		// 超级管理员拥有所有权限
		for _, role := range roleList {
			if role == "super_admin" {
				c.Next()
				return
			}
		}

		// 检查用户是否拥有所需权限
		perms, exists := c.Get("user_permissions")
		if !exists {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		permList, ok := perms.([]string)
		if !ok {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		for _, p := range permList {
			if p == permission {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "无权限访问，需要权限: "+permission)
		c.Abort()
	}
}

// RequireAnyPermission 任一权限即可通过
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("user_roles")
		if !exists {
			utils.Unauthorized(c, "未认证")
			c.Abort()
			return
		}

		roleList, ok := roles.([]string)
		if !ok {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		for _, role := range roleList {
			if role == "super_admin" {
				c.Next()
				return
			}
		}

		perms, exists := c.Get("user_permissions")
		if !exists {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		permList, ok := perms.([]string)
		if !ok {
			utils.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		permSet := make(map[string]bool)
		for _, p := range permList {
			permSet[p] = true
		}
		for _, need := range permissions {
			if permSet[need] {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "无权限访问")
		c.Abort()
	}
}

// AppAuth 应用识别中间件：通过 X-App-Key 请求头识别子应用
// 用于子应用接入场景，校验应用 API Key 与启用状态
func AppAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-App-Key")
		if apiKey == "" {
			utils.Unauthorized(c, "未提供应用标识（X-App-Key）")
			c.Abort()
			return
		}

		app, err := authService.GetAppByAPIKey(apiKey)
		if err != nil {
			utils.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		c.Set("app_id", app.ID)
		c.Set("app_code", app.Code)
		c.Next()
	}
}

// ============ 频率限制（公开接口防滥用） ============

type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int           // 时间窗口内最大请求数
	window   time.Duration // 时间窗口
}

var globalRateLimiter = &rateLimiter{
	requests: make(map[string][]time.Time),
	limit:    3000,        // 每窗口 3000 次
	window:   time.Minute, // 每分钟
}

func init() {
	// 后台定期清理过期频率限制记录，防止内存泄漏
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			globalRateLimiter.mu.Lock()
			now := time.Now()
			for k, v := range globalRateLimiter.requests {
				valid := v[:0]
				for _, t := range v {
					if now.Sub(t) < globalRateLimiter.window {
						valid = append(valid, t)
					}
				}
				if len(valid) == 0 {
					delete(globalRateLimiter.requests, k)
				} else {
					globalRateLimiter.requests[k] = valid
				}
			}
			globalRateLimiter.mu.Unlock()
		}
	}()
}

// RateLimit 基于 IP 的频率限制中间件（用于公开接口防滥用）
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path

		globalRateLimiter.mu.Lock()
		now := time.Now()
		key := ip + ":" + path

		// 清理过期记录
		times := globalRateLimiter.requests[key]
		validTimes := times[:0]
		for _, t := range times {
			if now.Sub(t) < globalRateLimiter.window {
				validTimes = append(validTimes, t)
			}
		}

		if len(validTimes) >= globalRateLimiter.limit {
			globalRateLimiter.mu.Unlock()
			Error429(c)
			c.Abort()
			return
		}

		globalRateLimiter.requests[key] = append(validTimes, now)

		// 定期清理全局 map 防止内存泄漏
		if len(globalRateLimiter.requests) > 10000 {
			for k, v := range globalRateLimiter.requests {
				if len(v) == 0 || now.Sub(v[len(v)-1]) > globalRateLimiter.window {
					delete(globalRateLimiter.requests, k)
				}
			}
		}

		globalRateLimiter.mu.Unlock()
		c.Next()
	}
}

// Error429 返回请求过于频繁错误
func Error429(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{
		"success":   false,
		"code":      "TOO_MANY_REQUESTS",
		"message":   "请求过于频繁，请稍后再试",
		"requestId": c.GetString("request_id"),
	})
}

// GetUserID 获取当前用户 ID
func GetUserID(c *gin.Context) uuid.UUID {
	id, _ := c.Get("user_id")
	if uid, ok := id.(uuid.UUID); ok {
		return uid
	}
	return uuid.Nil
}

// GetUserEmail 获取当前用户邮箱
func GetUserEmail(c *gin.Context) string {
	email, _ := c.Get("user_email")
	if e, ok := email.(string); ok {
		return e
	}
	return ""
}

// GetUserRoles 获取当前用户角色
func GetUserRoles(c *gin.Context) []string {
	roles, _ := c.Get("user_roles")
	if r, ok := roles.([]string); ok {
		return r
	}
	return []string{}
}

// GetClientIP 获取客户端 IP
func GetClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			ip = strings.TrimSpace(ips[0])
		}
	}
	return ip
}

// GetUserAgent 获取用户代理
func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

// AbortWithError 中止请求并返回错误
func AbortWithError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"success":   false,
		"code":      code,
		"message":   message,
		"requestId": c.GetString("request_id"),
	})
}
