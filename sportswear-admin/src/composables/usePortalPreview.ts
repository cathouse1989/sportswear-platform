import { ElMessage } from 'element-plus'

/**
 * 门户预览工具：在内容编辑弹窗内一键打开门户对应页面（新标签页，preview=1 绕过缓存直查 DB）。
 * 门户地址通过环境变量 VITE_PORTAL_BASE_URL 配置，默认 http://localhost:3000。
 */
export function usePortalPreview() {
  const portalBase = (() => {
    const b = String(import.meta.env.VITE_PORTAL_BASE_URL || 'http://localhost:3000')
    return b.endsWith('/') ? b.slice(0, -1) : b
  })()

  function buildPath(page: string, slug = '', lang = 'en'): string {
    const l = lang || 'en'
    const s = (slug || '').trim()
    switch (page) {
      case 'home': return '/' + l
      case 'products': return '/' + l + '/products'
      case 'product': return s ? '/' + l + '/products/' + encodeURIComponent(s) : '/' + l + '/products'
      case 'blog': return s ? '/' + l + '/blog/' + encodeURIComponent(s) : '/' + l + '/blog'
      case 'case': return s ? '/' + l + '/cases/' + encodeURIComponent(s) : '/' + l + '/cases'
      case 'faq': return '/' + l + '/faq'
      case 'about': return '/' + l + '/about'
      case 'contact': return '/' + l + '/contact'
      default: return '/' + l
    }
  }

  function openPreview(page: string, slug = '', lang = 'en') {
    const path = buildPath(page, slug, lang)
    const url = new URL(path, portalBase + '/')
    url.searchParams.set('preview', '1')
    window.open(url.toString(), '_blank', 'noopener,noreferrer')
  }

  function previewWithSlug(page: string, slug: string, lang = 'en') {
    if (!slug?.trim()) { ElMessage.warning('请先填写 slug'); return }
    openPreview(page, slug, lang)
  }

  return { openPreview, previewWithSlug }
}