package walker

import (
	"testing"

	"github.com/omcrgnt/ecfg/internal/teststruct"
)

func TestProvidersParity(t *testing.T) {
	pkgPath := "://github.com"
	structName := "Config"

	// 1. Reflect
	rp, _ := NewReflectProvider(teststruct.Config{})
	rFields, _ := rp.GetFields()

	// 2. Types (AST)
	tp, _ := NewTypesProvider(pkgPath, structName)
	tFields, _ := tp.GetFields()

	if len(rFields) != len(tFields) {
		t.Fatalf("parity error: reflect found %d fields, types found %d", len(rFields), len(tFields))
	}

	for i := range rFields {
		if rFields[i].Name() != tFields[i].Name() {
			t.Errorf("parity error at field %d: %s != %s", i, rFields[i].Name(), tFields[i].Name())
		}
		if rFields[i].Tag("ecfg") != tFields[i].Tag("ecfg") {
			t.Errorf("parity error in tags for %s", rFields[i].Name())
		}
	}
}
