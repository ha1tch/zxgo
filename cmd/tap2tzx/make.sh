#!/bin/sh

# Initialize module if go.mod doesn't exist
if [ ! -f go.mod ]; then
    go mod init tap2tzx
fi

# Get dependencies
go get gopkg.in/yaml.v3

# Build
GOOS=windows GOARCH=amd64 go build -o ../bin/tap2tzx.exe
GOOS=darwin  GOARCH=arm64 go build -o ../bin/tap2tzx.mac
GOOS=linux   GOARCH=amd64 go build -o ../bin/tap2tzx.linux

# Make executable
chmod +x ../bin/tap2tzx.*
