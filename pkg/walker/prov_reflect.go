package walker

import (
	"fmt"
	"reflect"
)

type reflectProvider struct {
	v reflect.Value // struct value (not pointer)
}

func (*reflectProvider) runtimeProvider() {}

// NewReflectProvider builds a RuntimeProvider for v.
// v must be a pointer to struct, e.g. &cfg.
func NewReflectProvider(v any) (Provider, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("%w, got %s", ErrPointerRequired, rv.Kind())
	}
	ev := rv.Elem()
	if ev.Kind() != reflect.Struct {
		return nil, fmt.Errorf("walker: expected pointer to struct, got pointer to %s", ev.Kind())
	}
	return &reflectProvider{v: ev}, nil
}

func (p *reflectProvider) GetFields() ([]Field, error) {
	t := p.v.Type()
	var fields []Field
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		fields = append(fields, &reflectField{f: sf, parent: p.v})
	}
	return fields, nil
}

func (p *reflectProvider) EntryName() string { return p.v.Type().Name() }

type reflectField struct {
	f      reflect.StructField
	parent reflect.Value
}

func (f *reflectField) Name() string          { return f.f.Name }
func (f *reflectField) Tag(key string) string { return f.f.Tag.Get(key) }
func (f *reflectField) Kind() reflect.Kind    { return elemKind(f.f.Type) }
func (f *reflectField) IsStruct() bool        { return f.Kind() == reflect.Struct }
func (f *reflectField) IsProto() bool         { return isProtoReflectType(f.f.Type) }

func (f *reflectField) Value() (reflect.Value, reflect.StructField, error) {
	return f.parent.FieldByIndex(f.f.Index), f.f, nil
}

func (f *reflectField) GetProvider() (Provider, error) {
	fv := f.parent.FieldByIndex(f.f.Index)
	for fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return nil, fmt.Errorf("walker: field %s is nil pointer", f.f.Name)
		}
		fv = fv.Elem()
	}
	if fv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("walker: field %s is not a struct", f.f.Name)
	}
	return &reflectProvider{v: fv}, nil
}

func (f *reflectField) ElemProvider() (Provider, error) {
	elem, ok := elemStructReflectType(f.f.Type)
	if !ok {
		return nil, ErrNotContainer
	}
	return &reflectTypeProvider{t: elem}, nil
}

func elemStructReflectType(t reflect.Type) (reflect.Type, bool) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		t = t.Elem()
	case reflect.Map:
		t = t.Elem()
	default:
		return nil, false
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, false
	}
	if isProtoReflectType(t) {
		return nil, false
	}
	return t, true
}

func elemKind(t reflect.Type) reflect.Kind {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind()
}

func structReflectType(t reflect.Type) (reflect.Type, bool) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, false
	}
	return t, true
}

