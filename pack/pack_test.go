package pack

import (
	"testing"
	"time"
)

func Test_GetInfoFile_ReturnsCorrectTimeFormat(t *testing.T) {
	goLaunchDate := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	actual := getInfoFileFormattedTime(goLaunchDate)
	expected := "2009-11-10T23:00:00Z"

	if expected != actual {
		t.Errorf("invalid infofile time format: expected %v actual %v", expected, actual)
	}
}
