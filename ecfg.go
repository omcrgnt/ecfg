package ecfg

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/omcrgnt/ecfg/pkg/walker"
)

// Parse — основная точка входа. Создает и заполняет структуру типа T данными из ENV.
func Parse[T any]() (*T, error) {
	var pathStack []string

	w := walker.New(
		walker.WithInitNilPointers(),
		walker.WithNodeHook(func(info walker.NodeInfo, next func() error) error {
			tag := info.Tag.Get("ecfg")

			// ПРАВИЛО: На первом уровне (корень) тег 'ecfg' обязателен для структур
			if len(pathStack) == 0 && tag == "" {
				return fmt.Errorf("ecfg: field %s at root must have 'ecfg' tag", info.Name)
			}

			// Определяем имя сегмента (тег или имя поля)
			name := tag
			if name == "" {
				name = info.Name
			}

			// Управляем стэком пути
			oldPath := pathStack
			pathStack = append(pathStack, strings.ToUpper(name))

			err := next() // Рекурсивный уход вглубь

			pathStack = oldPath // Возврат стэка
			return err
		}),
	)

	return walker.Process[T](w, func(ctx walker.FieldContext) error {
		tag := ctx.Field.Tag.Get("ecfg")

		// ПРАВИЛО: На первом уровне тег обязателен и для простых полей
		if len(pathStack) == 0 && tag == "" {
			return fmt.Errorf("ecfg: root field %s missing 'ecfg' tag", ctx.Field.Name)
		}

		name := tag
		if name == "" {
			name = ctx.Field.Name
		}

		// Сборка финального ключа ENV
		fullKey := strings.ToUpper(name)
		if len(pathStack) > 0 {
			fullKey = strings.Join(pathStack, "_") + "_" + fullKey
		}

		// Читаем значение
		val, ok := os.LookupEnv(fullKey)
		if !ok || val == "" {
			return nil
		}

		// Записываем значение в поле
		return setFieldValue(ctx.Value, val)
	})
}

// setFieldValue конвертирует строку из ENV в тип поля структуры
func setFieldValue(v reflect.Value, val string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Обработка time.Duration
		if v.Type().String() == "time.Duration" {
			d, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("parse duration: %w", err)
			}
			v.SetInt(int64(d))
			return nil
		}
		// Обычные целые числа
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int: %w", err)
		}
		v.SetInt(parsed)

	case reflect.Bool:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("parse bool: %w", err)
		}
		v.SetBool(parsed)

	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("parse float: %w", err)
		}
		v.SetFloat(parsed)

	default:
		return fmt.Errorf("unsupported type: %s", v.Type().String())
	}
	return nil
}
