#!/bin/bash -eu

cd "$(dirname "${0}")"/..

go get -d -t ./...
go install -ldflags "-X main.revision=$(git rev-parse HEAD) -X main.version=$(git tag -l 'v*.*.*' --points-at HEAD | sed -n '1s/^v//p')" github.com/petemoore/samfile/cmd/samfile
go test -v ./...
echo
echo "$(go env GOPATH)/bin/samfile built and tested successfully"
