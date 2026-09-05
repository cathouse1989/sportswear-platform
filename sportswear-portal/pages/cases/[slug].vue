<template>
  <div class="max-w-4xl mx-auto px-4 py-12 md:py-16">
    <!-- 返回列表 -->
    <NuxtLink :to="localePath('/cases')" class="inline-flex items-center gap-1.5 text-sm text-gray-500 hover:text-black mb-8 transition-colors">
      <span aria-hidden="true">←</span> {{ $t('case.back_to_list') }}
    </NuxtLink>

    <template v-if="caseItem">
      <!-- 元信息 -->
      <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-500 mb-5">
        <span v-if="caseItem.project_type" class="px-2.5 py-1 rounded-full bg-blue-50 text-blue-700 font-medium">{{ projectTypeLabel(caseItem.project_type) }}</span>
        <span v-if="caseItem.client_industry" class="px-2.5 py-1 rounded-full bg-gray-100 text-gray-600">{{ caseItem.client_industry }}</span>
        <time v-if="caseItem.published_at || caseItem.created_at">{{ formatDate(caseItem.published_at || caseItem.created_at) }}</time>
      </div>

      <h1 class="text-3xl md:text-4xl font-bold leading-tight mb-8 md:mb-10">{{ caseItem.title }}</h1>

      <img v-if="caseItem.cover_image" :src="caseItem.cover_image" :alt="caseItem.title" class="w-full rounded-xl object-cover mb-10 md:mb-12 max-h-[480px] bg-gray-100" />

      <!-- 项目概览 -->
      <div v-if="overviewItems.length" class="grid sm:grid-cols-3 gap-4 mb-10 md:mb-12">
        <div v-for="item in overviewItems" :key="item.label" class="rounded-xl border border-gray-100 bg-gray-50 p-4">
          <div class="text-xs uppercase tracking-wide text-gray-400 mb-1">{{ item.label }}</div>
          <div class="text-sm font-medium text-gray-800">{{ item.value }}</div>
        </div>
      </div>

      <!-- 案例正文分节 -->
      <div class="space-y-8 md:space-y-10">
        <section v-for="section in sections" :key="section.key">
          <h2 class="text-lg md:text-xl font-semibold mb-3 flex items-center gap-2.5">
            <span class="w-1 h-5 rounded-full bg-black inline-block"></span>{{ section.label }}
          </h2>
          <div class="text-gray-600 leading-relaxed prose prose-base max-w-none" v-html="sanitizeHtml(section.value)"></div>
        </section>
      </div>

      <!-- 询盘 CTA -->
      <div class="mt-12 md:mt-16 rounded-2xl bg-gray-50 border border-gray-100 p-6 md:p-10 text-center">
        <h2 class="text-xl md:text-2xl font-bold mb-3">{{ $t('case.cta_title') }}</h2>
        <p class="text-gray-600 mb-6 max-w-2xl mx-auto">{{ $t('case.cta_desc') }}</p>
        <NuxtLink :to="localePath({ path: '/contact', query: { case_title: caseItem.title, project_type: caseItem.project_type } })" class="inline-flex items-center gap-2 px-8 py-3.5 bg-black text-white rounded-full text-sm font-medium hover:bg-gray-800 transition min-h-[48px]">
          {{ $t('case.cta_button') }} <span aria-hidden="true">→</span>
        </NuxtLink>
      </div>

      <!-- 相关案例 -->
      <div v-if="caseItem?.project_type && (relatedLoading || relatedCases.length)" class="mt-14 md:mt-16">
        <h2 class="text-xl md:text-2xl font-bold mb-6">{{ $t('case.related') }}</h2>
        <!-- 骨架 -->
        <div v-if="relatedLoading" class="grid sm:grid-cols-2 gap-4 md:gap-6">
          <div v-for="i in 2" :key="i" class="animate-pulse bg-white rounded-xl overflow-hidden">
            <div class="aspect-[16/9] bg-gray-100" />
            <div class="p-4">
              <div class="h-4 bg-gray-100 rounded w-3/4 mb-2" />
              <div class="h-3 bg-gray-100 rounded w-1/2" />
            </div>
          </div>
        </div>
        <!-- 卡片 -->
        <div v-else class="grid sm:grid-cols-2 gap-4 md:gap-6">
          <NuxtLink
            v-for="r in relatedCases"
            :key="r.id"
            :to="localePath(`/cases/${r.slug}`)"
            class="group bg-white rounded-xl overflow-hidden shadow hover:shadow-lg transition flex flex-col active:scale-[0.98]"
          >
            <div class="aspect-[16/9] overflow-hidden bg-gray-100">
              <img v-if="r.cover_image" :src="r.cover_image" :alt="r.title" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
              <div v-else class="w-full h-full flex items-center justify-center text-gray-300 text-sm font-medium">{{ $t('nav.cases') }}</div>
            </div>
            <div class="p-4 flex flex-col flex-1">
              <h3 class="text-sm font-semibold line-clamp-2 mb-1">{{ r.title }}</h3>
              <p v-if="r.client_industry" class="text-xs text-gray-500 line-clamp-1">{{ r.client_industry }}</p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </template>

    <div v-else-if="!pending" class="text-center py-24">
      <p class="text-gray-400 mb-5">{{ $t('case.not_found') }}</p>
      <NuxtLink :to="localePath('/cases')" class="inline-flex items-center gap-1.5 text-blue-700 font-medium">
        <span aria-hidden="true">←</span> {{ $t('case.back_to_list') }}
      </NuxtLink>
    </div>
    <div v-else class="text-center py-24 text-gray-400">{{ $t('common.loading') }}</div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useApi()
