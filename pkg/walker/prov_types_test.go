package walker

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypesProvider(t *testing.T) {
	pkgPath := "github.com/omcrgnt/ecfg/internal/teststruct"

	p, err := NewTypesProvider(pkgPath, "Config")
	require.NoError(t, err)

	fields, err := p.GetFields()
	require.NoError(t, err)

	expected := map[string]struct {
		tag   string
		usage string
		kind  reflect.Kind
	}{
		"Port": {"PORT", "Server port", reflect.Int},
		"User": {"USER", "User name", reflect.String},
	}

	require.Len(t, fields, len(expected))

	for _, f := range fields {
		exp, ok := expected[f.Name()]
		require.True(t, ok, "unexpected field %s", f.Name())
		require.Equal(t, exp.tag, f.Tag("ecfg"), "ecfg tag for %s", f.Name())
		require.Equal(t, exp.usage, f.Tag("usage"), "usage tag for %s", f.Name())
		require.Equal(t, exp.kind, f.Kind(), "kind for %s", f.Name())
	}
}
