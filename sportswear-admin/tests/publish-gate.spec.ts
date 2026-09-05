import { describe, it, expect } from 'vitest'
import { checkProductGate, checkBlogGate, checkCaseGate, gateAlertMessage } from '@/utils/publish-gate'

describe('checkProductGate（P0-#2 发布质检门）', () => {
  const complete = {
    cover_image: 'https://cdn.example.com/a.jpg',
    name: 'Custom Yoga Leggings',
    brief: 'High-waist seamless leggings',
    description: 'Full description here',
    sample_moq: 1,
    production_moq: 300,
  }

  it('完整产品：ok=true，无阻断', () => {
    const gate = checkProductGate(complete)
    expect(gate.ok).toBe(true)
    expect(gate.errors).toEqual([])
  })

  it('缺封面图：阻断并列出缺失项', () => {
    const gate = checkProductGate({ ...complete, cover_image: '' })
    expect(gate.ok).toBe(false)
    expect(gate.errors).toContain('封面图')
  })

  it('缺名称/简介/描述：逐一列入 errors', () => {
    const gate = checkProductGate({ ...complete, name: '', brief: '  ', description: undefined })
    expect(gate.ok).toBe(false)
    expect(gate.errors).toContain('英文名称')
    expect(gate.errors).toContain('简介 brief')
    expect(gate.errors).toContain('描述 description')
  })

  it('MOQ 不合理（sample > production）：阻断', () => {
    const gate = checkProductGate({ ...complete, sample_moq: 500, production_moq: 100 })
    expect(gate.ok).toBe(false)
    expect(gate.errors.join()).toContain('MOQ')
  })

  it('缺 SEO 标题：不阻断但列入 warnings', () => {
    const gate = checkProductGate({ ...complete, seo_title: '' })
    expect(gate.ok).toBe(true)
    expect(gate.warnings).toContain('SEO 标题')
  })
})

describe('checkBlogGate', () => {
  const complete = { title: 'OEM vs ODM', slug: 'oem-vs-odm', content: 'body', cover_image: 'https://cdn/a.jpg' }

  it('完整博客：ok=true', () => {
    expect(checkBlogGate(complete).ok).toBe(true)
  })

  it('缺正文：阻断', () => {
    const gate = checkBlogGate({ ...complete, content: '' })
    expect(gate.ok).toBe(false)
    expect(gate.errors).toContain('正文')
  })

  it('缺封面图：仅警告不阻断', () => {
    const gate = checkBlogGate({ ...complete, cover_image: '' })
    expect(gate.ok).toBe(true)
    expect(gate.warnings).toContain('封面图')
  })
})

describe('checkCaseGate', () => {
  const complete = {
    title: 'Yoga Brand OEM',
    slug: 'yoga-brand-oem',
    solution: 'Full turnkey OEM solution',
    cover_image: 'https://cdn/a.jpg',
    client_industry: 'Fitness Brand',
  }

  it('完整案例：ok=true', () => {
    expect(checkCaseGate(complete).ok).toBe(true)
  })

  it('缺解决方案：阻断', () => {
    const gate = checkCaseGate({ ...complete, solution: '' })
    expect(gate.ok).toBe(false)
    expect(gate.errors).toContain('解决方案')
  })

  it('缺封面图/客户行业：仅警告不阻断', () => {
    const gate = checkCaseGate({ ...complete, cover_image: '', client_industry: '  ' })
    expect(gate.ok).toBe(true)
    expect(gate.warnings).toContain('封面图')
    expect(gate.warnings).toContain('客户行业')
  })
})

describe('gateAlertMessage', () => {
  it('拼接阻断与建议文案', () => {
    const gate = checkProductGate({ cover_image: '', seo_title: '' })
    const msg = gateAlertMessage(gate)
    expect(msg).toContain('无法发布')
    expect(msg).toContain('封面图')
    expect(msg).toContain('建议补充（不阻断）')
  })
})