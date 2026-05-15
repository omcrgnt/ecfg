package walker

import (
	"fmt"
	"reflect"
	"strconv"
)

type NodeKind int

const (
	KindStruct NodeKind = iota
	KindSlice
	KindMap
)

type NodeInfo struct {
	Name string
	Kind NodeKind
	Tag  reflect.StructTag
}

type FieldContext struct {
	Value reflect.Value
	Field reflect.StructField
}

type WalkFunc func(ctx FieldContext) error

// NodeHook дает полный контроль над обходом контейнера.
// next() запускает обход внутренностей этого узла.
type NodeHook func(info NodeInfo, next func() error) error

type Walker struct {
	initNilPointers bool
	nodeHook        NodeHook
}

type Option func(*Walker)

func WithInitNilPointers() Option {
	return func(w *Walker) { w.initNilPointers = true }
}

func WithNodeHook(hook NodeHook) Option {
	return func(w *Walker) { w.nodeHook = hook }
}

func New(opts ...Option) *Walker {
	w := &Walker{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Walker) Walk(v reflect.Value, fn WalkFunc) error {
	if v.Kind() != reflect.Pointer && !v.CanAddr() {
		return fmt.Errorf("walker: expected pointer or addressable value, got %v", v.Kind())
	}
	return w.recursiveWalk(v, reflect.StructField{}, fn)
}

func (w *Walker) recursiveWalk(v reflect.Value, field reflect.StructField, fn WalkFunc) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !w.initNilPointers {
				return nil
			}
			if v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			}
		}
		return w.recursiveWalk(v.Elem(), field, fn)
	}

	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			fv := v.Field(i)
			if !f.IsExported() {
				continue
			}

			next := func() error {
				return w.recursiveWalk(fv, f, fn)
			}

			if isContainer(fv) && w.nodeHook != nil {
				info := NodeInfo{Name: f.Name, Kind: KindStruct, Tag: f.Tag}
				if err := w.nodeHook(info, next); err != nil {
					return err
				}
			} else if err := next(); err != nil {
				return err
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			info := NodeInfo{Name: strconv.Itoa(i), Kind: KindSlice}

			next := func() error {
				return w.recursiveWalk(v.Index(i), reflect.StructField{}, fn)
			}

			if w.nodeHook != nil {
				if err := w.nodeHook(info, next); err != nil {
					return err
				}
			} else if err := next(); err != nil {
				return err
			}
		}

	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			keyStr := fmt.Sprintf("%v", iter.Key().Interface())
			info := NodeInfo{Name: keyStr, Kind: KindMap}

			next := func() error {
				return w.recursiveWalk(iter.Value(), reflect.StructField{}, fn)
			}

			if w.nodeHook != nil {
				if err := w.nodeHook(info, next); err != nil {
					return err
				}
			} else if err := next(); err != nil {
				return err
			}
		}

	default:
		return fn(FieldContext{Value: v, Field: field})
	}
	return nil
}

func isContainer(v reflect.Value) bool {
	k := v.Kind()
	if k == reflect.Ptr {
		if v.IsNil() {
			return isContainerType(v.Type().Elem().Kind())
		}
		return isContainer(v.Elem())
	}
	return isContainerType(k)
}

func isContainerType(k reflect.Kind) bool {
	return k == reflect.Struct || k == reflect.Slice || k == reflect.Array || k == reflect.Map
}

func Process[T any](w *Walker, fn WalkFunc) (*T, error) {
	var target T
	val := reflect.ValueOf(&target).Elem()
	if err := w.Walk(val, fn); err != nil {
		return nil, err
	}
	return &target, nil
}
