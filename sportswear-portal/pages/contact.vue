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
          <button type="submit" :disabled="submitting" class="w-full md:w-auto bg-[#D4A853] text-white px-8 py-3.5 rounded-lg font-semibold hover:bg-[#C49A3F] transition disabled:opacity-50 min-h-[48px] text-sm">
            {{ submitting ? 'Sending...' : $t('home.get_quote') }}
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

useHead({ title: 'Contact Us - Sportswear OEM/ODM' })

async function handleSubmit() {
  submitting.value = true
  successMsg.value = ''
  try {
    await api.submitLead({ ...form, project_type: 'contact' })
    successMsg.value = 'Thank you! We will contact you within 24 hours.'
    Object.assign(form, { name: '', email: '', phone: '', company: '', message: '' })
  } catch (e: any) {
    successMsg.value = 'Failed to send. Please try again later.'
  } finally { submitting.value = false }
}
</script>