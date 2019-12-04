package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kaakaa/mattermost-plugin-parser/server/mattermost"
	"github.com/kaakaa/mattermost-plugin-parser/server/pluginapi"
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/analysis"
	"gopkg.in/yaml.v2"
)

var (
	RepositoryURL string
	CommitID      string
	CommitedAt    string
)
var Analyzer = &analysis.Analyzer{
	Name: "StorePluginUsages",
	Doc:  "Plusing Usage Store",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		usages, ok := pass.ResultOf[pluginapi.Analyzer].(*pluginapi.PluginAPIUsages)
		if !ok {
			return nil, fmt.Errorf("Failed to get MattermostPluginAPI")
		}

		manifest, ok := pass.ResultOf[mattermost.PluginManifestAnalyzer].(*model.Manifest)
		if !ok {
			return nil, fmt.Errorf("Failed to get MattermostPluginAPI")
		}
		store(usages, manifest)
		return nil, nil
	},
	Requires: []*analysis.Analyzer{
		mattermost.PluginManifestAnalyzer,
		pluginapi.Analyzer,
	},
}

func store(usages *pluginapi.PluginAPIUsages, manifest *model.Manifest) {
	log.WithFields(log.Fields{
		"repository": RepositoryURL,
		"commitId":   CommitID,
		"commitedAt": CommitedAt,
	}).Info("mattermost-plugin-parser (server) starts.")

	// init database
	db, err := InitDB()
	if err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to init db.")
	}

	// Store repository information
	if err := InsertRepository(db, RepositoryURL, CommitID, CommitedAt); err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to insert repository info to db.")
	}

	// Store plugin manifest
	if err := InsertManifest(db, CommitID, manifest); err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to insert plugin manifest.")
	}

	// Store usages of plugin APIs/Helpers
	for _, v := range usages.APIs {
		if err := InsertUsage(db, CommitID, v.Name, v.Path, v.Position.Line, "server.api"); err != nil {
			log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to insert usage invoking plugin API.")
		}
	}

	// Store implements of plugin Hooks
	for _, v := range usages.Hooks {
		if err := InsertUsage(db, CommitID, v.Name, v.Path, v.Position.Line, "server.hooks"); err != nil {
			log.WithFields(log.Fields{"details": err.Error()}).Fatal("Failed to insert usage implementing plugin hooks.")
		}
	}
}

func parseMattermostPluginManifest(path string) (*model.Manifest, error) {
	switch {
	case fileExists(filepath.Join(path, "plugin.json")):
		file := filepath.Join(path, "plugin.json")
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		var manifest model.Manifest
		if err := json.Unmarshal(b, &manifest); err != nil {
			return nil, err
		}
		return &manifest, nil
	case fileExists(filepath.Join(path, "plugin.yaml")):
		file := filepath.Join(path, "plugin.yaml")
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		var manifest model.Manifest
		if err := yaml.Unmarshal(b, &manifest); err != nil {
			return nil, err
		}
		return &manifest, nil
	}
	return nil, errors.Errorf("There is no manifest file in %s.", path)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
