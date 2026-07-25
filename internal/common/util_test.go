package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupConfigParsesBgFlag(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempHome, "xdg"))

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"km", "start", "myenv", "--bg"}

	_, ops, err := SetupConfig()
	if err != nil {
		t.Fatalf("SetupConfig returned error: %v", err)
	}
	if !ops.Bg {
		t.Fatalf("expected ops.Bg to be true when --bg is passed, got false")
	}
}

func TestSetupConfigDefaultsBgFalse(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempHome, "xdg"))

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"km", "start", "myenv"}

	_, ops, err := SetupConfig()
	if err != nil {
		t.Fatalf("SetupConfig returned error: %v", err)
	}
	if ops.Bg {
		t.Fatalf("expected ops.Bg to be false when --bg is not passed, got true")
	}
}
