package config

import (
	"strings"
	"testing"
)

func TestDefaultSystemPrompt(t *testing.T) {
	prompt := DefaultSystemPrompt()
	
	if prompt == "" {
		t.Error("DefaultSystemPrompt() should not return empty string")
	}
	
	// Verify it contains key elements expected in a system prompt
	expectedElements := []string{
		"GitHub",
		"developer",
		"analysis",
		"assessment",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(strings.ToLower(prompt), strings.ToLower(element)) {
			t.Errorf("System prompt should contain '%s'", element)
		}
	}
}

func TestDefaultHTMLTemplate(t *testing.T) {
	template := DefaultHTMLTemplate()
	
	if template == "" {
		t.Error("DefaultHTMLTemplate() should not return empty string")
	}
	
	// Verify it contains essential HTML structure
	expectedElements := []string{
		"<!DOCTYPE html>",
		"<html",
		"<head>",
		"<body>",
		"{{.Username}}",
		"{{.Content}}",
		"{{.CSSStyles}}",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(template, element) {
			t.Errorf("HTML template should contain '%s'", element)
		}
	}
}

func TestDefaultCSSStyles(t *testing.T) {
	styles := DefaultCSSStyles()
	
	if styles == "" {
		t.Error("DefaultCSSStyles() should not return empty string")
	}
	
	// Verify it contains basic CSS elements
	expectedElements := []string{
		"body",
		"font-family",
		"margin",
		"h1",
		"h2",
		"table",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(styles, element) {
			t.Errorf("CSS styles should contain '%s'", element)
		}
	}
}

func TestTemplateConsistency(t *testing.T) {
	// Test that multiple calls return the same content
	prompt1 := DefaultSystemPrompt()
	prompt2 := DefaultSystemPrompt()
	
	if prompt1 != prompt2 {
		t.Error("DefaultSystemPrompt() should return consistent results")
	}
	
	template1 := DefaultHTMLTemplate()
	template2 := DefaultHTMLTemplate()
	
	if template1 != template2 {
		t.Error("DefaultHTMLTemplate() should return consistent results")
	}
	
	styles1 := DefaultCSSStyles()
	styles2 := DefaultCSSStyles()
	
	if styles1 != styles2 {
		t.Error("DefaultCSSStyles() should return consistent results")
	}
}
