package main

import (
	"os"

	"github.com/plexsystems/pacmod/internal/commands"
)

func main() {
	err := commands.NewDefaultCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
