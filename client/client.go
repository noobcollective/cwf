package client

// Package for client. I'm too tired to think of a better explanation.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"cwf/entities"
	"cwf/utilities"

	"github.com/pelletier/go-toml/v2"
)

var Version string
var baseURL string
var config entities.ClientConfig_t

// Init client
func initClient() bool {
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not retrieve home directory!")
		return false
	}

	configContent, err := os.ReadFile(userHome + "/.config/cwf/config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "No config file found! Check README for config example!\n")
		return false
	}

	if err := toml.Unmarshal(configContent, &config); err != nil {
		fmt.Println("Config file could not be parsed")
		return false
	}

	if emptyValues, ok := validateConfig(); !ok {
		fmt.Fprintf(os.Stderr, "Missing values in config: %s!\n", strings.Join(emptyValues, ", "))
		return false
	}

	var cwfProtocol string = "http://"
	if *config.Mothership.SSL {
		cwfProtocol = "https://"
	}

	baseURL = cwfProtocol + config.Mothership.IP + ":" + config.Mothership.Port + "/cwf/"
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
	} else if utilities.GetFlagValue[bool]("r") {
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

	res, err := makeRequest("POST", baseURL + "content/" + os.Args[1], bytes.NewBuffer(body))
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
	res, err := makeRequest("GET", baseURL + "content/" + os.Args[1], nil)
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

	res, err := makeRequest("DELETE", baseURL + "content/" + os.Args[2], nil)
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
	res, err := makeRequest("GET", baseURL + "register/" + config.Client.User, nil)
	if err != nil {
		return
	}

	defer res.Body.Close()
	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Reading response body! Error <%v>\n", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println(string(responseData))
		return
	}

	userHome, err := os.UserHomeDir()
	configFile, err := os.OpenFile(userHome + "/.config/cwf/config.toml", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config file! Error <%v>\n", err)
		return
	}

	config.Client.ID = string(responseData)
	if err := toml.NewEncoder(configFile).Encode(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving UUID to config file! Error <%v>\n", err)
		return
	}

	if err = configFile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to close config file: <%v>", err)
		return
	}

	fmt.Println("Successfully registered! Have fun using CWF!")
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

	if err := checkUserStatus(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request! Error <%v>\n", err)
		return nil, err
	}

	// Add needed headers
	req.Header.Set("Cwf-Cli-Req", "true")
	req.Header.Set("Cwf-Cli-Version", Version)
	req.Header.Set("Cwf-User-Name", config.Client.User)
	req.Header.Set("Cwf-User-Id", config.Client.ID)

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request! Err <%v>\n", err)
		return nil, err
	}

	return res, nil
}


// Checks if some input is given via pipe and returns result.
func fromPipe() bool {
	content, _ := os.Stdin.Stat()
	return content.Mode()&os.ModeCharDevice == 0
}


// Checks status of registration for current user.
// Return nil | error do be handled from caller.
func checkUserStatus() error {
	var err error = nil
	var isRegister bool = utilities.GetFlagValue[bool]("r")

	if !isRegister && config.Client.ID == "" {
		err = errors.New("ID not found! Try running with '-r' flag to register on server!")
	}

	return err
}


// Checks if there are mising values in config file.
// Returns empty fields and bool to check if config is valid.
func validateConfig() ([]string, bool) {
	var emptyValues []string

	if config.Mothership.IP == "" {
		emptyValues = append(emptyValues, "IP")
	}

	if config.Mothership.Port == "" {
		emptyValues = append(emptyValues, "Port")
	}

	if config.Mothership.SSL == nil {
		emptyValues = append(emptyValues, "SSL")
	}

	if config.Client.User == "" {
		emptyValues = append(emptyValues, "User")
	}

	return emptyValues, len(emptyValues) == 0
}
