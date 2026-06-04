package ecfg

import "github.com/omcrgnt/ecfg/internal/ecfgtool"

type parseOption func(*parseOptions)

type parseOptions struct {
	prefix string
}

// Parse loads *T from environment variables.
func Parse[T any](opts ...parseOption) (*T, error) {
	cfg := parseOptions{}
	for _, o := range opts {
		o(&cfg)
	}
	return ecfgtool.Parse[T](ecfgtool.Options{Prefix: cfg.prefix})
}

// WithPrefix adds a prefix to all environment variable names.
func WithPrefix(prefix string) parseOption {
	return func(o *parseOptions) {
		o.prefix = prefix
	}
}
