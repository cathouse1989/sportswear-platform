<template>
  <div class="product-page">
    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stat-row">
      <el-col :span="6">
        <el-card shadow="never" class="stat-card">
          <div class="stat-body">
            <div class="stat-info">
              <div class="stat-label">产品总数</div>
              <div class="stat-value">{{ total }}</div>
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
          <el-input
            v-model="keyword"
            placeholder="搜索 SKU / Slug..."
            clearable
            style="width: 220px"
            @keyup.enter="handleSearch"
          />
          <el-select v-model="status" placeholder="状态" clearable style="width: 120px" @change="handleSearch">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
            <el-option label="已下线" value="offline" />
          </el-select>
          <el-select v-model="gender" placeholder="性别" clearable style="width: 110px" @change="handleSearch">
            <el-option label="中性" value="unisex" />
            <el-option label="男装" value="male" />
            <el-option label="女装" value="female" />
            <el-option label="童装" value="kids" />
          </el-select>
          <el-select v-model="typeFilter" placeholder="类型" clearable style="width: 110px" @change="handleSearch">
            <el-option label="OEM" value="oem" />
            <el-option label="ODM" value="odm" />
            <el-option label="Both" value="both" />
          </el-select>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
        </div>
        <div class="toolbar-right">
          <el-button type="primary" @click="openCreateDialog">
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
            <el-tag size="small" :type="row.type === 'oem' ? 'info' : row.type === 'odm' ? 'warning' : 'primary'">{{ row.type?.toUpperCase() }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="性别" width="66" align="center">
          <template #default="{ row }">
            <span class="gender-tag">{{ genderLabel(row.gender) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="精选" width="56" align="center">
          <template #default="{ row }">
            <span v-if="row.is_featured" class="featured-star">★</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small" effect="plain">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="production_moq" label="MOQ" width="72" align="right" />
        <el-table-column label="操作" width="340" fixed="right">
          <template #default="{ row }: any">
            <div class="action-btns">
              <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
              <el-button size="small" @click="openImageDialog(row)">图片</el-button>
              <el-button size="small" @click="openSpecDialog(row)">规格</el-button>
              <el-button size="small" @click="openVideoDialog(row)">视频</el-button>
              <el-dropdown trigger="click" @command="(cmd: string) => handleMoreAction(cmd, row)">
                <el-button size="small">
                  更多
                  <el-icon style="margin-left:2px"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="publish" v-if="row.status !== 'published'">发布</el-dropdown-item>
                    <el-dropdown-item command="unpublish" v-else>下线</el-dropdown-item>
                    <el-dropdown-item command="customize">定制选项</el-dropdown-item>
                    <el-dropdown-item command="delete" divided style="color:#ef4444">删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑产品' : '新建产品'" width="840px" top="4vh" class="product-dialog" destroy-on-close>
      <el-form :model="form" label-width="100px" class="product-form">
        <el-tabs v-model="activeTab" class="product-tabs" type="border-card">
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
                    <el-input v-model="form.name" placeholder="产品英文名称（发布必填，前台展示用）" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="类型">
                    <el-select v-model="form.type" class="full-width">
                      <el-option label="OEM" value="oem" />
                      <el-option label="ODM" value="odm" />
                      <el-option label="Both" value="both" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="性别">
                    <el-select v-model="form.gender" class="full-width">
                      <el-option label="中性" value="unisex" />
                      <el-option label="男装" value="male" />
                      <el-option label="女装" value="female" />
                      <el-option label="童装" value="kids" />
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
                <div class="cover-upload">
                  <el-image v-if="form.cover_image" :src="form.cover_image" fit="cover" class="cover-preview" />
                  <div v-else class="cover-placeholder">点击选择图片</div>
                  <div class="cover-actions">
                    <el-button size="small" @click="selectCoverImage">从媒体库选择</el-button>
                    <el-button v-if="form.cover_image" size="small" type="danger" plain @click="form.cover_image = ''">移除</el-button>
                  </div>
                </div>
              </el-form-item>
              <el-form-item label="简述">
                <el-input v-model="form.brief" type="textarea" :rows="2" placeholder="产品简短描述" />
              </el-form-item>
            </div>
          </el-tab-pane>

          <el-tab-pane label="详细信息" name="detail">
            <div class="tab-pane-content">
              <el-form-item label="描述">
                <el-input v-model="form.description" type="textarea" :rows="5" placeholder="产品详细描述" />
              </el-form-item>
              <el-form-item label="特性">
                <el-input v-model="form.features" type="textarea" :rows="4" placeholder="产品特性，每行一项" />
              </el-form-item>
              <el-form-item label="用途">
                <el-input v-model="form.usage" type="textarea" :rows="3" placeholder="产品用途说明" />
              </el-form-item>
            </div>
          </el-tab-pane>

          <el-tab-pane label="规格参数" name="specs">
            <div class="tab-pane-content">
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
                  <el-form-item label="重量">
                    <el-input v-model="form.weight" placeholder="如：180g/m²" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="弹性">
                    <el-input v-model="form.elasticity" placeholder="如：Medium" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="合身度">
                    <el-input v-model="form.fit" placeholder="如：Regular" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="支撑等级">
                    <el-input v-model="form.support_level" placeholder="如：High" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="季节">
                    <el-input v-model="form.season" placeholder="如：All Season" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="尺码范围">
                    <el-input v-model="form.size_range" placeholder="如：XS-3XL" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>

          <el-tab-pane label="MOQ" name="moq">
            <div class="tab-pane-content">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="样品 MOQ">
                    <el-input-number v-model="form.sample_moq" :min="0" class="full-width" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="大货 MOQ">
                    <el-input-number v-model="form.production_moq" :min="0" class="full-width" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="颜色 MOQ">
                    <el-input-number v-model="form.color_moq" :min="0" class="full-width" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="尺码 MOQ">
                    <el-input-number v-model="form.size_moq" :min="0" class="full-width" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
          </el-tab-pane>

          <el-tab-pane label="翻译" name="translations">
            <div class="tab-pane-content">
              <el-alert
                type="info"
                :closable="false"
                show-icon
                class="translation-alert"
                title="英文为源语言：上方「基本信息 / 详细信息」中的内容即英文源，保存时自动写入英文翻译；其他语言可一键同步英文内容后微调。"
              />
              <div v-for="(t, idx) in form.translations" :key="idx" class="translation-item">
                <div class="translation-header">
                  <span class="translation-lang-badge">{{ t.language === 'zh' ? '中文' : t.language === 'es' ? 'Español' : t.language === 'fr' ? 'Français' : 'English' }}</span>
                  <el-tag :type="t.name ? 'success' : 'info'" size="small" effect="plain">{{ t.name ? '已翻译' : '未翻译' }}</el-tag>
                  <el-button size="small" text type="primary" @click="syncTranslationFromEn(idx)">同步英文内容</el-button>
                  <span class="translation-sort-label">该语言排序</span>
                  <el-input-number v-model="t.sort_order" :min="0" size="small" controls-position="right" style="width: 110px" />
                  <el-tooltip content="该语言站点中的展示顺序（越小越靠前），0 = 跟随全局排序" placement="top">
                    <span class="translation-sort-help">?</span>
                  </el-tooltip>
                  <el-button v-if="t.language !== 'en'" size="small" type="danger" text @click="removeTranslation(idx)">移除</el-button>
                </div>
                <el-form-item :label="'名称'">
                  <el-input v-model="t.name" />
                </el-form-item>
                <el-row :gutter="12">
                  <el-col :span="12">
                    <el-form-item :label="'简述'">
                      <el-input v-model="t.brief" type="textarea" :rows="2" />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item :label="'特性'">
                      <el-input v-model="t.features" type="textarea" :rows="2" />
                    </el-form-item>
                  </el-col>
                </el-row>
                <el-form-item :label="'描述'">
                  <el-input v-model="t.description" type="textarea" :rows="3" />
                </el-form-item>
                <el-form-item :label="'用途'">
                  <el-input v-model="t.usage" type="textarea" :rows="2" />
                </el-form-item>
                <el-divider v-if="idx < form.translations.length - 1" />
              </div>
              <el-button size="small" @click="addTranslation" class="add-translation-btn">+ 添加语言翻译</el-button>
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
                  <el-form-item label="SEO Title">
                    <el-input v-model="form.seo.title" placeholder="搜索引擎标题" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Keywords">
                    <el-input v-model="form.seo.keywords" placeholder="关键词，逗号分隔" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="SEO Description">
                <el-input v-model="form.seo.description" type="textarea" :rows="2" placeholder="搜索引擎描述" />
              </el-form-item>
              <el-divider />
              <div class="section-subtitle">Open Graph</div>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="OG Title">
                    <el-input v-model="form.seo.og_title" placeholder="社交分享标题" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="OG Image">
                    <el-input v-model="form.seo.og_image" placeholder="社交分享图片 URL" />
                  </el-form-item>
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
          <el-button @click="dialogVisible = false" size="large">取消</el-button>
          <el-button type="primary" @click="handleSave" size="large" :loading="saving">保存产品</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 图片管理对话框 -->
    <el-dialog v-model="imageDialogVisible" title="产品图片管理" width="750px" top="5vh" class="image-dialog" destroy-on-close>
      <template v-if="imageProductId">
        <div class="img-mgr-section">
          <div class="img-mgr-section-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
            封面图片
          </div>
          <div class="cover-section">
            <el-image v-if="imageCoverUrl" :src="imageCoverUrl" fit="cover" class="cover-preview-lg" />
            <div v-else class="cover-placeholder-lg">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
              <span>暂无封面</span>
            </div>
            <div class="cover-section-actions">
              <el-upload :show-file-list="false" :before-upload="beforeCoverUpload" accept="image/*">
                <el-button size="small" type="primary">上传封面图片</el-button>
              </el-upload>
            </div>
          </div>
        </div>

        <div class="img-mgr-section">
          <div class="img-mgr-section-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zm10 0a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"/></svg>
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
                    <el-button size="small" circle type="danger" @click="deleteImage(img, idx)" class="overlay-btn">
                      <template #icon><el-icon color="#fff"><Delete /></el-icon></template>
                    </el-button>
                  </el-tooltip>
                </div>
                <div v-if="img.url === imageCoverUrl" class="image-cover-badge">封面</div>
              </div>
              <el-input v-model="img.alt" size="small" placeholder="alt 描述（SEO / 无障碍），填颜色名可联动颜色选择器" class="image-alt-input" />
            </div>
            <el-upload class="image-upload-box" :show-file-list="false" :before-upload="beforeImageUpload" accept="image/*" multiple>
              <div class="upload-trigger">
                <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 5v14m-7-7h14"/></svg>
                <span>上传图片</span>
              </div>
            </el-upload>
          </div>
        </div>

        <div v-if="uploading" class="upload-progress">
          <el-progress :percentage="uploadProgress" :stroke-width="6" status="success" />
          <span>正在上传 {{ uploadProgress }}%...</span>
        </div>
      </template>
      <template #footer>
        <el-button @click="imageDialogVisible = false" size="large">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 规格管理对话框 -->
    <el-dialog v-model="specDialogVisible" title="产品规格管理" width="580px" class="spec-dialog" destroy-on-close>
      <template v-if="specProductId">
        <el-alert type="info" :closable="false" show-icon style="margin-bottom:12px;" title="规格值支持用逗号 / 斜杠分隔多个可选项（如 Color: Black, White, Navy），门户会渲染为可点选的颜色/尺码选择器；单一值则展示在规格表格中。" />
        <div style="display:flex; align-items:center; gap:8px; flex-wrap:wrap; margin-bottom:12px;">
          <span style="font-size:13px; color:#6b7280;">快速添加：</span>
          <el-button size="small" @click="addSpecPreset('Color', 'Black, White, Navy, Red')">颜色 Color</el-button>
          <el-button size="small" @click="addSpecPreset('Size', 'XS, S, M, L, XL, XXL')">尺码 Size</el-button>
          <el-button size="small" @click="addSpecPreset('Fabric', 'Polyester, Cotton, Nylon')">面料 Fabric</el-button>
        </div>
        <div class="spec-list">
          <div v-for="(spec, idx) in specs" :key="idx" class="spec-row">
            <el-input v-model="spec.name" placeholder="规格名称，如：面料" style="width: 180px" />
            <el-input v-model="spec.value" placeholder="规格值，如：涤纶 / 多选项用逗号分隔" style="width: 220px" />
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
        </div>
        <el-button size="small" @click="addSpec" class="mt-2">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:4px"><path d="M12 5v14m-7-7h14"/></svg>
          添加规格
        </el-button>
      </template>
      <template #footer>
        <el-button @click="specDialogVisible = false" size="large">取消</el-button>
        <el-button type="primary" @click="saveSpecs" size="large" :loading="savingSpecs">保存规格</el-button>
      </template>
    </el-dialog>

    <!-- 视频管理对话框 -->
    <el-dialog v-model="videoDialogVisible" title="产品视频管理" width="680px" class="video-dialog" destroy-on-close>
      <template v-if="videoProductId">
        <div class="video-list">
          <div v-for="(v, idx) in videos" :key="idx" class="video-row">
            <el-select v-model="v.type" placeholder="类型" style="width: 120px">
              <el-option v-for="opt in videoTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
            <el-input v-model="v.title" placeholder="标题" style="width: 140px" />
            <el-input v-model="v.url" placeholder="视频 URL（支持 YouTube 链接）" style="width: 220px" />
            <el-input v-model="v.cover" placeholder="封面 URL" style="width: 140px" />
            <el-button type="danger" size="small" @click="removeVideo(idx)" circle>
              <template #icon><el-icon><Delete /></el-icon></template>
            </el-button>
          </div>
        </div>
        <el-button size="small" @click="addVideo" class="mt-2">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:4px"><path d="M12 5v14m-7-7h14"/></svg>
          添加视频
        </el-button>
      </template>
      <template #footer>
        <el-button @click="videoDialogVisible = false" size="large">取消</el-button>
        <el-button type="primary" @click="saveVideos" size="large" :loading="savingVideos">保存视频</el-button>
      </template>
    </el-dialog>

    <!-- 定制管理对话框 -->
    <el-dialog v-model="customDialogVisible" title="产品定制能力管理" width="680px" class="custom-dialog" destroy-on-close>
      <template v-if="customProductId">
        <div class="custom-list">
          <div v-for="(c, idx) in customizations" :key="idx" class="custom-row">
            <el-select v-model="c.type" placeholder="定制类型" style="width: 160px">
              <el-option v-for="opt in customTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
            <el-switch v-model="c.is_enabled" active-text="启用" style="width: 90px" />
            <el-input v-model="c.note" placeholder="备注说明" style="width: 180px" />
            <el-button type="danger" size="small" @click="removeCustomization(idx)" circle>
              <template #icon><el-icon><Delete /></el-icon></template>
            </el-button>
          </div>
        </div>
        <el-button size="small" @click="addCustomization" class="mt-2">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:4px"><path d="M12 5v14m-7-7h14"/></svg>
          添加定制选项
        </el-button>
      </template>
      <template #footer>
        <el-button @click="customDialogVisible = false" size="large">取消</el-button>
        <el-button type="primary" @click="saveCustomizations" size="large" :loading="savingCustom">保存定制</el-button>
      </template>
    </el-dialog>

    <!-- 媒体选择对话框 -->
    <el-dialog v-model="mediaDialogVisible" title="选择图片" width="780px" top="5vh" class="media-dialog" destroy-on-close>
      <div class="media-select-container">
        <div class="media-select-toolbar">
          <span class="media-select-hint">选择已有图片或上传新图片</span>
          <el-upload :show-file-list="false" :before-upload="handleMediaUpload" accept="image/*">
            <el-button size="small" type="primary">上传新图片</el-button>
          </el-upload>
        </div>
        <div v-loading="mediaLoading" class="media-grid">
          <div v-for="m in mediaList" :key="m.id" class="media-item" :class="{ 'media-item-selected': selectedMediaId === m.id }" @click="selectedMediaId = m.id">
            <el-image :src="m.url" fit="cover" class="media-thumb" />
            <div class="media-item-check" v-if="selectedMediaId === m.id">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#fff" stroke-width="3"><path d="M5 13l4 4L19 7"/></svg>
            </div>
            <div class="media-item-label">{{ m.original_name?.slice(0, 20) }}</div>
          </div>
        </div>
        <div v-if="!mediaLoading && !mediaList.length" class="media-empty">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#d0d5dd" stroke-width="1.5"><path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
          <p>暂无媒体文件，请上传</p>
        </div>
        <el-pagination v-if="mediaTotal > mediaPageSize" v-model:current-page="mediaPage" :page-size="mediaPageSize" :total="mediaTotal" layout="prev, pager, next" small @current-change="loadMedia" class="media-pagination" />
      </div>
      <template #footer>
        <el-button @click="mediaDialogVisible = false" size="large">取消</el-button>
        <el-button type="primary" :disabled="!selectedMediaId" @click="confirmMediaSelect" size="large">选择图片</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Star, Delete, ArrowDown, ArrowUp, Search } from '@element-plus/icons-vue'
import { productApi, mediaApi, categoryApi, seriesApi, fabricApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import type { Product, Media, Category, Series, Fabric } from '@/types'
import { checkProductGate, gateAlertMessage } from '@/utils/publish-gate'

const products = ref<Product[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const status = ref('')
const gender = ref('')
const typeFilter = ref('')
const featuredFilter = ref('')

// 状态统计卡：来自后端全量统计接口，不受列表筛选/分页影响
const publishedCount = ref(0)
const draftCount = ref(0)
const offlineCount = ref(0)

async function loadProductStats() {
  try {
    const s = await productApi.stats()
    publishedCount.value = s.published
    draftCount.value = s.draft
    offlineCount.value = s.offline
  } catch {}
}

const dialogVisible = ref(false)
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
  translations: [] as Array<{ language: string; name: string; brief: string; description: string; features: string; usage: string; sort_order: number }>,
  seo: { title: '', description: '', keywords: '', og_title: '', og_description: '', og_image: '' },
  series_ids: [] as string[],
  fabric_ids: [] as string[],
})

// 图片管理
const imageDialogVisible = ref(false)
const imageProductId = ref('')
const imageCoverUrl = ref('')
const productImages = ref<Array<{ id: string; url: string; alt: string; sort_order: number }>>([])
// 打开对话框时拉取的完整产品详情快照（用于整体 PUT 时补齐标量字段）
const imageDetail = ref<any>(null)
const uploading = ref(false)
const uploadProgress = ref(0)

// 规格管理
const specDialogVisible = ref(false)
const specProductId = ref('')
const specs = ref<Array<{ name: string; value: string }>>([])
const specDetail = ref<any>(null)
const savingSpecs = ref(false)

// 视频管理
const videoDialogVisible = ref(false)
const videoProductId = ref('')
const videos = ref<Array<{ type: string; title: string; url: string; cover: string }>>([])
const videoDetail = ref<any>(null)
const savingVideos = ref(false)
const videoTypeOptions = [
  { label: '产品视频', value: 'product' },
  { label: '生产视频', value: 'production' },
  { label: '工艺视频', value: 'craft' },
  { label: '用途视频', value: 'usage' },
]

// 定制管理
const customDialogVisible = ref(false)
const customProductId = ref('')
const customizations = ref<Array<{ type: string; is_enabled: boolean; note: string }>>([])
const customDetail = ref<any>(null)
const savingCustom = ref(false)
const customTypeOptions = [
  { label: 'Logo 定制', value: 'logo' },
  { label: '颜色定制', value: 'color' },
  { label: '面料定制', value: 'fabric' },
  { label: '图案定制', value: 'pattern' },
  { label: '印花定制', value: 'print' },
  { label: '刺绣定制', value: 'embroidery' },
  { label: '标签定制', value: 'label' },
  { label: '吊牌定制', value: 'hangtag' },
  { label: '包装定制', value: 'packaging' },
  { label: '尺码定制', value: 'size' },
  { label: '合身度定制', value: 'fit' },
  { label: '拉链定制', value: 'zipper' },
  { label: '纽扣定制', value: 'button' },
  { label: '腰带定制', value: 'belt' },
  { label: '配件定制', value: 'accessory' },
]

// 媒体选择
const mediaDialogVisible = ref(false)
const mediaList = ref<Media[]>([])
const mediaLoading = ref(false)
const mediaPage = ref(1)
const mediaPageSize = ref(20)
const mediaTotal = ref(0)
const selectedMediaId = ref('')
let mediaResolve: ((url: string) => void) | null = null

function genderLabel(g: string) {
  return g === 'unisex' ? '中性' : g === 'male' ? '男装' : g === 'female' ? '女装' : '童装'
}

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
    if (featuredFilter.value) params.is_featured = featuredFilter.value

    const result = await productApi.list(params)
    products.value = result.items
    total.value = result.total
  } finally {
    loading.value = false
  }
  loadProductStats()
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
    translations: [],
    seo: { title: '', description: '', keywords: '', og_title: '', og_description: '', og_image: '' },
    series_ids: [],
    fabric_ids: [],
  })
  activeTab.value = 'basic'
}

