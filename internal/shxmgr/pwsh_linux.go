//go:build linux
// +build linux

package shxmgr

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"update-sh/internal/runner"
)

// PwshManager implements ShlexManagerImpl for PowerShell.
type PwshManager struct {
	PrimaryPackageManager string // Need to pass this info to the update method
}

// Update performs PowerShell (pwsh) updates via detected package managers.
func (p *PwshManager) Update(dryRun bool) error {
	log.Info().Msg("--- PowerShell (pwsh) Update ---")
	if !runner.CommandExists("pwsh") {
		log.Info().Msg("PowerShell (pwsh) is not installed. Skipping update.")
		log.Info().Msg("To install, visit: https://docs.microsoft.com/en-us/powershell/scripting/install/installing-powershell-on-linux")
		return nil
	}

	log.Info().Msg("PowerShell (pwsh) is already installed. Attempting to update via system package manager...")

	// Attempt to update via package manager if supported
	// We use the PrimaryPackageManager field from the struct, which needs to be set when creating PwshManager
	switch p.PrimaryPackageManager {
	case "apt":
		if err := runner.RunCommand("Update PowerShell (APT)", dryRun, "apt", nil, "install", "--only-upgrade", "powershell", "-y"); err != nil {
			log.Error().Err(err).Msg("Failed to update PowerShell via APT.")
			return err
		}
	case "dnf":
		if err := runner.RunCommand("Update PowerShell (DNF)", dryRun, "dnf", nil, "upgrade", "powershell", "-y"); err != nil {
			log.Error().Err(err).Msg("Failed to update PowerShell via DNF.")
			return err
		}
	case "pacman":
		if err := runner.RunCommand("Update PowerShell (Pacman)", dryRun, "pacman", nil, "-S", "powershell", "--noconfirm"); err != nil {
			log.Error().Err(err).Msg("Failed to update PowerShell via Pacman.")
			return err
		}
	case "zypper":
		if err := runner.RunCommand("Update PowerShell (Zypper)", dryRun, "zypper", nil, "update", "powershell", "-y"); err != nil {
			log.Error().Err(err).Msg("Failed to update PowerShell via Zypper.")
			return err
		}
	default:
		log.Info().Msg("No primary or configured common package manager found to update PowerShell automatically.")
		log.Info().Msg("Consider downloading the latest package from: https://github.com/PowerShell/PowerShell/releases")
	}

	// Update Oh My Posh CLI (can be cross-platform, but often installed via package managers or specific scripts)
	if runner.CommandExists("oh-my-posh") {
		log.Info().Msg("Found Oh My Posh CLI. Attempting to upgrade...")
		ompArgs := []string{"upgrade", "--force"}
		if err := runner.RunCommand("Upgrade Oh My Posh CLI", dryRun, "oh-my-posh", nil, ompArgs...); err != nil {
			log.Error().Err(err).Msg("Failed to upgrade Oh My Posh CLI.")
			return fmt.Errorf("failed to upgrade Oh My Posh CLI: %w", err)
		}
	} else {
		log.Debug().Msg("Oh My Posh CLI not found. Skipping Oh My Posh CLI update.")
	}

	log.Info().Msg("PowerShell update complete.")
	return nil
}
