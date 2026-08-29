<template>
  <el-card shadow="never">
    <div class="toolbar">
      <div class="page-title">
        <h3>轮播图管理</h3>
        <p class="page-desc">配置门户首页顶部翻页大图（图片、文案、按钮与展示参数）。</p>
      </div>
      <div class="spacer" />
      <el-button @click="loadHero" :loading="loading">刷新</el-button>
      <el-button type="primary" :loading="saving" :disabled="loading || !heroPageId" @click="saveHeroSlides">保存</el-button>
    </div>

    <el-alert
      v-if="!loading && !heroPageId"
      type="warning"
      :closable="false"
      show-icon
      title="未找到首页页面"
      description="请先在「页面管理」中创建并发布 slug 为 home（或类型为 home）的页面，再配置轮播图。"
      class="mb-alert"
    />

    <div v-loading="loading">
      <template v-if="heroPageId">
        <div class="hero-tip">
          共 {{ heroSlides.length }} 张轮播大图（门户首页顶部滑动展示）。支持添加 / 删除 / 排序，图片可在媒体库选择或粘贴 URL。
          文案留空时默认使用各语言内置文案。展示参数（自动播放等）与门户保持一致。
        </div>

        <div class="hero-settings">
          <div class="hero-settings-title">展示参数</div>
          <div class="hero-settings-grid">
            <div class="hero-setting-item">
              <span class="hero-label">自动播放</span>
              <el-switch v-model="heroSettings.autoplay" />
            </div>
            <div class="hero-setting-item">
              <span class="hero-label">间隔(ms)</span>
              <el-input-number v-model="heroSettings.interval_ms" :min="2000" :max="30000" :step="500" :disabled="!heroSettings.autoplay" controls-position="right" />
            </div>
            <div class="hero-setting-item">
              <span class="hero-label">切换动效</span>
              <el-select v-model="heroSettings.transition" style="width: 120px">
                <el-option label="淡入淡出" value="fade" />
                <el-option label="左右滑动" value="slide" />
              </el-select>
            </div>
            <div class="hero-setting-item">
              <span class="hero-label">指示点</span>
              <el-switch v-model="heroSettings.show_dots" />
            </div>
            <div class="hero-setting-item">
              <span class="hero-label">左右箭头</span>
              <el-switch v-model="heroSettings.show_arrows" />
            </div>
            <div class="hero-setting-item">
              <span class="hero-label">悬停暂停</span>
              <el-switch v-model="heroSettings.pause_on_hover" :disabled="!heroSettings.autoplay" />
            </div>
          </div>
        </div>

        <div class="hero-list">
          <div v-for="(s, idx) in heroSlides" :key="idx" class="hero-card">
            <div class="hero-index">{{ idx + 1 }}</div>
            <div class="hero-thumb">
              <el-image v-if="s.image" :src="s.image" fit="cover" style="width: 120px; height: 80px; border-radius: 6px" :preview-src-list="[s.image]" preview-teleported />
              <div v-else class="hero-thumb-empty">暂无图片</div>
            </div>
            <div class="hero-fields">
              <div class="hero-field">
                <span class="hero-label">图片</span>
                <div class="img-input">
                  <el-input v-model="s.image" placeholder="图片 URL（支持 /uploads/... 相对路径）" clearable />
                  <el-button size="small" @click="openMediaPicker(idx)">媒体库</el-button>
                </div>
              </div>
              <div class="hero-field">
                <span class="hero-label">标题</span>
                <el-input v-model="s.title" placeholder="轮播标题（留空用语言内置文案）" clearable />
              </div>
              <div class="hero-field">
                <span class="hero-label">副标题</span>
                <el-input v-model="s.subtitle" placeholder="轮播副标题（留空用语言内置文案）" clearable />
              </div>
              <div class="hero-field">
                <span class="hero-label">按钮</span>
                <el-input v-model="s.button_text" placeholder="按钮文字（默认 Get a Quote）" clearable style="width: 200px" />
                <el-input v-model="s.button_url" placeholder="按钮链接（如 /contact）" clearable style="width: 260px" />
              </div>
            </div>
            <div class="hero-ops">
              <el-button size="small" :disabled="idx === 0" @click="moveSlide(idx, -1)">↑</el-button>
              <el-button size="small" :disabled="idx === heroSlides.length - 1" @click="moveSlide(idx, 1)">↓</el-button>
              <el-button size="small" type="danger" @click="removeSlide(idx)">删除</el-button>
            </div>
          </div>
        </div>
        <div class="footer-actions">
          <el-button @click="addSlide">添加一张</el-button>
          <el-button type="primary" :loading="saving" @click="saveHeroSlides">保存</el-button>
        </div>
      </template>
    </div>

    <!-- 媒体库选择对话框 -->
    <el-dialog v-model="mediaDialogVisible" title="从媒体库选择图片" width="760px">
      <div class="media-toolbar">
        <el-input v-model="mediaKeyword" placeholder="搜索文件名" size="small" style="width: 200px" clearable @keyup.enter="searchMedia" />
        <el-button size="small" @click="searchMedia">查询</el-button>
      </div>
      <div class="media-grid" v-loading="mediaLoading">
        <div v-for="m in mediaItems" :key="m.id" class="media-cell">
          <el-image v-if="m.type === 'image'" :src="m.url" fit="cover" class="media-img" />
          <div v-else class="media-file">{{ m.type }}</div>
          <div class="media-name" :title="m.original_name">{{ m.original_name }}</div>
          <el-button size="small" type="primary" @click="pickMedia(m)">选用</el-button>
        </div>
        <div v-if="!mediaLoading && !mediaItems.length" class="media-empty">暂无图片，请先在「媒体管理」上传</div>
      </div>
      <el-pagination
        class="pagination"
        v-model:current-page="mediaPage"
        v-model:page-size="mediaPageSize"
        :total="mediaTotal"
        layout="total, prev, pager, next"
        @current-change="loadMedia"
      />
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { pageApi, mediaApi } from '@/api'
import type { Page, HeroSlide, HeroSettings } from '@/types'

