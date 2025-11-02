package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client Notion API 客户端
type Client struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

// NewClient 创建 Notion API 客户端
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.notion.com/v1",
	}
}

// request 执行 HTTP 请求
func (c *Client) request(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API 请求失败: %d - %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetPage 获取页面信息
func (c *Client) GetPage(pageID string) (*Page, error) {
	data, err := c.request("GET", "/pages/"+pageID, nil)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("解析页面数据失败: %w", err)
	}

	return &page, nil
}

// GetPageBlocks 获取页面块内容
func (c *Client) GetPageBlocks(pageID string) ([]Block, error) {
	endpoint := fmt.Sprintf("/blocks/%s/children", pageID)
	
	var allBlocks []Block
	var nextCursor string
	
	for {
		var url string
		if nextCursor != "" {
			url = fmt.Sprintf("%s?start_cursor=%s", endpoint, nextCursor)
		} else {
			url = endpoint
		}

		data, err := c.request("GET", url, nil)
		if err != nil {
			return nil, err
		}

		var response BlocksResponse
		if err := json.Unmarshal(data, &response); err != nil {
			return nil, fmt.Errorf("解析块数据失败: %w", err)
		}

		allBlocks = append(allBlocks, response.Results...)

		if !response.HasMore {
			break
		}
		nextCursor = response.NextCursor
	}

	return allBlocks, nil
}

// GetDatabase 获取数据库信息
func (c *Client) GetDatabase(databaseID string) (*Database, error) {
	data, err := c.request("GET", "/databases/"+databaseID, nil)
	if err != nil {
		return nil, err
	}

	var db Database
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("解析数据库数据失败: %w", err)
	}

	return &db, nil
}

// QueryDatabase 查询数据库
func (c *Client) QueryDatabase(databaseID string) ([]Page, error) {
	endpoint := fmt.Sprintf("/databases/%s/query", databaseID)
	
	var allPages []Page
	var nextCursor string
	
	for {
		var body map[string]interface{}
		if nextCursor != "" {
			body = map[string]interface{}{
				"start_cursor": nextCursor,
			}
		} else {
			body = make(map[string]interface{})
		}

		data, err := c.request("POST", endpoint, body)
		if err != nil {
			return nil, err
		}

		var response PagesResponse
		if err := json.Unmarshal(data, &response); err != nil {
			return nil, fmt.Errorf("解析数据库查询结果失败: %w", err)
		}

		allPages = append(allPages, response.Results...)

		if !response.HasMore {
			break
		}
		nextCursor = response.NextCursor
	}

	return allPages, nil
}

