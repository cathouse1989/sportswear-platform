<template>
  <el-card shadow="never">
    <el-tabs v-model="activeTab">
      <!-- ============ 后台操作日志（跨季度分表组合查询） ============ -->
      <el-tab-pane label="后台操作日志" name="operation">
        <div class="toolbar">
          <el-input v-model="criteria.module" placeholder="模块（如 product）" clearable style="width: 160px" @keyup.enter="handleSearch" />
          <el-select v-model="criteria.operation" placeholder="操作类型" clearable style="width: 140px" @change="handleSearch">
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="发布" value="publish" />
            <el-option label="下线" value="unpublish" />
            <el-option label="上传" value="upload" />
            <el-option label="登录" value="login" />
            <el-option label="登出" value="logout" />
          </el-select>
          <el-input v-model="criteria.entityType" placeholder="实体类型（如 lead）" clearable style="width: 160px" @keyup.enter="handleSearch" />
          <el-input v-model="criteria.ip" placeholder="操作 IP" clearable style="width: 140px" @keyup.enter="handleSearch" />
          <el-input v-model="criteria.userId" placeholder="用户 ID" clearable style="width: 140px" @keyup.enter="handleSearch" />
          <el-date-picker v-model="criteria.dateRange" type="daterange" value-format="YYYY-MM-DD" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" clearable style="width: 260px" @change="handleSearch" />
          <el-button type="primary" @click="handleSearch">查询</el-button>
        </div>

        <el-table :data="logs" v-loading="loading" stripe :row-key="opRowKey">
          <el-table-column prop="id" label="ID" width="180" />
          <el-table-column prop="user_id" label="用户 ID" width="140" />
          <el-table-column prop="operation" label="操作" width="100">
            <template #default="{ row }">{{ operationLabel(row.operation) }}</template>
          </el-table-column>
          <el-table-column prop="module" label="模块" width="120" />
          <el-table-column prop="entity_type" label="实体类型" width="120" />
          <el-table-column prop="description" label="描述" min-width="200" />
          <el-table-column prop="ip" label="IP" width="130" />
          <el-table-column prop="src" label="日志表" width="170">
            <template #default="{ row }">{{ row.src || 'operation_logs' }}</template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="180">
            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
          </el-table-column>
        </el-table>

        <el-pagination class="pagination" v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="loadData" @size-change="resetPage" />
      </el-tab-pane>

      <!-- ============ 门户访问日志（PV + 广告归因明细） ============ -->
      <el-tab-pane label="门户访问日志" name="visits">
        <div class="toolbar">
          <el-input v-model="vcriteria.ip" placeholder="IP" clearable style="width: 140px" @keyup.enter="handleVisitSearch" />
          <el-input v-model="vcriteria.country" placeholder="国家" clearable style="width: 120px" @keyup.enter="handleVisitSearch" />
          <el-input v-model="vcriteria.source" placeholder="来源（如 google）" clearable style="width: 160px" @keyup.enter="handleVisitSearch" />
          <el-select v-model="vcriteria.entityType" placeholder="实体类型" clearable style="width: 130px" @change="handleVisitSearch">
            <el-option label="product" value="product" />
            <el-option label="page" value="page" />
            <el-option label="blog" value="blog" />
            <el-option label="case" value="case" />
            <el-option label="home" value="home" />
            <el-option label="social_click" value="social_click" />
          </el-select>
          <el-select v-model="vcriteria.device" placeholder="设备" clearable style="width: 120px" @change="handleVisitSearch">
            <el-option label="desktop" value="desktop" />
            <el-option label="mobile" value="mobile" />
            <el-option label="tablet" value="tablet" />
            <el-option label="bot" value="bot" />
          </el-select>
          <el-input v-model="vcriteria.keyword" placeholder="关键词/路径" clearable style="width: 160px" @keyup.enter="handleVisitSearch" />
          <el-date-picker v-model="vcriteria.dateRange" type="daterange" value-format="YYYY-MM-DD" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" clearable style="width: 260px" @change="handleVisitSearch" />
          <el-button type="primary" @click="handleVisitSearch">查询</el-button>
        </div>

        <el-table :data="visits" v-loading="visitsLoading" stripe :row-key="visitRowKey">
          <el-table-column prop="created_at" label="时间" width="170">
            <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
          </el-table-column>
          <el-table-column prop="country" label="国家" width="100" />
          <el-table-column prop="language" label="语言" width="90" />
          <el-table-column prop="device" label="设备" width="90" />
          <el-table-column prop="device_model" label="设备型号" width="120" />
          <el-table-column prop="browser" label="浏览器" width="100" />
          <el-table-column prop="ip" label="IP" width="130" />
          <el-table-column prop="source" label="来源" width="100" />
          <el-table-column prop="medium" label="Medium" width="100" />
          <el-table-column prop="utm_campaign" label="广告Campaign" width="130" />
          <el-table-column prop="entity_type" label="实体类型" width="100" />
          <el-table-column prop="entity_name" label="产品/页面" min-width="160" />
          <el-table-column prop="path" label="路径" min-width="200" />
          <el-table-column prop="status" label="状态码" width="90" />
          <el-table-column prop="latency_ms" label="耗时(ms)" width="90" />
          <el-table-column prop="src" label="日志表" width="160">
            <template #default="{ row }">{{ row.src || 'visit_logs' }}</template>
          </el-table-column>
        </el-table>

        <el-pagination class="pagination" v-model:current-page="vpage" v-model:page-size="pageSize" :total="vtotal" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="searchVisits" @size-change="resetVisitPage" />
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { operationLogApi, analyticsApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import { formatDateTime } from '@/utils/format'

const activeTab = ref<'operation' | 'visits'>('operation')
const pageSize = useAdminPageSize()

// ============ 后台操作日志 ============
const criteria = ref({ module: '', operation: '', entityType: '', ip: '', userId: '', dateRange: null as null | [string, string] })
const logs = ref<any[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)

const opLabels: Record<string, string> = { create: '创建', update: '更新', delete: '删除', publish: '发布', unpublish: '下线', upload: '上传', login: '登录', logout: '登出' }
function operationLabel(op: string) { return opLabels[op] || op }
function formatTime(t: string) { return formatDateTime(t) }
function opRowKey(row: any) { return `${row.src || 'operation_logs'}-${row.id}` }

// 点击"查询"按钮/回车触发：重置到第一页再加载
function handleSearch() { page.value = 1; loadData() }
// 分页控件翻页触发：保留当前页码直接加载
function resetPage() { page.value = 1; loadData() }

async function loadData() {
  loading.value = true
  try {
    const params: any = { page: page.value, pageSize: pageSize.value, module: criteria.value.module, operation: criteria.value.operation, entity_type: criteria.value.entityType, ip: criteria.value.ip, user_id: criteria.value.userId }
    if (criteria.value.dateRange) { params.start_date = criteria.value.dateRange[0]; params.end_date = criteria.value.dateRange[1] }
    const res = await operationLogApi.list(params)
    logs.value = (res as any)?.items || []
    total.value = (res as any)?.total || 0
  } finally { loading.value = false }
}

// ============ 门户访问日志 ============
const vcriteria = ref({ ip: '', country: '', source: '', entityType: '', device: '', keyword: '', dateRange: null as null | [string, string] })
const visits = ref<any[]>([])
const visitsLoading = ref(false)
const vtotal = ref(0)
const vpage = ref(1)

function visitRowKey(row: any) { return `${row.src || 'visit_logs'}-${row.id}-${row.created_at}` }
// 点击"查询"按钮/回车触发：重置到第一页
function handleVisitSearch() { vpage.value = 1; searchVisits() }
// 分页控件翻页触发：保留当前页码
function resetVisitPage() { vpage.value = 1; searchVisits() }

async function searchVisits() {
  visitsLoading.value = true
  try {
    const params: any = { page: vpage.value, pageSize: pageSize.value, ip: vcriteria.value.ip, country: vcriteria.value.country, source: vcriteria.value.source, entity_type: vcriteria.value.entityType, device: vcriteria.value.device, keyword: vcriteria.value.keyword }
    if (vcriteria.value.dateRange) { params.start_date = vcriteria.value.dateRange[0]; params.end_date = vcriteria.value.dateRange[1] }
    const res = await analyticsApi.visitLogs(params)
    visits.value = (res as any)?.items || []
    vtotal.value = (res as any)?.total || 0
  } finally { visitsLoading.value = false }
}

onMounted(() => { loadData(); searchVisits() })
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 12px; align-items: center; margin-bottom: 16px; }
.pagination { margin-top: 16px; justify-content: flex-end; }
</style>

