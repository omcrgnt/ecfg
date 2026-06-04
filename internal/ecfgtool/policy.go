package ecfgtool

import (
	"reflect"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

func check(ctx visitCtx, isProto bool) error {
	f := ctx.walk.Field
	depth := ctx.walk.Depth

	switch kind := fieldKind(f); kind {
	case reflect.Slice, reflect.Map:
		return ErrUnsupportedContainer
	}

	if depth == 0 {
		if !isStructKind(f) {
			return ErrRootFieldKind
		}
		_, err := segment(0, parseEcfgTag(f.Tag), f.Name)
		return err
	}

	if depth == 1 && isStructKind(f) && !isProto {
		return ErrNestedBlockAtDepth1
	}

	if depth >= 2 && isStructKind(f) && !isProto {
		return ErrNestedBlockAtDepth1
	}

	return nil
}

func fieldKind(f walk.FieldDesc) reflect.Kind {
	if f.ReflectType != nil {
		return f.ReflectType.Kind()
	}
	if f.TypesType != nil {
		return walk.ReflectKind(f.TypesType)
	}
	return reflect.Invalid
}

func isStructKind(f walk.FieldDesc) bool {
	k := fieldKind(f)
	return k == reflect.Struct || k == reflect.Ptr
}

func parseEcfgTag(tag string) string {
	st := reflect.StructTag(tag)
	return st.Get("ecfg")
}
