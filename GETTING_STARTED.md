# 快速启动指南

本指南将帮助您快速启动 Notion Sync Blog 项目。

## 前置要求

- Go 1.21+
- MySQL 5.7+
- Redis 6.0+
- Notion Integration Token

## 1. 克隆项目

```bash
git clone <repository-url>
cd notion-sync-blog
```

## 2. 配置环境

### 2.1 复制环境变量文件

```bash
cp .env.example .env
```

根据您的环境修改 `.env` 文件中的配置。

### 2.2 复制配置文件

```bash
cp config.yaml.example config.yaml
```

根据您的环境修改 `config.yaml` 文件中的配置，特别是：
- 数据库连接信息
- Redis 连接信息
- Webhook Secret

## 3. 初始化数据库

### 方法一：使用自动脚本（推荐）

```bash
./init-db.sh config.yaml
```

### 方法二：手动执行

```bash
mysql -u root -p < backend/migrations/001_init.sql
```

## 4. 启动服务

### 4.1 使用 Docker Compose（推荐）

```bash
docker-compose up -d
```

这将启动：
- MySQL 数据库
- Redis 缓存
- 后端 API 服务
- 前端应用

### 4.2 手动启动

#### 启动 MySQL 和 Redis

```bash
# MySQL
sudo systemctl start mysql

# Redis
redis-server
```

#### 启动后端

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

#### 启动前端

```bash
cd frontend
npm install
npm run dev
```

## 5. 验证安装

### 5.1 检查后端 API

```bash
curl http://localhost:8080/api/projects
```

应该返回空的 `[]` 或项目列表。

### 5.2 检查前端

打开浏览器访问：`http://localhost:3000`

## 6. 配置 Notion Webhook

### 6.1 创建 Notion Integration

1. 访问 [Notion Developers](https://developers.notion.com/)
2. 创建新的 Integration
3. 复制 Internal Integration Token

### 6.2 设置 Webhook

1. 在 Notion 中创建 Webhook
2. 配置 Webhook URL: `http://your-server.com/api/webhook/notion`
3. 设置 Secret（在 config.yaml 中的 webhook.secret）
4. 选择要监听的事件

### 6.3 测试 Webhook

```bash
curl -X POST http://localhost:8080/api/webhook/notion \
  -H "Content-Type: application/json" \
  -d '{
    "object": "event",
    "entry": [{
      "id": "test-page-id",
      "event_type": "page.updated",
      "object": {"id": "test-page-id"},
      "time_stamp": 1234567890
    }]
  }'
```

## 7. 使用指南

### 7.1 创建项目

通过前端界面或 API 创建一个项目：

```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "notion_root_page_id": "your-page-id",
    "notion_root_page_title": "My Blog",
    "notion_token": "secret_xxx"
  }'
```

### 7.2 同步数据库

```bash
curl -X POST http://localhost:8080/api/projects/1/sync-database
```

### 7.3 同步页面

```bash
curl -X POST http://localhost:8080/api/projects/1/sync \
  -H "Content-Type: application/json" \
  -d '{
    "page_id": "your-page-id"
  }'
```

### 7.4 转换为 Astro 格式

```bash
curl -X POST http://localhost:8080/api/projects/1/convert
```

转换后的文件将保存在 `content/{project_id}/generated/astro/` 目录中。

## 8. 常见问题

### 8.1 数据库连接失败

检查：
- MySQL 服务是否运行
- 配置文件中的数据库连接信息是否正确
- 用户权限是否足够

### 8.2 Redis 连接失败

检查：
- Redis 服务是否运行
- 配置文件中的 Redis 连接信息是否正确

### 8.3 Notion API 调用失败

检查：
- Integration Token 是否正确
- Token 是否有访问页面的权限
- 页面 ID 是否正确

### 8.4 Webhook 签名验证失败

检查：
- config.yaml 中的 webhook.secret 是否正确
- Notion Webhook 的 Secret 是否与配置匹配

## 9. 目录结构

```
content/
  {project_id}/
    database/         # 数据库 JSON 缓存
    pages/            # 页面内容（.md 和 .json）
    generated/astro/  # 转换后的 Astro 文件
```

## 10. 进一步阅读

- [README.md](./README.md) - 详细项目说明
- [API 文档](./README.md#api-说明) - API 使用说明
- [Notion API 文档](https://developers.notion.com/)
- [Astro 文档](https://docs.astro.build/)
