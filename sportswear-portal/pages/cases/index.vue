<template>
  <div class="max-w-7xl mx-auto px-4 py-12">
    <h1 class="text-2xl md:text-3xl font-bold mb-4 md:mb-6">{{ $t('nav.cases') }}</h1>

    <!-- 项目类型筛选 -->
    <div class="flex flex-wrap gap-2 mb-8">
      <button
        v-for="chip in typeChips"
        :key="chip.value"
        @click="selectType(chip.value)"
        :class="[
          'px-4 py-2 rounded-full text-sm font-medium transition border min-h-[40px]',
          activeType === chip.value
            ? 'bg-black text-white border-black'
            : 'bg-white text-gray-700 border-gray-200 hover:border-black hover:text-black',
        ]"
      >
        {{ chip.label }}
      </button>
    </div>

    <!-- 骨架屏 -->
    <div v-if="loading" class="grid sm:grid-cols-2 md:grid-cols-3 gap-4 md:gap-6">
      <div v-for="i in 6" :key="i" class="animate-pulse bg-white rounded-xl overflow-hidden">
        <div class="aspect-[16/10] bg-gray-100" />
        <div class="p-4 md:p-6">
          <div class="h-3 bg-gray-100 rounded w-1/3 mb-3" />
          <div class="h-4 bg-gray-100 rounded w-3/4 mb-2" />
          <div class="h-3 bg-gray-100 rounded w-full" />
        </div>
      </div>
    </div>

    <!-- 案例卡片 -->
    <div v-show="!loading" class="grid sm:grid-cols-2 md:grid-cols-3 gap-4 md:gap-6">
      <NuxtLink
        v-for="c in cases"
        :key="c.id"
        :to="localePath(`/cases/${c.slug}`)"
        class="group bg-white rounded-xl overflow-hidden shadow hover:shadow-lg transition flex flex-col active:scale-[0.98]"
      >
        <div class="aspect-[16/10] overflow-hidden bg-gray-100">
          <img
            v-if="c.cover_image"
            :src="c.cover_image"
            :alt="c.title"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
          <div v-else class="w-full h-full flex items-center justify-center text-gray-300 text-sm font-medium">{{ $t('nav.cases') }}</div>
        </div>
        <div class="p-4 md:p-6 flex flex-col flex-1">
          <div class="flex flex-wrap items-center gap-2 mb-2">
            <span v-if="c.project_type" class="px-2 py-0.5 rounded-full bg-blue-50 text-blue-700 text-xs font-medium">{{ projectTypeLabel(c.project_type) }}</span>
            <span v-if="c.client_industry" class="text-xs text-gray-500">{{ c.client_industry }}</span>
          </div>
          <h3 class="text-sm md:text-base font-semibold mb-2 line-clamp-2">{{ c.title }}</h3>
          <p v-if="summary(c)" class="text-xs md:text-sm text-gray-500 line-clamp-2 mb-4">{{ summary(c) }}</p>
          <span class="mt-auto inline-flex items-center gap-1.5 text-sm font-medium text-black">
            {{ $t('case.view_case') }} <span aria-hidden="true">→</span>
          </span>
        </div>
      </NuxtLink>
    </div>

    <!-- 空态 -->
    <div v-if="!loading && !cases.length" class="text-center py-16 text-gray-400">{{ $t('case.empty') }}</div>

    <!-- 加载更多 -->
    <div v-if="!loading && cases.length" class="text-center mt-10 md:mt-12">
      <button
        v-if="hasMore"
        @click="loadMore"
        :disabled="loadingMore"
        class="inline-flex items-center gap-2 px-8 py-3.5 border-2 border-gray-200 rounded-full text-sm font-medium text-gray-700 hover:border-black hover:text-black transition disabled:opacity-50 min-h-[48px]"
      >
        <span v-if="loadingMore" class="inline-block w-4 h-4 border-2 border-gray-400 border-t-transparent rounded-full animate-spin" />
        {{ loadingMore ? $t('common.loading') : $t('common.load_more') }}
      </button>
      <p v-else class="text-sm text-gray-400">{{ $t('case.all_loaded') }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
const localePath = useLocalePath()
const api = useApi()
const { locale, t, te } = useI18n()
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
const activeType = ref('')

const typeChips = computed(() => [
  { label: t('case.all'), value: '' },
  { label: t('case.project_type_options.oem'), value: 'oem' },
  { label: t('case.project_type_options.odm'), value: 'odm' },
  { label: t('case.project_type_options.private_label'), value: 'private_label' },
])

// 本地化项目类型标签
function projectTypeLabel(value?: string) {
  if (!value) return ''
  const key = `case.project_type_options.${value}`
  return te(key) ? t(key) : value
}

// 卡片摘要：优先客户需求，回退解决方案；去除 HTML 标签
function summary(c: any) {
  const raw = c.client_need || c.solution || ''
  return (raw || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').replace(/\s+/g, ' ').trim()
}

async function fetchCases() {
  const res = await api.getCases({ page: page.value, pageSize: pageSize.value, project_type: activeType.value })
  return {
    items: res?.items || (Array.isArray(res) ? res : []),
    total: typeof res?.total === 'number' ? res.total : (Array.isArray(res) ? res.length : 0),
  }
}

// SSR 首屏：内联到 HTML 以提升 SEO
const { data: initialData } = await useAsyncData<any>(
  'cases-' + (locale.value || 'en'),
  () => fetchCases().catch(() => ({ items: [], total: 0 })),
)

if (initialData.value) {
  cases.value = initialData.value.items || []
  total.value = initialData.value.total || 0
  hasMore.value = cases.value.length > 0 && cases.value.length < total.value
}

// 切换项目类型：重置分页并重新加载（服务端筛选）
async function selectType(value: string) {
  if (activeType.value === value) return
  activeType.value = value
  page.value = 1
  cases.value = []
  total.value = 0
  hasMore.value = true
  loading.value = true
  try {
    const res = await fetchCases()
    cases.value = res.items
    total.value = res.total
    hasMore.value = cases.value.length > 0 && cases.value.length < total.value
  } finally { loading.value = false }
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  const nextPage = page.value + 1
  try {
    const res = await api.getCases({ page: nextPage, pageSize: pageSize.value, project_type: activeType.value })
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