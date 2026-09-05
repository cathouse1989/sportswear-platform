<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input v-model="keyword" placeholder="搜索流程名称" clearable style="width: 240px" @keyup.enter="handleSearch" />
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建流程</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column label="图片" width="76" align="center">
        <template #default="{ row }">
          <el-image v-if="row.image" :src="row.image" fit="cover" class="cover-thumb" :preview-src-list="[row.image]" preview-teleported />
          <div v-else class="cover-placeholder">—</div>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="180" />
      <el-table-column prop="sort_order" label="排序" width="80" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }"><el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag></template>
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
    <ProFormDialog v-model="dialogVisible" :title="editingId ? '编辑流程' : '新建流程'" :form="form" :rules="rules" @submit="handleSave">
      <el-form-item label="名称" prop="name"><el-input v-model="form.name" /></el-form-item>
      <el-form-item label="图片"><MediaPicker v-model="form.image" /></el-form-item>
      <el-form-item label="视频"><el-input v-model="form.video" placeholder="视频 URL（可选）" /></el-form-item>
      <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="form.description" type="textarea" /></el-form-item>
    </ProFormDialog>
  </el-card>
</template>
<script setup lang="ts">
import { onMounted } from 'vue'
import type { FormRules } from 'element-plus'
import { productionProcessApi } from '@/api'
import MediaPicker from '@/components/media/MediaPicker.vue'
import ProFormDialog from '@/components/pro/ProFormDialog.vue'
import { useCrud } from '@/composables/useCrud'

const emptyForm = () => ({ name: '', description: '', image: '', video: '', sort_order: 0 })
const {
  items, total, loading, page, pageSize, keyword,
  loadData, handleSearch,
  dialogVisible, editingId, form,
  openCreateDialog, openEditDialog, handleSave,
  handlePublish, handleUnpublish, handleDelete,
} = useCrud({ api: productionProcessApi, emptyForm, serverPagination: true })
const rules: FormRules = { name: [{ required: true, message: '请输入名称', trigger: 'blur' }] }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.cover-thumb { width: 44px; height: 44px; border-radius: 4px; }
.cover-placeholder { width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; color: #c0c4cc; background: #f5f7fa; border-radius: 4px; }
</style>