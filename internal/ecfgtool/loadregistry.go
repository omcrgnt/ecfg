package ecfgtool

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

// RegistryEntry is one config value stored in a registry (custom tags read by key).
type RegistryEntry interface {
	Value() any
	GetCustomTag(key string) (any, bool)
}

// LoadRegistry initializes config field values from the environment for each entry
// that has custom tag tagKey (value = ecfg segment, e.g. SERVICE_ITEM).
func LoadRegistry(walkEntries func(func(RegistryEntry) bool), tagKey string, opts Options) error {
	if tagKey == "" {
		return fmt.Errorf("ecfg: empty custom tag key")
	}

	var firstErr error
	walkEntries(func(e RegistryEntry) bool {
		if firstErr != nil {
			return false
		}
		if err := loadRegistryEntry(e, tagKey, opts); err != nil {
			firstErr = err
			return false
		}
		return true
	})
	return firstErr
}

func loadRegistryEntry(e RegistryEntry, tagKey string, opts Options) error {
	segVal, ok := e.GetCustomTag(tagKey)
	if !ok {
		return fmt.Errorf("ecfg: missing custom tag %q", tagKey)
	}
	seg, ok := segVal.(string)
	if !ok || seg == "" {
		return fmt.Errorf("ecfg: invalid custom tag value for key %q", tagKey)
	}

	cfg := e.Value()
	if err := validateConfigValue(cfg); err != nil {
		return err
	}

	rootField := reflect.TypeOf(cfg).String()
	eng, err := walk.NewEngineReflect(cfg)
	if err != nil {
		return fmt.Errorf("ecfg: %s: %w", rootField, err)
	}

	applyOpts := opts
	applyOpts.SkipUntaggedRoot = true

	return traverseRootBlocks(seg, rootField, eng, applyOpts, func(ctx visitCtx) error {
		if !ctx.isLeaf {
			return nil
		}

		env, err := lookupEnv(ctx.envKey)
		if err != nil {
			return err
		}

		parts := strings.Split(ctx.fieldPath, ".")
		if len(parts) < 2 {
			return fmt.Errorf("ecfg: invalid field path %q", ctx.fieldPath)
		}
		fieldName := parts[len(parts)-1]
		target, err := fieldValueByName(cfg, fieldName)
		if err != nil {
			return wrapFieldErr(ctx.fieldPath, ctx.envKey, err)
		}

		if err := setLeafValue(target, env, ctx.isProto); err != nil {
			return wrapFieldErr(ctx.fieldPath, ctx.envKey, err)
		}

		var validTarget any
		if ctx.isProto {
			v := target
			for v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			validTarget = v.Addr().Interface()
		} else {
			validTarget = target.Addr().Interface()
		}
		if err := validateLeaf(validTarget, ctx.isProto, ctx.usageText, ctx.fieldPath, ctx.envKey); err != nil {
			return err
		}
		return nil
	})
}

func validateConfigValue(v any) error {
	if v == nil {
		return fmt.Errorf("ecfg: nil config")
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		if rv.IsNil() {
			return fmt.Errorf("ecfg: nil config")
		}
	}
	return nil
}

func fieldValueByName(spec any, fieldName string) (reflect.Value, error) {
	rv := reflect.ValueOf(spec)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil spec")
		}
		rv = rv.Elem()
	}

	rv = rv.FieldByName(fieldName)
	if !rv.IsValid() {
		return reflect.Value{}, fmt.Errorf("no field %q on %T", fieldName, spec)
	}
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil pointer at %q", fieldName)
		}
		rv = rv.Elem()
	}

	if !rv.CanSet() {
		return reflect.Value{}, fmt.Errorf("field %q is not settable", fieldName)
	}
	return rv, nil
}
