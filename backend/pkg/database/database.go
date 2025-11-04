package database

import (
	"fmt"
	"time"

	"notion-sync-blog/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB 全局数据库实例
var DB *sqlx.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := cfg.GetDSN()

	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 设置连接池
	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.MaxIdleConns)
	DB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
