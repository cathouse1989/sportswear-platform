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
          <p class="text-gray-600 mb-4 leading-relaxed">
            Founded in 2008, we have grown from a small workshop into a leading OEM/ODM sportswear manufacturer
            serving global brands. Our 50,000 sqm facility houses advanced manufacturing equipment and a dedicated
            R&D team of 30+ professionals.
          </p>
          <p class="text-gray-600 mb-4 leading-relaxed">
            We specialize in producing high-performance activewear, yoga wear, running gear, team uniforms,
            and custom sportswear solutions for brands worldwide. Our commitment to quality and innovation
            has earned us long-term partnerships with clients in North America, Europe, and Australia.
          </p>
          <div class="grid grid-cols-3 gap-6 mt-8">
            <div class="text-center">
              <div class="text-3xl font-bold text-[#D4A853]">15+</div>
              <div class="text-sm text-gray-500 mt-1">Years Experience</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-[#D4A853]">500+</div>
              <div class="text-sm text-gray-500 mt-1">Global Clients</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-[#D4A853]">50K</div>
              <div class="text-sm text-gray-500 mt-1">sqm Facility</div>
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
            <h3 class="font-semibold text-sm">ISO 9001:2015</h3>
            <p class="text-xs text-gray-500 mt-1">Quality Management</p>
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
// SEO - 关于我们
useSeoHead({
  title: 'About Us - OEM/ODM Sportswear Manufacturer Since 2008',
  description:
    'Professional OEM/ODM sportswear manufacturer with 15+ years of experience. ISO 9001, BSCI, OEKO-TEX certified. 50,000 sqm facility serving 500+ global brands.',
  keywords:
    'sportswear manufacturer history, OEM factory China, ODM sportswear company, ISO certified activewear factory, garment manufacturing China',
  ogImage: '/images/og-about.jpg',
})

const api = useApi()
const localePath = useLocalePath()
const { locale } = useI18n()

const factorySteps = [
  { icon: '✂️', title: 'Cutting', desc: 'Computerized cutting machines with precision up to 0.1mm' },
  { icon: '🧵', title: 'Sewing', desc: '500+ skilled workers operating modern sewing lines' },
  { icon: '🖨️', title: 'Printing', desc: 'Sublimation, screen printing, and embroidery capabilities' },
  { icon: '🔬', title: 'Quality Control', desc: 'Multi-point inspection system at every production stage' },
  { icon: '📦', title: 'Packaging', desc: 'Professional packaging with your brand requirements' },
  { icon: '🚢', title: 'Shipping', desc: 'Global logistics partnerships for timely delivery' },
]



// SSR 阶段拉取认证列表并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存）
const { data: certifications } = await useAsyncData<any[]>(
  'certs-' + (locale.value || 'en'),
  () => api.getCertifications()
    .then((res: any) => res?.items || res || [])
    .catch(() => []),
)
</script>