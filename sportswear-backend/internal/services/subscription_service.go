package services

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// SubscriptionService 订阅服务（门户"订阅更新"Newsletter）
type SubscriptionService struct {
	db *gorm.DB
}

// NewSubscriptionService 创建订阅服务
func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

// SubscribeRequest 订阅请求
type SubscribeRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Language    string `json:"language"`
	Source      string `json:"source"`
	Medium      string `json:"medium"`
	Campaign    string `json:"campaign"`
	Keyword     string `json:"keyword"`
	LandingPage string `json:"landing_page"`
	Device      string `json:"device"`
}

// UpdateSubscriberRequest 更新订阅请求
type UpdateSubscriberRequest struct {
	Status models.SubscriberStatus `json:"status" binding:"required"`
}

// Subscribe 订阅（邮箱去重：已存在则重新激活为 subscribed，并刷新来源归因信息）
func (s *SubscriptionService) Subscribe(req *SubscribeRequest, ip, visitorID string) (*models.Subscriber, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		return nil, errors.New("邮箱不能为空")
	}

	now := time.Now()
	var sub models.Subscriber
	err := s.db.Where("email = ?", email).First(&sub).Error
	if err == nil {
		// 已存在：重新激活，并刷新来源归因
		sub.Status = models.SubscriberStatusSubscribed
		sub.SubscribedAt = &now
		sub.UnsubscribedAt = nil
		sub.Language = req.Language
		sub.Source = req.Source
		sub.Medium = req.Medium
		sub.Campaign = req.Campaign
		sub.Keyword = req.Keyword
		sub.LandingPage = req.LandingPage
		sub.Device = req.Device
		sub.IP = ip
		sub.VisitorID = visitorID
		if e := s.db.Save(&sub).Error; e != nil {
			return nil, e
		}
		return &sub, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	sub = models.Subscriber{
		Email:        email,
		Status:       models.SubscriberStatusSubscribed,
		Language:     req.Language,
		Source:       req.Source,
		Medium:       req.Medium,
		Campaign:     req.Campaign,
		Keyword:      req.Keyword,
		LandingPage:  req.LandingPage,
		Device:       req.Device,
		IP:           ip,
		VisitorID:    visitorID,
		SubscribedAt: &now,
	}
	if e := s.db.Create(&sub).Error; e != nil {
		return nil, e
	}
	return &sub, nil
}

// List 订阅列表（分页 + 状态/关键词过滤）
func (s *SubscriptionService) List(page, pageSize int, status, keyword string) ([]models.Subscriber, int64, error) {
	var subs []models.Subscriber
	var total int64

	q := s.db.Model(&models.Subscriber{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		like := "%" + strings.TrimSpace(keyword) + "%"
		q = q.Where("email ILIKE ? OR source ILIKE ? OR landing_page ILIKE ?", like, like, like)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if err := q.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&subs).Error; err != nil {
		return nil, 0, err
	}
	return subs, total, nil
}

// UpdateStatus 更新订阅状态（subscribed / unsubscribed）
func (s *SubscriptionService) UpdateStatus(id string, status models.SubscriberStatus) (*models.Subscriber, error) {
	if status != models.SubscriberStatusSubscribed && status != models.SubscriberStatusUnsubscribed {
		return nil, errors.New("订阅状态无效")
	}
	var sub models.Subscriber
	if err := s.db.First(&sub, "id = ?", id).Error; err != nil {
		return nil, errors.New("订阅记录不存在")
	}
	now := time.Now()
	sub.Status = status
	if status == models.SubscriberStatusSubscribed {
		sub.SubscribedAt = &now
		sub.UnsubscribedAt = nil
	} else {
		sub.UnsubscribedAt = &now
	}
	if err := s.db.Save(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

// Delete 删除订阅（软删除）
func (s *SubscriptionService) Delete(id string) error {
	result := s.db.Delete(&models.Subscriber{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("订阅记录不存在")
	}
	return nil
}

// Stats 订阅统计
func (s *SubscriptionService) Stats() (map[string]interface{}, error) {
	var total, subscribed, unsubscribed int64
	if err := s.db.Model(&models.Subscriber{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Subscriber{}).Where("status = ?", models.SubscriberStatusSubscribed).Count(&subscribed).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.Subscriber{}).Where("status = ?", models.SubscriberStatusUnsubscribed).Count(&unsubscribed).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"total":        total,
		"subscribed":   subscribed,
		"unsubscribed": unsubscribed,
	}, nil
}
