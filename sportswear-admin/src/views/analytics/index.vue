<template>
  <div>
    <!-- 工具栏：时间范围 + 刷新 -->
    <div class="toolbar">
      <el-radio-group v-model="days" size="small" style="margin-right: 12px">
        <el-radio-button :value="7">最近 7 天</el-radio-button>
        <el-radio-button :value="30">最近 30 天</el-radio-button>
        <el-radio-button :value="90">最近 90 天</el-radio-button>
      </el-radio-group>
      <el-button type="primary" @click="loadCharts" :loading="chartsLoading" plain>刷新</el-button>
    </div>

    <!-- 概览卡片 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="6" v-for="c in cards" :key="c.label">
        <el-card shadow="hover"><div class="stat-card"><div class="stat-value">{{ c.value }}</div><div class="stat-label">{{ c.label }}</div></div></el-card>
      </el-col>
    </el-row>

    <!-- 流量趋势 + 渠道来源 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12"><el-card shadow="hover"><template #header>流量趋势（访次 / 独立 IP）</template><div ref="trendRef" style="height:280px"></div></el-card></el-col>
      <el-col :span="12"><el-card shadow="hover"><template #header>渠道来源分布</template><div ref="sourceRef" style="height:280px"></div></el-card></el-col>
    </el-row>

    <!-- 热门页面 + 热门产品 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12"><el-card shadow="hover"><template #header>热门页面（浏览量 TOP 10）</template><div ref="pagesRef" style="height:320px"></div></el-card></el-col>
      <el-col :span="12"><el-card shadow="hover"><template #header>热门产品（浏览量 TOP 10）</template><div ref="productsRef" style="height:320px"></div></el-card></el-col>
    </el-row>

    <!-- 国家分布 + 设备分布 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12"><el-card shadow="hover"><template #header>国家/地区分布</template><div ref="countryRef" style="height:320px"></div></el-card></el-col>
      <el-col :span="12"><el-card shadow="hover"><template #header>设备类型分布</template><div ref="deviceRef" style="height:320px"></div></el-card></el-col>
    </el-row>

    <!-- UTM 广告归因 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="24"><el-card shadow="hover"><template #header>广告归因（UTM Campaign 点击来源）</template>
        <el-table :data="utmCampaigns" v-loading="utmLoading" stripe style="width:100%">
          <el-table-column prop="utm_campaign" label="Campaign" min-width="160" />
          <el-table-column prop="source" label="来源" width="120" />
          <el-table-column prop="medium" label="Medium" width="100" />
          <el-table-column prop="visits" label="访次" width="100" />
          <el-table-column prop="unique_visitors" label="独立访客" width="110" />
        </el-table></el-card></el-col>
    </el-row>

    <!-- 社交媒体点击 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="24"><el-card shadow="hover"><template #header>社交媒体点击（外部链接跳转）</template>
        <el-table :data="socialClicks" v-loading="socialLoading" stripe style="width:100%">
          <el-table-column prop="platform" label="平台" min-width="140" />
          <el-table-column prop="clicks" label="点击量" width="100" />
          <el-table-column prop="unique_visitors" label="独立访客" width="110" />
        </el-table></el-card></el-col>
    </el-row>

    <!-- 访问日志明细 -->
    <el-card shadow="hover" class="mt-20">
      <template #header><span>门户访问日志明细</span></template>
      <div class="toolbar">
        <el-input v-model="vc.ip" placeholder="IP" clearable style="width:140px" @keyup.enter="searchVisits" />
        <el-input v-model="vc.country" placeholder="国家" clearable style="width:120px" @keyup.enter="searchVisits" />
        <el-input v-model="vc.source" placeholder="来源（如 google）" clearable style="width:160px" @keyup.enter="searchVisits" />
        <el-select v-model="vc.entityType" placeholder="实体类型" clearable style="width:130px" @change="searchVisits">
          <el-option label="product" value="product" />
          <el-option label="page" value="page" />
          <el-option label="blog" value="blog" />
          <el-option label="case" value="case" />
          <el-option label="home" value="home" />
          <el-option label="social_click" value="social_click" />
        </el-select>
        <el-select v-model="vc.device" placeholder="设备" clearable style="width:120px" @change="searchVisits">
          <el-option label="desktop" value="desktop" />
          <el-option label="mobile" value="mobile" />
          <el-option label="tablet" value="tablet" />
          <el-option label="bot" value="bot" />
        </el-select>
        <el-input v-model="vc.keyword" placeholder="关键词/路径/实体" clearable style="width:180px" @keyup.enter="searchVisits" />
        <el-date-picker v-model="vc.dateRange" type="daterange" value-format="YYYY-MM-DD" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" clearable style="width:260px" @change="searchVisits" />
        <el-button type="primary" @click="searchVisits">查询</el-button>
      </div>
      <el-table :data="visits" v-loading="visitsLoading" stripe :row-key="visitRowKey" style="width:100%">
        <el-table-column prop="created_at" label="时间" width="170"><template #default="{row}">{{ formatTime(row.created_at) }}</template></el-table-column>
        <el-table-column prop="country" label="国家" width="110" />
        <el-table-column prop="language" label="语言" width="90" />
        <el-table-column prop="device" label="设备" width="90" />
        <el-table-column prop="device_model" label="设备型号" width="130" />
        <el-table-column prop="os" label="操作系统" width="110" />
        <el-table-column prop="browser" label="浏览器" width="110" />
        <el-table-column prop="ip" label="IP" width="140" />
        <el-table-column prop="source" label="来源" width="110" />
        <el-table-column prop="medium" label="Medium" width="100" />
        <el-table-column prop="utm_campaign" label="广告Campaign" width="140" />
        <el-table-column prop="entity_type" label="实体类型" width="110" />
        <el-table-column prop="entity_name" label="产品/页面" min-width="160" />
        <el-table-column prop="path" label="路径" min-width="200" />
        <el-table-column prop="status" label="状态码" width="90" />
        <el-table-column prop="latency_ms" label="耗时(ms)" width="90" />
      </el-table>
      <el-pagination class="pagination" v-model:current-page="vpage" v-model:page-size="pageSize" :total="vtotal" layout="total, sizes, prev, pager, next, jumper" :page-sizes="[10, 20, 50, 100]" @current-change="searchVisits" @size-change="resetVisitPage" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { analyticsApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import * as echarts from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([BarChart, LineChart, PieChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const days = ref(7)
const chartsLoading = ref(false)
const pageSize = useAdminPageSize()

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
const formatTime = (t: string) => (!t ? '-' : new Date(t).toLocaleString())
const initChart = (el: any, option: any): echarts.ECharts => { const ch = echarts.init(el); ch.setOption(option); charts.push(ch); return ch }
const getChart = (el: any): echarts.ECharts | null => (el ? (echarts.getInstanceByDom(el) as echarts.ECharts) : null)

async function loadCharts() {
  chartsLoading.value = true
  try {
    const overview = await analyticsApi.overview({ days: days.value })
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

    const sources = await analyticsApi.sources({ days: days.value })
    const sr = getChart(sourceRef.value) || initChart(sourceRef.value, {})
    sr.setOption({
      tooltip: { trigger: 'item' },
      legend: { left: 'center', top: 'bottom', data: (sources || []).map((i: any) => i.source) },
      series: [{ type: 'pie', radius: ['35%', '70%'], label: { formatter: '{b}: {c} ({d}%)' }, data: (sources || []).map((i: any) => ({ name: i.source, value: i.visits || i.count })) }],
    })

    const topPages = await analyticsApi.topPages({ days: days.value })
    const pr = getChart(pagesRef.value) || initChart(pagesRef.value, {})
    pr.setOption({ tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } }, yAxis: { type: 'category', data: (topPages || []).map((i: any) => i.path) }, xAxis: { type: 'value' }, series: [{ type: 'bar', data: (topPages || []).map((i: any) => i.views) }] })

    const topProducts = await analyticsApi.topProducts({ days: days.value })
    const pc = getChart(productsRef.value) || initChart(productsRef.value, {})
    pc.setOption({ tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } }, yAxis: { type: 'category', data: (topProducts || []).map((i: any) => i.name) }, xAxis: { type: 'value' }, series: [{ type: 'bar', data: (topProducts || []).map((i: any) => i.views) }] })

    const countries = await analyticsApi.countries({ days: days.value })
    const cr = getChart(countryRef.value) || initChart(countryRef.value, {})
    cr.setOption({ tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } }, yAxis: { type: 'category', inverse: true, data: (countries || []).slice(0, 10).map((i: any) => i.country || i.region) }, xAxis: { type: 'value' }, series: [{ type: 'bar', data: (countries || []).slice(0, 10).map((i: any) => i.visits || i.count) }] })

    const devices = await analyticsApi.devices({ days: days.value })
    cards.value[4].value = devices?.mobile_ratio ? ((Number(devices.mobile_ratio) * 100).toFixed(0) + '%') : '0%'
    const dr = getChart(deviceRef.value) || initChart(deviceRef.value, {})
    dr.setOption({ tooltip: { trigger: 'item' }, legend: { left: 'center', top: 'bottom', data: (devices?.devices || []).map((i: any) => i.device) }, series: [{ type: 'pie', radius: ['35%', '70%'], label: { formatter: '{b}: {c} ({d}%)' }, data: (devices?.devices || []).map((i: any) => ({ name: i.device, value: i.visits })) }] })

    utmCampaigns.value = (sources || []).map((i: any) => ({ utm_campaign: i.utm_campaign || i.campaign || '-', source: i.source, medium: i.medium, visits: i.visits || i.count, unique_visitors: i.unique_visitors || '' }))

    const social = await analyticsApi.socialClicks({ days: days.value })
    socialClicks.value = social || []
  } catch { /* ignore */ } finally { chartsLoading.value = false }
}

