package health

// HealthImpl defines the common interface for all system health checks.
type HealthImpl interface {
	// CheckHealth performs the health check operation for the specific OS.
	// dryRun: true if it's a dry run, false otherwise.
	CheckHealth(dryRun bool) error
}
