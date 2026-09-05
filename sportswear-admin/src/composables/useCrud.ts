import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAdminPageSize } from './useAdminPageSize'

/** 标准 CRUD API 形态（publish/unpublish 可选） */
export interface CrudApi {
  list: (params?: any) => Promise<any>
  create: (data: any) => Promise<any>
  update: (id: string, data: any) => Promise<any>
  delete: (id: string) => Promise<any>
  publish?: (id: string) => Promise<any>
  unpublish?: (id: string) => Promise<any>
}

export interface UseCrudOptions<T extends Record<string, any>> {
  api: CrudApi
  /** 表单初始值工厂：新建时重置、编辑时作为回填字段白名单 */
  emptyForm: () => T
  /** 关键字搜索字段（客户端过滤模式用） */
  searchFields?: string[]
  /** 实体名取值字段（删除确认文案） */
  nameKey?: string
  /** 服务端分页：api.list({ page, pageSize, keyword, ...extra }) -> { items, total } */
  serverPagination?: boolean
  /** 服务端分页时附加的查询参数（响应式 getter，如分类/语言） */
  extraParams?: () => Record<string, any>
  /** 保存前校验钩子：返回 false 则中止保存（如表单校验） */
  beforeSave?: () => Promise<boolean> | boolean
}

/**
 * 标准 CMS 列表 CRUD 组合式函数：
 * 统一「客户端分页 + 关键字过滤 + 弹窗表单 + 发布/下线/删除」样板，
 * 供 factories / certifications 等简单列表页复用（对齐《门户网页管理和后台多语言方案》Phase 2.3）。
 */
export function useCrud<T extends Record<string, any>>(options: UseCrudOptions<T>) {
  const { api, emptyForm, searchFields, nameKey = 'name', serverPagination = false, extraParams, beforeSave } = options

  const keyword = ref('')
  const page = ref(1)
  const pageSize = useAdminPageSize()
  const loading = ref(false)

  // 客户端模式：全量数据 + 本地过滤分页
  const allItems = ref<any[]>([])
  const filtered = computed(() => {
    const kw = keyword.value.trim().toLowerCase()
    if (!kw) return allItems.value
    const fields = searchFields && searchFields.length ? searchFields : [nameKey]
    return allItems.value.filter((item) =>
      fields.some((f) => String(item[f] ?? '').toLowerCase().includes(kw)),
    )
  })
  const clientItems = computed(() => filtered.value.slice((page.value - 1) * pageSize.value, page.value * pageSize.value))
  const clientTotal = computed(() => filtered.value.length)

  // 服务端模式：分页结果直取
  const serverItems = ref<any[]>([])
  const serverTotal = ref(0)

  const items = computed(() => (serverPagination ? serverItems.value : clientItems.value))
  const total = computed(() => (serverPagination ? serverTotal.value : clientTotal.value))

  async function loadData() {
    loading.value = true
    try {
      if (serverPagination) {
        const res = await api.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, ...(extraParams?.() || {}) })
        serverItems.value = Array.isArray(res) ? res : (res?.items || [])
        serverTotal.value = typeof res?.total === 'number' ? res.total : (Array.isArray(res) ? res.length : 0)
      } else {
        const res = await api.list()
        allItems.value = Array.isArray(res) ? res : (res?.items || res?.data || [])
      }
    } finally {
      loading.value = false
    }
  }
  function handleSearch() {
    page.value = 1
    if (serverPagination) loadData()
  }

  // 弹窗表单
  const dialogVisible = ref(false)
  const editingId = ref('')
  const saving = ref(false)
  const form = reactive(emptyForm()) as T

  function openCreateDialog() {
    editingId.value = ''
    Object.assign(form, emptyForm())
    dialogVisible.value = true
  }
  function openEditDialog(row: any) {
    editingId.value = row.id
    const blank = emptyForm()
    const next: Record<string, any> = {}
    for (const k of Object.keys(blank)) next[k] = row[k] ?? blank[k]
    Object.assign(form, next)
    dialogVisible.value = true
  }
  async function handleSave() {
    if (beforeSave) {
      const ok = await beforeSave()
      if (!ok) return
    }
    saving.value = true
    try {
      if (editingId.value) { await api.update(editingId.value, form) } else { await api.create(form) }
      ElMessage.success('保存成功')
      dialogVisible.value = false
      loadData()
    } catch { /* 失败提示由 HTTP 拦截器统一处理 */ } finally {
      saving.value = false
    }
  }
  async function handlePublish(row: any) {
    if (!api.publish) return
    await api.publish(row.id)
    ElMessage.success('已发布')
    loadData()
  }
  async function handleUnpublish(row: any) {
    if (!api.unpublish) return
    await api.unpublish(row.id)
    ElMessage.success('已下线')
    loadData()
  }
  async function handleDelete(row: any) {
    await ElMessageBox.confirm(`确定删除「${row[nameKey] ?? ''}」吗？`, '警告', { type: 'warning' })
    await api.delete(row.id)
    ElMessage.success('已删除')
    loadData()
  }

  return {
    items, total, loading, page, pageSize, keyword,
    loadData, handleSearch,
    dialogVisible, editingId, saving, form,
    openCreateDialog, openEditDialog, handleSave,
    handlePublish, handleUnpublish, handleDelete,
  }
}
