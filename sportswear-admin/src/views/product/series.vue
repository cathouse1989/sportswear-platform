<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-button type="primary" @click="openCreateDialog">新建系列</el-button>
    </div>
    <el-table :data="seriesList" v-loading="loading" stripe>
      <el-table-column prop="name" label="名称" min-width="180" />
      <el-table-column prop="slug" label="Slug" min-width="160" />
      <el-table-column prop="sort_order" label="排序" width="80" align="center" />
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="enumTag('product.status', row.status)" size="small">
            {{ enumLabel('product.status', row.status) }}
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
    <el-dialog v-model="dialogVisible" title="新建系列" width="480px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" maxlength="200" show-word-limit placeholder="系列名称" />
        </el-form-item>
        <el-form-item label="Slug" prop="slug">
          <el-input v-model="form.slug" maxlength="200" placeholder="URL 标识（如 performance）" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
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
import { seriesApi } from '@/api'
import { useEnumDict } from '@/composables/useEnumDict'
import type { Series } from '@/types'

const { ensureLoaded: loadEnumDict, label: enumLabel, tagType: enumTag } = useEnumDict()

const seriesList = ref<Series[]>([])
const loading = ref(false)
const saving = ref(false)

const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const form = reactive({ name: '', slug: '', sort_order: 0 })
const rules: FormRules = {
  name: [{ required: true, message: '请输入系列名称', trigger: 'blur' }],
  slug: [
    { required: true, message: '请输入 Slug', trigger: 'blur' },
    { pattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/, message: '仅支持小写字母、数字和连字符', trigger: 'blur' },
  ],
}

async function loadData() {
  loading.value = true
  try {
    seriesList.value = await seriesApi.list()
  } finally { loading.value = false }
}

function openCreateDialog() { Object.assign(form, { name: '', slug: '', sort_order: 0 }); dialogVisible.value = true }

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    await seriesApi.create({ ...form, is_active: true })
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch {} finally { saving.value = false }
}

async function handlePublish(row: Series) {
  await seriesApi.publish(row.id); ElMessage.success('已发布'); loadData()
}
async function handleUnpublish(row: Series) {
  await seriesApi.unpublish(row.id); ElMessage.success('已下线'); loadData()
}
onMounted(() => { loadEnumDict(); loadData() })
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
</style>