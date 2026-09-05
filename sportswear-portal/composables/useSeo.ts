import { useHead } from '#imports'
import type { MetaObject } from '@nuxt/schema'

export interface SeoOptions {
  /** 页面标题（不含站点后缀） */
  title?: string | (() => string)
  /** 页面描述（150-160 字最佳） */
  description?: string | (() => string)
  /** 关键词（逗号分隔） */
  keywords?: string | (() => string)
  /** Open Graph 图片 URL（建议 1200×630） */
  ogImage?: string | (() => string)
  /** Open Graph 类型 */
  ogType?: string
  /** 是否添加站点后缀到 title（默认 true） */
  appendSiteName?: boolean
  /** 结构化数据类型 */
  schema?: Record<string, any> | (() => Record<string, any>)
  /** 是否添加 noindex（未发布/预览页用） */
  noindex?: boolean
  /** alternate hreflang 覆盖 */
  alternates?: Array<{ href: string; hreflang: string }>
  /** canonical URL 覆盖 */
  canonical?: string
}

const SITE_NAME = 'OEM/ODM Sportswear Manufacturer'
const SITE_URL = 'https://sportswear-platform.com'
const DEFAULT_OG_IMAGE = '/images/og-default.jpg'
const DEFAULT_DESCRIPTION =
  'Professional OEM/ODM sportswear manufacturer. Custom sportswear, activewear, and athletic apparel for global brands. 15+ years of experience, ISO certified.'
const DEFAULT_KEYWORDS =
  'sportswear manufacturer, OEM, ODM, custom sportswear, activewear, athletic apparel, China factory, private label'

/**
 * 统一 SEO 设置 composable
 */
