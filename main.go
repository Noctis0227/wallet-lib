package main

import (
	"fmt"
	"github.com/Noctis0227/wallet-lib/app"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] != "help" {
		a := os.Args[1]
		app.StartCli(a)
	} else {
		printUsage()
	}
}

func printUsage() {
	fmt.Println("following is usage of command line")
	app.CliUsage(1)
}
