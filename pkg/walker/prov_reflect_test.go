package walker

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/ecfg/internal/teststruct"
	"github.com/stretchr/testify/require"
)

func TestNewReflectProvider_requiresPointer(t *testing.T) {
	_, err := NewReflectProvider(teststruct.Config{})
	require.ErrorIs(t, err, ErrPointerRequired)
}

func TestReflectProvider(t *testing.T) {
	var cfg teststruct.Config
	p, err := NewReflectProvider(&cfg)
	require.NoError(t, err)

	fields, err := p.GetFields()
	require.NoError(t, err)
	require.Len(t, fields, 2)

	f := fields[0]
	require.Equal(t, "Port", f.Name())
	require.Equal(t, "PORT", f.Tag("ecfg"))
	require.Equal(t, "Server port", f.Tag("usage"))
	require.Equal(t, reflect.Int, f.Kind())

	_, _, err = f.Value()
	require.NoError(t, err)
}
