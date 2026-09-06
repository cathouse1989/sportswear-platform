<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-radio-group v-model="mode" @change="onModeChange">
        <el-radio-button value="route">路由</el-radio-button>
        <el-radio-button value="entity">实体详情</el-radio-button>
      </el-radio-group>
      <template v-if="mode === 'route'">
        <el-select v-model="route" placeholder="选择路由" style="width: 240px" @change="loadData">
          <el-option v-for="r in ROUTES" :key="r.value" :label="r.label" :value="r.value" />
        </el-select>
      </template>
      <template v-else>
        <el-select v-model="entityType" placeholder="实体类型" style="width: 140px" @change="onEntityTypeChange">
          <el-option v-for="t in ENTITY_TYPES" :key="t.value" :label="t.label" :value="t.value" />
        </el-select>
        <el-select v-model="entityId" placeholder="选择实体" style="width: 260px" filterable @change="loadData">
          <el-option v-for="e in entities" :key="e.id" :label="e.label" :value="e.id" />
        </el-select>
      </template>
      <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
    </div>

    <el-alert
      v-if="(mode === 'route' && !route) || (mode === 'entity' && !entityId)"
      type="info"
      :closable="false"
      show-icon
      :title="mode === 'route'
        ? '请选择要配置 SEO 的路由（列表页/落地页）。英文为源语言，其他语言未填写时门户回退英文，再由页面默认值兜底。'
        : '请选择要配置 SEO 的实体（产品/博客/案例/分类详情）。英文为源语言，其他语言未填写时门户回退英文。'"
      class="mb"
    />

    <div v-loading="loading">
      <template v-if="(mode === 'route' && route) || (mode === 'entity' && entityId)">
        <el-divider content-position="left">源语言（English）</el-divider>
        <el-form label-width="110px">
          <el-form-item label="SEO 标题"><el-input v-model="form.en.title" maxlength="200" placeholder="如 Custom Sportswear Products" /></el-form-item>
          <el-form-item label="SEO 描述"><el-input v-model="form.en.description" type="textarea" :rows="2" maxlength="200" show-word-limit /></el-form-item>
          <el-form-item label="关键词"><el-input v-model="form.en.keywords" placeholder="逗号分隔" /></el-form-item>
          <el-form-item label="OG 标题"><el-input v-model="form.en.og_title" maxlength="200" /></el-form-item>
          <el-form-item label="OG 描述"><el-input v-model="form.en.og_description" type="textarea" :rows="2" /></el-form-item>
          <el-form-item label="OG 图片"><el-input v-model="form.en.og_image" placeholder="https://..." /></el-form-item>
        </el-form>

        <el-divider content-position="left">搜索结果预览（SERP）</el-divider>
        <SerpPreview :title="form.en.title" :description="form.en.description" :url="serpUrl" />

        <el-divider content-position="left">其他语言（未填写回退英文）</el-divider>
        <TransEditor
          v-model="transForm"
          :fields="FIELDS"
          :source="form.en"
          :langs="TRANS_LANGS"
          source-label="English"
        />
      </template>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { seoApi, portalHealthApi, productApi, blogApi, caseApi, categoryApi } from '@/api'
import { useRoute } from 'vue-router'
import TransEditor from '@/components/cms/TransEditor.vue'
import SerpPreview from '@/components/cms/SerpPreview.vue'

// 静态兜底路由（含历史 categories/series，保证后端不可用时下拉仍可用）
const FALLBACK_ROUTES = [
  { label: '首页', value: 'home' },
  { label: '产品列表', value: 'products' },
  { label: '博客列表', value: 'blog' },
  { label: '案例列表', value: 'cases' },
  { label: 'FAQ', value: 'faq' },
  { label: '联系我们', value: 'contact' },
  { label: '关于我们', value: 'about' },
  { label: '分类列表', value: 'categories' },
  { label: '系列列表', value: 'series' },
  { label: '隐私政策', value: 'privacy' },
]

// 动态路由下拉：后端路由注册表（预置 + 已发布页面派生），新建页面后自动出现
const ROUTES = ref<{ label: string; value: string }[]>(FALLBACK_ROUTES)

async function loadRoutes() {
  try {
    const list: any[] = await portalHealthApi.listRoutes()
    if (!list || !list.length) return
    const dynamic = list.filter((r) => r.route).map((r) => ({ value: r.route, label: r.label || r.route }))
    const dynamicValues = new Set(dynamic.map((d) => d.value))
    const legacy = FALLBACK_ROUTES.filter((f) => !dynamicValues.has(f.value))
    ROUTES.value = [...legacy, ...dynamic]
  } catch { /* 保留兜底 */ }
}

// 实体详情维度：类型 + 实体下拉
const ENTITY_TYPES = [
  { label: '产品', value: 'product' },
  { label: '博客', value: 'blog' },
  { label: '案例', value: 'case' },
  { label: '分类', value: 'category' },
]

