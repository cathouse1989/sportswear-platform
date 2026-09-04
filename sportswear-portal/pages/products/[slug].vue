<template>
  <div v-if="product" class="min-h-screen bg-[#FBF9F6]">
    <!-- Breadcrumb -->
    <div class="bg-white border-b border-[#EAE5DD]">
      <div class="max-w-7xl mx-auto px-4 lg:px-8 py-4">
        <Breadcrumb
          :items="[
            { name: $t('nav.home'), to: localePath('/') },
            { name: $t('product.products'), to: localePath('/products') },
            { name: product.name || product.sku },
          ]"
        />
      </div>
    </div>

    <!-- Product Main Section -->
    <div class="max-w-7xl mx-auto px-4 lg:px-8 py-8 lg:py-12">
      <div class="grid lg:grid-cols-2 gap-8 lg:gap-12">
        <!-- Left: Image Gallery -->
        <div>
          <div class="relative bg-white rounded-2xl overflow-hidden border border-[#EAE5DD] aspect-[4/5] mb-4 cursor-zoom-in" @click="openLightbox(currentImageIndex)">
            <img
              v-if="currentImage"
              :src="currentImage"
              :alt="`${product.sku || product.name} - Sportswear Product`"
              class="w-full h-full object-cover"
              fetchpriority="high"
              loading="eager"
            />
            <div v-else class="w-full h-full flex items-center justify-center text-8xl text-gray-200">📷</div>
            <!-- Image counter -->
            <div v-if="allImages.length > 1" class="absolute bottom-4 left-1/2 -translate-x-1/2 bg-black/60 text-white text-xs px-3 py-1 rounded-full">
              {{ currentImageIndex + 1 }} / {{ allImages.length }}
            </div>
            <!-- Prev/Next -->
            <button v-if="allImages.length > 1" @click.stop="prevImage" class="absolute left-3 top-1/2 -translate-y-1/2 w-10 h-10 bg-white/80 rounded-full flex items-center justify-center shadow hover:bg-white transition">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
            </button>
            <button v-if="allImages.length > 1" @click.stop="nextImage" class="absolute right-3 top-1/2 -translate-y-1/2 w-10 h-10 bg-white/80 rounded-full flex items-center justify-center shadow hover:bg-white transition">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
            </button>
          </div>
          <!-- Thumbnails -->
          <div v-if="allImages.length > 1" class="flex gap-3 overflow-x-auto pb-2 scrollbar-hide">
            <button
              v-for="(img, i) in allImages"
              :key="i"
              @click="currentImageIndex = i"
              class="shrink-0 w-20 h-20 rounded-xl overflow-hidden border-2 transition"
              :class="currentImageIndex === i ? 'border-[#D4A853]' : 'border-transparent hover:border-gray-300'"
            >
              <img :src="img" :alt="`${product.sku || product.name} thumbnail ${i + 1}`" class="w-full h-full object-cover" loading="lazy" />
            </button>
          </div>
        </div>

        <!-- Right: Product Info -->
        <div>
          <!-- Title & Actions -->
          <div class="flex items-start justify-between gap-4 mb-6">
            <div>
              <h1 class="text-3xl lg:text-4xl font-bold text-[#0D1B2A] mb-2">{{ product.name || product.sku }}</h1>
              <p class="text-lg text-gray-500">{{ product.sku }}</p>
            </div>
            <div class="flex gap-2 shrink-0">
              <!-- Copy Button -->
              <button @click="copyProductInfo" class="flex items-center gap-2 px-4 py-2.5 border border-[#EAE5DD] rounded-xl text-sm text-gray-600 hover:bg-white hover:border-gray-300 transition min-h-[44px]" :title="$t('product.detail.copy_info')">
                <svg v-if="!copied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
                <svg v-else class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
                <span>{{ copied ? $t('product.detail.copied') : $t('product.detail.copy_info') }}</span>
              </button>
            </div>
          </div>

          <!-- Brief -->
          <p v-if="product.brief" class="text-gray-600 mb-6 leading-relaxed">{{ product.brief }}</p>

          <!-- Key Specs Grid -->
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 mb-8">
            <div v-if="product.type" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.type') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.type }}</div>
            </div>
            <div v-if="product.gender" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.gender') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.gender }}</div>
            </div>
            <div v-if="product.material" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.material') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.material }}</div>
            </div>
            <div v-if="product.composition" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.composition') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.composition }}</div>
            </div>
            <div v-if="product.weight" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.weight') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.weight }}</div>
            </div>
            <div v-if="product.elasticity" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.elasticity') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.elasticity }}</div>
            </div>
            <div v-if="product.fit" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.fit') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.fit }}</div>
            </div>
            <div v-if="product.support_level" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.support_level') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.support_level }}</div>
            </div>
            <div v-if="product.season" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.season') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.season }}</div>
            </div>
            <div v-if="product.size_range" class="bg-white rounded-xl p-4 border border-[#EAE5DD]">
              <div class="text-[10px] uppercase tracking-wider text-gray-400 mb-1">{{ $t('product.detail.size_range') }}</div>
              <div class="font-semibold text-sm text-[#0D1B2A]">{{ product.size_range }}</div>
            </div>
          </div>

          <!-- MOQ -->
          <div class="bg-[#0D1B2A]/5 rounded-xl p-5 mb-8">
            <h3 class="text-sm font-semibold text-[#0D1B2A] mb-3">{{ $t('product.detail.moq') }}</h3>
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
              <div>
                <div class="text-xs text-gray-500">{{ $t('product.detail.sample_moq') }}</div>
                <div class="text-lg font-bold text-[#0D1B2A]">{{ product.sample_moq || 1 }}</div>
              </div>
              <div>
                <div class="text-xs text-gray-500">{{ $t('product.detail.production_moq') }}</div>
                <div class="text-lg font-bold text-[#0D1B2A]">{{ product.production_moq || 300 }}</div>
              </div>
              <div>
                <div class="text-xs text-gray-500">{{ $t('product.detail.color_moq') }}</div>
                <div class="text-lg font-bold text-[#0D1B2A]">{{ product.color_moq || 100 }}</div>
              </div>
              <div>
                <div class="text-xs text-gray-500">{{ $t('product.detail.size_moq') }}</div>
                <div class="text-lg font-bold text-[#0D1B2A]">{{ product.size_moq || 100 }}</div>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex flex-col sm:flex-row gap-3 mb-8">
            <NuxtLink
              :to="localePath('/contact') + '?product=' + encodeURIComponent(product.sku)"
              class="flex-1 inline-flex items-center justify-center gap-2 bg-[#D4A853] text-white px-8 py-3.5 rounded-xl font-semibold hover:bg-[#C49A3F] transition text-sm min-h-[48px]"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
              {{ $t('product.detail.inquiry_product') }}
            </NuxtLink>
            <a
              :href="`https://wa.me/${waNumber}?text=${encodeURIComponent(waMessage)}`"
              target="_blank"
              rel="noopener noreferrer"
              class="flex-1 inline-flex items-center justify-center gap-2 bg-[#25D366] text-white px-8 py-3.5 rounded-xl font-semibold hover:bg-[#20BD5A] transition text-sm min-h-[48px]"
            >
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z" /></svg>
              WhatsApp
            </a>
          </div>
        </div>
      </div>

      <!-- Description & Features -->
      <div class="grid lg:grid-cols-2 gap-8 lg:gap-12 mt-12">
        <div v-if="product.description" class="bg-white rounded-2xl p-6 lg:p-8 border border-[#EAE5DD]">
          <h2 class="text-xl font-bold text-[#0D1B2A] mb-4">{{ $t('product.detail.description') }}</h2>
          <div class="prose prose-sm max-w-none text-gray-600 leading-relaxed whitespace-pre-wrap">{{ product.description }}</div>
        </div>
        <div v-if="product.features" class="bg-white rounded-2xl p-6 lg:p-8 border border-[#EAE5DD]">
          <h2 class="text-xl font-bold text-[#0D1B2A] mb-4">{{ $t('product.detail.features') }}</h2>
          <div class="prose prose-sm max-w-none text-gray-600 leading-relaxed whitespace-pre-wrap">{{ product.features }}</div>
        </div>
        <div v-if="product.usage" class="bg-white rounded-2xl p-6 lg:p-8 border border-[#EAE5DD]">
          <h2 class="text-xl font-bold text-[#0D1B2A] mb-4">{{ $t('product.detail.usage') }}</h2>
          <div class="prose prose-sm max-w-none text-gray-600 leading-relaxed whitespace-pre-wrap">{{ product.usage }}</div>
        </div>
      </div>

      <!-- Specs Table -->
      <div v-if="product.specs?.length" class="mt-8 lg:mt-12">
        <h2 class="text-xl font-bold text-[#0D1B2A] mb-6">{{ $t('product.detail.specifications') }}</h2>
        <div class="bg-white rounded-2xl border border-[#EAE5DD] overflow-hidden">
          <table class="w-full">
            <tbody>
              <tr v-for="(spec, i) in product.specs" :key="spec.id || i" class="border-b border-[#EAE5DD] last:border-b-0">
                <td class="px-6 py-4 text-sm font-medium text-gray-500 bg-[#FBF9F6] w-1/3">{{ spec.name }}</td>
                <td class="px-6 py-4 text-sm text-[#0D1B2A]">{{ spec.value }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Customizations -->
      <div v-if="product.customizations?.length" class="mt-8 lg:mt-12">
        <h2 class="text-xl font-bold text-[#0D1B2A] mb-6">{{ $t('product.detail.customizations') }}</h2>
        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="(c, i) in product.customizations" :key="c.id || i" class="bg-white rounded-xl p-5 border border-[#EAE5DD] flex items-start gap-3">
            <div class="w-10 h-10 bg-[#D4A853]/10 rounded-lg flex items-center justify-center shrink-0">
              <svg class="w-5 h-5 text-[#D4A853]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
            </div>
            <div>
              <div class="font-semibold text-sm text-[#0D1B2A] capitalize">{{ c.type }}</div>
              <div v-if="c.note" class="text-xs text-gray-500 mt-1">{{ c.note }}</div>
              <div v-if="c.is_enabled" class="text-xs text-green-600 mt-1 font-medium">✓ Available</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Series & Fabrics -->
      <div v-if="product.series?.length || product.fabrics?.length" class="mt-8 lg:mt-12 grid lg:grid-cols-2 gap-8">
        <div v-if="product.series?.length" class="bg-white rounded-2xl p-6 lg:p-8 border border-[#EAE5DD]">
          <h2 class="text-xl font-bold text-[#0D1B2A] mb-4">{{ $t('product.detail.series') }}</h2>
          <div class="flex flex-wrap gap-2">
            <span v-for="s in product.series" :key="s.id" class="px-3 py-1.5 bg-[#FBF9F6] rounded-full text-sm text-gray-700 border border-[#EAE5DD]">{{ s.name }}</span>
          </div>
        </div>
        <div v-if="product.fabrics?.length" class="bg-white rounded-2xl p-6 lg:p-8 border border-[#EAE5DD]">
          <h2 class="text-xl font-bold text-[#0D1B2A] mb-4">{{ $t('product.detail.fabrics') }}</h2>
          <div class="flex flex-wrap gap-2">
            <span v-for="f in product.fabrics" :key="f.id" class="px-3 py-1.5 bg-[#FBF9F6] rounded-full text-sm text-gray-700 border border-[#EAE5DD]">{{ f.name }}</span>
          </div>
        </div>
      </div>

      <!-- Videos -->
      <div v-if="product.videos?.length" class="mt-8 lg:mt-12">
        <h2 class="text-xl font-bold text-[#0D1B2A] mb-6">{{ $t('product.detail.videos') }}</h2>
        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="(v, i) in product.videos" :key="v.id || i" class="bg-white rounded-2xl overflow-hidden border border-[#EAE5DD]">
            <div class="aspect-video bg-gray-100 relative group cursor-pointer" @click="openVideo(v)">
              <img v-if="v.cover" :src="v.cover" :alt="`${v.title || 'Product video'} - ${product.sku || ''}`" class="w-full h-full object-cover" loading="lazy" />
              <div class="absolute inset-0 flex items-center justify-center">
                <div class="w-14 h-14 bg-black/60 rounded-full flex items-center justify-center group-hover:bg-black/80 transition">
                  <svg class="w-6 h-6 text-white ml-0.5" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
                </div>
              </div>
            </div>
            <div v-if="v.title" class="p-3 text-sm font-medium text-gray-700">{{ v.title }}</div>
          </div>
        </div>
      </div>

      <!-- Gallery Images -->
      <div v-if="product.images?.length" class="mt-8 lg:mt-12">
        <h2 class="text-xl font-bold text-[#0D1B2A] mb-6">{{ $t('product.detail.gallery') }}</h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
          <div v-for="(img, i) in product.images" :key="img.id || i" class="bg-white rounded-2xl overflow-hidden border border-[#EAE5DD] aspect-square cursor-zoom-in group" @click="openLightbox(allImages.indexOf(img.url))">
            <img :src="img.url" :alt="img.alt || product.sku" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" loading="lazy" />
          </div>
        </div>
      </div>
    </div>

    <!-- Video Modal -->
    <transition name="fade">
      <div v-if="activeVideo" class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4" @click.self="activeVideo = null">
        <div class="bg-white rounded-2xl overflow-hidden max-w-3xl w-full shadow-2xl">
          <div class="relative">
            <button @click="activeVideo = null" class="absolute top-3 right-3 z-10 w-10 h-10 bg-black/50 rounded-full flex items-center justify-center text-white hover:bg-black/70 transition">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
            </button>
            <div class="aspect-video bg-gray-900">
              <iframe v-if="activeVideo.url" :src="videoEmbedUrl" class="w-full h-full" frameborder="0" allow="autoplay; fullscreen" allowfullscreen></iframe>
            </div>
          </div>
          <div v-if="activeVideo.title" class="p-4 text-sm font-medium text-gray-700">{{ activeVideo.title }}</div>
        </div>
      </div>
    </transition>

    <!-- 浮动图片查看器（Lightbox）：封面 + 图集多图全屏浮动展示 -->
    <transition name="fade">
      <div
        v-if="lightboxOpen && allImages.length"
        class="fixed inset-0 z-[60] bg-black/90 flex items-center justify-center select-none"
        @click.self="closeLightbox"
      >
        <!-- Close -->
        <button
          class="absolute top-4 right-4 z-10 w-10 h-10 rounded-full bg-white/10 hover:bg-white/25 text-white flex items-center justify-center transition"
          :aria-label="$t('common.close')"
          @click="closeLightbox"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
        <!-- Prev -->
        <button
          v-if="allImages.length > 1"
          class="absolute left-3 md:left-6 z-10 w-11 h-11 rounded-full bg-white/10 hover:bg-white/25 text-white flex items-center justify-center transition"
          aria-label="Previous image"
          @click.stop="lightboxPrev"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
        </button>
        <!-- 当前图片 -->
        <img
          :src="allImages[lightboxIndex]"
          :alt="`${product?.sku || ''} - Product image ${lightboxIndex + 1}`"
          class="max-h-[85vh] max-w-[92vw] object-contain rounded-lg shadow-2xl"
          @click.stop
        />
        <!-- Next -->
        <button
          v-if="allImages.length > 1"
          class="absolute right-3 md:right-6 z-10 w-11 h-11 rounded-full bg-white/10 hover:bg-white/25 text-white flex items-center justify-center transition"
          aria-label="Next image"
          @click.stop="lightboxNext"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
        </button>
        <!-- Counter -->
        <div v-if="allImages.length > 1" class="absolute bottom-5 left-1/2 -translate-x-1/2 bg-black/50 text-white text-xs px-3 py-1 rounded-full">
          {{ lightboxIndex + 1 }} / {{ allImages.length }}
        </div>
      </div>
    </transition>
  </div>

  <!-- Loading State -->
  <div v-else-if="loading" class="max-w-7xl mx-auto px-4 py-12">
    <div class="grid lg:grid-cols-2 gap-12 animate-pulse">
      <div class="aspect-[4/5] bg-gray-100 rounded-2xl"></div>
      <div class="space-y-4">
        <div class="h-8 bg-gray-100 rounded w-3/4"></div>
        <div class="h-4 bg-gray-100 rounded w-1/2"></div>
        <div class="h-4 bg-gray-100 rounded w-full"></div>
        <div class="h-4 bg-gray-100 rounded w-full"></div>
        <div class="grid grid-cols-2 gap-3 mt-6">
          <div v-for="i in 6" :key="i" class="h-20 bg-gray-100 rounded-xl"></div>
        </div>
      </div>
    </div>
  </div>

  <!-- Error State -->
  <div v-else class="max-w-7xl mx-auto px-4 py-24 text-center">
    <div class="text-6xl mb-4">🔍</div>
    <h2 class="text-2xl font-bold text-gray-900 mb-2">{{ $t('common.no_data') }}</h2>
    <p class="text-gray-500 mb-8">{{ $t('product.no_products') }}</p>
    <NuxtLink :to="localePath('/products')" class="inline-flex items-center gap-2 bg-[#D4A853] text-white px-6 py-3 rounded-xl font-semibold hover:bg-[#C49A3F] transition">
      {{ $t('product.products') }} →
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const localePath = useLocalePath()
const api = useApi()
const { locale } = useI18n()

const currentImageIndex = ref(0)
const copied = ref(false)
const activeVideo = ref<any>(null)

// ==================== 浮动图片查看器（Lightbox） ====================
// 点击主图 / 图集图片进入全屏浮动展示：支持前后切换、键盘 ←/→/Esc、点击遮罩关闭
const lightboxOpen = ref(false)
const lightboxIndex = ref(0)

function openLightbox(index: number) {
  if (!allImages.value.length) return
  lightboxIndex.value = Math.max(0, Math.min(index, allImages.value.length - 1))
  lightboxOpen.value = true
}

function closeLightbox() {
  lightboxOpen.value = false
}

function lightboxPrev() {
  lightboxIndex.value = (lightboxIndex.value - 1 + allImages.value.length) % allImages.value.length
}

function lightboxNext() {
  lightboxIndex.value = (lightboxIndex.value + 1) % allImages.value.length
}

function onLightboxKeydown(e: KeyboardEvent) {
  if (!lightboxOpen.value) return
  if (e.key === 'Escape') closeLightbox()
  else if (e.key === 'ArrowLeft') lightboxPrev()
  else if (e.key === 'ArrowRight') lightboxNext()
}

// 打开时锁定页面滚动，避免背景跟随滚动
watch(lightboxOpen, (open) => {
  if (import.meta.client) {
    document.body.style.overflow = open ? 'hidden' : ''
  }
})

const slug = route.params.slug as string

// SSR 阶段拉取产品详情 + 主题并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存）
const { data: product } = await useAsyncData<any>(
  'product-' + (locale.value || 'en') + '-' + slug,
  () => api.getProduct(slug).catch(() => null),
)
const { data: theme } = await useAsyncData<Record<string, any>>(
  'theme-' + (locale.value || 'en'),
  () => api.getTheme().catch(() => ({})),
)
const loading = computed(() => product.value == null)

// All images for gallery (cover + gallery images)
const allImages = computed(() => {
  const images: string[] = []
  if (product.value?.cover_image) {
    images.push(product.value.cover_image)
  }
  if (product.value?.images?.length) {
    product.value.images.forEach((img: any) => {
      if (img.url && !images.includes(img.url)) {
        images.push(img.url)
      }
    })
  }
  return images
})

const currentImage = computed(() => allImages.value[currentImageIndex.value] || null)

const waNumber = computed(() => theme.value?.whatsapp_number || '8612345678900')
const waMessage = computed(() => {
  const p = product.value
  if (!p) return 'Hello! I am interested in your sportswear products.'
  return `Hello! I am interested in your product: ${p.name || p.sku} (SKU: ${p.sku}). Please send me more information.`
})

const videoEmbedUrl = computed(() => {
  if (!activeVideo.value?.url) return ''
  const url = activeVideo.value.url
  // Convert YouTube URLs to embed format
  const ytMatch = url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]+)/)
  if (ytMatch) {
    return `https://www.youtube.com/embed/${ytMatch[1]}?autoplay=1`
  }
  return url
})

