package ecfgtool

import "errors"

var (
	ErrNotStruct             = errors.New("ecfg: target is not a struct")
	ErrRootFieldKind         = errors.New("ecfg: root field must be struct or *struct")
	ErrNestedBlockAtDepth1   = errors.New("ecfg: nested struct block at depth 1 is not allowed")
	ErrUnsupportedContainer  = errors.New("ecfg: slice and map are not supported")
	ErrMissingEcfgTag        = errors.New("ecfg: root field requires ecfg tag")
	ErrDuplicateEnvKey       = errors.New("ecfg: duplicate env key")
	ErrMissingEnv            = errors.New("ecfg: missing env variable")
	ErrEmptyEnv              = errors.New("ecfg: empty env variable")
	ErrInvalidProtoWrapper   = errors.New("ecfg: invalid proto wrapper")
	ErrIncompatibleLeaf      = errors.New("ecfg: leaf type must be proto message or implement Usage and Validator")
	ErrMissingUsage          = errors.New("ecfg: missing usage")
	ErrEmptyUsage            = errors.New("ecfg: empty usage")
)
