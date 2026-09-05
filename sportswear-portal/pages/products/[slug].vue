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
              :src="imgUrl(currentImage)"
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
              <img :src="imgUrl(img)" :alt="`${product.sku || product.name} thumbnail ${i + 1}`" class="w-full h-full object-cover" loading="lazy" />
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

          <!-- Selectable Options（颜色/尺码/面料等可点选规格） -->
          <div v-if="optionGroups.length" class="bg-white rounded-2xl p-5 lg:p-6 border border-[#EAE5DD] mb-8">
            <h3 class="text-sm font-semibold text-[#0D1B2A] mb-4 flex items-center gap-2">
              <svg class="w-4 h-4 text-[#D4A853]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" /></svg>
              {{ $t('product.detail.select_options') }}
            </h3>
            <div v-for="group in optionGroups" :key="group.name" class="mb-5 last:mb-0">
              <div class="text-xs font-medium text-gray-500 uppercase tracking-wider mb-2.5">
                {{ group.name }}
                <span v-if="selectedOptions[group.name]" class="normal-case tracking-normal text-[#D4A853] font-semibold">· {{ selectedOptions[group.name] }}</span>
              </div>
              <div class="flex flex-wrap gap-2.5">
                <button
                  v-for="opt in group.options"
                  :key="opt"
                  type="button"
                  @click="selectOption(group.name, opt)"
                  class="inline-flex items-center gap-2 px-3.5 py-2 rounded-xl border text-sm font-medium transition min-h-[40px]"
                  :class="selectedOptions[group.name] === opt ? 'border-[#0D1B2A] bg-[#0D1B2A] text-white' : 'border-[#EAE5DD] bg-white text-gray-700 hover:border-gray-400'"
                  :aria-pressed="selectedOptions[group.name] === opt"
                >
                  <span v-if="group.isColor" class="w-4 h-4 rounded-full border border-black/10 shrink-0" :style="{ backgroundColor: colorHex(opt) }"></span>
                  {{ opt }}
                </button>
              </div>
            </div>
            <div v-if="hasSelection" class="mt-5 pt-4 border-t border-[#EAE5DD] flex items-center justify-between gap-3 flex-wrap">
              <span class="text-sm text-gray-600">{{ $t('product.detail.your_selection') }}: <span class="font-semibold text-[#0D1B2A]">{{ selectionSummary }}</span></span>
              <button @click="clearSelection" class="text-xs text-gray-400 hover:text-gray-600 underline">{{ $t('product.detail.clear') }}</button>
            </div>
          </div>

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
              :to="localePath('/contact') + '?product=' + encodeURIComponent(product.sku) + (selectionSummary ? '&specs=' + encodeURIComponent(selectionSummary) : '')"
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
          <div class="relative" :class="{ 'max-h-64 overflow-hidden': expandable(product.description) && !expanded.description }">
            <div class="prose prose-sm max-w-none text-gray-600 leading-relaxed whitespace-pre-wrap">{{ product.description }}</div>
            <div v-if="expandable(product.description) && !expanded.description" class="absolute bottom-0 inset-x-0 h-20 bg-gradient-to-t from-white to-transparent pointer-events-none"></div>
          </div>
          <button v-if="expandable(product.description)" @click="toggleExpand('description')" class="mt-3 text-sm font-medium text-[#D4A853] hover:text-[#C49A3F] transition">
            {{ expanded.description ? $t('product.detail.show_less') : $t('product.detail.show_more') }}
          </button>
        </div>
        <div v-if="product.features" class="bg-white rounded-2xl p-6 lg:p-8 border border-[#EAE5DD]">
          <h2 class="text-xl font-bold text-[#0D1B2A] mb-4">{{ $t('product.detail.features') }}</h2>
          <div class="relative" :class="{ 'max-h-64 overflow-hidden': expandable(product.features) && !expanded.features }">
            <div class="prose prose-sm max-w-none text-gray-600 leading-relaxed whitespace-pre-wrap">{{ product.features }}</div>
            <div v-if="expandable(product.features) && !expanded.features" class="absolute bottom-0 inset-x-0 h-20 bg-gradient-to-t from-white to-transparent pointer-events-none"></div>
          </div>
          <button v-if="expandable(product.features)" @click="toggleExpand('features')" class="mt-3 text-sm font-medium text-[#D4A853] hover:text-[#C49A3F] transition">
            {{ expanded.features ? $t('product.detail.show_less') : $t('product.detail.show_more') }}
          </button>
        </div>
        <div v-if="product.usage" class="bg-white rounded-2xl p-6 lg:p-8 border border-[#EAE5DD]">
          <h2 class="text-xl font-bold text-[#0D1B2A] mb-4">{{ $t('product.detail.usage') }}</h2>
          <div class="relative" :class="{ 'max-h-64 overflow-hidden': expandable(product.usage) && !expanded.usage }">
            <div class="prose prose-sm max-w-none text-gray-600 leading-relaxed whitespace-pre-wrap">{{ product.usage }}</div>
            <div v-if="expandable(product.usage) && !expanded.usage" class="absolute bottom-0 inset-x-0 h-20 bg-gradient-to-t from-white to-transparent pointer-events-none"></div>
          </div>
          <button v-if="expandable(product.usage)" @click="toggleExpand('usage')" class="mt-3 text-sm font-medium text-[#D4A853] hover:text-[#C49A3F] transition">
            {{ expanded.usage ? $t('product.detail.show_less') : $t('product.detail.show_more') }}
          </button>
        </div>
      </div>

      <!-- Specs Table -->
      <div v-if="staticSpecs.length" class="mt-8 lg:mt-12">
        <h2 class="text-xl font-bold text-[#0D1B2A] mb-6">{{ $t('product.detail.specifications') }}</h2>
        <div class="bg-white rounded-2xl border border-[#EAE5DD] overflow-hidden">
          <table class="w-full">
            <tbody>
              <tr v-for="(spec, i) in staticSpecs" :key="spec.name + i" class="border-b border-[#EAE5DD] last:border-b-0">
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
          <div class="space-y-4">
            <div v-for="f in product.fabrics" :key="f.id" class="border border-[#EAE5DD] rounded-xl p-4">
              <div class="flex items-center justify-between gap-3 flex-wrap">
                <div class="font-semibold text-[#0D1B2A]">{{ f.name }}</div>
                <div v-if="f.code" class="text-xs text-gray-400 font-mono">{{ f.code }}</div>
              </div>
              <div v-if="f.composition" class="text-sm text-gray-600 mt-1">{{ f.composition }}</div>
              <div v-if="fabricTags(f).length" class="flex flex-wrap gap-1.5 mt-3">
                <span v-for="tag in fabricTags(f)" :key="tag.label" class="px-2 py-1 bg-[#FBF9F6] rounded text-xs text-gray-600 border border-[#EAE5DD]">{{ tag.label }}: {{ tag.value }}</span>
              </div>
              <div v-if="f.description" class="text-sm text-gray-500 mt-3 leading-relaxed">{{ f.description }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Videos -->
      <div v-if="product.videos?.length" class="mt-8 lg:mt-12">
        <h2 class="text-xl font-bold text-[#0D1B2A] mb-6">{{ $t('product.detail.videos') }}</h2>
        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div v-for="(v, i) in product.videos" :key="v.id || i" class="bg-white rounded-2xl overflow-hidden border border-[#EAE5DD]">
            <div class="aspect-video bg-gray-100 relative group cursor-pointer" @click="openVideo(v)">
              <img v-if="v.cover" :src="imgUrl(v.cover)" :alt="`${v.title || 'Product video'} - ${product.sku || ''}`" class="w-full h-full object-cover" loading="lazy" />
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
            <img :src="imgUrl(img.url)" :alt="img.alt || product.sku" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" loading="lazy" />
          </div>
        </div>
      </div>

      <!-- Related Products -->
      <div v-if="related?.length" class="mt-8 lg:mt-12">
        <h2 class="text-xl font-bold text-[#0D1B2A] mb-6">{{ $t('product.detail.related_products') }}</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 md:gap-6">
          <NuxtLink
            v-for="r in related"
            :key="r.id"
            :to="localePath('/products/' + r.slug)"
            class="group bg-white rounded-2xl overflow-hidden border border-[#EAE5DD] hover:shadow-lg transition-all active:scale-[0.98]"
          >
            <div class="aspect-[4/5] bg-gray-100 relative overflow-hidden">
              <img v-if="r.cover_image" :src="imgUrl(r.cover_image)" :alt="r.sku"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700" loading="lazy" />
              <img v-if="relatedHoverImg(r)" :src="relatedHoverImg(r)" :alt="`${r.sku} - detail`"
                class="absolute inset-0 w-full h-full object-cover opacity-0 group-hover:opacity-100 transition-opacity duration-300" loading="lazy" />
              <div v-if="!r.cover_image" class="w-full h-full flex items-center justify-center text-5xl opacity-30">🏋️</div>
              <div class="absolute bottom-3 left-3">
                <span class="px-2.5 py-1 bg-white/95 rounded-full text-[10px] font-semibold">{{ r.type }}</span>
              </div>
            </div>
            <div class="p-3">
              <h3 class="text-xs font-semibold text-[#0D1B2A] line-clamp-1">{{ r.name || r.sku }}</h3>
              <p class="text-[10px] text-gray-500 mt-1">MOQ {{ r.production_moq || 300 }}</p>
            </div>
          </NuxtLink>
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
          :src="imgUrl(allImages[lightboxIndex])"
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

// 相关产品（同分类优先，不足则用全量补充，用于底部推荐与站内 SEO 内链）
const { data: related } = await useAsyncData<any[]>(
  'related-' + (locale.value || 'en') + '-' + slug,
  async () => {
    try {
      const catId = product.value?.category_id
      const params: Record<string, any> = { page: 1, pageSize: 8 }
      if (catId) params.category_id = catId
      const res = await api.getProducts(params)
      let items = Array.isArray(res) ? res : (res?.items || [])
      items = items.filter((r: any) => r.slug !== slug)
      if (items.length < 4 && catId) {
        const all = await api.getProducts({ page: 1, pageSize: 12 })
        const allItems = Array.isArray(all) ? all : (all?.items || [])
        const seen = new Set(items.map((x: any) => x.id))
        for (const it of allItems) {
          if (it.slug !== slug && !seen.has(it.id)) {
            items.push(it)
            seen.add(it.id)
          }
          if (items.length >= 4) break
        }
      }
      return items.slice(0, 4)
    } catch {
      return []
    }
  },
)

const loading = computed(() => product.value == null)

function relatedHoverImg(r: any): string {
  return galleryImageUrls(r)[1] || ''
}

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

// ==================== 可点选规格（颜色/尺码/面料等选项） ====================
// 规格值支持用逗号 / 斜杠 / 竖线分隔多个可选项，门户据此渲染为可点选的选择器，
// 采购人员可快速点选规格，选择结果会带入「复制信息 / WhatsApp / 询盘链接」。
const COLOR_HEX: Record<string, string> = {
  black: '#1f2937',
  white: '#ffffff',
  red: '#dc2626',
  blue: '#2563eb',
  navy: '#1e3a8a',
  'navy blue': '#1e3a8a',
  green: '#16a34a',
  yellow: '#eab308',
  orange: '#f97316',
  pink: '#ec4899',
  purple: '#8b5cf6',
  gray: '#6b7280',
  grey: '#6b7280',
  brown: '#92400e',
  beige: '#e7d8c9',
  tan: '#d2b48c',
  maroon: '#7f1d1d',
  gold: '#d4a853',
  silver: '#c0c0c0',
  charcoal: '#374151',
  'light blue': '#93c5fd',
  'dark green': '#14532d',
  olive: '#708238',
  cream: '#fdf6e3',
  khaki: '#bdb76b',
  camel: '#c19a6b',
  burgundy: '#800020',
  teal: '#0d9488',
  lavender: '#c4b5fd',
}

function isColorGroup(name: string): boolean {
  return /color|colour|颜色|色/i.test(name)
}

function colorHex(name: string): string {
  const key = String(name || '').trim().toLowerCase()
  return COLOR_HEX[key] || '#e5e7eb'
}

// 解析 product.specs：同一「规格名」下拥有多个可选项时，视为可点选规格组
const optionGroups = computed<Array<{ name: string; options: string[]; isColor: boolean }>>(() => {
  const map = new Map<string, { name: string; options: string[]; isColor: boolean }>()
  for (const spec of product.value?.specs || []) {
    const name = String(spec.name || '').trim()
    const raw = String(spec.value || '').trim()
    if (!name || !raw) continue
    const options = raw.split(/[,，/／、|;；]+/).map((s: string) => s.trim()).filter(Boolean)
    const key = name.toLowerCase()
    if (!map.has(key)) map.set(key, { name, options: [], isColor: isColorGroup(name) })
    const group = map.get(key)!
    for (const o of options) if (!group.options.includes(o)) group.options.push(o)
  }
  return [...map.values()].filter(g => g.options.length > 1)
})

// 不可点选的单一值规格 → 继续在下方规格表格中展示
const staticSpecs = computed(() => {
  const selectable = new Set(optionGroups.value.map(g => g.name.toLowerCase()))
  const rows: Array<{ name: string; value: string }> = []
  const seen = new Set<string>()
  for (const spec of product.value?.specs || []) {
    const name = String(spec.name || '').trim()
    const raw = String(spec.value || '').trim()
    if (!name || !raw) continue
    if (selectable.has(name.toLowerCase())) continue
    const key = `${name}::${raw}`
    if (seen.has(key)) continue
    seen.add(key)
    rows.push({ name, value: raw })
  }
  return rows
})

const selectedOptions = reactive<Record<string, string>>({})
const hasSelection = computed(() => Object.keys(selectedOptions).length > 0)
const selectionSummary = computed(() =>
  Object.entries(selectedOptions).map(([k, v]) => `${k}: ${v}`).join('; ')
)

function selectOption(groupName: string, value: string) {
  if (selectedOptions[groupName] === value) {
    delete selectedOptions[groupName]
    // 取消选择后回到封面图
    currentImageIndex.value = 0
  } else {
    selectedOptions[groupName] = value
    // 若该选项在图集中有对应图片（alt / 文件名匹配），则联动切换主图
    const idx = imageIndexForOption(value)
    if (idx >= 0) currentImageIndex.value = idx
  }
}

// 根据选项值在图集中查找对应图片（优先 alt 匹配，其次 URL 文件名匹配）
function imageIndexForOption(opt: string): number {
  const q = String(opt || '').trim().toLowerCase()
  if (!q) return -1
  const images = product.value?.images || []
  // 1) alt 精确 / 包含匹配（后台图片管理可在 alt 中填写颜色名）
  for (const img of images) {
    const alt = String(img.alt || '').trim().toLowerCase()
    if (alt && (alt === q || alt.includes(q) || q.includes(alt))) {
      const idx = allImages.value.indexOf(img.url)
      if (idx >= 0) return idx
    }
  }
  // 2) URL 文件名包含匹配（如 .../black.jpg）
  for (const img of images) {
    try {
      const file = String(img.url || '').split('/').pop()?.toLowerCase() || ''
      if (file && file.includes(q.replace(/\s+/g, '-'))) {
        const idx = allImages.value.indexOf(img.url)
        if (idx >= 0) return idx
      }
    } catch { /* ignore */ }
  }
  return -1
}

function clearSelection() {
  Object.keys(selectedOptions).forEach(k => delete selectedOptions[k])
}

// ==================== 长文本展开/收起 + 面料详情 ====================
const expanded = reactive<Record<string, boolean>>({})

function expandable(text?: string): boolean {
  return (text?.length || 0) > 260
}

function toggleExpand(key: string) {
  expanded[key] = !expanded[key]
}

// 面料完整属性标签（采购快速了解面料性能）
function fabricTags(f: any): Array<{ label: string; value: string }> {
  const map: Array<[string, any]> = [
    ['Weight', f?.weight],
    ['Elasticity', f?.elasticity],
    ['Breathability', f?.breathability],
    ['Moisture Wicking', f?.moisture_wicking],
    ['Softness', f?.softness],
    ['Compression', f?.compression],
    ['UV Protection', f?.uv_protection],
    ['Eco Friendly', f?.eco_friendly],
  ]
  return map.filter(([, v]) => v).map(([label, v]) => ({ label, value: String(v) }))
}

const waNumber = computed(() => theme.value?.whatsapp_number || '8612345678900')
const waMessage = computed(() => {
  const p = product.value
  if (!p) return 'Hello! I am interested in your sportswear products.'
  const selection = selectionSummary.value ? `\nSelected options: ${selectionSummary.value}` : ''
  return `Hello! I am interested in your product: ${p.name || p.sku} (SKU: ${p.sku}).${selection}\nPlease send me more information.`
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
    selectionSummary.value ? `\nSelected Options: ${selectionSummary.value}` : '',
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