function openCreateDialog() {
  editingId.value = ''
  resetForm()
  dialogVisible.value = true
}

async function openEditDialog(row: Product) {
  editingId.value = row.id
  resetForm()

  try {
    const detail = await productApi.get(row.id)
    Object.assign(form, {
      sku: detail.sku,
      slug: detail.slug,
      name: detail.translations?.find((t: any) => t.language === 'en')?.name || '',
      category_id: detail.category_id || '',
      type: detail.type,
      gender: detail.gender,
      status: detail.status,
      is_featured: detail.is_featured,
      is_new: detail.is_new,
      sort_order: detail.sort_order,
      cover_image: detail.cover_image || '',
      brief: detail.brief || '',
      description: detail.description || '',
      features: detail.features || '',
      usage: detail.usage || '',
      material: detail.material || '',
      composition: detail.composition || '',
      weight: detail.weight || '',
      elasticity: detail.elasticity || '',
      fit: detail.fit || '',
      support_level: detail.support_level || '',
      season: detail.season || '',
      size_range: detail.size_range || '',
      sample_moq: detail.sample_moq || 1,
      production_moq: detail.production_moq || 300,
      color_moq: detail.color_moq || 100,
      size_moq: detail.size_moq || 100,
      translations: (detail.translations || []).filter((t: any) => t.language !== 'en').map((t: any) => ({
        language: t.language,
        name: t.name,
        brief: t.brief || '',
        description: t.description || '',
        features: t.features || '',
        usage: t.usage || '',
        sort_order: t.sort_order || 0,
      })),
      series_ids: (detail.series || []).map((s: any) => s.id),
      fabric_ids: (detail.fabrics || []).map((f: any) => f.id),
    })
    if (detail.seo) {
      form.seo = {
        title: detail.seo.title || '',
        description: detail.seo.description || '',
        keywords: detail.seo.keywords || '',
        og_title: detail.seo.og_title || '',
        og_description: detail.seo.og_description || '',
        og_image: detail.seo.og_image || '',
      }
    }
  } catch {
    Object.assign(form, {
      sku: row.sku,
      slug: row.slug,
      type: row.type,
      gender: row.gender,
      status: row.status,
      is_featured: row.is_featured,
      is_new: row.is_new,
      cover_image: row.cover_image || '',
      brief: row.brief || '',
      production_moq: row.production_moq || 300,
    })
  }
  dialogVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    const data = { ...form }
    // 产品主表无 name 列，英文名存储于 product_translations(language=en)。
    // 提交时自动并入 en 翻译，保证名称可被门户解析、发布质检可通过。
    const translations = [...form.translations]
    if (String(form.name || '').trim()) {
      const enEntry = {
        language: 'en',
        name: String(form.name).trim(),
        brief: form.brief || '',
        description: form.description || '',
        features: form.features || '',
        usage: form.usage || '',
        // en 为源语言，语言维度排序 0 = 跟随全局 sort_order
        sort_order: 0,
      }
      const enIdx = translations.findIndex(t => t.language === 'en')
      if (enIdx >= 0) translations[enIdx] = enEntry
      else translations.unshift(enEntry)
    }
    data.translations = translations
    if (editingId.value) {
      await productApi.update(editingId.value, data)
    } else {
      await productApi.create(data)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadData()
  } catch {
    // 错误已提示
  } finally {
    saving.value = false
  }
}

