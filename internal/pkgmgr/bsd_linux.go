//go:build linux
// +build linux

package pkgmgr

import (
	"update-sh/internal/runner" // Import runner for command execution

	"github.com/rs/zerolog/log" // Import zerolog for logging
)

// BSDManager implements PackageManagerImpl for BSD-like systems (like FreeBSD and OpenBSD)
// detected on a Linux environment (e.g., in a VM or WSL scenario where BSD tools might be present).
// Note: This file uses a `_linux.go` build tag, implying it's compiled on Linux.
// True BSD systems would ideally have their own `_freebsd.go` or `_openbsd.go` files
// if direct syscalls were involved, but for running commands, this is functional.
type BSDManager struct{}

// Update performs package management operations for BSD-like systems.
func (b *BSDManager) Update(dryRun bool) error {
	log.Info().Msg("--- BSD Package Management ---")

	// Check for FreeBSD's pkg
	if runner.CommandExists("pkg") {
		log.Info().Msg("Detected FreeBSD's 'pkg' package manager.")
		pkgArgs := []string{"upgrade", "-y"} // Upgrade all packages
		if err := runner.RunCommand("Update FreeBSD packages", dryRun, "pkg", nil, pkgArgs...); err != nil {
			log.Error().Err(err).Msg("Failed to update FreeBSD packages.")
			return err
		}

		pkgArgs = []string{"clean", "-a", "-y"} // Clean up unused packages and cache
		if err := runner.RunCommand("Clean FreeBSD pkg cache", dryRun, "pkg", nil, pkgArgs...); err != nil {
			log.Error().Err(err).Msg("Failed to clean FreeBSD pkg cache.")
			return err
		}

		log.Info().Msg("FreeBSD 'pkg' maintenance complete.")
		return nil // Return after successful FreeBSD update
	}

	// Check for OpenBSD's pkg_add
	if runner.CommandExists("pkg_add") {
		log.Info().Msg("Detected OpenBSD's 'pkg_add' package manager.")
		log.Info().Msg("OpenBSD 'pkg_add' does not have a simple 'update all' command.")
		log.Info().Msg("Consider running 'pkg_add -u' for specific packages or reinstalling.")
		log.Info().Msg("OpenBSD 'pkg_add' maintenance advisory complete.")
		return nil // Return after OpenBSD advisory
	}

	log.Debug().Msg("Neither 'pkg' (FreeBSD) nor 'pkg_add' (OpenBSD) package managers found. Skipping BSD package management.")
	return nil // No error if no BSD package manager is found/applicable
}
