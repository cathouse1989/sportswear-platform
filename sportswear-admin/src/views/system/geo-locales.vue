<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-button type="primary" @click="openCreate">新建映射</el-button>
      <div class="spacer" />
      <span class="hint">国家 → 语言/货币/时区映射，门户 GEO 推荐的依据</span>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="country" label="国家" min-width="150" />
      <el-table-column prop="country_iso2" label="ISO2" width="90" />
      <el-table-column prop="language" label="语言" width="90" />
      <el-table-column prop="currency" label="货币" width="90" />
      <el-table-column prop="timezone" label="时区" min-width="180" />
      <el-table-column label="启用" width="80" align="center">
        <template #default="{ row }"><el-switch v-model="row.is_active" @change="toggle(row)" /></template>
      </el-table-column>
      <el-table-column label="操作" width="140">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑映射' : '新建映射'" width="480px" :close-on-click-modal="false">
      <el-form :model="form" label-width="80px">
        <el-form-item label="国家" required><el-input v-model="form.country" /></el-form-item>
        <el-form-item label="ISO2" required><el-input v-model="form.country_iso2" placeholder="如 US / CN / FR" /></el-form-item>
        <el-form-item label="语言"><el-input v-model="form.language" placeholder="如 en / zh / fr" /></el-form-item>
        <el-form-item label="货币"><el-input v-model="form.currency" placeholder="如 USD / CNY / EUR" /></el-form-item>
        <el-form-item label="时区"><el-input v-model="form.timezone" placeholder="如 America/New_York" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.is_active" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { geoLocaleApi } from '@/api'

const items = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({ country: '', country_iso2: '', language: '', currency: '', timezone: '', is_active: true })

async function load() {
  loading.value = true
  try { items.value = await geoLocaleApi.list() } catch { items.value = [] } finally { loading.value = false }
}
function openCreate() {
  editingId.value = ''
  Object.assign(form, { country: '', country_iso2: '', language: '', currency: '', timezone: '', is_active: true })
  dialogVisible.value = true
}
function openEdit(row: any) {
  editingId.value = row.id
  Object.assign(form, { country: row.country, country_iso2: row.country_iso2, language: row.language, currency: row.currency, timezone: row.timezone, is_active: row.is_active })
  dialogVisible.value = true
}
async function save() {
  if (!form.country.trim() || !form.country_iso2.trim()) { ElMessage.warning('请填写国家和 ISO2'); return }
  saving.value = true
  try {
    await geoLocaleApi.upsert({ ...form })
    ElMessage.success('已保存')
    dialogVisible.value = false
    load()
  } catch { /* handled */ } finally { saving.value = false }
}
async function toggle(row: any) {
  try {
    await geoLocaleApi.upsert({ country: row.country, country_iso2: row.country_iso2, language: row.language, currency: row.currency, timezone: row.timezone, is_active: row.is_active })
  } catch { row.is_active = !row.is_active }
}
async function remove(row: any) {
  await ElMessageBox.confirm(`确定删除「${row.country}」的映射吗？`, '警告', { type: 'warning' })
  await geoLocaleApi.delete(row.id)
  ElMessage.success('已删除')
  load()
}
onMounted(load)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.hint { font-size: 13px; color: #909399; }
</style>