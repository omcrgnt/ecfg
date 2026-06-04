package walk

import (
	"fmt"
	"reflect"
)

// EngineReflect walks a live struct value via reflect.
type EngineReflect struct {
	v reflect.Value
}

// NewEngineReflect returns an Engine for the struct pointed to or held by target.
func NewEngineReflect(target any) (*EngineReflect, error) {
	v := reflect.ValueOf(target)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("walk: nil pointer")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("walk: %T is not a struct", target)
	}
	return &EngineReflect{v: v}, nil
}

func (e *EngineReflect) Fields() ([]FieldDesc, error) {
	t := e.v.Type()
	var out []FieldDesc
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		out = append(out, FieldDesc{
			Name:        sf.Name,
			Tag:         string(sf.Tag),
			ReflectType: sf.Type,
		})
	}
	return out, nil
}

func (e *EngineReflect) Child(desc FieldDesc) (Engine, error) {
	fv, _, err := e.FieldValue(desc)
	if err != nil {
		return nil, err
	}
	return childEngineReflect(fv, true)
}

// FieldValue returns the settable value and StructField for desc.
func (e *EngineReflect) FieldValue(desc FieldDesc) (reflect.Value, reflect.StructField, error) {
	sf, ok := e.v.Type().FieldByName(desc.Name)
	if !ok {
		return reflect.Value{}, reflect.StructField{}, fmt.Errorf("walk: field %s not found", desc.Name)
	}
	return e.v.FieldByIndex(sf.Index), sf, nil
}

func childEngineReflect(v reflect.Value, initPtr bool) (*EngineReflect, error) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !initPtr || !v.CanSet() {
				return nil, fmt.Errorf("walk: nil pointer")
			}
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, nil
	}
	return &EngineReflect{v: v}, nil
}

// InitPointerField allocates nil *struct field before visit.
func (e *EngineReflect) InitPointerField(desc FieldDesc) error {
	fv, _, err := e.FieldValue(desc)
	if err != nil {
		return err
	}
	if fv.Kind() == reflect.Ptr && fv.CanSet() && fv.IsNil() {
		elem := fv.Type().Elem()
		if elem.Kind() == reflect.Struct {
			fv.Set(reflect.New(elem))
		}
	}
	return nil
}
