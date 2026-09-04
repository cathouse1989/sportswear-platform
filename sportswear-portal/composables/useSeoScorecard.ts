// SEO 健康自查工具
// 在每个页面运行时检查组件的 SEO 要素是否完整，
// 结果可通过 console 输出（开发模式）或上报到监控系统。
//
// 用法：
// ```ts
// useSeoScorecard({
//   title: 'My Page Title',
//   description: 'Page description...',
//   hasH1: true,
//   hasCanonical: true,
// })
// ```

export interface SeoScorecardInput {
  title?: string
  description?: string
  hasH1?: boolean
  hasCanonical?: boolean
  ogTitle?: boolean
  ogDescription?: boolean
  ogImage?: boolean
  schemaPresent?: boolean
  breadcrumbPresent?: boolean
  pageType?: 'home' | 'product' | 'blog' | 'case' | 'about' | 'contact' | 'faq' | 'page'
}

export interface SeoScorecardResult {
  url: string
  locale: string
  pageType: string
  checks: {
    title: { pass: boolean; message: string }
    titleLength: { pass: boolean; message: string }
    description: { pass: boolean; message: string }
    descriptionLength: { pass: boolean; message: string }
    hasH1: { pass: boolean; message: string }
    canonical: { pass: boolean; message: string }
    ogTags: { pass: boolean; message: string }
    schema: { pass: boolean; message: string }
  }
  score: number
  scoreLevel: 'poor' | 'needs-improvement' | 'good' | 'excellent'
  timestamp: number
}

const TITLE_MAX_LENGTH = 60
const DESCRIPTION_MAX_LENGTH = 160
const DESCRIPTION_MIN_LENGTH = 50

/**
 * 在开发/调试模式下检测页面 SEO 要素并输出到控制台
 * 生产环境可关闭（通过环境变量控制）
 */
export function useSeoScorecard(input: SeoScorecardInput) {
  const route = useRoute()
  const { locale } = useI18n()
  const isDev = import.meta.dev || import.meta.client && window.location.hostname === 'localhost'

  const result: SeoScorecardResult = {
    url: route.path,
    locale: locale.value || 'en',
    pageType: input.pageType || 'page',
    checks: {
      title: {
        pass: !!input.title,
        message: input.title ? 'Title set' : '⚠️ Title is missing!',
      },
      titleLength: {
        pass: (input.title?.length || 0) <= TITLE_MAX_LENGTH,
        message: input.title && input.title.length > TITLE_MAX_LENGTH
          ? `⚠️ Title too long (${input.title.length} chars, max ${TITLE_MAX_LENGTH})`
          : input.title
            ? `✓ Title length OK (${input.title.length} chars)`
            : '⚠️ No title to check',
      },
      description: {
        pass: !!input.description,
        message: input.description ? 'Description set' : '⚠️ Meta description is missing!',
      },
      descriptionLength: {
        pass: input.description
          ? input.description.length >= DESCRIPTION_MIN_LENGTH && input.description.length <= DESCRIPTION_MAX_LENGTH
          : false,
        message: input.description
          ? input.description.length < DESCRIPTION_MIN_LENGTH
            ? `⚠️ Description too short (${input.description.length} chars, min ${DESCRIPTION_MIN_LENGTH})`
            : input.description.length > DESCRIPTION_MAX_LENGTH
              ? `⚠️ Description too long (${input.description.length} chars, max ${DESCRIPTION_MAX_LENGTH})`
              : `✓ Description length OK (${input.description.length} chars)`
          : '⚠️ No description to check',
      },
      hasH1: {
        pass: input.hasH1 !== false,
        message: input.hasH1 !== false ? '✓ H1 present' : '⚠️ H1 missing!',
      },
      canonical: {
        pass: input.hasCanonical !== false,
        message: input.hasCanonical !== false ? '✓ Canonical URL set' : '⚠️ Canonical URL missing!',
      },
      ogTags: {
        pass: !!(input.ogTitle && input.ogDescription && input.ogImage),
        message: input.ogTitle && input.ogDescription && input.ogImage
          ? '✓ OG tags complete'
          : '⚠️ OG tags incomplete (title/description/image)',
      },
      schema: {
        pass: input.schemaPresent !== false,
        message: input.schemaPresent !== false
          ? '✓ Structured data present'
          : '⚠️ Structured data (JSON-LD) missing!',
      },
    },
    score: 100,
    scoreLevel: 'excellent',
    timestamp: Date.now(),
  }

  // 计算评分
  const checks = Object.values(result.checks)
  const passed = checks.filter(c => c.pass).length
  result.score = Math.round((passed / checks.length) * 100)

  if (result.score === 100) result.scoreLevel = 'excellent'
  else if (result.score >= 80) result.scoreLevel = 'good'
  else if (result.score >= 60) result.scoreLevel = 'needs-improvement'
  else result.scoreLevel = 'poor'

  // 开发环境输出到控制台
  if (isDev) {
    const badge = result.scoreLevel === 'excellent' ? '✅' :
      result.scoreLevel === 'good' ? '🟢' :
      result.scoreLevel === 'needs-improvement' ? '🟡' : '🔴'

    console.group(`${badge} SEO Scorecard [${result.url}] — ${result.score}/100 (${result.scoreLevel})`)
    for (const [, check] of Object.entries(result.checks)) {
      const icon = check.pass ? '✅' : '⚠️'
      console.log(`${icon} ${check.message}`)
    }
    console.groupEnd()
  }

  return result
}