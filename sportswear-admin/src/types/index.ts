// 通用响应结构
export interface ApiResponse<T = any> {
  success: boolean
  code?: string
  message?: string
  data: T
  requestId?: string
}

export interface PageResult<T = any> {
  items: T[]
  page: number
  pageSize: number
  total: number
  totalPages: number
}

export interface User {
  id: string
  email: string
  name: string
  avatar?: string
  phone?: string
  status: string
  last_login_at?: string
  roles?: Role[]
  created_at: string
}

export interface Role {
  id: string
  name: string
  code: string
  description?: string
  is_active: boolean
  permissions?: Permission[]
}

export interface Permission {
  id: string
  name: string
  code: string
  module: string
  description?: string
}

export interface Product {
  id: string
  sku: string
  slug: string
  name?: string
  type: string
  gender: string
  status: string
  is_featured: boolean
  is_new: boolean
  sort_order: number
  cover_image: string
  brief: string
  description: string
  features?: string
  usage?: string
  material?: string
  composition?: string
  weight?: string
  elasticity?: string
  fit?: string
  support_level?: string
  season?: string
  size_range?: string
  sample_moq: number
  production_moq: number
  color_moq: number
  size_moq: number
  category_id?: string
  category?: Category
  translations?: ProductTranslation[]
  images?: ProductImage[]
  videos?: ProductVideo[]
  specs?: ProductSpec[]
  customizations?: ProductCustomization[]
  series?: Series[]
  fabrics?: Fabric[]
  seo?: SEO
  created_at: string
}

export interface ProductTranslation {
  id: string
  product_id: string
  language: string
  name: string
  brief: string
  description: string
  features?: string
  usage?: string
  /** 语言维度展示排序权重（0 = 跟随全局排序） */
  sort_order?: number
  status: string
}

export interface ProductImage {
  id: string
  product_id: string
  type: string
  url: string
  thumbnail?: string
  alt?: string
  sort_order: number
}

export interface ProductVideo {
  id: string
  product_id: string
  type: string
  url: string
  cover?: string
  title?: string
  sort_order: number
}

export interface ProductSpec {
  id: string
  product_id: string
  name: string
  value: string
  sort_order: number
}

export interface ProductCustomization {
  id: string
  product_id: string
  type: string
  is_enabled: boolean
  note?: string
}

export interface Category {
  id: string
  parent_id?: string
  name: string
  slug: string
  sort_order: number
  is_active: boolean
  image?: string
  children?: Category[]
}

export interface Series {
  id: string
  name: string
  slug: string
  status: string
  is_active: boolean
  sort_order: number
}

export interface Fabric {
  id: string
  name: string
  code: string
  status: string
  composition: string
  weight?: string
  elasticity?: string
  breathability?: string
  moisture_wicking?: string
  softness?: string
  compression?: string
  uv_protection?: string
  eco_friendly?: string
  description: string
  is_active: boolean
}

export interface SEO {
  id: string
  entity_type: string
  entity_id: string
  language?: string
  title?: string
  description?: string
  keywords?: string
  canonical?: string
  robots?: string
  og_title?: string
  og_description?: string
  og_image?: string
  schema_data?: string
}

export interface Page {
  id: string
  title: string
  slug: string
  type: string
  status: string
  template?: string
  sort_order: number
  modules?: PageModule[]
  translations?: PageTranslation[]
  seo?: SEO
  published_at?: string
  created_at?: string
}

export interface PageModule {
  id: string
  page_id: string
  type: string
  title: string
  sort_order: number
  is_visible: boolean
  config: string
}

export interface PageTranslation {
  id: string
  page_id: string
  language: string
  title: string
  content: string
  status: string
}

// 首页轮播图（hero banner）配置项
export interface HeroSlide {
  image: string
  title: string
  subtitle: string
  button_text: string
  button_url: string
}

// 首页轮播展示参数（与 slides 一并保存到 banner.config.settings）
export interface HeroSettings {
  autoplay?: boolean
  interval_ms?: number
  transition?: 'fade' | 'slide' | string
  show_dots?: boolean
  show_arrows?: boolean
  pause_on_hover?: boolean
}


