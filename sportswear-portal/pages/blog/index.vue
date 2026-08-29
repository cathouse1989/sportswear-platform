<template>
  <div class="max-w-7xl mx-auto px-4 py-12">
    <h1 class="text-2xl md:text-3xl font-bold mb-6 md:mb-8">{{ $t('nav.blog') }}</h1>
    <div class="grid sm:grid-cols-2 md:grid-cols-3 gap-4 md:gap-6">
      <NuxtLink v-for="b in blogs" :key="b.id" :to="localePath(`/blog/${b.slug}`)" class="bg-white rounded-xl overflow-hidden shadow hover:shadow-lg transition active:scale-[0.98]">
        <div class="h-36 md:h-44 bg-gray-100 flex items-center justify-center text-3xl md:text-4xl">📝</div>
        <div class="p-3 md:p-4">
          <p class="text-[10px] md:text-xs text-blue-600 mb-1">{{ b.category }}</p>
          <h3 class="text-sm md:text-base font-semibold mb-1 md:mb-2 line-clamp-2">{{ b.title }}</h3>
          <p class="text-xs md:text-sm text-gray-500 line-clamp-3">{{ b.content?.slice(0, 100) }}</p>
        </div>
      </NuxtLink>
    </div>
    <div v-if="!blogs?.length" class="text-center py-16 text-gray-400">{{ $t('common.no_data') || 'No articles yet.'  }}</div>
  </div>
</template>
<script setup lang="ts">
const localePath = useLocalePath()
const api = useApi()
const { locale } = useI18n()

useHead({ title: 'Blog - Sportswear OEM/ODM' })

// SSR 阶段拉取博客列表并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存命中完整内容）
const { data: blogs } = await useAsyncData<any[]>(
  'blogs-' + (locale.value || 'en'),
  () => api.getBlogs({ page: 1, pageSize: 20 })
    .then((res: any) => res?.items || res || [])
    .catch(() => []),
)
</script>