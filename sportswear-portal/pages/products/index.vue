<template>
  <div>
    <!-- Page Header -->
    <section class="bg-gradient-to-br from-gray-900 to-gray-800 text-white py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h1 class="text-4xl md:text-5xl font-bold mb-4">{{ $t('product.products') }}</h1>
        <p class="text-white/70 text-lg">{{ $t('product.products') }} — OEM/ODM {{ $t('home.hero_sub_1') }}</p>
      </div>
    </section>

    <div class="max-w-7xl mx-auto px-4 lg:px-8 py-12">
      <!-- Search Bar -->
      <div class="mb-6 md:mb-8">
        <div class="relative max-w-md">
          <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="$t('product.search')"
            class="w-full pl-12 pr-10 py-3.5 rounded-xl border border-[#EAE5DD] bg-white text-sm focus:outline-none focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] transition min-h-[48px]"
            @input="debouncedSearch"
          />
          <button v-if="searchQuery" @click="clearSearch" class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>
      </div>

      <!-- Filters -->
      <div class="flex flex-wrap items-center gap-4 mb-8 md:mb-10">
        <div class="flex flex-wrap gap-2 w-full md:w-auto overflow-x-auto md:overflow-visible pb-2 md:pb-0 -mx-1 px-1 md:mx-0 md:px-0 scrollbar-hide">
          <button
            @click="selectedCategory = ''; fetchProducts()"
            class="shrink-0 px-4 py-2.5 rounded-full text-sm font-medium transition border min-h-[40px]"
            :class="selectedCategory === '' ? 'bg-black text-white border-black' : 'bg-white text-gray-600 border-gray-200 hover:border-gray-400'"
          >
            {{ $t('product.all_categories') }}
          </button>
          <button
            v-for="cat in categories"
            :key="cat.id"
            @click="selectedCategory = cat.id; fetchProducts()"
            class="shrink-0 px-4 py-2.5 rounded-full text-sm font-medium transition border min-h-[40px]"
            :class="selectedCategory === cat.id ? 'bg-black text-white border-black' : 'bg-white text-gray-600 border-gray-200 hover:border-gray-400'"
          >
            {{ cat.name }}
          </button>
        </div>

        <select v-model="selectedCategory" @change="fetchProducts" class="px-4 py-2.5 rounded-full text-sm font-medium border border-gray-200 bg-white text-gray-600 focus:outline-none focus:border-black transition min-h-[40px]">
          <option value="">{{ $t('product.all_categories') }}</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>

        <select v-model="selectedGender" @change="fetchProducts" class="px-4 py-2.5 rounded-full text-sm font-medium border border-gray-200 bg-white text-gray-600 focus:outline-none focus:border-black transition min-h-[40px]">
          <option value="">{{ $t('product.detail.gender') }}: {{ $t('product.all_categories') }}</option>
          <option value="unisex">Unisex</option>
          <option value="male">Male</option>
          <option value="female">Female</option>
          <option value="kids">Kids</option>
        </select>

        <select v-model="selectedType" @change="fetchProducts" class="px-4 py-2.5 rounded-full text-sm font-medium border border-gray-200 bg-white text-gray-600 focus:outline-none focus:border-black transition min-h-[40px]">
          <option value="">{{ $t('product.detail.type') }}: {{ $t('product.all_categories') }}</option>
          <option value="oem">OEM</option>
          <option value="odm">ODM</option>
          <option value="both">Both</option>
        </select>
      </div>

      <!-- Products Grid -->
      <div v-if="!products.length && !loading" class="text-center py-16 md:py-20">
        <div class="text-5xl md:text-6xl mb-4">📦</div>
        <p class="text-gray-400 text-base md:text-lg">{{ $t('product.no_products') }}</p>
      </div>

      <div v-if="loading" class="grid grid-cols-2 gap-3 md:gap-6">
        <div v-for="i in 8" :key="i" class="animate-pulse">
          <div class="aspect-[4/5] bg-gray-100 rounded-xl md:rounded-2xl mb-3 md:mb-4"></div>
          <div class="h-4 bg-gray-100 rounded w-3/4 mb-2"></div>
          <div class="h-3 bg-gray-100 rounded w-1/2"></div>
        </div>
      </div>

      <div v-show="!loading" class="grid grid-cols-2 gap-3 md:gap-6">
        <NuxtLink
          v-for="p in products"
          :key="p.id"
          :to="localePath(`/products/${p.slug}`)"
          class="group bg-white rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300"
        >
          <div class="aspect-[4/5] bg-gray-100 relative overflow-hidden">
            <img v-if="p.cover_image" :src="imgUrl(p.cover_image)" :alt="p.sku"
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
              loading="lazy" />
            <div v-if="!p.cover_image" class="absolute inset-0 flex items-center justify-center text-6xl opacity-30">🏋️</div>
          </div>
          <div class="p-3 md:p-4">
            <h3 class="font-semibold text-gray-900 group-hover:text-black transition-colors text-sm md:text-base text-center">{{ p.name || p.sku }}</h3>
          </div>
        </NuxtLink>
      </div>

      <!-- 加载更多 / 结束提示 -->
      <div v-if="!loading && products.length" class="text-center mt-10 md:mt-12">
        <button
          v-if="hasMore"
          @click="loadMore"
          :disabled="loadingMore"
          class="inline-flex items-center gap-2 px-8 py-3.5 border-2 border-gray-200 rounded-full text-sm font-medium text-gray-700 hover:border-black hover:text-black transition disabled:opacity-50 min-h-[48px]"
        >
          {{ loadingMore ? $t('common.loading') : $t('product.load_more') }}
        </button>
        <p v-else class="text-sm text-gray-400">{{ $t('product.all_loaded') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const localePath = useLocalePath()
// SEO - 产品列表页
useSeoHead({
  title: 'Custom Sportswear Products - OEM/ODM Activewear Collection',
  description:
    'Browse our extensive catalog of custom sportswear products. OEM/ODM activewear, yoga wear, running gear, sports bras, leggings, and team uniforms for global brands.',
  keywords:
    'custom sportswear, OEM activewear, ODM products, sports bra manufacturer, custom leggings, yoga wear wholesale, athletic apparel catalog',
})

const api = useApi()
const route = useRoute()
const { locale } = useI18n()
const products = ref<any[]>([])
const categories = ref<any[]>([])
const selectedCategory = ref('')
const selectedGender = ref('')
const selectedType = ref('')
const searchQuery = ref('')
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = ref(Number(route.query._pageSize) || 20)
const hasMore = ref(true)
const total = ref(0)
let searchTimer: ReturnType<typeof setTimeout> | null = null

// SSR 首屏加载产品数据（内联到 HTML）
const { data: ssrData } = await useAsyncData<any>(
  'products-' + (locale.value || 'en'),
  () => api.getProducts({ page: 1, pageSize: pageSize.value, lang: locale.value })
    .then((res: any) => ({
      items: Array.isArray(res) ? res : (res?.items || []),
      total: typeof res?.total === 'number' ? res.total : 0,
    }))
    .catch(() => ({ items: [], total: 0 })),
  {
    getCachedData(key, nuxtApp) {
      return nuxtApp.isHydrating ? nuxtApp.payload.data[key] : undefined
    },
  },
)

// SSR 首屏加载分类
const { data: ssrCategories } = await useAsyncData<any[]>(
  'categories-' + (locale.value || 'en'),
  () => api.getCategories().catch(() => []),
)

// SSR 数据填充到响应式变量
if (ssrData.value) {
  products.value = ssrData.value.items || []
  total.value = ssrData.value.total || 0
  hasMore.value = products.value.length > 0 && products.value.length < total.value
}
if (ssrCategories.value) {
  categories.value = ssrCategories.value
}

onMounted(async () => {
  // 如果没有 SSR 数据（客户端导航时），才重新加载
  if (!products.value.length) {
    try {
      const [catData] = await Promise.all([
        api.getCategories(),
      ])
      categories.value = catData || []
      await fetchProducts()
    } catch (e) { console.error(e); loading.value = false }
  }
})

function debouncedSearch() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => fetchProducts(), 400)
}

