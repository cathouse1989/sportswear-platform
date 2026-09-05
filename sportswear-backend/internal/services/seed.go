package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-backend/internal/models"
)

// Seed 初始化示例数据（仅当数据为空时插入）
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		log.Println("[Seed] 已有产品数据，跳过初始种子")
		// 幂等回填已存在的演示产品（多语言 + 下级页面内容）
		ensureDemoProductContent(db)
		return
	}
	log.Println("[Seed] 创建示例数据...")

	// 分类
	cats := []models.Category{
		{Name: "Yoga Wear", Slug: "yoga-wear", SortOrder: 1, IsActive: true, Image: "https://images.unsplash.com/photo-1599901860904-17e6ed7083a0?w=400&h=600&fit=crop"},
		{Name: "Running Gear", Slug: "running-gear", SortOrder: 2, IsActive: true, Image: "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=400&h=600&fit=crop"},
		{Name: "Training Apparel", Slug: "training-apparel", SortOrder: 3, IsActive: true, Image: "https://images.unsplash.com/photo-1577221084712-45b0445d2b00?w=400&h=600&fit=crop"},
		{Name: "Team Uniforms", Slug: "team-uniforms", SortOrder: 4, IsActive: true, Image: "https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=400&h=600&fit=crop"},
		{Name: "Custom Design", Slug: "custom-design", SortOrder: 5, IsActive: true, Image: ""},
	}
	for i := range cats {
		db.Where("slug = ?", cats[i].Slug).FirstOrCreate(&cats[i])
	}

	// 获取分类ID映射
	var catList []models.Category
	db.Find(&catList)
	catMap := make(map[string]string)
	for _, c := range catList {
		catMap[c.Slug] = c.ID.String()
	}

	// 产品
	type productSeed struct {
		sku, slug, brief, desc, material, catSlug, gender, img string
		featured                                               bool
	}
	products := []productSeed{
		{"SW-YG-001", "premium-yoga-leggings", "High-waist yoga leggings with moisture-wicking fabric. Perfect for studio and outdoor practice.", "Our premium yoga leggings are crafted from custom-blended fabrics that offer exceptional stretch, recovery, and breathability.", "Nylon 75% + Spandex 25%", "yoga-wear", "female", "https://images.unsplash.com/photo-1599901860904-17e6ed7083a0?w=800&h=1000&fit=crop", true},
		{"SW-YT-002", "eco-friendly-yoga-top", "Lightweight yoga top made from recycled materials. Racerback design with built-in shelf bra.", "Made from 100% recycled polyester, this yoga top combines sustainability with performance.", "Recycled Polyester 88% + Spandex 12%", "yoga-wear", "female", "https://images.unsplash.com/photo-1596755389378-c31d21fd1273?w=800&h=1000&fit=crop", true},
		{"SW-RS-003", "professional-running-shorts", "Ultra-light running shorts with built-in liner and zip pocket. 4-way stretch fabric.", "Engineered for speed with compression liner support and laser-cut ventilation.", "Polyester 90% + Spandex 10%", "running-gear", "male", "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=800&h=1000&fit=crop", true},
		{"SW-RT-004", "reflective-running-jacket", "Windproof running jacket with 360° reflectivity. Stay visible, stay safe.", "Lightweight and packable with water-resistant coating. Ideal for early morning runs.", "Nylon 85% + Polyester 15%", "running-gear", "unisex", "https://images.unsplash.com/photo-1591047139829-d91aecb6caea?w=800&h=1000&fit=crop", true},
		{"SW-TT-005", "compression-training-tee", "Muscle compression tee with odor-resistant technology. Gradient compression.", "Anti-odor treatment keeps you fresh during intense workouts. Flatlock seams for zero distractions.", "Nylon 80% + Spandex 20%", "training-apparel", "male", "https://images.unsplash.com/photo-1577221084712-45b0445d2b00?w=800&h=1000&fit=crop", true},
		{"SW-TB-006", "training-duffle-bag", "Spacious 45L duffle with wet/dry compartment. Perfect for gym and travel.", "Multiple pockets with padded shoulder strap. Water-resistant bottom protects against wet floors.", "Polyester 600D", "training-apparel", "unisex", "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=800&h=1000&fit=crop", true},
		{"SW-UT-007", "pro-team-soccer-jersey", "Professional soccer jersey with heat-sealed crest. Moisture-wicking fabric.", "Heat-sealed club crest and sponsor logos that won't peel or crack. Tailored athletic fit.", "Polyester 100%", "team-uniforms", "male", "https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=800&h=1000&fit=crop", true},
		{"SW-UB-008", "custom-basketball-uniform", "Full custom basketball uniform with sublimation printing. Breathable mesh fabric.", "Complete set including jersey and shorts. Full sublimation ensures vibrant colors that never fade.", "Polyester Mesh 100%", "team-uniforms", "unisex", "https://images.unsplash.com/photo-1546519638-68e109498ffc?w=800&h=1000&fit=crop", true},
		{"SW-RS-009", "women-running-tights", "High-performance running tights with compression fit and reflective details.", "Zippered pockets, moisture-wicking fabric, and reflective elements for night safety.", "Nylon 70% + Spandex 30%", "running-gear", "female", "https://images.unsplash.com/photo-1599407950132-77e6e1e3b0c0?w=800&h=1000&fit=crop", true},
		{"SW-YG-010", "mens-yoga-pants", "Comfortable men's yoga pants with straight leg and elastic waistband.", "Breathable cotton-blend fabric with 4-way stretch for unrestricted movement.", "Cotton 60% + Polyester 35% + Spandex 5%", "yoga-wear", "male", "https://images.unsplash.com/photo-1591942632499-1b8bd0b1e3b0?w=800&h=1000&fit=crop", true},
	}

	for _, p := range products {
		product := models.Product{
			SKU: p.sku, Slug: p.slug, Brief: p.brief, Description: p.desc,
			Material: p.material, Gender: models.Gender(p.gender),
			CoverImage: p.img, Status: models.ProductStatusPublished,
			IsFeatured: p.featured, SortOrder: 1,
			SampleMOQ: 50, ProductionMOQ: 500, ColorMOQ: 200, SizeMOQ: 200,
		}
		if cid, ok := catMap[p.catSlug]; ok {
			product.CategoryID = strToUUIDPtr(cid)
		}
		db.Create(&product)

		// 产品图片
		db.Create(&models.ProductImage{
			ProductID: product.ID, Type: "main", URL: p.img,
			Alt: p.sku, SortOrder: 1,
		})

		// 翻译
		db.Create(&models.ProductTranslation{
			ProductID: product.ID, Language: "en",
			Name: p.sku, Brief: p.brief, Description: p.desc,
			Status: models.TranslationStatusPublished,
		})
	}

	// 博客
	blogs := []models.Blog{
		{Title: "The Future of Sustainable Sportswear: 2025 Trends", Slug: "sustainable-sportswear-2025", Category: "Industry Insights", Content: "The sportswear industry is undergoing a massive transformation towards sustainability. From recycled ocean plastics to biodegradable fabrics, manufacturers are innovating to reduce environmental impact. Key trends include closed-loop recycling systems, waterless dyeing technologies, plant-based performance fabrics, and circular design principles.", Status: "published", PublishedAt: &[]time.Time{time.Now()}[0]},
		{Title: "How to Choose the Right Sportswear Manufacturer", Slug: "choose-sportswear-manufacturer", Category: "Guide", Content: "Selecting the right manufacturing partner is crucial for your brand's success. This guide covers: evaluating factory certifications (ISO 9001, BSCI, OEKO-TEX), understanding MOQ requirements, assessing quality control processes, and building long-term partnerships.", Status: "published", PublishedAt: &[]time.Time{time.Now()}[0]},
		{Title: "Custom Sportswear Design: From Concept to Production", Slug: "custom-sportswear-design-process", Category: "Production", Content: "Ever wondered how your custom sportswear design becomes a reality? This behind-the-scenes look covers: initial consultation, fabric selection, pattern making, prototype development, bulk production, and quality control.", Status: "published", PublishedAt: &[]time.Time{time.Now()}[0]},
		{Title: "Understanding Fabric Technology in Performance Wear", Slug: "fabric-technology-performance-wear", Category: "Technology", Content: "Modern sportswear relies on advanced fabric technologies. This article explains moisture-wicking mechanisms, compression benefits, thermal regulation, anti-microbial treatments, UV protection ratings, and stretch properties.", Status: "published", PublishedAt: &[]time.Time{time.Now()}[0]},
	}
	for _, b := range blogs {
		db.Where("slug = ?", b.Slug).FirstOrCreate(&b)
	}

	// 认证
	certs := []models.Certification{
		{Name: "ISO 9001:2015", Description: "Quality Management System certified", Status: "published"},
		{Name: "BSCI Compliance", Description: "Amfori BSCI certified — ethical manufacturing", Status: "published"},
		{Name: "OEKO-TEX Standard 100", Description: "Product safety certification", Status: "published"},
		{Name: "SEDEX Registered", Description: "Supplier Ethical Data Exchange", Status: "published"},
	}
	for _, c := range certs {
		db.Where("name = ?", c.Name).FirstOrCreate(&c)
	}

	// FAQ
	faqs := []models.FAQ{
		{Question: "What is the minimum order quantity (MOQ)?", Answer: "Our standard MOQ is 300 units per design. Sample orders start at 50 units.", Category: "Orders", Language: "en", IsActive: true, SortOrder: 1},
		{Question: "How long does production take?", Answer: "Sample production: 7-10 days. Bulk production: 4-6 weeks. Rush orders available.", Category: "Production", Language: "en", IsActive: true, SortOrder: 2},
		{Question: "Can you customize fabric and colors?", Answer: "Yes, we offer full fabric customization including custom colors, compositions, and finishes.", Category: "Customization", Language: "en", IsActive: true, SortOrder: 3},
		{Question: "What certifications does your factory hold?", Answer: "ISO 9001:2015, BSCI, OEKO-TEX Standard 100, and SEDEX certified.", Category: "Quality", Language: "en", IsActive: true, SortOrder: 4},
		{Question: "Do you offer private label services?", Answer: "Absolutely. Full private label services including custom branding, packaging, and labeling.", Category: "Services", Language: "en", IsActive: true, SortOrder: 5},
	}
	for _, f := range faqs {
		db.Where("question = ?", f.Question).FirstOrCreate(&f)
	}

	log.Println("[Seed] 示例数据创建完成: 5个分类, 10个产品, 4篇博客, 4个认证, 5个FAQ")

	// 幂等回填演示产品的多语言翻译与下级页面内容
	ensureDemoProductContent(db)
}

