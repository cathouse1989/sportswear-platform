package models

// I18nEntry 国际化词条（i18n 词典）
// 服务端集中配置 UI 文案的多语言翻译，
// 前端通过 GET /api/v1/i18n?lang=en 获取对应语言词条，一键切换界面语言。
type I18nEntry struct {
	BaseModel
	Key       string `gorm:"type:varchar(255);index;not null" json:"key"`         // 词条键，如 nav.home、btn.contact、footer.copyright
	Language  string `gorm:"type:varchar(10);index;not null" json:"language"`     // 语言代码 en/zh/es/fr/de
	Value     string `gorm:"type:text" json:"value"`                              // 文案内容
	Module    string `gorm:"type:varchar(50);index;default:common" json:"module"` // 模块分组：common/nav/button/footer/form/error
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}
