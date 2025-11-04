package model

import (
	"time"
)

// Document 文档模型
type Document struct {
	ID              uint      `db:"id" json:"id"`
	ProjectID       uint      `db:"project_id" json:"project_id"`
	NotionPageID    string    `db:"notion_page_id" json:"notion_page_id"`
	NotionPageTitle string    `db:"notion_page_title" json:"notion_page_title"`
	ContentPath     string    `db:"content_path" json:"content_path"`
	Status          string    `db:"status" json:"status"` // draft, published
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`

	// 关联（用于查询时关联加载）
	Project *Project `db:"-" json:"project,omitempty"`
}

// TableName 指定表名
func (Document) TableName() string {
	return "notion_documents"
}
