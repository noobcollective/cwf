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
	"cwf/utilities"

	"gopkg.in/yaml.v3"
)

var baseURL string

// Init client
func initClient() bool {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not retrieve home directory!")
		return false
	}

	config, err := os.ReadFile(usrHome + "/.config/cwf/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "No config file found! Check README for config example! Error <%v>\n", err)
		return false
	}

	err = yaml.Unmarshal(config, &entities.MotherShip)
	if err != nil {
		fmt.Println("Config file could not be parsed")
		return false
	}

	if entities.MotherShip.MotherShipIP == "" || entities.MotherShip.MotherShipPort == "" {
		fmt.Println("IP address, Port or CWF File directory is not provided")
		return false
	}

	baseURL = "https://" + entities.MotherShip.MotherShipIP + ":" + entities.MotherShip.MotherShipPort + "/cwf"
	return true
}

// Start client and handle action types.
func StartClient() {
	if !initClient() {
		return
	}

	if fromPipe() {
		sendContent()
	} else if utilities.GetFlagValue[bool]("l") {
		listFiles()
	} else if utilities.GetFlagValue[bool]("d") {
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
		return
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println(string(bodyEncoded))
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
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create new Request! Error <%v>\n", err)
		return
	}

	q := req.URL.Query()
	if len(os.Args) > 2 {
		q.Add("dir", os.Args[2])
		req.URL.RawQuery = q.Encode()
	}

	res, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request! Error <%v>\n", err)
		return
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
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
	req, err := http.NewRequest("DELETE", baseURL + "/delete", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Creating a new request with method DELETE failed! Error <%v>\n", err)
		return
	}

	q := req.URL.Query()
	q.Add("path", os.Args[2])
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request! Error <%v>\n", err)
		return
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
		return
	}

	fmt.Println(string(responseData))
}

// Check if we are getting content from pipe.
func fromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode()&os.ModeCharDevice == 0
}
