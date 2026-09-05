package services

import (
	"encoding/json"
	"regexp"
	"strings"

	"sportswear-backend/internal/models"
)

// ==================== 内容本地化解析 ====================
//
// 语言回退链：目标语言翻译 → 主表默认内容（英文）
// 主表字段（Title/Brief/Description 等）保存默认语言内容，
// 翻译表保存其他语言内容；无目标语言翻译时回退到主表默认内容。

// LocalizeProduct 按语言解析产品翻译，覆盖到本地化展示字段
func LocalizeProduct(p *models.Product, lang string) {
	if p == nil {
		return
	}
	// 默认英文：主表字段即英文内容；Name 非 DB 列（gorm:"-"），须从 en 翻译表回填
	if lang == "" || lang == "en" {
		for _, t := range p.Translations {
			if t.Language == "en" && t.Status == models.TranslationStatusPublished {
				if t.Name != "" {
					p.Name = t.Name
				}
				break
			}
		}
		return
	}
	for _, t := range p.Translations {
		if t.Language == lang && t.Status == models.TranslationStatusPublished {
			if t.Name != "" {
				p.Name = t.Name
			}
			if t.Brief != "" {
				p.Brief = t.Brief
			}
			if t.Description != "" {
				p.Description = t.Description
			}
			if t.Features != "" {
				p.Features = t.Features
			}
			if t.Usage != "" {
				p.Usage = t.Usage
			}
			// 规格类字段按语言覆盖（无翻译时回退主表英文源）
			if t.Material != "" {
				p.Material = t.Material
			}
			if t.Composition != "" {
				p.Composition = t.Composition
			}
			if t.Weight != "" {
				p.Weight = t.Weight
			}
			if t.Elasticity != "" {
				p.Elasticity = t.Elasticity
			}
			if t.Fit != "" {
				p.Fit = t.Fit
			}
			if t.SupportLevel != "" {
				p.SupportLevel = t.SupportLevel
			}
			if t.Season != "" {
				p.Season = t.Season
			}
			if t.SizeRange != "" {
				p.SizeRange = t.SizeRange
			}
			break
		}
	}
	// 定制能力 note 按语言解析（JSONB 翻译字段，未配置时回退英文 Note）
	localizeCustomizations(p.Customizations, lang)
}

// resolveTranslated 解析 JSONB 翻译字段（如 {"zh":"...","es":"...","fr":"..."}），返回目标语言文案。
// 空值 / 英文 / 解析失败时返回空字符串（调用方回退主表英文源）。
func resolveTranslated(translationsJSON string, lang string) string {
	if translationsJSON == "" || lang == "" || lang == "en" {
		return ""
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(translationsJSON), &m); err != nil {
		return ""
	}
	return m[lang]
}

// localizeCustomizations 批量本地化定制能力的 note
func localizeCustomizations(customizations []models.ProductCustomization, lang string) {
	for i := range customizations {
		if note := resolveTranslated(customizations[i].Translations, lang); note != "" {
			customizations[i].Note = note
		}
	}
}

// LocalizeProducts 批量本地化产品
func LocalizeProducts(products []models.Product, lang string) {
	for i := range products {
		LocalizeProduct(&products[i], lang)
	}
}

// LocalizePage 按语言解析页面翻译
func LocalizePage(p *models.Page, lang string) {
	if p == nil {
		return
	}
	if lang == "" || lang == "en" {
		return
	}
	for _, t := range p.Translations {
		if t.Language == lang && t.Status == models.TranslationStatusPublished {
			if t.Title != "" {
				p.Title = t.Title
			}
			break
		}
	}
}

// LocalizeBlog 按语言解析博客翻译
func LocalizeBlog(b *models.Blog, lang string) {
	if b == nil {
		return
	}
	if lang == "" || lang == "en" {
		return
	}
	for _, t := range b.Translations {
		if t.Language == lang && t.Status == models.TranslationStatusPublished {
			if t.Title != "" {
				b.Title = t.Title
			}
			if t.Content != "" {
				b.Content = t.Content
			}
			break
		}
	}
}

// LocalizeBlogs 批量本地化博客
func LocalizeBlogs(blogs []models.Blog, lang string) {
	for i := range blogs {
		LocalizeBlog(&blogs[i], lang)
	}
}

// LocalizeCase 按语言解析案例翻译
func LocalizeCase(c *models.Case, lang string) {
	if c == nil {
		return
	}
	if lang == "" || lang == "en" {
		return
	}
	for _, t := range c.Translations {
		if t.Language == lang && t.Status == models.TranslationStatusPublished {
			if t.Title != "" {
				c.Title = t.Title
			}
			if t.ClientNeed != "" {
				c.ClientNeed = t.ClientNeed
			}
			if t.Problem != "" {
				c.Problem = t.Problem
			}
			if t.Solution != "" {
				c.Solution = t.Solution
			}
			if t.Process != "" {
				c.Process = t.Process
			}
			if t.Result != "" {
				c.Result = t.Result
			}
			break
		}
	}
}

// LocalizeCases 批量本地化案例
func LocalizeCases(cases []models.Case, lang string) {
	for i := range cases {
		LocalizeCase(&cases[i], lang)
	}
}

// 摘要提取用正则（编译一次，供列表接口复用，避免逐条重复编译）。
// 注意：Go regexp 基于 RE2，不支持反向引用（\1），故 script/style 分别匹配。
var (
	excerptScriptRe = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script>`)
	excerptStyleRe  = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style>`)
	excerptTagRe    = regexp.MustCompile(`(?s)<[^>]*>`)
	excerptSpaceRe  = regexp.MustCompile(`\s+`)
)

// ExcerptHTML 提取 HTML 纯文本摘要：去除标签/脚本，压缩空白并按字符截断。
// 用于列表接口瘦身——列表不下发全文 content，仅返回摘要。
func ExcerptHTML(html string, max int) string {
	text := excerptScriptRe.ReplaceAllString(html, "")
	text = excerptStyleRe.ReplaceAllString(text, "")
	text = excerptTagRe.ReplaceAllString(text, " ")
	// 解码常见 HTML 实体
	replacer := strings.NewReplacer("&nbsp;", " ", "&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", "\"", "&#39;", "'")
	text = replacer.Replace(text)
	text = excerptSpaceRe.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return strings.TrimSpace(string(runes[:max])) + "…"
}
