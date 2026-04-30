package cli

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"

	"github.com/carbonstop/carbonstop-cli/internal/apipath"
	"github.com/carbonstop/carbonstop-cli/internal/client"
	"github.com/carbonstop/carbonstop-cli/internal/config"
	"github.com/spf13/cobra"
)

func NewSearchFactorsCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search-factors",
		Short: "Search CCDB carbon emission factors",
		Long:  "Search the China Carbon Factor Database (CCDB) by keyword. Returns matching factors with metadata.",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			lang, _ := cmd.Flags().GetString("lang")
			runSearchFactors(getCfg(), name, lang, getRaw(cmd))
		},
	}
	cmd.Flags().String("name", "", "Search keyword (required)")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("lang", "zh", "Language: zh or en (default zh)")
	return cmd
}

func runSearchFactors(cfg *config.Config, name, lang string, raw bool) {
	validateAPIKey(cfg)

	signStr := "openclaw_ccdb" + name
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))

	body := map[string]string{
		"name": name,
		"sign": sign,
		"lang": lang,
	}
	payload, _ := json.Marshal(body)

	c := client.New(cfg)
	status, bodyStr, err := c.Post(apipath.SearchFactor(), payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	if status < 200 || status >= 300 {
		fmt.Fprintln(os.Stderr, "[carbonstop] http status", status)
		fmt.Println(bodyStr)
		os.Exit(ExitHTTP)
	}

	var parsed interface{}
	if json.Unmarshal([]byte(bodyStr), &parsed) != nil {
		fmt.Println(bodyStr)
		return
	}
	pretty, _ := json.MarshalIndent(parsed, "", "  ")
	fmt.Println(string(pretty))
}
