package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

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

// Record 记录操作日志
func (s *OperationLogService) Record(log *models.OperationLog) error {
	return s.db.Create(log).Error
}

// List 操作日志列表
func (s *OperationLogService) List(page, pageSize int, userID, module, operation string) ([]models.OperationLog, int64, error) {
	var logs []models.OperationLog
	var total int64

	query := s.db.Model(&models.OperationLog{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if operation != "" {
		query = query.Where("operation = ?", operation)
	}

	query.Count(&total)
	err := query.Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error
	return logs, total, err
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
