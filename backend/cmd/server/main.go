package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notion-sync-blog/config"
	"notion-sync-blog/internal/api"
	"notion-sync-blog/internal/cache"
	"notion-sync-blog/internal/repository"
	"notion-sync-blog/internal/service"
	"notion-sync-blog/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.CloseDB()

	// 注意：使用 sqlx 后，需要手动执行数据库迁移
	// 建议使用数据库迁移工具（如 golang-migrate）管理数据库结构
	log.Println("使用 sqlx，请确保数据库表已创建（参考 migrations 目录）")

	// 初始化 Redis
	redisCache, err := cache.NewCache(&cfg.Redis)
	if err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}
	defer redisCache.Close()

	// 初始化 Repository
	projectRepo := repository.NewProjectRepository()
	documentRepo := repository.NewDocumentRepository()
	databaseRepo := repository.NewNotionDatabaseRepository()
	syncLogRepo := repository.NewSyncLogRepository()

	// 初始化 Service
	projectService := service.NewProjectService(projectRepo, &cfg.Storage)
	syncService := service.NewSyncService(
		projectRepo,
		documentRepo,
		databaseRepo,
		syncLogRepo,
		redisCache,
		&cfg.Storage,
	)

	// 初始化路由
	router := api.SetupRouter(
		projectService,
		syncService,
		documentRepo,
		syncLogRepo,
		redisCache,
	)

	// 启动服务器
	serverAddr := cfg.Server.GetServerAddr()
	log.Printf("服务器启动在 %s", serverAddr)

	// 优雅关闭
	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("服务器正在关闭...")
}
