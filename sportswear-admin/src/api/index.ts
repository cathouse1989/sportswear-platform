import { http } from './client'
import type {
  DashboardStats,
  Lead,
  LeadFollowUp,
  Subscriber,
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
  SelfMedia,
  Navigation,
  Media,
  StorageSource,
  Currency,
  I18nEntry,
  Category,
  Series,
  Fabric,
} from '@/types'


// ============ 璁よ瘉 ============
export interface LoginResult {
  token: string
  user: User
  permissions: string[]
}

export const authApi = {
  login: (email: string, password: string) =>
    http.post<LoginResult>('/admin/auth/login', { email, password }),
  profile: () => http.get<User>('/admin/auth/profile'),
  // 登出：通知后端清除 HttpOnly 会话 Cookie（本地 token 由 store 清理）
  logout: () => http.post<{ logged_out: boolean }>('/admin/auth/logout'),
}

// ============ 浠〃鐩?============
export const dashboardApi = {
  stats: () => http.get<DashboardStats>('/admin/dashboard'),
}

// ============ 鐢ㄦ埛涓庤鑹?============
export const userApi = {
  list: (params?: any) => http.get<PageResult<User>>('/admin/users', params),
  create: (data: any) => http.post<User>('/admin/users', data),
  update: (id: string, data: any) => http.put<User>(`/admin/users/${id}`, data),
  delete: (id: string) => http.delete(`/admin/users/${id}`),
}

export const roleApi = {
  list: () => http.get<Role[]>('/admin/roles'),
  listPage: (params?: any) => http.get<PageResult<Role>>('/admin/roles/page', params),
  create: (data: any) => http.post<Role>('/admin/roles', data),
  update: (id: string, data: any) => http.put<Role>(`/admin/roles/${id}`, data),
  updateStatus: (id: string, isActive: boolean) =>
    http.put(`/admin/roles/${id}/status`, { is_active: isActive }),
  delete: (id: string) => http.delete(`/admin/roles/${id}`),
  permissions: () => http.get<Permission[]>('/admin/permissions'),
}

// ============ 浜у搧 ============
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

// ============ 鍒嗙被 ============
export const categoryApi = {
  list: () => http.get<Category[]>('/admin/categories'),
  create: (data: any) => http.post<Category>('/admin/categories', data),
  update: (id: string, data: any) => http.put<Category>(`/admin/categories/${id}`, data),
  delete: (id: string) => http.delete(`/admin/categories/${id}`),
}

// ============ 绯诲垪 ============
export const seriesApi = {
  list: () => http.get<Series[]>('/admin/series'),
  create: (data: any) => http.post<Series>('/admin/series', data),
  publish: (id: string) => http.post(`/admin/series/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/series/${id}/unpublish`),
}

// ============ 闈㈡枡 ============
export const fabricApi = {
  list: () => http.get<Fabric[]>('/admin/fabrics'),
  create: (data: any) => http.post<Fabric>('/admin/fabrics', data),
  publish: (id: string) => http.post(`/admin/fabrics/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/fabrics/${id}/unpublish`),
}

// ============ 椤甸潰 ============
export const pageApi = {
  list: (params?: any) => http.get<PageResult<Page>>('/admin/pages', params),
  get: (id: string) => http.get<Page>(`/admin/pages/${id}`),
  create: (data: any) => http.post<Page>('/admin/pages', data),
  update: (id: string, data: any) => http.put<Page>(`/admin/pages/${id}`, data),
  delete: (id: string) => http.delete(`/admin/pages/${id}`),
  publish: (id: string) => http.post(`/admin/pages/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/pages/${id}/unpublish`),
  // 杞挱鍥撅紙棣栭〉 hero banner锛夐厤缃細slides + 鍙€夊睍绀哄弬鏁?settings
  updateHeroSlides: (id: string, slides: HeroSlide[], settings?: HeroSettings) =>
    http.put<any>(`/admin/pages/${id}/hero`, { slides, settings: settings || {} }),
}

// ============ 瀵艰埅 ============
export const navigationApi = {
  list: (params?: any) => http.get<Navigation[]>('/admin/navigations', params),
  listAll: () => http.get<Navigation[]>('/admin/navigations', { type: '' }),
  create: (data: any) => http.post<Navigation>('/admin/navigations', data),
  update: (id: string, data: any) => http.put<Navigation>(`/admin/navigations/${id}`, data),
  delete: (id: string) => http.delete(`/admin/navigations/${id}`),
  batchSort: (items: { id: string; sort_order: number }[]) => http.put('/admin/navigations/sort', { items }),
  syncWithPages: () => http.post<{ updated: number }>('/admin/navigations/sync-with-pages'),
}

