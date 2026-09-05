<template>
  <div class="rich-editor" :class="{ 'is-focused': focused }" :data-v="version">
    <div v-if="editor" class="rich-editor__toolbar">
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('bold') }" title="加粗" @click="editor.chain().focus().toggleBold().run()"><b>B</b></button>
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('italic') }" title="斜体" @click="editor.chain().focus().toggleItalic().run()"><i>I</i></button>
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('underline') }" title="下划线" @click="editor.chain().focus().toggleUnderline().run()"><u>U</u></button>
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('strike') }" title="删除线" @click="editor.chain().focus().toggleStrike().run()"><s>S</s></button>
      <span class="rt-sep" />
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }" title="标题 1" @click="editor.chain().focus().toggleHeading({ level: 1 }).run()">H1</button>
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }" title="标题 2" @click="editor.chain().focus().toggleHeading({ level: 2 }).run()">H2</button>
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('heading', { level: 3 }) }" title="标题 3" @click="editor.chain().focus().toggleHeading({ level: 3 }).run()">H3</button>
      <span class="rt-sep" />
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('bulletList') }" title="无序列表" @click="editor.chain().focus().toggleBulletList().run()">• 列表</button>
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('orderedList') }" title="有序列表" @click="editor.chain().focus().toggleOrderedList().run()">1. 列表</button>
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('blockquote') }" title="引用" @click="editor.chain().focus().toggleBlockquote().run()">引用</button>
      <button type="button" class="rt-btn" :class="{ 'is-active': editor.isActive('codeBlock') }" title="代码块" @click="editor.chain().focus().toggleCodeBlock().run()">代码</button>
      <span class="rt-sep" />
      <button type="button" class="rt-btn" title="链接" @click="setLink">链接</button>
      <button type="button" class="rt-btn" title="水平分割线" @click="editor.chain().focus().setHorizontalRule().run()">—</button>
      <span class="rt-sep" />
      <button type="button" class="rt-btn" :disabled="!editor.can().undo()" title="撤销" @click="editor.chain().focus().undo().run()">↺</button>
      <button type="button" class="rt-btn" :disabled="!editor.can().redo()" title="重做" @click="editor.chain().focus().redo().run()">↻</button>
    </div>
    <EditorContent :editor="editor" class="rich-editor__content" :style="{ minHeight }" />
  </div>
</template>
<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import { Editor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import { ElMessageBox } from 'element-plus'

const props = withDefaults(defineProps<{ modelValue?: string; placeholder?: string; minHeight?: string }>(), {
  modelValue: '',
  placeholder: '请输入内容…',
  minHeight: '180px',
})
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()

const focused = ref(false)
// TipTap Editor 非响应式，通过递增 version 让工具栏 isActive 高亮实时刷新
const version = ref(0)

const editor = new Editor({
  content: props.modelValue,
  extensions: [
    StarterKit,
    Underline,
    Link.configure({ openOnClick: false, autolink: true }),
    Placeholder.configure({ placeholder: props.placeholder }),
  ],
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
    version.value++
  },
  onSelectionUpdate: () => { version.value++ },
  onFocus: () => { focused.value = true },
  onBlur: () => { focused.value = false },
})

// 外部值变化（编辑回填）时同步进编辑器，避免自身输入触发死循环
watch(() => props.modelValue, (val) => {
  if (val !== editor.getHTML()) {
    editor.commands.setContent(val || '', false)
  }
})

function setLink() {
  const prev = editor.getAttributes('link').href || ''
  ElMessageBox.prompt('请输入链接地址（http/https）', '插入链接', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputValue: prev,
    inputPattern: /^(https?:\/\/|\/|#|mailto:|tel:)/i,
    inputErrorMessage: '请输入 http(s):// 开头的链接',
  }).then(({ value }) => {
    const url = (value || '').trim()
    if (!url) {
      editor.chain().focus().unsetLink().run()
      return
    }
    editor.chain().focus().setLink({ href: url }).run()
  }).catch(() => {})
}

onBeforeUnmount(() => editor.destroy())
</script>
<style scoped>
.rich-editor {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: #fff;
  transition: border-color 0.2s;
}
.rich-editor.is-focused { border-color: #409eff; }
.rich-editor__toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  padding: 6px 8px;
  border-bottom: 1px solid #ebeef5;
  background: #fafafa;
}
.rt-btn {
  min-width: 28px;
  height: 28px;
  padding: 0 8px;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: #606266;
  font-size: 13px;
  cursor: pointer;
  line-height: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.rt-btn:hover { background: #e9e9eb; }
.rt-btn.is-active { background: #d9ecff; color: #409eff; }
.rt-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.rt-sep { width: 1px; height: 18px; background: #dcdfe6; margin: 5px 4px; }
.rich-editor__content {
  padding: 10px 12px;
  font-size: 14px;
  line-height: 1.7;
  color: #303133;
}
.rich-editor__content :deep(.ProseMirror) {
  outline: none;
}
.rich-editor__content :deep(.ProseMirror p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  color: #c0c4cc;
  float: left;
  height: 0;
  pointer-events: none;
}
</style>
