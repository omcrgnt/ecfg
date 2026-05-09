package ecfg

import (
	"unsafe"
)

type carrier struct {
	ptr     any
	value   any
	uptr    unsafe.Pointer
	efName  string
	efUsage string
}
