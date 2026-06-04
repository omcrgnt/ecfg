package ecfgtool

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func isProtoMessage(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}
	return reflect.PointerTo(typ).Implements(reflect.TypeOf((*proto.Message)(nil)).Elem())
}

func checkProtoWrapper(typ reflect.Type) error {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("%w: not a struct", ErrInvalidProtoWrapper)
	}
	n := 0
	var hasValue bool
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if !sf.IsExported() {
			continue
		}
		n++
		if sf.Name == "Value" {
			hasValue = true
		}
	}
	if n != 1 || !hasValue {
		return fmt.Errorf("%w: want exactly one exported field Value, got %d exported fields", ErrInvalidProtoWrapper, n)
	}
	return nil
}

func setProtoValue(msg proto.Message, env string) error {
	if err := checkProtoWrapper(reflect.TypeOf(msg)); err != nil {
		return err
	}
	rv := reflect.ValueOf(msg)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	valueField := rv.FieldByName("Value")
	if !valueField.IsValid() || !valueField.CanSet() {
		return fmt.Errorf("%w: Value field not settable", ErrInvalidProtoWrapper)
	}
	return setScalarValue(valueField, env)
}

func protoReflectType(typ reflect.Type) protoreflect.MessageType {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	m, ok := reflect.New(typ).Interface().(proto.Message)
	if !ok {
		return nil
	}
	return m.ProtoReflect().Type()
}
