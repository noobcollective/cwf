package main

import (
	"cwf/client"
	"cwf/server"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

// Server flags
var asDaemon = flag.Bool("serve", false, "Start as daemon.")
var filesDir = flag.String("filesdir", "/tmp/cwf/", "Directory to store cwf files.")
var port = flag.Int("port", 8787, "Port to serve on.")

// TODO: Set port and filesDir to shared variables (via config - but where?).

// Client flags
var list = flag.Bool("l", false, "List files.")
var deletion = flag.Bool("d", false, "Delete file.")

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Please use args or provide a filename!\n")
		return
	}

	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	flag.Parse()

	if !*asDaemon {
		client.StartClient()
		return
	}

	server.StartServer()
}
