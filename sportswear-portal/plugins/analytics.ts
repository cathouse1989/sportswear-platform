// Google Analytics 4 + Search Console / Webmaster Tools 集成
// - GA4 gtag 基于用户 cookie consent 合规加载（尊重 ANALYTICS 类别的同意状态）
// - 仅在客户端加载，SSR 不注入
//
// 环境变量：
//   NUXT_PUBLIC_GA_MEASUREMENT_ID — GA4 测量 ID（如 G-XXXXXXXXXX）
//   NUXT_PUBLIC_GOOGLE_SITE_VERIFICATION — Google Search Console 验证字符串
//
// 注意：Search Console 验证 meta 需要在 SSR head 中注入（让 Google 爬虫读到），
// 而 GA4 gtag 只应在客户端注入（且仅当用户同意后）。

export default defineNuxtPlugin({
  name: 'analytics',
  setup() {
    const config = useRuntimeConfig()
    const gaMeasurementId = config.public.gaMeasurementId as string || ''
    const siteVerification = config.public.googleSiteVerification as string || ''

    // 1. Google Search Console 验证 — 始终在 SSR head 中注入（爬虫需要）
    if (siteVerification) {
      useHead({
        meta: [
          { name: 'google-site-verification', content: siteVerification },
        ],
      })
    }

    // 2. GA4 gtag — 仅在浏览器端加载
    if (!import.meta.client) return
    if (!gaMeasurementId) return

    const loadGtag = () => {
      // 防止重复加载
      if (document.getElementById('ga-gtag')) return

      // 加载 gtag 脚本
      const script = document.createElement('script')
      script.id = 'ga-gtag'
      script.src = `https://www.googletagmanager.com/gtag/js?id=${gaMeasurementId}`
      script.async = true
      document.head.appendChild(script)

      // 初始化 gtag
      window.dataLayer = window.dataLayer || []
      // @ts-ignore
      function gtag() { window.dataLayer.push(arguments) }
      window.gtag = gtag as any
      gtag('js', new Date())
      gtag('config', gaMeasurementId, {
        anonymize_ip: true,         // IP 匿名化
        allow_google_signals: false, // 禁用广告功能
        allow_ad_personalization_signals: false,
      })

      // 设置默认拒绝所有非必要跟踪（等待用户同意）
      gtag('consent', 'default', {
        analytics_storage: 'denied',
        ad_storage: 'denied',
        ad_user_data: 'denied',
        ad_personalization: 'denied',
        functionality_storage: 'granted',
        security_storage: 'granted',
        personalization_storage: 'denied',
        wait_for_update: 500,
      })
    }

    // 3. 初始化 Web Vitals 监控
    if (hasAnalyticsConsent()) {
      try {
        const config = useRuntimeConfig()
        initWebVitals()
      } catch {}
    }

    // 4. 监听 cookie consent 状态变化
    const updateConsent = (categories: string[]) => {
      if (!window.gtag) return
      const analyticsGranted = categories.includes('analytics')
      const marketingGranted = categories.includes('marketing')
      window.gtag('consent', 'update', {
        analytics_storage: analyticsGranted ? 'granted' : 'denied',
        ad_storage: marketingGranted ? 'granted' : 'denied',
        ad_user_data: marketingGranted ? 'granted' : 'denied',
        ad_personalization: marketingGranted ? 'granted' : 'denied',
      })
    }

    // 加载 gtag（首次）
    const savedConsent = localStorage.getItem('sw_cookie_consent')
    if (savedConsent) {
      try {
        const state = JSON.parse(savedConsent)
        if (state.categories?.includes('analytics')) {
          loadGtag()
          // 延迟一帧确保 gtag 已初始化
          requestAnimationFrame(() => updateConsent(state.categories))
        }
      } catch {
        // 静默失败
      }
    }

    // 监听后续 consent 变化
    window.addEventListener('consent-changed', ((e: CustomEvent) => {
      const state = e.detail
      if (state?.categories?.includes('analytics')) {
        loadGtag()
        requestAnimationFrame(() => updateConsent(state.categories))
      } else if (window.gtag) {
        updateConsent(state?.categories || ['necessary'])
      }
    }) as EventListener)

    // 监听 withdraw
    window.addEventListener('consent-withdrawn', () => {
      if (window.gtag) {
        window.gtag('consent', 'update', {
          analytics_storage: 'denied',
          ad_storage: 'denied',
          ad_user_data: 'denied',
          ad_personalization: 'denied',
        })
      }
    })
  },
})

// 为 window 添加 gtag 类型声明
declare global {
  interface Window {
    dataLayer: any[]
    gtag: (...args: any[]) => void
  }
}