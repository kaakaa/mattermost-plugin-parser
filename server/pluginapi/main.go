package pluginapi

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/kaakaa/mattermost-plugin-parser/server/mattermost"
	"golang.org/x/tools/go/analysis"
)

// Target represents package root (e.g.: github.com/mattermost/mattermost-plugin-github)
var Target string

var Analyzer = &analysis.Analyzer{
	Name: "MattermostPluginAnalyzer",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		mattermost.PluginAPIAnalyzer,
	},
	ResultType: reflect.TypeOf((*PluginAPIUsages)(nil)),
}

const Doc = "Mattermost Plugin Analyzer"

type PluginAPIUsages struct {
	APIs  []Usage
	Hooks []Usage
}

type Usage struct {
	Name     string
	Path     string
	Position token.Position
}

func run(pass *analysis.Pass) (interface{}, error) {
	mmfuncs, ok := pass.ResultOf[mattermost.PluginAPIAnalyzer].(*mattermost.MattermostPluginFuncs)
	if !ok {
		return nil, fmt.Errorf("Failed to get MattermostPluginAPI")
	}

	res := &PluginAPIUsages{
		APIs:  []Usage{},
		Hooks: []Usage{},
	}
	findInvokingAPI(pass, mmfuncs, res)
	findImplementingHooks(pass, mmfuncs, res)
	return res, nil
}

func findInvokingAPI(pass *analysis.Pass, mmfuncs *mattermost.MattermostPluginFuncs, res *PluginAPIUsages) {
	for _, v := range pass.Files {
		ast.Inspect(v, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				var id *ast.Ident
				switch fun := call.Fun.(type) {
				case *ast.SelectorExpr:
					t := pass.TypesInfo.Types[fun.X].Type
					if t == nil {
						return true
					}
					nt, _ := t.(*types.Named)
					_, ok := mmfuncs.PluginFuncs[nt]
					if !ok {
						return true
					}
					id = fun.Sel
				default:
					return true
				}

				path := filepath.Join(pass.Pkg.Path(), filepath.Base(pass.Fset.File(n.Pos()).Name())) // PKG_NAME + FILE_NAME

				var fpath string
				tmpWorkdir := "github.com/kaakaa/mattermost-plugin-parser/.work"
				if strings.HasPrefix(path, tmpWorkdir) {
					fpath, _ = filepath.Rel(tmpWorkdir, path)
				} else {
					fpath, _ = filepath.Rel(Target, path)
				}
				res.APIs = append(res.APIs, Usage{
					Name:     id.Name,
					Path:     fpath,
					Position: pass.Fset.Position(id.Pos()),
				})
			}
			return true
		})
	}
}

type FuncCall struct {
	f *ast.FuncDecl
	t *ast.Ident
}

func findImplementingHooks(pass *analysis.Pass, mmfuncs *mattermost.MattermostPluginFuncs, res *PluginAPIUsages) {
	var pluginTypes []*ast.Ident
	var funcCalls []FuncCall
	for _, v := range pass.Files {
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
				// TODO: only StarExpr? needs TypeExpr or ValueExpr?
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

	for _, v := range funcCalls {
		for _, tp := range pluginTypes {
			if pass.TypesInfo.Uses[v.t] == pass.TypesInfo.Defs[tp] {
				for _, fn := range mmfuncs.HooksFuncs {
					if fn.Name() == v.f.Name.String() {
						path := filepath.Join(pass.Pkg.Path(), filepath.Base(pass.Fset.File(v.f.Pos()).Name())) // PKG_NAME + FILE_NAME
						fpath, _ := filepath.Rel(Target, path)
						res.Hooks = append(res.Hooks, Usage{
							Name:     v.f.Name.Name,
							Path:     fpath,
							Position: pass.Fset.Position(v.f.Pos()),
						})
					}
				}
			}
		}
	}
}
