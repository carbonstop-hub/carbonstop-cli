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

func NewAccountsCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Account list",
		Run: func(cmd *cobra.Command, args []string) {
			productID, _ := cmd.Flags().GetInt("product-id")
			pageNum, _ := cmd.Flags().GetInt("page-num")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			acctStatus, _ := cmd.Flags().GetInt("account-status")
			runAccounts(getCfg(), productID, pageNum, pageSize, acctStatus, getRaw(cmd))
		},
	}
	cmd.Flags().Int("product-id", 0, "Product ID (required)")
	cmd.MarkFlagRequired("product-id")
	cmd.Flags().Int("page-num", 1, "Page number (default 1)")
	cmd.Flags().Int("page-size", 10, "Page size (default 10)")
	cmd.Flags().Int("account-status", 3, "Account status (default 3)")
	return cmd
}

func NewAccountViewCmd(getCfg ConfigGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-view",
		Short: "Account detail",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetInt("id")
			groupType, _ := cmd.Flags().GetInt("group-type")
			lang, _ := cmd.Flags().GetString("lang")
			runAccountView(getCfg(), id, groupType, lang, getRaw(cmd))
		},
	}
	cmd.Flags().Int("id", 0, "Account ID (required)")
	cmd.MarkFlagRequired("id")
	cmd.Flags().Int("group-type", 0, "Group type (default 0)")
	cmd.Flags().String("lang", "zh", "Language (default zh)")
	return cmd
}

func runAccounts(cfg *config.Config, productID, pageNum, pageSize, acctStatus int, raw bool) {
	validateAPIKey(cfg)
	c := client.New(cfg)
	f := formatter.New(raw)
	query := map[string]string{
		"productId":     strconv.Itoa(productID),
		"pageNum":       strconv.Itoa(pageNum),
		"pageSize":      strconv.Itoa(pageSize),
		"accountStatus": strconv.Itoa(acctStatus),
	}
	st, body, err := c.Get(apipath.Accounts(), query)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(st, body)
	if st < 200 || st >= 300 {
		os.Exit(ExitHTTP)
	}
}

func runAccountView(cfg *config.Config, id, groupType int, lang string, raw bool) {
	validateAPIKey(cfg)
	c := client.New(cfg)
	f := formatter.New(raw)
	query := map[string]string{
		"id":        strconv.Itoa(id),
		"groupType": strconv.Itoa(groupType),
		"lang":      lang,
	}
	st, body, err := c.Get(apipath.AccountView(), query)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[carbonstop]", err)
		os.Exit(ExitTransport)
	}
	f.Print(st, body)
	if st < 200 || st >= 300 {
		os.Exit(ExitHTTP)
	}
}