// ====================================================================
// 演示产品内容回填（幂等）：为示例产品（SKU 前缀 SW-）补齐多语言翻译与下级页面内容
// ====================================================================

// demoCatWord 各分类在 en/zh/es/fr 下的展示名
var demoCatWordMap = map[string]map[string]string{
	"yoga-wear":        {"en": "Yoga Wear", "zh": "瑜伽服", "es": "ropa de yoga", "fr": "vêtements de yoga"},
	"running-gear":     {"en": "Running Gear", "zh": "跑步装备", "es": "equipo de running", "fr": "équipement de course"},
	"training-apparel": {"en": "Training Apparel", "zh": "训练服饰", "es": "ropa de entrenamiento", "fr": "vêtements d'entraînement"},
	"team-uniforms":    {"en": "Team Uniforms", "zh": "团队队服", "es": "uniformes de equipo", "fr": "uniformes d'équipe"},
	"custom-design":    {"en": "Custom Design", "zh": "定制设计", "es": "diseño personalizado", "fr": "design personnalisé"},
}

func demoCatWord(lang, catSlug string) string {
	if m, ok := demoCatWordMap[catSlug]; ok {
		if w := m[lang]; w != "" {
			return w
		}
	}
	return catSlug
}

// demoWeight 各分类示例克重
var demoWeight = map[string]string{
	"yoga-wear": "220g", "running-gear": "180g", "training-apparel": "200g",
	"team-uniforms": "190g", "custom-design": "200g",
}

