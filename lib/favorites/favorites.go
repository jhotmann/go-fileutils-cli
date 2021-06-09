package favorites

import (
	"errors"

	"github.com/jhotmann/go-fileutils-cli/lib/db"
	"github.com/jhotmann/go-fileutils-cli/lib/util"
)

func RunByName(name string) error {
	favorite := db.GetFavoriteByName(name)
	if favorite.Id == 0 {
		return errors.New("Favorite not found")
	}
	return util.RunCommand(favorite.Command, util.GetWorkingDir())
}
