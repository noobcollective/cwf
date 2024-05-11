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
var (
	asDaemon    = flag.Bool("serve", false, "Start as daemon.")
	filesDir    = flag.String("filesdir", "/tmp/cwf/", "Directory to store cwf files.")
	port        = flag.Int("port", 8787, "Port to serve on.")
	https       = flag.Bool("https", false, "Serve with SSL/TLS.")
	showVersion = flag.Bool("version", false, "Prints the program version")
	certPath    = flag.String("certpath", "", "Path where the SSL certificate is located.")
	certKey     = flag.String("keypath", "", "Path where the SSL key is located.")
	version     string
)

// TODO: Set port and filesDir to shared variables (via config - but where?).

// Client flags
var (
	list     = flag.Bool("l", false, "List files.")
	deletion = flag.Bool("d", false, "Delete file.")
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Please use args or provide a filename!\n")
		return
	}

	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	flag.Parse()

	if *showVersion {
		fmt.Println("CWF Version:", version)
		return
	}

	if !*asDaemon {
		client.StartClient()
		return
	}

	server.StartServer()
}
