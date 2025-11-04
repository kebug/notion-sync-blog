package model

import (
	"time"
)

// NotionDatabase 数据库模型
type NotionDatabase struct {
	ID               uint      `db:"id" json:"id"`
	ProjectID        uint      `db:"project_id" json:"project_id"`
	NotionDatabaseID string    `db:"notion_database_id" json:"notion_database_id"`
	CachedData       string    `db:"cached_data" json:"cached_data"` // JSON 格式
	LastSyncedAt     time.Time `db:"last_synced_at" json:"last_synced_at"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`

	// 关联（用于查询时关联加载）
	Project *Project `db:"-" json:"project,omitempty"`
}

// TableName 指定表名
func (NotionDatabase) TableName() string {
	return "notion_databases"
}
