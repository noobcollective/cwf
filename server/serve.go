package server

// Package to start the CWF server and handle all actions.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"

	"cwf/entities"
)

var FILE_SUFFIX string = ".cwf"

// Start the server and setup needed endpoits.
func StartServer() {
	// Endpoints
	http.HandleFunc("/cwf/get", handleStdout)
	http.HandleFunc("/cwf/copy", handleStdin)
	http.HandleFunc("/cwf/delete", handleDelete)
	http.HandleFunc("/cwf/list", handleList)

	// TODO: Make port either use global var or better via comline line or config file
	log.Fatal(http.ListenAndServe(":8787", nil))
}

// handleStdout is called on `GET` to return the saved content of a file.
func handleStdout(w http.ResponseWriter, r *http.Request) {
	if !allowedEndpoint(r.URL, "get") {
		writeRes(w, http.StatusForbidden, "Invalid endpoint!")
		return
	}

	file := r.URL.Query().Get("file")
	if file == "" {
		writeRes(w, http.StatusBadRequest, "No file name or path provided!")
		return
	}

	content, err := os.ReadFile(file + FILE_SUFFIX)
	if err != nil {
		writeRes(w, http.StatusNotFound, "File not found!")
		return
	}

	w.Write(content)
}

// handleStdin is called on `POST` to handle file saves.
// It also is able to create a directory, if a full path is sent.
func handleStdin(w http.ResponseWriter, r *http.Request) {
	if !allowedEndpoint(r.URL, "copy") {
		writeRes(w, http.StatusForbidden, "Invalid endpoint!")
		return
	}

	var body entities.CWFBody_t
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Resolve path and compare with configured basedir
	// Example confiuration: /tmp/cwf/
	// body.File = ../test.cwf resolves to: `/tmp` -> not allowed (not in basedir)
	// body.File = ../var/ resolves to: `/var` -> also not allowed (not in basedir)
	if strings.Contains(body.File, "..") {
		writeRes(w, http.StatusForbidden, "Not allowd!")
		return
	}

	if strings.Contains(body.File, "/") {
		dirs := strings.Split(body.File, "/")
		if len(dirs) > 2 {
			writeRes(w, http.StatusForbidden, "Not allowd!")
			return
		}

		if _, err := os.Stat(dirs[0]); os.IsNotExist(err) {
			err = os.Mkdir(dirs[0], os.ModePerm)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	err = os.WriteFile(body.File+FILE_SUFFIX, []byte(body.Content), 0644)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeRes(w, http.StatusOK, "Saved to: " + body.File)
}

// handleDelete is called on `DELETE` to clean the directory or file.
// TODO: Dir support needs to be implemented.
func handleDelete(w http.ResponseWriter, r *http.Request) {
	if !allowedEndpoint(r.URL, "delete") {
		writeRes(w, http.StatusForbidden, "Invalid endpoint!")
		return
	}

	var body entities.CWFBody_t
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.Contains(body.File, "..") {
		writeRes(w, http.StatusForbidden, "Not allowd!")
		return
	}

	file := r.URL.Query().Get("file")
	if file == "" {
		writeRes(w, http.StatusBadRequest, "No file name or path provided!")
		return
	}

	err = os.Remove(file + FILE_SUFFIX)
	if err != nil {
		writeRes(w, http.StatusNotFound, "File not found!")
		return
	}

	writeRes(w, http.StatusOK, "Deleted file: " + file)
}

// TODO: Work in progress currently i just print on the server, we need to return to the client
func handleList(w http.ResponseWriter, r *http.Request) {
	// TODO: This has been now written 5 times we should use a wrapper for this call
	// INFO: My implemenation = not really helpfull
	if !allowedEndpoint(r.URL, "list") {
		writeRes(w, http.StatusForbidden, "Invalid endpoint!")
		return
	}

	targetDir := r.URL.Query().Get("dir")

	if strings.Contains(targetDir, "..") {
		fmt.Println("------------------------")
		fmt.Println("Client tried to call \"..\" Called by: " + r.RemoteAddr)
		fmt.Println("------------------------")
		writeRes(w, http.StatusForbidden, "Not allowed!")
		return
	}

	content, err := os.ReadDir("./" + targetDir)
	if err != nil {
		fmt.Println("------------------------")
		fmt.Println(err)
		fmt.Println("Called by: " + r.RemoteAddr)
		fmt.Println("------------------------")
		writeRes(w, http.StatusNotFound, "No such directory!")
		return
	}

	sort.Slice(content, func(i int, j int) bool {
		fileI, _ := content[i].Info()
		fileJ, _ := content[j].Info()
		return fileI.ModTime().Before(fileJ.ModTime())
	})

	var response string
	for _, e := range content {
		modificationTime, _ := e.Info()
		if e.Type().IsDir() {
			response += "Dir: \t" + e.Name() + "\t\t Modified: " + modificationTime.ModTime().Format("2006-01-02 15:04:05") + "\n"
			continue
		} else if e.Type().IsRegular() {
			response += "File: \t" + e.Name() + "\t\t Modified: " + modificationTime.ModTime().Format("2006-01-02 15:04:05") + "\n"
			continue
		}
	}

	writeRes(w, http.StatusOK, response)
}

// Respond the go way.
func writeRes(w http.ResponseWriter, statuscode int, content string) {
	w.WriteHeader(statuscode)
	w.Write([]byte(content))
}

// Check if endpoint is allowed for current action.
func allowedEndpoint(filepath *url.URL, endpoint string) bool {
	return path.Base(filepath.Path) == endpoint
}
