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
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="loadData" />
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑词条' : '新建词条'" width="600px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="Key" required><el-input v-model="form.key" placeholder="如 home.hero_title_1 / home.get_quote" /></el-form-item>
        <el-form-item label="模块"><el-input v-model="form.module" placeholder="如 home / nav / common" /></el-form-item>
        <el-form-item label="文案">
          <el-tabs v-model="activeLang" class="lang-tabs" @tab-change="onLanguageChange">
            <el-tab-pane v-for="l in langTabs" :key="l" :name="l" :label="langLabels[l]">
              <el-input v-model="langValues[l]" type="textarea" :rows="3" :placeholder="'「' + langLabels[l] + '」文案'" />
            </el-tab-pane>
          </el-tabs>
          <div class="lang-tip">切换语言标签将自动保存当前语言文案，并载入该 Key 对应语言的已有文案，方便逐语言翻页配置。</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button @click="handleSaveNext">保存并下一个语言</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { i18nApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'

const items = ref<any[]>([]); const loading = ref(false); const total = ref(0); const page = ref(1); const pageSize = useAdminPageSize(); const keyword = ref(''); const lang = ref('')
const dialogVisible = ref(false); const editingId = ref(''); const form = reactive({ key: '', module: '' })
// 语言翻页：为同一 Key 逐语言编辑文案（en/zh/es/fr）
const langTabs = ['en', 'zh', 'es', 'fr']
const langLabels: Record<string, string> = { en: 'English', zh: '中文', es: 'Español', fr: 'Français' }
const activeLang = ref('en')
let prevLang = 'en'
const langValues = reactive<Record<string, string>>({ en: '', zh: '', es: '', fr: '' })
async function loadData() { loading.value = true; try { const r = await i18nApi.entries({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, language: lang.value }); items.value = r.items; total.value = r.total } finally { loading.value = false } }
function handleSearch() { page.value = 1; loadData() }
function openCreateDialog() {
  editingId.value = ''
  Object.assign(form, { key: '', module: '' })
  langTabs.forEach(l => { langValues[l] = '' })
  activeLang.value = 'en'
  prevLang = 'en'
  dialogVisible.value = true
}
async function openEditDialog(row: any) {
  editingId.value = row.id
  form.key = row.key
  form.module = row.module || ''
  langTabs.forEach(l => { langValues[l] = '' })
  // 预载该 Key 所有语言的已有文案
  try {
    const r = await i18nApi.entries({ page: 1, pageSize: 50, keyword: row.key, language: '' })
    for (const it of (r.items || [])) {
      if (it.key === row.key) langValues[it.language] = it.value || ''
    }
  } catch { /* handled */ }
  activeLang.value = row.language || 'en'
  prevLang = activeLang.value
  dialogVisible.value = true
}
/** 保存指定语言词条（按 key+language upsert），value 允许为空以清空 */
async function saveLang(l: string) {
  if (!form.key.trim()) { ElMessage.warning('请先填写 Key'); return false }
  try {
    await i18nApi.upsert({ key: form.key.trim(), language: l, value: (langValues[l] || '').trim(), module: form.module.trim() })
    return true
  } catch { return false }
}
/** 切换语言：自动保存上一语言 */
async function onLanguageChange(next: string | number) {
  const to = String(next)
  const old = prevLang
  if (old && old !== to) await saveLang(old)
  prevLang = to
  activeLang.value = to
}
/** 保存当前语言并关闭 */
async function handleSave() {
  const ok = await saveLang(activeLang.value)
  if (!ok) return
  ElMessage.success('保存成功')
  dialogVisible.value = false
  loadData()
}
/** 保存当前语言并切换到下一个语言（循环翻页） */
async function handleSaveNext() {
  const ok = await saveLang(activeLang.value)
  if (!ok) return
  const idx = langTabs.indexOf(activeLang.value)
  const next = langTabs[(idx + 1) % langTabs.length]
  prevLang = activeLang.value
  activeLang.value = next
  ElMessage.success(`已保存 ${langLabels[prevLang]} 文案，继续编辑 ${langLabels[next]}`)
  loadData()
}
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.lang-tabs { margin-bottom: 8px; }
.lang-tip { margin-top: 8px; font-size: 12px; color: #8b7d6b; line-height: 1.5; }
</style>