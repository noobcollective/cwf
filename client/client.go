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
	fmt.Println(os.Args)
	if stdInFromPipe() {
		sendContent()
		return
	}


	// 1) If no flag provided we want to paste from server
	res, err := http.Get("http://127.0.0.1:8787/cwf/get?file=" + os.Args[1])

	if err != nil {
		fmt.Println("Error getting content!")
		return
	}

	bodyEncoded, err := io.ReadAll(res.Body)
	bodyDecoded, err := base64.StdEncoding.DecodeString(string(bodyEncoded))
        if err != nil {
			fmt.Println("Failed to decode body!")
        }
	fmt.Println(string(bodyDecoded))
}


func sendContent() {
	content, err := io.ReadAll(os.Stdin)

	if err != nil {
		fmt.Println("Problem reading content.")
		return
	}

	encStr := base64.StdEncoding.EncodeToString(content)
	body, err := json.Marshal(entities.CWFBody_t{File: os.Args[1], Content: encStr})
	// TODO: Handle response correctly (e.g. file already exists -> prompt to override)
	res, err := http.Post("http://127.0.0.1:8787/cwf/copy", "application/json", bytes.NewBuffer(body))

	bodyStr, err := io.ReadAll(res.Body)

	fmt.Println(string(bodyStr))
}

// Check if we are getting content from pipe.
func stdInFromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode() & os.ModeCharDevice == 0
}
