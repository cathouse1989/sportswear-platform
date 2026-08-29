<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-select v-model="mediaType" placeholder="类型" clearable style="width: 140px" @change="loadData">
        <el-option label="图片" value="image" />
        <el-option label="视频" value="video" />
        <el-option label="文件" value="file" />
      </el-select>
      <el-select v-model="category" placeholder="分类" clearable style="width: 140px" @change="loadData">
        <el-option label="产品" value="product" />
        <el-option label="首页" value="home" />
        <el-option label="工厂" value="factory" />
        <el-option label="案例" value="case" />
        <el-option label="博客" value="blog" />
        <el-option label="公开" value="public" />
      </el-select>
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
          <el-image v-if="row.type === 'image'" :src="row.url" style="width: 48px; height: 48px" fit="cover" />
          <el-tag v-else size="small">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="original_name" label="文件名" min-width="200" />
      <el-table-column prop="file_type" label="类型" width="120" />
      <el-table-column :formatter="fmtSize" label="大小" width="100" />
      <el-table-column prop="category" label="分类" width="100" />
      <el-table-column prop="created_at" label="上传时间" width="170" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }: any">
          <el-button size="small" @click="copyUrl(row.url)">复制链接</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadData" />
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { mediaApi } from '@/api'
import type { Media } from '@/types'

const media = ref<Media[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const mediaType = ref('')
const category = ref('')

const token = localStorage.getItem('admin_token') || ''
const uploadUrl = `/api/v1/admin/media/upload`
const uploadHeaders = { Authorization: `Bearer ${token}` }
const acceptTypes = '.jpg,.jpeg,.png,.gif,.webp,.avif,.svg,.mp4,.webm,.pdf,.doc,.docx,.xls,.xlsx,.zip'

async function loadData() {
  loading.value = true
  try {
    const result = await mediaApi.list({ page: page.value, pageSize: pageSize.value, type: mediaType.value, category: category.value })
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

function onUploadSuccess() { ElMessage.success('上传成功'); loadData() }
function onUploadError() { ElMessage.error('上传失败') }
async function copyUrl(url: string) { try { await navigator.clipboard.writeText(url); ElMessage.success('已复制') } catch { ElMessage.warning('复制失败') } }
async function handleDelete(row: Media) { await ElMessageBox.confirm('确定删除此文件吗？', '警告', { type: 'warning' }); await mediaApi.delete(row.id); ElMessage.success('已删除'); loadData() }

onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>