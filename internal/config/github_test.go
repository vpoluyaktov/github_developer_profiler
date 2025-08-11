package config

import (
	"testing"
)

func TestDefaultGitHubConfig(t *testing.T) {
	cfg := DefaultGitHubConfig()
	
	if cfg == nil {
		t.Fatal("DefaultGitHubConfig() returned nil")
	}
	
	// Verify default values are reasonable
	if cfg.SampledRepoCount <= 0 {
		t.Error("SampledRepoCount should be positive")
	}
	
	if cfg.CommitsPerRepo <= 0 {
		t.Error("CommitsPerRepo should be positive")
	}
	
	if cfg.SampleFileCount <= 0 {
		t.Error("SampleFileCount should be positive")
	}
	
	if cfg.AnalysisYears <= 0 {
		t.Error("AnalysisYears should be positive")
	}
	
	if cfg.RandomSeed < 0 {
		t.Error("RandomSeed should be non-negative")
	}
	
	// Token should be empty by default
	if cfg.Token != "" {
		t.Error("Default token should be empty")
	}
}

func TestGitHubConfigValidation(t *testing.T) {
	cfg := DefaultGitHubConfig()
	
	// Test that all fields have sensible defaults
	testCases := []struct {
		name     string
		field    string
		value    interface{}
		minValue interface{}
	}{
		{"SampledRepoCount", "SampledRepoCount", cfg.SampledRepoCount, 1},
		{"CommitsPerRepo", "CommitsPerRepo", cfg.CommitsPerRepo, 1},
		{"SampleFileCount", "SampleFileCount", cfg.SampleFileCount, 1},
		{"AnalysisYears", "AnalysisYears", cfg.AnalysisYears, 1},
		{"RandomSeed", "RandomSeed", cfg.RandomSeed, int64(0)},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.value.(type) {
			case int:
				if minVal, ok := tc.minValue.(int); ok && v < minVal {
					t.Errorf("%s should be at least %d, got %d", tc.field, minVal, v)
				}
			case int64:
				if minVal, ok := tc.minValue.(int64); ok && v < minVal {
					t.Errorf("%s should be at least %d, got %d", tc.field, minVal, v)
				}
			}
		})
	}
}

func TestGitHubConfigDefaults(t *testing.T) {
	cfg := DefaultGitHubConfig()
	
	// Verify specific expected defaults based on the application's purpose
	expectedDefaults := map[string]interface{}{
		"IncludePrivateRepo": false, // Should default to false for privacy
		"SaveDebugJSON":      false, // Should default to false to avoid clutter
	}
	
	if cfg.IncludePrivateRepo != expectedDefaults["IncludePrivateRepo"].(bool) {
		t.Errorf("IncludePrivateRepo default should be %v, got %v", 
			expectedDefaults["IncludePrivateRepo"], cfg.IncludePrivateRepo)
	}
	
	if cfg.SaveDebugJSON != expectedDefaults["SaveDebugJSON"].(bool) {
		t.Errorf("SaveDebugJSON default should be %v, got %v", 
			expectedDefaults["SaveDebugJSON"], cfg.SaveDebugJSON)
	}
}
