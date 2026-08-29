<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-button type="primary" @click="openCreateDialog">新建角色</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" width="160" />
      <el-table-column prop="code" label="标识" width="140" />
      <el-table-column prop="description" label="描述" min-width="200" />
      <el-table-column label="启用" width="80">
        <template #default="{ row }">
          <el-switch :model-value="row.is_active" @change="(val: string | number | boolean) => handleToggle(row, val === true)" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" @click="openPermissionDialog(row)">权限</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑角色' : '新建角色'" width="520px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="标识" required><el-input v-model="form.code" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" @click="handleSave">保存</el-button></template>
    </el-dialog>
    <el-dialog v-model="permVisible" title="权限配置" width="600px">
      <el-tree :data="permTree" :props="{ label: 'name', children: 'children' }" show-checkbox node-key="code" :default-checked-keys="currentPerms" ref="treeRef" />
      <template #footer><el-button @click="permVisible = false">取消</el-button><el-button type="primary" @click="handleSavePerms">保存</el-button></template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { roleApi } from '@/api'

const items = ref<any[]>([]); const loading = ref(false); const dialogVisible = ref(false); const editingId = ref('')
const form = reactive({ name: '', code: '', description: '' })
const permVisible = ref(false); const currentRoleId = ref(''); const currentPerms = ref<string[]>([]); const permTree = ref<any[]>([]); const treeRef = ref<any>(null)

async function loadData() { loading.value = true; try { items.value = await roleApi.list() } finally { loading.value = false } }
function openCreateDialog() { editingId.value = ''; Object.assign(form, { name: '', code: '', description: '' }); dialogVisible.value = true }
function openEditDialog(row: any) { editingId.value = row.id; Object.assign(form, { name: row.name, code: row.code, description: row.description }); dialogVisible.value = true }
async function handleSave() { try { if (editingId.value) { await roleApi.update(editingId.value, form) } else { await roleApi.create(form) }; ElMessage.success('保存成功'); dialogVisible.value = false; loadData() } catch {} }
async function handleToggle(row: any, val: boolean) { await roleApi.updateStatus(row.id, val); ElMessage.success(val ? '已启用' : '已禁用'); loadData() }
async function openPermissionDialog(row: any) { currentRoleId.value = row.id; currentPerms.value = row.permissions?.map((p: any) => p.code) || []; await loadPermTree(); permVisible.value = true }
async function loadPermTree() { const perms = await roleApi.permissions(); const groups: Record<string, any> = {}; perms.forEach((p: any) => { if (!groups[p.module]) groups[p.module] = { name: p.module, children: [] }; groups[p.module].children.push(p) }); permTree.value = Object.values(groups); setTimeout(() => treeRef.value?.setCheckedKeys(currentPerms.value), 100) }
function handleSavePerms() { const checked = treeRef.value?.getCheckedKeys() || []; roleApi.update(currentRoleId.value, { permissions: checked }).then(() => { ElMessage.success('权限已更新'); permVisible.value = false; loadData() }) }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
</style>