package walker

import (
	"fmt"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/packages"
)

type typesProvider struct {
	st         *types.Struct
	structName string
	pkg        *packages.Package
}

func NewTypesProvider(pkgPath, structName string) (Provider, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports | packages.NeedDeps | packages.NeedSyntax | packages.NeedTypesInfo}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil || len(pkgs) == 0 {
		return nil, fmt.Errorf("load error: %v", err)
	}

	pkg := pkgs[0]
	obj := pkg.Types.Scope().Lookup(structName)
	st, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("not a struct")
	}

	return &typesProvider{st: st, structName: structName, pkg: pkg}, nil
}

func (p *typesProvider) GetFields() ([]Field, error) {
	var fields []Field
	for i := 0; i < p.st.NumFields(); i++ {
		f := p.st.Field(i)
		if !f.Exported() {
			continue
		}
		fields = append(fields, &typesField{f: f, t: p.st.Tag(i), pkg: p.pkg})
	}
	return fields, nil
}

func (p *typesProvider) EntryName() string { return p.structName }

type typesField struct {
	f   *types.Var
	t   string
	pkg *packages.Package
}

func (f *typesField) Name() string          { return f.f.Name() }
func (f *typesField) Tag(key string) string { return reflect.StructTag(f.t).Get(key) }
func (f *typesField) Kind() reflect.Kind {
	return underlyingKind(f.f.Type())
}
func (f *typesField) IsStruct() bool { return f.Kind() == reflect.Struct }
func (f *typesField) IsProto() bool {
	named, ok := f.f.Type().(*types.Named)
	if !ok {
		return false
	}
	for i := 0; i < named.NumMethods(); i++ {
		if named.Method(i).Name() == "ProtoMessage" {
			return true
		}
	}
	return false
}
func (f *typesField) GetProvider() (Provider, error) {
	st, ok := structTypesType(f.f.Type())
	if !ok {
		return nil, fmt.Errorf("walker: field %s is not a struct", f.f.Name())
	}
	return &typesProvider{st: st, structName: types.TypeString(f.f.Type(), nil), pkg: f.pkg}, nil
}

func structTypesType(t types.Type) (*types.Struct, bool) {
	for {
		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
			continue
		}
		break
	}
	st, ok := t.Underlying().(*types.Struct)
	return st, ok
}

func underlyingKind(t types.Type) reflect.Kind {
	// 1. Разворачиваем указатели (учитываем многоуровневые типа ***int)
	for {
		if ptr, ok := t.Underlying().(*types.Pointer); ok {
			t = ptr.Elem()
			continue
		}
		break
	}

	// 2. Определяем Kind на основе Underlying типа
	switch ut := t.Underlying().(type) {
	case *types.Basic:
		switch ut.Kind() {
		case types.String:
			return reflect.String
		case types.Bool:
			return reflect.Bool
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
			return reflect.Int // Для простоты сводим всё к Int
		case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return reflect.Uint
		case types.Float32, types.Float64:
			return reflect.Float64
		}
	case *types.Struct:
		return reflect.Struct
	case *types.Slice:
		return reflect.Slice
	case *types.Array:
		return reflect.Array
	case *types.Map:
		return reflect.Map
	}

	return reflect.Invalid
}
