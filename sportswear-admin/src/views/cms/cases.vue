<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索标题" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建案例</el-button>
    </div>
    <el-table :data="cases" v-loading="loading" stripe>
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column prop="client_industry" label="客户行业" width="140" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }"><el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="180">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadData" />

    <!-- 编辑弹窗（缺陷 B-01 修复：此前模板缺失导致新建/编辑点击无响应） -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑案例' : '新建案例'" width="640px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" maxlength="200" show-word-limit />
        </el-form-item>
        <el-form-item label="Slug" prop="slug">
          <el-input v-model="form.slug" placeholder="URL 标识，如 yoga-brand-oem" />
        </el-form-item>
        <el-form-item label="客户行业">
          <el-input v-model="form.client_industry" placeholder="如 Fitness Brand / Retail Chain" />
        </el-form-item>
        <el-form-item label="解决方案" prop="solution">
          <el-input v-model="form.solution" type="textarea" :rows="6" />
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
import { cmsApi } from '@/api/cms'
import type { Case } from '@/types'

const cases = ref<Case[]>([])
const loading = ref(false)
const saving = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')

const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({ title: '', slug: '', client_industry: '', solution: '' })
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
    const result = await cmsApi.cases.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value })
    cases.value = result.items
    total.value = result.total
  } finally { loading.value = false }
}
function handleSearch() { page.value = 1; loadData() }
function resetForm() { Object.assign(form, { title: '', slug: '', client_industry: '', solution: '' }) }
function openCreateDialog() { editingId.value = ''; resetForm(); dialogVisible.value = true }
function openEditDialog(row: Case) {
  editingId.value = row.id
  Object.assign(form, { title: row.title, slug: row.slug, client_industry: row.client_industry || '', solution: row.solution || '' })
  dialogVisible.value = true
}
async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (editingId.value) { await cmsApi.cases.update(editingId.value, form) } else { await cmsApi.cases.create(form) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch {} finally { saving.value = false }
}
async function handleDelete(row: Case) { await ElMessageBox.confirm(`确定删除案例 ${row.title} 吗？`, '警告', { type: 'warning' }); await cmsApi.cases.delete(row.id); ElMessage.success('已删除'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>