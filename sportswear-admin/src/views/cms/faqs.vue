<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索问题" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-select v-model="category" placeholder="分类" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="c in enumOptions('faq.category')" :key="c.value" :label="c.label" :value="c.value" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建 FAQ</el-button>
    </div>
    <el-table :data="faqs" v-loading="loading" stripe>
      <el-table-column prop="question" label="问题" min-width="300" />
      <el-table-column prop="category" label="分类" width="110">
        <template #default="{ row }">{{ enumLabel('faq.category', row.category) }}</template>
      </el-table-column>
      <el-table-column label="翻译" min-width="130">
        <template #default="{ row }">
          <template v-if="(row.translations || []).length">
            <el-tag v-for="t in row.translations" :key="t.id || t.language" size="small" type="info" class="lang-tag">{{ t.language }}</el-tag>
          </template>
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'" size="small">{{ row.is_active ? '启用' : '停用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="handleSearch" />

    <!-- 编辑弹窗（缺陷 B-01 修复：此前模板缺失导致新建/编辑点击无响应） -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑 FAQ' : '新建 FAQ'" width="640px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="问题" prop="question">
          <el-input v-model="form.question" maxlength="500" show-word-limit />
        </el-form-item>
        <el-form-item label="答案" prop="answer">
          <el-input v-model="form.answer" type="textarea" :rows="5" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="form.category" clearable style="width: 100%">
            <el-option v-for="c in enumOptions('faq.category')" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-divider content-position="left">多语言翻译（未填写回退英文）</el-divider>
        <TransEditor v-model="translations" :fields="TRANS_FIELDS" :source="{ question: form.question, answer: form.answer }" :langs="TRANS_LANGS" source-label="English" />
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.is_active" active-text="前台可见" inactive-text="已停用" />
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
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { faqApi } from '@/api'
import { useCrud } from '@/composables/useCrud'
import { useEnumDict } from '@/composables/useEnumDict'
import TransEditor from '@/components/cms/TransEditor.vue'
import { DEFAULT_TRANS_LANGS, createEmptyTranslations, translationsToRecord, translationsToPayload } from '@/composables/useTransRecord'

// FAQ 分类由「字典管理」驱动（faq.category）
const { ensureLoaded: loadEnumDict, options: enumOptions, label: enumLabel } = useEnumDict()

// 多语言：源语言 English（主表），目标语言 zh/es/fr（翻译表）
const TRANS_LANGS = DEFAULT_TRANS_LANGS
const TRANS_FIELDS = [
  { key: 'question', label: '问题' },
  { key: 'answer', label: '答案', type: 'textarea' as const, rows: 4 },
]
const TRANS_KEYS = TRANS_FIELDS.map((f) => f.key)
const translations = ref<Record<string, Record<string, string>>>({})

const category = ref('')
const formRef = ref<FormInstance>()
const rules: FormRules = {
  question: [{ required: true, message: '请输入问题', trigger: 'blur' }],
  answer: [{ required: true, message: '请输入答案', trigger: 'blur' }],
}

const emptyForm = () => ({ question: '', answer: '', category: '', sort_order: 0, is_active: true })
const {
  items: faqs, total, loading, page, pageSize, keyword,
  loadData, handleSearch,
  dialogVisible, editingId, form, saving,
  handleSave,
  openCreateDialog: crudCreate, openEditDialog: crudEdit,
} = useCrud({
  api: faqApi,
  emptyForm,
  serverPagination: true,
  extraParams: () => ({ category: category.value }),
  buildPayload: (): Record<string, any> => ({
    ...(form as any),
    translations: translationsToPayload(translations.value, TRANS_KEYS, TRANS_LANGS),
  }),
  // 表单校验作为保存前钩子，false 时中止保存
  beforeSave: async () => {
    const valid = await formRef.value?.validate().catch(() => false)
    return !!valid
  },
})

function openCreateDialog() { crudCreate(); translations.value = createEmptyTranslations(TRANS_KEYS) }
function openEditDialog(row: any) { crudEdit(row); translations.value = translationsToRecord(row.translations || [], TRANS_KEYS) }

// FAQ 删除确认不带实体名（问题文本可能过长）
async function handleDelete(row: any) {
  await ElMessageBox.confirm('确定删除此 FAQ 吗？', '警告', { type: 'warning' })
  await faqApi.delete(row.id)
  ElMessage.success('已删除'); loadData()
}
onMounted(() => { loadEnumDict(); loadData() })
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.lang-tag { margin-right: 4px; }
.muted { color: #c0c4cc; }
</style>