// Package walk provides generic struct tree traversal via reflect or go/types.
//
// It is policy-agnostic: no ecfg tags, env keys, or validation. Import it when
// you need struct walking only; for configuration from environment variables use
// the root ecfg package.
//
// Exported API is intentionally small. Implementation details and helpers are
// unexported.
package walk
