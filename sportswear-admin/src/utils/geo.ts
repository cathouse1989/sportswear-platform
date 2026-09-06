/**
 * 国家/地区代码 → 中文名称映射（用于展示 IP 解析出的国家 ISO2 代码）。
 * 与后端 internal/services/geo_service.go 的 countryLocaleMap 保持一致，
 * 未收录的代码原样返回（如 "US" → "美国"，"XX" → "XX"）。
 */
const COUNTRY_NAMES_CN: Record<string, string> = {
  US: '美国',
  GB: '英国',
  CN: '中国',
  HK: '中国香港',
  TW: '中国台湾',
  JP: '日本',
  KR: '韩国',
  DE: '德国',
  FR: '法国',
  ES: '西班牙',
  IT: '意大利',
  NL: '荷兰',
  RU: '俄罗斯',
  AE: '阿联酋',
  SA: '沙特阿拉伯',
  IN: '印度',
  AU: '澳大利亚',
  NZ: '新西兰',
  CA: '加拿大',
  BR: '巴西',
  MX: '墨西哥',
  ZA: '南非',
  SG: '新加坡',
  MY: '马来西亚',
  TH: '泰国',
  VN: '越南',
  ID: '印度尼西亚',
  PH: '菲律宾',
  TR: '土耳其',
  PL: '波兰',
}

/** 把 ISO2 国家代码转成中文名；空值返回空串，未知代码原样返回。 */
export function countryName(code?: string | null): string {
  if (!code) return ''
  const c = code.trim().toUpperCase()
  return COUNTRY_NAMES_CN[c] || code
}
