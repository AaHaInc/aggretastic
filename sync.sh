#!/bin/bash
cwd=$(pwd)
if [ ! -f sync ]; then
    go get -u -f github.com/dkonovenschi/aggretastic-sync
    cd $GOPATH/src/github.com/dkonovenschi/aggretastic-sync
    go build -o sync
    cp sync "$cwd"
    cd $cwd
fi

go get github.com/olivere/elastic
./sync
