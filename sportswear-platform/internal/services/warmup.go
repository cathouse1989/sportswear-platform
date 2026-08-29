package services

import (
	"sync"

	"gorm.io/gorm"

	"sportswear-platform/internal/models"
)

// WarmUpPublicCache 服务启动时的门户公开接口缓存预热。
// 实现"部署时生成缓存"的目标：服务一旦启动、数据与 Redis 就绪后，
// 立即将门户热点数据写入缓存，让首轮访问即可命中，避免冷启动回源变慢。
// 尽力而为：任意一项预热失败仅跳过该项，不影响服务启动与其余缓存。
func WarmUpPublicCache(db *gorm.DB, cache *CacheService) {
	if cache == nil || !cache.IsEnabled() {
		return
	}

	productSvc := NewProductService(db)
	cmsSvc := NewCMSService(db)
	langs := []string{"en", "zh", "es", "fr"}

	for _, lang := range langs {
		// 首页（聚合多类内容，成本最高，优先预热）
		if home := buildHomeData(cmsSvc, productSvc, lang); home != nil {
			cache.Set("cache:home:"+lang, home, CacheTTLShort)
		}

		// 无需本地化的静态列表（key 与 public handler 保持一致）
		if list, err := productSvc.ListCategories(); err == nil {
			cache.Set("cache:categories:"+lang, list, CacheTTLMedium)
		}
		if list, err := productSvc.ListPublishedSeries(); err == nil {
			cache.Set("cache:series:"+lang, list, CacheTTLMedium)
		}
		if list, err := productSvc.ListPublishedFabrics(); err == nil {
			cache.Set("cache:fabrics:"+lang, list, CacheTTLMedium)
		}
		if list, err := cmsSvc.ListPublishedFactories(); err == nil {
			cache.Set("cache:factories:"+lang, list, CacheTTLMedium)
		}
		if list, err := cmsSvc.ListPublishedCertifications(); err == nil {
			cache.Set("cache:certifications:"+lang, list, CacheTTLMedium)
		}
		if list, err := cmsSvc.ListPublishedProductionProcesses(); err == nil {
			cache.Set("cache:production-processes:"+lang, list, CacheTTLMedium)
		}
	}

	// 导航（低频变更，长 TTL）
	if list, err := cmsSvc.ListNavigations("header"); err == nil {
		cache.Set("cache:navigations:header", list, CacheTTLLong)
	}
}

// buildHomeData 组装与 public.GetHome 一致的首页数据并本地化（并发查询）。
func buildHomeData(cmsSvc *CMSService, productSvc *ProductService, lang string) map[string]interface{} {
	page, err := cmsSvc.GetPageBySlug("home")
	if err != nil {
		return nil
	}
	LocalizePage(page, lang)

	var (
		wg             sync.WaitGroup
		products       []models.Product
		categories     []models.Category
		blogs          []models.Blog
		cases          []models.Case
		certifications []models.Certification
		processes      []models.ProductionProcess
		factories      []models.Factory
	)

	wg.Add(7)
	go func() { defer wg.Done(); products, _ = productSvc.ListFeaturedProducts(8) }()
	go func() { defer wg.Done(); categories, _ = productSvc.ListCategories() }()
	go func() { defer wg.Done(); blogs, _, _ = cmsSvc.ListBlogs(1, 4, "", "published") }()
	go func() { defer wg.Done(); cases, _, _ = cmsSvc.ListPublishedCases(1, 4, "") }()
	go func() { defer wg.Done(); certifications, _ = cmsSvc.ListPublishedCertifications() }()
	go func() { defer wg.Done(); processes, _ = cmsSvc.ListPublishedProductionProcesses() }()
	go func() { defer wg.Done(); factories, _ = cmsSvc.ListPublishedFactories() }()
	wg.Wait()

	LocalizeProducts(products, lang)
	LocalizeBlogs(blogs, lang)
	LocalizeCases(cases, lang)

	return map[string]interface{}{
		"page":                 page,
		"featured_products":    products,
		"categories":           categories,
		"blogs":                blogs,
		"cases":                cases,
		"certifications":       certifications,
		"production_processes": processes,
		"factories":            factories,
	}
}