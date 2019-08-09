package pack

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Module represents a Go module
type Module struct {
	Path    string
	Name    string
	Version string
}

// ZipModule zips the module and outputs the result to the specified output path
func (m *Module) ZipModule(outPath string) error {
	zipOutputPath := filepath.Join(outPath, m.Version+".zip")

	zipFile, err := os.Create(zipOutputPath)
	if err != nil {
		return err
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
		return err
	}

	return zipWriter.Close()
}

func (m *Module) getZipPath(currentFilePath string) string {
	fileName := strings.TrimPrefix(currentFilePath, m.Path)
	moduleName := fmt.Sprintf("%s@%s", m.Name, m.Version)

	return filepath.Join(moduleName, fileName)
}
