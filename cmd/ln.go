package cmd

import (
	"fmt"
	"os"

	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/operation"

	"github.com/spf13/cobra"
)

var lnCmd = &cobra.Command{
	Use:     "ln {file(s) to link} {output template}",
	Short:   "Link files",
	Long:    `Link files with the power of templates`,
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"link", "mklink"},

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
			WithSoft(getBoolFlag(cmd, "soft", defaultOptions.Soft)).
			Initialize().
			Run(os.Args[1:])
	},
}

func init() {
	rootCmd.AddCommand(lnCmd)
	lnCmd.Flags().BoolP("force", "f", defaultOptions.Force, "Overwrite conflicts without prompt")
	lnCmd.Flags().BoolP("soft", "s", defaultOptions.Soft, "Create a soft link")
	lnCmd.Flags().Bool("simulate", defaultOptions.Simulate, "Simulate command and print outputs")
	lnCmd.Flags().String("sort", defaultOptions.Sort, "Sort files before running operations")
	lnCmd.Flags().BoolP("verbose", "v", defaultOptions.Verbose, "Verbose logging")
	lnCmd.Flags().BoolP("ignore-directories", "d", defaultOptions.IgnoreDirectories, "Do not move/rename directories")
	lnCmd.Flags().Bool("no-index", defaultOptions.NoIndex, "Do not automatically append an index when multiple operations result in the same file name")
	lnCmd.Flags().Bool("no-ext", defaultOptions.NoExt, "Do not automatically append the original file extension if one isn't supplied")
	lnCmd.Flags().Bool("no-mkdir", defaultOptions.NoMkdir, "Do not create any missing directories")
}
