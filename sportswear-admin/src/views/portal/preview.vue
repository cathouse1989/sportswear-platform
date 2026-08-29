<template>
  <div class="portal-preview">
    <!-- 预览工具栏 -->
    <div class="preview-toolbar">
      <div class="toolbar-left">
        <span class="toolbar-title">门户预览</span>
        <el-tag type="warning" size="small">预览模式 - 数据来自公开 API</el-tag>
      </div>
      <div class="toolbar-right">
        <el-select v-model="previewLang" size="small" style="width: 120px" @change="refreshPreview">
          <el-option v-for="l in languages" :key="l.code" :label="l.native_name || l.name" :value="l.code" />
        </el-select>
        <el-select v-model="previewPage" size="small" style="width: 160px" @change="onPageChange">
          <el-option label="首页" value="home" />
          <el-option label="产品列表" value="products" />
          <el-option label="产品详情" value="product-detail" />
          <el-option label="博客列表" value="blogs" />
          <el-option label="案例列表" value="cases" />
          <el-option label="FAQ" value="faqs" />
          <el-option label="工厂" value="factories" />
          <el-option label="认证" value="certifications" />
          <el-option label="生产流程" value="processes" />
        </el-select>
        <el-input
          v-if="previewPage === 'product-detail'"
          v-model="productSlug"
          placeholder="输入产品 Slug"
          size="small"
          style="width: 160px"
          clearable
          @keyup.enter="refreshPreview"
        />
        <el-button size="small" type="primary" @click="refreshPreview">刷新</el-button>
        <el-button size="small" @click="openInNewTab">新窗口打开</el-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="preview-loading">
      <el-skeleton :rows="10" animated />
    </div>

    <!-- 门户预览内容 -->
    <div v-else class="preview-content" :style="{ background: '#FBF9F6' }">
      <!-- ==================== 首页预览 ==================== -->
      <template v-if="previewPage === 'home' && homeData">
        <!-- Hero Section -->
        <section class="portal-hero" :style="{ background: `linear-gradient(135deg, ${theme.primary_color || '#0D1B2A'}, ${theme.secondary_color || '#1B2D44'})` }">
          <div class="portal-hero-bg">
            <div class="portal-hero-pattern"></div>
          </div>
          <div v-for="(s, idx) in heroSlides" :key="idx"
            class="portal-hero-slide" :class="{ 'portal-hero-slide-active': heroCur === idx }">
            <img v-if="s.image" :src="s.image" :alt="s.title" class="portal-hero-img" />
          </div>
          <div class="portal-hero-content">
            <div class="portal-hero-badge">
              <span class="portal-hero-badge-dot"></span>
              {{ t('home.hero_label') }}
            </div>
            <h1 class="portal-hero-title">{{ heroSlides[heroCur]?.title || homeData.page?.title || t('home.hero_title') }}</h1>
            <p v-if="heroSlides[heroCur]?.subtitle || homeData.page?.translations?.[0]?.content" class="portal-hero-subtitle">{{ heroSlides[heroCur]?.subtitle || homeData.page?.translations?.[0]?.content || t('home.hero_subtitle') }}</p>
            <div class="portal-hero-actions">
              <a class="portal-btn portal-btn-primary" @click="openInNewTab">
                {{ heroSlides[heroCur]?.button_text || t('home.get_quote') }}
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
              </a>
            </div>
          </div>
          <div v-if="heroSlides.length > 1" class="portal-hero-dots">
            <button v-for="(s, idx) in heroSlides" :key="idx" @click="heroCur = idx"
              :class="heroCur === idx ? 'portal-hero-dot portal-hero-dot-active' : 'portal-hero-dot'"></button>
          </div>
          <div class="portal-scroll-indicator">
            <span>Scroll</span>
            <div class="portal-scroll-line"></div>
          </div>
        </section>

        <!-- Trust Bar -->
        <section class="portal-trust-bar">
          <div class="portal-container">
            <div class="portal-trust-items">
              <span class="portal-trust-item"><svg width="16" height="16" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg> ISO 9001</span>
              <span class="portal-trust-item"><svg width="16" height="16" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg> BSCI</span>
              <span class="portal-trust-item"><svg width="16" height="16" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg> OEKO-TEX</span>
              <span class="portal-trust-item"><svg width="16" height="16" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg> 15+ Yrs Exp</span>
            </div>
          </div>
        </section>

        <!-- Featured Products -->
        <section class="portal-section">
          <div class="portal-container">
            <div class="portal-section-header">
              <div>
                <span class="portal-section-label">Products</span>
                <h2 class="portal-section-title">{{ t('home.featured_products') }}</h2>
              </div>
              <a class="portal-view-all" @click="openInNewTab">{{ t('home.view_all') }} →</a>
            </div>
            <div v-if="!homeData.featured_products?.length" class="portal-empty-state">
              <div class="portal-empty-icon">📦</div>
              <p>{{ t('common.no_data') }}</p>
            </div>
            <div v-else class="portal-product-grid">
              <div v-for="p in homeData.featured_products?.slice(0, 4)" :key="p.id" class="portal-product-card">
                <div class="portal-product-image">
                  <img v-if="p.cover_image" :src="p.cover_image" :alt="p.sku" />
                  <div v-else class="portal-placeholder-img">🏋️</div>
                  <div class="portal-product-overlay">
                    <span class="portal-product-type">{{ p.type || p.category }}</span>
                    <span class="portal-product-inquiry">{{ t('product.inquiry_now') }}</span>
                  </div>
                </div>
                <div class="portal-product-info">
                  <div class="portal-product-meta">
                    <h3 class="portal-product-name">{{ p.sku }}</h3>
                    <span class="portal-product-gender">{{ p.gender }}</span>
                  </div>
                  <p class="portal-product-brief">{{ p.brief || p.description?.slice(0, 80) }}</p>
                  <div class="portal-product-tags">
                    <span v-if="p.material" class="portal-tag">🧵 {{ p.material }}</span>
                    <span v-if="p.sample_moq" class="portal-tag">MOQ: {{ p.sample_moq }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Why Choose Us -->
        <section class="portal-section portal-section-alt">
          <div class="portal-container">
            <div class="portal-section-header portal-section-center">
              <span class="portal-section-label">{{ t('home.our_strength') }}</span>
              <h2 class="portal-section-title">{{ t('home.why_choose_us') }}</h2>
              <p class="portal-section-desc">{{ t('home.why_choose_us_desc') }}</p>
            </div>
            <div class="portal-strength-grid">
              <div v-for="(item, i) in strengths" :key="i" class="portal-strength-card">
                <div class="portal-strength-icon" v-html="item.icon"></div>
                <h3 class="portal-strength-title">{{ t(item.title) }}</h3>
                <p class="portal-strength-desc">{{ t(item.desc) }}</p>
              </div>
            </div>
          </div>
        </section>

        <!-- Categories -->
        <section v-if="homeData.categories?.length" class="portal-section">
          <div class="portal-container">
            <div class="portal-section-header portal-section-center">
              <span class="portal-section-label">Categories</span>
              <h2 class="portal-section-title">{{ t('home.categories') }}</h2>
            </div>
            <div class="portal-category-list">
              <a v-for="c in homeData.categories" :key="c.id" class="portal-category-tag" @click="openInNewTab">{{ c.name }}</a>
            </div>
          </div>
        </section>

        <!-- Blog -->
        <section v-if="homeData.blogs?.length" class="portal-section portal-section-alt">
          <div class="portal-container">
            <div class="portal-section-header">
              <div>
                <span class="portal-section-label">Insights</span>
                <h2 class="portal-section-title">{{ t('home.latest_blog') }}</h2>
              </div>
              <a class="portal-view-all" @click="openInNewTab">{{ t('home.view_all') }} →</a>
            </div>
            <div class="portal-blog-grid">
              <div v-for="b in homeData.blogs?.slice(0, 3)" :key="b.id" class="portal-blog-card">
                <div class="portal-blog-image">
                  <img v-if="b.cover_image" :src="b.cover_image" :alt="b.title" />
                  <div v-else class="portal-placeholder-img">📝</div>
                </div>
                <div class="portal-blog-content">
                  <div class="portal-blog-meta">
                    <span class="portal-blog-category">{{ b.category }}</span>
                    <span class="portal-blog-date">{{ b.created_at?.slice(0, 10) }}</span>
                  </div>
                  <h3 class="portal-blog-title">{{ b.title }}</h3>
                  <p class="portal-blog-excerpt">{{ b.content?.slice(0, 120) }}</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Certifications -->
        <section v-if="homeData.certifications?.length" class="portal-section">
          <div class="portal-container">
            <div class="portal-section-header portal-section-center">
              <span class="portal-section-label">Compliance</span>
              <h2 class="portal-section-title">{{ t('home.certifications') }}</h2>
            </div>
            <div class="portal-cert-grid">
              <div v-for="c in homeData.certifications" :key="c.id" class="portal-cert-card">
                <span class="portal-cert-icon">🏆</span>
                <div>
                  <div class="portal-cert-name">{{ c.name }}</div>
                  <div class="portal-cert-desc">{{ c.description?.slice(0, 60) }}</div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- CTA -->
        <section class="portal-cta" :style="{ background: theme.primary_color || '#0D1B2A' }">
          <div class="portal-cta-glow"></div>
          <div class="portal-cta-content">
            <h2 class="portal-cta-title">{{ t('home.cta_title') }}</h2>
            <p class="portal-cta-desc">{{ t('home.cta_desc') }}</p>
            <a class="portal-btn portal-btn-primary portal-btn-lg" @click="openInNewTab">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/></svg>
              {{ t('home.get_quote') }}
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
            </a>
          </div>
        </section>

        <!-- Footer -->
        <footer class="portal-footer" :style="{ background: theme.primary_color || '#0D1B2A' }">
          <div class="portal-footer-newsletter">
            <div class="portal-container">
              <div class="portal-newsletter-content">
                <div>
                  <h3>{{ t('footer.newsletter_title') }}</h3>
                  <p>{{ t('footer.newsletter_desc') }}</p>
                </div>
                <div class="portal-newsletter-form">
                  <input type="email" :placeholder="t('footer.newsletter_placeholder')" />
                  <button>{{ t('footer.subscribe') }}</button>
                </div>
              </div>
            </div>
          </div>
          <div class="portal-container">
            <div class="portal-footer-grid">
              <div class="portal-footer-col">
                <div class="portal-footer-logo">
                  <div class="portal-footer-logo-icon">S</div>
                  <div class="portal-footer-logo-text">
                    <span>SPORTSWEAR</span>
                    <span>Premium Mfg.</span>
                  </div>
                </div>
                <p class="portal-footer-about">Professional sportswear OEM/ODM manufacturer with 15+ years of experience. ISO 9001, BSCI, OEKO-TEX certified.</p>
                <div class="portal-footer-social">
                  <a v-for="s in socialLinks" :key="s.key" :href="s.url" target="_blank" rel="noopener noreferrer" class="portal-social-icon" :title="s.label" @click.prevent="handleSocialClick(s)">
                    <!-- YouTube -->
                    <svg v-if="s.icon === 'youtube'" width="16" height="16" fill="currentColor" viewBox="0 0 24 24"><path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/></svg>
                    <!-- Instagram -->
                    <svg v-else-if="s.icon === 'instagram'" width="16" height="16" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0C8.74 0 8.333.015 7.053.072 5.775.132 4.905.333 4.14.63c-.789.306-1.459.717-2.126 1.384S.935 3.35.63 4.14C.333 4.905.131 5.775.072 7.053.012 8.333 0 8.74 0 12s.015 3.667.072 4.947c.06 1.277.261 2.148.558 2.913.306.788.717 1.459 1.384 2.126.667.666 1.336 1.079 2.126 1.384.766.296 1.636.499 2.913.558C8.333 23.988 8.74 24 12 24s3.667-.015 4.947-.072c1.277-.06 2.148-.262 2.913-.558.788-.306 1.459-.718 2.126-1.384.666-.667 1.079-1.335 1.384-2.126.296-.765.499-1.636.558-2.913.06-1.28.072-1.687.072-4.947s-.015-3.667-.072-4.947c-.06-1.277-.262-2.149-.558-2.913-.306-.789-.718-1.459-1.384-2.126C21.319 1.347 20.651.935 19.86.63c-.765-.297-1.636-.499-2.913-.558C15.667.012 15.26 0 12 0zm0 2.16c3.203 0 3.585.016 4.85.071 1.17.055 1.805.249 2.227.415.562.217.96.477 1.382.896.419.42.679.819.896 1.381.164.422.36 1.057.413 2.227.057 1.266.07 1.646.07 4.85s-.015 3.585-.074 4.85c-.061 1.17-.256 1.805-.421 2.227-.224.562-.479.96-.899 1.382-.419.419-.824.679-1.38.896-.42.164-1.065.36-2.235.413-1.274.057-1.649.07-4.859.07-3.211 0-3.586-.015-4.859-.074-1.171-.061-1.816-.256-2.236-.421-.569-.224-.96-.479-1.379-.899-.421-.419-.69-.824-.9-1.38-.165-.42-.359-1.065-.42-2.235-.045-1.26-.061-1.649-.061-4.844 0-3.196.016-3.586.061-4.861.061-1.17.255-1.814.42-2.234.21-.57.479-.96.9-1.381.419-.419.81-.689 1.379-.898.42-.166 1.051-.361 2.221-.421 1.275-.045 1.65-.06 4.859-.06l.045.03zm0 3.678a6.162 6.162 0 1 0 0 12.324 6.162 6.162 0 1 0 0-12.324zM12 16c-2.21 0-4-1.79-4-4s1.79-4 4-4 4 1.79 4 4-1.79 4-4 4zm7.846-10.405a1.441 1.441 0 1 1-2.882 0 1.441 1.441 0 0 1 2.882 0z"/></svg>
                    <!-- 小红书 -->
                    <svg v-else-if="s.icon === 'xiaohongshu'" width="16" height="16" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm4.5 14h-9a.5.5 0 0 1-.5-.5v-7a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 .5.5v7a.5.5 0 0 1-.5.5z"/></svg>
                    <!-- Facebook -->
                    <svg v-else-if="s.icon === 'facebook'" width="16" height="16" fill="currentColor" viewBox="0 0 24 24"><path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/></svg>
                    <!-- Twitter/X -->
                    <svg v-else-if="s.icon === 'twitter'" width="16" height="16" fill="currentColor" viewBox="0 0 24 24"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>
                    <!-- LinkedIn -->
                    <svg v-else-if="s.icon === 'linkedin'" width="16" height="16" fill="currentColor" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
                  </a>
                </div>
              </div>
              <div class="portal-footer-col">
                <h4>{{ t('common.quick_links') }}</h4>
                <ul>
                  <li v-for="link in footerLinks" :key="link.path"><a @click="openInNewTab">{{ link.label }}</a></li>
                </ul>
              </div>
              <div class="portal-footer-col">
                <h4>{{ t('footer.products') }}</h4>
                <ul>
                  <li v-for="cat in footerCategories" :key="cat"><a @click="openInNewTab">{{ cat }}</a></li>
                </ul>
              </div>
              <div class="portal-footer-col">
                <h4>{{ t('common.contact_info') }}</h4>
                <ul class="portal-contact-list">
                  <li><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/></svg> info@sportswear.com</li>
                  <li><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/></svg> +86 123 4567 8900</li>
                  <li><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/><path d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/></svg> Guangzhou, China</li>
                </ul>
              </div>
            </div>
          </div>
          <div class="portal-footer-bottom">
            <div class="portal-container">
              <p>© 2026 Sportswear Manufacturer. {{ t('common.all_rights_reserved') }}.</p>
              <div class="portal-footer-links">
                <a @click="openInNewTab">FAQ</a>
                <a @click="openInNewTab">{{ t('common.contact_us') }}</a>
              </div>
            </div>
          </div>
        </footer>
      </template>

      <!-- ==================== 产品详情预览 ==================== -->
      <template v-else-if="previewPage === 'product-detail'">
        <section class="portal-section">
          <div class="portal-container">
            <div v-if="!productSlug" class="portal-empty-state">
              <div class="portal-empty-icon">🔍</div>
              <p>请输入产品 Slug 进行预览</p>
            </div>
            <div v-else-if="!productDetail" class="portal-empty-state">
              <div class="portal-empty-icon">📦</div>
              <p>未找到产品「{{ productSlug }}」</p>
            </div>
            <template v-else>
              <!-- 产品详情头部 -->
              <div class="pd-breadcrumb">
                <a @click="previewPage = 'products'; refreshPreview()">{{ t('product.products') }}</a>
                <span>/</span>
                <span>{{ productDetail.sku }}</span>
              </div>

              <div class="pd-layout">
                <!-- 左侧图片区 -->
                <div class="pd-gallery">
                  <div class="pd-main-image">
                    <img v-if="productDetail.cover_image" :src="productDetail.cover_image" :alt="productDetail.sku" />
                    <div v-else class="portal-placeholder-img" style="font-size:80px">📷</div>
                  </div>
                  <div v-if="productDetail.images?.length" class="pd-thumb-list">
                    <div
                      v-for="(img, idx) in productDetail.images"
                      :key="img.id"
                      class="pd-thumb-item"
                      :class="{ active: idx === activeImageIdx }"
                      @click="activeImageIdx = Number(idx)"
                    >
                      <img :src="img.url" :alt="'image ' + idx" />
                    </div>
                  </div>
                </div>

                <!-- 右侧信息区 -->
                <div class="pd-info">
                  <div class="pd-badges">
                    <span class="pd-badge pd-badge-type">{{ productDetail.type?.toUpperCase() }}</span>
                    <span v-if="productDetail.status === 'published'" class="pd-badge pd-badge-status">{{ t('common.published') }}</span>
                    <span v-if="productDetail.is_featured" class="pd-badge pd-badge-featured">Featured</span>
                    <span v-if="productDetail.is_new" class="pd-badge pd-badge-new">New</span>
                  </div>
                  <h1 class="pd-title">{{ productDetail.sku }}</h1>
                  <p class="pd-brief">{{ productDetail.brief || productDetail.description?.slice(0, 200) }}</p>

                  <!-- 规格参数 -->
                  <div class="pd-specs">
                    <div class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.gender') }}</span>
                      <span class="pd-spec-value">{{ productDetail.gender }}</span>
                    </div>
                    <div class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.material') }}</span>
                      <span class="pd-spec-value">{{ productDetail.material || '-' }}</span>
                    </div>
                    <div v-if="productDetail.composition" class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.composition') }}</span>
                      <span class="pd-spec-value">{{ productDetail.composition }}</span>
                    </div>
                    <div class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.sample_moq') }}</span>
                      <span class="pd-spec-value">{{ productDetail.sample_moq || '-' }}</span>
                    </div>
                    <div class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.production_moq') }}</span>
                      <span class="pd-spec-value">{{ productDetail.production_moq || '-' }}</span>
                    </div>
                    <div v-if="productDetail.weight" class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.weight') }}</span>
                      <span class="pd-spec-value">{{ productDetail.weight }}</span>
                    </div>
                    <div v-if="productDetail.season" class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.season') }}</span>
                      <span class="pd-spec-value">{{ productDetail.season }}</span>
                    </div>
                    <div v-if="productDetail.fit" class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.fit') }}</span>
                      <span class="pd-spec-value">{{ productDetail.fit }}</span>
                    </div>
                    <div v-if="productDetail.size_range" class="pd-spec">
                      <span class="pd-spec-label">{{ t('product.size_range') }}</span>
                      <span class="pd-spec-value">{{ productDetail.size_range }}</span>
                    </div>
                  </div>

                  <!-- 询盘按钮 -->
                  <div class="pd-actions">
                    <a class="portal-btn portal-btn-primary" @click="openInNewTab">
                      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/></svg>
                      {{ t('product.inquiry_now') }}
                    </a>
                  </div>
                </div>
              </div>

              <!-- 详细描述 -->
              <div v-if="productDetail.description" class="pd-description">
                <h3>{{ t('product.description') }}</h3>
                <div class="pd-description-content">{{ productDetail.description }}</div>
              </div>

              <!-- 自定义选项 -->
              <div v-if="productDetail.customizations?.length" class="pd-section">
                <h3>{{ t('product.customizations') }}</h3>
                <div class="pd-custom-grid">
                  <div v-for="c in productDetail.customizations" :key="c.id" class="pd-custom-item">
                    <span class="pd-custom-icon">{{ c.is_enabled ? '✅' : '❌' }}</span>
                    <div>
                      <div class="pd-custom-name">{{ c.type }}</div>
                      <div v-if="c.note" class="pd-custom-note">{{ c.note }}</div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 规格表 -->
              <div v-if="productDetail.specs?.length" class="pd-section">
                <h3>{{ t('product.specifications') }}</h3>
                <div class="pd-specs-grid">
                  <div v-for="s in productDetail.specs" :key="s.id" class="pd-specs-row">
                    <span class="pd-specs-name">{{ s.name }}</span>
                    <span class="pd-specs-val">{{ s.value }}</span>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </section>
      </template>

      <!-- ==================== 产品列表预览 ==================== -->
      <template v-else-if="previewPage === 'products' && productList">
        <section class="portal-section">
          <div class="portal-container">
            <div class="portal-section-header">
              <div>
                <span class="portal-section-label">{{ t('product.products') }}</span>
                <h2 class="portal-section-title">{{ t('product.products') }} ({{ productList.total }})</h2>
              </div>
            </div>
            <div v-if="!productList.items?.length" class="portal-empty-state">
              <div class="portal-empty-icon">📦</div>
              <p>{{ t('product.no_products') }}</p>
            </div>
            <div v-else class="portal-product-grid">
              <div v-for="p in productList.items" :key="p.id" class="portal-product-card">
                <div class="portal-product-image">
                  <img v-if="p.cover_image" :src="p.cover_image" :alt="p.sku" />
                  <div v-else class="portal-placeholder-img">📷</div>
                  <div class="portal-product-overlay">
                    <span class="portal-product-type">{{ p.type }}</span>
                    <span v-if="p.status === 'published'" class="portal-product-status">{{ t('common.published') }}</span>
                  </div>
                </div>
                <div class="portal-product-info">
                  <div class="portal-product-meta">
                    <h3 class="portal-product-name">{{ p.sku }}</h3>
                    <span class="portal-product-gender">{{ p.gender }}</span>
                  </div>
                  <p class="portal-product-brief">{{ p.brief || p.description?.slice(0, 80) }}</p>
                  <div class="portal-product-tags">
                    <span v-if="p.material" class="portal-tag">🧵 {{ p.material }}</span>
                    <span v-if="p.sample_moq" class="portal-tag">MOQ: {{ p.sample_moq }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </template>

      <!-- ==================== 其他列表预览 ==================== -->
      <template v-else>
        <section class="portal-section">
          <div class="portal-container">
            <div class="portal-section-header">
              <h2 class="portal-section-title">{{ pageTitle }}</h2>
            </div>
            <div v-if="!listData?.length" class="portal-empty-state">
              <div class="portal-empty-icon">📋</div>
              <p>{{ t('common.no_data') }}</p>
            </div>
            <div v-else class="portal-list-view">
              <div v-for="item in listData" :key="item.id" class="portal-list-item">
                <div class="portal-list-item-content">
                  <h4>{{ item.title || item.name || item.question }}</h4>
                  <p v-if="item.content || item.description || item.answer" class="portal-list-item-desc">{{ (item.content || item.description || item.answer)?.slice(0, 160) }}</p>
                  <div class="portal-list-item-tags">
                    <span v-if="item.category" class="portal-tag">{{ item.category }}</span>
                    <span v-if="item.status === 'published'" class="portal-tag portal-tag-success">已发布</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { publicApi } from '@/api'
import { ElMessage } from 'element-plus'

const previewLang = ref('en')
const previewPage = ref('home')
const loading = ref(false)
const languages = ref<Array<{ code: string; name: string; native_name: string }>>([])
const theme = ref<Record<string, any>>({})

// 数据
const homeData = ref<any>(null)
const productList = ref<any>(null)
const productDetail = ref<any>(null)
const listData = ref<any[]>([])
const productSlug = ref('')
const activeImageIdx = ref(0)

// 首页轮播图（数据来自 /public/home 的 hero_slides，与门户渲染保持一致）
const heroCur = ref(0)
const heroSlides = computed(() => (homeData.value?.hero_slides || []).map((s: any) => ({
  image: s?.image || '', title: s?.title || '', subtitle: s?.subtitle || '',
  button_text: s?.button_text || '', button_url: s?.button_url || '',
})))

function onPageChange() {
  productSlug.value = ''
  productDetail.value = null
  activeImageIdx.value = 0
  heroCur.value = 0
  refreshPreview()
}

const pageTitleMap: Record<string, string> = {
  blogs: 'Blogs', cases: 'Cases', faqs: 'FAQs',
  factories: 'Factories', certifications: 'Certifications', processes: 'Production Processes',
}
const pageTitle = ref('')

// 从后端 API 获取 i18n 词条（与门户保持同步）
const translations = ref<Record<string, string>>({})

async function loadTranslations() {
  try {
    const result = await publicApi.i18n(previewLang.value)
    translations.value = result.dictionary || result || {}
  } catch {
    translations.value = {}
  }
}

function t(key: string): string {
  return translations.value?.[key] || key
}

const strengths = [
  {
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>',
    title: 'home.quality_title', desc: 'home.quality_desc',
  },
  {
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M13 10V3L4 14h7v7l9-11h-7z"/></svg>',
    title: 'home.capacity_title', desc: 'home.capacity_desc',
  },
  {
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>',
    title: 'home.custom_title', desc: 'home.custom_desc',
  },
  {
    icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>',
    title: 'home.global_title', desc: 'home.global_desc',
  },
]

const footerLinks = [
  { path: '/', label: 'Home' },
  { path: '/products', label: 'Products' },
  { path: '/cases', label: 'Cases' },
  { path: '/blog', label: 'Blog' },
  { path: '/faq', label: 'FAQ' },
]

const footerCategories = ['Yoga Wear', 'Running Gear', 'Training Apparel', 'Team Uniforms', 'Custom Design']

const socialLinks = computed(() => [
  { key: 'social_youtube', label: 'YouTube', icon: 'youtube', url: theme.value.social_youtube || '#' },
  { key: 'social_instagram', label: 'Instagram', icon: 'instagram', url: theme.value.social_instagram || '#' },
  { key: 'social_xiaohongshu', label: '小红书', icon: 'xiaohongshu', url: theme.value.social_xiaohongshu || '#' },
  { key: 'social_facebook', label: 'Facebook', icon: 'facebook', url: theme.value.social_facebook || '#' },
  { key: 'social_twitter', label: 'Twitter/X', icon: 'twitter', url: theme.value.social_twitter || '#' },
  { key: 'social_linkedin', label: 'LinkedIn', icon: 'linkedin', url: theme.value.social_linkedin || '#' },
])

async function loadLanguages() {
  try {
    const result = await publicApi.languages()
    languages.value = result
  } catch {
    languages.value = [
      { code: 'en', name: 'English', native_name: 'English' },
      { code: 'zh', name: '中文', native_name: '中文' },
      { code: 'es', name: 'Español', native_name: 'Español' },
      { code: 'fr', name: 'Français', native_name: 'Français' },
    ]
  }
}

async function loadTheme() {
  try {
    const result = await publicApi.theme()
    theme.value = result
  } catch {}
}

async function refreshPreview() {
  loading.value = true
  try {
    // 同步加载 i18n 词条
    await loadTranslations()

    if (previewPage.value === 'home') {
      heroCur.value = 0
      homeData.value = await publicApi.home(previewLang.value)
    } else if (previewPage.value === 'products') {
      productList.value = await publicApi.products({ lang: previewLang.value, page: 1, pageSize: 8 })
    } else if (previewPage.value === 'product-detail') {
      if (productSlug.value) {
        try {
          productDetail.value = await publicApi.product(productSlug.value, previewLang.value)
        } catch {
          productDetail.value = null
        }
      }
    } else {
      const apiMap: Record<string, string> = {
        blogs: '/public/blogs', cases: '/public/cases', faqs: '/public/faqs',
        factories: '/public/factories', certifications: '/public/certifications', processes: '/public/production-processes',
      }
      pageTitle.value = pageTitleMap[previewPage.value] || ''
      const result = await publicApi.fetch(apiMap[previewPage.value] || '/public/home', previewLang.value)
      listData.value = result.items || result || []
    }
  } catch (e: any) {
    ElMessage.warning('加载预览数据失败: ' + (e.message || ''))
  } finally {
    loading.value = false
  }
}

function openInNewTab() {
  const baseUrl = window.location.origin
  const lang = previewLang.value
  const page = previewPage.value === 'home' ? '' : previewPage.value
  window.open(`${baseUrl}/portal/${lang}/${page}`, '_blank')
}

// 社交链接点击跟踪
function handleSocialClick(s: { key: string; icon: string; url: string }) {
  if (s.url && s.url !== '#') {
    publicApi.trackClick(s.icon, s.url).catch(() => {})
    window.open(s.url, '_blank', 'noopener,noreferrer')
  }
}

onMounted(async () => {
  await loadLanguages()
  await loadTheme()
  await refreshPreview()
})
</script>

<style scoped>
.portal-preview { min-height: calc(100vh - 120px); }

/* Toolbar */
.preview-toolbar {
  display: flex; justify-content: space-between; align-items: center;
  padding: 12px 20px; background: #fff; border-bottom: 1px solid #e5e7eb;
  position: sticky; top: 0; z-index: 100;
}
.toolbar-left { display: flex; align-items: center; gap: 12px; }
.toolbar-title { font-size: 16px; font-weight: 600; }
.toolbar-right { display: flex; align-items: center; gap: 8px; }
.preview-loading { padding: 40px; }

/* Portal Container */
.portal-container { max-width: 1280px; margin: 0 auto; padding: 0 32px; }

/* Hero Section */
.portal-hero {
  position: relative; min-height: 80vh; display: flex; align-items: center;
  overflow: hidden; padding: 80px 40px;
}
.portal-hero-bg {
  position: absolute; inset: 0;
  background: linear-gradient(135deg, #0D1B2A 0%, #1B2D44 100%);
}
.portal-hero-pattern {
  position: absolute; inset: 0;
  background-image: radial-gradient(circle at 25% 25%, rgba(212, 168, 83, 0.05) 0%, transparent 50%);
}
.portal-hero-slide {
  position: absolute; inset: 0; opacity: 0; transition: opacity 0.8s;
}
.portal-hero-slide-active { opacity: 1; }
.portal-hero-img { width: 100%; height: 100%; object-fit: cover; opacity: 0.5; }
.portal-hero-dots {
  position: absolute; bottom: 44px; left: 50%; transform: translateX(-50%);
  display: flex; gap: 8px; z-index: 6;
}
.portal-hero-dot {
  width: 10px; height: 10px; border-radius: 999px; border: none; padding: 0;
  background: rgba(255,255,255,0.35); cursor: pointer; transition: all 0.2s;
}
.portal-hero-dot-active { background: #D4A853; width: 26px; }
.portal-hero-content {
  position: relative; max-width: 720px; padding: 40px;
}
.portal-hero-badge {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 6px 16px; background: rgba(255,255,255,0.1);
  backdrop-filter: blur(8px); border-radius: 999px;
  font-size: 12px; color: #D4A853; text-transform: uppercase;
  letter-spacing: 0.1em; border: 1px solid rgba(212, 168, 83, 0.2);
  margin-bottom: 24px;
}
.portal-hero-badge-dot {
  width: 6px; height: 6px; background: #D4A853; border-radius: 50%;
}
.portal-hero-title {
  font-size: 56px; font-weight: 700; color: #fff;
  line-height: 1.1; margin-bottom: 20px; letter-spacing: -0.02em;
}
.portal-hero-subtitle {
  font-size: 20px; color: rgba(255,255,255,0.6);
  line-height: 1.6; margin-bottom: 32px;
}
.portal-hero-actions { display: flex; gap: 16px; }
.portal-scroll-indicator {
  position: absolute; bottom: 40px; right: 40px;
  display: flex; flex-direction: column; align-items: center; gap: 8px;
}
.portal-scroll-indicator span {
  font-size: 10px; color: rgba(255,255,255,0.3);
  letter-spacing: 0.2em; text-transform: uppercase;
}
.portal-scroll-line {
  width: 1px; height: 48px;
  background: linear-gradient(to bottom, rgba(212,168,83,0.5), transparent);
}

/* Buttons */
.portal-btn {
  display: inline-flex; align-items: center; gap: 12px;
  padding: 16px 32px; border-radius: 999px;
  font-size: 16px; font-weight: 600; cursor: pointer;
  transition: all 0.3s; border: none;
}
.portal-btn-primary {
  background: #D4A853; color: #fff;
  box-shadow: 0 8px 32px rgba(212, 168, 83, 0.2);
}
.portal-btn-primary:hover { background: #C49A3F; }
.portal-btn-lg { padding: 18px 40px; font-size: 18px; }
.portal-btn svg { transition: transform 0.3s; }
.portal-btn:hover svg { transform: translateX(4px); }

/* Trust Bar */
.portal-trust-bar {
  background: #fff; border-bottom: 1px solid #EAE5DD;
}
.portal-trust-items {
  display: flex; flex-wrap: wrap; justify-content: center;
  gap: 40px 60px; padding: 24px 0;
}
.portal-trust-item {
  display: flex; align-items: center; gap: 8px;
  font-size: 14px; color: #8B7D6B;
}
.portal-trust-item svg { color: #D4A853; flex-shrink: 0; }

/* Section */
.portal-section { padding: 80px 0; }
.portal-section-alt { background: #FAFAF8; }
.portal-section-header {
  display: flex; justify-content: space-between; align-items: flex-end;
  margin-bottom: 40px;
}
.portal-section-center { text-align: center; flex-direction: column; align-items: center; }
.portal-section-label {
  font-size: 12px; letter-spacing: 0.2em; color: #D4A853;
  text-transform: uppercase; font-weight: 600;
}
.portal-section-title {
  font-size: 36px; font-weight: 700; color: #0D1B2A;
  margin-top: 12px; line-height: 1.2;
}
.portal-section-desc { color: #9A8C7A; margin-top: 8px; }
.portal-view-all {
  color: #0D1B2A; font-size: 14px; font-weight: 500;
  cursor: pointer; transition: color 0.2s; white-space: nowrap;
}
.portal-view-all:hover { color: #D4A853; }

/* Empty State */
.portal-empty-state { text-align: center; padding: 60px 20px; color: #9A8C7A; }
.portal-empty-icon { font-size: 64px; margin-bottom: 16px; opacity: 0.5; }

/* Product Grid */
.portal-product-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 20px; }
.portal-product-card {
  background: #FAFAF8; border-radius: 16px; overflow: hidden;
  border: 1px solid #EAE5DD; transition: all 0.3s;
}
.portal-product-card:hover { border-color: rgba(212, 168, 83, 0.4); box-shadow: 0 8px 24px rgba(0,0,0,0.08); }
.portal-product-image {
  position: relative; aspect-ratio: 3/4;
  background: linear-gradient(135deg, #F5F0E8, #EAE5DD);
  display: flex; align-items: center; justify-content: center; overflow: hidden;
}
.portal-product-image img {
  width: 100%; height: 100%; object-fit: cover;
  transition: transform 0.5s;
}
.portal-product-card:hover .portal-product-image img { transform: scale(1.05); }
.portal-placeholder-img { font-size: 48px; opacity: 0.3; }
.portal-product-overlay {
  position: absolute; bottom: 16px; left: 16px; right: 16px;
  display: flex; justify-content: space-between; opacity: 0;
  transition: opacity 0.3s;
}
.portal-product-card:hover .portal-product-overlay { opacity: 1; }
.portal-product-type {
  padding: 4px 12px; background: rgba(255,255,255,0.95);
  backdrop-filter: blur(8px); border-radius: 999px;
  font-size: 12px; font-weight: 600; color: #0D1B2A;
}
.portal-product-inquiry {
  padding: 4px 12px; background: #D4A853; color: #fff;
  border-radius: 999px; font-size: 12px; font-weight: 600;
}
.portal-product-status {
  padding: 4px 12px; background: #10B981; color: #fff;
  border-radius: 999px; font-size: 12px; font-weight: 600;
}
.portal-product-info { padding: 20px; }
.portal-product-meta { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 8px; }
.portal-product-name { font-size: 16px; font-weight: 600; color: #0D1B2A; margin: 0; }
.portal-product-gender {
  font-size: 12px; background: #F5F0E8; padding: 2px 8px;
  border-radius: 999px; color: #8B7D6B; white-space: nowrap;
}
.portal-product-brief { font-size: 13px; color: #9A8C7A; line-height: 1.5; margin-bottom: 12px; }
.portal-product-tags { display: flex; gap: 8px; flex-wrap: wrap; }
.portal-tag {
  font-size: 12px; color: #8B7D6B; padding: 2px 8px;
  border-radius: 999px; background: #F5F0E8;
}
.portal-tag-success { background: #D1FAE5; color: #065F46; }

/* Strengths */
.portal-strength-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 20px; }
.portal-strength-card {
  padding: 24px; background: #fff; border-radius: 16px;
  border: 1px solid #EAE5DD; transition: all 0.3s;
}
.portal-strength-card:hover { border-color: rgba(212, 168, 83, 0.3); box-shadow: 0 4px 16px rgba(0,0,0,0.06); }
.portal-strength-icon {
  width: 48px; height: 48px; background: #0D1B2A; border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  margin-bottom: 16px; color: #fff; transition: background 0.3s;
}
.portal-strength-card:hover .portal-strength-icon { background: #D4A853; }
.portal-strength-title { font-size: 18px; font-weight: 700; color: #0D1B2A; margin-bottom: 8px; }
.portal-strength-desc { font-size: 14px; color: #9A8C7A; line-height: 1.6; }

/* Categories */
.portal-category-list { display: flex; flex-wrap: wrap; justify-content: center; gap: 12px; }
.portal-category-tag {
  padding: 12px 24px; background: #FAFAF8; border: 1px solid #EAE5DD;
  border-radius: 999px; font-size: 14px; font-weight: 500;
  color: #4A4A4A; cursor: pointer; transition: all 0.3s;
}
.portal-category-tag:hover {
  background: #0D1B2A; color: #fff; border-color: #0D1B2A;
}

/* Blog */
.portal-blog-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 24px; }
.portal-blog-card {
  background: #fff; border-radius: 16px; overflow: hidden;
  border: 1px solid #EAE5DD; transition: all 0.3s;
}
.portal-blog-card:hover { border-color: rgba(212, 168, 83, 0.3); box-shadow: 0 4px 16px rgba(0,0,0,0.06); }
.portal-blog-image {
  aspect-ratio: 16/9; background: linear-gradient(135deg, #F5F0E8, #EAE5DD);
  display: flex; align-items: center; justify-content: center;
  overflow: hidden;
}
.portal-blog-image img {
  width: 100%; height: 100%; object-fit: cover;
  transition: transform 0.5s;
}
.portal-blog-card:hover .portal-blog-image img { transform: scale(1.05); }
.portal-blog-content { padding: 20px; }
.portal-blog-meta { display: flex; align-items: center; gap: 12px; margin-bottom: 8px; }
.portal-blog-category { font-size: 12px; color: #D4A853; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; }
.portal-blog-date { font-size: 12px; color: #9A8C7A; }
.portal-blog-title { font-size: 18px; font-weight: 700; color: #0D1B2A; margin-bottom: 8px; line-height: 1.4; }
.portal-blog-excerpt { font-size: 14px; color: #9A8C7A; line-height: 1.6; }

/* Certifications */
.portal-cert-grid { display: flex; flex-wrap: wrap; justify-content: center; gap: 16px; }
.portal-cert-card {
  display: flex; align-items: center; gap: 12px;
  padding: 16px 24px; background: #FAFAF8; border: 1px solid #EAE5DD;
  border-radius: 16px; transition: all 0.3s;
}
.portal-cert-card:hover { border-color: rgba(212, 168, 83, 0.3); box-shadow: 0 4px 12px rgba(0,0,0,0.06); }
.portal-cert-icon { font-size: 24px; }
.portal-cert-name { font-size: 14px; font-weight: 600; color: #0D1B2A; }
.portal-cert-desc { font-size: 12px; color: #9A8C7A; margin-top: 2px; }

/* CTA */
.portal-cta {
  position: relative; padding: 80px 0; text-align: center;
  overflow: hidden;
}
.portal-cta-glow {
  position: absolute; top: 0; left: 50%; transform: translateX(-50%);
  width: 50%; height: 1px;
  background: linear-gradient(to right, transparent, rgba(212,168,83,0.5), transparent);
}
.portal-cta-content { position: relative; max-width: 640px; margin: 0 auto; padding: 0 32px; }
.portal-cta-title { font-size: 36px; font-weight: 700; color: #fff; margin-bottom: 12px; }
.portal-cta-desc { color: rgba(255,255,255,0.6); margin-bottom: 32px; font-size: 16px; }

/* Footer */
.portal-footer { color: #fff; }
.portal-footer-newsletter { border-bottom: 1px solid #1B2D44; }
.portal-newsletter-content {
  display: flex; align-items: center; justify-content: space-between;
  gap: 24px; padding: 48px 0;
}
.portal-newsletter-content h3 { font-size: 20px; font-weight: 700; margin-bottom: 4px; }
.portal-newsletter-content p { font-size: 14px; color: #9A8C7A; }
.portal-newsletter-form { display: flex; gap: 8px; flex-shrink: 0; }
.portal-newsletter-form input {
  padding: 12px 16px; border-radius: 12px; border: 1px solid #2A3F59;
  background: #1B2D44; color: #fff; font-size: 14px; width: 260px;
  outline: none; transition: border-color 0.2s;
}
.portal-newsletter-form input:focus { border-color: #D4A853; }
.portal-newsletter-form input::placeholder { color: #6B5D4B; }
.portal-newsletter-form button {
  padding: 12px 24px; background: #D4A853; color: #fff;
  border: none; border-radius: 12px; font-weight: 600;
  font-size: 14px; cursor: pointer; transition: background 0.2s;
}
.portal-newsletter-form button:hover { background: #C49A3F; }
.portal-footer-grid { display: grid; grid-template-columns: 1.5fr 1fr 1fr 1fr; gap: 40px; padding: 64px 0; }
.portal-footer-col h4 {
  font-size: 12px; font-weight: 600; text-transform: uppercase;
  letter-spacing: 0.15em; color: #D4A853; margin-bottom: 20px;
}
.portal-footer-col ul { list-style: none; padding: 0; margin: 0; }
.portal-footer-col ul li { margin-bottom: 12px; }
.portal-footer-col ul li a { color: #9A8C7A; font-size: 14px; cursor: pointer; transition: color 0.2s; }
.portal-footer-col ul li a:hover { color: #fff; }
.portal-footer-logo { display: flex; align-items: center; gap: 10px; margin-bottom: 20px; }
.portal-footer-logo-icon {
  width: 36px; height: 36px; background: #D4A853; border-radius: 12px;
  display: flex; align-items: center; justify-content: center;
  font-size: 16px; font-weight: 700; color: #0D1B2A;
}
.portal-footer-logo-text { display: flex; flex-direction: column; }
.portal-footer-logo-text span:first-child { font-size: 14px; font-weight: 700; }
.portal-footer-logo-text span:last-child { font-size: 9px; text-transform: uppercase; letter-spacing: 0.3em; color: #9A8C7A; }
.portal-footer-about { font-size: 14px; color: #9A8C7A; line-height: 1.6; margin-bottom: 24px; }
.portal-footer-social { display: flex; gap: 8px; }
.portal-social-icon {
  width: 36px; height: 36px; background: #1B2D44; border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 14px; font-weight: 700; color: #9A8C7A; text-decoration: none;
  transition: all 0.2s;
}
.portal-social-icon:hover { background: #D4A853; color: #0D1B2A; }
.portal-contact-list li {
  display: flex; align-items: center; gap: 12px;
  color: #9A8C7A; font-size: 14px;
}
.portal-contact-list li svg { flex-shrink: 0; color: #D4A853; }
.portal-footer-bottom { border-top: 1px solid #1B2D44; }
.portal-footer-bottom > div {
  display: flex; justify-content: space-between; align-items: center;
  padding: 24px 0; font-size: 14px; color: #6B5D4B;
}
.portal-footer-links { display: flex; gap: 24px; }
.portal-footer-links a { color: #6B5D4B; cursor: pointer; transition: color 0.2s; }
.portal-footer-links a:hover { color: #fff; }

/* Product Detail Styles */
.pd-breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #9A8C7A;
  margin-bottom: 32px;
}
.pd-breadcrumb a {
  color: #D4A853;
  cursor: pointer;
  transition: color 0.2s;
}
.pd-breadcrumb a:hover {
  color: #C49A3F;
}
.pd-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 48px;
  margin-bottom: 48px;
}
.pd-gallery {
  position: sticky;
  top: 140px;
  align-self: start;
}
.pd-main-image {
  width: 100%;
  aspect-ratio: 3/4;
  background: linear-gradient(135deg, #F5F0E8, #EAE5DD);
  border-radius: 16px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}
.pd-main-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.pd-thumb-list {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
}
.pd-thumb-item {
  width: 64px;
  height: 64px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
  flex-shrink: 0;
}
.pd-thumb-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.pd-thumb-item.active {
  border-color: #D4A853;
}
.pd-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}
.pd-badge {
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}
.pd-badge-type {
  background: #0D1B2A;
  color: #fff;
}
.pd-badge-status {
  background: #D1FAE5;
  color: #065F46;
}
.pd-badge-featured {
  background: #FEF3C7;
  color: #92400E;
}
.pd-badge-new {
  background: #DBEAFE;
  color: #1E40AF;
}
.pd-title {
  font-size: 32px;
  font-weight: 700;
  color: #0D1B2A;
  margin: 0 0 12px;
  line-height: 1.2;
}
.pd-brief {
  font-size: 15px;
  color: #9A8C7A;
  line-height: 1.6;
  margin-bottom: 24px;
}
.pd-specs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 24px;
}
.pd-spec {
  padding: 12px;
  background: #FAFAF8;
  border-radius: 8px;
}
.pd-spec-label {
  display: block;
  font-size: 11px;
  color: #9A8C7A;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 4px;
}
.pd-spec-value {
  font-size: 14px;
  font-weight: 600;
  color: #0D1B2A;
}
.pd-actions {
  margin-top: 24px;
}
.pd-description {
  margin-bottom: 48px;
  padding: 32px;
  background: #fff;
  border-radius: 16px;
  border: 1px solid #EAE5DD;
}
.pd-description h3 {
  font-size: 20px;
  font-weight: 700;
  color: #0D1B2A;
  margin-bottom: 16px;
}
.pd-description-content {
  font-size: 14px;
  color: #4A4A4A;
  line-height: 1.8;
  white-space: pre-wrap;
}
.pd-section {
  margin-bottom: 48px;
}
.pd-section h3 {
  font-size: 20px;
  font-weight: 700;
  color: #0D1B2A;
  margin-bottom: 16px;
}
.pd-custom-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}
.pd-custom-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #EAE5DD;
}
.pd-custom-icon {
  font-size: 20px;
}
.pd-custom-name {
  font-size: 14px;
  font-weight: 600;
  color: #0D1B2A;
}
.pd-custom-note {
  font-size: 12px;
  color: #9A8C7A;
  margin-top: 2px;
}
.pd-specs-grid {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #EAE5DD;
  overflow: hidden;
}
.pd-specs-row {
  display: flex;
  padding: 12px 20px;
  border-bottom: 1px solid #EAE5DD;
}
.pd-specs-row:last-child {
  border-bottom: none;
}
.pd-specs-name {
  width: 160px;
  font-size: 13px;
  color: #9A8C7A;
  font-weight: 500;
  flex-shrink: 0;
}
.pd-specs-val {
  font-size: 13px;
  color: #0D1B2A;
}

@media (max-width: 768px) {
  .pd-layout {
    grid-template-columns: 1fr;
    gap: 24px;
  }
  .pd-specs {
    grid-template-columns: 1fr;
  }
  .pd-gallery {
    position: static;
  }
}

/* List View */
.portal-list-view { display: flex; flex-direction: column; gap: 12px; max-width: 800px; }
.portal-list-item { background: #fff; border-radius: 12px; padding: 20px; border: 1px solid #EAE5DD; }
.portal-list-item-content h4 { margin: 0 0 8px; font-size: 16px; color: #0D1B2A; }
.portal-list-item-desc { color: #9A8C7A; font-size: 13px; line-height: 1.5; margin-bottom: 8px; }
.portal-list-item-tags { display: flex; gap: 8px; }

/* Responsive */
@media (max-width: 768px) {
  .portal-hero-title { font-size: 36px; }
  .portal-hero { min-height: 60vh; padding: 40px 20px; }
  .portal-section-title { font-size: 28px; }
  .portal-footer-grid { grid-template-columns: 1fr; gap: 32px; }
  .portal-newsletter-content { flex-direction: column; text-align: center; }
  .portal-newsletter-form { width: 100%; }
  .portal-newsletter-form input { flex: 1; }
  .portal-container { padding: 0 16px; }
  .portal-section { padding: 48px 0; }
}
</style>