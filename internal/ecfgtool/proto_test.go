package ecfgtool

import (
	"errors"
	"reflect"
	"testing"

	commonv1 "github.com/omcrgnt/proto/gen/go/common/v1"
)

func TestCheckProtoWrapper_valid(t *testing.T) {
	if err := checkProtoWrapper(reflect.TypeOf(commonv1.Port{})); err != nil {
		t.Fatal(err)
	}
}

func TestCheckProtoWrapper_invalid(t *testing.T) {
	type bad struct {
		Value int
		Extra int
	}
	if err := checkProtoWrapper(reflect.TypeOf(bad{})); !errors.Is(err, ErrInvalidProtoWrapper) {
		t.Fatalf("got %v", err)
	}
}

func TestProtoReflectType_doublePointer(t *testing.T) {
	var msg **commonv1.Port
	typ := reflect.TypeOf(msg)
	mt := protoReflectType(typ)
	if mt == nil {
		t.Fatal("expected message type for **Port")
	}
	if mt.Descriptor().FullName() != "common.v1.Port" {
		t.Fatalf("got %v", mt.Descriptor().FullName())
	}
}
