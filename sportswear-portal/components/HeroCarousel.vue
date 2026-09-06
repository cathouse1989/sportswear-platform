<template>
  <section
    class="hero-carousel relative min-h-screen overflow-hidden bg-[#0D1B2A] touch-pan-y"
    @touchstart="onTouchStart"
    @touchend="onTouchEnd"
    @mouseenter="onMouseEnter"
    @mouseleave="onMouseLeave"
  >
    <!-- 背景：z-0，不参与点击（点击跳转统一由文案区的 CTA「获取报价」按钮承担，避免全屏可点击背景抢占按钮点击区域） -->
    <div
      v-for="(s, i) in slides"
      :key="'bg-' + i"
      class="absolute inset-0 z-0 pointer-events-none"
      :class="slideClass(i)"
    >
      <img
        v-if="s.image && !imgFailed[i]"
        :src="s.image"
        :alt="s.title || ''"
        class="h-full w-full object-cover"
        :fetchpriority="i === 0 ? 'high' : 'auto'"
        :loading="i === 0 ? 'eager' : 'lazy'"
        @error="onImgError(i)"
      />
      <div class="absolute inset-0 bg-gradient-to-r from-[#0D1B2A]/95 via-[#0D1B2A]/70 to-transparent pointer-events-none" />
    </div>

    <!-- 文案：与箭头错开左右留白，避免盖住切图按钮 -->
    <div
      class="relative z-10 flex min-h-screen w-full items-center pointer-events-none"
      :class="hasNav ? 'px-14 sm:px-16 lg:px-24' : 'px-4 lg:px-8'"
    >
      <div class="mx-auto w-full max-w-7xl py-20">
        <div class="max-w-2xl">
          <div v-for="(s, i) in slides" :key="'c-' + i" :class="current === i ? '' : 'hidden'">
            <h1 class="mb-4 whitespace-pre-line text-4xl font-bold leading-tight text-white sm:text-5xl md:mb-6 md:text-7xl">
              {{ s.title }}
            </h1>
            <p v-if="s.subtitle" class="mb-8 whitespace-pre-line text-base text-white/60 sm:text-lg md:mb-10 md:text-xl">
              {{ s.subtitle }}
            </p>
            <NuxtLink
              :to="localePath(s.button_url || '/contact')"
              class="pointer-events-auto inline-flex min-h-[48px] items-center gap-3 rounded-full bg-[#D4A853] px-8 py-3.5 text-sm font-semibold text-white transition hover:bg-[#C49A3F] md:px-10 md:py-4 md:text-base"
            >{{ s.button_text || defaultButtonText }}</NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <!-- 左右箭头：Tailwind 实心样式保证构建产物可见；桌面右箭头避开 FloatingContact -->
    <div v-if="hasNav" class="pointer-events-none absolute inset-0 z-30">
      <button
        type="button"
        class="hero-arrow pointer-events-auto absolute left-2 top-1/2 z-30 flex h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full border-2 border-white bg-[#0D1B2A] text-white shadow-[0_8px_24px_rgba(0,0,0,0.45)] transition hover:bg-[#D4A853] hover:border-[#D4A853] sm:left-4 md:h-14 md:w-14"
        aria-label="Previous slide"
        @click.stop="goPrev"
      >
        <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <button
        type="button"
        class="hero-arrow pointer-events-auto absolute right-2 top-1/2 z-30 flex h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full border-2 border-white bg-[#0D1B2A] text-white shadow-[0_8px_24px_rgba(0,0,0,0.45)] transition hover:bg-[#D4A853] hover:border-[#D4A853] sm:right-4 md:right-16 md:h-14 md:w-14"
        aria-label="Next slide"
        @click.stop="goNext"
      >
        <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>

    <!-- 指示点 -->
    <div
      v-if="showDots && slideCount > 1"
      class="pointer-events-auto absolute bottom-12 left-1/2 z-30 flex -translate-x-1/2 gap-3"
    >
      <button
        v-for="(s, i) in slides"
        :key="'d-' + i"
        type="button"
        class="min-h-[10px] min-w-[10px] rounded-full transition-all"
        :class="current === i ? 'h-2.5 w-8 bg-[#D4A853]' : 'h-2.5 w-2.5 bg-white/40 hover:bg-white/70'"
        :aria-label="`Go to slide ${i + 1}`"
        @click.stop="goTo(i)"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
export interface HeroSlideItem {
  image?: string
  title?: string
  subtitle?: string
  button_text?: string
  button_url?: string
}

