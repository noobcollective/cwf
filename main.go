package main

import (
	"cwf/client"
	"cwf/server"
	"flag"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		panic("Please use args or provide a filename")
	}

	asDaemon := flag.Bool("serve", false, "Start as daemon.")

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
