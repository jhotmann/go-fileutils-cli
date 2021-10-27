package cmd

import (
	"fmt"
	"os"

	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/operation"

	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:     "cp {file(s) to copy} {output template}",
	Short:   "Copy files",
	Long:    `Copy files with the power of templates`,
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"copy"},

	Run: func(cmd *cobra.Command, args []string) {
		// output is the last non-flag argument
		outputTemplate, err := pongo2.FromString(args[len(args)-1])
		if err != nil {
			fmt.Println("Invalid Output: ", err.Error())
			os.Exit(1)
		}
		// all other non-flag arguments are input files
		inputFiles := args[0 : len(args)-1]
		// run it
		FilesToOperationsList(operation.OperationType.Cp, inputFiles, outputTemplate).
			WithForce(getBoolFlag(cmd, "force", defaultOptions.Force)).
			WithSimulate(getBoolFlag(cmd, "simulate", defaultOptions.Simulate)).
			WithVerbose(getBoolFlag(cmd, "verbose", defaultOptions.Verbose)).
			WithIgnoreDirectories(getBoolFlag(cmd, "ignore-directories", defaultOptions.Verbose)).
			WithNoIndex(getBoolFlag(cmd, "no-index", defaultOptions.NoIndex)).
			WithNoExt(getBoolFlag(cmd, "no-ext", defaultOptions.NoExt)).
			WithNoMkdir(getBoolFlag(cmd, "no-mkdir", defaultOptions.NoMkdir)).
			WithSort(getStringFlag(cmd, "sort", operation.AllowedSortValues, defaultOptions.Sort)).
			Initialize().
			Run(os.Args[1:])
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
	cpCmd.Flags().BoolP("force", "f", defaultOptions.Force, "Overwrite conflicts without prompt")
	cpCmd.Flags().BoolP("simulate", "s", defaultOptions.Simulate, "Simulate command and print outputs")
	cpCmd.Flags().String("sort", defaultOptions.Sort, "Sort files before running operations")
	cpCmd.Flags().BoolP("verbose", "v", defaultOptions.Verbose, "Verbose logging")
	cpCmd.Flags().BoolP("ignore-directories", "d", defaultOptions.IgnoreDirectories, "Do not move/rename directories")
	cpCmd.Flags().Bool("no-index", defaultOptions.NoIndex, "Do not automatically append an index when multiple operations result in the same file name")
	cpCmd.Flags().Bool("no-ext", defaultOptions.NoExt, "Do not automatically append the original file extension if one isn't supplied")
	cpCmd.Flags().Bool("no-mkdir", defaultOptions.NoMkdir, "Do not create any missing directories")
}
