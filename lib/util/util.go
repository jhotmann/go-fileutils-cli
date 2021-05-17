package util

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/spf13/cobra"
)

func GetBoolFlag(cmd *cobra.Command, name string, defaultValue bool) bool {
	ret, err := cmd.Flags().GetBool(name)
	if err != nil {
		return defaultValue
	}
	return ret
}

func GetStringFlag(cmd *cobra.Command, name string, allowedValues []string, defaultValue string) string {
	ret, err := cmd.Flags().GetString(name)
	if err != nil {
		return defaultValue
	}
	if allowedValues == nil {
		return ret
	}
	if IndexOf(ret, allowedValues) > -1 {
		return ret
	}
	return defaultValue
}

func IndexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp2.MustCompile("(.)([A-Z][a-z]+)", 0)
	var matchAllCap = regexp2.MustCompile("([a-z0-9])([A-Z])", 0)
	snake, _ := matchFirstCap.Replace(str, "${1}_${2}", -1, -1)
	snake, _ = matchAllCap.Replace(snake, "${1}_${2}", -1, -1)
	return strings.ToLower(snake)
}

func ZeroPad(num int, total int) string {
	desiredLength := len(fmt.Sprintf("%d", total))
	return fmt.Sprintf("%0*d", desiredLength, num)
}
