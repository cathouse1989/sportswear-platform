<template>
  <el-card shadow="never" v-loading="loading">
    <!-- 页面未初始化 -->
    <el-result
      v-if="!loading && !pageId"
      icon="warning"
      title="未找到「关于我们」页面"
      sub-title="请先初始化关于我们页面（slug=about），之后即可在此统一维护横幅、公司故事、数据统计、认证墙、生产流程与行动号召等内容。"
    >
      <template #extra>
        <el-button type="primary" :loading="initializing" @click="initPage">初始化关于我们页面</el-button>
      </template>
    </el-result>

    <template v-else-if="pageId">
      <div class="head">
        <div>
          <h3 class="head-title">关于我们</h3>
          <p class="head-sub">在此统一维护门户「关于我们」页的横幅、公司故事、数据统计、认证墙、生产流程与行动号召模块，保存后前台即时生效。</p>
        </div>
        <div class="head-actions">
          <el-button @click="loadPage" :loading="loading">刷新</el-button>
          <el-button v-if="form.status !== 'published'" type="success" :loading="saving" @click="handlePublish">发布</el-button>
          <el-button v-else type="warning" :loading="saving" @click="handleUnpublish">下线</el-button>
          <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
        </div>
      </div>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="模块内容" name="modules">
          <div class="module-toolbar">
            <el-button size="small" type="primary" @click="addModule">添加模块</el-button>
            <span class="muted">共 {{ modules.length }} 个模块，可上下调整顺序、控制显隐</span>
          </div>
          <div v-for="(m, i) in modules" :key="i" class="module-item">
            <div class="module-head">
              <span class="module-index">{{ i + 1 }}</span>
              <el-select v-model="m.type" size="small" style="width: 180px">
                <el-option v-for="t in ABOUT_MODULE_OPTIONS" :key="t.value" :label="t.label" :value="t.value" />
              </el-select>
              <el-input v-model="m.title" size="small" placeholder="模块标题（英文）" style="flex: 1" />
              <el-switch v-model="m.is_visible" size="small" active-text="显示" />
              <el-button-group>
                <el-button size="small" :disabled="i === 0" @click="moveModule(i, -1)">↑</el-button>
                <el-button size="small" :disabled="i === modules.length - 1" @click="moveModule(i, 1)">↓</el-button>
                <el-button size="small" type="danger" @click="removeModule(i)">删除</el-button>
              </el-button-group>
            </div>
            <ModuleConfigEditor v-model="m.config" :type="m.type" />
          </div>
          <el-empty v-if="!modules.length" description="暂无模块，点击「添加模块」开始搭建" :image-size="60" />
        </el-tab-pane>

        <el-tab-pane label="页面信息与多语言" name="basic">
          <el-form label-width="90px" style="max-width: 560px">
            <el-form-item label="页面标题"><el-input v-model="form.title" /></el-form-item>
            <el-form-item label="Slug"><el-input :model-value="form.slug" disabled /></el-form-item>
            <el-form-item label="状态">
              <el-tag :type="enumTag('content.status', form.status)" size="small">{{ enumLabel('content.status', form.status) }}</el-tag>
            </el-form-item>
          </el-form>
          <el-divider content-position="left">页面标题多语言（未填写回退英文）</el-divider>
          <TransEditor v-model="translations" :fields="TRANS_FIELDS" :source="{ title: form.title }" :langs="TRANS_LANGS" source-label="English" />
        </el-tab-pane>
      </el-tabs>
    </template>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { pageApi } from '@/api'
import ModuleConfigEditor from '@/components/cms/ModuleConfigEditor.vue'
import TransEditor from '@/components/cms/TransEditor.vue'
import { DEFAULT_TRANS_LANGS, createEmptyTranslations, translationsToRecord, translationsToPayload } from '@/composables/useTransRecord'
import { useEnumDict } from '@/composables/useEnumDict'

const { ensureLoaded: loadEnumDict, label: enumLabel, tagType: enumTag } = useEnumDict()

