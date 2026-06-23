package ecfgtool

import (
	"fmt"
	"reflect"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

// Options configures Parse and ApplySeeded.
type Options struct {
	Prefix           string
	SkipUntaggedRoot bool
}

// Parse loads *T from environment variables.
func Parse[T any](opts Options) (*T, error) {
	var zero T
	eng, err := walk.NewEngineReflect(&zero)
	if err != nil {
		return nil, wrapNotStruct(err)
	}
	if err := traverse(eng, opts, func(ctx visitCtx) error {
		if !ctx.isLeaf {
			return nil
		}
		env, err := lookupEnv(ctx.envKey)
		if err != nil {
			return err
		}
		if err := setLeafValue(ctx.value, env, ctx.isProto); err != nil {
			return wrapFieldErr(ctx.fieldPath, ctx.envKey, err)
		}
		var validTarget any
		if ctx.isProto {
			v := ctx.value
			for v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			validTarget = v.Addr().Interface()
		} else {
			validTarget = ctx.value.Addr().Interface()
		}
		if err := validateLeaf(validTarget, ctx.isProto, ctx.usageText, ctx.fieldPath, ctx.envKey); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &zero, nil
}

func wrapNotStruct(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %s", ErrNotStruct, err)
}

func wrapFieldErr(fieldPath, envKey string, err error) error {
	return fmt.Errorf("ecfg: %s (%s): %w", fieldPath, envKey, err)
}
