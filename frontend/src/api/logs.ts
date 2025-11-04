import { apiClient } from './client'
import type { SyncLog, ApiResponse } from '../types'

export const logApi = {
  // 获取同步日志
  async getLogs(
    projectId: number,
    params?: { offset?: number; limit?: number }
  ) {
    const response = await apiClient.get<ApiResponse<SyncLog[]>>(
      `/projects/${projectId}/logs`,
      { params }
    )
    return response
  },
}
