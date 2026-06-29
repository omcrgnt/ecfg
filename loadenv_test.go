package ecfg_test

import (
	"fmt"
	"testing"

	"github.com/omcrgnt/ecfg"
	"github.com/omcrgnt/res/unique"
)

type listLen int

func (listLen) Usage() string { return "max list len" }

func (l listLen) Validate() error {
	if l < 0 {
		return fmt.Errorf("max list len must be >= 0")
	}
	return nil
}

type itemConfig struct {
	MaxListLen listLen
}

type hostLeaf string

func (hostLeaf) Usage() string { return "host" }

func (hostLeaf) Validate() error { return nil }

type httpConfig struct {
	Host hostLeaf
}

func TestLoadEnv_happy(t *testing.T) {
	t.Cleanup(ecfg.ResetForTest)
	ecfg.SetPrefix("TEST")

	t.Setenv("TEST_SERVICE_ITEM_MAX_LIST_LEN", "100")

	reg := unique.New()
	cfg := &itemConfig{}
	if err := reg.AddWithCustomTag(cfg, ecfg.TagKey(), "SERVICE_ITEM"); err != nil {
		t.Fatal(err)
	}

	if err := ecfg.LoadEnv(reg); err != nil {
		t.Fatal(err)
	}
	if cfg.MaxListLen != 100 {
		t.Fatalf("got %d", cfg.MaxListLen)
	}
}

func TestLoadEnv_defaultPrefix(t *testing.T) {
	t.Cleanup(ecfg.ResetForTest)

	t.Setenv("APP_SERVICE_ITEM_MAX_LIST_LEN", "42")

	reg := unique.New()
	cfg := &itemConfig{}
	if err := reg.AddWithCustomTag(cfg, ecfg.TagKey(), "SERVICE_ITEM"); err != nil {
		t.Fatal(err)
	}

	if err := ecfg.LoadEnv(reg); err != nil {
		t.Fatal(err)
	}
	if cfg.MaxListLen != 42 {
		t.Fatalf("got %d", cfg.MaxListLen)
	}
}

func TestLoadEnv_missingCustomTag(t *testing.T) {
	t.Cleanup(ecfg.ResetForTest)

	reg := unique.New()
	if err := reg.Add(&itemConfig{}); err != nil {
		t.Fatal(err)
	}

	if err := ecfg.LoadEnv(reg); err == nil {
		t.Fatal("expected error for missing custom tag")
	}
}

func TestLoadEnv_nilRegistry(t *testing.T) {
	if err := ecfg.LoadEnv(nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadEnv_customTagKey(t *testing.T) {
	t.Cleanup(ecfg.ResetForTest)
	ecfg.SetTagKey("cfg")
	ecfg.SetPrefix("X")

	t.Setenv("X_BLOCK_HOST", "127.0.0.1")

	reg := unique.New()
	cfg := &httpConfig{}
	if err := reg.AddWithCustomTag(cfg, "cfg", "BLOCK"); err != nil {
		t.Fatal(err)
	}

	if err := ecfg.LoadEnv(reg); err != nil {
		t.Fatal(err)
	}
	if cfg.Host != "127.0.0.1" {
		t.Fatalf("got %q", cfg.Host)
	}
}

func TestLoadEnv_missingEnv(t *testing.T) {
	t.Cleanup(ecfg.ResetForTest)

	reg := unique.New()
	if err := reg.AddWithCustomTag(&itemConfig{}, ecfg.TagKey(), "SERVICE_ITEM"); err != nil {
		t.Fatal(err)
	}

	if err := ecfg.LoadEnv(reg); err == nil {
		t.Fatal("expected missing env error")
	}
}