// ============ 鍗氬 ============
export const blogApi = {
  list: (params?: any) => http.get<PageResult<Blog>>('/admin/blogs', params),
  get: (id: string) => http.get<Blog>(`/admin/blogs/${id}`),
  create: (data: any) => http.post<Blog>('/admin/blogs', data),
  update: (id: string, data: any) => http.put<Blog>(`/admin/blogs/${id}`, data),
  delete: (id: string) => http.delete(`/admin/blogs/${id}`),
  publish: (id: string) => http.post(`/admin/blogs/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/blogs/${id}/unpublish`),
}

// ============ 妗堜緥 ============
export const caseApi = {
  list: (params?: any) => http.get<PageResult<Case>>('/admin/cases', params),
  get: (id: string) => http.get<Case>(`/admin/cases/${id}`),
  create: (data: any) => http.post<Case>('/admin/cases', data),
  update: (id: string, data: any) => http.put<Case>(`/admin/cases/${id}`, data),
  delete: (id: string) => http.delete(`/admin/cases/${id}`),
  publish: (id: string) => http.post(`/admin/cases/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/cases/${id}/unpublish`),
}

// ============ FAQ ============
export const faqApi = {
  list: (params?: any) => http.get<PageResult<FAQ>>('/admin/faqs', params),
  create: (data: any) => http.post<FAQ>('/admin/faqs', data),
  update: (id: string, data: any) => http.put<FAQ>(`/admin/faqs/${id}`, data),
  delete: (id: string) => http.delete(`/admin/faqs/${id}`),
}

// ============ 宸ュ巶 ============
export const factoryApi = {
  list: (params?: any) => http.get<PageResult<Factory>>('/admin/factories', params),
  create: (data: any) => http.post<Factory>('/admin/factories', data),
  update: (id: string, data: any) => http.put<Factory>(`/admin/factories/${id}`, data),
  delete: (id: string) => http.delete(`/admin/factories/${id}`),
  publish: (id: string) => http.post(`/admin/factories/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/factories/${id}/unpublish`),
}

// ============ 璁よ瘉 ============
export const certificationApi = {
  list: (params?: any) => http.get<PageResult<Certification>>('/admin/certifications', params),
  create: (data: any) => http.post<Certification>('/admin/certifications', data),
  update: (id: string, data: any) => http.put<Certification>(`/admin/certifications/${id}`, data),
  delete: (id: string) => http.delete(`/admin/certifications/${id}`),
  publish: (id: string) => http.post(`/admin/certifications/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/certifications/${id}/unpublish`),
}

export const productionProcessApi = {
  list: (params?: any) => http.get<PageResult<ProductionProcess>>('/admin/production-processes', params),
  create: (data: any) => http.post<ProductionProcess>('/admin/production-processes', data),
  update: (id: string, data: any) => http.put<ProductionProcess>(`/admin/production-processes/${id}`, data),
  delete: (id: string) => http.delete(`/admin/production-processes/${id}`),
  publish: (id: string) => http.post(`/admin/production-processes/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/production-processes/${id}/unpublish`),
}

// ============ 濯掍綋 ============
// ============ 自媒体 ============
export const selfMediaApi = {
  list: (params?: any) => http.get<PageResult<SelfMedia>>('/admin/self-medias', params),
  create: (data: any) => http.post<SelfMedia>('/admin/self-medias', data),
  update: (id: string, data: any) => http.put<SelfMedia>(`/admin/self-medias/${id}`, data),
  delete: (id: string) => http.delete(`/admin/self-medias/${id}`),
  publish: (id: string) => http.post(`/admin/self-medias/${id}/publish`),
  unpublish: (id: string) => http.post(`/admin/self-medias/${id}/unpublish`),
  getConfig: () => http.get<{ max_display: number }>('/admin/self-medias/config'),
  updateConfig: (maxDisplay: number) => http.put<{ max_display: number }>('/admin/self-medias/config', { max_display: maxDisplay }),
}

export const mediaApi = {
  list: (params?: any) => http.get<PageResult<Media>>('/admin/media', params),
  get: (id: string) => http.get<Media>(`/admin/media/${id}`),
  create: (data: any) => http.post<Media>('/admin/media', data),
  upload: (data: FormData) => http.post<Media>('/admin/media/upload', data),
  update: (id: string, data: any) => http.put<Media>(`/admin/media/${id}`, data),
  delete: (id: string) => http.delete(`/admin/media/${id}`),
}

