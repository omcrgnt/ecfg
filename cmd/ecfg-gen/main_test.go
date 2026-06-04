package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_missingRequiredFlags(t *testing.T) {
	if code := run([]string{"-type", "AppConfig"}); code != 2 {
		t.Fatalf("exit code: %d", code)
	}
}

func TestRun_success(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.template")
	code := run([]string{
		"-type", "AppConfig",
		"-pkg", "github.com/omcrgnt/ecfg/internal/testdata",
		"-prefix", "CLI",
		"-o", out,
	})
	if code != 0 {
		t.Fatalf("exit code: %d", code)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "CLI_SERVER_LABEL=") {
		t.Fatalf("body: %s", body)
	}
}
