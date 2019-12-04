package mattermost

import (
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

type MattermostPluginFuncs struct {
	PluginFuncs map[*types.Named][]*types.Func
	HooksFuncs  []*types.Func
}

var PluginAPIAnalyzer = &analysis.Analyzer{
	Name: "PluginAPIAnalyzer",
	Doc:  Doc,
	Run: func(pass *analysis.Pass) (interface{}, error) {
		res := &MattermostPluginFuncs{
			PluginFuncs: map[*types.Named][]*types.Func{},
			HooksFuncs:  []*types.Func{},
		}
		findMattermostPluginAPI(pass, res)
		return res, nil
	},
	ResultType: reflect.TypeOf((*MattermostPluginFuncs)(nil)),
}

const Doc = "Mattermost Plugin API Analyzer"

func main() {
	singlechecker.Main(
		PluginAPIAnalyzer,
	)
}

func findMattermostPluginAPI(pass *analysis.Pass, res *MattermostPluginFuncs) {
	for _, imp := range pass.Pkg.Imports() {
		if imp.Path() == "github.com/mattermost/mattermost-server/plugin" {
			// Expected MattermostPlugin struct
			//
			// type MattermostPlugin struct {
			//   Ident1 Type1	// TypeN must be interface, not struct.
			//   Ident2 Type2
			//   ...
			// }
			pType := imp.Scope().Lookup("MattermostPlugin").Type().(*types.Named)
			// Get as struct for finding their fields
			pStruct := pType.Underlying().(*types.Struct)
			for i := 0; i < pStruct.NumFields(); i++ {
				// Get as named struct  from field
				field := pStruct.Field(i)
				namedType, ok := field.Type().(*types.Named)
				if !ok {
					continue
				}
				if _, ok = res.PluginFuncs[namedType]; !ok {
					res.PluginFuncs[namedType] = []*types.Func{}
				}
				// Get as interface for finding their funcs
				iface, ok := namedType.Underlying().(*types.Interface)
				if !ok {
					continue
				}
				for i := 0; i < iface.NumMethods(); i++ {
					res.PluginFuncs[namedType] = append(res.PluginFuncs[namedType], iface.Method(i))
				}
			}

			// Find funcs in Hooks
			hooksType, ok := imp.Scope().Lookup("Hooks").Type().(*types.Named)
			if !ok {
				return
			}
			iface, ok := hooksType.Underlying().(*types.Interface)
			if !ok {
				return
			}
			for i := 0; i < iface.NumMethods(); i++ {
				res.HooksFuncs = append(res.HooksFuncs, iface.Method(i))
			}
			return
		}
	}
}
