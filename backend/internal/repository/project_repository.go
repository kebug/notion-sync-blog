package repository

import (
	"database/sql"
	"errors"
	"notion-sync-blog/internal/model"
	"notion-sync-blog/pkg/database"
	"time"
)

// ProjectRepository 项目仓库接口
type ProjectRepository interface {
	Create(project *model.Project) error
	GetByID(id uint) (*model.Project, error)
	GetByNotionPageID(pageID string) (*model.Project, error)
	List(offset, limit int) ([]*model.Project, int64, error)
	Update(project *model.Project) error
	Delete(id uint) error
}

type projectRepository struct{}

// NewProjectRepository 创建项目仓库
func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

func (r *projectRepository) Create(project *model.Project) error {
	now := time.Now()
	project.CreatedAt = now
	project.UpdatedAt = now

	query := `INSERT INTO projects (notion_root_page_id, notion_root_page_title, notion_token, created_at, updated_at) 
			  VALUES (:notion_root_page_id, :notion_root_page_title, :notion_token, :created_at, :updated_at)`

	result, err := database.DB.NamedExec(query, project)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	project.ID = uint(lastID)
	return nil
}

func (r *projectRepository) GetByID(id uint) (*model.Project, error) {
	var project model.Project
	query := `SELECT * FROM projects WHERE id = ?`

	err := database.DB.Get(&project, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) GetByNotionPageID(pageID string) (*model.Project, error) {
	var project model.Project
	query := `SELECT * FROM projects WHERE notion_root_page_id = ? LIMIT 1`

	err := database.DB.Get(&project, query, pageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) List(offset, limit int) ([]*model.Project, int64, error) {
	var projects []*model.Project
	var total int64

	// 获取总数
	countQuery := `SELECT COUNT(*) FROM projects`
	err := database.DB.Get(&total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	query := `SELECT * FROM projects ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err = database.DB.Select(&projects, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) Update(project *model.Project) error {
	project.UpdatedAt = time.Now()

	query := `UPDATE projects SET 
			  notion_root_page_id = :notion_root_page_id, 
			  notion_root_page_title = :notion_root_page_title, 
			  notion_token = :notion_token, 
			  updated_at = :updated_at 
			  WHERE id = :id`

	_, err := database.DB.NamedExec(query, project)
	return err
}

func (r *projectRepository) Delete(id uint) error {
	query := `DELETE FROM projects WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}
