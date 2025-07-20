package config

import "github.com/spf13/viper"

// ConfigImpl defines the common interface for retrieving configuration values.
type ConfigImpl interface {
	// GetDefaultLogFile returns the default path for the application's log file.
	GetDefaultLogFile() string
	// GetDefaultUserID returns the default user UID for user-specific operations (primarily Linux).
	// On Windows, this might return an empty string or a non-applicable value.
	GetDefaultUserID() string
	// Add other common configuration methods here as needed for cross-platform settings.
}

// SetViperDefaults sets default values in the Viper configuration store using the platform-specific
// configuration manager (ConfigImpl). This should be called after viper.AddConfigPath and before
// viper.ReadInConfig.
func SetViperDefaults() {
	// Get the platform-specific configuration manager instance
	cfgManager := GetConfigManager()

	// Set default values using methods from the interface
	viper.SetDefault("log_file", cfgManager.GetDefaultLogFile())
	viper.SetDefault("user_id", cfgManager.GetDefaultUserID()) // Default UID for user-specific actions on Linux
	// Add other default config values here, also potentially fetched from cfgManager if they are platform-specific.
}