function handleMoreAction(cmd: string, row: Product) {
  if (cmd === 'publish') handlePublish(row)
  else if (cmd === 'unpublish') handleUnpublish(row)
  else if (cmd === 'customize') openCustomizationDialog(row)
  else if (cmd === 'delete') handleDelete(row)
}

// 翻译管理
// 英文为源语言：主表单「基本信息/详细信息」中的内容即英文源，
// 保存时自动写入 en 翻译（见 handleSave）；其他语言从这里同步后微调。
function addTranslation() {
  const langs = ['zh', 'es', 'fr'].filter(l => !form.translations.find(t => t.language === l))
  if (langs.length === 0) {
    ElMessage.warning('所有语言已添加')
    return
  }
  // 新增语言时自动同步英文源内容，避免从零录入
  form.translations.push({
    language: langs[0],
    name: form.name,
    brief: form.brief,
    description: form.description,
    features: form.features,
    usage: form.usage,
    sort_order: 0,
  })
}

// 一键同步英文源到指定语言词条（同步后可人工微调）
function syncTranslationFromEn(idx: number) {
  const t = form.translations[idx]
  if (!t) return
  t.name = form.name
  t.brief = form.brief
  t.description = form.description
  t.features = form.features
  t.usage = form.usage
  ElMessage.success('已同步英文内容，可在此基础上修改')
}

