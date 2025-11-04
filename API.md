# Notion Sync Blog API 文档

## 基础信息

- **基础 URL**: `http://localhost:8080/api`
- **内容类型**: `application/json`
- **认证方式**: 当前版本无需认证

## 通用响应格式

### 成功响应

```json
{
  "data": {},
  "message": "操作成功"
}
```

### 错误响应

```json
{
  "error": {
    "code": 400,
    "message": "错误描述",
    "type": "error_type"
  }
}
```

## 项目管理 API

### 1. 获取项目列表

**请求**:
```
GET /api/projects
```

**查询参数**:
- `offset` (int, optional): 偏移量，默认 0
- `limit` (int, optional): 每页数量，默认 20

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "notion_root_page_id": "page-uuid",
      "notion_root_page_title": "My Blog",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

### 2. 创建项目

**请求**:
```
POST /api/projects
Content-Type: application/json

{
  "notion_root_page_id": "page-uuid",
  "notion_root_page_title": "My Blog",
  "notion_token": "secret_xxx"
}
```

**响应**:
```json
{
  "id": 1,
  "notion_root_page_id": "page-uuid",
  "notion_root_page_title": "My Blog",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 3. 获取项目详情

**请求**:
```
GET /api/projects/{id}
```

**响应**:
```json
{
  "id": 1,
  "notion_root_page_id": "page-uuid",
  "notion_root_page_title": "My Blog",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 4. 更新项目

**请求**:
```
PUT /api/projects/{id}
Content-Type: application/json

{
  "title": "New Title"
}
```

**响应**:
```json
{
  "message": "更新成功"
}
```

### 5. 删除项目

**请求**:
```
DELETE /api/projects/{id}
```

**响应**:
```json
{
  "message": "删除成功"
}
```

## 文档管理 API

### 1. 获取文档列表

**请求**:
```
GET /api/projects/{id}/documents
```

**查询参数**:
- `offset` (int, optional): 偏移量，默认 0
- `limit` (int, optional): 每页数量，默认 20

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "project_id": 1,
      "notion_page_id": "page-uuid",
      "notion_page_title": "Article 1",
      "content_path": "./content/1/pages/page-uuid.md",
      "status": "draft",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

## 同步 API

### 1. 同步页面

**请求**:
```
POST /api/projects/{id}/sync
Content-Type: application/json

{
  "page_id": "page-uuid"
}
```

**响应**:
```json
{
  "message": "同步成功"
}
```

### 2. 同步数据库

**请求**:
```
POST /api/projects/{id}/sync-database
```

**响应**:
```json
{
  "message": "同步成功"
}
```

### 3. 转换为 Astro 格式

**请求**:
```
POST /api/projects/{id}/convert
```

**响应**:
```json
{
  "message": "转换成功"
}
```

## Webhook API

### 1. 接收 Notion Webhook

**请求**:
```
POST /api/webhook/notion
Content-Type: application/json
X-Notion-Signature: signature (如果配置了 secret)

{
  "object": "event",
  "entry": [
    {
      "id": "event-uuid",
      "event_type": "page.updated",
      "time_stamp": 1234567890,
      "object": {
        "id": "page-uuid"
      }
    }
  ]
}
```

**响应**:
```json
{
  "message": "处理成功"
}
```

## 日志 API

### 1. 获取同步日志

**请求**:
```
GET /api/projects/{id}/logs
```

**查询参数**:
- `offset` (int, optional): 偏移量，默认 0
- `limit` (int, optional): 每页数量，默认 50

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "project_id": 1,
      "event_type": "page.updated",
      "event_id": "event-uuid",
      "status": "success",
      "error_message": null,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

## 健康检查 API

### 1. 健康检查

**请求**:
```
GET /health
```

**响应**:
```json
{
  "status": "ok",
  "timestamp": 1234567890
}
```

### 2. 就绪检查

**请求**:
```
GET /ready
```

**响应**:
```json
{
  "status": "ready",
  "timestamp": 1234567890
}
```

## 错误代码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 204 | 无内容 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 状态值说明

### 文档状态 (status)
- `draft`: 草稿
- `published`: 已发布

### 同步状态 (status)
- `pending`: 处理中
- `success`: 成功
- `failed`: 失败

## 限制

- API 请求限制：当前无限制
- 文件大小限制：由服务器配置决定
- 同步频率：建议不超过每分钟 10 次

## 注意事项

1. Notion Token 会被加密存储，但 API 响应中不会返回
2. Webhook 需要配置正确的签名 Secret
3. 页面 ID 必须是有效的 Notion 页面 ID
4. 转换后的文件保存在 `content/{project_id}/generated/astro/` 目录
