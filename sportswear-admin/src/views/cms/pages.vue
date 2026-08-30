<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-input
        v-model="keyword"
        placeholder="搜索标题 / Slug"
        clearable
        style="width: 240px"
        @keyup.enter="handleSearch"
      />
      <el-select v-model="status" placeholder="状态" clearable style="width: 150px" @change="handleSearch">
        <el-option label="草稿" value="draft" />
        <el-option label="已发布" value="published" />
        <el-option label="已下线" value="offline" />
      </el-select>
      <el-button type="primary" @click="handleSearch">查询</el-button>
      <div class="spacer" />
      <el-button type="primary" @click="openCreateDialog">新建页面</el-button>
    </div>

    <el-table :data="pages" v-loading="loading" stripe>
      <el-table-column prop="title" label="标题" min-width="160" />
      <el-table-column prop="slug" label="Slug" min-width="140" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }"><el-tag size="small">{{ row.type }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="360">
        <template #default="{ row }: any">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="warning" @click="$router.push('/hero')" v-if="isHome(row)">轮播图</el-button>
          <el-button v-permission="'page:publish'" size="small" type="success" @click="handlePublish(row)" v-if="row.status !== 'published'">发布</el-button>
          <el-button v-permission="'page:publish'" size="small" type="warning" @click="handleUnpublish(row)" v-else>下线</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑页面' : '新建页面'" width="560px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题" required>
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="Slug" required>
          <el-input v-model="form.slug" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type">
            <el-option label="普通" value="normal" />
            <el-option label="首页" value="home" />
            <el-option label="产品" value="product" />
            <el-option label="OEM" value="oem" />
            <el-option label="ODM" value="odm" />
            <el-option label="工厂" value="factory" />
            <el-option label="FAQ" value="faq" />
            <el-option label="联系我们" value="contact" />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { pageApi } from '@/api'
import { useAdminPageSize } from '@/composables/useAdminPageSize'
import type { Page } from '@/types'

const pages = ref<Page[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = useAdminPageSize()
const keyword = ref('')
const status = ref('')

const dialogVisible = ref(false)
const editingId = ref('')
const form = reactive({ title: '', slug: '', type: 'normal', sort_order: 0 })

async function loadData() {
  loading.value = true
  try {
    const result = await pageApi.list({ page: page.value, pageSize: pageSize.value, keyword: keyword.value, status: status.value })
    pages.value = (result as any).items || result || []
    total.value = (result as any).total ?? pages.value.length
  } finally { loading.value = false }
}
function handleSearch() { page.value = 1; loadData() }
function openCreateDialog() { editingId.value = ''; Object.assign(form, { title: '', slug: '', type: 'normal', sort_order: 0 }); dialogVisible.value = true }
function openEditDialog(row: Page) { editingId.value = row.id; Object.assign(form, { title: row.title, slug: row.slug, type: row.type, sort_order: row.sort_order }); dialogVisible.value = true }
async function handleSave() {
  try {
    if (editingId.value) { await pageApi.update(editingId.value, form) } else { await pageApi.create(form) }
    ElMessage.success('保存成功'); dialogVisible.value = false; loadData()
  } catch { /* error handled */ }
}
async function handlePublish(row: Page) { await pageApi.publish(row.id); ElMessage.success('已发布'); loadData() }
async function handleUnpublish(row: Page) { await pageApi.unpublish(row.id); ElMessage.success('已下线'); loadData() }
async function handleDelete(row: Page) { await ElMessageBox.confirm(`确定删除页面 ${row.title} 吗？`, '警告', { type: 'warning' }); await pageApi.delete(row.id); ElMessage.success('已删除'); loadData() }

const isHome = (row: Page) => row.slug === 'home' || row.type === 'home'

onMounted(loadData)
</script>
<style scoped>
.toolbar { display: flex; gap: 12px; margin-bottom: 16px; align-items: center; }
.spacer { flex: 1; }
.pagination { margin-top: 16px; justify-content: flex-end; }

</style>