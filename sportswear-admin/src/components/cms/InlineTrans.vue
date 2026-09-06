<template>
  <div class="inline-trans">
    <el-input
      v-for="l in langs"
      :key="l.value"
      :model-value="modelValue[l.value] || ''"
      :placeholder="prefix ? `${prefix} (${l.label})` : l.label"
      @update:model-value="onInput(l.value, $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { DEFAULT_TRANS_LANGS, type TransLang } from '@/composables/useTransRecord'

// 子资源（规格/视频/定制）的"单字段多语言"紧凑输入组：
// v-model 为 Record<lang, string>，统一渲染 zh/es/fr 三个输入框，
// 取代此前散落的 *_zh / *_es / *_fr 平铺字段。
const props = withDefaults(
  defineProps<{
    modelValue?: Record<string, string>
    langs?: TransLang[]
    prefix?: string
  }>(),
  { modelValue: () => ({}), langs: () => DEFAULT_TRANS_LANGS, prefix: '' },
)
const emit = defineEmits<{ (e: 'update:modelValue', v: Record<string, string>): void }>()

function onInput(lang: string, val: string) {
  emit('update:modelValue', { ...(props.modelValue || {}), [lang]: val })
}
</script>

<style scoped>
.inline-trans { display: flex; gap: 8px; flex-wrap: wrap; }
.inline-trans .el-input { flex: 1; min-width: 130px; }
</style>
