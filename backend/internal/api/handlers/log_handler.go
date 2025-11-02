package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"notion-sync-blog/internal/repository"
)

// LogHandler 日志处理器
type LogHandler struct {
	syncLogRepo repository.SyncLogRepository
}

// NewLogHandler 创建日志处理器
func NewLogHandler(syncLogRepo repository.SyncLogRepository) *LogHandler {
	return &LogHandler{
		syncLogRepo: syncLogRepo,
	}
}

// ListLogs 获取同步日志列表
func (h *LogHandler) ListLogs(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	logs, total, err := h.syncLogRepo.ListByProjectID(uint(projectID), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
	})
}

