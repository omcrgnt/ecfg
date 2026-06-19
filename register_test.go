package ecfg_test

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/builder"
	"github.com/omcrgnt/ecfg"
	"github.com/omcrgnt/res"
)

type okModule struct{}

func (okModule) Build() (any, error) { return "built", nil }

func TestRegister_success(t *testing.T) {
	reg := res.New()
	cfg := struct {
		Skip int
		A    okModule
	}{
		A: okModule{},
	}

	if err := ecfg.Register(cfg, reg); err != nil {
		t.Fatal(err)
	}

	n := 0
	reg.WalkEntries(func(_ res.Entry) bool {
		n++
		return true
	})
	if n != 1 {
		t.Fatalf("expected 1 entry, got %d", n)
	}
}

func TestRegister_skipsNilPointer(t *testing.T) {
	reg := res.New()
	cfg := struct {
		M *okModule
		A okModule
	}{
		A: okModule{},
	}

	if err := ecfg.Register(cfg, reg); err != nil {
		t.Fatal(err)
	}

	n := 0
	reg.WalkEntries(func(_ res.Entry) bool {
		n++
		return true
	})
	if n != 1 {
		t.Fatalf("expected 1 entry, got %d", n)
	}
}

func TestRegister_skipsNonBuilder(t *testing.T) {
	reg := res.New()
	cfg := struct {
		Name string
	}{Name: "x"}

	if err := ecfg.Register(cfg, reg); err != nil {
		t.Fatal(err)
	}

	n := 0
	reg.WalkEntries(func(_ res.Entry) bool {
		n++
		return true
	})
	if n != 0 {
		t.Fatalf("expected 0 entries, got %d", n)
	}
}

func TestRegister_nilRegistry(t *testing.T) {
	err := ecfg.Register(struct{}{}, nil)
	if err == nil || err.Error() != "ecfg: nil registry" {
		t.Fatalf("expected nil registry error, got %v", err)
	}
}

func TestRegister_nilConfig(t *testing.T) {
	err := ecfg.Register(nil, res.New())
	if err == nil || err.Error() != "ecfg: nil config" {
		t.Fatalf("expected nil config error, got %v", err)
	}
}

func TestRegister_withBuilderBuild(t *testing.T) {
	reg := res.New()
	cfg := struct {
		A okModule
	}{A: okModule{}}

	if err := ecfg.Register(cfg, reg); err != nil {
		t.Fatal(err)
	}
	if err := builder.Build(reg); err != nil {
		t.Fatal(err)
	}

	got, err := reg.GetOneByType(reflect.TypeOf(""))
	if err != nil {
		t.Fatal(err)
	}
	if got != "built" {
		t.Fatalf("got %v", got)
	}
}

var _ builder.Builder = okModule{}
