import { describe, it, expect } from 'vitest'
import { MENU_CONFIG, type MenuItem } from '@/config/menu'

function hasUniqueSortOrder(items: MenuItem[]): boolean {
  const orders = items.map(i => i.sort_order)
  return new Set(orders).size === orders.length
}

function topLevelKeys(): string[] {
  return MENU_CONFIG.sort((a, b) => a.sort_order - b.sort_order).map(i => i.key)
}

describe('MENU_CONFIG - structure', () => {
  it('top-level sort_order unique', () => {
    expect(hasUniqueSortOrder(MENU_CONFIG)).toBe(true)
  })

  it('each group child sort_order unique', () => {
    for (const group of MENU_CONFIG.filter(i => i.children)) {
      expect(hasUniqueSortOrder(group.children!)).toBe(true)
    }
  })

  it('every item has key, label, index, sort_order', () => {
    function assertFields(items: MenuItem[]) {
      for (const item of items) {
        expect(item.key).toBeTruthy()
        expect(item.label).toBeTruthy()
        expect(item.index).toBeTruthy()
        expect(item.sort_order).toBeGreaterThan(0)
        if (item.children) assertFields(item.children)
      }
    }
    assertFields(MENU_CONFIG)
  })

  it('standalone items index starts with /', () => {
    for (const item of MENU_CONFIG.filter(i => !i.children)) {
      expect(item.index.startsWith('/')).toBe(true)
    }
  })

  it('group items index does not start with /', () => {
    for (const group of MENU_CONFIG.filter(i => i.children)) {
      expect(group.index.startsWith('/')).toBe(false)
    }
  })
})

describe('MENU_CONFIG - sort order', () => {
  it('top-level sorted by sort_order ascending', () => {
    const orders = MENU_CONFIG.map(i => i.sort_order)
    expect(orders).toEqual([...orders].sort((a, b) => a - b))
  })

  it('sales menu come after products in correct sequence', () => {
    const order = topLevelKeys()
    const productsIdx = order.indexOf('products')
    const leadsIdx = order.indexOf('leads')
    const quotesIdx = order.indexOf('quotes')
    const notificationsIdx = order.indexOf('notifications')
    expect(productsIdx).toBeGreaterThanOrEqual(0)
    expect(leadsIdx).toBeGreaterThan(productsIdx)
    expect(quotesIdx).toBeGreaterThan(leadsIdx)
    expect(notificationsIdx).toBeGreaterThan(quotesIdx)
  })

  it('insights group after system group', () => {
    const order = topLevelKeys()
    expect(order.indexOf('insights')).toBeGreaterThan(order.indexOf('system'))
  })

  it('operation-logs inside system group', () => {
    const system = MENU_CONFIG.find(i => i.key === 'system')!
    const opLogs = system.children!.find(i => i.key === 'operation-logs')
    expect(opLogs).toBeDefined()
  })

  it('insights group contains portal-preview and analytics', () => {
    const insights = MENU_CONFIG.find(i => i.key === 'insights')!
    expect(insights.children!.map(i => i.key)).toEqual(['portal-preview', 'analytics'])
  })
})

describe('MENU_CONFIG - permissions', () => {
  it('standalone items have permission', () => {
    for (const item of MENU_CONFIG.filter(i => !i.children)) {
      expect(item.permission).toBeTruthy()
    }
  })

  it('group children have permission', () => {
    for (const group of MENU_CONFIG.filter(i => i.children)) {
      for (const child of group.children!) {
        expect(child.permission).toBeTruthy()
      }
    }
  })
})
