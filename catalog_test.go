package ecfg_test

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/ecfg"
)

type catalogFixture struct {
	Block string `ecfg:"service_item"`
}

type catalogFixtureCfg struct {
	Block string `cfg:"SERVICE_ITEM"`
}

func TestCatalogSegment_defaultTagKey(t *testing.T) {
	ecfg.ResetForTest()
	t.Cleanup(ecfg.ResetForTest)

	sf := reflect.TypeOf(catalogFixture{}).Field(0)

	seg, err := ecfg.CatalogSegment(sf)
	if err != nil {
		t.Fatal(err)
	}
	if seg != "SERVICE_ITEM" {
		t.Fatalf("segment: got %q", seg)
	}
}

func TestCatalogSegment_customTagKey(t *testing.T) {
	ecfg.ResetForTest()
	t.Cleanup(ecfg.ResetForTest)

	ecfg.SetTagKey("cfg")

	sf := reflect.TypeOf(catalogFixtureCfg{}).Field(0)

	seg, err := ecfg.CatalogSegment(sf)
	if err != nil {
		t.Fatal(err)
	}
	if seg != "SERVICE_ITEM" {
		t.Fatalf("segment: got %q", seg)
	}
}

func TestCatalogSegment_missingTag(t *testing.T) {
	ecfg.ResetForTest()
	t.Cleanup(ecfg.ResetForTest)

	sf := reflect.StructField{Name: "X"}
	if _, err := ecfg.CatalogSegment(sf); err == nil {
		t.Fatal("expected error")
	}
}

func TestCatalogSegment_emptyTag(t *testing.T) {
	ecfg.ResetForTest()
	t.Cleanup(ecfg.ResetForTest)

	type emptyTag struct {
		Block string `ecfg:""`
	}
	sf := reflect.TypeOf(emptyTag{}).Field(0)

	if _, err := ecfg.CatalogSegment(sf); err == nil {
		t.Fatal("expected error")
	}
}

func TestCatalogSegment_ignoresTagOptionsAfterComma(t *testing.T) {
	ecfg.ResetForTest()
	t.Cleanup(ecfg.ResetForTest)

	type withOptions struct {
		Block string `ecfg:"APP,inline"`
	}
	sf := reflect.TypeOf(withOptions{}).Field(0)

	seg, err := ecfg.CatalogSegment(sf)
	if err != nil {
		t.Fatal(err)
	}
	if seg != "APP" {
		t.Fatalf("segment: got %q", seg)
	}
}
