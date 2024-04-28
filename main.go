package main

import (
	"fmt"
	"flag"

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

	fmt.Println("Give me that shit.")
}
