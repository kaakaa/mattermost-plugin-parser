package main

import (
	"go/ast"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/kaakaa/mattermost-plugin-parser/server/store"
)

func main() {
	if len(os.Args) != 5 {
		log.Fatal("Must be 4 arguments. [repository_url] [commit_id] [repo_dir] [commited_at]")
	}
	repository := os.Args[1]
	commitId := os.Args[2]
	basePath := os.Args[3]
	commitedAt := os.Args[4]

	// init database
	db, err := store.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.InsertRepository(db, repository, commitId, commitedAt); err != nil {
		log.Fatal(err)
	}

	// Parse Mattermost plugin package
	apiFuncs, hookFuncs, err := parseMattermostPluginInterface()
	if err != nil {
		log.Fatal(err)
	}

	// Parse target project
	fset := token.NewFileSet()
	var parsed []*ast.File
	files, err := parse(fset, basePath)
	if err != nil {
		log.Fatalf("parse error: %v", err)
	}
	parsed = append(parsed, files...)

	// Find and store calling plugin API
	called, err := findCalledPuluginAPI(fset, parsed)
	if err != nil {
		log.Fatalf("find error: %v", err)
	}
	for _, v := range called {
		for _, obj := range apiFuncs {
			if obj.Name == v.name {
				log.Printf("Call API Function: %s", v.name)
				p := strings.Replace(v.position.Filename, basePath+"/", "", -1)
				store.InsertUsage(db, commitId, v.name, p, v.position.Line, "server.api")
			}
		}
	}

	// Find and store declaring plugin Hooks
	decls, err := findDeclaredHooks(fset, parsed)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range decls {
		for _, obj := range hookFuncs {
			if obj.Name == v.name {
				log.Printf("Declare Plugin Hooks: %s", v.name)
				p := strings.Replace(v.position.Filename, basePath+"/", "", -1)
				store.InsertUsage(db, commitId, v.name, p, v.position.Line, "server.hooks")
			}
		}
	}

}
