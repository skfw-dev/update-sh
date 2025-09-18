package pkgmgr

import (
	"fmt"
	"strings"
	"update-sh/internal/runner"

	"github.com/rs/zerolog/log"
)

// CondaManager implements PackageManagerImpl for Conda on any platform.
// This struct will handle updating and cleaning Conda environments in a non-interactive way.
type CondaManager struct{}

// getExistingCondaChannels retrieves a list of channels currently configured in Conda.
// It returns a map for efficient lookups.
// This function assumes that the 'runner' package has a method to capture command output.
func getExistingCondaChannels(user string) (map[string]bool, error) {
	// The '--get channels' command prints the list of channels to stdout.
	args := []string{"config", "--get", "channels"}

	// Note: This requires a runner function that can capture command output.
	// For example: output, err := runner.RunUserCommandAndCaptureOutput(...)
	// We'll simulate this by assuming a successful command execution returns the channel list.
	log.Debug().Msg("Checking for existing Conda channels.")
	output, err := runner.RunUserCommandAndCaptureOutput("Get existing Conda channels", user, "conda", nil, args...)
	if err != nil {
		log.Warn().Err(err).Msg("Could not retrieve existing Conda channels. Will attempt to add all required channels.")
		return nil, err
	}

	channelsMap := make(map[string]bool)
	// Conda's output for '--get channels' is a list of lines, so we parse it line by line.
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Check if the line contains the channel configuration syntax.
		if strings.Contains(trimmedLine, "--add channels '") {
			// Extract the channel name, which is enclosed in single quotes.
			// The expected format is: --add channels 'channel_name'
			parts := strings.Split(trimmedLine, "'")
			if len(parts) >= 2 {
				channelName := parts[1]
				channelsMap[channelName] = true
			}
		}
	}
	return channelsMap, nil
}

// Update performs package updates and cleanup using Conda.
// It ensures required channels are present, updates all packages, and cleans the cache.
func (c *CondaManager) Update(dryRun bool) error {
	log.Info().Msg("--- Conda Package Management ---")

	// Check if the 'conda' command exists in the system's PATH.
	if !runner.CommandExists("conda") {
		log.Debug().Msg("Conda command not found. Skipping Conda package management.")
		return nil
	}

	// Get the target user to run the commands under, as Conda is typically a user-level install.
	user, err := runner.GetTargetUser()
	if err != nil {
		log.Error().Err(err).Msg("Cannot get target user for Conda operations.")
		return err
	}

	// Step 1: Add required channels after checking for their existence.
	log.Info().Msg("Checking and adding required Conda channels...")
	existingChannels, err := getExistingCondaChannels(user)
	if err != nil {
		// If we can't get the existing channels, we'll proceed by attempting to add them all.
		log.Warn().Msg("Failed to retrieve existing channels. Proceeding with a full channel addition attempt.")
	}

	// remove defaults channel to avoid conflicts with default channels in current conda installations
	requiredChannels := []string{"conda-forge", "pytorch", "nvidia", "pypi"}
	// requiredChannels := []string{"defaults", "conda-forge", "pytorch", "nvidia", "pypi"}
	for _, channel := range requiredChannels {
		if existingChannels != nil && existingChannels[channel] {
			log.Info().Msgf("Channel '%s' already exists. Skipping addition.", channel)
			continue
		}

		args := []string{"config", "--add", "channels", channel}
		log.Info().Msgf("Adding channel: %s", channel)
		if err := runner.RunUserCommand(fmt.Sprintf("Add %s channel", channel), dryRun, user, "conda", nil, args...); err != nil {
			log.Warn().Err(err).Msgf("Failed to add channel '%s'. Proceeding with others.", channel)
		}
	}

	// Step 2: Update all packages in the default/base environment.
	log.Info().Msg("Updating all Conda packages...")
	updateArgs := []string{"update", "--all", "--yes"}
	if err := runner.RunUserCommand("Update all Conda packages", dryRun, user, "conda", nil, updateArgs...); err != nil {
		log.Error().Err(err).Msg("Failed to update Conda packages.")
		return err
	}

	// Step 3: Clean up all caches and tarballs.
	log.Info().Msg("Performing Conda cleanup...")
	cleanArgs := []string{"clean", "--all", "--yes"}
	if err := runner.RunUserCommand("Clean Conda cache and old versions", dryRun, user, "conda", nil, cleanArgs...); err != nil {
		log.Warn().Err(err).Msg("Conda cleanup failed or found nothing to clean.")
	} else {
		log.Info().Msg("Conda cleanup complete.")
	}

	log.Info().Msg("Conda maintenance complete.")
	return nil
}
