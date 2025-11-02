package repository

import (
	"notion-sync-blog/internal/model"
	"notion-sync-blog/pkg/database"
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
	return database.DB.Create(document).Error
}

func (r *documentRepository) GetByID(id uint) (*model.Document, error) {
	var document model.Document
	if err := database.DB.Preload("Project").First(&document, id).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *documentRepository) GetByNotionPageID(projectID uint, pageID string) (*model.Document, error) {
	var document model.Document
	if err := database.DB.Where("project_id = ? AND notion_page_id = ?", projectID, pageID).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *documentRepository) ListByProjectID(projectID uint, offset, limit int) ([]*model.Document, int64, error) {
	var documents []*model.Document
	var total int64

	query := database.DB.Model(&model.Document{}).Where("project_id = ?", projectID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("updated_at DESC").Find(&documents).Error; err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

func (r *documentRepository) Update(document *model.Document) error {
	return database.DB.Save(document).Error
}

func (r *documentRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Document{}, id).Error
}

func (r *documentRepository) DeleteByProjectID(projectID uint) error {
	return database.DB.Where("project_id = ?", projectID).Delete(&model.Document{}).Error
}