// demoGallery 各分类的画廊示例图（Unsplash）
var demoGallery = map[string][]string{
	"yoga-wear": {
		"https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?w=800&h=1000&fit=crop",
		"https://images.unsplash.com/photo-1518611012118-696072aa579a?w=800&h=1000&fit=crop",
	},
	"running-gear": {
		"https://images.unsplash.com/photo-1476480862126-209bfaa8edc8?w=800&h=1000&fit=crop",
		"https://images.unsplash.com/photo-1571019613914-85f342c6a11e?w=800&h=1000&fit=crop",
	},
	"training-apparel": {
		"https://images.unsplash.com/photo-1576678927484-cc907957088c?w=800&h=1000&fit=crop",
		"https://images.unsplash.com/photo-1550259979-ed79b48d2a30?w=800&h=1000&fit=crop",
	},
	"team-uniforms": {
		"https://images.unsplash.com/photo-1517649763962-0c623066013b?w=800&h=1000&fit=crop",
		"https://images.unsplash.com/photo-1546519638-68e109498ffc?w=800&h=1000&fit=crop",
	},
	"custom-design": {
		"https://images.unsplash.com/photo-1558769132-cb1aea458c5e?w=800&h=1000&fit=crop",
		"https://images.unsplash.com/photo-1556905055-8f358a7a47b2?w=800&h=1000&fit=crop",
	},
}

