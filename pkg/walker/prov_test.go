package walker

import (
	"testing"

	"github.com/omcrgnt/ecfg/internal/teststruct"
	"github.com/stretchr/testify/require"
)

func TestProvidersParity_shallowRootFields(t *testing.T) {
	pkgPath := "github.com/omcrgnt/ecfg/internal/teststruct"
	structName := "Config"

	var cfg teststruct.Config
	rp, err := NewReflectProvider(&cfg)
	require.NoError(t, err)
	rFields, err := rp.GetFields()
	require.NoError(t, err)

	tp, err := NewTypesProvider(pkgPath, structName)
	require.NoError(t, err)
	tFields, err := tp.GetFields()
	require.NoError(t, err)

	require.Len(t, tFields, len(rFields), "field count at Config root")

	for i := range rFields {
		require.Equal(t, rFields[i].Name(), tFields[i].Name(), "field name at index %d", i)
		require.Equal(t, rFields[i].Tag("ecfg"), tFields[i].Tag("ecfg"), "ecfg tag for %s", rFields[i].Name())
	}
}

func TestWalk_nested_reflect(t *testing.T) {
	type Inner struct {
		Key string `ecfg:"KEY"`
	}
	type Outer struct {
		Module Inner `ecfg:"MODULE"`
	}

	outer := Outer{Module: Inner{Key: "secret"}}
	p, err := NewReflectProvider(&outer)
	require.NoError(t, err)

	w := New()
	err = w.Walk(p, func(f Field) error {
		rv, sf, err := f.Value()
		if err != nil {
			return err
		}
		if sf.Name == "Key" {
			require.Equal(t, "secret", rv.String())
		}
		return nil
	})
	require.NoError(t, err)
}

func TestWalk_ValueErrorOnTypesProvider(t *testing.T) {
	p, err := NewTypesProvider("github.com/omcrgnt/ecfg/internal/teststruct", "Config")
	require.NoError(t, err)

	w := New()
	err = w.Walk(p, func(f Field) error {
		_, _, err := f.Value()
		return err
	})
	require.ErrorIs(t, err, ErrNoRuntimeValue)
}

func TestWalk_types_sliceOfStruct(t *testing.T) {
	p, err := NewTypesProvider("github.com/omcrgnt/ecfg/internal/teststruct", "Nested")
	require.NoError(t, err)

	var names []string
	w := New()
	err = w.Walk(p, func(f Field) error {
		names = append(names, f.Name())
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, []string{"Items", "Key"}, names)
}

func TestElemProvider_notContainer(t *testing.T) {
	var cfg teststruct.Config
	p, err := NewReflectProvider(&cfg)
	require.NoError(t, err)
	fields, err := p.GetFields()
	require.NoError(t, err)
	_, err = fields[0].ElemProvider()
	require.ErrorIs(t, err, ErrNotContainer)
}

func TestWalk_reflect_setValue(t *testing.T) {
	type Inner struct {
		N int
	}
	type Outer struct {
		Items []Inner
	}

	outer := Outer{Items: []Inner{{}}}
	p, err := NewReflectProvider(&outer)
	require.NoError(t, err)

	w := New()
	err = w.Walk(p, func(f Field) error {
		rv, sf, err := f.Value()
		if err != nil {
			return err
		}
		if sf.Name == "N" {
			rv.SetInt(42)
		}
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 42, outer.Items[0].N)
}

func TestReflectProvider_PointerStruct(t *testing.T) {
	type Extra struct {
		ID int `ecfg:"ID"`
	}
	type Config struct {
		DB *Extra `ecfg:"DB"`
	}

	cfg := Config{DB: &Extra{}}
	p, err := NewReflectProvider(&cfg)
	require.NoError(t, err)

	fields, err := p.GetFields()
	require.NoError(t, err)

	sub, err := fields[0].GetProvider()
	require.NoError(t, err)
	require.Equal(t, "Extra", sub.EntryName())
}
