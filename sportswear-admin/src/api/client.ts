import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import type { ApiResponse } from '@/types'

const baseURL = '/api/v1'

const instance: AxiosInstance = axios.create({
  baseURL,
  timeout: 30000,
})

// 请求拦截器：注入 token
instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：统一处理错误
instance.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
      localStorage.removeItem('admin_permissions')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    } else {
      const msg = error.response?.data?.message || error.message || '请求失败'
      ElMessage.error(msg)
    }
    return Promise.reject(error)
  },
)

// 通用请求方法
export async function request<T = any>(config: AxiosRequestConfig): Promise<T> {
  const response = await instance.request<ApiResponse<T>>(config)
  const data = response.data
  if (!data.success) {
    throw new Error(data.message || '请求失败')
  }
  return data.data
}

export const http = {
  get: <T = any>(url: string, params?: any) => request<T>({ method: 'get', url, params }),
  post: <T = any>(url: string, data?: any) => request<T>({ method: 'post', url, data }),
  put: <T = any>(url: string, data?: any) => request<T>({ method: 'put', url, data }),
  delete: <T = any>(url: string) => request<T>({ method: 'delete', url }),
}

export default instance