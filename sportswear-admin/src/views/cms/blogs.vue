<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索标题" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-select v-model="category" placeholder="分类" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="c in CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建博客</el-button>
    </div>
    <el-table :data="blogs" v-loading="loading" stripe>
      <el-table-column label="封面" width="76" align="center">
        <template #default="{ row }">
          <el-image v-if="row.cover_image" :src="row.cover_image" fit="cover" class="cover-thumb" :preview-src-list="[row.cover_image]" preview-teleported />
          <div v-else class="cover-placeholder">—</div>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column label="分类" width="120">
        <template #default="{ row }">{{ categoryLabel(row.category) }}</template>
      </el-table-column>
      <el-table-column label="翻译" min-width="130">
        <template #default="{ row }">
          <template v-if="(row.translations || []).length">
            <el-tag v-for="t in row.translations" :key="t.id || t.language" size="small" type="info" class="lang-tag">{{ t.language }}</el-tag>
          </template>
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="260">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="success" @click="handlePublish(row)" v-if="row.status !== 'published'">发布</el-button>
          <el-button size="small" type="warning" @click="handleUnpublish(row)" v-else>下线</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="loadData" />

    <!-- 编辑弹窗（缺陷 B-01 修复：此前模板缺失导致新建/编辑点击无响应） -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑博客' : '新建博客'" width="760px" :close-on-click-modal="false" destroy-on-close>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基础信息" name="basic">
          <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
            <el-form-item label="标题" prop="title">
              <el-input v-model="form.title" maxlength="200" show-word-limit />
            </el-form-item>
            <el-form-item label="Slug" prop="slug">
              <el-input v-model="form.slug" placeholder="URL 标识，如 oem-vs-odm-guide" />
            </el-form-item>
            <el-form-item label="分类" prop="category">
              <el-select v-model="form.category" style="width: 100%">
                <el-option v-for="c in CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="作者">
              <el-input v-model="form.author" placeholder="可选，如 Admin" />
            </el-form-item>
            <el-form-item label="标签">
              <el-input v-model="form.tags" placeholder="逗号分隔，如 yoga,leggings" />
            </el-form-item>
            <el-form-item label="封面图">
              <MediaPicker v-model="form.cover_image" />
            </el-form-item>
            <el-form-item label="正文">
              <RichTextEditor v-model="form.content" placeholder="正文内容" min-height="220px" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="多语言翻译" name="translations">
          <el-alert
            type="info"
            show-icon
            :closable="false"
            class="trans-tip"
            title="默认语言（English）内容请在「基础信息」编辑；此处仅为其他语种配置翻译，未配置的语言将自动回退默认语言。"
          />
          <div class="trans-toolbar">
            <el-button size="small" type="primary" @click="addTranslation">添加翻译</el-button>
            <span class="trans-tip">切换门户语言后按对应语种展示，未翻译则回退默认语言（YouTube 式）</span>
          </div>
          <div v-for="(t, i) in translations" :key="i" class="translation-item">
            <div class="translation-head">
              <el-select v-model="t.language" size="small" style="width: 160px" placeholder="语言">
                <el-option v-for="l in TRANSLATABLE_LANGS" :key="l.value" :label="l.label" :value="l.value" />
              </el-select>
              <div class="spacer" />
              <el-button size="small" type="danger" @click="removeTranslation(i)">删除</el-button>
            </div>
            <el-input v-model="t.title" size="small" placeholder="翻译标题（留空则沿用默认语言标题）" class="trans-title" />
            <RichTextEditor v-model="t.content" placeholder="翻译正文（留空则沿用默认语言正文）" min-height="120px" class="trans-content" />
          </div>
          <el-empty v-if="!translations.length" description="暂无翻译，点击「添加翻译」配置多语言内容" :image-size="60" />
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
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { cmsApi } from '@/api/cms'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import MediaPicker from '@/components/media/MediaPicker.vue'
import RichTextEditor from '@/components/RichTextEditor.vue'
import type { Blog } from '@/types'
import { checkBlogGate, gateAlertMessage } from '@/utils/publish-gate'

// 分类统一定义：列表筛选与表单共用，避免两处硬编码漂移
const CATEGORIES = [
  { label: 'OEM 指南', value: 'oem_guide' },
  { label: 'ODM 指南', value: 'odm_guide' },
  { label: '面料', value: 'fabric' },
  { label: '趋势', value: 'trend' },
  { label: '行业', value: 'industry' },
]
function categoryLabel(value: string) { return CATEGORIES.find((c) => c.value === value)?.label || value }

// 与门户支持语言保持一致。
// 主语言（源语言）= en：主表 title/content 即英文，在「基础信息」编辑；
// 翻译 Tab 仅配置其他语言，未配置时门户回退主语言内容（YouTube 式多语言）。
const PRIMARY_LANG = 'en'
const TRANSLATABLE_LANGS = [
  { label: '中文', value: 'zh' },
  { label: 'Español', value: 'es' },
  { label: 'Français', value: 'fr' },
]

const blogs = ref<Blog[]>([])
const loading = ref(false)
const saving = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const category = ref('')

