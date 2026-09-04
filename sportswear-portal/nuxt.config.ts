export default defineNuxtConfig({
  devtools: { enabled: true },

  // 允许通过 IP 访问
  devServer: {
    host: '0.0.0.0',
  },

  modules: [
    '@nuxtjs/tailwindcss',
    '@nuxtjs/i18n',
    '@nuxt/image',
  ],

  // @nuxt/image 配置 — 自托管 IPX 图片优化
  image: {
    provider: 'ipx',
    format: ['webp', 'avif'],
    quality: 80,
    screens: {
      xs: 320,
      sm: 640,
      md: 768,
      lg: 1024,
      xl: 1280,
      xxl: 1536,
    },
    densities: [1, 2],
  },

  css: [
    '~/assets/css/main.css',
  ],

  // SSR 服务端渲染（SEO 核心）
  ssr: true,

  // 运行时配置（后端 API 地址 + 第三方集成）
  runtimeConfig: {
    // 仅服务端（SSR）使用的 API 地址；Docker 内通过 NUXT_API_SERVER 指向容器网络
    apiServer: process.env.NUXT_API_SERVER || 'http://localhost:8080/api/v1',
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',

      // Google Analytics 4 测量 ID（如 G-XXXXXXXXXX）
      gaMeasurementId: process.env.NUXT_PUBLIC_GA_MEASUREMENT_ID || '',

      // Google Search Console 站点验证字符串
      googleSiteVerification: process.env.NUXT_PUBLIC_GOOGLE_SITE_VERIFICATION || '',

      // 站点 URL（用于 canonical / OG / sitemap）
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'https://sportswear-platform.com',
    },
  },

  // 多语言配置（子路径方案：/en/products, /zh/products）
  i18n: {
    locales: [
      { code: 'en', name: 'English', flag: '🇬🇧', iso: 'en-US', file: 'en.json' },
      { code: 'zh', name: '中文', flag: '🇨🇳', iso: 'zh-CN', file: 'zh.json' },
      { code: 'es', name: 'Español', flag: '🇪🇸', iso: 'es-ES', file: 'es.json' },
      { code: 'fr', name: 'Français', flag: '🇫🇷', iso: 'fr-FR', file: 'fr.json' },
    ],
    defaultLocale: 'en',
    strategy: 'prefix',
    // v9 默认 restructureDir 为 "i18n"，langDir 会解析到 <root>/i18n/locales；
    // 本项目语言包位于项目根 <root>/locales，关闭 restructure 使其相对 srcDir 解析。
    restructureDir: false,
    langDir: 'locales',
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
    },
    seo: true,
    // 显式关闭 optimizeTranslationDirective 避免告警（v9 特性，与项目无冲突）
    bundle: {
      optimizeTranslationDirective: false,
    },
  },

  // 资源预连接和预加载提示（提升 LCP/font 加载性能）
  app: {
    head: {
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1',
      title: 'OEM/ODM Sportswear Manufacturer - Custom Sportswear Solutions',
      meta: [
        { name: 'description', content: 'Professional OEM/ODM sportswear manufacturer. Custom sportswear, activewear, and athletic apparel for global brands.' },
        { name: 'keywords', content: 'sportswear manufacturer, OEM, ODM, custom sportswear, activewear, athletic apparel, China factory' },
        // X-Robots-Tag（通过 meta 控制索引行为）
        { name: 'robots', content: 'index, follow, max-snippet:-1, max-image-preview:large' },
      ],
      // Favicon + 资源预连接
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAb5JREFUWEftlrFKA0EQhv8kKioWamNhYyPY2Fj6AD6Aj2ChYCFYCNpYWPgCgoWNjY2FhY0vYCEoKIqFhYIiFjYqKipC4szdHpvc3d5eQgRvlWY2O9/885/dmYUg/AMBf4cAQoAQAoQAMQwBwQ8BQoAQIIQAMQwBwQ8BQoAQIIQAMQwBwQ8BQoAQIIQAMQwBwQ9P4K5Ww26jwV5GxYvLZRxfXHDEtU0NZ9c3I+d6bwQykoQ3p4Kbyi2OigUkRUHhYh8WZyeQzWQgiiJ0XYeuaRBFEYIgwPf3N/q9Hjv6qdfRbjSQV1TkslkYHQOZTAaDwQCO6wKANwLX1SqyBQWFYvHHCzQMA5qmQZYlfNbr2NvZQbVa5Rr7QeB9CxgqVwqFQgirq2Zot1p4eXpCQpZxelEOPREvBCRJQrVahW3b+P76QjwRx/PLK7R6/ZET6IXA1cUF1tbXcXN9jUgkglarBQBIJBKYnpmBlU7j4fER2UwmpAQEIUAIEEKAGIYAZ9YQIIQAIUAME8KZ5c+aWULAz6a8nwAhQAgQQoAYhoDghwAhQAgQQoAYhoDghwAhQAgQQoAYhoDghwAhQAgQQoD4AcrFhE3SXf0PAAAAAElFTkSuQmCC' },
        // Google Fonts 预连接（当使用 Google Fonts 时生效）
        { rel: 'preconnect', href: 'https://fonts.googleapis.com', crossorigin: 'anonymous' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: 'anonymous' },
        // DNS 预取：CDN 域名（部署时替换为实际 CDN 域名）
        { rel: 'dns-prefetch', href: 'https://fonts.googleapis.com' },
        { rel: 'dns-prefetch', href: 'https://fonts.gstatic.com' },
      ],
    },
  },

  // 不配置 nitro.prerender —— /sitemap.xml 由后端 Go 服务提供，
  // 构建期无后端可用，预渲染会因 404 导致构建失败

  // 输出缓存（HTML 层缓存，P1）
  // 门户是 SSR + 客户端数据拉取，页面主体数据由客户端异步从后端 API 获取（已叠加 Redis 缓存），
  // 因此对稳定 locale 子路径启用 Nitro SWR，可大幅减少服务端渲染开销，命中后直接返回 HTML。
  // 静态资源与 /sitemap.xml 由后端 Go 提供；这里仅缓存各 locale 的渲染页。
  routeRules: {
    // 首页轮播等配置需即时生效，locale 根路径不走 SWR
    '/en': { swr: false, headers: { 'cache-control': 'no-store' } },
    '/zh': { swr: false, headers: { 'cache-control': 'no-store' } },
    '/es': { swr: false, headers: { 'cache-control': 'no-store' } },
    '/fr': { swr: false, headers: { 'cache-control': 'no-store' } },
    // 子页面 SWR（5 分钟）+ CDN 边缘缓存（1 小时），s-maxage 对 CDN 生效
    '/en/**': { swr: 300, headers: { 'cache-control': 'public, s-maxage=3600, max-age=0, stale-while-revalidate=300' } },
    '/zh/**': { swr: 300, headers: { 'cache-control': 'public, s-maxage=3600, max-age=0, stale-while-revalidate=300' } },
    '/es/**': { swr: 300, headers: { 'cache-control': 'public, s-maxage=3600, max-age=0, stale-while-revalidate=300' } },
    '/fr/**': { swr: 300, headers: { 'cache-control': 'public, s-maxage=3600, max-age=0, stale-while-revalidate=300' } },
    // 构建产物（_nuxt/*.js/css）做长缓存，改文件后哈希变化自动失效
    '/_nuxt/**': {
      headers: { 'cache-control': 'public, max-age=31536000, immutable' },
    },
    // 公开图片（用户上传的 /uploads/**）CDN 缓存 1 天
    '/uploads/**': {
      headers: { 'cache-control': 'public, max-age=86400, s-maxage=86400' },
    },
    // 站点地图短期缓存
    '/sitemap.xml': {
      headers: { 'cache-control': 'public, max-age=3600, s-maxage=3600' },
    },
    // robots.txt 1 天缓存
    '/robots.txt': {
      headers: { 'cache-control': 'public, max-age=86400, s-maxage=86400' },
    },
  },

  // estree-walker@3 是纯 ESM（exports 只含 import/types，无 require 条件）。
  // Nitro 若将其外部化到 .output/server/node_modules，运行期 CJS require 会抛
  // ERR_PACKAGE_PATH_NOT_EXPORTED: No "exports" main defined，导致 SSR 全部 500。
  // 这里强制内联打包，避免运行期对外部 ESM-only 包做 require。
  nitro: {
    externals: {
      inline: ['estree-walker'],
    },
    // 服务端压缩（gzip / brotli）
    compressPublicAssets: true,
    // 预渲染 sitemap fallback 页面路径
    prerender: {
      crawlLinks: false,
      routes: ['/robots.txt'],
    },
  },

  compatibilityDate: '2025-01-01',
})