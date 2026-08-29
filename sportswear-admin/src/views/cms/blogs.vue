<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索标题" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-select v-model="category" placeholder="分类" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="c in CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建博客</el-button>
    </div>
    <el-table :data="blogs" v-loading="loading" stripe>
      <el-table-column label="封面" width="76" align="center">
        <template #default="{ row }">
          <el-image v-if="row.cover_image" :src="row.cover_image" fit="cover" class="cover-thumb" :preview-src-list="[row.cover_image]" preview-teleported />
          <div v-else class="cover-placeholder">—</div>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column label="分类" width="120">
        <template #default="{ row }">{{ categoryLabel(row.category) }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }"><el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="220">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="success" @click="handlePublish(row)" v-if="row.status !== 'published'">发布</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="loadData" />

    <!-- 编辑弹窗（缺陷 B-01 修复：此前模板缺失导致新建/编辑点击无响应） -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑博客' : '新建博客'" width="640px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" maxlength="200" show-word-limit />
        </el-form-item>
        <el-form-item label="Slug" prop="slug">
          <el-input v-model="form.slug" placeholder="URL 标识，如 oem-vs-odm-guide" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="form.category" style="width: 100%">
            <el-option v-for="c in CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder="逗号分隔，如 yoga,leggings" />
        </el-form-item>
        <el-form-item label="封面图">
          <el-input v-model="form.cover_image" placeholder="图片 URL（媒体库直传将在媒体上传闭环任务中接入）" />
        </el-form-item>
        <el-form-item label="正文" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="8" placeholder="正文内容（富文本编辑器见打磨方案 P2-#12）" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { cmsApi } from '@/api/cms'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import type { Blog } from '@/types'
import { checkBlogGate, gateAlertMessage } from '@/utils/publish-gate'

// 分类统一定义：列表筛选与表单共用，避免两处硬编码漂移
const CATEGORIES = [
  { label: 'OEM 指南', value: 'oem_guide' },
  { label: 'ODM 指南', value: 'odm_guide' },
  { label: '面料', value: 'fabric' },
  { label: '趋势', value: 'trend' },
  { label: '行业', value: 'industry' },
]
function categoryLabel(value: string) { return CATEGORIES.find((c) => c.value === value)?.label || value }

const blogs = ref<Blog[]>([])
const loading = ref(false)
const saving = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const category = ref('')

const formRef = ref<FormInstance>()
const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({ title: '', slug: '', category: '', tags: '', cover_image: '', content: '' })
const rules: FormRules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' },
    { max: 200, message: '标题不能超过 200 字', trigger: 'blur' },
  ],
  slug: [
    { required: true, message: '请输入 Slug', trigger: 'blur' },
    { pattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/, message: '仅允许小写字母、数字和中划线', trigger: 'blur' },
  ],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  content: [{ required: true, message: '请输入正文', trigger: 'blur' }],
}

async function loadData() {
  loading.value = true
  try {
    const result = await cmsApi.blogs.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, category: category.value })
    blogs.value = result.items
    total.value = result.total
  } finally { loading.value = false }
}
function handleSearch() { page.value = 1; loadData() }
function resetForm() { Object.assign(form, { title: '', slug: '', category: '', tags: '', cover_image: '', content: '' }) }
function openCreateDialog() { editingId.value = ''; resetForm(); dialogVisible.value = true }
function openEditDialog(row: Blog) {
  editingId.value = row.id
  Object.assign(form, { title: row.title, slug: row.slug, category: row.category, tags: row.tags || '', cover_image: row.cover_image || '', content: row.content || '' })
  dialogVisible.value = true
}
async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (editingId.value) { await cmsApi.blogs.update(editingId.value, form) } else { await cmsApi.blogs.create(form) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch {} finally { saving.value = false }
}
async function handlePublish(row: Blog) {
  // 发布质检门（P0-#2）：标题/Slug/正文缺失时拦截并列出缺失项
  const gate = checkBlogGate(row)
  if (!gate.ok) {
    await ElMessageBox.alert(gateAlertMessage(gate), '无法发布', { type: 'warning', confirmButtonText: '知道了' }).catch(() => {})
    return
  }
  await cmsApi.blogs.publish(row.id)
  ElMessage.success('已发布')
  loadData()
}
async function handleDelete(row: Blog) { await ElMessageBox.confirm(`确定删除博客 ${row.title} 吗？`, '警告', { type: 'warning' }); await cmsApi.blogs.delete(row.id); ElMessage.success('已删除'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.cover-thumb { width: 44px; height: 44px; border-radius: 4px; }
.cover-placeholder { width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; color: #c0c4cc; background: #f5f7fa; border-radius: 4px; }
</style>