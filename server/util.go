package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Parse go files
// - Ignore `*_test.go` file
// - Ignore `vendor` and `build` directory
func parse(fset *token.FileSet, root string) ([]*ast.File, error) {
	var ret []*ast.File

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == "vendor" || info.Name() == "build" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}
		pkg, err := parser.ParseFile(fset, path, nil, parser.Mode(0))
		if err != nil {
			return err
		}
		ret = append(ret, pkg)
		return nil
	}); err != nil {
		return nil, err
	}
	return ret, nil
}
