package version

import "fmt"

// Version represents a software version with Major, Minor, and Patch components.
type Version struct {
	Major int
	Minor int
	Patch int // Optional, might not always be present for simple X.Y versions
}

// String returns the string representation of the Version.
func (v Version) String() string {
	// Only include Patch if it's set (e.g., non-zero or explicitly needed)
	if v.Patch != 0 || (v.Major == 0 && v.Minor == 0 && v.Patch == 0) { // For 0.0.0 case
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// IsAtLeast checks if this version is at least the specified required version (major.minor).
func (v Version) IsAtLeast(major, minor int) bool {
	if v.Major > major {
		return true
	}
	if v.Major < major {
		return false
	}
	// Major versions are equal, compare minor
	return v.Minor >= minor
}
