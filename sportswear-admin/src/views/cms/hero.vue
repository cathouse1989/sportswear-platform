<template>
  <div class="hero-page" v-loading="loading">
    <el-alert
      v-if="!loading && !pageId"
      type="warning"
      :closable="false"
      show-icon
      title="未找到首页页面"
      description="请先在「页面管理」中创建并发布 slug 为 home（或类型为 home）的页面，再配置轮播图。"
    />
    <HeroEditor v-else-if="pageId" :page-id="pageId" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { pageApi } from '@/api'
import type { Page } from '@/types'
import HeroEditor from '@/components/cms/HeroEditor.vue'

const loading = ref(false)
const pageId = ref('')
const isHome = (row: Page) => row.slug === 'home' || row.type === 'home'

async function findHomePage(): Promise<Page | null> {
  const tryList = async (status?: string) => {
    const result = await pageApi.list({ page: 1, pageSize: 100, status })
    const items: Page[] = (result as any).items || result || []
    return items.find(isHome) || null
  }
  return (await tryList('published')) || (await tryList()) || null
}

onMounted(async () => {
  loading.value = true
  try {
    const home = await findHomePage()
    pageId.value = home?.id || ''
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.hero-page { min-height: 360px; }
</style>
