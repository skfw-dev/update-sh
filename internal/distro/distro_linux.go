//go:build linux
// +build linux

package distro

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"update-sh/internal/runner"

	"github.com/rs/zerolog/log" // Changed to zerolog's log
)

// Distribution holds information about the detected Linux distribution.
type Distribution struct {
	ID                    string
	IDLike                string
	Family                string
	PrimaryPackageManager string
}

func (d *Distribution) GetID() string {
	return d.ID
}

func (d *Distribution) GetFamily() string {
	return d.Family
}

func (d *Distribution) GetIDLike() string {
	return d.IDLike
}

func (d *Distribution) GetPrimaryPackageManager() string {
	return d.PrimaryPackageManager
}

func (d *Distribution) String() string {
	return fmt.Sprintf("ID: %s, IDLike: %s, Family: %s, PrimaryPackageManager: %s",
		d.ID, d.IDLike, d.Family, d.PrimaryPackageManager)
}

// DetectDistro detects the Linux distribution and primary package manager.
func DetectDistro() (*Distribution, error) {
	log.Info().Msg("Detecting Linux distribution and primary package manager...")
	dist := &Distribution{
		ID:                    "unknown",
		IDLike:                "unknown",
		Family:                "unknown",
		PrimaryPackageManager: "unknown",
	}

	// Try lsb_release first
	if runner.CommandExists("lsb_release") {
		// Output of lsb_release -is
		cmdIs := exec.Command("lsb_release", "-is")
		outputIs, errIs := cmdIs.Output()
		if errIs == nil {
			dist.ID = strings.TrimSpace(strings.ToLower(string(outputIs)))
		} else {
			log.Debug().Err(errIs).Msg("lsb_release -is failed.")
		}

		// Output of lsb_release -as
		cmdAs := exec.Command("lsb_release", "-as")
		outputAs, errAs := cmdAs.Output()
		if errAs == nil {
			parts := strings.Fields(strings.ToLower(string(outputAs)))
			if len(parts) > 0 {
				dist.IDLike = parts[0]
			}
		} else {
			log.Debug().Err(errAs).Msg("lsb_release -as failed.")
		}
	}

	// Fallback to /etc/os-release if lsb_release fails or is not available
	if dist.ID == "unknown" || dist.ID == "" || strings.Contains(dist.ID, "not found") { // Check for "not found" in ID string
		if _, err := os.Stat("/etc/os-release"); err == nil {
			file, err := os.Open("/etc/os-release")
			if err != nil {
				return dist, fmt.Errorf("failed to open /etc/os-release: %w", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "ID=") {
					dist.ID = strings.Trim(strings.ToLower(strings.TrimPrefix(line, "ID=")), `"`)
				} else if strings.HasPrefix(line, "ID_LIKE=") {
					dist.IDLike = strings.Trim(strings.ToLower(strings.TrimPrefix(line, "ID_LIKE=")), `"`)
				}
			}
			if err := scanner.Err(); err != nil {
				return dist, fmt.Errorf("error reading /etc/os-release: %w", err)
			}
		}
	}

	switch {
	case slices.Contains([]string{"ubuntu", "debian", "linuxmint", "pop", "elementary", "mx"}, dist.ID) ||
		strings.Contains(dist.IDLike, "debian") ||
		strings.Contains(dist.IDLike, "ubuntu") ||
		strings.Contains(dist.IDLike, "linuxmint") ||
		strings.Contains(dist.IDLike, "pop") ||
		strings.Contains(dist.IDLike, "elementary") ||
		strings.Contains(dist.IDLike, "mx"):
		dist.ID = "debian"
		dist.Family = "debian"
		dist.PrimaryPackageManager = "apt"
	case slices.Contains([]string{"rhel", "fedora", "centos", "almalinux", "rocky"}, dist.ID) ||
		strings.Contains(dist.IDLike, "rhel") ||
		strings.Contains(dist.IDLike, "fedora") ||
		strings.Contains(dist.IDLike, "centos") ||
		strings.Contains(dist.IDLike, "almalinux") ||
		strings.Contains(dist.IDLike, "rocky"):
		dist.ID = "rhel"
		dist.Family = "rhel"
		dist.PrimaryPackageManager = "dnf"
	case slices.Contains([]string{"arch", "manjaro", "endeavouros"}, dist.ID) ||
		strings.Contains(dist.IDLike, "arch") ||
		strings.Contains(dist.IDLike, "manjaro") ||
		strings.Contains(dist.IDLike, "endeavouros"):
		dist.ID = "arch"
		dist.Family = "arch"
		dist.PrimaryPackageManager = "pacman"
	case slices.Contains([]string{"opensuse", "sles"}, dist.ID) ||
		strings.Contains(dist.IDLike, "suse") ||
		strings.Contains(dist.IDLike, "sles"):
		dist.ID = "suse"
		dist.Family = "suse"
		dist.PrimaryPackageManager = "zypper"
	case slices.Contains([]string{"gentoo"}, dist.ID) ||
		strings.Contains(dist.IDLike, "gentoo"):
		dist.ID = "gentoo"
		dist.Family = "gentoo"
		dist.PrimaryPackageManager = "portage"
	case slices.Contains([]string{"freebsd"}, dist.ID): // Although BSD, treat as a "distro" for this tool's purposes
		dist.ID = "freebsd"
		dist.Family = "bsd"
		dist.PrimaryPackageManager = "pkg"
	case slices.Contains([]string{"openbsd"}, dist.ID): // Although BSD, treat as a "distro" for this tool's purposes
		dist.ID = "openbsd"
		dist.Family = "bsd"
		dist.PrimaryPackageManager = "pkg_add"
	case strings.Contains(dist.IDLike, "bsd"): // Although BSD, treat as a "distro" for this tool's purposes
		dist.ID = "bsd" // Generic BSD
		dist.Family = "bsd"
		dist.PrimaryPackageManager = "generic_bsd_pkg"
	default:
		log.Info().Msg("Could not definitively determine distribution from lsb_release or /etc/os-release. Falling back to command and config checks.")
	}

	return dist, nil
}
