<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-select v-model="mediaType" placeholder="类型" clearable style="width: 140px" @change="loadData">
        <el-option label="图片" value="image" />
        <el-option label="视频" value="video" />
        <el-option label="文件" value="file" />
      </el-select>
      <el-select v-model="category" placeholder="分类" clearable style="width: 140px" @change="loadData">
        <el-option v-for="c in CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
      </el-select>
      <el-input v-model="keyword" placeholder="搜索文件名/标题" clearable style="width: 220px" @keyup.enter="loadData" />
      <el-button type="primary" @click="loadData">查询</el-button>
      <div class="spacer" />
      <el-upload
        :action="uploadUrl"
        :headers="uploadHeaders"
        :on-success="onUploadSuccess"
        :on-error="onUploadError"
        :show-file-list="false"
        :accept="acceptTypes"
      >
        <el-button type="primary">上传文件</el-button>
      </el-upload>
    </div>

    <el-table :data="media" v-loading="loading" stripe>
      <el-table-column label="预览" width="80">
        <template #default="{ row }">
          <el-image v-if="row.type === 'image'" :src="row.url" style="width: 48px; height: 48px" fit="cover" :preview-src-list="[row.url]" preview-teleported />
          <el-tag v-else size="small">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="文件信息" min-width="220">
        <template #default="{ row }">
          <div class="file-cell">
            <div class="file-name" :title="row.original_name">{{ row.original_name }}</div>
            <div class="file-meta">{{ row.file_type }} · {{ fmtSize(null, null, row.file_size) }}</div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="访问 URL / 存储路径" min-width="260">
        <template #default="{ row }">
          <div class="url-cell">
            <span class="url-text" :title="row.url">{{ row.url }}</span>
            <el-button link type="primary" size="small" @click="copyUrl(row.url)">复制</el-button>
          </div>
          <div class="url-cell">
            <span class="url-text" :title="row.path">{{ row.path }}</span>
            <el-button link type="primary" size="small" @click="copyUrl(row.path)">复制</el-button>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="分类" width="90">
        <template #default="{ row }"><el-tag size="small">{{ categoryLabel(row.category) }}</el-tag></template>
      </el-table-column>
      <el-table-column label="公开" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.is_public" type="success" size="small">公开</el-tag>
          <el-tag v-else type="info" size="small">私有</el-tag>
        </template>
      </el-table-column>
      <el-table-column :formatter="fmtTime" label="上传时间" width="160" />
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="loadData" />

    <!-- 在线编辑媒体数据 -->
    <el-dialog v-model="editDialogVisible" title="编辑媒体" width="520px" destroy-on-close>
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="预览">
          <el-image v-if="editForm.type === 'image'" :src="editForm.url" style="width: 120px; height: 90px" fit="cover" :preview-src-list="[editForm.url]" preview-teleported />
          <el-tag v-else size="small">{{ editForm.type }}</el-tag>
        </el-form-item>
        <el-form-item label="文件名">
          <el-input v-model="editForm.original_name" maxlength="200" />
        </el-form-item>
        <el-form-item label="访问 URL">
          <div class="readonly-path">
            <span class="path-text">{{ editForm.url }}</span>
            <el-button link type="primary" size="small" @click="copyUrl(editForm.url)">复制</el-button>
          </div>
        </el-form-item>
        <el-form-item label="存储路径">
          <div class="readonly-path">
            <span class="path-text">{{ editForm.path }}</span>
            <el-button link type="primary" size="small" @click="copyUrl(editForm.path)">复制</el-button>
          </div>
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="editForm.title" placeholder="图片标题（可选）" maxlength="200" />
        </el-form-item>
        <el-form-item label="替代文本">
          <el-input v-model="editForm.alt" placeholder="用于 SEO / 无障碍（可选）" maxlength="200" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editForm.description" type="textarea" :rows="3" placeholder="图片描述（可选）" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="editForm.category" style="width: 100%">
            <el-option v-for="c in CATEGORIES" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="公开访问">
          <el-switch v-model="editForm.is_public" active-text="公开" inactive-text="私有" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveEdit">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { mediaApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import { formatDateTime } from '@/utils/format'
