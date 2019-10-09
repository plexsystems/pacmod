package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/plexsystems/pacmod/pack"
	"github.com/spf13/cobra"
)

// NewPackCommand creates a new pack command which allows
// the user to package their Go modules
func NewPackCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "pack <version> <outputdirectory>",
		Short: "Package your Go module",
		Args:  cobra.MinimumNArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			return runPackCommand(args)
		},
	}

	return &cmd
}

func runPackCommand(args []string) error {
	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}

	version := args[0]
	outputDirectory := args[1]

	log.Printf("Packing module in path %s...", path)
	if err := pack.Module(path, version, outputDirectory); err != nil {
		return fmt.Errorf("could not package module: %w", err)
	}

	return nil
}
