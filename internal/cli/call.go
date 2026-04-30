package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/carbonstop/carbonstop-cli/internal/client"
	"github.com/carbonstop/carbonstop-cli/internal/config"
	"github.com/carbonstop/carbonstop-cli/internal/formatter"
	"github.com/spf13/cobra"
)

func NewCallCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call <method> <path>",
		Short: "Generic API call (arbitrary method + path)",
		Long: `Make an arbitrary HTTP call to the gateway API.

Examples:
  carbonstop call GET /api/some/endpoint -q id=123
  carbonstop call POST /api/some/endpoint -d '{"hello":"world"}'
  carbonstop call POST /api/some/endpoint -f payload.json`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			method := strings.ToUpper(args[0])
			path := args[1]
			data, _ := cmd.Flags().GetString("data")
			file, _ := cmd.Flags().GetString("file")
			queryStr, _ := cmd.Flags().GetStringToString("query")
			headerStr, _ := cmd.Flags().GetStringToString("header")
			runCall(getCfg(), method, path, data, file, queryStr, headerStr, getRaw(cmd))
		},
	}
	cmd.Flags().String("data", "", "Inline JSON request body")
	cmd.Flags().String("file", "", "Path to JSON file for request body ('-' for stdin)")
	cmd.Flags().StringToStringP("query", "q", nil, "Query parameters (key=value)")
	cmd.Flags().StringToStringP("header", "H", nil, "Extra headers (key=value)")
	return cmd
}

func runCall(cfg *config.Config, method, path, data, file string, query, extraHeaders map[string]string, raw bool) {
	validateAPIKey(cfg)
	payload, err := resolvePayload(data, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitArgs)
	}

	c := client.New(cfg)
	f := formatter.New(raw)

	var bodyReader io.Reader
	if payload != nil {
		bodyReader = bytes.NewReader(payload)
		if extraHeaders == nil {
			extraHeaders = make(map[string]string)
		}
		if extraHeaders["Content-Type"] == "" {
			extraHeaders["Content-Type"] = "application/json;charset=utf-8"
		}
	}
	status, body, err := c.Request(method, path, bodyReader, query, extraHeaders)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(status, body)
	if status < 200 || status >= 300 {
		os.Exit(ExitHTTP)
	}
}
