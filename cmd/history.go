package cmd

import (
	"errors"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/jhotmann/go-fileutils-cli/lib/db"
	"github.com/jhotmann/go-fileutils-cli/lib/util"
	"github.com/pterm/pterm"
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
		batches, err = db.GetBatches()
		if err != nil {
			panic(err)
		}
		printHistory(1)
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}

func printHistory(page int) {
	util.ClearTerm()
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("File", pterm.NewStyle(pterm.FgLightBlue)),
		pterm.NewLettersFromStringWithStyle("Utils", pterm.NewStyle(pterm.FgLightCyan))).
		Render()
	subset, pageBefore, pageAfter := batches.GetPage(page, 5)
	pterm.DefaultTable.WithHasHeader().WithData(subset.ToTableData()).Render()
	pterm.Println()
	pterm.Println("Input Options:")
	pterm.Println("  [#] Batch ID")
	if pageBefore {
		pterm.Println("  [P] for previous page")
	}
	if pageAfter {
		pterm.Println("  [N] for next page")
	}
	prompt := promptui.Prompt{
		Label: "Input",
		Validate: func(input string) error {
			if pageBefore && strings.ToLower(input) == "p" {
				return nil
			}
			if pageAfter && strings.ToLower(input) == "n" {
				return nil
			}
			number, err := strconv.ParseInt(input, 10, 0)
			if err != nil {
				return err
			}
			found := false
			for _, b := range subset {
				if b.Id == int(number) {
					found = true
					break
				}
			}
			if !found {
				return errors.New("Batch ID invalid")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	if result == "n" {
		printHistory(page + 1)
	} else if result == "p" {
		printHistory(page - 1)
	} else {
		intResult, _ := strconv.ParseInt(result, 10, 0)
		for _, b := range subset {
			if b.Id == int(intResult) {
				printBatch(b)
				break
			}
		}
	}
}

func printBatch(batch db.Batch) {
	operations, err := db.GetOperationsForBatch(batch.Id)
	if err != nil {
		panic(err)
	}
	util.ClearTerm()
	pterm.FgLightBlue.Printfln("Command: fu %s", batch.CommandString)
	pterm.Println()
	pterm.DefaultTable.WithHasHeader().WithData(operations.ToTableData(batch.WorkingDir)).Render()
}
