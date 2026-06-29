package ecfg

import "sync"

const (
	// DefaultPrefix is the env key prefix when [SetPrefix] was not called.
	DefaultPrefix = "APP"
	// DefaultTagKey is the registry custom tag key when [SetTagKey] was not called.
	DefaultTagKey = "ecfg"
)

var (
	cfgMu    sync.RWMutex
	prefix   = DefaultPrefix
	tagKey   = DefaultTagKey
)

// SetPrefix sets the env key prefix for [LoadEnv]. Empty string restores [DefaultPrefix].
func SetPrefix(p string) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	if p == "" {
		prefix = DefaultPrefix
		return
	}
	prefix = p
}

// SetTagKey sets the registry custom tag key for [LoadEnv]. Empty string restores [DefaultTagKey].
func SetTagKey(key string) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	if key == "" {
		tagKey = DefaultTagKey
		return
	}
	tagKey = key
}

// Prefix returns the current env key prefix for [LoadEnv].
func Prefix() string {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return prefix
}

// TagKey returns the current registry custom tag key for [LoadEnv] and side-registry [AddWithCustomTag].
func TagKey() string {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return tagKey
}

// ResetForTest restores [DefaultPrefix] and [DefaultTagKey].
func ResetForTest() {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	prefix = DefaultPrefix
	tagKey = DefaultTagKey
}
