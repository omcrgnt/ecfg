package walker

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/ecfg/internal/teststruct"
)

func TestReflectProvider(t *testing.T) {
	cfg := teststruct.Config{}
	p, err := NewReflectProvider(cfg)
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	fields, _ := p.GetFields()
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}

	f := fields[0]
	if f.Name() != "Port" || f.Tag("ecfg") != "PORT" || f.Tag("usage") != "Server port" {
		t.Errorf("unexpected field data: %+v", f)
	}

	if f.Kind() != reflect.Int {
		t.Errorf("expected int kind, got %v", f.Kind())
	}
}
