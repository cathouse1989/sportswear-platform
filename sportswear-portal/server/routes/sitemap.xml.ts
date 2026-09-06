// 多语言 Sitemap 动态生成
// 用 Nitro server route 动态生成 /sitemap.xml，包含所有页面的多语言版本
// 运行时请求后端 API 获取产品、博客、案例等动态内容，保证 sitemap 始终最新

const SITE_URL = 'https://sportswear-platform.com'

// 支持的语言
const LOCALES: Array<{ code: string; iso: string }> = [
  { code: 'en', iso: 'en-US' },
  { code: 'zh', iso: 'zh-CN' },
  { code: 'es', iso: 'es-ES' },
  { code: 'fr', iso: 'fr-FR' },
]

// 静态路由（所有语言版本）
const STATIC_ROUTES = [
  '',         // 首页
  '/about',
  '/contact',
  '/products',
  '/blog',
  '/cases',
  '/faq',
  '/privacy-policy',
]

interface SitemapUrl {
  loc: string
  lastmod?: string
  changefreq?: string
  priority?: number
  alternates?: Array<{ href: string; hreflang: string }>
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event)
  const apiBase = config.apiServer || config.public.apiBase || 'http://localhost:8080/api/v1'

  const urls: SitemapUrl[] = []

  // 1. 静态页面（每个语言版本）
  for (const locale of LOCALES) {
    for (const route of STATIC_ROUTES) {
      const url: SitemapUrl = {
        loc: `${SITE_URL}/${locale.code}${route}`,
        changefreq: route === '' ? 'daily' : 'weekly',
        priority: route === '' ? '1.0' : '0.8' as any,
        alternates: LOCALES.map((l) => ({
          href: `${SITE_URL}/${l.code}${route}`,
          hreflang: l.iso,
        })),
      }
      urls.push(url)
    }
  }

  // 2. 动态产品页面
  try {
    const productRes = await $fetch<{ success: boolean; data: { items?: any[]; total?: number } }>(
      `${apiBase}/public/products`,
      { params: { page: 1, pageSize: 1000, lang: 'en' } }
    )
    const products = productRes?.data?.items || []
    for (const locale of LOCALES) {
      for (const p of products) {
        if (!p.slug) continue
        urls.push({
          loc: `${SITE_URL}/${locale.code}/products/${p.slug}`,
          lastmod: p.updated_at || p.created_at || undefined,
          changefreq: 'weekly',
          priority: '0.9' as any,
          alternates: LOCALES.map((l) => ({
            href: `${SITE_URL}/${l.code}/products/${p.slug}`,
            hreflang: l.iso,
          })),
        })
      }
    }
  } catch (e) {
    // 后端不可用时跳过动态部分
    console.error('Sitemap: failed to fetch products', e)
  }

  // 2.5 动态分类页面（含子分类，递归展平）
  try {
    const catRes = await $fetch<{ success: boolean; data: any }>(
      `${apiBase}/public/categories`
    )
    const raw = catRes?.data
    const cats: any[] = []
    const flatten = (list: any[]) => {
      for (const c of list || []) {
        if (c.slug) cats.push(c)
        if (c.children?.length) flatten(c.children)
      }
    }
    flatten(Array.isArray(raw) ? raw : (raw?.items || []))
    for (const locale of LOCALES) {
      for (const c of cats) {
        urls.push({
          loc: `${SITE_URL}/${locale.code}/categories/${c.slug}`,
          changefreq: 'weekly',
          priority: '0.7' as any,
          alternates: LOCALES.map((l) => ({
            href: `${SITE_URL}/${l.code}/categories/${c.slug}`,
            hreflang: l.iso,
          })),
        })
      }
    }
  } catch (e) {
    console.error('Sitemap: failed to fetch categories', e)
  }

  // 2.6 动态系列页面
  try {
    const seriesRes = await $fetch<{ success: boolean; data: any }>(
      `${apiBase}/public/series`
    )
    const rawSeries = seriesRes?.data
    const series = Array.isArray(rawSeries) ? rawSeries : (rawSeries?.items || [])
    for (const locale of LOCALES) {
      for (const s of series) {
        if (!s.slug) continue
        urls.push({
          loc: `${SITE_URL}/${locale.code}/series/${s.slug}`,
          changefreq: 'weekly',
          priority: '0.7' as any,
          alternates: LOCALES.map((l) => ({
            href: `${SITE_URL}/${l.code}/series/${s.slug}`,
            hreflang: l.iso,
          })),
        })
      }
    }
  } catch (e) {
    console.error('Sitemap: failed to fetch series', e)
  }

  // 3. 动态博客页面
  try {
    const blogRes = await $fetch<{ success: boolean; data: { items?: any[]; total?: number } }>(
      `${apiBase}/public/blogs`,
      { params: { page: 1, pageSize: 1000, lang: 'en' } }
    )
    const blogs = blogRes?.data?.items || []
    for (const locale of LOCALES) {
      for (const b of blogs) {
        if (!b.slug) continue
        urls.push({
          loc: `${SITE_URL}/${locale.code}/blog/${b.slug}`,
          lastmod: b.updated_at || b.created_at || undefined,
          changefreq: 'monthly',
          priority: '0.7' as any,
          alternates: LOCALES.map((l) => ({
            href: `${SITE_URL}/${l.code}/blog/${b.slug}`,
            hreflang: l.iso,
          })),
        })
      }
    }
  } catch (e) {
    console.error('Sitemap: failed to fetch blogs', e)
  }

  // 4. 动态案例页面
  try {
    const caseRes = await $fetch<{ success: boolean; data: { items?: any[]; total?: number } }>(
      `${apiBase}/public/cases`,
      { params: { page: 1, pageSize: 1000, lang: 'en' } }
    )
    const cases = caseRes?.data?.items || []
    for (const locale of LOCALES) {
      for (const c of cases) {
        if (!c.slug) continue
        urls.push({
          loc: `${SITE_URL}/${locale.code}/cases/${c.slug}`,
          lastmod: c.updated_at || c.created_at || undefined,
          changefreq: 'monthly',
          priority: '0.6' as any,
          alternates: LOCALES.map((l) => ({
            href: `${SITE_URL}/${l.code}/cases/${c.slug}`,
            hreflang: l.iso,
          })),
        })
      }
    }
  } catch (e) {
    console.error('Sitemap: failed to fetch cases', e)
  }

  // 5. 动态自定义落地页（CMS 页面，走 /:slug 兜底路由）
  try {
    const pagesRes = await $fetch<{ success: boolean; data: any[] }>(
      `${apiBase}/public/pages`
    )
    const pages = pagesRes?.data || []
    for (const locale of LOCALES) {
      for (const pg of pages) {
        if (!pg.slug) continue
        urls.push({
          loc: `${SITE_URL}/${locale.code}/${pg.slug}`,
          lastmod: pg.updated_at || pg.published_at || undefined,
          changefreq: 'monthly',
          priority: '0.6' as any,
          alternates: LOCALES.map((l) => ({
            href: `${SITE_URL}/${l.code}/${pg.slug}`,
            hreflang: l.iso,
          })),
        })
      }
    }
  } catch (e) {
    console.error('Sitemap: failed to fetch landing pages', e)
  }

  // 生成 XML
  const xml = generateSitemapXml(urls)

  // 设置响应头
  setHeader(event, 'Content-Type', 'application/xml; charset=utf-8')
  setHeader(event, 'Cache-Control', 'public, max-age=3600, s-maxage=3600')

  return xml
})

