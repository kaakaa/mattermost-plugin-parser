#!/bin/bash

readonly PROCNAME=${0##*/}
function log() {
  local fname=${BASH_SOURCE[1]##*/}
  echo -e "$(date '+%Y-%m-%dT%H:%M:%S') ${PROCNAME} (${fname}:${BASH_LINENO[0]}:${FUNCNAME[1]}) $@"
}

FILENAME=repositories.txt
ROOTDIR=${PWD}
# WORKDIR=/tmp/mmplugin_parser

# rm -fr $WORKDIR
# mkdir -p $WORKDIR

while read URL; do
echo ""
echo "##################################################################"
echo "# Parse plugin: $URL #"
echo "##################################################################"
echo ""

log "Clone repository"


# NOTE:
# I want to use `go get` command for fetching plugin source code.
# However, many of plugin generated from `mattermost/mattermost-plugin-starter-template` have incorrect module name in go.mod.
# Their module name should be their package name, but they have used `mattermost-plugin-starter-template` as is.
# This incorrectness will cause the failure of `go get` command.
# 
# EXAMPLE:
# > $ go get github.com/kaakaa/mattermost-plugin-loudspeaker
# > go: finding github.com/kaakaa/mattermost-plugin-loudspeaker latest
# > go get: github.com/kaakaa/mattermost-plugin-loudspeaker@v0.0.0-20190716120632-260d2acffcf9: parsing go.mod:
# >         module declares its path as: github.com/mattermost/mattermost-plugin-starter-template
# >                 but was required as: github.com/kaakaa/mattermost-plugin-loudspeaker
# export GO111MODULE=off
# go get -u $URL
# WORKDIR=~/go/src/${URL}

TEMP_REPO=.work
mkdir $TEMP_REPO
git clone https://$URL $TEMP_REPO
WORKDIR=$PWD/$TEMP_REPO

COMMIT_ID=`git -C $WORKDIR rev-parse HEAD`
COMMITTED_AT=`git -C $WORKDIR show -s --format=%ci HEAD`

log "cloned $COMMIT_ID"

# server
cd $WORKDIR
# go install
if [ -e "go.mod" ]; then
  log "Found go.mod file"
  export GO111MODULE=on
  # go mod tidy -v
  # go install
  go build ./...
fi
if [ -e "./server" ]; then
  log "Found server directory"
  # cd $WORKDIR
fi
if [ -e "Gopkg.toml" ]; then
  log "Found Gopkg.toml"
  export GO111MODULE=auto
  dep ensure
fi

# cd $ROOTDIR

# FIXME: Remove this line after updating repository
echo "replace github.com/kaakaa/mattermost-plugin-parser => ../" >> go.mod

log ""
log "Command: go run *.go $URL $COMMIT_ID $WORKDIR $COMMITTED_AT"
log ""
go run $ROOTDIR/server/main.go \
  -MattermostPluginAnalyzer.package "$URL" \
  -PluginManifestAnalyzer.rootdir "$WORKDIR" \
  -StorePluginUsages.repository "$URL" \
  -StorePluginUsages.commitid "$COMMIT_ID" \
  -StorePluginUsages.commitedat "$COMMITTED_AT" \
  ./...

# webapp
log ""
log "Command: node index.js $URL $COMMIT_ID $WORKDIR $COMMITTED_AT"
log ""
node $ROOTDIR/webapp/index.js $URL $COMMIT_ID $WORKDIR "$COMMITTED_AT"

cd $ROOTDIR
rm -fr $WORKDIR

log ""
log "Done"
done < $FILENAME