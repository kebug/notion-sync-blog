// 项目类型
export interface Project {
  id: number
  notion_root_page_id: string
  notion_root_page_title: string
  created_at: string
  updated_at: string
}

// 文档类型
export interface Document {
  id: number
  project_id: number
  notion_page_id: string
  notion_page_title: string
  content_path: string
  status: 'draft' | 'published'
  created_at: string
  updated_at: string
  project?: Project
}

// 同步日志类型
export interface SyncLog {
  id: number
  project_id: number
  event_type: string
  event_id: string
  status: 'pending' | 'success' | 'failed'
  error_message?: string
  created_at: string
  project?: Project
}

// API 响应类型
export interface ApiResponse<T> {
  data: T
  total?: number
  message?: string
}

export interface ApiError {
  code: number
  message: string
  type?: string
}

// 分页参数
export interface PaginationParams {
  offset: number
  limit: number
}