export interface Blog {
  id: string
  title: string
  slug: string
  category: string
  tags: string
  author?: string
  cover_image: string
  content: string
  status: string
  translations?: BlogTranslation[]
  seo?: SEO
  published_at?: string
  created_at?: string
}

export interface BlogTranslation {
  id: string
  blog_id: string
  language: string
  title: string
  content: string
  status: string
}

export interface Case {
  id: string
  title: string
  slug: string
  client_industry: string
  project_type?: string
  products?: string
  client_need?: string
  problem?: string
  solution?: string
  process?: string
  result?: string
  cover_image?: string
  status: string
  translations?: CaseTranslation[]
  seo?: SEO
  published_at?: string
  created_at?: string
}

export interface CaseTranslation {
  id: string
  case_id: string
  language: string
  title: string
  client_need?: string
  problem?: string
  solution?: string
  process?: string
  result?: string
  status: string
}

export interface FAQ {
  id: string
  question: string
  answer: string
  category: string
  language: string
  sort_order: number
  is_active: boolean
}

export interface Factory {
  id: string
  name: string
  status: string
  location: string
  area?: string
  employees: number
  production_lines?: number
  equipment?: string
  monthly_capacity?: string
  annual_capacity?: string
  warehouse?: string
  quality_management?: string
  production_capability?: string
  description?: string
  image?: string
  is_active: boolean
}

export interface Certification {
  id: string
  name: string
  status: string
  code: string
  issue_date?: string
  expiry_date?: string
  image?: string
  pdf?: string
  description?: string
  is_active: boolean
}

export interface ProductionProcess {
  id: string
  name: string
  status: string
  description?: string
  image?: string
  video?: string
  sort_order: number
  is_active: boolean
}

export interface Navigation {
  id: string
  name: string
  type: string
  url: string
  target?: string
  sort_order: number
  is_visible: boolean
  parent_id?: string
  page_id?: string
  page?: Page
  page_status?: string  // 关联页面的状态（用于显示页面发布状态）
  children?: Navigation[]
}

export interface Media {
  id: string
  original_name: string
  file_name: string
  file_type: string
  file_size: number
  url: string
  path: string
  type: string
  source: string
  provider: string
  external_id?: string
  embed_url?: string
  category: string
  is_public: boolean
  title?: string
  alt?: string
  description?: string
  thumbnail?: string
  width?: number
  height?: number
  created_at: string
}

export interface StorageSource {
  id: string
  name: string
  code: string
  type: string
  protocol: string
  host: string
  port: string
  base_path: string
  bucket: string
  is_default: boolean
  is_active: boolean
}

export interface Lead {
  id: string
  name: string
  company: string
  email: string
  phone: string
  country: string
  project_type: string
  product_category: string
  quantity: number
  budget: string
  message: string
  attachments?: string
  source: string
  status: string
  score: number
  score_level: string
  assigned_to?: string
  next_follow_up?: string
  follow_ups?: LeadFollowUp[]
  created_at: string
}

export interface LeadFollowUp {
  id: string
  lead_id: string
  user_id: string
  method: string
  content: string
  next_follow_up?: string
  created_at: string
}

export interface Quote {
  id: string
  lead_id: string
  quote_number: string
  quantity: number
  unit_price: number
  total_price: number
  currency: string
  moq: number
  lead_time: string
  status: string
  valid_until?: string
}

export interface Notification {
  id: string
  type: string
  title: string
  content: string
  is_read: boolean
  entity_type: string
  created_at: string
}

export interface OperationLog {
  id: string
  user_id: string
  operation: string
  module: string
  entity_type: string
  description: string
  ip: string
  created_at: string
}

export interface DashboardStats {
  today_leads: number
  week_leads: number
  month_leads: number
  high_value_leads: number
  won_leads: number
  hot_products: Array<{ name: string; count: number }>
  hot_countries: Array<{ name: string; count: number }>
}

export interface Currency {
  id: string
  code: string
  name: string
  symbol: string
  rate: number
  is_default: boolean
  is_active: boolean
}

export interface I18nEntry {
  id: string
  key: string
  language: string
  value: string
  module: string
}