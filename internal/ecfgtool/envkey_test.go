package ecfgtool

import (
	"errors"
	"testing"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

func TestFieldNameToSegment(t *testing.T) {
	tests := map[string]string{
		"HTTPPort":   "HTTP_PORT",
		"HTTPServer": "HTTP_SERVER",
		"ID":         "ID",
		"Label":      "LABEL",
		"Server":     "SERVER",
		"DbHost":     "DB_HOST",
	}
	for in, want := range tests {
		if got := fieldNameToSegment(in); got != want {
			t.Errorf("%s: got %q want %q", in, got, want)
		}
	}
}

func TestSegment_rootRequiresTag(t *testing.T) {
	_, err := segment(0, "", "Server")
	if !errors.Is(err, ErrMissingEcfgTag) {
		t.Fatalf("got %v", err)
	}
}

func TestKeyRegistry_duplicate(t *testing.T) {
	r := newKeyRegistry()
	if err := r.add("A", "f1"); err != nil {
		t.Fatal(err)
	}
	if err := r.add("A", "f2"); !errors.Is(err, ErrDuplicateEnvKey) {
		t.Fatalf("got %v", err)
	}
}

func TestIsStructField_typesEngine(t *testing.T) {
	eng, err := walk.NewEngineTypes("github.com/omcrgnt/ecfg/internal/testdata", "AppConfig")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := eng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if !isStructField(fields[0]) {
		t.Fatal("Server should be struct")
	}
	child, err := eng.Child(fields[0])
	if err != nil {
		t.Fatal(err)
	}
	leafFields, err := child.Fields()
	if err != nil {
		t.Fatal(err)
	}
	if isStructField(leafFields[0]) {
		t.Fatal("Label should not be struct")
	}
}

func TestJoin_prefix(t *testing.T) {
	if got := join("app", "SERVER", "LABEL"); got != "APP_SERVER_LABEL" {
		t.Fatalf("got %q", got)
	}
	if got := join("", "SERVER", "LABEL"); got != "SERVER_LABEL" {
		t.Fatalf("got %q", got)
	}
}
