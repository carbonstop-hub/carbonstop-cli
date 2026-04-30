// Package version holds build-time injected values.
// These are set via -ldflags at build time.
package version

import (
	"encoding/base64"
	"strings"
)

var (
	// Version is the semantic version of the CLI.
	Version = "0.2.0-dev"

	// Commit is the git commit hash at build time.
	Commit = "none"

	// BuildTime is the UTC timestamp of the build.
	BuildTime = "unknown"

	// BaseURL is the gateway base URL, set via -ldflags.
	// In release builds it is base64-encoded to prevent plaintext extraction from the binary.
	// In local dev it may be a plaintext URL or "__CONFIG__".
	BaseURL = "__CONFIG__"
)

// GetBaseURL returns the decoded gateway base URL.
func GetBaseURL() string {
	if BaseURL == "" || BaseURL == "__CONFIG__" {
		return ""
	}
	// If it already looks like a URL, return as-is (local dev / env override).
	if strings.Contains(BaseURL, "://") {
		return BaseURL
	}
	// Release builds: base64-encoded to hide from strings/grep.
	decoded, err := base64.StdEncoding.DecodeString(BaseURL)
	if err != nil {
		return ""
	}
	return string(decoded)
}
