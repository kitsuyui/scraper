package main

import (
	"fmt"
	"os"

	"./commands"
)

var exit = os.Exit

func main() {
	if err := commands.CLIRoot.Execute(); err != nil {
		fmt.Println(err)
		exit(1)
	}
}
