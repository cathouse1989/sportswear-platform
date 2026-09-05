<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-select v-model="entityType" style="width: 200px" @change="loadData">
        <el-option v-for="e in ENTITIES" :key="e.value" :label="e.label" :value="e.value" />
      </el-select>
      <el-select v-model="missingLang" placeholder="筛选缺失语言" clearable style="width: 160px" @change="applyFilter">
        <el-option v-for="l in LANGS" :key="l.value" :label="l.label" :value="l.value" />
      </el-select>
      <div class="spacer" />
      <span class="stats">完整覆盖 {{ coveredCount }} / {{ rows.length }} · 缺失 {{ rows.length - coveredCount }}</span>
    </div>
    <el-table :data="filtered" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" min-width="220" show-overflow-tooltip />
      <el-table-column v-for="l in LANGS" :key="l.value" :label="l.label" width="96" align="center">
        <template #default="{ row }">
          <el-tag :type="row.coverage[l.value] ? 'success' : 'danger'" size="small" effect="plain">{{ row.coverage[l.value] ? '已填' : '缺失' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="缺失语言" min-width="150">
        <template #default="{ row }">
          <el-tag v-for="l in row.missing" :key="l" size="small" type="danger" class="lang-tag">{{ langLabel(l) }}</el-tag>
          <span v-if="!row.missing.length" class="muted">—</span>
        </template>
      </el-table-column>
    </el-table>
    <el-empty v-if="!loading && !filtered.length" description="暂无数据" :image-size="80" />
  </el-card>
</template>
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { productApi, blogApi, caseApi, pageApi, faqApi, factoryApi, certificationApi, productionProcessApi, navigationApi } from '@/api'

const LANGS = [
  { value: 'en', label: 'English' },
  { value: 'zh', label: '中文' },
  { value: 'es', label: 'Español' },
  { value: 'fr', label: 'Français' },
]
function langLabel(v: string) { return LANGS.find(l => l.value === v)?.label || v }

const ENTITIES = [
  { value: 'product', label: '产品', api: productApi.list, nameKey: 'name' },
  { value: 'blog', label: '博客', api: blogApi.list, nameKey: 'title' },
  { value: 'case', label: '案例', api: caseApi.list, nameKey: 'title' },
  { value: 'page', label: '页面', api: pageApi.list, nameKey: 'title' },
  { value: 'faq', label: 'FAQ', api: faqApi.list, nameKey: 'question' },
  { value: 'factory', label: '工厂', api: factoryApi.list, nameKey: 'name' },
  { value: 'certification', label: '认证', api: certificationApi.list, nameKey: 'name' },
  { value: 'process', label: '生产流程', api: productionProcessApi.list, nameKey: 'name' },
  { value: 'navigation', label: '导航', api: navigationApi.listAll, nameKey: 'name' },
]

const TEXT_FIELDS = ['name', 'title', 'question', 'brief', 'description', 'content', 'features', 'usage', 'client_need', 'problem', 'solution', 'process', 'result', 'quality_management', 'production_capability', 'answer']
function hasText(v: any) { return typeof v === 'string' && v.trim().length > 0 }
function transFilled(t: any) { return !!t && TEXT_FIELDS.some(f => hasText(t[f])) }

interface Row { name: string; coverage: Record<string, boolean>; missing: string[] }

const entityType = ref('product')
const missingLang = ref('')
const loading = ref(false)
const rows = ref<Row[]>([])

const coveredCount = computed(() => rows.value.filter(r => !r.missing.length).length)
const filtered = computed(() => {
  if (!missingLang.value) return rows.value
  return rows.value.filter(r => r.coverage[missingLang.value] === false)
})
function applyFilter() { /* computed 自动响应 */ }

async function loadData() {
  loading.value = true
  const ent = ENTITIES.find(e => e.value === entityType.value)!
  try {
    const res: any = await ent.api({ page: 1, pageSize: 500 })
    const list: any[] = Array.isArray(res) ? res : (res.items || res || [])
    const flat = flattenNav(list)
    rows.value = flat.map((item: any) => {
      const cov: Record<string, boolean> = {
        en: hasText(item[ent.nameKey]),
        zh: transFilled((item.translations || []).find((t: any) => t.language === 'zh')),
        es: transFilled((item.translations || []).find((t: any) => t.language === 'es')),
        fr: transFilled((item.translations || []).find((t: any) => t.language === 'fr')),
      }
      const missing = LANGS.filter(l => !cov[l.value]).map(l => l.value)
      return { name: item[ent.nameKey] || item.sku || '-', coverage: cov, missing }
    })
  } catch { rows.value = [] } finally { loading.value = false }
}
function flattenNav(items: any[]): any[] {
  const out: any[] = []
  const walk = (arr: any[]) => {
    for (const n of arr || []) {
      out.push(n)
      if (n.children?.length) walk(n.children)
    }
  }
  walk(items)
  return out
}
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.stats { font-size: 13px; color: #606266; }
.lang-tag { margin-right: 4px; }
.muted { color: #c0c4cc; }
</style>