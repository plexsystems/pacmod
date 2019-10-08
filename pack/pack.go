package pack

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Module packs the module at the given path and version then
// outputs the result to the specified output directory
func Module(path string, version string, outputDirectory string) error {
	moduleName, err := getModuleName(path)
	if err != nil {
		return fmt.Errorf("could not get module name: %w", err)
	}

	if err := createZipArchive(path, moduleName, version, outputDirectory); err != nil {
		return fmt.Errorf("could not create zip archive: %w", err)
	}

	if err := createInfoFile(version, outputDirectory); err != nil {
		return fmt.Errorf("could not create info file: %w", err)
	}

	if err := copyModuleFile(path, outputDirectory); err != nil {
		return fmt.Errorf("could not copy module file: %w", err)
	}

	return nil
}

func getModuleName(path string) (string, error) {
	moduleFilePath := filepath.Join(path, "go.mod")

	file, err := os.Open(moduleFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to open module file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if ok := scanner.Scan(); !ok {
		return "", fmt.Errorf("unable to get module header: %w", err)
	}

	moduleFileHeader := strings.Split(scanner.Text(), " ")
	if len(moduleFileHeader) <= 1 {
		return "", fmt.Errorf("unable to parse module header: %w", err)
	}

	return moduleFileHeader[1], nil
}

func createZipArchive(path string, moduleName string, version string, outputDirectory string) error {
	outputPath := filepath.Join(outputDirectory, version+".zip")

	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create empty zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(path, func(currentFilePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("unable to walk path: %w", err)
		}

		if skipFile, err := shouldSkipFile(fileInfo); skipFile {
			return err
		}

		zipPath := getZipPath(path, currentFilePath, moduleName, version)
		zipFileWriter, err := zipWriter.Create(zipPath)
		if err != nil {
			return fmt.Errorf("unable to add file to zip archive: %w", err)
		}

		file, err := os.Open(currentFilePath)
		if err != nil {
			return fmt.Errorf("unable to open file: %w", err)
		}
		defer file.Close()

		if _, err := io.Copy(zipFileWriter, file); err != nil {
			return fmt.Errorf("unable to copy file to zip archive: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to zip all files: %w", err)
	}

	return nil
}

func shouldSkipFile(fileInfo os.FileInfo) (bool, error) {
	// We do not want to include the .git directory in the archived module file
	// filepath.SkipDir tells the Walk() function to ignore everything inside of the directory
	if fileInfo.IsDir() && fileInfo.Name() == ".git" {
		return true, filepath.SkipDir
	}

	// Do not process directories or zip files
	// returning nil tells the Walk() function to ignore this file
	if fileInfo.IsDir() || filepath.Ext(fileInfo.Name()) == ".zip" {
		return true, nil
	}

	return false, nil
}

func getZipPath(modulePath string, currentFilePath string, moduleName string, version string) string {
	filePath := strings.TrimPrefix(currentFilePath, modulePath)
	return filepath.Join(fmt.Sprintf("%s@%s", moduleName, version), filePath)
}

func createInfoFile(version string, outputDirectory string) error {
	infoFilePath := filepath.Join(outputDirectory, version+".info")
	file, err := os.Create(infoFilePath)
	if err != nil {
		return fmt.Errorf("could not create info file: %w", err)
	}
	defer file.Close()

	type infoFile struct {
		Version string
		Time    string
	}

	currentTime := getInfoFileFormattedTime(time.Now())
	info := infoFile{
		Version: version,
		Time:    currentTime,
	}

	infoBytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("could not marshal info file: %w", err)
	}

	if _, err := file.Write(infoBytes); err != nil {
		return fmt.Errorf("could not write info file: %w", err)
	}

	return nil
}

func getInfoFileFormattedTime(currentTime time.Time) string {
	const infoFileTimeFormat = "2006-01-02T15:04:05Z"
	return currentTime.Format(infoFileTimeFormat)
}

func copyModuleFile(modulePath string, outputDirectory string) error {
	sourcePath := filepath.Join(modulePath, "go.mod")
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
