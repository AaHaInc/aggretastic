#!/bin/bash
cwd=$(pwd)
if [ ! -f sync ]; then
    go get gitlab.com/dmitry.konovenschi/aggretastic-sync
    cd $GOPATH/src/gitlab.com/dmitry.konovenschi/aggretastic-sync
    dep ensure
    go build -o sync
    cp sync "$cwd"
    cd $cwd
fi

go get github.com/olivere/elastic
./sync