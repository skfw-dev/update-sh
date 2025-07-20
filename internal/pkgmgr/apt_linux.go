//go:build linux
// +build linux

package pkgmgr

import (
	"bufio"
	"os/exec"
	"strings"
	"update-sh/internal/runner"

	"github.com/rs/zerolog/log" // Changed to zerolog's log
)

// APTManager implements PackageManagerImpl for APT.
type APTManager struct{}

// Update performs APT package management.
func (a *APTManager) Update(dryRun bool) error {
	log.Info().Msg("--- APT Package Management (Linux) ---")
	if !runner.CommandExists("apt") {
		log.Debug().Msg("APT not found. Skipping APT package management.")
		return nil
	}

	aptArgs := []string{"update", "-y"}
	if err := runner.RunCommand("Update APT package lists", dryRun, "apt", nil, aptArgs...); err != nil {
		return err
	}

	aptArgs = []string{"full-upgrade", "-y"}
	if err := runner.RunCommand("Perform full APT system upgrade", dryRun, "apt", nil, aptArgs...); err != nil {
		return err
	}

	aptArgs = []string{"autoremove", "--purge", "-y"}
	if err := runner.RunCommand("Remove unnecessary APT packages", dryRun, "apt", nil, aptArgs...); err != nil {
		return err
	}

	aptArgs = []string{"autoclean", "-y"}
	if err := runner.RunCommand("Clean up APT cache", dryRun, "apt", nil, aptArgs...); err != nil {
		return err
	}

	log.Info().Msg("APT maintenance complete.")

	// Check for partially removed dpkg packages
	a.checkPartiallyRemovedPackages(dryRun)

	return nil
}

// checkPartiallyRemovedPackages checks for partially removed dpkg packages on Linux.
func (a *APTManager) checkPartiallyRemovedPackages(dryRun bool) {
	log.Info().Msg("--- Checking for Partially Removed Packages (dpkg) ---")
	if dryRun {
		log.Info().Msg("Dry Run: Would check for partially removed dpkg packages.")
		return
	}

	if !runner.CommandExists("dpkg") {
		log.Debug().Msg("dpkg not found. Skipping check for partially removed packages.")
		return
	}

	cmd := exec.Command("dpkg", "--get-selections")
	output, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msg("Failed to run 'dpkg --get-selections'.")
		return
	}

	var packages []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "deinstall") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				packages = append(packages, parts[0])
			}
		}
	}

	if len(packages) > 0 {
		log.Info().Msg("Found partially deinstalled packages:")
		for _, pkg := range packages {
			log.Info().Msgf("  - %s", pkg)
		}
		log.Info().Msg("Consider running 'sudo apt autoremove --purge' if these are APT packages.")
	} else {
		log.Info().Msg("No partially deinstalled packages found.")
	}
}
