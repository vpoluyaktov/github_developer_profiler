package services

import (
	"testing"

	"dev_profiler/internal/config"
)

func TestNewGitHubService(t *testing.T) {
	// Test structure initialization
	cfg := &config.GitHubConfig{
		Token:            "test-token",
		SampledRepoCount: 5,
	}
	
	service := NewGitHubService(cfg)
	
	if service == nil {
		t.Fatal("NewGitHubService() returned nil")
	}
	
	if service.config != cfg {
		t.Error("Service config should match provided config")
	}
	
	if service.client == nil {
		t.Error("GitHub client should be initialized when token is provided")
	}
}

func TestNewGitHubServiceWithoutToken(t *testing.T) {
	cfg := &config.GitHubConfig{
		Token:            "",
		SampledRepoCount: 5,
	}
	
	service := NewGitHubService(cfg)
	
	if service == nil {
		t.Fatal("NewGitHubService() returned nil")
	}
	
	if service.client == nil {
		t.Error("GitHub client should still be initialized even without token")
	}
}

func TestIsCodeFile(t *testing.T) {
	cfg := &config.GitHubConfig{}
	service := NewGitHubService(cfg)
	
	testCases := []struct {
		path     string
		expected bool
	}{
		// Code files (based on actual implementation)
		{"main.go", true},
		{"script.py", true},
		{"app.js", true},
		{"app.ts", true},
		{"service.java", true},
		{"utils.cpp", true},
		{"main.c", true},
		{"app.cs", true},
		{"model.rb", true},
		{"controller.php", true},
		{"main.rs", true},
		{"app.kt", true},
		{"app.swift", true},
		{"app.scala", true},
		{"script.r", true},
		{"main.m", true},
		{"header.h", true},
		{"header.hpp", true},
		{"utils.cc", true},
		{"main.cxx", true},
		
		// Non-code files
		{"README.md", false},
		{"LICENSE", false},
		{"image.png", false},
		{"document.pdf", false},
		{"data.json", false},
		{"config.ini", false},
		{"log.txt", false},
		{".gitignore", false},
		{".env", false},
		{"CHANGELOG.md", false},
		{"style.css", false},
		{"index.html", false},
		{"component.tsx", false},
		{"config.yml", false},
		{"docker-compose.yaml", false},
		{"Dockerfile", false},
		{"Makefile", false},
		{"package.json", false},
		{"requirements.txt", false},
		{"pom.xml", false},
		{"build.gradle", false},
		
		// Edge cases
		{"", false},
		{"file_without_extension", false},
		{".hidden_file", false},
		{"path/to/main.go", true},
		{"deep/nested/path/script.py", true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := service.isCodeFile(tc.path)
			if result != tc.expected {
				t.Errorf("isCodeFile(%q) = %v, expected %v", tc.path, result, tc.expected)
			}
		})
	}
}

func TestDetectLanguage(t *testing.T) {
	cfg := &config.GitHubConfig{}
	service := NewGitHubService(cfg)
	
	testCases := []struct {
		path     string
		expected string
	}{
		// Supported languages (based on actual implementation)
		{"main.go", "Go"},
		{"script.py", "Python"},
		{"app.js", "JavaScript"},
		{"app.ts", "TypeScript"},
		{"service.java", "Java"},
		{"utils.cpp", "C++"},
		{"header.h", "C/C++ Header"},
		{"header.hpp", "C++ Header"},
		{"main.c", "C"},
		{"app.cs", "C#"},
		{"model.rb", "Ruby"},
		{"controller.php", "PHP"},
		{"main.rs", "Rust"},
		{"app.kt", "Kotlin"},
		{"app.swift", "Swift"},
		{"app.scala", "Scala"},
		{"script.r", "R"},
		{"main.m", "Objective-C"},
		
		// Unsupported extensions (should return "Unknown")
		{"empty_content", ""},
		{"component.tsx", "Unknown"},
		{"style.css", "Unknown"},
		{"index.html", "Unknown"},
		{"script.sh", "Unknown"},
		{"config.yml", "Unknown"},
		{"docker-compose.yaml", "Unknown"},
		{"Dockerfile", "Unknown"},
		{"Makefile", "Unknown"},
		{"package.json", "Unknown"},
		{"data.xml", "Unknown"},
		{"query.sql", "Unknown"},
		{"style.scss", "Unknown"},
		{"style.less", "Unknown"},
		{"config.toml", "Unknown"},
		{"notebook.ipynb", "Unknown"},
		{"unknown.xyz", "Unknown"},
		{"file_without_extension", "Unknown"},
		{"", "Unknown"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := service.detectLanguage(tc.path)
			if result != tc.expected {
				t.Errorf("detectLanguage(%q) = %q, expected %q", tc.path, result, tc.expected)
			}
		})
	}
}

