package server

// Package to start the CWF server and handle all actions.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"cwf/entities"
	"cwf/utilities"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"
)

const file_suffix string = ".cwf"

var filesDir string
var configPath string
var config entities.ServerConfig_t

// Global Variabel to hold users in memory
var users = make(map[string]entities.Account_t)

type cwfChecker_t struct {
	handler http.Handler
}

// Init server
func initServer() bool {
	// Load configuration file.
	configPath = utilities.GetFlagValue[string]("config")
	file, err := utilities.LoadConfig(configPath)
	if err != nil {
		return false
	}

	zap.L().Info("Reading allowed users from config")
	err = toml.Unmarshal(file, &config)
	if err != nil {
		zap.L().Error("Error deconding toml err: " + err.Error())
		return false
	}

	if emptyValues, ok := validateConfig(); !ok {
		fmt.Fprintf(os.Stderr, "Missing values in config: %s!\n", strings.Join(emptyValues, ", "))
		return false
	}

	filesDir = config.General.FilesDir
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(filesDir, 0777); err != nil {
			zap.L().Error(err.Error())
			return false
		}
	}

	zap.L().Info("Generating UUID's for Users")
	for i := range config.Accounts {
		user := &config.Accounts[i]
		id := uuid.New()
		if config.Accounts[i].Registered {
			users[user.Name] = *user
			continue
		}

		config.Accounts[i].ID = id.String()
		users[user.Name] = *user
	}

	tomlContent, err := toml.Marshal(config)
	if err != nil {
		zap.L().Error("Failed to parse config into string.")
		return false
	}

	err = utilities.WriteConfig(configPath, tomlContent)
	return err == nil
}

// Start the server and setup needed endpoints.
func StartServer() {
	if !initServer() {
		return
	}

	zap.L().Info("Welcome to CopyWithFriends on your Server!")

	var certsDir string = config.General.CertsDir
	var certPath string = certsDir + config.General.CertFile
	var keyPath string = certsDir + config.General.KeyFile

	mux := http.NewServeMux()
	mux.HandleFunc("GET /cwf/content/{pathname...}", handleGetContent)
	mux.HandleFunc("POST /cwf/content/{pathname...}", handlePostContent)
	mux.HandleFunc("DELETE /cwf/content/{pathname...}", handleDeleteContent)
	mux.HandleFunc("GET /cwf/list/{pathname...}", handleListContent)
	mux.HandleFunc("GET /cwf/register/{username}", handleAccountRegister)

	// Handler for 404
	mux.HandleFunc("/", handleNotFound)

	var port string = config.General.Port
	zap.L().Info("Serving on Port: " + port)
	if !*config.General.SSL {
		log.Fatal(http.ListenAndServe(":" + port, cwfChecker_t{mux}))
	} else {
		log.Fatal(http.ListenAndServeTLS(":" + port,
			certPath, keyPath, cwfChecker_t{mux}))
	}

}

// handleStdout is called on `GET` to return the saved content of a file.
func handleGetContent(writer http.ResponseWriter, req *http.Request) {
	zap.L().Info("Got GET on /content/")

	pathname := req.PathValue("pathname")
	if pathname == "" {
		zap.L().Info("No file name or path provided!")
		writeRes(writer, http.StatusBadRequest, "No file name or path provided!")
		return
	}

	content, err := os.ReadFile(filesDir + pathname + file_suffix)
	if err != nil {
		// Check if it is maybe a directory
		info, err := os.Stat(filesDir + pathname)
		if err != nil {
			zap.L().Warn("Failed to show stats of file for path: " + filesDir + pathname)
			writeRes(writer, http.StatusBadRequest, "No file name or path provided!")
			return
		}

		if info.IsDir() {
			zap.L().Warn("User tried typed Directory name. Either path is wrong or name Path: " + filesDir + pathname)
			writeRes(writer, http.StatusOK, "Requested File is a Directory, check your path/name")
			return
		}

		zap.L().Info("File not found! Filename: " + pathname)
		writeRes(writer, http.StatusNotFound, "File not found!")

		return
	}

	writer.Write(content)
}

