<template>
  <div>
    <!-- 工具栏：时间范围 + 刷新（最大 90 天，避免跨过多月度分表影响性能） -->
    <div class="toolbar">
      <el-radio-group v-model="days" size="small" style="margin-right: 12px">
        <el-radio-button :value="7">最近 7 天</el-radio-button>
        <el-radio-button :value="30">最近 30 天</el-radio-button>
        <el-radio-button :value="90">最近 90 天</el-radio-button>
      </el-radio-group>
      <el-button type="primary" @click="loadCharts" :loading="chartsLoading" plain>刷新</el-button>
      <span class="hint">数据按月分表存储，最大可查 90 天</span>
    </div>

    <!-- 概览卡片 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="6" v-for="c in cards" :key="c.label">
        <el-card shadow="hover"><div class="stat-card"><div class="stat-value">{{ c.value }}</div><div class="stat-label">{{ c.label }}</div></div></el-card>
      </el-col>
    </el-row>

    <!-- 流量趋势 + 渠道来源 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12"><el-card shadow="hover"><template #header>流量趋势（访次 / 独立 IP）</template><el-alert v-if="sectionErrors['overview']" :title="'数据加载失败: ' + sectionErrors['overview']" type="warning" show-icon :closable="false" style="margin-bottom:12px" /><div ref="trendRef" style="height:280px"></div></el-card></el-col>
      <el-col :span="12"><el-card shadow="hover"><template #header>渠道来源分布</template><el-alert v-if="sectionErrors['sources']" :title="'数据加载失败: ' + sectionErrors['sources']" type="warning" show-icon :closable="false" style="margin-bottom:12px" /><div ref="sourceRef" style="height:280px"></div></el-card></el-col>
    </el-row>

    <!-- 热门页面 + 热门产品 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12"><el-card shadow="hover"><template #header>热门页面（浏览量 TOP 10）</template><el-alert v-if="sectionErrors['topPages']" :title="'数据加载失败: ' + sectionErrors['topPages']" type="warning" show-icon :closable="false" style="margin-bottom:12px" /><div ref="pagesRef" style="height:320px"></div></el-card></el-col>
      <el-col :span="12"><el-card shadow="hover"><template #header>热门产品（浏览量 TOP 10）</template><el-alert v-if="sectionErrors['topProducts']" :title="'数据加载失败: ' + sectionErrors['topProducts']" type="warning" show-icon :closable="false" style="margin-bottom:12px" /><div ref="productsRef" style="height:320px"></div></el-card></el-col>
    </el-row>

    <!-- 国家分布 + 设备分布 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12"><el-card shadow="hover"><template #header>国家/地区分布</template><el-alert v-if="sectionErrors['countries']" :title="'数据加载失败: ' + sectionErrors['countries']" type="warning" show-icon :closable="false" style="margin-bottom:12px" /><div ref="countryRef" style="height:320px"></div></el-card></el-col>
      <el-col :span="12"><el-card shadow="hover"><template #header>设备类型分布</template><el-alert v-if="sectionErrors['devices']" :title="'数据加载失败: ' + sectionErrors['devices']" type="warning" show-icon :closable="false" style="margin-bottom:12px" /><div ref="deviceRef" style="height:320px"></div></el-card></el-col>
    </el-row>

    <!-- UTM 广告归因 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="24"><el-card shadow="hover"><template #header>广告归因（UTM Campaign 点击来源）</template>
        <el-alert v-if="sectionErrors['utmCampaigns']" :title="'数据加载失败: ' + sectionErrors['utmCampaigns']" type="warning" show-icon :closable="false" style="margin-bottom:12px" />
        <el-table :data="utmCampaigns" v-loading="utmLoading" stripe style="width:100%">
          <el-table-column prop="campaign" label="Campaign" min-width="160" />
          <el-table-column prop="source" label="来源" width="120" />
          <el-table-column prop="medium" label="Medium" width="100" />
          <el-table-column prop="visits" label="访次" width="100" />
          <el-table-column prop="unique_visitors" label="独立访客" width="110" />
        </el-table></el-card></el-col>
    </el-row>

    <!-- 社交媒体点击 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="24"><el-card shadow="hover"><template #header>社交媒体点击（外部链接跳转）</template>
        <el-alert v-if="sectionErrors['socialClicks']" :title="'数据加载失败: ' + sectionErrors['socialClicks']" type="warning" show-icon :closable="false" style="margin-bottom:12px" />
        <el-table :data="socialClicks" v-loading="socialLoading" stripe style="width:100%">
          <el-table-column prop="platform" label="平台" min-width="140" />
          <el-table-column prop="clicks" label="点击量" width="100" />
          <el-table-column prop="unique_visitors" label="独立访客" width="110" />
        </el-table></el-card></el-col>
    </el-row>

    <!-- 转化漏斗 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="24"><el-card shadow="hover"><template #header>转化漏斗（页面浏览 → 产品浏览 → 询盘）</template>
        <el-alert v-if="sectionErrors['funnel']" :title="'数据加载失败: ' + sectionErrors['funnel']" type="warning" show-icon :closable="false" style="margin-bottom:12px" />
        <div v-loading="funnelLoading" class="funnel-container">
          <div class="funnel-step">
            <div class="funnel-value">{{ funnelData.page_views || 0 }}</div>
            <div class="funnel-label">页面浏览</div>
          </div>
          <div class="funnel-arrow">→</div>
          <div class="funnel-step">
            <div class="funnel-value">{{ funnelData.product_views || 0 }}</div>
            <div class="funnel-label">产品浏览</div>
            <div class="funnel-rate">{{ (funnelData.product_rate || 0).toFixed(1) }}%</div>
          </div>
          <div class="funnel-arrow">→</div>
          <div class="funnel-step">
            <div class="funnel-value">{{ funnelData.leads || 0 }}</div>
            <div class="funnel-label">询盘提交</div>
            <div class="funnel-rate">{{ (funnelData.lead_rate || 0).toFixed(1) }}%</div>
          </div>
        </div>
      </el-card></el-col>
    </el-row>

    <!-- 访客旅程查询 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="24"><el-card shadow="hover"><template #header>访客旅程查询</template>
        <div class="journey-toolbar">
          <el-input v-model="journeyQuery" placeholder="输入 visitor_id 或 IP 查询" style="width: 300px" @keyup.enter="loadJourney" />
          <el-button type="primary" @click="loadJourney" :loading="journeyLoading">查询</el-button>
        </div>
        <el-table :data="journeyRows" v-loading="journeyLoading" stripe style="width:100%" size="small">
          <el-table-column prop="created_at" label="时间" width="170" :formatter="(r: any) => formatTime(r.created_at)" />
          <el-table-column prop="path" label="访问路径" min-width="200" show-overflow-tooltip />
          <el-table-column prop="entity_name" label="实体名称" min-width="150" show-overflow-tooltip />
          <el-table-column prop="source" label="来源" width="100" />
          <el-table-column prop="country" label="国家" width="80" />
          <el-table-column prop="device" label="设备" width="80" />
        </el-table>
        <el-empty v-if="!journeyLoading && journeyRows.length === 0 && journeyQueried" description="暂无访问记录" :image-size="60" />
      </el-card></el-col>
    </el-row>

    <!-- IP-询盘关联 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="24"><el-card shadow="hover"><template #header>IP-询盘关联</template>
        <div class="journey-toolbar">
          <el-input v-model="ipQuery" placeholder="输入 IP 地址查询关联询盘" style="width: 300px" @keyup.enter="loadIPLeads" />
          <el-button type="primary" @click="loadIPLeads" :loading="ipLeadsLoading">查询</el-button>
        </div>
        <el-table :data="ipLeadsRows" v-loading="ipLeadsLoading" stripe style="width:100%" size="small">
          <el-table-column prop="created_at" label="询盘时间" width="170" :formatter="(r: any) => formatTime(r.created_at)" />
          <el-table-column prop="name" label="姓名" width="100" />
          <el-table-column prop="company" label="公司" min-width="120" show-overflow-tooltip />
          <el-table-column prop="email" label="邮箱" min-width="160" />
          <el-table-column prop="country" label="国家" width="80" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag size="small" :type="row.status === 'new' ? 'primary' : 'info'">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="score" label="评分" width="60" />
          <el-table-column prop="source" label="来源" width="100" />
        </el-table>
        <el-empty v-if="!ipLeadsLoading && ipLeadsRows.length === 0 && ipQueried" description="该 IP 暂无关联询盘" :image-size="60" />
      </el-card></el-col>
    </el-row>

    <!-- 隐私合规洞察 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="24"><el-card shadow="hover"><template #header>隐私合规洞察（Cookie 同意与转化对比）</template>
        <el-alert v-if="sectionErrors['consent']" :title="'数据加载失败: ' + sectionErrors['consent']" type="warning" show-icon :closable="false" style="margin-bottom:12px" />
        <div v-loading="consentLoading" class="consent-container">
          <div class="consent-stats">
            <div class="consent-stat-item">
              <div class="consent-stat-value">{{ consentData.consent_rate?.toFixed(1) || 0 }}%</div>
              <div class="consent-stat-label">同意率</div>
            </div>
            <div class="consent-stat-item">
              <div class="consent-stat-value">{{ consentData.granted_count || 0 }}</div>
              <div class="consent-stat-label">同意人数</div>
            </div>
            <div class="consent-stat-item">
              <div class="consent-stat-value">{{ consentData.denied_count || 0 }}</div>
              <div class="consent-stat-label">拒绝人数</div>
            </div>
            <div class="consent-stat-item">
              <div class="consent-stat-value">{{ consentData.anonymized_ratio?.toFixed(1) || 0 }}%</div>
              <div class="consent-stat-label">匿名化占比</div>
            </div>
          </div>
          <el-divider />
          <div class="consent-comparison">
            <div class="consent-compare-item">
              <div class="consent-compare-label">同意用户转化率</div>
              <div class="consent-compare-value">{{ consentData.granted_conversion?.toFixed(2) || 0 }}%</div>
              <div class="consent-compare-detail">{{ consentData.granted_leads || 0 }} 个询盘 / {{ consentData.granted_count || 0 }} 人</div>
            </div>
            <div class="consent-compare-item">
              <div class="consent-compare-label">拒绝用户转化率</div>
              <div class="consent-compare-value">{{ consentData.denied_conversion?.toFixed(2) || 0 }}%</div>
              <div class="consent-compare-detail">{{ consentData.denied_leads || 0 }} 个询盘 / {{ consentData.denied_count || 0 }} 人</div>
            </div>
          </div>
        </div>
      </el-card></el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
