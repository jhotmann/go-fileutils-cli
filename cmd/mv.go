/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/lib/operation"
	"github.com/jhotmann/go-fileutils-cli/lib/options"

	"github.com/spf13/cobra"
)

// mvCmd represents the mv command
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
			panic(err)
		}
		// all other non-flag arguments are input files
		inputFiles := args[0 : len(args)-1]
		// parse options into our own struct
		opts := options.GetMoveOptions(cmd)
		// if rename alias used, set --no-move automatically
		if cmd.CalledAs() == "rename" {
			opts.NoMove = true
		}
		// create a list of operations for all input files
		operations := operation.FilesToOperationsList("move", inputFiles, outputTemplate)
		// filter out directories if --ignore-directories option passed
		if opts.IgnoreDirectories {
			operations = operations.RemoveDirectories()
		}
		// filter out repeat inputs (only applies to moves), sort, and convert output from template to string to PathObj
		operations = operations.RemoveDuplicateInputs().Sort(opts.Sort).RenderTemplates()
		if !opts.NoExt {
			operations = operations.PopulateBlankExtensions()
		}
		if opts.NoMove { // keep output directory same as input
			operations = operations.NoMove()
		}
		if !opts.Force { // don't care about conflicts
			operations = operations.FindConflicts()
		}
		if !opts.NoIndex { // auto-index conflicting outputs
			operations = operations.AddIndex()
		}
		operations.Run(opts.CommonOptions)
	},
}

func init() {
	rootCmd.AddCommand(mvCmd)
	mvCmd.Flags().BoolP("force", "f", options.Force, "Overwrite conflicts without prompt")
	mvCmd.Flags().BoolP("simulate", "s", options.Simulate, "Simulate command and print outputs")
	mvCmd.Flags().String("sort", options.Sort, "Sort files before running operations")
	mvCmd.Flags().BoolP("verbose", "v", options.Verbose, "Verbose logging")
	mvCmd.Flags().BoolP("ignore-directories", "d", options.IgnoreDirectories, "Do not move/rename directories")
	mvCmd.Flags().Bool("no-index", options.NoIndex, "Do not automatically append an index when multiple operations result in the same file name")
	mvCmd.Flags().Bool("no-move", options.NoMove, "Do not move files to a different directory")
	mvCmd.Flags().Bool("no-ext", options.NoExt, "Do not automatically append the original file extension if one isn't supplied")
	mvCmd.Flags().Bool("no-mkdir", options.NoMkdir, "Do not create any missing directories")
}
