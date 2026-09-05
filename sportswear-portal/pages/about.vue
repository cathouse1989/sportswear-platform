<template>
  <div>
    <!-- Hero -->
    <section class="bg-gradient-to-br from-gray-900 to-gray-800 text-white py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h1 class="text-4xl md:text-5xl font-bold mb-4">{{ $t('about.title') || 'About Us' }}</h1>
        <p class="text-white/70 text-lg max-w-2xl">{{ $t('about.subtitle') || 'Professional OEM/ODM sportswear manufacturer with 15+ years of experience.' }}</p>
      </div>
    </section>

    <!-- Company Overview -->
    <section class="max-w-7xl mx-auto px-4 lg:px-8 py-16 md:py-20">
      <div class="grid md:grid-cols-2 gap-12 items-center">
        <div>
          <h2 class="text-2xl md:text-3xl font-bold text-[#0D1B2A] mb-6">{{ $t('about.our_story') || 'Our Story' }}</h2>
          <p class="text-gray-600 mb-4 leading-relaxed">{{ $t('about.story_1') }}</p>
          <p class="text-gray-600 mb-4 leading-relaxed">{{ $t('about.story_2') }}</p>
          <div class="grid grid-cols-3 gap-6 mt-8">
            <div class="text-center">
              <div class="text-3xl font-bold text-[#D4A853]">15+</div>
              <div class="text-sm text-gray-500 mt-1">{{ $t('about.stat_years_label') }}</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-[#D4A853]">500+</div>
              <div class="text-sm text-gray-500 mt-1">{{ $t('about.stat_clients_label') }}</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-[#D4A853]">50K</div>
              <div class="text-sm text-gray-500 mt-1">{{ $t('about.stat_facility_label') }}</div>
            </div>
          </div>
        </div>
        <div class="bg-gray-100 rounded-2xl overflow-hidden aspect-[4/3] flex items-center justify-center text-6xl">
          🏭
        </div>
      </div>
    </section>

    <!-- Certifications -->
    <section class="bg-[#F5F0E8] py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h2 class="text-2xl md:text-3xl font-bold text-[#0D1B2A] mb-8 text-center">{{ $t('about.certifications') || 'Our Certifications' }}</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div v-for="cert in certifications" :key="cert.id" class="bg-white rounded-xl p-6 text-center shadow-sm hover:shadow-md transition">
            <div class="text-4xl mb-3">🏅</div>
            <h3 class="font-semibold text-sm">{{ cert.name }}</h3>
            <p v-if="cert.issuer" class="text-xs text-gray-500 mt-1">{{ cert.issuer }}</p>
          </div>
          <div v-if="!certifications?.length" v-for="i in 4" :key="i" class="bg-white rounded-xl p-6 text-center shadow-sm">
            <div class="text-4xl mb-3">🏅</div>
            <h3 class="font-semibold text-sm">{{ $t('about.cert_iso_name') }}</h3>
            <p class="text-xs text-gray-500 mt-1">{{ $t('about.cert_iso_issuer') }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Factory Tour -->
    <section class="max-w-7xl mx-auto px-4 lg:px-8 py-16 md:py-20">
      <h2 class="text-2xl md:text-3xl font-bold text-[#0D1B2A] mb-8 text-center">{{ $t('about.factory_tour') || 'Factory Tour' }}</h2>
      <div class="grid md:grid-cols-3 gap-6">
        <div v-for="step in factorySteps" :key="step.title" class="bg-white rounded-xl p-6 border border-[#EAE5DD] hover:shadow-md transition">
          <div class="text-3xl mb-4">{{ step.icon }}</div>
          <h3 class="font-semibold mb-2">{{ step.title }}</h3>
          <p class="text-sm text-gray-600">{{ step.desc }}</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
// SEO - 关于我们（多语言，随语言切换）
useSeoHead({
  title: () => t('about.seo_title'),
  description: () => t('about.seo_description'),
  keywords: () => t('about.seo_keywords'),
  ogImage: '/images/og-about.jpg',
})

const api = useApi()
const { locale, t } = useI18n()

const factorySteps = computed(() => [
  { icon: '✂️', title: t('about.step_cutting_title'), desc: t('about.step_cutting_desc') },
  { icon: '🧵', title: t('about.step_sewing_title'), desc: t('about.step_sewing_desc') },
  { icon: '🖨️', title: t('about.step_printing_title'), desc: t('about.step_printing_desc') },
  { icon: '🔬', title: t('about.step_qc_title'), desc: t('about.step_qc_desc') },
  { icon: '📦', title: t('about.step_packaging_title'), desc: t('about.step_packaging_desc') },
  { icon: '🚢', title: t('about.step_shipping_title'), desc: t('about.step_shipping_desc') },
])



// SSR 阶段拉取认证列表并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存）
const { data: certifications } = await useAsyncData<any[]>(
  'certs-' + (locale.value || 'en'),
  () => api.getCertifications()
    .then((res: any) => res?.items || res || [])
    .catch(() => []),
)
</script>