const utmCampaigns = ref<any[]>([])
const utmLoading = ref(false)
const socialClicks = ref<any[]>([])
const socialLoading = ref(false)

// Visit logs
const vc = ref({ ip: '', country: '', source: '', entityType: '', device: '', keyword: '', dateRange: null as null | [string, string] })
const visits = ref<any[]>([])
const visitsLoading = ref(false)
const vtotal = ref(0)
const vpage = ref(1)

function resetVisitPage() { vpage.value = 1; searchVisits() }

async function searchVisits() {
  vpage.value = 1
  visitsLoading.value = true
  try {
    const params: any = { page: vpage.value, pageSize: pageSize.value, ip: vc.value.ip, country: vc.value.country, source: vc.value.source, entity_type: vc.value.entityType, device: vc.value.device, keyword: vc.value.keyword }
    if (vc.value.dateRange) { params.start_date = vc.value.dateRange[0]; params.end_date = vc.value.dateRange[1] }
    const res = await analyticsApi.visitLogs(params)
    visits.value = res.items
    vtotal.value = res.total
  } finally { visitsLoading.value = false }
}

function visitRowKey(row: any) { return `${row.src || 'visit_logs'}-${row.id}-${row.created_at}` }

onMounted(() => { loadCharts(); searchVisits() })
watch(days, () => { loadCharts() })
onUnmounted(() => charts.forEach((c) => c.dispose()))
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 12px; align-items: center; }
.mt-20 { margin-top: 20px; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.stat-card { text-align: center; padding: 10px; }
.stat-value { font-size: 28px; font-weight: 700; color: #1e3a8a; }
.stat-label { font-size: 13px; color: #6b7280; margin-top: 6px; }
</style>
