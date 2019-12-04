package mattermost

import (
	"reflect"

	"github.com/mattermost/mattermost-server/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/analysis"
)

var RootDir string

var PluginManifestAnalyzer = &analysis.Analyzer{
	Name: "PluginManifestAnalyzer",
	Doc:  Doc,
	Run: func(pass *analysis.Pass) (interface{}, error) {
		manifest, path, err := model.FindManifest(RootDir)
		if err != nil {
			return nil, err
		}
		log.WithFields(log.Fields{
			"path": path,
		}).Info("Found plugin manifest file.")
		return manifest, nil
	},
	ResultType: reflect.TypeOf((*model.Manifest)(nil)),
}
