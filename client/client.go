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
	"os"

	"cwf/entities"
)

var baseURL string

// Start client and handle action types.
func StartClient() {
	baseURL = "http://" + entities.MotherShip.MotherShipIP + ":" + entities.MotherShip.MotherShipPort + "/cwf"

	if fromPipe() {
		sendContent()
	} else if getFlagValue("l") {
		listFiles()
	} else if getFlagValue("lt") {
		// listTree()
	} else if getFlagValue("d") {
		deleteFile()
	} else {
		getContent()
	}
}

// Send content to server to save it encoded in specified file.
func sendContent() {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading content from StdIn! Error <%v>\n", err)
		return
	}

	encStr := base64.StdEncoding.EncodeToString(content)
	body, err := json.Marshal(entities.CWFBody_t{File: os.Args[1], Content: encStr})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding data! Error <%v>\n", err)
		return
	}

	res, err := http.Post(baseURL+"/copy",
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request! Error <%v>\n", err)
		return
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Unexpected Status code. Expected <OK> got <%v>\n", res.StatusCode)
		return
	}

	fmt.Println(string(responseData))
}

// Get content of clipboard file.
func getContent() {
	res, err := http.Get(baseURL + "/get?file=" + os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting content! Err <%v>\n", err)
		return
	}

	bodyEncoded, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Unexpected Status code. Expected <OK> got <%v>\n", res.StatusCode)
		return
	}

	bodyDecoded, err := base64.StdEncoding.DecodeString(string(bodyEncoded))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode body! Error <%v>\n", err)
		return
	}

	fmt.Println(string(bodyDecoded))
}

// Get a list from server.
func listFiles() {
	requestUrl := baseURL + "/list"
	if len(os.Args) > 2 {
		requestUrl += "?dir=" + os.Args[2]
	}

	res, err := http.Get(requestUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request! Error <%v>\n", err)
		return
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Unexpected Status code. Expected <OK> got <%v>\n", res.StatusCode)
		return
	}

	fmt.Println(string(responseData))
}

// Delete a filename from server.
func deleteFile() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "No filename given to delete!\n")
		return
	}

	client := &http.Client{}
	requestUrl := baseURL + "/delete?file=" + os.Args[2]
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Creating a new request with method DELETE failed! Error <%v>\n", err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request! Error <%v>\n", err)
		return
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Unexpected Status code. Expected <OK> got <%v>\n", res.StatusCode)
		return
	}

	fmt.Println(string(responseData))
}

// Check if we are getting content from pipe.
func fromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode()&os.ModeCharDevice == 0
}

// Get value of a flag.
func getFlagValue(flagName string) bool {
	return flag.Lookup(flagName).Value.(flag.Getter).Get().(bool)
}
