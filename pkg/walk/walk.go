package walk

import (
	"errors"
	"reflect"
)

// Engine lists exported fields of one struct node.
type Engine interface {
	Fields() ([]FieldDesc, error)
	Child(FieldDesc) (Engine, error)
}

// StructWalk walks a struct tree rooted at eng.
func StructWalk(eng Engine, opts Options, visit func(VisitCtx) error) error {
	return walkStruct(eng, opts, visit, 0)
}

// VisitCtx is passed to visit for each exported field.
type VisitCtx struct {
	Depth  int
	Field  FieldDesc
	Engine Engine
}

// Options configures StructWalk.
type Options struct {
	InitPointers bool
	AfterField   func(VisitCtx)
}

// SkipDescend is returned from visit to avoid descending into a struct field.
func SkipDescend() error { return skipDescend }

func walkStruct(eng Engine, opts Options, visit func(VisitCtx) error, depth int) error {
	fields, err := eng.Fields()
	if err != nil {
		return err
	}
	for _, f := range fields {
		ctx := VisitCtx{Depth: depth, Field: f, Engine: eng}
		visitErr := visit(ctx)
		if visitErr != nil && !errors.Is(visitErr, skipDescend) {
			return visitErr
		}
		if isStructField(f) && !errors.Is(visitErr, skipDescend) {
			child, err := eng.Child(f)
			if err != nil {
				return err
			}
			if child != nil {
				if err := walkStruct(child, opts, visit, depth+1); err != nil {
					return err
				}
			}
		}
		if opts.AfterField != nil {
			opts.AfterField(ctx)
		}
	}
	return nil
}

func isStructField(f FieldDesc) bool {
	if f.ReflectType != nil {
		k := f.ReflectType.Kind()
		return k == reflect.Struct || k == reflect.Ptr
	}
	if f.TypesType != nil {
		k := reflectKind(f.TypesType)
		return k == reflect.Struct || k == reflect.Ptr
	}
	return false
}
