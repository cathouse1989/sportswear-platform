<template>
  <div class="enums-page">
    <div class="layout">
      <!-- 左侧：枚举类型列表 -->
      <el-card shadow="never" class="type-panel">
        <template #header>
          <div class="panel-header">
            <span class="font-semibold">枚举类型</span>
            <el-button size="small" type="primary" @click="openTypeDialog()">新建类型</el-button>
          </div>
        </template>
        <div v-loading="typeLoading" class="type-list">
          <div
            v-for="t in types"
            :key="t.id"
            class="type-item"
            :class="{ active: currentType?.id === t.id }"
            @click="selectType(t)"
          >
            <div class="type-main">
              <span class="type-name">{{ t.name }}</span>
              <el-tag v-if="t.is_system" size="small" type="warning" effect="plain">系统内置</el-tag>
            </div>
            <div class="type-code">{{ t.code }}</div>
            <div class="type-meta">{{ t.module }} · {{ t.item_count }} 项</div>
          </div>
          <el-empty v-if="!typeLoading && types.length === 0" description="暂无枚举类型" :image-size="60" />
        </div>
      </el-card>

      <!-- 右侧：枚举项列表 -->
      <el-card shadow="never" class="item-panel">
        <template #header>
          <div class="panel-header">
            <div>
              <span class="font-semibold">{{ currentType ? currentType.name : '请选择类型' }}</span>
              <span v-if="currentType" class="code-hint">{{ currentType.code }}</span>
              <span v-if="currentType?.i18n_prefix" class="prefix-hint">门户前缀：{{ currentType.i18n_prefix }}</span>
            </div>
            <div class="panel-actions">
              <el-button v-if="currentType" size="small" @click="openTypeDialog(currentType)">编辑类型</el-button>
              <el-button
                v-if="currentType && !currentType.is_system"
                size="small"
                type="danger"
                plain
                @click="handleDeleteType"
              >删除类型</el-button>
              <el-button v-if="currentType" size="small" @click="exportCsv">导出</el-button>
              <el-button v-if="currentType" size="small" @click="triggerImport">导入</el-button>
              <el-button v-if="currentType" size="small" type="primary" @click="openItemDialog()">添加枚举项</el-button>
            </div>
          </div>
        </template>

        <div v-loading="itemLoading">
          <el-empty v-if="!currentType" description="请在左侧选择枚举类型" :image-size="80" />
          <template v-else>
            <div class="filter-row">
              <el-switch v-model="onlyMissing" size="small" active-text="仅看缺译" />
              <span v-if="onlyMissing" class="filter-count">{{ filteredItems.length }} / {{ items.length }} 项缺译</span>
            </div>
            <el-table :data="filteredItems" stripe>
              <el-table-column prop="value" label="存储值" min-width="140">
                <template #default="{ row }"><code class="value-code">{{ row.value }}</code></template>
              </el-table-column>
              <el-table-column label="标签" min-width="140">
                <template #default="{ row }">
                  <el-tag :type="tagType(row.color)" disable-transitions>{{ row.label }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="多语言翻译" min-width="240">
                <template #default="{ row }">
                  <div class="trans-cell">
                    <span v-for="l in LANG_TABS" :key="l" class="trans-item">
                      <span class="trans-lang">{{ l }}</span>
                      <span v-if="row.translations?.[l]">{{ row.translations[l] }}</span>
                      <span v-else class="trans-missing">缺译</span>
                    </span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="sort_order" label="排序" width="70" />
              <el-table-column label="默认" width="70">
                <template #default="{ row }">
                  <el-tag v-if="row.is_default" size="small" type="success" effect="plain">默认</el-tag>
                  <span v-else class="dim">—</span>
                </template>
              </el-table-column>
              <el-table-column label="启用" width="70">
                <template #default="{ row }">
                  <el-switch :model-value="row.is_active" @change="(v: any) => toggleItem(row, v)" />
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150" fixed="right">
                <template #default="{ row }">
                  <el-button size="small" @click="openItemDialog(row)">编辑</el-button>
                  <el-button
                    v-if="!currentType?.is_system"
                    size="small"
                    type="danger"
                    @click="handleDeleteItem(row)"
                  >删除</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div v-if="currentType?.is_system" class="tip">系统内置类型的存储值不可改、条目不可删除，仅可调整文案、翻译、颜色、排序与启停。</div>
          </template>
        </div>
      </el-card>
    </div>

    <!-- 枚举项编辑弹窗 -->
    <el-dialog v-model="itemDialogVisible" :title="editingItemId ? '编辑枚举项' : '新建枚举项'" width="560px">
      <el-form :model="itemForm" label-width="90px">
        <el-form-item label="存储值" required>
          <el-input v-model="itemForm.value" :disabled="!!editingItemId" placeholder="如 oem / logo（保存后不可修改）" />
        </el-form-item>
        <el-form-item label="English" required>
          <el-input v-model="itemForm.label" placeholder="默认语言文案" />
        </el-form-item>
        <el-form-item label="中文">
          <el-input v-model="itemForm.zh" placeholder="中文翻译" />
        </el-form-item>
        <el-form-item label="Español">
          <el-input v-model="itemForm.es" placeholder="西班牙语翻译" />
        </el-form-item>
        <el-form-item label="Français">
          <el-input v-model="itemForm.fr" placeholder="法语翻译" />
        </el-form-item>
        <el-form-item label="标签颜色">
          <el-select v-model="itemForm.color" placeholder="选择颜色" clearable style="width: 100%">
            <el-option v-for="c in COLOR_OPTIONS" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <div class="form-row">
          <el-form-item label="排序">
            <el-input-number v-model="itemForm.sort_order" :min="0" controls-position="right" style="width: 140px" />
          </el-form-item>
          <el-form-item label="设为默认">
            <el-switch v-model="itemForm.is_default" />
          </el-form-item>
          <el-form-item label="启用">
            <el-switch v-model="itemForm.is_active" />
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="itemDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveItem">保存</el-button>
      </template>
    </el-dialog>

    <!-- 枚举类型编辑弹窗 -->
    <el-dialog v-model="typeDialogVisible" :title="editingTypeId ? '编辑枚举类型' : '新建枚举类型'" width="560px">
      <el-form :model="typeForm" label-width="100px">
        <el-form-item label="编码" required>
          <el-input v-model="typeForm.code" :disabled="editingTypeId ? isEditingSystemType : false" placeholder="如 product.type / order.status" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="typeForm.name" placeholder="如 产品类型" />
        </el-form-item>
        <el-form-item label="模块">
          <el-input v-model="typeForm.module" placeholder="如 product / order / system" />
        </el-form-item>
        <el-form-item label="门户前缀">
          <el-input v-model="typeForm.i18n_prefix" placeholder="门户 localizedEnum 的 key 前缀，如 product.type_options" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="typeForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <div class="form-row">
          <el-form-item label="排序">
            <el-input-number v-model="typeForm.sort_order" :min="0" controls-position="right" style="width: 140px" />
          </el-form-item>
          <el-form-item label="启用">
            <el-switch v-model="typeForm.is_active" />
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="typeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveType">保存</el-button>
      </template>
    </el-dialog>
  </div>
  <input ref="fileInput" type="file" accept=".csv,.txt" class="hidden-input" @change="onImportFile" />
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { enumApi } from '@/api'
import type { SysEnumType, SysEnumItem } from '@/types'

const LANG_TABS = ['zh', 'es', 'fr'] as const

const COLOR_OPTIONS = [
  { value: '', label: '默认' },
  { value: 'blue', label: '蓝' },
  { value: 'green', label: '绿' },
  { value: 'purple', label: '紫' },
  { value: 'orange', label: '橙' },
  { value: 'magenta', label: '品红' },
  { value: 'gold', label: '金' },
  { value: 'primary', label: '主题蓝' },
  { value: 'success', label: '成功绿' },
  { value: 'info', label: '灰' },
  { value: 'warning', label: '警告橙' },
  { value: 'danger', label: '危险红' },
]

type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'

const COLOR_TYPE_MAP: Record<string, TagType> = {
  blue: 'primary', green: 'success', purple: 'warning',
  orange: 'warning', magenta: 'danger', gold: 'warning',
  primary: 'primary', success: 'success', info: 'info',
  warning: 'warning', danger: 'danger',
}

function tagType(color?: string): TagType {
  if (!color) return 'info'
  return COLOR_TYPE_MAP[color] || 'info'
}

// ===== 类型 =====
const types = ref<SysEnumType[]>([])
const currentType = ref<SysEnumType | null>(null)
const typeLoading = ref(false)

const typeDialogVisible = ref(false)
const editingTypeId = ref('')
const typeForm = reactive({ code: '', name: '', module: '', i18n_prefix: '', description: '', sort_order: 0, is_active: true })
const isEditingSystemType = ref(false)

async function loadTypes() {
  typeLoading.value = true
  try {
    types.value = await enumApi.types()
    if (!currentType.value && types.value.length) selectType(types.value[0])
  } catch { /* handled */ } finally { typeLoading.value = false }
}

function selectType(t: SysEnumType) {
  currentType.value = t
  loadItems(t.id)
}

function openTypeDialog(t?: SysEnumType) {
  editingTypeId.value = t?.id || ''
  isEditingSystemType.value = !!t?.is_system
  Object.assign(typeForm, t
    ? { code: t.code, name: t.name, module: t.module, i18n_prefix: t.i18n_prefix, description: t.description, sort_order: t.sort_order, is_active: t.is_active }
    : { code: '', name: '', module: '', i18n_prefix: '', description: '', sort_order: 0, is_active: true })
  typeDialogVisible.value = true
}

async function saveType() {
  if (!typeForm.code.trim() || !typeForm.name.trim()) { ElMessage.warning('请填写编码和名称'); return }
  try {
    await enumApi.createType({ ...typeForm, code: typeForm.code.trim(), name: typeForm.name.trim() })
    ElMessage.success('类型已保存')
    typeDialogVisible.value = false
    await loadTypes()
  } catch { /* handled */ }
}

async function handleDeleteType() {
  if (!currentType.value) return
  try {
    await ElMessageBox.confirm(`确定删除类型「${currentType.value.name}」及其全部枚举项？`, '提示', { type: 'warning' })
  } catch { return }
  try {
    await enumApi.deleteType(currentType.value.id)
    ElMessage.success('已删除')
    currentType.value = null
    items.value = []
    await loadTypes()
  } catch { /* handled */ }
}

// ===== 条目 =====
const items = ref<SysEnumItem[]>([])
const onlyMissing = ref(false)
function isMissingItem(it: SysEnumItem) {
  return LANG_TABS.some((l) => !(it.translations?.[l]))
}
const filteredItems = computed(() => (onlyMissing.value ? items.value.filter(isMissingItem) : items.value))
const itemLoading = ref(false)

const itemDialogVisible = ref(false)
const editingItemId = ref('')
const itemForm = reactive({ value: '', label: '', zh: '', es: '', fr: '', color: '', sort_order: 0, is_default: false, is_active: true })

async function loadItems(typeId: string) {
  itemLoading.value = true
  try { items.value = await enumApi.items(typeId) } catch { items.value = [] } finally { itemLoading.value = false }
}

function openItemDialog(row?: any) {
  editingItemId.value = row?.id || ''
  Object.assign(itemForm, row
    ? { value: row.value, label: row.label, zh: row.translations?.zh || '', es: row.translations?.es || '', fr: row.translations?.fr || '', color: row.color, sort_order: row.sort_order, is_default: row.is_default, is_active: row.is_active }
    : { value: '', label: '', zh: '', es: '', fr: '', color: '', sort_order: 0, is_default: false, is_active: true })
  itemDialogVisible.value = true
}

async function saveItem() {
  if (!currentType.value) return
  if (!itemForm.value.trim() || !itemForm.label.trim()) { ElMessage.warning('请填写存储值和 English 文案'); return }
  const translations: Record<string, string> = {}
  for (const l of LANG_TABS) {
    const v = itemForm[l].trim()
    if (v) translations[l] = v
  }
  try {
    await enumApi.upsertItem({
      type_id: currentType.value.id,
      value: itemForm.value.trim(),
      label: itemForm.label.trim(),
      translations,
      color: itemForm.color,
      sort_order: itemForm.sort_order,
      is_default: itemForm.is_default,
      is_active: itemForm.is_active,
    })
    ElMessage.success('已保存')
    itemDialogVisible.value = false
    await loadItems(currentType.value.id)
    await loadTypes()
  } catch { /* handled */ }
}

async function toggleItem(row: any, v: boolean) {
  try {
    await enumApi.upsertItem({
      type_id: row.type_id, value: row.value, label: row.label, translations: row.translations,
      color: row.color, sort_order: row.sort_order, is_default: row.is_default, is_active: v,
    })
    row.is_active = v
    ElMessage.success(v ? '已启用' : '已停用')
  } catch { /* handled */ }
}

async function handleDeleteItem(row: any) {
  try {
    await ElMessageBox.confirm(`确定删除枚举项「${row.value}」？`, '提示', { type: 'warning' })
  } catch { return }
  try {
    await enumApi.deleteItem(row.id)
    ElMessage.success('已删除')
    if (currentType.value) await loadItems(currentType.value.id)
    await loadTypes()
  } catch { /* handled */ }
}

// ===== 批量导入 / 导出 =====
const fileInput = ref<HTMLInputElement>()

function csvEscape(v: string): string {
  const s = v ?? ''
  if (/[",\n\r]/.test(s)) return '"' + s.replace(/"/g, '""') + '"'
  return s
}

function exportCsv() {
  if (!currentType.value) return
  const header = ['value', 'label_en', 'zh', 'es', 'fr', 'color', 'sort_order', 'is_active', 'is_default']
  const rows = items.value.map((it) => [
    it.value, it.label,
    it.translations?.zh || '', it.translations?.es || '', it.translations?.fr || '',
    it.color || '', String(it.sort_order ?? 0), it.is_active ? '1' : '0', it.is_default ? '1' : '0',
  ])
  const csv = [header, ...rows].map((r) => r.map(csvEscape).join(',')).join('\r\n')
  const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${currentType.value.code.replace(/\./g, '_')}.csv`
  a.click()
  URL.revokeObjectURL(url)
  ElMessage.success('已导出 CSV')
}

function triggerImport() {
  fileInput.value?.click()
}

// 简易 CSV 解析（支持引号包裹与双引号转义）
function parseCsvLine(line: string): string[] {
  const out: string[] = []
  let cur = ''
  let inQuotes = false
  for (let i = 0; i < line.length; i++) {
    const ch = line[i]
    if (inQuotes) {
      if (ch === '"') {
        if (line[i + 1] === '"') { cur += '"'; i++ } else { inQuotes = false }
      } else { cur += ch }
    } else if (ch === '"') {
      inQuotes = true
    } else if (ch === ',') {
      out.push(cur); cur = ''
    } else {
      cur += ch
    }
  }
  out.push(cur)
  return out
}

async function onImportFile(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file || !currentType.value) return
  const text = await file.text()
  const lines = text.replace(/^\uFEFF/, '').split(/\r?\n/).map((l) => l.trim()).filter(Boolean)
  if (lines.length < 2) { ElMessage.warning('CSV 内容为空'); return }
  const header = parseCsvLine(lines[0])
  const vi = header.indexOf('value')
  const li = header.indexOf('label_en')
  if (vi < 0 || li < 0) { ElMessage.warning('CSV 缺少 value / label_en 列'); return }
  const zi = header.indexOf('zh'); const ei = header.indexOf('es'); const fi = header.indexOf('fr')
  const ci = header.indexOf('color'); const si = header.indexOf('sort_order')
  const ai = header.indexOf('is_active'); const di = header.indexOf('is_default')

  let ok = 0; let skip = 0
  for (let n = 1; n < lines.length; n++) {
    const cells = parseCsvLine(lines[n])
    const value = (cells[vi] || '').trim()
    const label = (cells[li] || '').trim()
    if (!value || !label) { skip++; continue }
    const translations: Record<string, string> = {}
    if (zi >= 0 && cells[zi]?.trim()) translations.zh = cells[zi].trim()
    if (ei >= 0 && cells[ei]?.trim()) translations.es = cells[ei].trim()
    if (fi >= 0 && cells[fi]?.trim()) translations.fr = cells[fi].trim()
    try {
      await enumApi.upsertItem({
        type_id: currentType.value.id,
        value, label, translations,
        color: (ci >= 0 ? cells[ci] : '')?.trim() || '',
        sort_order: si >= 0 ? Number(cells[si]) || 0 : 0,
        is_default: di >= 0 ? (cells[di]?.trim() === '1') : false,
        is_active: ai >= 0 ? (cells[ai]?.trim() !== '0') : true,
      })
      ok++
    } catch { skip++ }
  }
  ElMessage.success(`导入完成：成功 ${ok} 条，跳过 ${skip} 条`)
  await loadItems(currentType.value.id)
  await loadTypes()
}

onMounted(loadTypes)
</script>

<style scoped>
.layout { display: grid; grid-template-columns: 320px 1fr; gap: 16px; align-items: start; }
@media (max-width: 900px) { .layout { grid-template-columns: 1fr; } }
.panel-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.panel-actions { display: flex; gap: 8px; align-items: center; }
.code-hint { margin-left: 10px; font-size: 12px; color: #909399; }
.prefix-hint { margin-left: 10px; font-size: 12px; color: #b3a06b; }
.type-list { max-height: calc(100vh - 220px); overflow: auto; }
.type-item { padding: 10px 12px; border-radius: 8px; margin-bottom: 8px; cursor: pointer; border: 1px solid #ebeef5; transition: all 0.2s; }
.type-item:hover { border-color: #d4a853; }
.type-item.active { background: #0d1b2a; border-color: #0d1b2a; }
.type-item.active .type-name, .type-item.active .type-code { color: #fff; }
.type-item.active .type-meta { color: #b7c0c9; }
.type-main { display: flex; align-items: center; gap: 8px; }
.type-name { font-weight: 600; color: #303133; }
.type-code { font-size: 12px; color: #909399; margin-top: 2px; font-family: monospace; }
.type-meta { font-size: 12px; color: #c0c4cc; margin-top: 2px; }
.value-code { font-family: monospace; background: #f5f7fa; padding: 2px 6px; border-radius: 4px; font-size: 12px; }
.trans-cell { display: flex; flex-direction: column; gap: 2px; }
.trans-item { font-size: 12px; color: #606266; }
.trans-lang { display: inline-block; width: 24px; color: #909399; }
.trans-missing { color: #f56c6c; font-size: 12px; }
.filter-row { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.filter-count { font-size: 12px; color: #f56c6c; }
.dim { color: #c0c4cc; }
.tip { margin-top: 12px; font-size: 12px; color: #b3a06b; }
.form-row { display: flex; gap: 24px; }
.font-semibold { font-weight: 600; }
.hidden-input { display: none; }
</style>



