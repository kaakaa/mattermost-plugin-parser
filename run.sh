
FILENAME=repositories.txt
WORKDIR=${PWD}/mmplugin_parser

rm -fr $WORKDIR
mkdir -p $WORKDIR

while read URL; do
git clone $URL $WORKDIR
COMMIT_ID=`git -C $WORKDIR rev-parse HEAD`
COMMITTED_AT=`git -C $WORKDIR show -s --format=%ci HEAD`

echo "cloned: $COMMIT_ID"

# server
cd server
go run *.go $URL $COMMIT_ID $WORKDIR "$COMMITTED_AT"
cd ..

# webapp
cd webapp
node index.js $URL $COMMIT_ID $WORKDIR "$COMMITTED_AT"
cd ../

rm -fr $WORKDIR
done < $FILENAME