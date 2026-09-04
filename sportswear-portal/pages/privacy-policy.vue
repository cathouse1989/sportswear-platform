<template>
  <div class="privacy-page">
    <div class="container">
      <header class="page-header">
        <h1>{{ $t('privacy.title', 'Privacy Policy') }}</h1>
        <p class="last-updated">{{ $t('privacy.last_updated', 'Last Updated') }}: {{ lastUpdated }}</p>
      </header>

      <div class="content">
        <section class="section">
          <h2>{{ $t('privacy.sections.introduction', '1. Introduction') }}</h2>
          <p>{{ $t('privacy.intro_text', 'We respect your privacy and are committed to protecting your personal data.') }}</p>
        </section>

        <section class="section">
          <h2>{{ $t('privacy.sections.data_collected', '2. Information We Collect') }}</h2>
          <p>{{ $t('privacy.data_collected_intro', 'We may collect different kinds of personal data:') }}</p>
          <ul>
            <li><strong>{{ $t('privacy.identity_data', 'Identity Data') }}:</strong> {{ $t('privacy.identity_data_desc', 'Name, title, position.') }}</li>
            <li><strong>{{ $t('privacy.contact_data', 'Contact Data') }}:</strong> {{ $t('privacy.contact_data_desc', 'Email, phone, company, address.') }}</li>
            <li><strong>{{ $t('privacy.technical_data', 'Technical Data') }}:</strong> {{ $t('privacy.technical_data_desc', 'IP, browser, OS, usage patterns.') }}</li>
          </ul>
        </section>

        <section class="section">
          <h2>{{ $t('privacy.sections.cookies', '3. Cookies') }}</h2>
          <p>{{ $t('privacy.cookies_intro', 'We use cookies to distinguish you from other users. You can control cookies through our consent banner.') }}</p>
        </section>

        <section class="section">
          <h2>{{ $t('privacy.sections.data_use', '4. How We Use Your Data') }}</h2>
          <ul>
            <li>{{ $t('privacy.use_1', 'Respond to inquiries and provide support') }}</li>
            <li>{{ $t('privacy.use_2', 'Improve our website and services') }}</li>
            <li>{{ $t('privacy.use_3', 'Send marketing communications (with consent)') }}</li>
            <li>{{ $t('privacy.use_4', 'Comply with legal obligations') }}</li>
          </ul>
        </section>

        <section class="section">
          <h2>{{ $t('privacy.sections.data_retention', '5. Data Retention') }}</h2>
          <p>{{ $t('privacy.retention_text', 'We retain your data only as long as necessary for the purposes collected or as required by law.') }}</p>
        </section>

        <section class="section">
          <h2>{{ $t('privacy.sections.your_rights', '6. Your Rights') }}</h2>
          <ul>
            <li>{{ $t('privacy.right_access', 'Right of Access') }}</li>
            <li>{{ $t('privacy.right_rectification', 'Right to Rectification') }}</li>
            <li>{{ $t('privacy.right_erasure', 'Right to Erasure') }}</li>
            <li>{{ $t('privacy.right_restrict', 'Right to Restrict Processing') }}</li>
            <li>{{ $t('privacy.right_portability', 'Right to Data Portability') }}</li>
            <li>{{ $t('privacy.right_withdraw', 'Right to Withdraw Consent') }}</li>
          </ul>
        </section>

        <section class="section">
          <h2>{{ $t('privacy.sections.contact', '7. Contact Us') }}</h2>
          <ul class="contact-list">
            <li><strong>{{ $t('privacy.contact_email', 'Email') }}:</strong> <a :href="`mailto:${contactEmail}`">{{ contactEmail }}</a></li>
            <li><strong>{{ $t('privacy.contact_phone', 'Phone') }}:</strong> {{ contactPhone }}</li>
          </ul>
        </section>

        <section class="section regional-notice" v-if="regionalNotice">
          <h2>{{ regionalNotice.title }}</h2>
          <p>{{ regionalNotice.content }}</p>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const { t, locale } = useI18n()
const config = useRuntimeConfig()

const lastUpdated = computed(() => {
  return new Date().toLocaleDateString(locale.value, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
})

const contactEmail = computed(() => config.public.contactEmail || 'privacy@example.com')
const contactPhone = computed(() => config.public.contactPhone || '+1-XXX-XXX-XXXX')

const regionalNotice = computed(() => {
  const notices: Record<string, { title: string; content: string }> = {
    zh: {
      title: '8. 中国个人信息保护法特别说明',
      content: '根据《个人信息保护法》，我们处理您的个人信息前会征得您的明确同意。您有权查阅、复制、更正、补充、删除您的个人信息。',
    },
    es: {
      title: '8. Aviso RGPD (España/UE)',
      content: 'De acuerdo con el RGPD, necesitamos su consentimiento explícito antes de tratar sus datos personales. Puede presentar una reclamación ante la AEPD.',
    },
    fr: {
      title: '8. Avis RGPD (France/UE)',
      content: 'Conformément au RGPD, nous avons besoin de votre consentement explicite. Vous pouvez déposer une plainte auprès de la CNIL.',
    },
  }
  return notices[locale.value] || null
})

useSeoHead({
  title: t('privacy.title', 'Privacy Policy'),
  description: t('privacy.meta_description', 'Our privacy policy explains how we collect, use, and protect your personal data. Learn about your rights under GDPR, CCPA, and PIPL.'),
})
</script>

<style scoped>
.privacy-page {
  padding: 40px 0 80px;
  background: #f8f9fa;
  min-height: 100vh;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
  padding-bottom: 24px;
  border-bottom: 2px solid #e8eaed;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 12px;
}

.last-updated {
  color: #5f6368;
  font-size: 14px;
  margin: 0;
}

.content {
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.section {
  margin-bottom: 32px;
}

.section:last-child {
  margin-bottom: 0;
}

.section h2 {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e8eaed;
}

.section p {
  font-size: 15px;
  line-height: 1.7;
  color: #3c4043;
  margin: 0 0 12px;
}

.section ul {
  margin: 0 0 16px;
  padding-left: 24px;
}

.section li {
  font-size: 15px;
  line-height: 1.7;
  color: #3c4043;
  margin-bottom: 8px;
}

.section li strong {
  color: #1a1a1a;
}

.contact-list {
  list-style: none;
  padding: 0;
}

.contact-list li {
  margin-bottom: 8px;
}

.section a {
  color: #1a73e8;
  text-decoration: none;
}

.section a:hover {
  text-decoration: underline;
}

.regional-notice {
  background: #e8f0fe;
  border-radius: 8px;
  padding: 20px;
  margin-top: 32px;
}

.regional-notice h2 {
  color: #1a73e8;
  border-bottom-color: #d2e3fc;
}

@media (max-width: 640px) {
  .privacy-page {
    padding: 20px 0 60px;
  }

  .container {
    padding: 0 16px;
  }

  .content {
    padding: 24px 16px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .section h2 {
    font-size: 18px;
  }
}
</style>
