<template>
  <div>
    <section class="bg-gradient-to-br from-gray-900 to-gray-800 text-white py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h1 class="text-3xl md:text-4xl font-bold mb-3">{{ page?.title || 'Page' }}</h1>
        <p v-if="page?.type" class="text-white/60 text-sm uppercase tracking-wider">{{ page.type }}</p>
      </div>
    </section>

    <section v-if="page?.modules?.length" class="max-w-7xl mx-auto px-4 lg:px-8 py-12 md:py-16 space-y-12">
      <template v-for="m in visibleModules" :key="m.id">
        <div v-if="m.type === 'text'" class="prose max-w-none">
          <h2 v-if="m.title" class="text-2xl font-bold text-[#0D1B2A] mb-4">{{ m.title }}</h2>
          <div class="text-gray-600 leading-relaxed" v-html="m.config?.content || m.title"></div>
        </div>
        <div v-else-if="m.type === 'image'" class="text-center">
          <img v-if="m.config?.url" :src="m.config.url" :alt="m.title" class="rounded-2xl max-w-full mx-auto" />
          <p v-if="m.title" class="text-sm text-gray-500 mt-3">{{ m.title }}</p>
        </div>
        <div v-else-if="m.type === 'video'" class="aspect-video rounded-2xl overflow-hidden bg-black">
          <iframe v-if="m.config?.url" :src="m.config.url" class="w-full h-full" frameborder="0" allowfullscreen />
          <div v-else class="w-full h-full flex items-center justify-center text-white/50">No Video</div>
        </div>
        <div v-else-if="m.type === 'banner'" class="rounded-2xl overflow-hidden">
          <img v-if="m.config?.image || m.config?.url" :src="m.config?.image || m.config?.url" :alt="m.title" class="w-full h-auto max-h-[480px] object-cover rounded-2xl" />
        </div>
        <div v-else-if="['product_recommend','oem','odm','factory','production','case','certification','blog'].includes(m.type)" class="bg-white rounded-2xl p-8 border border-[#EAE5DD]">
          <h2 v-if="m.title" class="text-xl font-bold text-[#0D1B2A] mb-4">{{ m.title }}</h2>
          <div v-if="m.config?.description" class="text-gray-600 leading-relaxed" v-html="m.config.description"></div>
          <div v-else class="text-sm text-gray-400">Module type: {{ m.type }}</div>
        </div>
        <div v-else-if="m.type === 'contact'" class="bg-gradient-to-r from-[#F5F0E8] to-[#FBF9F6] rounded-2xl p-8 text-center">
          <h2 v-if="m.title" class="text-2xl font-bold text-[#0D1B2A] mb-3">{{ m.title }}</h2>
          <p class="text-gray-600 mb-6">{{ m.config?.description || 'Ready to start your project?' }}</p>
          <NuxtLink :to="localePath('/contact')" class="inline-block bg-[#D4A853] text-white px-8 py-3.5 rounded-full font-semibold hover:bg-[#C49A3F] transition shadow-sm">
            {{ $t('home.get_quote') }}
          </NuxtLink>
        </div>
        <div v-else class="bg-gray-50 rounded-2xl p-6 border border-dashed border-gray-300 text-center text-sm text-gray-400">
          Unknown module: {{ m.type }}
        </div>
      </template>
    </section>

    <section v-else-if="!pending" class="max-w-4xl mx-auto px-4 lg:px-8 py-16 md:py-20">
      <div class="bg-white rounded-2xl p-10 border border-[#EAE5DD] text-center">
        <p class="text-gray-500">This page has no content modules yet.</p>
      </div>
    </section>

    <section v-if="pending" class="max-w-7xl mx-auto px-4 lg:px-8 py-20 text-center text-gray-400">Loading...</section>

    <section v-if="!page && !pending" class="max-w-4xl mx-auto px-4 lg:px-8 py-20 text-center">
      <h2 class="text-xl font-bold text-[#0D1B2A] mb-2">{{ $t('common.not_found') || 'Page Not Found' }}</h2>
      <p class="text-gray-500 mb-6">The page you requested does not exist or has not been published yet.</p>
      <NuxtLink :to="localePath('/')" class="text-[#D4A853] font-medium hover:underline">Back to Home</NuxtLink>
    </section>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { t, locale } = useI18n()
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

