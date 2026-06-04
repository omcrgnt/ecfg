package testdata

import "fmt"

// AppConfig is the root configuration type for tests.
type AppConfig struct {
	Server ServerBlock `ecfg:"SERVER"`
}

// ServerBlock groups server settings.
type ServerBlock struct {
	Label Label `ecfg:"LABEL"`
}

// Label is a Go leaf wrapper.
type Label string

// Usage describes LABEL env.
func (Label) Usage() string {
	return "Метка приложения"
}

// Validate checks label value.
func (l Label) Validate() error {
	if l == "" {
		return fmt.Errorf("label empty")
	}
	return nil
}
