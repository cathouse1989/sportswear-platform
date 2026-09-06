<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-radio-group v-model="navType" @change="loadData">
        <el-radio-button v-for="o in enumOptions('navigation.type')" :key="o.value" :value="o.value">{{ o.label }}</el-radio-button>
      </el-radio-group>
      <el-input v-model="searchKeyword" placeholder="搜索导航名称" clearable style="width:200px" />
      <el-select v-model="statusFilter" placeholder="页面状态" clearable style="width:120px" @change="loadData">
        <el-option label="全部" value="" />
        <el-option label="已发布" value="published" />
        <el-option label="未发布" value="draft" />
        <el-option label="无关联页面" value="no_page" />
      </el-select>
      <div class="spacer" />
    </div>

    <el-table
      :data="filteredNavigations"
      v-loading="loading"
      row-key="id"
      stripe
      default-expand-all
      :tree-props="{ children: 'children' }"
    >
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="url" label="URL" min-width="150" />
      <el-table-column label="关联页面" width="180">
        <template #default="{ row }: any">
          <div v-if="row.page_id" class="page-info">
            <el-tag size="small" type="success">{{ row.page?.title || '已关联' }}</el-tag>
            <el-tag v-if="row.page_status" size="small" :type="enumTag('content.status', row.page_status)" class="page-status-tag">
              {{ enumLabel('content.status', row.page_status) }}
            </el-tag>
          </div>
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column label="SEO" width="80">
        <template #default="{ row }: any">
          <span v-if="isExternalUrl(row.url)" class="muted">—</span>
          <el-tag v-else-if="seoStatusMap[row.url]" size="small" type="success">已配置</el-tag>
          <el-tag v-else size="small" type="warning">未配置</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="target" label="打开方式" width="100">
        <template #default="{ row }: any">
          <el-tag size="small" :type="row.target === '_blank' ? 'warning' : 'info'">{{ row.target === '_blank' ? '新窗口' : '当前页' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="可见" width="70">
        <template #default="{ row }: any">
          <el-switch v-model="row.is_visible" size="small" @change="(val: string | number | boolean) => toggleVisible(row, Boolean(val))" />
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column label="操作" width="220">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" @click="handleMoveUp(row)" :disabled="isFirstSibling(row)">↑</el-button>
          <el-button size="small" @click="handleMoveDown(row)" :disabled="isLastSibling(row)">↓</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑导航' : '新建导航'" width="620px" :close-on-click-modal="false">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="如: OEM 服务" />
        </el-form-item>
        <el-form-item label="关联页面">
          <el-select v-model="form.page_id" placeholder="选择关联页面（可选）" clearable filterable style="width:100%"
            @change="onPageSelect">
            <el-option v-for="p in publishedPages" :key="p.id" :label="`${p.title} (${p.type})`" :value="p.id" />
          </el-select>
          <div class="form-tip">选择已发布页面后，名称和 URL 将自动填充，也可手动修改。</div>
        </el-form-item>
        <el-form-item label="URL">
          <el-input v-model="form.url" placeholder="如: /oem 或 https://..." />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type">
            <el-option v-for="o in enumOptions('navigation.type')" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="打开方式">
          <el-select v-model="form.target">
            <el-option label="当前页" value="_self" />
            <el-option label="新窗口" value="_blank" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="显示">
          <el-switch v-model="form.is_visible" />
        </el-form-item>
        <el-divider content-position="left">多语言翻译（未填写回退英文）</el-divider>
        <TransEditor v-model="translations" :fields="TRANS_FIELDS" :source="{ name: form.name }" :langs="TRANS_LANGS" source-label="English" />
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { navigationApi, pageApi, portalHealthApi } from '@/api'
import TransEditor from '@/components/cms/TransEditor.vue'
import { DEFAULT_TRANS_LANGS, createEmptyTranslations, translationsToRecord, translationsToPayload } from '@/composables/useTransRecord'
import type { Navigation, Page } from '@/types'
import { useEnumDict } from '@/composables/useEnumDict'
import { pagePath } from '@/constants/routes'

// 页面状态标签由「字典管理」驱动（content.status）
const { ensureLoaded: loadEnumDict, options: enumOptions, label: enumLabel, tagType: enumTag } = useEnumDict()

const navigations = ref<Navigation[]>([])
const loading = ref(false)
const navType = ref('header')
const searchKeyword = ref('')
const statusFilter = ref('')

// 用于搜索过滤的展平数据
const filteredNavigations = computed(() => {
  const kw = searchKeyword.value.toLowerCase().trim()
  const status = statusFilter.value
  
  // 先按状态过滤
  let items = navigations.value
  if (status) {
    items = items.filter(n => {
      if (status === 'no_page') return !n.page_id
      return n.page_status === status
    })
  }
  
  if (!kw) return items
  // 递归过滤：保留匹配的节点及其祖先
  const result: Navigation[] = []
  function collect(items: Navigation[], ancestors: Navigation[]): boolean {
    let hasMatch = false
    for (const item of items) {
      const path = [...ancestors, item]
      const childrenMatch = item.children?.length ? collect(item.children, path) : false
      if (item.name.toLowerCase().includes(kw) || childrenMatch) {
        // 确保祖先都在结果中
        for (const a of path) {
          if (!result.find(r => r.id === a.id)) result.push(a)
        }
        hasMatch = true
      }
    }
    return hasMatch
  }
  for (const item of items) {
    collect([item], [])
  }
  return result
})

async function loadData() {
  loading.value = true
  try {
    const result = await navigationApi.list({ type: navType.value })
    navigations.value = (result || []).map(n => ({
      ...n,
      page_status: n.page?.status || (n.page_id ? 'unknown' : undefined)
    })) as Navigation[]
  } catch { navigations.value = [] }
  finally { loading.value = false }
}

const publishedPages = ref<Page[]>([])
async function loadPublishedPages() {
  try {
    const result: any = await pageApi.list({ pageSize: 500, status: 'published' })
    publishedPages.value = (result?.items || result || []) as Page[]
  } catch { publishedPages.value = [] }
}

// 路由 path → 是否已配置 SEO（用于导航列表 SEO 状态列）
const seoStatusMap = ref<Record<string, boolean>>({})
async function loadRouteSeoStatus() {
  try {
    const list: any[] = await portalHealthApi.listRoutes()
    const map: Record<string, boolean> = {}
    for (const r of list || []) {
      if (r.path) map[r.path] = !!r.seo_configured
    }
    seoStatusMap.value = map
  } catch { seoStatusMap.value = {} }
}

function isExternalUrl(url: string): boolean {
  if (!url) return true
  return /^(https?:|mailto:|tel:|#)/i.test(url)
}

function onPageSelect(pageId: string) {
  if (!pageId) return
  const p = publishedPages.value.find(p => p.id === pageId)
  if (!p) return
  if (!form.name) form.name = p.title
  if (!form.url) form.url = pagePath(p.type, p.slug)
}

const dialogVisible = ref(false)
const editingId = ref('')
const parentId = ref<string | null>(null)
const saving = ref(false)
const form = reactive({
  name: '', url: '', type: 'header', target: '_self',
  sort_order: 0, is_visible: true, page_id: '',
})
const TRANS_LANGS = DEFAULT_TRANS_LANGS
const TRANS_FIELDS = [{ key: 'name', label: '名称' }]
const TRANS_KEYS = TRANS_FIELDS.map((f) => f.key)
const translations = ref<Record<string, Record<string, string>>>({})

function resetForm() {
  editingId.value = ''
  parentId.value = null
  Object.assign(form, { name: '', url: '', type: navType.value, target: '_self', sort_order: 0, is_visible: true, page_id: '' })
  translations.value = createEmptyTranslations(TRANS_KEYS)
}

async function openEditDialog(row: Navigation) {
  resetForm()
  editingId.value = row.id
  Object.assign(form, {
    name: row.name, url: row.url, type: row.type, target: row.target || '_self',
    sort_order: row.sort_order, is_visible: row.is_visible, page_id: row.page_id || '',
  })
  translations.value = translationsToRecord(row.translations || [], TRANS_KEYS)
  dialogVisible.value = true
}

async function handleSave() {
  if (!form.name.trim()) { ElMessage.warning('请输入名称'); return }
  if (!form.url.trim()) { ElMessage.warning('请输入 URL'); return }
  saving.value = true
  const payload: any = {
    name: form.name, url: form.url, type: form.type, target: form.target,
    sort_order: form.sort_order, is_visible: form.is_visible,
    page_id: form.page_id || null,
    translations: translationsToPayload(translations.value, TRANS_KEYS, TRANS_LANGS),
  }
  if (parentId.value) payload.parent_id = parentId.value
  try {
    if (editingId.value) {
      await navigationApi.update(editingId.value, payload)
      ElMessage.success('更新成功')
    } else {
      await navigationApi.create(payload)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch { /* handled */ }
  finally { saving.value = false }
}

function getSiblings(row: Navigation): Navigation[] {
  if (row.parent_id) {
    for (const n of navigations.value) {
      if (n.children?.length) {
        const sibs = n.children.filter(c => c.parent_id === row.parent_id)
        if (sibs.length) return sibs
      }
    }
  }
  return navigations.value.filter(n => !n.parent_id)
}

function isFirstSibling(row: Navigation): boolean {
  const sibs = getSiblings(row)
  return sibs.length > 0 && sibs[0].id === row.id
}

function isLastSibling(row: Navigation): boolean {
  const sibs = getSiblings(row)
  return sibs.length > 0 && sibs[sibs.length - 1].id === row.id
}

async function handleMoveUp(row: Navigation) {
  const sibs = getSiblings(row)
  const idx = sibs.findIndex(s => s.id === row.id)
  if (idx <= 0) return
  const items = [
    { id: sibs[idx].id, sort_order: sibs[idx - 1].sort_order },
    { id: sibs[idx - 1].id, sort_order: sibs[idx].sort_order },
  ]
  try {
    await navigationApi.batchSort(items)
    loadData()
  } catch { }
}

async function handleMoveDown(row: Navigation) {
  const sibs = getSiblings(row)
  const idx = sibs.findIndex(s => s.id === row.id)
  if (idx < 0 || idx >= sibs.length - 1) return
  const items = [
    { id: sibs[idx].id, sort_order: sibs[idx + 1].sort_order },
    { id: sibs[idx + 1].id, sort_order: sibs[idx].sort_order },
  ]
  try {
    await navigationApi.batchSort(items)
    loadData()
  } catch { }
}

async function toggleVisible(row: Navigation, val: boolean) {
  try {
    await navigationApi.update(row.id, { name: row.name, url: row.url, type: row.type, target: row.target || '_self', sort_order: row.sort_order, is_visible: val, page_id: row.page_id || null })
  } catch { row.is_visible = !val }
}

onMounted(() => {
  loadEnumDict()
  loadData()
  loadPublishedPages()
  loadRouteSeoStatus()
})
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.muted { color: #c0c4cc; }
.form-tip { margin-top: 4px; font-size: 12px; color: #909399; }
.page-info { display: flex; flex-direction: column; gap: 4px; }
.page-status-tag { align-self: flex-start; }
</style>