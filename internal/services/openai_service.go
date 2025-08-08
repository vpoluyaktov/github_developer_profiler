package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/russross/blackfriday/v2"
	"github.com/sashabaranov/go-openai"

	"dev_profiler/internal/dto"
)

// OpenAIService handles OpenAI API interactions
type OpenAIService struct {
	client *openai.Client
	config *dto.OpenAIConfig
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(config *dto.OpenAIConfig) *OpenAIService {
	var client *openai.Client
	if config.APIKey != "" {
		client = openai.NewClient(config.APIKey)
	}
	
	return &OpenAIService{
		client: client,
		config: config,
	}
}

// AnalyzeGitHubData sends GitHub audit data to OpenAI for analysis
func (s *OpenAIService) AnalyzeGitHubData(auditResult *dto.AuditResult) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("OpenAI client not initialized - API key required")
	}

	fmt.Printf("[DEBUG] Starting OpenAI analysis...\n")

	// Convert audit result to JSON for the prompt
	auditJSON, err := json.MarshalIndent(auditResult, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal audit result: %w", err)
	}

	fmt.Printf("[DEBUG] Audit JSON size: %d bytes\n", len(auditJSON))

	// Create the system prompt (from your original Python tool)
	systemPrompt := s.getSystemPrompt()
	fmt.Printf("[DEBUG] System prompt size: %d bytes\n", len(systemPrompt))
	
	// Create the user prompt with the audit data
	userPrompt := fmt.Sprintf("Please analyze the following GitHub user data and provide a comprehensive technical assessment:\n\n```json\n%s\n```", string(auditJSON))
	fmt.Printf("[DEBUG] User prompt size: %d bytes\n", len(userPrompt))
	fmt.Printf("[DEBUG] Total prompt size: %d bytes\n", len(systemPrompt)+len(userPrompt))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	fmt.Printf("[DEBUG] Sending request to OpenAI (model: %s)...\n", s.config.Model)

	// Create the chat completion request
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: s.config.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			MaxTokens:   4000,
			Temperature: 0.3,
		},
	)

	if err != nil {
		fmt.Printf("[DEBUG] OpenAI API error: %v\n", err)
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	fmt.Printf("[DEBUG] Received response from OpenAI\n")

	if len(resp.Choices) == 0 {
		fmt.Printf("[DEBUG] No choices in response\n")
		return "", fmt.Errorf("no response from OpenAI")
	}

	responseContent := resp.Choices[0].Message.Content
	fmt.Printf("[DEBUG] Response content size: %d bytes\n", len(responseContent))
	fmt.Printf("[DEBUG] OpenAI analysis completed successfully\n")

	return responseContent, nil
}

// getSystemPrompt returns the configurable system prompt from config
func (s *OpenAIService) getSystemPrompt() string {
	// Use the configurable system prompt from config, fallback to default if empty
	if s.config.SystemPrompt != "" {
		return s.config.SystemPrompt
	}
	// Fallback to default system prompt if not configured
	return dto.DefaultSystemPrompt()
}

// ConvertMarkdownToHTML converts markdown content to a complete HTML document using configurable templates
func (s *OpenAIService) ConvertMarkdownToHTML(markdownContent string, username string) string {
	// Convert markdown to HTML using blackfriday
	htmlBytes := blackfriday.Run([]byte(markdownContent))
	htmlContent := string(htmlBytes)
	
	// Get the configurable HTML template and CSS styles
	htmlTemplate := s.getHTMLTemplate()
	cssStyles := s.getCSSStyles()
	
	// Parse the HTML template
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		// Fallback to a simple template if parsing fails
		return s.fallbackHTMLGeneration(htmlContent, username)
	}
	
	// Prepare template data
	templateData := struct {
		Username  string
		Content   string
		CSSStyles string
	}{
		Username:  username,
		Content:   htmlContent,
		CSSStyles: cssStyles,
	}
	
	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		// Fallback to simple template if execution fails
		return s.fallbackHTMLGeneration(htmlContent, username)
	}
	
	return buf.String()
}

