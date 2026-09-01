package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sportswear-backend/internal/middleware"
	"sportswear-backend/internal/services"
	"sportswear-backend/internal/utils"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	ip := middleware.GetClientIP(c)
	resp, err := h.authService.Login(&req, ip)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 会话闭环：同时下发 HttpOnly Cookie（浏览器可见，
	// Authorization 丢失时中间件可回退读取，保障登录态不因 localStorage 清理而丢失）
	setSessionCookie(c, resp.Token, h.authService.SessionMaxAge())

	utils.Success(c, resp)
}

// setSessionCookie 下发会话 Cookie（HttpOnly，浏览器不可被 JS 读取）
func setSessionCookie(c *gin.Context, token string, maxAgeSeconds int) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.SessionCookieName, token, maxAgeSeconds, "/", "", false, true)
}

// clearSessionCookie 清除会话 Cookie
func clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.SessionCookieName, "", -1, "/", "", false, true)
}

// Logout 登出：清除会话 Cookie（无需认证；前端同时清理本地 token）
func (h *AuthHandler) Logout(c *gin.Context) {
	clearSessionCookie(c)
	utils.Success(c, gin.H{"logged_out": true})
}

// GetProfile 获取当前用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}
	utils.Success(c, user)
}

// ListUsers 用户列表
func (h *AuthHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	users, total, err := h.authService.ListUsers(page, pageSize, keyword)
	if err != nil {
		utils.InternalError(c, "获取用户列表失败")
		return
	}
	utils.SuccessPage(c, users, page, pageSize, total)
}

// CreateUser 创建用户
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.CreateUser(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, user)
}

// UpdateUser 更新用户
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.UpdateUser(c.Param("id"), &req, middleware.GetUserID(c))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, user)
}

// DeleteUser 删除用户
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	if err := h.authService.DeleteUser(c.Param("id"), middleware.GetUserID(c)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

// ListRoles 角色列表（全量，供角色下拉/分配使用）
func (h *AuthHandler) ListRoles(c *gin.Context) {
	roles, err := h.authService.ListRoles()
	if err != nil {
		utils.InternalError(c, "获取角色列表失败")
		return
	}
	utils.Success(c, roles)
}

// ListRolesPage 角色分页列表（支持名称/标识关键字模糊查询）
func (h *AuthHandler) ListRolesPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	roles, total, err := h.authService.ListRolesPage(page, pageSize, keyword)
	if err != nil {
		utils.InternalError(c, "获取角色列表失败")
		return
	}
	utils.SuccessPage(c, roles, page, pageSize, total)
}

// CreateRole 创建角色
func (h *AuthHandler) CreateRole(c *gin.Context) {
	var req services.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.authService.CreateRole(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, role)
}

// UpdateRole 更新角色（名称/标识/描述/启用状态/权限）
func (h *AuthHandler) UpdateRole(c *gin.Context) {
	var req services.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.authService.UpdateRole(c.Param("id"), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, role)
}

// DeleteRole 删除角色
func (h *AuthHandler) DeleteRole(c *gin.Context) {
	if err := h.authService.DeleteRole(c.Param("id")); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"deleted": true})
}

// ListPermissions 权限列表
func (h *AuthHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.authService.ListPermissions()
	if err != nil {
		utils.InternalError(c, "获取权限列表失败")
		return
	}
	utils.Success(c, permissions)
}

// ==================== 应用维度管理 ====================

// ListApps 应用列表
func (h *AuthHandler) ListApps(c *gin.Context) {
	apps, err := h.authService.ListApps()
	if err != nil {
		utils.InternalError(c, "获取应用列表失败")
		return
	}
	utils.Success(c, apps)
}

// CreateApp 创建应用
func (h *AuthHandler) CreateApp(c *gin.Context) {
	var req services.CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	app, err := h.authService.CreateApp(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, app)
}

// UpdateAppStatus 更新应用启用/禁用状态
func (h *AuthHandler) UpdateAppStatus(c *gin.Context) {
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.authService.UpdateAppStatus(c.Param("id"), req.IsActive); err != nil {
		utils.BadRequest(c, "更新应用状态失败")
		return
	}
	utils.Success(c, gin.H{"is_active": req.IsActive})
}

// AssignUserToApp 将用户授权到应用
func (h *AuthHandler) AssignUserToApp(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	appUser, err := h.authService.AssignUserToApp(c.Param("id"), req.UserID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, appUser)
}

// RemoveUserFromApp 将用户从应用移除
func (h *AuthHandler) RemoveUserFromApp(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.authService.RemoveUserFromApp(c.Param("id"), req.UserID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"removed": true})
}

// UpdateRoleStatus 更新角色启用/禁用状态
func (h *AuthHandler) UpdateRoleStatus(c *gin.Context) {
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.authService.UpdateRoleStatus(c.Param("id"), req.IsActive); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, gin.H{"is_active": req.IsActive})
}
