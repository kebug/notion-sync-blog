package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"notion-sync-blog/internal/repository"
)

// DocumentHandler 文档处理器
type DocumentHandler struct {
	documentRepo repository.DocumentRepository
}

// NewDocumentHandler 创建文档处理器
func NewDocumentHandler(documentRepo repository.DocumentRepository) *DocumentHandler {
	return &DocumentHandler{
		documentRepo: documentRepo,
	}
}

// ListDocuments 获取文档列表
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	documents, total, err := h.documentRepo.ListByProjectID(uint(projectID), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  documents,
		"total": total,
	})
}

