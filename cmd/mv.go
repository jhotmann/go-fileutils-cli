package cmd

import (
	"fmt"
	"os"

	"github.com/flosch/pongo2/v4"

	"github.com/spf13/cobra"
)

var mvCmd = &cobra.Command{
	Use:     "mv {file(s) to move} {output template}",
	Short:   "Move/Rename files",
	Long:    `Move/Rename files with the power of templates`,
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"move", "rename"},

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
		//opts := options.GetMoveOptions(cmd)
		// if rename alias used, set --no-move automatically
		if cmd.CalledAs() == "rename" {
			//opts.NoMove = true
		}
		// create a list of operations for all input files
		operations := FilesToOperationsList("move", inputFiles, outputTemplate)
		if len(operations) == 0 {
			fmt.Println("Error: no operations can be created from the input(s) specified")
			os.Exit(1)
		}
		// filter out directories if --ignore-directories option passed
		// if opts.IgnoreDirectories {
		// 	operations = operations.RemoveDirectories()
		// }
		// // filter out repeat inputs (only applies to moves), sort, and convert output from template to string to PathObj
		// operations = operations.RemoveDuplicateInputs().Sort(opts.Sort).RenderTemplates()
		// if !opts.NoExt {
		// 	operations = operations.PopulateBlankExtensions()
		// }
		// if opts.NoMove { // keep output directory same as input
		// 	operations = operations.NoMove()
		// }
		// if !opts.Force { // don't care about conflicts
		// 	operations = operations.FindConflicts()
		// }
		// if !opts.NoIndex { // auto-index conflicting outputs
		// 	operations = operations.AddIndex()
		// }
		// operations.Run(os.Args[1:], opts.CommonOptions)
	},
}

func init() {
	rootCmd.AddCommand(mvCmd)
	mvCmd.Flags().BoolP("force", "f", defaultOptions.Force, "Overwrite conflicts without prompt")
	mvCmd.Flags().BoolP("simulate", "s", defaultOptions.Simulate, "Simulate command and print outputs")
	mvCmd.Flags().String("sort", defaultOptions.Sort, "Sort files before running operations")
	mvCmd.Flags().BoolP("verbose", "v", defaultOptions.Verbose, "Verbose logging")
	mvCmd.Flags().BoolP("ignore-directories", "d", defaultOptions.IgnoreDirectories, "Do not move/rename directories")
	mvCmd.Flags().Bool("no-index", defaultOptions.NoIndex, "Do not automatically append an index when multiple operations result in the same file name")
	mvCmd.Flags().Bool("no-move", defaultOptions.NoMove, "Do not move files to a different directory")
	mvCmd.Flags().Bool("no-ext", defaultOptions.NoExt, "Do not automatically append the original file extension if one isn't supplied")
	mvCmd.Flags().Bool("no-mkdir", defaultOptions.NoMkdir, "Do not create any missing directories")
}
