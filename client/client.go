package client

// Package for client. I'm too tired to think of a better explanation.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"cwf/entities"
)


func StartClient() {
	if !stdInFromPipe() {
		fmt.Println("not from pipe... do we want to handle this?")
		return
	}

	content, err := io.ReadAll(os.Stdin)

	if err != nil {
		fmt.Println("Problem reading content.")
		return
	}

	encStr := base64.StdEncoding.EncodeToString(content)
	body, err := json.Marshal(entities.CWFBody_t{File: "test", Content: encStr})
	res, err := http.Post("http://127.0.0.1:8787/cwf/copy", "application/json", bytes.NewBuffer(body))

	fmt.Println(string(body))
	fmt.Println(res)
}


// Check if we are getting content from pipe.
func stdInFromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode() & os.ModeCharDevice == 0
}
