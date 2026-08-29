<template>
  <div>
    <!-- Hero：数据驱动轮播，展示参数来自后台 hero_settings -->
    <HeroCarousel
      :slides="heroSlides"
      :settings="heroSettings"
      :default-button-text="t('home.get_quote')"
    />

    <!-- Featured Products -->
    <section v-if="products?.length" class="py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <div class="flex justify-between items-end mb-8 md:mb-10">
          <h2 class="text-2xl md:text-3xl font-bold text-[#0D1B2A]">{{ $t('home.featured_products') }}</h2>
          <NuxtLink :to="localePath('/products')" class="text-sm text-[#D4A853] font-medium min-h-[44px] flex items-center">{{ $t('home.view_all') }} →</NuxtLink>
        </div>
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4 md:gap-6">
          <NuxtLink
            v-for="p in products.slice(0, 8)"
            :key="p.id"
            :to="localePath('/products/' + p.slug)"
            class="group bg-white rounded-2xl overflow-hidden border border-[#EAE5DD] hover:shadow-lg transition-all active:scale-[0.98]"
          >
            <div class="aspect-[3/4] bg-gray-100 relative overflow-hidden">
              <img
                v-if="p.cover_image"
                :src="p.cover_image"
                :alt="p.sku"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
                loading="lazy"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-4xl md:text-6xl">🏋️</div>
              <div class="absolute bottom-3 left-3 md:bottom-4 md:left-4">
                <span class="px-2.5 md:px-3 py-1 bg-white/95 rounded-full text-[10px] md:text-xs font-semibold">{{ p.type }}</span>
              </div>
            </div>
            <div class="p-3 md:p-5">
              <h3 class="text-xs md:text-sm font-semibold text-[#0D1B2A]">{{ p.sku }}</h3>
              <p class="text-[10px] md:text-sm text-gray-500 mt-0.5 md:mt-1">{{ p.brief?.slice(0, 60) }}</p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const localePath = useLocalePath()
const { t, locale } = useI18n()

const api = useApi()
const { data: homeData } = await useAsyncData<any>(
  'home-' + (locale.value || 'en'),
  () => api.getHome().catch(() => null),
)
const products = computed<any[]>(() => homeData.value?.featured_products || [])

const DEFAULT_IMAGES = [
  'https://images.unsplash.com/photo-1599901860904-17e6ed7083a0?w=1600&h=900&fit=crop',
  'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=1600&h=900&fit=crop',
  'https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=1600&h=900&fit=crop',
  'https://images.unsplash.com/photo-1577221084712-45b0445d2b00?w=1600&h=900&fit=crop',
]

// 优先后台 hero_slides；缺失时回退内置默认图/文案
const heroSlides = computed(() => {
  const raw = homeData.value?.hero_slides
  if (Array.isArray(raw) && raw.length) {
    return raw.map((s: any, i: number) => ({
      image: String(s?.image || ''),
      title: String(s?.title || t(`home.hero_title_${i + 1}`)),
      subtitle: String(s?.subtitle || t(`home.hero_sub_${i + 1}`)),
      button_text: String(s?.button_text || t('home.get_quote')),
      button_url: String(s?.button_url || '/contact'),
    }))
  }
  return [0, 1, 2, 3].map((i) => ({
    image: DEFAULT_IMAGES[i],
    title: t(`home.hero_title_${i + 1}`),
    subtitle: t(`home.hero_sub_${i + 1}`),
    button_text: t('home.get_quote'),
    button_url: '/contact',
  }))
})

// 展示参数：/public/home.hero_settings（后台可配）
const heroSettings = computed(() => {
  const s = homeData.value?.hero_settings || {}
  return {
    autoplay: s.autoplay !== false,
    interval_ms: Number(s.interval_ms) || 5000,
    transition: s.transition === 'slide' ? 'slide' : 'fade',
    show_dots: s.show_dots !== false,
    show_arrows: s.show_arrows === true,
    pause_on_hover: s.pause_on_hover !== false,
  }
})
</script>
