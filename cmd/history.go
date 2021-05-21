package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jhotmann/go-fileutils-cli/lib/db"
	"github.com/jhotmann/go-fileutils-cli/lib/history"
	"github.com/jhotmann/go-fileutils-cli/lib/util"
)

var (
	batches db.BatchList
	err     error
)

var historyCmd = &cobra.Command{
	Use:     "history",
	Short:   "View command history",
	Long:    `View command history and re-run, undo, or favorite commands`,
	Aliases: []string{"h"},

	Run: func(cmd *cobra.Command, args []string) {
		oldestFirst := util.GetBoolFlag(cmd, "oldest-first", false)
		itemsPerPage := util.GetIntFlag(cmd, "items-per-page", nil, 5)
		history.PrintHistory(1, itemsPerPage, oldestFirst)
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().BoolP("oldest-first", "o", false, "Order history from oldest to newest")
	historyCmd.Flags().IntP("items-per-page", "i", 5, "How many batches to display on each page")
}
