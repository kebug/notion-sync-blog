package repository

import (
	"database/sql"
	"errors"
	"notion-sync-blog/internal/model"
	"notion-sync-blog/pkg/database"
	"time"
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
	log.CreatedAt = time.Now()

	query := `INSERT INTO sync_logs (project_id, event_type, event_id, status, error_message, created_at) 
			  VALUES (:project_id, :event_type, :event_id, :status, :error_message, :created_at)`

	result, err := database.DB.NamedExec(query, log)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.ID = uint(lastID)
	return nil
}

func (r *syncLogRepository) GetByID(id uint) (*model.SyncLog, error) {
	var log model.SyncLog
	query := `SELECT * FROM sync_logs WHERE id = ?`

	err := database.DB.Get(&log, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	// 加载关联的 Project
	if err := r.loadProject(&log); err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *syncLogRepository) ListByProjectID(projectID uint, offset, limit int) ([]*model.SyncLog, int64, error) {
	var logs []*model.SyncLog
	var total int64

	// 获取总数
	countQuery := `SELECT COUNT(*) FROM sync_logs WHERE project_id = ?`
	err := database.DB.Get(&total, countQuery, projectID)
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	query := `SELECT * FROM sync_logs WHERE project_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err = database.DB.Select(&logs, query, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *syncLogRepository) Update(log *model.SyncLog) error {
	query := `UPDATE sync_logs SET 
			  project_id = :project_id, 
			  event_type = :event_type, 
			  event_id = :event_id, 
			  status = :status, 
			  error_message = :error_message 
			  WHERE id = :id`

	_, err := database.DB.NamedExec(query, log)
	return err
}

// loadProject 加载关联的项目信息
func (r *syncLogRepository) loadProject(log *model.SyncLog) error {
	var project model.Project
	query := `SELECT * FROM projects WHERE id = ?`

	err := database.DB.Get(&project, query, log.ProjectID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	log.Project = &project
	return nil
}
