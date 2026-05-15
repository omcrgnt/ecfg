package walker

import (
	"testing"

	"github.com/omcrgnt/ecfg/internal/teststruct"
)

func TestProvidersParity(t *testing.T) {
	pkgPath := "github.com/omcrgnt/ecfg/internal/teststruct"
	structName := "Config"

	rp, err := NewReflectProvider(teststruct.Config{})
	if err != nil {
		t.Fatalf("reflect provider: %v", err)
	}
	rFields, err := rp.GetFields()
	if err != nil {
		t.Fatalf("reflect fields: %v", err)
	}

	tp, err := NewTypesProvider(pkgPath, structName)
	if err != nil {
		t.Fatalf("types provider: %v", err)
	}
	tFields, err := tp.GetFields()
	if err != nil {
		t.Fatalf("types fields: %v", err)
	}

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

func TestWalkProvider_Nested(t *testing.T) {
	type Inner struct {
		Key string `ecfg:"KEY"`
	}
	type Outer struct {
		Module Inner `ecfg:"MODULE"`
	}

	p, err := NewReflectProvider(Outer{})
	if err != nil {
		t.Fatal(err)
	}

	var names []string
	err = WalkProvider(p, func(f Field) error {
		names = append(names, f.Name())
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"Module", "Key"}
	if len(names) != len(want) {
		t.Fatalf("got %d fields %v, want %v", len(names), names, want)
	}
	for i, n := range want {
		if names[i] != n {
			t.Errorf("field %d: got %s, want %s", i, names[i], n)
		}
	}
}

func TestReflectProvider_PointerStruct(t *testing.T) {
	type Extra struct {
		ID int `ecfg:"ID"`
	}
	type Config struct {
		DB *Extra `ecfg:"DB"`
	}

	p, err := NewReflectProvider(Config{})
	if err != nil {
		t.Fatal(err)
	}
	fields, err := p.GetFields()
	if err != nil {
		t.Fatal(err)
	}
	sub, err := fields[0].GetProvider()
	if err != nil {
		t.Fatalf("GetProvider for *Extra: %v", err)
	}
	if sub.EntryName() != "Extra" {
		t.Errorf("expected Extra, got %s", sub.EntryName())
	}
}
