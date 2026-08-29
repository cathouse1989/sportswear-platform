<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索名称" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建认证</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" min-width="180" />
      <el-table-column prop="code" label="编号" width="140" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }"><el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="220">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="success" @click="handlePublish(row)" v-if="row.status !== 'published'">发布</el-button>
          <el-button size="small" type="warning" @click="handleUnpublish(row)" v-else>下线</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @size-change="handleSearch" />
    <ProFormDialog v-model="dialogVisible" :title="editingId ? '编辑认证' : '新建认证'" :form="form" :rules="rules" @submit="handleSave">
      <el-form-item label="名称" prop="name"><el-input v-model="form.name" /></el-form-item>
      <el-form-item label="编号"><el-input v-model="form.code" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
    </ProFormDialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormRules } from 'element-plus'
import { certificationApi } from '@/api'
import ProFormDialog from '@/components/pro/ProFormDialog.vue'

const allItems = ref<any[]>([])
const items = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  const filtered = kw ? allItems.value.filter((i: any) => (i.name || '').toLowerCase().includes(kw) || (i.code || '').toLowerCase().includes(kw)) : allItems.value
  return filtered.slice((page.value - 1) * pageSize.value, page.value * pageSize.value)
})
const total = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  return kw ? allItems.value.filter((i: any) => (i.name || '').toLowerCase().includes(kw) || (i.code || '').toLowerCase().includes(kw)).length : allItems.value.length
})
const loading = ref(false); const page = ref(1); const pageSize = ref(20); const keyword = ref('')
async function loadData() { loading.value = true; try { allItems.value = await certificationApi.list() } finally { loading.value = false } }
const dialogVisible = ref(false); const editingId = ref(''); const form = reactive({ name: '', code: '', description: '' })
const rules: FormRules = { name: [{ required: true, message: '请输入名称', trigger: 'blur' }] }
function handleSearch() { page.value = 1; loadData() }
function openCreateDialog() { editingId.value = ''; Object.assign(form, { name: '', code: '', description: '' }); dialogVisible.value = true }
function openEditDialog(row: any) { editingId.value = row.id; Object.assign(form, { name: row.name, code: row.code, description: row.description }); dialogVisible.value = true }
async function handleSave() { try { if (editingId.value) { await certificationApi.update(editingId.value, form) } else { await certificationApi.create(form) }; ElMessage.success('保存成功'); dialogVisible.value = false; loadData() } catch {} }
async function handlePublish(row: any) { await certificationApi.publish(row.id); ElMessage.success('已发布'); loadData() }
async function handleUnpublish(row: any) { await certificationApi.unpublish(row.id); ElMessage.success('已下线'); loadData() }
async function handleDelete(row: any) { await ElMessageBox.confirm(`确定删除 ${row.name} 吗？`, '警告', { type: 'warning' }); await certificationApi.delete(row.id); ElMessage.success('已删除'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>