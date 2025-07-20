//go:build linux
// +build linux

package health

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"update-sh/internal/runner"

	"github.com/rs/zerolog/log"
)

// LinuxHealthManager implements HealthImpl for Linux systems.
type LinuxHealthManager struct{}

// CheckHealth performs comprehensive Linux health checks.
func (l *LinuxHealthManager) CheckHealth(dryRun bool) error {
	log.Info().Msg("--- Starting Linux System Health Checks ---")

	// Check System Init
	l.checkSystemInit(dryRun)

	log.Info().Msg("--- Linux System Health Checks Complete ---")
	return nil
}

// checkFailedSystemdUnitsSystem checks for failed systemd units (system scope) on Linux.
func (l *LinuxHealthManager) checkFailedSystemdUnitsSystem(dryRun bool) {
	log.Info().Msg("--- Checking for Failed Systemd Units (System Scope) ---")
	if dryRun {
		log.Info().Msg("Dry Run: Would check for failed system-scope systemd units.")
		return
	}

	if !runner.CommandExists("systemctl") {
		log.Debug().Msg("systemctl not found. Skipping systemd unit checks.")
		return
	}

	cmd := exec.Command("systemctl", "list-units", "--system", "--failed", "--no-pager", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		if len(output) == 0 && err.Error() == "exit status 1" {
			log.Info().Msg("No failed system-scope units found.")
			return
		}
		log.Error().Err(err).Msgf("Failed to check system-scope systemd units. Output:\n%s", strings.TrimSpace(string(output)))
		return
	}

	failedUnits := strings.TrimSpace(string(output))
	if failedUnits != "" {
		log.Info().Msg("Found failed system-scope units:")
		log.Info().Msg(failedUnits)
	} else {
		log.Info().Msg("No failed system-scope units found.")
	}
}

// checkFailedSystemdUnitsUser checks for failed systemd units (user scope) on Linux.
func (l *LinuxHealthManager) checkFailedSystemdUnitsUser(dryRun bool) {
	log.Info().Msg("--- Checking for Failed Systemd Units (User Scope) ---")
	if dryRun {
		log.Info().Msg("Dry Run: Would check for failed user-scope systemd units.")
		return
	}

	if !runner.CommandExists("systemctl") || !runner.CommandExists("dbus-launch") {
		log.Debug().Msg("systemctl or dbus-launch not found. Skipping user-scope systemd unit checks.")
		return
	}

	user, err := runner.GetTargetUser()
	if err != nil {
		log.Error().Err(err).Msg("Cannot check user-scope systemd units.")
		return
	}

	log.Info().Msgf("Attempting to check user-scope systemd units for user: %s", user)

	var failedUnitsUser string

	// Attempt to get the DBUS_SESSION_BUS_ADDRESS and XDG_RUNTIME_DIR for the user.
	cmd := exec.Command("sudo", "-u", user, "env")
	output, err := cmd.Output()
	if err != nil {
		log.Warn().Err(err).Msgf("Could not retrieve user environment for %s. Proceeding with common user session paths.", user)
	}

	userEnv := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			userEnv[parts[0]] = parts[1]
		}
	}

	dbusSessionBusAddress := userEnv["DBUS_SESSION_BUS_ADDRESS"]
	xdgRuntimeDir := userEnv["XDG_RUNTIME_DIR"]

	if dbusSessionBusAddress != "" && xdgRuntimeDir != "" {
		log.Debug().Msgf("Using retrieved DBus environment for user %s.", user)
		cmd = exec.Command("sudo", "-u", user, "systemctl", "--user", "list-units", "--failed", "--no-pager", "--no-legend")
		cmd.Env = append(cmd.Env, "DBUS_SESSION_BUS_ADDRESS="+dbusSessionBusAddress)
		cmd.Env = append(cmd.Env, "XDG_RUNTIME_DIR="+xdgRuntimeDir)
		output, err = cmd.Output()
		if err == nil {
			failedUnitsUser = strings.TrimSpace(string(output))
		}
	} else {
		log.Info().Msg("Warning: Could not retrieve DBUS_SESSION_BUS_ADDRESS or XDG_RUNTIME_DIR for user directly. Attempting with dbus-launch.")

		cmd = exec.Command("sudo", "-u", user, "dbus-launch", "systemctl", "--user", "list-units", "--failed", "--no-pager", "--no-legend")
		output, err = cmd.Output()
		if err == nil {
			failedUnitsUser = strings.TrimSpace(string(output))
		}
	}

	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		log.Error().Err(err).Msgf("Failed to check user-scope systemd units for %s. Output:\n%s", user, strings.TrimSpace(string(output)))
		return
	}

	if failedUnitsUser != "" {
		log.Info().Msgf("Found failed user-scope units for %s:", user)
		log.Info().Msg(failedUnitsUser)
	} else {
		log.Info().Msgf("No failed user-scope units found for %s (or could not connect to user D-Bus session).", user)
	}
}

// checkSystemInit determines and checks the primary system init system on Linux.
func (l *LinuxHealthManager) checkSystemInit(dryRun bool) {
	log.Info().Msg("--- Checking System Init System ---")
	initSystem := "Unknown"

	if _, err := os.Stat("/run/systemd/system"); err == nil {
		initSystem = "systemd"
		log.Info().Msg("Detected init system: systemd.")
		l.checkFailedSystemdUnitsSystem(dryRun)
		l.checkFailedSystemdUnitsUser(dryRun)
	} else if runner.CommandExists("initctl") {
		cmd := exec.Command("initctl", "--version")
		output, err := cmd.Output()
		if err == nil && strings.Contains(string(output), "Upstart") {
			initSystem = "Upstart"
			log.Info().Msg("Detected init system: Upstart.")
			log.Info().Msg("Upstart does not have a direct equivalent to 'list failed units' like systemd.")
			log.Info().Msg("You might want to check '/var/log/syslog' or 'dmesg' for Upstart service errors.")
		}
	} else if _, err := os.Stat("/etc/init.d/rcS"); err == nil {
		initSystem = "SysVinit"
		log.Info().Msg("Detected init system: SysVinit.")
		log.Info().Msg("SysVinit does not have a direct equivalent to 'list failed units' like systemd.")
		log.Info().Msg("You might want to check '/var/log/messages' or '/var/log/syslog' for service errors.")
	} else {
		log.Info().Msg("Could not definitively determine the primary init system.")
		log.Info().Msg("Common init systems include systemd, Upstart, and SysVinit.")
	}
	log.Info().Msgf("System init check complete. Detected: %s", initSystem)
}
