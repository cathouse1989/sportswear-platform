// ====================================================================
// 管理后台左侧导航菜单配置
//
// 设计原则：
//   - 数据驱动：菜单结构与排序通过此配置文件完全控制，调整顺序只需
//     改变 `sort_order` 字段，无需修改 Vue 模板。
//   - 每个项目都有 `sort_order` 字段，值越小越靠前。
//   - 顶层菜单与分组的子菜单各自按 `sort_order` 排序。
//   - 独立菜单项：当前用户拥有指定权限（any-of 语义）时显示。
//   - 分组菜单项：当 **任一** 子菜单可见时显示分组。
// ====================================================================

import type { Component } from 'vue'
import {
  DataAnalysis,
  Goods,
  Document,
  Picture,
  Message,
  User,
  Key,
  Setting,
  Monitor,
  Bell,
  Star,
  Promotion,
} from '@element-plus/icons-vue'

/** 权限类型：单一权限码或 any-of 数组，与 hasAnyPermission / can() 保持一致 */
export type MenuPermission = string | string[]

/** 菜单项定义 */
export interface MenuItem {
  /** 唯一键（用于 Vue :key） */
  key: string
  /** 显示文本 */
  label: string
  /** 路由路径 / 分组索引 */
  index: string
  /** 排序权重，值越小越靠前 */
  sort_order: number
  /** 图标组件 */
  icon?: Component
  /** 可见所需权限（any-of） */
  permission?: MenuPermission
  /** 子菜单（存在即为分组） */
  children?: MenuItem[]
}

/** 菜单配置 — 按 sort_order 决定显示顺序 */
export const MENU_CONFIG: MenuItem[] = [
  // ── 首页 ──
  {
    key: 'dashboard',
    label: '仪表盘',
    index: '/dashboard',
    sort_order: 10,
    icon: DataAnalysis,
    permission: ['dashboard:view', 'lead:view'],
  },

  // ── 业务管理 ──
  {
    key: 'products',
    label: '产品管理',
    index: 'products',
    sort_order: 20,
    icon: Goods,
    children: [
      { key: 'product-list', label: '产品列表', index: '/products', sort_order: 10, permission: 'product:view' },
      { key: 'product-categories', label: '分类管理', index: '/product-categories', sort_order: 20, permission: 'category:manage' },
      { key: 'product-series', label: '系列管理', index: '/product-series', sort_order: 30, permission: 'series:manage' },
      { key: 'product-fabrics', label: '面料管理', index: '/product-fabrics', sort_order: 40, permission: 'fabric:manage' },
    ],
  },
  {
    key: 'leads',
    label: '询盘管理',
    index: '/leads',
    sort_order: 30,
    icon: Message,
    permission: 'lead:view',
  },
  {
    key: 'quotes',
    label: '报价管理',
    index: '/quotes',
    sort_order: 40,
    icon: Star,
    permission: 'lead:view',
  },
  {
    key: 'notifications',
    label: '通知列表',
    index: '/notifications',
    sort_order: 50,
    icon: Bell,
    permission: 'lead:view',
  },
  {
    key: 'subscribers',
    label: '订阅管理',
    index: '/subscribers',
    sort_order: 55,
    icon: Promotion,
    permission: 'subscription:view',
  },

  // ── 内容管理 ──
  {
    key: 'content',
    label: '内容管理',
    index: 'cms',
    sort_order: 60,
    icon: Document,
    children: [
      { key: 'hero', label: '轮播图管理', index: '/hero', sort_order: 10, permission: 'page:view' },
      { key: 'navigations', label: '导航管理', index: '/navigations', sort_order: 20, permission: ['navigation:manage', 'page:view'] },
      { key: 'blogs', label: '博客管理', index: '/blogs', sort_order: 30, permission: ['blog:manage', 'page:view'] },
      { key: 'cases', label: '案例管理', index: '/cases', sort_order: 40, permission: ['case:manage', 'page:view'] },
      { key: 'faqs', label: 'FAQ 管理', index: '/faqs', sort_order: 50, permission: ['faq:manage', 'page:view'] },
      { key: 'factories', label: '工厂管理', index: '/factories', sort_order: 60, permission: ['factory:manage', 'page:view'] },
      { key: 'certifications', label: '认证管理', index: '/certifications', sort_order: 70, permission: ['certification:manage', 'page:view'] },
      { key: 'production-processes', label: '生产流程', index: '/production-processes', sort_order: 80, permission: ['production:manage', 'page:view'] },
      { key: 'self-medias', label: '自媒体管理', index: '/self-medias', sort_order: 90, permission: ['selfmedia:manage', 'page:view'] },
    ],
  },

  // ── 媒体管理 ──
  {
    key: 'media',
    label: '媒体管理',
    index: '/media',
    sort_order: 70,
    icon: Picture,
    permission: ['media:manage', 'media:upload'],
  },

  // ── 用户与权限 ──
  {
    key: 'user-management',
    label: '用户与权限',
    index: 'user-management',
    sort_order: 80,
    icon: User,
    children: [
      { key: 'users', label: '用户管理', index: '/users', sort_order: 10, permission: 'user:view' },
      { key: 'roles', label: '角色权限', index: '/roles', sort_order: 20, permission: 'role:manage' },
    ],
  },

  // ── 系统配置 ──
  {
    key: 'system',
    label: '系统配置',
    index: 'system',
    sort_order: 100,
    icon: Setting,
    children: [
      { key: 'theme', label: '主题配置', index: '/theme', sort_order: 10, permission: 'setting:manage' },
      { key: 'i18n', label: '词条管理', index: '/i18n', sort_order: 20, permission: 'language:manage' },
      { key: 'seo', label: 'SEO 管理', index: '/seo', sort_order: 25, permission: 'seo:manage' },
      { key: 'translation-coverage', label: '翻译完成度', index: '/translation-coverage', sort_order: 26, permission: ['language:manage', 'page:view'] },
      { key: 'storage-sources', label: '存储源配置', index: '/storage-sources', sort_order: 30, permission: 'setting:manage' },
      { key: 'portal-cache', label: '门户缓存', index: '/portal-cache', sort_order: 40, permission: 'setting:manage' },
      { key: 'trash', label: '回收站', index: '/trash', sort_order: 50, permission: 'setting:manage' },
      { key: 'operation-logs', label: '操作日志', index: '/operation-logs', sort_order: 60, permission: 'setting:manage' },
    ],
  },

  // ── 洞析 ──
  {
    key: 'insights',
    label: '洞析',
    index: 'insights',
    sort_order: 110,
    icon: Monitor,
    children: [
      { key: 'portal-preview', label: '门户预览', index: '/portal-preview', sort_order: 10, permission: ['dashboard:view', 'lead:view'] },
      { key: 'analytics', label: '数据分析', index: '/analytics', sort_order: 20, permission: ['dashboard:view', 'lead:view'] },
    ],
  },
]