const localePath = useLocalePath()
const { locale, t, te } = useI18n()
const slug = route.params.slug as string

// SSR 阶段拉取案例详情并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存）
const { data: caseItem, pending } = await useAsyncData<any>(
  'case-' + (locale.value || 'en') + '-' + slug,
  () => api.getCase(slug).catch(() => null),
)

// 本地化项目类型标签
function projectTypeLabel(value?: string) {
  if (!value) return ''
  const key = `case.project_type_options.${value}`
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

// 项目概览（仅展示含值的项）
const overviewItems = computed(() => {
  const c = caseItem.value || {}
  const items = [
    { label: t('case.client_industry'), value: c.client_industry },
    { label: t('case.project_type'), value: projectTypeLabel(c.project_type) },
    { label: t('case.products'), value: c.products },
  ]
  return items.filter((i) => i.value)
})

// 案例正文分节（仅展示含值的项）
const sectionKeys = ['client_need', 'problem', 'solution', 'process', 'result'] as const
const sections = computed(() => {
  const c = caseItem.value || {}
  return sectionKeys
    .map((key) => ({ key, label: t(`case.${key}`), value: c[key] }))
    .filter((s) => s.value)
})

// 相关案例（同项目类型，排除当前，客户端拉取）
const relatedCases = ref<any[]>([])
const relatedLoading = ref(true)
async function loadRelated() {
  if (!caseItem.value?.project_type) {
    relatedLoading.value = false
    return
  }
  try {
    const res = await api.getCases({ page: 1, pageSize: 3, project_type: caseItem.value.project_type })
    const items = res?.items || (Array.isArray(res) ? res : [])
    relatedCases.value = (items || []).filter((c: any) => c.id !== caseItem.value.id).slice(0, 2)
  } catch {
    relatedCases.value = []
  } finally {
    relatedLoading.value = false
  }
}
onMounted(loadRelated)

// 去除 HTML 标签生成纯文本摘要（供 meta description 使用）
function stripHtml(html?: string) {
  return (html || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').replace(/\s+/g, ' ').trim()
}

// 富文本安全渲染兜底（对齐博客详情页）
function sanitizeHtml(html?: string) {
  if (!html) return ''
  return html
    .replace(/<\s*(script|iframe|object|embed|style|link|meta)\b[^>]*>[\s\S]*?<\s*\/\s*\1\s*>/gi, '')
    .replace(/<\s*(script|iframe|object|embed|style|link|meta)\b[^>]*\/?\s*>/gi, '')
    .replace(/\son\w+\s*=\s*("[^"]*"|'[^']*'|[^\s>]+)/gi, '')
    .replace(/(href|src)\s*=\s*("|')?\s*javascript:[^"'\s>]*("|')?/gi, '$1=""')
}

// SEO - 案例详情页
useSeoHead({
  title: computed(() => caseItem.value?.seo?.title || caseItem.value?.title || 'Case Study'),
  description: computed(() => {
    const c = caseItem.value || {}
    return stripHtml(c.seo?.description || c.client_need || c.solution || '').slice(0, 160)
  }),
  keywords: computed(() => caseItem.value?.seo?.keywords || undefined),
  ogImage: computed(() => caseItem.value?.seo?.og_image || caseItem.value?.cover_image || undefined),
  ogType: 'article',
})
</script>
