// Package ecfg loads configuration from environment variables into a struct.
//
// Struct fields use the "ecfg" tag for the env name segment; nested structs build
// keys like MODULE_INNER_KEY. Optional WithPrefix prepends a global prefix.
package ecfg
