<template>
  <div class="module-config-editor">
    <el-alert v-if="!isKnown" type="info" :closable="false" class="mce-hint"
      title="该模块暂无可视化表单，请在下方 JSON 中配置（与既有模块兼容）。" />
    <el-alert v-else type="success" :closable="false" class="mce-hint" :title="hint" />

    <template v-if="isKnown">
      <!-- 短文本多语言字段 -->
      <div v-for="[k, label] in textFields" :key="k" class="mce-field">
        <div class="mce-label">{{ label }}（多语言，留空回退英文）</div>
        <div class="mce-langs">
          <el-input v-for="l in LANGS" :key="l.value" :model-value="trans(k)[l.value]"
            :placeholder="l.label" @update:model-value="setTrans(k, l.value, $event)" />
        </div>
      </div>

      <!-- 媒体字段（图片/视频） -->
      <div v-for="[k, label, mediaType] in mediaFields" :key="k" class="mce-field">
        <div class="mce-label">{{ label }}</div>
        <MediaPicker :model-value="cfgStr(k)" :media-type="mediaType" @update:model-value="set(k, $event)" />
      </div>

      <!-- 富文本多语言字段 -->
      <div v-for="[k, label] in richFields" :key="k" class="mce-field">
        <div class="mce-label-row">
          <span class="mce-label">{{ label }}（多语言）</span>
          <el-radio-group v-model="richLang" size="small">
            <el-radio-button v-for="l in LANGS" :key="l.value" :value="l.value">{{ l.label }}</el-radio-button>
          </el-radio-group>
        </div>
        <RichTextEditor :model-value="trans(k)[richLang]" :min-height="'160px'"
          @update:model-value="setTrans(k, richLang, $event)" />
      </div>

      <!-- 图片 -->
      <div v-if="type === 'about_hero' || type === 'about_story'" class="mce-field">
        <div class="mce-label">{{ type === 'about_hero' ? '背景图' : '配图' }}</div>
        <MediaPicker :model-value="cfgStr(type === 'about_hero' ? 'bg_image' : 'image')"
          @update:model-value="set(type === 'about_hero' ? 'bg_image' : 'image', $event)" />
      </div>

      <!-- 图文方向 -->
      <div v-if="type === 'about_story'" class="mce-field">
        <div class="mce-label">图片位置</div>
        <el-radio-group :model-value="cfgStr('image_side') || 'left'" @update:model-value="set('image_side', $event)">
          <el-radio-button value="left">图片在左</el-radio-button>
          <el-radio-button value="right">图片在右</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 按钮链接 -->
      <div v-if="type === 'about_cta'" class="mce-field">
        <div class="mce-label">按钮链接</div>
        <el-input :model-value="cfgStr('button_url')" placeholder="/contact" @update:model-value="set('button_url', $event)" />
      </div>

      <!-- 认证选择（about_certifications） -->
      <div v-if="type === 'about_certifications'" class="mce-field">
        <div class="mce-label">展示认证（留空则展示全部已发布认证）</div>
        <el-select v-model="refIds" multiple filterable clearable collapse-tags placeholder="选择认证" style="width: 100%">
          <el-option v-for="c in certOptions" :key="c.id" :label="c.name || c.code || c.id" :value="c.id" />
        </el-select>
      </div>

      <!-- 数据统计 items -->
      <div v-if="type === 'about_stats'" class="mce-field">
        <div class="mce-label-row">
          <span class="mce-label">统计数据</span>
          <el-button size="small" type="primary" plain @click="addItem('items', () => ({ value: '', label: emptyTrans() }))">添加一项</el-button>
        </div>
        <div v-for="(it, i) in items" :key="i" class="mce-array-item">
          <div class="mce-array-head"><span>#{{ i + 1 }}</span>
            <el-button size="small" text type="danger" @click="removeItem('items', i)">删除</el-button>
          </div>
          <el-input :model-value="it.value" placeholder="数值，如 15+ / 500+ / 50K" class="mce-gap" @update:model-value="patchItem('items', i, 'value', $event)" />
          <div class="mce-langs">
            <el-input v-for="l in LANGS" :key="l.value" :model-value="transIn(it.label, l.value)"
              :placeholder="'标签 ' + l.label" @update:model-value="patchItemTrans('items', i, 'label', l.value, $event)" />
          </div>
        </div>
      </div>

      <!-- 生产流程 steps -->
      <div v-if="type === 'about_process'" class="mce-field">
        <div class="mce-label-row">
          <span class="mce-label">流程步骤</span>
          <el-button size="small" type="primary" plain @click="addItem('steps', () => ({ icon: '✂️', title: emptyTrans(), desc: emptyTrans() }))">添加一步</el-button>
        </div>
        <div v-for="(st, i) in steps" :key="i" class="mce-array-item">
          <div class="mce-array-head"><span>步骤 {{ i + 1 }}</span>
            <el-button size="small" text type="danger" @click="removeItem('steps', i)">删除</el-button>
          </div>
          <el-input :model-value="st.icon" placeholder="图标（emoji）" class="mce-gap" @update:model-value="patchItem('steps', i, 'icon', $event)" />
          <div class="mce-langs">
            <el-input v-for="l in LANGS" :key="l.value" :model-value="transIn(st.title, l.value)"
              :placeholder="'标题 ' + l.label" @update:model-value="patchItemTrans('steps', i, 'title', l.value, $event)" />
          </div>
          <div class="mce-langs">
            <el-input v-for="l in LANGS" :key="l.value" :model-value="transIn(st.desc, l.value)"
              :placeholder="'描述 ' + l.label" @update:model-value="patchItemTrans('steps', i, 'desc', l.value, $event)" />
          </div>
        </div>
      </div>
    </template>

    <!-- 兜底 JSON 编辑（未知类型） -->
    <el-input v-if="!isKnown" type="textarea" :rows="4" :model-value="jsonText" placeholder='{"description":"..."}' @update:model-value="onJson" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import RichTextEditor from '@/components/RichTextEditor.vue'
