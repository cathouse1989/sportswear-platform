<template>
  <div>
    <!-- 统计卡 -->
    <div class="stats-row">
      <el-card shadow="never" class="stat-card">
        <div class="stat-value">{{ stats.total }}</div>
        <div class="stat-label">订阅总数</div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-value text-emerald">{{ stats.subscribed }}</div>
        <div class="stat-label">已订阅</div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="stat-value text-gray">{{ stats.unsubscribed }}</div>
        <div class="stat-label">已退订</div>
      </el-card>
    </div>

    <el-card shadow="never">
      <div class="toolbar">
        <el-input
          v-model="keyword"
          placeholder="搜索邮箱 / 来源 / 落地页"
          clearable
          style="width: 260px"
          @keyup.enter="handleSearch"
        />
        <el-select v-model="status" placeholder="状态" clearable style="width: 160px" @change="handleSearch">
          <el-option v-for="o in enumOptions('subscriber.status')" :key="o.value" :label="o.label" :value="o.value" />
        </el-select>
        <el-button type="primary" @click="handleSearch">查询</el-button>
      </div>

      <el-table :data="subscribers" v-loading="loading" stripe>
        <el-table-column prop="email" label="邮箱" min-width="220" />
        <el-table-column prop="language" label="语言" width="80" />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="enumTag('subscriber.status', row.status)" size="small">
              {{ enumLabel('subscriber.status', row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="source" label="来源" width="120">
          <template #default="{ row }">{{ row.source || '-' }}</template>
        </el-table-column>
        <el-table-column prop="landing_page" label="落地页" min-width="160" show-overflow-tooltip>
          <template #default="{ row }">{{ row.landing_page || '-' }}</template>
        </el-table-column>
        <el-table-column label="订阅时间" width="170">
          <template #default="{ row }">{{ formatDateTime(row.subscribed_at || row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }: any">
            <el-button size="small" @click="toggleStatus(row)">
              {{ row.status === 'subscribed' ? '退订' : '恢复订阅' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
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
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { subscriberApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import { useEnumDict } from '@/composables/useEnumDict'
import { formatDateTime } from '@/utils/format'
import type { Subscriber } from '@/types'

// 订阅状态由「字典管理」驱动（subscriber.status）
const { ensureLoaded: loadEnumDict, options: enumOptions, label: enumLabel, tagType: enumTag } = useEnumDict()

const subscribers = ref<Subscriber[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const status = ref('')

const stats = ref({ total: 0, subscribed: 0, unsubscribed: 0 })

async function loadStats() {
  try {
    stats.value = await subscriberApi.stats()
  } catch {
    /* 忽略统计失败，不阻塞列表 */
  }
}

async function loadData() {
  loading.value = true
  try {
    const result = await subscriberApi.list({
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value,
      status: status.value,
    })
    subscribers.value = result.items
    total.value = result.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

async function toggleStatus(row: Subscriber) {
  const next = row.status === 'subscribed' ? 'unsubscribed' : 'subscribed'
  const action = next === 'subscribed' ? '恢复订阅' : '退订'
  await ElMessageBox.confirm(`确定对 ${row.email} 执行「${action}」吗？`, '提示', { type: 'warning' })
  await subscriberApi.update(row.id, { status: next })
  ElMessage.success(`${action}成功`)
  loadData()
  loadStats()
}

async function handleDelete(row: Subscriber) {
  await ElMessageBox.confirm(`确定删除订阅 ${row.email} 吗？`, '警告', { type: 'warning' })
  await subscriberApi.delete(row.id)
  ElMessage.success('已删除')
  loadData()
  loadStats()
}

onMounted(() => {
  loadEnumDict()
  loadData()
  loadStats()
})
</script>

<style scoped>
.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 16px;
}
.stat-card {
  text-align: center;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
}
.stat-value.text-emerald {
  color: #10b981;
}
.stat-value.text-gray {
  color: #9ca3af;
}
.stat-label {
  margin-top: 4px;
  font-size: 13px;
  color: #6b7280;
}
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
