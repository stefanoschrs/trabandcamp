#!/bin/bash

mkdir -p build

### Linux ###
# (export GOOS=linux; export GOARCH=386; go build -o build/trabandcamp-$GOOS-$GOARCH trabandcamp.go)
(export GOOS=linux; export GOARCH=amd64; go build -o build/trabandcamp-$GOOS-$GOARCH trabandcamp.go)

### Windows ###
# (export GOOS=windows; export GOARCH=386; go build -o build/trabandcamp-$GOOS-$GOARCH.exe trabandcamp.go)
# (export GOOS=windows; export GOARCH=amd64; go build -o build/trabandcamp-$GOOS-$GOARCH.exe trabandcamp.go)
