import { ofetch } from 'ofetch'

export const useApi = () => {
  const config = useRuntimeConfig()
  const { locale } = useI18n()

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

  // 解包 ApiResponse，返回 data 字段
  async function unwrap<T>(url: string, opts?: any): Promise<T> {
    const res = await api<{ success: boolean; data: T }>(url, opts)
    return res.data
  }

  return {
    // 首页
    getHome: () => unwrap<any>('/public/home', { params: { lang: getLang() } }),

    // 产品（支持多维度筛选：gender, type, material, series_id, fabric_id, sort_by, sort_order）
    getProducts: (params?: any) => unwrap<any>('/public/products', { params: { lang: getLang(), ...params } }),
    getProduct: (slug: string) => unwrap<any>(`/public/products/${slug}`, { params: { lang: getLang() } }),

    // 页面
    getPage: (slug: string) => unwrap<any>(`/public/pages/${slug}`, { params: { lang: getLang() } }),

    // 博客
    getBlogs: (params?: any) => unwrap<any>('/public/blogs', { params: { lang: getLang(), ...params } }),
    getBlog: (slug: string) => unwrap<any>(`/public/blogs/${slug}`, { params: { lang: getLang() } }),

    // 案例
    getCases: (params?: any) => unwrap<any>('/public/cases', { params: { lang: getLang(), ...params } }),

    // FAQ
    getFaqs: (params?: any) => unwrap<any>('/public/faqs', { params: { lang: getLang(), ...params } }),

    // 分类/系列/面料
    getCategories: () => unwrap<any>('/public/categories'),
    getSeries: () => unwrap<any>('/public/series'),
    getFabrics: () => unwrap<any>('/public/fabrics'),

    // 工厂/认证/生产流程
    getFactories: () => unwrap<any>('/public/factories'),
    getCertifications: () => unwrap<any>('/public/certifications'),
    getProductionProcesses: () => unwrap<any>('/public/production-processes'),

    // 导航
    getNavigations: (type: string = 'header') => unwrap<any>('/public/navigations', { params: { type } }),

    // 主题
    getTheme: () => unwrap<any>('/public/theme'),

    // i18n 词典
    getI18n: () => unwrap<any>('/public/i18n', { params: { lang: getLang() } }),

    // 货币
    getCurrencies: () => unwrap<any>('/public/currencies'),

    // 语言
    getLanguages: () => unwrap<any>('/public/languages'),

    // 地理定位
    getGeo: (country?: string) => unwrap<any>('/public/geo', { params: { country } }),

    // 提交询盘
    submitLead: (data: any) => unwrap<any>('/public/leads', { method: 'POST', body: data }),

    // 外部链接点击跟踪（社交媒体跳转）
    trackClick: (target: string, url: string) => unwrap<any>('/public/click-track', { method: 'POST', body: { target, url } }),
  }
}
