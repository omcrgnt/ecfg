// Package ecfg loads environment variables into config values stored in a side-registry.
//
// Process defaults: env prefix [DefaultPrefix], custom tag key [DefaultTagKey].
// Override with [SetPrefix] and [SetTagKey] before [LoadEnv].
//
// [CatalogSegment] reads the catalog wire field tag named [TagKey].
//
// Standalone struct parsing without a registry is [github.com/omcrgnt/ecfg/config.Parse].
package ecfg
