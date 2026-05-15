package walker

import (
	"fmt"
	"reflect"
)

// reflectTypeProvider is a SchemaProvider for reflect.Type without values.
// Created by reflectField.ElemProvider for []struct schema descent.
type reflectTypeProvider struct {
	t reflect.Type
}

func (*reflectTypeProvider) schemaProvider() {}

func (p *reflectTypeProvider) GetFields() ([]Field, error) {
	if p.t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("walker: expected struct type, got %s", p.t.Kind())
	}
	var fields []Field
	for i := 0; i < p.t.NumField(); i++ {
		sf := p.t.Field(i)
		if !sf.IsExported() {
			continue
		}
		fields = append(fields, &reflectTypeField{t: p.t, f: sf})
	}
	return fields, nil
}

func (p *reflectTypeProvider) EntryName() string { return p.t.Name() }

type reflectTypeField struct {
	t reflect.Type
	f reflect.StructField
}

func (f *reflectTypeField) Name() string          { return f.f.Name }
func (f *reflectTypeField) Tag(key string) string { return f.f.Tag.Get(key) }
func (f *reflectTypeField) Kind() reflect.Kind    { return elemKind(f.f.Type) }
func (f *reflectTypeField) IsStruct() bool        { return f.Kind() == reflect.Struct }
func (f *reflectTypeField) IsProto() bool         { return isProtoReflectType(f.f.Type) }

func (f *reflectTypeField) Value() (reflect.Value, reflect.StructField, error) {
	return reflect.Value{}, reflect.StructField{}, ErrNoRuntimeValue
}

func (f *reflectTypeField) GetProvider() (Provider, error) {
	st, ok := structReflectType(f.f.Type)
	if !ok {
		return nil, fmt.Errorf("walker: field %s is not a struct", f.f.Name)
	}
	return &reflectTypeProvider{t: st}, nil
}

func (f *reflectTypeField) ElemProvider() (Provider, error) {
	elem, ok := elemStructReflectType(f.f.Type)
	if !ok {
		return nil, ErrNotContainer
	}
	return &reflectTypeProvider{t: elem}, nil
}
