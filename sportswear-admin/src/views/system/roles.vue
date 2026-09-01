<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input
        v-model="keyword"
        placeholder="搜索名称 / 标识"
        clearable
        style="width: 240px"
        @keyup.enter="handleSearch"
      />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button v-permission="'role:manage'" type="primary" @click="openCreateDialog">新建角色</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" width="200">
        <template #default="{ row }">
          <span>{{ row.name }}</span>
          <el-tag v-if="isBuiltin(row)" size="small" type="info" class="builtin-tag">内置</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="code" label="标识" width="160" />
      <el-table-column prop="description" label="描述" min-width="220" />
      <el-table-column label="启用" width="90">
        <template #default="{ row }">
          <el-switch
            :model-value="row.is_active"
            :disabled="row.code === 'super_admin'"
            @change="(val: string | number | boolean) => handleToggle(row, val === true)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }: any">
          <el-button v-permission="'role:manage'" size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button
            v-permission="'role:manage'"
            size="small"
            type="danger"
            :disabled="isBuiltin(row)"
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
      :page-sizes="[10,20,50,100]"
      @current-change="loadData"
      @size-change="loadData"
    />

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑角色' : '新建角色'" width="700px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" maxlength="100" />
        </el-form-item>
        <el-form-item label="标识" required>
          <el-input v-model="form.code" :disabled="editingCode === 'super_admin'" maxlength="50" />
          <div v-if="editingCode === 'super_admin'" class="form-tip">内置超级管理员标识不可修改</div>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" maxlength="500" />
        </el-form-item>
        <el-form-item v-if="editingId" label="启用">
          <el-switch v-model="form.is_active" :disabled="editingCode === 'super_admin'" />
          <div v-if="editingCode === 'super_admin'" class="form-tip">内置超级管理员角色不可禁用</div>
        </el-form-item>
        <el-form-item label="权限">
          <el-tree
            ref="treeRef"
            :data="permTree"
            :props="{ label: 'name', children: 'children' }"
            show-checkbox
            node-key="code"
            :default-checked-keys="currentPerms"
            class="perm-tree"
          />
          <div class="form-tip">勾选父级分组会同时选中其下全部权限；保存时父级+子级权限都会提交</div>
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
import { nextTick, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { roleApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import type { Permission, Role } from '@/types'

const items = ref<Role[]>([])
const loading = ref(false)
const keyword = ref('')
const page = ref(1)
const pageSize = useAdminPageSize()
const total = ref(0)

// 新建/编辑角色对话框
const dialogVisible = ref(false)
const editingId = ref('')
const editingCode = ref('')
const form = reactive({ name: '', code: '', description: '', is_active: true })
const treeRef = ref<any>(null)
const currentPerms = ref<string[]>([])

const permTree = ref<any[]>([])
// 内置角色集合（默认角色。super_admin 受后端 + 前端双重保护，不可禁用/改标识）
const builtinRoles = [
  'super_admin',
  'admin',
  'content_admin',
  'product_admin',
  'seo_admin',
  'sales_manager',
  'sales',
  'editor',
  'viewer',
]

function isBuiltin(row: any) {
  return builtinRoles.includes(row.code)
}

async function loadData() {
  loading.value = true
  try {
    const result = await roleApi.listPage({ page: page.value, pageSize: pageSize.value, keyword: keyword.value })
    items.value = result.items
    total.value = result.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

async function loadPermTree() {
  const perms = await roleApi.permissions()
  const groups: Record<string, any> = {}
  perms.forEach((p: Permission) => {
    if (!groups[p.module]) groups[p.module] = { name: moduleLabel(p.module), children: [] }
    groups[p.module].children.push({ id: p.id, name: p.name, code: p.code })
  })
  permTree.value = Object.values(groups)
}

function moduleLabel(module: string) {
  const map: Record<string, string> = {
    user: '用户与角色',
    product: '产品',
    cms: '内容',
    media: '媒体',
    lead: '询盘',
    seo: 'SEO',
    system: '系统',
  }
  return map[module] || module
}

let permTreeLoaded = false
async function loadPermTreeIfNeeded() {
  if (permTreeLoaded) return
  await loadPermTree()
  permTreeLoaded = true
}

// 汇总选中的权限码 = 勾选 + 半选（避免仅 getCheckedKeys 时父级勾选丢子权限）
function collectCheckedCodes(tree: any): string[] {
  if (!tree) return []
  const checked: string[] = tree.getCheckedKeys(false) || []
  const half: string[] = tree.getHalfCheckedKeys() || []
  return Array.from(new Set([...checked, ...half]))
}

// 权限码 -> 权限 ID（用于后端 CreateRole 的 permission_ids）
function codesToIds(codes: string[]): string[] {
  const map = new Map<string, string>()
  for (const group of permTree.value) {
    for (const child of group.children || []) map.set(child.code, child.id)
  }
  return codes.map((code) => map.get(code) || '').filter(Boolean)
}

function openCreateDialog() {
  editingId.value = ''
  editingCode.value = ''
  currentPerms.value = []
  Object.assign(form, { name: '', code: '', description: '', is_active: true })
  dialogVisible.value = true
  nextTick(async () => {
    await loadPermTreeIfNeeded()
    treeRef.value?.setCheckedKeys([])
  })
}

async function openEditDialog(row: any) {
  editingId.value = row.id
  editingCode.value = row.code
  currentPerms.value = (row.permissions || []).map((p: Permission) => p.code)
  Object.assign(form, {
    name: row.name,
    code: row.code,
    description: row.description,
    is_active: row.is_active,
  })
  dialogVisible.value = true
  nextTick(async () => {
    await loadPermTreeIfNeeded()
    treeRef.value?.setCheckedKeys(currentPerms.value)
  })
}

async function handleSave() {
  if (!form.name || !form.code) {
    ElMessage.warning('请填写名称与标识')
    return
  }
  try {
    const permissionCodes = collectCheckedCodes(treeRef.value)
    if (editingId.value) {
      await roleApi.update(editingId.value, {
        name: form.name,
        code: form.code,
        description: form.description,
        is_active: form.is_active,
        permissions: permissionCodes,
      })
    } else {
      await roleApi.create({
        name: form.name,
        code: form.code,
        description: form.description,
        permission_ids: codesToIds(permissionCodes),
      })
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadData()
  } catch {
    // 错误已提示
  }
}

async function handleToggle(row: any, val: boolean) {
  if (row.code === 'super_admin' && !val) {
    ElMessage.warning('不能禁用超级管理员角色')
    return
  }
  await roleApi.updateStatus(row.id, val)
  ElMessage.success(val ? '已启用' : '已禁用')
  loadData()
}

async function handleDelete(row: any) {
  if (isBuiltin(row)) {
    ElMessage.warning('内置角色不可删除')
    return
  }
  await ElMessageBox.confirm(
    `确定删除角色「${row.name}」吗？删除后该角色将从所有账号解除，其权限配置一并清除。`,
    '警告',
    { type: 'warning' }
  )
  await roleApi.delete(row.id)
  ElMessage.success('已删除')
  const lastPage = Math.max(1, Math.ceil((total.value - 1) / pageSize.value))
  if (page.value > lastPage) page.value = lastPage
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
.builtin-tag {
  margin-left: 8px;
}
.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
  margin-top: 4px;
}
.perm-tree {
  width: 100%;
  max-height: 320px;
  overflow: auto;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 8px;
}
</style>