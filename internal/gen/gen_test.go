package gen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/omcrgnt/ecfg/internal/gen"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.template")
	if err := gen.Run("AppConfig", "github.com/omcrgnt/ecfg/internal/testdata", "GEN", out); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "GEN_SERVER_LABEL=") {
		t.Fatalf("body: %s", body)
	}
}
