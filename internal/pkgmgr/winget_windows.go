//go:build windows
// +build windows

package pkgmgr

import (
	"update-sh/internal/runner"

	"github.com/rs/zerolog/log"
)

// WinGetManager implements PackageManagerImpl for Winget on Windows.
type WinGetManager struct{}

// Update performs package updates using Winget.
func (w *WinGetManager) Update(dryRun bool) error {
	log.Info().Msg("--- Winget Package Management (Windows) ---")
	if !runner.CommandExists("winget") {
		log.Debug().Msg("Winget not found. Skipping Winget package management.")
		return nil
	}

	// Winget upgrade flags:
	// --all: Upgrades all installed packages
	// --include-unknown: Includes packages that are not recognized by Winget
	// --silent: Runs in silent mode
	// --accept-package-agreements: Accepts package agreements
	// --accept-source-agreements: Accepts source agreements
	wingetArgs := []string{"upgrade", "--all", "--include-unknown", "--silent", "--accept-package-agreements", "--accept-source-agreements"}
	if err := runner.RunCommand("Update Winget packages", dryRun, "winget", nil, wingetArgs...); err != nil {
		return err
	}
	log.Info().Msg("Winget maintenance complete.")
	return nil
}
