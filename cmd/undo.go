package cmd

import (
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/jhotmann/go-fileutils-cli/db"
	"github.com/jhotmann/go-fileutils-cli/util"
	"github.com/pterm/pterm"
)

var undoCmd = &cobra.Command{
	Use:     "undo",
	Short:   "Undo the last undoable batch",
	Long:    `Undo the last undoable batch that hasn't already been undone`,
	Aliases: []string{"h"},

	Run: func(cmd *cobra.Command, args []string) {
		batch, err := db.GetLastNonUndone()
		if err != nil || batch.Id == 0 {
			pterm.Error.WithShowLineNumber(false).Println("No undoable batches found")
			os.Exit(1)
		}
		operations, err := db.GetOperationsForBatch(batch.Id)
		if err != nil {
			batch.Close()
			panic(err)
		}
		//util.ClearTerm()
		pterm.FgLightBlue.Printfln("Command: fu %s", batch.CommandString)
		pterm.Println()
		pterm.DefaultTable.WithHasHeader().WithData(operations.ToTableData(batch.WorkingDir)).Render()
		pterm.Println()
		prompt := promptui.Prompt{
			Label:     "Are you sure you want to undo this batch?",
			IsConfirm: true,
		}
		result, _ := prompt.Run()
		if util.IndexOf(strings.ToLower(result), []string{"y", "yes", "true", "1"}) > -1 {
			err = batch.Undo()
			batch.Close()
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
}
