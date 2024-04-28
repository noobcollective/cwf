package server
// Package to start the CWF server and handle all actions.

import (
	"net/http"
	"os"
	"log"
	"fmt"
	"path"
	"encoding/json"
	"strings"
)


// Typedef to parse the JSON body of the request correctly.
type JsonFileReq struct {
    File string
    Content string
}


// Start the server and setup needed endpoits.
func StartServer() {
	// Endpoints
	http.HandleFunc("/cwf/get", handleStdout)
	http.HandleFunc("/cwf/copy", handleStdin)
	http.HandleFunc("/cwf/clear", handleClear)

	log.Fatal(http.ListenAndServe(":8787", nil))
}


// handleStdout is called on `GET` to return the saved content of a file.
func handleStdout(w http.ResponseWriter, r *http.Request) {
	if (path.Base(r.URL.Path) != "get") {
		fmt.Fprintf(w, "Invalid endpoint!")
		return
	}

	file := r.URL.Query().Get("file")
	if (file == "") {
		fmt.Fprintf(w, "No file name or path provided!")
		return
	}

	content, err := os.ReadFile(file + ".cwf")
	if (err != nil) {
		fmt.Fprintf(w, "File not found!")
		return
	}

	fmt.Fprintf(w, string(content))
}


// handleStdin is called on `POST` to handle file saves.
// It also is able to create a directory, if a full path is sent.
func handleStdin(w http.ResponseWriter, r *http.Request) {
	if (path.Base(r.URL.Path) != "copy") {
		fmt.Fprintf(w, "Invalid endpoint!")
		return
	}

	var body JsonFileReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if (strings.ContainsAny("/", body.File)) {
		dir := strings.Split(body.File, "/")

		if (len(dir) > 2) {
			fmt.Fprintf(w, "Only one directory allowed!")
			return
		}

		err = os.Mkdir(dir[0], os.ModePerm)
		if (err != nil) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = os.WriteFile(body.File + ".cwf", []byte(body.Content), 0644)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Saved to %s", body.File)
}


// handleClear is called on `DELETE` to clean the directory or file.
// TODO: Dir support needs to be implemented.
func handleClear(w http.ResponseWriter, r *http.Request) {
	if (path.Base(r.URL.Path) != "clear") {
		fmt.Fprintf(w, "Invalid endpoint!")
		return
	}

	file := r.URL.Query().Get("file")
	if (file == "") {
		fmt.Fprintf(w, "No file name or path provided!")
		return
	}

	err := os.Remove(file + ".cwf")
	if (err != nil) {
		fmt.Fprintf(w, "File not found!")
		return
	}

	fmt.Fprintf(w, "Cleared!")
}