useHead({
  title: computed(() => page.value?.title || 'Page'),
  meta: [
    { name: 'description', content: computed(() => page.value?.seo?.description || page.value?.title || '') },
  ],
})
</script>
﻿<template>
  <div>
    <section class="bg-gradient-to-br from-gray-900 to-gray-800 text-white py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h1 class="text-3xl md:text-4xl font-bold mb-3">{{ page?.title || '\u9875\u9762' }}</h1>
        <p v-if="page?.type" class="text-white/60 text-sm uppercase tracking-wider">{{ page.type }}</p>
      </div>
    </section>

    <section v-if="page?.modules?.length" class="max-w-7xl mx-auto px-4 lg:px-8 py-12 md:py-16 space-y-12">
      <template v-for="m in visibleModules" :key="m.id">
        <div v-if="m.type === 'text'" class="prose max-w-none">
          <h2 v-if="m.title" class="text-2xl font-bold text-[#0D1B2A] mb-4">{{ m.title }}</h2>
          <div class="text-gray-600 leading-relaxed" v-html="m.config?.content || m.title"></div>
        </div>
        <div v-else-if="m.type === 'image'" class="text-center">
          <img v-if="m.config?.url" :src="m.config.url" :alt="m.title" class="rounded-2xl max-w-full mx-auto" />
          <p v-if="m.title" class="text-sm text-gray-500 mt-3">{{ m.title }}</p>
        </div>
        <div v-else-if="m.type === 'video'" class="aspect-video rounded-2xl overflow-hidden bg-black">
          <iframe v-if="m.config?.url" :src="m.config.url" class="w-full h-full" frameborder="0" allowfullscreen />
          <div v-else class="w-full h-full flex items-center justify-center text-white/50">🎬 暂无视频</div>
        </div>
        <div v-else-if="m.type === 'banner'" class="rounded-2xl overflow-hidden">
          <img v-if="m.config?.image || m.config?.url" :src="m.config?.image || m.config?.url" :alt="m.title" class="w-full h-auto max-h-[480px] object-cover rounded-2xl" />
        </div>
        <div v-else-if="['product_recommend','oem','odm','factory','production','case','certification','blog'].includes(m.type)" class="bg-white rounded-2xl p-8 border border-[#EAE5DD]">
          <h2 v-if="m.title" class="text-xl font-bold text-[#0D1B2A] mb-4">{{ m.title }}</h2>
          <div v-if="m.config?.description" class="text-gray-600 leading-relaxed" v-html="m.config.description"></div>
          <div v-else class="text-sm text-gray-400">模块类型：{{ m.type }}（待配置展示内容）</div>
        </div>
        <div v-else-if="m.type === 'contact'" class="bg-gradient-to-r from-[#F5F0E8] to-[#FBF9F6] rounded-2xl p-8 text-center">
          <h2 v-if="m.title" class="text-2xl font-bold text-[#0D1B2A] mb-3">{{ m.title }}</h2>
          <p class="text-gray-600 mb-6">{{ m.config?.description || 'Ready to start your project? Contact us today.' }}</p>
          <NuxtLink :to="localePath('/contact')" class="inline-block bg-[#D4A853] text-white px-8 py-3.5 rounded-full font-semibold hover:bg-[#C49A3F] transition shadow-sm">
            {{ $t('home.get_quote') }}
          </NuxtLink>
        </div>
        <div v-else class="bg-gray-50 rounded-2xl p-6 border border-dashed border-gray-300 text-center text-sm text-gray-400">
          模块类型「{{ m.type }}」— 暂无预览组件
        </div>
      </template>
    </section>

    <section v-else-if="!pending" class="max-w-4xl mx-auto px-4 lg:px-8 py-16 md:py-20">
      <div class="bg-white rounded-2xl p-10 border border-[#EAE5DD] text-center">
        <p class="text-gray-500">此页面暂无内容模块，请在管理后台「页面管理」中编辑。</p>
      </div>
    </section>

    <section v-if="pending" class="max-w-7xl mx-auto px-4 lg:px-8 py-20 text-center text-gray-400">加载中...</section>

    <section v-if="!page && !pending" class="max-w-4xl mx-auto px-4 lg:px-8 py-20 text-center">
      <h2 class="text-xl font-bold text-[#0D1B2A] mb-2">{{ $t('common.not_found') || 'Page Not Found' }}</h2>
      <p class="text-gray-500 mb-6">The page you requested does not exist or has not been published yet.</p>
      <NuxtLink :to="localePath('/')" class="text-[#D4A853] font-medium hover:underline">← {{ $t('nav.home') }}</NuxtLink>
    </section>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { t, locale } = useI18n()
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

useHead({
  title: computed(() => page.value?.title || 'Page'),
  meta: [
    { name: 'description', content: computed(() => page.value?.seo?.description || page.value?.title || '') },
  ],
})
</script>