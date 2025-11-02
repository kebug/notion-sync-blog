package model

import (
	"time"
)

// NotionDatabase 数据库模型
type NotionDatabase struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProjectID      uint      `gorm:"not null;index" json:"project_id"`
	NotionDatabaseID string  `gorm:"type:varchar(255);not null;index" json:"notion_database_id"`
	CachedData     string    `gorm:"type:longtext" json:"cached_data"` // JSON 格式
	LastSyncedAt   time.Time `json:"last_synced_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	
	// 关联
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 指定表名
func (NotionDatabase) TableName() string {
	return "notion_databases"
}

