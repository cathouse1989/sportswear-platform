<template>
  <div class="min-h-screen bg-[#FBF9F6]">
    <!-- Breadcrumb -->
    <div class="bg-white border-b border-[#EAE5DD]">
      <div class="max-w-7xl mx-auto px-4 lg:px-8 py-4">
        <Breadcrumb
          :items="[
            { name: $t('nav.home'), to: localePath('/') },
            { name: $t('product.products'), to: localePath('/products') },
            { name: categoryName },
          ]"
        />
      </div>
    </div>

    <!-- Header -->
    <section class="bg-gradient-to-br from-gray-900 to-gray-800 text-white py-16">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h1 class="text-3xl md:text-4xl font-bold mb-2">{{ categoryName }}</h1>
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
const { localizedCategory } = useLocalized()

const slug = route.params.slug as string

function findCategory(list: any[] | undefined, s: string): any {
  for (const c of list || []) {
    if (c.slug === s) return c
    if (c.children?.length) {
      const found = findCategory(c.children, s)
      if (found) return found
    }
  }
  return null
}

// 分类树（用于标题 / 面包屑）
const { data: catData, refresh: refreshCategories } = await useAsyncData<any[]>(
  'categories-' + (locale.value || 'en'),
  () => api.getCategories().catch(() => []),
)
const category = computed(() => findCategory(catData.value, slug))
const categoryName = computed(() => localizedCategory(category.value?.name || slug, slug))

// 分类下产品（自包含：内部重新拉分类匹配，避免与分类树刷新时序耦合）
const { data: productsData, refresh: refreshProducts } = await useAsyncData<any[]>(
  'cat-products-' + (locale.value || 'en') + '-' + slug,
  async () => {
    const cats = await api.getCategories().catch(() => [])
    const c = findCategory(cats, slug)
    if (!c?.id) return []
    const res = await api.getProducts({ category_id: c.id, page: 1, pageSize: 100 })
    return Array.isArray(res) ? res : (res?.items || [])
  },
)
const products = computed(() => productsData.value || [])

function hoverImg(p: any): string {
  return galleryImageUrls(p)[1] || ''
}

useSeoHead({
  title: computed(() => (categoryName.value ? `${categoryName.value} - Custom OEM/ODM Sportswear` : 'Category')),
  description: computed(() => `Browse ${categoryName.value || 'sportswear'} custom OEM/ODM sportswear products with low MOQ. Professional manufacturing for global brands.`),
  ogImage: computed(() => products.value[0]?.cover_image || undefined),
})

// 语言切换后（同一页面组件复用、不重新挂载时），按新语言重新拉取
watch(locale, () => {
  if (!import.meta.client) return
  refreshCategories()
  refreshProducts()
})
</script>