function clearSearch() {
  searchQuery.value = ''
  fetchProducts()
}

// 语言切换后（同一页面组件复用、不重新挂载时），显式按新语言重新拉取产品
watch(locale, () => {
  if (!import.meta.client) return
  selectedCategory.value = ''
  selectedGender.value = ''
  selectedType.value = ''
  searchQuery.value = ''
  fetchProducts()
})

async function fetchProducts() {
  loading.value = true
  page.value = 1
  products.value = []
  try {
    const params: Record<string, any> = {
      page: 1,
      pageSize: pageSize.value,
      category_id: selectedCategory.value || undefined,
      gender: selectedGender.value || undefined,
      type: selectedType.value || undefined,
    }
    if (searchQuery.value.trim()) {
      params.search = searchQuery.value.trim()
    }
    const res = await api.getProducts(params)
    // 兼容 PageResult { items, total } 和纯数组两种返回格式
    const items = Array.isArray(res) ? res : (res?.items || [])
    const totalCount = typeof res?.total === 'number' ? res.total : items.length
    products.value = items
    total.value = totalCount
    // 只有当前页数据量还没达到总数时才显示"加载更多"
    hasMore.value = items.length > 0 && items.length < totalCount
  } catch (e) {
    console.error('fetchProducts error:', e)
    products.value = []
    total.value = 0
    hasMore.value = false
  } finally { loading.value = false }
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  const nextPage = page.value + 1
  try {
    const params: Record<string, any> = {
      page: nextPage,
      pageSize: pageSize.value,
      category_id: selectedCategory.value || undefined,
      gender: selectedGender.value || undefined,
      type: selectedType.value || undefined,
    }
    if (searchQuery.value.trim()) {
      params.search = searchQuery.value.trim()
    }
    // 防止 RateLimit(30次/分钟) 拦截 -> 429
    const res = await api.getProducts(params)
    const newItems = Array.isArray(res) ? res : (res?.items || [])

    if (newItems.length > 0) {
      // 去重，防止同一页重复添加
      const existingIds = new Set(products.value.map((p: any) => p.id))
      const uniqueNew = newItems.filter((p: any) => !existingIds.has(p.id))
      if (uniqueNew.length > 0) {
        products.value.push(...uniqueNew)
        page.value = nextPage
      }
    }

    // 用最新 total 判断是否还有更多
    const latestTotal = typeof res?.total === 'number' ? res.total : total.value
    total.value = latestTotal
    hasMore.value = products.value.length > 0 && products.value.length < latestTotal
  } catch (e) {
    // 若 API 429/500 等出错，不隐藏已有数据，也不翻到空白页；保留 hasMore 以便重试
    console.error('loadMore error:', e)
    hasMore.value = true
  } finally { loadingMore.value = false }
}
</script>