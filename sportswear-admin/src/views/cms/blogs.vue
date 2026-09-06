<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索标题" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-select v-model="category" placeholder="分类" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="c in enumOptions('blog.category')" :key="c.value" :label="c.label" :value="c.value" />
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
        <template #default="{ row }">{{ enumLabel('blog.category', row.category) }}</template>
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
          <el-tag :type="enumTag('content.status', row.status)" size="small">{{ enumLabel('content.status', row.status) }}</el-tag>
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
                <el-option v-for="c in enumOptions('blog.category')" :key="c.value" :label="c.label" :value="c.value" />
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
          <TransEditor
            v-model="translations"
            :fields="TRANS_FIELDS"
            :source="{ title: form.title, content: form.content }"
            :langs="TRANSLATABLE_LANGS"
            source-label="English"
          />
        </el-tab-pane>
      </el-tabs>
      <template #footer>
        <el-button @click="previewItem">预览</el-button>
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
import { blogApi } from '@/api'
import { useCrud } from '@/composables/useCrud'
import { useEnumDict } from '@/composables/useEnumDict'
import { usePortalPreview } from '@/composables/usePortalPreview'
import MediaPicker from '@/components/media/MediaPicker.vue'
import RichTextEditor from '@/components/RichTextEditor.vue'
import TransEditor from '@/components/cms/TransEditor.vue'
import { createEmptyTranslations, translationsToRecord, translationsToPayload } from '@/composables/useTransRecord'
import type { Blog } from '@/types'
import { checkBlogGate, gateAlertMessage } from '@/utils/publish-gate'

// 分类/状态标签由「字典管理」驱动（blog.category / content.status）
const { ensureLoaded: loadEnumDict, options: enumOptions, label: enumLabel, tagType: enumTag } = useEnumDict()

// 与门户支持语言保持一致。
// 主语言（源语言）= en：主表 title/content 即英文，在「基础信息」编辑；
// 翻译 Tab 仅配置其他语言，未配置时门户回退主语言内容（YouTube 式多语言）。
const PRIMARY_LANG = 'en'
const TRANSLATABLE_LANGS = [
  { label: '中文', value: 'zh' },
  { label: 'Español', value: 'es' },
  { label: 'Français', value: 'fr' },
]

const category = ref('')

// 弹窗表单状态（声明在前，供 beforeSave/buildPayload 闭包引用）
const formRef = ref<FormInstance>()
const activeTab = ref('basic')
const form = reactive({ title: '', slug: '', category: '', author: '', tags: '', cover_image: '', content: '' })
const translations = ref<Record<string, Record<string, string>>>({})
const TRANS_FIELDS = [
  { key: 'title', label: '标题' },
  { key: 'content', label: '正文', type: 'richtext' as const, minHeight: '120px' },
]
const TRANS_KEYS = TRANS_FIELDS.map((f) => f.key)
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
// 富文本正文是否含实质内容（去除标签后非空）
function hasContent(html?: string) {
  return (html || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').trim().length > 0
}

// 列表分页/搜索 + 保存（校验 + 翻译 payload）复用 useCrud
const {
  items: blogs, total, loading, page, pageSize, keyword,
  loadData, handleSearch, saving, handleSave,
  dialogVisible, editingId,
} = useCrud({
  api: blogApi,
  emptyForm: () => ({}),
  serverPagination: true,
  extraParams: () => ({ category: category.value }),
  beforeSave: async () => {
    const valid = await formRef.value?.validate().catch(() => false)
    if (!valid) return false
    if (!hasContent(form.content)) { ElMessage.warning('请输入正文'); activeTab.value = 'basic'; return false }
    return true
  },
  buildPayload: () => ({
    ...form,
    translations: translationsToPayload(translations.value, TRANS_KEYS, TRANSLATABLE_LANGS),
  }),
})

function resetForm() { Object.assign(form, { title: '', slug: '', category: '', author: '', tags: '', cover_image: '', content: '' }) }
function openCreateDialog() { editingId.value = ''; activeTab.value = 'basic'; resetForm(); translations.value = createEmptyTranslations(TRANS_KEYS); dialogVisible.value = true }
async function openEditDialog(row: Blog) {
  editingId.value = row.id
  activeTab.value = 'basic'
  Object.assign(form, { title: row.title, slug: row.slug, category: row.category, author: row.author || '', tags: row.tags || '', cover_image: row.cover_image || '', content: row.content || '' })
  // 优先用列表已预加载的翻译；为空时回拉详情确保完整回填
  translations.value = translationsToRecord(row.translations || [], TRANS_KEYS)
  if (!(row.translations || []).length) {
    try {
      const detail = await blogApi.get(row.id)
      translations.value = translationsToRecord(detail.translations || [], TRANS_KEYS)
    } catch { /* 忽略：翻译为空也不影响主表编辑 */ }
  }
  dialogVisible.value = true
}

async function handlePublish(row: Blog) {
  // 发布质检门（P0-#2）：标题/Slug/正文缺失时拦截并列出缺失项
  const gate = checkBlogGate(row)
  if (!gate.ok) {
    await ElMessageBox.alert(gateAlertMessage(gate), '无法发布', { type: 'warning', confirmButtonText: '知道了' }).catch(() => {})
    return
  }
  await blogApi.publish(row.id)
  ElMessage.success('已发布')
  loadData()
}
async function handleUnpublish(row: Blog) {
  await ElMessageBox.confirm(`确定下线博客「${row.title}」吗？下线后前台不可见。`, '警告', { type: 'warning' })
  await blogApi.unpublish(row.id)
  ElMessage.success('已下线')
  loadData()
}
async function handleDelete(row: Blog) { await ElMessageBox.confirm(`确定删除博客 ${row.title} 吗？`, '警告', { type: 'warning' }); await blogApi.delete(row.id); ElMessage.success('已删除'); loadData() }
const { previewWithSlug } = usePortalPreview()
function previewItem() { previewWithSlug('blog', form.slug) }
onMounted(() => { loadEnumDict(); loadData() })
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