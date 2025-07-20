//go:build windows
// +build windows

package runner

import (
	"fmt"
	"os/exec"
	"os/user" // Still needed for user.Current() if you want to get current user

	"github.com/rs/zerolog/log"
	// No need to import "update-sh/internal/config" here for GetTargetUser,
	// as GetTargetUser is not applicable for Windows in this context.
)

// RunUserCommandWithOptions runs a command as a specific user on Windows.
// On Windows, we don't use sudo -u like on Linux.
// Instead, we just run the command directly as the current user.
func RunUserCommandWithOptions(opts *CommandOptions) error {
	if opts.DryRun {
		log.Info().Msgf("Dry Run: Would execute '%s': %s %v", opts.Description, opts.Name, opts.Args)
		return nil
	}

	log.Info().Msgf("%s (as user %s)...", opts.Description, opts.Name)

	// Build the command to run
	// On Windows, we don't use sudo -u like on Linux.
	// Instead, we just run the command directly as the current user.
	cmd := exec.Command(opts.Name, opts.Args...)
	if len(opts.Env) > 0 {
		cmd.Env = append(cmd.Env, opts.Env...)
	}

	// Use a transformer for encoding if specified
	decoder, err := makeDecoder(opts.Encoding)
	if err != nil {
		return fmt.Errorf("failed to create decoder for encoding %s: %w", opts.Encoding.String(), err)
	}

	if opts.User == "" {
		return fmt.Errorf("no user specified for running command")
	}

	// Custom zerolog console writer
	// cmd.Stdout = zerolog.ConsoleWriter{Out: log.Logger.Output(os.Stdout), TimeFormat: zerolog.TimeFormatUnix}
	// cmd.Stderr = zerolog.ConsoleWriter{Out: log.Logger.Output(os.Stderr), TimeFormat: zerolog.TimeFormatUnix}
	return streamAndWait(cmd, decoder, opts.Description, opts.User)
}

// RunUserCommand on Windows simply runs the command.
func RunUserCommand(description string, dryRun bool, user string, name string, env []string, arg ...string) error {
	opts := NewCommandOptions(description, dryRun, name, env, arg...)
	opts.User = user // Set the user for the command options
	return RunCommandWithOptions(opts)
}

// GetTargetUser is not directly applicable on Windows in the same way as Linux (UIDs).
// If you need the current Windows username, use os/user.Current().
func GetTargetUser() (string, error) {
	// On Windows, the concept of a numeric UID for a "target user" is not used for this purpose.
	// Package managers like Scoop are installed per-user and operate on the current user's context.
	// If the script needs to know the *current* user, os/user.Current() can be used.
	// For the purpose of "running as another user" for package managers, this function is not directly relevant.
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user on Windows: %w", err)
	}
	log.Debug().Msg("GetTargetUser called on Windows. Returning current user as target.")
	return currentUser.Username, nil
}
