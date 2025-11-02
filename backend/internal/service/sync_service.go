package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"notion-sync-blog/config"
	"notion-sync-blog/internal/cache"
	"notion-sync-blog/internal/converter"
	"notion-sync-blog/internal/model"
	"notion-sync-blog/internal/notion"
	"notion-sync-blog/internal/repository"
)

// SyncService 同步服务接口
type SyncService interface {
	SyncDatabase(projectID uint) error
	SyncPage(projectID uint, pageID string) error
	HandleWebhookEvent(projectID uint, eventType, eventID string, eventObject map[string]interface{}) error
	ConvertToAstro(projectID uint) error
}

type syncService struct {
	projectRepo    repository.ProjectRepository
	documentRepo   repository.DocumentRepository
	databaseRepo   repository.NotionDatabaseRepository
	syncLogRepo    repository.SyncLogRepository
	notionClient   *notion.Client
	cache          *cache.Cache
	converter      *converter.AstroConverter
	storageCfg     *config.StorageConfig
}

// NewSyncService 创建同步服务
func NewSyncService(
	projectRepo repository.ProjectRepository,
	documentRepo repository.DocumentRepository,
	databaseRepo repository.NotionDatabaseRepository,
	syncLogRepo repository.SyncLogRepository,
	cache *cache.Cache,
	storageCfg *config.StorageConfig,
) SyncService {
	return &syncService{
		projectRepo:  projectRepo,
		documentRepo: documentRepo,
		databaseRepo: databaseRepo,
		syncLogRepo:  syncLogRepo,
		cache:        cache,
		converter:    converter.NewAstroConverter(),
		storageCfg:   storageCfg,
	}
}

func (s *syncService) SyncDatabase(projectID uint) error {
	// 创建同步日志
	log := &model.SyncLog{
		ProjectID: projectID,
		EventType: "database.sync",
		Status:    "pending",
	}
	if err := s.syncLogRepo.Create(log); err != nil {
		return fmt.Errorf("创建同步日志失败: %w", err)
	}

	// 获取项目信息
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		s.updateLogError(log.ID, err)
		return err
	}

	// 创建 Notion 客户端
	s.notionClient = notion.NewClient(project.NotionToken)

	// 获取根页面下的数据库（简化：假设根页面包含数据库）
	// 这里需要根据实际需求查找数据库
	// 暂时跳过数据库同步的具体实现，因为需要知道如何从根页面查找数据库

	s.updateLogSuccess(log.ID)
	return nil
}

func (s *syncService) SyncPage(projectID uint, pageID string) error {
	// 创建同步日志
	log := &model.SyncLog{
		ProjectID: projectID,
		EventType: "page.sync",
		EventID:   pageID,
		Status:    "pending",
	}
	if err := s.syncLogRepo.Create(log); err != nil {
		return fmt.Errorf("创建同步日志失败: %w", err)
	}

	// 获取项目信息
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		s.updateLogError(log.ID, err)
		return err
	}

	// 创建 Notion 客户端
	client := notion.NewClient(project.NotionToken)

	// 获取页面信息
	page, err := client.GetPage(pageID)
	if err != nil {
		s.updateLogError(log.ID, err)
		return fmt.Errorf("获取页面失败: %w", err)
	}

	// 获取页面块
	blocks, err := client.GetPageBlocks(pageID)
	if err != nil {
		s.updateLogError(log.ID, err)
		return fmt.Errorf("获取页面块失败: %w", err)
	}

	// 保存到缓存
	cacheKey := cache.GetPageKey(projectID, pageID)
	pageData := map[string]interface{}{
		"page":   page,
		"blocks": blocks,
	}
	if err := s.cache.Set(cacheKey, pageData, 24*time.Hour); err != nil {
		// 缓存失败不影响主流程
		fmt.Printf("缓存页面失败: %v\n", err)
	}

	// 保存原始 JSON
	if err := s.savePageJSON(projectID, pageID, page, blocks); err != nil {
		s.updateLogError(log.ID, err)
		return fmt.Errorf("保存页面 JSON 失败: %w", err)
	}

	// 保存 Markdown（简化版）
	markdown, _ := s.converter.ConvertBlocks(blocks)
	if markdown != "" {
		mdPath := filepath.Join(s.storageCfg.BasePath, fmt.Sprintf("%d/pages", projectID), pageID+".md")
		os.WriteFile(mdPath, []byte(markdown), 0644)
	}

	// 更新或创建文档记录
	doc, _ := s.documentRepo.GetByNotionPageID(projectID, pageID)
	if doc == nil {
		doc = &model.Document{
			ProjectID:      projectID,
			NotionPageID:   pageID,
			NotionPageTitle: page.GetTitle(),
			ContentPath:    s.getPagePath(projectID, pageID),
			Status:         "draft",
		}
		if err := s.documentRepo.Create(doc); err != nil {
			s.updateLogError(log.ID, err)
			return fmt.Errorf("创建文档记录失败: %w", err)
		}
	} else {
		doc.NotionPageTitle = page.GetTitle()
		doc.ContentPath = s.getPagePath(projectID, pageID)
		if err := s.documentRepo.Update(doc); err != nil {
			s.updateLogError(log.ID, err)
			return fmt.Errorf("更新文档记录失败: %w", err)
		}
	}

	s.updateLogSuccess(log.ID)
	return nil
}

