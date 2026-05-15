package walker

import (
	"fmt"
	"reflect"
)

type reflectProvider struct {
	t reflect.Type
}

func NewReflectProvider(v interface{}) (Provider, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("walker: expected struct, got %s", t.Kind())
	}
	return &reflectProvider{t: t}, nil
}

func (p *reflectProvider) GetFields() ([]Field, error) {
	var fields []Field
	for i := 0; i < p.t.NumField(); i++ {
		f := p.t.Field(i)
		if !f.IsExported() {
			continue
		}
		fields = append(fields, &reflectField{f: f})
	}
	return fields, nil
}

func (p *reflectProvider) EntryName() string { return p.t.Name() }

type reflectField struct{ f reflect.StructField }

func (f *reflectField) Name() string          { return f.f.Name }
func (f *reflectField) Tag(key string) string { return f.f.Tag.Get(key) }
func (f *reflectField) Kind() reflect.Kind {
	return elemKind(f.f.Type)
}
func (f *reflectField) IsStruct() bool { return f.Kind() == reflect.Struct }
func (f *reflectField) IsProto() bool {
	t := f.f.Type
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	_, ok := reflect.PointerTo(t).MethodByName("ProtoMessage")
	return ok
}
func (f *reflectField) GetProvider() (Provider, error) {
	st, ok := structReflectType(f.f.Type)
	if !ok {
		return nil, fmt.Errorf("walker: field %s is not a struct", f.f.Name)
	}
	return &reflectProvider{t: st}, nil
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
