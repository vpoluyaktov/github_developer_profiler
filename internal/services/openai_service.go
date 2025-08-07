package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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

// ConvertMarkdownToHTML converts markdown content to a complete HTML document
func (s *OpenAIService) ConvertMarkdownToHTML(markdownContent string, username string) string {
	// Convert markdown to HTML using blackfriday
	htmlBytes := blackfriday.Run([]byte(markdownContent))
	htmlContent := string(htmlBytes)
	
	// Wrap in HTML document with improved styling
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GitHub Developer Assessment - %s</title>
    
    <!-- Prism.js CSS for syntax highlighting -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism.min.css" rel="stylesheet" />
    <link href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/line-numbers/prism-line-numbers.min.css" rel="stylesheet" />
    
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.7;
            color: #2c3e50;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
        }
        .container {
            background-color: white;
            padding: 50px;
            border-radius: 12px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.15);
            margin: 20px 0;
        }
        .header {
            text-align: center;
            margin-bottom: 40px;
            padding-bottom: 30px;
            border-bottom: 3px solid #3498db;
        }
        h1 { 
            color: #2c3e50; 
            font-size: 2.5em;
            margin-bottom: 10px;
            font-weight: 700;
        }
        h2 { 
            color: #34495e; 
            border-left: 4px solid #3498db;
            padding-left: 20px;
            margin: 40px 0 20px 0;
            font-size: 1.8em;
            font-weight: 600;
        }
        h3 { 
            color: #7f8c8d; 
            margin: 30px 0 15px 0;
            font-size: 1.4em;
            font-weight: 600;
        }
        h4 {
            color: #95a5a6;
            margin: 25px 0 10px 0;
            font-size: 1.2em;
            font-weight: 600;
        }
        p {
            margin: 12px 0;
            text-align: left;
            line-height: 1.6;
        }
        /* Better formatting for summary sections */
        h2 + p, h3 + p {
            margin-top: 8px;
        }
        /* Improve readability for recommendation text */
        p:has(strong) {
            margin: 16px 0;
            padding: 12px;
            background-color: #f8f9fa;
            border-left: 4px solid #007bff;
            border-radius: 4px;
        }
        ul, ol {
            margin: 15px 0;
            padding-left: 30px;
        }
        li {
            margin: 8px 0;
        }
        table {
            border-collapse: collapse;
            width: 100%%;
            margin: 25px 0;
            background-color: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        th {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            font-weight: 600;
            padding: 8px 10px;
            text-align: left;
            font-size: 0.9em;
        }
        td {
            padding: 6px 10px;
            border-bottom: 1px solid #ecf0f1;
            font-size: 0.9em;
            line-height: 1.4;
        }
        tr:nth-child(even) {
            background-color: #f8f9fa;
        }
        tr:hover {
            background-color: #e8f4fd;
        }
        /* Override Prism.js default styles for better integration */
        pre[class*="language-"] {
            background-color: #f8f9fa !important;
            border: 1px solid #dee2e6;
            border-radius: 6px;
            padding: 15px !important;
            overflow-x: auto;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
            font-size: 0.85em;
            line-height: 1.4;
            margin: 15px 0 !important;
            box-shadow: inset 0 1px 3px rgba(0,0,0,0.1);
        }
        code[class*="language-"] {
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
            font-size: 0.85em;
            line-height: 1.4;
        }
        /* Fallback for code blocks without language specification */
        pre:not([class*="language-"]) {
            background-color: #f8f9fa;
            color: #2c3e50;
            border: 1px solid #dee2e6;
            border-radius: 6px;
            padding: 15px;
            overflow-x: auto;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
            font-size: 0.85em;
            line-height: 1.4;
            margin: 15px 0;
            box-shadow: inset 0 1px 3px rgba(0,0,0,0.1);
        }
        pre:not([class*="language-"]) code {
            background-color: transparent;
            color: inherit;
            padding: 0;
            border-radius: 0;
            font-size: inherit;
        }
        /* Inline code styling */
        code:not([class*="language-"]) {
            background-color: #f1f3f4;
            color: #d73a49;
            padding: 2px 4px;
            border-radius: 3px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
            font-size: 0.85em;
        }
        .highlight {
            background-color: #fff3cd;
            padding: 15px;
            border-left: 4px solid #ffc107;
            margin: 20px 0;
            border-radius: 4px;
        }
        .print-button {
            position: fixed;
            top: 30px;
            right: 30px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 600;
            box-shadow: 0 4px 15px rgba(0,0,0,0.2);
            transition: all 0.3s ease;
        }
        .print-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(0,0,0,0.3);
        }
        @media print {
            body { 
                background: white;
                color: black;
            }
            .container { 
                box-shadow: none;
                padding: 20px;
            }
            .print-button { display: none; }
            h1, h2, h3 { color: black; }
        }
    </style>
</head>
<body>
    <button class="print-button" onclick="window.print()">🖨️ Print Report</button>
    <div class="container">
        <div class="header">
            <h1>GitHub Developer Assessment</h1>
            <p style="font-size: 1.2em; color: #7f8c8d; margin: 0;">Professional Technical Evaluation for <strong>%s</strong></p>
        </div>
        %s
    </div>
    
    <!-- Prism.js JavaScript for syntax highlighting -->
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-core.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/line-numbers/prism-line-numbers.min.js"></script>
    
    <script>
        // Initialize Prism.js after page load
        document.addEventListener('DOMContentLoaded', function() {
            Prism.highlightAll();
        });
    </script>
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
