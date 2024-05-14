package client

// Package for client. I'm too tired to think of a better explanation.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"cwf/entities"
	"cwf/utilities"

	"github.com/pelletier/go-toml/v2"
)

var baseURL string

// Init client
func initClient() bool {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not retrieve home directory!")
		return false
	}

	config, err := os.ReadFile(usrHome + "/.config/cwf/config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "No config file found! Check README for config example! Error <%v>\n", err)
		return false
	}

	err = toml.Unmarshal(config, &entities.MotherShip)
	if err != nil {
		fmt.Println("Config file could not be parsed")
		return false
	}

	if entities.MotherShip.MotherShipIP == "" || entities.MotherShip.MotherShipPort == "" {
		fmt.Println("IP address, Port or CWF File directory is not provided")
		return false
	}

	var cwfProtocol string = "http://"
	if entities.MotherShip.MotherShipSSL {
		cwfProtocol = "https://"
	}

	baseURL = cwfProtocol + entities.MotherShip.MotherShipIP + ":" + entities.MotherShip.MotherShipPort + "/cwf/"
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
	} else if utilities.GetFlagValue[string]("r") != "" {
		registerUser()
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
	body, err := json.Marshal(entities.CWFBody_t{Content: encStr})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding data! Error <%v>\n", err)
		return
	}

	res, err := makeRequest("POST", baseURL+"content/"+os.Args[1], bytes.NewBuffer(body))
	if err != nil {
		return
	}

	defer res.Body.Close()
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
		return
	}

	fmt.Println(string(responseData))
}

// Get content of clipboard file.
func getContent() {
	res, err := makeRequest("GET", baseURL+"content/"+os.Args[1], nil)
	if err != nil {
		return
	}

	defer res.Body.Close()
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
	reqUrl := baseURL + "list/"
	if len(os.Args) > 2 {
		reqUrl += os.Args[2]
	}

	res, err := makeRequest("GET", reqUrl, nil)
	if err != nil {
		return
	}

	defer res.Body.Close()
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

	res, err := makeRequest("DELETE", baseURL+"content/"+os.Args[2], nil)
	if err != nil {
		return
	}

	defer res.Body.Close()
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
		return
	}

	fmt.Println(string(responseData))
}

// Registers user with their name and stores the given UUID.
func registerUser() {
	userName := utilities.GetFlagValue[string]("r")
	if userName == "" {
		fmt.Fprintf(os.Stderr, "No username provided!")
		return
	}

	// userName := base64.StdEncoding.EncodeToString([]byte(os.Args[2]))
	// var jsonStr = []byte("{\"username\": " + userName + "\"}")
	// res, err := makeRequest("POST", baseURL + "user/register/", bytes.NewBuffer(jsonStr))
	res, err := makeRequest("GET", baseURL + "register/" + os.Args[2], nil)
	if err != nil {
		return
	}

	defer res.Body.Close()
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
		return
	}

	fmt.Println(string(responseData))
}

// Creates a request object and adds default headers.
// Returns (*http.Response, nil) when successful - (nil, error) otherwise.
func makeRequest(method string, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request! Error <%v>\n", err)
		return nil, err
	}

	// Add needed headers
	req.Header.Set("Cwf-Cli-Req", "true")
	req.Header.Set("Cwf-Cli-Version", "0.3.1")
	req.Header.Set("Cwf-User-Name", "<username>")
	req.Header.Set("Cwf-User-Nonce", "<uuid>")

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request! Err <%v>\n", err)
		return nil, err
	}

	return res, nil
}

// Check if we are getting content from pipe.
func fromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode()&os.ModeCharDevice == 0
}
