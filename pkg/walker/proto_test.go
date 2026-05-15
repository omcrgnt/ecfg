package walker

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"golang.org/x/tools/go/packages"
)

func TestIsProtoReflectType(t *testing.T) {
	require.True(t, isProtoReflectType(reflect.TypeOf(timestamppb.Timestamp{})))
	require.True(t, isProtoReflectType(reflect.TypeOf(&timestamppb.Timestamp{})))
	require.False(t, isProtoReflectType(reflect.TypeOf("")))
	require.False(t, isProtoReflectType(reflect.TypeOf(struct{}{})))
}

func TestIsProtoTypesType(t *testing.T) {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedTypes}, "google.golang.org/protobuf/types/known/timestamppb")
	require.NoError(t, err)
	require.NotEmpty(t, pkgs)

	obj := pkgs[0].Types.Scope().Lookup("Timestamp")
	require.NotNil(t, obj)
	require.True(t, isProtoTypesType(obj.Type()))
}
