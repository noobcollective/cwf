package main

import (
	"fmt"
	"flag"
)

func main() {

	asDaemon := flag.Bool("daemon", false, "Start as daemon.")
	flag.Parse()

	if *asDaemon {
		fmt.Println("Starting as daemon...")
	} else {
		fmt.Println("Running with logs...")
	}
}
