package pkgmgr

import (
	"fmt"
	"regexp"
	"strings"
	"update-sh/internal/runner"

	"github.com/rs/zerolog/log"
)

// This singleton variable holds the compiled regular expression.
// Compiling a regex is a computationally expensive operation, so compiling it once at
// package initialization and reusing it is a standard performance best practice.
// The regex pattern is designed to capture the channel name from a line
// formatted as "--add channels 'channel_name'".
var channelRegex = regexp.MustCompile(`--add channels '([^']+)'`)

// getExistingCondaChannels retrieves a list of channels currently configured in Conda.
// It returns a map for efficient lookups, where each key is a channel name.
// This function relies on the 'conda config --get channels' command and
// expects the output to be a list of `--add channels '...'` lines.
func getExistingCondaChannels(user string) (map[string]bool, error) {
	// We first construct the command to get existing channels, using the 'conda' executable
	// and specifying the required subcommand and arguments.
	cmdArgs := []string{"config", "--get", "channels"}

	// The command is executed using the provided `runner` interface to retrieve the
	// channel list. This abstraction ensures the function remains portable across
	// different operating systems.
	log.Debug().Msg("Attempting to retrieve existing conda channels.")
	output, err := runner.RunUserCommandAndCaptureOutput("Get existing conda channels", user, "conda", nil, cmdArgs...)

	// If the command fails to execute, we log a detailed warning with the error and output.
	// This helps diagnose common issues like Conda not being installed, a user's PATH
	// not being correctly configured, or a permissions' error. A wrapped error is returned
	// to the caller, preserving the original error details.
	if err != nil {
		log.Warn().Err(err).
			Str("output", output).
			Msgf("Failed to run 'conda config --get channels' for user '%s'. This might indicate that conda is not installed or a permissions issue.", user)
		return nil, fmt.Errorf("could not retrieve conda channels: %w", err)
	}

	// The command output is parsed using our pre-compiled regular expression. This approach
	// is not only performant but also robust, as it's designed to be resilient to minor
	// formatting variations in the command's output.
	matches := channelRegex.FindAllStringSubmatch(output, -1)

	// A map is created to store the channel names for efficient lookups.
	// We also check if no channels were found, which is a valid scenario if the Conda
	// configuration is minimal.
	channelsMap := make(map[string]bool)
	if len(matches) == 0 {
		log.Info().Msg("No channels found in conda configuration. The list may be empty.")
		return channelsMap, nil
	}

	// Each match is then iterated through to populate the map. We access the captured
	// group at index 1 to get the channel name. A final validation check ensures the
	// extracted name is not an empty string, logging a warning if an empty name is found
	// to aid in debugging unusual output formats.
	for _, match := range matches {
		if len(match) > 1 {
			channelName := strings.TrimSpace(match[1])
			if channelName != "" {
				channelsMap[channelName] = true
			} else {
				log.Warn().Msg("Detected an empty channel name in the conda configuration output.")
			}
		}
	}

	// Finally, a success message is logged, indicating the number of channels that were
	// successfully retrieved.
	log.Info().Msgf("Successfully retrieved %d existing conda channels.", len(channelsMap))

	return channelsMap, nil
}

// CondaManager implements PackageManagerImpl for Conda on any platform.
// This struct will handle updating and cleaning Conda environments in a non-interactive way.
type CondaManager struct{}

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
	requiredChannels := []string{"defaults", "conda-canary", "conda-forge", "pytorch", "nvidia", "pypi"}
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
