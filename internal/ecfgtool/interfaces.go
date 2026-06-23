package ecfgtool

// Usage provides human-readable hint text for env.md and validation errors.
type Usage interface {
	Usage() string
}

// Validator validates a leaf value after it was assigned from ENV.
type Validator interface {
	Validate() error
}
