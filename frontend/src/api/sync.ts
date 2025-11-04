import { apiClient } from './client'

export const syncApi = {
  // 同步页面
  async syncPage(projectId: number, pageId: string) {
    const response = await apiClient.post(`/projects/${projectId}/sync`, {
      page_id: pageId,
    })
    return response
  },

  // 同步数据库
  async syncDatabase(projectId: number) {
    const response = await apiClient.post(`/projects/${projectId}/sync-database`)
    return response
  },

  // 转换为 Astro 格式
  async convertToAstro(projectId: number) {
    const response = await apiClient.post(`/projects/${projectId}/convert`)
    return response
  },
}
