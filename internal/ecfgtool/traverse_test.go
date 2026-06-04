package ecfgtool

import (
	"testing"

	"github.com/omcrgnt/ecfg/internal/testdata"
	"github.com/omcrgnt/ecfg/pkg/walk"
)

func TestTraverse_protoSkipDescend(t *testing.T) {
	eng, err := walk.NewEngineReflect(&testdata.ProtoCfg{})
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	err = traverse(eng, Options{Prefix: "T"}, func(ctx visitCtx) error {
		names = append(names, ctx.walk.Field.Name)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Server", "Port"}
	if len(names) != len(want) {
		t.Fatalf("visited %v want %v", names, want)
	}
	for i, n := range want {
		if names[i] != n {
			t.Fatalf("visited %v want %v", names, want)
		}
	}
}

func TestTraverse_envKeyPath(t *testing.T) {
	eng, err := walk.NewEngineReflect(&testdata.AppConfig{})
	if err != nil {
		t.Fatal(err)
	}
	var key string
	err = traverse(eng, Options{Prefix: "APP"}, func(ctx visitCtx) error {
		if ctx.isLeaf {
			key = ctx.envKey
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if key != "APP_SERVER_LABEL" {
		t.Fatalf("got %q", key)
	}
}

func TestTraverse_initNilPointerBlock(t *testing.T) {
	eng, err := walk.NewEngineReflect(&testdata.PtrAppConfig{})
	if err != nil {
		t.Fatal(err)
	}
	var key string
	err = traverse(eng, Options{Prefix: "P"}, func(ctx visitCtx) error {
		if ctx.isLeaf {
			key = ctx.envKey
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if key != "P_SERVER_LABEL" {
		t.Fatalf("got %q", key)
	}
}
