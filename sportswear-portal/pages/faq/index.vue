<template>
  <div class="max-w-4xl mx-auto px-4 py-12">
    <h1 class="text-2xl md:text-3xl font-bold mb-6 md:mb-8">{{ $t('nav.faq') }}</h1>

    <!-- 骨架屏 -->
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 5" :key="i" class="animate-pulse border-b border-[#EAE5DD] py-4 md:py-5">
        <div class="h-4 bg-gray-100 rounded w-3/4" />
      </div>
    </div>

    <!-- FAQ 列表 -->
    <div v-show="!loading" v-for="item in faqs" :key="item.id" class="border-b border-[#EAE5DD] py-4 md:py-5">
      <button @click="toggle(item.id)" class="flex justify-between items-center w-full text-left font-medium text-sm md:text-base min-h-[48px] gap-4 active:text-[#D4A853] transition-colors">
        <span class="flex-1">{{ item.question }}</span>
        <span class="text-gray-400 shrink-0 text-lg">{{ openId === item.id ? '−' : '+' }}</span>
      </button>
      <transition name="accordion">
        <div v-if="openId === item.id" class="mt-3 text-gray-600 text-xs md:text-sm leading-relaxed pr-6">{{ item.answer }}</div>
      </transition>
    </div>

    <!-- 空态 -->
    <div v-if="!loading && !faqs?.length" class="text-center py-16 text-gray-400">No FAQs yet.</div>

    <!-- 加载更多 -->
    <div v-if="!loading && faqs.length" class="text-center mt-10 md:mt-12">
      <button
        v-if="hasMore"
        @click="loadMore"
        :disabled="loadingMore"
        class="inline-flex items-center gap-2 px-8 py-3.5 border-2 border-gray-200 rounded-full text-sm font-medium text-gray-700 hover:border-black hover:text-black transition disabled:opacity-50 min-h-[48px]"
      >
        <span v-if="loadingMore" class="inline-block w-4 h-4 border-2 border-gray-400 border-t-transparent rounded-full animate-spin" />
        {{ loadingMore ? $t('common.loading') : $t('product.load_more') }}
      </button>
      <p v-else class="text-sm text-gray-400">{{ $t('product.no_products') || 'No more FAQs' }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
const api = useApi()
const { locale } = useI18n()
const route = useRoute()
const openId = ref<string | null>(null)
function toggle(id: string) { openId.value = openId.value === id ? null : id }
useHead({ title: 'FAQ - Sportswear OEM/ODM' })

const faqs = ref<any[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = ref(Number(route.query._pageSize) || 20)
const hasMore = ref(true)
const total = ref(0)

// SSR 首屏：内联到 HTML 以提升 SEO
const { data: initialData } = await useAsyncData<any>(
  'faqs-' + (locale.value || 'en'),
  () => api.getFaqs({ page: 1, pageSize: pageSize.value })
    .then((res: any) => ({
      items: res?.items || (Array.isArray(res) ? res : []),
      total: typeof res?.total === 'number' ? res.total : (Array.isArray(res) ? res.length : 0),
    }))
    .catch(() => ({ items: [], total: 0 })),
)

if (initialData.value) {
  faqs.value = initialData.value.items || []
  total.value = initialData.value.total || 0
  hasMore.value = faqs.value.length > 0 && faqs.value.length < total.value
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  const nextPage = page.value + 1
  try {
    const res = await api.getFaqs({ page: nextPage, pageSize: pageSize.value })
    const newItems = Array.isArray(res) ? res : (res?.items || [])
    if (newItems.length > 0) {
      const existingIds = new Set(faqs.value.map((f: any) => f.id))
      const unique = newItems.filter((f: any) => !existingIds.has(f.id))
      if (unique.length > 0) {
        faqs.value.push(...unique)
        page.value = nextPage
      }
    }
    const latestTotal = typeof res?.total === 'number' ? res.total : total.value
    total.value = latestTotal
    hasMore.value = faqs.value.length > 0 && faqs.value.length < latestTotal
  } catch (e) {
    console.error('loadMore error:', e)
    hasMore.value = true
  } finally { loadingMore.value = false }
}
</script>