<template>
  <div class="media-picker">
    <div class="mp-row">
      <el-input :model-value="modelValue" placeholder="图片 URL 或从媒体库选择" clearable @update:model-value="onInput" />
      <el-button @click="openPicker">媒体库</el-button>
      <el-upload :http-request="customUpload" :show-file-list="false" :accept="uploadAccept">
        <el-button type="primary">上传</el-button>
      </el-upload>
    </div>

    <div v-if="modelValue" class="mp-preview">
      <el-image v-if="kind === 'image'" :src="modelValue" fit="cover" class="mp-preview-img" :preview-src-list="[modelValue]" preview-teleported />
      <div v-else class="mp-file-name" :title="modelValue">{{ modelValue }}</div>
      <el-button size="small" text type="danger" @click="emit('update:modelValue', '')">清空</el-button>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="780px" append-to-body class="media-picker-dialog">
      <div class="mp-toolbar">
        <el-input v-model="mediaKeyword" placeholder="搜索文件名/标题" style="width: 220px" clearable @keyup.enter="searchMedia" />
        <el-button @click="searchMedia">查询</el-button>
        <div class="mp-spacer"></div>
        <el-upload :http-request="customUpload" :show-file-list="false" :accept="uploadAccept">
          <el-button size="small" type="primary" plain>上传新图片</el-button>
        </el-upload>
      </div>
      <div class="mp-grid" v-loading="mediaLoading">
        <div v-for="m in mediaItems" :key="m.id" class="mp-cell" @click="pickMedia(m)">
          <el-image v-if="m.type === 'image'" :src="m.url" fit="cover" class="mp-img" />
          <div v-else class="mp-file">{{ m.type }}</div>
          <div class="mp-name" :title="m.original_name || m.title">{{ m.original_name || m.title }}</div>
        </div>
        <div v-if="!mediaLoading && !mediaItems.length" class="mp-empty">暂无图片，请先在「媒体管理」上传</div>
      </div>
      <el-pagination
        v-if="mediaTotal > mediaPageSize"
        class="mp-pagination"
        v-model:current-page="mediaPage"
        :page-size="mediaPageSize"
        :total="mediaTotal"
        layout="total, prev, pager, next"
        @current-change="loadMedia"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { mediaApi } from '@/api'

const props = defineProps<{ modelValue: string; mediaType?: string; accept?: string }>()
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()

const kind = props.mediaType || 'image'
const isVideo = kind === 'video'
const isFile = kind === 'file'
const dialogTitle = isVideo ? '从媒体库选择视频' : isFile ? '从媒体库选择文件' : '从媒体库选择图片'
const uploadAccept = props.accept || (isVideo || isFile
  ? '.jpg,.jpeg,.png,.gif,.webp,.avif,.svg,.mp4,.webm,.pdf,.doc,.docx,.xls,.xlsx,.zip'
  : '.jpg,.jpeg,.png,.gif,.webp,.avif,.svg')

function onInput(v: string) { emit('update:modelValue', v || '') }

const dialogVisible = ref(false)
const mediaKeyword = ref('')
const mediaItems = ref<any[]>([])
const mediaTotal = ref(0)
const mediaPage = ref(1)
const mediaPageSize = ref(12)
const mediaLoading = ref(false)

async function loadMedia() {
  mediaLoading.value = true
  try {
    const res: any = await mediaApi.list({ page: mediaPage.value, pageSize: mediaPageSize.value, type: props.mediaType || 'image', keyword: mediaKeyword.value || undefined })
    mediaItems.value = res.items || res || []
    mediaTotal.value = res.total ?? mediaItems.value.length
  } finally {
    mediaLoading.value = false
  }
}

function searchMedia() { mediaPage.value = 1; loadMedia() }
function openPicker() { mediaPage.value = 1; dialogVisible.value = true; loadMedia() }
function pickMedia(m: any) { emit('update:modelValue', m.url); dialogVisible.value = false; ElMessage.success('已选用媒体') }

async function customUpload(options: any) {
  const fd = new FormData()
  fd.append('file', options.file)
  try {
    const res: any = await mediaApi.upload(fd)
    emit('update:modelValue', res.url || '')
    ElMessage.success('上传成功')
    options.onSuccess?.(res)
    if (dialogVisible.value) loadMedia()
  } catch {
    ElMessage.error('上传失败')
    options.onError?.(new Error('upload failed'))
  }
}
</script>

<style scoped>
.mp-row { display: flex; gap: 8px; }
.mp-row .el-input { flex: 1; }
.mp-preview { display: flex; align-items: center; gap: 10px; margin-top: 10px; }
.mp-preview-img { width: 120px; height: 80px; border-radius: 6px; border: 1px solid #eae5dd; }
.mp-file-name { flex: 1; min-width: 0; font-size: 12px; color: #909399; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.mp-toolbar { display: flex; gap: 10px; margin-bottom: 12px; align-items: center; }
.mp-spacer { flex: 1; }
.mp-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 12px; min-height: 80px; }
.mp-cell { border: 1px solid #eae5dd; border-radius: 8px; padding: 8px; display: flex; flex-direction: column; align-items: center; gap: 6px; cursor: pointer; transition: border-color 0.2s; }
.mp-cell:hover { border-color: #d4a853; }
.mp-img { width: 100%; height: 90px; border-radius: 6px; }
.mp-file { width: 100%; height: 90px; display: flex; align-items: center; justify-content: center; background: #f5f5f5; border-radius: 6px; color: #999; font-size: 12px; }
.mp-name { font-size: 11px; color: #666; width: 100%; text-align: center; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.mp-empty { grid-column: 1 / -1; text-align: center; color: #999; padding: 30px 0; }
.mp-pagination { margin-top: 10px; justify-content: flex-end; }
</style>