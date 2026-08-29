<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="6" v-for="c in cards" :key="c.label">
        <el-card shadow="hover"><div class="stat-card"><div class="stat-value">{{ c.value }}</div><div class="stat-label">{{ c.label }}</div></div></el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <el-card shadow="hover"><template #header>流量趋势</template><div ref="trendRef" style="height:280px"></div></el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover"><template #header>来源分布</template><div ref="sourceRef" style="height:280px"></div></el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
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
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { analyticsApi } from '@/api'
// echarts 按需引入：仅打包用到的图表与组件，显著减小包体积
import * as echarts from 'echarts/core'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([BarChart, LineChart, PieChart, GridComponent, TooltipComponent, CanvasRenderer])

const cards = ref([
  { label: '总访问量', value: 0 }, { label: '独立 IP', value: 0 }, { label: '询盘转化', value: '0%' }, { label: '设备-移动端', value: '0%' },
])
const trendRef = ref<HTMLDivElement>(); const sourceRef = ref<HTMLDivElement>(); const pagesRef = ref<HTMLDivElement>(); const productsRef = ref<HTMLDivElement>(); const countryRef = ref<HTMLDivElement>()
const charts: echarts.ECharts[] = []

function initChart(el: HTMLDivElement, option: any) {
  const chart = echarts.init(el)
  chart.setOption(option)
  charts.push(chart)
}

onMounted(async () => {
  try {
    const overview = await analyticsApi.overview()
    if (overview) {
      cards.value[0].value = overview.total_visits || 0
      cards.value[1].value = overview.unique_ips || 0
      cards.value[2].value = (overview.conversion_rate * 100).toFixed(1) + '%' || '0%'
    }
    if (trendRef.value) initChart(trendRef.value, {
      xAxis: { type: 'category', data: [] }, yAxis: { type: 'value' },
      tooltip: { trigger: 'axis' }, series: [{ type: 'line', smooth: true, data: overview?.daily_trend || [] }]
    })
    if (sourceRef.value) initChart(sourceRef.value, {
      tooltip: { trigger: 'item' }, series: [{ type: 'pie', radius: ['40%', '70%'], data: [] }]
    })
    const topPages = await analyticsApi.topPages()
    if (pagesRef.value) initChart(pagesRef.value, { xAxis: { type: 'value' }, yAxis: { type: 'category', data: (topPages || []).map((i: any) => i.path) }, series: [{ type: 'bar', data: (topPages || []).map((i: any) => i.count) }] })
    const topProducts = await analyticsApi.topProducts()
    if (productsRef.value) initChart(productsRef.value, { xAxis: { type: 'value' }, yAxis: { type: 'category', data: (topProducts || []).map((i: any) => i.name) }, series: [{ type: 'bar', data: (topProducts || []).map((i: any) => i.count) }] })
    const countries = await analyticsApi.countries()
    // echarts 5 不再内置世界地图数据（需额外引入 ~200KB GeoJSON），改为横向柱状图展示国家分布 Top10
    if (countryRef.value) initChart(countryRef.value, {
      xAxis: { type: 'value' },
      yAxis: { type: 'category', inverse: true, data: (countries || []).slice(0, 10).map((i: any) => i.country) },
      tooltip: { trigger: 'axis' },
      series: [{ type: 'bar', data: (countries || []).slice(0, 10).map((i: any) => i.count) }],
    })
    const devices = await analyticsApi.devices()
    cards.value[3].value = (devices?.mobile_ratio * 100).toFixed(0) + '%' || '0%'
    const sources = await analyticsApi.sources()
    if (sourceRef.value) { const chart = echarts.getInstanceByDom(sourceRef.value); if (chart) chart.setOption({ series: [{ type: 'pie', radius: ['40%', '70%'], data: (sources || []).map((i: any) => ({ name: i.source, value: i.count })) }] }) }
  } catch {}
})

onUnmounted(() => charts.forEach(c => c.dispose()))
</script>
<style scoped>
.mt-20 { margin-top: 20px; }
.stat-card { text-align: center; padding: 10px; }
.stat-value { font-size: 28px; font-weight: 700; color: #1e3a8a; }
.stat-label { font-size: 13px; color: #6b7280; margin-top: 6px; }
</style>