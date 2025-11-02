package repository

import (
	"notion-sync-blog/internal/model"
	"notion-sync-blog/pkg/database"
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
	return database.DB.Create(db).Error
}

func (r *notionDatabaseRepository) GetByID(id uint) (*model.NotionDatabase, error) {
	var db model.NotionDatabase
	if err := database.DB.Preload("Project").First(&db, id).Error; err != nil {
		return nil, err
	}
	return &db, nil
}

func (r *notionDatabaseRepository) GetByNotionDatabaseID(projectID uint, dbID string) (*model.NotionDatabase, error) {
	var db model.NotionDatabase
	if err := database.DB.Where("project_id = ? AND notion_database_id = ?", projectID, dbID).First(&db).Error; err != nil {
		return nil, err
	}
	return &db, nil
}

func (r *notionDatabaseRepository) GetByProjectID(projectID uint) (*model.NotionDatabase, error) {
	var db model.NotionDatabase
	if err := database.DB.Where("project_id = ?", projectID).First(&db).Error; err != nil {
		return nil, err
	}
	return &db, nil
}

func (r *notionDatabaseRepository) Update(db *model.NotionDatabase) error {
	return database.DB.Save(db).Error
}

func (r *notionDatabaseRepository) Delete(id uint) error {
	return database.DB.Delete(&model.NotionDatabase{}, id).Error
}

