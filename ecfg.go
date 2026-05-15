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

type options struct {
	prefix string
}

type Option func(*options)

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

	var target T
	p, err := walker.NewReflectProvider(&target)
	if err != nil {
		return nil, err
	}
	if err := w.Walk(p, func(f walker.Field) error {
		rv, sf, err := f.Value()
		if err != nil {
			return err
		}

		tag := sf.Tag.Get("ecfg")
		if len(pathStack) == 0 && tag == "" {
			return fmt.Errorf("ecfg: root field %s missing 'ecfg' tag", sf.Name)
		}

		name := tag
		if name == "" {
			name = sf.Name
		}

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

		return setFieldValue(rv, val)
	}); err != nil {
		return nil, err
	}
	return &target, nil
}

func setFieldValue(v reflect.Value, val string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(val)
			if err != nil {
				return fmt.Errorf("ecfg: parse duration %q: %w", val, err)
			}
			v.SetInt(int64(d))
			return nil
		}
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("ecfg: parse int %q: %w", val, err)
		}
		v.SetInt(parsed)

	case reflect.Bool:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("ecfg: parse bool %q: %w", val, err)
		}
		v.SetBool(parsed)

	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("ecfg: parse float %q: %w", val, err)
		}
		v.SetFloat(parsed)

	default:
		return fmt.Errorf("ecfg: unsupported type %s", v.Type())
	}
	return nil
}
