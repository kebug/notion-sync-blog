<template>
  <div class="sync-logs">
    <n-card>
      <n-data-table
        :columns="columns"
        :data="logs"
        :loading="loading"
        :pagination="pagination"
      />
      <n-empty v-if="!loading && logs.length === 0" description="暂无日志" class="mt-4" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, h } from 'vue'
import { NDataTable, NCard, NEmpty, NTag, NTooltip, type DataTableColumns } from 'naive-ui'
import { logApi } from '../api/logs'
import type { SyncLog } from '../types'

interface Props {
  projectId: number
}

const props = defineProps<Props>()

// 响应式数据
const logs = ref<SyncLog[]>([])
const loading = ref(false)

const pagination = {
  page: 1,
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [20, 50, 100],
  onChange: (page: number) => {
    fetchLogs({ offset: (page - 1) * pagination.pageSize })
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    fetchLogs({ offset: 0, limit: pageSize })
  },
}

// 表格列定义
const columns: DataTableColumns<SyncLog> = [
  { title: 'ID', key: 'id', width: 80 },
  {
    title: '事件类型',
    key: 'event_type',
    render: (row) => {
      const typeMap: Record<string, { label: string; type: string }> = {
        'page.updated': { label: '页面更新', type: 'info' },
        'database.updated': { label: '数据库更新', type: 'warning' },
        'page.sync': { label: '页面同步', type: 'default' },
        'database.sync': { label: '数据库同步', type: 'default' },
      }
      const info = typeMap[row.event_type] || { label: row.event_type, type: 'default' }
      return h(NTag, { type: info.type as any }, { default: () => info.label })
    },
  },
  { title: '事件 ID', key: 'event_id', ellipsis: { tooltip: true } },
  {
    title: '状态',
    key: 'status',
    render: (row) => {
      const statusMap = {
        pending: { label: '处理中', type: 'warning' },
        success: { label: '成功', type: 'success' },
        failed: { label: '失败', type: 'error' },
      }
      const info = statusMap[row.status as keyof typeof statusMap]
      return h(
        NTag,
        { type: info?.type as any },
        { default: () => info?.label || row.status }
      )
    },
  },
  {
    title: '错误信息',
    key: 'error_message',
    render: (row) =>
      row.error_message
        ? h(
            NTooltip,
            { trigger: 'hover' },
            {
              trigger: () =>
                h(
                  NTag,
                  { type: 'error', size: 'small' },
                  { default: () => '错误' }
                ),
              default: () => row.error_message,
            }
          )
        : '-',
    ellipsis: { tooltip: true },
  },
  {
    title: '创建时间',
    key: 'created_at',
    render: (row) => new Date(row.created_at).toLocaleString(),
  },
]

// 获取日志列表
const fetchLogs = async (params?: { offset?: number; limit?: number }) => {
  loading.value = true
  try {
    const data = await logApi.getLogs(props.projectId, {
      offset: params?.offset || 0,
      limit: params?.limit || pagination.pageSize,
    })
    logs.value = data.data
  } catch (error) {
    console.error('获取日志列表失败:', error)
  } finally {
    loading.value = false
  }
}

watch(
  () => props.projectId,
  () => {
    if (props.projectId) {
      fetchLogs()
    }
  },
  { immediate: true }
)

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.sync-logs {
  margin-top: 16px;
}
</style>
