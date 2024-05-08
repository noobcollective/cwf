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
	} else if initServer() {
		server.StartServer(*port, *filesDir)
		zap.L().Info("Serving")
	}
}

func initServer() bool {
	if _, err := os.Stat(*filesDir); os.IsNotExist(err) {
		err := os.Mkdir(*filesDir, 0777)
		if err != nil {
			zap.L().Error(err.Error())
			return false
		}
	}

	return true
}

func initClient() bool {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not retrieve home directory!")
		return false
	}

	config, err := os.ReadFile(usrHome + "/.config/cwf/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "No config file found! Check README for config example! Error <%v>\n", err)
		return false
	}

	err = yaml.Unmarshal(config, &entities.MotherShip)
	if err != nil {
		fmt.Println("Config file could not be parsed")
		return false
	}

	if entities.MotherShip.MotherShipIP == "" || entities.MotherShip.MotherShipPort == "" {
		fmt.Println("IP address, Port or CWF File directory is not provided")
		return false
	}

	return true
}
