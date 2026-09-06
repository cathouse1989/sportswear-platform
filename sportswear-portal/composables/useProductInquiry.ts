// 产品询盘信息构建（公共逻辑）：
// 供「产品详情页复制信息 / WhatsApp 询盘 / 联系页预填 / 全局悬浮询盘」复用，
// 确保询盘与 WhatsApp 跳转始终携带完整产品信息（SKU / 规格 / MOQ / 定制等）。

// 构建完整产品信息文本（SKU / 名称 / 规格 / MOQ / 定制等）
export function buildProductInfoText(p: any, selectionSummary?: string): string {
  if (!p) return ''
  return [
    `SKU: ${p.sku}`,
    p.name ? `Name: ${p.name}` : '',
    p.type ? `Type: ${p.type}` : '',
    p.gender ? `Gender: ${p.gender}` : '',
    p.material ? `Material: ${p.material}` : '',
    p.composition ? `Composition: ${p.composition}` : '',
    p.weight ? `Weight: ${p.weight}` : '',
    p.elasticity ? `Elasticity: ${p.elasticity}` : '',
    p.fit ? `Fit: ${p.fit}` : '',
    p.season ? `Season: ${p.season}` : '',
    p.size_range ? `Size Range: ${p.size_range}` : '',
    p.brief ? `Brief: ${p.brief}` : '',
    p.description ? `Description: ${p.description}` : '',
    p.features ? `Features: ${p.features}` : '',
    p.usage ? `Usage: ${p.usage}` : '',
    `Sample MOQ: ${p.sample_moq || 1}`,
    `Production MOQ: ${p.production_moq || 300}`,
    `Color MOQ: ${p.color_moq || 100}`,
    `Size MOQ: ${p.size_moq || 100}`,
    p.specs?.length ? `\nSpecifications:\n${p.specs.map((s: any) => `  ${s.name}: ${s.value}`).join('\n')}` : '',
    selectionSummary ? `\nSelected Options: ${selectionSummary}` : '',
    p.customizations?.length ? `\nCustomization Options:\n${p.customizations.map((c: any) => `  ${c.type}: ${c.note || 'Available'}`).join('\n')}` : '',
  ].filter(Boolean).join('\n')
}

// 构建 WhatsApp 询盘消息：携带完整产品信息，便于销售在 WhatsApp 端直接识别产品
export function buildWhatsAppMessage(p: any, selectionSummary?: string): string {
  const info = buildProductInfoText(p, selectionSummary)
  if (!info) return 'Hello! I am interested in your sportswear products.'
  return `Hello! I am interested in the following product:\n\n${info}\n\nPlease send me more information (price, MOQ, lead time). Thank you!`
}
