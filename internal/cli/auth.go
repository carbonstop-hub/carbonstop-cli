package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/carbonstop/carbonstop-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewAuthCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage API authentication",
		Long:  "Login, logout, or check authentication status.",
	}
	cmd.AddCommand(newAuthLoginCmd(getCfg))
	cmd.AddCommand(newAuthStatusCmd(getCfg))
	return cmd
}

func newAuthLoginCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save API key to config profile",
		Run: func(cmd *cobra.Command, args []string) {
			apiKey, _ := cmd.Flags().GetString("api-key")
			runAuthLogin(getCfg(), apiKey)
		},
	}
	cmd.Flags().String("api-key", "", "API Key (PAT)")
	return cmd
}

func newAuthStatusCmd(getCfg ConfigGetter) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current auth configuration",
		Run: func(cmd *cobra.Command, args []string) {
			runAuthStatus(getCfg())
		},
	}
}

func runAuthLogin(cfg *config.Config, apiKey string) {
	if apiKey == "" {
		if cfg.APIKey != "" {
			key := cfg.APIKey
			if len(key) > 16 {
				fmt.Printf("Current key: %s...%s\n", key[:12], key[len(key)-4:])
			}
			fmt.Println("(leave blank to keep current key)")
		}
		fmt.Fprint(os.Stderr, "Enter API Key (PAT): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "[carbonstop] failed to read input:", err)
			os.Exit(ExitArgs)
		}
		apiKey = strings.TrimSpace(input)
		if apiKey == "" && cfg.APIKey != "" {
			fmt.Println("Keeping current key.")
			return
		}
	}
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "[carbonstop] API key is required. Get one at https://ccloud-d-test.carbonstop.com/")
		os.Exit(ExitArgs)
	}

	if err := cfg.SaveProfile(apiKey); err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop] failed to save config:", err)
		os.Exit(ExitArgs)
	}
	fmt.Printf("Saved profile '%s' to %s\n", cfg.Profile, config.ConfigFile)
}

func runAuthStatus(cfg *config.Config) {
	fmt.Println("Profile:     ", cfg.Profile)
	fmt.Println("Config file: ", config.ConfigFile)
	if cfg.APIKey != "" {
		key := cfg.APIKey
		if len(key) > 16 {
			fmt.Printf("API Key:      %s...%s\n", key[:12], key[len(key)-4:])
		} else if len(key) > 8 {
			fmt.Printf("API Key:      %s...\n", key[:8])
		} else {
			fmt.Println("API Key:      ***")
		}
	} else {
		fmt.Println("API Key:      (not set)")
		fmt.Println("\nGet started:")
		fmt.Println("  carbonstop auth login --api-key <your-key>")
		fmt.Println("  or: carbonstop auth login  (then paste your key)")
	}
}
