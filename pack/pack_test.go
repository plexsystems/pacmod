package pack

import (
	"testing"
	"time"
)

func Test_GetInfoFile_ReturnsCorrectInfo(t *testing.T) {
	goLaunchDate := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	actual := getInfoFileFormattedTime(goLaunchDate)
	expected := "2009-11-10T23:00:00Z"

	if expected != actual {
		t.Errorf("invalid infofile time format: expected %v actual %v", expected, actual)
	}
}

func Test_GetZipPath_PathAndModulePathsAreSame(t *testing.T) {
	modulePath := "/root/"
	currentFilePath := "/root/app.go"
	name := "root"
	version := "v1.0.0"

	actual := getZipPath(modulePath, currentFilePath, name, version)
	expected := "root@v1.0.0/app.go"

	if expected != actual {
		t.Errorf("invalid zip path: expected %v actual %v", expected, actual)
	}
}

func Test_GetZipPath_ModulePathIsChildOfPath(t *testing.T) {
	modulePath := "/root/repository/username/app"
	currentFilePath := "/root/repository/username/app/app.go"
	name := "repository/username/app"
	version := "v1.0.0"

	actual := getZipPath(modulePath, currentFilePath, name, version)
	expected := "repository/username/app@v1.0.0/app.go"

	if expected != actual {
		t.Errorf("invalid zip path: expected %v actual %v", expected, actual)
	}
}
