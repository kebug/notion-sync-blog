package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"notion-sync-blog/config"

	"github.com/go-redis/redis/v8"
)

// Cache Redis 缓存客户端
type Cache struct {
	client *redis.Client
	ctx    context.Context
}

// NewCache 创建 Redis 缓存实例
func NewCache(cfg *config.RedisConfig) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	ctx := context.Background()
	
	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	return &Cache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Close 关闭 Redis 连接
func (c *Cache) Close() error {
	return c.client.Close()
}

// Set 设置缓存
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}
	return c.client.Set(c.ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (c *Cache) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("缓存不存在")
		}
		return fmt.Errorf("获取缓存失败: %w", err)
	}
	return json.Unmarshal(data, dest)
}

// Delete 删除缓存
func (c *Cache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Exists 检查缓存是否存在
func (c *Cache) Exists(key string) (bool, error) {
	count, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetDatabaseKey 获取数据库缓存 key
func GetDatabaseKey(projectID uint, databaseID string) string {
	return fmt.Sprintf("notion:database:%d:%s", projectID, databaseID)
}

// GetPageKey 获取页面缓存 key
func GetPageKey(projectID uint, pageID string) string {
	return fmt.Sprintf("notion:page:%d:%s", projectID, pageID)
}

