import { http } from './client'
import type {
  DashboardStats,
  Lead,
  LeadFollowUp,
  Notification,
  OperationLog,
  PageResult,
  Permission,
  Product,
  Quote,
  Role,
  User,
  Page,
  HeroSlide,
  HeroSettings,
  Blog,
  Case,
  FAQ,
  Factory,
  Certification,
  ProductionProcess,
  Navigation,
  Media,
  StorageSource,
  Currency,
  I18nEntry,
  Category,
  Series,
  Fabric,
} from '@/types'


// ============ 认证 ============
export interface LoginResult {
  token: string
  user: User
  permissions: string[]
}

export const authApi = {
  login: (email: string, password: string) =>
    http.post<LoginResult>('/admin/auth/login', { email, password }),
  profile: () => http.get<User>('/admin/auth/profile'),
}

// ============ 仪表盘 ============
export const dashboardApi = {
  stats: () => http.get<DashboardStats>('/admin/dashboard'),
}

// ============ 用户与角色 ============
export const userApi = {
  list: (params?: any) => http.get<PageResult<User>>('/admin/users', params),
  create: (data: any) => http.post<User>('/admin/users', data),
  update: (id: string, data: any) => http.put<User>(`/admin/users/${id}`, data),
  delete: (id: string) => http.delete(`/admin/users/${id}`),
}

export const roleApi = {
  list: () => http.get<Role[]>('/admin/roles'),
  create: (data: any) => http.post<Role>('/admin/roles', data),
  update: (id: string, data: any) => http.put<Role>(`/admin/roles/${id}`, data),
  updateStatus: (id: string, isActive: boolean) =>
    http.put(`/admin/roles/${id}/status`, { is_active: isActive }),
  permissions: () => http.get<Permission[]>('/admin/permissions'),
}

// ============ 产品 ============
export const productApi = {
  list: (params?: any) => http.get<PageResult<Product>>('/admin/products', params),
  stats: () =>
    http.get<{ total: number; published: number; draft: number; offline: number; featured: number }>(
      '/admin/products/stats'
    ),
  get: (id: string) => http.get<Product>(`/admin/products/${id}`),
  create: (data: any) => http.post<Product>('/admin/products', data),
  update: (id: string, data: any) => http.put<Product>(`/admin/products/${id}`, data),
  delete: (id: string) => http.delete(`/admin/products/${id}`),
  publish: (id: string) => http.post(`/admin/products/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/products/${id}/unpublish`),
}

// ============ 分类 ============
export const categoryApi = {
  list: () => http.get<Category[]>('/admin/categories'),
  create: (data: any) => http.post<Category>('/admin/categories', data),
  update: (id: string, data: any) => http.put<Category>(`/admin/categories/${id}`, data),
  delete: (id: string) => http.delete(`/admin/categories/${id}`),
}

// ============ 系列 ============
export const seriesApi = {
  list: () => http.get<Series[]>('/admin/series'),
  create: (data: any) => http.post<Series>('/admin/series', data),
  publish: (id: string) => http.post(`/admin/series/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/series/${id}/unpublish`),
}

// ============ 面料 ============
export const fabricApi = {
  list: () => http.get<Fabric[]>('/admin/fabrics'),
  create: (data: any) => http.post<Fabric>('/admin/fabrics', data),
  publish: (id: string) => http.post(`/admin/fabrics/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/fabrics/${id}/unpublish`),
}

// ============ 页面 ============
export const pageApi = {
  list: (params?: any) => http.get<PageResult<Page>>('/admin/pages', params),
  get: (id: string) => http.get<Page>(`/admin/pages/${id}`),
  create: (data: any) => http.post<Page>('/admin/pages', data),
  update: (id: string, data: any) => http.put<Page>(`/admin/pages/${id}`, data),
  delete: (id: string) => http.delete(`/admin/pages/${id}`),
  publish: (id: string) => http.post(`/admin/pages/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/pages/${id}/unpublish`),
  // 轮播图（首页 hero banner）配置：slides + 可选展示参数 settings
  updateHeroSlides: (id: string, slides: HeroSlide[], settings?: HeroSettings) =>
    http.put<any>(`/admin/pages/${id}/hero`, settings ? { slides, settings } : { slides }),
}

// ============ 导航 ============
export const navigationApi = {
  list: (params?: any) => http.get<Navigation[]>('/admin/navigations', params),
  create: (data: any) => http.post<Navigation>('/admin/navigations', data),
  update: (id: string, data: any) => http.put<Navigation>(`/admin/navigations/${id}`, data),
  delete: (id: string) => http.delete(`/admin/navigations/${id}`),
}

