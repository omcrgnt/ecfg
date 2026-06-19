package ecfg

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/builder"
	"github.com/omcrgnt/res"
)

// Register adds first-level AppConfig fields that implement [builder.Builder]
// to reg without tags (explicit app configs). Skips unexported and nil fields.
func Register(cfg any, reg res.Registry) error {
	if reg == nil {
		return fmt.Errorf("ecfg: nil registry")
	}

	rv, err := structValue(cfg)
	if err != nil {
		return err
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fieldVal := rv.Field(i)
		if !fieldVal.CanInterface() {
			continue
		}
		if isNilValue(fieldVal) {
			continue
		}

		field := fieldVal.Interface()
		if _, ok := field.(builder.Builder); !ok {
			continue
		}

		if err := reg.Add(field); err != nil {
			return fmt.Errorf("ecfg: %s: %w", rt.Field(i).Name, err)
		}
	}

	return nil
}

func structValue(v any) (reflect.Value, error) {
	if v == nil {
		return reflect.Value{}, fmt.Errorf("ecfg: nil config")
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, fmt.Errorf("ecfg: nil config")
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("ecfg: want struct, got %s", rv.Kind())
	}

	return rv, nil
}

func isNilValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}
