package walk

import (
	"go/types"
	"reflect"
)

// FieldDesc describes one exported struct field at the current walk node.
type FieldDesc struct {
	Name string
	Tag  string
	// Exactly one source engine fills these:
	ReflectType reflect.Type
	TypesType   types.Type
}
