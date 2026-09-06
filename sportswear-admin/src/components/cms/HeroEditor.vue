<template>
  <div class="hero-editor">
    <div class="hero-head">
      <div>
        <h3>轮播图管理</h3>
        <p>在此管理门户首页轮播图片。自动播放、指示点、左右箭头等展示参数已默认开启，无需逐项配置；门户轮播文案按当前访问语言自动展示。</p>
      </div>
      <div class="hero-head-actions">
        <el-button @click="loadHero" :loading="loading">刷新</el-button>
        <el-button :disabled="loading" @click="addSlide">添加一张</el-button>
        <el-button type="primary" :loading="saving" :disabled="loading" @click="saveHeroSlides">保存配置</el-button>
      </div>
    </div>

    <el-alert
      v-if="!loading && !pageId"
      type="warning"
      :closable="false"
      show-icon
      title="未指定首页"
      class="mb-16"
    />

    <div v-loading="loading">
      <template v-if="pageId">
        <div class="preview-strip" v-if="heroSlides.length">
          <button
            v-for="(s, idx) in heroSlides"
            :key="'p-' + idx"
            type="button"
            class="preview-thumb"
            :class="{ active: previewIndex === idx }"
            @click="previewIndex = idx"
          >
            <img v-if="s.image" :src="s.image" alt="" />
            <span v-else>空</span>
            <em>{{ idx + 1 }}</em>
          </button>
        </div>

        <div class="cache-bar">
          <div class="cache-title">门户缓存</div>
          <el-switch
            v-model="cacheEnabled"
            :loading="cacheToggling"
            :disabled="!cacheRedisOk && !cacheEnabled"
            active-text="读缓存"
            inactive-text="读DB配置"
            @change="onCacheToggle"
          />
          <el-button
            size="small"
            type="primary"
            plain
            :loading="cacheRefreshing"
            :disabled="!cacheRedisOk"
            @click="handleCacheRefresh"
          >
            刷新缓存
          </el-button>
          <span class="cache-hint">
            {{ cacheRedisOk ? (cacheEnabled ? '门户当前读取 Redis 缓存，内容变更自动生效' : '门户当前直读数据库配置，改完刷新即可看到') : 'Redis 不可用，门户直读数据库配置' }}
          </span>
        </div>

        <div class="slide-list">
          <article v-for="(s, idx) in heroSlides" :key="idx" class="slide-card">
            <div class="slide-media">
              <el-image
                v-if="s.image"
                :src="s.image"
                fit="cover"
                class="slide-img"
                :preview-src-list="[s.image]"
                preview-teleported
              />
              <div v-else class="slide-img empty">暂无图片</div>
              <span class="slide-badge">{{ idx + 1 }}</span>
            </div>
            <div class="slide-body">
              <div class="field">
                <label>图片</label>
                <MediaPicker v-model="s.image" />
              </div>
              <div class="field">
                <label>按钮链接</label>
                <el-input v-model="s.button_url" placeholder="/contact" clearable />
              </div>
              <div class="field">
                <label>标题 / 副标题（多语言）</label>
                <el-tabs v-model="textLang" class="text-tabs">
                  <el-tab-pane v-for="l in textLangs" :key="l.value" :name="l.value" :label="l.label">
                    <el-input v-model="heroTexts[idx].title[l.value]" placeholder="标题（留空则不在门户显示）" class="text-input" />
                    <el-input v-model="heroTexts[idx].subtitle[l.value]" placeholder="副标题（留空则不在门户显示）" class="text-input" />
                  </el-tab-pane>
                </el-tabs>
              </div>
            </div>
            <div class="slide-ops">
              <el-button :disabled="idx === 0" @click="moveSlide(idx, -1)">上移</el-button>
              <el-button :disabled="idx === heroSlides.length - 1" @click="moveSlide(idx, 1)">下移</el-button>
              <el-button type="danger" plain @click="removeSlide(idx)">删除</el-button>
            </div>
          </article>
        </div>

        <div class="button-text-editor">
          <div class="field">
            <label>按钮文字（所有轮播图共用，多语言）</label>
            <el-tabs v-model="textLang">
              <el-tab-pane v-for="l in textLangs" :key="l.value" :name="l.value" :label="l.label">
                <el-input v-model="buttonTexts[l.value]" placeholder="如 Get a Quote / 获取报价" />
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </template>
    </div>

  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { pageApi, portalCacheApi, i18nApi } from '@/api'