func TestHasTestIndicators(t *testing.T) {
	cfg := &config.GitHubConfig{}
	service := NewGitHubService(cfg)
	
	testCases := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "contains test import",
			content:  "import testing\n\ndef test_function():\n    pass",
			expected: true,
		},
		{
			name:     "contains unittest",
			content:  "import unittest\n\nclass TestCase(unittest.TestCase):\n    pass",
			expected: true,
		},
		{
			name:     "contains pytest",
			content:  "import pytest\n\ndef test_something():\n    assert True",
			expected: true,
		},
		{
			name:     "contains jest",
			content:  "const { test, expect } = require('@jest/globals');\n\ntest('example', () => {});",
			expected: true,
		},
		{
			name:     "contains mocha",
			content:  "const mocha = require('mocha');\n\ndescribe('test suite', () => {});",
			expected: true,
		},
		{
			name:     "contains go testing",
			content:  "package main\n\nimport \"testing\"\n\nfunc TestExample(t *testing.T) {}",
			expected: true,
		},
		{
			name:     "contains test function",
			content:  "function testSomething() {\n    // test code\n}",
			expected: true,
		},
		{
			name:     "contains describe block",
			content:  "describe('component', () => {\n    it('should work', () => {});\n});",
			expected: true,
		},
		{
			name:     "contains it block",
			content:  "it('should do something', () => {\n    expect(true).toBe(true);\n});",
			expected: true,
		},
		{
			name:     "contains expect assertion",
			content:  "function verify() {\n    expect(result).toEqual(expected);\n}",
			expected: true,
		},
		{
			name:     "contains assert",
			content:  "def check_result():\n    assert result == expected",
			expected: true,
		},
		{
			name:     "no test indicators",
			content:  "function regularFunction() {\n    return 'hello world';\n}",
			expected: false,
		},
		{
			name:     "empty content",
			content:  "",
			expected: false,
		},
		{
			name:     "only comments",
			content:  "// This is just a comment\n/* Another comment */",
			expected: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.hasTestIndicators(tc.content)
			if result != tc.expected {
				t.Errorf("hasTestIndicators() = %v, expected %v for content: %q", result, tc.expected, tc.content)
			}
		})
	}
}

func TestAssessCodeQuality(t *testing.T) {
	cfg := &config.GitHubConfig{}
	service := NewGitHubService(cfg)
	
	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "high quality code",
			content:  "// Well documented function\n// Another comment\n// More documentation\nfunc calculateSum(a, b int) int {\n    return a + b\n}\n\n// Another documented function\nfunc multiply(x, y int) int {\n    return x * y\n}",
			expected: "Good",
		},
		{
			name:     "medium quality code",
			content:  "func calculate(a, b int) int {\n    return a + b\n}\n\n// some processing\nfunc process() {\n    doSomething()\n}",
			expected: "Fair",
		},
		{
			name:     "low quality code",
			content:  "func f(x,y int)int{return x+y}\nfunc g(){}\nfunc h(){}",
			expected: "Needs Improvement",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "Unknown",
		},
		{
			name:     "only whitespace",
			content:  "   \n\n   \t  ",
			expected: "Needs Improvement",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.assessCodeQuality(tc.content)
			if result != tc.expected {
				t.Errorf("assessCodeQuality() = %q, expected %q", result, tc.expected)
			}
		})
	}
}

func TestAssessComplexity(t *testing.T) {
	cfg := &config.GitHubConfig{}
	service := NewGitHubService(cfg)
	
	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "high complexity",
			content:  "if (condition) {\n    for (let i = 0; i < 10; i++) {\n        if (i % 2 === 0) {\n            while (true) {\n                switch (value) {\n                    case 1: break;\n                }\n            }\n        }\n    }\n}",
			expected: "High",
		},
		{
			name:     "medium complexity",
			content:  "if condition\nfor i in range(10)\n    if check\n        process()\n",
			expected: "Medium",
		},
		{
			name:     "low complexity",
			content:  "function add(a, b) {\n    return a + b;\n}\n\nfunction greet(name) {\n    return 'Hello ' + name;\n}",
			expected: "Low",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "Unknown",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.assessComplexity(tc.content)
			if result != tc.expected {
				t.Errorf("assessComplexity() = %q, expected %q", result, tc.expected)
			}
		})
	}
}

func TestGetSanitizedConfig(t *testing.T) {
	cfg := &config.GitHubConfig{
		Token:            "secret-token",
		SampledRepoCount: 10,
		CommitsPerRepo:   20,
	}
	
	service := NewGitHubService(cfg)
	sanitized := service.getSanitizedConfig()
	
	// Token should be redacted
	if sanitized.Token != "[REDACTED]" {
		t.Errorf("Token should be redacted, got: %s", sanitized.Token)
	}
	
	// Other fields should remain unchanged
	if sanitized.SampledRepoCount != cfg.SampledRepoCount {
		t.Error("SampledRepoCount should not be modified")
	}
	
	if sanitized.CommitsPerRepo != cfg.CommitsPerRepo {
		t.Error("CommitsPerRepo should not be modified")
	}
	
	// Original config should remain unchanged
	if cfg.Token != "secret-token" {
		t.Error("Original config token should not be modified")
	}
}
