package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/carbonstop/carbonstop-cli/internal/apipath"
	"github.com/carbonstop/carbonstop-cli/internal/client"
	"github.com/carbonstop/carbonstop-cli/internal/config"
	"github.com/carbonstop/carbonstop-cli/internal/formatter"
	"github.com/spf13/cobra"
)

func NewPingCmd(getCfg ConfigGetter) *cobra.Command {
	return &cobra.Command{
		Use:   "ping",
		Short: "Health check",
		Run: func(cmd *cobra.Command, args []string) {
			runPing(getCfg(), getRaw(cmd))
		},
	}
}

func NewWhoamiCmd(getCfg ConfigGetter) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show identity",
		Run: func(cmd *cobra.Command, args []string) {
			runWhoami(getCfg(), getRaw(cmd))
		},
	}
}

func NewEchoCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "echo",
		Short: "JSON echo test",
		Run: func(cmd *cobra.Command, args []string) {
			data, _ := cmd.Flags().GetString("data")
			file, _ := cmd.Flags().GetString("file")
			runEcho(getCfg(), data, file, getRaw(cmd))
		},
	}
	cmd.Flags().String("data", "", "Inline JSON body")
	cmd.Flags().String("file", "", "Path to JSON file ('-' for stdin)")
	return cmd
}

func NewPlusPingCmd(getCfg ConfigGetter) *cobra.Command {
	return &cobra.Command{
		Use:   "+ping",
		Short: "探活（快捷，带状态提示）",
		Run: func(cmd *cobra.Command, args []string) {
			runPlusPing(getCfg())
		},
	}
}

func NewPlusWhoamiCmd(getCfg ConfigGetter) *cobra.Command {
	return &cobra.Command{
		Use:   "+whoami",
		Short: "身份信息（快捷，可读输出）",
		Run: func(cmd *cobra.Command, args []string) {
			runPlusWhoami(getCfg())
		},
	}
}

func NewPlusModelCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "+model",
		Short: "一键建模（快捷，含阶段汇总）",
		Run: func(cmd *cobra.Command, args []string) {
			data, _ := cmd.Flags().GetString("data")
			file, _ := cmd.Flags().GetString("file")
			runPlusModel(getCfg(), data, file)
		},
	}
	cmd.Flags().String("data", "", "Inline JSON body")
	cmd.Flags().String("file", "", "Path to JSON file")
	return cmd
}

// ── Runner functions ──────────────────────────────────

func runPing(cfg *config.Config, raw bool) {
	validateAPIKey(cfg)
	c := client.New(cfg)
	f := formatter.New(raw)
	status, body, err := c.Get(apipath.Ping(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(status, body)
	if status < 200 || status >= 300 {
		os.Exit(ExitHTTP)
	}
}

func runWhoami(cfg *config.Config, raw bool) {
	validateAPIKey(cfg)
	c := client.New(cfg)
	f := formatter.New(raw)
	status, body, err := c.Get(apipath.Whoami(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(status, body)
	if status < 200 || status >= 300 {
		os.Exit(ExitHTTP)
	}
}

func runEcho(cfg *config.Config, data, file string, raw bool) {
	validateAPIKey(cfg)
	payload, err := resolvePayload(data, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitArgs)
	}
	if payload == nil {
		payload = []byte("{}")
	}
	c := client.New(cfg)
	f := formatter.New(raw)
	status, body, err := c.Post(apipath.Echo(), payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(status, body)
	if status < 200 || status >= 300 {
		os.Exit(ExitHTTP)
	}
}

func runPlusPing(cfg *config.Config) {
	validateAPIKey(cfg)
	c := client.New(cfg)
	_, body, err := c.Get(apipath.Ping(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	var data map[string]interface{}
	if json.Unmarshal([]byte(body), &data) == nil {
		if code, ok := data["code"].(float64); ok && int(code) == 200 {
			d := data["data"].(map[string]interface{})
			fmt.Println("网关连接正常")
			fmt.Printf("  状态: %v\n", d["status"])
			fmt.Printf("  时间: %v\n", d["ts"])
			return
		}
		m, _ := data["msg"].(string)
		fmt.Println("网关异常:", m)
		os.Exit(ExitHTTP)
	}
	fmt.Println(body)
}

func runPlusWhoami(cfg *config.Config) {
	validateAPIKey(cfg)
	c := client.New(cfg)
	_, body, err := c.Get(apipath.Whoami(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	var data map[string]interface{}
	if json.Unmarshal([]byte(body), &data) == nil {
		if code, ok := data["code"].(float64); ok && int(code) == 200 {
			d := data["data"].(map[string]interface{})
			fmt.Println("API Key 身份信息")
			fmt.Printf("  用户名:    %v\n", d["username"])
			fmt.Printf("  公司:      %v\n", d["companyName"])
			fmt.Printf("  CompanyId: %v\n", d["companyId"])
			fmt.Printf("  UserId:    %v\n", d["userId"])
			fmt.Printf("  认证方式:  %v\n", d["sourceType"])
			if rc, ok := d["roleCount"]; ok {
				fmt.Printf("  角色数:    %v\n", rc)
			}
			return
		}
		m, _ := data["msg"].(string)
		fmt.Println("查询失败:", m)
		os.Exit(ExitHTTP)
	}
	fmt.Println(body)
}

func runPlusModel(cfg *config.Config, data, file string) {
	validateAPIKey(cfg)
	payload, err := resolvePayload(data, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitArgs)
	}
	if payload == nil {
		fmt.Fprintln(os.Stderr, "[carbonstop] +model requires --data or --file")
		os.Exit(ExitArgs)
	}

	c := client.New(cfg)
	_, body, err := c.Post(apipath.AiModel(), payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}

	var result map[string]interface{}
	if json.Unmarshal([]byte(body), &result) == nil {
		if code, ok := result["code"].(float64); ok && int(code) == 200 {
			d := result["data"].(map[string]interface{})
			fmt.Println("建模成功")
			fmt.Printf("  核算ID: %v\n", d["accountId"])
			fmt.Printf("  产品ID: %v\n", d["productId"])
			fmt.Printf("  匹配率: %v\n", d["factorRate"])
			fmt.Println()

			// Stage summary
			stages := make(map[string]float64)
			emissions, _ := d["emissionList"].([]interface{})
			for _, item := range emissions {
				e := item.(map[string]interface{})
				stage := fmt.Sprint(e["stage_name"])
				total := parseFloat(fmt.Sprint(e["emissionTotal"]))
				stages[stage] += total
			}

			var grandTotal float64
			for _, v := range stages {
				grandTotal += v
			}

			fmt.Printf("%-12s %-18s %-8s\n", "阶段", "排放量(kgCO₂e)", "占比")
			fmt.Println(strings.Repeat("-", 40))
			for _, sname := range sortedKeys(stages) {
				stotal := stages[sname]
				pct := 0.0
				if grandTotal > 0 {
					pct = stotal / grandTotal * 100
				}
				fmt.Printf("%-12s %-18.6f %-8.2f%%\n", sname, stotal, pct)
			}
			fmt.Println(strings.Repeat("-", 40))
			fmt.Printf("%-12s %-18.6f 100.00%%\n", "总计", grandTotal)
			return
		}
		m, _ := result["msg"].(string)
		fmt.Println("建模失败:", m)
		os.Exit(ExitHTTP)
	}
	fmt.Println(body)
}

func getRaw(cmd *cobra.Command) bool {
	raw, _ := cmd.Flags().GetBool("raw")
	return raw
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func sortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
