//go:build windows
// +build windows

package distro

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// Distribution holds information about the detected Windows environment.
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

// DetectDistro detects the Windows environment.
// On Windows, this is simpler as there's no "distribution" in the Linux sense.
func DetectDistro() (*Distribution, error) {
	log.Info().Msg("Detecting Windows environment and primary package manager...")
	dist := &Distribution{
		ID:                    "windows",
		IDLike:                "windows",
		Family:                "windows",
		PrimaryPackageManager: "winget", // Assume Winget as the primary for now
	}
	log.Info().Msgf("Detected OS: Windows, Primary Package Manager: %s", dist.PrimaryPackageManager)
	return dist, nil
}
