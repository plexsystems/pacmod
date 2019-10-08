package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/plexsystems/pacmod/pack"
	"github.com/spf13/cobra"
)

// newPackCommand creates a new pack command
func newPackCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "pack <module> <version> <outputdirectory>",
		Short: "Package your Go module",
		Args:  cobra.MinimumNArgs(3),

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

	path = filepath.ToSlash(path)
	name := args[0]
	version := args[1]
	outputDirectory := args[2]

	module := pack.Module{
		Path:    path,
		Name:    name,
		Version: version,
	}

	log.Printf("Packing module %s...", name)
	if err := module.PackageModule(outputDirectory); err != nil {
		return fmt.Errorf("could not package module: %w", err)
	}

	return nil
}
