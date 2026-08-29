<template>
  <div class="max-w-7xl mx-auto px-4 py-12">
    <h1 class="text-2xl md:text-3xl font-bold mb-6 md:mb-8">{{ $t('nav.blog') }}</h1>

    <!-- 骨架屏 -->
    <div v-if="loading" class="grid sm:grid-cols-2 md:grid-cols-3 gap-4 md:gap-6">
      <div v-for="i in 6" :key="i" class="animate-pulse bg-white rounded-xl overflow-hidden">
        <div class="h-36 md:h-44 bg-gray-100" />
        <div class="p-3 md:p-4">
          <div class="h-3 bg-gray-100 rounded w-1/3 mb-2" />
          <div class="h-4 bg-gray-100 rounded w-3/4 mb-2" />
          <div class="h-3 bg-gray-100 rounded w-2/3" />
        </div>
      </div>
    </div>

    <!-- 博客卡片 -->
    <div v-show="!loading" class="grid sm:grid-cols-2 md:grid-cols-3 gap-4 md:gap-6">
      <NuxtLink v-for="b in blogs" :key="b.id" :to="localePath(`/blog/${b.slug}`)" class="bg-white rounded-xl overflow-hidden shadow hover:shadow-lg transition active:scale-[0.98]">
        <div class="h-36 md:h-44 bg-gray-100 flex items-center justify-center text-3xl md:text-4xl">📝</div>
        <div class="p-3 md:p-4">
          <p class="text-[10px] md:text-xs text-blue-600 mb-1">{{ b.category }}</p>
          <h3 class="text-sm md:text-base font-semibold mb-1 md:mb-2 line-clamp-2">{{ b.title }}</h3>
          <p class="text-xs md:text-sm text-gray-500 line-clamp-3">{{ b.content?.slice(0, 100) }}</p>
        </div>
      </NuxtLink>
    </div>

    <!-- 空态 -->
    <div v-if="!loading && !blogs?.length" class="text-center py-16 text-gray-400">{{ $t('common.no_data') || 'No articles yet.' }}</div>

    <!-- 加载更多 -->
    <div v-if="!loading && blogs.length" class="text-center mt-10 md:mt-12">
      <button
        v-if="hasMore"
        @click="loadMore"
        :disabled="loadingMore"
        class="inline-flex items-center gap-2 px-8 py-3.5 border-2 border-gray-200 rounded-full text-sm font-medium text-gray-700 hover:border-black hover:text-black transition disabled:opacity-50 min-h-[48px]"
      >
        <span v-if="loadingMore" class="inline-block w-4 h-4 border-2 border-gray-400 border-t-transparent rounded-full animate-spin" />
        {{ loadingMore ? $t('common.loading') : $t('product.load_more') }}
      </button>
      <p v-else class="text-sm text-gray-400">{{ $t('product.no_products') || 'No more articles' }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
const localePath = useLocalePath()
const api = useApi()
const { locale } = useI18n()
const route = useRoute()

useHead({ title: 'Blog - Sportswear OEM/ODM' })

const blogs = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = ref(Number(route.query._pageSize) || 20)
const hasMore = ref(true)
const total = ref(0)

// SSR 首屏：内联到 HTML 以提升 SEO
const { data: initialData } = await useAsyncData<any>(
  'blogs-' + (locale.value || 'en'),
  () => api.getBlogs({ page: 1, pageSize: pageSize.value })
    .then((res: any) => ({
      items: res?.items || (Array.isArray(res) ? res : []),
      total: typeof res?.total === 'number' ? res.total : (Array.isArray(res) ? res.length : 0),
    }))
    .catch(() => ({ items: [], total: 0 })),
)

if (initialData.value) {
  blogs.value = initialData.value.items || []
  total.value = initialData.value.total || 0
  hasMore.value = blogs.value.length > 0 && blogs.value.length < total.value
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  const nextPage = page.value + 1
  try {
    const res = await api.getBlogs({ page: nextPage, pageSize: pageSize.value })
    const newItems = Array.isArray(res) ? res : (res?.items || [])
    if (newItems.length > 0) {
      const existingIds = new Set(blogs.value.map((b: any) => b.id))
      const unique = newItems.filter((b: any) => !existingIds.has(b.id))
      if (unique.length > 0) {
        blogs.value.push(...unique)
        page.value = nextPage
      }
    }
    const latestTotal = typeof res?.total === 'number' ? res.total : total.value
    total.value = latestTotal
    hasMore.value = blogs.value.length > 0 && blogs.value.length < latestTotal
  } catch (e) {
    console.error('loadMore error:', e)
    hasMore.value = true
  } finally { loadingMore.value = false }
}
</script>