package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dlclark/regexp2"
	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/jhotmann/go-fileutils-cli/db"
	"github.com/jhotmann/go-fileutils-cli/util"
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
		oldestFirst := getBoolFlag(cmd, "oldest-first", false)
		itemsPerPage := getIntFlag(cmd, "items-per-page", nil, 5)
		printHistory(1, itemsPerPage, oldestFirst)
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().BoolP("oldest-first", "o", false, "Order history from oldest to newest")
	historyCmd.Flags().IntP("items-per-page", "i", 5, "How many batches to display on each page")
}

func getBatches() {
	batches, err = db.GetBatches()
	if err != nil {
		batches.Close()
		panic(err)
	}
}

func printHistory(page int, countPerPage int, oldestFirst bool) {
	if batches == nil {
		getBatches()
	}
	util.ClearTerm()
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("File", pterm.NewStyle(pterm.FgLightBlue)),
		pterm.NewLettersFromStringWithStyle("Utils", pterm.NewStyle(pterm.FgLightCyan))).
		Render()
	if !oldestFirst {
		batches = batches.Reverse()
	}
	subset, pageBefore, pageAfter := batches.GetPage(page, countPerPage)
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
		batches.Close()
		os.Exit(0)
	}
	if result == "n" {
		printHistory(page+1, countPerPage, true)
	} else if result == "p" {
		printHistory(page-1, countPerPage, true)
	} else {
		intResult, _ := strconv.ParseInt(result, 10, 0)
		for _, b := range subset {
			if b.Id == int(intResult) {
				printBatch(b, page, countPerPage)
				break
			}
		}
	}
}

func printBatch(batch db.Batch, returnPage int, countPerPage int) {
	if batches == nil {
		getBatches()
	}
	operations, err := db.GetOperationsForBatch(batch.Id)
	if err != nil {
		batches.Close()
		panic(err)
	}
	util.ClearTerm()
	pterm.FgLightBlue.Printfln("Command: fu %s", batch.CommandString)
	pterm.Println()
	pterm.DefaultTable.WithHasHeader().WithData(operations.ToTableData(batch.WorkingDir)).Render()
	pterm.Println()
	pterm.Println("Input Options:")
	pterm.Println("  [R] Re-run Command")
	pterm.Println("  [C] Copy Command")
	pterm.Println("  [U] Undo Command")
	pterm.Println("  [#'s] Undo Operation(s) (IDs separated by comma)")
	pterm.Println("  [F] Add to Favorites")
	pterm.Println("  [B] Back")
	prompt := promptui.Prompt{
		Label: "Input",
		Validate: func(input string) error {
			if util.IndexOf(strings.ToLower(input), []string{"r", "c", "u", "f", "b"}) > -1 {
				return nil
			}
			_, err := matchOperationsById(operations, input)
			return err
		},
	}
	result, err := prompt.Run()
	if err != nil {
		batches.Close()
		os.Exit(0)
	}
	switch strings.ToLower(result) {
	case "r": // re-run
		batch.Close()
		cmd := exec.Cmd{
			Path:   os.Args[0],
			Args:   append([]string{os.Args[0]}, batch.Command...),
			Dir:    batch.WorkingDir,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
		err = cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
		break
	case "c": // copy
		clipboard.WriteAll("fu " + batch.CommandString)
		batch.Close()
		break
	case "u": // undo
		batch.Undo()
		batch.Close()
		break
	case "f": // favorite
		// Add to favorites TODO
		batch.Close()
		break
	case "b": // back
		printHistory(returnPage, countPerPage, true)
		break
	default: // undo selected operations
		subset, _ := matchOperationsById(operations, result)
		subset.Undo(batch.CommandType, batch.WorkingDir)
		batch.Close()
	}
}

func matchOperationsById(operations db.OperationList, ids string) (db.OperationList, error) {
	uniqueIds := []int{}
	ret := db.OperationList{}
	var err error
	r := regexp2.MustCompile("^\\d+(,\\d+)*$", 0)
	matches, err := r.MatchString(ids)
	if err != nil {
		return ret, err
	}
	if !matches {
		return ret, errors.New("Invalid input")
	} else {
		allIdsFound := true
		numbers := strings.Split(ids, ",")
		for _, num := range numbers {
			idFound := false
			parsed, _ := strconv.ParseInt(num, 10, 0)
			for _, op := range operations {
				if int(parsed) == op.Id {
					idFound = true
					if util.IndexOfInt(int(parsed), uniqueIds) == -1 {
						ret = append(ret, op)
						uniqueIds = append(uniqueIds, int(parsed))
					}
					break
				}
			}
			if !idFound {
				allIdsFound = false
			}
		}
		if !allIdsFound {
			err = errors.New("Not all operation IDs valid")
		}
	}
	return ret, err
}
