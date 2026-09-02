<template>
  <div>
    <!-- 隐私偏好浮动入口：用户已做出同意决定且横幅隐藏时显示，可随时修改/撤回 -->
    <Teleport to="body">
      <button
        v-if="!showBanner"
        class="preferences-trigger"
        @click="openPreferences"
        :aria-label="$t('consent.preferences', 'Cookie Preferences')"
      >
        🍪 {{ $t('consent.preferences', 'Cookie Preferences') }}
      </button>
    </Teleport>

    <Teleport to="body">
      <Transition name="consent-slide">
        <div
          v-if="showBanner"
          class="cookie-consent-banner"
          role="dialog"
          aria-label="Cookie Consent"
        >
          <div class="consent-container">
            <div class="consent-content">
              <div class="consent-header">
                <h3 class="consent-title">{{ $t('consent.title', 'Cookie Consent') }}</h3>
                <span class="consent-badge">{{ regionalConfig.label }}</span>
              </div>
              <p class="consent-description">
                {{ $t('consent.description', regionalConfig.description) }}
              </p>

              <div v-if="showDetails" class="consent-details">
                <div class="cookie-category">
                  <label class="category-label">
                    <input type="checkbox" :checked="true" disabled class="category-checkbox" />
                    <span class="category-name">{{ $t('consent.necessary', 'Necessary Cookies') }}</span>
                    <span class="category-always">{{ $t('consent.always_active', 'Always Active') }}</span>
                  </label>
                  <p class="category-desc">{{ $t('consent.necessary_desc', 'Essential for website functionality.') }}</p>
                </div>

                <div class="cookie-category">
                  <label class="category-label">
                    <input type="checkbox" :checked="analyticsChecked" @change="toggleAnalytics" class="category-checkbox" />
                    <span class="category-name">{{ $t('consent.analytics', 'Analytics Cookies') }}</span>
                  </label>
                  <p class="category-desc">{{ $t('consent.analytics_desc', 'Help us understand visitor interactions.') }}</p>
                </div>

                <div class="cookie-category">
                  <label class="category-label">
                    <input type="checkbox" :checked="marketingChecked" @change="toggleMarketing" class="category-checkbox" />
                    <span class="category-name">{{ $t('consent.marketing', 'Marketing Cookies') }}</span>
                  </label>
                  <p class="category-desc">{{ $t('consent.marketing_desc', 'Used for advertising tracking.') }}</p>
                </div>
              </div>

              <button class="consent-toggle-details" @click="showDetails = !showDetails">
                {{ showDetails ? $t('consent.hide_details', 'Hide Details') : $t('consent.show_details', 'Show Details') }}
              </button>
            </div>

            <div class="consent-actions">
              <button class="consent-btn consent-btn-reject" @click="rejectAll">
                {{ $t('consent.reject', 'Reject All') }}
              </button>
              <button class="consent-btn consent-btn-customize" @click="saveCustomSelection">
                {{ $t('consent.save_selection', 'Save Selection') }}
              </button>
              <button class="consent-btn consent-btn-accept" @click="acceptAll">
                {{ $t('consent.accept', 'Accept All') }}
              </button>
            </div>

            <div class="consent-footer">
              <NuxtLink :to="localePath('/privacy-policy')" class="consent-link">
                {{ $t('consent.privacy_policy', 'Privacy Policy') }}
              </NuxtLink>
              <span class="consent-footer-sep">·</span>
              <button class="consent-withdraw" @click="withdrawConsent">
                {{ $t('consent.withdraw', 'Withdraw Consent') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useConsent, ConsentCategory } from '~/composables/useConsent'

const localePath = useLocalePath()

const {
  showBanner,
  regionalConfig,
  acceptAll,
  rejectAll,
  acceptCategories,
  hasConsent,
  initConsent,
  withdrawConsent,
} = useConsent()

const showDetails = ref(false)
const analyticsChecked = ref(false)
const marketingChecked = ref(false)

const toggleAnalytics = () => {
  analyticsChecked.value = !analyticsChecked.value
}

const toggleMarketing = () => {
  marketingChecked.value = !marketingChecked.value
}

const saveCustomSelection = () => {
  const categories: ConsentCategory[] = []
  if (analyticsChecked.value) categories.push(ConsentCategory.ANALYTICS)
  if (marketingChecked.value) categories.push(ConsentCategory.MARKETING)
  acceptCategories(categories)
}

/**
 * 打开隐私偏好设置（用户已同意后通过浮动按钮修改/撤回同意）
 */
const openPreferences = () => {
  analyticsChecked.value = hasConsent(ConsentCategory.ANALYTICS)
  marketingChecked.value = hasConsent(ConsentCategory.MARKETING)
  showBanner.value = true
}

onMounted(async () => {
  await initConsent()
  analyticsChecked.value = hasConsent(ConsentCategory.ANALYTICS)
  marketingChecked.value = hasConsent(ConsentCategory.MARKETING)
})
</script>

<style scoped>
/* 隐私偏好浮动入口（横幅隐藏时显示，左下角） */
.preferences-trigger {
  position: fixed;
  bottom: 20px;
  left: 20px;
  z-index: 9990;
  background: #fff;
  color: #1a73e8;
  border: 1px solid #dadce0;
  border-radius: 9999px;
  padding: 8px 16px;
  font-size: 13px;
  cursor: pointer;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.12);
  transition: all 0.2s;
}

