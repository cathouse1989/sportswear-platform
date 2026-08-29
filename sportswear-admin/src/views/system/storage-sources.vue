<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-button type="primary" @click="openCreateDialog">新建存储源</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" width="160" />
      <el-table-column prop="code" label="标识" width="100" />
      <el-table-column prop="protocol" label="协议" width="70" />
      <el-table-column prop="host" label="主机" min-width="160" />
      <el-table-column prop="port" label="端口" width="70" />
      <el-table-column prop="base_path" label="基础路径" width="140" />
      <el-table-column prop="bucket" label="Bucket" width="120" />
      <el-table-column label="默认" width="70">
        <template #default="{ row }"><el-tag v-if="row.is_default" size="small" type="success">默认</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑存储源' : '新建存储源'" width="640px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="标识" required><el-input v-model="form.code" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type"><el-option label="本地" value="local" /><el-option label="MinIO" value="minio" /><el-option label="外部" value="external" /></el-select>
        </el-form-item>
        <el-form-item label="协议"><el-select v-model="form.protocol"><el-option label="HTTP" value="http" /><el-option label="HTTPS" value="https" /></el-select></el-form-item>
        <el-form-item label="主机" required><el-input v-model="form.host" placeholder="IP 或域名" /></el-form-item>
        <el-form-item label="端口"><el-input v-model="form.port" /></el-form-item>
        <el-form-item label="基础路径"><el-input v-model="form.base_path" /></el-form-item>
        <el-form-item label="Bucket"><el-input v-model="form.bucket" /></el-form-item>
        <el-form-item label="默认"><el-switch v-model="form.is_default" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" @click="handleSave">保存</el-button></template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { storageSourceApi } from '@/api'

const items = ref<any[]>([]); const loading = ref(false); const dialogVisible = ref(false); const editingId = ref('')
const form = reactive({ name: '', code: '', type: 'local', protocol: 'http', host: '', port: '', base_path: '', bucket: '', is_default: false })
async function loadData() { loading.value = true; try { items.value = await storageSourceApi.list() } finally { loading.value = false } }
function openCreateDialog() { editingId.value = ''; Object.assign(form, { name: '', code: '', type: 'local', protocol: 'http', host: '', port: '', base_path: '', bucket: '', is_default: false }); dialogVisible.value = true }
function openEditDialog(row: any) { editingId.value = row.id; Object.assign(form, { name: row.name, code: row.code, type: row.type, protocol: row.protocol, host: row.host, port: row.port, base_path: row.base_path, bucket: row.bucket, is_default: row.is_default }); dialogVisible.value = true }
async function handleSave() { try { if (editingId.value) { await storageSourceApi.update(editingId.value, form) } else { await storageSourceApi.create(form) }; ElMessage.success('保存成功'); dialogVisible.value = false; loadData() } catch {} }
async function handleDelete(row: any) { await ElMessageBox.confirm(`确定删除存储源 ${row.name} 吗？`, '警告', { type: 'warning' }); await storageSourceApi.delete(row.id); ElMessage.success('已删除'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
</style>