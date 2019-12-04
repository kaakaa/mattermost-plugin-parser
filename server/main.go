package main

import (
	"github.com/kaakaa/mattermost-plugin-parser/server/mattermost"
	"github.com/kaakaa/mattermost-plugin-parser/server/pluginapi"
	"github.com/kaakaa/mattermost-plugin-parser/server/store"
	"golang.org/x/tools/go/analysis/multichecker"
)

func init() {
	mattermost.PluginManifestAnalyzer.Flags.StringVar(&mattermost.RootDir, "rootdir", mattermost.RootDir, "root dir path of package")
	pluginapi.Analyzer.Flags.StringVar(&pluginapi.Target, "package", pluginapi.Target, "root package name for parsing")

	store.Analyzer.Flags.StringVar(&store.RepositoryURL, "repository", store.RepositoryURL, "url of plugin repository")
	store.Analyzer.Flags.StringVar(&store.CommitID, "commitid", store.RepositoryURL, "commit id of plugin repository")
	store.Analyzer.Flags.StringVar(&store.CommitedAt, "commitedat", store.RepositoryURL, "date time when parsed soure are commited")
}

func main() {
	multichecker.Main(
		mattermost.PluginManifestAnalyzer,
		mattermost.PluginAPIAnalyzer,
		pluginapi.Analyzer,
		store.Analyzer,
	)
}
