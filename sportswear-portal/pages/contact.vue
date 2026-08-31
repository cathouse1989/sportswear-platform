<template>
  <div class="max-w-4xl mx-auto px-4 py-12">
    <h1 class="text-2xl md:text-3xl font-bold mb-6 md:mb-8">{{ $t('nav.contact') }}</h1>
    <div class="grid md:grid-cols-2 gap-8 md:gap-12">
      <div>
        <p class="text-gray-600 text-sm md:text-base mb-6">Send us a message and we'll get back to you within 24 hours.</p>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <input v-model="form.name" placeholder="Name *" required class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none min-h-[48px]" />
          <input v-model="form.email" type="email" placeholder="Email *" required class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none min-h-[48px]" />
          <input v-model="form.phone" placeholder="Phone" class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none min-h-[48px]" />
          <input v-model="form.company" placeholder="Company" class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none min-h-[48px]" />
          <textarea v-model="form.message" placeholder="Message *" required rows="5" class="w-full border border-gray-200 rounded-lg px-4 py-3.5 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none"></textarea>
          <!-- 附件上传：图片/文档/压缩包，单文件≤20MB，最多 5 个 -->
          <div>
            <label class="inline-flex items-center gap-2 cursor-pointer border border-dashed border-gray-300 rounded-lg px-4 py-3 text-sm text-gray-600 hover:border-[#D4A853] hover:text-[#D4A853] transition min-h-[48px]">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 10-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" /></svg>
              {{ uploading ? 'Uploading...' : 'Upload Attachment' }}
              <input ref="fileInputRef" type="file" multiple class="hidden" :accept="ACCEPT_TYPES" :disabled="uploading" @change="onFileChange" />
            </label>
            <p class="mt-1.5 text-xs text-gray-400">Only supports .jpg/.png/.gif/.webp/.pdf/.doc/.docx/.xls/.xlsx/.zip/.rar, max 20MB each, up to {{ MAX_FILES }} files.</p>
            <p v-if="fileError" class="mt-1 text-xs text-red-500">{{ fileError }}</p>
            <ul v-if="files.length" class="mt-2 space-y-1.5">
              <li v-for="(f, i) in files" :key="f.uid" class="flex items-center justify-between gap-2 bg-gray-50 rounded-md px-3 py-2 text-sm">
                <span class="truncate text-gray-600">{{ f.file.name }} <span class="text-gray-400">({{ formatSize(f.file.size) }})</span></span>
                <button type="button" class="shrink-0 text-gray-400 hover:text-red-500" @click="removeFile(i)">✕</button>
              </li>
            </ul>
          </div>
          <button type="submit" :disabled="submitting || uploading" class="w-full md:w-auto bg-[#D4A853] text-white px-8 py-3.5 rounded-lg font-semibold hover:bg-[#C49A3F] transition disabled:opacity-50 min-h-[48px] text-sm">
            {{ submitting || uploading ? 'Sending...' : $t('home.get_quote') }}
          </button>
        </form>
        <p v-if="successMsg" class="mt-4 text-green-600 text-sm">{{ successMsg }}</p>
      </div>
      <div>
        <h3 class="font-semibold text-base md:text-lg mb-4">Contact Information</h3>
        <div class="space-y-4 text-gray-600 text-sm md:text-base">
          <p><strong>Email:</strong> info@sportswear.com</p>
          <p><strong>Phone:</strong> +86 123 4567 8900</p>
          <p><strong>Address:</strong> Guangzhou, China</p>
          <p><strong>Working Hours:</strong> Mon-Fri 9:00-18:00 (GMT+8)</p>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
const api = useApi()
const submitting = ref(false)
const successMsg = ref('')
const form = reactive({ name: '', email: '', phone: '', company: '', message: '' })

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
      fileError.value = `Up to ${MAX_FILES} attachments allowed.`
      break
    }
    const ext = '.' + (file.name.split('.').pop() || '').toLowerCase()
    if (!ACCEPT_TYPES.split(',').includes(ext)) {
      fileError.value = `Unsupported file type: ${file.name}`
      continue
    }
    if (file.size > MAX_FILE_MB * 1024 * 1024) {
      fileError.value = `${file.name} exceeds the ${MAX_FILE_MB}MB limit.`
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

useHead({ title: 'Contact Us - Sportswear OEM/ODM' })

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
      project_type: 'contact',
      attachments: attachments.length ? JSON.stringify(attachments) : '',
    })
    successMsg.value = 'Thank you! We will contact you within 24 hours.'
    Object.assign(form, { name: '', email: '', phone: '', company: '', message: '' })
    files.value = []
  } catch (e: any) {
    successMsg.value = e?.data?.message || 'Failed to send. Please try again later.'
  } finally { uploading.value = false; submitting.value = false }
}
</script>