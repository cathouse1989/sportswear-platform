package services

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// TrashService 回收站服务
// 支持软删除数据的查看、恢复、彻底删除
type TrashService struct {
	db *gorm.DB
}

// NewTrashService 创建回收站服务
func NewTrashService(db *gorm.DB) *TrashService {
	return &TrashService{db: db}
}

// 支持软删除的实体类型注册表
var trashEntities = map[string]interface{}{
	"product":       &models.Product{},
	"category":      &models.Category{},
	"series":        &models.Series{},
	"fabric":        &models.Fabric{},
	"page":          &models.Page{},
	"blog":          &models.Blog{},
	"case":          &models.Case{},
	"faq":           &models.FAQ{},
	"factory":       &models.Factory{},
	"certification": &models.Certification{},
	"navigation":    &models.Navigation{},
	"media":         &models.Media{},
	"self_media":    &models.SelfMedia{},
	"lead":          &models.Lead{},
	"user":          &models.User{},
}

// ListTrashItem 已删除数据条目
type ListTrashItem struct {
	ID          string      `json:"id"`
	EntityType  string      `json:"entity_type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	DeletedAt   interface{} `json:"deleted_at"`
	CreatedAt   interface{} `json:"created_at"`
}

// ListTrash 列出回收站中指定类型的已删除数据
// entityType 为空时列出所有类型
func (s *TrashService) ListTrash(entityType string, page, pageSize int) ([]ListTrashItem, int64, error) {
	var items []ListTrashItem

	if entityType != "" {
		model, ok := trashEntities[entityType]
		if !ok {
			return nil, 0, errors.New("不支持的实体类型")
		}
		items = s.listByModel(model, entityType, page, pageSize)
		var total int64
		s.db.Unscoped().Model(model).Where("deleted_at IS NOT NULL").Count(&total)
		return items, total, nil
	}

	// 列出所有类型
	for et, model := range trashEntities {
		items = append(items, s.listByModel(model, et, page, pageSize)...)
	}
	return items, int64(len(items)), nil
}

func (s *TrashService) listByModel(model interface{}, entityType string, page, pageSize int) []ListTrashItem {
	var items []ListTrashItem

	query := s.db.Unscoped().Model(model).Where("deleted_at IS NOT NULL")

	// 根据实体类型选择标题字段
	titleField := "name"
	switch entityType {
	case "product", "page", "blog", "case":
		titleField = "title"
	case "user":
		titleField = "email"
	case "lead":
		titleField = "email"
	case "media":
		titleField = "original_name"
	case "navigation", "certification", "fabric", "series", "category", "factory", "faq":
		titleField = "name"
	}

	rows, err := query.Select(fmt.Sprintf("id, %s as title, deleted_at, created_at", titleField)).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("deleted_at DESC").Rows()
	if err != nil {
		return items
	}
	defer rows.Close()

	for rows.Next() {
		var item ListTrashItem
		item.EntityType = entityType
		var title string
		if err := rows.Scan(&item.ID, &title, &item.DeletedAt, &item.CreatedAt); err == nil {
			item.Title = title
			items = append(items, item)
		}
	}
	return items
}

// Restore 从回收站恢复数据
func (s *TrashService) Restore(entityType, id string) error {
	model, ok := trashEntities[entityType]
	if !ok {
		return errors.New("不支持的实体类型")
	}

	result := s.db.Unscoped().Model(model).Where("id = ? AND deleted_at IS NOT NULL", id).Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		return errors.New("未找到可恢复的数据")
	}

	// 恢复关联数据
	s.restoreRelations(entityType, id)
	return nil
}

// restoreRelations 恢复主实体的关联数据（翻译、图片等）
func (s *TrashService) restoreRelations(entityType, id string) {
	foreignKeyMap := map[string][]string{
		"product": {"product_translations:product_id", "product_images:product_id", "product_videos:product_id", "product_specs:product_id", "product_customizations:product_id"},
		"page":    {"page_translations:page_id", "page_modules:page_id"},
		"blog":    {"blog_translations:blog_id"},
		"case":    {"case_translations:case_id"},
		"faq":     {"faq_translations:faq_id"},
		"lead":    {"lead_follow_ups:lead_id"},
	}

	tables, ok := foreignKeyMap[entityType]
	if !ok {
		return
	}

	for _, ref := range tables {
		parts := splitRef(ref)
		if len(parts) != 2 {
			continue
		}
		s.db.Unscoped().Table(parts[0]).
			Where(parts[1]+" = ? AND deleted_at IS NOT NULL", id).
			Update("deleted_at", nil)
	}
}

func splitRef(ref string) []string {
	for i := 0; i < len(ref); i++ {
		if ref[i] == ':' {
			return []string{ref[:i], ref[i+1:]}
		}
	}
	return nil
}

// Purge 彻底删除（不可恢复）
func (s *TrashService) Purge(entityType, id string) error {
	model, ok := trashEntities[entityType]
	if !ok {
		return errors.New("不支持的实体类型")
	}

	// 先删除关联数据
	s.purgeRelations(entityType, id)

	result := s.db.Unscoped().Where("id = ?", id).Delete(model)
	if result.RowsAffected == 0 {
		return errors.New("未找到要删除的数据")
	}
	return nil
}

// purgeRelations 彻底删除关联数据
func (s *TrashService) purgeRelations(entityType, id string) {
	foreignKeyMap := map[string][]string{
		"product": {"product_translations:product_id", "product_images:product_id", "product_videos:product_id", "product_specs:product_id", "product_customizations:product_id"},
		"page":    {"page_translations:page_id", "page_modules:page_id"},
		"blog":    {"blog_translations:blog_id"},
		"case":    {"case_translations:case_id"},
		"faq":     {"faq_translations:faq_id"},
		"lead":    {"lead_follow_ups:lead_id"},
	}

	tables, ok := foreignKeyMap[entityType]
	if !ok {
		return
	}

	for _, ref := range tables {
		parts := splitRef(ref)
		if len(parts) != 2 {
			continue
		}
		s.db.Unscoped().Table(parts[0]).Where(parts[1]+" = ?", id).Delete(nil)
	}
}

// EmptyTrash 清空回收站（按类型或全部）
func (s *TrashService) EmptyTrash(entityType string) (int64, error) {
	total := int64(0)

	if entityType != "" {
		model, ok := trashEntities[entityType]
		if !ok {
			return 0, errors.New("不支持的实体类型")
		}
		result := s.db.Unscoped().Where("deleted_at IS NOT NULL").Delete(model)
		return result.RowsAffected, result.Error
	}

	for _, model := range trashEntities {
		result := s.db.Unscoped().Where("deleted_at IS NOT NULL").Delete(model)
		total += result.RowsAffected
	}
	return total, nil
}
