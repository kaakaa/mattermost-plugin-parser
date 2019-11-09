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
	/*
		imp := importer.ForCompiler(fset, runtime.Compiler, func(path string) (io.ReadCloser, error) {
			log.Println(path)
			file := filepath.Join(runtime.GOROOT(), "pkg", fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH), path+".a")
			log.Printf("file: %s", file)
			f, err := os.Open(file)
			if err != nil && os.IsNotExist(err) {
				log.Println("%%")
				pkg, err := build.ImportDir(path, build.IgnoreVendor)
				if err != nil {
					return nil, err
				}
				log.Printf("TEST: %#v", pkg)
				p := filepath.Join("..", "vendor", path)
				log.Println(p)
				f, err = os.Open(p)
				if err != nil {
					return nil, err
				}
			}
			return f, err
		})
	*/
	imp := importer.Default()
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
