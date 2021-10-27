package cmd

import (
	"errors"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/jhotmann/go-fileutils-cli/db"
	"github.com/jhotmann/go-fileutils-cli/util"
)

var favoriteCmd = &cobra.Command{
	Use:     "favorite",
	Short:   "View/run/edit favorited commands",
	Long:    `View/run/edit favorited commands`,
	Aliases: []string{"f", "fav", "favourite"},

	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(getStringFlag(cmd, "name", nil, ""))
		id := getIntFlag(cmd, "id", nil, 5)
		if name != "" {
			err := runByName(name)
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

func runByName(name string) error {
	favorite := db.GetFavoriteByName(name)
	if favorite.Id == 0 {
		return errors.New("Favorite not found")
	}
	return util.RunCommand(favorite.Command, util.GetWorkingDir())
}
