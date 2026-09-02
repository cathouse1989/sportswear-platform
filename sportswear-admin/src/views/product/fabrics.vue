<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-button type="primary" @click="openCreateDialog">新建面料</el-button>
    </div>
    <el-table :data="fabricList" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column prop="code" label="编码" width="140" />
      <el-table-column prop="composition" label="成分" min-width="200" />
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 'published' ? 'success' : row.status === 'offline' ? 'warning' : 'info'" size="small">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }: any">
          <el-button v-if="row.status !== 'published'" size="small" type="success" @click="handlePublish(row)">发布</el-button>
          <el-button v-else size="small" type="warning" @click="handleUnpublish(row)">下线</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建弹窗 -->
    <el-dialog v-model="dialogVisible" title="新建面料" width="520px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" maxlength="200" show-word-limit placeholder="面料名称" />
        </el-form-item>
        <el-form-item label="编码" prop="code">
          <el-input v-model="form.code" maxlength="50" placeholder="面料编码（如 POLY-Q235）" />
        </el-form-item>
        <el-form-item label="成分">
          <el-input v-model="form.composition" placeholder="如 90% Polyester, 10% Elastane" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
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
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { fabricApi } from '@/api'
import type { Fabric } from '@/types'

const fabricList = ref<Fabric[]>([])
const loading = ref(false)
const saving = ref(false)

const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const form = reactive({ name: '', code: '', composition: '', description: '' })
const rules: FormRules = {
  name: [{ required: true, message: '请输入面料名称', trigger: 'blur' }],
  code: [
    { required: true, message: '请输入面料编码', trigger: 'blur' },
    { pattern: /^[A-Za-z0-9][A-Za-z0-9_-]*$/, message: '仅支持字母、数字、下划线和连字符', trigger: 'blur' },
  ],
}

function statusLabel(s: string) {
  return s === 'published' ? '已发布' : s === 'offline' ? '已下线' : '草稿'
}

async function loadData() {
  loading.value = true
  try {
    fabricList.value = await fabricApi.list()
  } finally { loading.value = false }
}

function openCreateDialog() { Object.assign(form, { name: '', code: '', composition: '', description: '' }); dialogVisible.value = true }

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    await fabricApi.create({ ...form, is_active: true })
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch {} finally { saving.value = false }
}

async function handlePublish(row: Fabric) {
  await fabricApi.publish(row.id); ElMessage.success('已发布'); loadData()
}
async function handleUnpublish(row: Fabric) {
  await fabricApi.unpublish(row.id); ElMessage.success('已下线'); loadData()
}
onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
</style>