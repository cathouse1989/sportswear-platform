<template>
  <el-card shadow="never">
    <div class="toolbar">
      <div class="page-title">
        <h3>门户缓存</h3>
        <p class="page-desc">
          控制门户公开内容是否走 Redis 缓存。开关开启时公开接口读缓存，内容变更会自动失效缓存；关闭时始终查数据库，且不自动更新缓存，仅可手动刷新/发布。
        </p>
      </div>
      <div class="spacer" />
      <el-button @click="loadStatus" :loading="loading">刷新状态</el-button>
    </div>

    <el-descriptions :column="1" border v-loading="loading" class="status-box">
      <el-descriptions-item label="Redis 可用性">
        <el-tag :type="status.redis_ok ? 'success' : 'danger'" size="small">
          {{ status.redis_ok ? '可用' : '不可用' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="缓存开关">
        <div class="switch-row">
          <el-switch
            v-model="enabled"
            :loading="toggling"
            :disabled="!status.redis_ok && !enabled"
            active-text="走缓存"
            inactive-text="走数据库"
            @change="onToggle"
          />
          <span class="hint">{{ enabled ? '开启：读缓存，变更自动刷新缓存' : '关闭：直查库，仅可手动刷新缓存' }}</span>
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="当前读路径">
        <el-tag :type="status.using_cache ? 'success' : 'info'" size="small">
          {{ status.using_cache ? '缓存' : '数据库' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="最近刷新/发布时间">
        {{ status.last_refresh_at || '—' }}
      </el-descriptions-item>
    </el-descriptions>

    <div class="actions">
      <el-button type="primary" :loading="refreshing" :disabled="!status.redis_ok" @click="handleRefresh">
        手动刷新缓存
      </el-button>
      <el-button type="success" :loading="publishing" :disabled="!status.redis_ok" @click="handlePublish">
        发布更新门户
      </el-button>
    </div>
    <el-alert
      class="tip"
      type="info"
      :closable="false"
      show-icon
      title="说明"
      description="门户预览（preview=1）始终直查数据库，不受缓存影响。编辑内容后可在预览确认，再点「发布更新门户」让线上缓存与数据库一致。开关关闭时发布/刷新仍会写入 Redis 预热，但公开读路径仍走数据库。"
    />
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { portalCacheApi } from '@/api'

const loading = ref(false)
const toggling = ref(false)
const refreshing = ref(false)
const publishing = ref(false)
const enabled = ref(true)
const status = reactive({
  redis_ok: false,
  enabled: true,
  using_cache: false,
  last_refresh_at: null as string | null,
})

function applyStatus(data: any) {
  status.redis_ok = !!data?.redis_ok
  status.enabled = !!data?.enabled
  status.using_cache = !!data?.using_cache
  status.last_refresh_at = data?.last_refresh_at || null
  enabled.value = status.enabled
}

async function loadStatus() {
  loading.value = true
  try {
    const data = await portalCacheApi.status()
    applyStatus(data)
  } catch {
    /* handled */
  } finally {
    loading.value = false
  }
}

async function onToggle(val: string | number | boolean) {
  const next = !!val
  toggling.value = true
  try {
    const data = await portalCacheApi.setEnabled(next)
    applyStatus(data)
    ElMessage.success(next ? '已开启缓存，并自动刷新门户缓存' : '已关闭缓存，公开接口将直查数据库')
  } catch {
    enabled.value = !next
  } finally {
    toggling.value = false
  }
}

async function handleRefresh() {
  refreshing.value = true
  try {
    const data = await portalCacheApi.refresh()
    applyStatus(data)
    ElMessage.success('缓存已刷新')
  } catch {
    /* handled */
  } finally {
    refreshing.value = false
  }
}

async function handlePublish() {
  await ElMessageBox.confirm(
    '将根据当前数据库重建门户线上缓存，使预览中的最新内容对访客生效。是否继续？',
    '发布更新门户',
    { type: 'warning' },
  )
  publishing.value = true
  try {
    const data = await portalCacheApi.publish()
    applyStatus(data)
    ElMessage.success('门户缓存已发布更新')
  } catch {
    /* handled */
  } finally {
    publishing.value = false
  }
}

onMounted(loadStatus)
</script>

<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: flex-start; }
.page-title h3 { margin: 0; font-size: 16px; font-weight: 600; color: #0d1b2a; }
.page-desc { margin: 4px 0 0; font-size: 12px; color: #8b7d6b; max-width: 720px; line-height: 1.5; }
.spacer { flex: 1; }
.status-box { max-width: 720px; }
.switch-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.hint { font-size: 12px; color: #8b7d6b; }
.actions { margin-top: 20px; display: flex; gap: 12px; }
.tip { margin-top: 16px; max-width: 720px; }
</style>
