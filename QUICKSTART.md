# 快速开始指南

## 前置要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+
- Docker 和 Docker Compose（可选）

## 本地开发

### 1. 配置数据库和 Redis

确保 MySQL 和 Redis 服务已启动。

### 2. 配置应用

复制配置文件示例：
```bash
cp config.yaml.example config.yaml
```

编辑 `config.yaml`，设置数据库和 Redis 连接信息。

### 3. 初始化数据库

运行数据库迁移：
```bash
mysql -u root -p < backend/migrations/001_init.sql
```

### 4. 启动后端服务

```bash
cd backend
go mod download
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动。

### 5. 启动前端服务

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 `http://localhost:3000` 启动。

## Docker 部署

### 1. 配置

复制并编辑配置文件：
```bash
cp config.yaml.example config.yaml
# 编辑 config.yaml，确保数据库和 Redis 配置正确
```

### 2. 启动服务

```bash
docker-compose up -d
```

### 3. 查看日志

```bash
docker-compose logs -f
```

## API 使用示例

### 创建项目

```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "notion_root_page_id": "your-page-id",
    "notion_root_page_title": "My Blog",
    "notion_token": "your-notion-token"
  }'
```

### 同步页面

```bash
curl -X POST http://localhost:8080/api/projects/1/sync \
  -H "Content-Type: application/json" \
  -d '{
    "page_id": "your-page-id"
  }'
```

### 转换为 Astro 格式

```bash
curl -X POST http://localhost:8080/api/projects/1/convert
```

## 配置 Webhook

1. 在 Notion 中创建 Webhook，指向 `http://your-domain/api/webhook/notion`
2. 获取 Webhook Secret
3. 在 `config.yaml` 中配置 `webhook.secret`
4. 重启服务

## 注意事项

1. **Notion Token**: 需要创建 Notion Integration 并获取 Token
2. **权限**: 确保 Notion Integration 有访问相应页面的权限
3. **Webhook Secret**: 生产环境务必配置 Webhook Secret 以确保安全
4. **文件存储**: 内容文件存储在 `content/` 目录下，请确保有写入权限

