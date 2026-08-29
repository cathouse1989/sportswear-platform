<template>
  <div class="max-w-4xl mx-auto px-4 py-12">
    <h1 class="text-2xl md:text-3xl font-bold mb-6 md:mb-8">{{ $t('nav.faq') }}</h1>
    <div v-for="item in faqs" :key="item.id" class="border-b border-[#EAE5DD] py-4 md:py-5">
      <button @click="toggle(item.id)" class="flex justify-between items-center w-full text-left font-medium text-sm md:text-base min-h-[48px] gap-4 active:text-[#D4A853] transition-colors">
        <span class="flex-1">{{ item.question }}</span>
        <span class="text-gray-400 shrink-0 text-lg">{{ openId === item.id ? '−' : '+' }}</span>
      </button>
      <transition name="accordion">
        <div v-if="openId === item.id" class="mt-3 text-gray-600 text-xs md:text-sm leading-relaxed pr-6">{{ item.answer }}</div>
      </transition>
    </div>
    <div v-if="!faqs?.length" class="text-center py-16 text-gray-400">No FAQs yet.</div>
  </div>
</template>
<script setup lang="ts">
const api = useApi()
const { locale } = useI18n()
const openId = ref<string | null>(null)
function toggle(id: string) { openId.value = openId.value === id ? null : id }
useHead({ title: 'FAQ - Sportswear OEM/ODM' })

// SSR 阶段拉取 FAQ 列表并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存）
const { data: faqs } = await useAsyncData<any[]>(
  'faqs-' + (locale.value || 'en'),
  () => api.getFaqs({ page: 1, pageSize: 50 })
    .then((res: any) => res?.items || res || [])
    .catch(() => []),
)
</script>