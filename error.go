package ecfg

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrUnknownKind  = errors.New("unknown kind")
)

func errWrap(err error) error {
	return fmt.Errorf("failed to Parse: %w", err)
}
