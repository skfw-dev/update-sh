//go:build linux
// +build linux

package pkgmgr

import (
	"update-sh/internal/runner" // Import runner for command execution

	"github.com/rs/zerolog/log" // Import zerolog for logging
)

// DNFManager implements PackageManagerImpl for DNF.
type DNFManager struct{}

// Update performs DNF package management operations on Linux.
func (d *DNFManager) Update(dryRun bool) error {
	log.Info().Msg("--- DNF Package Management ---")
	if !runner.CommandExists("dnf") {
		log.Debug().Msg("DNF not found. Skipping DNF package management.")
		return nil // No error if DNF is not present
	}

	// Update DNF packages: 'dnf -y upgrade --refresh'
	// The '--refresh' option ensures that the metadata cache is updated before the upgrade.
	if err := runner.RunCommand("Update DNF packages", dryRun, "dnf", nil, "upgrade", "-y", "--refresh"); err != nil {
		log.Error().Err(err).Msg("Failed to update DNF packages.")
		return err
	}

	// Remove unnecessary DNF packages: 'dnf autoremove -y'
	// This command removes packages that were installed as dependencies but are no longer required.
	if err := runner.RunCommand("Remove unnecessary DNF packages (autoremove equivalent)", dryRun, "dnf", nil, "autoremove", "-y"); err != nil {
		// DNF autoremove might return an error if there are no packages to remove.
		// We'll log it as a warning/info rather than a critical error.
		log.Info().Err(err).Msg("No DNF packages to autoremove or failed during autoremove (check logs for details).")
	}

	// Clean DNF cache: 'dnf clean all'
	// This clears all cached packages, headers, and metadata.
	if err := runner.RunCommand("Clean DNF cache", dryRun, "dnf", nil, "clean", "all"); err != nil {
		log.Error().Err(err).Msg("Failed to clean DNF cache.")
		return err
	}

	log.Info().Msg("DNF maintenance complete.")
	return nil
}
