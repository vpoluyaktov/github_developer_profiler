package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	
	if cfg.GitHub == nil {
		t.Error("GitHub config should not be nil")
	}
	
	if cfg.OpenAI == nil {
		t.Error("OpenAI config should not be nil")
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}
	
	if path == "" {
		t.Error("Config path should not be empty")
	}
	
	if !filepath.IsAbs(path) {
		t.Error("Config path should be absolute")
	}
	
	expectedSuffix := filepath.Join(".config", ConfigDirName, ConfigFileName)
	// Normalize paths for cross-platform comparison
	normalizedPath := filepath.ToSlash(path)
	normalizedSuffix := filepath.ToSlash(expectedSuffix)
	if !strings.HasSuffix(normalizedPath, normalizedSuffix) {
		t.Errorf("Config path should end with %s, got %s", expectedSuffix, path)
	}
}

func TestGetConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() failed: %v", err)
	}
	
	if dir == "" {
		t.Error("Config dir should not be empty")
	}
	
	if !filepath.IsAbs(dir) {
		t.Error("Config dir should be absolute")
	}
	
	expectedSuffix := filepath.Join(".config", ConfigDirName)
	// Normalize paths for cross-platform comparison
	normalizedDir := filepath.ToSlash(dir)
	normalizedSuffix := filepath.ToSlash(expectedSuffix)
	if !strings.HasSuffix(normalizedDir, normalizedSuffix) {
		t.Errorf("Config dir should end with %s, got %s", expectedSuffix, dir)
	}
}

func TestConfigExists(t *testing.T) {
	// Test ConfigExists function - should not panic in any environment
	exists := ConfigExists()
	t.Logf("Config exists: %v", exists)
	
	// In CI/CD environments, config likely won't exist, which is expected
	// This test validates the function works without errors
}

func TestLoadConfig(t *testing.T) {
	// Test loading config - should work in CI/CD without local files
	cfg, err := LoadConfig()
	
	// In CI/CD, config file won't exist, so we expect default config
	if err != nil {
		// Error is expected when no config file exists
		t.Logf("LoadConfig error (expected in CI/CD): %v", err)
		// Should still return a valid default config
		if cfg == nil {
			t.Error("LoadConfig should return default config even on error")
			return
		}
	}
	
	// Validate the returned config structure
	if cfg == nil {
		t.Fatal("LoadConfig returned nil config")
	}
	if cfg.GitHub == nil {
		t.Error("GitHub config should not be nil")
	}
	if cfg.OpenAI == nil {
		t.Error("OpenAI config should not be nil")
	}
}

func TestLoadConfigWithoutFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	
	// Temporarily override home directory for testing
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)
	
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	
	// Should return default config when file doesn't exist
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil")
	}
	
	if cfg.GitHub == nil {
		t.Error("GitHub config should not be nil")
	}
	
	if cfg.OpenAI == nil {
		t.Error("OpenAI config should not be nil")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	
	// Temporarily override home directory for testing
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)
	
	// Create test config
	originalConfig := &Config{
		GitHub: &GitHubConfig{
			Token:              "test-token",
			SampledRepoCount:   10,
			CommitsPerRepo:     20,
			SampleFileCount:    15,
			AnalysisYears:      2,
			IncludePrivateRepo: true,
			RandomSeed:         42,
			SaveDebugJSON:      true,
		},
		OpenAI: &OpenAIConfig{
			APIKey:       "test-api-key",
			Model:        "gpt-4",
			SystemPrompt: "test prompt",
			HTMLTemplate: "test template",
			CSSStyles:    "test styles",
		},
	}
	
	// Save config
	err := SaveConfig(originalConfig)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}
	
	// Verify config file was created
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
	
	// Load config back
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	
	// Verify loaded config matches original (excluding encryption)
	if loadedConfig.GitHub.Token != originalConfig.GitHub.Token {
		t.Errorf("GitHub token mismatch. Expected: %s, Got: %s", originalConfig.GitHub.Token, loadedConfig.GitHub.Token)
	}
	
	if loadedConfig.OpenAI.APIKey != originalConfig.OpenAI.APIKey {
		t.Errorf("OpenAI API key mismatch. Expected: %s, Got: %s", originalConfig.OpenAI.APIKey, loadedConfig.OpenAI.APIKey)
	}
	
	if loadedConfig.GitHub.SampledRepoCount != originalConfig.GitHub.SampledRepoCount {
		t.Errorf("SampledRepoCount mismatch. Expected: %d, Got: %d", originalConfig.GitHub.SampledRepoCount, loadedConfig.GitHub.SampledRepoCount)
	}
	
	if loadedConfig.OpenAI.Model != originalConfig.OpenAI.Model {
		t.Errorf("Model mismatch. Expected: %s, Got: %s", originalConfig.OpenAI.Model, loadedConfig.OpenAI.Model)
	}
}

