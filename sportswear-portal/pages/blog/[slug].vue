<template>
  <div v-if="blog" class="max-w-4xl mx-auto px-4 py-12">
    <p class="text-sm text-blue-600 mb-2">{{ blog.category }}</p>
    <h1 class="text-3xl font-bold mb-4">{{ blog.title }}</h1>
    <div class="h-64 bg-gray-100 rounded-xl flex items-center justify-center text-6xl mb-8">📝</div>
    <div class="prose max-w-none" v-html="blog.content"></div>
  </div>
  <div v-else class="text-center py-24 text-gray-400">Loading...</div>
</template>
<script setup lang="ts">
const route = useRoute()
const api = useApi()
const { locale } = useI18n()
const slug = route.params.slug as string

// SSR 阶段拉取博客详情并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存）
const { data: blog } = await useAsyncData<any>(
  'blog-' + (locale.value || 'en') + '-' + slug,
  () => api.getBlog(slug).catch(() => null),
)

// SEO - 博客详情页
useSeoHead({
  title: computed(() => blog.value?.title ? `${blog.value.title} - Sportswear Manufacturing Blog` : 'Blog Article'),
  description: computed(() => blog.value?.content?.replace(/<[^>]*>/g, '').slice(0, 160) || ''),
  ogImage: computed(() => blog.value?.image || undefined),
  schema: computed(() => blog.value ? buildArticleSchema(blog.value) : undefined),
})
</script>