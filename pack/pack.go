package pack

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Module represents a Go module
type Module struct {
	Path    string
	Name    string
	Version string
}

// PackageModule packs the module and outputs the result to the specified output path
func (m Module) PackageModule(outputDirectory string) error {
	if err := m.createZipArchive(outputDirectory); err != nil {
		return fmt.Errorf("could not create zip archive: %w", err)
	}

	if err := m.createInfoFile(outputDirectory); err != nil {
		return fmt.Errorf("could not create info file: %w", err)
	}

	if err := m.copyModuleFile(outputDirectory); err != nil {
		return fmt.Errorf("could not copy module file: %w", err)
	}

	return nil
}

func (m Module) createZipArchive(outputDirectory string) error {
	outputPath := filepath.Join(outputDirectory, m.Version+".zip")

	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create empty zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	err = filepath.Walk(m.Path, func(currentFilePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.IsDir() && fileInfo.Name() == ".git" {
			return filepath.SkipDir
		}

		if fileInfo.IsDir() || filepath.Ext(currentFilePath) == ".zip" {
			return nil
		}

		file, err := os.Open(currentFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		zipPath := m.getZipPath(currentFilePath)
		zipFileWriter, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFileWriter, file)
		return err
	})
	if err != nil {
		return fmt.Errorf("could not zip all files: %w", err)
	}

	return zipWriter.Close()
}

func (m Module) getZipPath(currentFilePath string) string {
	fileName := strings.TrimPrefix(currentFilePath, m.Path)
	moduleName := fmt.Sprintf("%s@%s", m.Name, m.Version)

	return filepath.Join(moduleName, fileName)
}

func (m Module) createInfoFile(outputDirectory string) error {
	infoFilePath := filepath.Join(outputDirectory, m.Version+".info")
	file, err := os.Create(infoFilePath)
	if err != nil {
		return fmt.Errorf("could not create info file: %w", err)
	}
	defer file.Close()

	infoBytes, err := json.Marshal(struct {
		Version string
		Time    string
	}{
		Version: m.Version,
		Time:    time.Now().Format("2006-01-02T15:04:05Z"),
	})
	if err != nil {
		return fmt.Errorf("could not marshal info file: %w", err)
	}

	if _, err := file.Write(infoBytes); err != nil {
		return fmt.Errorf("could not write info file: %w", err)
	}

	return nil
}

func (m Module) copyModuleFile(outputDirectory string) error {
	sourcePath := filepath.Join(m.Path, "go.mod")
	destinationPath := filepath.Join(outputDirectory, "go.mod")

	sourceModule, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("could not open source module file: %w", err)
	}
	defer sourceModule.Close()

	destinationModule, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("could not create mod file: %w", err)
	}
	defer destinationModule.Close()

	if _, err := io.Copy(destinationModule, sourceModule); err != nil {
		return fmt.Errorf("could not copy module contents: %w", err)
	}

	return nil
}
