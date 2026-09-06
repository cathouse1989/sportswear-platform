<template>
  <div>
    <section class="bg-gradient-to-br from-gray-900 to-gray-800 text-white py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h1 class="text-3xl md:text-4xl font-bold mb-3">{{ page?.title || 'Page' }}</h1>
        <p v-if="page?.type" class="text-white/60 text-sm uppercase tracking-wider">{{ page.type }}</p>
      </div>
    </section>

    <section v-if="page?.modules?.length" class="max-w-7xl mx-auto px-4 lg:px-8 py-12 md:py-16">
      <PageModulesRenderer :modules="visibleModules" />
    </section>

    <section v-else-if="!pending" class="max-w-4xl mx-auto px-4 lg:px-8 py-16 md:py-20">
      <div class="bg-white rounded-2xl p-10 border border-[#EAE5DD] text-center">
        <p class="text-gray-500">{{ $t('common.no_content', 'This page has no content modules yet.') }}</p>
      </div>
    </section>

    <section v-if="pending" class="max-w-7xl mx-auto px-4 lg:px-8 py-20 text-center text-gray-400">{{ $t('common.loading', 'Loading...') }}</section>

    <section v-if="!page && !pending" class="max-w-4xl mx-auto px-4 lg:px-8 py-20 text-center">
      <h2 class="text-xl font-bold text-[#0D1B2A] mb-2">{{ $t('common.not_found', 'Page Not Found') }}</h2>
      <p class="text-gray-500 mb-6">{{ $t('common.not_found_desc', 'The page you requested does not exist or has not been published yet.') }}</p>
      <NuxtLink :to="localePath('/')" class="text-[#D4A853] font-medium hover:underline">{{ $t('common.back_home', 'Back to Home') }}</NuxtLink>
    </section>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { locale } = useI18n()
const localePath = useLocalePath()
const slug = computed(() => (route.params.slug as string) || '')

const api = useApi()

const { data: page, pending } = await useAsyncData<any>(
  'page-' + slug.value + '-' + (locale.value || 'en'),
  () => api.getPage(slug.value).catch(() => null),
  {
    getCachedData(key, nuxtApp) {
      return nuxtApp.isHydrating ? nuxtApp.payload.data[key] : undefined
    },
  },
)

const visibleModules = computed(() => {
  return (page.value?.modules || [])
    .filter((m: any) => m.is_visible !== false)
    .sort((a: any, b: any) => (a.sort_order ?? 0) - (b.sort_order ?? 0))
})

useSeoHead({
  title: computed(() => page.value?.title || 'Page'),
  description: computed(() => page.value?.seo?.description || page.value?.title || ''),
})
</script>
