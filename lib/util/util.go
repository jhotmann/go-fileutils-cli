package util

import (
	"fmt"
	"os"
	"os/exec"
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
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
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

func GetIntFlag(cmd *cobra.Command, name string, allowedValues []int, defaultValue int) int {
	ret, err := cmd.Flags().GetInt(name)
	if err != nil {
		return defaultValue
	}
	if allowedValues == nil {
		return ret
	}
	if IndexOfInt(ret, allowedValues) > -1 {
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

func IndexOfInt(num int, data []int) int {
	for k, v := range data {
		if num == v {
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
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func RunCommand(args []string, cwd string) error {
	cmd := exec.Cmd{
		Path:   os.Args[0],
		Args:   append([]string{os.Args[0]}, args...),
		Dir:    cwd,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmd.Run()
}
