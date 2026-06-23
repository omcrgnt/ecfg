package gen_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/omcrgnt/ecfg/internal/gen"
)

func TestRun_appResources(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.template")
	if err := gen.Run("AppResourcesFixture", "github.com/omcrgnt/ecfg/internal/testdata", "FIX", gen.Options{
		TemplatePath: out,
		MarkdownPath: filepath.Join(dir, "env.md"),
	}); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "FIX_SERVER_HTTP_ITEM_LABEL=") {
		t.Fatalf("body: %s", body)
	}
}
