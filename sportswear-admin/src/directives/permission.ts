import type { Directive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { hasAnyPermission } from '@/utils/permissions'

/**
 * 按钮级权限指令（P0-#3）：
 *   <el-button v-permission="'product:publish'">发布</el-button>
 *   <el-button v-permission="['user:delete', 'setting:manage']">删除</el-button>
 * 无权限时直接从 DOM 移除元素（any-of 语义，与路由守卫一致）。
 * 注意：权限矩阵见 platform/docs/权限矩阵.md，新增权限码需与后端同步。
 */
export const vPermission: Directive<HTMLElement, string | string[]> = {
  mounted(el, binding) {
    const store = useAuthStore()
    if (!hasAnyPermission(store, binding.value)) {
      el.parentNode?.removeChild(el)
    }
  },
}
