<template>
  <div class="product-hub">
    <el-tabs v-model="activeTab" class="hub-tabs" @tab-change="onTabChange">
      <el-tab-pane
        v-for="t in visibleTabs"
        :key="t.name"
        :name="t.name"
        :label="t.label"
        lazy
      >
        <component :is="t.component" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import ProductList from './list.vue'
import ProductCategories from './categories.vue'
import ProductSeries from './series.vue'
import ProductFabrics from './fabrics.vue'

interface ProductTab {
  name: string
  label: string
  permission: string
  component: any
}

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 产品管理四个维度：统一入口，标签页承载，按权限过滤显示。
// 拥有任一权限即可进入 /products（路由守卫 any-of 语义），标签按各自权限显隐。
const ALL_TABS: ProductTab[] = [
  { name: 'list', label: '产品列表', permission: 'product:view', component: ProductList },
  { name: 'categories', label: '分类管理', permission: 'category:manage', component: ProductCategories },
  { name: 'series', label: '系列管理', permission: 'series:manage', component: ProductSeries },
  { name: 'fabrics', label: '面料管理', permission: 'fabric:manage', component: ProductFabrics },
]

const visibleTabs = computed(() => ALL_TABS.filter((t) => authStore.hasPermission(t.permission)))

const activeTab = ref<string>(visibleTabs.value[0]?.name || 'list')

// 权限变化（如登录态恢复）时，确保当前标签始终可用
watch(visibleTabs, (tabs) => {
  if (!tabs.some((t) => t.name === activeTab.value)) {
    activeTab.value = tabs[0]?.name || 'list'
  }
})

function onTabChange(name: string | number) {
  router.replace({ query: { ...route.query, tab: String(name) } })
}

// 初始化：从 query.tab 恢复上次所在标签（仅当当前用户有权限）
const q = route.query.tab
if (typeof q === 'string' && visibleTabs.value.some((t) => t.name === q)) {
  activeTab.value = q
}
</script>

<style scoped>
.product-hub {
  padding: 0;
}
.hub-tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
  background: #fff;
  border: 1px solid #eef0f4;
  border-radius: 10px;
  padding: 0 12px;
}
.hub-tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}
.hub-tabs :deep(.el-tabs__item) {
  font-size: 14px;
  height: 48px;
  line-height: 48px;
}
.hub-tabs :deep(.el-tabs__item.is-active) {
  font-weight: 600;
  color: #4f46e5;
}
.hub-tabs :deep(.el-tabs__active-bar) {
  background: #4f46e5;
}
</style>
