package cmd

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/operation"
	"github.com/jhotmann/go-fileutils-cli/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func FilesToOperationsList(opType string, files []string, outputTemplate *pongo2.Template) operation.OperationList {
	operations := operation.OperationList{}

	var fileWg sync.WaitGroup
	fileWg.Add(len(files))

	for _, f := range files {
		go func(f string) {
			matches, err := filepath.Glob(f)
			if err != nil {
				panic(err)
			}
			if len(matches) == 0 {
				pterm.Warning.Printfln("%s does not match any existing files", f)
			}

			var opsWg sync.WaitGroup
			operationChannel := make(chan operation.Operation, 10)
			opsWg.Add(len(matches))

			for _, match := range matches {
				go func(match string) {
					var op operation.Operation
					op.Type = opType
					op.Input = operation.GetPathObj(match)
					op.OutputTemplate = outputTemplate
					stats, err := os.Stat(match)
					if err == nil {
						op.Stats = stats
					}
					op.Options = operation.DefaultOptions
					op.Skip = false
					//operations = append(operations, op)
					operationChannel <- op
				}(match)
			}
			go func() {
				for op := range operationChannel {
					operations = append(operations, op)
					opsWg.Done()
				}
			}()
			opsWg.Wait()
			fileWg.Done()
		}(f)
	}

	fileWg.Wait()
	return operations
}

var defaultOptions = operation.DefaultOptions

func getBoolFlag(cmd *cobra.Command, name string, defaultValue bool) bool {
	ret, err := cmd.Flags().GetBool(name)
	if err != nil {
		return defaultValue
	}
	return ret
}

func getStringFlag(cmd *cobra.Command, name string, allowedValues []string, defaultValue string) string {
	ret, err := cmd.Flags().GetString(name)
	if err != nil {
		return defaultValue
	}
	if allowedValues == nil {
		return ret
	}
	if util.IndexOf(ret, allowedValues) > -1 {
		return ret
	}
	return defaultValue
}

func getIntFlag(cmd *cobra.Command, name string, allowedValues []int, defaultValue int) int {
	ret, err := cmd.Flags().GetInt(name)
	if err != nil {
		return defaultValue
	}
	if allowedValues == nil {
		return ret
	}
	if util.IndexOfInt(ret, allowedValues) > -1 {
		return ret
	}
	return defaultValue
}
