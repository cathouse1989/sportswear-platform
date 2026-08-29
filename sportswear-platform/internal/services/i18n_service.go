package services

import (
	"gorm.io/gorm"

	"sportswear-platform/internal/models"
)

// I18nService 国际化服务
// 负责 UI 词条的配置管理，以及内容实体的按语言翻译解析
type I18nService struct {
	db *gorm.DB
}

// NewI18nService 创建国际化服务
func NewI18nService(db *gorm.DB) *I18nService {
	return &I18nService{db: db}
}

// ==================== 词条管理 ====================

// ListEntries 词条列表（后台管理）
func (s *I18nService) ListEntries(page, pageSize int, language, module, keyword string) ([]models.I18nEntry, int64, error) {
	var entries []models.I18nEntry
	var total int64

	query := s.db.Model(&models.I18nEntry{})
	if language != "" {
		query = query.Where("language = ?", language)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("key LIKE ? OR value LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	err := query.Order("module ASC, sort_order ASC, key ASC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&entries).Error
	return entries, total, err
}

// GetEntriesByLang 获取指定语言的全部词条（公开接口，前端一键切换）
func (s *I18nService) GetEntriesByLang(lang string) (map[string]string, error) {
	if lang == "" {
		lang = "en"
	}
	var entries []models.I18nEntry
	if err := s.db.Where("language = ? AND is_active = ?", lang, true).Find(&entries).Error; err != nil {
		return nil, err
	}

	// 若目标语言无词条，回退到英文
	if len(entries) == 0 && lang != "en" {
		if err := s.db.Where("language = ? AND is_active = ?", "en", true).Find(&entries).Error; err != nil {
			return nil, err
		}
	}

	result := make(map[string]string, len(entries))
	for _, e := range entries {
		result[e.Key] = e.Value
	}
	return result, nil
}

// UpsertEntry 创建或更新词条（按 key + language 唯一）
func (s *I18nService) UpsertEntry(req *I18nEntryRequest) (*models.I18nEntry, error) {
	var entry models.I18nEntry
	err := s.db.Where("key = ? AND language = ?", req.Key, req.Language).First(&entry).Error
	if err == nil {
		// 更新
		entry.Value = req.Value
		entry.Module = req.Module
		if err := s.db.Save(&entry).Error; err != nil {
			return nil, err
		}
		return &entry, nil
	}

	// 创建
	entry = models.I18nEntry{
		Key:      req.Key,
		Language: req.Language,
		Value:    req.Value,
		Module:   req.Module,
		IsActive: true,
	}
	if err := s.db.Create(&entry).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

// DeleteEntry 删除词条
func (s *I18nService) DeleteEntry(id string) error {
	return s.db.Delete(&models.I18nEntry{}, "id = ?", id).Error
}

// InitDefaults 初始化默认词条（英文 + 中文 + 西班牙语 + 法语）
func (s *I18nService) InitDefaults() error {
	var count int64
	s.db.Model(&models.I18nEntry{}).Count(&count)
	if count > 0 {
		return nil
	}
	entries := append(defaultI18nEntries, defaultI18nEntriesZh...)
	entries = append(entries, defaultI18nEntriesEs...)
	entries = append(entries, defaultI18nEntriesFr...)
	return s.db.Create(&entries).Error
}

// I18nEntryRequest 词条请求
type I18nEntryRequest struct {
	Key      string `json:"key" binding:"required"`
	Language string `json:"language" binding:"required"`
	Value    string `json:"value" binding:"required"`
	Module   string `json:"module"`
}

// ==================== 内容按语言解析 ====================

// ResolveLang 解析有效语言：若目标语言无配置则回退英文
func (s *I18nService) ResolveLang(lang string) string {
	if lang == "" {
		return "en"
	}
	var count int64
	s.db.Model(&models.Language{}).Where("code = ? AND is_active = ?", lang, true).Count(&count)
	if count == 0 {
		return "en"
	}
	return lang
}

// 默认英文词条（UI 文案）
var defaultI18nEntries = []models.I18nEntry{
	// common
	{Key: "common.learn_more", Language: "en", Value: "Learn More", Module: "common"},
	{Key: "common.view_all", Language: "en", Value: "View All", Module: "common"},
	{Key: "common.read_more", Language: "en", Value: "Read More", Module: "common"},
	{Key: "common.contact_us", Language: "en", Value: "Contact Us", Module: "common"},
	{Key: "common.get_quote", Language: "en", Value: "Get a Quote", Module: "common"},
	{Key: "common.send_inquiry", Language: "en", Value: "Send Inquiry", Module: "common"},
	{Key: "common.view_details", Language: "en", Value: "View Details", Module: "common"},
	{Key: "common.back_to_top", Language: "en", Value: "Back to Top", Module: "common"},
	// nav
	{Key: "nav.home", Language: "en", Value: "Home", Module: "nav"},
	{Key: "nav.products", Language: "en", Value: "Products", Module: "nav"},
	{Key: "nav.oem", Language: "en", Value: "OEM Service", Module: "nav"},
	{Key: "nav.odm", Language: "en", Value: "ODM Service", Module: "nav"},
	{Key: "nav.factory", Language: "en", Value: "Factory", Module: "nav"},
	{Key: "nav.cases", Language: "en", Value: "Cases", Module: "nav"},
	{Key: "nav.blog", Language: "en", Value: "Blog", Module: "nav"},
	{Key: "nav.faq", Language: "en", Value: "FAQ", Module: "nav"},
	{Key: "nav.contact", Language: "en", Value: "Contact", Module: "nav"},
	// button
	{Key: "btn.submit", Language: "en", Value: "Submit", Module: "button"},
	{Key: "btn.cancel", Language: "en", Value: "Cancel", Module: "button"},
	{Key: "btn.confirm", Language: "en", Value: "Confirm", Module: "button"},
	{Key: "btn.download", Language: "en", Value: "Download", Module: "button"},
	{Key: "btn.inquire_now", Language: "en", Value: "Inquire Now", Module: "button"},
	// form
	{Key: "form.name", Language: "en", Value: "Name", Module: "form"},
	{Key: "form.email", Language: "en", Value: "Email", Module: "form"},
	{Key: "form.phone", Language: "en", Value: "Phone", Module: "form"},
	{Key: "form.company", Language: "en", Value: "Company", Module: "form"},
	{Key: "form.country", Language: "en", Value: "Country", Module: "form"},
	{Key: "form.message", Language: "en", Value: "Message", Module: "form"},
	{Key: "form.quantity", Language: "en", Value: "Quantity", Module: "form"},
	// footer
	{Key: "footer.about", Language: "en", Value: "About Us", Module: "footer"},
	{Key: "footer.quick_links", Language: "en", Value: "Quick Links", Module: "footer"},
	{Key: "footer.contact_info", Language: "en", Value: "Contact Info", Module: "footer"},
	{Key: "footer.copyright", Language: "en", Value: "All Rights Reserved", Module: "footer"},
	// error
	{Key: "error.not_found", Language: "en", Value: "Not Found", Module: "error"},
	{Key: "error.load_failed", Language: "en", Value: "Failed to load data", Module: "error"},
}

// 默认中文词条
var defaultI18nEntriesZh = []models.I18nEntry{
	{Key: "common.learn_more", Language: "zh", Value: "了解更多", Module: "common"},
	{Key: "common.view_all", Language: "zh", Value: "查看全部", Module: "common"},
	{Key: "common.read_more", Language: "zh", Value: "阅读更多", Module: "common"},
	{Key: "common.contact_us", Language: "zh", Value: "联系我们", Module: "common"},
	{Key: "common.get_quote", Language: "zh", Value: "获取报价", Module: "common"},
	{Key: "common.send_inquiry", Language: "zh", Value: "发送询价", Module: "common"},
	{Key: "common.view_details", Language: "zh", Value: "查看详情", Module: "common"},
	{Key: "common.back_to_top", Language: "zh", Value: "返回顶部", Module: "common"},
	{Key: "nav.home", Language: "zh", Value: "首页", Module: "nav"},
	{Key: "nav.products", Language: "zh", Value: "产品中心", Module: "nav"},
	{Key: "nav.oem", Language: "zh", Value: "OEM 服务", Module: "nav"},
	{Key: "nav.odm", Language: "zh", Value: "ODM 服务", Module: "nav"},
	{Key: "nav.factory", Language: "zh", Value: "工厂展示", Module: "nav"},
	{Key: "nav.cases", Language: "zh", Value: "案例展示", Module: "nav"},
	{Key: "nav.blog", Language: "zh", Value: "博客", Module: "nav"},
	{Key: "nav.faq", Language: "zh", Value: "常见问题", Module: "nav"},
	{Key: "nav.contact", Language: "zh", Value: "联系我们", Module: "nav"},
	{Key: "btn.submit", Language: "zh", Value: "提交", Module: "button"},
	{Key: "btn.cancel", Language: "zh", Value: "取消", Module: "button"},
	{Key: "btn.confirm", Language: "zh", Value: "确认", Module: "button"},
	{Key: "btn.download", Language: "zh", Value: "下载", Module: "button"},
	{Key: "btn.inquire_now", Language: "zh", Value: "立即询价", Module: "button"},
	{Key: "form.name", Language: "zh", Value: "姓名", Module: "form"},
	{Key: "form.email", Language: "zh", Value: "邮箱", Module: "form"},
	{Key: "form.phone", Language: "zh", Value: "电话", Module: "form"},
	{Key: "form.company", Language: "zh", Value: "公司", Module: "form"},
	{Key: "form.country", Language: "zh", Value: "国家", Module: "form"},
	{Key: "form.message", Language: "zh", Value: "留言", Module: "form"},
	{Key: "form.quantity", Language: "zh", Value: "数量", Module: "form"},
	{Key: "footer.about", Language: "zh", Value: "关于我们", Module: "footer"},
	{Key: "footer.quick_links", Language: "zh", Value: "快速链接", Module: "footer"},
	{Key: "footer.contact_info", Language: "zh", Value: "联系方式", Module: "footer"},
	{Key: "footer.copyright", Language: "zh", Value: "保留所有权利", Module: "footer"},
	{Key: "error.not_found", Language: "zh", Value: "未找到", Module: "error"},
	{Key: "error.load_failed", Language: "zh", Value: "数据加载失败", Module: "error"},
}

// 默认西班牙语词条
var defaultI18nEntriesEs = []models.I18nEntry{
	{Key: "common.learn_more", Language: "es", Value: "Más Información", Module: "common"},
	{Key: "common.view_all", Language: "es", Value: "Ver Todo", Module: "common"},
	{Key: "common.read_more", Language: "es", Value: "Leer Más", Module: "common"},
	{Key: "common.contact_us", Language: "es", Value: "Contáctenos", Module: "common"},
	{Key: "common.get_quote", Language: "es", Value: "Solicitar Presupuesto", Module: "common"},
	{Key: "common.send_inquiry", Language: "es", Value: "Enviar Consulta", Module: "common"},
	{Key: "common.view_details", Language: "es", Value: "Ver Detalles", Module: "common"},
	{Key: "common.back_to_top", Language: "es", Value: "Volver Arriba", Module: "common"},
	{Key: "nav.home", Language: "es", Value: "Inicio", Module: "nav"},
	{Key: "nav.products", Language: "es", Value: "Productos", Module: "nav"},
	{Key: "nav.oem", Language: "es", Value: "Servicio OEM", Module: "nav"},
	{Key: "nav.odm", Language: "es", Value: "Servicio ODM", Module: "nav"},
	{Key: "nav.factory", Language: "es", Value: "Fábrica", Module: "nav"},
	{Key: "nav.cases", Language: "es", Value: "Casos", Module: "nav"},
	{Key: "nav.blog", Language: "es", Value: "Blog", Module: "nav"},
	{Key: "nav.faq", Language: "es", Value: "FAQ", Module: "nav"},
	{Key: "nav.contact", Language: "es", Value: "Contacto", Module: "nav"},
	{Key: "btn.submit", Language: "es", Value: "Enviar", Module: "button"},
	{Key: "btn.cancel", Language: "es", Value: "Cancelar", Module: "button"},
	{Key: "btn.confirm", Language: "es", Value: "Confirmar", Module: "button"},
	{Key: "btn.download", Language: "es", Value: "Descargar", Module: "button"},
	{Key: "btn.inquire_now", Language: "es", Value: "Consultar Ahora", Module: "button"},
	{Key: "form.name", Language: "es", Value: "Nombre", Module: "form"},
	{Key: "form.email", Language: "es", Value: "Correo Electrónico", Module: "form"},
	{Key: "form.phone", Language: "es", Value: "Teléfono", Module: "form"},
	{Key: "form.company", Language: "es", Value: "Empresa", Module: "form"},
	{Key: "form.country", Language: "es", Value: "País", Module: "form"},
	{Key: "form.message", Language: "es", Value: "Mensaje", Module: "form"},
	{Key: "form.quantity", Language: "es", Value: "Cantidad", Module: "form"},
	{Key: "footer.about", Language: "es", Value: "Sobre Nosotros", Module: "footer"},
	{Key: "footer.quick_links", Language: "es", Value: "Enlaces Rápidos", Module: "footer"},
	{Key: "footer.contact_info", Language: "es", Value: "Información de Contacto", Module: "footer"},
	{Key: "footer.copyright", Language: "es", Value: "Todos los Derechos Reservados", Module: "footer"},
	{Key: "error.not_found", Language: "es", Value: "No Encontrado", Module: "error"},
	{Key: "error.load_failed", Language: "es", Value: "Error al cargar datos", Module: "error"},
}

// 默认法语词条
var defaultI18nEntriesFr = []models.I18nEntry{
	{Key: "common.learn_more", Language: "fr", Value: "En Savoir Plus", Module: "common"},
	{Key: "common.view_all", Language: "fr", Value: "Voir Tout", Module: "common"},
	{Key: "common.read_more", Language: "fr", Value: "Lire Plus", Module: "common"},
	{Key: "common.contact_us", Language: "fr", Value: "Contactez-Nous", Module: "common"},
	{Key: "common.get_quote", Language: "fr", Value: "Demander un Devis", Module: "common"},
	{Key: "common.send_inquiry", Language: "fr", Value: "Envoyer la Demande", Module: "common"},
	{Key: "common.view_details", Language: "fr", Value: "Voir Détails", Module: "common"},
	{Key: "common.back_to_top", Language: "fr", Value: "Retour en Haut", Module: "common"},
	{Key: "nav.home", Language: "fr", Value: "Accueil", Module: "nav"},
	{Key: "nav.products", Language: "fr", Value: "Produits", Module: "nav"},
	{Key: "nav.oem", Language: "fr", Value: "Service OEM", Module: "nav"},
	{Key: "nav.odm", Language: "fr", Value: "Service ODM", Module: "nav"},
	{Key: "nav.factory", Language: "fr", Value: "Usine", Module: "nav"},
	{Key: "nav.cases", Language: "fr", Value: "Réalisations", Module: "nav"},
	{Key: "nav.blog", Language: "fr", Value: "Blog", Module: "nav"},
	{Key: "nav.faq", Language: "fr", Value: "FAQ", Module: "nav"},
	{Key: "nav.contact", Language: "fr", Value: "Contact", Module: "nav"},
	{Key: "btn.submit", Language: "fr", Value: "Soumettre", Module: "button"},
	{Key: "btn.cancel", Language: "fr", Value: "Annuler", Module: "button"},
	{Key: "btn.confirm", Language: "fr", Value: "Confirmer", Module: "button"},
	{Key: "btn.download", Language: "fr", Value: "Télécharger", Module: "button"},
	{Key: "btn.inquire_now", Language: "fr", Value: "Demander Maintenant", Module: "button"},
	{Key: "form.name", Language: "fr", Value: "Nom", Module: "form"},
	{Key: "form.email", Language: "fr", Value: "Email", Module: "form"},
	{Key: "form.phone", Language: "fr", Value: "Téléphone", Module: "form"},
	{Key: "form.company", Language: "fr", Value: "Société", Module: "form"},
	{Key: "form.country", Language: "fr", Value: "Pays", Module: "form"},
	{Key: "form.message", Language: "fr", Value: "Message", Module: "form"},
	{Key: "form.quantity", Language: "fr", Value: "Quantité", Module: "form"},
	{Key: "footer.about", Language: "fr", Value: "À Propos de Nous", Module: "footer"},
	{Key: "footer.quick_links", Language: "fr", Value: "Liens Rapides", Module: "footer"},
	{Key: "footer.contact_info", Language: "fr", Value: "Informations de Contact", Module: "footer"},
	{Key: "footer.copyright", Language: "fr", Value: "Tous Droits Réservés", Module: "footer"},
	{Key: "error.not_found", Language: "fr", Value: "Non Trouvé", Module: "error"},
	{Key: "error.load_failed", Language: "fr", Value: "Échec du chargement", Module: "error"},
}
