//go:build windows
// +build windows

package pkgmgr

import (
	"fmt"
	"update-sh/internal/runner"
	"update-sh/internal/shxmgr"

	"github.com/rs/zerolog/log"
	// Assuming this is needed for version checking in shxmgr
)

// ScoopManager implements PackageManagerImpl for Scoop on Windows.
type ScoopManager struct{}

// Update performs package updates using Scoop.
func (s *ScoopManager) Update(dryRun bool) error {
	log.Info().Msg("--- Scoop Package Management (Windows) ---")

	// Even if 'scoop' is not directly in PATH for cmd.exe, it might be available via PowerShell.
	// We'll proceed with PowerShell invocation.

	// First, ensure PowerShell executable is found and policy is set.
	psExe, psVersion, err := shxmgr.GetPowerShellExecutable()
	if err != nil {
		log.Error().Err(err).Msg("Failed to find a suitable PowerShell executable for Scoop. Please ensure PowerShell 7 (pwsh.exe) or Windows PowerShell is installed.")
		return fmt.Errorf("PowerShell not available for Scoop: %w", err)
	}
	log.Debug().Msgf("Using PowerShell executable '%s' (version %s) to run Scoop commands.", psExe, psVersion.String())

	// Ensure execution policy is set. This is critical for Scoop's PowerShell scripts.
	if err := shxmgr.SetExecutionPolicy(dryRun); err != nil {
		log.Error().Err(err).Msg("Failed to ensure PowerShell execution policy is set. Scoop operations might fail.")
		return err // Return error if policy check/set failed critically
	}

	user, err := runner.GetTargetUser() // Use the updated GetTargetUser
	if err != nil {
		log.Error().Err(err).Msg("Cannot update user-specific Scoop components.")
		return err // Return error for the interface
	}

	// Now, wrap Scoop commands in PowerShell.
	// Check if 'scoop' itself is callable within PowerShell.
	// This check is important as Scoop might not be installed or in the user's PowerShell profile.
	scoopArgs := []string{"-NoProfile", "-Command", "Get-Command scoop | Out-Null"}
	if err := runner.RunUserCommand("Check if Scoop is callable", dryRun, user, psExe, nil, scoopArgs...); err != nil {
		log.Warn().Msg("Scoop command not found when invoked via PowerShell. Skipping Scoop maintenance. Please ensure Scoop is correctly installed and its path is in your PowerShell profile.")
		return nil // Not a critical error if Scoop isn't installed
	}

	// Update Scoop itself: scoop update
	// Command: powershell.exe -NoProfile -Command "scoop update"
	log.Info().Msg("Updating Scoop core...")
	scoopArgs = []string{"-NoProfile", "-Command", "scoop update"}
	if err := runner.RunUserCommand("Update Scoop core", dryRun, user, psExe, nil, scoopArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to update Scoop core.")
		return err
	}

	// Update all installed Scoop packages: scoop update *
	// Command: powershell.exe -NoProfile -Command "scoop update *"
	log.Info().Msg("Updating all Scoop applications...")
	scoopArgs = []string{"-NoProfile", "-Command", "scoop update --all"}
	if err := runner.RunUserCommand("Update all Scoop applications", dryRun, user, psExe, nil, scoopArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to update all Scoop applications.")
		return err
	}

	// Scoop cleanup: scoop cleanup *
	// Command: powershell.exe -NoProfile -Command "scoop cleanup *"
	log.Info().Msg("Performing Scoop cleanup (removing old versions and shims)...")
	scoopArgs = []string{"-NoProfile", "-Command", "scoop cleanup --all"}
	if err := runner.RunUserCommand("Clean Scoop cache and old versions", dryRun, user, psExe, nil, scoopArgs...); err != nil {
		log.Warn().Err(err).Msg("Scoop cleanup failed or found nothing to clean.")
	} else {
		log.Info().Msg("Scoop cleanup complete.")
	}

	log.Info().Msg("Scoop maintenance complete.")
	return nil
}
