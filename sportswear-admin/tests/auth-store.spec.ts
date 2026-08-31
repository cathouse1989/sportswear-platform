import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

vi.mock('@/api', () => ({
  authApi: { login: vi.fn(), profile: vi.fn() },
}))

import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'
import type { User } from '@/types'

function makeProfile(roles: any[]): User {
  return {
    id: 'u1',
    email: 'tester@sportswear.com',
    name: '测试',
    status: 'active',
    created_at: '2026-01-01T00:00:00Z',
    roles,
  }
}

describe('authStore.fetchProfile（账号-角色-权限闭环）', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.mocked(authApi.profile).mockClear()
  })

  it('从 profile.roles[].permissions 收集权限码（后端 profile 需携带角色权限）', async () => {
    localStorage.setItem('admin_token', 'token')
    const store = useAuthStore()
    vi.mocked(authApi.profile).mockResolvedValue(
      makeProfile([
        {
          id: 'r1',
          name: '编辑人员',
          code: 'editor',
          is_active: true,
          permissions: [
            { id: 'p1', code: 'page:view', name: '页面查看', module: 'cms' },
            { id: 'p2', code: 'blog:manage', name: '博客管理', module: 'cms' },
          ],
        },
      ])
    )

    await store.fetchProfile()

    expect(store.permissions).toContain('page:view')
    expect(store.permissions).toContain('blog:manage')
    expect(store.permissions).not.toContain('*')
    expect(store.hasPermission('page:view')).toBe(true)
    expect(store.hasPermission('role:manage')).toBe(false)
  })

  it('禁用状态的角色不贡献权限码', async () => {
    localStorage.setItem('admin_token', 'token')
    const store = useAuthStore()
    vi.mocked(authApi.profile).mockResolvedValue(
      makeProfile([
        {
          id: 'r-disabled',
          name: '已禁用角色',
          code: 'legacy',
          is_active: false,
          permissions: [{ id: 'p1', code: 'role:manage', name: '角色管理', module: 'user' }],
        },
      ])
    )

    await store.fetchProfile()

    expect(store.permissions).not.toContain('role:manage')
    expect(store.hasPermission('role:manage')).toBe(false)
  })

  it('super_admin 角色获得通配符 *（全量权限）', async () => {
    localStorage.setItem('admin_token', 'token')
    const store = useAuthStore()
    vi.mocked(authApi.profile).mockResolvedValue(
      makeProfile([
        {
          id: 'r-sa',
          name: '超级管理员',
          code: 'super_admin',
          is_active: true,
          permissions: [],
        },
      ])
    )

    await store.fetchProfile()

    expect(store.permissions).toContain('*')
    expect(store.hasPermission('anything:at:all')).toBe(true)
  })

  it('单角色 + 多角色权限去重收集', async () => {
    localStorage.setItem('admin_token', 'token')
    const store = useAuthStore()
    vi.mocked(authApi.profile).mockResolvedValue(
      makeProfile([
        {
          id: 'r1',
          name: '销售',
          code: 'sales',
          is_active: true,
          permissions: [{ id: 'p1', code: 'lead:view', name: '询盘查看', module: 'lead' }],
        },
        {
          id: 'r2',
          name: '销售经理',
          code: 'sales_manager',
          is_active: true,
          permissions: [
            { id: 'p1', code: 'lead:view', name: '询盘查看', module: 'lead' },
            { id: 'p2', code: 'dashboard:view', name: '数据统计', module: 'system' },
          ],
        },
      ])
    )

    await store.fetchProfile()

    expect(store.permissions.filter((p) => p === 'lead:view')).toHaveLength(1)
    expect(store.hasPermission('dashboard:view')).toBe(true)
  })
})