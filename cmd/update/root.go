package update

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"update-sh/internal/config" // Import config package
	"update-sh/internal/logger" // Alias to avoid conflict with zerolog's log
)

var (
	cfgFile string
	verbose bool
	quiet   bool
	dryRun  bool
)

// Declare a global instance of the config manager
var appConfig config.ConfigImpl

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "update-sh",
	Short: "A comprehensive system maintenance tool.",
	Long: `update-sh performs comprehensive system maintenance across various operating systems.
It handles package updates for common package managers and checks for system health.

Usage:
  [sudo|doas] update-sh [OPTIONS]

Example:
  sudo update-sh -v --dry-run
  sudo update-sh --zsh-update --pwsh-update
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Set viper defaults
		config.SetViperDefaults()

		// Initialize the platform-specific config manager
		appConfig = config.GetConfigManager()

		// Initialize logging based on flags and the default log file from config
		logger.Init(verbose, quiet, viper.GetString("log_file")) // Use viper.GetString for log_file

		// Check for mutually exclusive flags
		if verbose && quiet {
			log.Fatal().Msg("Options -v/--verbose and -q/--quiet are mutually exclusive. Please choose only one.")
		}

		// Bind flags to Viper after parsing, to ensure --config is processed first
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			log.Error().Err(err).Msg("Failed to bind flags to Viper.")
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is given, run the default maintenance (same as `run.go` logic)
		performMaintenance(dryRun, viper.GetBool("init-check"), viper.GetBool("zsh-update"), viper.GetBool("pwsh-update"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) // Cobra handles its own errors by printing
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global Flags (Persistent means available to all subcommands)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.update-sh.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output.")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all output except errors and final summary.")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Perform a dry run without making any changes to the system.")

	// Local Flags (Only for the root command if no subcommand is specified)
	rootCmd.Flags().BoolP("init-check", "i", false, "Only perform systemd/init checks, no package management.")
	rootCmd.Flags().BoolP("zsh-update", "z", false, "Update Oh My Zsh and Powerlevel10k.")
	rootCmd.Flags().BoolP("pwsh-update", "p", false, "Update PowerShell (pwsh).")

	// Initialize appConfig here to get default log file for viper.SetDefault
	// This is safe because GetConfigManager is idempotent (uses sync.Once)
	appConfig = config.GetConfigManager()

	// Bind flags to Viper (before PersistentPreRunE, so it can load config file based on --config)
	// These are default values if not provided via flags or config file
	viper.SetDefault("init-check", false)
	viper.SetDefault("zsh-update", false)
	viper.SetDefault("pwsh-update", false)
	viper.SetDefault("log_file", appConfig.GetDefaultLogFile()) // Use value from the config manager
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".update-sh" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".update-sh")
		viper.SetConfigType("yaml") // or json, toml etc.
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		// Only log if config file is explicitly specified via flag but not found
		if cfgFile != "" {
			log.Warn().Msgf("Could not read config file %s: %v", cfgFile, err)
		} else {
			log.Debug().Msg("No config file found, using defaults and environment variables.")
		}
	}
}
