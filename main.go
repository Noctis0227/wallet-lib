package main

import (
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/app"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] != "help" {
		a := os.Args[0]
		app.StartCli(a)
	} else {
		printUsage()
	}
}

func printUsage() {
	fmt.Println("following is usage of command line")
	app.CliUsage(1)
}