import type { Media } from '@/types'

const CATEGORIES = [
  { label: '产品', value: 'product' },
  { label: '首页', value: 'home' },
  { label: 'OEM', value: 'oem' },
  { label: 'ODM', value: 'odm' },
  { label: '工厂', value: 'factory' },
  { label: '生产流程', value: 'production' },
  { label: '案例', value: 'case' },
  { label: '博客', value: 'blog' },
  { label: '认证', value: 'certification' },
  { label: '公开', value: 'public' },
]
function categoryLabel(value: string) { return CATEGORIES.find((c) => c.value === value)?.label || value }

const media = ref<Media[]>([])
const loading = ref(false)
const saving = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const mediaType = ref('')
const category = ref('')
const keyword = ref('')

const token = localStorage.getItem('admin_token') || ''
const uploadUrl = `/api/v1/admin/media/upload`
const uploadHeaders = { Authorization: `Bearer ${token}` }
const acceptTypes = '.jpg,.jpeg,.png,.gif,.webp,.avif,.svg,.mp4,.webm,.pdf,.doc,.docx,.xls,.xlsx,.zip'

const editDialogVisible = ref(false)
const editingId = ref('')
const editForm = reactive({ url: '', path: '', type: '', original_name: '', title: '', alt: '', description: '', category: '', is_public: true })

async function loadData() {
  loading.value = true
  try {
    const result = await mediaApi.list({ page: page.value, pageSize: pageSize.value, type: mediaType.value, category: category.value, keyword: keyword.value || undefined })
    media.value = (result as any).items || []
    total.value = (result as any).total ?? 0
  } finally { loading.value = false }
}

function fmtSize(_: any, __: any, val: number) {
  if (!val) return '-'
  if (val < 1024) return `${val}B`
  if (val < 1024 * 1024) return `${(val / 1024).toFixed(1)}KB`
  return `${(val / (1024 * 1024)).toFixed(1)}MB`
}

function fmtTime(_: any, __: any, val: string) {
  if (!val) return '-'
  return formatDateTime(val)
}

function onUploadSuccess() { ElMessage.success('上传成功'); loadData() }
function onUploadError() { ElMessage.error('上传失败') }
async function copyUrl(url: string) { try { await navigator.clipboard.writeText(url); ElMessage.success('已复制') } catch { ElMessage.warning('复制失败') } }
async function handleDelete(row: Media) { await ElMessageBox.confirm('确定删除此文件吗？', '警告', { type: 'warning' }); await mediaApi.delete(row.id); ElMessage.success('已删除'); loadData() }

function openEdit(row: Media) {
  editingId.value = row.id
  Object.assign(editForm, {
    url: row.url,
    path: row.path || '',
    type: row.type,
    original_name: row.original_name || '',
    title: row.title || '',
    alt: row.alt || '',
    description: row.description || '',
    category: row.category || 'public',
    is_public: !!row.is_public,
  })
  editDialogVisible.value = true
}

async function handleSaveEdit() {
  saving.value = true
  try {
    await mediaApi.update(editingId.value, {
      original_name: editForm.original_name,
      title: editForm.title,
      alt: editForm.alt,
      description: editForm.description,
      category: editForm.category,
      is_public: editForm.is_public,
    })
    ElMessage.success('媒体数据已更新')
    editDialogVisible.value = false
    loadData()
  } catch { /* 错误已在拦截器提示 */ }
  finally { saving.value = false }
}

onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.url-cell { display: flex; align-items: center; gap: 6px; min-width: 0; }
.url-cell + .url-cell { margin-top: 2px; }
.url-text { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 12px; color: #6b7280; direction: rtl; text-align: left; }
.file-cell { min-width: 0; }
.file-name { font-size: 13px; color: #303133; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-meta { font-size: 12px; color: #9ca3af; margin-top: 2px; }
</style>