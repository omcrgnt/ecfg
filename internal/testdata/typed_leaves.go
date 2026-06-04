package testdata

import (
	"fmt"
	"time"
)

// IntCfg uses an int-based leaf.
type IntCfg struct {
	Block IntBlock `ecfg:"BLOCK"`
}

// IntBlock holds Port leaf.
type IntBlock struct {
	Port Port `ecfg:"PORT"`
}

// Port is a named int leaf.
type Port int

// Usage describes PORT.
func (Port) Usage() string { return "TCP port" }

// Validate checks port range.
func (p Port) Validate() error {
	if p <= 0 || p > 65535 {
		return fmt.Errorf("port out of range: %d", p)
	}
	return nil
}

// BoolCfg uses a bool leaf.
type BoolCfg struct {
	Block BoolBlock `ecfg:"BLOCK"`
}

// BoolBlock holds Enabled leaf.
type BoolBlock struct {
	Enabled Enabled `ecfg:"ENABLED"`
}

// Enabled is a bool leaf.
type Enabled bool

// Usage describes ENABLED.
func (Enabled) Usage() string { return "Feature enabled" }

// Validate is a no-op for tests.
func (Enabled) Validate() error { return nil }

// FloatCfg uses a float leaf.
type FloatCfg struct {
	Block FloatBlock `ecfg:"BLOCK"`
}

// FloatBlock holds Score.
type FloatBlock struct {
	Score Score `ecfg:"SCORE"`
}

// Score is a float64 leaf.
type Score float64

// Usage describes SCORE.
func (Score) Usage() string { return "Score value" }

// Validate is a no-op.
func (Score) Validate() error { return nil }

// DurationCfg uses a duration leaf.
type DurationCfg struct {
	Block DurationBlock `ecfg:"BLOCK"`
}

// DurationBlock holds Timeout.
type DurationBlock struct {
	Timeout Timeout `ecfg:"TIMEOUT"`
}

// Timeout is a time.Duration leaf.
type Timeout time.Duration

// Usage describes TIMEOUT.
func (Timeout) Usage() string { return "Request timeout" }

// Validate is a no-op.
func (Timeout) Validate() error { return nil }

// RatioCfg uses an unsupported complex scalar leaf.
type RatioCfg struct {
	Block RatioBlock `ecfg:"BLOCK"`
}

// RatioBlock holds Ratio.
type RatioBlock struct {
	R Ratio `ecfg:"RATIO"`
}

// Ratio is complex128 (unsupported scalar).
type Ratio complex128

// Usage describes RATIO.
func (Ratio) Usage() string { return "Ratio" }

// Validate is a no-op.
func (Ratio) Validate() error { return nil }

// StrictLabel fails validation on a sentinel value.
type StrictLabel string

// Usage for strict label.
func (StrictLabel) Usage() string { return "Strict label" }

// Validate rejects the literal "reject".
func (l StrictLabel) Validate() error {
	if l == "reject" {
		return fmt.Errorf("rejected")
	}
	return nil
}

// StrictCfg uses StrictLabel.
type StrictCfg struct {
	Block StrictBlock `ecfg:"BLOCK"`
}

// StrictBlock holds StrictLabel.
type StrictBlock struct {
	Label StrictLabel `ecfg:"LABEL"`
}

// EmptyUsageLabel has whitespace-only usage text.
type EmptyUsageLabel string

// Usage returns empty text.
func (EmptyUsageLabel) Usage() string { return "   " }

// Validate is a no-op.
func (EmptyUsageLabel) Validate() error { return nil }

// EmptyUsageCfg uses EmptyUsageLabel.
type EmptyUsageCfg struct {
	Block EmptyUsageBlock `ecfg:"BLOCK"`
}

// EmptyUsageBlock holds EmptyUsageLabel.
type EmptyUsageBlock struct {
	Label EmptyUsageLabel `ecfg:"LABEL"`
}
