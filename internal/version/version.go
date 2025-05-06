package version

import "fmt"

var (
	// Version is the current version of the CLI
	Version = "dev"
	// Commit is the git commit hash
	Commit = "unknown"
	// BuildTime is the time the binary was built
	BuildTime = "unknown"
)

// String returns the full version string
func String() string {
	return fmt.Sprintf("titanium version %s (commit: %s, built: %s)", Version, Commit, BuildTime)
}
