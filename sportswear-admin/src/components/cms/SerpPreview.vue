<template>
  <div class="serp-preview">
    <div class="serp-title">{{ shortTitle }}</div>
    <div class="serp-url">{{ url }}</div>
    <div class="serp-desc">{{ shortDesc }}</div>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ title?: string; url?: string; description?: string }>()
const shortTitle = computed(() => {
  const t = (props.title || '').trim()
  return t.length > 60 ? t.slice(0, 57) + '…' : (t || '（未填写标题）')
})
const shortDesc = computed(() => {
  const d = (props.description || '').trim()
  return d.length > 160 ? d.slice(0, 157) + '…' : (d || '（未填写描述）')
})
</script>
<style scoped>
.serp-preview { font-family: Arial, sans-serif; padding: 12px 14px; border: 1px solid #e4e7ed; border-radius: 8px; background: #fff; }
.serp-title { color: #1a0dab; font-size: 18px; line-height: 1.3; margin-bottom: 2px; word-break: break-all; }
.serp-url { color: #006621; font-size: 13px; margin-bottom: 4px; word-break: break-all; }
.serp-desc { color: #545454; font-size: 13px; line-height: 1.45; word-break: break-all; }
</style>