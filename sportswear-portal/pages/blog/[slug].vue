<template>
  <div class="max-w-3xl mx-auto px-4 py-12 md:py-16">
    <!-- 返回列表 -->
    <NuxtLink :to="localePath('/blog')" class="inline-flex items-center gap-1.5 text-sm text-gray-500 hover:text-black mb-8 transition-colors">
      <span aria-hidden="true">←</span> {{ $t('blog.back_to_list') }}
    </NuxtLink>

    <template v-if="blog">
      <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-gray-500 mb-5">
        <span class="px-2.5 py-1 rounded-full bg-blue-50 text-blue-700 font-medium">{{ categoryLabel(blog.category) }}</span>
        <time>{{ formatDate(blog.published_at || blog.created_at) }}</time>
        <span v-if="blog.author">· {{ blog.author }}</span>
      </div>
      <h1 class="text-3xl md:text-4xl font-bold leading-tight mb-8 md:mb-10">{{ blog.title }}</h1>
      <img v-if="blog.cover_image" :src="blog.cover_image" :alt="blog.title" class="w-full rounded-xl object-cover mb-10 md:mb-12 max-h-[440px] bg-gray-100" />
      <div class="prose prose-lg max-w-none" v-html="sanitizeHtml(blog.content)"></div>
      <div v-if="tags.length" class="mt-10 md:mt-12 flex flex-wrap gap-2">
        <span v-for="tag in tags" :key="tag" class="px-3 py-1 rounded-full bg-gray-100 text-gray-600 text-xs">#{{ tag }}</span>
      </div>

      <!-- 相关文章 -->
      <div v-if="blog.category && (relatedLoading || relatedBlogs.length)" class="mt-14 md:mt-16 border-t border-gray-100 pt-10 md:pt-12">
        <h2 class="text-xl md:text-2xl font-bold mb-6">{{ $t('blog.related') }}</h2>
        <!-- 骨架 -->
        <div v-if="relatedLoading" class="grid sm:grid-cols-2 gap-4 md:gap-6">
          <div v-for="i in 2" :key="i" class="animate-pulse bg-white rounded-xl overflow-hidden">
            <div class="aspect-[16/9] bg-gray-100" />
            <div class="p-4">
              <div class="h-4 bg-gray-100 rounded w-3/4 mb-2" />
              <div class="h-3 bg-gray-100 rounded w-1/2" />
            </div>
          </div>
        </div>
        <!-- 卡片 -->
        <div v-else class="grid sm:grid-cols-2 gap-4 md:gap-6">
          <NuxtLink
            v-for="r in relatedBlogs"
            :key="r.id"
            :to="localePath(`/blog/${r.slug}`)"
            class="group bg-white rounded-xl overflow-hidden shadow hover:shadow-lg transition flex flex-col active:scale-[0.98]"
          >
            <div class="aspect-[16/9] overflow-hidden bg-gray-100">
              <img v-if="r.cover_image" :src="r.cover_image" :alt="r.title" class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
              <div v-else class="w-full h-full flex items-center justify-center text-4xl">📝</div>
            </div>
            <div class="p-4 flex flex-col flex-1">
              <h3 class="text-sm font-semibold line-clamp-2 mb-1">{{ r.title }}</h3>
              <p v-if="r.summary" class="text-xs text-gray-500 line-clamp-2 mt-1">{{ r.summary }}</p>
            </div>
          </NuxtLink>
        </div>
      </div>
    </template>

    <div v-else-if="!pending" class="text-center py-24">
      <p class="text-gray-400 mb-5">{{ $t('blog.not_found') }}</p>
      <NuxtLink :to="localePath('/blog')" class="inline-flex items-center gap-1.5 text-blue-700 font-medium">
        <span aria-hidden="true">←</span> {{ $t('blog.back_to_list') }}
      </NuxtLink>
    </div>
    <div v-else class="text-center py-24 text-gray-400">{{ $t('common.loading') }}</div>
  </div>
</template>
<script setup lang="ts">
const route = useRoute()
const api = useApi()
const localePath = useLocalePath()
const { locale, t, te } = useI18n()
const slug = route.params.slug as string

// SSR 阶段拉取博客详情并内联进 HTML（提升 SEO + 配合 SWR HTML 缓存）
const { data: blog, pending } = await useAsyncData<any>(
  'blog-' + (locale.value || 'en') + '-' + slug,
  () => api.getBlog(slug).catch(() => null),
)

const tags = computed(() => {
  const raw = blog.value?.tags || ''
  return raw.split(',').map((s: string) => s.trim()).filter(Boolean)
})

// 相关文章（同分类，排除当前，客户端拉取）
const relatedBlogs = ref<any[]>([])
const relatedLoading = ref(true)
async function loadRelated() {
  if (!blog.value?.category) {
    relatedLoading.value = false
    return
  }
  try {
    const res = await api.getBlogs({ page: 1, pageSize: 3, category: blog.value.category })
    const items = res?.items || (Array.isArray(res) ? res : [])
    relatedBlogs.value = (items || []).filter((b: any) => b.id !== blog.value.id).slice(0, 2)
  } catch {
    relatedBlogs.value = []
  } finally {
    relatedLoading.value = false
  }
}
onMounted(loadRelated)

// 本地化分类标签
function categoryLabel(value?: string) {
  if (!value) return ''
  const key = `blog.categories.${value}`
  return te(key) ? t(key) : value
}

// 本地化日期
function formatDate(value?: string) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  try {
    return new Intl.DateTimeFormat(locale.value || 'en', { year: 'numeric', month: 'long', day: 'numeric' }).format(d)
  } catch {
    return d.toISOString().slice(0, 10)
  }
}

// 去除 HTML 标签生成摘要
function stripHtml(html?: string) {
  return (html || '').replace(/<[^>]*>/g, ' ').replace(/&nbsp;/g, ' ').replace(/\s+/g, ' ').trim()
}

// 富文本安全渲染兜底：去除脚本/内嵌对象/事件属性/javascript: 协议，
// 兼容历史纯 textarea 录入的 HTML 数据（TipTap 输出本身已在白名单内）。
function sanitizeHtml(html?: string) {
  if (!html) return ''
  return html
    .replace(/<\s*(script|iframe|object|embed|style|link|meta)\b[^>]*>[\s\S]*?<\s*\/\s*\1\s*>/gi, '')
    .replace(/<\s*(script|iframe|object|embed|style|link|meta)\b[^>]*\/?\s*>/gi, '')
    .replace(/\son\w+\s*=\s*("[^"]*"|'[^']*'|[^\s>]+)/gi, '')
    .replace(/(href|src)\s*=\s*("|')?\s*javascript:[^"'\s>]*("|')?/gi, '$1=""')
}

// SEO - 博客详情页
useSeoHead({
  title: computed(() => blog.value?.title ? `${blog.value.title} - Sportswear Manufacturing Blog` : 'Blog Article'),
  description: computed(() => stripHtml(blog.value?.content).slice(0, 160) || ''),
  ogImage: computed(() => blog.value?.cover_image || undefined),
  ogType: 'article',
  schema: computed(() => blog.value ? buildArticleSchema({ ...blog.value, image: blog.value.cover_image }) : undefined),
})
</script>