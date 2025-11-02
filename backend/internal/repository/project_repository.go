package repository

import (
	"notion-sync-blog/internal/model"
	"notion-sync-blog/pkg/database"
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
	return database.DB.Create(project).Error
}

func (r *projectRepository) GetByID(id uint) (*model.Project, error) {
	var project model.Project
	if err := database.DB.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetByNotionPageID(pageID string) (*model.Project, error) {
	var project model.Project
	if err := database.DB.Where("notion_root_page_id = ?", pageID).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) List(offset, limit int) ([]*model.Project, int64, error) {
	var projects []*model.Project
	var total int64

	if err := database.DB.Model(&model.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Offset(offset).Limit(limit).Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *projectRepository) Update(project *model.Project) error {
	return database.DB.Save(project).Error
}

func (r *projectRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Project{}, id).Error
}

