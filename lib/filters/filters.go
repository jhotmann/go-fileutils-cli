package filters

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/flosch/pongo2/v4"
	"github.com/iancoleman/strcase"
	"github.com/jhotmann/go-fileutils-cli/lib/util"
)

func init() {
	pongo2.ReplaceFilter("date", DateFilter)
	pongo2.ReplaceFilter("time", DateFilter)
	pongo2.ReplaceFilter("title", titleFilter)
	pongo2.RegisterFilter("pascal", titleFilter)
	pongo2.RegisterFilter("snake", snakeFilter)
	pongo2.RegisterFilter("camel", camelFilter)
	pongo2.RegisterFilter("kebab", kebabFilter)
	pongo2.RegisterFilter("replace", replaceFilter)
	pongo2.RegisterFilter("regexReplace", regexReplaceFilter)
	pongo2.RegisterFilter("with", withFilter)
	pongo2.RegisterFilter("match", matchFilter)
	pongo2.RegisterFilter("index", indexFilter)
	pongo2.RegisterFilter("pad", padFilter)
}

func titleFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(strcase.ToCamel(in.String())), nil
}

func snakeFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(strcase.ToSnake(in.String())), nil
}

func camelFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(strcase.ToLowerCamel(in.String())), nil
}

func kebabFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(strcase.ToKebab(in.String())), nil
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

func replaceFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(strings.ReplaceAll(in.String(), param.String(), "--REPLACEME--")), nil
}

func regexReplaceFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return pongo2.AsValue(""), nil
	}
	re := regexp2.MustCompile(param.String(), 0)
	out, err := re.Replace(in.String(), "--REPLACEME--", -1, -1)
	if err != nil {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(out), nil
}

func withFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(strings.ReplaceAll(in.String(), "--REPLACEME--", param.String())), nil
}

func matchFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return pongo2.AsValue(""), nil
	}
	re := regexp2.MustCompile(param.String(), 0)
	var matches []string
	m, _ := re.FindStringMatch(in.String())
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	switch len(matches) {
	case 0:
		return pongo2.AsValue(""), nil
	case 1:
		return pongo2.AsValue(matches[0]), nil
	default:
		return pongo2.AsValue(matches), nil
	}
}

func indexFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.CanSlice() || in.Len() <= 0 {
		return in, nil
	}
	i, err := strconv.ParseInt(param.String(), 10, 8)
	if err != nil {
		return in, nil
	}
	return in.Index(int(i)), nil
}

func padFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(util.ZeroPadString(in.String(), param.String())), nil
}
