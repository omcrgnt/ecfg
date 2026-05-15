package walker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type bareProvider struct{}

func (bareProvider) GetFields() ([]Field, error) { return nil, nil }
func (bareProvider) EntryName() string           { return "bare" }

func TestWalker_unknownProvider(t *testing.T) {
	w := New()
	err := w.Walk(bareProvider{}, func(Field) error { return nil })
	require.Error(t, err)
	require.Contains(t, err.Error(), "RuntimeProvider")
}

func TestProvider_markerInterfaces(t *testing.T) {
	var _ RuntimeProvider = (*reflectProvider)(nil)
	var _ SchemaProvider = (*typesProvider)(nil)
	var _ SchemaProvider = (*reflectTypeProvider)(nil)
}
