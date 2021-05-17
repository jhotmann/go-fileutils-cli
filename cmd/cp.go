package cmd

import (
	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/lib/operation"
	"github.com/jhotmann/go-fileutils-cli/lib/options"

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
			panic(err)
		}
		// all other non-flag arguments are input files
		inputFiles := args[0 : len(args)-1]
		// parse options into our own struct
		opts := options.GetCommonOptions(cmd)
		// create a list of operations for all input files
		operations := operation.FilesToOperationsList("copy", inputFiles, outputTemplate)
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
		operations.Run(opts)
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
	cpCmd.Flags().BoolP("force", "f", options.Force, "Overwrite conflicts without prompt")
	cpCmd.Flags().BoolP("simulate", "s", options.Simulate, "Simulate command and print outputs")
	cpCmd.Flags().String("sort", options.Sort, "Sort files before running operations")
	cpCmd.Flags().BoolP("verbose", "v", options.Verbose, "Verbose logging")
	cpCmd.Flags().BoolP("ignore-directories", "d", options.IgnoreDirectories, "Do not move/rename directories")
	cpCmd.Flags().Bool("no-index", options.NoIndex, "Do not automatically append an index when multiple operations result in the same file name")
	cpCmd.Flags().Bool("no-move", options.NoMove, "Do not move files to a different directory")
	cpCmd.Flags().Bool("no-ext", options.NoExt, "Do not automatically append the original file extension if one isn't supplied")
	cpCmd.Flags().Bool("no-mkdir", options.NoMkdir, "Do not create any missing directories")
}
