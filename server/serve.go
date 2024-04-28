package server

import (
	"io"
	"net/http"
	"os"
	"log"
	"fmt"
	"path"
	// "time"
	// "syscall"
	// "context"
	// "os/signal"
)


func StartServer() {
	// Endpoints
	http.HandleFunc("/cwf/get", handleStdout)
	http.HandleFunc("/cwf/copy/{file}", handleStdin)
	http.HandleFunc("/cwf/clear", handleClear)

	log.Fatal(http.ListenAndServe(":8787", nil))
}


func handleStdout(w http.ResponseWriter, r *http.Request) {
	if (path.Base(r.URL.Path) != "get") {
		fmt.Fprintf(w, "Invalid endpoint!")
	}

	file := r.URL.Query().Get("file")
	if (file == "") {
		fmt.Fprintf(w, "No file name or path provided!")
		return
	}

	content, err := os.ReadFile(file)
	if (err != nil) {
		fmt.Fprintf(w, "File not found!")
		return
	}

	fmt.Fprintf(w, string(content))
}


func handleStdin(w http.ResponseWriter, r *http.Request) {
	if (path.Base(r.URL.Path) != "copy") {
		fmt.Fprintf(w, "Invalid endpoint!")
	}

	file := r.PathValue("file")
	if (file == "") {
		fmt.Fprintf(w, "No filename specified!")
	}

	body, err := io.ReadAll(r.Body)
	if (err != nil) {
		fmt.Fprintf(w, "Problem reading content!")
	}

	fmt.Println(string(body))
}


func handleClear(http.ResponseWriter, *http.Request) {
	fmt.Println("Clearing...")
}
