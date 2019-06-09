package main

import (
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
)

type called struct {
	name     string
	position token.Position
}

func findCalledPuluginAPI(fset *token.FileSet, parsed []*ast.File) ([]called, error) {
	imp := importer.ForCompiler(fset, "source", nil)
	// imp := importer.Default()
	conf := types.Config{Importer: imp}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	_, err := conf.Check("main", fset, parsed, info)
	if err != nil {
		return nil, err
	}

	var ret []called
	for _, v := range parsed {
		ast.Inspect(v, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				var id *ast.Ident
				switch fun := call.Fun.(type) {
				case *ast.SelectorExpr:
					t := info.Types[fun.X].Type
					if t == nil {
						return true
					}
					if t.String() != "github.com/mattermost/mattermost-server/plugin.API" {
						return true
					}
					id = fun.Sel
				default:
					return true
				}
				p := fset.Position(id.Pos())
				ret = append(ret, called{
					name:     id.Name,
					position: p,
				})
			}
			return true
		})
	}
	return ret, nil
}