function generateSitemapXml(urls: SitemapUrl[]): string {
  let xml = `<?xml version="1.0" encoding="UTF-8"?>\n`
  xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"\n`
  xml += `  xmlns:xhtml="http://www.w3.org/1999/xhtml">\n`

  for (const url of urls) {
    xml += `  <url>\n`
    xml += `    <loc>${escapeXml(url.loc)}</loc>\n`

    if (url.lastmod) {
      const d = new Date(url.lastmod)
      if (!isNaN(d.getTime())) {
        xml += `    <lastmod>${d.toISOString().split('T')[0]}</lastmod>\n`
      }
    }
    if (url.changefreq) {
      xml += `    <changefreq>${url.changefreq}</changefreq>\n`
    }
    if (url.priority) {
      xml += `    <priority>${url.priority}</priority>\n`
    }

    // xhtml:link for hreflang alternates
    if (url.alternates) {
      const seen = new Set<string>()
      for (const alt of url.alternates) {
        if (seen.has(alt.hreflang)) continue
        seen.add(alt.hreflang)
        xml += `    <xhtml:link rel="alternate" hreflang="${alt.hreflang}" href="${escapeXml(alt.href)}" />\n`
      }
      // x-default (default to English)
      if (!seen.has('x-default')) {
        xml += `    <xhtml:link rel="alternate" hreflang="x-default" href="${escapeXml(url.loc.replace(/\/[a-z]{2}\//, '/en/'))}" />\n`
      }
    }

    xml += `  </url>\n`
  }

  xml += `</urlset>`
  return xml
}

function escapeXml(s: string): string {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;')
}