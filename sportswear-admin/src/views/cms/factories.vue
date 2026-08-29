<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索工厂名称" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建工厂</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" min-width="180" />
      <el-table-column prop="location" label="位置" min-width="160" />
      <el-table-column prop="employees" label="员工数" width="100" />
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
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="(val: number) => { page = val }" @size-change="handleSearch" />
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑工厂' : '新建工厂'" width="560px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="位置"><el-input v-model="form.location" /></el-form-item>
        <el-form-item label="员工数"><el-input-number v-model="form.employees" :min="0" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" @click="handleSave">保存</el-button></template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { factoryApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'

const allItems = ref<any[]>([])
const items = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  const filtered = kw ? allItems.value.filter((i: any) => (i.name || '').toLowerCase().includes(kw)) : allItems.value
  return filtered.slice((page.value - 1) * pageSize.value, page.value * pageSize.value)
})
const total = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  return kw ? allItems.value.filter((i: any) => (i.name || '').toLowerCase().includes(kw)).length : allItems.value.length
})
const loading = ref(false)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({ name: '', location: '', employees: 0, description: '' })

async function loadData() { loading.value = true; try { allItems.value = await factoryApi.list() } finally { loading.value = false } }
function handleSearch() { page.value = 1 }
function openCreateDialog() { editingId.value = ''; Object.assign(form, { name: '', location: '', employees: 0, description: '' }); dialogVisible.value = true }
function openEditDialog(row: any) { editingId.value = row.id; Object.assign(form, { name: row.name, location: row.location, employees: row.employees, description: row.description }); dialogVisible.value = true }
async function handleSave() { try { if (editingId.value) { await factoryApi.update(editingId.value, form) } else { await factoryApi.create(form) }; ElMessage.success('保存成功'); dialogVisible.value = false; loadData() } catch {} }
async function handlePublish(row: any) { await factoryApi.publish(row.id); ElMessage.success('已发布'); loadData() }
async function handleUnpublish(row: any) { await factoryApi.unpublish(row.id); ElMessage.success('已下线'); loadData() }
async function handleDelete(row: any) { await ElMessageBox.confirm(`确定删除 ${row.name} 吗？`, '警告', { type: 'warning' }); await factoryApi.delete(row.id); ElMessage.success('已删除'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>