import type { HeroSlide } from '@/types'
import MediaPicker from '@/components/media/MediaPicker.vue'

const props = defineProps<{ pageId: string }>()
const emit = defineEmits<{ saved: [] }>()

const pageId = ref(props.pageId)
watch(() => props.pageId, (v) => { pageId.value = v; loadHero() })

const heroSlides = ref<HeroSlide[]>([])
const previewIndex = ref(0)

// 轮播文案（标题/副标题/按钮文字）四语言，保存时写入 i18n_entries（home.hero_title_N / home.hero_sub_N / home.get_quote）
const textLangs = [
  { value: 'en', label: 'English' },
  { value: 'zh', label: '中文' },
  { value: 'es', label: 'Español' },
  { value: 'fr', label: 'Français' },
]
const textLang = ref('en')
const heroTexts = ref<Array<{ title: Record<string, string>; subtitle: Record<string, string> }>>([])
const buttonTexts = ref<Record<string, string>>({ en: '', zh: '', es: '', fr: '' })

async function loadHeroTexts() {
  try {
    const res = await i18nApi.entries({ page: 1, pageSize: 200, keyword: 'home.hero_', language: '' })
    const dict: Record<string, Record<string, string>> = {}
    for (const it of (res as any).items || []) {
      if (!dict[it.key]) dict[it.key] = {}
      dict[it.key][it.language] = it.value || ''
    }
    heroTexts.value = heroSlides.value.map((_, i) => ({
      title: { en: '', zh: '', es: '', fr: '', ...(dict[`home.hero_title_${i + 1}`] || {}) },
      subtitle: { en: '', zh: '', es: '', fr: '', ...(dict[`home.hero_sub_${i + 1}`] || {}) },
    }))
    buttonTexts.value = { en: '', zh: '', es: '', fr: '', ...(dict['home.get_quote'] || {}) }
  } catch { /* 词条加载失败不阻塞轮播图编辑 */ }
}

// 展示参数固定为默认开启（自动播放/指示点/箭头/悬停暂停等不再暴露开关）
const DEFAULT_HERO_SETTINGS = {
  autoplay: true,
  interval_ms: 5000,
  transition: 'fade' as const,
  show_dots: true,
  show_arrows: true,
  pause_on_hover: true,
}

