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

// Настройки парсинга
type options struct {
	prefix string
}

type Option func(*options)

// WithPrefix добавляет префикс ко всем переменным окружения (например, "APP_")
func WithPrefix(p string) Option {
	return func(o *options) {
		o.prefix = strings.ToUpper(strings.TrimSuffix(p, "_"))
	}
}

func Parse[T any](opts ...Option) (*T, error) {
	cfgOpts := &options{}
	for _, opt := range opts {
		opt(cfgOpts)
	}

	var pathStack []string

	w := walker.New(
		walker.WithInitNilPointers(),
		walker.WithNodeHook(func(info walker.NodeInfo, next func() error) error {
			tag := info.Tag.Get("ecfg")

			// Правило первого уровня (теперь учитываем, есть ли префикс)
			if len(pathStack) == 0 && tag == "" {
				return fmt.Errorf("ecfg: field %s at root must have 'ecfg' tag", info.Name)
			}

			name := tag
			if name == "" {
				name = info.Name
			}

			oldPath := pathStack
			pathStack = append(pathStack, strings.ToUpper(name))
			err := next()
			pathStack = oldPath
			return err
		}),
	)

	return walker.Process[T](w, func(ctx walker.FieldContext) error {
		tag := ctx.Field.Tag.Get("ecfg")
		if len(pathStack) == 0 && tag == "" {
			return fmt.Errorf("ecfg: root field %s missing 'ecfg' tag", ctx.Field.Name)
		}

		name := tag
		if name == "" {
			name = ctx.Field.Name
		}

		// Сборка ключа с учетом префикса
		parts := make([]string, 0, len(pathStack)+2)
		if cfgOpts.prefix != "" {
			parts = append(parts, cfgOpts.prefix)
		}
		parts = append(parts, pathStack...)
		parts = append(parts, strings.ToUpper(name))

		fullKey := strings.Join(parts, "_")

		val, ok := os.LookupEnv(fullKey)
		if !ok || val == "" {
			return nil
		}

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
