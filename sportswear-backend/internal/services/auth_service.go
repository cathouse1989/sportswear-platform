package services

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/config"
	"sportswear-backend/internal/models"
	"sportswear-backend/internal/utils"
)

// AuthService 认证服务
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest, ip string) (*LoginResponse, error) {
	// 去除首尾空白：防止直接调用接口/前端遗漏时，复制粘贴带入空格导致「邮箱或密码错误」
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	if user.Status != "active" {
		return nil, errors.New("账号已被禁用")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("邮箱或密码错误")
	}

	// 加载角色(仅启用的角色)
	var roles []models.Role
	if err := s.db.Distinct().
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.is_active = ?", user.ID, true).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	// 生成 Token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Name, roleCodes, s.cfg.JWT.Secret, s.cfg.JWT.ExpireHours)
	if err != nil {
		return nil, err
	}

	// 更新登录信息
	now := time.Now()
	s.db.Model(&user).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	})

	// 隐藏密码
	user.Password = ""

	return &LoginResponse{
		Token: token,
		User:  &user,
	}, nil
}

// GetUserByID 获取用户
// 关键闭环：profile 必须返回「启用角色 + 角色权限」，前端 authStore 才能收集到权限码，
// 与后端 Auth 中间件(实时从数据库加载权限)保持一致。
// SessionMaxAge 会话 Cookie 有效期(秒)，与 JWT 过期时间对齐
func (s *AuthService) SessionMaxAge() int {
	return s.cfg.JWT.ExpireHours * 3600
}

func (s *AuthService) GetUserByID(id interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.
		Preload("Roles", "is_active = ?", true).
		Preload("Roles.Permissions").
		First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}

// CreateUser 创建用户
func (s *AuthService) CreateUser(req *CreateUserRequest) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Phone:    req.Phone,
		Status:   "active",
	}

	// 创建时支持直接指定账号状态(仅接受 active/disabled，默认 active)
	if req.Status != "" {
		if req.Status == "active" || req.Status == "disabled" {
			user.Status = req.Status
		}
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		if err := s.db.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
			return nil, err
		}
		s.db.Model(&user).Association("Roles").Replace(roles)
	}

	user.Password = ""
	return &user, nil
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Name     string   `json:"name" binding:"required"`
	Phone    string   `json:"phone"`
	Status   string   `json:"status"`   // active / disabled，可选，默认 active
	RoleIDs  []string `json:"role_ids"` // 分配的角色 ID 列表(账号-角色闭环)
}

// UpdateUser 更新用户
// operatorID：当前登录用户 ID，用于账号自保护(不能被自己禁掉导致系统锁死)。
func (s *AuthService) UpdateUser(id string, req *UpdateUserRequest, operatorID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// ---------- 账号自保护(防止系统锁死) ----------
	// 1. 不能禁用当前登录账号
	if req.Status == "disabled" && operatorID != uuid.Nil && user.ID == operatorID {
		return nil, errors.New("不能禁用当前登录账号")
	}
	// 2. 最后一个超级管理员：不能禁用、不能移除其 super_admin 角色
	if s.isSuperAdminAccount(user.ID) && s.superAdminAccountCount() <= 1 {
		if req.Status == "disabled" {
			return nil, errors.New("不能禁用最后一个超级管理员账号")
		}
		if req.RoleIDs != nil && !containsRoleCode(req.RoleIDs, s.superAdminRoleID()) {
			return nil, errors.New("不能移除最后一个超级管理员的管理员角色")
		}
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		updates["password"] = hashed
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 更新角色
	if req.RoleIDs != nil {
		var roles []models.Role
		if err := s.db.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
			return nil, err
		}
		s.db.Model(&user).Association("Roles").Replace(roles)
	}

	user.Password = ""
	return &user, nil
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Status   string   `json:"status"`
	Password string   `json:"password"`
	RoleIDs  []string `json:"role_ids"`
}

// DeleteUser 删除用户
// operatorID：当前登录用户 ID，防止删除自己 / 最后一个超级管理员导致系统锁死。
func (s *AuthService) DeleteUser(id string, operatorID uuid.UUID) error {
	if operatorID != uuid.Nil && operatorID.String() == id {
		return errors.New("不能删除当前登录账号")
	}

	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return errors.New("用户不存在")
	}
	if s.isSuperAdminAccount(user.ID) && s.superAdminAccountCount() <= 1 {
		return errors.New("不能删除最后一个超级管理员账号")
	}

	return s.db.Delete(&models.User{}, "id = ?", id).Error
}

