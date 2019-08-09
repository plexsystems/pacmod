package main

import (
	"github.com/plexsystems/pacmod/internal/commands"
	"os"
)

func main() {
	err := commands.NewDefaultCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