const loading = ref(false)
const saving = ref(false)
const initializing = ref(false)
const pageId = ref('')
const activeTab = ref('modules')

const form = reactive({ title: '', slug: 'about', type: 'normal', template: '', sort_order: 0, status: 'draft' })
const modules = ref<any[]>([])

const TRANS_LANGS = DEFAULT_TRANS_LANGS
const TRANS_FIELDS = [{ key: 'title', label: '标题' }]
const PAGE_TRANS_FIELDS = TRANS_FIELDS.map((f) => f.key)
const translations = ref<Record<string, Record<string, string>>>({})

const ABOUT_MODULE_OPTIONS = [
  { value: 'about_hero', label: '横幅 Hero' },
  { value: 'about_story', label: '公司故事' },
  { value: 'about_stats', label: '数据统计' },
  { value: 'about_certifications', label: '认证墙' },
  { value: 'about_process', label: '生产流程' },
  { value: 'about_cta', label: '行动号召' },
]

// 「关于我们」默认 6 个模块（与 seed-about-page.sql 保持一致），用于首次初始化
const DEFAULT_ABOUT_MODULES = [
  {
    type: 'about_hero',
    title: 'About Us',
    config: {
      title: { en: 'About Us', zh: '关于我们', es: 'Sobre Nosotros', fr: 'À Propos' },
      subtitle: {
        en: 'Professional OEM/ODM sportswear manufacturer with 15+ years of experience.',
        zh: '拥有15年以上经验的专业OEM/ODM运动服装制造商。',
        es: 'Fabricante profesional OEM/ODM de ropa deportiva con más de 15 años de experiencia.',
        fr: `Fabricant professionnel OEM/ODM de vêtements de sport avec plus de 15 ans d'expérience.`,
      },
      bg_image: '',
    },
  },
  {
    type: 'about_story',
    title: 'Our Story',
    config: {
      title: { en: 'Our Story', zh: '我们的故事', es: 'Nuestra Historia', fr: 'Notre Histoire' },
      body: {
        en: '<p>Founded in 2008, we have grown from a small workshop into a leading OEM/ODM sportswear manufacturer serving global brands.</p><p>We specialize in high-performance activewear, yoga wear, running gear, team uniforms and custom sportswear solutions for brands worldwide.</p>',
        zh: '<p>我们成立于 2008 年，已从一间小作坊成长为服务全球品牌的领先 OEM/ODM 运动服装制造商。</p><p>我们专注于高性能运动服、瑜伽服、跑步装备、团队队服及定制运动服装解决方案。</p>',
        es: '<p>Fundada en 2008, hemos pasado de ser un pequeño taller a un fabricante líder OEM/ODM de ropa deportiva.</p><p>Nos especializamos en ropa deportiva de alto rendimiento, yoga, running, uniformes y soluciones personalizadas.</p>',
        fr: `<p>Fondée en 2008, nous sommes passés d'un petit atelier à un fabricant leader OEM/ODM de vêtements de sport.</p><p>Nous sommes spécialisés dans les vêtements de sport haute performance, yoga, running et uniformes.</p>`,
      },
      image: '',
      image_side: 'left',
    },
  },
  {
    type: 'about_stats',
    title: '',
    config: {
      items: [
        { value: '15+', label: { en: 'Years Experience', zh: '年行业经验', es: 'Años de Experiencia', fr: `Ans d'Expérience` } },
        { value: '500+', label: { en: 'Global Clients', zh: '全球客户', es: 'Clientes Globales', fr: 'Clients Mondiaux' } },
        { value: '50K', label: { en: 'sqm Facility', zh: '平方米工厂', es: 'm² de Instalaciones', fr: `m² d'Installations` } },
      ],
    },
  },
  {
    type: 'about_certifications',
    title: 'Our Certifications',
    config: {
      title: { en: 'Our Certifications', zh: '资质认证', es: 'Certificaciones', fr: 'Certifications' },
      ref_ids: [],
    },
  },
  {
    type: 'about_process',
    title: 'Factory Tour',
    config: {
      title: { en: 'Factory Tour', zh: '工厂参观', es: 'Visita a la Fábrica', fr: `Visite de l'Usine` },
      steps: [
        { icon: '✂️', title: { en: 'Cutting', zh: '裁剪', es: 'Corte', fr: 'Coupe' }, desc: { en: 'Computerized cutting with 0.1mm precision', zh: '电脑裁剪机，精度高达 0.1mm', es: 'Corte computarizado con precisión de 0.1mm', fr: 'Découpe informatisée avec une précision de 0.1mm' } },
        { icon: '🧵', title: { en: 'Sewing', zh: '缝制', es: 'Costura', fr: 'Couture' }, desc: { en: '500+ skilled workers on modern lines', zh: '500 余名熟练工人操作现代化流水线', es: '500+ trabajadores calificados', fr: 'Plus de 500 ouvriers qualifiés' } },
        { icon: '🖨️', title: { en: 'Printing', zh: '印花', es: 'Impresión', fr: 'Impression' }, desc: { en: 'Sublimation, screen print & embroidery', zh: '升华、丝印与刺绣', es: 'Sublimación, serigrafía y bordado', fr: 'Sublimation, sérigraphie et broderie' } },
        { icon: '🔬', title: { en: 'Quality Control', zh: '质量检验', es: 'Control de Calidad', fr: 'Contrôle Qualité' }, desc: { en: 'Multi-point inspection at every stage', zh: '每道工序多点检验', es: 'Inspección multipunto en cada etapa', fr: `Inspection multi-points à chaque étape` } },
        { icon: '📦', title: { en: 'Packaging', zh: '包装', es: 'Embalaje', fr: 'Emballage' }, desc: { en: 'Professional packaging to your brand', zh: '按品牌要求专业包装', es: 'Embalaje profesional según su marca', fr: 'Emballage professionnel selon votre marque' } },
        { icon: '🚢', title: { en: 'Shipping', zh: '运输', es: 'Envío', fr: 'Expédition' }, desc: { en: 'Global logistics for on-time delivery', zh: '全球物流，准时交付', es: 'Logística global para entrega puntual', fr: 'Logistique mondiale pour livraison ponctuelle' } },
      ],
    },
  },
  {
    type: 'about_cta',
    title: 'Ready to Start?',
    config: {
      title: { en: 'Ready to Start Your Project?', zh: '准备开始您的项目？', es: '¿Listo para empezar?', fr: 'Prêt à démarrer votre projet ?' },
      description: {
        en: `Let's talk about your custom sportswear needs.`,
        zh: '让我们聊聊您的定制运动服装需求。',
        es: 'Hablemos de sus necesidades de ropa deportiva.',
        fr: 'Parlons de vos besoins en vêtements de sport.',
      },
      button_text: { en: 'Get a Quote', zh: '获取报价', es: 'Solicitar Cotización', fr: 'Obtenir un Devis' },
      button_url: '/contact',
    },
  },
]

