package pack

import (
	"os"
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

func Test_ShouldSkipFile(t *testing.T) {
	testCases := []struct {
		name     string
		isDir    bool
		expected bool
	}{
		{"", true, true},
		{".git", true, true},
		{".zip", false, true},
	}

	for _, testCase := range testCases {
		infoFile := fakeInfoFile{
			name:  testCase.name,
			isDir: testCase.isDir,
		}

		actual, _ := shouldSkipFile(infoFile)
		expected := testCase.expected

		if expected != actual {
			t.Errorf("file %v was not skipped. expected: %v actual: %v", infoFile, expected, actual)
		}
	}
}

type fakeInfoFile struct {
	name  string
	isDir bool
}

var _ os.FileInfo = fakeInfoFile{}

func (f fakeInfoFile) Name() string {
	return f.name
}

func (f fakeInfoFile) IsDir() bool {
	return f.isDir
}

func (f fakeInfoFile) Mode() os.FileMode {
	return 0
}

func (f fakeInfoFile) ModTime() time.Time {
	return time.Now()
}

func (f fakeInfoFile) Sys() interface{} {
	return nil
}

func (f fakeInfoFile) Size() int64 {
	return 0
}
