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
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
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
// 位置以「距视口右 / 下边缘的像素间距」(anchor) 保存，而非绝对 left/top 像素。
// 这样在窗口 resize / 移动端浏览时地址栏收起等场景下，按钮能始终锚定在对应
// 角落而不会因视口高度抖动反复向上/向下漂移（之前以绝对 top 像素 + Math.min
// 夹逼会导致按钮越滚越靠边且无法回弹）。
const STORAGE_KEY = 'sw_cookie_trigger_pos'
const EDGE_MARGIN = 20
const btnSize = 56 // 按钮不可见 / 未测量时的兜底尺寸
const triggerRef = ref<HTMLElement | null>(null)
const isDragging = ref(false)
const dragStartPos = ref({ x: 0, y: 0 })
// 鼠标在按钮内部的按下偏移量，用于拖动时平滑跟随、不跳动
const dragOffsetPos = ref({ x: btnSize / 2, y: btnSize / 2 })
// 渲染用：按钮左上角的 left/top（viewport 像素），由 anchor 派生
const triggerPos = ref({ x: 0, y: 0 })
// 锚点：距离右 / 下边缘的间距（px），持久化存储的源头
const anchor = ref<{ right: number; bottom: number }>({ right: EDGE_MARGIN, bottom: EDGE_MARGIN })
const hasMoved = ref(false)
const mounted = ref(false)

// 默认锚点：右下角
const getDefaultAnchor = (): { right: number; bottom: number } => ({
  right: EDGE_MARGIN,
  bottom: EDGE_MARGIN,
})

// 从 localStorage 读取保存的锚点
const loadSavedAnchor = (): { right: number; bottom: number } => {
  if (import.meta.server) return getDefaultAnchor()
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const a = JSON.parse(saved)
      if (
        a &&
        typeof a.right === 'number' &&
        typeof a.bottom === 'number' &&
        a.right >= 0 &&
        a.right <= Math.max(0, window.innerWidth - btnSize) &&
        a.bottom >= 0 &&
        a.bottom <= Math.max(0, window.innerHeight - btnSize)
      ) {
        return { right: a.right, bottom: a.bottom }
      }
    }
  } catch { /* ignore */ }
  return getDefaultAnchor()
}

// 根据锚点 + 当前视口 + 按钮实际尺寸，重新计算 left/top
const applyAnchor = () => {
  if (!mounted.value) return
  const el = triggerRef.value
  const w = el ? el.offsetWidth : btnSize
  const h = el ? el.offsetHeight : btnSize
  triggerPos.value = {
    x: Math.max(0, window.innerWidth - w - anchor.value.right),
    y: Math.max(0, window.innerHeight - h - anchor.value.bottom),
  }
}

const triggerStyle = computed(() => {
  if (!mounted.value) return {}
  return {
    left: `${triggerPos.value.x}px`,
    top: `${triggerPos.value.y}px`,
    right: 'auto',
    bottom: 'auto',
  }
})

const startDrag = (e: MouseEvent | TouchEvent) => {
  isDragging.value = true
  hasMoved.value = false

  const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX
  const clientY = 'touches' in e ? e.touches[0].clientY : e.clientY
  dragStartPos.value = { x: clientX, y: clientY }

  // 记录鼠标相对于按钮左上角的偏移，拖动时手位不跳动
  const el = triggerRef.value
  if (el) {
    const rect = el.getBoundingClientRect()
    dragOffsetPos.value = { x: clientX - rect.left, y: clientY - rect.top }
  }

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

  // 限制在视口范围内（使用按钮实际尺寸）
  const el = triggerRef.value
  const w = el ? el.offsetWidth : btnSize
  const h = el ? el.offsetHeight : btnSize

  let newX = clientX - dragOffsetPos.value.x
  let newY = clientY - dragOffsetPos.value.y
  newX = Math.max(0, Math.min(window.innerWidth - w, newX))
  newY = Math.max(0, Math.min(window.innerHeight - h, newY))

  triggerPos.value = { x: newX, y: newY }

  // 同步锚点，保证 resize 后按钮停留在用户拖放的位置而非反弹
  anchor.value = {
    right: Math.max(0, window.innerWidth - w - newX),
    bottom: Math.max(0, window.innerHeight - h - newY),
  }

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

  // 落地时重新校准锚点（基于最终位置 + 实际尺寸），并持久化
  const el = triggerRef.value
  const w = el ? el.offsetWidth : btnSize
  const h = el ? el.offsetHeight : btnSize
  anchor.value = {
    right: Math.max(0, window.innerWidth - w - triggerPos.value.x),
    bottom: Math.max(0, window.innerHeight - h - triggerPos.value.y),
  }
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(anchor.value))
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

// 按钮从横幅隐藏→显示（banner 收起）时，根据锚点校准一次位置
watch(
  () => !showBanner.value,
  (visible) => {
    if (visible && mounted.value) {
      nextTick(() => applyAnchor())
    }
  },
)

onMounted(async () => {
  await initConsent()
  analyticsChecked.value = hasConsent(ConsentCategory.ANALYTICS)
  marketingChecked.value = hasConsent(ConsentCategory.MARKETING)
  // 初始化锚点与位置
  anchor.value = loadSavedAnchor()
  mounted.value = true
  applyAnchor()
  // 窗口大小变化时，根据锚点重新计算位置，保持按钮固定在视口角落，不漂移
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
  if (!isDragging.value) applyAnchor()
}
</script>

<style scoped>
/* 隐私偏好浮动入口（横幅隐藏时显示，可拖动） */
.preferences-trigger {
    position: fixed;
  /* SSR/首次渲染兜底：固定在视口右下角；mounted 后内联 left/top 接管 */
  right: 20px;
  bottom: 20px;
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
