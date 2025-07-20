//go:build linux
// +build linux

package pkgmgr

import (
	"os/exec"
	"strings"

	"update-sh/internal/runner" // Import runner for command execution

	"github.com/rs/zerolog/log" // Import zerolog for logging
)

// PacmanManager implements PackageManagerImpl for Pacman.
type PacmanManager struct{}

// Update performs Pacman package management operations on Linux.
func (p *PacmanManager) Update(dryRun bool) error {
	log.Info().Msg("--- Pacman Package Management ---")
	if !runner.CommandExists("pacman") {
		log.Debug().Msg("Pacman not found. Skipping Pacman package management.")
		return nil // No error if Pacman is not present
	}

	// Update Pacman packages: 'pacman -Syu --noconfirm'
	// -S: Sync packages
	// -y: Refresh package databases
	// -u: Upgrade installed packages
	// --noconfirm: Skip confirmation prompts
	pacmanArgs := []string{"-Syu", "--noconfirm"}
	if err := runner.RunCommand("Update Pacman packages", dryRun, "pacman", nil, pacmanArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to update Pacman packages.")
		return err
	}

	// Remove orphaned Pacman packages
	// Orphaned packages are those that were installed as dependencies but are no longer needed by any explicitly installed package.
	if dryRun {
		log.Info().Msg("Dry Run: Would remove orphaned Pacman packages.")
	} else {
		log.Info().Msg("Removing orphaned Pacman packages...")
		// First, list orphaned packages: 'pacman -Qtdq'
		// -Q: Query the package database
		// -t: Limit to packages that are no longer required by any installed package
		// -d: Limit to dependencies
		// -q: Only show package names
		cmd := exec.Command("pacman", "-Qtdq")
		output, err := cmd.Output()
		if err == nil && len(strings.TrimSpace(string(output))) > 0 {
			// If there are orphaned packages, remove them: 'pacman -Rns --noconfirm'
			// -R: Remove packages
			// -n: Do not save configuration files
			// -s: Remove dependencies that are no longer required by any installed package
			// --noconfirm: Skip confirmation prompts
			orphanedPackages := strings.Fields(strings.TrimSpace(string(output)))
			pacmanArgs = append([]string{"-Rns", "--noconfirm"}, orphanedPackages...)
			if err := runner.RunCommand("Remove orphaned Pacman packages", dryRun, "pacman", nil, pacmanArgs...); err != nil {
				log.Error().Err(err).Msg("Failed to remove orphaned Pacman packages.")
			} else {
				log.Debug().Msg("Pacman orphaned packages removed.")
			}
		} else if err != nil {
			// Log error if pacman -Qtdq itself failed, but not if there are simply no orphaned packages
			log.Warn().Err(err).Msg("Failed to query orphaned Pacman packages (might be nothing to remove).")
		} else {
			log.Info().Msg("No Pacman orphaned packages to remove.")
		}
	}

	// Clean Pacman cache: 'pacman -Sc --noconfirm'
	// -S: Sync packages (used with -c for cleanup context)
	// -c: Clean the package cache. Using -c twice means remove all downloaded packages not currently installed.
	// --noconfirm: Skip confirmation prompts
	pacmanArgs = []string{"-Sc", "--noconfirm"}
	if err := runner.RunCommand("Clean Pacman cache", dryRun, "pacman", nil, pacmanArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to clean Pacman cache.")
		return err
	}

	log.Info().Msg("Pacman maintenance complete.")
	return nil
}
