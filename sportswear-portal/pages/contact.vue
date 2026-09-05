<template>
  <div>
    <!-- 品牌 + 服务 Hero -->
    <section class="bg-gradient-to-br from-[#0D1B2A] to-[#1B2D44] text-white py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4 lg:px-8">
        <h1 class="text-3xl md:text-5xl font-bold mb-3">{{ $t('contact.title') }}</h1>
        <p class="text-white/70 text-base md:text-lg mb-6">{{ $t('contact.subtitle') }}</p>
        <div class="flex flex-wrap items-center gap-4">
          <span class="text-[#D4A853] font-bold text-xl">{{ brandName }}</span>
          <span v-if="brandSubtitle" class="text-white/50 text-sm uppercase tracking-[0.2em]">{{ brandSubtitle }}</span>
        </div>
      </div>
    </section>

    <div class="max-w-7xl mx-auto px-4 lg:px-8 py-12 md:py-16">
      <!-- 品牌定位 -->
      <div class="max-w-3xl mb-10">
        <h2 class="text-2xl md:text-3xl font-bold text-[#0D1B2A] mb-4">{{ $t('contact.brand_title') }}</h2>
        <p class="text-gray-600 text-sm md:text-base leading-relaxed">{{ $t('contact.brand_desc') }}</p>
      </div>

      <!-- 服务亮点 -->
      <h3 class="text-xl md:text-2xl font-bold text-[#0D1B2A] mb-6">{{ $t('contact.services_title') }}</h3>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-4 md:gap-6 mb-12">
        <div v-for="s in services" :key="s.key" class="bg-white border border-[#EAE5DD] rounded-2xl p-5 hover:shadow-lg transition">
          <div class="text-3xl mb-3">{{ s.icon }}</div>
          <h4 class="font-semibold text-[#0D1B2A] text-sm md:text-base mb-1.5">{{ $t(`contact.${s.key}_title`) }}</h4>
          <p class="text-gray-500 text-xs md:text-sm">{{ $t(`contact.${s.key}_desc`) }}</p>
        </div>
      </div>

      <!-- 表单 + 联系信息 -->
      <div class="grid md:grid-cols-2 gap-8 md:gap-12">
        <div>
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <input v-model="form.name" :placeholder="$t('contact.name')" required class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none min-h-[48px]" />
            <input v-model="form.email" type="email" :placeholder="$t('contact.email')" required class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none min-h-[48px]" />
            <input v-model="form.phone" :placeholder="$t('contact.phone')" class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none min-h-[48px]" />
            <input v-model="form.company" :placeholder="$t('contact.company')" class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none min-h-[48px]" />
            <textarea v-model="form.message" :placeholder="$t('contact.message')" required rows="5" class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none"></textarea>
            <!-- 附件上传：图片/文档/压缩包，单文件≤20MB，最多 5 个 -->
            <div>
              <label class="inline-flex items-center gap-2 cursor-pointer border border-dashed border-gray-300 rounded-lg px-4 py-3 text-sm text-gray-600 hover:border-[#D4A853] hover:text-[#D4A853] transition min-h-[48px]">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 10-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" /></svg>
                {{ uploading ? $t('contact.uploading') : $t('contact.attachment') }}
                <input ref="fileInputRef" type="file" multiple class="hidden" :accept="ACCEPT_TYPES" :disabled="uploading" @change="onFileChange" />
              </label>
              <p class="mt-1.5 text-xs text-gray-400">{{ $t('contact.attachment_hint') }}</p>
              <p v-if="fileError" class="mt-1 text-xs text-red-500">{{ fileError }}</p>
              <ul v-if="files.length" class="mt-2 space-y-1.5">
                <li v-for="(f, i) in files" :key="f.uid" class="flex items-center justify-between gap-2 bg-gray-50 rounded-md px-3 py-2 text-sm">
                  <span class="truncate text-gray-600">{{ f.file.name }} <span class="text-gray-400">({{ formatSize(f.file.size) }})</span></span>
                  <button type="button" class="shrink-0 text-gray-400 hover:text-red-500" @click="removeFile(i)">✕</button>
                </li>
              </ul>
            </div>
            <button type="submit" :disabled="submitting || uploading" class="w-full md:w-auto bg-[#D4A853] text-white px-8 py-3.5 rounded-lg font-semibold hover:bg-[#C49A3F] transition disabled:opacity-50 min-h-[48px] text-sm">
              {{ submitting || uploading ? $t('contact.sending') : $t('contact.submit') }}
            </button>
          </form>
          <p v-if="successMsg" class="mt-4 text-green-600 text-sm">{{ successMsg }}</p>
        </div>
        <div>
          <h3 class="font-semibold text-base md:text-lg mb-4">{{ $t('contact.info_title') }}</h3>
          <div class="space-y-4 text-gray-600 text-sm md:text-base">
            <p><strong>{{ $t('contact.email_label') }}:</strong> {{ $t('contact.email_value') }}</p>
            <p><strong>{{ $t('contact.phone_label') }}:</strong> {{ $t('contact.phone_value') }}</p>
            <p><strong>{{ $t('contact.address_label') }}:</strong> {{ $t('contact.address_value') }}</p>
            <p><strong>{{ $t('contact.hours_label') }}:</strong> {{ $t('contact.hours_value') }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
const api = useApi()
const { t } = useI18n()
const route = useRoute()
const submitting = ref(false)
const successMsg = ref('')
const form = reactive({ name: '', email: '', phone: '', company: '', message: '' })

// 品牌名 / 副标题（优先主题配置，回退默认）
const brandName = ref('SPORTSWEAR')
const brandSubtitle = ref('Premium Mfg.')

// 服务亮点
const services = [
  { key: 'service_oem', icon: '🏭' },
  { key: 'service_odm', icon: '🎨' },
  { key: 'service_custom', icon: '✂️' },
  { key: 'service_moq', icon: '📦' },
  { key: 'service_factory', icon: '🏗️' },
  { key: 'service_cert', icon: '✅' },
]

// ==================== 附件上传（询盘闭环：附件随询盘提交，后台详情可查看） ====================
const ACCEPT_TYPES = '.jpg,.jpeg,.png,.gif,.webp,.pdf,.doc,.docx,.xls,.xlsx,.zip,.rar'
const MAX_FILE_MB = 20
const MAX_FILES = 5
const fileInputRef = ref<HTMLInputElement | null>(null)
const files = ref<Array<{ uid: number; file: File }>>([])
const fileError = ref('')
const uploading = ref(false)
let uidSeq = 0

function onFileChange(e: Event) {
  fileError.value = ''
  const input = e.target as HTMLInputElement
  const picked = Array.from(input.files || [])
  input.value = '' // 允许重复选择同一文件
  for (const file of picked) {
    if (files.value.length >= MAX_FILES) {
      fileError.value = t('contact.file_max')
      break
    }
    const ext = '.' + (file.name.split('.').pop() || '').toLowerCase()
    if (!ACCEPT_TYPES.split(',').includes(ext)) {
      fileError.value = `${t('contact.file_type')}: ${file.name}`
      continue
    }
    if (file.size > MAX_FILE_MB * 1024 * 1024) {
      fileError.value = `${file.name} ${t('contact.file_size')}`
      continue
    }
    files.value.push({ uid: ++uidSeq, file })
  }
}

function removeFile(idx: number) {
  files.value.splice(idx, 1)
  fileError.value = ''
}

function formatSize(bytes: number) {
  if (bytes >= 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + 'MB'
  if (bytes >= 1024) return (bytes / 1024).toFixed(0) + 'KB'
  return bytes + 'B'
}

// 加载主题品牌信息（品牌名 / 副标题）
onMounted(async () => {
  // 来自案例详情的 CTA：预填留言，便于销售识别来源案例
  if (route.query.case_title) {
    form.message = t('contact.case_prefill', { case: String(route.query.case_title) })
  }
  try {
    const theme = await api.getTheme()
    if (theme?.brand_name) brandName.value = theme.brand_name
    if (theme?.brand_subtitle) brandSubtitle.value = theme.brand_subtitle
  } catch { /* ignore */ }
})

// SEO - 联系我们
useSeoHead({
  title: 'Contact Us - OEM/ODM Sportswear Manufacturer',
  description:
    'Get in touch with our OEM/ODM sportswear manufacturing team. Request a quote, ask about customization options, or discuss your activewear project. We respond within 24 hours.',
  keywords:
    'contact sportswear manufacturer, OEM inquiry, ODM quote, activewear supplier, China garment factory contact',
})

async function handleSubmit() {
  submitting.value = true
  successMsg.value = ''
  try {
    // 先逐个上传附件，收集 URL 后随询盘一起提交（attachments 为 JSON 数组字符串）
    uploading.value = !!files.value.length
    const attachments: Array<{ url: string; name: string; size: number; type: string }> = []
    for (const f of files.value) {
      const res = await api.uploadLeadAttachment(f.file)
      attachments.push({
        url: res.url,
        name: res.file_name || f.file.name,
        size: res.file_size ?? f.file.size,
        type: res.file_type || '',
      })
    }
    await api.submitLead({
      ...form,
      project_type: (route.query.project_type as string) || 'contact',
      attachments: attachments.length ? JSON.stringify(attachments) : '',
    })
    successMsg.value = t('contact.success')
    Object.assign(form, { name: '', email: '', phone: '', company: '', message: '' })
    files.value = []
  } catch (e: any) {
    successMsg.value = e?.data?.message || t('contact.error')
  } finally { uploading.value = false; submitting.value = false }
}
</script>