package main

import (
	"cwf/client"
	"cwf/server"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

// General variables & flags.
var (
	version string
	showVersion = flag.Bool("version", false, "Prints the program version")
)

// Server flags
var (
	startServer = flag.Bool("serve", false, "Start the cwf server.")
	configPath  = flag.String("config", "/etc/cwf/config.toml", "Path to config file for server usage.")
)

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

	client.Version = version
	server.Version = version
	flag.Parse()

	if *showVersion {
		fmt.Println("CWF Version:", version)
		return
	}

	if !*startServer {
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
