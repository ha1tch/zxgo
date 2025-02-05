#!/bin/bash
mkdir -p bin

echo "Building all binaries for zxgotools..."
pushd cmd/maketap
GOOS=windows GOARCH=amd64 go build -x -o ../../bin/maketap.exe          maketap.go
GOOS=windows GOARCH=386   go build -x -o ../../bin/maketap.win32.exe    maketap.go
GOOS=linux   GOARCH=amd64 go build -x -o ../../bin/maketap.linux        maketap.go
GOOS=linux   GOARCH=386   go build -x -o ../../bin/maketap.linux32      maketap.go
GOOS=linux   GOARCH=arm   go build -x -o ../../bin/maketap.rpi          maketap.go
GOOS=linux   GOARCH=arm64 go build -x -o ../../bin/maketap.rpi64        maketap.go
GOOS=darwin  GOARCH=arm64 go build -x -o ../../bin/maketap.mac          maketap.go
popd

pushd cmd/loadtap
GOOS=windows GOARCH=amd64 go build -x -o ../../bin/loadtap.exe          loadtap.go
GOOS=windows GOARCH=386   go build -x -o ../../bin/loadtap.win32.exe    loadtap.go
GOOS=linux   GOARCH=amd64 go build -x -o ../../bin/loadtap.linux        loadtap.go
GOOS=linux   GOARCH=386   go build -x -o ../../bin/loadtap.linux32      loadtap.go
GOOS=linux   GOARCH=arm   go build -x -o ../../bin/loadtap.rpi          loadtap.go
GOOS=linux   GOARCH=arm64 go build -x -o ../../bin/loadtap.rpi64        loadtap.go
GOOS=darwin  GOARCH=arm64 go build -x -o ../../bin/loadtap.mac          loadtap.go
popd

pushd cmd/totap
GOOS=windows GOARCH=amd64 go build -x -o ../../bin/totap.exe          totap.go
GOOS=windows GOARCH=386   go build -x -o ../../bin/totap.win32.exe    totap.go
GOOS=linux   GOARCH=amd64 go build -x -o ../../bin/totap.linux        totap.go
GOOS=linux   GOARCH=386   go build -x -o ../../bin/totap.linux32      totap.go
GOOS=linux   GOARCH=arm   go build -x -o ../../bin/totap.rpi          totap.go
GOOS=linux   GOARCH=arm64 go build -x -o ../../bin/totap.rpi64        totap.go
GOOS=darwin  GOARCH=arm64 go build -x -o ../../bin/totap.mac          totap.go
popd

(pushd cmd/tap2tzx && ./mk.sh)
popd