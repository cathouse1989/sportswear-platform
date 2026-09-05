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
                :src="imgUrl(p.cover_image)"
                :alt="p.sku"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
                loading="lazy"
              />
              <img
                v-if="hoverImg(p)"
                :src="hoverImg(p)"
                :alt="`${p.sku} - detail`"
                class="absolute inset-0 w-full h-full object-cover opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                loading="lazy"
              />
              <div v-if="!p.cover_image" class="w-full h-full flex items-center justify-center text-4xl md:text-6xl">🏋️</div>
              <div class="absolute bottom-3 left-3 md:bottom-4 md:left-4">
                <span class="px-2.5 md:px-3 py-1 bg-white/95 rounded-full text-[10px] md:text-xs font-semibold">{{ localizedEnum('product.type_options', p.type) }}</span>
              </div>
            </div>
            <div class="p-3 md:p-5">
              <h3 class="text-xs md:text-sm font-semibold text-[#0D1B2A] line-clamp-1">{{ p.name || p.sku }}</h3>
              <p class="text-[10px] md:text-sm text-gray-500 mt-0.5 md:mt-1 line-clamp-2">{{ p.brief?.slice(0, 60) }}</p>
              <p class="text-[10px] md:text-xs text-gray-400 mt-1 md:mt-2">MOQ {{ p.production_moq || 300 }}</p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- Blog：首页博客信息区块 -->
    <section v-if="blogs?.length" class="py-16 md:py-20 bg-white border-t border-[#EAE5DD]">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <div class="flex justify-between items-end mb-8 md:mb-10">
          <h2 class="text-2xl md:text-3xl font-bold text-[#0D1B2A]">{{ $t('blog.title') }}</h2>
          <NuxtLink :to="localePath('/blog')" class="text-sm text-[#D4A853] font-medium min-h-[44px] flex items-center">{{ $t('blog.read_more') }} →</NuxtLink>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6">
          <NuxtLink
            v-for="b in blogs.slice(0, 3)"
            :key="b.id"
            :to="localePath('/blog/' + b.slug)"
            class="group bg-white rounded-2xl overflow-hidden border border-[#EAE5DD] hover:shadow-lg transition-all active:scale-[0.98]"
          >
            <div class="aspect-[16/9] bg-gray-100 relative overflow-hidden">
              <img
                v-if="b.cover_image"
                :src="imgUrl(b.cover_image)"
                :alt="b.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
                loading="lazy"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-4xl md:text-5xl">📝</div>
            </div>
            <div class="p-4 md:p-5">
              <div class="flex items-center gap-2 text-[10px] md:text-xs text-gray-400 mb-2">
                <span class="px-2 py-0.5 rounded-full bg-[#F5F0E8] text-[#9A6B1F] font-medium">{{ categoryLabel(b.category) }}</span>
                <time>{{ formatDate(b.published_at || b.created_at) }}</time>
              </div>
              <h3 class="text-sm md:text-base font-semibold text-[#0D1B2A] line-clamp-2 group-hover:text-[#D4A853] transition-colors">{{ b.title }}</h3>
              <p class="text-xs md:text-sm text-gray-500 mt-2 line-clamp-2">{{ b.summary || stripHtml(b.content) }}</p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const localePath = useLocalePath()
const { t, te, locale } = useI18n()
const { localizedEnum } = useLocalized()
const api = useApi()

// SEO - 首页
const homeSeo = await useRouteSeo('home')
useSeoHead({
  title: '', // 首页不需要额外 title，composable 会自动用 SITE_NAME
  description: () =>
    homeSeo.value?.description ||
    'Leading OEM/ODM sportswear manufacturer specializing in custom activewear, athletic apparel, and team uniforms. 15+ years of experience serving global brands with ISO-certified quality.',
  keywords: () =>
    homeSeo.value?.keywords ||
    'sportswear manufacturer, OEM sportswear, ODM activewear, custom athletic apparel, China sportswear factory, private label activewear',
  ogImage: '/images/og-home.jpg',
})

// 首页数据（SSR 内联到 HTML）
const { data: homeData, refresh: refreshHome } = await useAsyncData<any>(
  'home-' + (locale.value || 'en'),
  () => api.getHome().catch(() => null),
  {
    getCachedData(key, nuxtApp) {
      return nuxtApp.isHydrating ? nuxtApp.payload.data[key] : undefined
    },
  },
)

const { data: themeData } = await useAsyncData<Record<string, any>>(
  'theme-' + (locale.value || 'en'),
  () => api.getTheme().catch(() => ({})),
  {
    getCachedData(key, nuxtApp) {
      return nuxtApp.isHydrating ? nuxtApp.payload.data[key] : undefined
    },
  },
)

const products = computed<any[]>(() => homeData.value?.featured_products || [])
const blogs = computed<any[]>(() => homeData.value?.blogs || [])

// 语言切换后（同一页面组件复用、不重新挂载时），显式按新语言重新拉取首页聚合数据
watch(locale, () => {
  if (!import.meta.client) return
  refreshHome()
})

// 精选产品卡片 hover 第二张图
function hoverImg(p: any): string {
  return galleryImageUrls(p)[1] || ''
}

