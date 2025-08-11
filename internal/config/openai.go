package config

// OpenAIConfig holds OpenAI API configuration
type OpenAIConfig struct {
	APIKey       string `json:"api_key"`
	Model        string `json:"model"`
	SystemPrompt string `json:"system_prompt"`
	HTMLTemplate string `json:"html_template"`
	CSSStyles    string `json:"css_styles"`
}

// DefaultOpenAIConfig returns default OpenAI configuration
func DefaultOpenAIConfig() *OpenAIConfig {
	return &OpenAIConfig{
		APIKey:       "",
		Model:        "gpt-4o",
		SystemPrompt: DefaultSystemPrompt(),
		HTMLTemplate: DefaultHTMLTemplate(),
		CSSStyles:    DefaultCSSStyles(),
	}
}
