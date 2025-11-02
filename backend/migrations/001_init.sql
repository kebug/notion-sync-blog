-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS notion_sync_blog DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE notion_sync_blog;

-- 创建 projects 表
CREATE TABLE IF NOT EXISTS projects (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    notion_root_page_id VARCHAR(255) NOT NULL,
    notion_root_page_title VARCHAR(500) DEFAULT '',
    notion_token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_notion_root_page_id (notion_root_page_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建 notion_documents 表
CREATE TABLE IF NOT EXISTS notion_documents (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL,
    notion_page_id VARCHAR(255) NOT NULL,
    notion_page_title VARCHAR(500) DEFAULT '',
    content_path VARCHAR(1000) DEFAULT '',
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    INDEX idx_notion_page_id (notion_page_id),
    INDEX idx_status (status),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建 notion_databases 表
CREATE TABLE IF NOT EXISTS notion_databases (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL,
    notion_database_id VARCHAR(255) NOT NULL,
    cached_data LONGTEXT,
    last_synced_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    INDEX idx_notion_database_id (notion_database_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建 sync_logs 表
CREATE TABLE IF NOT EXISTS sync_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    project_id BIGINT UNSIGNED NOT NULL,
    event_type VARCHAR(100) DEFAULT '',
    event_id VARCHAR(255) DEFAULT '',
    status VARCHAR(50) DEFAULT 'pending',
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    INDEX idx_event_type (event_type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