// 本地化博客分类标签：优先词条，回退原始值
function categoryLabel(value?: string) {
  if (!value) return ''
  const key = `blog.categories.${value}`
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

// 去除 HTML 标签生成摘要
function stripHtml(html?: string) {
  return (html || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').replace(/\s+/g, ' ').trim()
}

// 轮播文案词条字典（来自后台「词条管理」i18n_entries，按语言翻页配置）
// Key：home.hero_title_1..N / home.hero_sub_1..N / home.get_quote
const { data: heroDict } = await useAsyncData<Record<string, string>>(
  'hero-i18n-' + (locale.value || 'en'),
  async () => {
    const res = await api.getI18n().catch(() => null)
    return res?.dictionary || {}
  },
  {
    getCachedData(key, nuxtApp) {
      return nuxtApp.isHydrating ? nuxtApp.payload.data[key] : undefined
    },
  },
)

// 本地英雄轮播图片（存储于 public/images/hero/）
// 使用本地图片而非外部 URL，提升加载速度与可靠性，并便于版本管理
// 如需替换图片：1）放入新图片到 public/images/hero/ 2）更新此处路径即可
// 原始图片来源：Unsplash（免费商业用途，详情见 public/images/hero/README.md）
const DEFAULT_IMAGES = [
  '/images/hero/slide-1-sportswear-manufacturing.jpg',
  '/images/hero/slide-2-running-training.jpg',
  '/images/hero/slide-3-team-uniforms.jpg',
  '/images/hero/slide-4-concept-to-product.jpg',
]

// 从词条字典取文案；未配置时回退 slides 存量文本 / 静态语言包
// 换行符统一为 \n，模板配合 whitespace-pre-line 渲染多行
const pickDict = (key: string, fallback: string) => {
  const v = heroDict.value?.[key]
  return typeof v === 'string' && v.trim()
    ? v.replace(/\r\n?/g, '\n').trim()
    : fallback
}

// 优先后台 hero_slides 图片；文本按当前语言取自（词条字典 → slides 存量 → 语言包）
const heroSlides = computed(() => {
  const raw = homeData.value?.hero_slides
  if (Array.isArray(raw) && raw.length) {
    return raw.map((s: any, i: number) => ({
      image: String(s?.image || DEFAULT_IMAGES[i] || ''),
      title: pickDict(`home.hero_title_${i + 1}`, String(s?.title || t(`home.hero_title_${i + 1}`))),
      subtitle: pickDict(`home.hero_sub_${i + 1}`, String(s?.subtitle || t(`home.hero_sub_${i + 1}`))),
      button_text: pickDict('home.get_quote', String(s?.button_text || t('home.get_quote'))),
      button_url: String(s?.button_url || '/contact'),
    }))
  }
  return [0, 1, 2, 3].map((i) => ({
    image: DEFAULT_IMAGES[i],
    title: pickDict(`home.hero_title_${i + 1}`, t(`home.hero_title_${i + 1}`)),
    subtitle: pickDict(`home.hero_sub_${i + 1}`, t(`home.hero_sub_${i + 1}`)),
    button_text: pickDict('home.get_quote', t('home.get_quote')),
    button_url: '/contact',
  }))
})

/** 主题兜底值：theme 表 hero_* 配置 → 数据库不存在时回退硬编码默认值 */
const themeHero = computed(() => {
  const t = themeData.value || {}
  const toBool = (key: string, def: boolean) => {
    if (!(key in t)) return def
    const v = t[key]
    if (typeof v === 'boolean') return v
    if (typeof v === 'string') return v.toLowerCase() === 'true'
    return def
  }
  const toNum = (key: string, def: number) => {
    if (!(key in t)) return def
    const n = Number(t[key])
    return Number.isFinite(n) ? n : def
  }
  const toString = (key: string, def: string) => (typeof t[key] === 'string' ? t[key] : def)
  return {
    autoplay: toBool('hero_autoplay', true),
    interval_ms: toNum('hero_interval_ms', 5000),
    transition: toString('hero_transition', 'fade'),
    show_dots: toBool('hero_show_dots', true),
    pause_on_hover: toBool('hero_pause_on_hover', true),
  }
})

// 展示参数：banner 模块 settings > theme 配置 > 硬编码默认
const heroSettings = computed(() => {
  const s = homeData.value?.hero_settings || {}
  const th = themeHero.value
  const pick = (v: unknown, fallback: boolean) => {
    if (v === undefined || v === null) return fallback
    if (v === false || v === 0 || v === '0') return false
    if (typeof v === 'string' && v.toLowerCase() === 'false') return false
    return true
  }
  return {
    autoplay: pick(s.autoplay, th.autoplay),
    interval_ms: Number(s.interval_ms) || th.interval_ms || 5000,
    transition: (s.transition === 'slide' || s.transition === 'fade') ? s.transition : th.transition,
    show_dots: pick(s.show_dots, th.show_dots),
    show_arrows: true,
    pause_on_hover: pick(s.pause_on_hover, th.pause_on_hover),
  }
})
</script>
