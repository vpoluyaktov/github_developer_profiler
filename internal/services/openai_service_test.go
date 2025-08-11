package services

import (
	"testing"

	"dev_profiler/internal/config"
	"dev_profiler/internal/dto"
)

func TestNewOpenAIService(t *testing.T) {
	cfg := &config.OpenAIConfig{
		APIKey: "test-key",
		Model:  "gpt-4",
	}
	
	service := NewOpenAIService(cfg)
	
	if service == nil {
		t.Fatal("NewOpenAIService() returned nil")
	}
	
	if service.config != cfg {
		t.Error("Service config should match provided config")
	}
}

func TestNewOpenAIServiceWithEmptyKey(t *testing.T) {
	cfg := &config.OpenAIConfig{
		APIKey: "",
		Model:  "gpt-4",
	}
	
	service := NewOpenAIService(cfg)
	
	if service == nil {
		t.Fatal("NewOpenAIService() returned nil")
	}
	
	// Client should be nil when no API key is provided
	if service.client != nil {
		t.Error("Client should be nil when API key is empty")
	}
}

func TestGetSystemPrompt(t *testing.T) {
	testCases := []struct {
		name           string
		configPrompt   string
		expectedResult string
	}{
		{
			name:           "with custom prompt",
			configPrompt:   "Custom system prompt",
			expectedResult: "Custom system prompt",
		},
		{
			name:           "with empty prompt",
			configPrompt:   "",
			expectedResult: config.DefaultSystemPrompt(),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.OpenAIConfig{
				SystemPrompt: tc.configPrompt,
			}
			
			service := NewOpenAIService(cfg)
			result := service.getSystemPrompt()
			
			if result != tc.expectedResult {
				t.Errorf("Expected: %s, Got: %s", tc.expectedResult, result)
			}
		})
	}
}

func TestGetHTMLTemplate(t *testing.T) {
	testCases := []struct {
		name           string
		configTemplate string
		expectedResult string
	}{
		{
			name:           "with custom template",
			configTemplate: "<html>Custom template</html>",
			expectedResult: "<html>Custom template</html>",
		},
		{
			name:           "with empty template",
			configTemplate: "",
			expectedResult: config.DefaultHTMLTemplate(),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.OpenAIConfig{
				HTMLTemplate: tc.configTemplate,
			}
			
			service := NewOpenAIService(cfg)
			result := service.getHTMLTemplate()
			
			if result != tc.expectedResult {
				t.Errorf("Expected: %s, Got: %s", tc.expectedResult, result)
			}
		})
	}
}

func TestGetCSSStyles(t *testing.T) {
	testCases := []struct {
		name           string
		configStyles   string
		expectedResult string
	}{
		{
			name:           "with custom styles",
			configStyles:   "body { color: red; }",
			expectedResult: "body { color: red; }",
		},
		{
			name:           "with empty styles",
			configStyles:   "",
			expectedResult: config.DefaultCSSStyles(),
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.OpenAIConfig{
				CSSStyles: tc.configStyles,
			}
			
			service := NewOpenAIService(cfg)
			result := service.getCSSStyles()
			
			if result != tc.expectedResult {
				t.Errorf("Expected: %s, Got: %s", tc.expectedResult, result)
			}
		})
	}
}

func TestFallbackHTMLGeneration(t *testing.T) {
	cfg := &config.OpenAIConfig{}
	service := NewOpenAIService(cfg)
	
	htmlContent := "<h1>Test Content</h1>"
	username := "testuser"
	
	result := service.fallbackHTMLGeneration(htmlContent, username)
	
	// Verify result contains expected elements
	if !contains(result, "<!DOCTYPE html>") {
		t.Error("Result should contain DOCTYPE declaration")
	}
	
	if !contains(result, htmlContent) {
		t.Error("Result should contain the provided HTML content")
	}
	
	if !contains(result, username) {
		t.Error("Result should contain the username")
	}
	
	if !contains(result, "<title>") {
		t.Error("Result should contain a title tag")
	}
	
	if !contains(result, "<style>") {
		t.Error("Result should contain CSS styles")
	}
}

func TestConvertMarkdownToHTML(t *testing.T) {
	cfg := &config.OpenAIConfig{
		HTMLTemplate: `<!DOCTYPE html><html><head><title>{{.Username}}</title><style>{{.CSSStyles}}</style></head><body>{{.Content}}</body></html>`,
		CSSStyles:    "body { margin: 0; }",
	}
	
	service := NewOpenAIService(cfg)
	
	markdownContent := "# Test Header\n\nThis is **bold** text."
	username := "testuser"
	
	result := service.ConvertMarkdownToHTML(markdownContent, username)
	
	// Verify result contains expected elements
	if !contains(result, "<!DOCTYPE html>") {
		t.Error("Result should contain DOCTYPE declaration")
	}
	
	if !contains(result, username) {
		t.Error("Result should contain the username")
	}
	
	if !contains(result, "body { margin: 0; }") {
		t.Error("Result should contain the CSS styles")
	}
	
	// Should contain converted markdown
	if !contains(result, "<h1>") {
		t.Error("Markdown header should be converted to HTML")
	}
	
	if !contains(result, "<strong>") || !contains(result, "</strong>") {
		t.Error("Markdown bold should be converted to HTML")
	}
}

func TestConvertMarkdownToHTMLWithInvalidTemplate(t *testing.T) {
	cfg := &config.OpenAIConfig{
		HTMLTemplate: `{{invalid template syntax`,
		CSSStyles:    "body { margin: 0; }",
	}
	
	service := NewOpenAIService(cfg)
	
	markdownContent := "# Test Header"
	username := "testuser"
	
	result := service.ConvertMarkdownToHTML(markdownContent, username)
	
	// Should fallback to simple HTML generation
	if !contains(result, "<!DOCTYPE html>") {
		t.Error("Fallback should contain DOCTYPE declaration")
	}
	
	if !contains(result, username) {
		t.Error("Fallback should contain the username")
	}
}

func TestAnalyzeGitHubDataWithoutClient(t *testing.T) {
	cfg := &config.OpenAIConfig{
		APIKey: "", // Empty API key means no client
	}
	
	service := NewOpenAIService(cfg)
	
	auditResult := &dto.AuditResult{
		UserInfo: dto.UserInfo{Username: "testuser"},
	}
	
	_, err := service.AnalyzeGitHubData(auditResult)
	if err == nil {
		t.Error("AnalyzeGitHubData() should fail when client is not initialized")
	}
	
	expectedErrMsg := "OpenAI client not initialized - API key required"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message: %s, Got: %s", expectedErrMsg, err.Error())
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
