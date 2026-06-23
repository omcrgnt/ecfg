package ecfgtool

import (
	"testing"
)

func TestCollectTemplateEntries_appResources(t *testing.T) {
	const (
		pkgPath  = "github.com/omcrgnt/ecfg/internal/testdata"
		typeName = "AppResourcesFixture"
		prefix   = "FIX"
	)
	entries, err := CollectTemplateEntries(pkgPath, typeName, prefix)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 5 {
		t.Fatalf("entries: %+v", entries)
	}
	keys := map[string]bool{}
	for _, e := range entries {
		keys[e.EnvKey] = true
		if e.Usage == "" {
			t.Fatalf("empty usage for %s", e.EnvKey)
		}
	}
	for _, want := range []string{
		"FIX_APP_SHUTDOWN_TIMEOUT",
		"FIX_SERVICE_ITEM_MAX_LIST_LEN",
		"FIX_SERVER_HTTP_ITEM_LABEL",
		"FIX_SERVER_HTTP_ITEM_HOST",
		"FIX_SERVER_HTTP_ITEM_PORT",
	} {
		if !keys[want] {
			t.Fatalf("missing %s in %+v", want, entries)
		}
	}
}
