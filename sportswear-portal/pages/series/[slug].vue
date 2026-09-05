<template>
  <div class="min-h-screen bg-[#FBF9F6]">
    <!-- Breadcrumb -->
    <div class="bg-white border-b border-[#EAE5DD]">
      <div class="max-w-7xl mx-auto px-4 lg:px-8 py-4">
        <Breadcrumb
          :items="[
            { name: $t('nav.home'), to: localePath('/') },
            { name: $t('product.products'), to: localePath('/products') },
            { name: seriesName },
          ]"
        />
      </div>
    </div>

    <!-- Header -->
    <section class="bg-gradient-to-br from-gray-900 to-gray-800 text-white py-16">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h1 class="text-3xl md:text-4xl font-bold mb-2">{{ seriesName }}</h1>
        <p class="text-white/70 text-lg">{{ $t('product.products') }} · {{ products.length }}</p>
      </div>
    </section>

    <!-- Products Grid -->
    <div class="max-w-7xl mx-auto px-4 lg:px-8 py-12">
      <div v-if="products.length" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 md:gap-6">
        <NuxtLink
          v-for="p in products"
          :key="p.id"
          :to="localePath('/products/' + p.slug)"
          class="group bg-white rounded-2xl overflow-hidden border border-[#EAE5DD] hover:shadow-lg transition-all active:scale-[0.98]"
        >
          <div class="aspect-[4/5] bg-gray-100 relative overflow-hidden">
            <img v-if="p.cover_image" :src="imgUrl(p.cover_image)" :alt="p.sku"
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700" loading="lazy" />
            <img v-if="hoverImg(p)" :src="hoverImg(p)" :alt="`${p.sku} - detail`"
              class="absolute inset-0 w-full h-full object-cover opacity-0 group-hover:opacity-100 transition-opacity duration-300" loading="lazy" />
            <div v-if="!p.cover_image" class="w-full h-full flex items-center justify-center text-5xl opacity-30">🏋️</div>
            <div class="absolute bottom-3 left-3">
              <span class="px-2.5 py-1 bg-white/95 rounded-full text-[10px] font-semibold">{{ p.type }}</span>
            </div>
          </div>
          <div class="p-3">
            <h3 class="text-xs font-semibold text-[#0D1B2A] line-clamp-1">{{ p.name || p.sku }}</h3>
            <p class="text-[10px] text-gray-500 mt-1">MOQ {{ p.production_moq || 300 }}</p>
          </div>
        </NuxtLink>
      </div>
      <div v-else class="text-center py-20">
        <div class="text-6xl mb-4">📦</div>
        <p class="text-gray-400">{{ $t('product.no_products') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const localePath = useLocalePath()
const api = useApi()
const { locale } = useI18n()

const slug = route.params.slug as string

function asArray(list: any): any[] {
  return Array.isArray(list) ? list : (list?.items || [])
}

// 系列列表（用于标题 / 面包屑）
const { data: seriesData, refresh: refreshSeries } = await useAsyncData<any[]>(
  'series-' + (locale.value || 'en'),
  () => api.getSeries().catch(() => []),
)
const series = computed(() => asArray(seriesData.value).find((s: any) => s.slug === slug))
const seriesName = computed(() => series.value?.name || slug)

// 系列下产品（自包含：内部重新拉系列匹配，避免与系列列表刷新时序耦合）
const { data: productsData, refresh: refreshProducts } = await useAsyncData<any[]>(
  'series-products-' + (locale.value || 'en') + '-' + slug,
  async () => {
    const list = asArray(await api.getSeries().catch(() => []))
    const s = list.find((x: any) => x.slug === slug)
    if (!s?.id) return []
    const res = await api.getProducts({ series_id: s.id, page: 1, pageSize: 100 })
    return Array.isArray(res) ? res : (res?.items || [])
  },
)
const products = computed(() => productsData.value || [])

function hoverImg(p: any): string {
  return galleryImageUrls(p)[1] || ''
}

useSeoHead({
  title: computed(() => (seriesName.value ? `${seriesName.value} - Custom OEM/ODM Sportswear` : 'Series')),
  description: computed(() => `Explore the ${seriesName.value || 'sportswear'} collection. Custom OEM/ODM sportswear with low MOQ for global brands.`),
  ogImage: computed(() => products.value[0]?.cover_image || undefined),
})

// 语言切换后（同一页面组件复用、不重新挂载时），按新语言重新拉取
watch(locale, () => {
  if (!import.meta.client) return
  refreshSeries()
  refreshProducts()
})
</script>