async function findAboutPage(): Promise<any | null> {
  const tryList = async (status?: string) => {
    const res: any = await pageApi.list({ page: 1, pageSize: 100, status })
    const items = res?.items || res || []
    return items.find((p: any) => p.slug === 'about') || null
  }
  return (await tryList('published')) || (await tryList()) || null
}

async function loadPage() {
  loading.value = true
  try {
    const page = await findAboutPage()
    pageId.value = page?.id || ''
    if (pageId.value) {
      const detail: any = await pageApi.get(pageId.value)
      Object.assign(form, {
        title: detail.title || '',
        slug: detail.slug || 'about',
        type: detail.type || 'normal',
        template: detail.template || '',
        sort_order: detail.sort_order ?? 0,
        status: detail.status || 'draft',
      })
      modules.value = (detail.modules || []).map((m: any) => ({
        type: m.type,
        title: m.title || '',
        sort_order: m.sort_order ?? 0,
        is_visible: m.is_visible !== false,
        config: safeParseConfig(m.config),
      }))
      translations.value = translationsToRecord(detail.translations || [], PAGE_TRANS_FIELDS)
    }
  } finally {
    loading.value = false
  }
}

function safeParseConfig(raw: any): Record<string, any> {
  if (!raw) return {}
  if (typeof raw === 'object') return raw
  try { return JSON.parse(raw) } catch { return {} }
}

