package models

import (
	"time"

	"github.com/google/uuid"
)

// PageVersion 页面版本（体系化发布的核心：草稿与线上分离）
// 工作流：后台编辑 → 保存为草稿版本（不影响线上）→ 预览确认 → 发布该版本 → 前台立即生效
type PageVersion struct {
	BaseModel
	PageID  uuid.UUID `gorm:"type:uuid;index;not null" json:"page_id"`
	Version int       `gorm:"not null" json:"version"`                            // 版本号，递增
	Status  string    `gorm:"type:varchar(20);default:draft;index" json:"status"` // draft / published / archived
	// 完整快照：页面基础信息 + 全部模块配置（JSONB）
	Snapshot    string     `gorm:"type:jsonb" json:"snapshot"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedBy   *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	Note        string     `gorm:"type:varchar(500)" json:"note"` // 版本备注
}

// ThemeConfig 门户主题配置（全局展示效果）
// 后台可配置：主题色、字体、头部样式、按钮风格、圆角、间距等
// 前端读取后渲染整体视觉风格
type ThemeConfig struct {
	BaseModel
	Key      string `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value    string `gorm:"type:text" json:"value"`              // JSON 配置值
	Group    string `gorm:"type:varchar(50);index" json:"group"` // color, typography, layout, header, footer, button
	Name     string `gorm:"type:varchar(100)" json:"name"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// 默认主题配置项
var DefaultThemeConfigs = []ThemeConfig{
	{Key: "primary_color", Value: `"#2563eb"`, Group: "color", Name: "主色调", IsActive: true},
	{Key: "secondary_color", Value: `"#1e40af"`, Group: "color", Name: "辅助色", IsActive: true},
	{Key: "accent_color", Value: `"#f59e0b"`, Group: "color", Name: "强调色（CTA 按钮）", IsActive: true},
	{Key: "font_family", Value: `"Inter, system-ui, sans-serif"`, Group: "typography", Name: "字体", IsActive: true},
	{Key: "heading_size", Value: `"large"`, Group: "typography", Name: "标题尺寸", IsActive: true},
	{Key: "header_style", Value: `"sticky"`, Group: "header", Name: "导航栏样式（sticky/transparent/solid）", IsActive: true},
	{Key: "header_height", Value: `72`, Group: "header", Name: "导航栏高度(px)", IsActive: true},
	// 品牌 Logo（管理后台顶栏使用；logo_url 为空时前端回退显示 logo_alt 文字）
	{Key: "logo_url", Value: ``, Group: "header", Name: "Logo 图片地址", IsActive: true},
	{Key: "logo_alt", Value: `"Sportswear"`, Group: "header", Name: "Logo 文字（无图片时显示）", IsActive: true},
	{Key: "button_radius", Value: `8`, Group: "button", Name: "按钮圆角(px)", IsActive: true},
	{Key: "card_radius", Value: `12`, Group: "layout", Name: "卡片圆角(px)", IsActive: true},
	{Key: "section_spacing", Value: `80`, Group: "layout", Name: "区块间距(px)", IsActive: true},
	{Key: "container_width", Value: `1280`, Group: "layout", Name: "内容区最大宽度(px)", IsActive: true},
	{Key: "footer_columns", Value: `4`, Group: "footer", Name: "Footer 栏数", IsActive: true},
	{Key: "whatsapp_number", Value: `"8612345678900"`, Group: "contact", Name: "WhatsApp 号码（含国家代码，不含+号）", IsActive: true},
	{Key: "whatsapp_message", Value: `"Hello! I am interested in your sportswear products."`, Group: "contact", Name: "WhatsApp 默认消息", IsActive: true},
	{Key: "contact_email", Value: `"info@sportswear.com"`, Group: "contact", Name: "联系邮箱", IsActive: true},
	{Key: "contact_phone", Value: `"+86 123 4567 8900"`, Group: "contact", Name: "联系电话", IsActive: true},
	// 社交媒体链接
	{Key: "social_youtube", Value: `"https://youtube.com/@sportswear"`, Group: "social", Name: "YouTube 链接", IsActive: true},
	{Key: "social_instagram", Value: `"https://instagram.com/sportswear"`, Group: "social", Name: "Instagram 链接", IsActive: true},
	{Key: "social_xiaohongshu", Value: `"https://xiaohongshu.com/user/sportswear"`, Group: "social", Name: "小红书链接", IsActive: true},
	{Key: "social_facebook", Value: `"https://facebook.com/sportswear"`, Group: "social", Name: "Facebook 链接", IsActive: true},
	{Key: "social_twitter", Value: `"https://twitter.com/sportswear"`, Group: "social", Name: "Twitter/X 链接", IsActive: true},
	{Key: "social_linkedin", Value: `"https://linkedin.com/company/sportswear"`, Group: "social", Name: "LinkedIn 链接", IsActive: true},
	// 翻页配置
	{Key: "products_page_size", Value: "20", Group: "pagination", Name: "产品列表每页条数", IsActive: true},
	{Key: "blogs_page_size", Value: "20", Group: "pagination", Name: "博客列表每页条数", IsActive: true},
	{Key: "cases_page_size", Value: "20", Group: "pagination", Name: "案例列表每页条数", IsActive: true},
	{Key: "faqs_page_size", Value: "20", Group: "pagination", Name: "FAQ 每页条数", IsActive: true},
	{Key: "admin_page_size", Value: "20", Group: "pagination", Name: "管理后台默认每页条数", IsActive: true},
	// Hero 轮播配置
	{Key: "hero_autoplay", Value: "true", Group: "hero", Name: "首页轮播自动播放", IsActive: true},
	{Key: "hero_interval_ms", Value: "5000", Group: "hero", Name: "轮播间隔(毫秒)", IsActive: true},
	{Key: "hero_transition", Value: "fade", Group: "hero", Name: "轮播过渡动画(fade/slide)", IsActive: true},
	{Key: "hero_show_dots", Value: "true", Group: "hero", Name: "显示轮播指示点", IsActive: true},
	{Key: "hero_show_arrows", Value: "true", Group: "hero", Name: "显示左右箭头", IsActive: true},
	{Key: "hero_pause_on_hover", Value: "true", Group: "hero", Name: "悬停暂停轮播", IsActive: true},
}
