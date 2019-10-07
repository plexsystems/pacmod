package pack

import "testing"

func Test_GetZipPath_PathAndModulePathsAreSame(t *testing.T) {
	module := Module{
		Path:    "/root/",
		Name:    "root",
		Version: "v1.0.0",
	}

	actual := module.getZipPath("/root/app.go")
	expected := "root@v1.0.0/app.go"

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func Test_GetZipPath_ModulePathChildOfPath(t *testing.T) {
	module := Module{
		Path:    "/root/repository/username/app",
		Name:    "repository/username/app",
		Version: "v1.0.0",
	}
	actual := module.getZipPath("/root/repository/username/app/app.go")
	expected := "repository/username/app@v1.0.0/app.go"

	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
