<template>
  <div class="project-detail-page">
    <n-space vertical size="large">
      <n-page-header :title="project?.notion_root_page_title || '项目详情'" @back="handleBack">
        <template #extra>
          <n-space>
            <n-button @click="syncDatabase" :loading="syncing">
              同步数据库
            </n-button>
            <n-button type="primary" @click="convertToAstro" :loading="converting">
              转换为 Astro
            </n-button>
          </n-space>
        </template>
      </n-page-header>

      <n-tabs type="line" animated>
        <n-tab-pane name="documents" tab="文档列表">
          <DocumentList :project-id="projectId" />
        </n-tab-pane>
        <n-tab-pane name="logs" tab="同步日志">
          <SyncLogs :project-id="projectId" />
        </n-tab-pane>
        <n-tab-pane name="settings" tab="项目设置">
          <n-card title="项目信息" :bordered="false">
            <n-descriptions :column="2" bordered>
              <n-descriptions-item label="ID">
                {{ project?.id }}
              </n-descriptions-item>
              <n-descriptions-item label="项目标题">
                {{ project?.notion_root_page_title }}
              </n-descriptions-item>
              <n-descriptions-item label="Notion 页面 ID">
                {{ project?.notion_root_page_id }}
              </n-descriptions-item>
              <n-descriptions-item label="创建时间">
                {{ formatDate(project?.created_at) }}
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-tab-pane>
      </n-tabs>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NButton,
  NCard,
  NSpace,
  NPageHeader,
  NTabs,
  NTabPane,
  NDescriptions,
  NDescriptionsItem,
  useMessage,
} from 'naive-ui'
import { projectApi } from '../api/projects'
import { syncApi } from '../api/sync'
import type { Project } from '../types'
import DocumentList from '../components/DocumentList.vue'
import SyncLogs from '../components/SyncLogs.vue'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const projectId = computed(() => Number(route.params.id))
const project = ref<Project | null>(null)
const syncing = ref(false)
const converting = ref(false)

// 获取项目信息
const fetchProject = async () => {
  try {
    const data = await projectApi.getProject(projectId.value)
    project.value = data
  } catch (error) {
    message.error('获取项目信息失败')
  }
}

// 同步数据库
const syncDatabase = async () => {
  syncing.value = true
  try {
    await syncApi.syncDatabase(projectId.value)
    message.success('同步请求已发送，请查看日志')
  } catch (error) {
    message.error('同步失败')
  } finally {
    syncing.value = false
  }
}

// 转换为 Astro
const convertToAstro = async () => {
  converting.value = true
  try {
    await syncApi.convertToAstro(projectId.value)
    message.success('转换请求已发送，请查看日志')
  } catch (error) {
    message.error('转换失败')
  } finally {
    converting.value = false
  }
}

// 返回
const handleBack = () => {
  router.push('/')
}

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  fetchProject()
})
</script>

<style scoped>
.project-detail-page {
  padding: 24px;
}
</style>
