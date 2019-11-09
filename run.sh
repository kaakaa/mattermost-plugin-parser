#!/bin/bash -x

FILENAME=repositories.txt
ROOTDIR=${PWD}
# WORKDIR=/tmp/mmplugin_parser

# rm -fr $WORKDIR
# mkdir -p $WORKDIR

while read URL; do
set +x
echo ""
echo "##################################################################"
echo "# Parse plugin: $URL #"
echo "##################################################################"
echo ""
set -x

export GO111MODULE=off
# git clone $URL $WORKDIR
go get -u $URL
WORKDIR=~/go/src/${URL}
COMMIT_ID=`git -C $WORKDIR rev-parse HEAD`
COMMITTED_AT=`git -C $WORKDIR show -s --format=%ci HEAD`

echo "cloned: $COMMIT_ID"

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

echo "Command: go run *.go $URL $COMMIT_ID $WORKDIR $COMMITTED_AT"
go run $ROOTDIR/server/*.go $URL $COMMIT_ID $WORKDIR "$COMMITTED_AT"

# webapp
echo "Command: node index.js $URL $COMMIT_ID $WORKDIR $COMMITTED_AT"
node $ROOTDIR/webapp/index.js $URL $COMMIT_ID $WORKDIR "$COMMITTED_AT"

cd $ROOTDIR
rm -fr $WORKDIR
done < $FILENAME