package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/database"
	"sportswear-backend/internal/models"
)

// ==================== 通知服务 ====================

// NotificationService 通知服务
type NotificationService struct {
	db *gorm.DB
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// Create 创建通知
func (s *NotificationService) Create(notification *models.Notification) error {
	return s.db.Create(notification).Error
}

// CreateForAllAdmins 为所有管理员创建通知（新询盘提醒）
func (s *NotificationService) NotifyAdmins(notifType models.NotificationType, title, content string, entityType string, entityID *uuid.UUID) error {
	var adminUserIDs []uuid.UUID
	// 查找拥有 lead:view 权限的角色对应的用户（简化：查找所有有角色的用户）
	s.db.Table("users").
		Select("users.id").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("users.status = ?", "active").
		Distinct().
		Scan(&adminUserIDs)

	for _, uid := range adminUserIDs {
		uidCopy := uid
		notif := models.Notification{
			Type:       notifType,
			Title:      title,
			Content:    content,
			UserID:     &uidCopy,
			EntityType: entityType,
			EntityID:   entityID,
		}
		s.db.Create(&notif)
	}
	return nil
}

// List 通知列表（按用户）
func (s *NotificationService) List(userID uuid.UUID, page, pageSize int, onlyUnread bool) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := s.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	if onlyUnread {
		query = query.Where("is_read = ?", false)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&notifications).Error
	return notifications, total, err
}

// UnreadCount 未读数量
func (s *NotificationService) UnreadCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

// MarkRead 标记已读
func (s *NotificationService) MarkRead(id string, userID uuid.UUID) error {
	now := time.Now()
	result := s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": now})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("通知不存在")
	}
	return nil
}

// MarkAllRead 全部标记已读
func (s *NotificationService) MarkAllRead(userID uuid.UUID) error {
	now := time.Now()
	return s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error
}

// ==================== 操作日志服务 ====================

// OperationLogService 操作日志服务
type OperationLogService struct {
	db *gorm.DB
}

// NewOperationLogService 创建操作日志服务
func NewOperationLogService(db *gorm.DB) *OperationLogService {
	return &OperationLogService{db: db}
}

// Record 记录操作日志（按月分表写入 operation_logs_YYYYMM）
func (s *OperationLogService) Record(log *models.OperationLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	tableName := database.MonthTableName("operation_logs", log.CreatedAt)
	if err := database.EnsureTable(s.db, tableName, &models.OperationLog{}); err != nil {
		return err
	}
	return s.db.Table(tableName).Create(log).Error
}

