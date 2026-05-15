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

func (*typesProvider) schemaProvider() {}

// NewTypesProvider loads a struct type from pkgPath via go/packages and returns a SchemaProvider.
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
	return isProtoTypesType(f.f.Type())
}
func (f *typesField) Value() (reflect.Value, reflect.StructField, error) {
	return reflect.Value{}, reflect.StructField{}, ErrNoRuntimeValue
}

func (f *typesField) GetProvider() (Provider, error) {
	st, ok := structTypesType(f.f.Type())
	if !ok {
		return nil, fmt.Errorf("walker: field %s is not a struct", f.f.Name())
	}
	return &typesProvider{st: st, structName: types.TypeString(f.f.Type(), nil), pkg: f.pkg}, nil
}

func (f *typesField) ElemProvider() (Provider, error) {
	st, named, ok := elemStructTypesNamed(f.f.Type())
	if !ok {
		return nil, ErrNotContainer
	}
	if isProtoTypesType(named) {
		return nil, ErrNotContainer
	}
	name := named.Obj().Name()
	if name == "" {
		name = types.TypeString(named, nil)
	}
	return &typesProvider{st: st, structName: name, pkg: f.pkg}, nil
}

func elemStructTypesNamed(t types.Type) (*types.Struct, *types.Named, bool) {
	t = types.Unalias(t)
	for {
		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
			continue
		}
		break
	}
	switch ut := t.Underlying().(type) {
	case *types.Slice:
		t = ut.Elem()
	case *types.Array:
		t = ut.Elem()
	case *types.Map:
		t = ut.Elem()
	default:
		return nil, nil, false
	}
	for {
		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
			continue
		}
		break
	}
	named, ok := t.(*types.Named)
	if !ok {
		return nil, nil, false
	}
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil, nil, false
	}
	return st, named, true
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

// underlyingKind maps go/types to reflect.Kind, preserving int/uint width.
func underlyingKind(t types.Type) reflect.Kind {
	for {
		if ptr, ok := t.Underlying().(*types.Pointer); ok {
			t = ptr.Elem()
			continue
		}
		break
	}

	switch ut := t.Underlying().(type) {
	case *types.Basic:
		return basicKind(ut.Kind())
	case *types.Struct:
		return reflect.Struct
	case *types.Slice:
		return reflect.Slice
	case *types.Array:
		return reflect.Array
	case *types.Map:
		return reflect.Map
	default:
		return reflect.Invalid
	}
}

func basicKind(k types.BasicKind) reflect.Kind {
	switch k {
	case types.Bool:
		return reflect.Bool
	case types.Int:
		return reflect.Int
	case types.Int8:
		return reflect.Int8
	case types.Int16:
		return reflect.Int16
	case types.Int32:
		return reflect.Int32
	case types.Int64:
		return reflect.Int64
	case types.Uint:
		return reflect.Uint
	case types.Uint8:
		return reflect.Uint8
	case types.Uint16:
		return reflect.Uint16
	case types.Uint32:
		return reflect.Uint32
	case types.Uint64:
		return reflect.Uint64
	case types.Float32:
		return reflect.Float32
	case types.Float64:
		return reflect.Float64
	case types.String:
		return reflect.String
	case types.UnsafePointer:
		return reflect.UnsafePointer
	default:
		return reflect.Invalid
	}
}
