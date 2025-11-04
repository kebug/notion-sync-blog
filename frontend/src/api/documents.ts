import { apiClient } from './client'
import type { Document, ApiResponse } from '../types'

export const documentApi = {
  // 获取项目文档列表
  async getDocuments(
    projectId: number,
    params?: { offset?: number; limit?: number }
  ) {
    const response = await apiClient.get<ApiResponse<Document[]>>(
      `/projects/${projectId}/documents`,
      { params }
    )
    return response
  },
}
