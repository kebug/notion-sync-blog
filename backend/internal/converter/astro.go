package converter

import (
	"fmt"
	"strings"
	"time"

	"notion-sync-blog/internal/notion"
)

// AstroConverter Astro 格式转换器
type AstroConverter struct{}

// NewAstroConverter 创建 Astro 转换器
func NewAstroConverter() *AstroConverter {
	return &AstroConverter{}
}

// ConvertPage 转换页面为 Astro 格式
func (c *AstroConverter) ConvertPage(page *notion.Page, blocks []notion.Block, status string) (string, error) {
	var builder strings.Builder

	// Front Matter
	builder.WriteString("---\n")
	builder.WriteString(fmt.Sprintf("title: %q\n", escapeYAML(page.GetTitle())))
	
	if page.CreatedTime != "" {
		if created, err := parseTime(page.CreatedTime); err == nil {
			builder.WriteString(fmt.Sprintf("created: %s\n", created.Format(time.RFC3339)))
		}
	}
	
	if page.LastEditedTime != "" {
		if updated, err := parseTime(page.LastEditedTime); err == nil {
			builder.WriteString(fmt.Sprintf("updated: %s\n", updated.Format(time.RFC3339)))
		}
	}
	
	builder.WriteString(fmt.Sprintf("status: %q\n", status))
	builder.WriteString("---\n\n")

	// 转换块内容为 Markdown
	content, err := c.convertBlocks(blocks)
	if err != nil {
		return "", fmt.Errorf("转换块内容失败: %w", err)
	}

	builder.WriteString(content)
	return builder.String(), nil
}

// ConvertBlocks 转换块列表为 Markdown（公共方法）
func (c *AstroConverter) ConvertBlocks(blocks []notion.Block) (string, error) {
	return c.convertBlocks(blocks)
}

// convertBlocks 转换块列表为 Markdown
func (c *AstroConverter) convertBlocks(blocks []notion.Block) (string, error) {
	var builder strings.Builder

	for _, block := range blocks {
		markdown, err := c.convertBlock(block)
		if err != nil {
			continue // 跳过无法转换的块
		}
		builder.WriteString(markdown)
		builder.WriteString("\n")

		// 递归处理子块
		if block.HasChildren && len(block.Children) > 0 {
			childContent, err := c.convertBlocks(block.Children)
			if err == nil {
				builder.WriteString(childContent)
			}
		}
	}

	return builder.String(), nil
}

// convertBlock 转换单个块为 Markdown
func (c *AstroConverter) convertBlock(block notion.Block) (string, error) {
	switch block.Type {
	case "heading_1":
		return c.convertHeading(block, 1), nil
	case "heading_2":
		return c.convertHeading(block, 2), nil
	case "heading_3":
		return c.convertHeading(block, 3), nil
	case "paragraph":
		return c.convertParagraph(block), nil
	case "bulleted_list_item":
		return c.convertBulletedList(block), nil
	case "numbered_list_item":
		return c.convertNumberedList(block), nil
	case "code":
		return c.convertCode(block), nil
	case "image":
		return c.convertImage(block), nil
	default:
		return "", fmt.Errorf("不支持的块类型: %s", block.Type)
	}
}

// convertHeading 转换标题
func (c *AstroConverter) convertHeading(block notion.Block, level int) string {
	text := c.extractRichText(block, "rich_text")
	prefix := strings.Repeat("#", level)
	return fmt.Sprintf("%s %s", prefix, text)
}

// convertParagraph 转换段落
func (c *AstroConverter) convertParagraph(block notion.Block) string {
	return c.extractRichText(block, "rich_text")
}

// convertBulletedList 转换无序列表
func (c *AstroConverter) convertBulletedList(block notion.Block) string {
	text := c.extractRichText(block, "rich_text")
	return fmt.Sprintf("- %s", text)
}

// convertNumberedList 转换有序列表
func (c *AstroConverter) convertNumberedList(block notion.Block) string {
	text := c.extractRichText(block, "rich_text")
	return fmt.Sprintf("1. %s", text)
}

// convertCode 转换代码块
func (c *AstroConverter) convertCode(block notion.Block) string {
	code := c.extractRichText(block, "rich_text")
	language := c.extractString(block, "language")
	
	if language != "" {
		return fmt.Sprintf("```%s\n%s\n```", language, code)
	}
	return fmt.Sprintf("```\n%s\n```", code)
}

// convertImage 转换图片
func (c *AstroConverter) convertImage(block notion.Block) string {
	imageURL := c.extractImageURL(block)
	caption := c.extractRichText(block, "caption")
	
	if caption != "" {
		return fmt.Sprintf("![%s](%s)", caption, imageURL)
	}
	return fmt.Sprintf("![image](%s)", imageURL)
}

// extractRichText 提取富文本
func (c *AstroConverter) extractRichText(block notion.Block, field string) string {
	if block.Content == nil {
		return ""
	}

	richTextArray, ok := block.Content[field].([]interface{})
	if !ok {
		return ""
	}

	var texts []string
	for _, item := range richTextArray {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if plainText, ok := itemMap["plain_text"].(string); ok {
				texts = append(texts, plainText)
			}
		}
	}

	return strings.Join(texts, "")
}

// extractString 提取字符串字段
func (c *AstroConverter) extractString(block notion.Block, field string) string {
	if block.Content == nil {
		return ""
	}

	if str, ok := block.Content[field].(string); ok {
		return str
	}
	return ""
}

// extractImageURL 提取图片 URL
func (c *AstroConverter) extractImageURL(block notion.Block) string {
	if block.Content == nil {
		return ""
	}

	if fileObj, ok := block.Content["file"].(map[string]interface{}); ok {
		if url, ok := fileObj["url"].(string); ok {
			return url
		}
	}

	if externalObj, ok := block.Content["external"].(map[string]interface{}); ok {
		if url, ok := externalObj["url"].(string); ok {
			return url
		}
	}

	return ""
}

// escapeYAML 转义 YAML 字符串
func escapeYAML(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

// parseTime 解析时间字符串
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

