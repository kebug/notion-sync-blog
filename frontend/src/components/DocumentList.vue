<template>
  <div class="document-list">
    <n-card>
      <n-data-table
        :columns="columns"
        :data="documents"
        :loading="loading"
        :pagination="pagination"
      />
      <n-empty v-if="!loading && documents.length === 0" description="暂无文档" class="mt-4" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, h } from 'vue'
import { NDataTable, NCard, NEmpty, NButton, NTag, type DataTableColumns } from 'naive-ui'
import { documentApi } from '../api/documents'
import type { Document } from '../types'

interface Props {
  projectId: number
}

const props = defineProps<Props>()

// 响应式数据
const documents = ref<Document[]>([])
const loading = ref(false)

const pagination = {
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    fetchDocuments({ offset: (page - 1) * pagination.pageSize })
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    fetchDocuments({ offset: 0, limit: pageSize })
  },
}

// 表格列定义
const columns: DataTableColumns<Document> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '页面标题', key: 'notion_page_title' },
  { title: 'Notion 页面 ID', key: 'notion_page_id' },
  {
    title: '状态',
    key: 'status',
    render: (row) =>
      h(
        NTag,
        {
          type: row.status === 'published' ? 'success' : 'default',
        },
        { default: () => (row.status === 'published' ? '已发布' : '草稿') }
      ),
  },
  {
    title: '内容路径',
    key: 'content_path',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '更新时间',
    key: 'updated_at',
    render: (row) => new Date(row.updated_at).toLocaleString(),
  },
]

// 获取文档列表
const fetchDocuments = async (params?: { offset?: number; limit?: number }) => {
  loading.value = true
  try {
    const data = await documentApi.getDocuments(props.projectId, {
      offset: params?.offset || 0,
      limit: params?.limit || pagination.pageSize,
    })
    documents.value = data.data
  } catch (error) {
    console.error('获取文档列表失败:', error)
  } finally {
    loading.value = false
  }
}

watch(
  () => props.projectId,
  () => {
    if (props.projectId) {
      fetchDocuments()
    }
  },
  { immediate: true }
)

onMounted(() => {
  fetchDocuments()
})
</script>

<style scoped>
.document-list {
  margin-top: 16px;
}
</style>
