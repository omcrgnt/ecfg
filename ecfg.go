package ecfg

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Parse(source any) error {
	return parse(source)
}

func parse(source any) error {
	if err := checkSourceType(source); err != nil {
		return errWrap(err)
	}

	if err := walkThroughStruct(source); err != nil {
		return errWrap(err)
	}

	return nil
}

const tagName = "ecfg"

// source must be pointer to struct
func walkThroughStruct(source any) error {
	var typ = reflect.TypeOf(source)
	var val = reflect.ValueOf(source)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := range typ.NumField() {
		var structField = typ.Field(i)
		if !structField.IsExported() {
			continue
		}
		var envPrefix = strings.ToUpper(structField.Tag.Get(tagName))
		if envPrefix == "" {
			return errors.New("empsty ecfg tag")
		}
		visitField(structField, val.Field(i), envPrefix)
	}
	return nil
}

func visitField(structField reflect.StructField, value reflect.Value, envPrefix string) {
	fmt.Printf("%#v %#v %#v\n", structField, value, envPrefix)
	newIfNilPointer(structField, value)

	switch structField.Type.Kind() {
	case reflect.Invalid:
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Uintptr:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.Slice:
	case reflect.String:
		fmt.Println("::::::::::::::::", structField)
	case reflect.Struct:
		walkThroughStruct(value.Addr())
	case reflect.UnsafePointer:
	default:
	}

}

func newIfNilPointer(structField reflect.StructField, value reflect.Value) {
	switch structField.Type.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			value.Set(reflect.New(structField.Type.Elem()).Elem().Addr())
		}
	default:
		fmt.Println("not ponter")
	}
}
