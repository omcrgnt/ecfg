package usage

import (
	"errors"
	"testing"
)

func TestGoUsageFromAST_label(t *testing.T) {
	text, err := GoUsageFromAST("github.com/omcrgnt/ecfg/internal/testdata", "Label")
	if err != nil {
		t.Fatal(err)
	}
	if text == "" {
		t.Fatal("expected usage text")
	}
}

func TestGoUsageFromAST_missingType(t *testing.T) {
	_, err := GoUsageFromAST("github.com/omcrgnt/ecfg/internal/testdata", "NoSuchType")
	if !errors.Is(err, errMissingUsage) {
		t.Fatalf("got %v", err)
	}
}

func TestGoUsageFromAST_packageNotFound(t *testing.T) {
	_, err := GoUsageFromAST("github.com/no/such/package/xyz", "Label")
	if err == nil {
		t.Fatal("expected error")
	}
}
