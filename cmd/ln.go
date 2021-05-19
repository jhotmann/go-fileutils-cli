package cmd

import (
	"fmt"
	"os"

	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/lib/operation"
	"github.com/jhotmann/go-fileutils-cli/lib/options"

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
		// parse options into our own struct
		opts := options.GetLinkOptions(cmd)
		// create a list of operations for all input files
		var operations operation.OperationList
		if opts.Soft {
			operations = operation.FilesToOperationsList("link-soft", inputFiles, outputTemplate)
		} else {
			operations = operation.FilesToOperationsList("link-hard", inputFiles, outputTemplate)
		}
		// filter out directories if --ignore-directories option passed
		if opts.IgnoreDirectories {
			operations = operations.RemoveDirectories()
		}
		// filter out repeat inputs (only applies to moves), sort, and convert output from template to string to PathObj
		operations = operations.RemoveDuplicateInputs().Sort(opts.Sort).RenderTemplates()
		if !opts.NoExt {
			operations = operations.PopulateBlankExtensions()
		}
		if !opts.Force { // don't care about conflicts
			operations = operations.FindConflicts()
		}
		if !opts.NoIndex { // auto-index conflicting outputs
			operations = operations.AddIndex()
		}
		operations.Run(os.Args[1:], opts.CommonOptions)
	},
}

func init() {
	rootCmd.AddCommand(lnCmd)
	lnCmd.Flags().BoolP("force", "f", options.Force, "Overwrite conflicts without prompt")
	lnCmd.Flags().BoolP("soft", "s", options.Soft, "Create a soft link")
	lnCmd.Flags().Bool("simulate", options.Simulate, "Simulate command and print outputs")
	lnCmd.Flags().String("sort", options.Sort, "Sort files before running operations")
	lnCmd.Flags().BoolP("verbose", "v", options.Verbose, "Verbose logging")
	lnCmd.Flags().BoolP("ignore-directories", "d", options.IgnoreDirectories, "Do not move/rename directories")
	lnCmd.Flags().Bool("no-index", options.NoIndex, "Do not automatically append an index when multiple operations result in the same file name")
	lnCmd.Flags().Bool("no-move", options.NoMove, "Do not move files to a different directory")
	lnCmd.Flags().Bool("no-ext", options.NoExt, "Do not automatically append the original file extension if one isn't supplied")
	lnCmd.Flags().Bool("no-mkdir", options.NoMkdir, "Do not create any missing directories")
}
