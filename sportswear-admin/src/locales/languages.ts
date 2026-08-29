// 支持的语言列表 — 与门户前端 /public/languages 保持一致
export const supportedLanguages = [
  { code: 'en', name: 'English', native_name: 'English' },
  { code: 'zh', name: 'Chinese', native_name: '中文' },
  { code: 'es', name: 'Español', native_name: 'Español' },
  { code: 'fr', name: 'Français', native_name: 'Français' },
] as const

export type SupportedLang = (typeof supportedLanguages)[number]['code']

export const fallbackLanguage: SupportedLang = 'en'
