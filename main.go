package main

import (
	"fmt"
	"flag"

	"cwf/client"
	"cwf/server"
)

func main() {

	asDaemon := flag.Bool("serve", false, "Start as daemon.")
	flag.Parse()

	if *asDaemon {
		server.StartServer()
		fmt.Println("Serving")
		return
	}

	client.StartClient()
}
