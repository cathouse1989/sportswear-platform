// 分类 / 系列名称的本地化展示：
// 优先取 i18n 词条（product.categories.<slug> / product.series.<slug>），
// 未配置时回退后端返回的英文原文，保证自定义数据不丢文案。
export function useLocalized() {
  const { t, te } = useI18n()

  function localizedCategory(name: string, slug?: string): string {
    if (slug && te(`product.categories.${slug}`)) {
      return t(`product.categories.${slug}`)
    }
    return name
  }

  function localizedSeries(name: string, slug?: string): string {
    if (slug && te(`product.series.${slug}`)) {
      return t(`product.series.${slug}`)
    }
    return name
  }

  return { localizedCategory, localizedSeries }
}
