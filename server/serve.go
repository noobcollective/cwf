package server

// Package to start the CWF server and handle all actions.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	http.HandleFunc("/cwf/clear", handleClear)
	http.HandleFunc("/cwf/list", handleList)

	// TODO: Make port either use global var or better via comline line or config file
	log.Fatal(http.ListenAndServe(":8787", nil))
}

// handleStdout is called on `GET` to return the saved content of a file.
func handleStdout(w http.ResponseWriter, r *http.Request) {
	if path.Base(r.URL.Path) != "get" {
		fmt.Fprintf(w, "Invalid endpoint!")
		return
	}

	file := r.URL.Query().Get("file")
	if file == "" {
		fmt.Fprintf(w, "No file name or path provided!")
		return
	}

	content, err := os.ReadFile(file + FILE_SUFFIX)
	if err != nil {
		fmt.Fprintf(w, "File not found!")
		return
	}

	fmt.Fprintf(w, string(content))
}

// handleStdin is called on `POST` to handle file saves.
// It also is able to create a directory, if a full path is sent.
func handleStdin(w http.ResponseWriter, r *http.Request) {
	if path.Base(r.URL.Path) != "copy" {
		fmt.Fprintf(w, "Invalid endpoint!")
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
		fmt.Fprintf(w, "Not allowed!")
		return
	}

	if strings.Contains(body.File, "/") {
		dirs := strings.Split(body.File, "/")
		if len(dirs) > 2 {
			fmt.Fprintf(w, "Not allowed!")
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

	fmt.Fprintf(w, "Saved to: "+body.File)
}

// handleClear is called on `DELETE` to clean the directory or file.
// TODO: Dir support needs to be implemented.
func handleClear(w http.ResponseWriter, r *http.Request) {
	if path.Base(r.URL.Path) != "clear" {
		fmt.Fprintf(w, "Invalid endpoint!")
		return
	}

	file := r.URL.Query().Get("file")
	if file == "" {
		fmt.Fprintf(w, "No file name or path provided!")
		return
	}

	err := os.Remove(file + FILE_SUFFIX)
	if err != nil {
		fmt.Fprintf(w, "File not found!")
		return
	}

	fmt.Fprintf(w, "Cleared!")
}

// TODO: Work in progress currently i just print on the server, we need to return to the client
func handleList(w http.ResponseWriter, r *http.Request) {
	// TODO: This has been now written 5 times we should use a wrapper for this call
	if path.Base(r.URL.Path) != "list" {
		fmt.Fprintf(w, "Invalid endpoint!")
		return
	}

	content, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(content, func(i int, j int) bool {
		fileI, _ := content[i].Info()
		fileJ, _ := content[j].Info()
		return fileI.ModTime().Before(fileJ.ModTime())
	})

	for _, e := range content {
		fmt.Println(e.Name())
	}
}
