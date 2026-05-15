package walker_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/omcrgnt/ecfg/pkg/walker"
	"github.com/stretchr/testify/require"
)

type SubConfig struct {
	Value string
}

type RootConfig struct {
	Str      string
	Ptr      *SubConfig
	Items    []string
	Settings map[string]int
}

func mustReflectProvider(t *testing.T, v any) walker.Provider {
	t.Helper()
	p, err := walker.NewReflectProvider(v)
	require.NoError(t, err)
	return p
}

func TestWalker_OptionsAndInit(t *testing.T) {
	t.Run("should NOT init nil pointer by default", func(t *testing.T) {
		w := walker.New()
		cfg := &RootConfig{}
		err := w.Walk(mustReflectProvider(t, cfg), func(walker.Field) error { return nil })
		require.NoError(t, err)
		require.Nil(t, cfg.Ptr)
	})

	t.Run("should init nil pointer with option", func(t *testing.T) {
		w := walker.New(walker.WithInitNilPointers())
		cfg := &RootConfig{}
		err := w.Walk(mustReflectProvider(t, cfg), func(walker.Field) error { return nil })
		require.NoError(t, err)
		require.NotNil(t, cfg.Ptr)
	})
}

func TestWalker_NodeHookAndHierarchy(t *testing.T) {
	var path []string
	cfg := RootConfig{
		Ptr:   &SubConfig{Value: "test"},
		Items: []string{"a", "b"},
	}

	w := walker.New(
		walker.WithInitNilPointers(),
		walker.WithNodeHook(func(info walker.NodeInfo, next func() error) error {
			path = append(path, info.Name)
			err := next()
			path = path[:len(path)-1]
			return err
		}),
	)

	collectedPaths := make(map[string]bool)
	err := w.Walk(mustReflectProvider(t, &cfg), func(f walker.Field) error {
		_, sf, err := f.Value()
		if err != nil {
			return err
		}
		full := strings.Join(path, ".")
		if sf.Name != "" {
			if full != "" {
				full += "."
			}
			full += sf.Name
		}
		collectedPaths[full] = true
		return nil
	})
	require.NoError(t, err)

	expected := []string{"Ptr.Value", "Items.0", "Items.1"}
	for _, p := range expected {
		require.True(t, collectedPaths[p], "expected path %s was not visited", p)
	}
}

func TestWalker_Negative(t *testing.T) {
	t.Run("should propagate error from NodeHook", func(t *testing.T) {
		customErr := errors.New("abort walk")
		w := walker.New(walker.WithNodeHook(func(info walker.NodeInfo, next func() error) error {
			return customErr
		}))

		cfg := RootConfig{Ptr: &SubConfig{}}
		err := w.Walk(mustReflectProvider(t, &cfg), func(walker.Field) error { return nil })
		require.ErrorIs(t, err, customErr)
	})

	t.Run("should propagate error from handler", func(t *testing.T) {
		customErr := errors.New("leaf error")
		w := walker.New()
		cfg := RootConfig{Str: "val"}
		err := w.Walk(mustReflectProvider(t, &cfg), func(walker.Field) error {
			return customErr
		})
		require.ErrorIs(t, err, customErr)
	})
}

func TestWalker_MapAndSlice(t *testing.T) {
	cfg := RootConfig{
		Settings: map[string]int{"one": 1, "two": 2},
		Str:      "test",
	}
	count := 0
	w := walker.New()

	err := w.Walk(mustReflectProvider(t, &cfg), func(walker.Field) error {
		count++
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, 3, count, "Str + 2 map elements")
}
