<template>
  <el-card shadow="never">
    <div class="toolbar">
      <span class="title">门户闭环健康检查</span>
      <div class="spacer" />
      <el-button type="primary" :loading="loading" @click="loadData">重新检查</el-button>
    </div>

    <el-alert
      type="info"
      show-icon
      :closable="false"
      class="mb"
      title="一键检测「导航死链 / 页面漏配导航 / 路由漏配 SEO / 导航路径漂移」等闭环断点，帮助快速定位并修复。"
    />

    <div v-if="report" class="stat-row">
      <div class="stat-card" :class="{ bad: report.total_issues > 0 }">
        <div class="stat-num">{{ report.total_issues }}</div>
        <div class="stat-label">待处理问题</div>
      </div>
      <div class="stat-card"><div class="stat-num">{{ report.routes }}</div><div class="stat-label">路由总数</div></div>
      <div class="stat-card"><div class="stat-num">{{ report.nav_count }}</div><div class="stat-label">导航总数</div></div>
      <div class="stat-card"><div class="stat-num">{{ report.page_count }}</div><div class="stat-label">已发布页面</div></div>
    </div>

    <el-empty v-if="report && report.total_issues === 0" description="门户闭环完整，未发现问题" />

    <el-tabs v-else-if="report" v-model="activeType">
      <el-tab-pane v-for="g in groups" :key="g.type" :label="`${g.label} (${g.items.length})`" :name="g.type">
        <el-table :data="g.items" stripe size="small" empty-text="该类问题已清零">
          <el-table-column prop="title" label="问题" min-width="180" />
          <el-table-column label="路径 / 路由" width="160">
            <template #default="{ row }">{{ row.path || row.route || '—' }}</template>
          </el-table-column>
          <el-table-column prop="detail" label="说明" min-width="280" show-overflow-tooltip />
          <el-table-column label="操作" width="150">
            <template #default="{ row }">
              <el-button v-if="row.type === 'missing_seo'" size="small" type="primary" @click="goSeo(row)">去配置 SEO</el-button>
              <el-button v-else-if="row.type === 'missing_nav'" size="small" type="primary" @click="goNav">去加导航</el-button>
              <el-button v-else size="small" type="primary" @click="goNav">去修导航</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { portalHealthApi } from '@/api'
import type { PortalHealthReport } from '@/types'

const router = useRouter()
const loading = ref(false)
const report = ref<PortalHealthReport | null>(null)
const activeType = ref('dead_link')

const TYPE_LABELS: Record<string, string> = {
  dead_link: '导航死链',
  missing_nav: '缺导航',
  missing_seo: '缺 SEO',
  nav_drift: '路径漂移',
}

const groups = computed(() => {
  if (!report.value) return []
  const order = ['dead_link', 'missing_nav', 'missing_seo', 'nav_drift']
  return order.map((t) => ({
    type: t,
    label: TYPE_LABELS[t] || t,
    items: (report.value!.issues || []).filter((i) => i.type === t),
  }))
})

async function loadData() {
  loading.value = true
  try {
    report.value = await portalHealthApi.healthCheck()
  } catch {
    report.value = null
  } finally {
    loading.value = false
  }
}

function goSeo(row: any) {
  router.push({ path: '/seo', query: row.route ? { route: row.route } : {} })
}
function goNav() {
  router.push('/navigations')
}

onMounted(loadData)
</script>

<style scoped>
.toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
.title { font-size: 16px; font-weight: 600; }
.spacer { flex: 1; }
.mb { margin-bottom: 16px; }
.stat-row { display: flex; gap: 16px; margin-bottom: 16px; }
.stat-card {
  flex: 1; background: #f5f7fa; border: 1px solid #e4e7ed; border-radius: 8px;
  padding: 16px 20px; text-align: center;
}
.stat-card.bad { background: #fef0f0; border-color: #fbc4c4; }
.stat-num { font-size: 26px; font-weight: 700; line-height: 1.2; }
.stat-card.bad .stat-num { color: #f56c6c; }
.stat-label { font-size: 13px; color: #909399; margin-top: 4px; }
</style>
