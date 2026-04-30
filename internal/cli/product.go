package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/carbonstop/carbonstop-cli/internal/apipath"
	"github.com/carbonstop/carbonstop-cli/internal/client"
	"github.com/carbonstop/carbonstop-cli/internal/config"
	"github.com/carbonstop/carbonstop-cli/internal/formatter"
	"github.com/spf13/cobra"
)

func NewProductsCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "products",
		Short: "Product list",
		Run: func(cmd *cobra.Command, args []string) {
			pageNum, _ := cmd.Flags().GetInt("page-num")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			status, _ := cmd.Flags().GetInt("status")
			search, _ := cmd.Flags().GetString("search")
			runProducts(getCfg(), pageNum, pageSize, status, search, getRaw(cmd))
		},
	}
	cmd.Flags().Int("page-num", 1, "Page number (default 1)")
	cmd.Flags().Int("page-size", 12, "Page size (default 12)")
	cmd.Flags().Int("status", 2, "Status filter (default 2)")
	cmd.Flags().String("search", "", "Search product name")
	return cmd
}

func NewProductInfoCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product-info",
		Short: "Product detail",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetInt("id")
			if id <= 0 {
				fmt.Fprintln(os.Stderr, "[carbonstop] --id is required")
				os.Exit(ExitArgs)
			}
			runProductInfo(getCfg(), id, getRaw(cmd))
		},
	}
	cmd.Flags().Int("id", 0, "Product ID (required)")
	cmd.MarkFlagRequired("id")
	return cmd
}

func runProducts(cfg *config.Config, pageNum, pageSize, status int, search string, raw bool) {
	validateAPIKey(cfg)
	c := client.New(cfg)
	f := formatter.New(raw)
	query := map[string]string{
		"pageNum":  strconv.Itoa(pageNum),
		"pageSize": strconv.Itoa(pageSize),
		"status":   strconv.Itoa(status),
	}
	if search != "" {
		query["search"] = search
	}
	st, body, err := c.Get(apipath.Products(), query)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(st, body)
	if st < 200 || st >= 300 {
		os.Exit(ExitHTTP)
	}
}

func runProductInfo(cfg *config.Config, id int, raw bool) {
	validateAPIKey(cfg)
	c := client.New(cfg)
	f := formatter.New(raw)
	st, body, err := c.Get(apipath.ProductInfo(), map[string]string{"id": strconv.Itoa(id)})
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(st, body)
	if st < 200 || st >= 300 {
		os.Exit(ExitHTTP)
	}
}
