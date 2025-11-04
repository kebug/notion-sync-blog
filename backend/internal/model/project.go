package model

import (
	"time"
)

// Project 项目模型
type Project struct {
	ID                  uint      `db:"id" json:"id"`
	NotionRootPageID    string    `db:"notion_root_page_id" json:"notion_root_page_id"`
	NotionRootPageTitle string    `db:"notion_root_page_title" json:"notion_root_page_title"`
	NotionToken         string    `db:"notion_token" json:"-"` // 不返回给前端
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}
