package main

import (
	"fmt"
	"git.diabin.com/BlockChain/wallet-lib/app"
	"os"
)

var version = "kahf v0.1.6"

func main() {

	if len(os.Args) > 1 && os.Args[1] != "help" {
		a := os.Args[1]

		app.StartCli(a)

		//if a == "--console" {
		//	app.StartConsole()
		//} else if a == "--serv" {
		//	app.StartService()
		//} else if a == "--version" {
		//	fmt.Println(version)
		//} else if a == "--h" {
		//	printUsage()
		//} else {
		//	app.StartCli(a)
		//}
	} else {
		printUsage()
	}
}

func printUsage() {
	fmt.Println("--serv:  start the wallet service")
	fmt.Println()
	fmt.Println("--console:  start console interaction mode")
	fmt.Println()
	fmt.Println("--version:  show kahf version")
	fmt.Println()
	fmt.Println("following is usage of command line")
	app.CliUsage(1)
}
