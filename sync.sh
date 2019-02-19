#!/bin/bash
cwd=$(pwd)
go get -u github.com/olivere/elastic

if [ ! -f sync ]; then
    cd synchronizer
    go build -o sync
    cp sync ../
    cd ../
fi

if [ -d vendor ]; then
  mv vendor _vendor
fi

go get github.com/olivere/elastic
./sync

if [ -d _vendor ]; then
  mv _vendor vendor
fi
