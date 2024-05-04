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
var list = flag.String("l", "", "List files.")

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

	if entities.MotherShip.MotherShipIP == "" || entities.MotherShip.MotherShipPort == "" {
		panic("IP address to Server is not provided")
	}
}

func main() {
	if len(os.Args) == 1 {
		panic("Please use args or provide a filename")
	}

	zap.L().Info("Welcome to CopyWithFriends -> cwf")

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
