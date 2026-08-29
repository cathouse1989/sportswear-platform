import type { useAuthStore } from '@/stores/auth'

type AuthStore = ReturnType<typeof useAuthStore>

/**
 * any-of 权限判定：命中数组中任一权限码即通过。
 * 路由守卫（meta.permission）与 v-permission 指令共用的唯一实现，
 * 通配符 '*'（超级管理员）由 authStore.hasPermission 内部处理。
 */
export function hasAnyPermission(store: AuthStore, required: string | string[]): boolean {
  const codes = Array.isArray(required) ? required : [required]
  return codes.some((code) => store.hasPermission(code))
}
