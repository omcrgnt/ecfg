package ecfgtool

import (
	"errors"
	"reflect"
	"testing"

	"github.com/omcrgnt/ecfg/internal/testdata"
	"github.com/omcrgnt/ecfg/pkg/walk"
)

func TestCheck_depth2Struct(t *testing.T) {
	ctx := visitCtx{
		walk: walk.VisitCtx{
			Depth: 2,
			Field: walk.FieldDesc{ReflectType: reflect.TypeOf(struct{ X int }{})},
		},
	}
	err := check(ctx, false)
	if !errors.Is(err, ErrNestedBlockAtDepth1) {
		t.Fatalf("got %v want ErrNestedBlockAtDepth1", err)
	}
}

func TestParse_mapContainer(t *testing.T) {
	_, err := Parse[testdata.BadMap](Options{Prefix: "T"})
	if !errors.Is(err, ErrUnsupportedContainer) {
		t.Fatalf("got %v want ErrUnsupportedContainer", err)
	}
}