func (s *syncService) HandleWebhookEvent(projectID uint, eventType, eventID string, eventObject map[string]interface{}) error {
	// 创建同步日志
	log := &model.SyncLog{
		ProjectID: projectID,
		EventType: eventType,
		EventID:   eventID,
		Status:    "pending",
	}
	if err := s.syncLogRepo.Create(log); err != nil {
		return fmt.Errorf("创建同步日志失败: %w", err)
	}

	switch eventType {
	case "page.updated":
		// 从事件对象中提取页面 ID
		if pageID, ok := eventObject["id"].(string); ok {
			if err := s.SyncPage(projectID, pageID); err != nil {
				s.updateLogError(log.ID, err)
				return err
			}
		}
	case "database.updated":
		if err := s.SyncDatabase(projectID); err != nil {
			s.updateLogError(log.ID, err)
			return err
		}
	default:
		s.updateLogError(log.ID, fmt.Errorf("不支持的事件类型: %s", eventType))
		return fmt.Errorf("不支持的事件类型: %s", eventType)
	}

	return nil
}

func (s *syncService) ConvertToAstro(projectID uint) error {
	// 获取项目所有文档
	documents, _, err := s.documentRepo.ListByProjectID(projectID, 0, 1000)
	if err != nil {
		return fmt.Errorf("获取文档列表失败: %w", err)
	}

	// 获取项目信息
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return fmt.Errorf("获取项目失败: %w", err)
	}

	client := notion.NewClient(project.NotionToken)

	for _, doc := range documents {
		// 获取页面块
		blocks, err := client.GetPageBlocks(doc.NotionPageID)
		if err != nil {
			fmt.Printf("获取页面块失败 %s: %v\n", doc.NotionPageID, err)
			continue
		}

		// 获取页面信息
		page, err := client.GetPage(doc.NotionPageID)
		if err != nil {
			fmt.Printf("获取页面失败 %s: %v\n", doc.NotionPageID, err)
			continue
		}

		// 转换为 Astro 格式
		astroContent, err := s.converter.ConvertPage(page, blocks, doc.Status)
		if err != nil {
			fmt.Printf("转换页面失败 %s: %v\n", doc.NotionPageID, err)
			continue
		}

		// 保存 Astro 文件
		astroPath := filepath.Join(s.storageCfg.BasePath, fmt.Sprintf("%d/generated/astro", projectID), doc.NotionPageID+".md")
		if err := os.WriteFile(astroPath, []byte(astroContent), 0644); err != nil {
			fmt.Printf("保存 Astro 文件失败 %s: %v\n", astroPath, err)
			continue
		}
	}

	return nil
}

// 辅助方法

func (s *syncService) savePageJSON(projectID uint, pageID string, page *notion.Page, blocks []notion.Block) error {
	pageData := map[string]interface{}{
		"page":   page,
		"blocks": blocks,
	}

	jsonData, err := json.MarshalIndent(pageData, "", "  ")
	if err != nil {
		return err
	}

	jsonPath := filepath.Join(s.storageCfg.BasePath, fmt.Sprintf("%d/pages", projectID), pageID+".json")
	return os.WriteFile(jsonPath, jsonData, 0644)
}


func (s *syncService) getPagePath(projectID uint, pageID string) string {
	return filepath.Join(s.storageCfg.BasePath, fmt.Sprintf("%d/pages", projectID), pageID+".md")
}

func (s *syncService) updateLogSuccess(logID uint) {
	log, _ := s.syncLogRepo.GetByID(logID)
	if log != nil {
		log.Status = "success"
		s.syncLogRepo.Update(log)
	}
}

func (s *syncService) updateLogError(logID uint, err error) {
	log, _ := s.syncLogRepo.GetByID(logID)
	if log != nil {
		log.Status = "failed"
		log.ErrorMessage = err.Error()
		s.syncLogRepo.Update(log)
	}
}

