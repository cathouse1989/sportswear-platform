<template>
  <!-- 统一 Logo 组件：支持默认 SVG / 自定义图片 URL 两种模式 -->
  <div class="flex items-center gap-2.5 group">
    <!-- 自定义图片 Logo（优先） -->
    <template v-if="logoSrc">
      <img
        :src="logoSrc"
        :alt="logoAlt || 'Sportswear'"
        class="h-9 w-auto rounded-xl shadow-sm group-hover:shadow-md transition-all duration-300"
      />
      <div v-if="!hideText" class="leading-none select-none">
        <div
          class="text-sm font-bold tracking-tight"
          :class="variant === 'dark' ? 'text-white' : 'text-[#0D1B2A]'"
        >{{ brandName || 'SPORTSWEAR' }}</div>
        <div
          class="text-[9px] uppercase tracking-[0.3em]"
          :class="variant === 'dark' ? 'text-[#9A8C7A]' : 'text-[#9A8C7A]'"
        >{{ brandSubtitle || 'Premium Mfg.' }}</div>
      </div>
    </template>

    <!-- 默认 SVG Logo（无自定义图片时） -->
    <template v-else>
      <div
        class="w-9 h-9 rounded-xl flex items-center justify-center shadow-sm group-hover:shadow-md transition-all duration-300"
        :class="variant === 'dark' ? 'bg-[#D4A853]' : 'bg-[#0D1B2A]'"
      >
        <svg viewBox="0 0 36 36" class="w-5.5 h-5.5" fill="none">
          <path
            d="M8 24c0-2.5 1.5-4.5 4-5.5l-1-1C9.5 16 8 13.5 8 11c0-3 2.5-5.5 6-5.5 3 0 5 2 5.5 4.5l-3 1C16 9.5 15 8.5 14 8.5c-1.5 0-2.5 1-2.5 2.5 0 1.5 1 2.5 2.5 3.5l4 2c2.5 1 4 3 4 5.5 0 3-2.5 5.5-6 5.5-3.5 0-6-2-6.5-5l3-1c0.5 1.5 1.5 3 3.5 3 1.5 0 2.5-1 2.5-2.5z"
            :class="variant === 'dark' ? 'fill-[#0D1B2A]' : 'fill-[#D4A853]'"
            class="transition-colors duration-300"
          />
          <path
            d="M22 9l4-1M24 12l5-2"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
            :class="variant === 'dark' ? 'text-[#0D1B2A]/40' : 'text-[#D4A853]/50'"
          />
        </svg>
      </div>
      <div v-if="!hideText" class="leading-none select-none">
        <div
          class="text-sm font-bold tracking-tight"
          :class="variant === 'dark' ? 'text-white' : 'text-[#0D1B2A]'"
        >{{ brandName || 'SPORTSWEAR' }}</div>
        <div
          class="text-[9px] uppercase tracking-[0.3em]"
          :class="variant === 'dark' ? 'text-[#9A8C7A]' : 'text-[#9A8C7A]'"
        >{{ brandSubtitle || 'Premium Mfg.' }}</div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  variant?: 'light' | 'dark'
  logoSrc?: string
  logoAlt?: string
  brandName?: string
  brandSubtitle?: string
  hideText?: boolean
}>()
</script>