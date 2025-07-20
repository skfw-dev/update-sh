//go:build windows
// +build windows

package health

import (
	"update-sh/internal/runner"

	"github.com/rs/zerolog/log"
)

// WindowsHealthManager implements HealthImpl for Windows systems.
type WindowsHealthManager struct{}

// CheckHealth performs comprehensive Windows health checks.
func (w *WindowsHealthManager) CheckHealth(dryRun bool) error {
	log.Info().Msg("--- Starting Windows System Health Checks ---")

	// Example: Check DISM health
	w.checkDismHealth(dryRun)

	// Example: Check SFC (System File Checker)
	w.checkSfcIntegrity(dryRun)

	// Placeholder for other Windows-specific checks (e.g., Event Viewer logs, drive health)
	log.Info().Msg("Windows system health checks complete.")
	return nil
}

// checkDismHealth performs a DISM /RestoreHealth check on Windows.
func (w *WindowsHealthManager) checkDismHealth(dryRun bool) {
	log.Info().Msg("Checking Windows component store health with DISM...")
	if dryRun {
		log.Info().Msg("Dry Run: Would check DISM health.")
		return
	}

	if runner.CommandExists("dism") {
		dismArgs := []string{"/Online", "/Cleanup-Image", "/RestoreHealth"}
		if err := runner.RunCommand("Check DISM health", dryRun, "dism", nil, dismArgs...); err != nil {
			log.Error().Err(err).Msg("Failed to check/restore Windows component store health with DISM.")
		} else {
			log.Info().Msg("DISM health check complete.")
		}
	} else {
		log.Debug().Msg("DISM not found. Skipping DISM health check.")
	}
}

// checkSfcIntegrity performs an SFC /scannow check on Windows.
func (w *WindowsHealthManager) checkSfcIntegrity(dryRun bool) {
	log.Info().Msg("Checking system file integrity with SFC...")
	if dryRun {
		log.Info().Msg("Dry Run: Would check SFC integrity.")
		return
	}

	if runner.CommandExists("sfc") {
		sfcArgs := []string{"/scannow"}
		// if err := runner.RunCommand("Check SFC integrity", dryRun, "sfc", nil, sfcArgs...); err != nil {
		// 	log.Error().Err(err).Msg("Failed to check system file integrity with SFC.")
		// } else {
		// 	log.Info().Msg("SFC integrity check complete.")
		// }
		opts := runner.NewCommandOptions("Check SFC integrity", dryRun, "sfc", nil, sfcArgs...)
		opts.Encoding = runner.UTF16LE // Use UTF-16 Little Endian for Windows SFC output
		opts.User = "SYSTEM"           // SFC typically runs as SYSTEM user
		if err := runner.RunCommandWithOptions(opts); err != nil {
			log.Error().Err(err).Msg("Failed to check system file integrity with SFC.")
		} else {
			log.Info().Msg("SFC integrity check complete.")
		}
	} else {
		log.Debug().Msg("SFC not found. Skipping SFC integrity check.")
	}
}
