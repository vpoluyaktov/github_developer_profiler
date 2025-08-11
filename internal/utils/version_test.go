package utils

import (
	"strings"
	"testing"
)

func TestAppName(t *testing.T) {
	if AppName == "" {
		t.Error("AppName should not be empty")
	}
}

func TestAppVersion(t *testing.T) {
	if AppVersion == "" {
		t.Error("AppVersion should not be empty")
	}
	
	// Version should follow semantic versioning pattern (basic check)
	if len(AppVersion) < 5 { // Minimum "x.y.z"
		t.Error("AppVersion should follow semantic versioning format")
	}
}

func TestGetDetailedVersion(t *testing.T) {
	version := GetDetailedVersion()
	
	if version == "" {
		t.Error("GetDetailedVersion() should not return empty string")
	}
	
	// Should contain app name and version
	if !strings.Contains(version, AppName) {
		t.Error("Detailed version should contain app name")
	}
	
	if !strings.Contains(version, AppVersion) {
		t.Error("Detailed version should contain app version")
	}
}
