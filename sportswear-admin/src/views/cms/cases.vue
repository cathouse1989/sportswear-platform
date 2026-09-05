<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索标题" clearable style="width: 240px" @keyup.enter="handleSearch" />
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
                  <el-input v-model="form.client_need" type="textarea" :rows="3" placeholder="客户的品牌定位与采购诉求" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="问题挑战">
                  <el-input v-model="form.problem" type="textarea" :rows="3" placeholder="项目面临的核心难点" />
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="解决方案" prop="solution">
                  <el-input v-model="form.solution" type="textarea" :rows="4" placeholder="我们如何为客户解决上述问题" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="实施过程">
                  <el-input v-model="form.process" type="textarea" :rows="3" placeholder="打样、生产、质检、交付等关键节点" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="项目成果">
                  <el-input v-model="form.result" type="textarea" :rows="3" placeholder="交付成果与客户反馈" />
                </el-form-item>
              </el-col>
            </el-row>
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
            <span class="trans-tip">切换门户语言后按对应语种展示，未翻译则回退默认语言</span>
          </div>
          <div v-for="(t, i) in translations" :key="i" class="translation-item">
            <div class="translation-head">
              <el-select v-model="t.language" size="small" style="width: 160px" placeholder="语言">
                <el-option v-for="l in TRANSLATABLE_LANGS" :key="l.value" :label="l.label" :value="l.value" />
              </el-select>
              <span v-if="t.language" class="lang-badge" :class="`lang-${t.language}`">{{ langTag(t.language).label }}</span>
              <div class="spacer" />
              <el-button size="small" type="danger" text @click="removeTranslation(i)">删除</el-button>
            </div>
            <el-input v-model="t.title" size="small" placeholder="翻译标题（留空则沿用默认语言标题）" class="trans-field" />
            <el-row :gutter="12">
              <el-col :span="12">
                <div class="trans-label">客户需求</div>
                <el-input v-model="t.client_need" type="textarea" :rows="3" placeholder="留空则沿用默认语言" />
              </el-col>
              <el-col :span="12">
                <div class="trans-label">问题挑战</div>
                <el-input v-model="t.problem" type="textarea" :rows="3" placeholder="留空则沿用默认语言" />
              </el-col>
              <el-col :span="12">
                <div class="trans-label">解决方案</div>
                <el-input v-model="t.solution" type="textarea" :rows="3" placeholder="留空则沿用默认语言" />
              </el-col>
              <el-col :span="12">
                <div class="trans-label">实施过程</div>
                <el-input v-model="t.process" type="textarea" :rows="3" placeholder="留空则沿用默认语言" />
              </el-col>
              <el-col :span="24">
                <div class="trans-label">项目成果</div>
                <el-input v-model="t.result" type="textarea" :rows="3" placeholder="留空则沿用默认语言" />
              </el-col>
            </el-row>
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
import { caseApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import MediaPicker from '@/components/media/MediaPicker.vue'
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

interface TranslationForm {
  language: string
  title: string
  client_need: string
  problem: string
  solution: string
  process: string
  result: string
}

const cases = ref<Case[]>([])
const loading = ref(false)
const saving = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')

const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const activeTab = ref('basic')
const editingId = ref('')
const form = reactive({
  title: '', slug: '', client_industry: '', project_type: '', products: '',
  cover_image: '', client_need: '', problem: '', solution: '', process: '', result: '',
})
const translations = ref<TranslationForm[]>([])
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

async function loadData() {
  loading.value = true
  try {
    const result = await caseApi.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value })
    cases.value = result.items
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
function resetForm() {
  Object.assign(form, {
    title: '', slug: '', client_industry: '', project_type: '', products: '',
    cover_image: '', client_need: '', problem: '', solution: '', process: '', result: '',
  })
}
function openCreateDialog() {
  editingId.value = ''
  activeTab.value = 'basic'
  resetForm()
  translations.value = []
  dialogVisible.value = true
}
function openEditDialog(row: Case) {
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
  // 列表已预加载翻译；为空时回拉详情确保完整回填
  translations.value = (row.translations || []).map((t: any) => ({
    language: t.language,
    title: t.title || '',
    client_need: t.client_need || '',
    problem: t.problem || '',
    solution: t.solution || '',
    process: t.process || '',
    result: t.result || '',
  }))
  if (!translations.value.length) {
    caseApi.get(row.id).then((detail) => {
      translations.value = (detail.translations || []).map((t: any) => ({
        language: t.language,
        title: t.title || '',
        client_need: t.client_need || '',
        problem: t.problem || '',
        solution: t.solution || '',
        process: t.process || '',
        result: t.result || '',
      }))
    }).catch(() => { /* 翻译为空不影响主表编辑 */ })
  }
  dialogVisible.value = true
}
function addTranslation() {
  const used = new Set(translations.value.map((t) => t.language))
  const lang = TRANSLATABLE_LANGS.find((l) => !used.has(l.value))?.value || ''
  translations.value.push({ language: lang, title: '', client_need: '', problem: '', solution: '', process: '', result: '' })
}
function removeTranslation(i: number) { translations.value.splice(i, 1) }
function hasText(v?: string) { return Boolean((v || '').trim()) }
function translationHasContent(t: TranslationForm) {
  return hasText(t.title) || hasText(t.client_need) || hasText(t.problem) || hasText(t.solution) || hasText(t.process) || hasText(t.result)
}
async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  for (const t of translations.value) {
    if (!t.language) { ElMessage.warning('翻译语言不能为空'); activeTab.value = 'translations'; return }
  }
  // 过滤空翻译：所有字段均为空的行不提交（避免生成无效空翻译记录）
  const validTranslations = translations.value.filter((t) => t.language && translationHasContent(t))
  saving.value = true
  try {
    const payload = {
      ...form,
      translations: validTranslations.map((t) => ({
        language: t.language,
        title: t.title || '',
        client_need: t.client_need || '',
        problem: t.problem || '',
        solution: t.solution || '',
        process: t.process || '',
        result: t.result || '',
      })),
    }
    if (editingId.value) { await caseApi.update(editingId.value, payload) } else { await caseApi.create(payload) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch {} finally { saving.value = false }
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