// SEO - 产品详情页（动态数据驱动）
const seoTitle = computed(() =>
  product.value?.name ? `${product.value.name} - Premium OEM/ODM Sportswear` : `${slug} - Sportswear Product`
)
const seoDescription = computed(() =>
  product.value?.brief?.slice(0, 160) || product.value?.description?.slice(0, 160) || ''
)
const seoKeywords = computed(() =>
  [product.value?.sku, product.value?.type, product.value?.material, 'sportswear OEM', 'ODM', 'custom activewear']
    .filter(Boolean).join(', ')
)

useSeoHead({
  title: seoTitle,
  description: seoDescription,
  keywords: seoKeywords,
  ogImage: computed(() => product.value?.cover_image || undefined),
  schema: computed(() =>
    product.value ? buildProductSchema(product.value) : undefined
  ),
})

function prevImage() {
  currentImageIndex.value = (currentImageIndex.value - 1 + allImages.value.length) % allImages.value.length
}

function nextImage() {
  currentImageIndex.value = (currentImageIndex.value + 1) % allImages.value.length
}

function openVideo(v: any) {
  activeVideo.value = v
}

async function copyProductInfo() {
  const p = product.value
  if (!p) return

  const info = [
    `SKU: ${p.sku}`,
    p.name ? `Name: ${p.name}` : '',
    p.type ? `Type: ${p.type}` : '',
    p.gender ? `Gender: ${p.gender}` : '',
    p.material ? `Material: ${p.material}` : '',
    p.composition ? `Composition: ${p.composition}` : '',
    p.weight ? `Weight: ${p.weight}` : '',
    p.elasticity ? `Elasticity: ${p.elasticity}` : '',
    p.fit ? `Fit: ${p.fit}` : '',
    p.season ? `Season: ${p.season}` : '',
    p.size_range ? `Size Range: ${p.size_range}` : '',
    p.brief ? `Brief: ${p.brief}` : '',
    p.description ? `Description: ${p.description}` : '',
    p.features ? `Features: ${p.features}` : '',
    p.usage ? `Usage: ${p.usage}` : '',
    `Sample MOQ: ${p.sample_moq || 1}`,
    `Production MOQ: ${p.production_moq || 300}`,
    `Color MOQ: ${p.color_moq || 100}`,
    `Size MOQ: ${p.size_moq || 100}`,
    p.specs?.length ? `\nSpecifications:\n${p.specs.map((s: any) => `  ${s.name}: ${s.value}`).join('\n')}` : '',
    p.customizations?.length ? `\nCustomization Options:\n${p.customizations.map((c: any) => `  ${c.type}: ${c.note || 'Available'}`).join('\n')}` : '',
    `\nURL: ${window.location.href}`,
  ].filter(Boolean).join('\n')

  try {
    await navigator.clipboard.writeText(info)
    copied.value = true
    setTimeout(() => { copied.value = false }, 3000)
  } catch {
    // Fallback for older browsers
    const textarea = document.createElement('textarea')
    textarea.value = info
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    copied.value = true
    setTimeout(() => { copied.value = false }, 3000)
  }
}

onMounted(() => {
  // 产品数据已在 SSR 阶段通过 useAsyncData 拉取并 hydration 复用，无需再次请求
  // 浮动图片查看器：键盘导航（←/→ 切换，Esc 关闭）
  if (import.meta.client) {
    window.addEventListener('keydown', onLightboxKeydown)
  }
})

onBeforeUnmount(() => {
  if (import.meta.client) {
    window.removeEventListener('keydown', onLightboxKeydown)
    document.body.style.overflow = ''
  }
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.scrollbar-hide::-webkit-scrollbar { display: none; }
.scrollbar-hide { -ms-overflow-style: none; scrollbar-width: none; }
</style>