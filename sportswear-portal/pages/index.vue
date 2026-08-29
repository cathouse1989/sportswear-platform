<template>
  <div>
    <!-- Hero -->
    <section
      class="relative min-h-screen flex items-center bg-[#0D1B2A] overflow-hidden touch-pan-y"
      @touchstart="onTouchStart"
      @touchend="onTouchEnd"
      @mouseenter="paused = true"
      @mouseleave="paused = false"
    >
      <div v-for="(s, i) in heroSlides" :key="i"
        class="absolute inset-0 transition-opacity duration-1000"
        :class="currentSlide === i ? 'opacity-100' : 'opacity-0'">
        <img v-if="s.image && !imgFailed[i]"
          :src="s.image" :alt="s.title"
          class="w-full h-full object-cover"
          @error="onImgError(i)" />
        <div class="absolute inset-0 bg-gradient-to-r from-[#0D1B2A]/95 via-[#0D1B2A]/70 to-transparent"></div>
      </div>
      <div class="relative max-w-7xl mx-auto px-4 lg:px-8 py-20 w-full">
        <div class="max-w-2xl">
          <div v-for="(s, i) in heroSlides" :key="i" :class="currentSlide === i ? '' : 'hidden'">
            <h1 class="text-4xl sm:text-5xl md:text-7xl font-bold text-white mb-4 md:mb-6 leading-tight">{{ s.title }}</h1>
            <p v-if="s.subtitle" class="text-base sm:text-lg md:text-xl text-white/60 mb-8 md:mb-10">{{ s.subtitle }}</p>
            <NuxtLink :to="buttonTarget(i)"
              class="inline-flex items-center gap-3 bg-[#D4A853] text-white px-8 md:px-10 py-3.5 md:py-4 rounded-full font-semibold hover:bg-[#C49A3F] transition text-sm md:text-base min-h-[48px]">
              {{ buttonText(i) }}
            </NuxtLink>
          </div>
        </div>
      </div>
      <div class="absolute bottom-12 left-1/2 -translate-x-1/2 flex gap-3">
        <button v-for="(s, i) in heroSlides" :key="i" @click="currentSlide = i"
          class="w-2.5 h-2.5 rounded-full transition-all min-h-[10px]"
          :class="currentSlide === i ? 'bg-[#D4A853] w-8' : 'bg-white/30'"></button>
      </div>
    </section>

    <!-- Featured Products -->
    <section v-if="products?.length" class="py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <div class="flex justify-between items-end mb-8 md:mb-10">
          <h2 class="text-2xl md:text-3xl font-bold text-[#0D1B2A]">{{ $t('home.featured_products') }}</h2>
          <NuxtLink :to="localePath('/products')" class="text-sm text-[#D4A853] font-medium min-h-[44px] flex items-center">{{ $t('home.view_all') }} →</NuxtLink>
        </div>
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4 md:gap-6">
          <NuxtLink v-for="p in products.slice(0, 8)" :key="p.id"
            :to="localePath('/products/' + p.slug)"
            class="group bg-white rounded-2xl overflow-hidden border border-[#EAE5DD] hover:shadow-lg transition-all active:scale-[0.98]">
            <div class="aspect-[3/4] bg-gray-100 relative overflow-hidden">
              <img v-if="p.cover_image" :src="p.cover_image" :alt="p.sku"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
                loading="lazy" />
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

// 首页数据改为服务端 SSR 拉取（数据内联）：
//   - SSR 阶段请求后端并输出到 HTML → 提升 SEO（HTML 直接携带真实内容）
//   - 配合 nuxt.config 的 routeRules（/en/** /zh/** ... SWR），HTML 缓存命中即返回完整首页
//   - 客户端 hydration 直接复用，不再二次请求
const api = useApi()
const { data: homeData } = await useAsyncData<any>(
  'home-' + (locale.value || 'en'),
  () => api.getHome().catch(() => null),
)
const products = computed<any[]>(() => homeData.value?.featured_products || [])

const currentSlide = ref(0)
const paused = ref(false)
const touchStartX = ref(0)
const touchEndX = ref(0)

// 图片加载失败集合：外链/相对路径失效时隐藏该图，轮播区显示纯渐变背景，避免白屏
const imgFailed = ref<number[]>([])

// 兜底默认轮播图（后台未配置任何轮播图时的四张内置大图）
const DEFAULT_IMAGES = [
  'https://images.unsplash.com/photo-1599901860904-17e6ed7083a0?w=1600&h=900&fit=crop',
  'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=1600&h=900&fit=crop',
  'https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=1600&h=900&fit=crop',
  'https://images.unsplash.com/photo-1577221084712-45b0445d2b00?w=1600&h=900&fit=crop',
]

// 轮播图数据驱动：优先使用后台「页面管理 → 轮播图」配置（/public/home 的 hero_slides），
// 数据缺失或无效时回退到内置默认文案/图片。
// 单张配置中 title/subtitle/按钮缺失时回退到当前语言内置文案。
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

function onImgError(i: number) {
  if (!imgFailed.value.includes(i)) imgFailed.value.push(i)
}
function buttonText(i: number) { return (heroSlides.value[i]?.button_text as string) || t('home.get_quote') }
function buttonTarget(i: number) { return localePath((heroSlides.value[i]?.button_url as string) || '/contact') }

function onTouchStart(e: TouchEvent) {
  touchStartX.value = e.changedTouches[0].screenX
}

function onTouchEnd(e: TouchEvent) {
  touchEndX.value = e.changedTouches[0].screenX
  const diff = touchStartX.value - touchEndX.value
  const threshold = 50
  if (Math.abs(diff) > threshold) {
    if (diff > 0) {
      // 向左滑动 -> 下一张
      currentSlide.value = (currentSlide.value + 1) % heroSlides.length
    } else {
      // 向右滑动 -> 上一张
      currentSlide.value = (currentSlide.value - 1 + heroSlides.length) % heroSlides.length
    }
  }
}

let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  // 纯客户端行为（hero 轮播），SSR 阶段不执行；悬停暂停，路由离开时清理定时器
  timer = setInterval(() => {
    if (paused.value || !heroSlides.length) return
    currentSlide.value = (currentSlide.value + 1) % heroSlides.length
  }, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>