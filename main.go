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
	version string = "dev"
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

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"/var/log/cwf/cwf.log",
	}

	zap.ReplaceGlobals(zap.Must(cfg.Build()))

	server.StartServer()
}
