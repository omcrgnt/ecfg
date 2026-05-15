package walker

import (
	"go/types"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
	"golang.org/x/tools/go/packages"
)

var protoMessageReflect = reflect.TypeOf((*proto.Message)(nil)).Elem()

var (
	protoMessageIface *types.Interface
	protoMessageOnce  sync.Once
)

func protoMessageInterface() *types.Interface {
	protoMessageOnce.Do(func() {
		pkgs, err := packages.Load(&packages.Config{
			Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo,
		}, "google.golang.org/protobuf/proto")
		if err != nil || len(pkgs) == 0 || pkgs[0].Types == nil {
			return
		}
		obj := pkgs[0].Types.Scope().Lookup("Message")
		if obj == nil {
			return
		}
		iface, ok := obj.Type().Underlying().(*types.Interface)
		if ok {
			protoMessageIface = iface
		}
	})
	return protoMessageIface
}

func isProtoReflectType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	return reflect.PointerTo(t).Implements(protoMessageReflect)
}

func isProtoTypesType(t types.Type) bool {
	if iface := protoMessageInterface(); iface != nil {
		if isProtoTypesImplements(t, iface) {
			return true
		}
	}
	return isProtoTypesMethodSet(t)
}

func isProtoTypesImplements(t types.Type, iface *types.Interface) bool {
	st := structTypesNamed(t)
	if st == nil {
		return false
	}
	return types.Implements(types.NewPointer(st), iface)
}

// isProtoTypesMethodSet is a fallback when proto.Message is not loaded in types.Importer.
// Pointer-receiver methods are not on the struct method set, so we check *T.
func isProtoTypesMethodSet(t types.Type) bool {
	st := structTypesNamed(t)
	if st == nil {
		return false
	}
	mset := types.NewMethodSet(types.NewPointer(st))
	for i := 0; i < mset.Len(); i++ {
		switch mset.At(i).Obj().Name() {
		case "ProtoReflect", "ProtoMessage":
			return true
		}
	}
	return false
}

func structTypesNamed(t types.Type) *types.Named {
	t = types.Unalias(t)
	for {
		if p, ok := t.(*types.Pointer); ok {
			t = p.Elem()
			continue
		}
		break
	}
	named, ok := t.(*types.Named)
	if !ok {
		return nil
	}
	if _, ok := named.Underlying().(*types.Struct); !ok {
		return nil
	}
	return named
}
