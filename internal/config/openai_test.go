package config

import (
	"testing"
)

func TestDefaultOpenAIConfig(t *testing.T) {
	cfg := DefaultOpenAIConfig()
	
	if cfg == nil {
		t.Fatal("DefaultOpenAIConfig() returned nil")
	}
	
	// Verify default model is set
	if cfg.Model == "" {
		t.Error("Default model should not be empty")
	}
	
	// API key should be empty by default
	if cfg.APIKey != "" {
		t.Error("Default API key should be empty")
	}
	
	// System prompt should have a default value
	if cfg.SystemPrompt == "" {
		t.Error("Default system prompt should not be empty")
	}
	
	// HTML template should have a default value
	if cfg.HTMLTemplate == "" {
		t.Error("Default HTML template should not be empty")
	}
	
	// CSS styles should have a default value
	if cfg.CSSStyles == "" {
		t.Error("Default CSS styles should not be empty")
	}
}

func TestOpenAIConfigDefaults(t *testing.T) {
	cfg := DefaultOpenAIConfig()
	
	// Verify the model is a reasonable default
	expectedModels := []string{"gpt-4", "gpt-4-turbo", "gpt-3.5-turbo", "gpt-4o", "gpt-4.1"}
	modelValid := false
	for _, model := range expectedModels {
		if cfg.Model == model {
			modelValid = true
			break
		}
	}
	
	if !modelValid {
		t.Errorf("Default model should be one of %v, got %s", expectedModels, cfg.Model)
	}
}
