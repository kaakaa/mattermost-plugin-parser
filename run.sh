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

export GO111MODULE=off
# git clone $URL $WORKDIR
go get -u $URL
WORKDIR=~/go/src/${URL}
COMMIT_ID=`git -C $WORKDIR rev-parse HEAD`
COMMITTED_AT=`git -C $WORKDIR show -s --format=%ci HEAD`

log "cloned $COMMIT_ID"

# server
cd $WORKDIR
if [ -e "go.mod" ]; then
    export GO111MODULE=on
    go mod tidy -v
    go mod download
fi
if [ -e "server" ]; then
    cd server
fi
if [ -e "Gopkg.toml" ]; then
    export GO111MODULE=auto
    dep ensure
fi

log ""
log "Command: go run *.go $URL $COMMIT_ID $WORKDIR $COMMITTED_AT"
log ""
go run $ROOTDIR/server/*.go $URL $COMMIT_ID $WORKDIR "$COMMITTED_AT"

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