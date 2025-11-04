package model

import (
	"time"
)

// SyncLog 同步日志模型
type SyncLog struct {
	ID           uint      `db:"id" json:"id"`
	ProjectID    uint      `db:"project_id" json:"project_id"`
	EventType    string    `db:"event_type" json:"event_type"` // page.updated, database.updated 等
	EventID      string    `db:"event_id" json:"event_id"`
	Status       string    `db:"status" json:"status"` // success, failed, pending
	ErrorMessage string    `db:"error_message" json:"error_message,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`

	// 关联（用于查询时关联加载）
	Project *Project `db:"-" json:"project,omitempty"`
}

// TableName 指定表名
func (SyncLog) TableName() string {
	return "sync_logs"
}