const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const activeTab = ref('basic')
const editingId = ref('')
const form = reactive({ title: '', slug: '', category: '', author: '', tags: '', cover_image: '', content: '' })
const translations = ref<Array<{ language: string; title: string; content: string }>>([])
const rules: FormRules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { max: 200, message: '标题不能超过 200 字', trigger: 'blur' },
  ],
  slug: [
    { required: true, message: '请输入 Slug', trigger: 'blur' },
    { pattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/, message: '仅允许小写字母、数字和中划线', trigger: 'blur' },
  ],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
}

async function loadData() {
  loading.value = true
  try {
    const result = await cmsApi.blogs.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, category: category.value })
    blogs.value = result.items
    total.value = result.total
  } finally { loading.value = false }
}
function statusTag(status: string) {
  if (status === 'published') return 'success'
  if (status === 'offline') return 'warning'
  return 'info'
}
function statusLabel(status: string) {
  const map: Record<string, string> = { draft: '草稿', review: '审核中', published: '已发布', offline: '已下线', archived: '已归档' }
  return map[status] || status
}
function handleSearch() { page.value = 1; loadData() }
function resetForm() { Object.assign(form, { title: '', slug: '', category: '', author: '', tags: '', cover_image: '', content: '' }) }
function openCreateDialog() { editingId.value = ''; activeTab.value = 'basic'; resetForm(); translations.value = []; dialogVisible.value = true }
async function openEditDialog(row: Blog) {
  editingId.value = row.id
  activeTab.value = 'basic'
  Object.assign(form, { title: row.title, slug: row.slug, category: row.category, author: row.author || '', tags: row.tags || '', cover_image: row.cover_image || '', content: row.content || '' })
  // 优先用列表已预加载的翻译；为空时回拉详情确保完整回填
  translations.value = (row.translations || []).map((t: any) => ({ language: t.language, title: t.title || '', content: t.content || '' }))
  if (!translations.value.length) {
    try {
      const detail = await cmsApi.blogs.get(row.id)
      translations.value = (detail.translations || []).map((t: any) => ({ language: t.language, title: t.title || '', content: t.content || '' }))
    } catch { /* 忽略：翻译为空也不影响主表编辑 */ }
  }
  dialogVisible.value = true
}
function addTranslation() {
  const used = new Set(translations.value.map((t) => t.language))
  const lang = TRANSLATABLE_LANGS.find((l) => !used.has(l.value))?.value || ''
  translations.value.push({ language: lang, title: '', content: '' })
}
function removeTranslation(i: number) { translations.value.splice(i, 1) }
// 富文本正文是否含实质内容（去除标签后非空）
function hasContent(html?: string) {
  return (html || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').trim().length > 0
}
async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  if (!hasContent(form.content)) { ElMessage.warning('请输入正文'); activeTab.value = 'basic'; return }
  for (const t of translations.value) {
    if (!t.language) { ElMessage.warning('翻译语言不能为空'); activeTab.value = 'translations'; return }
  }
  // 过滤空翻译：标题与正文都为空的行不提交（避免生成无效的空翻译记录）
  const validTranslations = translations.value.filter((t) => t.language && (t.title.trim() || hasContent(t.content)))
  saving.value = true
  try {
    const payload = {
      ...form,
      translations: validTranslations.map((t) => ({ language: t.language, title: t.title || '', content: t.content || '' })),
    }
    if (editingId.value) { await cmsApi.blogs.update(editingId.value, payload) } else { await cmsApi.blogs.create(payload) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch {} finally { saving.value = false }
}
async function handlePublish(row: Blog) {
  // 发布质检门（P0-#2）：标题/Slug/正文缺失时拦截并列出缺失项
  const gate = checkBlogGate(row)
  if (!gate.ok) {
    await ElMessageBox.alert(gateAlertMessage(gate), '无法发布', { type: 'warning', confirmButtonText: '知道了' }).catch(() => {})
    return
  }
  await cmsApi.blogs.publish(row.id)
  ElMessage.success('已发布')
  loadData()
}
async function handleUnpublish(row: Blog) {
  await ElMessageBox.confirm(`确定下线博客「${row.title}」吗？下线后前台不可见。`, '警告', { type: 'warning' })
  await cmsApi.blogs.unpublish(row.id)
  ElMessage.success('已下线')
  loadData()
}
async function handleDelete(row: Blog) { await ElMessageBox.confirm(`确定删除博客 ${row.title} 吗？`, '警告', { type: 'warning' }); await cmsApi.blogs.delete(row.id); ElMessage.success('已删除'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.cover-thumb { width: 44px; height: 44px; border-radius: 4px; }
.cover-placeholder { width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; color: #c0c4cc; background: #f5f7fa; border-radius: 4px; }
.lang-tag { margin-right: 4px; }
.muted { color: #c0c4cc; }
.trans-toolbar { display: flex; gap: 12px; align-items: center; margin-bottom: 12px; }
.trans-tip { color: #909399; font-size: 12px; }
.translation-item { border: 1px solid #e4e7ed; border-radius: 6px; padding: 10px 12px; margin-bottom: 10px; background: #fafbfc; }
.translation-head { display: flex; gap: 8px; align-items: center; margin-bottom: 8px; }
.trans-title { margin-bottom: 8px; }
</style>