// ============ 首页轮播图配置 ============
const isHome = (row: Page) => row.slug === 'home' || row.type === 'home'

const heroPageId = ref('')
const heroSlides = ref<HeroSlide[]>([])
const heroSettings = reactive<Required<Pick<HeroSettings, 'autoplay' | 'interval_ms' | 'transition' | 'show_dots' | 'show_arrows' | 'pause_on_hover'>>>({
  autoplay: true,
  interval_ms: 5000,
  transition: 'fade',
  show_dots: true,
  show_arrows: true,
  pause_on_hover: true,
})

function defaultHeroSettings() {
  return {
    autoplay: true,
    interval_ms: 5000,
    transition: 'fade' as const,
    show_dots: true,
    show_arrows: true,
    pause_on_hover: true,
  }
}

function applyHeroSettings(raw?: Partial<HeroSettings> | null) {
  const d = defaultHeroSettings()
  const src = raw || {}
  heroSettings.autoplay = src.autoplay !== false
  heroSettings.interval_ms = Number(src.interval_ms) > 0 ? Number(src.interval_ms) : d.interval_ms
  heroSettings.transition = src.transition === 'slide' ? 'slide' : 'fade'
  heroSettings.show_dots = src.show_dots !== false
  heroSettings.show_arrows = src.show_arrows !== false
  heroSettings.pause_on_hover = src.pause_on_hover !== false
}

// 默认四张轮播图（与门户内置默认一致，后台可在此基础上编辑）

