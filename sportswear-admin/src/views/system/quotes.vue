<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-select v-model="status" placeholder="状态" clearable style="width: 160px" @change="handleSearch">
        <el-option label="草稿" value="draft" />
        <el-option label="已发送" value="sent" />
        <el-option label="已查看" value="viewed" />
        <el-option label="已接受" value="accepted" />
        <el-option label="已拒绝" value="rejected" />
        <el-option label="已过期" value="expired" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建报价</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="quote_number" label="编号" width="180" />
      <el-table-column label="金额" width="120">
        <template #default="{ row }">{{ row.currency }} {{ row.total_price }}</template>
      </el-table-column>
      <el-table-column prop="quantity" label="数量" width="80" />
      <el-table-column prop="moq" label="MOQ" width="80" />
      <el-table-column prop="lead_time" label="交期" width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }"><el-tag size="small">{{ row.status }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadData" />
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑报价' : '新建报价'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="Lead ID" required><el-input v-model="form.lead_id" /></el-form-item>
        <el-form-item label="数量"><el-input-number v-model="form.quantity" :min="0" /></el-form-item>
        <el-form-item label="单价"><el-input-number v-model="form.unit_price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="总价"><el-input-number v-model="form.total_price" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="货币"><el-select v-model="form.currency"><el-option label="USD" value="USD" /><el-option label="EUR" value="EUR" /><el-option label="CNY" value="CNY" /></el-select></el-form-item>
        <el-form-item label="MOQ"><el-input-number v-model="form.moq" :min="0" /></el-form-item>
        <el-form-item label="交期"><el-input v-model="form.lead_time" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" @click="handleSave">保存</el-button></template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { quoteApi } from '@/api'

const items = ref<any[]>([]); const loading = ref(false); const total = ref(0); const page = ref(1); const pageSize = ref(20); const status = ref('')
const dialogVisible = ref(false); const editingId = ref('')
const form = reactive({ lead_id: '', quantity: 0, unit_price: 0, total_price: 0, currency: 'USD', moq: 0, lead_time: '' })
async function loadData() { loading.value = true; try { const r = await quoteApi.list({ page: page.value, pageSize: pageSize.value, status: status.value }); items.value = r.items; total.value = r.total } finally { loading.value = false } }
function handleSearch() { page.value = 1; loadData() }
function openCreateDialog() { editingId.value = ''; Object.assign(form, { lead_id: '', quantity: 0, unit_price: 0, total_price: 0, currency: 'USD', moq: 0, lead_time: '' }); dialogVisible.value = true }
function openEditDialog(row: any) { editingId.value = row.id; Object.assign(form, { lead_id: row.lead_id, quantity: row.quantity, unit_price: row.unit_price, total_price: row.total_price, currency: row.currency, moq: row.moq, lead_time: row.lead_time }); dialogVisible.value = true }
async function handleSave() { try { if (editingId.value) { await quoteApi.update(editingId.value, form) } else { await quoteApi.create(form) }; ElMessage.success('保存成功'); dialogVisible.value = false; loadData() } catch {} }
async function handleDelete(row: any) { await ElMessageBox.confirm('确定删除此报价吗？', '警告', { type: 'warning' }); await quoteApi.delete(row.id); ElMessage.success('已删除'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>