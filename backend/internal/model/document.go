package model

import (
	"time"
)

// Document 文档模型
type Document struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProjectID      uint      `gorm:"not null;index" json:"project_id"`
	NotionPageID   string    `gorm:"type:varchar(255);not null;index" json:"notion_page_id"`
	NotionPageTitle string   `gorm:"type:varchar(500)" json:"notion_page_title"`
	ContentPath    string    `gorm:"type:varchar(1000)" json:"content_path"`
	Status         string    `gorm:"type:varchar(50);default:'draft'" json:"status"` // draft, published
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	
	// 关联
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 指定表名
func (Document) TableName() string {
	return "notion_documents"
}

