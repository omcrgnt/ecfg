package walker

import (
	"go/types"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicKind_preservesWidth(t *testing.T) {
	cases := []struct {
		basic types.BasicKind
		want  reflect.Kind
	}{
		{types.Int, reflect.Int},
		{types.Int32, reflect.Int32},
		{types.Uint8, reflect.Uint8},
		{types.Float32, reflect.Float32},
		{types.String, reflect.String},
		{types.Bool, reflect.Bool},
	}
	for _, tc := range cases {
		require.Equal(t, tc.want, basicKind(tc.basic), "basic %v", tc.basic)
	}
}
