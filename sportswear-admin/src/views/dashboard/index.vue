<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="6" v-for="card in statCards" :key="card.label">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-value">{{ card.value }}</div>
            <div class="stat-label">{{ card.label }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>热门产品 Top5</template>
          <div v-for="(item, index) in hotProducts" :key="item.name" class="rank-item">
            <span class="rank">{{ index + 1 }}</span>
            <span class="name">{{ item.name }}</span>
            <span class="count">{{ item.count }} 条询盘</span>
          </div>
          <el-empty v-if="hotProducts.length === 0" description="暂无数据" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>热门国家 Top5</template>
          <div v-for="(item, index) in hotCountries" :key="item.name" class="rank-item">
            <span class="rank">{{ index + 1 }}</span>
            <span class="name">{{ item.name }}</span>
            <span class="count">{{ item.count }} 条询盘</span>
          </div>
          <el-empty v-if="hotCountries.length === 0" description="暂无数据" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { dashboardApi } from '@/api'
import type { DashboardStats } from '@/types'

const stats = ref<DashboardStats | null>(null)

const statCards = computed(() => [
  { label: '今日询盘', value: stats.value?.today_leads ?? 0 },
  { label: '本周询盘', value: stats.value?.week_leads ?? 0 },
  { label: '本月询盘', value: stats.value?.month_leads ?? 0 },
  { label: '高价值询盘', value: stats.value?.high_value_leads ?? 0 },
])

const hotProducts = computed(() => stats.value?.hot_products ?? [])
const hotCountries = computed(() => stats.value?.hot_countries ?? [])

onMounted(async () => {
  try {
    stats.value = await dashboardApi.stats()
  } catch {
    // 忽略加载失败
  }
})
</script>

<style scoped>
.mt-20 {
  margin-top: 20px;
}
.stat-card {
  text-align: center;
  padding: 10px;
}
.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #1e3a8a;
}
.stat-label {
  font-size: 13px;
  color: #6b7280;
  margin-top: 8px;
}
.rank-item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f3f4f6;
}
.rank {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #eef2ff;
  color: #1e3a8a;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  margin-right: 12px;
}
.name {
  flex: 1;
}
.count {
  color: #6b7280;
  font-size: 13px;
}
</style>