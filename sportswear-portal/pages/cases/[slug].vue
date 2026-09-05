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
          <p class="text-gray-600 leading-relaxed whitespace-pre-line">{{ section.value }}</p>
        </section>
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

// SEO - 案例详情页
useSeoHead({
  title: computed(() => (caseItem.value?.title ? caseItem.value.title : 'Case Study')),
  description: computed(() => {
    const c = caseItem.value || {}
    return (c.client_need || c.solution || '').slice(0, 160)
  }),
  ogImage: computed(() => caseItem.value?.cover_image || undefined),
  ogType: 'article',
})
</script>
