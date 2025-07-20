package shxmgr

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"update-sh/internal/runner"
	"update-sh/internal/version"

	"github.com/rs/zerolog/log"
)

// GetPowerShellExecutable determines the preferred PowerShell executable (pwsh.exe for v7+, or powershell.exe).
// It returns the path to the executable and its parsed version.Version struct.
func GetPowerShellExecutable() (string, version.Version, error) {
	// Helper to parse major.minor from a PowerShell output line
	parseVersion := func(output string) (version.Version, error) {
		lines := strings.Split(output, "\n")
		var parts []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				parts = append(parts, trimmed)
			}
		}

		if len(parts) >= 2 { // Expect at least Major and Minor
			major, err := strconv.Atoi(parts[0])
			if err != nil {
				return version.Version{}, fmt.Errorf("invalid major version: %w", err)
			}
			minor, err := strconv.Atoi(parts[1])
			if err != nil {
				return version.Version{}, fmt.Errorf("invalid minor version: %w", err)
			}

			patch := 0
			if len(parts) >= 3 {
				p, err := strconv.Atoi(parts[2])
				if err == nil {
					patch = p
				}
			}

			return version.Version{Major: major, Minor: minor, Patch: patch}, nil
		}

		return version.Version{}, fmt.Errorf("unexpected PowerShell version output format: %s", output)
	}

	// 1. Try pwsh.exe (PowerShell Core / PowerShell 7+)
	if runner.CommandExists("pwsh") {
		// Get Major, Minor, Patch on separate lines for easy parsing
		psArgs := []string{"-NoProfile", "-Command", "$PSVersionTable.PSVersion.Major;$PSVersionTable.PSVersion.Minor;$PSVersionTable.PSVersion.Patch"}
		cmd := exec.Command("pwsh", psArgs...)
		output, err := cmd.CombinedOutput()
		if err == nil {
			psVersion, parseErr := parseVersion(string(output))
			if parseErr == nil {
				if psVersion.IsAtLeast(7, 0) { // Prioritize pwsh if it's v7 or higher
					log.Debug().Msgf("Detected PowerShell Core (pwsh.exe) version %s", psVersion.String())
					return "pwsh.exe", psVersion, nil
				}
			} else {
				log.Debug().Err(parseErr).Msgf("Failed to parse pwsh.exe version from '%s'", strings.TrimSpace(string(output)))
			}
		} else {
			log.Debug().Err(err).Msgf("Failed to get pwsh.exe version info from command output '%s'", strings.TrimSpace(string(output)))
		}
		log.Debug().Msg("pwsh.exe found but not v7+ or version check failed, falling back to powershell.exe if needed.")
	}

	// 2. Fallback to powershell.exe (Windows PowerShell 5.1)
	if runner.CommandExists("powershell.exe") {
		psArgs := []string{"-NoProfile", "-Command", "$PSVersionTable.PSVersion.Major;$PSVersionTable.PSVersion.Minor;$PSVersionTable.PSVersion.Patch"}
		cmd := exec.Command("powershell.exe", psArgs...)
		output, err := cmd.CombinedOutput()
		if err == nil {
			psVersion, parseErr := parseVersion(string(output))
			if parseErr == nil {
				log.Debug().Msgf("Detected Windows PowerShell (powershell.exe) version %s", psVersion.String())
				return "powershell.exe", psVersion, nil
			} else {
				log.Warn().Err(parseErr).Msgf("Failed to parse powershell.exe version from '%s'", strings.TrimSpace(string(output)))
			}
		} else {
			log.Warn().Err(err).Msgf("Failed to get powershell.exe version info from command output '%s'", strings.TrimSpace(string(output)))
		}
		// Return 0.0.0 if version check fails but powershell.exe is found
		return "powershell.exe", version.Version{}, nil
	}

	return "", version.Version{}, fmt.Errorf("neither pwsh.exe (PowerShell 7+) nor powershell.exe (Windows PowerShell 5.1) found")
}

// SetExecutionPolicy checks and sets the PowerShell execution policy for the current user.
// It prioritizes PowerShell 7+ (pwsh.exe) if available.
func SetExecutionPolicy(dryRun bool) error {
	if dryRun {
		log.Info().Msg("Dry Run: Would check and set PowerShell execution policy.")
		return nil
	}

	log.Debug().Msg("Checking and setting PowerShell execution policy...")

	psExe, psVersion, err := GetPowerShellExecutable()
	if err != nil {
		log.Error().Err(err).Msg("Cannot find a suitable PowerShell executable.")
		log.Error().Msg("Please install PowerShell 7 (pwsh.exe) from: https://aka.ms/powershell-release?tag=stable")
		return err
	}

	log.Info().Msgf("Using PowerShell executable: %s (version %s) to check/set execution policy.", psExe, psVersion.String())

	// Get the current execution policy for the CurrentUser scope
	psArgs := []string{"-NoProfile", "-Command", "Get-ExecutionPolicy -Scope CurrentUser -ErrorAction SilentlyContinue | Out-String -Stream"}
	getPolicyCmd := exec.Command(psExe, psArgs...)
	var stdout, stderr strings.Builder
	getPolicyCmd.Stdout = &stdout
	getPolicyCmd.Stderr = &stderr

	err = getPolicyCmd.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get current PowerShell execution policy. Stdout: '%s', Stderr: '%s'", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()))
		return fmt.Errorf("failed to get PowerShell execution policy: %w, stderr: %s", err, strings.TrimSpace(stderr.String()))
	}
	currentPolicy := strings.TrimSpace(stdout.String())
	log.Debug().Msgf("Current CurrentUser execution policy: '%s'", currentPolicy)

	// Define policies that are suitable for running scripts
	suitablePolicies := []string{"RemoteSigned", "Unrestricted", "Bypass"}
	isSuitable := false
	for _, p := range suitablePolicies {
		if strings.EqualFold(currentPolicy, p) { // Case-insensitive comparison
			isSuitable = true
			break
		}
	}

	if !isSuitable {
		log.Warn().Msgf("CurrentUser execution policy is '%s'. Setting to 'RemoteSigned' for package manager script execution.", currentPolicy)
		// Set the execution policy to RemoteSigned for the CurrentUser scope
		psArgs := []string{"-NoProfile", "-Command", "Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force -ErrorAction Stop | Out-String -Stream"}
		setPolicyCmd := exec.Command(psExe, psArgs...)
		stdout.Reset() // Clear buffers for next command
		stderr.Reset()
		setPolicyCmd.Stdout = &stdout
		setPolicyCmd.Stderr = &stderr

		setError := setPolicyCmd.Run()
		if setError != nil {
			log.Error().Err(setError).Msgf("Failed to set CurrentUser execution policy. Stdout: '%s', Stderr: '%s'", strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()))
			return fmt.Errorf("failed to set PowerShell execution policy: %w, stderr: %s", setError, strings.TrimSpace(stderr.String()))
		}
		log.Info().Msg("Execution policy for CurrentUser successfully set to 'RemoteSigned'.")
	} else {
		log.Debug().Msgf("CurrentUser execution policy is '%s', which is suitable for package managers.", currentPolicy)
	}
	return nil
}
