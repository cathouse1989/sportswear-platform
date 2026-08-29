import { describe, it, expect, vi, beforeEach } from 'vitest'

// 回归用例（缺陷 B-11）：回收站 restore/purge 必须携带 entity_type（后端 binding:"required"）
vi.mock('@/api/client', () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import { http } from '@/api/client'
import { trashApi } from '@/api'

beforeEach(() => {
  vi.clearAllMocks()
})

describe('trashApi 契约（B-11 回归）', () => {
  it('restore 以 body 携带 entity_type', () => {
    trashApi.restore('uuid-1', 'product')
    expect(http.post).toHaveBeenCalledWith('/admin/trash/uuid-1/restore', { entity_type: 'product' })
  })

  it('purge 以 body 携带 entity_type', () => {
    trashApi.purge('uuid-2', 'factory')
    expect(http.post).toHaveBeenCalledWith('/admin/trash/uuid-2/purge', { entity_type: 'factory' })
  })

  it('list 透传分页与类型过滤参数', () => {
    trashApi.list({ page: 2, pageSize: 20, entity_type: 'blog' })
    expect(http.get).toHaveBeenCalledWith('/admin/trash', { page: 2, pageSize: 20, entity_type: 'blog' })
  })

  it('empty 清空不携带参数', () => {
    trashApi.empty()
    expect(http.delete).toHaveBeenCalledWith('/admin/trash')
  })
})