// resolveTables 解析查询涉及的表：存量基础表（未分表前数据）+ 时间范围内的月度表（仅已存在）
func (s *OperationLogService) resolveTables(start, end time.Time) ([]string, error) {
	if start.IsZero() {
		start = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	if end.IsZero() {
		end = time.Now()
	}
	if end.Before(start) {
		start, end = end, start
	}
	candidates := database.MonthTablesBetween("operation_logs", start, end)
	existing, err := database.ExistingTables(s.db, "operation_logs")
	if err != nil {
		return nil, err
	}
	exMap := make(map[string]bool, len(existing))
	for _, n := range existing {
		exMap[n] = true
	}
	tables := []string{"operation_logs"} // 存量基础表（历史数据）
	for _, n := range candidates {
		if exMap[n] {
			tables = append(tables, n)
		}
	}
	return tables, nil
}

// List 操作日志列表（按时间范围组合查询对应季度分表）
func (s *OperationLogService) List(page, pageSize int, userID, module, operation, entityType, ip string, start, end time.Time) ([]map[string]interface{}, int64, error) {
	tables, err := s.resolveTables(start, end)
	if err != nil {
		return nil, 0, err
	}

	var conds []string
	var args []interface{}
	if userID != "" {
		conds = append(conds, "user_id = ?")
		args = append(args, userID)
	}
	if module != "" {
		conds = append(conds, "module = ?")
		args = append(args, module)
	}
	if operation != "" {
		conds = append(conds, "operation = ?")
		args = append(args, operation)
	}
	if entityType != "" {
		conds = append(conds, "entity_type = ?")
		args = append(args, entityType)
	}
	if ip != "" {
		conds = append(conds, "ip = ?")
		args = append(args, ip)
	}
	if !start.IsZero() {
		conds = append(conds, "created_at >= ?")
		args = append(args, start)
	}
	if !end.IsZero() {
		conds = append(conds, "created_at <= ?")
		args = append(args, end)
	}
	where := ""
	if len(conds) > 0 {
		where = " WHERE " + strings.Join(conds, " AND ")
	}

	selectFmt := `SELECT "id","user_id","operation","module","entity_type","entity_id","before","after","ip","user_agent","description","created_at","updated_at","deleted_at",'__T__' AS src FROM __T__` + where
	unionSQL, unionArgs := database.UnionAll(tables, selectFmt, args)

	var total int64
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM (%s) t`, unionSQL)
	if err := s.db.Raw(countSQL, unionArgs...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	listArgs := make([]interface{}, 0, len(unionArgs)+2)
	listArgs = append(listArgs, unionArgs...)
	listArgs = append(listArgs, pageSize, (page-1)*pageSize)
	listSQL := fmt.Sprintf(`SELECT * FROM (%s) t ORDER BY created_at DESC LIMIT ? OFFSET ?`, unionSQL)

	var rows []map[string]interface{}
	if err := s.db.Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ==================== 报价服务 ====================

// QuoteService 报价服务
type QuoteService struct {
	db *gorm.DB
}

// NewQuoteService 创建报价服务
func NewQuoteService(db *gorm.DB) *QuoteService {
	return &QuoteService{db: db}
}

// ListQuotes 报价列表
func (s *QuoteService) ListQuotes(page, pageSize int, leadID, status string) ([]models.Quote, int64, error) {
	var quotes []models.Quote
	var total int64

	query := s.db.Model(&models.Quote{})
	if leadID != "" {
		query = query.Where("lead_id = ?", leadID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&quotes).Error
	return quotes, total, err
}

// GetQuote 获取报价
func (s *QuoteService) GetQuote(id string) (*models.Quote, error) {
	var quote models.Quote
	if err := s.db.First(&quote, "id = ?", id).Error; err != nil {
		return nil, errors.New("报价不存在")
	}
	return &quote, nil
}

// CreateQuote 创建报价
func (s *QuoteService) CreateQuote(req *QuoteRequest) (*models.Quote, error) {
	var quote models.Quote
	s.db.Where("lead_id = ?", req.LeadID).Order("created_at DESC").First(&quote)
	seq := len(s.listQuoteNumbers(req.LeadID)) + 1

	quote = models.Quote{
		LeadID:         req.LeadID,
		QuoteNumber:    fmt.Sprintf("Q-%s-%03d", time.Now().Format("20060102"), seq),
		ProductID:      req.ProductID,
		Quantity:       req.Quantity,
		UnitPrice:      req.UnitPrice,
		TotalPrice:     req.TotalPrice,
		Currency:       req.Currency,
		MOQ:            req.MOQ,
		LeadTime:       req.LeadTime,
		PaymentTerms:   req.PaymentTerms,
		ShippingMethod: req.ShippingMethod,
		ValidUntil:     req.ValidUntil,
		Status:         req.Status,
		Notes:          req.Notes,
	}
	if quote.Status == "" {
		quote.Status = "draft"
	}
	if err := s.db.Create(&quote).Error; err != nil {
		return nil, err
	}

	// 创建通知
	var notifErr error
	_ = notifErr

	return &quote, nil
}

func (s *QuoteService) listQuoteNumbers(leadID uuid.UUID) []string {
	var nums []string
	s.db.Model(&models.Quote{}).Where("lead_id = ?", leadID).Pluck("quote_number", &nums)
	return nums
}

// UpdateQuote 更新报价
func (s *QuoteService) UpdateQuote(id string, req *QuoteRequest) (*models.Quote, error) {
	var quote models.Quote
	if err := s.db.First(&quote, "id = ?", id).Error; err != nil {
		return nil, errors.New("报价不存在")
	}

	updates := map[string]interface{}{
		"quantity":        req.Quantity,
		"unit_price":      req.UnitPrice,
		"total_price":     req.TotalPrice,
		"currency":        req.Currency,
		"moq":             req.MOQ,
		"lead_time":       req.LeadTime,
		"payment_terms":   req.PaymentTerms,
		"shipping_method": req.ShippingMethod,
		"valid_until":     req.ValidUntil,
		"status":          req.Status,
		"notes":           req.Notes,
	}
	if err := s.db.Model(&quote).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetQuote(id)
}

// DeleteQuote 删除报价
func (s *QuoteService) DeleteQuote(id string) error {
	return s.db.Delete(&models.Quote{}, "id = ?", id).Error
}

// QuoteRequest 报价请求
type QuoteRequest struct {
	LeadID         uuid.UUID  `json:"lead_id" binding:"required"`
	ProductID      *uuid.UUID `json:"product_id"`
	Quantity       int        `json:"quantity"`
	UnitPrice      float64    `json:"unit_price"`
	TotalPrice     float64    `json:"total_price"`
	Currency       string     `json:"currency"`
	MOQ            int        `json:"moq"`
	LeadTime       string     `json:"lead_time"`
	PaymentTerms   string     `json:"payment_terms"`
	ShippingMethod string     `json:"shipping_method"`
	ValidUntil     *time.Time `json:"valid_until"`
	Status         string     `json:"status"`
	Notes          string     `json:"notes"`
}
