package main

import (
	"go/ast"
	"go/token"
	"os"
	"strings"

	"github.com/kaakaa/mattermost-plugin-parser/server/store"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if len(os.Args) != 5 {
		log.WithFields(log.Fields{
			"given": os.Args,
		}).Fatal("Must be 4 arguments. [repository_url] [commit_id] [repo_dir] [commited_at].")
	}
	repository := os.Args[1]
	commitId := os.Args[2]
	basePath := os.Args[3]
	commitedAt := os.Args[4]

	log.WithFields(log.Fields{
		"repository": repository,
		"commitId":   commitId,
		"basePath":   basePath,
		"commitedAt": commitedAt,
	}).Info("mattermost-plugin-parser (server) starts.")

	// init database
	db, err := store.InitDB()
	if err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to init db.")
	}
	if err := store.InsertRepository(db, repository, commitId, commitedAt); err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to insert repository info to db.")
	}

	// Parse manifest file
	manifest, err := parseMattermostPluginManifest(basePath)
	if err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to parse plugin manifest.")
	}
	if err := store.InsertManifest(db, commitId, manifest); err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to insert plugin manifest.")
	}

	// Parse Mattermost plugin package
	apiFuncs, hookFuncs, err := parseMattermostPluginInterface()
	if err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to parse Mattermost Plugin interface.")
	}

	// Parse target project
	fset := token.NewFileSet()
	var parsed []*ast.File
	files, err := parse(fset, basePath)
	if err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to parse plugin go files.")
	}
	parsed = append(parsed, files...)

	// Find and store calling plugin API
	log.Info("Start finding usages invoking plugin API.")
	called, err := findCalledPuluginAPI(fset, parsed)
	if err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Error occurs when finding invoking plugin API.")
	}
	for _, v := range called {
		for _, obj := range apiFuncs {
			if obj.Name == v.name {
				log.WithFields(log.Fields{"method": v.name}).Info("  Found")
				p := strings.Replace(v.position.Filename, basePath+"/", "", -1)
				if err := store.InsertUsage(db, commitId, v.name, p, v.position.Line, "server.api"); err != nil {
					log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to insert usage invoking plugin API.")
				}
			}
		}
	}

	// Find and store declaring plugin Hooks
	log.Info("Start finding usages implementing plugin hooks.")
	decls, err := findDeclaredHooks(fset, parsed)
	if err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Error occurs when finding implementing plugin hooks.")
	}
	for _, v := range decls {
		for _, obj := range hookFuncs {
			if obj.Name == v.name {
				log.WithFields(log.Fields{"method": v.name}).Info("  Found")
				p := strings.Replace(v.position.Filename, basePath+"/", "", -1)
				if err := store.InsertUsage(db, commitId, v.name, p, v.position.Line, "server.hooks"); err != nil {
					log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to insert usage implementing plugin hooks.")
				}
			}
		}
	}

}
