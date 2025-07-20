//go:build windows
// +build windows

package pkgmgr

import (
	"update-sh/internal/runner"

	"github.com/rs/zerolog/log"
)

// ChocolateyManager implements PackageManagerImpl for Chocolatey on Windows.
type ChocolateyManager struct{}

// Update performs package updates using Chocolatey.
func (c *ChocolateyManager) Update(dryRun bool) error {
	log.Info().Msg("--- Chocolatey Package Management (Windows) ---")
	if !runner.CommandExists("choco") {
		log.Debug().Msg("Chocolatey not found. Skipping Chocolatey package management.")
		return nil
	}

	// choco upgrade all -y: Upgrades all packages, accepts confirmation
	chocoArgs := []string{"upgrade", "all", "-y"}
	if err := runner.RunCommand("Update Chocolatey packages", dryRun, "choco", nil, chocoArgs...); err != nil {
		return err
	}

	// choco clean -y: Cleans up old package files
	chocoArgs = []string{"cache", "remove", "-y"}
	if err := runner.RunCommand("Clean Chocolatey cache", dryRun, "choco", nil, chocoArgs...); err != nil {
		log.Warn().Msg("Failed to clean Chocolatey cache or no cache to clean.")
	}
	log.Info().Msg("Chocolatey maintenance complete.")
	return nil
}
