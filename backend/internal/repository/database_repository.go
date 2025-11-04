package repository

import (
	"database/sql"
	"errors"
	"notion-sync-blog/internal/model"
	"notion-sync-blog/pkg/database"
	"time"
)

// NotionDatabaseRepository 数据库仓库接口
type NotionDatabaseRepository interface {
	Create(db *model.NotionDatabase) error
	GetByID(id uint) (*model.NotionDatabase, error)
	GetByNotionDatabaseID(projectID uint, dbID string) (*model.NotionDatabase, error)
	GetByProjectID(projectID uint) (*model.NotionDatabase, error)
	Update(db *model.NotionDatabase) error
	Delete(id uint) error
}

type notionDatabaseRepository struct{}

// NewNotionDatabaseRepository 创建数据库仓库
func NewNotionDatabaseRepository() NotionDatabaseRepository {
	return &notionDatabaseRepository{}
}

func (r *notionDatabaseRepository) Create(db *model.NotionDatabase) error {
	now := time.Now()
	db.CreatedAt = now
	db.UpdatedAt = now

	query := `INSERT INTO notion_databases (project_id, notion_database_id, cached_data, last_synced_at, created_at, updated_at) 
			  VALUES (:project_id, :notion_database_id, :cached_data, :last_synced_at, :created_at, :updated_at)`

	result, err := database.DB.NamedExec(query, db)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	db.ID = uint(lastID)
	return nil
}

func (r *notionDatabaseRepository) GetByID(id uint) (*model.NotionDatabase, error) {
	var db model.NotionDatabase
	query := `SELECT * FROM notion_databases WHERE id = ?`

	err := database.DB.Get(&db, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	// 加载关联的 Project
	if err := r.loadProject(&db); err != nil {
		return nil, err
	}

	return &db, nil
}

func (r *notionDatabaseRepository) GetByNotionDatabaseID(projectID uint, dbID string) (*model.NotionDatabase, error) {
	var db model.NotionDatabase
	query := `SELECT * FROM notion_databases WHERE project_id = ? AND notion_database_id = ?`

	err := database.DB.Get(&db, query, projectID, dbID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &db, nil
}

func (r *notionDatabaseRepository) GetByProjectID(projectID uint) (*model.NotionDatabase, error) {
	var db model.NotionDatabase
	query := `SELECT * FROM notion_databases WHERE project_id = ? LIMIT 1`

	err := database.DB.Get(&db, query, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &db, nil
}

func (r *notionDatabaseRepository) Update(db *model.NotionDatabase) error {
	db.UpdatedAt = time.Now()

	query := `UPDATE notion_databases SET 
			  project_id = :project_id, 
			  notion_database_id = :notion_database_id, 
			  cached_data = :cached_data, 
			  last_synced_at = :last_synced_at, 
			  updated_at = :updated_at 
			  WHERE id = :id`

	_, err := database.DB.NamedExec(query, db)
	return err
}

func (r *notionDatabaseRepository) Delete(id uint) error {
	query := `DELETE FROM notion_databases WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}

// loadProject 加载关联的项目信息
func (r *notionDatabaseRepository) loadProject(db *model.NotionDatabase) error {
	var project model.Project
	query := `SELECT * FROM projects WHERE id = ?`

	err := database.DB.Get(&project, query, db.ProjectID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	db.Project = &project
	return nil
}
