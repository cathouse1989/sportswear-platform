// Core Web Vitals 监控 composable
// 使用 web-vitals 库在浏览器端采集 LCP / INP (FID) / CLS，
// 上报到后端 API /public/analytics/web-vitals 用于监控面板
//
// 仅在用户同意 ANALYTICS cookie 后才采集并上报
//
// 安装依赖：npm install web-vitals

export interface WebVitalMetric {
  name: 'LCP' | 'FID' | 'INP' | 'CLS' | 'TTFB'
  value: number
  rating: 'good' | 'needs-improvement' | 'poor'
  delta: number
  id: string
  url?: string
  timestamp: number
}

const CONSENT_STORAGE_KEY = 'sw_cookie_consent'

/**
 * 检查用户是否已同意分析 Cookie
 */
function hasAnalyticsConsent(): boolean {
  if (import.meta.server) return false
  try {
    const raw = localStorage.getItem(CONSENT_STORAGE_KEY)
    if (!raw) return false
    const state = JSON.parse(raw)
    return state.categories?.includes('analytics') ?? false
  } catch {
    return false
  }
}

/**
 * 获取页面加载性能基准数据
 */
function getNavigationTiming(): Record<string, number> {
  if (import.meta.server || !performance?.getEntriesByType) return {}
  try {
    const nav = performance.getEntriesByType('navigation')?.[0] as PerformanceNavigationTiming | undefined
    if (!nav) return {}
    return {
      dns: nav.domainLookupEnd - nav.domainLookupStart,
      tcp: nav.connectEnd - nav.connectStart,
      tls: nav.secureConnectionStart > 0 ? nav.connectEnd - nav.secureConnectionStart : 0,
      ttfb: nav.responseStart - nav.requestStart,
      dom_ready: nav.domComplete - nav.domInteractive,
      load: nav.loadEventEnd - nav.loadEventStart,
    }
  } catch {
    return {}
  }
}

/**
 * 上报 Web Vitals 到后端
 */
async function reportMetric(metric: WebVitalMetric) {
  if (!hasAnalyticsConsent()) return

  // 合并导航时间数据（仅首次上报时附带）
  const navTiming = getNavigationTiming()

  try {
    const config = useRuntimeConfig()
    const base = config.public.apiBase as string || 'http://localhost:8080/api/v1'
    await $fetch(`${base}/public/analytics/web-vitals`, {
      method: 'POST',
      body: {
        ...metric,
        navigation_timing: Object.keys(navTiming).length > 0 ? navTiming : undefined,
        user_agent: navigator.userAgent,
        viewport: `${window.innerWidth}x${window.innerHeight}`,
      },
    })
  } catch {
    // 静默失败，不影响用户体验
  }
}

/**
 * 初始化 Web Vitals 监控
 *
 * 在 layout 或 app.vue 中调用：
 * ```ts
 * if (import.meta.client) {
 *   initWebVitals()
 * }
 * ```
 */
export function initWebVitals() {
  if (import.meta.server) return
  if (!hasAnalyticsConsent()) return

  // 动态导入 web-vitals
  // 注意：需要在项目中先安装 web-vitals 包
  import('web-vitals').then(({ onLCP, onINP, onCLS, onTTFB }) => {
    const ratingThresholds: Record<string, { good: number; poor: number }> = {
      LCP: { good: 2500, poor: 4000 },
      INP: { good: 200, poor: 500 },
      CLS: { good: 0.1, poor: 0.25 },
      TTFB: { good: 800, poor: 1800 },
    }

    const getRating = (name: string, value: number): 'good' | 'needs-improvement' | 'poor' => {
      const t = ratingThresholds[name]
      if (!t) return 'needs-improvement'
      if (value <= t.good) return 'good'
      if (value <= t.poor) return 'needs-improvement'
      return 'poor'
    }

    const send = (m: { name: string; value: number; delta: number; id: string }) => {
      const metric: WebVitalMetric = {
        name: m.name as WebVitalMetric['name'],
        value: m.name === 'CLS' ? m.value * 1000 : m.value, // CLS 转千分数便于展示
        rating: getRating(m.name, m.value),
        delta: m.delta,
        id: m.id,
        url: window.location.pathname,
        timestamp: Date.now(),
      }
      reportMetric(metric)
    }

    onLCP(send)
    onINP(send)
    onCLS(send)
    onTTFB(send)
  }).catch(() => {
    // web-vitals 未安装，静默跳过
    console.warn('[WebVitals] web-vitals library not installed. Skipping monitoring.')
  })
}

/**
 * 手动触发 Web Vitals 上报（用于自定义场景）
 */
export function reportWebVital(name: WebVitalMetric['name'], value: number) {
  const metric: WebVitalMetric = {
    name,
    value,
    rating: 'needs-improvement',
    delta: 0,
    id: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
    url: window.location.pathname,
    timestamp: Date.now(),
  }
  reportMetric(metric)
}