package cmd

import (
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/jhotmann/go-fileutils-cli/lib/favorites"
	"github.com/jhotmann/go-fileutils-cli/lib/util"
)

var favoriteCmd = &cobra.Command{
	Use:     "favorite",
	Short:   "View/run/edit favorited commands",
	Long:    `View/run/edit favorited commands`,
	Aliases: []string{"f", "fav", "favourite"},

	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(util.GetStringFlag(cmd, "name", nil, ""))
		id := util.GetIntFlag(cmd, "id", nil, 5)
		if name != "" {
			err := favorites.RunByName(name)
			if err != nil {
				pterm.Error.WithShowLineNumber(false).Printfln("Favorite with name %s not found", name)
			}
		} else if id != 0 {
			//run by id
		} else {
			//display favorites
		}
	},
}

func init() {
	rootCmd.AddCommand(favoriteCmd)
	favoriteCmd.Flags().StringP("name", "n", "", "Run a favorited command by name")
	favoriteCmd.Flags().IntP("id", "i", 0, "Run a favorited command by ID")
}
