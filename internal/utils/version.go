package utils

import (
	"fmt"
	"runtime"
)

const (
	AppName    = "GitHub Developer Profiler"
	AppVersion = "1.0.0"
)

// GetDetailedVersion returns detailed version information
func GetDetailedVersion() string {
	return fmt.Sprintf("%s v%s (built with %s)", AppName, AppVersion, runtime.Version())
}
