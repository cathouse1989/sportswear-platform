# 英雄轮播图片维护指南

## 目录

本目录存储门户落地页 (`pages/index.vue`) 英雄轮播 HeroCarousel 的背景图片。

图片通过 `DEFAULT_IMAGES` 数组在 `pages/index.vue` 中引用，路径为 `/images/hero/<filename>`。

> **优先级**：后台 `hero_slides` 配置的图片 > 本目录的本地默认图片。
> 如后台返回轮播图片 URL，系统会优先使用后台配置的图片；本目录图片仅作兜底方案。

---

## 图片目录

| 文件名 | 尺寸 | 大小 | 对应轮播 | 说明 |
|--------|------|------|----------|------|
| `slide-1-sportswear-manufacturing.jpg` | 1600×900 | ~150KB | Slide 1 | 运动服制造 / 全球品牌 OEM-ODM |
| `slide-2-running-training.jpg` | 1600×900 | ~175KB | Slide 2 | 跑步与训练服装 / 技术面料 |
| `slide-3-team-uniforms.jpg` | 1600×900 | ~382KB | Slide 3 | 专业团队 uniforms / 全色打印 |
| `slide-4-concept-to-product.jpg` | 1600×900 | ~385KB | Slide 4 | 从概念到产品 / 全程制造支持 |
| `alt-factory-athlete.jpg` | 1600×900 | ~305KB | 备用 | 工厂/运动员备用图 |
| `alt-yoga-gear.jpg` | 1600×900 | ~172KB | 备用 | 瑜伽装备备用图 |
| `alt-yoga-studio.jpg` | 1600×900 | ~212KB | 备用 | 瑜伽室备用图 |

---

## 替换图片

1. **将新图片放入本目录**，建议使用 1600×900 的横版 JPG/PNG 图片。
2. **更新 `pages/index.vue`** 中的 `DEFAULT_IMAGES` 数组路径。

```javascript
const DEFAULT_IMAGES = [
  '/images/hero/slide-1-sportswear-manufacturing.jpg',
  '/images/hero/slide-2-running-training.jpg',
  '/images/hero/slide-3-team-uniforms.jpg',
  '/images/hero/slide-4-concept-to-product.jpg',
]
```

3. **提交到版本控制**，图片会随代码一起部署。

---

## 调整轮播数量

如需增加或减少轮播图数量：

- **增加**：在 `DEFAULT_IMAGES` 中添加新路径，对应的文案词条 `home.hero_title_5`、`home.hero_sub_5` 会自动回退到语言包默认值。
- **减少**：从 `DEFAULT_IMAGES` 中删除路径即可。

> 文案词条在 `locales/en.json`、`locales/zh.json` 等语言包中配置，
> key 格式为 `home.hero_title_{N}` / `home.hero_sub_{N}`。
> 后台「词条管理」可动态覆盖这些文案。

---

## 图片来源与版权

所有图片均从 **Unsplash** 下载，遵循 Unsplash 许可协议（可自由用于商业和非商业用途，无需署名）。

原始 Unsplash 图片 ID：

| 文件名 | Unsplash Photo ID |
|--------|-------------------|
| slide-1 | `photo-1599901860904-17e6ed7083a0` |
| slide-2 | `photo-1571019613454-1cb2f99b2d8b` |
| slide-3 | `photo-1517466787929-bc90951d0974` |
| slide-4 | `photo-1577221084712-45b0445d2b00` |
| alt-factory-athlete | `photo-1581091226825-a6a2a5aee158` |
| alt-yoga-gear | `photo-1575052814086-f385e2e2ad1b` |
| alt-yoga-studio | `photo-1601925260368-ae2f83cf8b7f` |

---

## 图片命名规范

- `slide-N-<theme-description>.jpg` — 主要轮播图（N = 1, 2, 3, ...）
- `alt-<theme-description>.jpg` — 备用图片，可用于手动切换或 A/B 测试

命名中的 `<theme-description>` 应简要描述图片内容，便于维护时快速识别。