// getHTMLTemplate returns the configurable HTML template from config
func (s *OpenAIService) getHTMLTemplate() string {
	// Use the configurable HTML template from config, fallback to default if empty
	if s.config.HTMLTemplate != "" {
		return s.config.HTMLTemplate
	}
	// Fallback to default HTML template if not configured
	return dto.DefaultHTMLTemplate()
}

// getCSSStyles returns the configurable CSS styles from config
func (s *OpenAIService) getCSSStyles() string {
	// Use the configurable CSS styles from config, fallback to default if empty
	if s.config.CSSStyles != "" {
		return s.config.CSSStyles
	}
	// Fallback to default CSS styles if not configured
	return dto.DefaultCSSStyles()
}

// fallbackHTMLGeneration provides a simple HTML generation fallback
func (s *OpenAIService) fallbackHTMLGeneration(htmlContent, username string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GitHub Developer Assessment - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
        .container { max-width: 1200px; margin: 0 auto; }
        h1, h2, h3 { color: #333; }
        table { border-collapse: collapse; width: 100%%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="container">
        <h1>GitHub Developer Assessment - %s</h1>
        %s
    </div>
</body>
</html>`, username, username, htmlContent)
}

// Helper functions for HTML processing (no longer needed since OpenAI returns HTML directly)
func convertHeaders(html string) string {
	lines := strings.Split(html, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			lines[i] = "<h1>" + strings.TrimPrefix(line, "# ") + "</h1>"
		} else if strings.HasPrefix(line, "## ") {
			lines[i] = "<h2>" + strings.TrimPrefix(line, "## ") + "</h2>"
		} else if strings.HasPrefix(line, "### ") {
			lines[i] = "<h3>" + strings.TrimPrefix(line, "### ") + "</h3>"
		}
	}
	return strings.Join(lines, "\n")
}

func convertTables(html string) string {
	lines := strings.Split(html, "\n")
	var result []string
	inTable := false
	
	for i, line := range lines {
		if strings.Contains(line, "|") && !inTable {
			// Start of table
			inTable = true
			result = append(result, "<table>")
			if i+1 < len(lines) && strings.Contains(lines[i+1], "---") {
				// Header row
				result = append(result, "<thead>")
				result = append(result, convertTableRow(line, "th"))
				result = append(result, "</thead>")
				result = append(result, "<tbody>")
			} else {
				result = append(result, "<tbody>")
				result = append(result, convertTableRow(line, "td"))
			}
		} else if strings.Contains(line, "|") && inTable {
			if !strings.Contains(line, "---") {
				result = append(result, convertTableRow(line, "td"))
			}
		} else if inTable && !strings.Contains(line, "|") {
			// End of table
			inTable = false
			result = append(result, "</tbody>")
			result = append(result, "</table>")
			result = append(result, line)
		} else {
			result = append(result, line)
		}
	}
	
	if inTable {
		result = append(result, "</tbody>")
		result = append(result, "</table>")
	}
	
	return strings.Join(result, "\n")
}

func convertTableRow(line, cellType string) string {
	cells := strings.Split(line, "|")
	var result []string
	result = append(result, "<tr>")
	for _, cell := range cells {
		cell = strings.TrimSpace(cell)
		if cell != "" {
			result = append(result, fmt.Sprintf("<%s>%s</%s>", cellType, cell, cellType))
		}
	}
	result = append(result, "</tr>")
	return strings.Join(result, "")
}

func convertCodeBlocks(html string) string {
	// Convert code blocks (```language ... ```)
	lines := strings.Split(html, "\n")
	var result []string
	inCodeBlock := false
	
	for _, line := range lines {
		if strings.HasPrefix(line, "```") && !inCodeBlock {
			inCodeBlock = true
			result = append(result, "<pre><code>")
		} else if strings.HasPrefix(line, "```") && inCodeBlock {
			inCodeBlock = false
			result = append(result, "</code></pre>")
		} else if inCodeBlock {
			result = append(result, line)
		} else {
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}
