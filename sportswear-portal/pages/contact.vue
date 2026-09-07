<template>
  <div>
    <!-- ==================== 面包屑 ==================== -->
    <div class="border-b border-[#EAE5DD] bg-white">
      <div class="mx-auto max-w-7xl px-4 py-4 lg:px-8">
        <Breadcrumb
          :items="[
            { name: $t('nav.home'), to: localePath('/') },
            { name: $t('contact.title') },
          ]"
        />
      </div>
    </div>

    <!-- ==================== Hero ==================== -->
    <section class="relative overflow-hidden bg-[#0D1B2A] text-white">
      <div class="pointer-events-none absolute inset-0">
        <div class="absolute -right-24 -top-24 h-96 w-96 rounded-full bg-[#D4A853]/10 blur-3xl"></div>
        <div class="absolute -bottom-32 -left-24 h-[28rem] w-[28rem] rounded-full bg-[#1B2D44]/70 blur-3xl"></div>
        <div class="absolute inset-0 opacity-[0.05]" style="background-image: radial-gradient(#D4A853 1px, transparent 1px); background-size: 30px 30px;"></div>
        <div class="absolute inset-x-0 bottom-0 h-px bg-gradient-to-r from-transparent via-[#D4A853]/60 to-transparent"></div>
      </div>

      <div class="relative mx-auto max-w-7xl px-4 py-16 md:py-24 lg:px-8">
        <div class="max-w-3xl">
          <div class="mb-5 inline-flex items-center gap-3">
            <span class="h-px w-10 bg-[#D4A853]"></span>
            <span class="text-xs font-semibold uppercase tracking-[0.25em] text-[#D4A853]">{{ $t('contact.eyebrow') }}</span>
          </div>
          <h1 class="text-4xl font-bold leading-tight md:text-6xl">{{ $t('contact.title') }}</h1>
          <p class="mb-8 mt-5 max-w-2xl text-base leading-relaxed text-white/70 md:text-lg">{{ $t('contact.subtitle') }}</p>

          <div class="mb-8 flex flex-wrap items-center gap-4">
            <span class="text-2xl font-bold text-[#D4A853]">{{ brandName }}</span>
            <span v-if="brandSubtitle" class="hidden border-l border-white/15 pl-4 text-xs uppercase tracking-[0.25em] text-white/40 sm:inline">{{ brandSubtitle }}</span>
          </div>

          <div class="flex flex-wrap gap-x-6 gap-y-3 text-sm text-white/60">
            <span class="inline-flex items-center gap-2">
              <svg class="h-4 w-4 text-[#D4A853]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              {{ $t('contact.trust_response') }}
            </span>
            <span class="inline-flex items-center gap-2">
              <svg class="h-4 w-4 text-[#D4A853]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25" /></svg>
              {{ $t('contact.trust_oem') }}
            </span>
            <span class="inline-flex items-center gap-2">
              <svg class="h-4 w-4 text-[#D4A853]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
              {{ $t('contact.trust_cert') }}
            </span>
          </div>
        </div>
      </div>
    </section>

    <!-- ==================== 品牌定位 ==================== -->
    <section class="mx-auto max-w-7xl px-4 pt-12 md:pt-16 lg:px-8">
      <div class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-[#0D1B2A] to-[#1B2D44] p-8 text-white md:p-12">
        <div class="pointer-events-none absolute -right-16 -top-16 h-64 w-64 rounded-full bg-[#D4A853]/10 blur-2xl"></div>
        <div class="relative grid items-center gap-8 md:grid-cols-[1fr_auto]">
          <div class="max-w-2xl">
            <h2 class="mb-4 text-2xl font-bold md:text-3xl">{{ $t('contact.brand_title') }}</h2>
            <p class="text-sm leading-relaxed text-white/70 md:text-base">{{ $t('contact.brand_desc') }}</p>
          </div>
          <button type="button" @click="scrollToForm" class="inline-flex min-h-[48px] items-center justify-center gap-2 whitespace-nowrap rounded-full bg-[#D4A853] px-8 py-3.5 text-sm font-semibold text-white transition hover:bg-[#C49A3F]">
            {{ $t('contact.form_title') }}
            <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M19 14l-7 7m0 0l-7-7m7 7V3" /></svg>
          </button>
        </div>
      </div>
    </section>

      <!-- ==================== 服务亮点 ==================== -->
    <section class="mx-auto max-w-7xl px-4 py-12 md:py-16 lg:px-8">
      <div class="mx-auto mb-10 max-w-2xl text-center">
        <h2 class="text-2xl font-bold text-[#0D1B2A] md:text-3xl">{{ $t('contact.services_title') }}</h2>
        <div class="mt-4 flex items-center justify-center gap-2">
          <span class="h-1 w-6 rounded-full bg-[#D4A853]"></span>
          <span class="h-1 w-2 rounded-full bg-[#D4A853]/40"></span>
          <span class="h-1 w-1 rounded-full bg-[#D4A853]/30"></span>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4 md:grid-cols-3 md:gap-6">
        <div v-for="s in services" :key="s.key" class="group rounded-2xl border border-[#EAE5DD] bg-white p-6 transition-all duration-300 hover:-translate-y-1 hover:border-[#D4A853]/60 hover:shadow-xl">
          <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl border border-[#EAE5DD] bg-[#FBF9F6] transition-colors group-hover:border-[#D4A853] group-hover:bg-[#D4A853]">
            <svg class="h-6 w-6 text-[#D4A853] transition-colors group-hover:text-white" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" :d="s.icon" /></svg>
          </div>
          <h3 class="mb-1.5 text-sm font-semibold text-[#0D1B2A] md:text-base">{{ $t(`contact.${s.key}_title`) }}</h3>
          <p class="text-xs leading-relaxed text-gray-500 md:text-sm">{{ $t(`contact.${s.key}_desc`) }}</p>
        </div>
      </div>
    </section>

    <!-- ==================== 询盘问题模板 ==================== -->
    <section class="mx-auto max-w-7xl px-4 pb-12 md:pb-16 lg:px-8">
      <div class="rounded-3xl border border-[#EAE5DD] bg-white p-6 md:p-10">
        <div class="mb-8 flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div class="max-w-2xl">
            <h2 class="text-2xl font-bold text-[#0D1B2A] md:text-3xl">{{ $t('contact.templates.title') }}</h2>
            <p class="mt-2 text-sm text-gray-500">{{ $t('contact.templates.hint') }}</p>
          </div>
          <span class="hidden shrink-0 items-center gap-2 rounded-full bg-[#FBF9F6] px-4 py-2 text-xs font-medium text-gray-500 md:inline-flex">
            <svg class="h-4 w-4 text-[#D4A853]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" /></svg>
            {{ $t('contact.trust_response') }}
          </span>
        </div>
        <div v-for="group in templateGroups" :key="group.category" class="mb-6 last:mb-0">
          <p class="mb-3 text-xs font-semibold uppercase tracking-[0.15em] text-gray-400">{{ $t('faq.categories.' + group.category) }}</p>
          <div class="flex flex-wrap gap-2.5">
            <button
              v-for="tpl in group.templates"
              :key="tpl.key"
              type="button"
              @click="applyTemplate(tpl)"
              class="inline-flex min-h-[44px] items-center gap-2 rounded-full border px-4 py-2.5 text-left text-sm transition"
              :class="selectedTemplate === tpl.key
                ? 'border-[#D4A853] bg-[#D4A853] text-white shadow-sm'
                : 'border-gray-200 bg-white text-gray-700 hover:border-[#D4A853] hover:text-[#D4A853]'"
            >
              <svg class="h-4 w-4 shrink-0" :class="selectedTemplate === tpl.key ? 'text-white' : 'text-[#D4A853]'" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path v-if="selectedTemplate === tpl.key" stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" d="M12 5v14m-7-7h14" />
              </svg>
              {{ tpl.text }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- ==================== 表单 + 联系信息 ==================== -->
    <section ref="formSectionRef" class="mx-auto max-w-7xl scroll-mt-24 px-4 pb-16 md:pb-24 lg:px-8">
      <div class="mx-auto mb-10 max-w-2xl text-center">
        <h2 class="text-2xl font-bold text-[#0D1B2A] md:text-3xl">{{ $t('contact.form_title') }}</h2>
        <p class="mt-3 text-sm text-gray-500 md:text-base">{{ $t('contact.form_subtitle') }}</p>
      </div>

      <div class="grid gap-8 lg:grid-cols-5 lg:gap-10">
        <!-- 表单 -->
        <div class="rounded-3xl border border-[#EAE5DD] bg-white p-6 md:p-8 lg:col-span-3">
          <form ref="formRef" @submit.prevent="handleSubmit" class="space-y-5">
            <div class="grid gap-5 sm:grid-cols-2">
              <div class="relative">
                <span class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
                  <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" /></svg>
                </span>
                <input v-model="form.name" :placeholder="$t('contact.name')" required class="min-h-[48px] w-full rounded-xl border border-gray-200 py-3.5 pl-12 pr-4 text-sm text-gray-800 transition focus:border-[#D4A853] focus:outline-none focus:ring-4 focus:ring-[#D4A853]/15" />
              </div>
              <div class="relative">
                <span class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
                  <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" /></svg>
                </span>
                <input v-model="form.email" type="email" :placeholder="$t('contact.email')" required class="min-h-[48px] w-full rounded-xl border border-gray-200 py-3.5 pl-12 pr-4 text-sm text-gray-800 transition focus:border-[#D4A853] focus:outline-none focus:ring-4 focus:ring-[#D4A853]/15" />
              </div>
            </div>
            <div class="grid gap-5 sm:grid-cols-2">
              <div class="relative">
                <span class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
                  <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 002.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-.282.376-.769.542-1.21.38a12.035 12.035 0 01-7.143-7.143c-.162-.441.004-.928.38-1.21l1.293-.97c.363-.271.527-.734.417-1.173L6.963 3.102a1.125 1.125 0 00-1.091-.852H4.5A2.25 2.25 0 002.25 4.5v2.25z" /></svg>
                </span>
                <input v-model="form.phone" :placeholder="$t('contact.phone')" class="min-h-[48px] w-full rounded-xl border border-gray-200 py-3.5 pl-12 pr-4 text-sm text-gray-800 transition focus:border-[#D4A853] focus:outline-none focus:ring-4 focus:ring-[#D4A853]/15" />
              </div>
              <div class="relative">
                <span class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
                  <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 21h16.5M4.5 3h15M5.25 3v18m13.5-18v18M9 6.75h1.5m-1.5 3h1.5m-1.5 3h1.5m3-6H15m-1.5 3H15m-1.5 3H15M9 21v-3.375c0-.621.504-1.125 1.125-1.125h3.75c.621 0 1.125.504 1.125 1.125V21" /></svg>
                </span>
                <input v-model="form.company" :placeholder="$t('contact.company')" class="min-h-[48px] w-full rounded-xl border border-gray-200 py-3.5 pl-12 pr-4 text-sm text-gray-800 transition focus:border-[#D4A853] focus:outline-none focus:ring-4 focus:ring-[#D4A853]/15" />
              </div>
            </div>
            <div>
              <textarea v-model="form.message" :placeholder="$t('contact.message')" required rows="5" class="w-full rounded-xl border border-gray-200 px-4 py-3.5 text-sm text-gray-800 transition focus:border-[#D4A853] focus:outline-none focus:ring-4 focus:ring-[#D4A853]/15"></textarea>
            </div>
            <!-- 附件上传：图片/文档/压缩包，单文件≤20MB，最多 5 个 -->
            <div>
              <label class="inline-flex min-h-[48px] cursor-pointer items-center gap-2 rounded-xl border border-dashed border-gray-300 px-4 py-3 text-sm text-gray-600 transition hover:border-[#D4A853] hover:text-[#D4A853]">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 10-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" /></svg>
                {{ uploading ? $t('contact.uploading') : $t('contact.attachment') }}
                <input ref="fileInputRef" type="file" multiple class="hidden" :accept="ACCEPT_TYPES" :disabled="uploading" @change="onFileChange" />
              </label>
              <p class="mt-1.5 text-xs text-gray-400">{{ $t('contact.attachment_hint') }}</p>
              <p v-if="fileError" class="mt-1 text-xs text-red-500">{{ fileError }}</p>
              <ul v-if="files.length" class="mt-2 space-y-1.5">
                <li v-for="(f, i) in files" :key="f.uid" class="flex items-center justify-between gap-2 rounded-md bg-gray-50 px-3 py-2 text-sm">
                  <span class="truncate text-gray-600">{{ f.file.name }} <span class="text-gray-400">({{ formatSize(f.file.size) }})</span></span>
                  <button type="button" class="shrink-0 text-gray-400 hover:text-red-500" @click="removeFile(i)">✕</button>
                </li>
              </ul>
            </div>

            <button type="submit" :disabled="submitting || uploading" class="inline-flex min-h-[48px] w-full items-center justify-center gap-2 rounded-xl bg-[#D4A853] px-8 py-3.5 text-sm font-semibold text-white transition hover:bg-[#C49A3F] disabled:opacity-50 md:w-auto">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" /></svg>
              {{ submitting || uploading ? $t('contact.sending') : $t('contact.submit') }}
            </button>
          </form>
          <p v-if="successMsg" class="mt-4 flex items-center gap-2 text-sm text-emerald-600">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            {{ successMsg }}
          </p>
        </div>
        <!-- 联系信息 -->
        <div class="lg:col-span-2">
          <div class="overflow-hidden rounded-3xl bg-gradient-to-br from-[#0D1B2A] to-[#1B2D44] text-white">
            <div class="border-b border-white/10 p-6 md:p-8">
              <h3 class="text-lg font-bold md:text-xl">{{ $t('contact.info_title') }}</h3>
              <p class="mt-2 text-sm text-white/60">{{ $t('contact.info_desc') }}</p>
            </div>
            <ul class="space-y-6 p-6 md:p-8">
              <li class="flex items-start gap-4">
                <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#D4A853]/15">
                  <svg class="h-5 w-5 text-[#D4A853]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" /></svg>
                </span>
                <div class="min-w-0">
                  <p class="text-xs uppercase tracking-wide text-white/40">{{ $t('contact.email_label') }}</p>
                  <a :href="contactEmailHref" class="mt-0.5 block break-all text-sm font-medium text-white transition hover:text-[#D4A853]">{{ $t('contact.email_value') }}</a>
                </div>
              </li>
              <li class="flex items-start gap-4">
                <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#D4A853]/15">
                  <svg class="h-5 w-5 text-[#D4A853]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 002.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-.282.376-.769.542-1.21.38a12.035 12.035 0 01-7.143-7.143c-.162-.441.004-.928.38-1.21l1.293-.97c.363-.271.527-.734.417-1.173L6.963 3.102a1.125 1.125 0 00-1.091-.852H4.5A2.25 2.25 0 002.25 4.5v2.25z" /></svg>
                </span>
                <div class="min-w-0">
                  <p class="text-xs uppercase tracking-wide text-white/40">{{ $t('contact.phone_label') }}</p>
                  <a :href="contactPhoneHref" class="mt-0.5 block text-sm font-medium text-white transition hover:text-[#D4A853]">{{ $t('contact.phone_value') }}</a>
                </div>
              </li>
              <li class="flex items-start gap-4">
                <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#D4A853]/15">
                  <svg class="h-5 w-5 text-[#D4A853]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z" /><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z" /></svg>
                </span>
                <div class="min-w-0">
                  <p class="text-xs uppercase tracking-wide text-white/40">{{ $t('contact.address_label') }}</p>
                  <p class="mt-0.5 text-sm font-medium text-white">{{ $t('contact.address_value') }}</p>
                </div>
              </li>
              <li class="flex items-start gap-4">
                <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-[#D4A853]/15">
                  <svg class="h-5 w-5 text-[#D4A853]" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </span>
                <div class="min-w-0">
                  <p class="text-xs uppercase tracking-wide text-white/40">{{ $t('contact.hours_label') }}</p>
                  <p class="mt-0.5 text-sm font-medium text-white">{{ $t('contact.hours_value') }}</p>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
<script setup lang="ts">
const api = useApi()
const { t, locale } = useI18n()
const route = useRoute()
const localePath = useLocalePath()
const submitting = ref(false)
const successMsg = ref('')
const form = reactive({ name: '', email: '', phone: '', company: '', message: '' })

// 联系信息链接（mailto / tel，文案走多语言词条，与 Footer 保持一致）
const contactEmailHref = computed(() => 'mailto:' + t('contact.email_value'))
const contactPhoneHref = computed(() => 'tel:' + t('contact.phone_value').replace(/[^\d+]/g, ''))

// 来自产品详情页的询盘参数：结构化产品字段（提交时回传，用于询盘评分与归类）
const productId = ref<string | null>(null)
const productCategory = ref('')

// 品牌名 / 副标题（优先主题配置，回退默认）
const brandName = ref('SPORTSWEAR')
const brandSubtitle = ref('Premium Mfg.')

// 服务亮点（图标为 Heroicons outline path，随卡片 hover 切换品牌色）
const services = [
  { key: 'service_oem', icon: 'M3.75 21h16.5M4.5 3h15M5.25 3v18m13.5-18v18M9 6.75h1.5m-1.5 3h1.5m-1.5 3h1.5m3-6H15m-1.5 3H15m-1.5 3H15M9 21v-3.375c0-.621.504-1.125 1.125-1.125h3.75c.621 0 1.125.504 1.125 1.125V21' },
  { key: 'service_odm', icon: 'M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L6.832 19.82a4.5 4.5 0 01-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 011.13-1.897L16.863 4.487zm0 0L19.5 7.125' },
  { key: 'service_custom', icon: 'M10.5 6h9.75M10.5 6a1.5 1.5 0 11-3 0m3 0a1.5 1.5 0 10-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-9.75 0h9.75' },
  { key: 'service_moq', icon: 'M15.75 10.5V6a3.75 3.75 0 10-7.5 0v4.5m11.356-1.993l1.263 12c.07.665-.45 1.243-1.119 1.243H4.25a1.125 1.125 0 01-1.12-1.243l1.264-12A1.125 1.125 0 015.513 7.5h12.974c.576 0 1.059.435 1.119 1.007zM8.625 10.5a.375.375 0 11-.75 0 .375.375 0 01.75 0zm7.5 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z' },
  { key: 'service_factory', icon: 'M2.25 21h19.5m-18-18v18m10.5-18v18m6-13.5V21M6.75 6.75h.75m-.75 3h.75m-.75 3h.75m3-6h.75m-.75 3h.75m-.75 3h.75M6.75 21v-3.375c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21M3 3h12m-.75 4.5H21m-3.75 3.75h.008v.008h-.008v-.008zm0 3h.008v.008h-.008v-.008zm0 3h.008v.008h-.008v-.008z' },
  { key: 'service_cert', icon: 'M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z' },
]

// 询盘问题模板：优先后台「询盘模板」模型（/public/inquiry-templates，可运营增删/排序/上下架），
// 接口失败或未配置时降级到静态词条分组（contact.templates.*，方案 B 兜底）。
const { data: templateData } = await useAsyncData<any>(
  'inquiry-templates-' + (locale.value || 'en'),
  () => api.getInquiryTemplates().catch(() => null),
)

const templateGroups = computed(() => {
  const items = templateData.value
  if (Array.isArray(items) && items.length) {
    const grouped: Record<string, Array<{ key: string; text: string; projectType: string }>> = {}
    for (const it of items) {
      const cat = it.category || 'other'
      if (!grouped[cat]) grouped[cat] = []
      grouped[cat].push({ key: it.id, text: it.question, projectType: it.project_type || '' })
    }
    return Object.keys(grouped).map((category) => ({ category, templates: grouped[category] }))
  }
  // 静态降级（按 FAQ 分类分组，project_type 走本地映射）
  return inquiryTemplateGroups.map((g) => ({
    category: g.category,
    templates: g.items.map((key) => ({
      key,
      text: t(`contact.templates.${key}`),
      projectType: PROJECT_TYPE_BY_CATEGORY[g.category] || '',
    })),
  }))
})

const selectedTemplate = ref('')
const templateProjectType = ref('')
const formRef = ref<HTMLFormElement | null>(null)
const formSectionRef = ref<HTMLElement | null>(null)

// 品牌定位卡片 CTA：平滑滚动到下方询盘表单
function scrollToForm() {
  nextTick(() => {
    formSectionRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  })
}

function applyTemplate(tpl: { key: string; text: string; projectType: string }) {
  selectedTemplate.value = tpl.key
  form.message = tpl.text
  templateProjectType.value = tpl.projectType || ''
  nextTick(() => {
    formRef.value?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  })
}

// 静态兜底分组（接口失败时使用）
const inquiryTemplateGroups = [
  { category: 'moq', items: ['q_moq'] },
  { category: 'oem', items: ['q_oem'] },
  { category: 'odm', items: ['q_odm'] },
  { category: 'sample', items: ['q_sample'] },
  { category: 'payment', items: ['q_price', 'q_payment'] },
  { category: 'fabric', items: ['q_fabric'] },
  { category: 'logistics', items: ['q_shipping'] },
  { category: 'production', items: ['q_leadtime'] },
  { category: 'quality', items: ['q_quality'] },
]
const PROJECT_TYPE_BY_CATEGORY: Record<string, string> = {
  oem: 'oem',
  odm: 'odm',
  private_label: 'private_label',
}

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
  } else if (route.query.product_info) {
    // 来自产品详情页的询盘：预填完整产品信息，便于销售识别具体产品
    form.message = String(route.query.product_info)
  } else if (route.query.product) {
    // 兼容旧链接（仅携带 SKU / 名称 / 规格）
    form.message = `I am interested in product: ${String(route.query.product)}` +
      (route.query.specs ? `\nSelected options: ${String(route.query.specs)}` : '')
  }
  // 结构化产品字段：随询盘提交，触发后端「已选产品」评分加成
  if (route.query.product_id) productId.value = String(route.query.product_id)
  if (route.query.product_category) productCategory.value = String(route.query.product_category)
  try {
    const theme = await api.getTheme()
    if (theme?.brand_name) brandName.value = theme.brand_name
    if (theme?.brand_subtitle) brandSubtitle.value = theme.brand_subtitle
  } catch { /* ignore */ }
})

// SEO - 联系我们
const contactSeo = await useRouteSeo('contact')
useSeoHead({
  title: () => contactSeo.value?.title || 'Contact Us - OEM/ODM Sportswear Manufacturer',
  description: () =>
    contactSeo.value?.description ||
    'Get in touch with our OEM/ODM sportswear manufacturing team. Request a quote, ask about customization options, or discuss your activewear project. We respond within 24 hours.',
  keywords: () =>
    contactSeo.value?.keywords ||
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
      project_type: templateProjectType.value || (route.query.project_type as string) || 'contact',
      product_id: productId.value || undefined,
      product_category: productCategory.value || undefined,
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