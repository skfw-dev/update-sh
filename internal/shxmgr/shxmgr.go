package shxmgr

// ShlexManagerImpl defines the common interface for all shell-related update operations.
type ShlexManagerImpl interface {
	// Update performs the update operation for the specific shell component.
	// dryRun: true if it's a dry run, false otherwise.
	Update(dryRun bool) error
}
