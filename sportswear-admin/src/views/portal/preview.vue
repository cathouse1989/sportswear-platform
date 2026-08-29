<template>
  <div class="portal-preview">
    <div class="preview-toolbar">
      <div class="toolbar-left">
        <span class="toolbar-title">门户预览</span>
        <el-tag type="warning" size="small">iframe 真门户 · preview=1</el-tag>
      </div>
      <div class="toolbar-right">
        <el-select v-model="previewLang" size="small" style="width: 120px" @change="refreshPreview">
          <el-option v-for="l in languages" :key="l.code" :label="l.native_name || l.name" :value="l.code" />
        </el-select>
        <el-select v-model="previewPage" size="small" style="width: 160px" @change="onPageChange">
          <el-option v-for="opt in pageOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
        </el-select>
        <el-input
          v-if="needsSlug"
          v-model="pathSlug"
          :placeholder="slugPlaceholder"
          size="small"
          style="width: 160px"
          clearable
          @keyup.enter="refreshPreview"
        />
        <el-radio-group v-model="device" size="small">
          <el-radio-button value="desktop">桌面</el-radio-button>
          <el-radio-button value="tablet">平板</el-radio-button>
          <el-radio-button value="mobile">手机</el-radio-button>
        </el-radio-group>
        <el-button size="small" type="primary" :loading="iframeLoading" @click="refreshPreview">刷新</el-button>
        <el-button size="small" @click="openInNewTab">新窗口打开</el-button>
      </div>
    </div>
    <div v-if="!portalBase" class="preview-hint">
      <el-alert type="warning" show-icon :closable="false" title="未配置门户地址"
        description="请设置 VITE_PORTAL_BASE_URL（例如 http://localhost:3000）并重启 dev server。" />
    </div>
    <div v-else class="preview-frame-wrap">
      <div class="preview-frame-inner" :style="frameStyle">
        <div v-if="iframeLoading" class="preview-loading-mask"><el-skeleton :rows="8" animated /></div>
        <iframe :key="iframeKey" class="preview-iframe" :src="iframeSrc" title="门户预览" @load="onIframeLoad" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { publicApi } from '@/api'
import { ElMessage } from 'element-plus'

const portalBase = (() => { const b = String(import.meta.env.VITE_PORTAL_BASE_URL || 'http://localhost:3000'); return b.endsWith('/') ? b.slice(0, -1) : b })()
const previewLang = ref('en')
const previewPage = ref('home')
const pathSlug = ref('')
const device = ref<'desktop' | 'tablet' | 'mobile'>('desktop')
const languages = ref<Array<{ code: string; name?: string; native_name?: string }>>([])
const iframeLoading = ref(true)
const iframeKey = ref(0)

const pageOptions = [
  { label: '首页', value: 'home' },
  { label: '产品列表', value: 'products' },
  { label: '产品详情', value: 'product-detail' },
  { label: '博客列表', value: 'blogs' },
  { label: '博客详情', value: 'blog-detail' },
  { label: '案例列表', value: 'cases' },
  { label: 'FAQ', value: 'faqs' },
  { label: '关于我们', value: 'about' },
  { label: '联系我们', value: 'contact' },
]

const needsSlug = computed(() => previewPage.value === 'product-detail' || previewPage.value === 'blog-detail')
const slugPlaceholder = computed(() => (previewPage.value === 'blog-detail' ? '输入博客 slug' : '输入产品 slug'))
const deviceWidth: Record<string, string> = { desktop: '100%', tablet: '768px', mobile: '390px' }
const frameStyle = computed(() => ({ width: deviceWidth[device.value] || '100%', maxWidth: '100%' }))

function buildPortalPath(lang: string, page: string, slug: string): string {
  const l = lang || 'en'
  const s = (slug || '').trim()
  switch (page) {
    case 'home': return '/' + l
    case 'products': return '/' + l + '/products'
    case 'product-detail': return s ? '/' + l + '/products/' + encodeURIComponent(s) : '/' + l + '/products'
    case 'blogs': return '/' + l + '/blog'
    case 'blog-detail': return s ? '/' + l + '/blog/' + encodeURIComponent(s) : '/' + l + '/blog'
    case 'cases': return '/' + l + '/cases'
    case 'faqs': return '/' + l + '/faq'
    case 'about': return '/' + l + '/about'
    case 'contact': return '/' + l + '/contact'
    default: return '/' + l
  }
}

const iframeSrc = computed(() => {
  if (!portalBase) return ''
  const path = buildPortalPath(previewLang.value, previewPage.value, pathSlug.value)
  const url = new URL(path, portalBase + '/')
  url.searchParams.set('preview', '1')
  url.searchParams.set('_ts', String(iframeKey.value))
  return url.toString()
})

function refreshPreview() {
  if (needsSlug.value && !pathSlug.value.trim()) {
    ElMessage.info(previewPage.value === 'blog-detail' ? '未填 slug 时将预览博客列表' : '未填 slug 时将预览产品列表')
  }
  iframeLoading.value = true
  iframeKey.value += 1
}
function onPageChange() { pathSlug.value = ''; refreshPreview() }
function onIframeLoad() { iframeLoading.value = false }
function openInNewTab() {
  if (!portalBase) { ElMessage.warning('请先配置 VITE_PORTAL_BASE_URL'); return }
  const path = buildPortalPath(previewLang.value, previewPage.value, pathSlug.value)
  const url = new URL(path, portalBase + '/')
  url.searchParams.set('preview', '1')
  window.open(url.toString(), '_blank', 'noopener,noreferrer')
}
async function loadLanguages() {
  try {
    const result = await publicApi.languages()
    languages.value = result?.length ? result : [
      { code: 'en', name: 'English', native_name: 'English' },
      { code: 'zh', name: '中文', native_name: '中文' },
      { code: 'es', name: 'Español', native_name: 'Español' },
      { code: 'fr', name: 'Français', native_name: 'Français' },
    ]
  } catch {
    languages.value = [
      { code: 'en', name: 'English', native_name: 'English' },
      { code: 'zh', name: '中文', native_name: '中文' },
      { code: 'es', name: 'Español', native_name: 'Español' },
      { code: 'fr', name: 'Français', native_name: 'Français' },
    ]
  }
}
onMounted(async () => { await loadLanguages(); refreshPreview() })
</script>

<style scoped>
.portal-preview { min-height: calc(100vh - 120px); display: flex; flex-direction: column; background: #f3f4f6; }
.preview-toolbar {
  display: flex; justify-content: space-between; align-items: center; gap: 12px; flex-wrap: wrap;
  padding: 12px 20px; background: #fff; border-bottom: 1px solid #e5e7eb; position: sticky; top: 0; z-index: 100;
}
.toolbar-left { display: flex; align-items: center; gap: 12px; }
.toolbar-title { font-size: 16px; font-weight: 600; }
.toolbar-right { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
.preview-hint { padding: 24px; }
.preview-frame-wrap { flex: 1; display: flex; justify-content: center; padding: 16px; min-height: 0; }
.preview-frame-inner {
  position: relative; height: calc(100vh - 180px); min-height: 480px; background: #fff;
  border: 1px solid #e5e7eb; border-radius: 12px; overflow: hidden;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08); transition: width 0.25s ease;
}
.preview-iframe { width: 100%; height: 100%; border: 0; display: block; background: #FBF9F6; }
.preview-loading-mask { position: absolute; inset: 0; z-index: 2; padding: 24px; background: rgba(255,255,255,0.92); }
</style>
