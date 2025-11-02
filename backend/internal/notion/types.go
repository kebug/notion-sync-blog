package notion

import "encoding/json"

// Page Notion 页面
type Page struct {
	ID         string                 `json:"id"`
	Object     string                 `json:"object"`
	CreatedTime string                `json:"created_time"`
	LastEditedTime string             `json:"last_edited_time"`
	Properties map[string]interface{} `json:"properties"`
	Title      []RichText             `json:"title,omitempty"`
}

// Database Notion 数据库
type Database struct {
	ID         string                 `json:"id"`
	Object     string                 `json:"object"`
	Title      []RichText             `json:"title"`
	Properties map[string]interface{} `json:"properties"`
}

// Block Notion 块
type Block struct {
	ID         string                 `json:"id"`
	Object     string                 `json:"object"`
	Type       string                 `json:"type"`
	CreatedTime string                `json:"created_time"`
	LastEditedTime string             `json:"last_edited_time"`
	HasChildren bool                  `json:"has_children"`
	Children   []Block                `json:"children,omitempty"`
	Content    map[string]interface{} `json:"-"`
}

// UnmarshalJSON 自定义反序列化
func (b *Block) UnmarshalJSON(data []byte) error {
	type Alias Block
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// 提取类型特定的内容
	if content, ok := raw[b.Type].(map[string]interface{}); ok {
		b.Content = content
	}

	return nil
}

// RichText 富文本
type RichText struct {
	Type        string      `json:"type"`
	PlainText   string      `json:"plain_text"`
	Href        *string     `json:"href,omitempty"`
	Annotations *Annotation `json:"annotations,omitempty"`
}

// Annotation 文本注解
type Annotation struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

// BlocksResponse 块列表响应
type BlocksResponse struct {
	Object     string  `json:"object"`
	Results    []Block `json:"results"`
	NextCursor string  `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
}

// PagesResponse 页面列表响应
type PagesResponse struct {
	Object     string  `json:"object"`
	Results    []Page  `json:"results"`
	NextCursor string  `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
}

// GetTitle 获取页面标题
func (p *Page) GetTitle() string {
	for _, title := range p.Title {
		return title.PlainText
	}
	// 尝试从 properties 中获取标题
	if titleProp, ok := p.Properties["title"]; ok {
		if titleObj, ok := titleProp.(map[string]interface{}); ok {
			if titleArray, ok := titleObj["title"].([]interface{}); ok {
				for _, item := range titleArray {
					if itemObj, ok := item.(map[string]interface{}); ok {
						if plainText, ok := itemObj["plain_text"].(string); ok {
							return plainText
						}
					}
				}
			}
		}
	}
	return "Untitled"
}