function addModule() {
  modules.value.push({ type: 'text', title: '', sort_order: modules.value.length + 1, is_visible: true, config: {} })
}
function removeModule(i: number) { modules.value.splice(i, 1); reindexModules() }
function moveModule(i: number, dir: number) {
  const j = i + dir
  if (j < 0 || j >= modules.value.length) return
  const arr = modules.value
  ;[arr[i], arr[j]] = [arr[j], arr[i]]
  reindexModules()
}
function reindexModules() { modules.value.forEach((m: any, idx: number) => { m.sort_order = idx + 1 }) }

function buildPayload() {
  reindexModules()
  return {
    title: form.title,
    slug: form.slug,
    type: form.type,
    template: form.template,
    sort_order: form.sort_order,
    modules: modules.value.map((m: any) => ({
      type: m.type,
      title: m.title || '',
      sort_order: m.sort_order ?? 0,
      is_visible: m.is_visible !== false,
      config: JSON.stringify(m.config || {}),
    })),
    translations: translationsToPayload(translations.value, PAGE_TRANS_FIELDS, TRANS_LANGS),
  }
}

async function handleSave() {
  if (!pageId.value) return
  if (!form.title.trim()) { ElMessage.warning('请输入页面标题'); return }
  saving.value = true
  try {
    await pageApi.update(pageId.value, buildPayload())
    ElMessage.success('保存成功，前台即时生效')
    loadPage()
  } finally { saving.value = false }
}

async function handlePublish() {
  if (!pageId.value) return
  saving.value = true
  try { await pageApi.publish(pageId.value); ElMessage.success('已发布'); loadPage() } finally { saving.value = false }
}
async function handleUnpublish() {
  if (!pageId.value) return
  saving.value = true
  try { await pageApi.unpublish(pageId.value); ElMessage.success('已下线'); loadPage() } finally { saving.value = false }
}

async function initPage() {
  initializing.value = true
  try {
    const created: any = await pageApi.create({
      title: 'About Us',
      slug: 'about',
      type: 'normal',
      status: 'published',
      template: '',
      sort_order: 0,
      modules: DEFAULT_ABOUT_MODULES.map((m: any, i: number) => ({
        type: m.type,
        title: m.title,
        sort_order: i + 1,
        is_visible: true,
        config: JSON.stringify(m.config),
      })),
      translations: [
        { language: 'zh', title: '关于我们', content: '' },
        { language: 'es', title: 'Sobre Nosotros', content: '' },
        { language: 'fr', title: 'À Propos', content: '' },
      ],
    })
    pageId.value = created?.id || ''
    ElMessage.success('关于我们页面已初始化')
    loadPage()
  } finally { initializing.value = false }
}

onMounted(() => { loadEnumDict(); loadPage() })
</script>

<style scoped>
.head { display: flex; align-items: flex-start; justify-content: space-between; gap: 16px; margin-bottom: 16px; }
.head-title { margin: 0 0 6px; font-size: 18px; font-weight: 700; color: #0d1b2a; }
.head-sub { margin: 0; font-size: 13px; color: #909399; line-height: 1.6; }
.head-actions { display: flex; gap: 8px; flex-shrink: 0; }
.module-toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.module-item { border: 1px solid #e4e7ed; border-radius: 6px; padding: 10px 12px; margin-bottom: 10px; background: #fafbfc; }
.module-head { display: flex; gap: 8px; align-items: center; margin-bottom: 8px; }
.module-index {
  display: inline-block; width: 20px; height: 20px; line-height: 20px; text-align: center;
  background: #409eff; color: #fff; border-radius: 50%; font-size: 12px; flex-shrink: 0;
}
.muted { color: #909399; font-size: 12px; }
@media (max-width: 960px) { .head { flex-direction: column; } .head-actions { flex-wrap: wrap; } }
</style>