// 注：门户访问日志明细的筛选查询已在「操作日志」页面的「门户访问日志」Tab 提供，此处不再重复展示
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { analyticsApi } from '@/api'
import * as echarts from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { formatDateTime } from '@/utils/format'

echarts.use([BarChart, LineChart, PieChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const days = ref(7)
const chartsLoading = ref(false)

// 各模块错误状态（非阻断：某模块失败不影响其他模块渲染）
const sectionErrors = ref<Record<string, string>>({})
function setModuleError(module: string, err: unknown) {
  const msg = err instanceof Error ? err.message : '加载失败'
  sectionErrors.value[module] = msg
  console.warn(`[analytics] ${module}:`, err)
}
function clearModuleError(module: string) {
  delete sectionErrors.value[module]
}

const cards = ref([
  { label: '总访问量', value: 0 },
  { label: '独立 IP', value: 0 },
  { label: 'Bot 流量', value: 0 },
  { label: '询盘转化率', value: '0%' },
  { label: '移动端占比', value: '0%' },
])

const trendRef = ref(); const sourceRef = ref(); const pagesRef = ref(); const productsRef = ref(); const countryRef = ref(); const deviceRef = ref()
const charts: echarts.ECharts[] = []
const fmt = (v: any, def: number = 0): number => (v ? Number(v) : def)
const formatTime = (iso?: string | null) => formatDateTime(iso)
const initChart = (el: any, option: any): echarts.ECharts => { const ch = echarts.init(el); ch.setOption(option); charts.push(ch); return ch }
const getChart = (el: any): echarts.ECharts | null => (el ? (echarts.getInstanceByDom(el) as echarts.ECharts) : null)

async function loadCharts() {
  chartsLoading.value = true
  // 清除旧错误
  ;['overview', 'sources', 'topPages', 'topProducts', 'countries', 'devices', 'utmCampaigns', 'socialClicks'].forEach(m => clearModuleError(m))
  try {
    const overview = await analyticsApi.overview({ days: days.value })
    clearModuleError('overview')
    cards.value[0].value = fmt(overview?.total_visits)
    cards.value[1].value = fmt(overview?.unique_ips)
    cards.value[2].value = fmt(overview?.bot_visits)
    cards.value[4].value = overview?.mobile_ratio != null ? (Number(overview.mobile_ratio) * 100).toFixed(0) + '%' : '0%'
    cards.value[3].value = overview?.conversion_rate != null ? Number(overview.conversion_rate).toFixed(2) + '%' : '0%'

    const trend = (overview?.daily_trend || []) as any[]
    const tr = getChart(trendRef.value) || initChart(trendRef.value, {})
    tr.setOption({
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      legend: { data: ['访次', '独立 IP'] },
      xAxis: { type: 'category', data: trend.map((d: any) => d.date) },
      yAxis: [{ type: 'value', name: '访次' }, { type: 'value', name: '独立 IP' }],
      series: [
        { name: '访次', type: 'bar', data: trend.map((d: any) => d.visits) },
        { name: '独立 IP', type: 'line', yAxisIndex: 1, smooth: true, data: trend.map((d: any) => d.unique_visitors || d.unique_ips) },
      ],
    })
  } catch (e) { setModuleError('overview', e) }

  try {
    const sources = await analyticsApi.sources({ days: days.value })
    clearModuleError('sources')
    const sr = getChart(sourceRef.value) || initChart(sourceRef.value, {})
    sr.setOption({
      tooltip: { trigger: 'item' },
      legend: { left: 'center', top: 'bottom', data: (sources || []).map((i: any) => i.source) },
      series: [{ type: 'pie', radius: ['35%', '70%'], label: { formatter: '{b}: {c} ({d}%)' }, data: (sources || []).map((i: any) => ({ name: i.source, value: i.visits || i.count })) }],
    })
  } catch (e) { setModuleError('sources', e) }

  try {
    const topPages = await analyticsApi.topPages({ days: days.value })
    clearModuleError('topPages')
    const pr = getChart(pagesRef.value) || initChart(pagesRef.value, {})
    pr.setOption({ tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } }, yAxis: { type: 'category', data: (topPages || []).map((i: any) => i.path) }, xAxis: { type: 'value' }, series: [{ type: 'bar', data: (topPages || []).map((i: any) => i.views) }] })
  } catch (e) { setModuleError('topPages', e) }

  try {
    const topProducts = await analyticsApi.topProducts({ days: days.value })
    clearModuleError('topProducts')
    const pc = getChart(productsRef.value) || initChart(productsRef.value, {})
    pc.setOption({ tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } }, yAxis: { type: 'category', data: (topProducts || []).map((i: any) => i.name) }, xAxis: { type: 'value' }, series: [{ type: 'bar', data: (topProducts || []).map((i: any) => i.views) }] })
  } catch (e) { setModuleError('topProducts', e) }

  try {
    const countries = await analyticsApi.countries({ days: days.value })
    clearModuleError('countries')
    const cr = getChart(countryRef.value) || initChart(countryRef.value, {})
    cr.setOption({ tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } }, yAxis: { type: 'category', inverse: true, data: (countries || []).slice(0, 10).map((i: any) => i.country || i.region) }, xAxis: { type: 'value' }, series: [{ type: 'bar', data: (countries || []).slice(0, 10).map((i: any) => i.visits || i.count) }] })
  } catch (e) { setModuleError('countries', e) }

  try {
    const devices = await analyticsApi.devices({ days: days.value })
    clearModuleError('devices')
    cards.value[4].value = devices?.mobile_ratio ? ((Number(devices.mobile_ratio) * 100).toFixed(0) + '%') : '0%'
    const dr = getChart(deviceRef.value) || initChart(deviceRef.value, {})
    dr.setOption({ tooltip: { trigger: 'item' }, legend: { left: 'center', top: 'bottom', data: (devices?.devices || []).map((i: any) => i.device) }, series: [{ type: 'pie', radius: ['35%', '70%'], label: { formatter: '{b}: {c} ({d}%)' }, data: (devices?.devices || []).map((i: any) => ({ name: i.device, value: i.visits })) }] })
  } catch (e) { setModuleError('devices', e) }

  try {
    const utm = await analyticsApi.utmCampaigns({ days: days.value })
    clearModuleError('utmCampaigns')
    utmCampaigns.value = utm || []
  } catch (e) { setModuleError('utmCampaigns', e) }
  try {
    const social = await analyticsApi.socialClicks({ days: days.value })
    clearModuleError('socialClicks')
    socialClicks.value = social || []
  } catch (e) { setModuleError('socialClicks', e) }

  chartsLoading.value = false
}

