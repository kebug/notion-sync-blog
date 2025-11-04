package repository

import (
	"database/sql"
	"errors"
	"notion-sync-blog/internal/model"
	"notion-sync-blog/pkg/database"
	"time"
)

// DocumentRepository 文档仓库接口
type DocumentRepository interface {
	Create(document *model.Document) error
	GetByID(id uint) (*model.Document, error)
	GetByNotionPageID(projectID uint, pageID string) (*model.Document, error)
	ListByProjectID(projectID uint, offset, limit int) ([]*model.Document, int64, error)
	Update(document *model.Document) error
	Delete(id uint) error
	DeleteByProjectID(projectID uint) error
}

type documentRepository struct{}

// NewDocumentRepository 创建文档仓库
func NewDocumentRepository() DocumentRepository {
	return &documentRepository{}
}

func (r *documentRepository) Create(document *model.Document) error {
	now := time.Now()
	document.CreatedAt = now
	document.UpdatedAt = now

	query := `INSERT INTO notion_documents (project_id, notion_page_id, notion_page_title, content_path, status, created_at, updated_at) 
			  VALUES (:project_id, :notion_page_id, :notion_page_title, :content_path, :status, :created_at, :updated_at)`

	result, err := database.DB.NamedExec(query, document)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	document.ID = uint(lastID)
	return nil
}

func (r *documentRepository) GetByID(id uint) (*model.Document, error) {
	var document model.Document
	query := `SELECT * FROM notion_documents WHERE id = ?`

	err := database.DB.Get(&document, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	// 加载关联的 Project
	if err := r.loadProject(&document); err != nil {
		return nil, err
	}

	return &document, nil
}

func (r *documentRepository) GetByNotionPageID(projectID uint, pageID string) (*model.Document, error) {
	var document model.Document
	query := `SELECT * FROM notion_documents WHERE project_id = ? AND notion_page_id = ?`

	err := database.DB.Get(&document, query, projectID, pageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &document, nil
}

func (r *documentRepository) ListByProjectID(projectID uint, offset, limit int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	// 获取总数
	countQuery := `SELECT COUNT(*) FROM notion_documents WHERE project_id = ?`
	err := database.DB.Get(&total, countQuery, projectID)
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	query := `SELECT * FROM notion_documents WHERE project_id = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?`
	err = database.DB.Select(&documents, query, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

func (r *documentRepository) Update(document *model.Document) error {
	document.UpdatedAt = time.Now()

	query := `UPDATE notion_documents SET 
			  project_id = :project_id, 
			  notion_page_id = :notion_page_id, 
			  notion_page_title = :notion_page_title, 
			  content_path = :content_path, 
			  status = :status, 
			  updated_at = :updated_at 
			  WHERE id = :id`

	_, err := database.DB.NamedExec(query, document)
	return err
}

func (r *documentRepository) Delete(id uint) error {
	query := `DELETE FROM notion_documents WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}

func (r *documentRepository) DeleteByProjectID(projectID uint) error {
	query := `DELETE FROM notion_documents WHERE project_id = ?`
	_, err := database.DB.Exec(query, projectID)
	return err
}

// loadProject 加载关联的项目信息
func (r *documentRepository) loadProject(document *model.Document) error {
	var project model.Project
	query := `SELECT * FROM projects WHERE id = ?`

	err := database.DB.Get(&project, query, document.ProjectID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	document.Project = &project
	return nil
}
