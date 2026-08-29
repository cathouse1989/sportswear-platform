<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索 key / 文案" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-select v-model="lang" placeholder="语言" clearable style="width: 120px" @change="handleSearch">
        <el-option label="English" value="en" />
        <el-option label="中文" value="zh" />
        <el-option label="Español" value="es" />
        <el-option label="Français" value="fr" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建词条</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="key" label="Key" min-width="200" />
      <el-table-column prop="language" label="语言" width="80" />
      <el-table-column prop="value" label="文案" min-width="260" />
      <el-table-column prop="module" label="模块" width="100" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadData" />
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑词条' : '新建词条'" width="560px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="Key" required><el-input v-model="form.key" /></el-form-item>
        <el-form-item label="语言" required>
          <el-select v-model="form.language">
            <el-option label="English" value="en" />
            <el-option label="中文" value="zh" />
            <el-option label="Español" value="es" />
            <el-option label="Français" value="fr" />
          </el-select>
        </el-form-item>
        <el-form-item label="文案" required><el-input v-model="form.value" type="textarea" /></el-form-item>
        <el-form-item label="模块"><el-input v-model="form.module" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" @click="handleSave">保存</el-button></template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { i18nApi } from '@/api'

const items = ref<any[]>([]); const loading = ref(false); const total = ref(0); const page = ref(1); const pageSize = ref(20); const keyword = ref(''); const lang = ref('')
const dialogVisible = ref(false); const editingId = ref(''); const form = reactive({ key: '', language: 'en', value: '', module: '' })
async function loadData() { loading.value = true; try { const r = await i18nApi.entries({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, language: lang.value }); items.value = r.items; total.value = r.total } finally { loading.value = false } }
function handleSearch() { page.value = 1; loadData() }
function openCreateDialog() { editingId.value = ''; Object.assign(form, { key: '', language: 'en', value: '', module: '' }); dialogVisible.value = true }
function openEditDialog(row: any) { editingId.value = row.id; Object.assign(form, { key: row.key, language: row.language, value: row.value, module: row.module }); dialogVisible.value = true }
async function handleSave() { try { await i18nApi.upsert(form); ElMessage.success('保存成功'); dialogVisible.value = false; loadData() } catch {} }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>