package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"

	"github.com/spf13/cobra"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		fmt.Println("\033[H\033[2J")
		// cmd := exec.Command("clear") //Linux example, its tested
		// cmd.Stdout = os.Stdout
		// cmd.Run()
	}
	clear["darwin"] = func() {
		fmt.Println("\033[H\033[2J")
	}
	clear["freebsd"] = func() {
		fmt.Println("\033[H\033[2J")
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func GetUserDir() string {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return user.HomeDir
}

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

func ClearTerm() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
