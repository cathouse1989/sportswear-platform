<template>
  <div>
    <!-- Logo 配置区块 -->
    <el-card shadow="never" class="mb-4">
      <template #header><span class="font-semibold">商标 Logo 配置</span></template>
      <el-form label-width="140px" label-position="left">
        <el-form-item label="Logo 图片 URL">
          <MediaPicker v-model="logoUrl" class="w-full" />
        </el-form-item>
        <el-form-item label="Logo 预览">
          <div class="flex items-center gap-4">
            <img v-if="logoUrl" :src="logoUrl" alt="Logo" class="h-12 w-auto rounded-lg shadow-sm" />
            <span v-else class="text-gray-400 text-sm">未配置，将使用默认 SVG Logo</span>
          </div>
        </el-form-item>
        <el-form-item label="系统名称" title="用于管理后台侧栏显示">
          <el-input v-model="adminSystemName" placeholder="管理后台侧栏 Logo 旁显示的名称，例如：OEM/ODM 管理系统" size="large" />
          <div class="mt-1 text-xs text-gray-400">仅管理后台侧栏显示，与门户品牌名称独立</div>
        </el-form-item>
        <el-form-item label="品牌名称">
          <el-input v-model="brandName" placeholder="例如：SPORTSWEAR（门户头部/页脚显示）" size="large" />
        </el-form-item>
        <el-form-item label="品牌副标题">
          <el-input v-model="brandSubtitle" placeholder="例如：Premium Mfg." size="large" />
        </el-form-item>
        <el-form-item label="Logo 替代文本">
          <el-input v-model="logoAlt" placeholder="例如：Sportswear Manufacturer" size="large" />
        </el-form-item>
        <el-form-item label="Favicon URL">
          <MediaPicker v-model="faviconUrl" class="w-full" />
          <div class="mt-2 text-xs text-gray-400">未配置时使用默认 favicon，建议使用 32x32 的 SVG 或 PNG</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" @click="handleSaveLogo" :loading="savingLogo">保存 Logo 配置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 原有主题配置表 -->
    <el-card shadow="never">
      <template #header><span class="font-semibold">其他主题配置项</span></template>
      <el-table :data="items" v-loading="loading" stripe max-height="400">
        <el-table-column prop="key" label="配置项" width="200" />
        <el-table-column prop="name" label="名称" width="160" />
        <el-table-column prop="group" label="分组" width="100" />
        <el-table-column label="值" min-width="200">
          <template #default="{ row }">
            <el-input v-if="['primary_color','secondary_color','accent_color'].includes(row.key)" v-model="row.value" size="small">
              <template #append><el-color-picker v-model="row.value" size="small" /></template>
            </el-input>
            <el-input-number
              v-else-if="isNumberKey(row.key)"
              :model-value="Number(row.value)"
              :min="0"
              size="small"
              controls-position="right"
              style="width: 100%"
              @update:model-value="(v) => (row.value = String(v))"
            />
            <el-input v-else v-model="row.value" size="small" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleSave(row)">保存</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { themeApi } from '@/api'
import MediaPicker from '@/components/media/MediaPicker.vue'

const items = ref<any[]>([]); const loading = ref(false)
const logoUrl = ref('')
const faviconUrl = ref('')
const brandName = ref('')
const brandSubtitle = ref('')
const adminSystemName = ref('')
const logoAlt = ref('')
const savingLogo = ref(false)

// 数值型配置项：以数字输入框渲染（含「自媒体最多展示数量」）
const NUMBER_KEYS = [
  'header_height', 'footer_columns', 'section_spacing', 'container_width',
  'button_radius', 'card_radius',
  'products_page_size', 'blogs_page_size', 'cases_page_size', 'faqs_page_size', 'admin_page_size',
  'hero_interval_ms', 'self_media_max_display',
]
const isNumberKey = (key: string) => NUMBER_KEYS.includes(key)

async function loadData() {
  loading.value = true
  try {
    const allItems = await themeApi.adminList()
    items.value = allItems
    logoUrl.value = allItems.find((i: any) => i.key === 'logo_url')?.value || ''
    faviconUrl.value = allItems.find((i: any) => i.key === 'favicon_url')?.value || ''
    brandName.value = allItems.find((i: any) => i.key === 'brand_name')?.value || ''
    brandSubtitle.value = allItems.find((i: any) => i.key === 'brand_subtitle')?.value || ''
    adminSystemName.value = allItems.find((i: any) => i.key === 'admin_system_name')?.value || ''
    logoAlt.value = allItems.find((i: any) => i.key === 'logo_alt')?.value || ''
  } finally { loading.value = false }
}

async function handleSaveLogo() {
  savingLogo.value = true
  try {
    const logoKeys: Record<string, string> = {
      logo_url: logoUrl.value,
      favicon_url: faviconUrl.value,
      brand_name: brandName.value,
      brand_subtitle: brandSubtitle.value,
      admin_system_name: adminSystemName.value,
      logo_alt: logoAlt.value,
    }
    await Promise.all(
      Object.entries(logoKeys).map(([key, value]) =>
        themeApi.update(key, value)
      )
    )
    ElMessage.success('Logo 配置已保存')
  } catch {
    ElMessage.error('保存失败')
  } finally { savingLogo.value = false }
}

async function handleSave(row: any) {
  await themeApi.update(row.key, row.value)
  ElMessage.success('已更新')
}

onMounted(loadData)
</script>