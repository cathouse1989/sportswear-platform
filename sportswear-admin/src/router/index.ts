import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { hasAnyPermission } from '@/utils/permissions'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '登录' },
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/dashboard/index.vue'), meta: { title: '仪表盘', permission: 'dashboard:view' } },
      { path: 'products', name: 'Products', component: () => import('@/views/product/index.vue'), meta: { title: '产品管理', permission: 'product:view' } },
      { path: 'leads', name: 'Leads', component: () => import('@/views/lead/index.vue'), meta: { title: '询盘管理', permission: 'lead:view' } },
      { path: 'users', name: 'Users', component: () => import('@/views/system/users.vue'), meta: { title: '用户管理', permission: 'user:view' } },
      { path: 'hero', name: 'Hero', component: () => import('@/views/cms/hero.vue'), meta: { title: '轮播图管理', permission: 'page:view' } },
      // { path: 'pages', name: 'Pages', component: () => import('@/views/cms/pages.vue'), meta: { title: '页面管理', permission: 'page:view' } },
      { path: 'navigations', name: 'Navigations', component: () => import('@/views/cms/navigations.vue'), meta: { title: '导航管理', permission: ['navigation:manage', 'page:view'] } },
      { path: 'blogs', name: 'Blogs', component: () => import('@/views/cms/blogs.vue'), meta: { title: '博客管理', permission: ['blog:manage', 'page:view'] } },
      { path: 'cases', name: 'Cases', component: () => import('@/views/cms/cases.vue'), meta: { title: '案例管理', permission: ['case:manage', 'page:view'] } },
      { path: 'faqs', name: 'FAQs', component: () => import('@/views/cms/faqs.vue'), meta: { title: 'FAQ 管理', permission: ['faq:manage', 'page:view'] } },
      { path: 'factories', name: 'Factories', component: () => import('@/views/cms/factories.vue'), meta: { title: '工厂管理', permission: ['factory:manage', 'page:view'] } },
      { path: 'certifications', name: 'Certifications', component: () => import('@/views/cms/certifications.vue'), meta: { title: '认证管理', permission: ['certification:manage', 'page:view'] } },
      { path: 'media', name: 'Media', component: () => import('@/views/media/index.vue'), meta: { title: '媒体管理', permission: ['media:manage', 'media:upload'] } },
      { path: 'roles', name: 'Roles', component: () => import('@/views/system/roles.vue'), meta: { title: '角色权限', permission: 'role:manage' } },
      { path: 'theme', name: 'Theme', component: () => import('@/views/system/theme.vue'), meta: { title: '主题配置', permission: 'setting:manage' } },
      { path: 'i18n', name: 'I18n', component: () => import('@/views/system/i18n.vue'), meta: { title: '词条管理', permission: 'language:manage' } },
      { path: 'portal-cache', name: 'PortalCache', component: () => import('@/views/system/portal-cache.vue'), meta: { title: '门户缓存', permission: 'setting:manage' } },
      { path: 'trash', name: 'Trash', component: () => import('@/views/system/trash.vue'), meta: { title: '回收站', permission: 'setting:manage' } },
      { path: 'storage-sources', name: 'StorageSources', component: () => import('@/views/system/storage-sources.vue'), meta: { title: '存储源配置', permission: 'setting:manage' } },
      { path: 'quotes', name: 'Quotes', component: () => import('@/views/system/quotes.vue'), meta: { title: '报价管理', permission: 'lead:view' } },
      { path: 'notifications', name: 'Notifications', component: () => import('@/views/system/notifications.vue'), meta: { title: '通知列表', permission: 'lead:view' } },
      { path: 'portal-preview', name: 'PortalPreview', component: () => import('@/views/portal/preview.vue'), meta: { title: '门户预览', permission: 'dashboard:view' } },
      { path: 'analytics', name: 'Analytics', component: () => import('@/views/analytics/index.vue'), meta: { title: '数据分析', permission: ['dashboard:view', 'lead:view'] } },
      { path: 'operation-logs', name: 'OperationLogs', component: () => import('@/views/system/operation-logs.vue'), meta: { title: '操作日志', permission: 'setting:manage' } },
    ],
  },
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach(async (to, _from, next) => {
  document.title = `${to.meta.title || ''} - OEM/ODM 管理后台`
  const authStore = useAuthStore()
  if (to.path === '/login') { next(); return }
  if (!authStore.token) { next('/login'); return }
  if (!authStore.user) { try { await authStore.fetchProfile() } catch { authStore.logout(); next('/login'); return } }
  const permission = to.meta.permission as string | string[] | undefined
  if (permission && !hasAnyPermission(authStore, permission)) {
    next('/dashboard')
    return
  }
  next()
})

export default router