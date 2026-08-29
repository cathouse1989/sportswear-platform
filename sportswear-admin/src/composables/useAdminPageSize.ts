/**
 * 从主题配置中读取管理后台列表默认每页条数。
 * 用法：const pageSize = useAdminPageSize()  // Ref<number>
 */
import { ref, type Ref } from 'vue'
import { themeApi } from '@/api'

let cached: number | null = null
let pending: Promise<number> | null = null

export function useAdminPageSize(): Ref<number> {
  const size = ref(20)

  async function load() {
    if (cached !== null) { size.value = cached; return }
    if (!pending) {
      pending = themeApi.adminList()
        .then((items: any[]) => {
          const found = items?.find((i: any) => i.key === 'admin_page_size')
          cached = Number(found?.value) || 20
          return cached
        })
        .catch(() => { cached = 20; return cached })
        .finally(() => { pending = null })
    }
    size.value = await pending
  }

  load()
  return size
}