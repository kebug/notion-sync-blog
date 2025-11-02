package notion

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// VerifySignature 验证 Notion Webhook 签名
func VerifySignature(secret, body []byte, signature string) error {
	if secret == nil || len(secret) == 0 {
		return errors.New("webhook secret 未配置")
	}

	// 计算 HMAC-SHA256
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// 比较签名
	if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
		return fmt.Errorf("签名验证失败: 期望 %s, 实际 %s", expectedSignature, signature)
	}

	return nil
}

// WebhookEvent Webhook 事件
type WebhookEvent struct {
	Object     string                 `json:"object"`
	Entry      []WebhookEntry         `json:"entry"`
}

// WebhookEntry Webhook 条目
type WebhookEntry struct {
	ID        string                 `json:"id"`
	TimeStamp int64                  `json:"time_stamp"`
	EventType string                 `json:"event_type"`
	Object    map[string]interface{} `json:"object"`
}

