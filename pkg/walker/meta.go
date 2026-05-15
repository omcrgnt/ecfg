package walker

import (
	"errors"
	"fmt"
	"reflect"
)

// walkFields visits every field in schema mode: handler on each field, then descent.
func walkFields(p Provider, fn Handler) error {
	fields, err := p.GetFields()
	if err != nil {
		return fmt.Errorf("walker: %s: %w", p.EntryName(), err)
	}

	for _, f := range fields {
		if f.Kind() == reflect.Invalid {
			return fmt.Errorf("walker: %s: field %s has invalid type", p.EntryName(), f.Name())
		}

		if err := fn(f); err != nil {
			return err
		}

		if err := descendFields(f, fn); err != nil {
			return err
		}
	}
	return nil
}

// descendFields follows nested structs or container element types (no slice indices).
func descendFields(f Field, fn Handler) error {
	if f.IsProto() {
		return nil
	}

	if f.IsStruct() {
		sub, err := f.GetProvider()
		if err != nil {
			return err
		}
		return walkFields(sub, fn)
	}

	sub, err := f.ElemProvider()
	if err != nil {
		if errors.Is(err, ErrNotContainer) {
			return nil
		}
		return err
	}
	return walkFields(sub, fn)
}