const mode = ref<'route' | 'entity'>('route')
const entityType = ref('product')
const entityId = ref('')
const entities = ref<{ id: string; label: string }[]>([])

async function loadEntities() {
  entities.value = []
  entityId.value = ''
  if (mode.value !== 'entity') return
  try {
    if (entityType.value === 'product') {
      const res: any = await productApi.list({ page: 1, pageSize: 500 })
      entities.value = (res?.items || res || []).map((p: any) => ({ id: p.id, label: p.name || p.slug || p.id }))
    } else if (entityType.value === 'blog') {
      const res: any = await blogApi.list({ page: 1, pageSize: 500 })
      entities.value = (res?.items || res || []).map((b: any) => ({ id: b.id, label: b.title || b.slug || b.id }))
    } else if (entityType.value === 'case') {
      const res: any = await caseApi.list({ page: 1, pageSize: 500 })
      entities.value = (res?.items || res || []).map((c: any) => ({ id: c.id, label: c.title || c.slug || c.id }))
    } else if (entityType.value === 'category') {
      const res: any = await categoryApi.list()
      entities.value = (res || []).map((c: any) => ({ id: c.id, label: c.name || c.slug || c.id }))
    }
  } catch { entities.value = [] }
}

function onModeChange() {
  route.value = ''
  entityId.value = ''
  reset()
  if (mode.value === 'entity') loadEntities()
}

function onEntityTypeChange() {
  entityId.value = ''
  reset()
  loadEntities()
}

const TRANS_LANGS = [
  { value: 'zh', label: '中文' },
  { value: 'es', label: 'Español' },
  { value: 'fr', label: 'Français' },
]

const FIELDS = [
  { key: 'title', label: 'SEO 标题' },
  { key: 'description', label: 'SEO 描述', type: 'textarea' as const, rows: 2 },
  { key: 'keywords', label: '关键词' },
  { key: 'og_title', label: 'OG 标题' },
  { key: 'og_description', label: 'OG 描述', type: 'textarea' as const, rows: 2 },
  { key: 'og_image', label: 'OG 图片' },
]

const blank = () => ({ title: '', description: '', keywords: '', og_title: '', og_description: '', og_image: '' })

const route = ref('')
const currentRoute = useRoute()
const serpUrl = computed(() => {
  if (mode.value === 'entity') return 'https://sportswear-platform.com/' + entityType.value + '/' + entityId.value
  return 'https://sportswear-platform.com/' + (route.value || '')
})
const loading = ref(false)
const saving = ref(false)
const form = reactive<Record<string, Record<string, string>>>({ en: blank(), zh: blank(), es: blank(), fr: blank() })
const transForm = ref<Record<string, Record<string, string>>>({})

function applyItem(it: any) {
  return { title: it.title || '', description: it.description || '', keywords: it.keywords || '', og_title: it.og_title || '', og_description: it.og_description || '', og_image: it.og_image || '' }
}
function reset() {
  form.en = blank(); form.zh = blank(); form.es = blank(); form.fr = blank()
  transForm.value = { zh: { ...blank() }, es: { ...blank() }, fr: { ...blank() } }
}

async function loadData() {
  reset()
  if (mode.value === 'route' && !route.value) return
  if (mode.value === 'entity' && (!entityType.value || !entityId.value)) return
  loading.value = true
  try {
    const res = mode.value === 'route'
      ? await seoApi.listRoutes(route.value)
      : await seoApi.listEntities(entityType.value, entityId.value)
    for (const it of (res as any).items || []) {
      if (form[it.language]) form[it.language] = applyItem(it)
    }
    transForm.value = { zh: { ...form.zh }, es: { ...form.es }, fr: { ...form.fr } }
  } catch { /* 忽略 */ } finally { loading.value = false }
}

async function handleSave() {
  if (mode.value === 'route' && !route.value) { ElMessage.warning('请先选择路由'); return }
  if (mode.value === 'entity' && (!entityType.value || !entityId.value)) { ElMessage.warning('请先选择实体'); return }
  for (const l of ['zh', 'es', 'fr']) form[l] = { ...(transForm.value[l] || blank()) }
  const entries = ['en', 'zh', 'es', 'fr'].map((l) => ({ language: l, ...form[l] }))
  saving.value = true
  try {
    if (mode.value === 'route') {
      await seoApi.saveRoutes(route.value, entries)
    } else {
      await seoApi.saveEntities(entityType.value, entityId.value, entries)
    }
    ElMessage.success('已保存')
    loadData()
  } catch { /* handled */ } finally { saving.value = false }
}

onMounted(() => {
  loadRoutes()
  const q = currentRoute.query.route
  if (typeof q === 'string' && q) { route.value = q; loadData() }
})
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.mb { margin-bottom: 16px; }
</style>