package commands

import (
	"github.com/hugocarreira/go-decent-copy"
	"github.com/plexsystems/pacmod/pkg/pack"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func newPackCommand() *cobra.Command {
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
		return err
	}

	name := args[0]
	version := args[1]

	module := pack.Module{
		Path:    path,
		Name:    name,
		Version: version,
	}

	outputDirectory := filepath.Join(module.Path, module.Version)
	err = os.Mkdir(outputDirectory, 0777)
	if err != nil {
		return err
	}

	log.Printf("Packing module %s into output directory %s", module.Name, outputDirectory)

	log.Println("Creating archive...")
	err = module.ZipModule(outputDirectory)
	if err != nil {
		return err
	}

	log.Println("Creating info...")
	info := pack.Info{
		Version: version,
	}
	err = info.CreateInfo(outputDirectory)
	if err != nil {
		return err
	}

	log.Println("Copying mod...")
	return copyModuleFile(module.Path, outputDirectory)
}

func copyModuleFile(source string, destination string) error {
	sourceModFile := filepath.Join(source, "go.mod")
	destinationModFile := filepath.Join(destination, "go.mod")

	err := decentcopy.Copy(sourceModFile, destinationModFile)
	if err != nil {
		return err
	}

	return nil
}