// handlePostContent is called on `POST` to handle file saves.
// It also is able to create a directory, if a full path is sent.
func handlePostContent(writer http.ResponseWriter, req *http.Request) {
	zap.L().Info("Got POST on /content/")

	pathname := req.PathValue("pathname")
	if pathname == "" {
		zap.L().Info("No file name or path provided!")
		writeRes(writer, http.StatusBadRequest, "No file name or path provided!")
		return
	}

	var body entities.CWFBody_t
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		zap.L().Error("Failed to decode request body! Error: " + err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.Contains(pathname, "/") {
		dirs := strings.Split(pathname, "/")
		if len(dirs) > 2 {
			writeRes(writer, http.StatusBadRequest, "Directory depth must not exceed 2 levels")
			return
		}

		if _, err := os.Stat(filesDir + dirs[0]); os.IsNotExist(err) {
			if err := os.Mkdir(filesDir + dirs[0], os.ModePerm); err != nil {
				zap.L().Error("Error while creating new directory: " + err.Error())
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	if err := os.WriteFile(filesDir + pathname + file_suffix, []byte(body.Content), 0644); err != nil {
		zap.L().Error("Error while creating/writing file! Error: " + err.Error())
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writeRes(writer, http.StatusOK, "Saved to: " + pathname)
}

// handleDelete is called on `DELETE` to clean the directory or file.
func handleDeleteContent(writer http.ResponseWriter, req *http.Request) {
	zap.L().Info("Got DELETE on /content/")

	target := req.PathValue("pathname")
	if target == "" {
		zap.L().Warn("No file or path provided")
		writeRes(writer, http.StatusBadRequest, "No file name or path provided!")
		return
	}

	var err error
	var msg string
	if strings.HasSuffix(target, "/") {
		err = os.RemoveAll(filesDir + target)
		msg = "Deleted directory: " + strings.TrimSuffix(target, "/")
	} else {
		err = os.Remove(filesDir + target + file_suffix)
		msg = "Deleted file: " + target
	}

	if err != nil {
		zap.L().Error(err.Error())
		writeRes(writer, http.StatusNotFound, "File or directory not found!")
		return
	}

	writeRes(writer, http.StatusOK, msg)
}

// Function to return all files/directories in given query parameter
// If no further pathname is provided we list files in root folder
func handleListContent(writer http.ResponseWriter, req *http.Request) {
	zap.L().Info("Got GET on /list/")

	targetDir := req.PathValue("pathname")
	targetDir = filesDir + targetDir

	content, err := os.ReadDir(targetDir)
	if err != nil {
		zap.L().Warn(err.Error() + "Called By: " + req.RemoteAddr)
		writeRes(writer, http.StatusNotFound, "No such directory!")
		return
	}

	sort.Slice(content, func(i int, j int) bool {
		fileI, _ := content[i].Info()
		fileJ, _ := content[j].Info()
		return fileI.ModTime().Before(fileJ.ModTime())
	})

	maxNameLen := 0
	for _, entry := range content {
		if len(entry.Name()) > maxNameLen {
			maxNameLen = len(entry.Name())
		}
	}

	maxNameLen += 1

	var response string
	response += fmt.Sprintf("Type    Name" + fmt.Sprintf("%*s", maxNameLen-4, "") + "Modified\n")

	var entryType string

	for _, e := range content {
		if !e.IsDir() && !strings.HasSuffix(e.Name(), ".cwf") {
			continue
		}

		modificationTime, _ := e.Info()
		if e.Type().IsDir() {
			entryType = "Dir"
		} else if e.Type().IsRegular() {
			entryType = "File"
		}

		response += fmt.Sprintf("%-7s%-*s%s\n", entryType, maxNameLen, e.Name(), modificationTime.ModTime().Format("2006-01-02 15:04:05"))
	}

	writeRes(writer, http.StatusOK, response)
}

// handleRegisterAccount for exchanging ID with client
func handleAccountRegister(writer http.ResponseWriter, req *http.Request) {
	zap.L().Info("Got GET on /cwf/register")

	username := req.PathValue("username")
	if username == "" {
		zap.L().Info("No  username provided!")
		writeRes(writer, http.StatusBadRequest, "No username provided!")
		return
	}

	val, ok := users[username]
	if !ok {
		zap.L().Error("Unknown user: " + username)
		http.Error(writer, "Unknown user: "+username, http.StatusBadRequest)
		return
	}

	if val.Registered {
		zap.L().Info("User already registered")
		writeRes(writer, http.StatusBadRequest, "User already registered")
		return
	}

	file, err := utilities.LoadConfig(configPath)
	if err != nil {
		return
	}

	err = toml.Unmarshal(file, &config)
	if err != nil {
		zap.L().Error("Error deconding toml err: " + err.Error())
		return
	}

	for i := range config.Accounts {
		user := &config.Accounts[i]
		if user.Name == username {
			user.Registered = true
			users[user.Name] = *user
			break
		}
	}

	tomlContent, err := toml.Marshal(config)
	if err != nil {
		zap.L().Error("Failed toml marshal")
		writeRes(writer, http.StatusInternalServerError, "Error in handling registration request.")
		return
	}

	err = utilities.WriteConfig(configPath, tomlContent)
	if err != nil {
		return
	}

	// Returning uuid client must use this from now on
	writeRes(writer, http.StatusOK, users[username].ID)
}

// Default handler for 404 pages
func handleNotFound(writer http.ResponseWriter, req *http.Request) {
	zap.L().Warn("User called Endpoint: '" + req.URL.String() + "'!")
	writeRes(writer, http.StatusMethodNotAllowed, "YOU ARE A BAD BOY, ONLY USE cwf client for making requests")
	// TODO: We should probabyl ban/block such ip addresses which try acces endpoints without the cwf client
}

// Sets default header and writes back the response.
func writeRes(writer http.ResponseWriter, statuscode int, content string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(statuscode)
	writer.Write([]byte(content))
}

// Prehandler for all CWF request.
// Checks various headers to determine if usage is safe.
func (checker cwfChecker_t) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	// Check for needed header.
	if _, ok := req.Header["Cwf-Cli-Req"]; !ok {
		http.Error(writer, "Not authorized!", http.StatusForbidden)
		return
	}

	// Match CWF version of server against client.
	cliVersion := req.Header.Get("Cwf-Cli-Version")
	if cliVersion == "" || cliVersion != "0.3.1" {
		zap.L().Warn("Got version: " + cliVersion)
		http.Error(writer, "No version found or version mismatch!", http.StatusBadRequest)
		return
	}

	// If we call register we have to skip the checks
	if strings.Contains(req.URL.Path, "register/") {
		zap.L().Info("User called register skipping checks")
		checker.handler.ServeHTTP(writer, req)
		return
	}

	// Check for user ID.
	userName := req.Header.Get("Cwf-User-Name")
	user, ok := users[userName]
	if !ok {
		zap.L().Warn("User not found! Please register")
		http.Error(writer, "User not found! Please register ", http.StatusForbidden)
		return
	}

	// Check for user ID.
	userID := req.Header.Get("Cwf-User-Id")
	if user.ID != userID {
		http.Error(writer, "Wrong ID", http.StatusForbidden)
		return
	}

	checker.handler.ServeHTTP(writer, req)
}

// Checks if there are mising values in config file.
// Returns empty fields and bool to check if config is valid.
func validateConfig() ([]string, bool) {
	var emptyValues []string

	if config.General.Port == "" {
		emptyValues = append(emptyValues, "Port")
	}

	if config.General.FilesDir == "" {
		emptyValues = append(emptyValues, "FilesDir")
	}

	if config.General.SSL == nil {
		emptyValues = append(emptyValues, "SSL")
	} else if *config.General.SSL {
		if config.General.CertsDir == "" {
			emptyValues = append(emptyValues, "CertsDir")
		}

		if config.General.CertFile == "" {
			emptyValues = append(emptyValues, "CertFile")
		}

		if config.General.KeyFile == "" {
			emptyValues = append(emptyValues, "Keyfile")
		}
	}

	return emptyValues, len(emptyValues) == 0
}
