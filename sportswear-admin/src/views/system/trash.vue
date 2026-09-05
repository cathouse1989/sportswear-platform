<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-select v-model="entityType" placeholder="全部类型" clearable filterable style="width: 220px" @change="handleSearch">
        <el-option v-for="t in ENTITY_TYPES" :key="t.value" :label="`${t.label}（${t.value}）`" :value="t.value" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="danger" @click="handleEmptyTrash">清空回收站</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column prop="entity_type" label="类型" width="150">
        <template #default="{ row }">{{ entityLabel(row.entity_type) }}</template>
      </el-table-column>
      <el-table-column prop="title" label="标题" min-width="220" />
      <el-table-column prop="deleted_at" label="删除时间" width="170" />
      <el-table-column label="操作" width="180">
        <template #default="{ row }">
          <el-button size="small" type="success" @click="handleRestore(row)">恢复</el-button>
          <el-button size="small" type="danger" @click="handlePurge(row)">彻底删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="loadData" />
  </el-card>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { trashApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'

// 与后端 trash_service.go 的 trashEntities 注册表保持一致（缺陷 B-02 修复：此前硬编码仅 5 种）
const ENTITY_TYPES = [
  { label: '产品', value: 'product' },
  { label: '产品分类', value: 'category' },
  { label: '产品系列', value: 'series' },
  { label: '面料', value: 'fabric' },
  { label: '页面', value: 'page' },
  { label: '博客', value: 'blog' },
  { label: '案例', value: 'case' },
  { label: 'FAQ', value: 'faq' },
  { label: '工厂', value: 'factory' },
  { label: '认证', value: 'certification' },
  { label: '生产流程', value: 'production_process' },
  { label: '导航', value: 'navigation' },
  { label: '媒体', value: 'media' },
  { label: '自媒体', value: 'self_media' },
  { label: '询盘', value: 'lead' },
  { label: '用户', value: 'user' },
]
function entityLabel(value: string) { return ENTITY_TYPES.find((t) => t.value === value)?.label || value }

const items = ref<any[]>([]); const loading = ref(false); const total = ref(0); const page = ref(1); const pageSize = useAdminPageSize(); const entityType = ref('')
async function loadData() { loading.value = true; try { const r = await trashApi.list({ page: page.value, pageSize: pageSize.value, entity_type: entityType.value }); items.value = r.items; total.value = r.total } finally { loading.value = false } }
function handleSearch() { page.value = 1; loadData() }
// 恢复/彻底删除必须携带 entity_type（后端 binding:"required"，缺失会 400，缺陷 B-11 修复）
async function handleRestore(row: any) { await trashApi.restore(row.id, row.entity_type); ElMessage.success('已恢复'); loadData() }
async function handlePurge(row: any) { await ElMessageBox.confirm('确定彻底删除吗？不可恢复', '警告', { type: 'warning' }); await trashApi.purge(row.id, row.entity_type); ElMessage.success('已删除'); loadData() }
async function handleEmptyTrash() { await ElMessageBox.confirm('确定清空回收站吗？', '警告', { type: 'warning' }); await trashApi.empty(); ElMessage.success('已清空'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>