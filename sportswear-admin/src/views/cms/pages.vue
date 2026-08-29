<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input
        v-model="keyword"
        placeholder="搜索标题 / Slug"
        clearable
        style="width: 240px"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="status" placeholder="状态" clearable style="width: 150px" @change="handleSearch">
        <el-option label="草稿" value="draft" />
        <el-option label="已发布" value="published" />
        <el-option label="已下线" value="offline" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建页面</el-button>
    </div>

    <el-table :data="pages" v-loading="loading" stripe>
      <el-table-column prop="title" label="标题" min-width="160" />
      <el-table-column prop="slug" label="Slug" min-width="140" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }"><el-tag size="small">{{ row.type }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="360">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="warning" @click="openHeroDialog(row)" v-if="isHome(row)">轮播图</el-button>
          <el-button v-permission="'page:publish'" size="small" type="success" @click="handlePublish(row)" v-if="row.status !== 'published'">发布</el-button>
          <el-button v-permission="'page:publish'" size="small" type="warning" @click="handleUnpublish(row)" v-else>下线</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="loadData"
    />

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑页面' : '新建页面'" width="560px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="Slug" required>
          <el-input v-model="form.slug" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type">
            <el-option label="普通" value="normal" />
            <el-option label="首页" value="home" />
            <el-option label="产品" value="product" />
            <el-option label="OEM" value="oem" />
            <el-option label="ODM" value="odm" />
            <el-option label="工厂" value="factory" />
            <el-option label="FAQ" value="faq" />
            <el-option label="联系我们" value="contact" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- 首页轮播图配置对话框 -->
    <el-dialog v-model="heroDialogVisible" title="首页轮播图配置" width="960px" top="5vh">
      <div v-if="heroLoading" v-loading="true" class="hero-loading" />
      <template v-else>
        <div class="hero-tip">
          共 {{ heroSlides.length }} 张轮播大图（门户首页顶部滑动展示）。支持添加 / 删除 / 排序，图片可在媒体库选择或粘贴 URL。
          文案留空时默认使用各语言内置文案。
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
      </template>
      <template #footer>
        <el-button @click="heroDialogVisible = false">取消</el-button>
        <el-button @click="addSlide">添加一张</el-button>
        <el-button type="primary" :disabled="heroLoading" @click="saveHeroSlides">保存</el-button>
      </template>
    </el-dialog>

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
import { ElMessage, ElMessageBox } from 'element-plus'
import { pageApi, mediaApi } from '@/api'
import type { Page, HeroSlide } from '@/types'

const pages = ref<Page[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const status = ref('')

const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({ title: '', slug: '', type: 'normal', sort_order: 0 })

async function loadData() {
  loading.value = true
  try {
    const result = await pageApi.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, status: status.value })
    pages.value = (result as any).items || result || []
    total.value = (result as any).total ?? pages.value.length
  } finally { loading.value = false }
}
function handleSearch() { page.value = 1; loadData() }
function openCreateDialog() { editingId.value = ''; Object.assign(form, { title: '', slug: '', type: 'normal', sort_order: 0 }); dialogVisible.value = true }
function openEditDialog(row: Page) { editingId.value = row.id; Object.assign(form, { title: row.title, slug: row.slug, type: row.type, sort_order: row.sort_order }); dialogVisible.value = true }
async function handleSave() {
  try {
    if (editingId.value) { await pageApi.update(editingId.value, form) } else { await pageApi.create(form) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch { /* error handled */ }
}
async function handlePublish(row: Page) { await pageApi.publish(row.id); ElMessage.success('已发布'); loadData() }
async function handleUnpublish(row: Page) { await pageApi.unpublish(row.id); ElMessage.success('已下线'); loadData() }
async function handleDelete(row: Page) { await ElMessageBox.confirm(`确定删除页面 ${row.title} 吗？`, '警告', { type: 'warning' }); await pageApi.delete(row.id); ElMessage.success('已删除'); loadData() }

// ============ 首页轮播图配置 ============
const isHome = (row: Page) => row.slug === 'home' || row.type === 'home'

const heroDialogVisible = ref(false)
const heroLoading = ref(false)
const heroPageId = ref('')
const heroSlides = ref<HeroSlide[]>([])

// 默认四张轮播图（与门户内置默认一致，后台可在此基础上编辑）
function defaultHeroSlides(): HeroSlide[] {
  return [
    { image: 'https://images.unsplash.com/photo-1599901860904-17e6ed7083a0?w=1600&h=900&fit=crop', title: 'Custom Sportswear Manufacturer', subtitle: 'OEM & ODM Solutions for Global Brands', button_text: 'Get a Quote', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=1600&h=900&fit=crop', title: 'Advanced Running & Training Apparel', subtitle: 'Technical fabrics for peak performance', button_text: 'Get a Quote', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=1600&h=900&fit=crop', title: 'Professional Team Uniforms', subtitle: 'Full custom sublimation printing', button_text: 'Get a Quote', button_url: '/contact' },
    { image: 'https://images.unsplash.com/photo-1577221084712-45b0445d2b00?w=1600&h=900&fit=crop', title: 'From Concept to Product', subtitle: 'End-to-end manufacturing support', button_text: 'Get a Quote', button_url: '/contact' },
  ]
}

async function openHeroDialog(row: Page) {
  heroPageId.value = row.id
  heroDialogVisible.value = true
  heroLoading.value = true
  try {
    const page = await pageApi.get(row.id)
    const banner = (page.modules || []).find((m: any) => m.type === 'banner')
    const parsed = banner?.config ? parseBannerConfig(banner.config) : []
    heroSlides.value = parsed.length ? parsed : defaultHeroSlides()
  } catch {
    heroSlides.value = defaultHeroSlides()
  } finally {
    heroLoading.value = false
  }
}

// 解析 banner 模块 config：兼容 {"slides":[...]} 与旧的单对象 {"image":"...",...} 两种结构
function parseBannerConfig(config: string | unknown): HeroSlide[] {
  if (typeof config !== 'string') return []
  try {
    const obj = JSON.parse(config)
    const norm = (s: any): HeroSlide => ({
      image: s?.image || '', title: s?.title || '', subtitle: s?.subtitle || '',
      button_text: s?.button_text || '', button_url: s?.button_url || '',
    })
    if (obj && Array.isArray(obj.slides)) return (obj.slides as any[]).map(norm)
    if (obj && !Array.isArray(obj) && typeof obj === 'object' && 'image' in obj) return [norm(obj)]
  } catch { /* ignore */ }
  return []
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
  const valid = heroSlides.value.filter(s => s.image && s.image.trim())
  if (!valid.length) { ElMessage.warning('请至少为一张轮播图填写图片地址'); return }
  try {
    await pageApi.updateHeroSlides(heroPageId.value, valid)
    ElMessage.success('轮播图保存成功，刷新门户/门户预览即可查看')
    heroDialogVisible.value = false
  } catch { /* error handled */ }
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

onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }

/* 轮播图配置 */
.hero-loading { min-height: 120px; }
.hero-tip {
  color: #8b7d6b; font-size: 12px; line-height: 1.5;
  background: #fbf9f6; border: 1px solid #eae5dd; border-radius: 8px;
  padding: 10px 14px; margin-bottom: 14px;
}
.hero-list { display: flex; flex-direction: column; gap: 12px; max-height: 60vh; overflow-y: auto; }
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