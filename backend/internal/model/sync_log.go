package model

import (
	"time"
)

// SyncLog 同步日志模型
type SyncLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProjectID    uint      `gorm:"not null;index" json:"project_id"`
	EventType    string    `gorm:"type:varchar(100)" json:"event_type"` // page.updated, database.updated 等
	EventID      string    `gorm:"type:varchar(255)" json:"event_id"`
	Status       string    `gorm:"type:varchar(50);default:'pending'" json:"status"` // success, failed, pending
	ErrorMessage string    `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	
	// 关联
	Project Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName 指定表名
func (SyncLog) TableName() string {
	return "sync_logs"
}