const utmCampaigns = ref<any[]>([])
const utmLoading = ref(false)
const socialClicks = ref<any[]>([])
const socialLoading = ref(false)

// ============ 转化漏斗 ============
const funnelLoading = ref(false)
const funnelData = ref<any>({})

async function loadFunnel() {
  funnelLoading.value = true
  try {
    funnelData.value = await analyticsApi.conversionFunnel({ days: days.value })
    clearModuleError('funnel')
  } catch (e) { setModuleError('funnel', e) } finally { funnelLoading.value = false }
}

// ============ 访客旅程 ============
const journeyQuery = ref('')
const journeyLoading = ref(false)
const journeyRows = ref<any[]>([])
const journeyQueried = ref(false)

async function loadJourney() {
  if (!journeyQuery.value.trim()) return
  journeyLoading.value = true
  journeyQueried.value = true
  try {
    const rows = await analyticsApi.visitorJourney({ visitor_id: journeyQuery.value.trim(), days: 90 })
    clearModuleError('journey')
    journeyRows.value = rows || []
  } catch (e) { setModuleError('journey', e); journeyRows.value = [] } finally { journeyLoading.value = false }
}

// ============ IP-询盘关联 ============
const ipQuery = ref('')
const ipLeadsLoading = ref(false)
const ipLeadsRows = ref<any[]>([])
const ipQueried = ref(false)

