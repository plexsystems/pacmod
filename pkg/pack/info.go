package pack

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Info represents the JSON fields in the info file
type Info struct {
	Version string `json:"Version"`
	time    string `json:"Time"`
}

// CreateInfo creates an info file associated with a module
func (i *Info) CreateInfo(outPath string) error {
	infoFilePath := filepath.Join(outPath, i.Version+".info")
	file, err := os.Create(infoFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	infoBytes, err := json.Marshal(struct {
		Version string
		Time    string
	}{
		Version: i.Version,
		Time:    time.Now().Format("2006-01-02T15:04:05Z"),
	})
	if err != nil {
		return err
	}

	_, err = file.Write(infoBytes)
	if err != nil {
		return err
	}

	return nil
}
