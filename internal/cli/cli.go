// Package cli provides the command-line interface commands.
package cli

import (
	"os"

	"github.com/carbonstop/carbonstop-cli/internal/config"
	"github.com/carbonstop/carbonstop-cli/internal/formatter"
)

// ConfigGetter resolves config at execution time (after flags are parsed).
type ConfigGetter func() *config.Config

// CLI holds shared state for all commands.
type CLI struct {
	Config *config.Config
}

// Exit codes
const (
	ExitOK        = 0
	ExitArgs      = 2
	ExitHTTP      = 3
	ExitTransport = 4
)

// validateAPIKey checks if the API key is configured.
func validateAPIKey(cfg *config.Config) bool {
	if cfg.APIKey == "" {
		formatter.Info("missing api key. Get one at https://ccloud-d-test.carbonstop.com/ then run:")
			formatter.Info("  carbonstop auth login --api-key <your-key>")
		os.Exit(ExitArgs)
	}
	return true
}