// humanizeSlug 将 slug 转成可读英文名："premium-yoga-leggings" -> "Premium Yoga Leggings"
func humanizeSlug(slug string) string {
	parts := strings.Split(slug, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}

// localizedDemoCopy 按语言生成演示文案（brief/description/features/usage）
func localizedDemoCopy(lang, catSlug, material, enBrief, enDesc string) (string, string, string, string) {
	cat := demoCatWord(lang, catSlug)
	switch lang {
	case "zh":
		return fmt.Sprintf("专业%s定制运动服装，采用%s面料，支持颜色、尺码、Logo 等全维度定制，低起订量。", cat, material),
			fmt.Sprintf("我们提供高品质%s定制制造服务。本产品采用 %s 面料，支持颜色、尺码、Logo、版型等定制，欢迎索取样品与报价。", cat, material),
			fmt.Sprintf("- 高品质 %s 面料\n- 支持 OEM/ODM 全定制\n- 低起订量，快速打样\n- 15 年以上制造经验", material),
			"适用于运动健身、团队训练、瑜伽、跑步等场景。"
	case "es":
		return fmt.Sprintf("Ropa deportiva personalizada de %s con tejido %s. Personalización total de color, talla y logo, MOQ bajo.", cat, material),
			fmt.Sprintf("Ofrecemos fabricación personalizada de %s de alta calidad. Este producto utiliza %s, con personalización de color, talla, logo y patrón.", cat, material),
			fmt.Sprintf("- Tejido %s de alta calidad\n- Personalización OEM/ODM completa\n- MOQ bajo con muestras rápidas\n- Más de 15 años de experiencia", material),
			"Ideal para fitness, entrenamiento en equipo, yoga, running y más."
	case "fr":
		return fmt.Sprintf("Vêtements de sport personnalisés (%s) en tissu %s. Personnalisation complète couleur, taille et logo, MOQ bas.", cat, material),
			fmt.Sprintf("Nous offrons une fabrication personnalisée de %s de haute qualité. Ce produit utilise %s, avec personnalisation couleur, taille, logo et patron.", cat, material),
			fmt.Sprintf("- Tissu %s de haute qualité\n- Personnalisation OEM/ODM complète\n- MOQ bas avec échantillons rapides\n- Plus de 15 ans d'expérience", material),
			"Idéal pour fitness, entraînement en équipe, yoga, course et plus."
	default: // en
		return enBrief, enDesc,
			fmt.Sprintf("- High-quality %s fabric\n- Full OEM/ODM customization\n- Low MOQ with fast sampling\n- 15+ years manufacturing experience", material),
			"Ideal for fitness, team training, yoga, running and more."
	}
}

// buildDemoTranslation 构造一条演示翻译
func buildDemoTranslation(lang, catSlug, material, enName, enBrief, enDesc string) models.ProductTranslation {
	brief, desc, features, usage := localizedDemoCopy(lang, catSlug, material, enBrief, enDesc)
	return models.ProductTranslation{
		Language: lang, Name: enName, Brief: brief, Description: desc,
		Features: features, Usage: usage, Status: models.TranslationStatusPublished,
	}
}

// ensureProductTranslation 缺失则创建；已存在但为裸种子数据（name=SKU 或关键字段空）则补缺修正
func ensureProductTranslation(db *gorm.DB, productID uuid.UUID, sku string, t models.ProductTranslation) {
	var existing models.ProductTranslation
	if err := db.Where("product_id = ? AND language = ?", productID, t.Language).First(&existing).Error; err != nil {
		t.ProductID = productID
		db.Create(&t)
		return
	}
	updates := map[string]interface{}{}
	if existing.Name == "" || existing.Name == sku {
		updates["name"] = t.Name
	}
	if existing.Brief == "" {
		updates["brief"] = t.Brief
	}
	if existing.Description == "" {
		updates["description"] = t.Description
	}
	if existing.Features == "" {
		updates["features"] = t.Features
	}
	if existing.Usage == "" {
		updates["usage"] = t.Usage
	}
	if existing.Status != models.TranslationStatusPublished {
		updates["status"] = models.TranslationStatusPublished
	}
	if len(updates) > 0 {
		db.Model(&existing).Updates(updates)
	}
}

// ensureDemoSeries 确保各分类演示系列存在，返回 category-slug -> series-id 映射
func ensureDemoSeries(db *gorm.DB) map[string]string {
	seriesByCat := map[string]string{}
	defs := []struct{ cat, name, slug string }{
		{"yoga-wear", "Zen Studio Collection", "zen-studio-collection"},
		{"running-gear", "Aero Run Series", "aero-run-series"},
		{"training-apparel", "Pro Training Line", "pro-training-line"},
		{"team-uniforms", "Team Pro Series", "team-pro-series"},
		{"custom-design", "Custom Atelier", "custom-atelier"},
	}
	for _, d := range defs {
		var s models.Series
		if err := db.Where("slug = ?", d.slug).First(&s).Error; err != nil {
			s = models.Series{Name: d.name, Slug: d.slug, SortOrder: 1, Status: models.ProductStatusPublished, IsActive: true}
			db.Create(&s)
		} else if s.Status != models.ProductStatusPublished {
			db.Model(&s).Update("status", models.ProductStatusPublished)
		}
		seriesByCat[d.cat] = s.ID.String()
	}
	return seriesByCat
}

// enrichDemoProduct 为单个演示产品补齐主表属性、多语言翻译、规格、定制、图集与系列
func enrichDemoProduct(db *gorm.DB, p *models.Product, catSlug, seriesID string) {
	material := p.Material
	if material == "" {
		material = "high-performance fabric"
	}
	enName := humanizeSlug(p.Slug)
	_, _, enFeatures, enUsage := localizedDemoCopy("en", catSlug, material, p.Brief, p.Description)

	// 1) 主表标量属性：缺失时补齐
	updates := map[string]interface{}{}
	if p.Features == "" {
		updates["features"] = enFeatures
	}
	if p.Usage == "" {
		updates["usage"] = enUsage
	}
	if p.Composition == "" {
		updates["composition"] = material
	}
	if p.Weight == "" {
		updates["weight"] = demoWeight[catSlug]
	}
	if p.Elasticity == "" {
		updates["elasticity"] = "4-Way Stretch"
	}
	if p.Fit == "" {
		updates["fit"] = "Regular Fit"
	}
	if p.SupportLevel == "" {
		updates["support_level"] = "Medium"
	}
	if p.Season == "" {
		updates["season"] = "All Season"
	}
	if p.SizeRange == "" {
		updates["size_range"] = "XS-XXL"
	}
	if len(updates) > 0 {
		db.Model(p).Updates(updates)
	}

	// 2) 多语言翻译（en/zh/es/fr）
	for _, lang := range []string{"en", "zh", "es", "fr"} {
		ensureProductTranslation(db, p.ID, p.SKU, buildDemoTranslation(lang, catSlug, material, enName, p.Brief, p.Description))
	}

	// 3) 规格：缺失时补齐
	var specCount int64
	db.Model(&models.ProductSpec{}).Where("product_id = ?", p.ID).Count(&specCount)
	if specCount == 0 {
		specs := []models.ProductSpec{
			{Name: "Fabric", Value: material},
			{Name: "Composition", Value: material},
			{Name: "Weight", Value: demoWeight[catSlug]},
			{Name: "Fit", Value: "Regular Fit"},
			{Name: "Size Range", Value: "XS-XXL"},
			{Name: "MOQ", Value: fmt.Sprintf("%d pcs", p.ProductionMOQ)},
		}
		for i := range specs {
			specs[i].ProductID = p.ID
			specs[i].SortOrder = i
			db.Create(&specs[i])
		}
	}

	// 4) 定制能力：缺失时补齐
	var cusCount int64
	db.Model(&models.ProductCustomization{}).Where("product_id = ?", p.ID).Count(&cusCount)
	if cusCount == 0 {
		customizations := []models.ProductCustomization{
			{Type: "logo", IsEnabled: true, Note: "Custom logo via embroidery, print or heat transfer"},
			{Type: "color", IsEnabled: true, Note: "Custom colors with Pantone matching"},
			{Type: "fabric", IsEnabled: true, Note: "Fabric sourcing & custom blends"},
			{Type: "pattern", IsEnabled: true, Note: "Custom patterns & cuts"},
			{Type: "packaging", IsEnabled: true, Note: "Private label & custom packaging"},
		}
		for i := range customizations {
			customizations[i].ProductID = p.ID
			db.Create(&customizations[i])
		}
	}

	// 5) 图集：主图之外补齐画廊图
	var imgCount int64
	db.Model(&models.ProductImage{}).Where("product_id = ?", p.ID).Count(&imgCount)
	if imgCount < 3 {
		for _, u := range demoGallery[catSlug] {
			db.Create(&models.ProductImage{ProductID: p.ID, Type: "gallery", URL: u, Alt: p.SKU, SortOrder: int(imgCount) + 1})
			imgCount++
		}
	}

	// 6) 系列关联
	if seriesID != "" {
		db.Exec("INSERT INTO product_series (product_id, series_id) VALUES (?, ?) ON CONFLICT DO NOTHING", p.ID, seriesID)
	}
}

// ensureDemoProductContent 为所有演示产品（SKU 前缀 SW-）幂等回填多语言与下级页面内容
func ensureDemoProductContent(db *gorm.DB) {
	var products []models.Product
	if err := db.Where("sku LIKE ?", "SW-%").Find(&products).Error; err != nil {
		return
	}
	if len(products) == 0 {
		return
	}

	var cats []models.Category
	db.Find(&cats)
	catSlugByID := map[string]string{}
	for _, c := range cats {
		catSlugByID[c.ID.String()] = c.Slug
	}

	seriesByCat := ensureDemoSeries(db)

	for i := range products {
		p := &products[i]
		catSlug := "custom-design"
		if p.CategoryID != nil {
			if s, ok := catSlugByID[p.CategoryID.String()]; ok && s != "" {
				catSlug = s
			}
		}
		enrichDemoProduct(db, p, catSlug, seriesByCat[catSlug])
	}
	log.Printf("[Seed] 演示产品多语言/下级页面内容已就绪（共 %d 个）", len(products))
}

// strToUUIDPtr 转换字符串为*uuid.UUID
func strToUUIDPtr(s string) *uuid.UUID {
	uid, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &uid
}
