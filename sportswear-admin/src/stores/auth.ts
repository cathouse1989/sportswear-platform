import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('admin_token') || '')
  const user = ref<User | null>(null)
  const permissions = ref<string[]>([])

  async function login(email: string, password: string) {
    const result = await authApi.login(email, password)
    token.value = result.token
    localStorage.setItem('admin_token', result.token)
    await fetchProfile()
    return result
  }

  async function fetchProfile() {
    if (!token.value) return
    const profile = await authApi.profile()
    user.value = profile

    // 收集权限码（只收集启用状态的角色；禁用角色不产生任何权限，与后端 Auth 中间件一致）
    const perms: string[] = []
    let isSuperAdmin = false
    if (profile.roles) {
      for (const role of profile.roles) {
        if (role.is_active === false) continue
        // 超级管理员拥有所有权限
        if (role.code === 'super_admin') {
          isSuperAdmin = true
        }
        if (role.permissions) {
          for (const p of role.permissions) {
            if (!perms.includes(p.code)) perms.push(p.code)
          }
        }
      }
    }
    // 超级管理员：设置通配符权限
    if (isSuperAdmin) {
      perms.push('*')
    }
    permissions.value = perms
  }

  function hasPermission(code: string): boolean {
    // 通配符权限（超级管理员）拥有所有权限
    if (permissions.value.includes('*')) return true
    return permissions.value.includes(code)
  }

  async function logout() {
    // 先通知后端清除 HttpOnly 会话 Cookie（失败不阻塞本地清理）
    try {
      await authApi.logout()
    } catch {
      // 忽略登出接口异常（如 token 已过期、网络中断）
    }
    token.value = ''
    user.value = null
    permissions.value = []
    localStorage.removeItem('admin_token')
  }

  return { token, user, permissions, login, fetchProfile, hasPermission, logout }
})