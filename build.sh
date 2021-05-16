#!/usr/bin/env bash
set -u

version=`git log -n 1 --date=short --pretty=format:"%cd-%h"`
echo "* ver: $version"

echo "* go mod tidy"
go mod tidy

echo "* go generate"
go generate ./...

echo "* go fmt"
go fmt ./...

echo "* go vet"
go vet ./...

echo "* staticcheck"
staticcheck ./...

echo "* go test"
go test -v ./...

#echo "* Build kaban: "
#env GOARCH=amd64 GOOS=linux go build \
#  -ldflags "-s -X main.version=$version" -race \
#  -o kaban.exe cmd/kaban/main.go
#if [ $? -eq 0 ]; then
#  echo "  OK"
#else
#  echo "  NG"
#fi
