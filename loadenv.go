package ecfg

import (
	"fmt"

	"github.com/omcrgnt/ecfg/internal/ecfgtool"
	"github.com/omcrgnt/res"
	"github.com/omcrgnt/res/unique"
)

// LoadEnv initializes config field values from the environment for each entry in reg
// that has a custom tag [TagKey] (tag value = ecfg segment, e.g. SERVICE_ITEM).
func LoadEnv(reg *unique.Registry) error {
	if reg == nil {
		return fmt.Errorf("ecfg: nil registry")
	}
	return ecfgtool.LoadRegistry(func(yield func(ecfgtool.RegistryEntry) bool) {
		reg.WalkEntries(func(e res.Entry) bool {
			return yield(e)
		})
	}, TagKey(), ecfgtool.Options{Prefix: Prefix()})
}
