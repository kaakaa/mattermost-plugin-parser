package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	ApiFileURL   = "https://raw.githubusercontent.com/mattermost/mattermost-server/master/plugin/api.go"
	HooksFileURL = "https://raw.githubusercontent.com/mattermost/mattermost-server/master/plugin/hooks.go"
)

func parseMattermostPluginInterface() ([]*ast.Object, []*ast.Object, error) {
	apiFuncs, err := parseInterface(ApiFileURL, "API")
	if err != nil {
		return nil, nil, err
	}
	hooks, err := parseInterface(HooksFileURL, "Hooks")
	if err != nil {
		return nil, nil, err
	}
	return apiFuncs, hooks, nil
}

func parseInterface(url, iName string) ([]*ast.Object, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to request %v", resp.Status)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	if err := ioutil.WriteFile(f.Name(), b, os.ModePerm); err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, f.Name(), nil, parser.Mode(0))
	if err != nil {
		return nil, err
	}

	var ret []*ast.Object
	ast.Inspect(file, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.GenDecl:
			decl := n.(*ast.GenDecl)
			if decl.Tok != token.TYPE {
				return true
			}
			for _, spec := range decl.Specs {
				tspec := spec.(*ast.TypeSpec)
				if tspec.Name.Name != iName {
					return true
				}
				iface, ok := tspec.Type.(*ast.InterfaceType)
				if !ok {
					return true
				}
				for _, field := range iface.Methods.List {
					if len(field.Names) != 1 {
						return true
					}
					ident := field.Names[0]
					if ident.Obj.Kind != ast.Fun || !ident.IsExported() {
						return true
					}
					ret = append(ret, ident.Obj)
				}
			}
		}
		return true
	})
	return ret, nil
}
