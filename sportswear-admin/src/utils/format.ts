/**
 * 统一时间格式化工具（管理后台）
 *
 * 时间处理遵循《sportswear-backend/docs/全球时区处理方案.md》三大原则：
 * 1. 存储统一 UTC；2. 传输统一 RFC3339（Z 后缀）；3. 展示按访客（浏览器）时区。
 *
 * 后端 API 返回的都是 UTC 时间，管理后台所有展示给运营人员的时间，
 * 一律通过本工具按【浏览器本地时区】格式化，避免异地运营看到 UTC 原始字符串产生歧义。
 */

/** 获取浏览器时区 IANA 名称（如 Asia/Shanghai / America/New_York），失败回退 UTC。 */
export function browserTimeZone(): string {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC'
  } catch {
    return 'UTC'
  }
}

/**
 * 浏览器时区相对 UTC 的偏移文案（如 UTC+08:00 / UTC-05:00）。
 * 基于 Date.getTimezoneOffset()：东八区返回 480 → UTC+08:00。
 */
export function timeZoneOffset(date: Date = new Date()): string {
  const minutes = -date.getTimezoneOffset()
  const sign = minutes >= 0 ? '+' : '-'
  const abs = Math.abs(minutes)
  const hh = String(Math.floor(abs / 60)).padStart(2, '0')
  const mm = String(abs % 60).padStart(2, '0')
  return `UTC${sign}${hh}:${mm}`
}

/**
 * 按指定时区（默认浏览器时区）把后端 RFC3339（UTC，Z 后缀）时间格式化为
 * `YYYY-MM-DD HH:mm:ss`。空值 / 非法值返回 fallback（默认 '-'）。
 * @param iso      后端返回的时间字符串（ISO8601 / RFC3339，带 Z 后缀）
 * @param fallback 空值或非法值时的占位文案
 * @param timeZone IANA 时区名（默认浏览器时区；主要供测试指定固定时区）
 */
export function formatDateTime(
  iso?: string | null,
  fallback = '-',
  timeZone: string = browserTimeZone(),
): string {
  if (!iso) return fallback
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return fallback
  try {
    const parts = new Intl.DateTimeFormat('en-US', {
      timeZone,
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hourCycle: 'h23',
    }).formatToParts(d)
    const values: Record<string, string> = {}
    for (const part of parts) {
      if (part.type !== 'literal') values[part.type] = part.value
    }
    const year = values.year || String(d.getUTCFullYear())
    if (!values.month || !values.day || !values.hour || !values.minute || !values.second) {
      throw new Error('incomplete date parts')
    }
    return `${year}-${values.month}-${values.day} ${values.hour}:${values.minute}:${values.second}`
  } catch {
    // 时区名非法等异常情况：回退用浏览器本地时间直接拼接
    const p = (n: number) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
  }
}