package walker

import "errors"

var (
	// ErrNoRuntimeValue is returned by Field.Value on schema providers.
	ErrNoRuntimeValue = errors.New("walker: provider has no runtime values; use NewReflectProvider with a pointer to struct")

	// ErrPointerRequired is returned when v in NewReflectProvider is not *struct.
	ErrPointerRequired = errors.New("walker: reflect provider requires a pointer to struct")

	// ErrNotContainer is returned by ElemProvider when the field is not []struct, [N]struct, or map[K]struct.
	ErrNotContainer = errors.New("walker: field is not a container with struct element type")
)