export interface HeroSettingsInput {
  autoplay?: boolean
  interval_ms?: number
  transition?: string
  show_dots?: boolean
  show_arrows?: boolean
  pause_on_hover?: boolean
}

const props = withDefaults(
  defineProps<{
    slides: HeroSlideItem[]
    settings?: HeroSettingsInput | null
    defaultButtonText?: string
  }>(),
  { settings: null, defaultButtonText: 'Get a Quote' },
)

const localePath = useLocalePath()
const current = ref(0)
const paused = ref(false)
const imgFailed = ref<number[]>([])
const touchStartX = ref(0)
const touchEndX = ref(0)

/** 后台/主题开关：仅显式 false / "false" / 0 时关闭，缺省视为开启 */
function flagOn(v: unknown, defaultOn = true) {
  if (v === undefined || v === null || v === '') return defaultOn
  if (v === false || v === 0 || v === '0') return false
  if (typeof v === 'string' && v.toLowerCase() === 'false') return false
  return true
}

const autoplay = computed(() => flagOn(props.settings?.autoplay, true))
const intervalMs = computed(() => {
  const n = Number(props.settings?.interval_ms || 5000)
  if (!Number.isFinite(n) || n < 2000) return 5000
  if (n > 30000) return 30000
  return n
})
const transition = computed(() => (props.settings?.transition === 'slide' ? 'slide' : 'fade'))
const showDots = computed(() => flagOn(props.settings?.show_dots, true))
const pauseOnHover = computed(() => flagOn(props.settings?.pause_on_hover, true))
const slideCount = computed(() => props.slides?.length || 0)
const hasNav = computed(() => slideCount.value > 1)

function slideClass(i: number) {
  if (transition.value === 'slide') {
    const n = slideCount.value
    if (n <= 1) return ['translate-x-0 opacity-100']
    const offset = ((i - current.value) % n + n) % n
    if (offset === 0) return ['z-10 translate-x-0 opacity-100 transition-transform duration-700 ease-out']
    if (offset === n - 1) return ['z-0 -translate-x-full opacity-0 transition-transform duration-700 ease-out']
    return ['z-0 translate-x-full opacity-0 transition-transform duration-700 ease-out']
  }
  return [
    'transition-opacity duration-1000',
    current.value === i ? 'opacity-100' : 'opacity-0',
  ]
}

function onImgError(i: number) {
  if (!imgFailed.value.includes(i)) imgFailed.value.push(i)
}
function goTo(i: number) {
  if (!slideCount.value) return
  current.value = ((i % slideCount.value) + slideCount.value) % slideCount.value
}
function goNext() {
  if (!slideCount.value) return
  current.value = (current.value + 1) % slideCount.value
}
function goPrev() {
  if (!slideCount.value) return
  current.value = (current.value - 1 + slideCount.value) % slideCount.value
}
function onTouchStart(e: TouchEvent) {
  touchStartX.value = e.changedTouches[0]?.screenX || 0
}
function onTouchEnd(e: TouchEvent) {
  touchEndX.value = e.changedTouches[0]?.screenX || 0
  const diff = touchStartX.value - touchEndX.value
  if (Math.abs(diff) > 50) {
    if (diff > 0) goNext()
    else goPrev()
  }
}
function onMouseEnter() {
  if (pauseOnHover.value) paused.value = true
}
function onMouseLeave() {
  paused.value = false
}

let timer: ReturnType<typeof setInterval> | null = null
function clearTimer() {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}
function startTimer() {
  clearTimer()
  if (!import.meta.client) return
  if (slideCount.value <= 1) return
  if (window.matchMedia?.('(prefers-reduced-motion: reduce)').matches) return
  if (!autoplay.value) return
  timer = setInterval(() => {
    if (paused.value) return
    goNext()
  }, intervalMs.value)
}

watch(
  () => [autoplay.value, intervalMs.value, slideCount.value],
  () => {
    if (current.value >= slideCount.value) current.value = 0
    startTimer()
  },
)

let mounted = false
onMounted(() => {
  mounted = true
  startTimer()
})
watchEffect(() => {
  if (mounted && slideCount.value > 1 && autoplay.value && !timer) {
    startTimer()
  }
})
onUnmounted(() => {
  mounted = false
  clearTimer()
})
</script>

<style scoped>
.hero-arrow {
  padding: 0;
  margin: 0;
  appearance: none;
  -webkit-appearance: none;
  cursor: pointer;
}
</style>