function removeTranslation(idx: number) {
  form.translations.splice(idx, 1)
}

// 图片管理
// 构建完整的产品 PUT payload：子资源保存时携带全部标量字段，
// 保证「整体 PUT」语义（后端 UpdateProduct 全量更新标量）。
function productPayload(detail: any): Record<string, any> {
  return {
    sku: detail?.sku,
    slug: detail?.slug,
    category_id: detail?.category_id || null,
    type: detail?.type || 'both',
    gender: detail?.gender || 'unisex',
    status: detail?.status || 'draft',
    is_featured: detail?.is_featured ?? false,
    is_new: detail?.is_new ?? false,
    sort_order: detail?.sort_order || 0,
    cover_image: detail?.cover_image || '',
    brief: detail?.brief || '',
    description: detail?.description || '',
    features: detail?.features || '',
    usage: detail?.usage || '',
    material: detail?.material || '',
    composition: detail?.composition || '',
    weight: detail?.weight || '',
    elasticity: detail?.elasticity || '',
    fit: detail?.fit || '',
    support_level: detail?.support_level || '',
    season: detail?.season || '',
    size_range: detail?.size_range || '',
    sample_moq: detail?.sample_moq || 1,
    production_moq: detail?.production_moq || 300,
    color_moq: detail?.color_moq || 100,
    size_moq: detail?.size_moq || 100,
    translations: (detail?.translations || []).map((t: any) => ({
      language: t.language,
      name: t.name,
      brief: t.brief || '',
      description: t.description || '',
      features: t.features || '',
      usage: t.usage || '',
    })),
    seo: detail?.seo
      ? {
          language: detail.seo.language || 'en',
          title: detail.seo.title || '',
          description: detail.seo.description || '',
          keywords: detail.seo.keywords || '',
          canonical: detail.seo.canonical || '',
          robots: detail.seo.robots || '',
          og_title: detail.seo.og_title || '',
          og_description: detail.seo.og_description || '',
          og_image: detail.seo.og_image || '',
          schema_data: detail.seo.schema_data || '',
        }
      : null,
    series_ids: (detail?.series || []).map((s: any) => s.id),
    fabric_ids: (detail?.fabrics || []).map((f: any) => f.id),
  }
}

