<template>
  <div>
    <n-space vertical>
      <n-h2>项目管理</n-h2>
      <n-button type="primary" @click="showCreateModal = true">
        创建项目
      </n-button>
      <n-table :data="projects" :loading="loading">
        <thead>
          <tr>
            <th>ID</th>
            <th>页面标题</th>
            <th>根页面 ID</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="project in projects" :key="project.id">
            <td>{{ project.id }}</td>
            <td>{{ project.notion_root_page_title }}</td>
            <td>{{ project.notion_root_page_id }}</td>
            <td>{{ project.created_at }}</td>
            <td>
              <n-space>
                <n-button size="small" @click="viewProject(project.id)">
                  查看
                </n-button>
                <n-button size="small" @click="syncProject(project.id)">
                  同步
                </n-button>
              </n-space>
            </td>
          </tr>
        </tbody>
      </n-table>
    </n-space>

    <n-modal v-model:show="showCreateModal">
      <n-card style="width: 600px" title="创建项目">
        <n-form>
          <n-form-item label="Notion 根页面 ID">
            <n-input v-model:value="newProject.notion_root_page_id" />
          </n-form-item>
          <n-form-item label="页面标题">
            <n-input v-model:value="newProject.notion_root_page_title" />
          </n-form-item>
          <n-form-item label="Notion Token">
            <n-input v-model:value="newProject.notion_token" type="password" />
          </n-form-item>
          <n-form-item>
            <n-button @click="createProject">创建</n-button>
            <n-button @click="showCreateModal = false">取消</n-button>
          </n-form-item>
        </n-form>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NH2, NSpace, NButton, NTable, NCard, NModal, NForm, NFormItem, NInput } from 'naive-ui'
import client from '../api/client'

const projects = ref([])
const loading = ref(false)
const showCreateModal = ref(false)
const newProject = ref({
  notion_root_page_id: '',
  notion_root_page_title: '',
  notion_token: '',
})

const loadProjects = async () => {
  loading.value = true
  try {
    const res = await client.get('/projects')
    projects.value = res.data.data
  } catch (error) {
    console.error('加载项目失败:', error)
  } finally {
    loading.value = false
  }
}

const createProject = async () => {
  try {
    await client.post('/projects', newProject.value)
    showCreateModal.value = false
    loadProjects()
  } catch (error) {
    console.error('创建项目失败:', error)
  }
}

const viewProject = (id: number) => {
  // TODO: 实现查看项目详情
  console.log('查看项目:', id)
}

const syncProject = async (id: number) => {
  try {
    await client.post(`/projects/${id}/sync-database`)
    alert('同步成功')
  } catch (error) {
    console.error('同步失败:', error)
  }
}

onMounted(() => {
  loadProjects()
})
</script>

