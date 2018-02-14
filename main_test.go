package main

import (
	"testing"

	"github.com/kitsuyui/scraper/commands"
)

func TestMain(t *testing.T) {
	exit = func(i int) {
		if i != 1 {
			t.Error("exit code Must be 1")
		}
	}
	commands.CLIRoot.SetArgs([]string{"invalid-command"})
	main()
}