// 整体 PUT 保存产品图片列表（连同全部标量字段）
async function saveProductImages() {
  if (!imageProductId.value || !imageDetail.value) return
  const images = productImages.value.map((img, i) => ({
    type: 'gallery',
    url: img.url,
    thumbnail: (img as any).thumbnail || '',
    alt: (img as any).alt || '',
    sort_order: i,
  }))
  await productApi.update(imageProductId.value, { ...productPayload(imageDetail.value), images })
}

async function openImageDialog(row: Product) {
  imageProductId.value = row.id
  imageDetail.value = row
  imageCoverUrl.value = row.cover_image || ''
  productImages.value = []
  uploading.value = false
  uploadProgress.value = 0
  imageDialogVisible.value = true

  try {
    const detail = await productApi.get(row.id)
    imageDetail.value = detail
    if (detail.images?.length) {
      productImages.value = detail.images
        .filter(i => i.type !== 'cover')
        .map(i => ({ id: i.id, url: i.url, alt: (i as any).alt || '', sort_order: i.sort_order }))
      if (!imageCoverUrl.value) {
        const cover = detail.images.find(i => i.type === 'cover')
        if (cover) imageCoverUrl.value = cover.url
      }
    }
  } catch {}
}

async function selectCoverImage() {
  const url = await openMediaSelector()
  if (url) form.cover_image = url
}