// ==================== 账号自保护辅助 ====================

// isSuperAdminAccount 判断用户是否拥有 super_admin 角色
func (s *AuthService) isSuperAdminAccount(userID uuid.UUID) bool {
	var count int64
	s.db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.code = ? AND roles.deleted_at IS NULL", userID, "super_admin").
		Count(&count)
	return count > 0
}

// superAdminAccountCount 拥有 super_admin 角色的有效用户数
func (s *AuthService) superAdminAccountCount() int64 {
	var count int64
	s.db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Joins("JOIN users ON users.id = user_roles.user_id").
		Where("roles.code = ? AND users.status = ? AND roles.deleted_at IS NULL AND users.deleted_at IS NULL", "super_admin", "active").
		Count(&count)
	return count
}

// superAdminRoleID 获取 super_admin 角色的 ID(不存在返回零值)
func (s *AuthService) superAdminRoleID() uuid.UUID {
	var role models.Role
	if err := s.db.Where("code = ?", "super_admin").First(&role).Error; err != nil {
		return uuid.Nil
	}
	return role.ID
}

// containsRoleCode 判断 ID 列表中是否包含目标角色 ID(字符串比对)
func containsRoleCode(roleIDs []string, target uuid.UUID) bool {
	if target == uuid.Nil {
		return false
	}
	for _, id := range roleIDs {
		if id == target.String() {
			return true
		}
	}
	return false
}

