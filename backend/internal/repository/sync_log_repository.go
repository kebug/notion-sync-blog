package repository

import (
	"notion-sync-blog/internal/model"
	"notion-sync-blog/pkg/database"
)

// SyncLogRepository 同步日志仓库接口
type SyncLogRepository interface {
	Create(log *model.SyncLog) error
	GetByID(id uint) (*model.SyncLog, error)
	ListByProjectID(projectID uint, offset, limit int) ([]*model.SyncLog, int64, error)
	Update(log *model.SyncLog) error
}

type syncLogRepository struct{}

// NewSyncLogRepository 创建同步日志仓库
func NewSyncLogRepository() SyncLogRepository {
	return &syncLogRepository{}
}

func (r *syncLogRepository) Create(log *model.SyncLog) error {
	return database.DB.Create(log).Error
}

func (r *syncLogRepository) GetByID(id uint) (*model.SyncLog, error) {
	var log model.SyncLog
	if err := database.DB.Preload("Project").First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *syncLogRepository) ListByProjectID(projectID uint, offset, limit int) ([]*model.SyncLog, int64, error) {
	var logs []*model.SyncLog
	var total int64

	query := database.DB.Model(&model.SyncLog{}).Where("project_id = ?", projectID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *syncLogRepository) Update(log *model.SyncLog) error {
	return database.DB.Save(log).Error
}

