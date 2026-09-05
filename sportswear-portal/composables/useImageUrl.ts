// 图片 URL 归一化 composable
// 后端 API 返回的 cover_image / image 字段可能是：
//   1. http://localhost:8080/uploads/xxx.jpg  (本地开发)
//   2. http://backend:8080/uploads/xxx.jpg    (Docker 内部)
//   3. https://images.unsplash.com/...         (第三方外链，如 seed 数据)
//   4. /uploads/xxx.jpg                        (相对路径)
//
// 本 composable 将情况 1/2 转为相对路径 /uploads/xxx.jpg，
// 使浏览器通过门户自身的 /uploads 反向代理去获取图片。
// 这样无论门户部署在什么域名/隧道后面，图片都能正常显示。

const BACKEND_PATTERNS = [
  /^https?:\/\/localhost(?::\d+)?\//i,
  /^https?:\/\/127\.0\.0\.1(?::\d+)?\//i,
  /^https?:\/\/backend(?::\d+)?\//i,
  /^https?:\/\/host\.docker\.internal(?::\d+)?\//i,
]

/**
 * 归一化图片 URL：将后端内部绝对地址转为相对路径
 */
export function normalizeImageUrl(url: string | null | undefined): string {
  if (!url) return ''

  // 如果是第三方外链（http/https 但不匹配后端模式），直接返回
  if (url.startsWith('http://') || url.startsWith('https://')) {
    // 检查是否匹配后端内部地址模式
    for (const pattern of BACKEND_PATTERNS) {
      if (pattern.test(url)) {
        // 提取路径部分（去掉 protocol://host:port）
        try {
          const parsed = new URL(url)
          return parsed.pathname + parsed.search
        } catch {
          return url
        }
      }
    }
    // 外部 URL（如 Unsplash），保持原样
    return url
  }

  // 已经是相对路径，直接返回
  return url
}

/**
 * 归一化图片对象数组（用于 product.images / gallery 等）
 */
export function normalizeImages(images: Array<{ url?: string; [key: string]: any }> | undefined | null): Array<{ url?: string; [key: string]: any }> {
  if (!images) return []
  return images.map(img => ({
    ...img,
    url: img.url ? normalizeImageUrl(img.url) : img.url,
  }))
}

/**
 * 在模板中使用的便捷方法
 *
 * 用法：
 * ```vue
 * <img :src="imgUrl(product.cover_image)" :alt="product.sku" />
 * ```
 */
export function imgUrl(url: string | null | undefined): string {
  return normalizeImageUrl(url)
}

/**
 * 取产品展示用图片 URL 列表（封面 + 图集，归一化 + 去重）。
 * 用于产品卡片 hover 切换第二张图、图集数量徽章等。
 */
export function galleryImageUrls(product: any): string[] {
  const list: string[] = []
  if (product?.cover_image) {
    const u = normalizeImageUrl(product.cover_image)
    if (u) list.push(u)
  }
  for (const img of product?.images || []) {
    const u = normalizeImageUrl(img?.url || '')
    if (u && !list.includes(u)) list.push(u)
  }
  return list
}