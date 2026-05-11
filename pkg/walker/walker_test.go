package walker_test

import (
	"reflect"
	"testing"

	"github.com/omcrgnt/ecfg/pkg/walker"
)

func TestWalker_Complex(t *testing.T) {
	// 1. Описываем тестовую структуру (включая дженерик и указатели)
	type SubConfig struct {
		Port int `ecfg:"PORT"`
	}

	type GenericConfig[T any] struct {
		Payload T
	}

	type Config struct {
		AppID   string                `ecfg:"APP_ID"`
		Server  *SubConfig            // nil-указатель
		General GenericConfig[string] // дженерик
	}

	// 2. Стек для сборки пути (как это будет в ecfg)
	var pathStack []string
	collected := make(map[string]interface{})

	w := walker.New(
		walker.WithInitNilPointers(),
		walker.WithOnEnter(
			func(info walker.NodeInfo) {
				pathStack = append(pathStack, info.Name)
			}),
		walker.WithOnExit(
			func(info walker.NodeInfo) {
				pathStack = pathStack[:len(pathStack)-1]
			}),
	)

	// 3. Запускаем процесс
	cfg, err := walker.Process[Config](w, func(ctx walker.FieldContext) error {
		// Собираем полный путь для проверки
		fullPath := ""
		for _, p := range pathStack {
			fullPath += p + "_"
		}
		fullPath += ctx.Field.Name

		collected[fullPath] = ctx.Value.Interface()
		return nil
	})

	// 4. Проверки
	if err != nil {
		t.Fatalf("Walker failed: %v", err)
	}

	// Проверяем, что nil-указатель Server был инициализирован
	if cfg.Server == nil {
		t.Error("Expected cfg.Server to be initialized, but got nil")
	}

	// Проверяем собранные пути
	expectedPaths := []string{
		"AppID",
		"Server_Port",
		"General_Payload",
	}

	for _, p := range expectedPaths {
		if _, ok := collected[p]; !ok {
			t.Errorf("Path %s was not collected", p)
		}
	}
}

func TestWalker_Slice(t *testing.T) {
	type Item struct {
		Name string
	}
	type Root struct {
		Items []Item
	}

	// Инициализируем слайс с данными
	cfg := Root{
		Items: []Item{{Name: "first"}, {Name: "second"}},
	}

	paths := []string{}
	w := walker.New(walker.WithOnEnter(
		func(info walker.NodeInfo) { paths = append(paths, info.Name) },
	),
	)

	err := w.Walk(reflect.ValueOf(&cfg).Elem(), func(ctx walker.FieldContext) error {
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	// Проверяем, что индексы попали в путь: Items -> 0 -> Name, Items -> 1 -> Name
	// Но так как у нас логика OnEnter только для контейнеров:
	foundIndex := false
	for _, p := range paths {
		if p == "0" || p == "1" {
			foundIndex = true
		}
	}
	if !foundIndex {
		t.Error("Indices not found in paths for slice elements")
	}
}
