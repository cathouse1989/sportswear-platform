<template>
  <div class="min-h-screen flex flex-col bg-[#FBF9F6]">
    <!-- 顶部公告条 -->
    <div class="bg-[#0D1B2A] text-[#D4A853] text-xs text-center py-2.5 px-4 tracking-wider uppercase font-light">
      {{ $t('common.free_shipping') }}
    </div>

    <!-- Header -->
    <header class="sticky top-0 z-40 bg-white/95 backdrop-blur-lg border-b border-[#EAE5DD]">
      <div class="max-w-7xl mx-auto px-4 lg:px-8 h-[68px] flex items-center justify-between">
        <!-- Logo -->
        <NuxtLink :to="localePath('/')" class="group">
          <SportswearLogo
            variant="light"
            :logoSrc="theme.logo_url"
            :logoAlt="theme.logo_alt"
            :brandName="theme.brand_name"
            :brandSubtitle="theme.brand_subtitle"
          />
        </NuxtLink>

        <!-- Desktop Nav -->
        <nav class="hidden lg:flex items-center gap-1">
          <template v-for="item in navTree" :key="item.id || item.path">
            <!-- 无子级：普通链接 -->
            <a
              v-if="!item.children?.length"
              :href="navHref(item.path, item.isExternal)"
              :target="item.target === '_blank' ? '_blank' : undefined"
              :rel="item.target === '_blank' ? 'noopener noreferrer' : undefined"
              class="text-sm text-[#4A4A4A] hover:text-[#0D1B2A] transition-colors relative px-3 py-1 after:absolute after:bottom-0 after:left-3 after:right-3 after:h-[2px] after:w-0 after:bg-[#D4A853] after:transition-all after:duration-300 hover:after:w-[calc(100%-24px)]"
            >
              {{ item.label }}
            </a>
            <!-- 有子级：下拉菜单 -->
            <div v-else class="relative group">
              <button class="flex items-center gap-1 text-sm text-[#4A4A4A] hover:text-[#0D1B2A] transition-colors px-3 py-1 rounded-lg hover:bg-[#F5F0E8]">
                {{ item.label }}
                <svg class="w-2.5 h-2.5 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" /></svg>
              </button>
              <div class="absolute left-0 top-full mt-1 bg-white border border-[#EAE5DD] rounded-xl shadow-2xl py-2 min-w-[200px] z-50 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200">
                <a
                  v-for="child in item.children"
                  :key="child.id"
                  :href="navHref(child.path || '/', child.isExternal)"
                  :target="child.target === '_blank' ? '_blank' : undefined"
                  :rel="child.target === '_blank' ? 'noopener noreferrer' : undefined"
                  class="block px-5 py-3 text-sm text-[#4A4A4A] hover:bg-[#FBF9F6] hover:text-[#0D1B2A] transition-colors"
                >
                  {{ child.label }}
                </a>
              </div>
            </div>
          </template>
        </nav>

        <!-- Right -->
        <div class="flex items-center gap-3">
          <!-- Language Switcher -->
          <div class="relative" @mouseenter="showLang = true" @mouseleave="showLang = false">
            <button class="flex items-center gap-2 text-sm text-[#9A8C7A] hover:text-[#0D1B2A] transition px-3 py-1.5 rounded-lg hover:bg-[#F5F0E8]">
              <span class="text-base leading-none">{{ langFlag }}</span>
              <span class="text-xs font-semibold uppercase">{{ locale }}</span>
              <svg class="w-2.5 h-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" /></svg>
            </button>
            <transition name="fade">
              <div v-if="showLang" class="absolute right-0 top-full mt-2 bg-white border border-[#EAE5DD] rounded-xl shadow-2xl py-1.5 min-w-[180px] z-50">
                <button v-for="l in locales" :key="l.code" @click="switchLocale(l.code)"
                  class="w-full text-left px-4 py-3 text-sm hover:bg-[#FBF9F6] flex items-center gap-3 transition"
                  :class="locale === l.code ? 'text-[#0D1B2A] font-semibold' : 'text-[#9A8C7A]'">
                  <span class="text-lg">{{ l.flag || '🌐' }}</span>                  <span>{{ l.name }}</span>
                  <svg v-if="locale === l.code" class="ml-auto w-4 h-4 text-[#D4A853]" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
                </button>
              </div>
            </transition>
          </div>

          <NuxtLink :to="localePath('/contact')" class="hidden md:inline-flex items-center gap-2 bg-[#D4A853] text-white px-5 py-2.5 rounded-full text-sm font-semibold hover:bg-[#C49A3F] transition-all shadow-sm">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" /></svg>
            {{ $t('home.get_quote') }}
          </NuxtLink>

          <button @click="mobileMenuOpen = !mobileMenuOpen" class="lg:hidden p-2 -mr-2">
            <svg class="w-6 h-6 text-[#0D1B2A]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 6h16M4 12h16M4 18h16" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <!-- Mobile Drawer Menu：Teleport 到 body，避免 header 的 backdrop-filter 使 fixed 相对 header 定位而被裁剪 -->
    <Teleport to="body">
      <transition name="drawer-fade">
        <div v-if="mobileMenuOpen" class="lg:hidden fixed inset-0 z-50 bg-black/40 backdrop-blur-sm" @click="mobileMenuOpen = false"></div>
      </transition>

      <transition name="drawer">
        <div v-if="mobileMenuOpen" class="lg:hidden fixed top-0 right-0 bottom-0 z-[60] w-[86vw] max-w-[92vw] sm:w-[360px] sm:max-w-[380px] md:w-[420px] md:max-w-[440px] bg-white shadow-2xl flex flex-col">
          <div class="flex items-center justify-between px-6 py-4 border-b border-[#EAE5DD]">
            <span class="text-sm font-bold text-[#0D1B2A]">{{ $t('nav.menu') || 'Menu' }}</span>
            <button @click="mobileMenuOpen = false" class="p-2 -mr-2 hover:bg-[#FBF9F6] rounded-lg transition">
              <svg class="w-5 h-5 text-[#0D1B2A]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="flex-1 overflow-y-auto px-4 py-6 space-y-1">
            <template v-for="item in navTree" :key="item.id || item.path">
              <template v-if="!item.children?.length">
                <a :href="navHref(item.path, item.isExternal)"
                  :target="item.target === '_blank' ? '_blank' : undefined"
                  :rel="item.target === '_blank' ? 'noopener noreferrer' : undefined"
                  @click="mobileMenuOpen = false"
                  class="flex items-center gap-3 px-4 py-3.5 rounded-xl text-sm text-[#4A4A4A] hover:bg-[#FBF9F6] active:bg-[#F0EBE3] transition min-h-[48px]">
                  <span>{{ item.label }}</span>
                  <svg class="ml-auto w-4 h-4 text-[#C4B8A8]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </a>
              </template>
              <template v-else>
                <button @click="toggleExpand(item.id || item.path)"
                  class="flex items-center gap-3 px-4 py-3.5 rounded-xl text-sm text-[#4A4A4A] hover:bg-[#FBF9F6] active:bg-[#F0EBE3] transition min-h-[48px] w-full">
                  <span>{{ item.label }}</span>
                  <svg class="ml-auto w-4 h-4 text-[#C4B8A8] transition-transform" :class="expandedItems[item.id || item.path] ? 'rotate-90' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
                <div v-show="expandedItems[item.id || item.path]" class="pl-8 space-y-1 border-l-2 border-[#EAE5DD] ml-6">
                  <a v-for="child in item.children" :key="child.id"
                    :href="navHref(child.path || '/', child.isExternal)"
                    :target="child.target === '_blank' ? '_blank' : undefined"
                    :rel="child.target === '_blank' ? 'noopener noreferrer' : undefined"
                    @click="mobileMenuOpen = false"
                    class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm text-[#4A4A4A] hover:bg-[#FBF9F6] active:bg-[#F0EBE3] transition min-h-[48px]">
                    <span>{{ child.label }}</span>
                  </a>
                </div>
              </template>
            </template>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Main -->
    <main class="flex-1 pb-[52px] md:pb-0">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="bg-[#0D1B2A] text-white">
      <div class="border-b border-[#1B2D44]">
        <div class="max-w-7xl mx-auto px-4 lg:px-8 py-12">
          <div class="grid lg:grid-cols-2 gap-6 items-center">
            <div>
              <h3 class="text-xl font-bold mb-1">{{ $t('footer.newsletter_title') }}</h3>
              <p class="text-[#9A8C7A] text-sm">{{ $t('footer.newsletter_desc') }}</p>
            </div>
            <div>
              <form @submit.prevent="handleNewsletter" class="flex gap-2">
                <input v-model="newsletterEmail" type="email" :placeholder="$t('footer.newsletter_placeholder')" required
                  class="flex-1 px-4 py-3 rounded-xl bg-[#1B2D44] border border-[#2A3F59] text-white placeholder-[#6B5D4B] focus:outline-none focus:border-[#D4A853] transition text-sm" />
                <button type="submit" :disabled="newsletterState === 'loading'"
                  class="px-6 py-3 bg-[#D4A853] text-white rounded-xl font-semibold hover:bg-[#C49A3F] transition text-sm whitespace-nowrap disabled:opacity-60">
                  {{ $t('footer.subscribe') }}
                </button>
              </form>
              <p v-if="newsletterState === 'success'" class="mt-2 text-sm text-emerald-400">{{ $t('footer.newsletter_success') }}</p>
              <p v-else-if="newsletterState === 'error'" class="mt-2 text-sm text-red-400">{{ $t('footer.newsletter_error') }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="max-w-7xl mx-auto px-4 lg:px-8 py-16">
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-10">
          <div>
            <SportswearLogo
              variant="dark"
              :logoSrc="theme.logo_url"
              :logoAlt="theme.logo_alt"
              :brandName="theme.brand_name"
              :brandSubtitle="theme.brand_subtitle"
            />
            <p class="text-[#9A8C7A] text-sm leading-relaxed mb-6 whitespace-pre-line">{{ $t('footer.about_desc') }}</p>
            <div class="flex gap-2 flex-wrap">
              <a v-for="s in socialLinks" :key="s.key" :href="s.url" target="_blank" rel="noopener noreferrer" class="w-9 h-9 bg-[#1B2D44] rounded-full flex items-center justify-center hover:bg-[#D4A853] transition-colors group" :title="s.label" @click.prevent="handleSocialClick(s)">
                <!-- YouTube -->
                <svg v-if="s.icon === 'youtube'" class="w-4 h-4 text-[#9A8C7A] group-hover:text-[#0D1B2A] transition-colors" fill="currentColor" viewBox="0 0 24 24"><path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/></svg>
                <!-- Instagram -->
                <svg v-else-if="s.icon === 'instagram'" class="w-4 h-4 text-[#9A8C7A] group-hover:text-[#0D1B2A] transition-colors" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0C8.74 0 8.333.015 7.053.072 5.775.132 4.905.333 4.14.63c-.789.306-1.459.717-2.126 1.384S.935 3.35.63 4.14C.333 4.905.131 5.775.072 7.053.012 8.333 0 8.74 0 12s.015 3.667.072 4.947c.06 1.277.261 2.148.558 2.913.306.788.717 1.459 1.384 2.126.667.666 1.336 1.079 2.126 1.384.766.296 1.636.499 2.913.558C8.333 23.988 8.74 24 12 24s3.667-.015 4.947-.072c1.277-.06 2.148-.262 2.913-.558.788-.306 1.459-.718 2.126-1.384.666-.667 1.079-1.335 1.384-2.126.296-.765.499-1.636.558-2.913.06-1.28.072-1.687.072-4.947s-.015-3.667-.072-4.947c-.06-1.277-.262-2.149-.558-2.913-.306-.789-.718-1.459-1.384-2.126C21.319 1.347 20.651.935 19.86.63c-.765-.297-1.636-.499-2.913-.558C15.667.012 15.26 0 12 0zm0 2.16c3.203 0 3.585.016 4.85.071 1.17.055 1.805.249 2.227.415.562.217.96.477 1.382.896.419.42.679.819.896 1.381.164.422.36 1.057.413 2.227.057 1.266.07 1.646.07 4.85s-.015 3.585-.074 4.85c-.061 1.17-.256 1.805-.421 2.227-.224.562-.479.96-.899 1.382-.419.419-.824.679-1.38.896-.42.164-1.065.36-2.235.413-1.274.057-1.649.07-4.859.07-3.211 0-3.586-.015-4.859-.074-1.171-.061-1.816-.256-2.236-.421-.569-.224-.96-.479-1.379-.899-.421-.419-.69-.824-.9-1.38-.165-.42-.359-1.065-.42-2.235-.045-1.26-.061-1.649-.061-4.844 0-3.196.016-3.586.061-4.861.061-1.17.255-1.814.42-2.234.21-.57.479-.96.9-1.381.419-.419.81-.689 1.379-.898.42-.166 1.051-.361 2.221-.421 1.275-.045 1.65-.06 4.859-.06l.045.03zm0 3.678a6.162 6.162 0 1 0 0 12.324 6.162 6.162 0 1 0 0-12.324zM12 16c-2.21 0-4-1.79-4-4s1.79-4 4-4 4 1.79 4 4-1.79 4-4 4zm7.846-10.405a1.441 1.441 0 1 1-2.882 0 1.441 1.441 0 0 1 2.882 0z"/></svg>
                <!-- 小红书 -->
                <svg v-else-if="s.icon === 'xiaohongshu'" class="w-4 h-4 text-[#9A8C7A] group-hover:text-[#0D1B2A] transition-colors" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm4.5 14h-9a.5.5 0 0 1-.5-.5v-7a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-.5.5z"/></svg>
                <!-- Facebook -->
                <svg v-else-if="s.icon === 'facebook'" class="w-4 h-4 text-[#9A8C7A] group-hover:text-[#0D1B2A] transition-colors" fill="currentColor" viewBox="0 0 24 24"><path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/></svg>
                <!-- Twitter/X -->
                <svg v-else-if="s.icon === 'twitter'" class="w-4 h-4 text-[#9A8C7A] group-hover:text-[#0D1B2A] transition-colors" fill="currentColor" viewBox="0 0 24 24"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>
                <!-- LinkedIn -->
                <svg v-else-if="s.icon === 'linkedin'" class="w-4 h-4 text-[#9A8C7A] group-hover:text-[#0D1B2A] transition-colors" fill="currentColor" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
              </a>
            </div>
          </div>

          <div>
            <h4 class="text-xs font-semibold uppercase tracking-[0.15em] text-[#D4A853] mb-5">{{ $t('common.quick_links') }}</h4>
            <ul class="space-y-3">
              <li v-for="link in footerLinks" :key="link.path">
                <a :href="navHref(link.path, link.isExternal)" :target="link.target === '_blank' ? '_blank' : undefined"
                  :rel="link.target === '_blank' ? 'noopener noreferrer' : undefined"
                  class="text-sm text-[#9A8C7A] hover:text-white transition">{{ link.label }}</a>
              </li>
            </ul>
          </div>

          <div>
            <h4 class="text-xs font-semibold uppercase tracking-[0.15em] text-[#D4A853] mb-5">{{ $t('footer.products') }}</h4>
            <ul class="space-y-3">
              <li v-for="cat in footerCategories" :key="cat.key">
                <NuxtLink :to="localePath('/products')" class="text-sm text-[#9A8C7A] hover:text-white transition">{{ catLabel(cat) }}</NuxtLink>
              </li>
            </ul>
          </div>

          <div>
            <h4 class="text-xs font-semibold uppercase tracking-[0.15em] text-[#D4A853] mb-5">{{ $t('common.contact_info') }}</h4>
            <ul class="space-y-3 text-sm text-[#9A8C7A]">
              <li class="flex items-center gap-3">
                <svg class="w-4 h-4 flex-shrink-0 text-[#D4A853]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>
                <a href="mailto:info@sportswear.com" class="hover:text-white">info@sportswear.com</a>
              </li>
              <li class="flex items-center gap-3">
                <svg class="w-4 h-4 flex-shrink-0 text-[#D4A853]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg>
                <a href="tel:+8612345678900" class="hover:text-white">+86 123 4567 8900</a>
              </li>
              <li class="flex items-center gap-3">
                <svg class="w-4 h-4 flex-shrink-0 text-[#D4A853]" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                Guangzhou, China
              </li>
            </ul>
          </div>
        </div>
      </div>

      <div class="border-t border-[#1B2D44]">
        <div class="max-w-7xl mx-auto px-4 lg:px-8 py-6 flex flex-col md:flex-row justify-between items-center gap-4">
          <p class="text-sm text-[#6B5D4B]">© {{ new Date().getFullYear() }} Sportswear Manufacturer. {{ $t('common.all_rights_reserved') }}.</p>
          <div class="flex gap-6 text-sm text-[#6B5D4B]">
            <NuxtLink :to="localePath('/faq')" class="hover:text-white transition">FAQ</NuxtLink>
            <NuxtLink :to="localePath('/contact')" class="hover:text-white transition">{{ $t('common.contact_us') }}</NuxtLink>
          </div>
        </div>
      </div>
    </footer>

    <FloatingContact />

    <!-- Cookie Consent Banner（隐私合规） -->
    <CookieConsent />
  </div>
</template>

<script setup lang="ts">
const { locale, locales, setLocale, t, te } = useI18n()
const localePath = useLocalePath()
const switchLocalePath = useSwitchLocalePath()
const showLang = ref(false)
const mobileMenuOpen = ref(false)
const newsletterEmail = ref('')
const newsletterState = ref<'idle' | 'loading' | 'success' | 'error'>('idle')
const theme = ref<Record<string, any>>({})

// 默认导航（API 失败时的回退值）
const DEFAULT_NAV = [
  { path: '/', label: 'nav.home', target: '_self', isExternal: false },
  { path: '/products', label: 'nav.products', target: '_self', isExternal: false },
  { path: '/about', label: 'nav.about', target: '_self', isExternal: false },
  { path: '/cases', label: 'nav.cases', target: '_self', isExternal: false },
  { path: '/blog', label: 'nav.blog', target: '_self', isExternal: false },
  { path: '/faq', label: 'nav.faq', target: '_self', isExternal: false },
  { path: '/contact', label: 'nav.contact', target: '_self', isExternal: false },
]

const DEFAULT_FOOTER = [
  { path: '/', label: 'nav.home', target: '_self', isExternal: false },
  { path: '/products', label: 'nav.products', target: '_self', isExternal: false },
  { path: '/about', label: 'nav.about', target: '_self', isExternal: false },
  { path: '/cases', label: 'nav.cases', target: '_self', isExternal: false },
  { path: '/blog', label: 'nav.blog', target: '_self', isExternal: false },
  { path: '/faq', label: 'nav.faq', target: '_self', isExternal: false },
]

// 动态导航数据（SSR + 客户端）
const { getNavigations, getTheme, trackClick, subscribe } = useApi()
const { data: headerNavData } = await useAsyncData<any[]>(
  'nav-header',
  () => getNavigations('header').catch(() => null),
  {
    getCachedData(key, nuxtApp) {
      return nuxtApp.isHydrating ? nuxtApp.payload.data[key] : undefined
    },
  },
)
const { data: footerNavData } = await useAsyncData<any[]>(
  'nav-footer',
  () => getNavigations('footer').catch(() => null),
  {
    getCachedData(key, nuxtApp) {
      return nuxtApp.isHydrating ? nuxtApp.payload.data[key] : undefined
    },
  },
)

function navToItem(n: any) {
  return { path: n.url || '/', label: navLabel(n.name, n.url), target: n.target || '_self', isExternal: isExternalUrl(n.url) }
}

function isExternalUrl(url: string): boolean {
  if (!url) return false
  return /^(https?:|mailto:|tel:|#)/i.test(url)
}

// 导航多语言：后台「导航管理」的 name 是运营维护的展示文案（通常只有英文），
// 这里优先按 URL 映射到 nav.* 词条 key 取多语言文案（静态语言包 + 后台词条管理均可覆盖）；
// 无对应词条时回退 DB 原文，保证自定义导航项不丢文案。
const NAV_KEY_BY_PATH: Record<string, string> = {
  '/': 'nav.home',
  '/products': 'nav.products',
  '/about': 'nav.about',
  '/cases': 'nav.cases',
  '/blog': 'nav.blog',
  '/faq': 'nav.faq',
  '/contact': 'nav.contact',
  '/oem': 'nav.oem',
  '/odm': 'nav.odm',
  '/factory': 'nav.factory',
}

function navLabel(name: string, url?: string): string {
  const rawName = (name || '').trim()
  // name 本身就是词条 key（如 nav.home）时直接翻译
  if (/^nav\.[\w.]+$/.test(rawName)) return te(rawName) ? t(rawName) : rawName
  const path = (url || '').replace(/\/+$/, '') || '/'
  const key = NAV_KEY_BY_PATH[path] || ''
  return key && te(key) ? t(key) : rawName
}

function flattenNav(items: any[]): any[] {
  const result: any[] = []
  for (const item of items) {
    if (item.is_visible !== false) result.push(item)
    if (item.children?.length) result.push(...flattenNav(item.children))
  }
  return result
}

// 将后端导航树（name/url/children）映射为模板所需的 {label, path, children} 结构，
// 与移动端 navToItem 保持一致，避免 $t(undefined) 触发 vue-i18n INVALID_ARGUMENT（SSR 500）。
function mapNavTree(items: any[]): any[] {
  return (items || [])
    .filter((n: any) => n.is_visible !== false)
    .map((n: any) => ({
      id: n.id,
      label: navLabel(n.name, n.url),
      path: n.url || '/',
      target: n.target || '_self',
      isExternal: isExternalUrl(n.url),
      children: n.children?.length ? mapNavTree(n.children) : [],
    }))
}

// 导航链接地址解析：外部 URL / 锚点直接返回，内部路径加语言前缀
function navHref(path: string, isExternal: boolean): string {
  return isExternal ? path : localePath(path)
}

// 桌面端：保留树结构以支持下拉菜单
const navTree = computed(() => {
  if (headerNavData.value?.length) return mapNavTree(headerNavData.value)
  return DEFAULT_NAV.map(n => ({ id: n.path, label: navLabel(n.label, n.path), path: n.path, target: n.target, isExternal: n.isExternal, children: [] }))
})

// 移动端：手风琴展开/折叠状态
const expandedItems = reactive<Record<string, boolean>>({})
function toggleExpand(key: string) {
  expandedItems[key] = !expandedItems[key]
}

// 移动端 / Footer：展平为列表
const navItems = computed(() => {
  if (headerNavData.value?.length) return flattenNav(headerNavData.value).map(navToItem)
  return DEFAULT_NAV.map(n => ({ ...n, label: navLabel(n.label, n.path) }))
})

const footerLinks = computed(() => {
  if (footerNavData.value?.length) return flattenNav(footerNavData.value).map(navToItem)
  return DEFAULT_FOOTER.map(n => ({ ...n, label: navLabel(n.label, n.path) }))
})

// Footer 产品分类：走词条（footer.cat_*），后台词条管理可覆盖；无词条时回退英文默认
const footerCategories = [
  { key: 'footer.cat_yoga_wear', fallback: 'Yoga Wear' },
  { key: 'footer.cat_running_gear', fallback: 'Running Gear' },
  { key: 'footer.cat_training_apparel', fallback: 'Training Apparel' },
  { key: 'footer.cat_team_uniforms', fallback: 'Team Uniforms' },
  { key: 'footer.cat_custom_design', fallback: 'Custom Design' },
]
const catLabel = (c: { key: string; fallback: string }) => (te(c.key) ? t(c.key) : c.fallback)

const socialLinks = computed(() => [
  { key: 'social_youtube', label: 'YouTube', icon: 'youtube', url: theme.value.social_youtube || '#' },
  { key: 'social_instagram', label: 'Instagram', icon: 'instagram', url: theme.value.social_instagram || '#' },
  { key: 'social_xiaohongshu', label: '小红书', icon: 'xiaohongshu', url: theme.value.social_xiaohongshu || '#' },
  { key: 'social_facebook', label: 'Facebook', icon: 'facebook', url: theme.value.social_facebook || '#' },
  { key: 'social_twitter', label: 'Twitter/X', icon: 'twitter', url: theme.value.social_twitter || '#' },
  { key: 'social_linkedin', label: 'LinkedIn', icon: 'linkedin', url: theme.value.social_linkedin || '#' },
])

const langFlag = computed(() => {
  const current = (locales.value || []).find((l) => l.code === locale.value)
  return current?.flag || '🌐'
})

function switchLocale(code: string) {
  showLang.value = false
  // 跳转到对应语言前缀路由，触发页面重新渲染 + 数据按新语言重新拉取；
  // 同时保证 URL/hreflang/SEO 与当前语言一致。
  const target = switchLocalePath(code)
  if (target) {
    navigateTo(target)
  } else {
    setLocale(code)
  }
}

async function handleNewsletter() {
  const email = newsletterEmail.value.trim()
  if (!email) return
  newsletterState.value = 'loading'
  try {
    await subscribe(email)
    newsletterEmail.value = ''
    newsletterState.value = 'success'
  } catch {
    newsletterState.value = 'error'
  }
}

// 加载主题配置（含社交链接和Logo）
onMounted(async () => {
  try {
    theme.value = await getTheme()
    // 如果主题配置了 favicon_url，动态更新浏览器标签页图标
    if (theme.value.favicon_url) {
      const existing = document.querySelector('link[rel="icon"]')
      if (existing) {
        existing.setAttribute('href', theme.value.favicon_url)
      } else {
        const link = document.createElement('link')
        link.rel = 'icon'
        link.href = theme.value.favicon_url
        document.head.appendChild(link)
      }
    }
  } catch (e) {
    console.error('Failed to load theme config', e)
  }
})

// 社交链接点击跟踪：先记录点击事件，再跳转
function handleSocialClick(s: { key: string; icon: string; url: string }) {
  if (s.url && s.url !== '#') {
    // 异步记录点击，不阻塞跳转
    trackClick(s.icon, s.url).catch(() => {})
    window.open(s.url, '_blank', 'noopener,noreferrer')
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.15s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.drawer-fade-enter-active, .drawer-fade-leave-active { transition: opacity 0.3s ease; }
.drawer-fade-enter-from, .drawer-fade-leave-to { opacity: 0; }

.drawer-enter-active { transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.drawer-leave-active { transition: transform 0.2s ease-in; }
.drawer-enter-from, .drawer-leave-to { transform: translateX(100%); }
</style>
