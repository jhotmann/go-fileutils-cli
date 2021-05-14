package filters

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/flosch/pongo2/v4"
	"github.com/jhotmann/go-fileutils-cli/lib/util"
)

func init() {
	pongo2.ReplaceFilter("date", DateFilter)
	pongo2.ReplaceFilter("time", DateFilter)
	pongo2.RegisterFilter("snake", snakeFilter)
}

func snakeFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(util.ToSnakeCase(in.String())), nil
}

func DateFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, isTime := in.Interface().(time.Time)
	if !isTime {
		return nil, &pongo2.Error{
			Sender:    "filter:date",
			OrigError: errors.New("filter input argument must be of type 'time.Time'"),
		}
	}
	fmt.Println(unicodeToGoDateFormat(param.String()))
	return pongo2.AsValue(t.Format(unicodeToGoDateFormat(param.String()))), nil
}

func unicodeToGoDateFormat(s string) string {
	mRegex := regexp2.MustCompile(`(?<![AP])M`, 0)
	eeeRegex := regexp2.MustCompile(`E{1,3}`, 0)
	aRegex := regexp2.MustCompile(`(?<!J)a{1,4}(?!n)`, 0)

	ret := strings.ReplaceAll(s, "yyyy", "2006")
	ret = strings.ReplaceAll(ret, "yyy", "006")
	ret = strings.ReplaceAll(ret, "yy", "06")
	ret = strings.ReplaceAll(ret, "MMMM", "January")
	ret = strings.ReplaceAll(ret, "MMM", "Jan")
	ret = strings.ReplaceAll(ret, "MM", "01")
	ret = strings.ReplaceAll(ret, "dd", "02")
	ret = strings.ReplaceAll(ret, "d", "2")
	ret = strings.ReplaceAll(ret, "EEEE", "Monday")
	ret, _ = eeeRegex.Replace(ret, "Mon", -1, -1)
	ret = strings.ReplaceAll(ret, "HH", "15")
	ret = strings.ReplaceAll(ret, "H", "15")
	ret = strings.ReplaceAll(ret, "hh", "03")
	ret = strings.ReplaceAll(ret, "h", "3")
	ret = strings.ReplaceAll(ret, "mm", "04")
	ret = strings.ReplaceAll(ret, "m", "4")
	ret = strings.ReplaceAll(ret, "ss", "05")
	ret = strings.ReplaceAll(ret, "s", "5")
	ret, _ = mRegex.Replace(ret, "1", -1, -1)
	ret, _ = aRegex.Replace(ret, "PM", -1, -1)
	return ret
}
