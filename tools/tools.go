package tools

import (
	"fmt"
	"os"
)

// Check if error occured and exit if so.
func ExitOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

// Exit with message.
func ExitWithMsg(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