// ============ 博客 ============
export const blogApi = {
  list: (params?: any) => http.get<PageResult<Blog>>('/admin/blogs', params),
  get: (id: string) => http.get<Blog>(`/admin/blogs/${id}`),
  create: (data: any) => http.post<Blog>('/admin/blogs', data),
  update: (id: string, data: any) => http.put<Blog>(`/admin/blogs/${id}`, data),
  delete: (id: string) => http.delete(`/admin/blogs/${id}`),
  publish: (id: string) => http.post(`/admin/blogs/${id}/publish`),
}

// ============ 案例 ============
export const caseApi = {
  list: (params?: any) => http.get<PageResult<Case>>('/admin/cases', params),
  get: (id: string) => http.get<Case>(`/admin/cases/${id}`),
  create: (data: any) => http.post<Case>('/admin/cases', data),
  update: (id: string, data: any) => http.put<Case>(`/admin/cases/${id}`, data),
  delete: (id: string) => http.delete(`/admin/cases/${id}`),
}

// ============ FAQ ============
export const faqApi = {
  list: (params?: any) => http.get<PageResult<FAQ>>('/admin/faqs', params),
  create: (data: any) => http.post<FAQ>('/admin/faqs', data),
  update: (id: string, data: any) => http.put<FAQ>(`/admin/faqs/${id}`, data),
  delete: (id: string) => http.delete(`/admin/faqs/${id}`),
}

// ============ 工厂 ============
export const factoryApi = {
  list: (params?: any) => http.get<Factory[]>('/admin/factories', params),
  create: (data: any) => http.post<Factory>('/admin/factories', data),
  update: (id: string, data: any) => http.put<Factory>(`/admin/factories/${id}`, data),
  delete: (id: string) => http.delete(`/admin/factories/${id}`),
  publish: (id: string) => http.post(`/admin/factories/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/factories/${id}/unpublish`),
}

// ============ 认证 ============
export const certificationApi = {
  list: (params?: any) => http.get<Certification[]>('/admin/certifications', params),
  create: (data: any) => http.post<Certification>('/admin/certifications', data),
  update: (id: string, data: any) => http.put<Certification>(`/admin/certifications/${id}`, data),
  delete: (id: string) => http.delete(`/admin/certifications/${id}`),
  publish: (id: string) => http.post(`/admin/certifications/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/certifications/${id}/unpublish`),
}

// ============ 生产流程 ============
export const productionProcessApi = {
  list: (params?: any) => http.get<ProductionProcess[]>('/admin/production-processes', params),
  create: (data: any) => http.post<ProductionProcess>('/admin/production-processes', data),
  update: (id: string, data: any) => http.put<ProductionProcess>(`/admin/production-processes/${id}`, data),
  delete: (id: string) => http.delete(`/admin/production-processes/${id}`),
  publish: (id: string) => http.post(`/admin/production-processes/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/production-processes/${id}/unpublish`),
}

// ============ 媒体 ============
export const mediaApi = {
  list: (params?: any) => http.get<PageResult<Media>>('/admin/media', params),
  get: (id: string) => http.get<Media>(`/admin/media/${id}`),
  create: (data: any) => http.post<Media>('/admin/media', data),
  upload: (data: FormData) => http.post<Media>('/admin/media/upload', data),
  update: (id: string, data: any) => http.put<Media>(`/admin/media/${id}`, data),
  delete: (id: string) => http.delete(`/admin/media/${id}`),
}

// ============ 询盘 ============
export const leadApi = {
  list: (params?: any) => http.get<PageResult<Lead>>('/admin/leads', params),
  get: (id: string) => http.get<Lead>(`/admin/leads/${id}`),
  update: (id: string, data: any) => http.put<Lead>(`/admin/leads/${id}`, data),
  delete: (id: string) => http.delete(`/admin/leads/${id}`),
  addFollowUp: (id: string, data: any) =>
    http.post<LeadFollowUp>(`/admin/leads/${id}/followups`, data),
}

// ============ 报价 ============
export const quoteApi = {
  list: (params?: any) => http.get<PageResult<Quote>>('/admin/quotes', params),
  get: (id: string) => http.get<Quote>(`/admin/quotes/${id}`),
  create: (data: any) => http.post<Quote>('/admin/quotes', data),
  update: (id: string, data: any) => http.put<Quote>(`/admin/quotes/${id}`, data),
  delete: (id: string) => http.delete(`/admin/quotes/${id}`),
}

// ============ 通知 ============
export const notificationApi = {
  list: (params?: any) => http.get<PageResult<Notification>>('/admin/notifications', params),
  unreadCount: () => http.get<{ count: number }>('/admin/notifications/unread-count'),
  markRead: (id: string) => http.post(`/admin/notifications/${id}/read`),
  markAllRead: () => http.post('/admin/notifications/read-all'),
}

