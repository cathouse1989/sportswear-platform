// Cookie 同意管理 Composable
// 实现 GDPR/CCPA/PIPL 等多地区隐私法规的同意管理

import { ref, computed } from 'vue'
import { useApi } from '~/composables/useApi'

export enum ConsentCategory {
  NECESSARY = 'necessary',
  ANALYTICS = 'analytics',
  MARKETING = 'marketing',
}

export interface ConsentState {
  version: string
  timestamp: number
  categories: ConsentCategory[]
  country: string
}

const CONSENT_STORAGE_KEY = 'sw_cookie_consent'
const CONSENT_VERSION = '1.0'

const GDPR_COUNTRIES = [
  'AT', 'BE', 'BG', 'HR', 'CY', 'CZ', 'DK', 'EE', 'FI', 'FR',
  'DE', 'GR', 'HU', 'IE', 'IT', 'LV', 'LT', 'LU', 'MT', 'NL',
  'PL', 'PT', 'RO', 'SK', 'SI', 'ES', 'SE',
  'GB', 'CH', 'NO', 'IS', 'LI',
]

const PIPL_COUNTRIES = ['CN']
const CCPA_COUNTRIES = ['US']
const LGPD_COUNTRIES = ['BR']

export type RegionalMode = 'opt-in' | 'opt-out' | 'explicit'

export interface RegionalConfig {
  mode: RegionalMode
  label: string
  description: string
  defaultCategories: ConsentCategory[]
  requiresBanner: boolean
}

export function getRegionalConfig(countryISO2: string): RegionalConfig {
  const code = (countryISO2 || '').toUpperCase()

  if (GDPR_COUNTRIES.includes(code)) {
    return {
      mode: 'opt-in',
      label: 'GDPR',
      description: 'According to GDPR, we need your explicit consent before setting non-essential cookies.',
      defaultCategories: [ConsentCategory.NECESSARY],
      requiresBanner: true,
    }
  }

  if (PIPL_COUNTRIES.includes(code)) {
    return {
      mode: 'explicit',
      label: 'PIPL',
      description: '根据《个人信息保护法》，我们需要在处理您的个人信息前获得您的明确同意。',
      defaultCategories: [ConsentCategory.NECESSARY],
      requiresBanner: true,
    }
  }

  if (LGPD_COUNTRIES.includes(code)) {
    return {
      mode: 'opt-in',
      label: 'LGPD',
      description: 'De acordo com a LGPD, precisamos do seu consentimento antes de definir cookies não essenciais.',
      defaultCategories: [ConsentCategory.NECESSARY],
      requiresBanner: true,
    }
  }

  if (CCPA_COUNTRIES.includes(code)) {
    return {
      mode: 'opt-out',
      label: 'CCPA',
      description: 'We use cookies to enhance your experience. You can opt out of non-essential cookies below.',
      defaultCategories: [ConsentCategory.NECESSARY, ConsentCategory.ANALYTICS],
      requiresBanner: true,
    }
  }

  return {
    mode: 'opt-in',
    label: 'Privacy',
    description: 'We use cookies to improve your browsing experience and analyze site traffic.',
    defaultCategories: [ConsentCategory.NECESSARY],
    requiresBanner: true,
  }
}

export const useConsent = () => {
  const consentState = ref<ConsentState | null>(null)
  const showBanner = ref(false)
  const userCountry = ref<string>('')

  const regionalConfig = computed(() => getRegionalConfig(userCountry.value))

  const loadConsent = (): ConsentState | null => {
    if (import.meta.server) return null
    try {
      const raw = localStorage.getItem(CONSENT_STORAGE_KEY)
      if (!raw) return null
      const state = JSON.parse(raw) as ConsentState
      if (state.version !== CONSENT_VERSION) {
        localStorage.removeItem(CONSENT_STORAGE_KEY)
        return null
      }
      return state
    } catch {
      return null
    }
  }

  const saveConsent = (categories: ConsentCategory[]) => {
    if (import.meta.server) return
    const state: ConsentState = {
      version: CONSENT_VERSION,
      timestamp: Date.now(),
      categories,
      country: userCountry.value,
    }
    localStorage.setItem(CONSENT_STORAGE_KEY, JSON.stringify(state))
    consentState.value = state
    window.dispatchEvent(new CustomEvent('consent-changed', { detail: state }))
  }

  const acceptAll = () => {
    saveConsent([ConsentCategory.NECESSARY, ConsentCategory.ANALYTICS, ConsentCategory.MARKETING])
    showBanner.value = false
  }

  const rejectAll = () => {
    saveConsent([ConsentCategory.NECESSARY])
    showBanner.value = false
  }

  const acceptCategories = (categories: ConsentCategory[]) => {
    saveConsent([ConsentCategory.NECESSARY, ...categories])
    showBanner.value = false
  }

  const hasConsent = (category: ConsentCategory): boolean => {
    if (category === ConsentCategory.NECESSARY) return true
    return consentState.value?.categories.includes(category) ?? false
  }

  const getConsentedCategories = (): ConsentCategory[] => {
    return consentState.value?.categories || [ConsentCategory.NECESSARY]
  }

  const getConsentHeader = (): string => {
    return getConsentedCategories().join(',')
  }

  const withdrawConsent = () => {
    if (import.meta.server) return
    localStorage.removeItem(CONSENT_STORAGE_KEY)
    consentState.value = null
    showBanner.value = true
    window.dispatchEvent(new CustomEvent('consent-withdrawn'))
  }

  const initConsent = async (country?: string) => {
    if (country) {
      userCountry.value = country
    } else {
      try {
        const api = useApi()
        const geo = await api.getGeo()
        if (geo?.country_iso2) {
          userCountry.value = geo.country_iso2
        }
      } catch {}
    }
    const saved = loadConsent()
    if (saved) {
      consentState.value = saved
      showBanner.value = false
    } else {
      showBanner.value = regionalConfig.value.requiresBanner
    }
  }

  return {
    consentState,
    showBanner,
    userCountry,
    regionalConfig,
    loadConsent,
    saveConsent,
    acceptAll,
    rejectAll,
    acceptCategories,
    hasConsent,
    getConsentedCategories,
    getConsentHeader,
    withdrawConsent,
    initConsent,
  }
}


