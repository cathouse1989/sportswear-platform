import { ofetch } from 'ofetch'

export const useApi = () => {
  const config = useRuntimeConfig()
  const { locale } = useI18n()
  const route = useRoute()

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

  // 解包 ApiResponse，返回 data 字段
  async function unwrap<T>(url: string, opts?: any): Promise<T> {
    const res = await api<{ success: boolean; data: T }>(url, opts)
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
    trackClick: (target: string, url: string) => unwrap<any>('/public/click-track', { method: 'POST', body: { target, url } }),
  }
}
