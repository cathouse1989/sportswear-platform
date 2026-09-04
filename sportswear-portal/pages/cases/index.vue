<template>
  <div class="max-w-7xl mx-auto px-4 py-12">
    <h1 class="text-2xl md:text-3xl font-bold mb-6 md:mb-8">{{ $t('nav.cases') }}</h1>

    <!-- 骨架屏 -->
    <div v-if="loading" class="grid sm:grid-cols-2 md:grid-cols-3 gap-4 md:gap-6">
      <div v-for="i in 6" :key="i" class="animate-pulse bg-white rounded-xl overflow-hidden p-4 md:p-6">
        <div class="h-4 bg-gray-100 rounded w-3/4 mb-2" />
        <div class="h-3 bg-gray-100 rounded w-1/2" />
      </div>
    </div>

    <!-- 案例卡片 -->
    <div v-show="!loading" class="grid sm:grid-cols-2 md:grid-cols-3 gap-4 md:gap-6">
      <NuxtLink v-for="c in cases" :key="c.id" :to="localePath(`/cases/${c.slug}`)" class="bg-white rounded-xl overflow-hidden shadow hover:shadow-lg transition p-4 md:p-6 active:scale-[0.98]">
        <h3 class="text-sm md:text-base font-semibold mb-2 line-clamp-2">{{ c.title }}</h3>
        <p class="text-xs md:text-sm text-gray-500">{{ c.client_industry }}</p>
      </NuxtLink>
    </div>

    <!-- 空态 -->
    <div v-if="!loading && !cases?.length" class="text-center py-16 text-gray-400">No cases yet.</div>

    <!-- 加载更多 -->
    <div v-if="!loading && cases.length" class="text-center mt-10 md:mt-12">
      <button
        v-if="hasMore"
        @click="loadMore"
        :disabled="loadingMore"
        class="inline-flex items-center gap-2 px-8 py-3.5 border-2 border-gray-200 rounded-full text-sm font-medium text-gray-700 hover:border-black hover:text-black transition disabled:opacity-50 min-h-[48px]"
      >
        <span v-if="loadingMore" class="inline-block w-4 h-4 border-2 border-gray-400 border-t-transparent rounded-full animate-spin" />
        {{ loadingMore ? $t('common.loading') : $t('product.load_more') }}
      </button>
      <p v-else class="text-sm text-gray-400">{{ $t('product.no_products') || 'No more cases' }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
const localePath = useLocalePath()
const api = useApi()
const { locale } = useI18n()
const route = useRoute()
// SEO - 案例列表
useSeoHead({
  title: 'Case Studies - OEM/ODM Sportswear Manufacturing',
  description:
    'Explore our successful OEM/ODM sportswear projects. See how we help global brands bring their activewear visions to life with quality manufacturing and innovative solutions.',
  keywords:
    'sportswear case studies, OEM projects, ODM manufacturing cases, activewear brand partnerships, garment factory portfolio',
})

const cases = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = ref(Number(route.query._pageSize) || 20)
const hasMore = ref(true)
const total = ref(0)

// SSR 首屏：内联到 HTML 以提升 SEO
const { data: initialData } = await useAsyncData<any>(
  'cases-' + (locale.value || 'en'),
  () => api.getCases({ page: 1, pageSize: pageSize.value })
    .then((res: any) => ({
      items: res?.items || (Array.isArray(res) ? res : []),
      total: typeof res?.total === 'number' ? res.total : (Array.isArray(res) ? res.length : 0),
    }))
    .catch(() => ({ items: [], total: 0 })),
)

if (initialData.value) {
  cases.value = initialData.value.items || []
  total.value = initialData.value.total || 0
  hasMore.value = cases.value.length > 0 && cases.value.length < total.value
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  const nextPage = page.value + 1
  try {
    const res = await api.getCases({ page: nextPage, pageSize: pageSize.value })
    const newItems = Array.isArray(res) ? res : (res?.items || [])
    if (newItems.length > 0) {
      const existingIds = new Set(cases.value.map((c: any) => c.id))
      const unique = newItems.filter((c: any) => !existingIds.has(c.id))
      if (unique.length > 0) {
        cases.value.push(...unique)
        page.value = nextPage
      }
    }
    const latestTotal = typeof res?.total === 'number' ? res.total : total.value
    total.value = latestTotal
    hasMore.value = cases.value.length > 0 && cases.value.length < latestTotal
  } catch (e) {
    console.error('loadMore error:', e)
    hasMore.value = true
  } finally { loadingMore.value = false }
}
</script>