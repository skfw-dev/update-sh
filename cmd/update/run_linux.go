//go:build linux
// +build linux

package update

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
	"update-sh/internal/distro"
	"update-sh/internal/health"
	"update-sh/internal/pkgmgr"
	"update-sh/internal/runner"
	"update-sh/internal/shxmgr"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func isAdmin() (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, fmt.Errorf("failed to get current user: %w", err)
	}

	// Check if UID == 0 (root)
	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return false, fmt.Errorf("invalid UID: %w", err)
	}
	return uid == 0, nil
}

func acquireRoot() {
	admin, err := isAdmin()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to determine administrative status.")
	}

	if !admin {
		log.Info().Msg("Script is not running as root. Attempting to re-run with sudo...")

		selfPath, err := os.Executable()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get executable path.")
		}

		args := os.Args // Preserve original arguments

		// Attempt to re-execute the script with sudo
		if runner.CommandExists("sudo") { // Check if sudo command exists
			log.Info().Msg("Using sudo to elevate privileges.")
			sudoArgs := []string{"sudo", selfPath}
			sudoArgs = append(sudoArgs, args[1:]...) // Add original arguments
			err = syscall.Exec("/usr/bin/sudo", sudoArgs, os.Environ())
			if err != nil {
				log.Warn().Err(err).Msg("Failed to re-run with sudo.")
			}
		} else {
			log.Info().Msg("sudo command not found.")
		}

		// If sudo failed or wasn't found, attempt with doas
		if err != nil || !runner.CommandExists("sudo") { // Only try doas if sudo failed or wasn't there
			if runner.CommandExists("doas") { // Check if doas command exists
				log.Info().Msg("Using doas to elevate privileges.")
				doasArgs := []string{"doas", selfPath}
				doasArgs = append(doasArgs, args[1:]...) // Add original arguments
				err = syscall.Exec("/usr/bin/doas", doasArgs, os.Environ())
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to re-run with doas. Neither 'sudo' nor 'doas' found or failed to execute. Please run this script as root.")
				}
			} else {
				log.Fatal().Msg("doas command not found. Neither 'sudo' nor 'doas' found. Please run this script as root.")
			}
		}

		// If we reach here, it means either sudo/doas was found and exec'd (which replaces the process),
		// or they weren't found/failed, and the fatal error would have been logged.
		// This line should ideally not be reached if exec is successful.
		log.Fatal().Msg("Failed to acquire root privileges. Please ensure sudo or doas is installed and configured, or run as root.")
	} else {
		log.Info().Msg("Script is already running as root.")
	}
}

// performLinuxPackageUpdates runs all Linux-specific package manager updates.
func performLinuxPackageUpdates(dryRun bool, primaryPkgManager string) {
	var managersToRun []pkgmgr.PackageManagerImpl

	// Prioritize based on detected primary package manager.
	switch primaryPkgManager {
	case "apt":
		managersToRun = append(managersToRun, &pkgmgr.APTManager{})
	case "dnf":
		managersToRun = append(managersToRun, &pkgmgr.DNFManager{})
	case "pacman":
		managersToRun = append(managersToRun, &pkgmgr.PacmanManager{})
	case "zypper":
		managersToRun = append(managersToRun, &pkgmgr.ZypperManager{})
	case "pkg", "pkg_add", "generic_bsd_pkg": // Handle BSD package managers for Linux builds (e.g., WSL)
		managersToRun = append(managersToRun, &pkgmgr.BSDManager{})
	default:
		// If the primary package manager isn't definitively detected,
		// attempt to run common Linux package managers. Each manager will
		// internally check if its corresponding command exists.
		log.Info().Msg("Primary Linux package manager not definitively detected. Attempting common Linux package managers.")
		managersToRun = append(managersToRun,
			&pkgmgr.APTManager{},
			&pkgmgr.DNFManager{},
			&pkgmgr.PacmanManager{},
			&pkgmgr.ZypperManager{},
			&pkgmgr.BSDManager{},
		)
	}

	// Snap and Flatpak are universal Linux package managers (cross-distro),
	// so always attempt to run their updates if their commands exist.
	// Their implementations (e.g., `pkgmgr/snap_linux.go`) already have the `_linux.go` tag.
	managersToRun = append(managersToRun, &pkgmgr.SnapManager{}, &pkgmgr.FlatpakManager{})

	// Execute all collected package managers.
	for _, pm := range managersToRun {
		if err := pm.Update(dryRun); err != nil {
			// Log an error if a specific package manager update fails.
			log.Error().Err(err).Msgf("Linux package manager update failed for %T.", pm)
		}
	}
}

func performMaintenance(dryRun, initCheckOnly, zshUpdateEnabled, pwshUpdateEnabled bool) {
	log.Info().Msg("Starting comprehensive system maintenance script.")
	log.Info().Msgf("Log file: %s", viper.GetString("log_file"))

	// Acquire root privileges based on the OS. This function is defined in run_linux.go or run_windows.go
	acquireRoot()

	// Set non-interactive mode for Debian-based systems (Linux-specific)
	if runtime.GOOS == "linux" {
		os.Setenv("DEBIAN_FRONTEND", "noninteractive")
	}

	// Detect distribution and primary package manager first
	d, err := distro.DetectDistro()
	if err != nil {
		log.Error().Err(err).Msg("Error detecting distribution.")
		d = &distro.Distribution{
			ID:                    "unknown",
			IDLike:                "unknown",
			PrimaryPackageManager: "unknown",
		}
	}
	log.Info().Msgf("Detected OS: %s, Distribution ID: %s, Family: %s, Suggested Primary Package Manager: %s", runtime.GOOS, d.ID, d.Family, d.PrimaryPackageManager)

	// --- System Health Checks ---
	healthManager := &health.LinuxHealthManager{}
	if err := healthManager.CheckHealth(dryRun); err != nil {
		log.Error().Err(err).Msgf("System health check failed for %T.", healthManager)
	}

	// --- Shell-specific Updates ---
	var shlexManagersToRun []shxmgr.ShlexManagerImpl

	if zshUpdateEnabled {
		shlexManagersToRun = append(shlexManagersToRun, &shxmgr.ZshManager{})
	} else {
		log.Info().Msg("Skipping Zsh update. Use '-z' to enable.")
	}

	if pwshUpdateEnabled {
		// Create PwshManager and pass the detected primary package manager
		shlexManagersToRun = append(shlexManagersToRun, &shxmgr.PwshManager{PrimaryPackageManager: d.PrimaryPackageManager})
	} else {
		log.Info().Msg("Skipping PowerShell update. Use '-p' to enable.")
	}

	// Execute all collected shell managers
	for _, shlexManager := range shlexManagersToRun {
		if err := shlexManager.Update(dryRun); err != nil {
			log.Error().Err(err).Msgf("Shell component update failed for %T.", shlexManager)
		}
	}

	// --- Core Package Manager Updates (Platform-specific calls) ---
	if !initCheckOnly {
		log.Info().Msg("--- Starting Core Package Manager Updates ---")
		performLinuxPackageUpdates(dryRun, d.PrimaryPackageManager)
		log.Info().Msg("--- Core Package Manager Updates Complete ---")
	} else {
		log.Info().Msg("Skipping core package management updates due to '--init-check' flag.")
	}

	log.Info().Msg("Comprehensive system maintenance complete.")
	if dryRun {
		log.Info().Msg("Remember: This was a DRY RUN. No changes were applied.")
	}
}
