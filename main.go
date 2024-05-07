package main

import (
	"cwf/client"
	"cwf/entities"
	"cwf/server"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var asDaemon = flag.Bool("serve", false, "Start as daemon.")

var list = flag.Bool("l", false, "List files.")

// var listTree = flag.Bool("lt", false, "List files in tree.")
var deletion = flag.Bool("d", false, "Delete file.")

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Please use args or provide a filename!\n")
		return
	}

	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	flag.Parse()

	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not retrieve home directory! Error <%v>\n", err)
		return
	}

	config, err := os.ReadFile(userHome + "/.config/cwf/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "No config file found! Check README for config example! Error <%v>\n", err)
		return
	}

	err = yaml.Unmarshal(config, &entities.MotherShip)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config file could not be parsed! Error <%v>\n", err)
		return
	}

	if (*asDaemon && entities.MotherShip.MotherShipPort == "") || (!*asDaemon && (entities.MotherShip.MotherShipIP == "" || entities.MotherShip.MotherShipPort == "")) {
		fmt.Fprintf(os.Stderr, "IP address or Port is not provided\n")
		return
	}

	if *asDaemon {
		server.StartServer()
		zap.L().Info("Serving on Port: " + entities.MotherShip.MotherShipPort)
		return
	}

	client.StartClient()
}
