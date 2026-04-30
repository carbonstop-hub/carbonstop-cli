package main

import (
	"fmt"
	"os"

	"github.com/carbonstop/carbonstop-cli/internal/cli"
	"github.com/carbonstop/carbonstop-cli/internal/config"
	"github.com/carbonstop/carbonstop-cli/internal/version"
	"github.com/spf13/cobra"
)

var (
	flagBaseURL string
	flagAPIKey  string
	flagTimeout int
	flagProfile string
	flagRaw     bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "carbonstop",
	Short: "Carbonstop CLI — gateway client for Carbonstop API",
	Long: `carbonstop is a command-line tool for interacting with the Carbonstop Gateway API.

Get started:
  1. Register at https://ccloud-d-test.carbonstop.com/
  2. Create an API Key (PAT)
  3. Run: carbonstop auth login --api-key <your-key>
  4. Try:  carbonstop +ping

Commands:
  - Raw API: products, accounts, ai-model, ping, whoami, echo, call, search-factors
  - Shortcuts (+): +ping, +whoami, +model
  - Auth: auth login, auth status

Configure via env vars, config file, or CLI flags.
  CARBONSTOP_API_KEY, CARBONSTOP_BASE_URL, CARBONSTOP_TIMEOUT`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version.Version, version.Commit, version.BuildTime),
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagBaseURL, "base-url", "", "Gateway base URL")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "API Key (PAT)")
	rootCmd.PersistentFlags().IntVar(&flagTimeout, "timeout", 0, "Request timeout in seconds")
	rootCmd.PersistentFlags().StringVarP(&flagProfile, "profile", "p", "", "Config profile name")
	rootCmd.PersistentFlags().BoolVar(&flagRaw, "raw", false, "Output raw JSON (no pretty-print)")

	getConfig := func() *config.Config {
		return config.New(flagBaseURL, flagAPIKey, flagTimeout, flagProfile)
	}

	rootCmd.AddCommand(cli.NewPingCmd(getConfig))
	rootCmd.AddCommand(cli.NewWhoamiCmd(getConfig))
	rootCmd.AddCommand(cli.NewEchoCmd(getConfig))
	rootCmd.AddCommand(cli.NewCallCmd(getConfig))

	rootCmd.AddCommand(cli.NewProductsCmd(getConfig))
	rootCmd.AddCommand(cli.NewProductInfoCmd(getConfig))
	rootCmd.AddCommand(cli.NewAccountsCmd(getConfig))
	rootCmd.AddCommand(cli.NewAccountViewCmd(getConfig))
	rootCmd.AddCommand(cli.NewAiModelCmd(getConfig))
	rootCmd.AddCommand(cli.NewSearchFactorsCmd(getConfig))

	rootCmd.AddCommand(cli.NewPlusPingCmd(getConfig))
	rootCmd.AddCommand(cli.NewPlusWhoamiCmd(getConfig))
	rootCmd.AddCommand(cli.NewPlusModelCmd(getConfig))

	rootCmd.AddCommand(cli.NewAuthCmd(getConfig))
}
