//go:build linux
// +build linux

package pkgmgr

import (
	"update-sh/internal/runner" // Import runner for command execution

	"github.com/rs/zerolog/log" // Import zerolog for logging
)

// SnapManager implements PackageManagerImpl for Snap.
type SnapManager struct{}

// Update performs Snap package management operations on Linux.
func (s *SnapManager) Update(dryRun bool) error {
	log.Info().Msg("--- Snap Package Management ---")
	if !runner.CommandExists("snap") {
		log.Debug().Msg("Snap not found. Skipping Snap package management.")
		return nil // No error if Snap is not present
	}

	// Update Snap packages: 'snap refresh'
	// The 'refresh' command updates a snap to the latest version.
	snapArgs := []string{"refresh"}
	if err := runner.RunCommand("Update Snap packages", dryRun, "snap", nil, snapArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to update Snap packages.")
		return err
	}

	log.Info().Msg("Snap maintenance complete.")
	return nil
}
