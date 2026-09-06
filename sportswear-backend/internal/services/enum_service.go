package services

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// EnumService 枚举字典服务
// 作为数据库枚举字段的唯一事实源，后台可配置值域 + 多语言翻译 + 排序 + 颜色 + 启停。
type EnumService struct {
	db *gorm.DB
}

// NewEnumService 创建枚举字典服务
func NewEnumService(db *gorm.DB) *EnumService {
	return &EnumService{db: db}
}

// ==================== 视图 / 请求 DTO ====================

// EnumTypeView 枚举类型视图（含条目数）
type EnumTypeView struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	I18nPrefix  string `json:"i18n_prefix"`
	Description string `json:"description"`
	IsSystem    bool   `json:"is_system"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
	ItemCount   int64  `json:"item_count"`
}

// EnumItemView 枚举项视图（translations 已解析为 map，便于前端直接使用）
type EnumItemView struct {
	ID           string            `json:"id"`
	TypeID       string            `json:"type_id"`
	Value        string            `json:"value"`
	Label        string            `json:"label"`
	Translations map[string]string `json:"translations"`
	Color        string            `json:"color"`
	IsDefault    bool              `json:"is_default"`
	IsActive     bool              `json:"is_active"`
	SortOrder    int               `json:"sort_order"`
}

// EnumTypeRequest 枚举类型请求
type EnumTypeRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Module      string `json:"module"`
	I18nPrefix  string `json:"i18n_prefix"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
}

// EnumItemRequest 枚举项请求
type EnumItemRequest struct {
	TypeID       string            `json:"type_id" binding:"required"`
	Value        string            `json:"value" binding:"required"`
	Label        string            `json:"label" binding:"required"`
	Translations map[string]string `json:"translations"`
	Color        string            `json:"color"`
	IsDefault    bool              `json:"is_default"`
	IsActive     *bool             `json:"is_active"`
	SortOrder    int               `json:"sort_order"`
}

// ==================== 类型管理 ====================

// ListTypes 枚举类型列表（按 module/sort_order 排序）
func (s *EnumService) ListTypes() ([]EnumTypeView, error) {
	var types []models.SysEnumType
	if err := s.db.Order("module ASC, sort_order ASC, created_at ASC").Find(&types).Error; err != nil {
		return nil, err
	}
	views := make([]EnumTypeView, 0, len(types))
	for _, t := range types {
		var cnt int64
		s.db.Model(&models.SysEnumItem{}).Where("type_id = ?", t.ID).Count(&cnt)
		views = append(views, EnumTypeView{
			ID:          t.ID.String(),
			Code:        t.Code,
			Name:        t.Name,
			Module:      t.Module,
			I18nPrefix:  t.I18nPrefix,
			Description: t.Description,
			IsSystem:    t.IsSystem,
			IsActive:    t.IsActive,
			SortOrder:   t.SortOrder,
			ItemCount:   cnt,
		})
	}
	return views, nil
}

// UpsertType 创建或更新枚举类型（按 code 唯一）
func (s *EnumService) UpsertType(req *EnumTypeRequest) (*models.SysEnumType, error) {
	var t models.SysEnumType
	err := s.db.Where("code = ?", req.Code).First(&t).Error
	if err == nil {
		t.Name = req.Name
		t.Module = req.Module
		t.I18nPrefix = req.I18nPrefix
		t.Description = req.Description
		if req.IsActive != nil {
			t.IsActive = *req.IsActive
		}
		t.SortOrder = req.SortOrder
		if err := s.db.Save(&t).Error; err != nil {
			return nil, err
		}
		return &t, nil
	}

	t = models.SysEnumType{
		Code:        req.Code,
		Name:        req.Name,
		Module:      req.Module,
		I18nPrefix:  req.I18nPrefix,
		Description: req.Description,
		IsSystem:    false,
		IsActive:    boolDefault(req.IsActive, true),
		SortOrder:   req.SortOrder,
	}
	if err := s.db.Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// DeleteType 删除枚举类型（系统内置类型拒绝删除，级联删除其条目）
func (s *EnumService) DeleteType(id string) error {
	var t models.SysEnumType
	if err := s.db.First(&t, "id = ?", id).Error; err != nil {
		return errors.New("枚举类型不存在")
	}
	if t.IsSystem {
		return errors.New("系统内置枚举类型不可删除")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("type_id = ?", t.ID).Delete(&models.SysEnumItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&t).Error
	})
}

