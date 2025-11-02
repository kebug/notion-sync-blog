package model

import (
	"time"
)

// Project 项目模型
type Project struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	NotionRootPageID  string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"notion_root_page_id"`
	NotionRootPageTitle string  `gorm:"type:varchar(500)" json:"notion_root_page_title"`
	NotionToken       string    `gorm:"type:text;not null" json:"-"` // 不返回给前端
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

