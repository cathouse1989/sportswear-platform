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
        <el-table :data="socialClicks" v-loading="socialLoading" stripe style="width:100%">
          <el-table-column prop="platform" label="平台" min-width="140" />
          <el-table-column prop="clicks" label="点击量" width="100" />
          <el-table-column prop="unique_visitors" label="独立访客" width="110" />
        </el-table></el-card></el-col>
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

echarts.use([BarChart, LineChart, PieChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const days = ref(7)
const chartsLoading = ref(false)

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

    const utm = await analyticsApi.utmCampaigns({ days: days.value })
    utmCampaigns.value = utm || []

    const social = await analyticsApi.socialClicks({ days: days.value })
    socialClicks.value = social || []
  } catch { /* ignore */ } finally { chartsLoading.value = false }
}

const utmCampaigns = ref<any[]>([])
const utmLoading = ref(false)
const socialClicks = ref<any[]>([])
const socialLoading = ref(false)

onMounted(() => { loadCharts() })
watch(days, () => { loadCharts() })
onUnmounted(() => charts.forEach((c) => c.dispose()))
</script>

<style scoped>
.toolbar { display: flex; flex-wrap: wrap; gap: 12px; align-items: center; }
.mt-20 { margin-top: 20px; }
.stat-card { text-align: center; padding: 10px; }
.stat-value { font-size: 28px; font-weight: 700; color: #1e3a8a; }
.stat-label { font-size: 13px; color: #6b7280; margin-top: 6px; }
</style>