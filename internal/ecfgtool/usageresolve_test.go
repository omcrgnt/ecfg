package ecfgtool

import (
	"errors"
	"reflect"
	"testing"

	"github.com/omcrgnt/ecfg/internal/testdata"
	"github.com/omcrgnt/ecfg/pkg/walk"
	commonv1 "github.com/omcrgnt/proto/gen/go/common/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestResolveUsage_goLeaf(t *testing.T) {
	text, err := resolveUsage(walk.FieldDesc{ReflectType: reflect.TypeOf(testdata.Label(""))}, false)
	if err != nil || text == "" {
		t.Fatalf("got %q %v", text, err)
	}
}

func TestResolveUsage_goIncompatible(t *testing.T) {
	_, err := resolveUsage(walk.FieldDesc{ReflectType: reflect.TypeOf(0)}, false)
	if !errors.Is(err, ErrIncompatibleLeaf) {
		t.Fatalf("got %v", err)
	}
}

func TestResolveUsage_emptyUsage(t *testing.T) {
	_, err := resolveUsage(walk.FieldDesc{ReflectType: reflect.TypeOf(testdata.EmptyUsageLabel(""))}, false)
	if !errors.Is(err, ErrEmptyUsage) {
		t.Fatalf("got %v want ErrEmptyUsage", err)
	}
}

func TestResolveUsage_proto(t *testing.T) {
	text, err := resolveUsage(walk.FieldDesc{ReflectType: reflect.TypeOf(&commonv1.Port{})}, true)
	if err != nil || text == "" {
		t.Fatalf("got %q %v", text, err)
	}
}

func TestResolveUsageInput_missing(t *testing.T) {
	_, err := resolveUsageInput(usageInput{})
	if !errors.Is(err, ErrMissingUsage) {
		t.Fatalf("got %v", err)
	}
}

func TestResolveUsage_typesEngineLeaf(t *testing.T) {
	eng, err := walk.NewEngineTypes("github.com/omcrgnt/ecfg/internal/testdata", "AppConfig")
	if err != nil {
		t.Fatal(err)
	}
	fields, err := eng.Fields()
	if err != nil {
		t.Fatal(err)
	}
	child, err := eng.Child(fields[0])
	if err != nil {
		t.Fatal(err)
	}
	leafFields, err := child.Fields()
	if err != nil {
		t.Fatal(err)
	}
	text, err := resolveUsage(leafFields[0], false)
	if err != nil || text == "" {
		t.Fatalf("got %q %v", text, err)
	}
}

func TestResolveUsage_protoDoublePointerReflect(t *testing.T) {
	var p **commonv1.Port
	text, err := resolveUsage(walk.FieldDesc{ReflectType: reflect.TypeOf(p)}, true)
	if err != nil || text == "" {
		t.Fatalf("got %q %v", text, err)
	}
}

func TestResolveUsage_protoMissingExtension(t *testing.T) {
	_, err := resolveUsage(walk.FieldDesc{
		ReflectType: reflect.TypeOf(&wrapperspb.UInt32Value{}),
	}, true)
	if !errors.Is(err, ErrMissingUsage) {
		t.Fatalf("got %v want ErrMissingUsage", err)
	}
}
