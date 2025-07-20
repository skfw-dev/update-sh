package shxmgr

import (
	// Needed for fmt.Errorf
	// Needed for strings.TrimSpace

	"update-sh/internal/runner"

	"github.com/rs/zerolog/log"
	// Assuming version package is imported for GetPowerShellExecutable
)

// PwshManager implements ShlexManagerImpl for PowerShell.
type PwshManager struct {
	PrimaryPackageManager string // Need to pass this info to the update method
}

// Update performs PowerShell (pwsh) updates via detected package managers.
func (p *PwshManager) Update(dryRun bool) error {
	log.Info().Msg("--- PowerShell (pwsh) Update ---")

	// First, determine the PowerShell executable to use.
	// This is needed for the 'scoop' case, but also generally good for logging.
	psExe, psVersion, err := GetPowerShellExecutable()
	if err != nil {
		log.Error().Err(err).Msg("Failed to find a suitable PowerShell executable. Skipping PowerShell update via package manager.")
		log.Info().Msg("To install PowerShell 7, visit: https://aka.ms/powershell-release?tag=stable")
		return nil // Not a critical error for the whole script if we can't update pwsh itself
	}
	log.Info().Msgf("Detected PowerShell executable: %s (version %s).", psExe, psVersion.String())

	// Check if pwsh.exe (PowerShell 7+) is explicitly installed.
	// If not, provide guidance.
	if !runner.CommandExists("pwsh") {
		log.Info().Msg("PowerShell 7 (pwsh.exe) is not found. Attempting to update Windows PowerShell (powershell.exe) if applicable.")
		log.Info().Msg("For the best experience, consider installing PowerShell 7 from: https://aka.ms/powershell-release?tag=stable")
	}

	log.Info().Msg("Attempting to update PowerShell via system package manager...")

	// Attempt to update via package manager if supported
	// We use the PrimaryPackageManager field from the struct, which needs to be set when creating PwshManager
	switch p.PrimaryPackageManager {
	case "winget":
		// Winget command structure: winget upgrade <package_id>
		wingetArgs := []string{"upgrade", "Microsoft.PowerShell", "--silent", "--accept-package-agreements", "--accept-source-agreements"}
		if err := runner.RunCommand("Update PowerShell (Winget)", dryRun, "winget", nil, wingetArgs...); err != nil {
			log.Error().Err(err).Msg("Failed to update PowerShell via Winget.")
			return err
		}
	case "chocolatey":
		// Chocolatey command structure: choco upgrade powershell-core -y
		chocoArgs := []string{"upgrade", "powershell-core", "-y"}
		if err := runner.RunCommand("Update PowerShell (Chocolatey)", dryRun, "choco", nil, chocoArgs...); err != nil {
			log.Error().Err(err).Msg("Failed to update PowerShell via Chocolatey.")
			return err
		}
	case "scoop":
		// Corrected: Scoop commands must be run via PowerShell.
		// Command: pwsh.exe -NoProfile -Command "scoop update pwsh"
		log.Info().Msg("Attempting to update PowerShell via Scoop (using PowerShell executable).")
		scoopArgs := []string{"-NoProfile", "-Command", "scoop update pwsh"}
		if err := runner.RunCommand("Update PowerShell (Scoop)", dryRun, psExe, nil, scoopArgs...); err != nil { // Use powerShellExe here
			log.Error().Err(err).Msg("Failed to update PowerShell via Scoop.")
			return err
		}
	default:
		log.Info().Msg("No primary or configured common package manager found to update PowerShell automatically.")
		log.Info().Msg("Consider downloading the latest package from: https://aka.ms/powershell-release?tag=stable")
	}

	log.Info().Msg("PowerShell update complete.")
	return nil
}
