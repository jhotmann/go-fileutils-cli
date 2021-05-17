package util

import (
	"fmt"

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

func ZeroPad(num int, total int) string {
	desiredLength := len(fmt.Sprintf("%d", total))
	return fmt.Sprintf("%0*d", desiredLength, num)
}

func ZeroPadString(in string, total string) string {
	desiredLength := len(total)
	return fmt.Sprintf("%0*s", desiredLength, in)
}
