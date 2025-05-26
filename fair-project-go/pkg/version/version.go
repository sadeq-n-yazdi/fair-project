package version

import (
	"fmt"
)

// Version information
const (
	// Major version when you make incompatible API changes
	Major = 0
	// Minor version when you add functionality in a backwards compatible manner
	Minor = 0
	// Patch version when you make backwards compatible bug fixes
	Patch = 1
)

// BranchHash is the git branch hash, set at build time
var BranchHash = "unknown"

// String returns the version string in the format "Major.Minor.Patch+BranchHash"
func String() string {
	if BranchHash == "unknown" {
		return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
	}
	return fmt.Sprintf("%d.%d.%d+%s", Major, Minor, Patch, BranchHash)
}

// Map returns the version information as a map
func Map() map[string]interface{} {
	return map[string]interface{}{
		"version":     String(),
		"major":       Major,
		"minor":       Minor,
		"patch":       Patch,
		"branch_hash": BranchHash,
	}
}
