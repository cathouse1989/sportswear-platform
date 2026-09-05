<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索问题" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-select v-model="category" placeholder="分类" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="c in CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
      </el-select>
      <el-select v-model="language" placeholder="语言" clearable style="width: 140px" @change="handleSearch">
        <el-option v-for="l in LANGUAGES" :key="l.value" :label="l.label" :value="l.value" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建 FAQ</el-button>
    </div>
    <el-table :data="faqs" v-loading="loading" stripe>
      <el-table-column prop="question" label="问题" min-width="300" />
      <el-table-column prop="category" label="分类" width="110">
        <template #default="{ row }">{{ categoryLabel(row.category) }}</template>
      </el-table-column>
      <el-table-column prop="language" label="语言" width="80" />
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
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="loadData" />

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
            <el-option v-for="c in CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="语言">
          <el-select v-model="form.language" style="width: 100%">
            <el-option v-for="l in LANGUAGES" :key="l.value" :label="l.label" :value="l.value" />
          </el-select>
        </el-form-item>
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
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { faqApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import type { FAQ } from '@/types'

// 分类统一定义：列表筛选与表单共用
const CATEGORIES = [
  { label: 'MOQ', value: 'moq' },
  { label: 'OEM', value: 'oem' },
  { label: 'ODM', value: 'odm' },
  { label: '打样', value: 'sample' },
  { label: '支付', value: 'payment' },
  { label: '生产', value: 'production' },
  { label: '物流', value: 'logistics' },
  { label: '面料', value: 'fabric' },
  { label: '质量', value: 'quality' },
  { label: '认证', value: 'certification' },
]
function categoryLabel(value: string) { return CATEGORIES.find((c) => c.value === value)?.label || value }

// 与门户支持语言保持一致
const LANGUAGES = [
  { label: 'English', value: 'en' },
  { label: '中文', value: 'zh' },
  { label: 'Español', value: 'es' },
  { label: 'Français', value: 'fr' },
]

const faqs = ref<FAQ[]>([])
const loading = ref(false)
const saving = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const category = ref('')
const language = ref('')

const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({ question: '', answer: '', category: '', language: 'en', sort_order: 0, is_active: true })
const rules: FormRules = {
  question: [{ required: true, message: '请输入问题', trigger: 'blur' }],
  answer: [{ required: true, message: '请输入答案', trigger: 'blur' }],
}

async function loadData() {
  loading.value = true
  try {
    const result = await faqApi.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, category: category.value, language: language.value })
    faqs.value = result.items
    total.value = result.total
  } finally { loading.value = false }
}
function handleSearch() { page.value = 1; loadData() }
function resetForm() { Object.assign(form, { question: '', answer: '', category: '', language: 'en', sort_order: 0, is_active: true }) }
function openCreateDialog() { editingId.value = ''; resetForm(); dialogVisible.value = true }
function openEditDialog(row: FAQ) {
  editingId.value = row.id
  Object.assign(form, { question: row.question, answer: row.answer, category: row.category || '', language: row.language || 'en', sort_order: row.sort_order ?? 0, is_active: row.is_active })
  dialogVisible.value = true
}
async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (editingId.value) { await faqApi.update(editingId.value, form) } else { await faqApi.create(form) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch {} finally { saving.value = false }
}
async function handleDelete(row: FAQ) { await ElMessageBox.confirm('确定删除此 FAQ 吗？', '警告', { type: 'warning' }); await faqApi.delete(row.id); ElMessage.success('已删除'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>