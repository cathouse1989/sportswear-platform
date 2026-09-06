<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索名称 / Slug" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建分类</el-button>
    </div>
    <el-table :data="categories" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" min-width="180" />
      <el-table-column prop="slug" label="Slug" min-width="160" />
      <el-table-column prop="sort_order" label="排序" width="80" align="center" />
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'" size="small">{{ row.is_active ? '启用' : '停用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="handleSearch" @size-change="handleSearch" />

    <!-- 新建/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑分类' : '新建分类'" width="560px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" maxlength="200" show-word-limit placeholder="分类名称" />
        </el-form-item>
        <el-form-item label="Slug" prop="slug">
          <el-input v-model="form.slug" maxlength="200" placeholder="URL 标识（如 yoga-wear）" />
        </el-form-item>
        <el-form-item label="父级分类">
          <el-select v-model="form.parent_id" clearable placeholder="无（顶级分类）" style="width: 100%">
            <el-option v-for="c in parentOptions" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="图片">
          <MediaPicker v-model="form.image" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.is_active" active-text="启用" />
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
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { categoryApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import type { Category } from '@/types'
import MediaPicker from '@/components/media/MediaPicker.vue'

// 分类列表接口不分页，前端自行分页筛选
const allCategories = ref<Category[]>([])
const loading = ref(false)
const saving = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')

const categories = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  const filtered = kw
    ? allCategories.value.filter(c => c.name.toLowerCase().includes(kw) || c.slug.toLowerCase().includes(kw))
    : allCategories.value
  total.value = filtered.length
  const start = (page.value - 1) * pageSize.value
  return filtered.slice(start, start + pageSize.value)
})

const parentOptions = computed(() => allCategories.value.filter(c => c.id !== editingId.value))

const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({ name: '', slug: '', parent_id: '', sort_order: 0, image: '', is_active: true })
const rules: FormRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }],
  slug: [
    { required: true, message: '请输入 Slug', trigger: 'blur' },
    { pattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/, message: '仅支持小写字母、数字和连字符', trigger: 'blur' },
  ],
}

async function loadData() {
  loading.value = true
  try {
    allCategories.value = await categoryApi.list()
  } finally { loading.value = false }
}
function handleSearch() { page.value = 1 }
function resetForm() { Object.assign(form, { name: '', slug: '', parent_id: '', sort_order: 0, image: '', is_active: true }) }
function openCreateDialog() { editingId.value = ''; resetForm(); dialogVisible.value = true }
function openEditDialog(row: Category) {
  editingId.value = row.id
  Object.assign(form, {
    name: row.name,
    slug: row.slug,
    parent_id: row.parent_id || '',
    sort_order: row.sort_order ?? 0,
    image: row.image || '',
    is_active: row.is_active ?? true,
  })
  dialogVisible.value = true
}
async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    const data = { ...form, parent_id: form.parent_id || null }
    if (editingId.value) { await categoryApi.update(editingId.value, data) } else { await categoryApi.create(data) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch {} finally { saving.value = false }
}
async function handleDelete(row: Category) {
  await ElMessageBox.confirm(`确定删除分类「${row.name}」吗？`, '确认删除', { type: 'warning', confirmButtonText: '确定删除', cancelButtonText: '取消' })
  await categoryApi.delete(row.id); ElMessage.success('已删除'); loadData()
}
onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>