async function beforeCoverUpload(file: File) {
  try {
    const formData = new FormData()
    formData.append('file', file)
    const result = await mediaApi.upload(formData)
    const url = result.url
    if (imageProductId.value) {
      await productApi.update(imageProductId.value, { ...productPayload(imageDetail.value), cover_image: url })
    }
    imageCoverUrl.value = url
    ElMessage.success('封面上传成功')
    loadData()
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
    if (imageProductId.value) {
      productImages.value.push({ id: Date.now().toString(), url: result.url, alt: '', sort_order: productImages.value.length })
      await saveProductImages()
    }
    uploadProgress.value = 100
    ElMessage.success('图片上传成功')
    loadData()
  } catch {
    ElMessage.error('图片上传失败')
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
  return false
}

async function setAsCover(img: { id: string; url: string }) {
  try {
    if (imageProductId.value) {
      await productApi.update(imageProductId.value, { ...productPayload(imageDetail.value), cover_image: img.url })
    }
    imageCoverUrl.value = img.url
    ElMessage.success('已设为封面')
    loadData()
  } catch {
    ElMessage.error('设置封面失败')
  }
}

async function deleteImage(img: { id: string; url: string }, idx: number) {
  try {
    await ElMessageBox.confirm('确定删除该图片吗？', '警告', { type: 'warning' })
    productImages.value.splice(idx, 1)
    await saveProductImages()
    if (imageCoverUrl.value === img.url) {
      imageCoverUrl.value = productImages.value[0]?.url || ''
    }
    ElMessage.success('图片已删除')
    loadData()
  } catch {}
}

// 图片排序：dir -1 前移 / 1 后移，保存后按数组顺序持久化 sort_order，
// 门户详情页的图集 / 浮动查看器即按此顺序展示
async function moveImage(idx: number, dir: number) {
  const target = idx + dir
  if (target < 0 || target >= productImages.value.length) return
  const list = productImages.value
  const [item] = list.splice(idx, 1)
  list.splice(target, 0, item)
  try {
    await saveProductImages()
    ElMessage.success('顺序已更新')
  } catch {
    ElMessage.error('顺序保存失败')
  }
}

// 规格管理
async function openSpecDialog(row: Product) {
  specProductId.value = row.id
  specDetail.value = row
  specs.value = []
  specDialogVisible.value = true
  try {
    const detail = await productApi.get(row.id)
    specDetail.value = detail
    if (detail.specs?.length) {
      specs.value = detail.specs.map(s => ({ name: s.name, value: s.value }))
    }
  } catch {}
}

function addSpec() {
  specs.value.push({ name: '', value: '' })
}

function addSpecPreset(name: string, value: string) {
  specs.value.push({ name, value })
}

function removeSpec(idx: number) {
  specs.value.splice(idx, 1)
}

// 规格排序：dir -1 上移 / 1 下移，保存时按数组顺序持久化 sort_order
function moveSpec(idx: number, dir: number) {
  const target = idx + dir
  if (target < 0 || target >= specs.value.length) return
  const list = specs.value
  const [item] = list.splice(idx, 1)
  list.splice(target, 0, item)
}

async function saveSpecs() {
  savingSpecs.value = true
  try {
    await productApi.update(specProductId.value, { ...productPayload(specDetail.value), specs: specs.value.map((s, i) => ({ name: s.name, value: s.value, sort_order: i })) })
    ElMessage.success('规格保存成功')
    specDialogVisible.value = false
    loadData()
  } catch {
    ElMessage.error('规格保存失败')
  } finally {
    savingSpecs.value = false
  }
}

// 视频管理
async function openVideoDialog(row: Product) {
  videoProductId.value = row.id
  videoDetail.value = row
  videos.value = []
  videoDialogVisible.value = true
  try {
    const detail = await productApi.get(row.id)
    videoDetail.value = detail
    if (detail.videos?.length) {
      videos.value = detail.videos.map((v: any) => ({
        type: v.type || 'product',
        title: v.title || '',
        url: v.url || '',
        cover: v.cover || '',
      }))
    }
  } catch {}
}

function addVideo() {
  videos.value.push({ type: 'product', title: '', url: '', cover: '' })
}

function removeVideo(idx: number) {
  videos.value.splice(idx, 1)
}

async function saveVideos() {
  const valid = videos.value.every((v) => String(v.url || '').trim() !== '')
  if (!valid) {
    ElMessage.error('每个视频都必须填写 URL')
    return
  }
  savingVideos.value = true
  try {
    await productApi.update(videoProductId.value, {
      ...productPayload(videoDetail.value),
      videos: videos.value.map((v, i) => ({ type: v.type, url: v.url, cover: v.cover, title: v.title, sort_order: i })),
    })
    ElMessage.success('视频保存成功')
    videoDialogVisible.value = false
    loadData()
  } catch {
    ElMessage.error('视频保存失败')
  } finally {
    savingVideos.value = false
  }
}

// 定制管理
async function openCustomizationDialog(row: Product) {
  customProductId.value = row.id
  customDetail.value = row
  customizations.value = []
  customDialogVisible.value = true
  try {
    const detail = await productApi.get(row.id)
    customDetail.value = detail
    if (detail.customizations?.length) {
      customizations.value = detail.customizations.map(c => ({
        type: c.type,
        is_enabled: c.is_enabled,
        note: c.note || '',
      }))
    }
  } catch {}
}

function addCustomization() {
  customizations.value.push({ type: 'logo', is_enabled: true, note: '' })
}

function removeCustomization(idx: number) {
  customizations.value.splice(idx, 1)
}

async function saveCustomizations() {
  savingCustom.value = true
  try {
    await productApi.update(customProductId.value, { ...productPayload(customDetail.value), customizations: customizations.value })
    ElMessage.success('定制选项保存成功')
    customDialogVisible.value = false
    loadData()
  } catch {
    ElMessage.error('定制选项保存失败')
  } finally {
    savingCustom.value = false
  }
}

// 媒体选择器
function openMediaSelector(): Promise<string> {
  return new Promise((resolve) => {
    mediaResolve = resolve
    selectedMediaId.value = ''
    mediaPage.value = 1
    mediaList.value = []
    mediaDialogVisible.value = true
    loadMedia()
  })
}

async function loadMedia() {
  mediaLoading.value = true
  try {
    const result = await mediaApi.list({ page: mediaPage.value, pageSize: mediaPageSize.value, type: 'image' })
    mediaList.value = result.items
    mediaTotal.value = result.total
  } catch {}
  finally { mediaLoading.value = false }
}

async function handleMediaUpload(file: File) {
  try {
    const formData = new FormData()
    formData.append('file', file)
    await mediaApi.upload(formData)
    ElMessage.success('上传成功')
    loadMedia()
  } catch {
    ElMessage.error('上传失败')
  }
  return false
}

function confirmMediaSelect() {
  const media = mediaList.value.find(m => m.id === selectedMediaId.value)
  if (media && mediaResolve) {
    mediaResolve(media.url)
    mediaResolve = null
  }
  mediaDialogVisible.value = false
}

async function handlePublish(row: Product) {
  // 发布质检门（P0-#2）：必填项缺失时拦截并列出缺失项。
  // 注意：admin 列表接口不本地化，row.name 恒为空，英文名须从 en 翻译解析。
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
    type: 'warning',
    confirmButtonText: '确定删除',
    cancelButtonText: '取消',
  })
  await productApi.delete(row.id)
  ElMessage.success('已删除')
  loadData()
}

