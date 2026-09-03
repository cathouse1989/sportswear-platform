<template>
  <div>
    <!-- 隐私偏好浮动入口：用户已做出同意决定且横幅隐藏时显示，可随时修改/撤回 -->
    <Teleport to="body">
      <button
        v-if="!showBanner"
        ref="triggerRef"
        class="preferences-trigger"
        :class="{ 'dragging': isDragging }"
        :style="triggerStyle"
        @click="handleClick"
        @mousedown="startDrag"
        @touchstart.passive="startDrag"
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
import { ref, onMounted, onUnmounted, computed } from 'vue'
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

// ============ 浮动按钮拖动逻辑 ============
const STORAGE_KEY = 'sw_cookie_trigger_pos'
const triggerRef = ref<HTMLElement | null>(null)
const isDragging = ref(false)
const dragStartPos = ref({ x: 0, y: 0 })
const triggerStartPos = ref({ x: 0, y: 0 })
const triggerPos = ref({ x: 0, y: 0 })
const hasMoved = ref(false)

// 默认位置：右下角
const getDefaultPos = () => {
  if (import.meta.server) return { x: 20, y: 20 }
  const btnSize = 56
  const margin = 20
  return {
    x: window.innerWidth - btnSize - margin,
    y: window.innerHeight - btnSize - margin,
  }
}

// 从 localStorage 读取保存位置
const loadSavedPos = () => {
  if (import.meta.server) return getDefaultPos()
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const pos = JSON.parse(saved)
      // 确保位置在视口内
      if (pos.x >= 0 && pos.x <= window.innerWidth - 56 &&
          pos.y >= 0 && pos.y <= window.innerHeight - 56) {
        return pos
      }
    }
  } catch { /* ignore */ }
  return getDefaultPos()
}

const triggerStyle = computed(() => ({
  left: `${triggerPos.value.x}px`,
  top: `${triggerPos.value.y}px`,
}))

const startDrag = (e: MouseEvent | TouchEvent) => {
  isDragging.value = true
  hasMoved.value = false

  const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX
  const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY
  dragStartPos.value = { x: clientX, y: clientY }
  triggerStartPos.value = { ...triggerPos.value }

  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
  document.addEventListener('touchmove', onDrag, { passive: false })
  document.addEventListener('touchend', stopDrag)
}

const onDrag = (e: MouseEvent | TouchEvent) => {
  if (!isDragging.value) return

  const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX
  const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY

  const dx = clientX - dragStartPos.value.x
  const dy = clientY - dragStartPos.value.y

  // 移动超过 5px 才算拖动（避免误触点击）
  if (Math.abs(dx) > 5 || Math.abs(dy) > 5) {
    hasMoved.value = true
  }

  // 限制在视口范围内
  const btnSize = 56
  const newX = Math.max(0, Math.min(window.innerWidth - btnSize, triggerStartPos.value.x + dx))
  const newY = Math.max(0, Math.min(window.innerHeight - btnSize, triggerStartPos.value.y + dy))

  triggerPos.value = { x: newX, y: newY }

  // 触摸时阻止页面滚动
  if ('touches' in e && e.cancelable) {
    e.preventDefault()
  }
}

const stopDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
  document.removeEventListener('touchmove', onDrag)
  document.removeEventListener('touchend', stopDrag)

  // 保存位置到 localStorage
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(triggerPos.value))
  } catch { /* ignore */ }
}

const handleClick = () => {
  // 如果拖动过，不触发点击
  if (hasMoved.value) return
  openPreferences()
}

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
  // 初始化位置
  triggerPos.value = loadSavedPos()
  // 窗口大小变化时调整位置
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
  document.removeEventListener('touchmove', onDrag)
  document.removeEventListener('touchend', stopDrag)
})

const handleResize = () => {
  // 窗口变化时确保按钮在视口内
  const btnSize = 56
  triggerPos.value = {
    x: Math.min(triggerPos.value.x, window.innerWidth - btnSize),
    y: Math.min(triggerPos.value.y, window.innerHeight - btnSize),
  }
}
</script>

<style scoped>
/* 隐私偏好浮动入口（横幅隐藏时显示，可拖动） */
.preferences-trigger {
  position: fixed;
  /* 默认位置由内联样式 left/top 控制，bottom/right 作为 fallback */
  bottom: auto;
  right: auto;
  left: auto;
  top: auto;
  z-index: 9990;
  background: #fff;
  color: #1a73e8;
  border: 1px solid #dadce0;
  border-radius: 9999px;
  padding: 8px 16px;
  font-size: 13px;
  cursor: pointer;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.12);
  transition: box-shadow 0.2s, background 0.2s;
  user-select: none;
  -webkit-user-select: none;
  touch-action: none;
}

.preferences-trigger:hover {
  background: #e8f0fe;
  border-color: #1a73e8;
}

.preferences-trigger.dragging {
  cursor: grabbing;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  opacity: 0.9;
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
