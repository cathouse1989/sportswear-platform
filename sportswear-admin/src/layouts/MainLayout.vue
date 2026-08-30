<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <template v-if="logoUrl">
          <img :src="logoUrl" :alt="logoAlt || 'Logo'" class="logo-img" />
        </template>
        <template v-else>
          <h2>OEM/ODM 后台</h2>
        </template>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#001529"
        text-color="#ffffffa6"
        active-text-color="#ffffff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/products">
          <el-icon><Goods /></el-icon>
          <span>产品管理</span>
        </el-menu-item>
        <el-sub-menu index="cms">
          <template #title><el-icon><Document /></el-icon><span>内容管理</span></template>
          <el-menu-item index="/pages">页面管理</el-menu-item>
          <el-menu-item index="/hero">轮播图管理</el-menu-item>
          <el-menu-item index="/blogs">博客管理</el-menu-item>
          <el-menu-item index="/cases">案例管理</el-menu-item>
          <el-menu-item index="/faqs">FAQ 管理</el-menu-item>
          <el-menu-item index="/factories">工厂管理</el-menu-item>
          <el-menu-item index="/certifications">认证管理</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/media">
          <el-icon><Picture /></el-icon>
          <span>媒体管理</span>
        </el-menu-item>
        <el-menu-item index="/leads">
          <el-icon><Message /></el-icon>
          <span>询盘管理</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/roles">
          <el-icon><Key /></el-icon>
          <span>角色权限</span>
        </el-menu-item>
        <el-sub-menu index="system">
          <template #title><el-icon><Setting /></el-icon><span>系统配置</span></template>
          <el-menu-item index="/theme">主题配置</el-menu-item>
          <el-menu-item index="/i18n">词条管理</el-menu-item>
          <el-menu-item index="/storage-sources">存储源配置</el-menu-item>
          <el-menu-item index="/quotes">报价管理</el-menu-item>
          <el-menu-item index="/notifications">通知列表</el-menu-item>
          <el-menu-item index="/portal-cache">门户缓存</el-menu-item>
          <el-menu-item index="/trash">回收站</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/portal-preview">
          <el-icon><Monitor /></el-icon>
          <span>门户预览</span>
        </el-menu-item>
        <el-menu-item index="/analytics">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据分析</span>
        </el-menu-item>
        <el-menu-item index="/operation-logs">
          <el-icon><Document /></el-icon>
          <span>操作日志</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <span class="breadcrumb">{{ currentTitle }}</span>
        </div>
        <div class="header-right">
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notify-badge">
            <el-icon :size="20"><Bell /></el-icon>
          </el-badge>
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="30">{{ userInitial }}</el-avatar>
              <span class="username">{{ userName }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { notificationApi } from '@/api'
import { http } from '@/api/client'
// 图标按需引入（替换 main.ts 的全局注册，减小首屏包体积）
import { DataAnalysis, Goods, Document, Picture, Message, User, Key, Setting, Monitor, Bell } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const unreadCount = ref(0)
const logoUrl = ref('')
const logoAlt = ref('')

const activeMenu = computed(() => route.path)
const currentTitle = computed(() => (route.meta.title as string) || '')
const userName = computed(() => authStore.user?.name || authStore.user?.email || '')
const userInitial = computed(() => userName.value.charAt(0).toUpperCase())

onMounted(async () => {
  try {
    const result = await notificationApi.unreadCount()
    unreadCount.value = result.count
  } catch {
    // 忽略未读计数加载失败
  }
  // 加载 Logo 配置
  try {
    const themeItems = await http.get<any[]>('/admin/theme')
    const item = themeItems.find((i: any) => i.key === 'logo_url')
    if (item?.value) logoUrl.value = item.value
    const altItem = themeItems.find((i: any) => i.key === 'logo_alt')
    if (altItem?.value) logoAlt.value = altItem.value
  } catch {}
})

async function handleCommand(command: string) {
  if (command === 'logout') {
    await ElMessageBox.confirm('确定退出登录吗？', '提示', { type: 'warning' })
    authStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}
.sidebar {
  background-color: #001529;
}
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  border-bottom: 1px solid #ffffff14;
  padding: 0 16px;
}
.logo h2 {
  font-size: 16px;
  margin: 0;
  white-space: nowrap;
}
.logo-img {
  max-height: 36px;
  max-width: 180px;
  object-fit: contain;
}
.sidebar :deep(.el-menu) {
  border-right: none;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e5e7eb;
  background: #fff;
}
.header-left .breadcrumb {
  font-size: 16px;
  font-weight: 500;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}
.notify-badge {
  cursor: pointer;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}
.username {
  font-size: 14px;
}
.main-content {
  background: #f5f6f8;
  padding: 20px;
}
</style>