function statusType(s: string) {
  return s === 'published' ? 'success' : s === 'draft' ? 'info' : 'warning'
}
function statusLabel(s: string) {
  return s === 'published' ? '已发布' : s === 'draft' ? '草稿' : '已下线'
}

onMounted(() => {
  loadData()
  loadCategories()
  loadSeries()
  loadFabrics()
})
</script>

<style scoped>
/* ============ 页面布局 ============ */
.product-page {
  padding: 0;
}

/* ============ 统计卡片 ============ */
.stat-row {
  margin-bottom: 20px;
}
.stat-card {
  border-radius: 12px;
  border: 1px solid #eef0f4;
  transition: all 0.25s;
}
.stat-card:hover {
  border-color: #d1d5db;
  box-shadow: 0 2px 12px rgba(0,0,0,0.04);
}
.stat-body {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
}
.stat-label {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 6px;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  line-height: 1;
}
.stat-value-success { color: #059669; }
.stat-value-warning { color: #d97706; }
.stat-value-danger { color: #dc2626; }
.stat-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #eef2ff;
  color: #4f46e5;
  display: flex;
  align-items: center;
  justify-content: center;
}
.stat-icon-green {
  background: #ecfdf5;
  color: #059669;
}
.stat-icon-orange {
  background: #fffbeb;
  color: #d97706;
}
.stat-icon-red {
  background: #fef2f2;
  color: #dc2626;
}

/* ============ 主卡片 ============ */
.main-card {
  border-radius: 12px;
  border: 1px solid #eef0f4;
}
.main-card :deep(.el-card__body) {
  padding: 20px;
}

/* ============ 工具栏 ============ */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 12px;
  flex-wrap: wrap;
}
.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* ============ 表格 ============ */
.product-table {
  border-radius: 8px;
}
.product-table :deep(.el-table__header-wrapper) {
  border-bottom: 1px solid #eef0f4;
}
.product-table :deep(.el-table__body tr:hover > td) {
  background-color: #f8faff;
}
.product-table :deep(.el-table__cell) {
  padding: 10px 0;
}
/* 封面上传 */
.cover-cell {
  display: flex;
  justify-content: center;
}
.cover-thumb {
  width: 36px;
  height: 36px;
  border-radius: 6px;
  display: block;
}
.cover-placeholder-cell {
  width: 36px;
  height: 36px;
  background: #f3f4f6;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}
.gender-tag {
  font-size: 12px;
  color: #6b7280;
  background: #f3f4f6;
  padding: 2px 8px;
  border-radius: 4px;
}
.featured-star {
  color: #f59e0b;
  font-size: 18px;
}
.action-btns {
  display: flex;
  gap: 6px;
  flex-wrap: nowrap;
}
.action-btns .el-button {
  padding: 5px 10px;
  font-size: 12px;
}

/* ============ 分页 ============ */
.pagination {
  margin-top: 20px;
  justify-content: flex-end;
}
.pagination :deep(.el-pagination__total) {
  margin-right: auto;
}

/* ============ 产品编辑对话框 ============ */
.product-dialog :deep(.el-dialog__body) {
  padding: 0;
}
.product-tabs {
  border: none;
}
.product-tabs :deep(.el-tabs__header) {
  background: #f8fafc;
  border-bottom: 1px solid #eef0f4;
  padding: 0 20px;
  margin: 0;
}
.product-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0;
}
.product-tabs :deep(.el-tabs__item) {
  font-size: 13px;
  padding: 0 16px;
  height: 44px;
  line-height: 44px;
}
.product-tabs :deep(.el-tabs__item.is-active) {
  font-weight: 600;
  color: #4f46e5;
}
.product-tabs :deep(.el-tabs__active-bar) {
  background: #4f46e5;
}
.tab-pane-content {
  padding: 24px 20px 20px;
  min-height: 300px;
}
.section-subtitle {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 16px;
}
.full-width {
  width: 100%;
}
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 封面上传 */
.cover-upload {
  display: flex;
  align-items: center;
  gap: 16px;
}
.cover-preview {
  width: 120px;
  height: 120px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  object-fit: cover;
}
.cover-placeholder {
  width: 120px;
  height: 120px;
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  background: #f9fafb;
}
.cover-placeholder:hover {
  border-color: #6366f1;
  color: #6366f1;
  background: #eef2ff;
}
.cover-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 翻译 */
.translation-alert {
  margin-bottom: 12px;
}
.translation-item {
  background: #f9fafb;
  border: 1px solid #eef0f4;
  border-radius: 10px;
  padding: 16px;
  margin-bottom: 12px;
}
.translation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}
.translation-lang-badge {
  font-size: 12px;
  font-weight: 600;
  background: #eef2ff;
  color: #4f46e5;
  padding: 3px 10px;
  border-radius: 6px;
}
.translation-sort-label {
  font-size: 12px;
  color: #6b7280;
  margin-left: auto;
}
.translation-sort-help {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #e5e7eb;
  color: #6b7280;
  font-size: 11px;
  cursor: help;
}
.add-translation-btn {
  margin-top: 8px;
}

