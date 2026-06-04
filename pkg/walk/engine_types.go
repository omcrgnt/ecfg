package walk

import (
	"fmt"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/packages"
)

type engineTypes struct {
	st         *types.Struct
	structName string
	pkg        *types.Package
}

// NewEngineTypes loads pkgPath.typeName and returns an Engine.
func NewEngineTypes(pkgPath, typeName string) (Engine, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedModule,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("walk: package %s not found", pkgPath)
	}
	pkg := pkgs[0]
	for _, e := range pkg.Errors {
		if e.Kind == packages.TypeError || e.Kind == packages.ParseError {
			return nil, fmt.Errorf("walk: %s: %s", pkgPath, e.Msg)
		}
	}
	obj := pkg.Types.Scope().Lookup(typeName)
	if obj == nil {
		return nil, fmt.Errorf("walk: type %s not found in %s", typeName, pkgPath)
	}
	st, ok := types.Unalias(obj.Type()).Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("walk: %s is not a struct", typeName)
	}
	return &engineTypes{st: st, structName: typeName, pkg: pkg.Types}, nil
}

func (e *engineTypes) Fields() ([]FieldDesc, error) {
	var out []FieldDesc
	for i := 0; i < e.st.NumFields(); i++ {
		f := e.st.Field(i)
		if !f.Exported() {
			continue
		}
		out = append(out, FieldDesc{
			Name:      f.Name(),
			Tag:       e.st.Tag(i),
			TypesType: f.Type(),
		})
	}
	return out, nil
}

func (e *engineTypes) Child(desc FieldDesc) (Engine, error) {
	st, ok := structTypesType(desc.TypesType)
	if !ok {
		return nil, nil
	}
	return &engineTypes{st: st, structName: desc.Name, pkg: e.pkg}, nil
}

func structTypesType(t types.Type) (*types.Struct, bool) {
	t = types.Unalias(t)
	for {
		if p, ok := t.(*types.Pointer); ok {
			t = p.Elem()
			continue
		}
		break
	}
	st, ok := t.Underlying().(*types.Struct)
	return st, ok
}

// ReflectKind approximates reflect.Kind for a types.Type.
func ReflectKind(t types.Type) reflect.Kind {
	return reflectKind(t)
}

func reflectKind(t types.Type) reflect.Kind {
	t = types.Unalias(t)
	switch t := t.(type) {
	case *types.Pointer:
		elem := types.Unalias(t.Elem())
		if _, ok := elem.Underlying().(*types.Struct); ok {
			return reflect.Ptr
		}
		return underlyingReflectKind(elem)
	case *types.Slice:
		return reflect.Slice
	case *types.Map:
		return reflect.Map
	case *types.Basic:
		return basicReflectKind(t)
	case *types.Struct:
		return reflect.Struct
	case *types.Named:
		return reflectKind(t.Underlying())
	default:
		return reflect.Invalid
	}
}

func underlyingReflectKind(t types.Type) reflect.Kind {
	switch types.Unalias(t).(type) {
	case *types.Basic:
		return basicReflectKind(types.Unalias(t).(*types.Basic))
	case *types.Struct:
		return reflect.Struct
	case *types.Named:
		return reflectKind(t)
	default:
		return reflect.Invalid
	}
}

func basicReflectKind(b *types.Basic) reflect.Kind {
	switch b.Kind() {
	case types.Bool:
		return reflect.Bool
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		return reflect.Int
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
		return reflect.Uint
	case types.Float32:
		return reflect.Float32
	case types.Float64:
		return reflect.Float64
	case types.String:
		return reflect.String
	default:
		return reflect.Invalid
	}
}
