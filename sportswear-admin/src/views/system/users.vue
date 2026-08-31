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
      <el-button v-permission="'user:create'" type="primary" @click="openCreateDialog">新建用户</el-button>
    </div>

    <el-table :data="users" v-loading="loading" stripe>
      <el-table-column prop="email" label="邮箱" min-width="200" />
      <el-table-column prop="name" label="姓名" min-width="120" />
      <el-table-column label="角色" min-width="240">
        <template #default="{ row }">
          <template v-if="row.roles?.length">
            <el-tag
              v-for="r in row.roles"
              :key="r.id"
              size="small"
              class="role-tag"
              :type="r.is_active === false ? 'info' : r.code === 'super_admin' ? 'danger' : 'primary'"
            >
              {{ r.name }}
            </el-tag>
          </template>
          <el-tag v-else size="small" type="warning">未分配角色</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
            {{ row.status === 'active' ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最后登录" width="170">
        <template #default="{ row }">
          <span v-if="row.last_login_at">{{ formatTime(row.last_login_at) }}</span>
          <span v-else>—</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }: any">
          <el-button v-permission="'user:update'" size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button
            v-permission="'user:delete'"
            size="small"
            type="danger"
            :disabled="isSelf(row)"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑用户' : '新建用户'" width="580px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="邮箱" required>
          <el-input v-model="form.email" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="密码" :required="!editingId">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            :placeholder="editingId ? '留空则不修改密码（如需重置请填写新密码）' : '设置登录密码'"
          />
        </el-form-item>
        <el-form-item label="角色">
          <el-select
            v-model="form.role_ids"
            multiple
            clearable
            placeholder="选择角色（不选则账号无任何后台权限）"
            style="width: 100%"
          >
            <el-option
              v-for="r in roleOptions"
              :key="r.id"
              :label="r.name + (r.is_active === false ? '（已禁用）' : '')"
              :value="r.id"
            />
          </el-select>
          <div class="form-tip">角色配置决定该账号可见菜单与操作按钮，保存后即时生效</div>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="form.status"
            active-value="active"
            inactive-value="disabled"
            :disabled="editingId === currentUserId"
          />
          <div v-if="editingId === currentUserId" class="form-tip">不能禁用当前登录账号</div>
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
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { roleApi, userApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import type { Role, User } from '@/types'

const authStore = useAuthStore()
const currentUserId = computed(() => authStore.user?.id || '')

const users = ref<User[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')

const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({
  email: '',
  name: '',
  password: '',
  status: 'active',
  role_ids: [] as string[],
})

// 角色选项（账号-角色闭环：创建/编辑用户时分配角色）
const roles = ref<Role[]>([])
const roleOptions = computed(() => roles.value)

// 无 role:manage 权限时 GET /admin/roles 也会通过（后端已放开给 user:create/update），
// 读取失败则仅能编辑已有关联角色（从 row.roles 回填）
async function loadRoles() {
  try {
    const list = await roleApi.list()
    roles.value = list || []
  } catch {
    roles.value = []
  }
}

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

function formatTime(iso?: string) {
  if (!iso) return '—'
  const d = new Date(iso)
  return isNaN(d.getTime()) ? '—' : d.toLocaleString()
}

function isSelf(row: User) {
  return !!currentUserId.value && row.id === currentUserId.value
}

function openCreateDialog() {
  editingId.value = ''
  Object.assign(form, { email: '', name: '', password: '', status: 'active', role_ids: [] })
  dialogVisible.value = true
}

function openEditDialog(row: User) {
  editingId.value = row.id
  const roleIds = (row.roles || []).map((r) => r.id)
  Object.assign(form, { email: row.email, name: row.name, password: '', status: row.status, role_ids: roleIds })
  // 已关联的禁用角色也要保留为已选项（防止选中项显示为原始 ID）
  for (const r of row.roles || []) {
    if (!roles.value.some((x) => x.id === r.id)) {
      roles.value.push(r)
    }
  }
  dialogVisible.value = true
}

async function handleSave() {
  if (!form.email || !form.email.includes('@')) {
    ElMessage.warning('请输入有效的邮箱')
    return
  }
  if (!editingId.value && !form.password) {
    ElMessage.warning('请设置登录密码')
    return
  }
  try {
    if (editingId.value) {
      await userApi.update(editingId.value, {
        email: form.email,
        name: form.name,
        password: form.password || '',
        status: form.status,
        role_ids: form.role_ids,
      })
    } else {
      await userApi.create({
        email: form.email,
        name: form.name,
        password: form.password,
        status: form.status,
        role_ids: form.role_ids,
      })
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadData()
  } catch {
    // 错误已提示
  }
}

async function handleDelete(row: User) {
  if (isSelf(row)) {
    ElMessage.warning('不能删除当前登录账号')
    return
  }
  const isSuper = (row.roles || []).some((r) => r.code === 'super_admin')
  await ElMessageBox.confirm(
    isSuper
      ? `该账号是超级管理员，删除后系统必须仍有至少一个超级管理员可用。确定删除 ${row.email} 吗？`
      : `确定删除用户 ${row.email} 吗？`,
    '警告',
    { type: 'warning' }
  )
  await userApi.delete(row.id)
  ElMessage.success('已删除')
  loadData()
}

onMounted(() => {
  loadData()
  loadRoles()
})
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
.role-tag {
  margin-right: 6px;
}
.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
  margin-top: 4px;
}
</style>