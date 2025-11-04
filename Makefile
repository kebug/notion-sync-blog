# Notion Sync Blog Makefile

# 变量定义
BINARY_NAME=notion-sync-blog
BACKEND_DIR=backend
FRONTEND_DIR=frontend
DOCKER_COMPOSE_FILE=docker-compose.yml

# 默认目标
.PHONY: all
all: build

# 构建后端
.PHONY: build-backend
build-backend:
	@echo "构建后端..."
	cd $(BACKEND_DIR) && go build -o ../$(BINARY_NAME) cmd/server/main.go

# 构建前端
.PHONY: build-frontend
build-frontend:
	@echo "构建前端..."
	cd $(FRONTEND_DIR) && npm run build

# 构建所有
.PHONY: build
build: build-backend build-frontend

# 运行后端开发服务器
.PHONY: run-backend
run-backend:
	@echo "启动后端开发服务器..."
	cd $(BACKEND_DIR) && go run cmd/server/main.go

# 运行前端开发服务器
.PHONY: run-frontend
run-frontend:
	@echo "启动前端开发服务器..."
	cd $(FRONTEND_DIR) && npm run dev

# 启动开发环境（前端和后端）
.PHONY: dev
dev: run-backend

# 安装依赖
.PHONY: install
install:
	@echo "安装后端依赖..."
	cd $(BACKEND_DIR) && go mod download
	@echo "安装前端依赖..."
	cd $(FRONTEND_DIR) && npm install

# 运行测试
.PHONY: test
test:
	@echo "运行后端测试..."
	cd $(BACKEND_DIR) && go test -v ./...

# 初始化数据库
.PHONY: db-init
db-init:
	@echo "初始化数据库..."
	./init-db.sh config.yaml

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	rm -f $(BINARY_NAME)
	cd $(BACKEND_DIR) && go clean
	cd $(FRONTEND_DIR) && rm -rf dist node_modules

# 使用 Docker Compose 启动
.PHONY: up
up:
	docker-compose up -d

# 停止 Docker Compose
.PHONY: down
down:
	docker-compose down

# 重启 Docker Compose
.PHONY: restart
restart:
	docker-compose restart

# 查看日志
.PHONY: logs
logs:
	docker-compose logs -f

# 迁移数据库
.PHONY: migrate
migrate:
	@echo "执行数据库迁移..."
	@if [ -f "$(DOCKER_COMPOSE_FILE)" ]; then \
		docker-compose exec mysql mysql -u root -p$$MYSQL_ROOT_PASSWORD < backend/migrations/001_init.sql; \
	else \
		mysql -u root -p < backend/migrations/001_init.sql; \
	fi

# 生成 Swagger 文档（如果使用 swagger）
.PHONY: swagger
swagger:
	@echo "生成 Swagger 文档..."
	cd $(BACKEND_DIR) && go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/server/main.go --dir internal/api --output docs

# 格式代码
.PHONY: fmt
fmt:
	@echo "格式化 Go 代码..."
	cd $(BACKEND_DIR) && gofmt -w .
	@echo "格式化前端代码..."
	cd $(FRONTEND_DIR) && npx prettier --write "src/**/*.{js,ts,vue,css,scss}"

# 代码检查
.PHONY: lint
lint:
	@echo "运行 Go 代码检查..."
	cd $(BACKEND_DIR) && go vet ./...
	@echo "运行前端代码检查..."
	cd $(FRONTEND_DIR) && npx eslint "src/**/*.{js,ts,vue}"

# 帮助信息
.PHONY: help
help:
	@echo "可用的 Make 命令:"
	@echo "  build          - 构建前后端"
	@echo "  build-backend  - 构建后端"
	@echo "  build-frontend - 构建前端"
	@echo "  dev            - 启动开发服务器"
	@echo "  install        - 安装所有依赖"
	@echo "  test           - 运行测试"
	@echo "  db-init        - 初始化数据库"
	@echo "  clean          - 清理构建文件"
	@echo "  up             - 使用 Docker 启动"
	@echo "  down           - 停止 Docker"
	@echo "  restart        - 重启 Docker"
	@echo "  logs           - 查看 Docker 日志"
	@echo "  migrate        - 执行数据库迁移"
	@echo "  fmt            - 格式化代码"
	@echo "  lint           - 代码检查"
	@echo "  help           - 显示帮助信息"
