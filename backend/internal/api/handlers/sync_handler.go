package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"notion-sync-blog/internal/service"
	"notion-sync-blog/internal/notion"

	"github.com/gin-gonic/gin"
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

	// 从上下文或配置中获取签名（简化实现）
	// 在实际生产中，应该从配置中获取
	signature := c.GetHeader("X-Notion-Signature")
	if signature != "" {
		// 使用默认 secret（实际生产中应从配置获取）
		secret := []byte("your-webhook-secret")
		if err := notion.VerifySignature(secret, body, signature); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "签名验证失败"})
			return
		}
	}

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
		// 从事件对象中提取页面或数据库 ID
		var objectID string
		if obj, ok := entry.Object["id"].(string); ok {
			objectID = obj
		}

		// 根据事件类型和对象 ID 确定项目 ID
		// 实际生产环境中，应该维护一个映射表或通过数据库查询
		// 这里简化处理，假设从路径参数或事件数据中可以确定项目
		projectID, err := h.resolveProjectFromEvent(entry.EventType, objectID, entry.Object)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无法确定项目ID: %v", err)})
			return
		}

		if err := h.syncService.HandleWebhookEvent(projectID, entry.EventType, entry.ID, entry.Object); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "处理成功"})
}

// resolveProjectFromEvent 从事件中解析项目ID
// 实际生产中应该查询数据库根据页面或数据库ID找到对应的项目
func (h *SyncHandler) resolveProjectFromEvent(eventType, objectID string, eventObject map[string]interface{}) (uint, error) {
	// 简化实现：从查询参数中获取 project_id
	// 在实际生产中，应该查询数据库找到包含该页面或数据库的项目
	if projectIDStr := eventObject["project_id"]; projectIDStr != nil {
		if projectID, ok := projectIDStr.(string); ok {
			if id, err := strconv.ParseUint(projectID, 10, 32); err == nil {
				return uint(id), nil
			}
		} else if projectID, ok := projectIDStr.(float64); ok {
			return uint(projectID), nil
		}
	}

	// 如果无法从事件中获取，返回错误
	// 生产环境应该根据objectID查询数据库
	return 0, fmt.Errorf("无法从事件中提取项目ID: event_type=%s, object_id=%s", eventType, objectID)
}
