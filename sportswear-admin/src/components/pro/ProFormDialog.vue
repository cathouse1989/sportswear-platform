<template>
  <el-dialog
    :model-value="modelValue"
    :title="title"
    :width="width"
    :destroy-on-close="destroyOnClose"
    @update:model-value="(v: boolean) => emit('update:modelValue', v)"
  >
    <el-form ref="formRef" :model="form" :rules="rules" :label-width="labelWidth">
      <slot />
    </el-form>
    <template #footer>
      <slot name="footer" :saving="saving" :close="close">
        <el-button @click="close">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleConfirm">保存</el-button>
      </slot>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

defineOptions({ name: 'ProFormDialog' })

const props = withDefaults(
  defineProps<{
    modelValue: boolean
    title?: string
    width?: string
    labelWidth?: string
    form: Record<string, any>
    rules?: FormRules
    destroyOnClose?: boolean
  }>(),
  { title: '', width: '560px', labelWidth: '90px', destroyOnClose: true },
)

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'submit', form: Record<string, any>): Promise<void> | void
}>()

const formRef = ref<FormInstance>()
const saving = ref(false)

async function handleConfirm() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    // await emit 会等待父组件 submit 处理器返回的 Promise，实现"异步保存 + 按钮 loading"统一节奏
    await emit('submit', props.form)
  } finally {
    saving.value = false
  }
}

function close() {
  emit('update:modelValue', false)
}
</script>