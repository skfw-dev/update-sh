package pkgmgr

// PackageManagerImpl defines the common interface for all package managers.
type PackageManagerImpl interface {
	// Update performs the update operation for the specific package manager.
	// dryRun: true if it's a dry run, false otherwise.
	Update(dryRun bool) error
}
