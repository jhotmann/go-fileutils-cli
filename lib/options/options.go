package options

import (
	"github.com/jhotmann/go-fileutils-cli/lib/util"
	"github.com/spf13/cobra"
)

// Default Values
var (
	Force             = false
	Simulate          = false
	Sort              = "none"
	Verbose           = false
	IgnoreDirectories = false
	NoIndex           = false
	NoMove            = false
	NoExt             = false
	NoMkdir           = false
	Soft              = false
	AllowedSortValues = []string{"none", "alphabet", "reverse-alphabet", "date", "reverse-date", "size", "reverse-size"}
)

type CommonOptions struct {
	Force             bool
	Simulate          bool
	Sort              string
	Verbose           bool
	IgnoreDirectories bool
	NoIndex           bool
	NoExt             bool
	NoMkdir           bool
}

type MoveOptions struct {
	CommonOptions
	NoMove bool
}

type LinkOptions struct {
	CommonOptions
	Soft bool
}

func GetCommonOptions(cmd *cobra.Command) CommonOptions {
	var common = CommonOptions{
		Force:             util.GetBoolFlag(cmd, "force", Force),
		Simulate:          util.GetBoolFlag(cmd, "simulate", Simulate),
		Sort:              util.GetStringFlag(cmd, "sort", AllowedSortValues, Sort),
		Verbose:           util.GetBoolFlag(cmd, "verbose", Verbose),
		IgnoreDirectories: util.GetBoolFlag(cmd, "ignore-directories", IgnoreDirectories),
		NoIndex:           util.GetBoolFlag(cmd, "no-index", NoIndex),
		NoExt:             util.GetBoolFlag(cmd, "no-ext", NoExt),
		NoMkdir:           util.GetBoolFlag(cmd, "no-mkdir", NoMkdir),
	}
	return common
}

func GetMoveOptions(cmd *cobra.Command) MoveOptions {
	var opts = MoveOptions{
		CommonOptions: GetCommonOptions(cmd),
		NoMove:        util.GetBoolFlag(cmd, "no-move", NoMove),
	}
	return opts
}

func GetLinkOptions(cmd *cobra.Command) LinkOptions {
	var opts = LinkOptions{
		CommonOptions: GetCommonOptions(cmd),
		Soft:          util.GetBoolFlag(cmd, "soft", Soft),
	}
	return opts
}
