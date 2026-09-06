<template>
  <div class="space-y-12">
    <template v-for="(m, idx) in parsed" :key="m.id || idx">
      <div v-if="m.type === 'text'" class="prose max-w-none">
        <h2 v-if="m.title" class="text-2xl font-bold text-[#0D1B2A] mb-4">{{ m.title }}</h2>
        <div class="text-gray-600 leading-relaxed" v-html="loc(m.config?.content, locale) || m.title"></div>
      </div>
      <div v-else-if="m.type === 'image'" class="text-center">
        <img v-if="m.config?.url" :src="imgUrl(m.config.url)" :alt="m.title || 'Image'" class="rounded-2xl max-w-full mx-auto" loading="lazy" />
        <p v-if="m.title" class="text-sm text-gray-500 mt-3">{{ m.title }}</p>
      </div>
      <div v-else-if="m.type === 'video'" class="aspect-video rounded-2xl overflow-hidden bg-black">
        <iframe v-if="m.config?.url" :src="m.config.url" class="w-full h-full" frameborder="0" allowfullscreen />
        <div v-else class="w-full h-full flex items-center justify-center text-white/50">{{ $t('common.no_video', 'No Video') }}</div>
      </div>
      <div v-else-if="m.type === 'banner'" class="rounded-2xl overflow-hidden">
        <img v-if="m.config?.image || m.config?.url" :src="imgUrl(m.config?.image || m.config?.url)" :alt="m.title || 'Banner'" class="w-full h-auto max-h-[480px] object-cover rounded-2xl" loading="lazy" />
      </div>
      <div v-else-if="['product_recommend','oem','odm','factory','production','case','certification','blog'].includes(m.type)" class="bg-white rounded-2xl p-8 border border-[#EAE5DD]">
        <h2 v-if="m.title" class="text-xl font-bold text-[#0D1B2A] mb-4">{{ m.title }}</h2>
        <div v-if="m.config?.description" class="text-gray-600 leading-relaxed" v-html="m.config.description"></div>
        <div v-else class="text-sm text-gray-400">{{ $t('common.module_type', 'Module type') }}: {{ m.type }}</div>
      </div>
      <div v-else-if="m.type === 'contact'" class="bg-gradient-to-r from-[#F5F0E8] to-[#FBF9F6] rounded-2xl p-8 text-center">
        <h2 v-if="m.title" class="text-2xl font-bold text-[#0D1B2A] mb-3">{{ m.title }}</h2>
        <p class="text-gray-600 mb-6">{{ m.config?.description || $t('common.ready_to_start', 'Ready to start your project?') }}</p>
        <NuxtLink :to="localePath('/contact')" class="inline-block bg-[#D4A853] text-white px-8 py-3.5 rounded-full font-semibold hover:bg-[#C49A3F] transition shadow-sm">
          {{ $t('home.get_quote') }}
        </NuxtLink>
      </div>
      <div v-else-if="m.type === 'about_hero'" class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-gray-900 to-gray-800 text-white">
        <img v-if="m.config?.bg_image" :src="imgUrl(m.config.bg_image)" class="absolute inset-0 w-full h-full object-cover opacity-30" loading="lazy" alt="" />
        <div class="relative max-w-3xl px-8 py-16 md:py-24">
          <h2 v-if="loc(m.config?.title, locale) || m.title" class="text-3xl md:text-4xl font-bold mb-4">{{ loc(m.config?.title, locale) || m.title }}</h2>
          <p v-if="loc(m.config?.subtitle, locale)" class="text-white/80 text-lg">{{ loc(m.config.subtitle, locale) }}</p>
        </div>
      </div>
      <div v-else-if="m.type === 'about_story'" class="grid md:grid-cols-2 gap-10 items-center">
        <div :class="m.config?.image_side === 'right' ? 'md:order-2' : ''">
          <h2 v-if="loc(m.config?.title, locale) || m.title" class="text-2xl md:text-3xl font-bold text-[#0D1B2A] mb-6">{{ loc(m.config?.title, locale) || m.title }}</h2>
          <div v-if="loc(m.config?.body, locale)" class="text-gray-600 leading-relaxed prose max-w-none" v-html="loc(m.config.body, locale)"></div>
        </div>
        <div v-if="m.config?.image" class="bg-gray-100 rounded-2xl overflow-hidden aspect-[4/3]" :class="m.config?.image_side === 'right' ? 'md:order-1' : ''">
          <img :src="imgUrl(m.config.image)" class="w-full h-full object-cover" loading="lazy" :alt="loc(m.config?.title, locale) || m.title || 'Story'" />
        </div>
      </div>
      <div v-else-if="m.type === 'about_stats'" class="grid grid-cols-3 gap-6">
        <div v-for="(it, i) in (m.config?.items || [])" :key="i" class="text-center">
          <div class="text-3xl font-bold text-[#D4A853]">{{ it.value }}</div>
          <div class="text-sm text-gray-500 mt-1">{{ loc(it.label, locale) }}</div>
        </div>
      </div>
      <div v-else-if="m.type === 'about_certifications'" class="bg-[#F5F0E8] rounded-2xl p-8">
        <h2 v-if="loc(m.config?.title, locale) || m.title" class="text-2xl md:text-3xl font-bold text-[#0D1B2A] mb-8 text-center">{{ loc(m.config?.title, locale) || m.title }}</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div v-for="cert in filteredCerts(m)" :key="cert.id" class="bg-white rounded-xl p-6 text-center shadow-sm">
            <div class="text-4xl mb-3">🏅</div>
            <h3 class="font-semibold text-sm">{{ cert.name }}</h3>
            <p v-if="cert.description" class="text-xs text-gray-500 mt-1">{{ cert.description }}</p>
          </div>
        </div>
      </div>
      <div v-else-if="m.type === 'about_process'" class="grid md:grid-cols-3 gap-6">
        <div v-for="(st, i) in (m.config?.steps || [])" :key="i" class="bg-white rounded-xl p-6 border border-[#EAE5DD]">
          <div class="text-3xl mb-4">{{ st.icon }}</div>
          <h3 class="font-semibold mb-2">{{ loc(st.title, locale) }}</h3>
          <p class="text-sm text-gray-600">{{ loc(st.desc, locale) }}</p>
        </div>
      </div>
      <div v-else-if="m.type === 'about_cta'" class="bg-gradient-to-r from-[#F5F0E8] to-[#FBF9F6] rounded-2xl p-8 text-center">
        <h2 v-if="loc(m.config?.title, locale) || m.title" class="text-2xl font-bold text-[#0D1B2A] mb-3">{{ loc(m.config?.title, locale) || m.title }}</h2>
        <p v-if="loc(m.config?.description, locale)" class="text-gray-600 mb-6">{{ loc(m.config.description, locale) }}</p>
        <NuxtLink v-if="m.config?.button_url" :to="m.config.button_url" class="inline-block bg-[#D4A853] text-white px-8 py-3.5 rounded-full font-semibold hover:bg-[#C49A3F] transition shadow-sm">
          {{ loc(m.config?.button_text, locale) || $t('home.get_quote') }}
        </NuxtLink>
      </div>
      <div v-else class="bg-gray-50 rounded-2xl p-6 border border-dashed border-gray-300 text-center text-sm text-gray-400">
        {{ $t('common.unknown_module', 'Unknown module') }}: {{ m.type }}
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ modules?: any[] }>()
const { locale } = useI18n()
const localePath = useLocalePath()
const api = useApi()

