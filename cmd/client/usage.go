package main

import "fmt"

func printUsage(progName string) {
	fmt.Println(progName + " <command> <folder> <filename> [<localfile>]")
	fmt.Println("Available commands:")
	fmt.Println(" - upload (needs local file)")
	fmt.Println(" - delete")
}
