package main

import (
	"cwf/client"
	"cwf/entities"
	"cwf/server"
	"flag"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var asDaemon = flag.Bool("serve", false, "Start as daemon.")
var list = flag.Bool("l", false, "List files.")
var listTree = flag.Bool("lt", false, "List files in tree.")
var deletion = flag.Bool("d", false, "Delete file.")

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))

	// TODO: This should not be hardcoded i guess
	config, err := os.ReadFile("./config/config.yaml")

	if err != nil {
		panic("No config file found")
	}

	err = yaml.Unmarshal(config, &entities.MotherShip)
	if err != nil {
		panic("Config file could not be parsed")
	}

	if entities.MotherShip.MotherShipIP == "" || entities.MotherShip.MotherShipPort == "" || entities.MotherShip.MotherShipCWFDirectory == "" {
		panic("IP address, Port or CWF File directory is not provided")
	}
}

func main() {
	if len(os.Args) == 1 {
		panic("Please use args or provide a filename")
	}

	//listFiles := flag.Bool("l", false, "List all clipboard filenames")
	flag.Parse()
	//fmt.Println(listFiles)

	if *asDaemon {
		server.StartServer()
		zap.L().Info("Serving")
		return
	}

	client.StartClient()
}
