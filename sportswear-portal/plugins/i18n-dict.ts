// 门户动态词条（后台「词条管理」i18n_entries → GET /public/i18n）合并进 vue-i18n：
// - 运营在后台按语言（en/zh/es/fr）配置的文案对门户所有 $t 文案即时生效（覆盖静态语言包）；
// - 后台未配置的 Key 继续走静态语言包（locales/*.json）兜底，不会丢文案；
// - 词条 value 支持换行：统一 \r\n → \n，门户模板配合 whitespace-pre-line 渲染多行。
export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig()

  /**
   * vue-i18n 消息格式特殊字符转义：
   * `@` 是"链接消息"前缀、`|` 是复数分隔符，纯文本文案（如邮箱 cathouse1989@gmail.com）
   * 含裸 `@` 会触发 "Message compilation error: Invalid linked format"，
   * 生产环境 SSR 直接 500。后台词条一律按纯文本维护，故在合并进 vue-i18n 前统一转义为字面量；
   * 已按 vue-i18n 语法转义过的 {'x'} 先还原再转义，避免二次转义。
   */
  const escapeI18nSpecials = (value: string) =>
    value
      .replace(/\{'(.)'\}/g, '$1')
      .replace(/[@|]/g, (ch) => `{'${ch}'}`)

  /** 将 "a.b.c" 扁平 key 还原为嵌套结构，保证 vue-i18n 可解析；换行符归一化 */
  const nestDict = (dict: Record<string, string>) => {
    const messages: Record<string, any> = {}
    for (const [key, value] of Object.entries(dict)) {
      if (typeof value !== 'string' || !value.trim()) continue
      const parts = key.split('.')
      let node = messages
      for (let i = 0; i < parts.length - 1; i++) {
        node = node[parts[i]] ||= {}
      }
      node[parts[parts.length - 1]] = escapeI18nSpecials(value.replace(/\r\n?/g, '\n'))
    }
    return messages
  }

  const applyDict = async () => {
    const i18n: any = (nuxtApp as any).$i18n
    if (!i18n?.mergeLocaleMessage) return
    const lang = (typeof i18n.locale === 'string' ? i18n.locale : i18n.locale?.value) || 'en'
    // SSR 走容器网络内的 apiServer；浏览器走公共 API 地址
    const base = import.meta.server
      ? config.apiServer || config.public.apiBase
      : config.public.apiBase
    try {
      const res = await $fetch<{ success: boolean; data: { dictionary?: Record<string, string> } }>(
        `${base}/public/i18n`,
        { params: { lang } },
      )
      const messages = nestDict(res?.data?.dictionary || {})
      if (Object.keys(messages).length) i18n.mergeLocaleMessage(lang, messages)
    } catch {
      // 后端不可用时忽略，静态语言包兜底
    }
  }

  // 首次加载（SSR/客户端水合前）拉取当前语言词条；SPA 内切换语言时重新拉取
  nuxtApp.hook('app:created', async () => {
    await applyDict()
  })
  nuxtApp.hook('i18n:localeSwitched', () => {
    applyDict()
  })
})