/* 规格/定制 */
.spec-list, .custom-list, .video-list {
  margin-bottom: 8px;
}
.spec-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 10px;
}
.video-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 10px;
  flex-wrap: wrap;
}
.custom-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 10px;
  flex-wrap: wrap;
}
.mt-2 {
  margin-top: 8px;
}

/* ============ 图片管理对话框 ============ */
.image-dialog :deep(.el-dialog__body) {
  padding: 24px;
}
.img-mgr-section {
  margin-bottom: 28px;
}
.img-mgr-section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.img-mgr-section-title svg {
  color: #6b7280;
}
.cover-section {
  display: flex;
  align-items: center;
  gap: 16px;
}
.cover-preview-lg {
  width: 160px;
  height: 160px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  object-fit: cover;
}
.cover-placeholder-lg {
  width: 160px;
  height: 160px;
  border: 2px dashed #d1d5db;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  gap: 8px;
  font-size: 13px;
  background: #f9fafb;
}
.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}
.image-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e5e7eb;
  transition: border-color 0.2s;
}
.image-item:hover {
  border-color: #6366f1;
}
.image-item-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.image-alt-input {
  width: 100%;
}
.image-item-wrap .image-item-overlay {
  gap: 6px;
  flex-wrap: wrap;
  padding: 0 4px;
}
.image-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.image-item-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  opacity: 0;
  transition: opacity 0.25s;
}
.image-item:hover .image-item-overlay {
  opacity: 1;
}
.overlay-btn {
  border: none;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(4px);
}
.overlay-btn:hover {
  background: rgba(255,255,255,0.35);
}
.image-cover-badge {
  position: absolute;
  top: 6px;
  left: 6px;
  background: #f59e0b;
  color: #fff;
  font-size: 10px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
}
.image-upload-box {
  aspect-ratio: 1;
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  background: #f9fafb;
}
.image-upload-box:hover {
  border-color: #6366f1;
  background: #eef2ff;
}
.upload-trigger {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: #9ca3af;
  font-size: 13px;
}
.upload-progress {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-top: 16px;
  padding: 12px 16px;
  background: #f0fdf4;
  border-radius: 8px;
}
.upload-progress span {
  font-size: 13px;
  color: #059669;
  white-space: nowrap;
}

/* ============ 媒体选择器 ============ */
.media-select-container {
  min-height: 320px;
}
.media-select-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.media-select-hint {
  font-size: 13px;
  color: #6b7280;
}
.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}
.media-item {
  cursor: pointer;
  border: 2px solid transparent;
  border-radius: 8px;
  overflow: hidden;
  transition: border-color 0.2s;
  position: relative;
}
.media-item:hover {
  border-color: #6366f1;
}
.media-item-selected {
  border-color: #6366f1;
}
.media-item-check {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 22px;
  height: 22px;
  background: #6366f1;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.media-thumb {
  width: 100%;
  aspect-ratio: 1;
  object-fit: cover;
  display: block;
}
.media-item-label {
  font-size: 11px;
  color: #6b7280;
  padding: 4px 6px;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  background: #f9fafb;
}
.media-empty {
  text-align: center;
  padding: 48px 20px;
  color: #9ca3af;
}
.media-empty svg {
  margin-bottom: 12px;
}
.media-empty p {
  margin: 0;
  font-size: 14px;
}
.media-pagination {
  justify-content: center;
}
</style>