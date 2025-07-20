//go:build linux
// +build linux

package pkgmgr

import (
	"update-sh/internal/runner" // Import runner for command execution

	"github.com/rs/zerolog/log" // Import zerolog for logging
)

// ZypperManager implements PackageManagerImpl for Zypper.
type ZypperManager struct{}

// Update performs Zypper package management operations on Linux.
func (z *ZypperManager) Update(dryRun bool) error {
	log.Info().Msg("--- Zypper Package Management ---")
	if !runner.CommandExists("zypper") {
		log.Debug().Msg("Zypper not found. Skipping Zypper package management.")
		return nil // No error if Zypper is not present
	}

	// Refresh Zypper repositories: 'zypper refresh'
	// This ensures that the local package metadata is up-to-date with the repositories.
	zypperArgs := []string{"refresh"}
	if err := runner.RunCommand("Refresh Zypper repositories", dryRun, "zypper", nil, zypperArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to refresh Zypper repositories.")
		return err
	}

	// Update Zypper packages: 'zypper update -y'
	// This upgrades all installed packages to their latest available versions.
	zypperArgs = []string{"update", "-y"}
	if err := runner.RunCommand("Update Zypper packages", dryRun, "zypper", nil, zypperArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to update Zypper packages.")
		return err
	}

	// Note on autoremove equivalent: Zypper does not have a direct 'autoremove all unneeded'
	// command like APT's `autoremove`. Unneeded dependencies are generally handled during
	// `zypper remove` or `zypper purge`.
	log.Info().Msg("Zypper does not have a direct 'autoremove all unneeded' equivalent like apt or dnf.")

	// Clean Zypper cache: 'zypper clean --all'
	// This clears all cached packages, metadata, and temporary files.
	zypperArgs = []string{"clean", "--all"}
	if err := runner.RunCommand("Clean Zypper cache", dryRun, "zypper", nil, zypperArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to clean Zypper cache.")
		return err
	}

	log.Info().Msg("Zypper maintenance complete.")
	return nil
}
