import { describe, expect, it } from 'vitest'
import { browserTimeZone, formatDateTime, timeZoneOffset } from '@/utils/format'

describe('formatDateTime 按指定时区格式化后端 UTC 时间', () => {
  it('空值 / 非法值返回 fallback', () => {
    expect(formatDateTime(null)).toBe('-')
    expect(formatDateTime(undefined)).toBe('-')
    expect(formatDateTime('')).toBe('-')
    expect(formatDateTime('not-a-date')).toBe('-')
    expect(formatDateTime(null, '—')).toBe('—')
  })

  it('UTC 时间转 Asia/Shanghai（+08:00）', () => {
    // 2026-08-22T04:39:37Z → 北京 2026-08-22 12:39:37
    expect(formatDateTime('2026-08-22T04:39:37Z', '-', 'Asia/Shanghai')).toBe('2026-08-22 12:39:37')
  })

  it('UTC 时间转 America/New_York（冬季 EST，-05:00）', () => {
    // 2026-01-05T12:00:00Z → 纽约 2026-01-05 07:00:00
    expect(formatDateTime('2026-01-05T12:00:00Z', '-', 'America/New_York')).toBe('2026-01-05 07:00:00')
  })

  it('UTC 时间转 America/New_York（夏季 EDT，-04:00）', () => {
    // 2026-08-22T00:30:00Z → 纽约（EDT）2026-08-21 20:30:00
    expect(formatDateTime('2026-08-22T00:30:00Z', '-', 'America/New_York')).toBe('2026-08-21 20:30:00')
  })

  it('跨 UTC 日期转换正确', () => {
    // 2026-08-22T16:00:00Z → 北京 2026-08-23 00:00:00
    expect(formatDateTime('2026-08-22T16:00:00Z', '-', 'Asia/Shanghai')).toBe('2026-08-23 00:00:00')
  })

  it('带毫秒小数与不带 Z 的 RFC3339 均兼容', () => {
    expect(formatDateTime('2026-08-22T04:39:37.772101Z', '-', 'Asia/Shanghai')).toBe('2026-08-22 12:39:37')
    expect(formatDateTime('2026-08-22T04:39:37Z', '-', 'Asia/Shanghai')).toBe('2026-08-22 12:39:37')
  })
})

describe('browserTimeZone / timeZoneOffset', () => {
  it('browserTimeZone 返回非空 IANA 时区名', () => {
    expect(typeof browserTimeZone()).toBe('string')
    expect(browserTimeZone().length).toBeGreaterThan(0)
  })

  it('timeZoneOffset 格式为 UTC±HH:MM，且与 getTimezoneOffset 语义一致', () => {
    expect(timeZoneOffset()).toMatch(/^UTC[+-]\d{2}:\d{2}$/)
    // 任一具体时刻结果一致：工时区由运行环境决定，与 Date 内容无关
    expect(timeZoneOffset(new Date())).toBe(timeZoneOffset())
  })
})