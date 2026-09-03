package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
	"sportswear-backend/internal/utils"
)

// LeadService 询盘服务
type LeadService struct {
	db          *gorm.DB
	mailService *MailService
}

// NewLeadService 创建询盘服务
func NewLeadService(db *gorm.DB) *LeadService {
	return &LeadService{
		db:          db,
		mailService: NewMailService(),
	}
}

// CreateLead 创建询盘
func (s *LeadService) CreateLead(req *LeadRequest, ip, visitorID string) (*models.Lead, error) {
	lead := models.Lead{
		Name:            req.Name,
		Company:         req.Company,
		Email:           req.Email,
		Phone:           req.Phone,
		WhatsApp:        req.WhatsApp,
		Country:         req.Country,
		CompanyWebsite:  req.CompanyWebsite,
		ProjectType:     req.ProjectType,
		ProductCategory: req.ProductCategory,
		ProductID:       utils.StringPtrToUUIDPtr(req.ProductID),
		Quantity:        req.Quantity,
		Budget:          req.Budget,
		TargetDate:      req.TargetDate,
		Message:         req.Message,
		Attachments:     req.Attachments,
		Source:          req.Source,
		Medium:          req.Medium,
		Campaign:        req.Campaign,
		Keyword:         req.Keyword,
		FirstPage:       req.FirstPage,
		LandingPage:     req.LandingPage,
		Device:          req.Device,
		Language:        req.Language,
		VisitorID:       visitorID,
		IP:              ip,
		Status:          models.LeadStatusNew,
	}

	// 计算询盘评分
	lead.Score = s.calculateScore(req)
	lead.ScoreLevel = s.getScoreLevel(lead.Score)

	if err := s.db.Create(&lead).Error; err != nil {
		return nil, err
	}

	// 为新询盘创建通知和邮件提醒（异步）
	go s.notifyNewLead(&lead)

	return &lead, nil
}

// notifyNewLead 为新询盘创建通知
func (s *LeadService) notifyNewLead(lead *models.Lead) {
	var adminUserIDs []uuid.UUID
	s.db.Table("users").
		Select("users.id").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("users.status = ?", "active").
		Distinct().
		Scan(&adminUserIDs)

	content := lead.Name
	if lead.Company != "" {
		content += " (" + lead.Company + ")"
	}
	if lead.Country != "" {
		content += " - " + lead.Country
	}

	for _, uid := range adminUserIDs {
		uidCopy := uid
		s.db.Create(&models.Notification{
			Type:       models.NotificationNewLead,
			Title:      "新询盘",
			Content:    content,
			UserID:     &uidCopy,
			EntityType: "lead",
			EntityID:   &lead.ID,
		})
	}

	// 发送邮件通知给管理员（如果配置了 SMTP）
	if s.mailService.IsEnabled() {
		adminEmails := s.getAdminEmails()
		if len(adminEmails) > 0 {
			scoreLevel := ""
			switch lead.ScoreLevel {
			case models.LeadScoreHigh:
				scoreLevel = "高价值"
			case models.LeadScoreMedium:
				scoreLevel = "中等"
			case models.LeadScoreNormal:
				scoreLevel = "普通"
			default:
				scoreLevel = "低"
			}
			s.mailService.SendLeadNotification(adminEmails, &LeadNotificationData{
				Name:        lead.Name,
				Email:       lead.Email,
				Phone:       lead.Phone,
				Company:     lead.Company,
				Country:     lead.Country,
				Message:     lead.Message,
				ProjectType: lead.ProjectType,
				Score:       lead.Score,
				ScoreLevel:  scoreLevel,
				Time:        lead.CreatedAt.Format("2006-01-02 15:04:05"),
				AdminURL:    os.Getenv("ADMIN_URL"),
			})
		}
	}
}

// getAdminEmails 获取管理员邮箱列表
func (s *LeadService) getAdminEmails() []string {
	var emails []string
	s.db.Table("users").
		Select("users.email").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("users.status = ? AND users.email != ''", "active").
		Distinct().
		Scan(&emails)
	return emails
}

// LeadRequest 询盘请求
type LeadRequest struct {
	Name            string  `json:"name" binding:"required"`
	Company         string  `json:"company"`
	Email           string  `json:"email" binding:"required,email"`
	Phone           string  `json:"phone"`
	WhatsApp        string  `json:"whatsapp"`
	Country         string  `json:"country"`
	CompanyWebsite  string  `json:"company_website"`
	ProjectType     string  `json:"project_type"`
	ProductCategory string  `json:"product_category"`
	ProductID       *string `json:"product_id"`
	Quantity        int     `json:"quantity"`
	Budget          string  `json:"budget"`
	TargetDate      string  `json:"target_date"`
	Message         string  `json:"message"`
	Attachments     string  `json:"attachments"`
	Source          string  `json:"source"`
	Medium          string  `json:"medium"`
	Campaign        string  `json:"campaign"`
	Keyword         string  `json:"keyword"`
	FirstPage       string  `json:"first_page"`
	LandingPage     string  `json:"landing_page"`
	Device          string  `json:"device"`
	Language        string  `json:"language"`
}

