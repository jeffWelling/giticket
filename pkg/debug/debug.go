package debug

import "fmt"

// DebugMessage takes a flag and a message, and prints the message with
// fmt.Println based on flag
func DebugMessage(debug bool, msg string) {
	if debug {
		fmt.Println(msg)
	}
}
