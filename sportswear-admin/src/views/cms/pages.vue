<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input
        v-model="keyword"
        placeholder="搜索标题 / Slug"
        clearable
        style="width: 240px"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="status" placeholder="状态" clearable style="width: 150px" @change="handleSearch">
        <el-option label="草稿" value="draft" />
        <el-option label="审核中" value="review" />
        <el-option label="已发布" value="published" />
        <el-option label="已下线" value="offline" />
        <el-option label="已归档" value="archived" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建页面</el-button>
    </div>

    <el-table :data="pages" v-loading="loading" stripe>
      <el-table-column prop="title" label="标题" min-width="150" />
      <el-table-column prop="slug" label="Slug" min-width="120" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }"><el-tag size="small">{{ row.type }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)" size="small">{{ row.status || 'draft' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="翻译" min-width="130">
        <template #default="{ row }">
          <template v-if="(row.translations || []).length">
            <el-tag v-for="t in row.translations" :key="t.id" size="small" type="info" class="lang-tag">{{ t.language }}</el-tag>
          </template>
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column label="导航状态" width="120">
        <template #default="{ row }: any">
          <template v-if="navStatusMap[row.id]?.length">
            <el-tag size="small" type="success">{{ navStatusLabel(row) }}</el-tag>
          </template>
          <template v-else>
            <el-button v-if="row.status === 'published'" size="small" link type="primary" @click="handleAddToNav(row)">加入导航</el-button>
            <span v-else class="muted">—</span>
          </template>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="430">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button v-if="canPreview(row)" size="small" @click="handlePreview(row)">预览</el-button>
          <el-button size="small" type="warning" @click="$router.push('/hero')" v-if="isHome(row)">轮播图</el-button>
          <el-button v-permission="'page:publish'" size="small" type="success" @click="handlePublish(row)" v-if="row.status !== 'published'">发布</el-button>
          <el-button v-permission="'page:publish'" size="small" type="warning" @click="handleUnpublish(row)" v-else>下线</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[10, 20, 50, 100]"
      @current-change="loadData"
      @size-change="loadData"
    />

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑页面' : '新建页面'" width="820px" :close-on-click-modal="false">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基础信息" name="basic">
          <el-form :model="form" label-width="90px">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="Slug" required>
          <el-input v-model="form.slug" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type">
            <el-option label="普通" value="normal" />
            <el-option label="首页" value="home" />
            <el-option label="产品" value="product" />
            <el-option label="OEM" value="oem" />
            <el-option label="ODM" value="odm" />
            <el-option label="工厂" value="factory" />
            <el-option label="FAQ" value="faq" />
            <el-option label="联系我们" value="contact" />
          </el-select>
        </el-form-item>
            <el-form-item label="模板">
              <el-input v-model="form.template" placeholder="可选，如 default" />
            </el-form-item>
            <el-form-item label="排序">
              <el-input-number v-model="form.sort_order" :min="0" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="模块编辑" name="modules" :disabled="!editingId">
          <el-alert
            v-if="modules.some((m: any) => m.type === 'banner')"
            type="info"
            show-icon
            :closable="false"
            class="banner-tip"
            title="轮播 Banner 模块建议通过列表中的「轮播图」按钮可视化编辑，直接修改 Config 可能影响轮播展示。"
          />
          <div class="module-toolbar">
            <el-button size="small" type="primary" @click="addModule">添加模块</el-button>
          </div>
          <div v-for="(m, i) in modules" :key="i" class="module-item">
            <div class="module-head">
              <span class="module-index">{{ i + 1 }}</span>
              <el-select v-model="m.type" size="small" style="width: 150px">
                <el-option v-for="t in MODULE_TYPES" :key="t.value" :label="t.label" :value="t.value" />
              </el-select>
              <el-input v-model="m.title" size="small" placeholder="模块标题" style="flex: 1" />
              <el-switch v-model="m.is_visible" size="small" active-text="显示" />
              <el-button-group>
                <el-button size="small" :disabled="i === 0" @click="moveModule(i, -1)">↑</el-button>
                <el-button size="small" :disabled="i === modules.length - 1" @click="moveModule(i, 1)">↓</el-button>
                <el-button size="small" type="danger" @click="removeModule(i)">删除</el-button>
              </el-button-group>
            </div>
            <el-input
              v-model="m.config"
              type="textarea"
              :rows="3"
              class="module-config"
              placeholder='模块配置（JSON，可选），如 {"description": "..."}'
            />
          </div>
          <el-empty v-if="!modules.length" description="暂无模块，点击「添加模块」开始搭建页面" :image-size="60" />
        </el-tab-pane>

        <el-tab-pane label="多语言翻译" name="translations">
          <div class="module-toolbar">
            <el-button size="small" type="primary" @click="addTranslation">添加翻译</el-button>
          </div>
          <div v-for="(t, i) in translations" :key="i" class="translation-item">
            <div class="translation-head">
              <el-select v-model="t.language" size="small" placeholder="语言" style="width: 120px">
                <el-option v-for="lang in languages" :key="lang" :label="lang" :value="lang" />
              </el-select>
              <el-input v-model="t.title" size="small" placeholder="翻译后的标题" style="flex: 1" />
              <el-button size="small" type="danger" @click="removeTranslation(i)">删除</el-button>
            </div>
            <el-input v-model="t.content" type="textarea" :rows="3" placeholder="翻译后的正文内容（可选）" />
          </div>
          <el-empty v-if="!translations.length" description="暂无翻译，点击「添加翻译」维护多语言内容" :image-size="60" />
        </el-tab-pane>

        <el-tab-pane v-if="editingId" label="版本历史" name="versions" lazy>
          <div class="version-toolbar">
            <el-input v-model="draftNote" size="small" placeholder="版本备注（可选）" style="width: 260px" />
            <el-button size="small" type="primary" @click="saveDraftVersion">将当前内容存为草稿版本</el-button>
          </div>
          <el-table :data="versions" v-loading="versionsLoading" size="small">
            <el-table-column prop="version" label="版本" width="80">
              <template #default="{ row }">v{{ row.version }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag size="small" :type="row.status === 'published' ? 'success' : row.status === 'draft' ? 'warning' : 'info'">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="note" label="备注" min-width="150" show-overflow-tooltip />
            <el-table-column prop="created_at" label="创建时间" width="160" />
            <el-table-column prop="published_at" label="发布时间" width="160" />
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button v-permission="'page:publish'" v-if="row.status !== 'published'" size="small" type="success" @click="publishVersion(row)">发布</el-button>
                <el-button v-permission="'page:publish'" v-else size="small" @click="rollbackVersion(row)">回滚</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!versions.length && !versionsLoading" description="暂无版本记录" :image-size="60" />
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { pageApi, portalApi, publicApi, navigationApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import type { Page, Navigation } from '@/types'

const router = useRouter()
const pages = ref<Page[]>([])
const navStatusMap = ref<Record<string, Navigation[]>>({})   // page_id → nav items
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const status = ref('')

// 模块类型（对齐需求文档 §5.5 与后端 PageModule 注释）
const MODULE_TYPES = [
  { value: 'banner', label: '轮播 Banner' },
  { value: 'text', label: '富文本' },
  { value: 'image', label: '图片' },
  { value: 'video', label: '视频' },
  { value: 'product_recommend', label: '产品推荐' },
  { value: 'oem', label: 'OEM 服务' },
  { value: 'odm', label: 'ODM 服务' },
  { value: 'factory', label: '工厂展示' },
  { value: 'production', label: '生产流程' },
  { value: 'case', label: '案例' },
  { value: 'certification', label: '认证' },
  { value: 'blog', label: '博客' },
  { value: 'contact', label: '联系表单' },
]
const languages = ref<string[]>(['en', 'zh', 'es', 'fr'])

// ---------- 编辑器状态 ----------
const dialogVisible = ref(false)
const editingId = ref('')
const activeTab = ref('basic')
const saving = ref(false)
const form = reactive({ title: '', slug: '', type: 'normal', template: '', sort_order: 0 })
const modules = ref<any[]>([])
const translations = ref<any[]>([])

async function loadData() {
  loading.value = true
  try {
    const result = await pageApi.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, status: status.value })
    pages.value = (result as any).items || result || []
    total.value = (result as any).total ?? pages.value.length
    // 批量加载导航关联状态
    loadNavStatus()
  } finally { loading.value = false }
}

async function loadNavStatus() {
  try {
    const allNavs = await navigationApi.listAll()
    const map: Record<string, Navigation[]> = {}
    for (const n of (allNavs || [])) {
      if (n.page_id) {
        if (!map[n.page_id]) map[n.page_id] = []
        map[n.page_id].push(n)
      }
    }
    navStatusMap.value = map
  } catch { navStatusMap.value = {} }
}
function handleSearch() { page.value = 1; loadData() }

const isHome = (row: Page) => row.slug === 'home' || row.type === 'home'
const statusTag = (s: string) => (s === 'published' ? 'success' : s === 'review' ? 'warning' : 'info')
const navStatusLabel = (row: Page) => {
  const navs = navStatusMap.value[row.id]
  if (!navs || navs.length === 0) return ''
  const types = [...new Set(navs.map(n => n.type))]
  return types.join(' / ')
}

const PAGE_TYPE_URL_MAP: Record<string, string> = {
  home: '/', product: '/products', oem: '/oem', odm: '/odm',
  factory: '/factory', blog: '/blog', case: '/cases', faq: '/faq',
  contact: '/contact', production: '/production',
  private_label: '/private-label', normal: '',
}

// 页面类型 → 门户预览页映射（portal 无通用 CMS 页面路由，仅映射有对应路由的类型）
const PREVIEW_PAGE_MAP: Record<string, string> = { home: 'home', product: 'products', blog: 'blogs', case: 'cases', faq: 'faqs', contact: 'contact' }
const canPreview = (row: Page) => !!PREVIEW_PAGE_MAP[row.type]
function handlePreview(row: Page) {
  router.push({ path: '/portal-preview', query: { page: PREVIEW_PAGE_MAP[row.type] || 'home' } })
}

function resetEditor() {
  editingId.value = ''
  activeTab.value = 'basic'
  Object.assign(form, { title: '', slug: '', type: 'normal', template: '', sort_order: 0 })
  modules.value = []
  translations.value = []
  versions.value = []
  draftNote.value = ''
}
function openCreateDialog() { resetEditor(); dialogVisible.value = true }
async function openEditDialog(row: Page) {
  resetEditor()
  editingId.value = row.id
  dialogVisible.value = true
  try {
    const detail: any = await pageApi.get(row.id)
    Object.assign(form, { title: detail.title, slug: detail.slug, type: detail.type, template: detail.template || '', sort_order: detail.sort_order ?? 0 })
    modules.value = (detail.modules || []).map((m: any) => ({ type: m.type, title: m.title || '', sort_order: m.sort_order ?? 0, is_visible: m.is_visible !== false, config: m.config || '' }))
    translations.value = (detail.translations || []).map((t: any) => ({ language: t.language, title: t.title || '', content: t.content || '' }))
  } catch { /* error handled */ }
}
function addModule() {
  modules.value.push({ type: 'text', title: '', sort_order: modules.value.length + 1, is_visible: true, config: '' })
}
function removeModule(i: number) { modules.value.splice(i, 1); reindexModules() }
function moveModule(i: number, dir: number) {
  const j = i + dir
  if (j < 0 || j >= modules.value.length) return
  const arr = modules.value
  ;[arr[i], arr[j]] = [arr[j], arr[i]]
  reindexModules()
}
function reindexModules() { modules.value.forEach((m: any, idx: number) => { m.sort_order = idx + 1 }) }
function moduleTypeLabel(type: string) { return MODULE_TYPES.find(t => t.value === type)?.label || type }

function addTranslation() {
  const used = new Set(translations.value.map((t: any) => t.language))
  const lang = languages.value.find(l => !used.has(l)) || ''
  translations.value.push({ language: lang, title: '', content: '' })
}
function removeTranslation(i: number) { translations.value.splice(i, 1) }

function validate(): boolean {
  if (!form.title.trim()) { ElMessage.warning('请输入标题'); activeTab.value = 'basic'; return false }
  if (!form.slug.trim()) { ElMessage.warning('请输入 Slug'); activeTab.value = 'basic'; return false }
  for (const m of modules.value) {
    if (!m.type) { ElMessage.warning('模块类型不能为空'); activeTab.value = 'modules'; return false }
    if (m.config && m.config.trim()) {
      try { JSON.parse(m.config) } catch {
        ElMessage.warning(`模块「${moduleTypeLabel(m.type)}」的 Config 不是合法 JSON`); activeTab.value = 'modules'; return false
      }
    }
  }
  for (const t of translations.value) {
    if (!t.language) { ElMessage.warning('翻译语言不能为空'); activeTab.value = 'translations'; return false }
  }
  return true
}
function buildPayload() {
  reindexModules()
  return {
    ...form,
    modules: modules.value.map((m: any) => ({
      type: m.type, title: m.title || '', sort_order: m.sort_order ?? 0,
      is_visible: m.is_visible !== false, config: (m.config || '').trim(),
    })),
    translations: translations.value.map((t: any) => ({ language: t.language, title: t.title || '', content: t.content || '' })),
  }
}
async function handleSave() {
  if (!validate()) return
  saving.value = true
  try {
    const payload = buildPayload()
    if (editingId.value) { await pageApi.update(editingId.value, payload) } else { await pageApi.create(payload) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch { /* error handled */ } finally { saving.value = false }
}

// ---------- 行操作 ----------
async function handlePublish(row: Page) { await pageApi.publish(row.id); ElMessage.success('已发布'); loadData() }
async function handleUnpublish(row: Page) { await pageApi.unpublish(row.id); ElMessage.success('已下线'); loadData() }
async function handleDelete(row: Page) { await ElMessageBox.confirm(`确定删除页面 ${row.title} 吗？`, '警告', { type: 'warning' }); await pageApi.delete(row.id); ElMessage.success('已删除'); loadData() }

async function handleAddToNav(row: Page) {
  if (row.status !== 'published') {
    ElMessage.warning('请先发布页面再加入导航')
    return
  }
  try {
    const url = PAGE_TYPE_URL_MAP[row.type as string] || '/' + row.slug
    await navigationApi.create({
      name: row.title, type: 'header', url,
      target: '_self', sort_order: 99, is_visible: true, page_id: row.id,
    })
    ElMessage.success(`已将「${row.title}」加入 Header 导航`)
    loadNavStatus()
  } catch { }
}

// ---------- 版本历史 ----------
const versions = ref<any[]>([])
const versionsLoading = ref(false)
const draftNote = ref('')
watch(activeTab, tab => { if (tab === 'versions') loadVersions() })

async function loadVersions() {
  if (!editingId.value) return
  versionsLoading.value = true
  try { versions.value = (await portalApi.listVersions(editingId.value)) || [] } catch { versions.value = [] } finally { versionsLoading.value = false }
}
async function saveDraftVersion() {
  if (!editingId.value) return
  if (!validate()) return
  try {
    await portalApi.saveDraft(editingId.value, { snapshot: buildPayload(), note: draftNote.value })
    ElMessage.success('草稿版本已保存'); draftNote.value = ''; loadVersions()
  } catch { /* error handled */ }
}
async function publishVersion(v: any) {
  await ElMessageBox.confirm(`确定发布版本 v${v.version} 吗？发布后前台立即生效。`, '发布版本', { type: 'warning' })
  await portalApi.publishVersion(v.id)
  ElMessage.success('版本已发布'); loadVersions(); loadData()
}
async function rollbackVersion(v: any) {
  await ElMessageBox.confirm(`确定回滚到版本 v${v.version} 吗？将复制该版本快照并立即发布。`, '回滚版本', { type: 'warning' })
  await portalApi.rollbackVersion(v.id)
  ElMessage.success('已回滚'); loadVersions(); loadData()
}

async function init() {
  loadData()
  try {
    const langs: any = await publicApi.languages()
    const codes = (langs || []).map((l: any) => l.code).filter(Boolean)
    if (codes.length) languages.value = codes
  } catch { /* 使用默认语言列表 */ }
}
onMounted(init)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.lang-tag { margin-right: 4px; }
.muted { color: #c0c4cc; }
.module-toolbar { margin-bottom: 12px; }
.module-item, .translation-item {
  border: 1px solid #e4e7ed; border-radius: 6px; padding: 10px 12px; margin-bottom: 10px; background: #fafbfc;
}
.module-head, .translation-head { display: flex; gap: 8px; align-items: center; margin-bottom: 8px; }
.module-index {
  display: inline-block; width: 20px; height: 20px; line-height: 20px; text-align: center;
  background: #409eff; color: #fff; border-radius: 50%; font-size: 12px; flex-shrink: 0;
}
.module-config { margin-top: 4px; }
.banner-tip { margin-bottom: 12px; }
.version-toolbar { display: flex; gap: 8px; margin-bottom: 12px; }
</style>