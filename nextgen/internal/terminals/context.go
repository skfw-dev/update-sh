package terminals

type Context struct {
	DryRun         bool     `json:"dryRun" yaml:"dry_run"`
	ForceInstall   bool     `json:"forceInstall" yaml:"force_install"`
	NonInteractive bool     `json:"nonInteractive" yaml:"non_interactive"`
	Packages       []string `json:"packages" yaml:"packages"`
}
