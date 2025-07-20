//go:build linux
// +build linux

package config

import (
	"sync"
)

// LinuxConfigManager implements ConfigImpl for Linux systems.
type LinuxConfigManager struct{}

// GetDefaultLogFile returns the default log file path for Linux.
func (l *LinuxConfigManager) GetDefaultLogFile() string {
	return "/var/log/system-maintenance.log"
}

// GetDefaultUserUID returns the default user UID for Linux.
func (l *LinuxConfigManager) GetDefaultUserID() string {
	return "1000" // Common default UID for the first non-root user on Linux
}

var configManagerOnce sync.Once
var currentConfigManager ConfigImpl

func GetConfigManager() ConfigImpl {
	return &LinuxConfigManager{}
}

func IsLinux() bool {
	return true // This function can be used to check if the current OS is Linux
}

func IsWindows() bool {
	return false // This function can be used to check if the current OS is Windows
}
