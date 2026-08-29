<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input
        v-model="module"
        placeholder="模块"
        clearable
        style="width: 160px"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="operation" placeholder="操作类型" clearable style="width: 160px" @change="handleSearch">
        <el-option label="创建" value="create" />
        <el-option label="更新" value="update" />
        <el-option label="删除" value="delete" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
    </div>

    <el-table :data="logs" v-loading="loading" stripe>
      <el-table-column prop="module" label="模块" width="120" />
      <el-table-column prop="operation" label="操作" width="100">
        <template #default="{ row }">
          <el-tag size="small">{{ operationLabel(row.operation) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="entity_type" label="实体类型" width="140" />
      <el-table-column prop="description" label="描述" min-width="220" />
      <el-table-column prop="ip" label="IP" width="130" />
      <el-table-column prop="created_at" label="时间" width="180">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      :page-sizes="[10, 20, 50, 100]"
      @current-change="loadData"
      @size-change="loadData"
    />
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { operationLogApi } from '@/api'
import type { OperationLog } from '@/types'

const logs = ref<OperationLog[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const module = ref('')
const operation = ref('')

async function loadData() {
  loading.value = true
  try {
    const result = await operationLogApi.list({
      page: page.value,
      pageSize: pageSize.value,
      module: module.value,
      operation: operation.value,
    })
    logs.value = result.items
    total.value = result.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

function operationLabel(op: string) {
  const map: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    publish: '发布',
    unpublish: '下线',
    upload: '上传',
    login: '登录',
    logout: '登出',
  }
  return map[op] || op
}

function formatTime(t: string) {
  if (!t) return '-'
  return new Date(t).toLocaleString()
}

onMounted(loadData)
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
}
.pagination {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>