package ecfg

import (
	"fmt"
	"strings"

	"github.com/omcrgnt/ecfg/pkg/walker"
)

func Parse[T any]() (*T, error) {
	var path []string

	w := walker.New(
		walker.WithInitNilPointers(),
		walker.WithOnEnter(func(info walker.NodeInfo) {
			// Проверка для структур первого уровня
			name := info.Tag.Get("ecfg")
			if len(path) == 0 && name == "" {
				// В реальном коде тут можно либо паниковать,
				// либо прокидывать ошибку через состояние, если нужно
				panic(fmt.Sprintf("field %s must have ecfg tag", info.Name))
			}

			if name == "" {
				name = info.Name
			}
			path = append(path, name)
		}),
		walker.WithOnExit(func(info walker.NodeInfo) {
			path = path[:len(path)-1]
		}),
	)

	return walker.Process[T](w, func(ctx walker.FieldContext) error {
		// Проверка для "листьев" первого уровня
		if len(path) == 0 && ctx.Field.Tag.Get("ecfg") == "" {
			return fmt.Errorf("root field %s missing ecfg tag", ctx.Field.Name)
		}

		fmt.Println(strings.Join(path, "_"))
		return nil
	})
}
