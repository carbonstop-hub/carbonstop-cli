package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/carbonstop/carbonstop-cli/internal/apipath"
	"github.com/carbonstop/carbonstop-cli/internal/client"
	"github.com/carbonstop/carbonstop-cli/internal/config"
	"github.com/carbonstop/carbonstop-cli/internal/formatter"
	"github.com/spf13/cobra"
)

func NewAiModelCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai-model",
		Short: "One-click AI carbon footprint modeling",
		Run: func(cmd *cobra.Command, args []string) {
			data, _ := cmd.Flags().GetString("data")
			file, _ := cmd.Flags().GetString("file")
			runAiModel(getCfg(), data, file, getRaw(cmd))
		},
	}
	cmd.Flags().String("data", "", "Inline JSON body")
	cmd.Flags().String("file", "", "Path to JSON file ('-' for stdin)")
	return cmd
}

func runAiModel(cfg *config.Config, data, file string, raw bool) {
	validateAPIKey(cfg)
	payload, err := resolvePayload(data, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitArgs)
	}
	if payload == nil {
		fmt.Fprintln(os.Stderr, "[carbonstop] ai-model requires --data or --file")
		os.Exit(ExitArgs)
	}

	c := client.New(cfg)
	f := formatter.New(raw)
	status, body, err := c.Post(apipath.AiModel(), payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(status, body)
	if status < 200 || status >= 300 {
		os.Exit(ExitHTTP)
	}
}

func resolvePayload(data, file string) ([]byte, error) {
	if data != "" && file != "" {
		return nil, fmt.Errorf("--data and --file are mutually exclusive")
	}

	var raw string
	if data != "" {
		raw = data
	} else if file != "" {
		if file == "-" {
			stdinBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("reading stdin: %w", err)
			}
			raw = strings.TrimSpace(string(stdinBytes))
		} else {
			fileBytes, err := os.ReadFile(file)
			if err != nil {
				return nil, fmt.Errorf("reading file: %w", err)
			}
			raw = strings.TrimSpace(string(fileBytes))
		}
	} else {
		return nil, nil
	}

	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return nil, fmt.Errorf("payload is not valid JSON: %w", err)
	}
	return []byte(raw), nil
}
