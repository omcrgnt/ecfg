package walker

import (
	"fmt"
	"reflect"
	"strconv"
)

// NodeKind classifies a container segment in NodeHook paths.
type NodeKind int

const (
	KindStruct NodeKind = iota
	KindSlice
	KindMap
)

// NodeInfo is passed to NodeHook for struct fields and slice/map segments.
type NodeInfo struct {
	Name string
	Kind NodeKind
	Tag  reflect.StructTag
}

// NodeHook wraps descent into a container. Call next() to visit children.
// Used to build logical paths (e.g. env key prefixes) without changing walk order.
type NodeHook func(info NodeInfo, next func() error) error

// Walker configures struct traversal.
type Walker struct {
	initNilPointers bool
	nodeHook        NodeHook
}

// Option configures Walker.
type Option func(*Walker)

// WithInitNilPointers allocates nil pointer fields before descending (runtime only).
func WithInitNilPointers() Option {
	return func(w *Walker) { w.initNilPointers = true }
}

// WithNodeHook registers a hook for struct fields and slice/map segments.
func WithNodeHook(hook NodeHook) Option {
	return func(w *Walker) { w.nodeHook = hook }
}

// New builds a Walker with the given options.
func New(opts ...Option) *Walker {
	w := &Walker{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Walk traverses the tree rooted at p. p must implement RuntimeProvider or SchemaProvider.
func (w *Walker) Walk(p Provider, fn Handler) error {
	switch {
	case isRuntimeProvider(p):
		return w.walkValues(p, fn)
	case isSchemaProvider(p):
		return walkFields(p, fn)
	default:
		return fmt.Errorf("walker: provider %T must implement RuntimeProvider or SchemaProvider", p)
	}
}

func isRuntimeProvider(p Provider) bool {
	_, ok := p.(RuntimeProvider)
	return ok
}

func isSchemaProvider(p Provider) bool {
	_, ok := p.(SchemaProvider)
	return ok
}

func (w *Walker) walkValues(p Provider, fn Handler) error {
	fields, err := p.GetFields()
	if err != nil {
		return fmt.Errorf("walker: %s: %w", p.EntryName(), err)
	}

	for _, field := range fields {
		val, sf, err := field.Value()
		if err != nil {
			return err
		}
		if err := w.walkValue(val, sf, field, fn); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkValue(v reflect.Value, sf reflect.StructField, field Field, fn Handler) error {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !w.initNilPointers {
				return nil
			}
			if v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			}
		}
		return w.walkValue(v.Elem(), sf, field, fn)
	}

	if sf.Name != "" && w.nodeHook != nil && isContainer(v) {
		info := NodeInfo{Name: sf.Name, Kind: nodeKind(v), Tag: sf.Tag}
		return w.nodeHook(info, func() error {
			return w.walkContent(v, field, fn)
		})
	}

	return w.walkContent(v, field, fn)
}

func (w *Walker) walkContent(v reflect.Value, field Field, fn Handler) error {
	switch v.Kind() {
	case reflect.Struct:
		if field != nil && field.IsProto() {
			return fn(field)
		}

		sub, err := w.structProvider(v, field)
		if err != nil {
			return err
		}
		return w.walkValues(sub, fn)

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			idx := i
			info := NodeInfo{Name: strconv.Itoa(idx), Kind: KindSlice}
			next := func() error {
				return w.walkContent(v.Index(idx), nil, fn)
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
			val := iter.Value()
			next := func() error {
				return w.walkContent(val, nil, fn)
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
		if field == nil {
			field = &leafField{v: v}
		}
		return fn(field)
	}
	return nil
}

// leafField stands in for a slice/map element that is not a struct field (no tags).
type leafField struct {
	v reflect.Value
}

func (f *leafField) Name() string          { return "" }
func (f *leafField) Tag(key string) string { return "" }
func (f *leafField) IsStruct() bool        { return false }
func (f *leafField) IsProto() bool         { return false }
func (f *leafField) Kind() reflect.Kind    { return f.v.Kind() }
func (f *leafField) GetProvider() (Provider, error) {
	return nil, ErrNotContainer
}
func (f *leafField) ElemProvider() (Provider, error) {
	return nil, ErrNotContainer
}
func (f *leafField) Value() (reflect.Value, reflect.StructField, error) {
	return f.v, reflect.StructField{}, nil
}

func nodeKind(v reflect.Value) NodeKind {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return KindSlice
	case reflect.Map:
		return KindMap
	default:
		return KindStruct
	}
}

func (w *Walker) structProvider(v reflect.Value, field Field) (Provider, error) {
	if field != nil {
		return field.GetProvider()
	}
	return &reflectProvider{v: v}, nil
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
