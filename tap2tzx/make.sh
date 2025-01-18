#!/bin/sh

# Initialize module if go.mod doesn't exist
if [ ! -f go.mod ]; then
    go mod init tap2tzx
fi

# Get dependencies
go get gopkg.in/yaml.v3

# Build
go build -o tap2tzx

# Make executable
chmod +x tap2tzx