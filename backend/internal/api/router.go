package api

import (
	"notion-sync-blog/internal/api/handlers"
	"notion-sync-blog/internal/cache"
	"notion-sync-blog/internal/middleware"
	"notion-sync-blog/internal/repository"
	"notion-sync-blog/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(
	projectService service.ProjectService,
	syncService service.SyncService,
	documentRepo repository.DocumentRepository,
	syncLogRepo repository.SyncLogRepository,
	cache *cache.Cache,
) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware())

	// API 路由
	api := router.Group("/api")
	{
		// 项目相关
		projectHandler := handlers.NewProjectHandler(projectService)
		api.GET("/projects", projectHandler.ListProjects)
		api.POST("/projects", projectHandler.CreateProject)
		api.GET("/projects/:id", projectHandler.GetProject)
		api.PUT("/projects/:id", projectHandler.UpdateProject)
		api.DELETE("/projects/:id", projectHandler.DeleteProject)

		// 同步相关
		syncHandler := handlers.NewSyncHandler(syncService)
		api.POST("/projects/:id/sync", syncHandler.SyncPage)
		api.POST("/projects/:id/sync-database", syncHandler.SyncDatabase)
		api.POST("/projects/:id/convert", syncHandler.ConvertToAstro)
		api.POST("/webhook/notion", syncHandler.HandleWebhook)

		// 文档相关
		documentHandler := handlers.NewDocumentHandler(documentRepo)
		api.GET("/projects/:id/documents", documentHandler.ListDocuments)

		// 日志相关
		logHandler := handlers.NewLogHandler(syncLogRepo)
		api.GET("/projects/:id/logs", logHandler.ListLogs)
	}

	return router
}
