package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/plexsystems/pacmod/pack"
	"github.com/spf13/cobra"
)

// NewPackCommand creates a new pack command
func NewPackCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "pack <module> <version>",
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

	name := args[0]
	version := args[1]

	module := pack.Module{
		Path:    path,
		Name:    name,
		Version: version,
	}

	outputDirectory := filepath.Join(module.Path, module.Version)
	if err := os.Mkdir(outputDirectory, 0777); err != nil {
		return fmt.Errorf("could not create output directory: %w", err)
	}

	log.Println("Creating zip archive...")
	if err := pack.ZipModule(outputDirectory); err != nil {
		return fmt.Errorf("could not create archive: %w", err)
	}

	log.Println("Creating info file...")
	if err := pack.CreateInfoFile(version, outputDirectory); err != nil {
		return fmt.Errorf("could not create info file: %w", err)
	}

	log.Println("Copying mod file...")
	if err := copyModuleFile(module.Path, outputDirectory); err != nil {
		return fmt.Errorf("could not copy mod file: %w", err)
	}

	return nil
}

func copyModuleFile(source string, destination string) error {
	sourcePath := filepath.Join(source, "go.mod")
	destinationPath := filepath.Join(destination, "go.mod")

	sourceModule, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer sourceModule.Close()

	destinationModule, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("could not create mod file: %w", err)
	}
	defer destinationModule.Close()

	if _, err := io.Copy(sourceModule, destinationModule); err != nil {
		return fmt.Errorf("could not copy module contents: %w", err)
	}

	return nil
}
