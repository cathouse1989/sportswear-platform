<template>
  <section
    class="relative min-h-screen flex items-center bg-[#0D1B2A] overflow-hidden touch-pan-y"
    @touchstart="onTouchStart"
    @touchend="onTouchEnd"
    @mouseenter="onMouseEnter"
    @mouseleave="onMouseLeave"
  >
    <div v-for="(s, i) in slides" :key="i" class="absolute inset-0" :class="slideClass(i)">
      <img
        v-if="s.image && !imgFailed[i]"
        :src="s.image"
        :alt="s.title"
        class="w-full h-full object-cover"
        :fetchpriority="i === 0 ? 'high' : 'auto'"
        :loading="i === 0 ? 'eager' : 'lazy'"
        @error="onImgError(i)"
      />
      <div class="absolute inset-0 bg-gradient-to-r from-[#0D1B2A]/95 via-[#0D1B2A]/70 to-transparent" />
    </div>

    <div class="relative max-w-7xl mx-auto px-4 lg:px-8 py-20 w-full">
      <div class="max-w-2xl">
        <div v-for="(s, i) in slides" :key="'c-' + i" :class="current === i ? '' : 'hidden'">
          <h1 class="text-4xl sm:text-5xl md:text-7xl font-bold text-white mb-4 md:mb-6 leading-tight">{{ s.title }}</h1>
          <p v-if="s.subtitle" class="text-base sm:text-lg md:text-xl text-white/60 mb-8 md:mb-10">{{ s.subtitle }}</p>
          <NuxtLink
            :to="localePath(s.button_url || '/contact')"
            class="inline-flex items-center gap-3 bg-[#D4A853] text-white px-8 md:px-10 py-3.5 md:py-4 rounded-full font-semibold hover:bg-[#C49A3F] transition text-sm md:text-base min-h-[48px]"
          >{{ s.button_text || defaultButtonText }}</NuxtLink>
        </div>
      </div>
    </div>

    <template v-if="showArrows && slides.length > 1">
      <button type="button" class="absolute left-4 top-1/2 -translate-y-1/2 w-11 h-11 rounded-full bg-white/15 hover:bg-white/30 text-white items-center justify-center transition hidden md:flex" aria-label="Previous slide" @click="goPrev">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
      </button>
      <button type="button" class="absolute right-4 top-1/2 -translate-y-1/2 w-11 h-11 rounded-full bg-white/15 hover:bg-white/30 text-white items-center justify-center transition hidden md:flex" aria-label="Next slide" @click="goNext">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
      </button>
    </template>

    <div v-if="showDots && slides.length > 1" class="absolute bottom-12 left-1/2 -translate-x-1/2 flex gap-3">
      <button
        v-for="(s, i) in slides"
        :key="'d-' + i"
        type="button"
        class="w-2.5 h-2.5 rounded-full transition-all min-h-[10px]"
        :class="current === i ? 'bg-[#D4A853] w-8' : 'bg-white/30'"
        :aria-label="`Go to slide ${i + 1}`"
        @click="goTo(i)"
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

const autoplay = computed(() => props.settings?.autoplay !== false)
const intervalMs = computed(() => {
  const n = Number(props.settings?.interval_ms || 5000)
  if (!Number.isFinite(n) || n < 2000) return 5000
  if (n > 30000) return 30000
  return n
})
const transition = computed(() => (props.settings?.transition === 'slide' ? 'slide' : 'fade'))
const showDots = computed(() => props.settings?.show_dots !== false)
const showArrows = computed(() => props.settings?.show_arrows === true)
const pauseOnHover = computed(() => props.settings?.pause_on_hover !== false)
const slideCount = computed(() => props.slides?.length || 0)

function slideClass(i: number) {
  if (transition.value === 'slide') {
    return [
      'transition-transform duration-700 ease-out',
      current.value === i ? 'translate-x-0 opacity-100 z-10' : 'translate-x-full opacity-0 z-0',
    ]
  }
  return [
    'transition-opacity duration-1000',
    current.value === i ? 'opacity-100 z-10' : 'opacity-0 z-0',
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
  if (!autoplay.value || slideCount.value <= 1) return
  if (window.matchMedia?.('(prefers-reduced-motion: reduce)').matches) return
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
onMounted(startTimer)
onUnmounted(clearTimer)

</script>
