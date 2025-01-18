#!/bin/bash
mkdir -p bin

echo "Building all binaries for zxgotools..."
pushd cmd/maketap
GOOS=windows GOARCH=amd64 go build -x -o ../../bin/maketap.exe       maketap.go
GOOS=windows GOARCH=386   go build -x -o ../../bin/maketap.win32.exe maketap.go
GOOS=linux   GOARCH=amd64 go build -x -o ../../bin/maketap.linux     maketap.go
GOOS=linux   GOARCH=386   go build -x -o ../../bin/maketap.linux32   maketap.go
GOOS=darwin  GOARCH=arm64 go build -x -o ../../bin/maketap.mac       maketap.go
popd

pushd cmd/loadtap
GOOS=windows GOARCH=amd64 go build -x -o ../../bin/loadtap.exe          loadtap.go
GOOS=windows GOARCH=386   go build -x -o ../../bin/loadtap.win32.exe    loadtap.go
GOOS=linux   GOARCH=amd64 go build -x -o ../../bin/loadtap.linux        loadtap.go
GOOS=linux   GOARCH=386   go build -x -o ../../bin/loadtap.linux32      loadtap.go
GOOS=darwin  GOARCH=arm64 go build -x -o ../../bin/loadtap.mac          loadtap.go
popd

(pushd cmd/tap2tzx && ./mk.sh)
popd

