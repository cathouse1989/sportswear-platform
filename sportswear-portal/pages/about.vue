<template>
  <div>
    <section v-if="page?.modules?.length" class="max-w-7xl mx-auto px-4 lg:px-8 py-12 md:py-16">
      <PageModulesRenderer :modules="aboutModules" />
    </section>

    <section v-else-if="pending" class="max-w-7xl mx-auto px-4 lg:px-8 py-20 text-center text-gray-400">{{ $t('common.loading', 'Loading...') }}</section>

    <section v-else class="max-w-4xl mx-auto px-4 lg:px-8 py-16 md:py-20">
      <div class="bg-white rounded-2xl p-10 border border-[#EAE5DD] text-center">
        <p class="text-gray-500">{{ $t('common.no_content', 'This page has no content modules yet.') }}</p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const { locale, t } = useI18n()

// 读取 CMS「关于我们」页面（slug=about），排版与内容由后台「内容管理-页面」维护
const { data: page, pending } = await useAsyncData<any>(
  'page-about-' + (locale.value || 'en'),
  () => api.getPage('about').catch(() => null),
)

const aboutModules = computed(() => {
  return (page.value?.modules || [])
    .filter((m: any) => m.is_visible !== false)
    .sort((a: any, b: any) => (a.sort_order ?? 0) - (b.sort_order ?? 0))
})

// SEO（独立词条体系，随语言切换）
const aboutSeo = await useRouteSeo('about')
useSeoHead({
  title: () => aboutSeo.value?.title || t('about.seo_title'),
  description: () => aboutSeo.value?.description || t('about.seo_description'),
  keywords: () => aboutSeo.value?.keywords || t('about.seo_keywords'),
  ogImage: '/images/og-about.jpg',
})
</script>