function defaultHeroSlides(): HeroSlide[] {
  return [
    { image: 'https://images.unsplash.com/photo-1599901860904-17e6ed7083a0?w=1600&h=900&fit=crop', title: '', subtitle: '', button_text: '', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=1600&h=900&fit=crop', title: '', subtitle: '', button_text: '', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=1600&h=900&fit=crop', title: '', subtitle: '', button_text: '', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1577221084712-45b0445d2b00?w=1600&h=900&fit=crop', title: '', subtitle: '', button_text: '', button_url: '/contact' },
  ]
}

function parseBannerConfig(config: string | unknown): HeroSlide[] {
  let obj: any = config
  if (typeof config === 'string') {
    try {
      obj = JSON.parse(config)
    } catch {
      return []
    }
  }
  if (typeof obj === 'string') {
    try {
      obj = JSON.parse(obj)
    } catch {
      return []
    }
  }
  if (!obj || typeof obj !== 'object') return []
  const norm = (s: any): HeroSlide => ({
    image: s?.image || '', title: s?.title || '', subtitle: s?.subtitle || '',
    button_text: s?.button_text || '', button_url: s?.button_url || '',
  })
  if (Array.isArray(obj.slides)) return (obj.slides as any[]).map(norm)
  if ('image' in obj) return [norm(obj)]
  return []
}

const loading = ref(false)
const saving = ref(false)

async function loadHero() {
  loading.value = true
  try {
    if (!pageId.value) { heroSlides.value = []; return }
    const page = await pageApi.get(pageId.value)
    const banner = (page.modules || []).find((m: any) => m.type === 'banner')
    const parsed = banner?.config ? parseBannerConfig(banner.config) : []
    heroSlides.value = parsed.length ? parsed : defaultHeroSlides()
    await loadHeroTexts()
    previewIndex.value = 0
  } catch {
    heroSlides.value = defaultHeroSlides()
  } finally {
    loading.value = false
  }
}

const addSlide = () => heroSlides.value.push({ image: '', title: '', subtitle: '', button_text: '', button_url: '/contact' })
const removeSlide = (idx: number) => { heroSlides.value.splice(idx, 1); if (!heroSlides.value.length) addSlide() }
function moveSlide(idx: number, dir: number) {
  const arr = heroSlides.value
  const to = idx + dir
  if (to < 0 || to >= arr.length) return
  const cur = arr[idx]
  arr.splice(idx, 1)
  arr.splice(to, 0, cur)
}

async function saveHeroSlides() {
  if (!pageId.value) { ElMessage.warning('未找到首页页面'); return }
  const valid = heroSlides.value
    .filter(s => s.image && s.image.trim())
    .map(s => ({ image: s.image, title: '', subtitle: '', button_text: '', button_url: s.button_url }))
  if (!valid.length) { ElMessage.warning('请至少为一张轮播图填写图片地址'); return }
  saving.value = true
  try {
    await pageApi.updateHeroSlides(pageId.value, valid, { ...DEFAULT_HERO_SETTINGS })
    await saveHeroTexts()
    ElMessage.success('已保存')
    emit('saved')
    await loadHero()
  } catch { /* handled */ }
  finally { saving.value = false }
}

async function saveHeroTexts() {
  const ops: Array<Promise<any>> = []
  heroTexts.value.forEach((t, i) => {
    for (const l of textLangs) {
      if ((t.title[l.value] ?? '').trim()) ops.push(i18nApi.upsert({ key: `home.hero_title_${i + 1}`, language: l.value, value: t.title[l.value].trim(), module: 'home' }))
      if ((t.subtitle[l.value] ?? '').trim()) ops.push(i18nApi.upsert({ key: `home.hero_sub_${i + 1}`, language: l.value, value: t.subtitle[l.value].trim(), module: 'home' }))
    }
  })
  for (const l of textLangs) {
    if ((buttonTexts.value[l.value] ?? '').trim()) ops.push(i18nApi.upsert({ key: 'home.get_quote', language: l.value, value: buttonTexts.value[l.value].trim(), module: 'home' }))
  }
  if (ops.length) await Promise.all(ops)
}

// ============ 门户缓存（读缓存 / 读DB配置 / 刷新缓存） ============
const cacheEnabled = ref(true)
const cacheRedisOk = ref(true)
const cacheToggling = ref(false)
const cacheRefreshing = ref(false)

function applyCacheStatus(data: any) {
  cacheRedisOk.value = !!data?.redis_ok
  cacheEnabled.value = !!data?.enabled
}

async function loadCacheStatus() {
  try {
    const data = await portalCacheApi.status()
    applyCacheStatus(data)
  } catch { /* 缓存状态不可用时（如无 Redis）不阻塞编辑 */ }
}

async function onCacheToggle(val: string | number | boolean) {
  const next = !!val
  cacheToggling.value = true
  try {
    const data = await portalCacheApi.setEnabled(next)
    applyCacheStatus(data)
    ElMessage.success(next ? '已切换为读缓存' : '已切换为读DB配置（直查数据库）')
  } catch {
    cacheEnabled.value = !next
  } finally {
    cacheToggling.value = false
  }
}

async function handleCacheRefresh() {
  cacheRefreshing.value = true
  try {
    const data = await portalCacheApi.refresh()
    applyCacheStatus(data)
    ElMessage.success('缓存已刷新')
  } catch { /* handled */ }
  finally { cacheRefreshing.value = false }
}

onMounted(() => { loadHero(); loadCacheStatus() })
</script>

<style scoped>
.hero-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 18px;
}
.hero-head h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: #0d1b2a;
}
.hero-head p {
  margin: 6px 0 0;
  font-size: 13px;
  color: #8b7d6b;
  line-height: 1.5;
  max-width: 640px;
}
.hero-head-actions { display: flex; gap: 8px; flex-shrink: 0; }
.mb-16 { margin-bottom: 16px; }
.pagination { margin-top: 16px; justify-content: flex-end; }