// ==================== 条目管理 ====================

// ListItems 指定类型的枚举项列表
func (s *EnumService) ListItems(typeID string) ([]EnumItemView, error) {
	var items []models.SysEnumItem
	if err := s.db.Where("type_id = ?", typeID).Order("sort_order ASC, created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	views := make([]EnumItemView, 0, len(items))
	for _, it := range items {
		views = append(views, toItemView(it))
	}
	return views, nil
}

// UpsertItem 创建或更新枚举项（按 type_id + value 唯一）
func (s *EnumService) UpsertItem(req *EnumItemRequest) (*models.SysEnumItem, error) {
	uid, err := uuid.Parse(req.TypeID)
	if err != nil {
		return nil, errors.New("枚举类型 ID 非法")
	}
	var t models.SysEnumType
	if err := s.db.First(&t, "id = ?", uid).Error; err != nil {
		return nil, errors.New("枚举类型不存在")
	}

	tr, err := marshalTranslations(req.Translations)
	if err != nil {
		return nil, errors.New("翻译字段格式非法")
	}

	var item models.SysEnumItem
	err = s.db.Where("type_id = ? AND value = ?", uid, req.Value).First(&item).Error
	if err == nil {
		item.Label = req.Label
		item.Translations = tr
		item.Color = req.Color
		item.IsDefault = req.IsDefault
		if req.IsActive != nil {
			item.IsActive = *req.IsActive
		}
		item.SortOrder = req.SortOrder
		if err := s.db.Save(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}

	item = models.SysEnumItem{
		TypeID:       uid,
		Value:        req.Value,
		Label:        req.Label,
		Translations: tr,
		Color:        req.Color,
		IsDefault:    req.IsDefault,
		IsActive:     boolDefault(req.IsActive, true),
		SortOrder:    req.SortOrder,
	}
	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// DeleteItem 删除枚举项（系统内置类型下的条目拒绝删除，仅可停用）
func (s *EnumService) DeleteItem(id string) error {
	var item models.SysEnumItem
	if err := s.db.First(&item, "id = ?", id).Error; err != nil {
		return errors.New("枚举项不存在")
	}
	var t models.SysEnumType
	if err := s.db.First(&t, "id = ?", item.TypeID).Error; err == nil && t.IsSystem {
		return errors.New("系统内置枚举项不可删除，可停用")
	}
	return s.db.Delete(&item).Error
}

// ==================== 公开字典 ====================

// GetPublicEnums 获取指定语言的枚举字典（平铺 key，如 product.type_options.oem → OEM）
// key 规则：i18n_prefix（空则用 code）+ "." + value；label 回退顺序：translations[lang] → label(en) → value。
func (s *EnumService) GetPublicEnums(lang string) (map[string]string, error) {
	if lang == "" {
		lang = "en"
	}
	var types []models.SysEnumType
	if err := s.db.Where("is_active = ?", true).Find(&types).Error; err != nil {
		return nil, err
	}

	typeIDs := make([]uuid.UUID, 0, len(types))
	prefixByType := make(map[uuid.UUID]string, len(types))
	for _, t := range types {
		typeIDs = append(typeIDs, t.ID)
		p := t.I18nPrefix
		if p == "" {
			p = t.Code
		}
		prefixByType[t.ID] = p
	}

	var items []models.SysEnumItem
	if len(typeIDs) > 0 {
		if err := s.db.Where("type_id IN ? AND is_active = ?", typeIDs, true).Find(&items).Error; err != nil {
			return nil, err
		}
	}

	result := make(map[string]string, len(items))
	for _, it := range items {
		prefix := prefixByType[it.TypeID]
		result[prefix+"."+it.Value] = resolveItemLabel(it, lang)
	}
	return result, nil
}

// ==================== 初始化默认数据 ====================

// InitDefaults 初始化内置枚举字典（幂等：类型/条目不存在时才创建，不覆盖运营已改配置）
func (s *EnumService) InitDefaults() error {
	for _, seed := range defaultEnumSeeds {
		var t models.SysEnumType
		err := s.db.Where("code = ?", seed.Code).First(&t).Error
		if err != nil {
			t = models.SysEnumType{
				Code:        seed.Code,
				Name:        seed.Name,
				Module:      seed.Module,
				I18nPrefix:  seed.I18nPrefix,
				Description: seed.Description,
				IsSystem:    true,
				IsActive:    true,
				SortOrder:   seed.SortOrder,
			}
			if err := s.db.Create(&t).Error; err != nil {
				return err
			}
		}
		for _, it := range seed.Items {
			var count int64
			s.db.Model(&models.SysEnumItem{}).Where("type_id = ? AND value = ?", t.ID, it.Value).Count(&count)
			if count > 0 {
				continue
			}
			tr, _ := marshalTranslations(it.Translations)
			item := models.SysEnumItem{
				TypeID:       t.ID,
				Value:        it.Value,
				Label:        it.Label,
				Translations: tr,
				Color:        it.Color,
				IsDefault:    it.IsDefault,
				IsActive:     true,
				SortOrder:    it.SortOrder,
			}
			if err := s.db.Create(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// ==================== 辅助函数 ====================

func toItemView(it models.SysEnumItem) EnumItemView {
	return EnumItemView{
		ID:           it.ID.String(),
		TypeID:       it.TypeID.String(),
		Value:        it.Value,
		Label:        it.Label,
		Translations: parseTranslations(it.Translations),
		Color:        it.Color,
		IsDefault:    it.IsDefault,
		IsActive:     it.IsActive,
		SortOrder:    it.SortOrder,
	}
}

func parseTranslations(raw string) map[string]string {
	m := map[string]string{}
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &m)
	}
	return m
}

func marshalTranslations(m map[string]string) (string, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}", err
	}
	return string(b), nil
}

func resolveItemLabel(item models.SysEnumItem, lang string) string {
	tr := parseTranslations(item.Translations)
	if v, ok := tr[lang]; ok && v != "" {
		return v
	}
	if item.Label != "" {
		return item.Label
	}
	return item.Value
}

func boolDefault(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

// ==================== 内置枚举种子 ====================

type enumItemSeed struct {
	Value        string
	Label        string
	Translations map[string]string
	Color        string
	IsDefault    bool
	SortOrder    int
}

type enumTypeSeed struct {
	Code        string
	Name        string
	Module      string
	I18nPrefix  string
	Description string
	SortOrder   int
	Items       []enumItemSeed
}

// 种子数据翻译对齐门户静态语言包 locales/{en,zh,es,fr}.json，保证「字典优先、静态兜底」一致。
var defaultEnumSeeds = []enumTypeSeed{
	{
		Code: "product.type", Name: "产品类型", Module: "product", I18nPrefix: "product.type_options",
		Description: "产品 OEM/ODM 类型", SortOrder: 10,
		Items: []enumItemSeed{
			{Value: "oem", Label: "OEM", Translations: map[string]string{"zh": "OEM", "es": "OEM", "fr": "OEM"}, Color: "blue", SortOrder: 1},
			{Value: "odm", Label: "ODM", Translations: map[string]string{"zh": "ODM", "es": "ODM", "fr": "ODM"}, Color: "purple", SortOrder: 2},
			{Value: "both", Label: "OEM/ODM", Translations: map[string]string{"zh": "OEM/ODM", "es": "OEM/ODM", "fr": "OEM/ODM"}, Color: "gold", IsDefault: true, SortOrder: 3},
		},
	},
	{
		Code: "product.gender", Name: "适用性别", Module: "product", I18nPrefix: "product.gender_options",
		Description: "产品适用性别", SortOrder: 20,
		Items: []enumItemSeed{
			{Value: "unisex", Label: "Unisex", Translations: map[string]string{"zh": "男女皆宜", "es": "Unisex", "fr": "Unisexe"}, Color: "green", IsDefault: true, SortOrder: 1},
			{Value: "male", Label: "Male", Translations: map[string]string{"zh": "男款", "es": "Hombre", "fr": "Homme"}, Color: "blue", SortOrder: 2},
			{Value: "female", Label: "Female", Translations: map[string]string{"zh": "女款", "es": "Mujer", "fr": "Femme"}, Color: "magenta", SortOrder: 3},
			{Value: "kids", Label: "Kids", Translations: map[string]string{"zh": "儿童", "es": "Niños", "fr": "Enfants"}, Color: "orange", SortOrder: 4},
		},
	},
	{
		Code: "product.status", Name: "产品状态", Module: "product", I18nPrefix: "product.status_options",
		Description: "产品发布状态", SortOrder: 30,
		Items: []enumItemSeed{
			{Value: "draft", Label: "Draft", Translations: map[string]string{"zh": "草稿", "es": "Borrador", "fr": "Brouillon"}, Color: "info", SortOrder: 1},
			{Value: "published", Label: "Published", Translations: map[string]string{"zh": "已发布", "es": "Publicado", "fr": "Publié"}, Color: "success", SortOrder: 2},
			{Value: "offline", Label: "Offline", Translations: map[string]string{"zh": "已下线", "es": "Fuera de línea", "fr": "Hors ligne"}, Color: "danger", SortOrder: 3},
		},
	},
	{
		Code: "product.customization.type", Name: "定制类型", Module: "product", I18nPrefix: "product.customization_options",
		Description: "产品定制能力类型", SortOrder: 40,
		Items: []enumItemSeed{
			{Value: "logo", Label: "Logo", Translations: map[string]string{"zh": "Logo", "es": "Logo", "fr": "Logo"}, SortOrder: 1},
			{Value: "color", Label: "Color", Translations: map[string]string{"zh": "颜色", "es": "Color", "fr": "Couleur"}, SortOrder: 2},
			{Value: "fabric", Label: "Fabric", Translations: map[string]string{"zh": "面料", "es": "Tela", "fr": "Tissu"}, SortOrder: 3},
			{Value: "pattern", Label: "Pattern", Translations: map[string]string{"zh": "图案", "es": "Patrón", "fr": "Motif"}, SortOrder: 4},
			{Value: "print", Label: "Print", Translations: map[string]string{"zh": "印花", "es": "Estampado", "fr": "Impression"}, SortOrder: 5},
			{Value: "embroidery", Label: "Embroidery", Translations: map[string]string{"zh": "刺绣", "es": "Bordado", "fr": "Broderie"}, SortOrder: 6},
			{Value: "label", Label: "Label", Translations: map[string]string{"zh": "标签", "es": "Etiqueta", "fr": "Étiquette"}, SortOrder: 7},
			{Value: "hangtag", Label: "Hangtag", Translations: map[string]string{"zh": "吊牌", "es": "Etiqueta colgante", "fr": "Étiquette volante"}, SortOrder: 8},
			{Value: "packaging", Label: "Packaging", Translations: map[string]string{"zh": "包装", "es": "Embalaje", "fr": "Emballage"}, SortOrder: 9},
			{Value: "size", Label: "Size", Translations: map[string]string{"zh": "尺码", "es": "Talla", "fr": "Taille"}, SortOrder: 10},
			{Value: "fit", Label: "Fit", Translations: map[string]string{"zh": "版型", "es": "Ajuste", "fr": "Ajustement"}, SortOrder: 11},
			{Value: "zipper", Label: "Zipper", Translations: map[string]string{"zh": "拉链", "es": "Cremallera", "fr": "Fermeture éclair"}, SortOrder: 12},
			{Value: "button", Label: "Button", Translations: map[string]string{"zh": "纽扣", "es": "Botón", "fr": "Bouton"}, SortOrder: 13},
			{Value: "belt", Label: "Belt", Translations: map[string]string{"zh": "腰带", "es": "Cinturón", "fr": "Ceinture"}, SortOrder: 14},
			{Value: "accessory", Label: "Accessory", Translations: map[string]string{"zh": "配饰", "es": "Accesorio", "fr": "Accessoire"}, SortOrder: 15},
		},
	},
	{
		Code: "product.spec.name", Name: "规格名称", Module: "product", I18nPrefix: "product.spec_options",
		Description: "产品规格项名称", SortOrder: 50,
		Items: []enumItemSeed{
			{Value: "fabric", Label: "Fabric", Translations: map[string]string{"zh": "面料", "es": "Tela", "fr": "Tissu"}, SortOrder: 1},
			{Value: "moq", Label: "MOQ", Translations: map[string]string{"zh": "起订量", "es": "MOQ", "fr": "MOQ"}, SortOrder: 2},
			{Value: "fit", Label: "Fit", Translations: map[string]string{"zh": "合身度", "es": "Ajuste", "fr": "Ajustement"}, SortOrder: 3},
			{Value: "composition", Label: "Composition", Translations: map[string]string{"zh": "成分", "es": "Composición", "fr": "Composition"}, SortOrder: 4},
			{Value: "size range", Label: "Size Range", Translations: map[string]string{"zh": "尺码范围", "es": "Rango de Tallas", "fr": "Gamme de Tailles"}, SortOrder: 5},
			{Value: "weight", Label: "Weight", Translations: map[string]string{"zh": "重量", "es": "Peso", "fr": "Poids"}, SortOrder: 6},
		},
	},
	{
		Code: "product.video.type", Name: "视频类型", Module: "product", I18nPrefix: "product.video_type_options",
		Description: "产品视频类型", SortOrder: 60,
		Items: []enumItemSeed{
			{Value: "product", Label: "Product Video", Translations: map[string]string{"zh": "产品视频", "es": "Vídeo de producto", "fr": "Vidéo produit"}, SortOrder: 1},
			{Value: "production", Label: "Production Video", Translations: map[string]string{"zh": "生产视频", "es": "Vídeo de producción", "fr": "Vidéo de production"}, SortOrder: 2},
			{Value: "craft", Label: "Craft Video", Translations: map[string]string{"zh": "工艺视频", "es": "Vídeo de proceso", "fr": "Vidéo d'artisanat"}, SortOrder: 3},
			{Value: "usage", Label: "Usage Video", Translations: map[string]string{"zh": "用途视频", "es": "Vídeo de uso", "fr": "Vidéo d'utilisation"}, SortOrder: 4},
		},
	},
	{
		Code: "content.status", Name: "内容状态", Module: "content", I18nPrefix: "content.status_options",
		Description: "博客/案例/页面通用状态", SortOrder: 70,
		Items: []enumItemSeed{
			{Value: "draft", Label: "Draft", Translations: map[string]string{"zh": "草稿", "es": "Borrador", "fr": "Brouillon"}, Color: "info", SortOrder: 1},
			{Value: "review", Label: "Review", Translations: map[string]string{"zh": "审核中", "es": "En revisión", "fr": "En révision"}, Color: "warning", SortOrder: 2},
			{Value: "published", Label: "Published", Translations: map[string]string{"zh": "已发布", "es": "Publicado", "fr": "Publié"}, Color: "success", SortOrder: 3},
			{Value: "offline", Label: "Offline", Translations: map[string]string{"zh": "已下线", "es": "Fuera de línea", "fr": "Hors ligne"}, Color: "danger", SortOrder: 4},
			{Value: "archived", Label: "Archived", Translations: map[string]string{"zh": "已归档", "es": "Archivado", "fr": "Archivé"}, Color: "info", SortOrder: 5},
		},
	},
	{
		Code: "blog.category", Name: "博客分类", Module: "content", I18nPrefix: "blog.category_options",
		Description: "博客文章分类", SortOrder: 80,
		Items: []enumItemSeed{
			{Value: "oem_guide", Label: "OEM Guide", Translations: map[string]string{"zh": "OEM 指南", "es": "Guía OEM", "fr": "Guide OEM"}, SortOrder: 1},
			{Value: "odm_guide", Label: "ODM Guide", Translations: map[string]string{"zh": "ODM 指南", "es": "Guía ODM", "fr": "Guide ODM"}, SortOrder: 2},
			{Value: "fabric", Label: "Fabric", Translations: map[string]string{"zh": "面料", "es": "Tela", "fr": "Tissu"}, SortOrder: 3},
			{Value: "trend", Label: "Trend", Translations: map[string]string{"zh": "趋势", "es": "Tendencia", "fr": "Tendance"}, SortOrder: 4},
			{Value: "industry", Label: "Industry", Translations: map[string]string{"zh": "行业", "es": "Industria", "fr": "Industrie"}, SortOrder: 5},
		},
	},
	{
		Code: "case.project_type", Name: "项目类型", Module: "content", I18nPrefix: "case.project_type_options",
		Description: "案例项目类型", SortOrder: 90,
		Items: []enumItemSeed{
			{Value: "oem", Label: "OEM", Translations: map[string]string{"zh": "OEM", "es": "OEM", "fr": "OEM"}, SortOrder: 1},
			{Value: "odm", Label: "ODM", Translations: map[string]string{"zh": "ODM", "es": "ODM", "fr": "ODM"}, SortOrder: 2},
			{Value: "private_label", Label: "Private Label", Translations: map[string]string{"zh": "自有品牌", "es": "Marca propia", "fr": "Marque privée"}, SortOrder: 3},
		},
	},
	{
		Code: "faq.category", Name: "FAQ 分类", Module: "content", I18nPrefix: "faq.category_options",
		Description: "常见问题分类", SortOrder: 100,
		Items: []enumItemSeed{
			{Value: "moq", Label: "MOQ", Translations: map[string]string{"zh": "MOQ", "es": "MOQ", "fr": "MOQ"}, SortOrder: 1},
			{Value: "oem", Label: "OEM", Translations: map[string]string{"zh": "OEM", "es": "OEM", "fr": "OEM"}, SortOrder: 2},
			{Value: "odm", Label: "ODM", Translations: map[string]string{"zh": "ODM", "es": "ODM", "fr": "ODM"}, SortOrder: 3},
			{Value: "sample", Label: "Sampling", Translations: map[string]string{"zh": "打样", "es": "Muestreo", "fr": "Échantillonnage"}, SortOrder: 4},
			{Value: "payment", Label: "Payment", Translations: map[string]string{"zh": "支付", "es": "Pago", "fr": "Paiement"}, SortOrder: 5},
			{Value: "production", Label: "Production", Translations: map[string]string{"zh": "生产", "es": "Producción", "fr": "Production"}, SortOrder: 6},
			{Value: "logistics", Label: "Logistics", Translations: map[string]string{"zh": "物流", "es": "Logística", "fr": "Logistique"}, SortOrder: 7},
			{Value: "fabric", Label: "Fabric", Translations: map[string]string{"zh": "面料", "es": "Tela", "fr": "Tissu"}, SortOrder: 8},
			{Value: "quality", Label: "Quality", Translations: map[string]string{"zh": "质量", "es": "Calidad", "fr": "Qualité"}, SortOrder: 9},
			{Value: "certification", Label: "Certification", Translations: map[string]string{"zh": "认证", "es": "Certificación", "fr": "Certification"}, SortOrder: 10},
		},
	},
	{
		Code: "media.type", Name: "媒体类型", Module: "media", I18nPrefix: "media.type_options",
		Description: "媒体文件类型", SortOrder: 110,
		Items: []enumItemSeed{
			{Value: "image", Label: "Image", Translations: map[string]string{"zh": "图片", "es": "Imagen", "fr": "Image"}, SortOrder: 1},
			{Value: "video", Label: "Video", Translations: map[string]string{"zh": "视频", "es": "Vídeo", "fr": "Vidéo"}, SortOrder: 2},
			{Value: "file", Label: "File", Translations: map[string]string{"zh": "文件", "es": "Archivo", "fr": "Fichier"}, SortOrder: 3},
		},
	},
	{
		Code: "media.category", Name: "媒体分类", Module: "media", I18nPrefix: "media.category_options",
		Description: "媒体文件分类", SortOrder: 120,
		Items: []enumItemSeed{
			{Value: "product", Label: "Product", Translations: map[string]string{"zh": "产品", "es": "Producto", "fr": "Produit"}, SortOrder: 1},
			{Value: "home", Label: "Home", Translations: map[string]string{"zh": "首页", "es": "Inicio", "fr": "Accueil"}, SortOrder: 2},
			{Value: "oem", Label: "OEM", Translations: map[string]string{"zh": "OEM", "es": "OEM", "fr": "OEM"}, SortOrder: 3},
			{Value: "odm", Label: "ODM", Translations: map[string]string{"zh": "ODM", "es": "ODM", "fr": "ODM"}, SortOrder: 4},
			{Value: "factory", Label: "Factory", Translations: map[string]string{"zh": "工厂", "es": "Fábrica", "fr": "Usine"}, SortOrder: 5},
			{Value: "production", Label: "Production", Translations: map[string]string{"zh": "生产流程", "es": "Producción", "fr": "Production"}, SortOrder: 6},
			{Value: "case", Label: "Case", Translations: map[string]string{"zh": "案例", "es": "Caso", "fr": "Étude de cas"}, SortOrder: 7},
			{Value: "blog", Label: "Blog", Translations: map[string]string{"zh": "博客", "es": "Blog", "fr": "Blog"}, SortOrder: 8},
			{Value: "certification", Label: "Certification", Translations: map[string]string{"zh": "认证", "es": "Certificación", "fr": "Certification"}, SortOrder: 9},
			{Value: "public", Label: "Public", Translations: map[string]string{"zh": "公开", "es": "Público", "fr": "Public"}, SortOrder: 10},
		},
	},
	{
		Code: "page.module.type", Name: "页面模块类型", Module: "content", I18nPrefix: "page.module_type_options",
		Description: "页面编辑器模块类型", SortOrder: 130,
		Items: []enumItemSeed{
			{Value: "banner", Label: "Banner", Translations: map[string]string{"zh": "轮播 Banner", "es": "Banner", "fr": "Bannière"}, SortOrder: 1},
			{Value: "text", Label: "Rich Text", Translations: map[string]string{"zh": "富文本", "es": "Texto enriquecido", "fr": "Texte enrichi"}, SortOrder: 2},
			{Value: "image", Label: "Image", Translations: map[string]string{"zh": "图片", "es": "Imagen", "fr": "Image"}, SortOrder: 3},
			{Value: "video", Label: "Video", Translations: map[string]string{"zh": "视频", "es": "Vídeo", "fr": "Vidéo"}, SortOrder: 4},
			{Value: "product_recommend", Label: "Product Recommend", Translations: map[string]string{"zh": "产品推荐", "es": "Recomendación de producto", "fr": "Recommandation de produit"}, SortOrder: 5},
			{Value: "oem", Label: "OEM Service", Translations: map[string]string{"zh": "OEM 服务", "es": "Servicio OEM", "fr": "Service OEM"}, SortOrder: 6},
			{Value: "odm", Label: "ODM Service", Translations: map[string]string{"zh": "ODM 服务", "es": "Servicio ODM", "fr": "Service ODM"}, SortOrder: 7},
			{Value: "factory", Label: "Factory", Translations: map[string]string{"zh": "工厂展示", "es": "Fábrica", "fr": "Usine"}, SortOrder: 8},
			{Value: "production", Label: "Production", Translations: map[string]string{"zh": "生产流程", "es": "Producción", "fr": "Production"}, SortOrder: 9},
			{Value: "case", Label: "Case", Translations: map[string]string{"zh": "案例", "es": "Caso", "fr": "Étude de cas"}, SortOrder: 10},
			{Value: "certification", Label: "Certification", Translations: map[string]string{"zh": "认证", "es": "Certificación", "fr": "Certification"}, SortOrder: 11},
			{Value: "blog", Label: "Blog", Translations: map[string]string{"zh": "博客", "es": "Blog", "fr": "Blog"}, SortOrder: 12},
			{Value: "contact", Label: "Contact Form", Translations: map[string]string{"zh": "联系表单", "es": "Formulario de contacto", "fr": "Formulaire de contact"}, SortOrder: 13},
		},
	},
	{
		Code: "lead.status", Name: "询盘状态", Module: "crm", I18nPrefix: "lead.status_options",
		Description: "询盘销售跟进状态", SortOrder: 140,
		Items: []enumItemSeed{
			{Value: "new", Label: "New", Translations: map[string]string{"zh": "新询盘", "es": "Nueva consulta", "fr": "Nouvelle demande"}, Color: "primary", SortOrder: 1},
			{Value: "contacted", Label: "Contacted", Translations: map[string]string{"zh": "已联系", "es": "Contactado", "fr": "Contacté"}, Color: "info", SortOrder: 2},
			{Value: "confirmed", Label: "Confirmed", Translations: map[string]string{"zh": "已确认", "es": "Confirmado", "fr": "Confirmé"}, Color: "warning", SortOrder: 3},
			{Value: "quoted", Label: "Quoted", Translations: map[string]string{"zh": "已报价", "es": "Cotizado", "fr": "Devis envoyé"}, Color: "warning", SortOrder: 4},
			{Value: "sampling", Label: "Sampling", Translations: map[string]string{"zh": "打样中", "es": "Muestreo", "fr": "Échantillonnage"}, Color: "warning", SortOrder: 5},
			{Value: "negotiating", Label: "Negotiating", Translations: map[string]string{"zh": "洽谈中", "es": "Negociando", "fr": "En négociation"}, Color: "warning", SortOrder: 6},
			{Value: "won", Label: "Won", Translations: map[string]string{"zh": "已成交", "es": "Ganado", "fr": "Gagné"}, Color: "success", SortOrder: 7},
			{Value: "lost", Label: "Lost", Translations: map[string]string{"zh": "已流失", "es": "Perdido", "fr": "Perdu"}, Color: "danger", SortOrder: 8},
			{Value: "spam", Label: "Spam", Translations: map[string]string{"zh": "垃圾", "es": "Spam", "fr": "Spam"}, Color: "info", SortOrder: 9},
		},
	},
	{
		Code: "subscriber.status", Name: "订阅状态", Module: "crm", I18nPrefix: "subscriber.status_options",
		Description: "Newsletter 订阅状态", SortOrder: 150,
		Items: []enumItemSeed{
			{Value: "subscribed", Label: "Subscribed", Translations: map[string]string{"zh": "已订阅", "es": "Suscrito", "fr": "Abonné"}, Color: "success", SortOrder: 1},
			{Value: "unsubscribed", Label: "Unsubscribed", Translations: map[string]string{"zh": "已退订", "es": "Dado de baja", "fr": "Désabonné"}, Color: "info", SortOrder: 2},
		},
	},
	{
		Code: "user.status", Name: "用户状态", Module: "system", I18nPrefix: "user.status_options",
		Description: "后台账号启停状态", SortOrder: 160,
		Items: []enumItemSeed{
			{Value: "active", Label: "Active", Translations: map[string]string{"zh": "启用", "es": "Activo", "fr": "Actif"}, Color: "success", SortOrder: 1},
			{Value: "disabled", Label: "Disabled", Translations: map[string]string{"zh": "禁用", "es": "Deshabilitado", "fr": "Désactivé"}, Color: "danger", SortOrder: 2},
		},
	},
	{
		Code: "self_media.platform", Name: "自媒体平台", Module: "content", I18nPrefix: "self_media.platform_options",
		Description: "门户「关注我们」支持的第三方平台", SortOrder: 170,
		Items: []enumItemSeed{
			{Value: "wechat", Label: "WeChat", Translations: map[string]string{"zh": "微信公众号", "es": "WeChat", "fr": "WeChat"}, SortOrder: 1},
			{Value: "weibo", Label: "Weibo", Translations: map[string]string{"zh": "微博", "es": "Weibo", "fr": "Weibo"}, SortOrder: 2},
			{Value: "douyin", Label: "Douyin", Translations: map[string]string{"zh": "抖音", "es": "Douyin", "fr": "Douyin"}, SortOrder: 3},
			{Value: "xiaohongshu", Label: "Xiaohongshu", Translations: map[string]string{"zh": "小红书", "es": "Xiaohongshu", "fr": "Xiaohongshu"}, SortOrder: 4},
			{Value: "video", Label: "Video Channel", Translations: map[string]string{"zh": "视频号", "es": "Canal de vídeo", "fr": "Chaîne vidéo"}, SortOrder: 5},
			{Value: "bilibili", Label: "Bilibili", Translations: map[string]string{"zh": "B站", "es": "Bilibili", "fr": "Bilibili"}, SortOrder: 6},
			{Value: "youtube", Label: "YouTube", Translations: map[string]string{"zh": "YouTube", "es": "YouTube", "fr": "YouTube"}, SortOrder: 7},
			{Value: "instagram", Label: "Instagram", Translations: map[string]string{"zh": "Instagram", "es": "Instagram", "fr": "Instagram"}, SortOrder: 8},
			{Value: "facebook", Label: "Facebook", Translations: map[string]string{"zh": "Facebook", "es": "Facebook", "fr": "Facebook"}, SortOrder: 9},
			{Value: "twitter", Label: "Twitter/X", Translations: map[string]string{"zh": "Twitter/X", "es": "Twitter/X", "fr": "Twitter/X"}, SortOrder: 10},
			{Value: "linkedin", Label: "LinkedIn", Translations: map[string]string{"zh": "LinkedIn", "es": "LinkedIn", "fr": "LinkedIn"}, SortOrder: 11},
			{Value: "tiktok", Label: "TikTok", Translations: map[string]string{"zh": "TikTok", "es": "TikTok", "fr": "TikTok"}, SortOrder: 12},
			{Value: "other", Label: "Other", Translations: map[string]string{"zh": "其他", "es": "Otro", "fr": "Autre"}, SortOrder: 13},
		},
	},
	{
		Code: "navigation.type", Name: "导航类型", Module: "content", I18nPrefix: "navigation.type_options",
		Description: "导航位置类型", SortOrder: 180,
		Items: []enumItemSeed{
			{Value: "header", Label: "Header", Translations: map[string]string{"zh": "Header 导航", "es": "Header", "fr": "Header"}, SortOrder: 1},
			{Value: "footer", Label: "Footer", Translations: map[string]string{"zh": "Footer 导航", "es": "Footer", "fr": "Footer"}, SortOrder: 2},
		},
	},
}