// ListUsers 用户列表
func (s *AuthService) ListUsers(page, pageSize int, keyword string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})
	if keyword != "" {
		query = query.Where("email LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Preload("Roles").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&users).Error

	for i := range users {
		users[i].Password = ""
	}

	return users, total, err
}

// ListRoles 角色列表(全量，供角色下拉/分配使用)
func (s *AuthService) ListRoles() ([]models.Role, error) {
	var roles []models.Role
	err := s.db.Preload("Permissions").Order("created_at ASC").Find(&roles).Error
	return roles, err
}

// ListRolesPage 角色分页列表(支持名称/标识/描述关键字模糊查询)
func (s *AuthService) ListRolesPage(page, pageSize int, keyword string) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	query := s.db.Model(&models.Role{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Permissions").Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at ASC").Find(&roles).Error
	return roles, total, err
}

// CreateRole 创建角色
func (s *AuthService) CreateRole(req *CreateRoleRequest) (*models.Role, error) {
	// 角色标识唯一性检查(友好报错)
	var count int64
	s.db.Model(&models.Role{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return nil, errors.New("角色标识已存在: " + req.Code)
	}

	role := models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, err
	}

	if len(req.PermissionIDs) > 0 {
		var permissions []models.Permission
		if err := s.db.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
			return nil, err
		}
		s.db.Model(&role).Association("Permissions").Replace(permissions)
	}

	return &role, nil
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name          string   `json:"name" binding:"required"`
	Code          string   `json:"code" binding:"required"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

// UpdateRoleRequest 更新角色请求(支持部分更新：字段非空才更新)
// - permissions:  权限码数组(前端角色权限树提交的是 code)
// - permission_ids: 权限 ID 数组(兼容创建角色参数风格)
type UpdateRoleRequest struct {
	Name          string   `json:"name"`
	Code          string   `json:"code"`
	Description   string   `json:"description"`
	IsActive      *bool    `json:"is_active"`
	Permissions   []string `json:"permissions"`
	PermissionIDs []string `json:"permission_ids"`
}

// UpdateRole 更新角色基本信息与权限
func (s *AuthService) UpdateRole(id string, req *UpdateRoleRequest) (*models.Role, error) {
	var role models.Role
	if err := s.db.First(&role, "id = ?", id).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	// 内置超级管理员保护：不能改名/改标识/禁用(否则系统会锁死)
	if role.Code == "super_admin" {
		if req.Code != "" && req.Code != "super_admin" {
			return nil, errors.New("不能修改超级管理员角色的标识")
		}
		if req.IsActive != nil && !*req.IsActive {
			return nil, errors.New("不能禁用超级管理员角色")
		}
	}

	// 更新基础字段(仅非空字段)
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		// 角色标识唯一性检查(排除自身)
		var count int64
		s.db.Model(&models.Role{}).Where("code = ? AND id <> ?", req.Code, role.ID).Count(&count)
		if count > 0 {
			return nil, errors.New("角色标识已存在: " + req.Code)
		}
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if len(updates) > 0 {
		if err := s.db.Model(&role).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// 更新权限：permissions(权限码)与 permission_ids(权限 ID)互斥，两者显式声明才生效(支持清空)
	if req.PermissionIDs != nil || req.Permissions != nil {
		var ids []string
		if req.PermissionIDs != nil {
			ids = req.PermissionIDs
		} else {
			var permissions []models.Permission
			if err := s.db.Where("code IN ?", req.Permissions).Find(&permissions).Error; err != nil {
				return nil, err
			}
			for _, p := range permissions {
				ids = append(ids, p.ID.String())
			}
		}
		if err := s.replaceRolePermissions(&role, ids); err != nil {
			return nil, err
		}
	}

	// 返回最新角色(含权限)
	if err := s.db.Preload("Permissions").First(&role, "id = ?", role.ID).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// replaceRolePermissions 替换角色权限关联(先清空再按 ID 重建)
func (s *AuthService) replaceRolePermissions(role *models.Role, permissionIDs []string) error {
	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		if err := s.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}
	}
	// 空列表 = 清空该角色的所有权限
	return s.db.Model(role).Association("Permissions").Replace(permissions)
}

// UpdateRoleStatus 更新角色启用/禁用状态
// 禁用后：该角色用户登录时不再获得该角色，权限立即失效
func (s *AuthService) UpdateRoleStatus(id string, isActive bool) error {
	var role models.Role
	if err := s.db.First(&role, "id = ?", id).Error; err != nil {
		return errors.New("角色不存在")
	}
	if role.Code == "super_admin" && !isActive {
		return errors.New("不能禁用超级管理员角色")
	}
	return s.db.Model(&role).Update("is_active", isActive).Error
}

// DeleteRole 删除角色
// - super_admin 内置角色受保护不可删除(否则系统可能锁死)
// - 事务级联清理 role_permissions / user_roles 关联(保证一致性)
func (s *AuthService) DeleteRole(id string) error {
	var role models.Role
	if err := s.db.First(&role, "id = ?", id).Error; err != nil {
		return errors.New("角色不存在")
	}
	if role.Code == "super_admin" {
		return errors.New("不能删除超级管理员角色")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", role.ID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Role{}, "id = ?", role.ID).Error
	})
}

// ListPermissions 权限列表
func (s *AuthService) ListPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := s.db.Order("module ASC, name ASC").Find(&permissions).Error
	return permissions, err
}

// GetUserPermissions 获取用户所有权限码(基于角色)
func (s *AuthService) GetUserPermissions(userID string) ([]string, error) {
	var permissions []models.Permission
	err := s.db.Distinct().
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, len(permissions))
	for _, p := range permissions {
		codes = append(codes, p.Code)
	}
	return codes, nil
}

// GetUserRoleCodes 获取用户角色码(仅启用状态的角色)
func (s *AuthService) GetUserRoleCodes(userID string) ([]string, error) {
	var roles []models.Role
	err := s.db.Distinct().
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.is_active = ?", userID, true).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, len(roles))
	for _, r := range roles {
		codes = append(codes, r.Code)
	}
	return codes, nil
}

// ==================== 应用维度管理 ====================

// ListApps 应用列表
func (s *AuthService) ListApps() ([]models.App, error) {
	var apps []models.App
	err := s.db.Order("sort_order ASC").Find(&apps).Error
	return apps, err
}

// CreateAppRequest 创建应用请求
type CreateAppRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	CallbackURL string `json:"callback_url"`
	Desc        string `json:"desc"`
	SortOrder   int    `json:"sort_order"`
}

// CreateApp 创建应用(自动生成 API Key)
func (s *AuthService) CreateApp(req *CreateAppRequest) (*models.App, error) {
	app := models.App{
		Code:        req.Code,
		Name:        req.Name,
		APIKey:      "sk_" + utils.GenerateRandomString(32),
		Secret:      utils.GenerateRandomString(48),
		CallbackURL: req.CallbackURL,
		Desc:        req.Desc,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}
	if err := s.db.Create(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// UpdateAppStatus 更新应用启用/禁用状态
func (s *AuthService) UpdateAppStatus(id string, isActive bool) error {
	return s.db.Model(&models.App{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// GetAppByAPIKey 通过 API Key 获取应用
func (s *AuthService) GetAppByAPIKey(apiKey string) (*models.App, error) {
	var app models.App
	err := s.db.Where("api_key = ? AND is_active = ?", apiKey, true).First(&app).Error
	if err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}
	return &app, nil
}

// AssignUserToApp 将用户授权到应用
func (s *AuthService) AssignUserToApp(appID, userID string) (*models.AppUser, error) {
	aID, err := uuid.Parse(appID)
	if err != nil {
		return nil, errors.New("应用 ID 无效")
	}
	uID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("用户 ID 无效")
	}

	// 检查是否已存在
	var count int64
	s.db.Model(&models.AppUser{}).Where("app_id = ? AND user_id = ?", aID, uID).Count(&count)
	if count > 0 {
		return nil, errors.New("用户已在该应用中")
	}

	appUser := models.AppUser{
		AppID:  aID,
		UserID: uID,
		Status: "active",
	}
	if err := s.db.Create(&appUser).Error; err != nil {
		return nil, err
	}
	return &appUser, nil
}

// RemoveUserFromApp 将用户从应用移除
func (s *AuthService) RemoveUserFromApp(appID, userID string) error {
	aID, err := uuid.Parse(appID)
	if err != nil {
		return errors.New("应用 ID 无效")
	}
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("用户 ID 无效")
	}
	return s.db.Where("app_id = ? AND user_id = ?", aID, uID).Delete(&models.AppUser{}).Error
}

// defaultPermissions 默认权限清单(与 docs/权限矩阵.md 保持一致；新增权限码在此补充即可幂等入库)
func (s *AuthService) defaultPermissions() []models.Permission {
	return []models.Permission{
		{Name: "用户查看", Code: "user:view", Module: "user"},
		{Name: "用户新增", Code: "user:create", Module: "user"},
		{Name: "用户修改", Code: "user:update", Module: "user"},
		{Name: "用户删除", Code: "user:delete", Module: "user"},
		{Name: "角色管理", Code: "role:manage", Module: "user"},
		{Name: "产品查看", Code: "product:view", Module: "product"},
		{Name: "产品新增", Code: "product:create", Module: "product"},
		{Name: "产品修改", Code: "product:update", Module: "product"},
		{Name: "产品删除", Code: "product:delete", Module: "product"},
		{Name: "产品发布", Code: "product:publish", Module: "product"},
		{Name: "分类管理", Code: "category:manage", Module: "product"},
		{Name: "系列管理", Code: "series:manage", Module: "product"},
		{Name: "面料管理", Code: "fabric:manage", Module: "product"},
		{Name: "页面查看", Code: "page:view", Module: "cms"},
		{Name: "页面修改", Code: "page:update", Module: "cms"},
		{Name: "页面发布", Code: "page:publish", Module: "cms"},
		{Name: "导航管理", Code: "navigation:manage", Module: "cms"},
		{Name: "博客管理", Code: "blog:manage", Module: "cms"},
		{Name: "案例管理", Code: "case:manage", Module: "cms"},
		{Name: "FAQ管理", Code: "faq:manage", Module: "cms"},
		{Name: "工厂管理", Code: "factory:manage", Module: "cms"},
		{Name: "认证管理", Code: "certification:manage", Module: "cms"},
		{Name: "生产流程", Code: "production:manage", Module: "cms"},
		{Name: "自媒体管理", Code: "selfmedia:manage", Module: "cms"},
		{Name: "媒体上传", Code: "media:upload", Module: "media"},
		{Name: "媒体管理", Code: "media:manage", Module: "media"},
		{Name: "询盘查看", Code: "lead:view", Module: "lead"},
		{Name: "询盘修改", Code: "lead:update", Module: "lead"},
		{Name: "询盘跟进", Code: "lead:followup", Module: "lead"},
		{Name: "订阅查看", Code: "subscription:view", Module: "lead"},
		{Name: "订阅管理", Code: "subscription:update", Module: "lead"},
		{Name: "SEO管理", Code: "seo:manage", Module: "seo"},
		{Name: "语言管理", Code: "language:manage", Module: "system"},
		{Name: "系统配置", Code: "setting:manage", Module: "system"},
		{Name: "操作日志", Code: "log:view", Module: "system"},
		{Name: "数据统计", Code: "dashboard:view", Module: "system"},
	}
}

// ensureDefaultPermissions 幂等补齐默认权限：按 code 逐条检查，缺失即创建(已有数据库升级时也会补齐)
func (s *AuthService) ensureDefaultPermissions() {
	for _, p := range s.defaultPermissions() {
		var count int64
		s.db.Model(&models.Permission{}).Where("code = ?", p.Code).Count(&count)
		if count == 0 {
			s.db.Create(&p)
		}
	}
}

// ensureDefaultRolePermissions 幂等补齐默认角色的扩展权限(仅追加缺失项，不覆盖已有自定义权限)
func (s *AuthService) ensureDefaultRolePermissions() {
	extra := map[string][]string{
		"admin":         {"production:manage", "subscription:view", "subscription:update", "selfmedia:manage"},
		"content_admin": {"production:manage", "subscription:view", "selfmedia:manage"},
		"sales_manager": {"subscription:view"},
		"sales":         {"subscription:view"},
		"viewer":        {"subscription:view"},
	}
	for roleCode, codes := range extra {
		var role models.Role
		if err := s.db.Preload("Permissions").Where("code = ?", roleCode).First(&role).Error; err != nil {
			continue
		}
		var perms []models.Permission
		if err := s.db.Where("code IN ?", codes).Find(&perms).Error; err != nil || len(perms) == 0 {
			continue
		}
		existing := make(map[string]bool, len(role.Permissions))
		for _, p := range role.Permissions {
			existing[p.Code] = true
		}
		var toAppend []models.Permission
		for _, p := range perms {
			if !existing[p.Code] {
				toAppend = append(toAppend, p)
			}
		}
		if len(toAppend) > 0 {
			s.db.Model(&role).Association("Permissions").Append(toAppend)
		}
	}
}

// InitDefaultData 初始化默认数据
func (s *AuthService) InitDefaultData() error {
	// 创建默认语言
	var langCount int64
	s.db.Model(&models.Language{}).Count(&langCount)
	if langCount == 0 {
		languages := []models.Language{
			{Code: "en", Name: "English", NativeName: "English", IsDefault: true, IsActive: true, SortOrder: 1, Flag: "🇬🇧"},
			{Code: "zh", Name: "Chinese", NativeName: "中文", IsActive: true, SortOrder: 2, Flag: "🇨🇳"},
			{Code: "es", Name: "Spanish", NativeName: "Español", IsActive: false, SortOrder: 3, Flag: "🇪🇸"},
			{Code: "fr", Name: "French", NativeName: "Français", IsActive: false, SortOrder: 4, Flag: "🇫🇷"},
			{Code: "de", Name: "German", NativeName: "Deutsch", IsActive: false, SortOrder: 5, Flag: "🇩🇪"},
		}
		s.db.Create(&languages)
	}

	// 创建默认权限(幂等：按权限码逐条 FirstOrCreate，已有数据库也能补齐新增权限码)
	s.ensureDefaultPermissions()

	// 创建默认角色并分配权限
	var roleCount int64
	s.db.Model(&models.Role{}).Count(&roleCount)
	if roleCount == 0 {
		roles := []models.Role{
			{Name: "超级管理员", Code: "super_admin", Description: "拥有所有权限"},
			{Name: "管理员", Code: "admin", Description: "系统管理员"},
			{Name: "内容管理员", Code: "content_admin", Description: "管理网站内容"},
			{Name: "产品管理员", Code: "product_admin", Description: "管理产品"},
			{Name: "SEO管理员", Code: "seo_admin", Description: "管理SEO"},
			{Name: "销售经理", Code: "sales_manager", Description: "销售经理"},
			{Name: "销售人员", Code: "sales", Description: "销售人员"},
			{Name: "编辑人员", Code: "editor", Description: "内容编辑"},
			{Name: "只读用户", Code: "viewer", Description: "只读访问"},
		}
		s.db.Create(&roles)
	}

	// 角色权限关联为空时，为各角色分配权限(兼容已有数据库)
	var rolePermCount int64
	s.db.Model(&models.RolePermission{}).Count(&rolePermCount)
	if rolePermCount == 0 {
		s.assignRolePermissions()
	}

	// 幂等补齐默认角色的扩展权限(如 production:manage)，仅追加、不覆盖已有自定义权限
	s.ensureDefaultRolePermissions()

	// 创建默认导航
	var navCount int64
	s.db.Model(&models.Navigation{}).Count(&navCount)
	if navCount == 0 {
		// header 与 footer 顶级导航顺序保持一致：
		// 首页 / 产品中心 / 案例展示 / 关于我们 / 博客 / 常见问题 / 联系我们
		navigations := []models.Navigation{
			{Name: "Home", Type: "header", URL: "/", SortOrder: 1, IsVisible: true},
			{Name: "Products", Type: "header", URL: "/products", SortOrder: 2, IsVisible: true},
			{Name: "Cases", Type: "header", URL: "/cases", SortOrder: 3, IsVisible: true},
			{Name: "About Us", Type: "header", URL: "/about", SortOrder: 4, IsVisible: true},
			{Name: "Blog", Type: "header", URL: "/blog", SortOrder: 5, IsVisible: true},
			{Name: "FAQ", Type: "header", URL: "/faq", SortOrder: 6, IsVisible: true},
			{Name: "Contact Us", Type: "header", URL: "/contact", SortOrder: 7, IsVisible: true},
			{Name: "Home", Type: "footer", URL: "/", SortOrder: 1, IsVisible: true},
			{Name: "Products", Type: "footer", URL: "/products", SortOrder: 2, IsVisible: true},
			{Name: "Cases", Type: "footer", URL: "/cases", SortOrder: 3, IsVisible: true},
			{Name: "About Us", Type: "footer", URL: "/about", SortOrder: 4, IsVisible: true},
			{Name: "Blog", Type: "footer", URL: "/blog", SortOrder: 5, IsVisible: true},
			{Name: "FAQ", Type: "footer", URL: "/faq", SortOrder: 6, IsVisible: true},
			{Name: "Contact Us", Type: "footer", URL: "/contact", SortOrder: 7, IsVisible: true},
		}
		s.db.Create(&navigations)
	}

	// 创建默认首页
	var homeCount int64
	s.db.Model(&models.Page{}).Where("slug = ?", "home").Count(&homeCount)
	if homeCount == 0 {
		now := time.Now()
		homePage := models.Page{
			Title:       "Home",
			Slug:        "home",
			Type:        models.PageTypeHome,
			Status:      models.ContentStatusPublished,
			SortOrder:   1,
			PublishedAt: &now,
		}
		if err := s.db.Create(&homePage).Error; err == nil {
			// 创建首页模块
			modules := []models.PageModule{
				{PageID: homePage.ID, Type: "banner", Title: "Hero Banner", SortOrder: 1, IsVisible: true,
					Config: `{"slides":[
						{"image":"https://images.unsplash.com/photo-1599901860904-17e6ed7083a0?w=1600&h=900&fit=crop","title":"Custom Sportswear Manufacturer","subtitle":"OEM & ODM Solutions for Global Brands","button_text":"Get a Quote","button_url":"/contact"},
						{"image":"https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=1600&h=900&fit=crop","title":"Advanced Running & Training Apparel","subtitle":"Technical fabrics for peak performance","button_text":"Get a Quote","button_url":"/contact"},
						{"image":"https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=1600&h=900&fit=crop","title":"Professional Team Uniforms","subtitle":"Full custom sublimation printing","button_text":"Get a Quote","button_url":"/contact"},
						{"image":"https://images.unsplash.com/photo-1577221084712-45b0445d2b00?w=1600&h=900&fit=crop","title":"From Concept to Product","subtitle":"End-to-end manufacturing support","button_text":"Get a Quote","button_url":"/contact"}
					]}`},
				{PageID: homePage.ID, Type: "product_recommend", Title: "Featured Products", SortOrder: 2, IsVisible: true,
					Config: `{"count":8,"source":"featured"}`},
				{PageID: homePage.ID, Type: "categories", Title: "Product Categories", SortOrder: 3, IsVisible: true,
					Config: `{}`},
				{PageID: homePage.ID, Type: "oem", Title: "OEM Service", SortOrder: 4, IsVisible: true,
					Config: `{"title":"OEM Manufacturing","description":"Your designs, our expertise. Full OEM service from tech pack to delivery."}`},
				{PageID: homePage.ID, Type: "odm", Title: "ODM Service", SortOrder: 5, IsVisible: true,
					Config: `{"title":"ODM Development","description":"From concept to product. Our design team develops your collection."}`},
				{PageID: homePage.ID, Type: "factory", Title: "Our Factory", SortOrder: 6, IsVisible: true,
					Config: `{}`},
				{PageID: homePage.ID, Type: "production", Title: "Production Process", SortOrder: 7, IsVisible: true,
					Config: `{}`},
				{PageID: homePage.ID, Type: "case", Title: "Client Cases", SortOrder: 8, IsVisible: true,
					Config: `{"count":4}`},
				{PageID: homePage.ID, Type: "certification", Title: "Certifications", SortOrder: 9, IsVisible: true,
					Config: `{}`},
				{PageID: homePage.ID, Type: "blog", Title: "Latest Articles", SortOrder: 10, IsVisible: true,
					Config: `{"count":4}`},
				{PageID: homePage.ID, Type: "contact", Title: "Contact Us", SortOrder: 11, IsVisible: true,
					Config: `{}`},
			}
			s.db.Create(&modules)

			// 创建首页英文翻译
			s.db.Create(&models.PageTranslation{
				PageID:   homePage.ID,
				Language: "en",
				Title:    "Home",
				Status:   models.TranslationStatusPublished,
			})
		}
	}

	// 创建默认管理员
	var adminCount int64
	s.db.Model(&models.User{}).Where("email = ?", s.cfg.Admin.Email).Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, _ := utils.HashPassword(s.cfg.Admin.Password)
		admin := models.User{
			Email:    s.cfg.Admin.Email,
			Password: hashedPassword,
			Name:     "超级管理员",
			Status:   "active",
		}
		if err := s.db.Create(&admin).Error; err != nil {
			return err
		}

		// 分配超级管理员角色
		var superAdmin models.Role
		if err := s.db.Where("code = ?", "super_admin").First(&superAdmin).Error; err == nil {
			s.db.Model(&admin).Association("Roles").Append(&superAdmin)
		}
	}

	return nil
}

// assignRolePermissions 为默认角色分配权限
func (s *AuthService) assignRolePermissions() {
	// 角色权限映射
	rolePermissions := map[string][]string{
		"admin": { // 管理员：除用户/角色管理外的所有权限
			"product:view", "product:create", "product:update", "product:delete", "product:publish",
			"category:manage", "series:manage", "fabric:manage",
			"page:view", "page:update", "page:publish", "navigation:manage",
			"blog:manage", "case:manage", "faq:manage",
			"factory:manage", "certification:manage", "production:manage",
			"selfmedia:manage",
			"media:upload", "media:manage",
			"lead:view", "lead:update", "lead:followup",
			"subscription:view", "subscription:update",
			"seo:manage", "dashboard:view",
		},
		"content_admin": { // 内容管理员：CMS 相关(页面管理只读，通过导航管理+轮播图管理操作)
			"page:view", "navigation:manage",
			"blog:manage", "case:manage", "faq:manage",
			"factory:manage", "certification:manage", "production:manage",
			"selfmedia:manage",
			"media:upload", "media:manage",
			"subscription:view",
		},
		"product_admin": { // 产品管理员：产品相关
			"product:view", "product:create", "product:update", "product:delete", "product:publish",
			"category:manage", "series:manage", "fabric:manage",
			"media:upload", "media:manage",
		},
		"seo_admin": { // SEO 管理员
			"page:view", "seo:manage",
			"product:view", "blog:manage", "case:manage",
		},
		"sales_manager": { // 销售经理：询盘全部权限 + 统计
			"lead:view", "lead:update", "lead:followup", "dashboard:view",
			"subscription:view",
			"product:view",
		},
		"sales": { // 销售人员：询盘查看和跟进
			"lead:view", "lead:followup",
			"subscription:view",
			"product:view",
		},
		"editor": { // 编辑人员：内容编辑
			"page:view", "page:update",
			"blog:manage", "case:manage", "faq:manage",
			"media:upload",
		},
		"viewer": { // 只读用户
			"product:view", "page:view", "lead:view", "subscription:view",
		},
	}

	for roleCode, permCodes := range rolePermissions {
		var role models.Role
		if err := s.db.Where("code = ?", roleCode).First(&role).Error; err != nil {
			continue
		}

		var permissions []models.Permission
		if err := s.db.Where("code IN ?", permCodes).Find(&permissions).Error; err != nil {
			continue
		}

		s.db.Model(&role).Association("Permissions").Replace(permissions)
	}
}
