<template>
  <div class="product-page">
    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-info">
              <div class="stat-label">产品总数</div>
              <div class="stat-value">{{ totalCount }}</div>
            </div>
            <div class="stat-icon-box">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/></svg>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-info">
              <div class="stat-label">已发布</div>
              <div class="stat-value stat-value-success">{{ publishedCount }}</div>
            </div>
            <div class="stat-icon-box stat-icon-green">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-info">
              <div class="stat-label">草稿</div>
              <div class="stat-value stat-value-warning">{{ draftCount }}</div>
            </div>
            <div class="stat-icon-box stat-icon-orange">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-info">
              <div class="stat-label">已下线</div>
              <div class="stat-value stat-value-danger">{{ offlineCount }}</div>
            </div>
            <div class="stat-icon-box stat-icon-red">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/></svg>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 主卡片 -->
    <el-card shadow="never" class="main-card">
      <!-- 工具栏 -->
      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="keyword" placeholder="搜索 SKU / Slug..." clearable style="width: 220px" @keyup.enter="handleSearch" />
          <el-select v-model="status" placeholder="状态" clearable style="width: 120px" @change="handleSearch">
            <el-option v-for="o in enumOptions('product.status')" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-select v-model="gender" placeholder="性别" clearable style="width: 110px" @change="handleSearch">
            <el-option v-for="o in enumOptions('product.gender')" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-select v-model="typeFilter" placeholder="类型" clearable style="width: 110px" @change="handleSearch">
            <el-option v-for="o in enumOptions('product.type')" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
        </div>
        <div class="toolbar-right">
          <el-button v-permission="'product:create'" type="primary" @click="openCreateDialog">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:4px"><path d="M12 5v14m-7-7h14"/></svg>
            新建产品
          </el-button>
        </div>
      </div>
      <!-- 表格 -->
      <el-table :data="products" v-loading="loading" stripe class="product-table" :header-cell-style="{ background: '#fafafa', color: '#374151', fontWeight: 600 }">
        <el-table-column label="封面" width="68" align="center">
          <template #default="{ row }">
            <div class="cover-cell">
              <el-image v-if="row.cover_image" :src="row.cover_image" fit="cover" class="cover-thumb" />
              <div v-else class="cover-placeholder-cell">📷</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="sku" label="SKU" width="130" />
        <el-table-column prop="slug" label="Slug" min-width="150" />
        <el-table-column label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="enumTag('product.type', row.type)">{{ enumLabel('product.type', row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="性别" width="66" align="center">
          <template #default="{ row }">
            <span class="gender-tag">{{ enumLabel('product.gender', row.gender) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="精选" width="56" align="center">
          <template #default="{ row }">
            <span v-if="row.is_featured" class="featured-star">★</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="enumTag('product.status', row.status)" size="small" effect="plain">{{ enumLabel('product.status', row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="production_moq" label="MOQ" width="72" align="right" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }: any">
            <div class="action-btns">
              <el-button v-permission="'product:update'" size="small" @click="openEditDialog(row)">编辑</el-button>
              <el-button v-permission="'product:publish'" v-if="row.status !== 'published'" size="small" type="success" @click="handlePublish(row)">发布</el-button>
              <el-button v-permission="'product:publish'" v-else size="small" type="warning" @click="handleUnpublish(row)">下线</el-button>
              <el-button v-permission="'product:delete'" size="small" type="danger" @click="handleDelete(row)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        :page-sizes="[10, 20, 50, 100]"
        @current-change="loadData"
        @size-change="loadData"
      />
    </el-card>


    <!-- 编辑抽屉：基本信息 / 详细信息 / MOQ / 图片 / 规格 / 视频 / 定制 / 翻译 / 系列面料 / SEO 集中管理 -->
    <el-drawer
      v-model="drawerVisible"
      :title="editingId ? '编辑产品' : '新建产品'"
      size="92%"
      :close-on-click-modal="false"
      class="product-drawer"
    >
      <el-form :model="form" label-width="110px" class="product-form">
        <el-tabs v-model="activeTab" class="drawer-tabs">
          <el-tab-pane label="基本信息" name="basic">
            <div class="tab-pane-content">
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="SKU" required>
                    <el-input v-model="form.sku" placeholder="产品 SKU 编号" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="Slug" required>
                    <el-input v-model="form.slug" placeholder="URL 标识" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="英文名称" required>
                    <el-input v-model="form.name" placeholder="产品英文名称（发布必填）" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="类型">
                    <el-select v-model="form.type" class="full-width">
                      <el-option v-for="o in enumOptions('product.type')" :key="o.value" :label="o.label" :value="o.value" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="性别">
                    <el-select v-model="form.gender" class="full-width">
                      <el-option v-for="o in enumOptions('product.gender')" :key="o.value" :label="o.label" :value="o.value" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="分类">
                    <el-select v-model="form.category_id" clearable placeholder="选择分类" class="full-width">
                      <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="排序">
                    <el-input-number v-model="form.sort_order" :min="0" class="full-width" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="精选">
                    <el-switch v-model="form.is_featured" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="新品">
                    <el-switch v-model="form.is_new" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="封面图片">
                <MediaPicker v-model="form.cover_image" />
              </el-form-item>
              <el-form-item label="简述">
                <el-input v-model="form.brief" type="textarea" :rows="2" placeholder="产品简短描述" />
              </el-form-item>
            </div>
          </el-tab-pane>
          <el-tab-pane label="详细信息" name="detail">
            <div class="tab-pane-content">
              <el-form-item label="描述">
                <el-input v-model="form.description" type="textarea" :rows="4" placeholder="产品详细描述" />
              </el-form-item>
              <el-form-item label="特性">
                <el-input v-model="form.features" type="textarea" :rows="3" placeholder="产品特性，每行一项" />
              </el-form-item>
              <el-form-item label="用途">
                <el-input v-model="form.usage" type="textarea" :rows="2" placeholder="产品用途说明" />
              </el-form-item>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="面料">
                    <el-input v-model="form.material" placeholder="如：Polyester, Cotton" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="成分">
                    <el-input v-model="form.composition" placeholder="如：90% Polyester, 10% Elastane" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="重量"><el-input v-model="form.weight" placeholder="如：180g/m²" /></el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="弹性"><el-input v-model="form.elasticity" placeholder="如：Medium" /></el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="合身度"><el-input v-model="form.fit" placeholder="如：Regular" /></el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="支撑等级"><el-input v-model="form.support_level" placeholder="如：High" /></el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="季节"><el-input v-model="form.season" placeholder="如：All Season" /></el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="尺码范围"><el-input v-model="form.size_range" placeholder="如：XS-3XL" /></el-form-item>
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>

          <el-tab-pane label="MOQ" name="moq">
            <div class="tab-pane-content">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="样品 MOQ"><el-input-number v-model="form.sample_moq" :min="0" class="full-width" /></el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="大货 MOQ"><el-input-number v-model="form.production_moq" :min="0" class="full-width" /></el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="颜色 MOQ"><el-input-number v-model="form.color_moq" :min="0" class="full-width" /></el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="尺码 MOQ"><el-input-number v-model="form.size_moq" :min="0" class="full-width" /></el-form-item>
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>


          <el-tab-pane label="图片" name="images">
            <div class="tab-pane-content">
              <div class="img-mgr-section">
                <div class="img-mgr-section-title">封面图片</div>
                <div class="cover-section">
                  <el-image v-if="form.cover_image" :src="form.cover_image" fit="cover" class="cover-preview-lg" />
                  <div v-else class="cover-placeholder-lg"><span>暂无封面</span></div>
                  <div class="cover-section-actions">
                    <el-upload :show-file-list="false" :before-upload="beforeCoverUpload" accept="image/*">
                      <el-button size="small" type="primary">上传封面图片</el-button>
                    </el-upload>
                    <div class="cover-tip">上传后即更新「基本信息」中的封面，保存时统一提交。</div>
                  </div>
                </div>
              </div>
              <div class="img-mgr-section">
                <div class="img-mgr-section-title">
                  产品图片
                  <el-tag size="small" type="info" effect="plain">{{ productImages.length }} 张</el-tag>
                </div>
                <div class="image-grid">
                  <div v-for="(img, idx) in productImages" :key="img.id" class="image-item-wrap">
                    <div class="image-item">
                      <el-image :src="img.url" fit="cover" class="image-thumb" />
                      <div class="image-item-overlay">
                        <el-tooltip content="设为封面" placement="top">
                          <el-button size="small" circle @click="setAsCover(img)" class="overlay-btn">
                            <template #icon><el-icon color="#fff"><Star /></el-icon></template>
                          </el-button>
                        </el-tooltip>
                        <el-tooltip content="前移" placement="top">
                          <el-button size="small" circle @click="moveImage(idx, -1)" :disabled="idx === 0" class="overlay-btn">
                            <template #icon><el-icon color="#fff"><ArrowUp /></el-icon></template>
                          </el-button>
                        </el-tooltip>
                        <el-tooltip content="后移" placement="top">
                          <el-button size="small" circle @click="moveImage(idx, 1)" :disabled="idx === productImages.length - 1" class="overlay-btn">
                            <template #icon><el-icon color="#fff"><ArrowDown /></el-icon></template>
                          </el-button>
                        </el-tooltip>
                        <el-tooltip content="删除" placement="top">
                          <el-button size="small" circle @click="deleteImage(idx)" class="overlay-btn">
                            <template #icon><el-icon color="#fff"><Delete /></el-icon></template>
                          </el-button>
                        </el-tooltip>
                      </div>
                    </div>
                    <el-input v-model="img.alt" placeholder="Alt 文案" size="small" class="image-alt-input" />
                  </div>
                  <el-upload :show-file-list="false" :before-upload="beforeImageUpload" accept="image/*">
                    <div class="image-upload-box">
                      <div class="upload-trigger">
                        <span>＋</span>
                        <span>上传图片</span>
                      </div>
                    </div>
                  </el-upload>
                </div>
                <div v-if="uploading" class="upload-progress">
                  <el-progress :percentage="uploadProgress" />
                </div>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="规格" name="specs">
            <div class="tab-pane-content">
              <el-alert type="info" :closable="false" show-icon style="margin-bottom:12px;" title="规格值支持用逗号 / 斜杠分隔多个可选项（如 Color: Black, White, Navy），门户会渲染为可点选的颜色/尺码选择器；单一值则展示在规格表格中。" />
              <div style="display:flex; align-items:center; gap:8px; flex-wrap:wrap; margin-bottom:12px;">
                <span style="font-size:13px; color:#6b7280;">快速添加：</span>
                <el-button size="small" @click="addSpecPreset('Color', 'Black, White, Navy, Red')">颜色 Color</el-button>
                <el-button size="small" @click="addSpecPreset('Size', 'XS, S, M, L, XL, XXL')">尺码 Size</el-button>
                <el-button size="small" @click="addSpecPreset('Fabric', 'Polyester, Cotton, Nylon')">面料 Fabric</el-button>
              </div>
              <div class="spec-list">
                <div v-for="(spec, idx) in specs" :key="idx" class="spec-item">
                  <div class="spec-row">
                    <el-input v-model="spec.name" placeholder="规格名称，如：面料" style="width: 180px" />
                    <el-input v-model="spec.value" placeholder="规格值 (英文)" style="width: 220px" />
                    <el-button size="small" circle @click="moveSpec(idx, -1)" :disabled="idx === 0" title="上移">
                      <template #icon><el-icon><ArrowUp /></el-icon></template>
                    </el-button>
                    <el-button size="small" circle @click="moveSpec(idx, 1)" :disabled="idx === specs.length - 1" title="下移">
                      <template #icon><el-icon><ArrowDown /></el-icon></template>
                    </el-button>
                    <el-button type="danger" size="small" @click="removeSpec(idx)" circle>
                      <template #icon><el-icon><Delete /></el-icon></template>
                    </el-button>
                  </div>
                  <div class="spec-trans">
                    <InlineTrans v-model="spec.value_trans" prefix="值" />
                  </div>
                </div>
              </div>
              <el-button size="small" @click="addSpec" class="mt-2">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:4px"><path d="M12 5v14m-7-7h14"/></svg>
                添加规格
              </el-button>
            </div>
          </el-tab-pane>

          <el-tab-pane label="视频" name="videos">
            <div class="tab-pane-content">
              <div class="video-list">
                <div v-for="(v, idx) in videos" :key="idx" class="video-item">
                  <div class="video-row">
                    <el-select v-model="v.type" placeholder="类型" style="width: 120px">
                      <el-option v-for="opt in enumOptions('product.video.type')" :key="opt.value" :label="opt.label" :value="opt.value" />
                    </el-select>
                    <el-input v-model="v.title" placeholder="标题 (英文)" style="width: 140px" />
                    <el-input v-model="v.url" placeholder="视频 URL（支持 YouTube 链接）" style="width: 220px" />
                    <el-button type="danger" size="small" @click="removeVideo(idx)" circle>
                      <template #icon><el-icon><Delete /></el-icon></template>
                    </el-button>
                  </div>
                  <div class="video-cover">
                    <span class="video-cover-label">封面图</span>
                    <MediaPicker v-model="v.cover" class="video-cover-picker" />
                  </div>
                  <div class="video-trans">
                    <InlineTrans v-model="v.title_trans" prefix="标题" />
                  </div>
                </div>
              </div>
              <el-button size="small" @click="addVideo" class="mt-2">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:4px"><path d="M12 5v14m-7-7h14"/></svg>
                添加视频
              </el-button>
            </div>
          </el-tab-pane>
          <el-tab-pane label="定制" name="customizations">
            <div class="tab-pane-content">
              <div class="custom-list">
                <div v-for="(c, idx) in customizations" :key="idx" class="custom-item">
                  <div class="custom-row">
                    <el-select v-model="c.type" placeholder="定制类型" style="width: 160px">
                      <el-option v-for="opt in enumOptions('product.customization.type')" :key="opt.value" :label="opt.label" :value="opt.value" />
                    </el-select>
                    <el-switch v-model="c.is_enabled" active-text="启用" style="width: 90px" />
                    <el-button type="danger" size="small" @click="removeCustomization(idx)" circle>
                      <template #icon><el-icon><Delete /></el-icon></template>
                    </el-button>
                  </div>
                  <div class="custom-trans">
                    <el-input v-model="c.note" placeholder="备注 (英文)" />
                    <InlineTrans v-model="c.note_trans" prefix="备注" />
                  </div>
                </div>
              </div>
              <el-button size="small" @click="addCustomization" class="mt-2">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:4px"><path d="M12 5v14m-7-7h14"/></svg>
                添加定制选项
              </el-button>
            </div>
          </el-tab-pane>

          <el-tab-pane label="翻译" name="translations">
            <div class="tab-pane-content">
              <el-alert type="info" :closable="false" show-icon class="translation-alert" title="英文为源语言：上方「基本信息 / 详细信息」中的内容即英文源，保存时自动写入英文翻译；其他语言可从英文一键复制后微调。" />
              <TransEditor v-model="translations" :fields="TRANS_FIELDS" :source="transSource" :langs="TRANS_LANGS" source-label="English" />
              <el-divider content-position="left">各语言排序（越小越靠前，0 = 跟随全局排序）</el-divider>
              <el-row :gutter="12">
                <el-col v-for="l in TRANS_LANGS" :key="l.value" :span="8">
                  <el-form-item :label="l.label + '排序'">
                    <el-input-number v-model="transSortOrders[l.value]" :min="0" controls-position="right" style="width: 100%" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>
          <el-tab-pane label="系列/面料" name="series">
            <div class="tab-pane-content">
              <el-form-item label="所属系列">
                <el-select v-model="form.series_ids" multiple clearable placeholder="选择系列" class="full-width">
                  <el-option v-for="s in seriesList" :key="s.id" :label="s.name" :value="s.id" />
                </el-select>
              </el-form-item>
              <el-form-item label="可选面料">
                <el-select v-model="form.fabric_ids" multiple clearable placeholder="选择面料" class="full-width">
                  <el-option v-for="f in fabricList" :key="f.id" :label="f.name + ' (' + f.code + ')'" :value="f.id" />
                </el-select>
              </el-form-item>
            </div>
          </el-tab-pane>

          <el-tab-pane label="SEO" name="seo">
            <div class="tab-pane-content">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="SEO Title"><el-input v-model="form.seo.title" placeholder="搜索引擎标题" /></el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Keywords"><el-input v-model="form.seo.keywords" placeholder="关键词，逗号分隔" /></el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="SEO Description">
                <el-input v-model="form.seo.description" type="textarea" :rows="2" placeholder="搜索引擎描述" />
              </el-form-item>
              <el-divider />
              <div class="section-subtitle">Open Graph</div>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="OG Title"><el-input v-model="form.seo.og_title" placeholder="社交分享标题" /></el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="OG Image"><MediaPicker v-model="form.seo.og_image" /></el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="OG Description">
                <el-input v-model="form.seo.og_description" type="textarea" :rows="2" placeholder="社交分享描述" />
              </el-form-item>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="drawerVisible = false" size="large">取消</el-button>
          <el-button type="primary" @click="handleSave" size="large" :loading="saving">保存产品</el-button>
        </div>
      </template>
    </el-drawer>
  </div>




</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Star, Delete, ArrowDown, ArrowUp, Search } from '@element-plus/icons-vue'
import { productApi, mediaApi, categoryApi, seriesApi, fabricApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import { useEnumDict } from '@/composables/useEnumDict'
import type { Product, Category, Series, Fabric } from '@/types'
import { checkProductGate, gateAlertMessage } from '@/utils/publish-gate'
import MediaPicker from '@/components/media/MediaPicker.vue'
import TransEditor from '@/components/cms/TransEditor.vue'
import InlineTrans from '@/components/cms/InlineTrans.vue'
import { DEFAULT_TRANS_LANGS, createEmptyTranslations, translationsToRecord, translationsToPayload, parseJsonTrans, buildJsonTrans } from '@/composables/useTransRecord'

const { ensureLoaded: loadEnumDict, options: enumOptions, label: enumLabel, tagType: enumTag } = useEnumDict()

const products = ref<Product[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const status = ref('')
const gender = ref('')
const typeFilter = ref('')

// 状态统计卡：来自后端全量统计接口，不受列表筛选/分页影响
const totalCount = ref(0)
const publishedCount = ref(0)
const draftCount = ref(0)
const offlineCount = ref(0)

const drawerVisible = ref(false)
const editingId = ref('')
const saving = ref(false)
const activeTab = ref('basic')
const categories = ref<Category[]>([])
const seriesList = ref<Series[]>([])
const fabricList = ref<Fabric[]>([])

const form = reactive({
  sku: '',
  slug: '',
  name: '',
  category_id: '',
  type: 'both',
  gender: 'unisex',
  status: 'draft',
  is_featured: false,
  is_new: false,
  sort_order: 0,
  cover_image: '',
  brief: '',
  description: '',
  features: '',
  usage: '',
  material: '',
  composition: '',
  weight: '',
  elasticity: '',
  fit: '',
  support_level: '',
  season: '',
  size_range: '',
  sample_moq: 1,
  production_moq: 300,
  color_moq: 100,
  size_moq: 100,
  seo: { title: '', description: '', keywords: '', og_title: '', og_description: '', og_image: '' },
  series_ids: [] as string[],
  fabric_ids: [] as string[],
})

// 图片（封面走 form.cover_image，图集暂存在本地，保存时统一提交）
const productImages = ref<Array<{ id: string; url: string; alt: string; sort_order: number; thumbnail?: string }>>([])
const uploading = ref(false)
const uploadProgress = ref(0)

// 规格
const specs = ref<Array<{ name: string; value: string; value_trans: Record<string, string> }>>([])

// 视频
const videos = ref<Array<{ type: string; title: string; url: string; cover: string; title_trans: Record<string, string> }>>([])

// 定制
const customizations = ref<Array<{ type: string; is_enabled: boolean; note: string; note_trans: Record<string, string> }>>([])

// 翻译
const PRODUCT_TRANS_FIELDS = [
  'name', 'brief', 'description', 'features', 'usage',
  'material', 'composition', 'weight', 'elasticity', 'fit',
  'support_level', 'season', 'size_range',
]
const TRANS_LANGS = DEFAULT_TRANS_LANGS
const TRANS_FIELDS = [
  { key: 'name', label: '名称' },
  { key: 'brief', label: '简述', type: 'textarea' as const, rows: 2 },
  { key: 'description', label: '描述', type: 'textarea' as const, rows: 3 },
  { key: 'features', label: '特性', type: 'textarea' as const, rows: 2 },
  { key: 'usage', label: '用途', type: 'textarea' as const, rows: 2 },
  { key: 'material', label: '面料' },
  { key: 'composition', label: '成分' },
  { key: 'weight', label: '重量' },
  { key: 'elasticity', label: '弹性' },
  { key: 'fit', label: '合身度' },
  { key: 'support_level', label: '支撑等级' },
  { key: 'season', label: '季节' },
  { key: 'size_range', label: '尺码范围' },
]
const translations = ref<Record<string, Record<string, string>>>({})
const transSortOrders = ref<Record<string, number>>({ zh: 0, es: 0, fr: 0 })
const transSource = computed(() => ({
  name: form.name, brief: form.brief, description: form.description,
  features: form.features, usage: form.usage, material: form.material,
  composition: form.composition, weight: form.weight, elasticity: form.elasticity,
  fit: form.fit, support_level: form.support_level, season: form.season,
  size_range: form.size_range,
}))


async function loadData() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: page.value,
      pageSize: pageSize.value,
      keyword: keyword.value,
      status: status.value,
    }
    if (gender.value) params.gender = gender.value
    if (typeFilter.value) params.type = typeFilter.value
    const result = await productApi.list(params)
    products.value = result.items
    total.value = result.total
  } finally {
    loading.value = false
  }
  loadProductStats()
}

async function loadProductStats() {
  try {
    const s = await productApi.stats()
    totalCount.value = s.total
    publishedCount.value = s.published
    draftCount.value = s.draft
    offlineCount.value = s.offline
  } catch {}
}

async function loadCategories() {
  try { categories.value = await categoryApi.list() } catch {}
}

async function loadSeries() {
  try { seriesList.value = await seriesApi.list() } catch {}
}

async function loadFabrics() {
  try { fabricList.value = await fabricApi.list() } catch {}
}

function handleSearch() {
  page.value = 1
  loadData()
}

function resetForm() {
  Object.assign(form, {
    sku: '', slug: '', name: '', category_id: '', type: 'both', gender: 'unisex', status: 'draft',
    is_featured: false, is_new: false, sort_order: 0, cover_image: '', brief: '', description: '',
    features: '', usage: '', material: '', composition: '', weight: '', elasticity: '', fit: '',
    support_level: '', season: '', size_range: '', sample_moq: 1, production_moq: 300,
    color_moq: 100, size_moq: 100,
    seo: { title: '', description: '', keywords: '', og_title: '', og_description: '', og_image: '' },
    series_ids: [], fabric_ids: [],
  })
  translations.value = createEmptyTranslations(PRODUCT_TRANS_FIELDS)
  transSortOrders.value = { zh: 0, es: 0, fr: 0 }
  productImages.value = []
  specs.value = []
  videos.value = []
  customizations.value = []
  activeTab.value = 'basic'
}

function openCreateDialog() {
  editingId.value = ''
  resetForm()
  drawerVisible.value = true
}

async function openEditDialog(row: Product) {
  editingId.value = row.id
  resetForm()
  try {
    const detail = await productApi.get(row.id)
    Object.assign(form, {
      sku: detail.sku, slug: detail.slug,
      name: detail.translations?.find((t: any) => t.language === 'en')?.name || '',
      category_id: detail.category_id || '', type: detail.type, gender: detail.gender,
      status: detail.status, is_featured: detail.is_featured, is_new: detail.is_new,
      sort_order: detail.sort_order, cover_image: detail.cover_image || '',
      brief: detail.brief || '', description: detail.description || '',
      features: detail.features || '', usage: detail.usage || '',
      material: detail.material || '', composition: detail.composition || '',
      weight: detail.weight || '', elasticity: detail.elasticity || '', fit: detail.fit || '',
      support_level: detail.support_level || '', season: detail.season || '', size_range: detail.size_range || '',
      sample_moq: detail.sample_moq || 1, production_moq: detail.production_moq || 300,
      color_moq: detail.color_moq || 100, size_moq: detail.size_moq || 100,
      series_ids: (detail.series || []).map((s: any) => s.id),
      fabric_ids: (detail.fabrics || []).map((f: any) => f.id),
    })
    const transList = (detail.translations || []).filter((t: any) => t.language !== 'en')
    translations.value = translationsToRecord(transList, PRODUCT_TRANS_FIELDS)
    transSortOrders.value = { zh: 0, es: 0, fr: 0 }
    for (const t of transList) {
      if (t.language && transSortOrders.value[t.language] !== undefined) transSortOrders.value[t.language] = t.sort_order || 0
    }
    if (detail.seo) {
      form.seo = {
        title: detail.seo.title || '', description: detail.seo.description || '',
        keywords: detail.seo.keywords || '', og_title: detail.seo.og_title || '',
        og_description: detail.seo.og_description || '', og_image: detail.seo.og_image || '',
      }
    }
    // 子资源回填
    if (detail.images?.length) {
      productImages.value = detail.images
        .filter((i: any) => i.type !== 'cover')
        .map((i: any) => ({ id: i.id, url: i.url, alt: i.alt || '', sort_order: i.sort_order, thumbnail: i.thumbnail || '' }))
    }
    if (detail.specs?.length) {
      specs.value = detail.specs.map((s: any) => ({ name: s.name, value: s.value, value_trans: parseJsonTrans(s.translations) }))
    }
    if (detail.videos?.length) {
      videos.value = detail.videos.map((v: any) => ({
        type: v.type || 'product', title: v.title || '', url: v.url || '',
        cover: v.cover || '', title_trans: parseJsonTrans(v.translations),
      }))
    }
    if (detail.customizations?.length) {
      customizations.value = detail.customizations.map((c: any) => ({
        type: c.type, is_enabled: c.is_enabled, note: c.note || '', note_trans: parseJsonTrans(c.translations),
      }))
    }
  } catch {
    Object.assign(form, {
      sku: row.sku, slug: row.slug, type: row.type, gender: row.gender,
      status: row.status, is_featured: row.is_featured, is_new: row.is_new,
      cover_image: row.cover_image || '', brief: row.brief || '', production_moq: row.production_moq || 300,
    })
  }
  drawerVisible.value = true
}

// 构建完整提交 payload：主表 + 翻译 + SEO + 系列/面料 + 图片/规格/视频/定制，一次整体保存
function buildPayload(): Record<string, any> {
  const data: Record<string, any> = { ...form }
  const trans: Array<Record<string, any>> = translationsToPayload(translations.value, PRODUCT_TRANS_FIELDS, TRANS_LANGS)
    .map((t) => ({ ...t, sort_order: transSortOrders.value[t.language] ?? 0 }))
  if (String(form.name || '').trim()) {
    trans.unshift({
      language: 'en',
      name: String(form.name).trim(),
      brief: form.brief || '', description: form.description || '',
      features: form.features || '', usage: form.usage || '',
      material: form.material || '', composition: form.composition || '',
      weight: form.weight || '', elasticity: form.elasticity || '', fit: form.fit || '',
      support_level: form.support_level || '', season: form.season || '', size_range: form.size_range || '',
      sort_order: 0,
    })
  }
  data.translations = trans
  data.images = productImages.value.map((img, i) => ({
    type: 'gallery', url: img.url, thumbnail: img.thumbnail || '', alt: img.alt || '', sort_order: i,
  }))
  data.specs = specs.value.map((s: any, i) => {
    const { value_trans, ...rest } = s
    return { ...rest, sort_order: i, translations: buildJsonTrans(value_trans) }
  })
  data.videos = videos.value.map((v: any, i) => {
    const { title_trans, ...rest } = v
    return { ...rest, sort_order: i, translations: buildJsonTrans(title_trans) }
  })
  data.customizations = customizations.value.map((c: any) => {
    const { note_trans, ...rest } = c
    return { ...rest, translations: buildJsonTrans(note_trans) }
  })
  return data
}

async function handleSave() {
  saving.value = true
  try {
    const data = buildPayload()
    if (editingId.value) await productApi.update(editingId.value, data)
    else await productApi.create(data)
    ElMessage.success('保存成功')
    drawerVisible.value = false
    loadData()
  } catch {
    // 错误已提示
  } finally {
    saving.value = false
  }
}



// 图片：上传后仅暂存到本地，点「保存产品」时统一提交
async function beforeCoverUpload(file: File) {
  try {
    const formData = new FormData()
    formData.append('file', file)
    const result = await mediaApi.upload(formData)
    form.cover_image = result.url
    ElMessage.success('封面上传成功，保存后生效')
  } catch {
    ElMessage.error('封面上传失败')
  }
  return false
}

async function beforeImageUpload(file: File) {
  try {
    uploading.value = true
    uploadProgress.value = 0
    const formData = new FormData()
    formData.append('file', file)
    const result = await mediaApi.upload(formData)
    productImages.value.push({
      id: Date.now().toString(),
      url: result.url,
      alt: '',
      sort_order: productImages.value.length,
      thumbnail: result.thumbnail || '',
    })
    uploadProgress.value = 100
    ElMessage.success('图片上传成功，保存后生效')
  } catch {
    ElMessage.error('图片上传失败')
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
  return false
}

function setAsCover(img: { id: string; url: string }) {
  form.cover_image = img.url
  ElMessage.success('已设为封面，保存后生效')
}

function deleteImage(idx: number) {
  productImages.value.splice(idx, 1)
}

function moveImage(idx: number, dir: number) {
  const target = idx + dir
  if (target < 0 || target >= productImages.value.length) return
  const list = productImages.value
  const [item] = list.splice(idx, 1)
  list.splice(target, 0, item)
}

// 规格
function addSpec() {
  specs.value.push({ name: '', value: '', value_trans: {} })
}
function addSpecPreset(name: string, value: string) {
  specs.value.push({ name, value, value_trans: {} })
}
function removeSpec(idx: number) {
  specs.value.splice(idx, 1)
}
function moveSpec(idx: number, dir: number) {
  const target = idx + dir
  if (target < 0 || target >= specs.value.length) return
  const list = specs.value
  const [item] = list.splice(idx, 1)
  list.splice(target, 0, item)
}

// 视频
function addVideo() {
  videos.value.push({ type: 'product', title: '', url: '', cover: '', title_trans: {} })
}
function removeVideo(idx: number) {
  videos.value.splice(idx, 1)
}

// 定制
function addCustomization() {
  customizations.value.push({ type: 'logo', is_enabled: true, note: '', note_trans: {} })
}
function removeCustomization(idx: number) {
  customizations.value.splice(idx, 1)
}

// 发布/下线/删除
async function handlePublish(row: Product) {
  const enName = (row as any).translations?.find((t: any) => t.language === 'en')?.name || ''
  const gate = checkProductGate({ ...row, name: row.name || enName, seo_title: row.seo?.title })
  if (!gate.ok) {
    await ElMessageBox.alert(gateAlertMessage(gate), '无法发布', { type: 'warning', confirmButtonText: '知道了' }).catch(() => {})
    return
  }
  await productApi.publish(row.id)
  ElMessage.success('已发布')
  loadData()
}

async function handleUnpublish(row: Product) {
  await productApi.unpublish(row.id)
  ElMessage.success('已下线')
  loadData()
}

async function handleDelete(row: Product) {
  await ElMessageBox.confirm(`确定删除产品 ${row.sku} 吗？此操作可通过回收站恢复。`, '确认删除', {
    type: 'warning', confirmButtonText: '确定删除', cancelButtonText: '取消',
  })
  await productApi.delete(row.id)
  ElMessage.success('已删除')
  loadData()
}

onMounted(() => {
  loadEnumDict()
  loadData()
  loadCategories()
  loadSeries()
  loadFabrics()
})


</script>

<style scoped>
.product-page { padding: 0; }
.stat-row { margin-bottom: 20px; }
.stat-card { border-radius: 12px; border: 1px solid #eef0f4; transition: all 0.25s; }
.stat-card:hover { border-color: #d1d5db; box-shadow: 0 2px 12px rgba(0,0,0,0.04); }
.stat-body { display: flex; align-items: center; justify-content: space-between; padding: 4px 0; }
.stat-label { font-size: 13px; color: #6b7280; margin-bottom: 6px; }
.stat-value { font-size: 28px; font-weight: 700; color: #1e293b; line-height: 1; }
.stat-value-success { color: #059669; }
.stat-value-warning { color: #d97706; }
.stat-value-danger { color: #dc2626; }
.stat-icon-box { width: 44px; height: 44px; border-radius: 12px; background: #eef2ff; color: #4f46e5; display: flex; align-items: center; justify-content: center; }
.stat-icon-green { background: #ecfdf5; color: #059669; }
.stat-icon-orange { background: #fffbeb; color: #d97706; }
.stat-icon-red { background: #fef2f2; color: #dc2626; }
.main-card { border-radius: 12px; border: 1px solid #eef0f4; }
.main-card :deep(.el-card__body) { padding: 20px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 12px; flex-wrap: wrap; }
.toolbar-left { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.toolbar-right { display: flex; align-items: center; gap: 8px; }
.product-table { border-radius: 8px; }
.product-table :deep(.el-table__header-wrapper) { border-bottom: 1px solid #eef0f4; }
.product-table :deep(.el-table__body tr:hover > td) { background-color: #f8faff; }
.product-table :deep(.el-table__cell) { padding: 10px 0; }
.cover-cell { display: flex; justify-content: center; }
.cover-thumb { width: 36px; height: 36px; border-radius: 6px; display: block; }
.cover-placeholder-cell { width: 36px; height: 36px; background: #f3f4f6; border-radius: 6px; display: flex; align-items: center; justify-content: center; font-size: 16px; }
.gender-tag { font-size: 12px; color: #6b7280; background: #f3f4f6; padding: 2px 8px; border-radius: 4px; }
.featured-star { color: #f59e0b; font-size: 18px; }
.action-btns { display: flex; gap: 6px; flex-wrap: nowrap; }
.action-btns .el-button { padding: 5px 10px; font-size: 12px; }
.pagination { margin-top: 20px; justify-content: flex-end; }
.pagination :deep(.el-pagination__total) { margin-right: auto; }

/* 抽屉 */
.product-drawer :deep(.el-drawer__body) { padding: 0; overflow: auto; }
.product-form { padding: 0; }
.drawer-tabs { border: none; }
.drawer-tabs :deep(.el-tabs__header) { background: #f8fafc; border-bottom: 1px solid #eef0f4; padding: 0 20px; margin: 0; }
.drawer-tabs :deep(.el-tabs__nav-wrap) { padding: 0; }
.drawer-tabs :deep(.el-tabs__item) { font-size: 13px; padding: 0 16px; height: 44px; line-height: 44px; }
.drawer-tabs :deep(.el-tabs__item.is-active) { font-weight: 600; color: #4f46e5; }
.drawer-tabs :deep(.el-tabs__active-bar) { background: #4f46e5; }
.tab-pane-content { padding: 24px 20px 20px; min-height: 300px; }
.section-subtitle { font-size: 14px; font-weight: 600; color: #374151; margin-bottom: 16px; }
.full-width { width: 100%; }
.dialog-footer { display: flex; justify-content: flex-end; gap: 12px; }
.translation-alert { margin-bottom: 12px; }

/* 图片 */
.img-mgr-section { margin-bottom: 28px; }
.img-mgr-section-title { font-size: 14px; font-weight: 600; color: #1e293b; margin-bottom: 14px; display: flex; align-items: center; gap: 8px; }
.cover-section { display: flex; align-items: center; gap: 16px; }
.cover-preview-lg { width: 160px; height: 160px; border-radius: 10px; border: 1px solid #e5e7eb; object-fit: cover; }
.cover-placeholder-lg { width: 160px; height: 160px; border: 2px dashed #d1d5db; border-radius: 10px; display: flex; align-items: center; justify-content: center; color: #9ca3af; font-size: 13px; background: #f9fafb; }
.cover-section-actions { display: flex; flex-direction: column; gap: 8px; }
.cover-tip { font-size: 12px; color: #909399; }
.image-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 12px; }
.image-item { position: relative; aspect-ratio: 1; border-radius: 8px; overflow: hidden; border: 1px solid #e5e7eb; transition: border-color 0.2s; }
.image-item:hover { border-color: #6366f1; }
.image-item-wrap { display: flex; flex-direction: column; gap: 4px; }
.image-thumb { width: 100%; height: 100%; object-fit: cover; }
.image-item-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.45); display: flex; align-items: center; justify-content: center; gap: 6px; opacity: 0; transition: opacity 0.25s; flex-wrap: wrap; padding: 0 4px; }
.image-item:hover .image-item-overlay { opacity: 1; }
.overlay-btn { border: none; background: rgba(255,255,255,0.2); backdrop-filter: blur(4px); }
.overlay-btn:hover { background: rgba(255,255,255,0.35); }
.image-alt-input { width: 100%; }
.image-upload-box { aspect-ratio: 1; border: 2px dashed #d1d5db; border-radius: 8px; display: flex; align-items: center; justify-content: center; cursor: pointer; transition: all 0.2s; background: #f9fafb; }
.image-upload-box:hover { border-color: #6366f1; background: #eef2ff; }
.upload-trigger { display: flex; flex-direction: column; align-items: center; gap: 8px; color: #9ca3af; font-size: 13px; }
.upload-progress { display: flex; align-items: center; gap: 14px; margin-top: 16px; padding: 12px 16px; background: #f0fdf4; border-radius: 8px; }

/* 规格 / 视频 / 定制 */
.spec-item { margin-bottom: 12px; padding: 10px; border: 1px solid #f0ede8; border-radius: 8px; }
.spec-row { display: flex; gap: 8px; align-items: center; margin-bottom: 10px; flex-wrap: wrap; }
.spec-trans { display: flex; gap: 8px; margin-top: 8px; flex-wrap: wrap; }
.spec-trans .el-input { flex: 1; min-width: 140px; }
.video-item { margin-bottom: 12px; padding: 10px; border: 1px solid #f0ede8; border-radius: 8px; }
.video-row { display: flex; gap: 8px; align-items: center; margin-bottom: 10px; flex-wrap: wrap; }
.video-cover { display: flex; align-items: center; gap: 8px; margin-top: 8px; }
.video-cover-label { font-size: 12px; color: #909399; flex-shrink: 0; }
.video-cover-picker { flex: 1; min-width: 0; }
.video-trans { display: flex; gap: 8px; margin-top: 8px; flex-wrap: wrap; }
.video-trans .el-input { flex: 1; min-width: 140px; }
.custom-item { margin-bottom: 12px; padding: 10px; border: 1px solid #f0ede8; border-radius: 8px; }
.custom-row { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; }
.custom-trans { display: flex; gap: 8px; margin-top: 8px; flex-wrap: wrap; }
.custom-trans .el-input { flex: 1; min-width: 140px; }
.mt-2 { margin-top: 8px; }

</style>
