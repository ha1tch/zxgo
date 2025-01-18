#!/bin/bash
mkdir -p bin

# Build maketap for different platforms
echo "Building maketap..."
GOOS=windows GOARCH=amd64 go build -o bin/maketap.exe   ./cmd/maketap
GOOS=windows GOARCH=386   go build -o bin/maketap.win32.exe ./cmd/maketap
GOOS=linux   GOARCH=amd64 go build -o bin/maketap.linux ./cmd/maketap
GOOS=linux   GOARCH=386   go build -o bin/maketap.linux32 ./cmd/maketap
GOOS=darwin  GOARCH=arm64 go build -o bin/maketap.mac   ./cmd/maketap

# Build loadtap for different platforms
echo "Building loadtap..."
GOOS=windows GOARCH=amd64 go build -o bin/loadtap.exe   ./cmd/loadtap
GOOS=windows GOARCH=386   go build -o bin/loadtap.win32.exe ./cmd/loadtap
GOOS=linux   GOARCH=amd64 go build -o bin/loadtap.linux ./cmd/loadtap
GOOS=linux   GOARCH=386   go build -o bin/loadtap.linux32 ./cmd/loadtap
GOOS=darwin  GOARCH=arm64 go build -o bin/loadtap.mac   ./cmd/loadtap

# Build tap2tzx (from its own directory)
echo "Building tap2tzx..."
(cd cmd/tap2tzx && ./mk.sh)