// ============ 操作日志 ============
export const operationLogApi = {
  list: (params?: any) =>
    http.get<PageResult<OperationLog>>('/admin/operation-logs', params),
}

// ============ 数据分析 ============
export const analyticsApi = {
  overview: (params?: any) => http.get<any>('/admin/analytics/overview', params),
  topPages: (params?: any) => http.get<any>('/admin/analytics/top-pages', params),
  topProducts: (params?: any) => http.get<any>('/admin/analytics/top-products', params),
  sources: (params?: any) => http.get<any>('/admin/analytics/sources', params),
  countries: (params?: any) => http.get<any>('/admin/analytics/countries', params),
  devices: (params?: any) => http.get<any>('/admin/analytics/devices', params),
}

// ============ 国际化 ============
export const i18nApi = {
  entries: (params?: any) => http.get<PageResult<I18nEntry>>('/admin/i18n/entries', params),
  upsert: (data: any) => http.post<I18nEntry>('/admin/i18n/entries', data),
  delete: (id: string) => http.delete(`/admin/i18n/entries/${id}`),
}

// ============ 主题配置 ============
export const themeApi = {
  get: () => http.get<any>('/public/theme'),
  adminList: () => http.get<any[]>('/admin/theme'),
  update: (key: string, value: any) => http.put<any>(`/admin/theme/${key}`, { value }),
}

// ============ 门户公开接口（后台预览用） ============
// 注意：预览读取统一带 preview=1，后端将强制绕过 Redis 缓存直查 DB，
// 保证"门户预览"始终反映数据库最新状态（发布/下线立即可见），且不污染线上缓存。
export const publicApi = {
  i18n: (lang?: string) => http.get<any>('/public/i18n', lang ? { lang, preview: 1 } : { preview: 1 }),
  languages: () => http.get<any[]>('/public/languages'),
  theme: () => http.get<any>('/public/theme'),
  home: (lang?: string) => http.get<any>('/public/home', lang ? { lang, preview: 1 } : { preview: 1 }),
  products: (params?: any) => http.get<any>('/public/products', { ...params, preview: 1 }),
  product: (slug: string, lang?: string) =>
    http.get<any>(`/public/products/${slug}`, lang ? { lang, preview: 1 } : { preview: 1 }),
  // 预览页通用读取（path 需含 /public 前缀，如 /public/blogs、/public/home）
  fetch: (path: string, lang?: string) =>
    http.get<any>(path, lang ? { lang, preview: 1 } : { preview: 1 }),
  trackClick: (target: string, url: string) =>
    http.post('/public/click-track', { target, url }),
}

// ============ 存储源 ============
export const storageSourceApi = {
  list: () => http.get<StorageSource[]>('/admin/storage-sources'),
  get: (id: string) => http.get<StorageSource>(`/admin/storage-sources/${id}`),
  create: (data: any) => http.post<StorageSource>('/admin/storage-sources', data),
  update: (id: string, data: any) => http.put<StorageSource>(`/admin/storage-sources/${id}`, data),
  delete: (id: string) => http.delete(`/admin/storage-sources/${id}`),
}

// ============ 回收站 ============
// 注意：后端 RestoreItem/PurgeItem 均要求 body 携带 entity_type（binding:"required"），缺失会 400
export const trashApi = {
  list: (params?: any) => http.get<PageResult<any>>('/admin/trash', params),
  restore: (id: string, entityType: string) => http.post(`/admin/trash/${id}/restore`, { entity_type: entityType }),
  purge: (id: string, entityType: string) => http.post(`/admin/trash/${id}/purge`, { entity_type: entityType }),
  empty: () => http.delete('/admin/trash'),
}

// ============ 门户版本 ============
export const portalApi = {
  saveDraft: (pageId: string, data: any) => http.post<any>(`/admin/pages/${pageId}/versions`, data),
  listVersions: (pageId: string) => http.get<any[]>(`/admin/pages/${pageId}/versions`),
  publishVersion: (versionId: string) => http.post(`/admin/versions/${versionId}/publish`),
  rollbackVersion: (versionId: string) => http.post(`/admin/versions/${versionId}/rollback`),
}

// ============ SEO ============
export const seoApi = {
  upsert: (data: any) => http.post<any>('/admin/seo', data),
}

// ============ 货币 ============
export const currencyApi = {
  list: () => http.get<Currency[]>('/public/currencies'),
  convert: (params: any) => http.get<any>('/public/currencies/convert', params),
}