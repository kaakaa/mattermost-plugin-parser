package main

import (
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
)

type FuncCall struct {
	f *ast.FuncDecl
	t *ast.Ident
}

func findDeclaredHooks(fset *token.FileSet, parsed []*ast.File) ([]called, error) {
	imp := importer.Default()
	conf := types.Config{Importer: imp}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
	}
	_, err := conf.Check("main", fset, parsed, info)
	if err != nil {
		return nil, err
	}

	var pluginTypes []*ast.Ident
	var funcCalls []FuncCall
	for _, v := range parsed {
		ast.Inspect(v, func(n ast.Node) bool {
			switch n.(type) {
			case *ast.GenDecl:
				decl := n.(*ast.GenDecl)
				// Filter to get only declared struct
				if decl.Tok != token.TYPE {
					return true
				}
				specs := decl.Specs
				if len(specs) == 0 {
					return true
				}
				// Check each struct (e.g.: type( Struct1 {} Struct2 {} ))
				for _, s := range specs {
					ts := s.(*ast.TypeSpec)
					if st, ok := ts.Type.(*ast.StructType); ok {
						// Get fields of struct
						for _, f := range st.Fields.List {
							switch f.Type.(type) {
							case *ast.SelectorExpr:
								t := f.Type.(*ast.SelectorExpr)
								if x, ok := t.X.(*ast.Ident); ok {
									// Filter to get only type which is struct inluding `plugin.MattermostPlugin` type
									if x.Name == "plugin" && t.Sel.Name == "MattermostPlugin" {
										pluginTypes = append(pluginTypes, ts.Name)
									}
								}
							}
						}
					}
				}
			case *ast.FuncDecl:
				f := n.(*ast.FuncDecl)
				// Ignore method that doesn't have any recievers
				if f.Recv == nil {
					return true
				}

				receivers := f.Recv.List
				// Ignore unexpected cases
				if len(receivers) != 1 {
					return true
				}

				recv := receivers[0]
				if expr, ok := recv.Type.(*ast.StarExpr); ok {
					if ident, ok := expr.X.(*ast.Ident); ok {
						funcCalls = append(funcCalls, FuncCall{
							f: f,
							t: ident,
						})
					}
				}
			}
			return true
		})
	}

	var ret []called
	for _, v := range funcCalls {
		for _, tp := range pluginTypes {
			if info.Uses[v.t] == info.Defs[tp] {
				position := fset.Position(v.t.Pos())
				ret = append(ret, called{
					name:     v.f.Name.Name,
					position: position,
				})
			}
		}
	}
	return ret, nil
}
