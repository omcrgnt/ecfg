package config

import "github.com/omcrgnt/ecfg/internal/ecfgtool"

type Option func(*options)

type options struct {
	prefix string
}

// Parse loads *T from environment variables.
func Parse[T any](opts ...Option) (*T, error) {
	cfg := options{}
	for _, o := range opts {
		o(&cfg)
	}
	return ecfgtool.Parse[T](ecfgtool.Options{Prefix: cfg.prefix})
}

// WithPrefix adds a prefix to all environment variable names.
func WithPrefix(prefix string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}
