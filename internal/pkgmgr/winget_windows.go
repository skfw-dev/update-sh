//go:build windows
// +build windows

package pkgmgr

import (
	"update-sh/internal/runner"

	"github.com/rs/zerolog/log"
)

// WinGetManager implements PackageManagerImpl for WinGet on Windows.
type WinGetManager struct{}

// Update performs package updates using WinGet.
func (w *WinGetManager) Update(dryRun bool) error {
	log.Info().Msg("--- WinGet Package Management (Windows) ---")
	if !runner.CommandExists("winget") {
		log.Debug().Msg("WinGet not found. Skipping WinGet package management.")
		return nil
	}

	// WinGet upgrade flags:
	// --all: Upgrades all installed packages
	// --include-unknown: Includes packages that are not recognized by WinGet
	// --silent: Runs in silent mode
	// --accept-package-agreements: Accepts package agreements
	// --accept-source-agreements: Accepts source agreements
	wingetArgs := []string{"upgrade", "--all", "--include-unknown", "--silent", "--accept-package-agreements", "--accept-source-agreements"}
	if err := runner.RunCommand("Update WinGet packages", dryRun, "winget", nil, wingetArgs...); err != nil {
		return err
	}
	log.Info().Msg("WinGet maintenance complete.")
	return nil
}
