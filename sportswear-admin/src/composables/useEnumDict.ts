import { ref } from 'vue'
import { enumApi } from '@/api'
import type { SysEnumType, SysEnumItem } from '@/types'

// 语义色 → Element Plus el-tag type 映射
type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'

const COLOR_TYPE_MAP: Record<string, TagType> = {
  blue: 'primary', green: 'success', purple: 'warning',
  orange: 'warning', magenta: 'danger', gold: 'warning',
  primary: 'primary', success: 'success', info: 'info',
  warning: 'warning', danger: 'danger',
}

// 模块级缓存：多个页面共享同一份枚举字典，避免重复请求
const types = ref<SysEnumType[]>([])
const itemsByCode = ref<Record<string, SysEnumItem[]>>({})
let loaded = false

/** 拉取全部枚举类型 + 条目（幂等，失败时静默降级，不影响主流程） */
async function ensureLoaded() {
  if (loaded) return
  try {
    types.value = await enumApi.types()
    const map: Record<string, SysEnumItem[]> = {}
    for (const t of types.value) {
      map[t.code] = await enumApi.items(t.id)
    }
    itemsByCode.value = map
    loaded = true
  } catch {
    // 枚举字典加载失败：下拉/标签回退空，不影响页面主流程
  }
}

/** 后台默认中文界面：优先 zh 翻译，回退 en（label），再回退 value */
function itemLabel(it: SysEnumItem): string {
  return it.translations?.zh || it.label || it.value
}

/**
 * 枚举字典读取（后台收编硬编码 label 用）：
 * 与「字典管理」/enums 同源，枚举值域/翻译/颜色在后台一处维护，全后台生效。
 */
export function useEnumDict() {
  return {
    ensureLoaded,
    /** 下拉选项（按 sort_order，过滤停用项） */
    options(code: string): Array<{ label: string; value: string }> {
      const items = itemsByCode.value[code] || []
      return items
        .filter((i) => i.is_active)
        .sort((a, b) => a.sort_order - b.sort_order)
        .map((i) => ({ label: itemLabel(i), value: i.value }))
    },
    /** 标签文案（zh 优先回退 en，找不到回退原值） */
    label(code: string, value?: string | null): string {
      if (!value) return ''
      const items = itemsByCode.value[code] || []
      const it = items.find((i) => i.value === value)
      return it ? itemLabel(it) : value
    },
    /** Element Plus el-tag type（映射语义色） */
    tagType(code: string, value?: string | null): TagType {
      if (!value) return 'info'
      const items = itemsByCode.value[code] || []
      const it = items.find((i) => i.value === value)
      const raw = it?.color || 'info'
      return COLOR_TYPE_MAP[raw] || 'info'
    },
  }
}