.preview-strip {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding-bottom: 12px;
  margin-bottom: 14px;
}
.preview-thumb {
  position: relative;
  width: 112px;
  height: 64px;
  border: 2px solid #eae5dd;
  border-radius: 10px;
  overflow: hidden;
  padding: 0;
  background: #f4f1ec;
  cursor: pointer;
  flex-shrink: 0;
}
.preview-thumb.active { border-color: #d4a853; box-shadow: 0 0 0 3px rgba(212, 168, 83, 0.25); }
.preview-thumb img { width: 100%; height: 100%; object-fit: cover; display: block; }
.preview-thumb span, .preview-thumb em {
  position: absolute;
  color: #fff;
  font-size: 11px;
}
.preview-thumb span { inset: 0; display: flex; align-items: center; justify-content: center; color: #999; }
.preview-thumb em {
  right: 4px; bottom: 4px;
  background: rgba(13, 27, 42, 0.7);
  border-radius: 999px;
  padding: 0 6px;
  font-style: normal;
}

.cache-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  border: 1px solid #eae5dd;
  border-radius: 14px;
  padding: 12px 18px;
  background: linear-gradient(180deg, #fff 0%, #fbf9f6 100%);
  margin-bottom: 16px;
}
.cache-title { font-size: 14px; font-weight: 700; color: #0d1b2a; }
.cache-hint { font-size: 12px; color: #8b7d6b; max-width: 520px; line-height: 1.5; }

.slide-list { display: flex; flex-direction: column; gap: 14px; }
.slide-card {
  display: grid;
  grid-template-columns: 220px 1fr auto;
  gap: 16px;
  align-items: stretch;
  border: 1px solid #eae5dd;
  border-radius: 14px;
  padding: 14px;
  background: #fff;
  box-shadow: 0 6px 18px rgba(13, 27, 42, 0.04);
}
.slide-media { position: relative; }
.slide-img {
  width: 100%;
  height: 140px;
  border-radius: 10px;
  overflow: hidden;
  display: block;
}
.slide-img.empty {
  background: #f3eee6;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #aaa;
  font-size: 13px;
}
.slide-badge {
  position: absolute;
  top: 8px;
  left: 8px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #d4a853;
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}
.slide-body { display: flex; flex-direction: column; gap: 10px; min-width: 0; }
.field { display: flex; flex-direction: column; gap: 6px; flex: 1; }
.field label { font-size: 12px; color: #6b7280; }
.field-row { display: flex; gap: 12px; }
.img-input { display: flex; gap: 8px; }
.text-tip {
  font-size: 12px;
  color: #8b7d6b;
  background: #fbf9f6;
  border: 1px dashed #eae5dd;
  border-radius: 8px;
  padding: 8px 10px;
  line-height: 1.6;
}
.slide-ops { display: flex; flex-direction: column; gap: 8px; justify-content: center; }

@media (max-width: 960px) {
  .slide-card { grid-template-columns: 1fr; }
  .slide-ops { flex-direction: row; }
  .hero-head { flex-direction: column; }
}
.text-tabs :deep(.el-tabs__content) { padding: 6px 0 0; }
.text-input { margin-bottom: 8px; }
.text-input:last-child { margin-bottom: 0; }
.button-text-editor {
  margin-top: 4px;
  border: 1px solid #eae5dd;
  border-radius: 14px;
  padding: 14px;
  background: #fff;
  box-shadow: 0 6px 18px rgba(13, 27, 42, 0.04);
}
</style>