// calculateScore 计算询盘评分
func (s *LeadService) calculateScore(req *LeadRequest) int {
	score := 0

	// 有公司网站 +10
	if req.CompanyWebsite != "" {
		score += 10
	}

	// 企业邮箱 +10
	if containsAny(req.Email, []string{"@company.com", "@gmail.com", "@yahoo.com", "@hotmail.com"}) {
		score += 10
	}

	// 数量 > 1000 +20
	if req.Quantity > 1000 {
		score += 20
	} else if req.Quantity > 100 {
		score += 10
	}

	// 有预算 +10
	if req.Budget != "" {
		score += 10
	}

	// 明确交付时间 +10
	if req.TargetDate != "" {
		score += 10
	}

	// 已经选择产品 +10
	if req.ProductID != nil || req.ProductCategory != "" {
		score += 10
	}

	// 有 WhatsApp +5
	if req.WhatsApp != "" {
		score += 5
	}

	// 有详细留言 +5
	if len(req.Message) > 50 {
		score += 5
	}

	return score
}

// getScoreLevel 获取评分等级
func (s *LeadService) getScoreLevel(score int) models.LeadScore {
	switch {
	case score >= 80:
		return models.LeadScoreHigh
	case score >= 50:
		return models.LeadScoreMedium
	case score >= 20:
		return models.LeadScoreNormal
	default:
		return models.LeadScoreLow
	}
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if len(s) > 0 && len(sub) > 0 && len(s) >= len(sub) && s[len(s)-len(sub):] == sub {
			return true
		}
	}
	return false
}

// BackfillLeadID 回写 visit_logs 的 lead_id（关联该访客的所有访问记录）
// 按月分表逐一更新，避免跨表 UPDATE 的复杂度
func (s *LeadService) BackfillLeadID(visitorID string, leadID uuid.UUID) {
	if visitorID == "" || visitorID == "unknown" || visitorID == "anonymous" {
		return
	}
	// 查询最近 90 天的月度表，回写 leadID
	s.db.Exec(`
		UPDATE visit_logs SET lead_id = ? 
		WHERE visitor_id = ? AND lead_id IS NULL AND created_at >= ?
	`, leadID, visitorID, time.Now().AddDate(0, 0, -90))
	// 同时更新月度分表（visit_logs_YYYYMM）
	// 注意：这里简化处理，实际按月表名循环更新
	for i := 0; i < 3; i++ {
		tableName := fmt.Sprintf("visit_logs_%s", time.Now().AddDate(0, -i, 0).Format("200601"))
		s.db.Exec(fmt.Sprintf(`
			UPDATE %s SET lead_id = ? 
			WHERE visitor_id = ? AND lead_id IS NULL
		`, tableName), leadID, visitorID)
	}
}

