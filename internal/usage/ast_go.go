package usage

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

var errMissingUsage = errors.New("ecfg: missing usage")

// GoUsageFromAST finds Usage() string return for a named type in pkgPath.
func GoUsageFromAST(pkgPath, typeName string) (string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedFiles | packages.NeedDeps,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", errMissingUsage
	}
	pkg := pkgs[0]
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Name.Name != "Usage" {
				continue
			}
			if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
				continue
			}
			recvName, ok := recvTypeName(fn.Recv.List[0].Type)
			if !ok || recvName != typeName {
				continue
			}
			if lit := stringReturn(fn.Body); lit != "" {
				return lit, nil
			}
		}
	}
	return "", errMissingUsage
}

func recvTypeName(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, true
	case *ast.StarExpr:
		if id, ok := t.X.(*ast.Ident); ok {
			return id.Name, true
		}
	}
	return "", false
}

func stringReturn(body *ast.BlockStmt) string {
	if body == nil || len(body.List) != 1 {
		return ""
	}
	ret, ok := body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return ""
	}
	switch v := ret.Results[0].(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			return strings.Trim(v.Value, `"`)
		}
	}
	return ""
}
