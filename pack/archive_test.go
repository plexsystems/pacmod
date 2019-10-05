package pack

import (
	"testing"
)

func TestGetZipPath(t *testing.T) {
	module := Module{
		Path:    "/linux",
		Name:    "modulename",
		Version: "v1.0.0",
	}

	actual := module.getZipPath("/linux/file.go")
	expected := "modulename@v1.0.0\\file.go"

	if actual != expected {
		t.Errorf("Linux: expected %v, got %v", expected, actual)
	}

	module = Module{
		Path:    "C:\\windows",
		Name:    "modulename",
		Version: "v1.0.0",
	}

	actual = module.getZipPath("C:\\windows\\file.go")
	expected = "modulename@v1.0.0\\file.go"

	if actual != expected {
		t.Errorf("Windows expected %v, got %v", expected, actual)
	}
}
