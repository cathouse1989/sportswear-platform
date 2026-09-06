// ====================================================================
// 门户「路由」单一事实来源（前端侧）
//
// 与后端 services/portal_health.go 的 PagePath / RouteKey / presetRoutes
// 保持一致，用于 SEO 管理下拉兜底、导航 URL 派生。
// 运行时优先使用 GET /admin/portal/routes 动态数据，本文件作为兜底与类型契约。
// ====================================================================

/** 页面类型 + slug → 门户访问路径 */
export function pagePath(type: string, slug: string): string {
  switch (type) {
    case 'home':
      return '/'
    case 'product':
    case 'product_category':
      return '/products'
    case 'blog':
      return '/blog'
    case 'case':
      return '/cases'
    case 'faq':
      return '/faq'
    case 'contact':
      return '/contact'
    default: {
      const s = (slug || '').trim()
      if (!s) return '/'
      return '/' + s.replace(/^\/+/, '')
    }
  }
}

/** 页面类型 + slug → 路由 key（route，SEO 管理维度） */
export function routeFromPage(type: string, slug: string): string {
  switch (type) {
    case 'home':
      return 'home'
    case 'product':
    case 'product_category':
      return 'products'
    case 'blog':
      return 'blog'
    case 'case':
      return 'cases'
    case 'faq':
      return 'faq'
    case 'contact':
      return 'contact'
    default:
      return slug || ''
  }
}

/** 系统预置路由（SEO 管理下拉兜底，接口不可用时使用） */
export const PRESET_ROUTES = [
  { route: 'home', path: '/', label: '首页' },
  { route: 'products', path: '/products', label: '产品中心' },
  { route: 'cases', path: '/cases', label: '案例展示' },
  { route: 'about', path: '/about', label: '关于我们' },
  { route: 'blog', path: '/blog', label: '博客' },
  { route: 'faq', path: '/faq', label: '常见问题' },
  { route: 'contact', path: '/contact', label: '联系我们' },
  { route: 'privacy', path: '/privacy-policy', label: '隐私政策' },
]
