import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

vi.mock('@/api', () => ({
  authApi: { login: vi.fn(), profile: vi.fn() },
}))

import { useAuthStore } from '@/stores/auth'
import { hasAnyPermission } from '@/utils/permissions'

describe('authStore.hasPermission', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('精确匹配权限码', () => {
    const store = useAuthStore()
    store.permissions = ['product:view', 'lead:view']
    expect(store.hasPermission('product:view')).toBe(true)
    expect(store.hasPermission('product:publish')).toBe(false)
  })

  it('超级管理员通配符 * 拥有任意权限', () => {
    const store = useAuthStore()
    store.permissions = ['*']
    expect(store.hasPermission('setting:manage')).toBe(true)
    expect(store.hasPermission('anything:at:all')).toBe(true)
  })
})

describe('hasAnyPermission（any-of 语义，路由守卫与 v-permission 共用）', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('字符串形式：命中即通过', () => {
    const store = useAuthStore()
    store.permissions = ['page:view']
    expect(hasAnyPermission(store, 'page:view')).toBe(true)
    expect(hasAnyPermission(store, 'page:publish')).toBe(false)
  })

  it('数组形式：任一命中即通过', () => {
    const store = useAuthStore()
    store.permissions = ['blog:manage']
    expect(hasAnyPermission(store, ['blog:manage', 'page:view'])).toBe(true)
    expect(hasAnyPermission(store, ['case:manage', 'page:view'])).toBe(false)
  })

  it('通配符在 any-of 中同样生效', () => {
    const store = useAuthStore()
    store.permissions = ['*']
    expect(hasAnyPermission(store, ['user:delete', 'role:manage'])).toBe(true)
  })
})