.preferences-trigger:hover {
  background: #e8f0fe;
  border-color: #1a73e8;
}

.cookie-consent-banner {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 9999;
  background: #fff;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.15);
  border-top: 3px solid #1a73e8;
}

.consent-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px 24px;
}

.consent-content {
  margin-bottom: 16px;
}

.consent-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.consent-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.consent-badge {
  background: #e8f0fe;
  color: #1a73e8;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.consent-description {
  font-size: 14px;
  color: #5f6368;
  line-height: 1.5;
  margin: 0 0 12px;
}

.consent-details {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 12px;
}

.cookie-category {
  margin-bottom: 12px;
}

.cookie-category:last-child {
  margin-bottom: 0;
}

.category-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.category-checkbox {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.category-checkbox:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.category-name {
  font-weight: 500;
  color: #1a1a1a;
  font-size: 14px;
}

.category-always {
  background: #e8eaed;
  color: #5f6368;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  margin-left: auto;
}

.category-desc {
  font-size: 12px;
  color: #5f6368;
  margin: 4px 0 0 28px;
  line-height: 1.4;
}

.consent-toggle-details {
  background: none;
  border: none;
  color: #1a73e8;
  font-size: 13px;
  cursor: pointer;
  padding: 0;
  text-decoration: underline;
}

.consent-toggle-details:hover {
  color: #1557b0;
}

.consent-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.consent-btn {
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.consent-btn-reject {
  background: #fff;
  border-color: #dadce0;
  color: #5f6368;
}

.consent-btn-reject:hover {
  background: #f8f9fa;
  border-color: #5f6368;
}

.consent-btn-customize {
  background: #e8f0fe;
  color: #1a73e8;
  border-color: #e8f0fe;
}

.consent-btn-customize:hover {
  background: #d2e3fc;
}

.consent-btn-accept {
  background: #1a73e8;
  color: #fff;
}

.consent-btn-accept:hover {
  background: #1557b0;
}

.consent-footer {
  text-align: center;
}

.consent-link {
  color: #1a73e8;
  font-size: 13px;
  text-decoration: none;
}

.consent-link:hover {
  text-decoration: underline;
}

.consent-footer-sep {
  color: #dadce0;
  margin: 0 8px;
}

.consent-withdraw {
  background: none;
  border: none;
  color: #5f6368;
  font-size: 13px;
  cursor: pointer;
  text-decoration: underline;
  padding: 0;
}

.consent-withdraw:hover {
  color: #d93025;
}

.consent-slide-enter-active,
.consent-slide-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.consent-slide-enter-from,
.consent-slide-leave-to {
  transform: translateY(100%);
  opacity: 0;
}

@media (max-width: 640px) {
  .consent-container {
    padding: 16px;
  }

  .consent-actions {
    flex-direction: column;
  }

  .consent-btn {
    width: 100%;
    text-align: center;
  }

  .consent-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