async function loadIPLeads() {
  if (!ipQuery.value.trim()) return
  ipLeadsLoading.value = true
  ipQueried.value = true
  try {
    const rows = await analyticsApi.ipLeads({ ip: ipQuery.value.trim(), days: 90 })
    clearModuleError('ipLeads')
    ipLeadsRows.value = rows || []
  } catch (e) { setModuleError('ipLeads', e); ipLeadsRows.value = [] } finally { ipLeadsLoading.value = false }
}

// ============ 隐私合规洞察 ============
const consentLoading = ref(false)
const consentData = ref<any>({})

async function loadConsentInsights() {
  consentLoading.value = true
  try {
    consentData.value = await analyticsApi.consentInsights({ days: days.value })
    clearModuleError('consent')
  } catch (e) { setModuleError('consent', e) } finally { consentLoading.value = false }
}

onMounted(() => { loadCharts(); loadFunnel(); loadConsentInsights() })
watch(days, () => { loadCharts(); loadFunnel(); loadConsentInsights() })
onUnmounted(() => charts.forEach((c) => c.dispose()))
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 12px; align-items: center; }
.toolbar .hint { font-size: 12px; color: #909399; }
.mt-20 { margin-top: 20px; }
.journey-toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.funnel-container { display: flex; align-items: center; justify-content: center; gap: 20px; padding: 20px; }
.funnel-step { text-align: center; min-width: 120px; }
.funnel-value { font-size: 28px; font-weight: 700; color: #1e3a8a; }
.funnel-label { font-size: 14px; color: #6b7280; margin-top: 4px; }
.funnel-rate { font-size: 12px; color: #67c23a; margin-top: 2px; }
.funnel-arrow { font-size: 24px; color: #909399; }
.consent-container { padding: 10px 0; }
.consent-stats { display: flex; justify-content: space-around; padding: 10px 0; }
.consent-stat-item { text-align: center; }
.consent-stat-value { font-size: 24px; font-weight: 700; color: #1e3a8a; }
.consent-stat-label { font-size: 13px; color: #6b7280; margin-top: 4px; }
.consent-comparison { display: flex; justify-content: center; gap: 60px; padding: 20px 0; }
.consent-compare-item { text-align: center; }
.consent-compare-label { font-size: 14px; color: #6b7280; margin-bottom: 8px; }
.consent-compare-value { font-size: 28px; font-weight: 700; color: #67c23a; }
.consent-compare-detail { font-size: 12px; color: #909399; margin-top: 4px; }
.stat-card { text-align: center; padding: 10px; }
.stat-value { font-size: 28px; font-weight: 700; color: #1e3a8a; }
.stat-label { font-size: 13px; color: #6b7280; margin-top: 6px; }
</style>