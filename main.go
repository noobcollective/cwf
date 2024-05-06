package main

import (
	"cwf/client"
	"cwf/entities"
	"cwf/server"
	"fmt"
	"flag"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var asDaemon = flag.Bool("serve", false, "Start as daemon.")
var list = flag.Bool("l", false, "List files.")
var listTree = flag.Bool("lt", false, "List files in tree.")
var deletion = flag.Bool("d", false, "Delete file.")

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please use args or provide a filename!")
		return
	}

	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	flag.Parse()

	usrHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not retrieve home directory!")
		return
	}

	config, err := os.ReadFile(usrHome + "/.config/cwf/config.yaml")
	if err != nil {
		fmt.Println("No config file found!")
		return
	}

	err = yaml.Unmarshal(config, &entities.MotherShip)
	if err != nil {
		fmt.Println("Config file could not be parsed")
		return
	}

	if entities.MotherShip.MotherShipIP == "" || entities.MotherShip.MotherShipPort == "" {
		fmt.Println("IP address, Port or CWF File directory is not provided")
		return
	}

	if *asDaemon {
		server.StartServer()
		zap.L().Info("Serving")
		return
	}

	client.StartClient()
}
