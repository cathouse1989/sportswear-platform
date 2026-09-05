<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索标题 / Slug" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-select v-model="projectType" placeholder="项目类型" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="p in PROJECT_TYPES" :key="p.value" :label="p.label" :value="p.value" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建案例</el-button>
    </div>
    <el-table :data="cases" v-loading="loading" stripe>
      <el-table-column label="封面" width="76" align="center">
        <template #default="{ row }">
          <el-image v-if="row.cover_image" :src="row.cover_image" fit="cover" class="cover-thumb" :preview-src-list="[row.cover_image]" preview-teleported />
          <div v-else class="cover-placeholder">—</div>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="client_industry" label="客户行业" width="140" show-overflow-tooltip />
      <el-table-column label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.project_type" size="small" effect="plain" type="primary">{{ projectTypeLabel(row.project_type) }}</el-tag>
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column label="翻译" min-width="150">
        <template #default="{ row }">
          <template v-if="(row.translations || []).length">
            <el-tag v-for="t in row.translations" :key="t.id || t.language" size="small" :type="langTag(t.language).type" class="lang-tag">{{ langTag(t.language).label }}</el-tag>
          </template>
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100" align="center">
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

    <!-- 编辑弹窗：基础信息 + 多语言翻译 -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑案例' : '新建案例'" width="820px" :close-on-click-modal="false" destroy-on-close>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基础信息" name="basic">
          <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="标题" prop="title">
                  <el-input v-model="form.title" maxlength="200" show-word-limit />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Slug" prop="slug">
                  <el-input v-model="form.slug" placeholder="URL 标识，如 yoga-brand-oem" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="客户行业">
                  <el-input v-model="form.client_industry" placeholder="如 Fitness Brand / Retail Chain" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="项目类型">
                  <el-select v-model="form.project_type" style="width: 100%" placeholder="请选择项目类型">
                    <el-option v-for="p in PROJECT_TYPES" :key="p.value" :label="p.label" :value="p.value" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="产品">
                  <el-input v-model="form.products" placeholder="如 Yoga Leggings, Sports Bra" />
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="封面图">
                  <MediaPicker v-model="form.cover_image" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="客户需求">
                  <RichTextEditor v-model="form.client_need" placeholder="客户的品牌定位与采购诉求" min-height="110px" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="问题挑战">
                  <RichTextEditor v-model="form.problem" placeholder="项目面临的核心难点" min-height="110px" />
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="解决方案" prop="solution">
                  <RichTextEditor v-model="form.solution" placeholder="我们如何为客户解决上述问题" min-height="130px" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="实施过程">
                  <RichTextEditor v-model="form.process" placeholder="打样、生产、质检、交付等关键节点" min-height="110px" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="项目成果">
                  <RichTextEditor v-model="form.result" placeholder="交付成果与客户反馈" min-height="110px" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="多语言翻译" name="translations">
          <TransEditor
            v-model="translations"
            :fields="TRANS_FIELDS"
            :source="{ title: form.title, client_need: form.client_need, problem: form.problem, solution: form.solution, process: form.process, result: form.result }"
            :langs="TRANSLATABLE_LANGS"
            source-label="English"
          />
        </el-tab-pane>
        <el-tab-pane label="SEO" name="seo">
          <el-form label-width="110px">
            <el-form-item label="SEO 标题">
              <el-input v-model="form.seo.title" maxlength="200" placeholder="留空则沿用案例标题" />
            </el-form-item>
            <el-form-item label="关键词">
              <el-input v-model="form.seo.keywords" placeholder="逗号分隔，如 yoga,oem,leggings" />
            </el-form-item>
            <el-form-item label="SEO 描述">
              <el-input v-model="form.seo.description" type="textarea" :rows="2" maxlength="200" show-word-limit placeholder="留空则自动取案例摘要" />
            </el-form-item>
            <el-divider content-position="left">Open Graph（社交分享）</el-divider>
            <el-form-item label="OG 标题">
              <el-input v-model="form.seo.og_title" maxlength="200" placeholder="社交分享标题" />
            </el-form-item>
            <el-form-item label="OG 描述">
              <el-input v-model="form.seo.og_description" type="textarea" :rows="2" placeholder="社交分享描述" />
            </el-form-item>
            <el-form-item label="OG 图片">
              <el-input v-model="form.seo.og_image" placeholder="社交分享图片 URL" />
            </el-form-item>
          </el-form>
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
import { caseApi } from '@/api'
import { useCrud } from '@/composables/useCrud'
import MediaPicker from '@/components/media/MediaPicker.vue'
import RichTextEditor from '@/components/RichTextEditor.vue'
import TransEditor from '@/components/cms/TransEditor.vue'
import type { Case } from '@/types'
import { checkCaseGate, gateAlertMessage } from '@/utils/publish-gate'

