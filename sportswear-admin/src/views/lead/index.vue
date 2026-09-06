<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input
        v-model="keyword"
        placeholder="搜索姓名 / 邮箱 / 公司"
        clearable
        style="width: 260px"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="status" placeholder="状态" clearable style="width: 160px" @change="handleSearch">
        <el-option v-for="o in enumOptions('lead.status')" :key="o.value" :label="o.label" :value="o.value" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <el-button :loading="backfilling" @click="handleBackfill">回填 IP 国家</el-button>
    </div>

    <el-table :data="leads" v-loading="loading" stripe>
      <el-table-column prop="name" label="姓名" width="120" />
      <el-table-column prop="company" label="公司" min-width="140" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column prop="country" label="国家" width="110" />
      <el-table-column label="IP / IP国家" min-width="150">
        <template #default="{ row }">
          <div>{{ row.ip || '-' }}</div>
          <div class="muted">{{ countryName(row.ip_country) || '-' }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="110">
        <template #default="{ row }">
          <el-tag size="small">{{ enumLabel('lead.status', row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="score" label="评分" width="90">
        <template #default="{ row }">
          <el-tag :type="scoreType(row.score)" size="small">{{ row.score }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160">
        <template #default="{ row }: any">
          <el-button size="small" @click="openDetail(row)">详情</el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="询盘详情" width="720px">
      <template v-if="current">
        <div class="detail-actions">
          <el-button size="small" type="primary" plain :disabled="!current.ip" @click="openIPVisits(current)">
            查看 IP 访问记录
          </el-button>
        </div>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="姓名">{{ current.name }}</el-descriptions-item>
          <el-descriptions-item label="公司">{{ current.company || '-' }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ current.email }}</el-descriptions-item>
          <el-descriptions-item label="电话">{{ current.phone || '-' }}</el-descriptions-item>
          <el-descriptions-item label="国家">{{ current.country || '-' }}</el-descriptions-item>
          <el-descriptions-item label="IP 国家">{{ countryName(current.ip_country) || '-' }}</el-descriptions-item>
          <el-descriptions-item label="IP" :span="2">{{ current.ip || '-' }}</el-descriptions-item>
          <el-descriptions-item label="项目类型">{{ current.project_type || '-' }}</el-descriptions-item>
          <el-descriptions-item label="产品分类">{{ current.product_category || '-' }}</el-descriptions-item>
          <el-descriptions-item label="数量">{{ current.quantity || 0 }}</el-descriptions-item>
          <el-descriptions-item label="来源">{{ current.source || '-' }}</el-descriptions-item>
          <el-descriptions-item label="评分">{{ current.score }}</el-descriptions-item>
          <el-descriptions-item label="留言" :span="2">{{ current.message || '-' }}</el-descriptions-item>
          <el-descriptions-item label="附件" :span="2">
            <template v-if="attachmentList.length">
              <div class="attachments">
                <template v-for="(att, i) in attachmentList" :key="i">
                  <el-image
                    v-if="isImage(att.url)"
                    :src="att.url"
                    :preview-src-list="imageUrls"
                    :initial-index="imageUrls.indexOf(att.url)"
                    fit="cover"
                    class="attachment-image"
                    preview-teleported
                  />
                  <el-link v-else :href="att.url" target="_blank" type="primary">
                    {{ att.name || att.url }}<span v-if="att.size" class="muted-size">（{{ formatAttSize(att.size) }}）</span>
                  </el-link>
                </template>
              </div>
            </template>
            <span v-else>-</span>
          </el-descriptions-item>
        </el-descriptions>

        <el-divider content-position="left">跟进记录</el-divider>
        <div v-if="current.follow_ups?.length">
          <div v-for="fu in current.follow_ups" :key="fu.id" class="followup-item">
            <div class="followup-content">{{ fu.content }}</div>
            <div class="followup-meta">{{ fu.method }} · {{ formatDateTime(fu.created_at) }}</div>
          </div>
        </div>
        <el-empty v-else description="暂无跟进记录" :image-size="60" />
      </template>
    </el-dialog>

    <!-- IP 访问记录抽屉（现在 / 过去 / 历史访问轨迹） -->
    <el-drawer v-model="ipVisitsVisible" title="IP 访问记录" size="60%">
      <template v-if="ipVisitsSummary">
        <el-descriptions :column="3" border size="small" class="ip-summary">
          <el-descriptions-item label="IP">{{ ipVisitsSummary.ip || '-' }}</el-descriptions-item>
          <el-descriptions-item label="国家">{{ countryName(ipVisitsSummary.country) || ipVisitsSummary.country || '-' }}</el-descriptions-item>
          <el-descriptions-item label="关联询盘">{{ ipVisitsSummary.lead_count || 0 }}</el-descriptions-item>
          <el-descriptions-item label="首次访问">{{ formatDateTime(ipVisitsSummary.first_visit_at) }}</el-descriptions-item>
          <el-descriptions-item label="最近访问">{{ formatDateTime(ipVisitsSummary.last_visit_at) }}</el-descriptions-item>
          <el-descriptions-item label="浏览 / 独立访客">{{ ipVisitsSummary.page_views || 0 }} / {{ ipVisitsSummary.unique_visitors || 0 }}</el-descriptions-item>
        </el-descriptions>
      </template>

      <el-table :data="ipVisitsRows" v-loading="ipVisitsLoading" stripe size="small">
        <el-table-column prop="created_at" label="时间" width="170" :formatter="(r: any) => formatDateTime(r.created_at)" />
        <el-table-column prop="path" label="访问路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="entity_name" label="实体" min-width="140" show-overflow-tooltip />
        <el-table-column prop="source" label="来源" width="90" />
        <el-table-column prop="device" label="设备" width="80" />
        <el-table-column prop="country" label="国家" width="80" />
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="ipVisitsPage"
        v-model:page-size="ipVisitsPageSize"
        :total="ipVisitsTotal"
        layout="total, prev, pager, next"
        :page-sizes="[10, 20, 50]"
        @current-change="loadIPVisits"
        @size-change="loadIPVisits"
      />
    </el-drawer>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { leadApi, analyticsApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import { useEnumDict } from '@/composables/useEnumDict'
import { formatDateTime } from '@/utils/format'
import { countryName } from '@/utils/geo'
import type { Lead } from '@/types'

const leads = ref<Lead[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const status = ref('')

const detailVisible = ref(false)
const current = ref<Lead | null>(null)

// ==================== 回填 IP 国家 ====================
const backfilling = ref(false)
async function handleBackfill() {
  backfilling.value = true
  try {
    const res = await leadApi.backfillIPCountry()
    ElMessage.success(`已回填 ${res.updated} 条询盘的 IP 国家`)
    await loadData()
  } finally {
    backfilling.value = false
  }
}

// ==================== IP 访问记录（现在 / 过去 / 历史） ====================
const ipVisitsVisible = ref(false)
const ipVisitsLoading = ref(false)
const ipVisitsSummary = ref<any>(null)
const ipVisitsRows = ref<any[]>([])
const ipVisitsTotal = ref(0)
const ipVisitsPage = ref(1)
const ipVisitsPageSize = ref(20)
const ipVisitsIp = ref('')

function openIPVisits(row: Lead) {
  if (!row.ip) return
  ipVisitsIp.value = row.ip
  ipVisitsPage.value = 1
  ipVisitsVisible.value = true
  loadIPVisits()
}

async function loadIPVisits() {
  if (!ipVisitsIp.value) return
  ipVisitsLoading.value = true
  try {
    const res = await analyticsApi.ipVisits({
      ip: ipVisitsIp.value,
      days: 90,
      page: ipVisitsPage.value,
      pageSize: ipVisitsPageSize.value,
    })
    ipVisitsSummary.value = res.summary || null
    ipVisitsRows.value = res.items || []
    ipVisitsTotal.value = res.total || 0
  } finally {
    ipVisitsLoading.value = false
  }
}

// ==================== 附件展示（门户询盘上传的图片/文件，JSON 数组字符串） ====================
interface LeadAttachment {
  url: string
  name?: string
  size?: number
  type?: string
}

const attachmentList = computed<LeadAttachment[]>(() => {
  const raw = current.value?.attachments
  if (!raw) return []
  try {
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed.filter((i) => i && typeof i.url === 'string') : []
  } catch {
    return []
  }
})

const imageUrls = computed(() =>
  attachmentList.value.filter((a) => isImage(a.url)).map((a) => a.url),
)

const IMAGE_EXTS = ['.jpg', '.jpeg', '.png', '.gif', '.webp']
function isImage(url: string) {
  const path = url.split('?')[0].toLowerCase()
  return IMAGE_EXTS.some((ext) => path.endsWith(ext))
}

function formatAttSize(bytes: number) {
  if (bytes >= 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + 'MB'
  if (bytes >= 1024) return (bytes / 1024).toFixed(0) + 'KB'
  return bytes + 'B'
}

// 询盘状态由「字典管理」驱动（lead.status）
const { ensureLoaded: loadEnumDict, options: enumOptions, label: enumLabel } = useEnumDict()

async function loadData() {
  loading.value = true
  try {
    const result = await leadApi.list({
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value,
      status: status.value,
    })
    leads.value = result.items
    total.value = result.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

async function openDetail(row: Lead) {
  current.value = await leadApi.get(row.id)
  detailVisible.value = true
}

async function handleDelete(row: Lead) {
  await ElMessageBox.confirm(`确定删除询盘 ${row.name} 吗？`, '警告', { type: 'warning' })
  await leadApi.delete(row.id)
  ElMessage.success('已删除')
  loadData()
}

function scoreType(score: number) {
  return score >= 80 ? 'danger' : score >= 50 ? 'warning' : 'info'
}

onMounted(() => { loadEnumDict(); loadData() })
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
.followup-item {
  padding: 10px;
  background: #f9fafb;
  border-radius: 6px;
  margin-bottom: 8px;
}
.followup-content {
  font-size: 14px;
}
.followup-meta {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 4px;
}
.attachments {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}
.attachment-image {
  width: 72px;
  height: 72px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
  cursor: pointer;
}
.muted-size {
  color: #9ca3af;
  font-size: 12px;
}
.muted {
  color: #9ca3af;
  font-size: 12px;
}
.detail-actions {
  margin-bottom: 12px;
}
.ip-summary {
  margin-bottom: 16px;
}
</style>