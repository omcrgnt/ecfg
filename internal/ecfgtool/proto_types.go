package ecfgtool

import (
	"go/types"
)

func isProtoTypesType(t types.Type) bool {
	named := typesNamed(t)
	if named == nil {
		return false
	}
	mset := types.NewMethodSet(types.NewPointer(named))
	for i := 0; i < mset.Len(); i++ {
		switch mset.At(i).Obj().Name() {
		case "ProtoReflect", "ProtoMessage":
			return true
		}
	}
	return false
}

func typesNamed(t types.Type) *types.Named {
	t = types.Unalias(t)
	for {
		if p, ok := t.(*types.Pointer); ok {
			t = p.Elem()
			continue
		}
		break
	}
	n, ok := t.(*types.Named)
	if !ok {
		return nil
	}
	if _, ok := n.Underlying().(*types.Struct); !ok {
		return nil
	}
	return n
}
