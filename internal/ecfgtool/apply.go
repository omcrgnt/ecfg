package ecfgtool

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

// ApplySeeded loads environment into specs from [builder.Seed] (shape from resource [BuildConfig], not wire fields).
func ApplySeeded(appResources any, seedMap map[string]any, opts Options) error {
	if len(seedMap) == 0 {
		return fmt.Errorf("ecfg: empty seed map")
	}

	rv, rt, err := appResourcesStruct(appResources)
	if err != nil {
		return wrapNotStruct(err)
	}

	applyOpts := opts
	applyOpts.SkipUntaggedRoot = true

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		if !sf.IsExported() {
			continue
		}
		ecfgTag := parseEcfgTag(string(sf.Tag))
		if applyOpts.SkipUntaggedRoot && ecfgTag == "" {
			continue
		}

		spec, ok := seedMap[sf.Name]
		if !ok {
			continue
		}

		eng, err := walk.NewEngineReflect(spec)
		if err != nil {
			return fmt.Errorf("ecfg: %s: %w", sf.Name, err)
		}

		if err := traverseRootBlocks(ecfgTag, sf.Name, eng, applyOpts, func(ctx visitCtx) error {
			if !ctx.isLeaf {
				return nil
			}

			env, err := lookupEnv(ctx.envKey)
			if err != nil {
				return err
			}

			spec, ok := seedMap[ctx.rootGroup]
			if !ok {
				rootField := strings.SplitN(ctx.fieldPath, ".", 2)[0]
				spec, ok = seedMap[rootField]
			}
			if !ok {
				return fmt.Errorf("ecfg: no spec for block %s", ctx.rootGroup)
			}

			parts := strings.Split(ctx.fieldPath, ".")
			if len(parts) < 2 {
				return fmt.Errorf("ecfg: invalid field path %q", ctx.fieldPath)
			}
			fieldName := parts[len(parts)-1]
			target, err := fieldValueByName(spec, fieldName)
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
		}); err != nil {
			return err
		}
	}

	_ = rv // wire values are not read; only field tags and names matter
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
