package ecfgtool

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/omcrgnt/ecfg/pkg/walk"
)

func loadRootStruct(pkgPath, typeName string) (*packages.Package, *types.Struct, map[string]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedDeps | packages.NeedModule,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(pkgs) == 0 {
		return nil, nil, nil, fmt.Errorf("walk: package %s not found", pkgPath)
	}
	pkg := pkgs[0]
	for _, e := range pkg.Errors {
		if e.Kind == packages.TypeError || e.Kind == packages.ParseError {
			return nil, nil, nil, fmt.Errorf("walk: %s: %s", pkgPath, e.Msg)
		}
	}
	obj := pkg.Types.Scope().Lookup(typeName)
	if obj == nil {
		return nil, nil, nil, fmt.Errorf("walk: type %s not found in %s", typeName, pkgPath)
	}
	st, ok := types.Unalias(obj.Type()).Underlying().(*types.Struct)
	if !ok {
		return nil, nil, nil, fmt.Errorf("walk: %s is not a struct", typeName)
	}
	all := indexPackages(pkg)
	return pkg, st, all, nil
}

func indexPackages(root *packages.Package) map[string]*packages.Package {
	all := make(map[string]*packages.Package)
	var visit func(*packages.Package)
	visit = func(p *packages.Package) {
		if p == nil || all[p.PkgPath] != nil {
			return
		}
		all[p.PkgPath] = p
		for _, imp := range p.Imports {
			visit(imp)
		}
	}
	visit(root)
	return all
}

func engineForRootField(pkgs map[string]*packages.Package, field *types.Var) (walk.Engine, error) {
	if specT := builderSpecTypesType(pkgs, field.Type()); specT != nil {
		return engineFromTypesType(specT)
	}
	return engineFromTypesType(field.Type())
}

func engineFromTypesType(t types.Type) (walk.Engine, error) {
	st, ok := structFromTypes(t)
	if !ok {
		return nil, fmt.Errorf("ecfg: %s is not a struct block", t)
	}
	pkg := typesPkg(t)
	if pkg == nil {
		return nil, fmt.Errorf("ecfg: package for %s not loaded", t)
	}
	return &engineTypesBlock{st: st, pkg: pkg}, nil
}

func typesPkg(t types.Type) *types.Package {
	if n := namedType(t); n != nil {
		return n.Obj().Pkg()
	}
	return nil
}

type engineTypesBlock struct {
	st  *types.Struct
	pkg *types.Package
}

func (e *engineTypesBlock) Fields() ([]walk.FieldDesc, error) {
	var out []walk.FieldDesc
	for i := 0; i < e.st.NumFields(); i++ {
		f := e.st.Field(i)
		if !f.Exported() {
			continue
		}
		out = append(out, walk.FieldDesc{
			Name:      f.Name(),
			Tag:       e.st.Tag(i),
			TypesType: f.Type(),
		})
	}
	return out, nil
}

func (e *engineTypesBlock) Child(desc walk.FieldDesc) (walk.Engine, error) {
	st, ok := structFromTypes(desc.TypesType)
	if !ok {
		return nil, nil
	}
	return &engineTypesBlock{st: st, pkg: e.pkg}, nil
}

func structFromTypes(t types.Type) (*types.Struct, bool) {
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

func builderSpecTypesType(pkgs map[string]*packages.Package, wire types.Type) types.Type {
	named := namedType(wire)
	if named == nil {
		return nil
	}
	typePkg := pkgs[named.Obj().Pkg().Path()]
	if typePkg == nil {
		typePkg = loadPackageSyntax(named.Obj().Pkg().Path())
		if typePkg != nil {
			pkgs[typePkg.PkgPath] = typePkg
		}
	}
	if typePkg == nil {
		return nil
	}
	for i := 0; i < named.NumMethods(); i++ {
		m := named.Method(i)
		if m.Name() != "BuildConfig" {
			continue
		}
		if specT := specTypeFromBuildConfigAST(typePkg, named, m); specT != nil {
			return specT
		}
	}
	return nil
}

func loadPackageSyntax(pkgPath string) *packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedDeps,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil || len(pkgs) == 0 {
		return nil
	}
	return pkgs[0]
}

func namedType(t types.Type) *types.Named {
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
	return n
}

func specTypeFromBuildConfigAST(pkg *packages.Package, recv *types.Named, fn *types.Func) types.Type {
	recvName := recv.Obj().Name()
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Name.Name != "BuildConfig" || fd.Body == nil {
				continue
			}
			if fd.Recv == nil || len(fd.Recv.List) != 1 {
				continue
			}
			if !receiverMatches(fd.Recv.List[0].Type, recvName) {
				continue
			}
			if t := firstConcreteBuilderReturn(pkg, fd.Body); t != nil {
				return t
			}
		}
	}
	return nil
}

func receiverMatches(expr ast.Expr, recvName string) bool {
	switch t := expr.(type) {
	case *ast.StarExpr:
		ident, ok := t.X.(*ast.Ident)
		return ok && ident.Name == recvName
	case *ast.Ident:
		return t.Name == recvName
	default:
		return false
	}
}

func firstConcreteBuilderReturn(pkg *packages.Package, body *ast.BlockStmt) types.Type {
	for _, stmt := range body.List {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok || len(ret.Results) == 0 {
			continue
		}
		if t := typeOfExpr(pkg, ret.Results[0]); t != nil {
			return t
		}
	}
	return nil
}

func typeOfExpr(pkg *packages.Package, expr ast.Expr) types.Type {
	switch e := expr.(type) {
	case *ast.UnaryExpr:
		if e.Op.String() == "&" {
			return typeOfExpr(pkg, e.X)
		}
	case *ast.CompositeLit:
		if t := pkg.TypesInfo.TypeOf(e.Type); t != nil {
			return t
		}
	case *ast.Ident:
		if obj := pkg.TypesInfo.ObjectOf(e); obj != nil {
			return obj.Type()
		}
	}
	return nil
}
