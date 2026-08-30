import { resolve } from 'node:path'

export default defineNuxtConfig({
  devtools: { enabled: true },

  // 允许通过 IP 访问
  devServer: {
    host: '0.0.0.0',
  },

  modules: [
    '@nuxtjs/tailwindcss',
    '@nuxtjs/i18n',
  ],

  css: [
    '~/assets/css/main.css',
  ],

  // SSR 服务端渲染（SEO 核心）
  ssr: true,

  // 运行时配置（后端 API 地址）
  runtimeConfig: {
    // 仅服务端（SSR）使用的 API 地址；Docker 内通过 NUXT_API_SERVER 指向容器网络
    apiServer: process.env.NUXT_API_SERVER || 'http://localhost:8080/api/v1',
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1',
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
    strategy: 'prefix', // 子路径前缀：/en/, /zh/, /es/, /fr/
    langDir: resolve(__dirname, 'locales'),
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: 'i18n_redirected',
      redirectOn: 'root',
    },
    seo: true, // 自动生成 hreflang 标签
  },

  // SEO 优化
  app: {
    head: {
      charset: 'utf-8',
      viewport: 'width=device-width, initial-scale=1',
      title: 'OEM/ODM Sportswear Manufacturer - Custom Sportswear Solutions',
      meta: [
        { name: 'description', content: 'Professional OEM/ODM sportswear manufacturer. Custom sportswear, activewear, and athletic apparel for global brands.' },
        { name: 'keywords', content: 'sportswear manufacturer, OEM, ODM, custom sportswear, activewear, athletic apparel, China factory' },
      ],
      // Favicon: SVG logo
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAb5JREFUWEftlrFKA0EQhv8kKioWamNhYyPY2Fj6AD6Aj2ChYCFYCNpYWPgCgoWNjY2FhY0vYCEoKIqFhYIiFjYqKipC4szdHpvc3d5eQgRvlWY2O9/885/dmYUg/AMBf4cAQoAQAoQAMQwBwQ8BQoAQIIQAMQwBwQ8BQoAQIIQAMQwBwQ8BQoAQIIQAMQwBwQ9P4K5Ww26jwV5GxYvLZRxfXHDEtU0NZ9c3I+d6bwQykoQ3p4Kbyi2OigUkRUHhYh8WZyeQzWQgiiJ0XYeuaRBFEYIgwPf3N/q9Hjv6qdfRbjSQV1TkslkYHQOZTAaDwQCO6wKANwLX1SqyBQWFYvHHCzQMA5qmQZYlfNbr2NvZQbVa5Rr7QeB9CxgqVwqFQgirq2Zot1p4eXpCQpZxelEOPREvBCRJQrVahW3b+P76QjwRx/PLK7R6/ZET6IXA1cUF1tbXcXN9jUgkglarBQBIJBKYnpmBlU7j4fER2UwmpAQEIUAIEEKAGIYAZ9YQIIQAIUAME8KZ5c+aWULAz6a8nwAhQAgQQoAYhoDghwAhQAgQQoAYhoDghwAhQAgQQoAYhoDghwAhQAgQQoD4AcrFhE3SXf0PAAAAAElFTkSuQmCC' },
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
    '/en/**': { swr: 300 },
    '/zh/**': { swr: 300 },
    '/es/**': { swr: 300 },
    '/fr/**': { swr: 300 },
    // 构建产物（_nuxt/*.js/css）做长缓存，改文件后哈希变化自动失效
    '/_nuxt/**': {
      headers: { 'cache-control': 'public, max-age=31536000, immutable' },
    },
  },

  compatibilityDate: '2025-01-01',
})