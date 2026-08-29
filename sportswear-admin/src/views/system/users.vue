<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input
        v-model="keyword"
        placeholder="搜索邮箱 / 姓名"
        clearable
        style="width: 240px"
        @keyup.enter="handleSearch"
      />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建用户</el-button>
    </div>

    <el-table :data="users" v-loading="loading" stripe>
      <el-table-column prop="email" label="邮箱" min-width="200" />
      <el-table-column prop="name" label="姓名" min-width="120" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
            {{ row.status === 'active' ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button v-permission="'user:delete'" size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[10, 20, 50, 100]"
      @current-change="loadData"
      @size-change="loadData"
    />

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑用户' : '新建用户'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="邮箱" required>
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="密码" :required="!editingId">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" active-value="active" inactive-value="disabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userApi } from '@/api'
import type { User } from '@/types'

const users = ref<User[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')

const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({
  email: '',
  name: '',
  password: '',
  status: 'active',
})

async function loadData() {
  loading.value = true
  try {
    const result = await userApi.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value })
    users.value = result.items
    total.value = result.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

function openCreateDialog() {
  editingId.value = ''
  Object.assign(form, { email: '', name: '', password: '', status: 'active' })
  dialogVisible.value = true
}

function openEditDialog(row: User) {
  editingId.value = row.id
  Object.assign(form, { email: row.email, name: row.name, password: '', status: row.status })
  dialogVisible.value = true
}

async function handleSave() {
  try {
    if (editingId.value) {
      await userApi.update(editingId.value, { ...form })
    } else {
      await userApi.create({ ...form })
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadData()
  } catch {
    // 错误已提示
  }
}

async function handleDelete(row: User) {
  await ElMessageBox.confirm(`确定删除用户 ${row.email} 吗？`, '警告', { type: 'warning' })
  await userApi.delete(row.id)
  ElMessage.success('已删除')
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
}
.spacer {
  flex: 1;
}
.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>