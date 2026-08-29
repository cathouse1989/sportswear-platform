package services

import (
	"sportswear-platform/internal/models"
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
	if lang == "" || lang == "en" {
		return // 默认英文，直接用主表字段
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
			break
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
