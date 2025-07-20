//go:build linux
// +build linux

package runner

import (
	"errors"
	"fmt"
	"os/exec"
	"os/user"

	// Needed for potential syscall.Credential if you ever go that route

	"github.com/rs/zerolog/log"

	"update-sh/internal/config"
)

func RunUserCommandWithOptions(opts *CommandOptions) error {
	if opts.DryRun {
		log.Info().Msgf("Dry Run: Would execute '%s' as user '%s': %s %v", opts.Description, opts.User, opts.Name, opts.Args)
		return nil
	}

	log.Info().Msgf("%s (as user %s)...", opts.Description, opts.User)

	// Construct the command to run via sudo -u
	sudoArgs := []string{"-u", opts.Name}
	sudoArgs = append(sudoArgs, opts.Args...)

	// Build the command with sudo
	cmd := exec.Command("sudo", sudoArgs...)
	if len(opts.Env) > 0 {
		cmd.Env = append(cmd.Env, opts.Env...)
	}

	// Use a transformer for encoding if specified
	decoder, err := makeDecoder(opts.Encoding)
	if err != nil {
		return fmt.Errorf("failed to create decoder for encoding %s: %w", opts.Encoding.String(), err)
	}

	if opts.User == "" {
		return errors.New("no user specified for running command")
	}

	// Custom zerolog console writer
	// cmd.Stdout = zerolog.ConsoleWriter{Out: log.Logger.Output(os.Stdout), TimeFormat: zerolog.TimeFormatUnix}
	// cmd.Stderr = zerolog.ConsoleWriter{Out: log.Logger.Output(os.Stderr), TimeFormat: zerolog.TimeFormatUnix}
	return streamAndWait(cmd, decoder, opts.Description, opts.User)
}

// RunUserCommand executes a command as a specific user on Linux/Unix-like systems.
// It typically uses 'sudo -u' to change user context.
func RunUserCommand(description string, dryRun bool, user string, name string, env []string, arg ...string) error {
	opts := NewCommandOptions(description, dryRun, name, env, arg...)
	opts.User = user // Set the user for the command options
	return RunCommandWithOptions(opts)
}

// GetTargetUser retrieves the username for a given UID on Linux/Unix-like systems.
func GetTargetUser() (string, error) {
	// Get the platform-specific config manager
	cfg := config.GetConfigManager()

	var userID string
	// Use the default UID from the platform-specific config
	userID = cfg.GetDefaultUserID()
	if userID == "" {
		return "", fmt.Errorf("default user ID not applicable or not found for this OS")
	}

	targetUser, err := user.LookupId(userID)
	if err != nil {
		return "", fmt.Errorf("no user found for ID %s: %w", userID, err)
	}
	return targetUser.Username, nil
}
