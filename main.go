package main

import (
	"buanu/cmdsum"
	_ "embed"
	"fmt"
	"log"
	"os"
)

//go:embed help.txt
var help string

func main() {
	if len(os.Args) < 2 {
		fmt.Print(help)
		os.Exit(1)
	}
	subcommand := os.Args[1]
	switch subcommand {
	case "sum":
		cmdsum.Run()
	default:
		log.Fatalf("%q is not a valid subcommand.", subcommand)
	}
}
