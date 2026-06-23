package ecfgtool

// Proto descriptors must be linked into binaries that call [CollectTemplateEntries]
// on configs with proto wrapper leaves (codegen, ecfg-gen). Runtime Apply uses reflect
// on live values and does not depend on this import.
import _ "github.com/omcrgnt/proto/gen/go/common/v1"
