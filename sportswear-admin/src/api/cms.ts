import { http } from './client'
import type { Blog, Case, FAQ, Page } from '@/types'

export const cmsApi = {
  // 页面
  pages: {
    list: (params?: any) => http.get<{ items: Page[]; total: number }>('/admin/pages', params),
    get: (id: string) => http.get<Page>(`/admin/pages/${id}`),
    create: (data: any) => http.post<Page>('/admin/pages', data),
    update: (id: string, data: any) => http.put<Page>(`/admin/pages/${id}`, data),
    delete: (id: string) => http.delete(`/admin/pages/${id}`),
    publish: (id: string) => http.post(`/admin/pages/${id}/publish`),
    unpublish: (id: string) => http.post(`/admin/pages/${id}/unpublish`),
  },
  // 博客
  blogs: {
    list: (params?: any) => http.get<{ items: Blog[]; total: number }>('/admin/blogs', params),
    get: (id: string) => http.get<Blog>(`/admin/blogs/${id}`),
    create: (data: any) => http.post<Blog>('/admin/blogs', data),
    update: (id: string, data: any) => http.put<Blog>(`/admin/blogs/${id}`, data),
    delete: (id: string) => http.delete(`/admin/blogs/${id}`),
    publish: (id: string) => http.post(`/admin/blogs/${id}/publish`),
  },
  // 案例
  cases: {
    list: (params?: any) => http.get<{ items: Case[]; total: number }>('/admin/cases', params),
    get: (id: string) => http.get<Case>(`/admin/cases/${id}`),
    create: (data: any) => http.post<Case>('/admin/cases', data),
    update: (id: string, data: any) => http.put<Case>(`/admin/cases/${id}`, data),
    delete: (id: string) => http.delete(`/admin/cases/${id}`),
  },
  // FAQ
  faqs: {
    list: (params?: any) => http.get<{ items: FAQ[]; total: number }>('/admin/faqs', params),
    create: (data: any) => http.post<FAQ>('/admin/faqs', data),
    update: (id: string, data: any) => http.put<FAQ>(`/admin/faqs/${id}`, data),
    delete: (id: string) => http.delete(`/admin/faqs/${id}`),
  },
}