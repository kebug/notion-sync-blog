package service

import (
	"fmt"
	"os"
	"path/filepath"

	"notion-sync-blog/config"
	"notion-sync-blog/internal/model"
	"notion-sync-blog/internal/repository"
)

// ProjectService 项目服务接口
type ProjectService interface {
	CreateProject(rootPageID, rootPageTitle, token string) (*model.Project, error)
	GetProject(id uint) (*model.Project, error)
	ListProjects(offset, limit int) ([]*model.Project, int64, error)
	UpdateProject(id uint, title string) error
	DeleteProject(id uint) error
	GetProjectToken(projectID uint) (string, error)
}

type projectService struct {
	projectRepo repository.ProjectRepository
	storageCfg  *config.StorageConfig
}

// NewProjectService 创建项目服务
func NewProjectService(
	projectRepo repository.ProjectRepository,
	storageCfg *config.StorageConfig,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		storageCfg:  storageCfg,
	}
}

func (s *projectService) CreateProject(rootPageID, rootPageTitle, token string) (*model.Project, error) {
	// 检查是否已存在
	existing, _ := s.projectRepo.GetByNotionPageID(rootPageID)
	if existing != nil {
		return nil, fmt.Errorf("项目已存在")
	}

	project := &model.Project{
		NotionRootPageID:    rootPageID,
		NotionRootPageTitle: rootPageTitle,
		NotionToken:         token, // TODO: 加密存储
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, fmt.Errorf("创建项目失败: %w", err)
	}

	// 创建项目存储目录
	if err := s.createProjectStorage(project.ID); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	return project, nil
}

func (s *projectService) GetProject(id uint) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取项目失败: %w", err)
	}
	// 不返回 token
	project.NotionToken = ""
	return project, nil
}

func (s *projectService) ListProjects(offset, limit int) ([]*model.Project, int64, error) {
	projects, total, err := s.projectRepo.List(offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("获取项目列表失败: %w", err)
	}

	// 清除 token
	for _, p := range projects {
		p.NotionToken = ""
	}

	return projects, total, nil
}

func (s *projectService) UpdateProject(id uint, title string) error {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("项目不存在: %w", err)
	}

	project.NotionRootPageTitle = title
	if err := s.projectRepo.Update(project); err != nil {
		return fmt.Errorf("更新项目失败: %w", err)
	}

	return nil
}

func (s *projectService) DeleteProject(id uint) error {
	// 删除存储目录
	if err := s.deleteProjectStorage(id); err != nil {
		return fmt.Errorf("删除存储目录失败: %w", err)
	}

	if err := s.projectRepo.Delete(id); err != nil {
		return fmt.Errorf("删除项目失败: %w", err)
	}

	return nil
}

func (s *projectService) GetProjectToken(projectID uint) (string, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return "", fmt.Errorf("项目不存在: %w", err)
	}
	return project.NotionToken, nil
}

func (s *projectService) createProjectStorage(projectID uint) error {
	basePath := filepath.Join(s.storageCfg.BasePath, fmt.Sprintf("%d", projectID))
	dirs := []string{
		filepath.Join(basePath, "database"),
		filepath.Join(basePath, "pages"),
		filepath.Join(basePath, "generated", "astro"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (s *projectService) deleteProjectStorage(projectID uint) error {
	basePath := filepath.Join(s.storageCfg.BasePath, fmt.Sprintf("%d", projectID))
	return os.RemoveAll(basePath)
}

