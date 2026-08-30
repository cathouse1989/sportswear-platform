<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-radio-group v-model="navType" @change="loadData">
        <el-radio-button value="header">Header 导航</el-radio-button>
        <el-radio-button value="footer">Footer 导航</el-radio-button>
      </el-radio-group>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog(null)">新建导航</el-button>
      <el-button type="success" plain @click="handleSync">一键同步已发布页面</el-button>
    </div>

    <el-table
      :data="treeData"
      v-loading="loading"
      row-key="id"
      stripe
      default-expand-all
      :tree-props="{ children: 'children' }"
    >
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="url" label="URL" min-width="150" />
      <el-table-column label="关联页面" width="130">
        <template #default="{ row }: any">
          <el-tag v-if="row.page_id" size="small" type="success">{{ row.page?.title || '已关联' }}</el-tag>
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column prop="target" label="打开方式" width="100">
        <template #default="{ row }: any">
          <el-tag size="small" :type="row.target === '_blank' ? 'warning' : 'info'">{{ row.target === '_blank' ? '新窗口' : '当前页' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="可见" width="70">
        <template #default="{ row }: any">
          <el-switch v-model="row.is_visible" size="small" @change="(val: boolean) => toggleVisible(row, val)" />
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column label="操作" width="280">
        <template #default="{ row }: any">
          <el-button size="small" @click="openCreateDialog(row.id)">添加子级</el-button>
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="warning" @click="handleMoveUp(row)">上移</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
            <el-option label="Header" value="header" />
            <el-option label="Footer" value="footer" />
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
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { navigationApi, pageApi } from '@/api'
import type { Navigation, Page } from '@/types'

const navigations = ref<Navigation[]>([])
const loading = ref(false)
const navType = ref('header')

const treeData = ref<Navigation[]>([])

function flattenTree(items: Navigation[], depth = 0): Navigation[] {
  const result: Navigation[] = []
  for (const item of items) {
    result.push(item)
    if (item.children?.length) result.push(...flattenTree(item.children, depth + 1))
  }
  return result
}

async function loadData() {
  loading.value = true
  try {
    const result = await navigationApi.list({ type: navType.value })
    navigations.value = result || []
    treeData.value = flattenTree(navigations.value)
  } catch { navigations.value = []; treeData.value = [] }
  finally { loading.value = false }
}

const publishedPages = ref<Page[]>([])
async function loadPublishedPages() {
  try {
    const result: any = await pageApi.list({ pageSize: 500, status: 'published' })
    publishedPages.value = (result?.items || result || []) as Page[]
  } catch { publishedPages.value = [] }
}

const PAGE_TYPE_URL_MAP: Record<string, string> = {
  home: '/', product: '/products', product_category: '/products',
  oem: '/oem', odm: '/odm', private_label: '/private-label',
  factory: '/factory', production: '/production',
  blog: '/blog', case: '/cases', faq: '/faq',
  contact: '/contact', seo_landing: '/landing',
}

function onPageSelect(pageId: string) {
  if (!pageId) return
  const p = publishedPages.value.find(p => p.id === pageId)
  if (!p) return
  if (!form.name) form.name = p.title
  if (!form.url) form.url = PAGE_TYPE_URL_MAP[p.type as string] || '/' + p.slug
}

const dialogVisible = ref(false)
const editingId = ref('')
const parentId = ref<string | null>(null)
const saving = ref(false)
const form = reactive({
  name: '', url: '', type: 'header', target: '_self',
  sort_order: 0, is_visible: true, page_id: '',
})

function resetForm() {
  editingId.value = ''
  parentId.value = null
  Object.assign(form, { name: '', url: '', type: navType.value, target: '_self', sort_order: 0, is_visible: true, page_id: '' })
}

function openCreateDialog(pid: string | null) {
  resetForm()
  parentId.value = pid
  if (pid) {
    for (const n of navigations.value) {
      if (n.id === pid) { form.type = n.type; break }
      if (n.children?.length) {
        const child = n.children.find(c => c.id === pid)
        if (child) { form.type = child.type; break }
      }
    }
  }
  dialogVisible.value = true
}

async function openEditDialog(row: Navigation) {
  resetForm()
  editingId.value = row.id
  Object.assign(form, {
    name: row.name, url: row.url, type: row.type, target: row.target || '_self',
    sort_order: row.sort_order, is_visible: row.is_visible, page_id: row.page_id || '',
  })
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

async function toggleVisible(row: Navigation, val: boolean) {
  try {
    await navigationApi.update(row.id, { name: row.name, url: row.url, type: row.type, target: row.target || '_self', sort_order: row.sort_order, is_visible: val, page_id: row.page_id || null })
  } catch { row.is_visible = !val }
}

async function handleMoveUp(row: Navigation) {
  try {
    await navigationApi.update(row.id, { name: row.name, url: row.url, type: row.type, target: row.target || '_self', sort_order: Math.max(0, row.sort_order - 1), is_visible: row.is_visible, page_id: row.page_id || null })
    loadData()
  } catch { }
}

async function handleDelete(row: Navigation) {
  await ElMessageBox.confirm(`确定删除导航「${row.name}」吗？`, '警告', { type: 'warning' })
  await navigationApi.delete(row.id)
  ElMessage.success('已删除')
  loadData()
}

async function handleSync() {
  await ElMessageBox.confirm(
    '将扫描所有「已发布」页面，为尚未关联导航的页面自动创建导航条目。已有关联不会被覆盖。',
    '一键同步', { type: 'info' }
  )
  try {
    const allPages: any = await pageApi.list({ pageSize: 500, status: 'published' })
    const pages: Page[] = (allPages?.items || allPages || []) as Page[]
    const allNavs = await navigationApi.listAll()
    const usedIds = new Set((allNavs || []).filter((n: Navigation) => n.page_id).map((n: Navigation) => n.page_id))
    let created = 0
    for (const p of pages) {
      if (usedIds.has(p.id)) continue
      const url = PAGE_TYPE_URL_MAP[p.type as string] || ''
      if (!url) continue
      await navigationApi.create({
        name: p.title, type: navType.value, url,
        target: '_self', sort_order: 99, is_visible: true, page_id: p.id,
      })
      created++
    }
    ElMessage.success(`已同步创建 ${created} 条导航`)
    loadData()
  } catch { }
}

onMounted(() => {
  loadData()
  loadPublishedPages()
})
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.muted { color: #c0c4cc; }
.form-tip { margin-top: 4px; font-size: 12px; color: #909399; }
</style>