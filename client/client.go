package client

// Package for client. I'm too tired to think of a better explanation.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"os"

	"cwf/entities"
)

var flagLookup = map[string]string{
	"-l":  "list",
	"-lt": "list-tree",
}

// Start client and handle action types.
func StartClient() {
	if fromPipe() {
		sendContent()
		return
	}

	flag := checkArgsFlags()
	switch flag {
		case "-l", "-lt":
			listFiles(flag)
		case "-d":
			deleteFile()

		// If no flags provided we want to print out content.
		case "":
		default:
			getContent()
	}
}

// Send content to server to save it encoded in specified file.
func sendContent() {
	if !checkArgsFilename() { return }
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic("Error reading content from StdIn")
	}

	encStr := base64.StdEncoding.EncodeToString(content)
	body, err := json.Marshal(entities.CWFBody_t{File: os.Args[1], Content: encStr})
	if err != nil {
		panic("Error encoding data.")
	}

	res, err := http.Post("http://127.0.0.1:8787/cwf/copy",
		"application/json", bytes.NewBuffer(body))
	// TODO: Handle response correctly (e.g. file already exists -> prompt to override)
	if err != nil {
		panic("Error sending request.")
	}

	bodyStr, err := io.ReadAll(res.Body)
	fmt.Println(string(bodyStr))
}

// Get content of clipboard file.
func getContent() {
	if !checkArgsFilename() { return }
	res, err := http.Get("http://127.0.0.1:8787/cwf/get?file=" + os.Args[1])
	if err != nil {
		panic("Error getting content!")
	}

	bodyEncoded, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(string(bodyEncoded))
		return
	}

	bodyDecoded, err := base64.StdEncoding.DecodeString(string(bodyEncoded))
	if err != nil {
		panic("Failed to decode body!")
	}

	fmt.Println(string(bodyDecoded))
}

// Get a list from server.
func listFiles(flag string) {
	requestUrl := "http://127.0.0.1:8787/cwf/" + flagLookup[flag]

	res, err := http.Get(requestUrl)
	if err != nil {
		panic("Error sending request.")
	}

	bodyEncoded, err := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Println(string(bodyEncoded))
		return
	}

	// bodyDecoded, err := base64.StdEncoding.DecodeString(string(bodyEncoded))
	// if err != nil { panic("Failed to decode body!") }

	fmt.Println(string(bodyEncoded))
}

// Delete a filename from server.
func deleteFile() {
	if len(os.Args) < 3 {
		fmt.Println("No filename given to delete.")
		return
	}

	client := &http.Client{}
	requestUrl := "http://127.0.0.1:8787/cwf/delete?file=" + os.Args[2]
	req, err := http.NewRequest("DELETE", requestUrl, nil)

	res, err := client.Do(req)
	if err != nil { panic("Error sending request.") }

	fmt.Println(res)
}

// Check if we are getting content from pipe.
func fromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode()&os.ModeCharDevice == 0
}

// Check arguments for filename -> panics if no filename is given via args.
func checkArgsFilename() bool {
	if !strings.Contains(os.Args[1], "-") && len(os.Args) >= 2 {
		return true
	}

	fmt.Println("No filename given.")
	return false
}

// Check if flags have been passed and return the passed flag.
func checkArgsFlags() string {
	if strings.Contains(os.Args[1], "-") && len(os.Args) >= 2 {
		return os.Args[1]
	}

	return ""
}
