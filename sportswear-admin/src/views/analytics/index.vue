<template>
  <div>
    <!-- 工具栏：时间范围 + 概览/明细切换 -->
    <div class="toolbar">
      <el-radio-group v-model="days" size="small" style="margin-right: 12px">
        <el-radio-button :value="7">最近 7 天</el-radio-button>
        <el-radio-button :value="30">最近 30 天</el-radio-button>
        <el-radio-button :value="90">最近 90 天</el-radio-button>
      </el-radio-group>
      <el-tabs v-model="tab" class="mt-16">
        <el-tab-pane label="流量概览" name="charts" />
        <el-tab-pane label="访问日志明细" name="list" />
      </el-tabs>
    </div>

    <!-- 总览卡片 -->
    <el-row :gutter="20" v-if="tab === 'charts'">
      <el-col :span="6" v-for="c in cards" :key="c.label">
        <el-card shadow="hover"><div class="stat-card"><div class="stat-value">{{ c.value }}</div><div class="stat-label">{{ c.label }}</div></div></el-card>
      </el-col>
    </el-row>

    <!-- 图表 -->
    <el-row :gutter="20" class="mt-20" v-if="tab === 'charts'">
      <el-col :span="12">
        <el-card shadow="hover"><template #header>流量趋势</template><div ref="trendRef" style="height:280px"></div></el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover"><template #header>来源分布</template><div ref="sourceRef" style="height:280px"></div></el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" class="mt-20" v-if="tab === 'charts'">
      <el-col :span="8">
        <el-card shadow="hover"><template #header>热门页面</template><div ref="pagesRef" style="height:240px"></div></el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover"><template #header>热门产品</template><div ref="productsRef" style="height:240px"></div></el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover"><template #header>国家分布</template><div ref="countryRef" style="height:240px"></div></el-card>
      </el-col>
    </el-row>

    <!-- 访问日志明细 -->
    <el-card shadow="hover" v-if="tab === 'list'" class="mt-20">
      <div class="toolbar">
        <el-input v-model="vc.ip" placeholder="IP" clearable style="width:140px" @keyup.enter="searchVisits" />
        <el-input v-model="vc.country" placeholder="国家" clearable style="width:120px" @keyup.enter="searchVisits" />
        <el-input v-model="vc.source" placeholder="来源" clearable style="width:150px" @keyup.enter="searchVisits" />
        <el-input v-model="vc.keyword" placeholder="关键词/路径" clearable style="width:180px" @keyup.enter="searchVisits" />
        <el-date-picker v-model="vc.dateRange" type="daterange" value-format="YYYY-MM-DD" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" clearable style="width:260px" @change="searchVisits" />
        <el-button type="primary" @click="searchVisits">查询</el-button>
      </div>
      <el-table :data="visits" v-loading="visitsLoading" stripe :row-key="(r:any)=>'v'+r.src+'-'+r.id">
        <el-table-column prop="created_at" label="时间" width="160"><template #default="{row}">{{ formatTime(row.created_at) }}</template></el-table-column>
        <el-table-column prop="country" label="国家" width="100" />
        <el-table-column prop="language" label="语言" width="80" />
        <el-table-column prop="device" label="设备" width="80" />
        <el-table-column prop="device_model" label="设备型号" width="130" />
        <el-table-column prop="browser" label="浏览器" width="100" />
        <el-table-column prop="os" label="操作系统" width="110" />
        <el-table-column prop="source" label="来源" width="100" />
        <el-table-column prop="medium" label="Medium" width="100" />
        <el-table-column prop="utm_campaign" label="广告Campaign" width="130" />
        <el-table-column prop="entity_name" label="产品/页面" min-width="140" />
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="path" label="路径" min-width="200" />
        <el-table-column prop="status" label="状态码" width="90" />
        <el-table-column prop="latency_ms" label="耗时(ms)" width="90" />
      </el-table>
      <el-pagination class="pagination" v-model:current-page="vpage" v-model:page-size="pageSize" :total="vtotal" layout="total, sizes, prev, pager, next" :page-sizes="[10, 20, 50, 100]" @current-change="searchVisits" @size-change="searchVisits" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import { analyticsApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import * as echarts from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([BarChart, LineChart, PieChart, GridComponent, TooltipComponent, CanvasRenderer])

const days = ref(7)
const tab = ref<'charts' | 'list'>('charts')
const cards = ref([
  { label: '总访问量', value: 0 },
  { label: '独立 IP', value: 0 },
  { label: '询盘转化', value: '0%' },
  { label: '设备-移动端', value: '0%' },
])

const trendRef = ref(); const sourceRef = ref(); const pagesRef = ref(); const productsRef = ref(); const countryRef = ref()
const charts: echarts.ECharts[] = []
const pageSize = useAdminPageSize()

function fmt(v: any, def = 0) { return v ? v : def }
function formatTime(t: string) { return !t ? '-' : new Date(t).toLocaleString() }

function initChart(el: any, option: any): echarts.ECharts {
  const chart = echarts.init(el)
  chart.setOption(option)
  charts.push(chart)
  return chart
}

async function loadCharts() {
  try {
    const overview = await analyticsApi.overview({ days: days.value })
    cards.value[0].value = fmt(overview?.total_visits)
    cards.value[1].value = fmt(overview?.unique_ips)
    const trend = (overview?.daily_trend || []) as any[]
    if (trendRef.value) {
      const c = (echarts.getInstanceByDom(trendRef.value) as echarts.ECharts) || initChart(trendRef.value, {})
      c.setOption({
        xAxis: { type: 'category', data: trend.map((d: any) => d.date) },
        yAxis: { type: 'value' }, tooltip: { trigger: 'axis' },
        series: [{ type: 'line', smooth: true, data: trend.map((d: any) => d.visits) }],
      })
    }
    if (sourceRef.value) initChart(sourceRef.value, {
      tooltip: { trigger: 'item' },
      series: [{ type: 'pie', radius: ['40%', '70%'], data: [] }],
    })
    const topPages = await analyticsApi.topPages({ days: days.value })
    if (pagesRef.value) initChart(pagesRef.value, {
      xAxis: { type: 'value' }, yAxis: { type: 'category', data: (topPages || []).map((i: any) => i.path) },
      tooltip: { trigger: 'axis' }, series: [{ type: 'bar', data: (topPages || []).map((i: any) => i.views) }],
    })
    const topProducts = await analyticsApi.topProducts({ days: days.value })
    if (productsRef.value) initChart(productsRef.value, {
      xAxis: { type: 'value' }, yAxis: { type: 'category', data: (topProducts || []).map((i: any) => i.name) },
      tooltip: { trigger: 'axis' }, series: [{ type: 'bar', data: (topProducts || []).map((i: any) => i.views) }],
    })
    const countries = await analyticsApi.countries({ days: days.value })
    if (countryRef.value) initChart(countryRef.value, {
      xAxis: { type: 'value' },
      yAxis: { type: 'category', inverse: true, data: (countries || []).slice(0, 10).map((i: any) => i.country || i.region) },
      tooltip: { trigger: 'axis' },
      series: [{ type: 'bar', data: (countries || []).slice(0, 10).map((i: any) => i.visits || i.count) }],
    })
    const devices = await analyticsApi.devices({ days: days.value })
    cards.value[3].value = devices?.mobile_ratio ? ((devices.mobile_ratio * 100).toFixed(0) + '%') : '0%'
    const sources = await analyticsApi.sources({ days: days.value })
    if (sourceRef.value) {
      const c = echarts.getInstanceByDom(sourceRef.value) as echarts.ECharts
      if (c) c.setOption({
        series: [{ type: 'pie', radius: ['40%', '70%'], data: (sources || []).map((i: any) => ({ name: i.source, value: i.visits || i.count })) }],
      })
    }
  } catch { /* 忽略 */ }
}

const vc = ref({ ip: '', country: '', source: '', keyword: '', dateRange: null as null | [string, string] })
const visits = ref<any[]>([])
const visitsLoading = ref(false)
const vtotal = ref(0)
const vpage = ref(1)

async function searchVisits() {
  vpage.value = 1
  visitsLoading.value = true
  try {
    const params: any = { page: vpage.value, pageSize: pageSize.value, ip: vc.value.ip, country: vc.value.country, source: vc.value.source, keyword: vc.value.keyword }
    if (vc.value.dateRange) { params.start_date = vc.value.dateRange[0]; params.end_date = vc.value.dateRange[1] }
    const res = await analyticsApi.visitLogs(params)
    visits.value = res.items
    vtotal.value = res.total
  } finally {
    visitsLoading.value = false
  }
}

onMounted(() => { if (tab.value === 'charts') loadCharts(); else searchVisits() })
watch(days, () => { if (tab.value === 'charts') loadCharts() })
watch(tab, () => { if (tab.value === 'list') { vpage.value = 1; searchVisits() } })

onUnmounted(() => charts.forEach((c) => c.dispose()))
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 12px; align-items: center; }
.mt-16 { margin-top: 16px; }
.mt-20 { margin-top: 20px; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.stat-card { text-align: center; padding: 10px; }
.stat-value { font-size: 28px; font-weight: 700; color: #1e3a8a; }
.stat-label { font-size: 13px; color: #6b7280; margin-top: 6px; }
</style>