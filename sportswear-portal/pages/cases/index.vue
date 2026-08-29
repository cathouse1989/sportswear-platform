<template>
  <div class="max-w-7xl mx-auto px-4 py-12">
    <h1 class="text-2xl md:text-3xl font-bold mb-6 md:mb-8">{{ $t('nav.cases') }}</h1>
    <div class="grid sm:grid-cols-2 md:grid-cols-3 gap-4 md:gap-6">
      <NuxtLink v-for="c in cases" :key="c.id" :to="localePath(`/cases/${c.slug}`)" class="bg-white rounded-xl overflow-hidden shadow hover:shadow-lg transition p-4 md:p-6 active:scale-[0.98]">
        <h3 class="text-sm md:text-base font-semibold mb-2 line-clamp-2">{{ c.title }}</h3>
        <p class="text-xs md:text-sm text-gray-500">{{ c.client_industry }}</p>
      </NuxtLink>
    </div>
    <div v-if="!cases?.length" class="text-center py-16 text-gray-400">No cases yet.</div>
  </div>
</template>
<script setup lang="ts">
const localePath = useLocalePath()
const api = useApi()
const { locale } = useI18n()
useHead({ title: 'Cases - Sportswear OEM/ODM' })

// SSR 阶段拉取案例列表并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存）
const { data: cases } = await useAsyncData<any[]>(
  'cases-' + (locale.value || 'en'),
  () => api.getCases({ page: 1, pageSize: 20 })
    .then((res: any) => res?.items || res || [])
    .catch(() => []),
)
</script>