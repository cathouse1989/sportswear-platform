/**
 * 内容发布质检门（P0-#2）
 * 对齐《管理后台产品说明书》第七章质检基线：
 * - severity=error   缺失 → 阻断发布（标题/正文/封面/MOQ 等门户核心要素）
 * - severity=warning 缺失 → 提示但不拦截（SEO 建议项 / 非核心素材）
 * 接入点：各实体 handlePublish 在调用 publish 接口前执行检查。
 */

export type PublishSeverity = 'error' | 'warning'

export interface PublishGateItem {
  label: string
  ok: boolean
  severity: PublishSeverity
}

export interface PublishGate {
  /** true = 无阻断项（warning 不影响 ok） */
  ok: boolean
  items: PublishGateItem[]
  errors: string[]
  warnings: string[]
}

const passItem = (label: string): PublishGateItem => ({ label, ok: true, severity: 'error' })
const errorItem = (label: string): PublishGateItem => ({ label, ok: false, severity: 'error' })
const warnItem = (label: string): PublishGateItem => ({ label, ok: false, severity: 'warning' })

const hasText = (v?: string | null) => Boolean((v || '').trim())

function buildGate(items: PublishGateItem[]): PublishGate {
  return {
    ok: items.every((i) => i.ok || i.severity === 'warning'),
    items,
    errors: items.filter((i) => !i.ok && i.severity === 'error').map((i) => i.label),
    warnings: items.filter((i) => !i.ok && i.severity === 'warning').map((i) => i.label),
  }
}

/** 产品发布门：封面/英文名/简介/描述/MOQ 必填缺失阻断；SEO 标题缺失警告。列表接口全字段返回，可安全强检。 */
export function checkProductGate(p: {
  cover_image?: string | null
  name?: string | null
  brief?: string | null
  description?: string | null
  sample_moq?: number | null
  production_moq?: number | null
  seo_title?: string | null
}): PublishGate {
  const sample = Number(p.sample_moq ?? 0)
  const prod = Number(p.production_moq ?? 0)
  return buildGate([
    hasText(p.cover_image) ? passItem('封面图') : errorItem('封面图'),
    hasText(p.name) ? passItem('英文名称') : errorItem('英文名称'),
    hasText(p.brief) ? passItem('简介 brief') : errorItem('简介 brief'),
    hasText(p.description) ? passItem('描述 description') : errorItem('描述 description'),
    sample >= 1 && prod >= 1 && sample <= prod ? passItem('MOQ 范围（sample ≤ production）') : errorItem('MOQ 范围（sample ≤ production）'),
    hasText(p.seo_title) ? passItem('SEO 标题') : warnItem('SEO 标题'),
  ])
}

/** 博客发布门：标题/Slug/正文缺失阻断；封面图缺失仅警告（博客无图也可发布）。 */
export function checkBlogGate(b: {
  title?: string | null
  slug?: string | null
  content?: string | null
  cover_image?: string | null
}): PublishGate {
  return buildGate([
    hasText(b.title) ? passItem('标题') : errorItem('标题'),
    hasText(b.slug) ? passItem('Slug') : errorItem('Slug'),
    hasText(b.content) ? passItem('正文') : errorItem('正文'),
    hasText(b.cover_image) ? passItem('封面图') : warnItem('封面图'),
  ])
}

/** 生成发布拦截提示文案（供 ElMessageBox.alert 使用） */
export function gateAlertMessage(gate: PublishGate): string {
  const lines: string[] = []
  if (gate.errors.length) {
    lines.push(`该内容缺少以下必填信息，无法发布：\n\n• ${gate.errors.join('\n• ')}`)
  }
  if (gate.warnings.length) {
    lines.push(`建议补充（不阻断）：\n• ${gate.warnings.join('\n• ')}`)
  }
  return lines.join('\n\n')
}