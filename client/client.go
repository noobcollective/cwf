package client
// Package for client. I'm too tired to think of a better explanation.

import (
	"fmt"
	"os"
	"bufio"
)


func StartClient() {
	if !stdInFromPipe() {
		fmt.Println("not from pipe... not implemented yet")
		return
	}

	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	// Read line by line
	for scanner.Scan() {
		// Output line via stdOut
		_, err := fmt.Fprintln(os.Stdout, scanner.Text())
		if err != nil {
			fmt.Println(err)
		}
	}
}


// Check if we are getting content from pipe.
func stdInFromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode() & os.ModeCharDevice == 0
}
