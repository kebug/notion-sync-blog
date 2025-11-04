import { apiClient } from './client'
import type { Project, ApiResponse } from '../types'

export const projectApi = {
  // 获取项目列表
  async getProjects(params?: { offset?: number; limit?: number }) {
    const response = await apiClient.get<ApiResponse<Project[]>>('/projects', {
      params,
    })
    return response
  },

  // 获取项目详情
  async getProject(id: number) {
    const response = await apiClient.get<Project>(`/projects/${id}`)
    return response
  },

  // 创建项目
  async createProject(data: {
    notion_root_page_id: string
    notion_root_page_title?: string
    notion_token: string
  }) {
    const response = await apiClient.post<Project>('/projects', data)
    return response
  },

  // 更新项目
  async updateProject(id: number, data: { title: string }) {
    const response = await apiClient.put(`/projects/${id}`, data)
    return response
  },

  // 删除项目
  async deleteProject(id: number) {
    const response = await apiClient.delete(`/projects/${id}`)
    return response
  },
}
