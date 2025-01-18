#!/bin/sh

# Ensure we're in the tap2tzx directory
cd "$(dirname "$0")"

# Initialize module if go.mod doesn't exist
if [ ! -f go.mod ]; then
    go mod init tap2tzx
fi

# Get dependencies
go get gopkg.in/yaml.v3

# Build into the main bin directory (two levels up)
echo "Building tap2tzx..."
GOOS=windows GOARCH=amd64 go build -o ../../bin/tap2tzx.exe
GOOS=windows GOARCH=386   go build -o ../../bin/tap2tzx.win32.exe
GOOS=linux   GOARCH=amd64 go build -o ../../bin/tap2tzx.linux
GOOS=darwin  GOARCH=arm64 go build -o ../../bin/tap2tzx.mac

# Make executable
chmod +x ../../bin/tap2tzx.*