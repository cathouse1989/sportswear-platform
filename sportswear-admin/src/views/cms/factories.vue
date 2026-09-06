<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索工厂名称" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建工厂</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column label="图片" width="76" align="center">
        <template #default="{ row }">
          <el-image v-if="row.image" :src="row.image" fit="cover" class="cover-thumb" :preview-src-list="[row.image]" preview-teleported />
          <div v-else class="cover-placeholder">—</div>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="180" />
      <el-table-column label="翻译" min-width="120">
        <template #default="{ row }">
          <el-tag v-for="t in row.translations || []" :key="t.id || t.language" size="small" type="info" class="lang-tag">{{ t.language }}</el-tag>
          <span v-if="!(row.translations || []).length" class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column prop="location" label="位置" min-width="160" />
      <el-table-column prop="employees" label="员工数" width="100" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }"><el-tag :type="enumTag('content.status', row.status)" size="small">{{ enumLabel('content.status', row.status) }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="220">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="success" @click="handlePublish(row)" v-if="row.status !== 'published'">发布</el-button>
          <el-button size="small" type="warning" @click="handleUnpublish(row)" v-else>下线</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="handleSearch" />
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑工厂' : '新建工厂'" width="680px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称" required><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="图片"><MediaPicker v-model="form.image" /></el-form-item>
        <el-form-item label="位置"><el-input v-model="form.location" /></el-form-item>
        <el-form-item label="员工数"><el-input-number v-model="form.employees" :min="0" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
        <el-divider content-position="left">多语言翻译（未填写回退英文）</el-divider>
        <TransEditor v-model="translations" :fields="TRANS_FIELDS" :source="{ name: form.name, description: form.description }" :langs="TRANS_LANGS" source-label="English" />
      </el-form>
      <template #footer><el-button @click="dialogVisible = false">取消</el-button><el-button type="primary" @click="handleSave">保存</el-button></template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { factoryApi } from '@/api'
import MediaPicker from '@/components/media/MediaPicker.vue'
import TransEditor from '@/components/cms/TransEditor.vue'
import { useCrud } from '@/composables/useCrud'
import { useEnumDict } from '@/composables/useEnumDict'

// 状态标签由「字典管理」驱动（content.status）
const { ensureLoaded: loadEnumDict, label: enumLabel, tagType: enumTag } = useEnumDict()

const TRANS_LANGS = [
  { value: 'zh', label: '中文' },
  { value: 'es', label: 'Español' },
  { value: 'fr', label: 'Français' },
]
const TRANS_FIELDS = [
  { key: 'name', label: '名称' },
  { key: 'description', label: '描述', type: 'textarea' as const, rows: 2 },
]
const emptyForm = () => ({ name: '', image: '', location: '', employees: 0, description: '' })
const translations = ref<Record<string, Record<string, string>>>({})
function emptyTranslations(): Record<string, Record<string, string>> {
  return { zh: { name: '', description: '' }, es: { name: '', description: '' }, fr: { name: '', description: '' } }
}
function translationsToRecord(list: any[]) {
  const next = emptyTranslations()
  for (const t of list || []) {
    const lang = t.language
    if (lang && lang !== 'en' && next[lang]) {
      next[lang] = { name: t.name || '', description: t.description || '' }
    }
  }
  return next
}

const {
  items, total, loading, page, pageSize, keyword,
  loadData, handleSearch,
  dialogVisible, editingId, form,
  handleSave, handlePublish, handleUnpublish, handleDelete,
  openCreateDialog: crudCreate, openEditDialog: crudEdit,
} = useCrud({
  api: factoryApi,
  emptyForm,
  serverPagination: true,
  buildPayload: (): Record<string, any> => ({
    ...(form as any),
    translations: Object.entries(translations.value)
      .filter(([, t]) => (t.name || '').trim() || (t.description || '').trim())
      .map(([language, t]) => ({ language, name: t.name || '', description: t.description || '' })),
  }),
})
function openCreateDialog() { crudCreate(); translations.value = emptyTranslations() }
function openEditDialog(row: any) { crudEdit(row); translations.value = translationsToRecord(row.translations || []) }
onMounted(() => { loadEnumDict(); loadData() })
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.cover-thumb { width: 44px; height: 44px; border-radius: 4px; }
.cover-placeholder { width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; color: #c0c4cc; background: #f5f7fa; border-radius: 4px; }
.lang-tag { margin-right: 4px; }
.muted { color: #c0c4cc; }
</style>