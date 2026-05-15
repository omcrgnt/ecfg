package walker

import "reflect"

// Field is one exported struct field in the current Provider scope.
type Field interface {
	Name() string
	Tag(key string) string
	// IsStruct reports a nested struct (after peeling pointers).
	IsStruct() bool
	// IsProto reports a protobuf message; Walk treats it as a leaf.
	IsProto() bool
	Kind() reflect.Kind
	// GetProvider descends into a nested struct value or type.
	GetProvider() (Provider, error)
	// ElemProvider returns the element struct type for []T, [N]T, or map[K]T
	// when T is a struct. Used by schema walk only; runtime uses walkContent.
	ElemProvider() (Provider, error)
	// Value returns the field value and StructField. Schema providers return ErrNoRuntimeValue.
	Value() (reflect.Value, reflect.StructField, error)
}

// Provider lists exported fields of one struct (runtime or schema).
type Provider interface {
	GetFields() ([]Field, error)
	EntryName() string
}

// SchemaProvider is implemented by types-only providers (AST, reflect.Type).
// Walk uses walkFields: handler runs per field, Value() is not available.
type SchemaProvider interface {
	Provider
	schemaProvider()
}

// RuntimeProvider is implemented by NewReflectProvider.
// Walk uses walkValues and can enter slice/map elements by index or key.
type RuntimeProvider interface {
	Provider
	runtimeProvider()
}

// Handler is called for each leaf and, in schema mode, for each field before descent.
type Handler func(f Field) error
