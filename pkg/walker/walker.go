package walker

import (
	"fmt"
	"reflect"
)

func Walk(p Provider, fn func(f Field) error) error {
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

		if f.IsStruct() && !f.IsProto() {
			subProv, err := f.GetProvider()
			if err != nil {
				return err
			}
			if err := Walk(subProv, fn); err != nil {
				return err
			}
		}
	}
	return nil
}
