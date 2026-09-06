/**
 * 多语言翻译记录的统一辅助函数。
 *
 * 背景：CMS 各实体（博客/案例/工厂/认证/生产流程/导航/页面/产品）此前在各自视图里
 * 重复实现 `emptyTranslations()` / `translationsToRecord()` / 构建 payload 的逻辑，
 * 字段列表略有差异导致"不是同种方式维护翻译"。这里收敛为单一实现，供各页面复用。
 *
 * 约定：
 * - 源语言固定为 English（en），主表字段即英文源；
 * - 目标语言默认 zh/es/fr（可通过 langs 参数扩展）；
 * - translations 以 `Record<lang, Record<field, string>>` 形式在视图内维护，
 *   提交时用 translationsToPayload 转回后端期望的 `{ language, ...fields }[]`。
 */

export interface TransLang {
  value: string
  label: string
}

export const DEFAULT_TRANS_LANGS: TransLang[] = [
  { value: 'zh', label: '中文' },
  { value: 'es', label: 'Español' },
  { value: 'fr', label: 'Français' },
]

/** 判断字段是否"有内容"：富文本需剥离标签后判断，避免空 <p> 被视为已填 */
export function hasTransText(v?: string): boolean {
  return Boolean((v || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').trim())
}

/** 生成空翻译记录（每种语言 × 每个字段 = 空字符串） */
export function createEmptyTranslations(
  fields: string[],
  langs: TransLang[] = DEFAULT_TRANS_LANGS,
): Record<string, Record<string, string>> {
  const blank: Record<string, string> = {}
  for (const f of fields) blank[f] = ''
  const next: Record<string, Record<string, string>> = {}
  for (const l of langs) next[l.value] = { ...blank }
  return next
}

/** 将后端返回的翻译列表 `{ language, ...fields }[]` 转回 Record 形式（跳过源语言 en） */
export function translationsToRecord(
  list: any[],
  fields: string[],
  langs: TransLang[] = DEFAULT_TRANS_LANGS,
): Record<string, Record<string, string>> {
  const next = createEmptyTranslations(fields, langs)
  for (const t of list || []) {
    const lang = t?.language
    if (lang && lang !== 'en' && next[lang]) {
      for (const f of fields) next[lang][f] = t[f] || ''
    }
  }
  return next
}

/** 将 Record 形式转为后端期望的翻译列表 `{ language, ...fields }[]`（只保留有内容的语言） */
export function translationsToPayload(
  record: Record<string, Record<string, string>>,
  fields: string[],
  langs: TransLang[] = DEFAULT_TRANS_LANGS,
): Array<Record<string, string>> {
  const out: Array<Record<string, string>> = []
  for (const l of langs) {
    const t = record[l.value]
    if (!t) continue
    if (!fields.some((f) => hasTransText(t[f]))) continue
    const item: Record<string, string> = { language: l.value }
    for (const f of fields) item[f] = t[f] || ''
    out.push(item)
  }
  return out
}

/**
 * 子资源（规格/视频/定制等）的"单字段多语言"以 JSONB 字符串存于 translations 字段，
 * 形如 `{"zh":"...","es":"...","fr":"..."}`。下面两个函数负责 JSONB ↔ Record 互转，
 * 供 InlineTrans 组件与提交/回填逻辑统一使用。
 */

/** 解析 JSONB translations 字符串为 Record（非法/空值安全回退空对象） */
export function parseJsonTrans(json?: string): Record<string, string> {
  if (!json) return {}
  try {
    const o = JSON.parse(json)
    return o && typeof o === 'object' && !Array.isArray(o) ? o : {}
  } catch {
    return {}
  }
}

/** 将 Record 序列化为 JSONB translations 字符串（仅保留非空值） */
export function buildJsonTrans(record: Record<string, string>): string {
  const obj: Record<string, string> = {}
  for (const [k, v] of Object.entries(record || {})) if (v) obj[k] = v
  return JSON.stringify(obj)
}
