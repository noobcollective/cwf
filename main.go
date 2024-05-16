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
	certsDir    = flag.String("certsdir", "/etc/crypts/", "Path where the SSL certificate and key are located.")
	certFile    = flag.String("certfile", "", "Filename of the SSL certificate.")
	keyFile     = flag.String("keyfile", "", "Filename of the SSL key.")
	version     string
)

// TODO: Set port and filesDir to shared variables (via config - but where?).

// Client flags
var (
	list     = flag.Bool("l", false, "List files.")
	deletion = flag.Bool("d", false, "Delete file.")
	register = flag.Bool("r", false, "Register a user on server for the first time.")
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Please use args or provide a filename!\n")
		return
	}

	flag.Parse()

	if *showVersion {
		fmt.Println("CWF Version:", version)
		return
	}

	if !*asDaemon {
		client.StartClient()
		return
	}

	usrHome, err := os.UserHomeDir()
	if err != nil {
		zap.L().Error("Could not retrieve home directory!")
		return
	}

	// TODO: Check with iCulture whats best place for logs or move it as a flag or config file...
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		usrHome + "/.config/cwf/cwf.log",
	}

	zap.ReplaceGlobals(zap.Must(cfg.Build()))

	server.StartServer()
}