// ============ 璇㈢洏 ============
export const leadApi = {
  list: (params?: any) => http.get<PageResult<Lead>>('/admin/leads', params),
  get: (id: string) => http.get<Lead>(`/admin/leads/${id}`),
  update: (id: string, data: any) => http.put<Lead>(`/admin/leads/${id}`, data),
  delete: (id: string) => http.delete(`/admin/leads/${id}`),
  addFollowUp: (id: string, data: any) =>
    http.post<LeadFollowUp>(`/admin/leads/${id}/followups`, data),
}

// ============ 订阅 ============
export const subscriberApi = {
  stats: () =>
    http.get<{ total: number; subscribed: number; unsubscribed: number }>('/admin/subscribers/stats'),
  list: (params?: any) => http.get<PageResult<Subscriber>>('/admin/subscribers', params),
  update: (id: string, data: any) => http.put<Subscriber>(`/admin/subscribers/${id}`, data),
  delete: (id: string) => http.delete(`/admin/subscribers/${id}`),
}

// ============ 鎶ヤ环 ============
export const quoteApi = {
  list: (params?: any) => http.get<PageResult<Quote>>('/admin/quotes', params),
  get: (id: string) => http.get<Quote>(`/admin/quotes/${id}`),
  create: (data: any) => http.post<Quote>('/admin/quotes', data),
  update: (id: string, data: any) => http.put<Quote>(`/admin/quotes/${id}`, data),
  delete: (id: string) => http.delete(`/admin/quotes/${id}`),
}

// ============ 閫氱煡 ============
export const notificationApi = {
  list: (params?: any) => http.get<PageResult<Notification>>('/admin/notifications', params),
  unreadCount: () => http.get<{ count: number }>('/admin/notifications/unread-count'),
  markRead: (id: string) => http.post(`/admin/notifications/${id}/read`),
  markAllRead: () => http.post('/admin/notifications/read-all'),
}

// ============ 鎿嶄綔鏃ュ織 ============
export const operationLogApi = {
  list: (params?: any) =>
    http.get<PageResult<OperationLog>>('/admin/operation-logs', params),
}

// ============ 鏁版嵁鍒嗘瀽 ============
export const analyticsApi = {
  overview: (params?: any) => http.get<any>('/admin/analytics/overview', params),
  topPages: (params?: any) => http.get<any>('/admin/analytics/top-pages', params),
  topProducts: (params?: any) => http.get<any>('/admin/analytics/top-products', params),
  sources: (params?: any) => http.get<any>('/admin/analytics/sources', params),
  countries: (params?: any) => http.get<any>('/admin/analytics/countries', params),
  devices: (params?: any) => http.get<any>('/admin/analytics/devices', params),
  socialClicks: (params?: any) => http.get<any[]>('/admin/analytics/social-clicks', params),
  utmCampaigns: (params?: any) => http.get<any[]>('/admin/analytics/utm-campaigns', params),
  // 门户访问日志明细（IP/国家/来源/实体/时间范围筛选，跨月度分表组合查询）
  visitLogs: (params?: any) => http.get<PageResult<any>>('/admin/analytics/visit-logs', params),
  // 访客旅程（某个访客的完整访问路径）
  visitorJourney: (params: { visitor_id: string; days?: number }) =>
    http.get<any[]>('/admin/analytics/visitor-journey', params),
  // IP-询盘关联（某个 IP 提交的询盘列表）
  ipLeads: (params: { ip: string; days?: number }) =>
    http.get<any[]>('/admin/analytics/ip-leads', params),
  // 转化漏斗（访问 → 产品浏览 → 询盘）
  conversionFunnel: (params?: { days?: number }) =>
    http.get<any>('/admin/analytics/conversion-funnel', params),
  // 隐私合规洞察（同意/未同意用户的转化对比）
  consentInsights: (params?: { days?: number }) =>
    http.get<any>('/admin/analytics/consent-insights', params),
}

// ============ 鍥介檯鍖?============
export const i18nApi = {
  entries: (params?: any) => http.get<PageResult<I18nEntry>>('/admin/i18n/entries', params),
  upsert: (data: any) => http.post<I18nEntry>('/admin/i18n/entries', data),
  delete: (id: string) => http.delete(`/admin/i18n/entries/${id}`),
}

// ============ 涓婚閰嶇疆 ============
export const themeApi = {
  get: () => http.get<any>('/public/theme'),
  adminList: () => http.get<any[]>('/admin/theme'),
  update: (key: string, value: any) => http.put<any>(`/admin/theme/${key}`, { value }),
}

