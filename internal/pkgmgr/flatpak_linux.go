//go:build linux
// +build linux

package pkgmgr

import (
	"update-sh/internal/runner" // Import runner for command execution

	"github.com/rs/zerolog/log" // Import zerolog for logging
)

// FlatpakManager implements PackageManagerImpl for Flatpak.
type FlatpakManager struct{}

// Update performs Flatpak package management operations on Linux.
func (f *FlatpakManager) Update(dryRun bool) error {
	log.Info().Msg("--- Flatpak Package Management ---")
	if !runner.CommandExists("flatpak") {
		log.Debug().Msg("Flatpak not found. Skipping Flatpak package management.")
		return nil // No error if Flatpak is not present
	}

	if dryRun {
		log.Info().Msg("Dry Run: Would update Flatpak packages.")
		return nil
	}

	// Flatpak applications can be installed system-wide or user-specific.
	// When run as root (sudo), `flatpak update` by default updates system-wide Flatpaks.
	// To update user-installed Flatpaks, it should be run as that specific user.
	// We'll prioritize running as the original invoking user if SUDO_USER is available.
	log.Info().Msg("Running Flatpak update as root (primarily for system-wide Flatpaks).")

	flatpakArgs := []string{"update", "-y"}
	if err := runner.RunCommand("Update Flatpak packages", dryRun, "flatpak", nil, flatpakArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to update Flatpak packages as root.")
		return err
	}
	log.Info().Msg("Flatpak packages updated as root.")

	// Flatpak cleanup (uninstalling unused runtimes and extensions)
	log.Info().Msg("Performing Flatpak cleanup...")
	flatpakArgs = []string{"uninstall", "--unused", "-y"}
	if err := runner.RunCommand("Clean Flatpak unused data", dryRun, "flatpak", nil, flatpakArgs...); err != nil {
		// Cleanup might not find anything to remove, which isn't an error.
		// Log as info/warn if it fails for other reasons.
		log.Warn().Err(err).Msg("Flatpak cleanup failed or found nothing to uninstall.")
	} else {
		log.Info().Msg("Flatpak cleanup complete.")
	}

	log.Info().Msg("Flatpak maintenance complete.")
	return nil
}
