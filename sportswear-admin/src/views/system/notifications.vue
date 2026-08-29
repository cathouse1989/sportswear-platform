<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-switch v-model="onlyUnread" active-text="仅未读" @change="loadData" />
      <el-button type="primary" @click="handleMarkAllRead" :disabled="!hasUnread">全部标记已读</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe>
      <el-table-column label="状态" width="70">
        <template #default="{ row }"><el-badge v-if="!row.is_read" is-dot /></template>
      </el-table-column>
      <el-table-column prop="title" label="标题" width="160" />
      <el-table-column prop="content" label="内容" min-width="260" />
      <el-table-column prop="type" label="类型" width="100" />
      <el-table-column prop="created_at" label="时间" width="170" />
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button v-if="!row.is_read" size="small" @click="handleMarkRead(row)">标为已读</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="loadData" />
  </el-card>
</template>
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { notificationApi } from '@/api'

const items = ref<any[]>([]); const loading = ref(false); const total = ref(0); const page = ref(1); const pageSize = ref(20); const onlyUnread = ref(false)
const hasUnread = computed(() => items.value.some(i => !i.is_read))
async function loadData() { loading.value = true; try { const r = await notificationApi.list({ page: page.value, pageSize: pageSize.value, only_unread: onlyUnread.value }); items.value = r.items; total.value = r.total } finally { loading.value = false } }
async function handleMarkRead(row: any) { await notificationApi.markRead(row.id); ElMessage.success('已标记'); loadData() }
async function handleMarkAllRead() { await notificationApi.markAllRead(); ElMessage.success('全部已读'); loadData() }
onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>