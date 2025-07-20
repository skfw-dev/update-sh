//go:build linux
// +build linux

package shxmgr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"update-sh/internal/runner"

	"github.com/rs/zerolog/log" // Changed to zerolog's log
)

// ZshManager implements ShlexManagerImpl for Zsh-related components on Linux.
type ZshManager struct{}

// Update performs updates for Oh My Zsh, Powerlevel10k, and Oh My Posh CLI on Linux.
func (z *ZshManager) Update(dryRun bool) error {
	log.Info().Msg("--- Zsh (Oh My Zsh & Powerlevel10k) Update (Linux) ---")

	user, err := runner.GetTargetUser() // Use the updated GetTargetUser
	if err != nil {
		log.Error().Err(err).Msg("Cannot update user-specific Zsh components.")
		return err // Return error for the interface
	}

	cmdUserHomeDir := exec.Command("sudo", "-u", user, "printenv", "HOME")
	userHomeDirBytes, err := cmdUserHomeDir.Output()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get home directory for user %s.", user)
		return err // Return error
	}
	homeDir := strings.TrimSpace(string(userHomeDirBytes))

	ohMyZshPath := filepath.Join(homeDir, ".oh-my-zsh")
	powerlevel10kPath := filepath.Join(ohMyZshPath, "custom", "themes", "powerlevel10k")

	log.Info().Msgf("Checking Zsh components for user: %s in %s", user, homeDir)

	if !runner.CommandExists("git") {
		log.Error().Msg("'git' is not installed. Required for Zsh component updates. Skipping.")
		return fmt.Errorf("'git' is not installed, required for Zsh component updates") // Return specific error
	}

	// Update Oh My Zsh
	log.Info().Msg("Attempting to update Oh My Zsh using 'omz update'...")
	if err := runner.RunUserCommand("Update Oh My Zsh", dryRun, user, "zsh", nil, "-i", "-c", "omz update --unattended"); err == nil {
		log.Debug().Msg("Oh My Zsh updated using 'omz update'.")
	} else {
		log.Warn().Err(err).Msg("Failed to update Oh My Zsh using 'omz update'. Attempting 'git pull'.")
		if err := runner.RunUserCommand("Update Oh My Zsh (git pull)", dryRun, user, "git", nil, "-C", ohMyZshPath, "pull"); err != nil {
			log.Error().Err(err).Msg("Failed to update Oh My Zsh using 'git pull'.")
			// Decide if this is a fatal error or if other updates can proceed.
			// For now, let's allow it to continue but mark the overall update as failed if this part fails.
			return fmt.Errorf("failed to update Oh My Zsh: %w", err)
		} else {
			log.Debug().Msg("Oh My Zsh updated using 'git pull'.")
		}
	}

	// Update Powerlevel10k
	if _, err := os.Stat(powerlevel10kPath); err == nil {
		log.Info().Msgf("Found Powerlevel10k theme at %s.", powerlevel10kPath)
		if err := runner.RunUserCommand("Update Powerlevel10k", dryRun, user, "git", nil, "-C", powerlevel10kPath, "pull"); err != nil {
			log.Error().Err(err).Msg("Failed to update Powerlevel10k.")
			return fmt.Errorf("failed to update Powerlevel10k: %w", err)
		}
	} else {
		log.Debug().Msgf("Powerlevel10k not found at %s. Skipping Powerlevel10k update.", powerlevel10kPath)
	}

	// Update Oh My Posh CLI (can be cross-platform, but often installed via package managers or specific scripts)
	if runner.CommandExists("oh-my-posh") {
		log.Info().Msg("Found Oh My Posh CLI. Attempting to upgrade...")
		if err := runner.RunCommand("Upgrade Oh My Posh CLI", dryRun, "oh-my-posh", nil, "upgrade", "--force"); err != nil {
			log.Error().Err(err).Msg("Failed to upgrade Oh My Posh CLI.")
			return fmt.Errorf("failed to upgrade Oh My Posh CLI: %w", err)
		}
	} else {
		log.Debug().Msg("Oh My Posh CLI not found. Skipping Oh My Posh CLI update.")
	}
	log.Info().Msg("Zsh components update complete.")
	return nil
}
