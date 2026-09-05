<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索名称 / 平台 / 账号" clearable style="width: 260px" @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建自媒体</el-button>
    </div>

    <div class="max-display-bar">
      <span class="max-display-label">门户最多展示数量：</span>
      <el-input-number v-model="maxDisplay" :min="1" :max="50" size="small" />
      <el-button size="small" type="primary" :loading="savingMax" @click="saveMaxDisplay">保存</el-button>
      <span class="max-display-hint">已发布的自媒体按排序取前 N 条展示在门户「关注我们」区块</span>
    </div>

    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column label="图片" width="76" align="center">
        <template #default="{ row }">
          <el-image v-if="row.image" :src="row.image" fit="cover" class="cover-thumb" :preview-src-list="[row.image]" preview-teleported />
          <div v-else class="cover-placeholder">—</div>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column label="平台" width="120">
        <template #default="{ row }">{{ platformLabel(row.platform) }}</template>
      </el-table-column>
      <el-table-column prop="account" label="账号" min-width="140" />
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status === 'published' ? '展示' : '不展示' }}</el-tag>
        </template>
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

    <ProFormDialog v-model="dialogVisible" :title="editingId ? '编辑自媒体' : '新建自媒体'" :form="form" :rules="rules" @submit="handleSave">
      <el-form-item label="名称" prop="name"><el-input v-model="form.name" /></el-form-item>
      <el-form-item label="平台">
        <el-select v-model="form.platform" placeholder="选择平台" style="width: 100%">
          <el-option v-for="p in PLATFORMS" :key="p.value" :label="p.label" :value="p.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="账号"><el-input v-model="form.account" placeholder="账号 ID / 昵称（可选）" /></el-form-item>
      <el-form-item label="链接"><el-input v-model="form.url" placeholder="跳转链接（可选）" /></el-form-item>
      <el-form-item label="二维码/头像"><MediaPicker v-model="form.image" /></el-form-item>
      <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
      <el-form-item label="简介"><el-input v-model="form.description" type="textarea" /></el-form-item>
    </ProFormDialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { selfMediaApi } from '@/api'
import MediaPicker from '@/components/media/MediaPicker.vue'
import ProFormDialog from '@/components/pro/ProFormDialog.vue'
import { useCrud } from '@/composables/useCrud'

const PLATFORMS = [
  { value: 'wechat', label: '微信公众号' },
  { value: 'weibo', label: '微博' },
  { value: 'douyin', label: '抖音' },
  { value: 'xiaohongshu', label: '小红书' },
  { value: 'video', label: '视频号' },
  { value: 'bilibili', label: 'B站' },
  { value: 'youtube', label: 'YouTube' },
  { value: 'instagram', label: 'Instagram' },
  { value: 'facebook', label: 'Facebook' },
  { value: 'twitter', label: 'Twitter/X' },
  { value: 'linkedin', label: 'LinkedIn' },
  { value: 'tiktok', label: 'TikTok' },
  { value: 'other', label: '其他' },
]
const platformLabel = (v: string) => PLATFORMS.find((p) => p.value === v)?.label || v || '—'

const emptyForm = () => ({ name: '', platform: '', account: '', url: '', image: '', description: '', sort_order: 0 })
const {
  items, total, loading, page, pageSize, keyword,
  loadData, handleSearch,
  dialogVisible, editingId, form,
  openCreateDialog, openEditDialog, handleSave,
  handlePublish, handleUnpublish, handleDelete,
} = useCrud({ api: selfMediaApi, emptyForm, serverPagination: true })
const rules: FormRules = { name: [{ required: true, message: '请输入名称', trigger: 'blur' }] }

// 最多展示数量（theme_configs.self_media_max_display）
const maxDisplay = ref(6)
const savingMax = ref(false)
async function loadMaxDisplay() {
  try {
    const res = await selfMediaApi.getConfig()
    const n = Number(res?.max_display)
    if (Number.isFinite(n) && n > 0) maxDisplay.value = n
  } catch { /* 读取失败用默认值 */ }
}
async function saveMaxDisplay() {
  savingMax.value = true
  try {
    await selfMediaApi.updateConfig(maxDisplay.value)
    ElMessage.success('已保存最多展示数量')
  } catch {
    ElMessage.error('保存失败')
  } finally {
    savingMax.value = false
  }
}

onMounted(() => { loadData(); loadMaxDisplay() })
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 12px; align-items: center; }
.spacer { flex: 1; }
.max-display-bar { display: flex; gap: 8px; align-items: center; padding: 8px 12px; background: #f5f7fa; border-radius: 6px; margin-bottom: 16px; }
.max-display-label { font-size: 13px; color: #374151; }
.max-display-hint { font-size: 12px; color: #9ca3af; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.cover-thumb { width: 44px; height: 44px; border-radius: 4px; }
.cover-placeholder { width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; color: #c0c4cc; background: #f5f7fa; border-radius: 4px; }
</style>