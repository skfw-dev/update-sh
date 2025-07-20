//go:build windows
// +build windows

package config

import (
	"os" // Needed for os.Getenv
	"sync"
)

// WindowsConfigManager implements ConfigImpl for Windows systems.
type WindowsConfigManager struct{}

// GetDefaultLogFile returns the default log file path for Windows.
func (w *WindowsConfigManager) GetDefaultLogFile() string {
	// A common location for application logs on Windows is in ProgramData or user's AppData.
	// For simplicity, let's use a path relative to the current user's APPDATA or TEMP.
	// In a real application, you might use ProgramData or a more robust path.
	if appData := os.Getenv("APPDATA"); appData != "" {
		return appData + "\\system-maintenance\\logs\\system-maintenance.log"
	}
	return os.TempDir() + "\\system-maintenance.log" // Fallback to Temp directory
}

// GetDefaultUserID returns an empty string for Windows as UID is not applicable.
func (w *WindowsConfigManager) GetDefaultUserID() string {
	return "" // UID concept is not directly applicable on Windows
}

var configManagerOnce sync.Once
var currentConfigManager ConfigImpl

func GetConfigManager() ConfigImpl {
	return &WindowsConfigManager{}
}

func IsLinux() bool {
	return false // This function can be used to check if the current OS is Linux
}

func IsWindows() bool {
	return true // This function can be used to check if the current OS is Windows
}
