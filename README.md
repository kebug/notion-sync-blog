# NOTION SYNC BLOG

## 项目介绍

在notion的一个页面下，编辑文档后，该服务通过notion的webhook接受通知，然后使用notion api，将最新的内容下载到本地，自动部署博客。


## notion格式

在notion中创建一个数据库用于管理需要发布的文档以及文档的基本信息，服务根据该数据库管理具体的内容。



## 项目技术栈
- **后端**: Golang, Gin, MySQL, Redis
- **前端**: Vue 3, TypeScript, Naive UI
- **缓存**: Redis（缓存 Notion 数据库内容和页面内容）
- **配置**: YAML 格式配置文件


## 处理流程
1. 前端页面选择 Notion 中的根页面作为项目
2. 服务下载该页面下的管理数据库到 Redis 缓存（固定格式）
3. 接收 Notion Webhook 事件（使用 Notion 签名验证确保安全性）
4. 数据库更新就更新数据库到 Redis 缓存和本地存储
5. 页面更新时，检查页面是否在数据库范围中，如果是则更新文档到本地
6. 管理页面可以将文档转换为 Astro 静态博客格式（目前仅支持 Astro，转换基础必备信息）



## notion webhook
Notion Webhook API 见文档：https://developers.notion.com/reference/webhooks-events-delivery

### 安全验证
- 使用 Notion 签名验证确保 Webhook 请求的合法性
- 验证算法：HMAC-SHA256

## 数据库设计

### MySQL 表结构

#### projects 表
存储 Notion 项目配置
- `id`: 主键
- `notion_root_page_id`: Notion 根页面 ID
- `notion_root_page_title`: 根页面标题
- `notion_token`: Notion Integration Token（加密存储）
- `created_at`: 创建时间
- `updated_at`: 更新时间

#### notion_documents 表
存储同步的文档信息
- `id`: 主键
- `project_id`: 关联项目 ID
- `notion_page_id`: Notion 页面 ID
- `notion_page_title`: 页面标题
- `content_path`: 本地内容文件路径
- `status`: 状态（draft/published）
- `created_at`: 创建时间
- `updated_at`: 更新时间

#### notion_databases 表
存储管理数据库的缓存信息
- `id`: 主键
- `project_id`: 关联项目 ID
- `notion_database_id`: Notion 数据库 ID
- `cached_data`: 缓存的数据库内容（JSON 格式）
- `last_synced_at`: 最后同步时间

#### sync_logs 表
同步日志记录
- `id`: 主键
- `project_id`: 关联项目 ID
- `event_type`: 事件类型（page.updated, database.updated 等）
- `event_id`: Notion 事件 ID
- `status`: 处理状态（success/failed）
- `error_message`: 错误信息（可选）
- `created_at`: 创建时间

## 缓存策略

### Redis 缓存
- **缓存数据库内容**: 缓存 Notion 数据库查询结果
- **缓存页面内容**: 缓存 Notion 页面块内容
- 缓存 key 格式：
  - `notion:database:{project_id}:{database_id}`
  - `notion:page:{project_id}:{page_id}`

## 文件存储结构

```
content/
  {project_id}/
    database/
      {database_id}.json      # 数据库 JSON 缓存
    pages/
      {page_id}.md            # Markdown 格式
      {page_id}.json          # 原始 JSON 格式
    generated/
      astro/                  # Astro 格式输出
        {page_slug}.md        # 转换后的 Astro 格式文件
```

## 格式转换（Astro）

目前仅支持转换为 Astro 格式，转换以下基础信息：
- 页面标题 → Front Matter `title`
- 创建时间 → Front Matter `created`
- 更新时间 → Front Matter `updated`
- 状态 → Front Matter `status`
- Notion 块内容 → Markdown 正文
  - 标题块 → Markdown 标题
  - 段落块 → Markdown 段落
  - 列表块 → Markdown 列表
  - 代码块 → Markdown 代码块
  - 图片块 → Markdown 图片链接

## 配置文件

使用 YAML 格式配置文件，主要配置项：
- MySQL 连接信息
- Redis 连接信息
- 服务器端口
- 文件存储路径
- CORS 配置
- 日志配置

## 项目结构

```
notion-sync-blog/
├── backend/              # Go 后端服务
│   ├── cmd/
│   │   └── server/      # 主程序入口
│   ├── internal/
│   │   ├── api/         # API 路由处理
│   │   ├── service/     # 业务逻辑层
│   │   ├── repository/  # 数据访问层
│   │   ├── model/       # 数据模型
│   │   ├── cache/       # Redis 缓存封装
│   │   ├── notion/      # Notion API 封装
│   │   └── converter/  # 格式转换器
│   ├── pkg/             # 公共包
│   ├── config/          # 配置管理
│   ├── migrations/      # 数据库迁移文件
│   └── go.mod
├── frontend/            # Vue 前端
│   ├── src/
│   │   ├── views/       # 页面视图
│   │   ├── components/  # 组件
│   │   ├── api/         # API 调用
│   │   └── utils/       # 工具函数
│   └── package.json
├── config.yaml.example  # 配置文件示例
├── .env.example         # 环境变量示例
├── docker-compose.yml   # Docker 编排文件
└── README.md

```

## API 说明

### Webhook 端点
- `POST /api/webhook/notion` - 接收 Notion Webhook 事件

### 项目管理 API
- `GET /api/projects` - 获取项目列表
- `POST /api/projects` - 创建项目
- `GET /api/projects/:id` - 获取项目详情
- `PUT /api/projects/:id` - 更新项目
- `DELETE /api/projects/:id` - 删除项目

### 文档管理 API
- `GET /api/projects/:id/documents` - 获取文档列表
- `POST /api/projects/:id/sync` - 手动触发同步
- `POST /api/projects/:id/convert` - 转换为 Astro 格式

### 日志 API
- `GET /api/projects/:id/logs` - 获取同步日志