function defaultHeroSlides(): HeroSlide[] {
  return [
    { image: 'https://images.unsplash.com/photo-1599901860904-17e6ed7083a0?w=1600&h=900&fit=crop', title: 'Custom Sportswear Manufacturer', subtitle: 'OEM & ODM Solutions for Global Brands', button_text: 'Get a Quote', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=1600&h=900&fit=crop', title: 'Advanced Running & Training Apparel', subtitle: 'Technical fabrics for peak performance', button_text: 'Get a Quote', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=1600&h=900&fit=crop', title: 'Professional Team Uniforms', subtitle: 'Full custom sublimation printing', button_text: 'Get a Quote', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1577221084712-45b0445d2b00?w=1600&h=900&fit=crop', title: 'From Concept to Product', subtitle: 'End-to-end manufacturing support', button_text: 'Get a Quote', button_url: '/contact' },
  ]
}

const loading = ref(false)
const saving = ref(false)

async function findHomePage(): Promise<Page | null> {
  const tryList = async (status?: string) => {
    const result = await pageApi.list({ page: 1, pageSize: 100, status })
    const items: Page[] = (result as any).items || result || []
    return items.find(isHome) || null
  }
  return (await tryList('published')) || (await tryList()) || null
}

async function loadHero() {
  loading.value = true
  heroPageId.value = ''
  applyHeroSettings(defaultHeroSettings())
  try {
    const home = await findHomePage()
    if (!home) { heroSlides.value = []; return }
    heroPageId.value = home.id
    const page = await pageApi.get(home.id)
    const banner = (page.modules || []).find((m: any) => m.type === 'banner')
    const parsed = banner?.config ? parseBannerConfig(banner.config) : { slides: [], settings: null }
    heroSlides.value = parsed.slides.length ? parsed.slides : defaultHeroSlides()
    applyHeroSettings(parsed.settings)
  } catch {
    heroSlides.value = defaultHeroSlides()
    applyHeroSettings(defaultHeroSettings())
  } finally {
    loading.value = false
  }
}


// 解析 banner 模块 config：兼容 slides 数组 / 旧单对象，以及可选 settings
function parseBannerConfig(config: string | unknown): { slides: HeroSlide[]; settings: Partial<HeroSettings> | null } {
  if (typeof config !== 'string') return { slides: [], settings: null }
  try {
    const obj = JSON.parse(config)
    const norm = (s: any): HeroSlide => ({
      image: s?.image || '', title: s?.title || '', subtitle: s?.subtitle || '',
      button_text: s?.button_text || '', button_url: s?.button_url || '',
    })
    let slides: HeroSlide[] = []
    if (obj && Array.isArray(obj.slides)) slides = (obj.slides as any[]).map(norm)
    else if (obj && !Array.isArray(obj) && typeof obj === 'object' && 'image' in obj) slides = [norm(obj)]
    const settings = obj && typeof obj === 'object' && obj.settings && typeof obj.settings === 'object'
      ? (obj.settings as Partial<HeroSettings>)
      : null
    return { slides, settings }
  } catch { /* ignore */ }
  return { slides: [], settings: null }
}


const addSlide = () => heroSlides.value.push({ image: '', title: '', subtitle: '', button_text: 'Get a Quote', button_url: '/contact' })
const removeSlide = (idx: number) => { heroSlides.value.splice(idx, 1); if (!heroSlides.value.length) addSlide() }
function moveSlide(idx: number, dir: number) {
  const arr = heroSlides.value
  const to = idx + dir
  if (to < 0 || to >= arr.length) return
  const t = arr[idx]; arr[idx] = arr[to]; arr[to] = t
}

async function saveHeroSlides() {
  if (!heroPageId.value) { ElMessage.warning('未找到首页页面，无法保存'); return }
  const valid = heroSlides.value.filter(s => s.image && s.image.trim())
  if (!valid.length) { ElMessage.warning('请至少为一张轮播图填写图片地址'); return }
  saving.value = true
  try {
    await pageApi.updateHeroSlides(heroPageId.value, valid, {
      autoplay: heroSettings.autoplay,
      interval_ms: heroSettings.interval_ms,
      transition: heroSettings.transition,
      show_dots: heroSettings.show_dots,
      show_arrows: heroSettings.show_arrows,
      pause_on_hover: heroSettings.pause_on_hover,
    })
    ElMessage.success('轮播图保存成功，刷新门户/门户预览即可查看')
  } catch { /* error handled */ }
  finally { saving.value = false }
}



// ============ 媒体库选择 ============
const mediaDialogVisible = ref(false)
const mediaKeyword = ref('')
const mediaItems = ref<any[]>([])
const mediaTotal = ref(0)
const mediaPage = ref(1)
const mediaPageSize = ref(12)
const mediaLoading = ref(false)
const editingSlideIndex = ref(-1)

async function loadMedia() {
  mediaLoading.value = true
  try {
    const res = await mediaApi.list({ page: mediaPage.value, pageSize: mediaPageSize.value, type: 'image', keyword: mediaKeyword.value || undefined })
    mediaItems.value = (res as any).items || res || []
    mediaTotal.value = (res as any).total ?? mediaItems.value.length
  } finally { mediaLoading.value = false }
}
function searchMedia() { mediaPage.value = 1; loadMedia() }
function openMediaPicker(idx: number) {
  editingSlideIndex.value = idx
  mediaPage.value = 1
  mediaDialogVisible.value = true
  loadMedia()
}
function pickMedia(m: any) {
  if (editingSlideIndex.value >= 0) heroSlides.value[editingSlideIndex.value].image = m.url
  mediaDialogVisible.value = false
}


onMounted(loadHero)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: flex-start; }
.page-title h3 { margin: 0; font-size: 16px; font-weight: 600; color: #0d1b2a; }
.page-desc { margin: 4px 0 0; font-size: 12px; color: #8b7d6b; }
.spacer { flex: 1; }
.mb-alert { margin-bottom: 16px; }
.pagination { margin-top: 16px; justify-content: flex-end; }
.footer-actions { display: flex; gap: 12px; margin-top: 16px; justify-content: flex-end; }
/* 轮播图配置 */
.hero-loading { min-height: 120px; }
.hero-tip {
  color: #8b7d6b; font-size: 12px; line-height: 1.5;
  background: #fbf9f6; border: 1px solid #eae5dd; border-radius: 8px;
  padding: 10px 14px; margin-bottom: 14px;
}
.hero-settings {
  border: 1px solid #eae5dd; border-radius: 10px; padding: 12px 14px;
  background: #fff; margin-bottom: 14px;
}
.hero-settings-title {
  font-size: 13px; font-weight: 600; color: #0d1b2a; margin-bottom: 10px;
}
.hero-settings-grid {
  display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 10px 16px; align-items: center;
}
.hero-setting-item { display: flex; align-items: center; gap: 10px; }
.hero-list { display: flex; flex-direction: column; gap: 12px; }

.hero-card {
  display: flex; gap: 14px; align-items: flex-start;
  border: 1px solid #eae5dd; border-radius: 10px; padding: 14px;
  background: #fff;
}
.hero-index {
  width: 22px; height: 22px; flex-shrink: 0; border-radius: 50%;
  background: #d4a853; color: #fff; font-size: 12px; font-weight: 600;
  display: flex; align-items: center; justify-content: center; margin-top: 30px;
}
.hero-thumb { flex-shrink: 0; }
.hero-thumb-empty {
  width: 120px; height: 80px; border-radius: 6px; background: #f0eded;
  display: flex; align-items: center; justify-content: center;
  color: #aaa; font-size: 12px;
}
.hero-fields { flex: 1; display: flex; flex-direction: column; gap: 10px; }
.hero-field { display: flex; align-items: center; gap: 10px; }
.hero-label { width: 52px; flex-shrink: 0; color: #555; font-size: 13px; }
.img-input { display: flex; gap: 8px; flex: 1; }
.hero-ops { display: flex; flex-direction: column; gap: 6px; flex-shrink: 0; margin-top: 8px; }

/* 媒体库选择 */
.media-toolbar { display: flex; gap: 10px; margin-bottom: 12px; }
.media-grid {
  display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px; min-height: 80px;
}
.media-cell {
  border: 1px solid #eae5dd; border-radius: 8px; padding: 8px;
  display: flex; flex-direction: column; align-items: center; gap: 6px;
}
.media-img { width: 100%; height: 90px; border-radius: 6px; }
.media-file {
  width: 100%; height: 90px; display: flex; align-items: center; justify-content: center;
  background: #f5f5f5; border-radius: 6px; color: #999; font-size: 12px;
}
.media-name {
  font-size: 11px; color: #666; width: 100%; text-align: center;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.media-empty { grid-column: 1 / -1; text-align: center; color: #999; padding: 30px 0; }
</style>
