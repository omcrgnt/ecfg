package ecfg

import "github.com/omcrgnt/ecfg/internal/ecfgtool"

var (
	ErrNotStruct            = ecfgtool.ErrNotStruct
	ErrRootFieldKind        = ecfgtool.ErrRootFieldKind
	ErrNestedBlockAtDepth1  = ecfgtool.ErrNestedBlockAtDepth1
	ErrUnsupportedContainer = ecfgtool.ErrUnsupportedContainer
	ErrMissingEcfgTag       = ecfgtool.ErrMissingEcfgTag
	ErrDuplicateEnvKey      = ecfgtool.ErrDuplicateEnvKey
	ErrMissingEnv           = ecfgtool.ErrMissingEnv
	ErrEmptyEnv             = ecfgtool.ErrEmptyEnv
	ErrInvalidProtoWrapper  = ecfgtool.ErrInvalidProtoWrapper
	ErrIncompatibleLeaf     = ecfgtool.ErrIncompatibleLeaf
	ErrMissingUsage         = ecfgtool.ErrMissingUsage
	ErrEmptyUsage           = ecfgtool.ErrEmptyUsage
)
