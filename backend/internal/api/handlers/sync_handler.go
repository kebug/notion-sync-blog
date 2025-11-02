package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"notion-sync-blog/internal/service"
)

// SyncHandler 同步处理器
type SyncHandler struct {
	syncService service.SyncService
}

// NewSyncHandler 创建同步处理器
func NewSyncHandler(syncService service.SyncService) *SyncHandler {
	return &SyncHandler{
		syncService: syncService,
	}
}

// SyncPage 同步页面
func (h *SyncHandler) SyncPage(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	var req struct {
		PageID string `json:"page_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.syncService.SyncPage(uint(projectID), req.PageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "同步成功"})
}

// SyncDatabase 同步数据库
func (h *SyncHandler) SyncDatabase(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	if err := h.syncService.SyncDatabase(uint(projectID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "同步成功"})
}

// ConvertToAstro 转换为 Astro 格式
func (h *SyncHandler) ConvertToAstro(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	if err := h.syncService.ConvertToAstro(uint(projectID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "转换成功"})
}

// HandleWebhook 处理 Notion Webhook
func (h *SyncHandler) HandleWebhook(c *gin.Context) {
	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取请求体失败"})
		return
	}

	// TODO: 验证 Notion 签名
	// secret := []byte(cfg.Webhook.Secret)
	// signature := c.GetHeader("X-Notion-Signature")
	// if err := notion.VerifySignature(secret, body, signature); err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
	// 	return
	// }

	// 解析 Webhook 事件
	var webhookEvent struct {
		Object string `json:"object"`
		Entry  []struct {
			ID        string                 `json:"id"`
			TimeStamp int64                  `json:"time_stamp"`
			EventType string                 `json:"event_type"`
			Object    map[string]interface{} `json:"object"`
		} `json:"entry"`
	}

	if err := c.ShouldBindJSON(&webhookEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理每个事件
	for _, entry := range webhookEvent.Entry {
		// 从事件对象中提取项目 ID（需要根据实际 Webhook 结构调整）
		// 这里简化处理，假设能从事件中获取 project_id
		projectID := uint(1) // TODO: 从事件中获取实际项目 ID

		if err := h.syncService.HandleWebhookEvent(projectID, entry.EventType, entry.ID, entry.Object); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "处理成功"})
}

