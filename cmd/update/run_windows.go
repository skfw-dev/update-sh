//go:build windows
// +build windows

package update

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"update-sh/internal/distro"
	"update-sh/internal/health"
	"update-sh/internal/pkgmgr"
	"update-sh/internal/shxmgr"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sys/windows"
)

// isAdmin function remains the same as previously corrected.
func isAdmin() (bool, error) {
	// Start impersonation to get an impersonation token
	if err := windows.ImpersonateSelf(windows.SecurityImpersonation); err != nil {
		return false, fmt.Errorf("failed to impersonate self: %w", err)
	}

	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)

	if err != nil {
		return false, fmt.Errorf("failed to initialize SID: %w", err)
	}
	defer windows.FreeSid(sid)

	var token windows.Token
	if err := windows.OpenThreadToken(windows.CurrentThread(), windows.TOKEN_QUERY, true, &token); err != nil {
		return false, fmt.Errorf("failed to open thread token: %w", err)
	}
	defer token.Close()

	isAdmin, err := token.IsMember(sid)
	if err != nil {
		return false, fmt.Errorf("failed to check token membership: %w", err)
	}
	return isAdmin, nil
}

// acquireRoot handles privilege elevation on Windows.
func acquireRoot() {
	admin, err := isAdmin()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to determine administrator status.")
	}

	if !admin {
		log.Info().Msg("Script is not running with administrator privileges. Attempting to re-run as administrator...")

		verb := "runas"
		exe, err := os.Executable()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get current executable path.")
		}
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get current working directory.")
		}
		args := strings.Join(os.Args[1:], " ")

		ptrVerb, _ := syscall.UTF16PtrFromString(verb)
		ptrExe, _ := syscall.UTF16PtrFromString(exe)
		ptrArgs, _ := syscall.UTF16PtrFromString(args)
		ptrCwd, _ := syscall.UTF16PtrFromString(cwd)

		if err := windows.ShellExecute(0, ptrVerb, ptrExe, ptrArgs, ptrCwd, windows.SW_SHOWNORMAL); err != nil {
			log.Fatal().Err(err).Msg("Failed to re-run script as administrator.")
		}

		os.Exit(0)
	} else {
		log.Info().Msg("Script is already running with administrator privileges.")
	}
}

// performWindowsPackageUpdates runs all Windows-specific package manager updates.
func performWindowsPackageUpdates(dryRun bool) {
	var packageManagersToRun []pkgmgr.PackageManagerImpl

	// Add Windows-specific package managers.
	// These managers will internally check if their respective commands (winget, choco) exist.
	packageManagersToRun = append(packageManagersToRun, &pkgmgr.WinGetManager{})
	packageManagersToRun = append(packageManagersToRun, &pkgmgr.ChocolateyManager{})
	packageManagersToRun = append(packageManagersToRun, &pkgmgr.ScoopManager{})
	packageManagersToRun = append(packageManagersToRun, &pkgmgr.CondaManager{})

	// Execute all collected package managers.
	for _, packageManager := range packageManagersToRun {
		if err := packageManager.Update(dryRun); err != nil {
			// Log an error if a specific package manager update fails.
			// %T prints the type of the manager (e.g., *pkgmgr.WinGetManager).
			log.Error().Err(err).Msgf("Windows package manager update failed for %T.", packageManager)
		}
	}
}

func performMaintenance(dryRun, initCheckOnly, zshUpdateEnabled, pwshUpdateEnabled bool) {
	log.Info().Msg("Starting comprehensive system maintenance script.")
	log.Info().Msgf("Log file: %s", viper.GetString("log_file"))

	// Acquire root privileges based on the OS. This function is defined in run_linux.go or run_windows.go
	acquireRoot()

	// Detect distribution and primary package manager first
	d, err := distro.DetectDistro()
	if err != nil {
		log.Error().Err(err).Msg("Error detecting distribution.")
		d = &distro.Distribution{
			ID:                    "unknown",
			IDLike:                "unknown",
			PrimaryPackageManager: "unknown",
		}
	}
	log.Info().Msgf("Detected OS: %s, Distribution ID: %s, Family: %s, Suggested Primary Package Manager: %s", runtime.GOOS, d.ID, d.Family, d.PrimaryPackageManager)

	// --- System Health Checks ---
	healthManager := &health.WindowsHealthManager{}
	if err := healthManager.CheckHealth(dryRun); err != nil {
		log.Error().Err(err).Msgf("System health check failed for %T.", healthManager)
	}

	// --- Shell-specific Updates ---
	var shlexManagersToRun []shxmgr.ShlexManagerImpl

	if zshUpdateEnabled {
		log.Warn().Msg("Zsh update is a Linux-specific feature. Skipping on non-Linux OS.")
	} else {
		log.Info().Msg("Skipping Zsh update. Use '-z' to enable.")
	}

	if pwshUpdateEnabled {
		// Create PwshManager and pass the detected primary package manager
		shlexManagersToRun = append(shlexManagersToRun, &shxmgr.PwshManager{PrimaryPackageManager: d.PrimaryPackageManager})
	} else {
		log.Info().Msg("Skipping PowerShell update. Use '-p' to enable.")
	}

	// Execute all collected shell managers
	for _, shlexManager := range shlexManagersToRun {
		if err := shlexManager.Update(dryRun); err != nil {
			log.Error().Err(err).Msgf("Shell component update failed for %T.", shlexManager)
		}
	}

	// --- Core Package Manager Updates (Platform-specific calls) ---
	if !initCheckOnly {
		log.Info().Msg("--- Starting Core Package Manager Updates ---")
		performWindowsPackageUpdates(dryRun)
		log.Info().Msg("--- Core Package Manager Updates Complete ---")
	} else {
		log.Info().Msg("Skipping core package management updates due to '--init-check' flag.")
	}

	log.Info().Msg("Comprehensive system maintenance complete.")
	if dryRun {
		log.Info().Msg("Remember: This was a DRY RUN. No changes were applied.")
	}
}
