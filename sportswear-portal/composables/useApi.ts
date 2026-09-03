import { ofetch } from 'ofetch'

const CONSENT_STORAGE_KEY = 'sw_cookie_consent'
const VISITOR_ID_KEY = 'sw_visitor_id'

// 生成或读取访客唯一标识（存 localStorage，跨会话持久化）
const getVisitorId = (): string => {
  if (import.meta.server) return ''
  let id = localStorage.getItem(VISITOR_ID_KEY)
  if (!id) {
    id = crypto.randomUUID()
    localStorage.setItem(VISITOR_ID_KEY, id)
  }
  return id
}

export const useApi = () => {
  const config = useRuntimeConfig()
  const { locale } = useI18n()
  const route = useRoute()

  // ==================== 流量归因（UTM/来源追踪，用于广告转化分析） ====================
  // 策略：
  // 1. 每次请求附带当前 URL 的 utm_* 参数（SSR 首屏天然携带）；
  // 2. 首次命中的归因信息写入 Cookie（sw_attribution），后续 SPA 内页面跳转持续透传；
  // 3. 后端中间件读取 X-UTM-* 头落库到 visit_logs / leads。
  // 合规：仅在用户同意 ANALYTICS 类别后才写入归因 Cookie
  const ATTRIBUTION_COOKIE = 'sw_attribution'

  /**
   * 检查用户是否已同意分析 Cookie（直接读 localStorage，避免与 useConsent 循环依赖）
   */
  const hasAnalyticsConsent = (): boolean => {
    if (import.meta.server) return true
    try {
      const raw = localStorage.getItem(CONSENT_STORAGE_KEY)
      if (!raw) return false
      const state = JSON.parse(raw)
      return state.categories?.includes('analytics') ?? false
    } catch {
      return false
    }
  }

  const readStoredAttribution = () => {
    if (import.meta.server) return {}
    try {
      const raw = document.cookie.split('; ').find((c) => c.startsWith(ATTRIBUTION_COOKIE + '='))
      if (!raw) return {}
      const val = decodeURIComponent(raw.split('=').slice(1).join('='))
      return JSON.parse(val || '{}') as Record<string, string>
    } catch {
      return {}
    }
  }

  const writeStoredAttribution = (data: Record<string, string>) => {
    if (import.meta.server) return
    // 合规检查：仅在用户同意分析 Cookie 后写入
    if (!hasAnalyticsConsent()) return
    try {
      document.cookie = `${ATTRIBUTION_COOKIE}=${encodeURIComponent(JSON.stringify(data))}; path=/; max-age=31536000; SameSite=Lax`
    } catch {
      /* ignore */
    }
  }

  const currentAttribution = () => {
    const q = route.query
    const pick = (key: string) => {
      const v = q[key]
      if (typeof v === 'string') return v
      if (Array.isArray(v)) return v[0] || ''
      return ''
    }
    const fresh: Record<string, string> = {
      source: pick('utm_source'),
      medium: pick('utm_medium'),
      campaign: pick('utm_campaign'),
      content: pick('utm_content'),
      term: pick('utm_term'),
      ref: pick('ref'),
    }
    const stored = readStoredAttribution()
    const merged: Record<string, string> = { ...stored }
    let hasFresh = false
    for (const k of Object.keys(fresh)) {
      if (fresh[k]) {
        merged[k] = fresh[k]
        hasFresh = true
      }
    }
    // 首次访问无 UTM 时，用 document.referrer 记录来源网站（便于"网站推送来源"归因）
    if (import.meta.client && !merged.ref && !fresh.ref && document.referrer) {
      merged.ref = document.referrer
    }
    if (hasFresh && import.meta.client) writeStoredAttribution(merged)
    return merged
  }

  const utmHeaders = (attr: Record<string, string>) => {
    const out: Record<string, string> = {}
    const mapping: Array<[string, string]> = [
      ['source', 'X-UTM-Source'],
      ['medium', 'X-UTM-Medium'],
      ['campaign', 'X-UTM-Campaign'],
      ['content', 'X-UTM-Content'],
      ['term', 'X-UTM-Term'],
    ]
    for (const [key, header] of mapping) {
      if (attr[key]) out[header] = attr[key]
    }
    // 附带同意状态头，供后端判断是否记录个人数据
    out['X-Consent-Analytics'] = hasAnalyticsConsent() ? 'true' : 'false'
    // 附带访客唯一标识，用于关联访问记录与询盘
    const visitorId = getVisitorId()
    if (visitorId) out['X-Visitor-ID'] = visitorId
    return out
  }

  // 管理后台 iframe 预览会带 ?preview=1；门户请求公开 API 时透传，强制后端绕过缓存
  const isPreview = () => {
    const q = route.query?.preview
    return q === '1' || q === 'true' || q === 'yes'
  }

  const api = ofetch.create({
    // SSR（服务端渲染）在容器网络内访问后端；浏览器端走公共 API 地址
    baseURL: import.meta.server ? config.apiServer || config.public.apiBase : config.public.apiBase,
    params: {
      lang: locale.value,
    },
    onResponseError({ response }) {
      console.error(`API Error: ${response.status}`, response._data)
    },
  })

  const getLang = () => locale.value || 'en'

  const withPreview = (params: Record<string, any> = {}) => {
    const out = { ...params }
    if (isPreview()) out.preview = 1
    return out
  }

  // 解包 ApiResponse，返回 data 字段；自动附带 UTM 归因请求头
  async function unwrap<T>(url: string, opts?: any): Promise<T> {
    const res = await api<{ success: boolean; data: T }>(url, {
      ...opts,
      headers: { ...(opts?.headers || {}), ...utmHeaders(currentAttribution()) },
    })
    return res.data
  }

  return {
    getHome: () => unwrap<any>('/public/home', { params: withPreview({ lang: getLang() }) }),
    getProducts: (params?: any) => unwrap<any>('/public/products', { params: withPreview({ lang: getLang(), ...params }) }),
    getProduct: (slug: string) => unwrap<any>(`/public/products/${slug}`, { params: withPreview({ lang: getLang() }) }),
    getPage: (slug: string) => unwrap<any>(`/public/pages/${slug}`, { params: withPreview({ lang: getLang() }) }),
    getBlogs: (params?: any) => unwrap<any>('/public/blogs', { params: withPreview({ lang: getLang(), ...params }) }),
    getBlog: (slug: string) => unwrap<any>(`/public/blogs/${slug}`, { params: withPreview({ lang: getLang() }) }),
    getCases: (params?: any) => unwrap<any>('/public/cases', { params: withPreview({ lang: getLang(), ...params }) }),
    getFaqs: (params?: any) => unwrap<any>('/public/faqs', { params: withPreview({ lang: getLang(), ...params }) }),
    getCategories: () => unwrap<any>('/public/categories', { params: withPreview() }),
    getSeries: () => unwrap<any>('/public/series', { params: withPreview() }),
    getFabrics: () => unwrap<any>('/public/fabrics', { params: withPreview() }),
    getFactories: () => unwrap<any>('/public/factories', { params: withPreview() }),
    getCertifications: () => unwrap<any>('/public/certifications', { params: withPreview() }),
    getNavigations: (type: string = 'header') => unwrap<any>('/public/navigations', { params: withPreview({ type }) }),
    getTheme: () => unwrap<any>('/public/theme', { params: withPreview() }),
    getI18n: () => unwrap<any>('/public/i18n', { params: withPreview({ lang: getLang() }) }),
    getCurrencies: () => unwrap<any>('/public/currencies', { params: withPreview() }),
    getLanguages: () => unwrap<any>('/public/languages', { params: withPreview() }),
    getGeo: (country?: string) => unwrap<any>('/public/geo', { params: withPreview({ country }) }),
    submitLead: (data: any) => unwrap<any>('/public/leads', { method: 'POST', body: data }),
    uploadLeadAttachment: (file: File) => {
      const fd = new FormData()
      fd.append('file', file)
      return unwrap<any>('/public/uploads/lead-attachment', { method: 'POST', body: fd })
    },
    trackClick: (target: string, url: string) => unwrap<any>('/public/click-track', { method: 'POST', body: { target, url } }),
  }
}
