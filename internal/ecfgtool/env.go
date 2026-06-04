package ecfgtool

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

func lookupEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrMissingEnv, key)
	}
	if strings.TrimSpace(val) == "" {
		return "", fmt.Errorf("%w: %s", ErrEmptyEnv, key)
	}
	return val, nil
}

func setLeafValue(v reflect.Value, env string, isProto bool) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return fmt.Errorf("ecfg: nil pointer")
		}
		v = v.Elem()
	}
	if isProto {
		msg, ok := v.Addr().Interface().(proto.Message)
		if !ok {
			if p, ok := interface{}(v).(proto.Message); ok {
				msg = p
			} else {
				return ErrIncompatibleLeaf
			}
		}
		return setProtoValue(msg, env)
	}
	return setGoLeafValue(v, env)
}

func setGoLeafValue(v reflect.Value, env string) error {
	typ := v.Type()
	if !implementsUsage(typ) || !implementsValidator(typ) {
		return ErrIncompatibleLeaf
	}
	if v.Kind() == reflect.String {
		v.SetString(env)
		return nil
	}
	if v.Type().Kind() == reflect.String {
		v.Set(reflect.ValueOf(env).Convert(v.Type()))
		return nil
	}
	return setScalarValue(v, env)
}

func setScalarValue(v reflect.Value, env string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(env)
	case reflect.Bool:
		b, err := strconv.ParseBool(env)
		if err != nil {
			return fmt.Errorf("ecfg: parse bool: %w", err)
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isDurationType(v.Type()) {
			d, err := time.ParseDuration(env)
			if err != nil {
				return fmt.Errorf("ecfg: parse duration: %w", err)
			}
			v.SetInt(int64(d))
			return nil
		}
		n, err := strconv.ParseInt(env, 10, 64)
		if err != nil {
			return fmt.Errorf("ecfg: parse int: %w", err)
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(env, 10, 64)
		if err != nil {
			return fmt.Errorf("ecfg: parse uint: %w", err)
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(env, 64)
		if err != nil {
			return fmt.Errorf("ecfg: parse float: %w", err)
		}
		v.SetFloat(f)
	default:
		// named type: convert via underlying if possible
		if v.CanConvert(reflect.TypeOf("")) {
			v.Set(reflect.ValueOf(env).Convert(v.Type()))
			return nil
		}
		return fmt.Errorf("ecfg: unsupported scalar kind %s", v.Kind())
	}
	return nil
}

func isDurationType(t reflect.Type) bool {
	return t.ConvertibleTo(reflect.TypeOf(time.Duration(0))) && t.Kind() == reflect.Int64
}

func implementsUsage(typ reflect.Type) bool {
	return reflect.PointerTo(typ).Implements(reflect.TypeOf((*Usage)(nil)).Elem()) ||
		typ.Implements(reflect.TypeOf((*Usage)(nil)).Elem())
}

func implementsValidator(typ reflect.Type) bool {
	return reflect.PointerTo(typ).Implements(reflect.TypeOf((*Validator)(nil)).Elem()) ||
		typ.Implements(reflect.TypeOf((*Validator)(nil)).Elem())
}