// ListLeads 询盘列表
func (s *LeadService) ListLeads(page, pageSize int, status, keyword string) ([]models.Lead, int64, error) {
	var leads []models.Lead
	var total int64

	query := s.db.Model(&models.Lead{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR company LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Preload("FollowUps").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&leads).Error

	return leads, total, err
}

// GetLead 获取询盘
func (s *LeadService) GetLead(id string) (*models.Lead, error) {
	var lead models.Lead
	err := s.db.Preload("FollowUps").
		First(&lead, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("询盘不存在")
	}
	return &lead, nil
}

// UpdateLead 更新询盘
func (s *LeadService) UpdateLead(id string, req *UpdateLeadRequest) (*models.Lead, error) {
	var lead models.Lead
	if err := s.db.First(&lead, "id = ?", id).Error; err != nil {
		return nil, errors.New("询盘不存在")
	}

	updates := map[string]interface{}{}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.AssignedTo != nil {
		updates["assigned_to"] = req.AssignedTo
	}
	if req.NextFollowUp != nil {
		updates["next_follow_up"] = req.NextFollowUp
	}
	if req.Notes != "" {
		updates["message"] = req.Notes
	}

	if len(updates) > 0 {
		if err := s.db.Model(&lead).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.GetLead(id)
}

// UpdateLeadRequest 更新询盘请求
type UpdateLeadRequest struct {
	Status       models.LeadStatus `json:"status"`
	AssignedTo   *string           `json:"assigned_to"`
	NextFollowUp *time.Time        `json:"next_follow_up"`
	Notes        string            `json:"notes"`
}

// AddFollowUp 添加跟进记录
func (s *LeadService) AddFollowUp(leadID string, userID string, req *FollowUpRequest) (*models.LeadFollowUp, error) {
	var lead models.Lead
	if err := s.db.First(&lead, "id = ?", leadID).Error; err != nil {
		return nil, errors.New("询盘不存在")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("用户 ID 无效")
	}

	followUp := models.LeadFollowUp{
		LeadID:       lead.ID,
		UserID:       uid,
		Method:       req.Method,
		Content:      req.Content,
		NextFollowUp: req.NextFollowUp,
	}

	if err := s.db.Create(&followUp).Error; err != nil {
		return nil, err
	}

	// 更新询盘状态和下次跟进时间
	updates := map[string]interface{}{
		"status": models.LeadStatusContacted,
	}
	if req.NextFollowUp != nil {
		updates["next_follow_up"] = req.NextFollowUp
	}
	s.db.Model(&lead).Updates(updates)

	return &followUp, nil
}

// FollowUpRequest 跟进请求
type FollowUpRequest struct {
	Method       models.FollowUpMethod `json:"method" binding:"required"`
	Content      string                `json:"content" binding:"required"`
	NextFollowUp *time.Time            `json:"next_follow_up"`
}

// DeleteLead 删除询盘（级联软删除跟进记录）
func (s *LeadService) DeleteLead(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("lead_id = ?", id).Delete(&models.LeadFollowUp{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Lead{}, "id = ?", id).Error
	})
}

// GetDashboardStats 获取仪表盘统计
func (s *LeadService) GetDashboardStats() (map[string]interface{}, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -7)
	monthStart := todayStart.AddDate(0, -1, 0)

	var todayCount, weekCount, monthCount, highValueCount, wonCount int64

	s.db.Model(&models.Lead{}).Where("created_at >= ?", todayStart).Count(&todayCount)
	s.db.Model(&models.Lead{}).Where("created_at >= ?", weekStart).Count(&weekCount)
	s.db.Model(&models.Lead{}).Where("created_at >= ?", monthStart).Count(&monthCount)
	s.db.Model(&models.Lead{}).Where("score >= ?", 80).Count(&highValueCount)
	s.db.Model(&models.Lead{}).Where("status = ?", models.LeadStatusWon).Count(&wonCount)

	// 热门产品
	var hotProducts []map[string]interface{}
	s.db.Model(&models.Lead{}).
		Select("product_category as name, COUNT(*) as count").
		Where("product_category != ''").
		Group("product_category").
		Order("count DESC").
		Limit(5).
		Scan(&hotProducts)

	// 热门国家
	var hotCountries []map[string]interface{}
	s.db.Model(&models.Lead{}).
		Select("country as name, COUNT(*) as count").
		Where("country != ''").
		Group("country").
		Order("count DESC").
		Limit(5).
		Scan(&hotCountries)

	return map[string]interface{}{
		"today_leads":      todayCount,
		"week_leads":       weekCount,
		"month_leads":      monthCount,
		"high_value_leads": highValueCount,
		"won_leads":        wonCount,
		"hot_products":     hotProducts,
		"hot_countries":    hotCountries,
	}, nil
}

// ExportByEmail 根据邮箱导出关联的所有个人数据（GDPR 数据可携权 / PIPL 查询权）
func (s *LeadService) ExportByEmail(email string) (map[string]interface{}, error) {
	var leads []models.Lead
	if err := s.db.Where("email = ?", email).Order("created_at DESC").Find(&leads).Error; err != nil {
		return nil, err
	}
	if len(leads) == 0 {
		return nil, errors.New("未找到对应的个人数据")
	}

	// 构建导出数据结构（包含询盘及其跟进记录）
	exportLeads := make([]map[string]interface{}, 0, len(leads))
	for _, lead := range leads {
		var followUps []models.LeadFollowUp
		s.db.Where("lead_id = ?", lead.ID).Order("created_at ASC").Find(&followUps)

		exportLeads = append(exportLeads, map[string]interface{}{
			"id":          lead.ID,
			"created_at":  lead.CreatedAt,
			"name":        lead.Name,
			"company":     lead.Company,
			"email":       lead.Email,
			"phone":       lead.Phone,
			"whatsapp":    lead.WhatsApp,
			"country":     lead.Country,
			"project_type": lead.ProjectType,
			"message":     lead.Message,
			"source":      lead.Source,
			"medium":      lead.Medium,
			"status":      lead.Status,
			"follow_ups":  followUps,
		})
	}

	return map[string]interface{}{
		"email": email,
		"leads": exportLeads,
	}, nil
}

// AnonymizeByEmail 根据邮箱匿名化或删除个人数据（GDPR 被遗忘权 / PIPL 删除权）
func (s *LeadService) AnonymizeByEmail(email string, reason string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 查找所有关联询盘
		var leads []models.Lead
		if err := tx.Where("email = ?", email).Find(&leads).Error; err != nil {
			return err
		}
		if len(leads) == 0 {
			return errors.New("未找到对应的个人数据")
		}

		for _, lead := range leads {
			// 删除关联的跟进记录
			tx.Where("lead_id = ?", lead.ID).Delete(&models.LeadFollowUp{})

			// 匿名化个人数据（保留业务统计不影响的分析字段）
			updates := map[string]interface{}{
				"name":             "[已删除]",
				"company":          "[已删除]",
				"email":            "anonymized_" + lead.ID.String() + "@deleted.local",
				"phone":            "[已删除]",
				"whatsapp":         "[已删除]",
				"company_website":  "[已删除]",
				"message":          "[已删除]",
				"ip":               "0.0.0.0",
				"status":           models.LeadStatusLost,
				"updated_at":       time.Now(),
			}
			if reason != "" {
				updates["notes"] = "删除原因: " + reason
			}
			tx.Model(&lead).Updates(updates)
		}

		return nil
	})
}
