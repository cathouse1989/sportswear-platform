// 枚举字典值域分组读取：后台「字典管理」为唯一配置入口，
// 门户列表页（博客分类/案例类型/FAQ 分类）按 i18n_prefix 动态取值域，避免硬编码值域与字典脱节。
export function useEnumGroupValues(prefix: string) {
  const api = useApi()
  const { data } = useAsyncData<Record<string, string[]>>(
    'enum-groups',
    () =>
      api
        .getEnumGroups()
        .catch(() => ({ groups: {} }))
        .then((r: any) => r?.groups || {}),
  )
  return computed<string[]>(() => data.value?.[prefix] || [])
}
