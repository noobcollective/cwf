package client

// Package for client. I'm too tired to think of a better explanation.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"flag"
	"io"
	"net/http"
	"os"

	"cwf/entities"
)

var flagLookup = map[string]string {
	"-l": "list",
	"-lt": "list-tree",
	"-d": "delete",
}

// Start client and handle action types.
func StartClient() {
	if fromPipe() {
		sendContent()
		return
	}

	fmt.Println(flag.Lookup("l").Value)
	// If no flags provided we want to print out content.
	// getContent()
}


// Send content to server to save it encoded in specified file.
func sendContent() {
	checkArgsFilename()
	content, err := io.ReadAll(os.Stdin)
	if err != nil { panic("Error reading content from StdIn") }

	encStr := base64.StdEncoding.EncodeToString(content)
	body, err := json.Marshal(entities.CWFBody_t{File: os.Args[1], Content: encStr})
	if err != nil { panic("Error encoding data.") }

	res, err := http.Post("http://127.0.0.1:8787/cwf/copy",
					"application/json", bytes.NewBuffer(body))
	// TODO: Handle response correctly (e.g. file already exists -> prompt to override)
	if err != nil { panic("Error sending request.") }

	bodyStr, err := io.ReadAll(res.Body)
	fmt.Println(string(bodyStr))
}


// Get content of clipboard file.
func getContent() {
	checkArgsFilename()
	res, err := http.Get("http://127.0.0.1:8787/cwf/get?file=" + os.Args[1])
	if err != nil { panic("Error getting content!") }

	bodyEncoded, err := io.ReadAll(res.Body)
	bodyDecoded, err := base64.StdEncoding.DecodeString(string(bodyEncoded))
	if err != nil { panic("Failed to decode body!") }

	fmt.Println(string(bodyDecoded))
}


// Check for flags and decide what to do.
// INFO: WIP - not all eventes are handled yet.
func doFlaggedAction(flag string) {
	requestUrl := "127.0.0.1:8787/cwf/" + flagLookup[flag]

	if flag == "-d" {
		checkArgsFilename()
		requestUrl += "?file=" + os.Args[1]
	}
	res, err := http.Get(requestUrl)
	if err != nil { panic("Error sending request.") }

	fmt.Println(res)
}


// Check if we are getting content from pipe.
func fromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode() & os.ModeCharDevice == 0
}


// Check arguments for filename -> panics if no filename is given via args.
func checkArgsFilename() {
	if (len(os.Args) < 2) { panic("No filename given.") }
}