import MediaPicker from '@/components/media/MediaPicker.vue'
import { certificationApi } from '@/api'

const LANGS = [
  { value: 'en', label: 'EN' },
  { value: 'zh', label: '中文' },
  { value: 'es', label: 'ES' },
  { value: 'fr', label: 'FR' },
]

// 各模块类型的"短文本多语言字段"（字段名 + 中文标签）
const TEXT_FIELDS: Record<string, Array<[string, string]>> = {
  about_hero: [['title', '标题'], ['subtitle', '副标题']],
  about_story: [['title', '标题']],
  about_certifications: [['title', '标题']],
  about_process: [['title', '标题']],
  about_cta: [['title', '标题'], ['description', '描述'], ['button_text', '按钮文字']],
}
// 富文本多语言字段
const RICH_FIELDS: Record<string, Array<[string, string]>> = {
  text: [['content', '正文']],
  about_story: [['body', '正文']],
}
// 媒体字段：[字段名, 标签, mediaType]，mediaType 传给 MediaPicker
const MEDIA_FIELDS: Record<string, Array<[string, string, string]>> = {
  image: [['url', '图片', 'image']],
  video: [['url', '视频', 'video']],
  banner: [['image', '图片', 'image']],
}
const KNOWN_TYPES = new Set<string>([...Object.keys(TEXT_FIELDS), ...Object.keys(RICH_FIELDS), ...Object.keys(MEDIA_FIELDS), 'about_stats'])

const HINTS: Record<string, string> = {
  text: '富文本正文（多语言）。',
  image: '图片：从媒体库选择图片 URL。',
  video: '视频：视频地址（支持 YouTube / 本地视频）。',
  banner: '横幅图：大图展示。',
  about_hero: '首屏横幅：副标题多语言 + 背景图。',
  about_story: '公司故事：标题/正文多语言 + 配图（可选左右）。',
  about_stats: '数据统计：数值 + 标签（多语言）。',
  about_certifications: '认证墙：标题多语言；认证数据前台自动读取「认证管理」。',
  about_process: '生产流程：标题多语言 + 步骤（图标/标题/描述）。',
  about_cta: '行动号召：标题/描述/按钮文字多语言 + 链接。',
}

const props = withDefaults(defineProps<{ type?: string; modelValue?: Record<string, any> }>(), {
  type: 'text',
  modelValue: () => ({}),
})
const emit = defineEmits<{ (e: 'update:modelValue', v: Record<string, any>): void }>()