func TestSaveConfigWithEmptyTokens(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	
	// Temporarily override home directory for testing
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)
	
	// Create config with empty tokens
	cfg := &Config{
		GitHub: &GitHubConfig{
			Token:            "",
			SampledRepoCount: 5,
		},
		OpenAI: &OpenAIConfig{
			APIKey: "",
			Model:  "gpt-3.5-turbo",
		},
	}
	
	// Save config
	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() with empty tokens failed: %v", err)
	}
	
	// Load config back
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	
	// Verify empty tokens remain empty
	if loadedConfig.GitHub.Token != "" {
		t.Error("Empty GitHub token should remain empty")
	}
	
	if loadedConfig.OpenAI.APIKey != "" {
		t.Error("Empty OpenAI API key should remain empty")
	}
}

func TestLoadConfigWithInvalidJSON(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	
	// Temporarily override home directory for testing
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)
	
	// Create config directory
	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() failed: %v", err)
	}
	
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	
	// Write invalid JSON to config file
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}
	
	invalidJSON := `{"invalid": json}`
	err = os.WriteFile(configPath, []byte(invalidJSON), 0600)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}
	
	// Try to load config
	_, err = LoadConfig()
	if err == nil {
		t.Error("LoadConfig() should fail with invalid JSON")
	}
}

func TestSaveConfigCreatesDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	
	// Temporarily override home directory for testing (cross-platform)
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()
	
	// Set both HOME and USERPROFILE for cross-platform compatibility
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	
	// Get config directory path
	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() failed: %v", err)
	}
	
	// Remove config directory if it exists (for clean test)
	os.RemoveAll(configDir)
	
	// Verify config directory doesn't exist now
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		t.Logf("Config directory exists, removing for clean test: %s", configDir)
		os.RemoveAll(configDir)
	}
	
	// Save config
	cfg := DefaultConfig()
	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}
	
	// Verify config directory was created
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Config directory should be created by SaveConfig()")
	}
}

func TestConfigJSONStructure(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	
	// Temporarily override home directory for testing
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)
	
	// Save default config
	cfg := DefaultConfig()
	cfg.GitHub.Token = "test-token"
	cfg.OpenAI.APIKey = "test-key"
	
	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}
	
	// Read the raw JSON file
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	
	// Parse JSON to verify structure
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		t.Fatalf("Config file contains invalid JSON: %v", err)
	}
	
	// Verify top-level structure
	if _, exists := jsonData["github"]; !exists {
		t.Error("Config JSON should contain 'github' section")
	}
	
	if _, exists := jsonData["openai"]; !exists {
		t.Error("Config JSON should contain 'openai' section")
	}
	
	// Verify tokens are encrypted (should be base64 strings, not plain text)
	githubSection := jsonData["github"].(map[string]interface{})
	if token, exists := githubSection["token"]; exists && token != "" {
		tokenStr := token.(string)
		if tokenStr == "test-token" {
			t.Error("GitHub token should be encrypted in saved file")
		}
	}
}
