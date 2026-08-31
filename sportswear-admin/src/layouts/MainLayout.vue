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
        <template
          v-for="item in visibleMenu"
          :key="item.key"
        >
          <el-menu-item
            v-if="!item.children"
            :index="item.index"
          >
            <el-icon v-if="item.icon">
              <component :is="item.icon" />
            </el-icon>
            <span>{{ item.label }}</span>
          </el-menu-item>
          <el-sub-menu
            v-else
            :index="item.index"
          >
            <template #title>
              <el-icon v-if="item.icon">
                <component :is="item.icon" />
              </el-icon>
              <span>{{ item.label }}</span>
            </template>
            <el-menu-item
              v-for="child in item.children!"
              :key="child.key"
              :index="child.index"
            >
              {{ child.label }}
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <span class="breadcrumb">{{ currentTitle }}</span>
        </div>
        <div class="header-right">
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notify-badge" @click="goToNotifications">
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
import { hasAnyPermission } from '@/utils/permissions'
import { notificationApi } from '@/api'
import { http } from '@/api/client'
// 图标按需引入（替换 main.ts 的全局注册，减小首屏包体积）
import { Bell } from '@element-plus/icons-vue'
import { MENU_CONFIG, type MenuItem } from '@/config/menu'

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

// 菜单权限过滤：与路由 meta.permission（v-permission 指令）同一套 any-of 语义
function can(required: string | string[]) {
  return hasAnyPermission(authStore, required)
}
// 可见菜单：过滤无权限项 + 按 sort_order 排序
// - 独立菜单项：拥有指定权限时显示
// - 分组菜单项：任一子菜单可见时显示（分组本身不单设权限）
const visibleMenu = computed(() => {
  return MENU_CONFIG
    .filter((item: MenuItem) => {
      if (item.children) {
        return item.children.some(c => can(c.permission as string | string[]))
      }
      return can(item.permission as string | string[])
    })
    .sort((a, b) => a.sort_order - b.sort_order)
    .map((item: MenuItem) => {
      if (item.children) {
        return {
          ...item,
          children: item.children
            .filter(c => can(c.permission as string | string[]))
            .sort((a, b) => a.sort_order - b.sort_order),
        }
      }
      return item
    })
})

onMounted(async () => {
  try {
    const result = await notificationApi.unreadCount()
    unreadCount.value = result.count
  } catch {
    // 忽略未读计数加载失败
  }
  // 加载 Logo 配置（含 Favicon）
  try {
    const themeItems = await http.get<any[]>('/admin/theme')
    const logoItem = themeItems.find((i: any) => i.key === 'logo_url')
    if (logoItem?.value) logoUrl.value = logoItem.value
    const altItem = themeItems.find((i: any) => i.key === 'logo_alt')
    if (altItem?.value) logoAlt.value = altItem.value
    
    // 动态更新浏览器 Tab 图标（Favicon）
    const faviconItem = themeItems.find((i: any) => i.key === 'favicon_url')
    if (faviconItem?.value) {
      const faviconUrl = faviconItem.value
      const existingLink = document.querySelector('link[rel="icon"]')
      if (existingLink) {
        existingLink.setAttribute('href', faviconUrl)
      } else {
        const link = document.createElement('link')
        link.rel = 'icon'
        link.href = faviconUrl
        document.head.appendChild(link)
      }
    }
  } catch {
    // 忽略配置加载失败
  }
})

async function handleCommand(command: string) {
  if (command === 'logout') {
    await ElMessageBox.confirm('确定退出登录吗？', '提示', { type: 'warning' })
    authStore.logout()
    router.push('/login')
  }
}

function goToNotifications() {
  router.push('/notifications')
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
  display: inline-flex;
  align-items: center;
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