const config = computed(() => props.modelValue || {})
const isKnown = computed(() => KNOWN_TYPES.has(props.type || ''))
const textFields = computed(() => TEXT_FIELDS[props.type || ''] || [])
const richFields = computed(() => RICH_FIELDS[props.type || ''] || [])
const mediaFields = computed(() => MEDIA_FIELDS[props.type || ''] || [])
const hint = computed(() => HINTS[props.type || ''] || '')
const richLang = ref('en')

// 认证墙：拉取已发布认证供运营勾选引用（config.ref_ids）
const certOptions = ref<any[]>([])
const refIds = computed<string[]>({
  get: () => (Array.isArray(config.value?.ref_ids) ? config.value.ref_ids : []),
  set: (v: string[]) => set('ref_ids', v),
})
async function loadCerts() {
  try {
    const res: any = await certificationApi.list({ page: 1, pageSize: 200 })
    certOptions.value = res?.items || res || []
  } catch { certOptions.value = [] }
}
watch(() => props.type, (t) => { if (t === 'about_certifications') loadCerts() }, { immediate: true })

function emptyTrans(): Record<string, string> { return { en: '', zh: '', es: '', fr: '' } }
function set(key: string, val: any) { emit('update:modelValue', { ...config.value, [key]: val }) }
function cfgStr(key: string): string { return (config.value?.[key] || '') as string }
function trans(key: string): Record<string, string> {
  const v = config.value?.[key]
  const o = v && typeof v === 'object' && !Array.isArray(v) ? v : {}
  return { en: o.en || '', zh: o.zh || '', es: o.es || '', fr: o.fr || '' }
}
function setTrans(key: string, lang: string, val: string) {
  const cur = trans(key); cur[lang] = val; set(key, { ...cur })
}
function transIn(o: any, lang: string): string { return (o && typeof o === 'object' ? o[lang] : '') || '' }

const items = computed(() => (Array.isArray(config.value?.items) ? config.value.items : []))
const steps = computed(() => (Array.isArray(config.value?.steps) ? config.value.steps : []))
function addItem(key: string, template: () => any) {
  set(key, [...(Array.isArray(config.value?.[key]) ? config.value[key] : []), template()])
}
function removeItem(key: string, idx: number) {
  const arr = [...(Array.isArray(config.value?.[key]) ? config.value[key] : [])]
  arr.splice(idx, 1); set(key, arr)
}
function patchItem(key: string, idx: number, field: string, val: any) {
  const arr = [...(Array.isArray(config.value?.[key]) ? config.value[key] : [])]
  arr[idx] = { ...arr[idx], [field]: val }; set(key, arr)
}
function patchItemTrans(key: string, idx: number, field: string, lang: string, val: string) {
  const arr = [...(Array.isArray(config.value?.[key]) ? config.value[key] : [])]
  const cur = arr[idx]?.[field] && typeof arr[idx][field] === 'object' ? { ...arr[idx][field] } : {}
  cur[lang] = val
  arr[idx] = { ...arr[idx], [field]: cur }; set(key, arr)
}

// 兜底 JSON 编辑（未知类型）
const jsonText = ref('')
function onJson(v: string) {
  try { emit('update:modelValue', v.trim() ? JSON.parse(v) : {}) } catch { /* 非法 JSON 忽略 */ }
}
watch(() => props.modelValue, (v) => {
  if (!isKnown.value) jsonText.value = v ? JSON.stringify(v, null, 2) : ''
}, { immediate: true })
</script>

<style scoped>
.module-config-editor { padding: 4px 0; }
.mce-hint { margin-bottom: 12px; }
.mce-field { margin-bottom: 14px; }
.mce-label { font-size: 12px; color: #909399; margin-bottom: 6px; }
.mce-label-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.mce-label-row .mce-label { margin-bottom: 0; }
.mce-langs { display: flex; gap: 8px; flex-wrap: wrap; }
.mce-langs .el-input { flex: 1; min-width: 120px; }
.mce-array-item {
  border: 1px solid #ebeef5; border-radius: 6px; padding: 10px 12px; margin-bottom: 10px; background: #fafbfc;
}
.mce-array-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; font-size: 13px; color: #606266; }
.mce-gap { margin-bottom: 8px; }
</style>