// ============ 闂ㄦ埛鍏紑鎺ュ彛锛堝悗鍙伴瑙堢敤锛?============
// 娉ㄦ剰锛氶瑙堣鍙栫粺涓€甯?preview=1锛屽悗绔皢寮哄埗缁曡繃 Redis 缂撳瓨鐩存煡 DB锛?
// 淇濊瘉"闂ㄦ埛棰勮"濮嬬粓鍙嶆槧鏁版嵁搴撴渶鏂扮姸鎬侊紙鍙戝竷/涓嬬嚎绔嬪嵆鍙锛夛紝涓斾笉姹℃煋绾夸笂缂撳瓨銆?
export const publicApi = {
  i18n: (lang?: string) => http.get<any>('/public/i18n', lang ? { lang, preview: 1 } : { preview: 1 }),
  languages: () => http.get<any[]>('/public/languages'),
  theme: () => http.get<any>('/public/theme'),
  home: (lang?: string) => http.get<any>('/public/home', lang ? { lang, preview: 1 } : { preview: 1 }),
  products: (params?: any) => http.get<any>('/public/products', { ...params, preview: 1 }),
  product: (slug: string, lang?: string) =>
    http.get<any>(`/public/products/${slug}`, lang ? { lang, preview: 1 } : { preview: 1 }),
  // 棰勮椤甸€氱敤璇诲彇锛坧ath 闇€鍚?/public 鍓嶇紑锛屽 /public/blogs銆?public/home锛?
  fetch: (path: string, lang?: string) =>
    http.get<any>(path, lang ? { lang, preview: 1 } : { preview: 1 }),
  trackClick: (target: string, url: string) =>
    http.post('/public/click-track', { target, url }),
}

// ============ 瀛樺偍婧?============
export const storageSourceApi = {
  list: () => http.get<StorageSource[]>('/admin/storage-sources'),
  get: (id: string) => http.get<StorageSource>(`/admin/storage-sources/${id}`),
  create: (data: any) => http.post<StorageSource>('/admin/storage-sources', data),
  update: (id: string, data: any) => http.put<StorageSource>(`/admin/storage-sources/${id}`, data),
  delete: (id: string) => http.delete(`/admin/storage-sources/${id}`),
}

// ============ 鍥炴敹绔?============
// 娉ㄦ剰锛氬悗绔?RestoreItem/PurgeItem 鍧囪姹?body 鎼哄甫 entity_type锛坆inding:"required"锛夛紝缂哄け浼?400
export const trashApi = {
  list: (params?: any) => http.get<PageResult<any>>('/admin/trash', params),
  restore: (id: string, entityType: string) => http.post(`/admin/trash/${id}/restore`, { entity_type: entityType }),
  purge: (id: string, entityType: string) => http.post(`/admin/trash/${id}/purge`, { entity_type: entityType }),
  empty: () => http.delete('/admin/trash'),
}

// ============ 闂ㄦ埛缂撳瓨 ============
export const portalCacheApi = {
  status: () => http.get<{ redis_ok: boolean; enabled: boolean; using_cache: boolean; last_refresh_at?: string | null }>('/admin/portal-cache'),
  setEnabled: (enabled: boolean) => http.put<any>('/admin/portal-cache/enabled', { enabled }),
  refresh: () => http.post<any>('/admin/portal-cache/refresh'),
  publish: () => http.post<any>('/admin/portal-cache/publish'),
}

// ============ 闂ㄦ埛鐗堟湰 ============
export const portalApi = {
  saveDraft: (pageId: string, data: any) => http.post<any>(`/admin/pages/${pageId}/versions`, data),
  listVersions: (pageId: string) => http.get<any[]>(`/admin/pages/${pageId}/versions`),
  publishVersion: (versionId: string) => http.post(`/admin/versions/${versionId}/publish`),
  rollbackVersion: (versionId: string) => http.post(`/admin/versions/${versionId}/rollback`),
}

// ============ SEO ============
export const seoApi = {
  upsert: (data: any) => http.post<any>('/admin/seo', data),
  listRoutes: (route: string) => http.get<any>('/admin/seo/routes', { route }),
  saveRoutes: (route: string, entries: any[]) => http.post<any>('/admin/seo/routes', { route, entries }),
}

// ============ 璐у竵 ============
export const currencyApi = {
  list: () => http.get<Currency[]>('/public/currencies'),
  convert: (params: any) => http.get<any>('/public/currencies/convert', params),
}

// ============ 国家本地化映射 ============
export const geoLocaleApi = {
  list: () => http.get<any[]>('/admin/geo-locales'),
  upsert: (data: any) => http.post<any>('/admin/geo-locales', data),
  delete: (id: string) => http.delete(`/admin/geo-locales/${id}`),
}
// ============ AI 翻译 ============
export const aiTranslateApi = {
  translate: (texts: string[], targetLang: string) => http.post<any>('/admin/ai/translate', { texts, target_lang: targetLang }),
}