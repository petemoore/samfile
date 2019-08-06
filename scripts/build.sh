#!/bin/bash -e

# This script is needed because `go test -covermode=atomic` cover doesn't
# currently support being run against multiple packages

cd "$(dirname "${0}")"/..

go get -d -t ./...
go install -ldflags "-X main.revision=$(git rev-parse HEAD) -X main.version=$(git tag -l 'v*.*.*' --points-at HEAD | sed -n '1s/^v//p')"
"$(go env GOPATH)/bin/samfile" --help
