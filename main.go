package main

import (
	"cwf/client"
	"cwf/entities"
	"cwf/server"
	"cwf/tools"
	"fmt"
	"flag"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var asDaemon = flag.Bool("serve", false, "Start CWF server.")
var list = flag.Bool("l", false, "List files.")
var listTree = flag.Bool("lt", false, "List files in tree.")
var deletion = flag.Bool("d", false, "Delete file.")

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	flag.Parse()

	// Deferred function to print error without stacktrace.
	defer func() {
        if r := recover(); r != nil {
            fmt.Println(r)
            os.Exit(1)
        }
    }()

	usrHome, err := os.UserHomeDir()
	tools.ExitOnError(err, "Could not retrieve home directory!")

	config, err := os.ReadFile(usrHome + "/.config/cwf/config.yaml")
	tools.ExitOnError(err, "No config file found!")

	err = yaml.Unmarshal(config, &entities.MotherShip)
	tools.ExitOnError(err, "Config file could not be parsed!")

	if entities.MotherShip.MotherShipIP == "" || entities.MotherShip.MotherShipPort == "" {
		tools.ExitWithMsg("IP adress or port not provided in config file.")
	}
}

func main() {
	if len(os.Args) == 1 {
		tools.ExitWithMsg("Please use args or provide a filename!")
	}

	if *asDaemon {
		server.StartServer()
		zap.L().Info("Serving")
		return
	}

	client.StartClient()
}