// 项目类型统一定义：列表展示与表单共用
const PROJECT_TYPES = [
  { label: 'OEM', value: 'oem' },
  { label: 'ODM', value: 'odm' },
  { label: '自有品牌', value: 'private_label' },
]
function projectTypeLabel(value?: string) {
  return PROJECT_TYPES.find((p) => p.value === value)?.label || value || '—'
}

// 主语言（源语言）= en：主表字段即英文，在「基础信息」编辑；
// 翻译 Tab 仅配置其他语言，未配置时门户回退主语言内容（YouTube 式多语言）。
const TRANSLATABLE_LANGS = [
  { label: '中文', value: 'zh' },
  { label: 'Español', value: 'es' },
  { label: 'Français', value: 'fr' },
]
const LANG_META: Record<string, { label: string; type: 'primary' | 'success' | 'warning' | 'info' | 'danger' }> = {
  en: { label: 'EN', type: 'info' },
  zh: { label: '中文', type: 'success' },
  es: { label: 'ES', type: 'warning' },
  fr: { label: 'FR', type: 'danger' },
}
function langTag(language: string) {
  return LANG_META[language] || { label: language, type: 'info' as const }
}



const projectType = ref('')

// 弹窗表单状态（声明在前，供 beforeSave/buildPayload 闭包引用）
const formRef = ref<FormInstance>()
const activeTab = ref('basic')
const form = reactive({
  title: '', slug: '', client_industry: '', project_type: '', products: '',
  cover_image: '', client_need: '', problem: '', solution: '', process: '', result: '',
  seo: { language: 'en', title: '', description: '', keywords: '', og_title: '', og_description: '', og_image: '' },
})
const translations = ref<Record<string, Record<string, string>>>({})
const TRANS_FIELDS = [
  { key: 'title', label: '标题' },
  { key: 'client_need', label: '客户需求', type: 'richtext' as const, minHeight: '90px' },
  { key: 'problem', label: '问题挑战', type: 'richtext' as const, minHeight: '90px' },
  { key: 'solution', label: '解决方案', type: 'richtext' as const, minHeight: '90px' },
  { key: 'process', label: '实施过程', type: 'richtext' as const, minHeight: '90px' },
  { key: 'result', label: '项目成果', type: 'richtext' as const, minHeight: '90px' },
]
function emptyTranslations(): Record<string, Record<string, string>> {
  const blank = { title: '', client_need: '', problem: '', solution: '', process: '', result: '' }
  return { zh: { ...blank }, es: { ...blank }, fr: { ...blank } }
}
function translationsToRecord(list: Array<any>) {
  const next = emptyTranslations()
  for (const t of list || []) {
    const lang = t.language
    if (lang && lang !== 'en' && next[lang]) {
      next[lang] = { title: t.title || '', client_need: t.client_need || '', problem: t.problem || '', solution: t.solution || '', process: t.process || '', result: t.result || '' }
    }
  }
  return next
}
const rules: FormRules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { max: 200, message: '标题不能超过 200 字', trigger: 'blur' },
  ],
  slug: [
    { required: true, message: '请输入 Slug', trigger: 'blur' },
    { pattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/, message: '仅允许小写字母、数字和中划线', trigger: 'blur' },
  ],
}
function hasText(v?: string) { return Boolean((v || '').trim()) }
function hasHtmlContent(v?: string) { return Boolean((v || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').trim()) }
function translationHasContent(t: Record<string, string>) {
  return hasText(t.title) || hasHtmlContent(t.client_need) || hasHtmlContent(t.problem) || hasHtmlContent(t.solution) || hasHtmlContent(t.process) || hasHtmlContent(t.result)
}

// 列表分页/搜索 + 保存（校验 + 翻译/SEO payload）复用 useCrud
const {
  items: cases, total, loading, page, pageSize, keyword,
  loadData, handleSearch, saving, handleSave,
  dialogVisible, editingId,
} = useCrud({
  api: caseApi,
  emptyForm: () => ({}),
  serverPagination: true,
  extraParams: () => ({ project_type: projectType.value }),
  beforeSave: async () => {
    const valid = await formRef.value?.validate().catch(() => false)
    if (!valid) return false
    return true
  },
  buildPayload: () => {
    const seoHasContent = ['title', 'description', 'keywords', 'og_title', 'og_description', 'og_image'].some((f) => hasText((form.seo as any)[f]))
    return {
      ...form,
      translations: Object.entries(translations.value)
        .filter(([, t]) => translationHasContent(t))
        .map(([language, t]) => ({ language, ...t })),
      seo: seoHasContent ? { ...form.seo } : null,
    }
  },
})

function statusTag(status: string) {
  if (status === 'published') return 'success'
  if (status === 'offline') return 'warning'
  return 'info'
}
function statusLabel(status: string) {
  const map: Record<string, string> = { draft: '草稿', review: '审核中', published: '已发布', offline: '已下线', archived: '已归档' }
  return map[status] || status
}
function resetForm() {
  Object.assign(form, {
    title: '', slug: '', client_industry: '', project_type: '', products: '',
    cover_image: '', client_need: '', problem: '', solution: '', process: '', result: '',
  })
  Object.assign(form.seo, { language: 'en', title: '', description: '', keywords: '', og_title: '', og_description: '', og_image: '' })
}
function openCreateDialog() {
  editingId.value = ''
  activeTab.value = 'basic'
  resetForm()
  translations.value = emptyTranslations()
  dialogVisible.value = true
}
async function openEditDialog(row: Case) {
  editingId.value = row.id
  activeTab.value = 'basic'
  Object.assign(form, {
    title: row.title, slug: row.slug,
    client_industry: row.client_industry || '',
    project_type: row.project_type || '',
    products: row.products || '',
    cover_image: row.cover_image || '',
    client_need: row.client_need || '',
    problem: row.problem || '',
    solution: row.solution || '',
    process: row.process || '',
    result: row.result || '',
  })
  // 翻译先用列表预加载值，随后用详情覆盖（同时回填 SEO）
  translations.value = translationsToRecord(row.translations || [])
  try {
    const detail = await caseApi.get(row.id)
    translations.value = translationsToRecord(detail.translations || [])
    Object.assign(form.seo, {
      language: detail.seo?.language || 'en',
      title: detail.seo?.title || '',
      description: detail.seo?.description || '',
      keywords: detail.seo?.keywords || '',
      og_title: detail.seo?.og_title || '',
      og_description: detail.seo?.og_description || '',
      og_image: detail.seo?.og_image || '',
    })
  } catch { /* 详情回填失败不影响编辑主表 */ }
  dialogVisible.value = true
}

async function handlePublish(row: Case) {
  // 发布质检门：标题/Slug/解决方案缺失时拦截并列出缺失项
  const gate = checkCaseGate({
    title: row.title,
    slug: row.slug,
    solution: row.solution,
    cover_image: row.cover_image,
    client_industry: row.client_industry,
  })
  if (!gate.ok) {
    await ElMessageBox.alert(gateAlertMessage(gate), '无法发布', { type: 'warning', confirmButtonText: '知道了' }).catch(() => {})
    return
  }
  await caseApi.publish(row.id)
  ElMessage.success('已发布')
  loadData()
}
async function handleUnpublish(row: Case) {
  await ElMessageBox.confirm(`确定下线案例「${row.title}」吗？下线后前台不可见。`, '警告', { type: 'warning' })
  await caseApi.unpublish(row.id)
  ElMessage.success('已下线')
  loadData()
}
async function handleDelete(row: Case) {
  await ElMessageBox.confirm(`确定删除案例「${row.title}」吗？`, '警告', { type: 'warning' })
  await caseApi.delete(row.id)
  ElMessage.success('已删除')
  loadData()
}
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

.trans-tip { margin-bottom: 12px; }
.trans-toolbar { display: flex; gap: 12px; align-items: center; margin-bottom: 12px; }
.trans-toolbar .trans-tip { margin-bottom: 0; color: #909399; font-size: 12px; }
.translation-item { border: 1px solid #e4e7ed; border-radius: 8px; padding: 12px 14px; margin-bottom: 12px; background: #fafbfc; }
.translation-head { display: flex; gap: 8px; align-items: center; margin-bottom: 10px; }
.trans-field { margin-bottom: 10px; }
.trans-label { font-size: 12px; color: #909399; margin-bottom: 4px; }
.lang-badge {
  display: inline-flex; align-items: center; padding: 1px 8px; border-radius: 10px;
  font-size: 12px; line-height: 18px; color: #606266; background: #f0f2f5;
}
.lang-zh { color: #67c23a; background: #f0f9eb; }
.lang-es { color: #e6a23c; background: #fdf6ec; }
.lang-fr { color: #f56c6c; background: #fef0f0; }
</style>