export function useSeoHead(opts: SeoOptions) {
  const { locale } = useI18n()
  const localePath = useLocalePath()
  const route = useRoute()
  // 可从 runtimeConfig 读取自定义站点 URL（环境变量 NUXT_PUBLIC_SITE_URL）
  const siteUrl = useRuntimeConfig().public.siteUrl as string || SITE_URL

  const siteName = opts.appendSiteName !== false ? ` | ${SITE_NAME}` : ''

  const resolve = <T>(v: T | (() => T)): T =>
    typeof v === 'function' ? (v as () => T)() : v

  const resolvedTitle = computed(() => {
    const t = resolve(opts.title)
    if (!t) return SITE_NAME
    if (t.includes(SITE_NAME)) return t
    // 首页直接返回站点名
    if (route.path === '/' ||
        route.path.replace(/^\/[a-z]{2}/, '') === '' ||
        route.path.replace(/^\/[a-z]{2}\//, '/') === '/') {
      return SITE_NAME
    }
    return `${t}${siteName}`
  })

  const resolvedDescription = computed(() =>
    resolve(opts.description) || DEFAULT_DESCRIPTION
  )

  const resolvedKeywords = computed(() =>
    resolve(opts.keywords) || DEFAULT_KEYWORDS
  )

  const resolvedOgImage = computed(() => resolve(opts.ogImage) || DEFAULT_OG_IMAGE)
const alternates = computed(() => {
    if (opts.alternates) return opts.alternates
    return undefined
  })

  const headObj: MetaObject = {
    title: resolvedTitle.value as string,
    meta: [
      { name: 'description', content: resolvedDescription.value as string },
      { name: 'keywords', content: resolvedKeywords.value as string },

      // Open Graph
      { property: 'og:title', content: resolvedTitle.value as string },
      { property: 'og:description', content: resolvedDescription.value as string },
      { property: 'og:image', content: resolvedOgImage.value as string },
      { property: 'og:type', content: opts.ogType || 'website' },
      { property: 'og:url', content: `${siteUrl}${route.path}` },
      { property: 'og:site_name', content: SITE_NAME },

      // Twitter Card
      { name: 'twitter:card', content: 'summary_large_image' },
      { name: 'twitter:title', content: resolvedTitle.value as string },
      { name: 'twitter:description', content: resolvedDescription.value as string },
      { name: 'twitter:image', content: resolvedOgImage.value as string },
    ],
    link: [],
  }

  // noindex (预览页/未发布内容)
  if (opts.noindex) {
    (headObj.meta as any[]).push({ name: 'robots', content: 'noindex, nofollow' })
  }

  // canonical
  if (opts.canonical) {
    (headObj.link as any[]).push({ rel: 'canonical', href: opts.canonical })
  }

  // alternates
  if (alternates.value) {
    for (const alt of alternates.value) {
      (headObj.link as any[]).push({ rel: 'alternate', href: alt.href, hreflang: alt.hreflang })
    }
  }

  useHead(headObj)

  // 注入 JSON-LD 结构化数据
  const resolvedSchema = computed(() => resolve(opts.schema))
  if (resolvedSchema.value) {
    useHead({
      script: [
        {
          type: 'application/ld+json',
          children: JSON.stringify(resolvedSchema.value),
        },
      ],
    })
  }

  // SEO 健康自查（开发模式下输出到控制台）
  if (import.meta.client || import.meta.dev) {
    try {
      useSeoScorecard({
        title: resolvedTitle.value as string,
        description: resolvedDescription.value as string,
        hasH1: true,
        hasCanonical: !opts.canonical,
        ogTitle: true,
        ogDescription: true,
        ogImage: !!opts.ogImage,
        schemaPresent: !!opts.schema,
        breadcrumbPresent: true,
      })
    } catch {
      // Scorecard 仅为辅助工具，不影响核心功能
    }
  }

  return {
    resolvedTitle,
    resolvedDescription,
  }
}

/**
 * 生成产品页 Product Schema
 */
export function buildProductSchema(product: any): Record<string, any> {
  const images: string[] = []
  if (product.cover_image) images.push(product.cover_image)
  for (const img of product.images || []) {
    if (img.url && !images.includes(img.url)) images.push(img.url)
  }
  return {
    '@context': 'https://schema.org',
    '@type': 'Product',
    name: product.name || product.sku,
    description: product.brief || product.description || '',
    sku: product.sku,
    image: images.length ? images : undefined,
    category: product.type || undefined,
    brand: {
      '@type': 'Brand',
      name: SITE_NAME,
    },
    offers: {
      '@type': 'Offer',
      availability: 'https://schema.org/InStock',
      price: product.price || undefined,
      priceCurrency: 'USD',
    },
  }
}

/**
 * 生成博客文章 Article Schema
 */
export function buildArticleSchema(article: any): Record<string, any> {
  return {
    '@context': 'https://schema.org',
    '@type': 'Article',
    headline: article.title,
    description: article.content?.slice(0, 160) || '',
    image: article.image || undefined,
    datePublished: article.created_at || article.published_at || undefined,
    dateModified: article.updated_at || undefined,
    author: {
      '@type': 'Organization',
      name: SITE_NAME,
    },
  }
}

/**
 * 生成 FAQ Schema
 */
export function buildFaqSchema(faqs: Array<{ question: string; answer: string }>): Record<string, any> {
  return {
    '@context': 'https://schema.org',
    '@type': 'FAQPage',
    mainEntity: faqs.map((f) => ({
      '@type': 'Question',
      name: f.question,
      acceptedAnswer: {
        '@type': 'Answer',
        text: f.answer,
      },
    })),
  }
}

/**
 * 生成面包屑 BreadcrumbList Schema
 */
export function buildBreadcrumbSchema(
  items: Array<{ name: string; item: string }>
): Record<string, any> {
  return {
    '@context': 'https://schema.org',
    '@type': 'BreadcrumbList',
    itemListElement: items.map((item, index) => ({
      '@type': 'ListItem',
      position: index + 1,
      name: item.name,
      item: item.item,
    })),
  }
}
