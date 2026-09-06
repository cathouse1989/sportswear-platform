<template>
  <div class="trans-editor">
    <div class="trans-editor__bar">
      <div class="trans-editor__langs">
        <button
          v-for="l in langs"
          :key="l.value"
          type="button"
          class="trans-editor__lang"
          :class="{ 'is-active': activeLang === l.value }"
          @click="activeLang = l.value"
        >
          <span class="trans-editor__lang-name">{{ l.label }}</span>
          <span class="trans-editor__coverage" :class="coverage(l.value).missing ? 'is-missing' : 'is-done'">
            {{ coverage(l.value).done }}/{{ fields.length }}
          </span>
        </button>
      </div>
      <el-button size="small" plain type="primary" @click="copyFromSource(false)">
        从 {{ sourceLabel }} 复制（仅填空缺）
      </el-button>
      <el-button size="small" plain type="success" :loading="aiTranslating" @click="aiTranslate">
        AI 翻译
      </el-button>
    </div>

    <el-alert
      v-if="showTip"
      type="info"
      :closable="false"
      show-icon
      title="未填写的字段将自动回退源语言（English）；覆盖度标记当前语言已填字段数。"
      class="trans-editor__tip"
    />

    <div class="trans-editor__fields">
      <div v-for="f in fields" :key="f.key" class="trans-editor__field">
        <div class="trans-editor__label">
          <span>{{ f.label }}</span>
          <el-tag v-if="!fieldFilled(activeLang, f.key)" size="small" type="danger" effect="plain">未填</el-tag>
        </div>
        <el-input
          v-if="f.type === 'textarea'"
          :model-value="getVal(activeLang, f.key)"
          :rows="f.rows || 3"
          :placeholder="f.placeholder || '留空则沿用源语言'"
          @update:model-value="onInput(activeLang, f.key, $event)"
        />
        <RichTextEditor
          v-else-if="f.type === 'richtext'"
          :model-value="getVal(activeLang, f.key)"
          :placeholder="f.placeholder || '留空则沿用源语言'"
          :min-height="f.minHeight || '120px'"
          @update:model-value="onInput(activeLang, f.key, $event)"
        />
        <el-input
          v-else
          :model-value="getVal(activeLang, f.key)"
          :placeholder="f.placeholder || '留空则沿用源语言'"
          @update:model-value="onInput(activeLang, f.key, $event)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import RichTextEditor from '@/components/RichTextEditor.vue'
import { aiTranslateApi } from '@/api'
import { DEFAULT_TRANS_LANGS } from '@/composables/useTransRecord'

interface TransField {
  key: string
  label: string
  type?: 'input' | 'textarea' | 'richtext'
  rows?: number
  placeholder?: string
  minHeight?: string
}
interface TransLang {
  value: string
  label: string
}

const props = withDefaults(
  defineProps<{
    fields: TransField[]
    modelValue: Record<string, Record<string, string>>
    source: Record<string, string>
    langs?: TransLang[]
    sourceLabel?: string
    showTip?: boolean
  }>(),
  {
    langs: () => DEFAULT_TRANS_LANGS,
    sourceLabel: 'English',
    showTip: true,
  },
)
const emit = defineEmits<{ (e: 'update:modelValue', v: Record<string, Record<string, string>>): void }>()

const activeLang = ref(props.langs[0]?.value || 'zh')

function clone(): Record<string, Record<string, string>> {
  const next: Record<string, Record<string, string>> = {}
  for (const l of props.langs) next[l.value] = { ...(props.modelValue?.[l.value] || {}) }
  return next
}
function getVal(lang: string, key: string): string {
  return (props.modelValue?.[lang]?.[key]) || ''
}
function onInput(lang: string, key: string, val: string) {
  const next = clone()
  if (!next[lang]) next[lang] = {}
  next[lang][key] = val
  emit('update:modelValue', next)
}
function isFilled(v?: string): boolean {
  return Boolean((v || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').trim())
}
function fieldFilled(lang: string, key: string): boolean {
  return isFilled(getVal(lang, key))
}
function coverage(lang: string) {
  let done = 0
  for (const f of props.fields) if (fieldFilled(lang, f.key)) done++
  return { done, missing: done < props.fields.length }
}
function copyFromSource(overwriteAll: boolean) {
  const next = clone()
  let copied = 0
  for (const l of props.langs) {
    if (!next[l.value]) next[l.value] = {}
    for (const f of props.fields) {
      const src = props.source?.[f.key]
      if (!isFilled(src)) continue
      if (overwriteAll || !isFilled(next[l.value][f.key])) {
        next[l.value][f.key] = src
        copied++
      }
    }
  }
  emit('update:modelValue', next)
  ElMessage.success(overwriteAll ? '已覆盖全部语言字段' : `已填充 ${copied} 个空缺字段`)
}

const aiTranslating = ref(false)
async function aiTranslate() {
  const srcFields = props.fields.filter((f) => isFilled(props.source[f.key]))
  const texts = srcFields.map((f) => props.source[f.key])
  if (!texts.length) { ElMessage.warning('源语言无内容可翻译'); return }
  aiTranslating.value = true
  try {
    const res = await aiTranslateApi.translate(texts, activeLang.value)
    const trans = (res as any).translations || []
    const next = clone()
    if (!next[activeLang.value]) next[activeLang.value] = {}
    srcFields.forEach((f, i) => { if (trans[i]) next[activeLang.value][f.key] = trans[i] })
    emit('update:modelValue', next)
    ElMessage.success('AI 翻译完成')
  } catch { ElMessage.warning('AI 翻译失败，可改用「从 English 复制」') } finally { aiTranslating.value = false }
}
</script>

<style scoped>
.trans-editor__bar { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 12px; flex-wrap: wrap; }
.trans-editor__langs { display: flex; gap: 8px; flex-wrap: wrap; }
.trans-editor__lang {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 5px 12px; border: 1px solid #dcdfe6; border-radius: 6px;
  background: #fff; color: #606266; font-size: 13px; cursor: pointer; transition: all .2s;
}
.trans-editor__lang:hover { border-color: #409eff; color: #409eff; }
.trans-editor__lang.is-active { border-color: #409eff; color: #409eff; background: #ecf5ff; }
.trans-editor__coverage { font-size: 11px; padding: 0 5px; border-radius: 9px; line-height: 16px; }
.trans-editor__coverage.is-done { color: #67c23a; background: #f0f9eb; }
.trans-editor__coverage.is-missing { color: #e6a23c; background: #fdf6ec; }
.trans-editor__tip { margin-bottom: 12px; }
.trans-editor__field { margin-bottom: 14px; }
.trans-editor__label { display: flex; align-items: center; gap: 8px; font-size: 12px; color: #909399; margin-bottom: 4px; }
</style>