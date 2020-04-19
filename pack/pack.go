package pack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/zip"
)

// Module packs the module at the given path and version then
// outputs the result to the specified output directory
func Module(path string, version string, outputDirectory string) error {
	moduleFile, err := getModuleFile(path, version)
	if err != nil {
		return fmt.Errorf("get module file: %w", err)
	}

	if err := createZipArchive(path, moduleFile, outputDirectory); err != nil {
		return fmt.Errorf("create zip archive: %w", err)
	}

	if err := createInfoFile(moduleFile, outputDirectory); err != nil {
		return fmt.Errorf("create info file: %w", err)
	}

	if err := copyModuleFile(path, moduleFile, outputDirectory); err != nil {
		return fmt.Errorf("copy module file: %w", err)
	}

	return nil
}

func getModuleFile(path string, version string) (*modfile.File, error) {
	path = filepath.Join(path, "go.mod")
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open module file: %w", err)
	}
	defer file.Close()

	moduleBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read module file: %w", err)
	}

	moduleFile, err := modfile.Parse(path, moduleBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("parse module file: %w", err)
	}

	if moduleFile.Module == nil {
		return nil, fmt.Errorf("parsing module returned nil module")
	}

	moduleFile.Module.Mod.Version = version

	return moduleFile, nil
}

func createZipArchive(path string, moduleFile *modfile.File, outputDirectory string) error {
	outputPath := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".zip")

	var zipContents bytes.Buffer
	if err := zip.CreateFromDir(&zipContents, moduleFile.Module.Mod, path); err != nil {
		return fmt.Errorf("create zip from dir: %w", err)
	}

	if err := ioutil.WriteFile(outputPath, zipContents.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing zip file: %w", err)
	}

	return nil
}

func createInfoFile(moduleFile *modfile.File, outputDirectory string) error {
	infoFilePath := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".info")
	file, err := os.Create(infoFilePath)
	if err != nil {
		return fmt.Errorf("create info file: %w", err)
	}
	defer file.Close()

	type infoFile struct {
		Version string
		Time    string
	}

	currentTime := getInfoFileFormattedTime(time.Now())
	info := infoFile{
		Version: moduleFile.Module.Mod.Version,
		Time:    currentTime,
	}

	infoBytes, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal info file: %w", err)
	}

	if _, err := file.Write(infoBytes); err != nil {
		return fmt.Errorf("write info file: %w", err)
	}

	return nil
}

func getInfoFileFormattedTime(currentTime time.Time) string {
	const infoFileTimeFormat = "2006-01-02T15:04:05Z"
	return currentTime.Format(infoFileTimeFormat)
}

func copyModuleFile(path string, moduleFile *modfile.File, outputDirectory string) error {
	if outputDirectory == "." {
		return nil
	}

	sourcePath := filepath.Join(path, "go.mod")
	destinationPath := filepath.Join(outputDirectory, moduleFile.Module.Mod.Version+".mod")

	if sourcePath == destinationPath {
		return nil
	}

	moduleContents, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("read module file: %w", err)
	}

	if err := ioutil.WriteFile(destinationPath, moduleContents, 0644); err != nil {
		return fmt.Errorf("write module file: %w", err)
	}

	return nil
}
