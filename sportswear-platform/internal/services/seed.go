package services

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sportswear-platform/internal/models"
)

// Seed 初始化示例数据（仅当数据为空时插入）
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		log.Println("[Seed] 已有产品数据，跳过")
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
}

// strToUUIDPtr 转换字符串为*uuid.UUID
func strToUUIDPtr(s string) *uuid.UUID {
	uid, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &uid
}