// config 可能是 JSON 字符串（后端 jsonb 直出）或对象，统一解析
function parseConfig(raw: any): Record<string, any> {
  if (!raw) return {}
  if (typeof raw === 'object') return raw
  try { return JSON.parse(raw) } catch { return {} }
}
const parsed = computed(() => (props.modules || []).map((m: any) => ({ ...m, config: parseConfig(m.config) })))

// 多语言字段取值：兼容旧字符串 + 新多语言对象 { en, zh, es, fr }
function loc(v: any, lang: string): string {
  if (!v) return ''
  if (typeof v === 'string') return v
  if (typeof v === 'object') return v[lang] || v.en || v.zh || v.es || v.fr || ''
  return ''
}

// 认证数据（about_certifications 模块使用，SSR 内联）
const { data: certifications } = await useAsyncData<any[]>(
  'certs-' + (locale.value || 'en'),
  () => api.getCertifications().then((res: any) => res?.items || res || []).catch(() => []),
)

// 认证墙按引用过滤：config.ref_ids 有值则只展示勾选的认证，否则展示全部
function filteredCerts(m: any): any[] {
  const all = certifications.value || []
  const refs = m?.config?.ref_ids
  if (Array.isArray(refs) && refs.length) return all.filter((c: any) => refs.includes(c.id))
  return all
}
</script>
