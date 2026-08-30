<template>
  <div>
    <!-- Desktop: Floating Sidebar (hidden on mobile) -->
    <div class="fixed right-0 top-1/2 -translate-y-1/2 z-40 flex-col gap-1 hidden md:flex pointer-events-auto">
      <!-- WhatsApp -->
      <a
        :href="waUrl"
        target="_blank"
        rel="noopener noreferrer"
        class="flex items-center gap-2 bg-[#25D366] text-white px-3 py-3 rounded-l-lg shadow-lg hover:bg-[#20BD5A] transition-all min-h-[44px]"
        @mouseenter="showWa = true"
        @mouseleave="showWa = false"
      >
        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
          <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z" />
        </svg>
        <span v-if="showWa" class="text-sm whitespace-nowrap">{{ $t('floating.whatsapp') }}</span>
      </a>

      <!-- Inquiry Button -->
      <button
        @click="showPopup = true"
        class="flex items-center gap-2 bg-[#0D1B2A] text-white px-3 py-4 rounded-l-lg shadow-lg hover:bg-[#1B2D44] transition-all min-h-[44px]"
        @mouseenter="showInquiry = true"
        @mouseleave="showInquiry = false"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
        <span v-if="showInquiry" class="text-sm whitespace-nowrap">{{ $t('floating.inquiry') }}</span>
      </button>
    </div>

    <!-- Mobile: Bottom Action Bar (hidden on desktop) -->
    <div class="fixed bottom-0 left-0 right-0 z-50 flex md:hidden bg-white border-t border-gray-200 shadow-2xl safe-area-bottom">
      <a
        :href="waUrl"
        target="_blank"
        rel="noopener noreferrer"
        class="flex-1 flex items-center justify-center gap-2 py-3 bg-[#25D366] text-white text-sm font-semibold min-h-[52px] active:bg-[#20BD5A] transition-colors"
      >
        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
          <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z" />
        </svg>
        <span>WhatsApp</span>
      </a>
      <button
        @click="showPopup = true"
        class="flex-1 flex items-center justify-center gap-2 py-3 bg-[#0D1B2A] text-white text-sm font-semibold min-h-[52px] active:bg-[#1B2D44] transition-colors"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
        <span>Inquiry</span>
      </button>
    </div>

    <!-- Inquiry Popup -->
    <transition name="fade">
      <div v-if="showPopup" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm" @click.self="dismissPopup">
        <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-3 p-5 sm:mx-4 sm:p-8 relative max-h-[92vh] overflow-y-auto">
          <button @click="dismissPopup" class="absolute top-4 right-4 text-gray-400 hover:text-gray-600">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <div class="text-center mb-6">
            <div class="w-16 h-16 bg-[#D4A853]/10 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-[#D4A853]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-[#0D1B2A]">{{ $t('floating.send_us_inquiry') }}</h3>
            <p class="text-gray-500 mt-2 text-sm">{{ $t('floating.get_free_quote') }}</p>
          </div>
          <form @submit.prevent="submitLead" class="space-y-3">
            <input v-model="form.name" :placeholder="$t('floating.your_name')" required class="w-full border border-gray-200 rounded-lg px-4 py-3 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none" />
            <input v-model="form.email" type="email" :placeholder="$t('floating.your_email')" required class="w-full border border-gray-200 rounded-lg px-4 py-3 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none" />
            <textarea v-model="form.message" :placeholder="$t('floating.your_message')" rows="3" class="w-full border border-gray-200 rounded-lg px-4 py-3 text-sm focus:ring-2 focus:ring-[#D4A853]/30 focus:border-[#D4A853] focus:outline-none"></textarea>
            <button type="submit" :disabled="sending" class="w-full bg-[#D4A853] text-white py-3 rounded-lg font-semibold hover:bg-[#C49A3F] transition disabled:opacity-50">
              {{ sending ? $t('floating.sending') : $t('floating.send_message') }}
            </button>
          </form>
          <p v-if="success" class="mt-3 text-center text-green-600 text-sm">{{ success }}</p>
          <p v-if="error" class="mt-3 text-center text-red-500 text-sm">{{ error }}</p>
          <button @click="dismissPopup" class="w-full text-center text-sm text-gray-400 hover:text-gray-600 mt-4">
            {{ $t('floating.maybe_later') }}
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const showWa = ref(false)
const showInquiry = ref(false)
const showPopup = ref(false)
const sending = ref(false)
const success = ref('')
const error = ref('')
const form = reactive({ name: '', email: '', message: '' })
const waNum = ref('8612345678900')
const waMsg = ref('Hello! I am interested in your sportswear products.')

const waUrl = computed(() => 'https://wa.me/' + waNum.value + '?text=' + encodeURIComponent(waMsg.value))

// 使用 localStorage 记录弹窗关闭状态，7天内不再自动弹出
const POPUP_STORAGE_KEY = 'sportswear_inquiry_popup'

function isPopupDismissed(): boolean {
  try {
    const val = localStorage.getItem(POPUP_STORAGE_KEY)
    if (!val) return false
    const dismissedAt = parseInt(val, 10)
    if (isNaN(dismissedAt)) return false
    // 7 天有效期
    const now = Date.now()
    return (now - dismissedAt) < 7 * 24 * 60 * 60 * 1000
  } catch {
    return false
  }
}

function dismissPopup() {
  showPopup.value = false
  try {
    localStorage.setItem(POPUP_STORAGE_KEY, String(Date.now()))
  } catch {}
}

onMounted(async () => {
  try {
    const theme = await api.getTheme()
    if (theme?.whatsapp_number) waNum.value = theme.whatsapp_number
    if (theme?.whatsapp_message) waMsg.value = theme.whatsapp_message
  } catch (e) {}
  // 仅在桌面端才显示自动弹窗，移动端不自动弹出
  const isMobile = window.innerWidth < 768
  if (!isMobile && !isPopupDismissed()) {
    setTimeout(() => { showPopup.value = true }, 8000)
  }
})

async function submitLead() {
  sending.value = true
  success.value = ''
  error.value = ''
  try {
    await api.submitLead({ name: form.name, email: form.email, message: form.message, source: 'floating' })
    success.value = 'Thank you! We will contact you within 24 hours.'
    form.name = ''; form.email = ''; form.message = ''
    setTimeout(() => { dismissPopup(); success.value = '' }, 2000)
  } catch (e) {
    error.value = 'Failed to send. Please try again.'
  } finally { sending.value = false }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>