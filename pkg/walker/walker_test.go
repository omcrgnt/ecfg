package walker_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/omcrgnt/ecfg/pkg/walker"
)

// --- Тестовые структуры ---

type SubConfig struct {
	Value string
}

type RootConfig struct {
	Str      string
	Ptr      *SubConfig     // Для проверки WithInitNilPointers
	Items    []string       // Слайс
	Settings map[string]int // Мапа
}

// --- Тесты ---

func TestWalker_OptionsAndInit(t *testing.T) {
	t.Run("should NOT init nil pointer by default", func(t *testing.T) {
		w := walker.New()
		cfg := &RootConfig{}
		_ = w.Walk(reflect.ValueOf(cfg).Elem(), func(ctx walker.FieldContext) error { return nil })

		if cfg.Ptr != nil {
			t.Error("Expected nil pointer to remain nil")
		}
	})

	t.Run("should init nil pointer with option", func(t *testing.T) {
		w := walker.New(walker.WithInitNilPointers())
		cfg := &RootConfig{}
		_ = w.Walk(reflect.ValueOf(cfg).Elem(), func(ctx walker.FieldContext) error { return nil })

		if cfg.Ptr == nil {
			t.Error("Expected nil pointer to be initialized")
		}
	})
}

func TestWalker_NodeHookAndHierarchy(t *testing.T) {
	t.Run("should correctly track path through NodeHook", func(t *testing.T) {
		var path []string
		type SubConfig struct {
			Value string
		}
		type RootConfig struct {
			Ptr   *SubConfig
			Items []string
		}

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
		err := w.Walk(reflect.ValueOf(&cfg).Elem(), func(ctx walker.FieldContext) error {
			full := strings.Join(path, ".")
			if ctx.Field.Name != "" {
				if full != "" {
					full += "."
				}
				full += ctx.Field.Name
			}
			collectedPaths[full] = true
			return nil
		})

		if err != nil {
			t.Fatal(err)
		}

		// Теперь реально используем ожидаемые пути
		expected := []string{
			"Ptr.Value",
			"Items.0",
			"Items.1",
		}

		for _, p := range expected {
			if !collectedPaths[p] {
				t.Errorf("Expected path %s was not visited", p)
			}
		}
	})
}

func TestWalker_Negative(t *testing.T) {
	t.Run("should fail if input is not a pointer", func(t *testing.T) {
		w := walker.New()
		cfg := RootConfig{}
		err := w.Walk(reflect.ValueOf(cfg), func(ctx walker.FieldContext) error { return nil })
		if err == nil {
			t.Error("Expected error for non-pointer input")
		}
	})

	t.Run("should propagate error from NodeHook", func(t *testing.T) {
		customErr := errors.New("abort walk")
		w := walker.New(walker.WithNodeHook(func(info walker.NodeInfo, next func() error) error {
			return customErr
		}))

		cfg := RootConfig{Ptr: &SubConfig{}}
		err := w.Walk(reflect.ValueOf(&cfg).Elem(), func(ctx walker.FieldContext) error { return nil })

		if !errors.Is(err, customErr) {
			t.Errorf("Expected error '%v', got '%v'", customErr, err)
		}
	})

	t.Run("should propagate error from WalkFunc", func(t *testing.T) {
		customErr := errors.New("leaf error")
		w := walker.New()
		cfg := RootConfig{Str: "val"}
		err := w.Walk(reflect.ValueOf(&cfg).Elem(), func(ctx walker.FieldContext) error {
			return customErr
		})

		if !errors.Is(err, customErr) {
			t.Errorf("Expected error '%v', got '%v'", customErr, err)
		}
	})
}

func TestWalker_MapAndSlice(t *testing.T) {
	t.Run("should visit all map elements", func(t *testing.T) {
		cfg := RootConfig{
			Settings: map[string]int{"one": 1, "two": 2},
			Str:      "test",
		}
		count := 0
		w := walker.New()

		// Обходим и просто считаем каждый Leaf (Str + 2 элемента мапы)
		_ = w.Walk(reflect.ValueOf(&cfg).Elem(), func(ctx walker.FieldContext) error {
			count++
			return nil
		})

		// В RootConfig у нас: Str (1) + Items (0, т.к. пуст) + Settings (2 элемента) + Ptr (0, т.к. nil)
		// Итого ожидаем 3 вызова WalkFunc
		if count != 3 {
			t.Errorf("Expected 3 leaves, got %d", count)
		}
	})
}

func TestWalker_DebugPrint(t *testing.T) {
	type Extra struct {
		ID int `ecfg:"ID" usage:"ID сущности"`
	}

	type ComplexConfig struct {
		AppID   string        `ecfg:"APP_ID" usage:"Идентификатор приложения"`
		Timeout time.Duration `ecfg:"TIMEOUT" usage:"Таймаут операций"`
		Enabled bool          `ecfg:"ENABLED"`
		// Слайс структур
		Clusters []Extra `ecfg:"CLUSTERS"`
		// Мапа структур
		Services map[string]Extra `ecfg:"SERVICES"`
		// Указатель на структуру
		Database *Extra `ecfg:"DB"`
	}

	cfg := &ComplexConfig{
		AppID:   "my-app",
		Timeout: 5 * time.Second,
		Clusters: []Extra{
			{ID: 1},
			{ID: 2},
		},
		Services: map[string]Extra{
			"auth": {ID: 100},
		},
		// Database оставляем nil для проверки WithInitNilPointers
	}

	var path []string
	w := walker.New(
		walker.WithInitNilPointers(),
		walker.WithNodeHook(func(info walker.NodeInfo, next func() error) error {
			path = append(path, info.Name)
			fmt.Printf("➡️  ENTER: Name=%-10s Kind=%-10v Path=%v\n", info.Name, info.Kind, path)

			err := next()

			path = path[:len(path)-1]
			fmt.Printf("⬅️  EXIT : Name=%-10s Path=%v\n", info.Name, path)
			return err
		}),
	)

	fmt.Println("\n=== START WALKING ===")
	err := w.Walk(reflect.ValueOf(cfg).Elem(), func(ctx walker.FieldContext) error {
		fullPath := strings.Join(path, ".")
		if ctx.Field.Name != "" {
			if fullPath != "" {
				fullPath += "."
			}
			fullPath += ctx.Field.Name
		}

		fmt.Printf("  📍 LEAF: Field=%-10s Kind=%-10v Type=%-15s Path=%-20s Tag=%v\n",
			ctx.Field.Name,
			ctx.Value.Kind(),
			ctx.Value.Type().String(),
			fullPath,
			ctx.Field.Tag,
		)
		return nil
	})
	fmt.Println("=== END WALKING ===")

	if err != nil {
		t.Fatal(err)
	}
}
