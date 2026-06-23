package ecfg

import (
	"fmt"

	"github.com/omcrgnt/builder"
	"github.com/omcrgnt/ecfg/internal/ecfgtool"
	"github.com/omcrgnt/res"
)

// Apply loads environment variables into config specs registered by [builder.Seed].
func Apply(reg res.Registry, appResources any, opts ...parseOption) error {
	if reg == nil {
		return fmt.Errorf("ecfg: nil registry")
	}

	seedMap, ok := builder.SeedMapFor(reg)
	if !ok {
		return fmt.Errorf("ecfg: no seed map for registry")
	}

	cfg := parseOptions{}
	for _, o := range opts {
		o(&cfg)
	}

	return ecfgtool.ApplySeeded(appResources, seedMap, ecfgtool.Options{Prefix: cfg.prefix})
}
