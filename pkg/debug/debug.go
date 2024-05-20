package debug

import "fmt"

func DebugMessage(debug bool, msg string) {
	if debug {
		fmt.Println(msg)
	}
}
