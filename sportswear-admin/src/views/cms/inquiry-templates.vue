<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索问题" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-select v-model="category" placeholder="分类" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="c in enumOptions('faq.category')" :key="c.value" :label="c.label" :value="c.value" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建模板</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="question" label="问题（英文源）" min-width="320" show-overflow-tooltip />
      <el-table-column prop="category" label="分类" width="110">
        <template #default="{ row }">{{ enumLabel('faq.category', row.category) }}</template>
      </el-table-column>
      <el-table-column label="询盘类型" width="130">
        <template #default="{ row }">
          <el-tag v-if="row.project_type" size="small" effect="plain" type="primary">{{ enumLabel('case.project_type', row.project_type) }}</el-tag>
          <span v-else class="muted">contact</span>
        </template>
      </el-table-column>
      <el-table-column label="翻译" min-width="120">
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑模板' : '新建模板'" width="640px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="问题" prop="question">
          <el-input v-model="form.question" type="textarea" :rows="3" maxlength="500" show-word-limit />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="form.category" clearable style="width: 100%">
            <el-option v-for="c in enumOptions('faq.category')" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="询盘类型">
          <el-select v-model="form.project_type" clearable placeholder="默认 contact" style="width: 100%">
            <el-option v-for="p in enumOptions('case.project_type')" :key="p.value" :label="p.label" :value="p.value" />
          </el-select>
        </el-form-item>
        <el-divider content-position="left">多语言翻译（未填写回退英文）</el-divider>
        <TransEditor v-model="translations" :fields="TRANS_FIELDS" :source="{ question: form.question }" :langs="TRANS_LANGS" source-label="English" />
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
import { inquiryTemplateApi } from '@/api'
import { useCrud } from '@/composables/useCrud'
import { useEnumDict } from '@/composables/useEnumDict'
import TransEditor from '@/components/cms/TransEditor.vue'
import { DEFAULT_TRANS_LANGS, createEmptyTranslations, translationsToRecord, translationsToPayload } from '@/composables/useTransRecord'

// 分类/询盘类型由「字典管理」驱动（faq.category / case.project_type）
const { ensureLoaded: loadEnumDict, options: enumOptions, label: enumLabel } = useEnumDict()

// 多语言：源语言 English（主表），目标语言 zh/es/fr（翻译表）
const TRANS_LANGS = DEFAULT_TRANS_LANGS
const TRANS_FIELDS = [
  { key: 'question', label: '问题', type: 'textarea' as const, rows: 3 },
]
const TRANS_KEYS = TRANS_FIELDS.map((f) => f.key)
const translations = ref<Record<string, Record<string, string>>>({})

const category = ref('')
const formRef = ref<FormInstance>()
const rules: FormRules = {
  question: [{ required: true, message: '请输入问题', trigger: 'blur' }],
}

const emptyForm = () => ({ question: '', category: '', project_type: '', sort_order: 0, is_active: true })
const {
  items, total, loading, page, pageSize, keyword,
  loadData, handleSearch,
  dialogVisible, editingId, form, saving,
  handleSave,
  openCreateDialog: crudCreate, openEditDialog: crudEdit,
} = useCrud({
  api: inquiryTemplateApi,
  emptyForm,
  serverPagination: true,
  extraParams: () => ({ category: category.value }),
  buildPayload: (): Record<string, any> => ({
    ...(form as any),
    translations: translationsToPayload(translations.value, TRANS_KEYS, TRANS_LANGS),
  }),
  beforeSave: async () => {
    const valid = await formRef.value?.validate().catch(() => false)
    return !!valid
  },
})

function openCreateDialog() { crudCreate(); translations.value = createEmptyTranslations(TRANS_KEYS) }
function openEditDialog(row: any) { crudEdit(row); translations.value = translationsToRecord(row.translations || [], TRANS_KEYS) }

async function handleDelete(row: any) {
  await ElMessageBox.confirm('确定删除此询盘模板吗？', '警告', { type: 'warning' })
  await inquiryTemplateApi.delete(row.id)
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
