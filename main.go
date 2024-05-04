package main

import (
	"cwf/client"
	"cwf/server"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

var asDaemon = flag.Bool("serve", false, "Start as daemon.")
var list = flag.String("l", "", "List files.")

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	if len(os.Args) == 1 {
		panic("Please use args or provide a filename")
	}

	zap.L().Info("Hello from Zap!")

	//listFiles := flag.Bool("l", false, "List all clipboard filenames")
	flag.Parse()
	//fmt.Println(listFiles)

	if *asDaemon {
		server.StartServer()
		fmt.Println("Serving")
		return
	}

	client.StartClient()
}
