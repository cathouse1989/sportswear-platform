<template>
  <div class="max-w-3xl mx-auto px-4 py-12 md:py-16">
    <!-- 页头 -->
    <header class="mb-10 md:mb-14">
      <h1 class="text-3xl md:text-4xl font-bold tracking-tight">{{ $t('blog.title') }}</h1>
      <p class="mt-3 text-gray-500 text-base leading-relaxed">{{ $t('blog.subtitle') }}</p>
    </header>

    <!-- 骨架屏 -->
    <div v-if="loading" class="divide-y divide-gray-100">
      <div v-for="i in 4" :key="i" class="animate-pulse flex flex-col sm:flex-row gap-6 py-8 md:py-10">
        <div class="w-full sm:w-52 md:w-60 h-40 sm:h-32 bg-gray-100 rounded-lg shrink-0" />
        <div class="flex-1 space-y-3">
          <div class="h-3 bg-gray-100 rounded w-24" />
          <div class="h-5 bg-gray-100 rounded w-3/4" />
          <div class="h-4 bg-gray-100 rounded w-full" />
          <div class="h-4 bg-gray-100 rounded w-5/6" />
        </div>
      </div>
    </div>

    <!-- 博客列表：一行一条 -->
    <div v-show="!loading" class="divide-y divide-gray-100">
      <NuxtLink
        v-for="b in blogs"
        :key="b.id"
        :to="localePath(`/blog/${b.slug}`)"
        class="group flex flex-col sm:flex-row gap-6 py-8 md:py-10 -mx-4 px-4 rounded-xl hover:bg-gray-50/70 transition-colors"
      >
        <div class="w-full sm:w-52 md:w-60 shrink-0">
          <img v-if="b.cover_image" :src="b.cover_image" :alt="b.title" loading="lazy" class="w-full h-40 sm:h-32 object-cover rounded-lg bg-gray-100" />
          <div v-else class="w-full h-40 sm:h-32 rounded-lg bg-gray-100 flex items-center justify-center text-4xl">📝</div>
        </div>
        <div class="flex-1 flex flex-col justify-center min-w-0">
          <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-500 mb-3">
            <span class="px-2.5 py-1 rounded-full bg-blue-50 text-blue-700 font-medium">{{ categoryLabel(b.category) }}</span>
            <time>{{ formatDate(b.published_at || b.created_at) }}</time>
          </div>
          <h2 class="text-lg md:text-xl font-semibold leading-snug text-gray-900 group-hover:text-blue-700 transition-colors line-clamp-2">{{ b.title }}</h2>
          <p class="mt-2.5 text-sm md:text-base text-gray-500 leading-relaxed line-clamp-2">{{ stripHtml(b.content) }}</p>
          <span class="mt-4 inline-flex items-center gap-1.5 text-sm font-medium text-blue-700">
            {{ $t('blog.read_more') }}
            <span aria-hidden="true">→</span>
          </span>
        </div>
      </NuxtLink>
    </div>

    <!-- 空态 -->
    <div v-if="!loading && !blogs?.length" class="text-center py-20 text-gray-400">{{ $t('blog.empty') }}</div>

    <!-- 加载更多 -->
    <div v-if="!loading && blogs.length" class="text-center mt-12 md:mt-16">
      <button
        v-if="hasMore"
        @click="loadMore"
        :disabled="loadingMore"
        class="inline-flex items-center gap-2 px-8 py-3.5 border-2 border-gray-200 rounded-full text-sm font-medium text-gray-700 hover:border-black hover:text-black transition disabled:opacity-50 min-h-[48px]"
      >
        <span v-if="loadingMore" class="inline-block w-4 h-4 border-2 border-gray-400 border-t-transparent rounded-full animate-spin" />
        {{ loadingMore ? $t('common.loading') : $t('product.load_more') }}
      </button>
      <p v-else class="text-sm text-gray-400">{{ $t('blog.all_loaded') }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
const localePath = useLocalePath()
const api = useApi()
const { locale, t, te } = useI18n()
const route = useRoute()

// SEO - 博客列表页
useSeoHead({
  title: 'Sportswear Manufacturing Blog - Industry Insights & Tips',
  description:
    'Expert insights on sportswear manufacturing, OEM/ODM processes, fabric technology, and industry trends. Learn how to choose the right manufacturer for your activewear brand.',
  keywords:
    'sportswear blog, OEM manufacturing insights, activewear industry, garment factory, sportswear manufacturing tips',
})

const blogs = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = ref(Number(route.query._pageSize) || 20)
const hasMore = ref(true)
const total = ref(0)

// 本地化分类标签：优先词条，回退原始 code
function categoryLabel(value?: string) {
  if (!value) return ''
  const key = `blog.categories.${value}`
  return te(key) ? t(key) : value
}

// 本地化日期
function formatDate(value?: string) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  try {
    return new Intl.DateTimeFormat(locale.value || 'en', { year: 'numeric', month: 'long', day: 'numeric' }).format(d)
  } catch {
    return d.toISOString().slice(0, 10)
  }
}

// 去除 HTML 标签生成摘要
function stripHtml(html?: string) {
  return (html || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').replace(/\s+/g, ' ').trim()
}

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